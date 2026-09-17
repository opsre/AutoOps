// 飞书文档知识库 Service 层（方案A：纯代理模式 + OAuth 用户授权）
// 职责：
//   - 从 KeyManage 读取飞书凭证 → 获取 app access token → 透传飞书 API（团队知识库）
//   - OAuth 2.0 授权流程：生成授权URL → 处理回调 → 存储用户 token
//   - 用户 token 自动刷新：在 access_token 过期前静默刷新
package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	ccDao "dodevops-api/api/configcenter/dao"
	tokenDao "dodevops-api/api/tool/dao1"
	tokenModel "dodevops-api/api/tool/model"
	"dodevops-api/common/config"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/common/util/feishu"
	"dodevops-api/pkg/log"

	"github.com/gin-gonic/gin"
)

const feishuKeyType = 7 // KeyManage.KeyType = 7 代表飞书应用

// IFeishuService 接口定义
type IFeishuService interface {
	// 应用级（团队知识库，使用 tenant_access_token）
	GetFeishuApps(c *gin.Context)
	GetWikiSpaces(c *gin.Context, keyManageID uint, openID string)
	GetWikiNodes(c *gin.Context, keyManageID uint, spaceID, parentNodeToken, openID string)
	GetDocContent(c *gin.Context, keyManageID uint, docToken, openID string)
	SearchDocs(c *gin.Context, keyManageID uint, keyword string, page, pageSize int, openID string)
	Validate(c *gin.Context, keyManageID uint)

	// OAuth 用户授权（个人知识库）
	GetOAuthURL(c *gin.Context, keyManageID uint, frontendRedirect string)
	HandleOAuthCallback(c *gin.Context, code, state string)
	ListUserTokens(c *gin.Context, keyManageID uint)
	RevokeUserToken(c *gin.Context, keyManageID uint, openID string)
}

type feishuServiceImpl struct {
	keyManageDao *ccDao.KeyManageDao
}

// -------- Token 获取 --------

// getAppAccessToken 获取 tenant_access_token（应用级，Redis 缓存）
func (s *feishuServiceImpl) getAppAccessToken(keyManageID uint) (string, error) {
	keyManage, err := s.keyManageDao.GetByID(keyManageID)
	if err != nil {
		return "", fmt.Errorf("密钥记录不存在(id=%d): %w", keyManageID, err)
	}
	if keyManage.KeyType != feishuKeyType {
		return "", fmt.Errorf("密钥类型不匹配，期望飞书应用(%d)，实际为: %d", feishuKeyType, keyManage.KeyType)
	}
	appID, appSecret, err := keyManage.DecryptKeys()
	if err != nil {
		return "", fmt.Errorf("密钥解密失败: %w", err)
	}
	return feishu.GetAppAccessToken(keyManageID, appID, appSecret)
}

// getEffectiveToken 获取有效 token：
//   - openID 为空 → 使用 tenant_access_token（团队知识库）
//   - openID 非空 → 使用 user_access_token（个人知识库），自动刷新
func (s *feishuServiceImpl) getEffectiveToken(keyManageID uint, openID string) (string, error) {
	if openID == "" {
		return s.getAppAccessToken(keyManageID)
	}
	return s.getUserAccessToken(keyManageID, openID)
}

// getUserAccessToken 获取用户 access_token，若即将过期则自动刷新
func (s *feishuServiceImpl) getUserAccessToken(keyManageID uint, openID string) (string, error) {
	record, err := tokenDao.GetFeishuUserToken(keyManageID, openID)
	if err != nil {
		return "", fmt.Errorf("用户授权记录不存在(openId=%s)，请重新授权: %w", openID, err)
	}

	// refresh_token 已过期 → 需要重新授权
	if time.Now().After(record.RefreshExpiresAt) {
		return "", fmt.Errorf("用户授权已过期(openId=%s)，请重新授权", openID)
	}

	// access_token 剩余时间 < 30 分钟 → 静默刷新
	if time.Now().Add(30 * time.Minute).After(record.ExpiresAt) {
		log.Log().Infof("[Feishu] user token 即将过期，自动刷新 openId=%s", openID)
		if refreshErr := s.refreshUserToken(keyManageID, record); refreshErr != nil {
			log.Log().Warnf("[Feishu] 自动刷新 token 失败: %v", refreshErr)
			// 刷新失败但 token 还没过期，暂时继续使用旧 token
			if time.Now().After(record.ExpiresAt) {
				return "", refreshErr
			}
		}
	}

	accessToken, err := util.AESDecrypt(record.AccessToken)
	if err != nil {
		return "", fmt.Errorf("用户 token 解密失败: %w", err)
	}
	return accessToken, nil
}

// refreshUserToken 用 refresh_token 刷新并持久化新 token
func (s *feishuServiceImpl) refreshUserToken(keyManageID uint, record *tokenModel.FeishuUserToken) error {
	appToken, err := s.getAppAccessToken(keyManageID)
	if err != nil {
		return err
	}
	refreshToken, err := util.AESDecrypt(record.RefreshToken)
	if err != nil {
		return fmt.Errorf("refresh_token 解密失败: %w", err)
	}
	resp, err := feishu.RefreshUserToken(appToken, refreshToken)
	if err != nil {
		return err
	}
	return s.persistToken(record, resp)
}

// persistToken 加密并保存/更新用户 token
func (s *feishuServiceImpl) persistToken(record *tokenModel.FeishuUserToken, resp *feishu.OAuthTokenResp) error {
	encAccess, err := util.AESEncrypt(resp.Data.AccessToken)
	if err != nil {
		return fmt.Errorf("access_token 加密失败: %w", err)
	}
	encRefresh, err := util.AESEncrypt(resp.Data.RefreshToken)
	if err != nil {
		return fmt.Errorf("refresh_token 加密失败: %w", err)
	}
	record.AccessToken = encAccess
	record.RefreshToken = encRefresh
	record.ExpiresAt = time.Now().Add(time.Duration(resp.Data.ExpiresIn) * time.Second)
	record.RefreshExpiresAt = time.Now().Add(time.Duration(resp.Data.RefreshExpiresIn) * time.Second)
	if resp.Data.Name != "" {
		record.UserName = resp.Data.Name
	}
	if resp.Data.AvatarURL != "" {
		record.AvatarURL = resp.Data.AvatarURL
	}
	return tokenDao.UpdateFeishuUserToken(record)
}

// -------- 应用级 Wiki API --------

// GetFeishuApps 获取所有飞书应用配置列表（KeyType=7，脱敏）
func (s *feishuServiceImpl) GetFeishuApps(c *gin.Context) {
	apps, err := s.keyManageDao.GetByType(feishuKeyType)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "获取飞书应用列表失败: "+err.Error())
		return
	}
	type AppVO struct {
		ID      uint   `json:"id"`
		KeyName string `json:"keyName"`
		Remark  string `json:"remark"`
	}
	list := make([]AppVO, 0, len(apps))
	for _, app := range apps {
		list = append(list, AppVO{ID: app.ID, KeyName: app.KeyName, Remark: app.Remark})
	}
	result.Success(c, list)
}

// GetWikiSpaces 获取 Wiki 空间列表
func (s *feishuServiceImpl) GetWikiSpaces(c *gin.Context, keyManageID uint, openID string) {
	token, err := s.getEffectiveToken(keyManageID, openID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	spaces, err := feishu.GetWikiSpaces(token)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "获取 Wiki 空间失败: "+err.Error())
		return
	}
	result.Success(c, spaces)
}

// GetWikiNodes 获取 Wiki 节点列表
func (s *feishuServiceImpl) GetWikiNodes(c *gin.Context, keyManageID uint, spaceID, parentNodeToken, openID string) {
	token, err := s.getEffectiveToken(keyManageID, openID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	nodes, err := feishu.GetWikiNodes(token, spaceID, parentNodeToken)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "获取 Wiki 节点失败: "+err.Error())
		return
	}
	result.Success(c, nodes)
}

// GetDocContent 获取飞书文档 Block 内容
func (s *feishuServiceImpl) GetDocContent(c *gin.Context, keyManageID uint, docToken, openID string) {
	token, err := s.getEffectiveToken(keyManageID, openID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	content, err := feishu.GetDocumentContent(token, docToken)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "获取文档内容失败: "+err.Error())
		return
	}
	result.Success(c, content)
}

// SearchDocs 全文搜索飞书文档
func (s *feishuServiceImpl) SearchDocs(c *gin.Context, keyManageID uint, keyword string, page, pageSize int, openID string) {
	token, err := s.getEffectiveToken(keyManageID, openID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	items, total, err := feishu.SearchDocuments(token, keyword, page, pageSize)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "搜索失败: "+err.Error())
		return
	}
	result.SuccessWithPage(c, items, int64(total), page, pageSize)
}

// Validate 诊断接口
func (s *feishuServiceImpl) Validate(c *gin.Context, keyManageID uint) {
	token, err := s.getAppAccessToken(keyManageID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "获取 token 失败: "+err.Error())
		return
	}
	raw, err := feishu.GetWikiSpacesRaw(token)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "请求飞书失败: "+err.Error())
		return
	}
	result.Success(c, map[string]interface{}{
		"tokenOk":       true,
		"tokenPrefix":   token[:10] + "...",
		"feishuRawResp": raw,
	})
}

// -------- OAuth 用户授权 --------

// GetOAuthURL 生成飞书 OAuth 授权链接
func (s *feishuServiceImpl) GetOAuthURL(c *gin.Context, keyManageID uint, frontendRedirect string) {
	keyManage, err := s.keyManageDao.GetByID(keyManageID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "密钥记录不存在: "+err.Error())
		return
	}
	appID, _, err := keyManage.DecryptKeys()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "密钥解密失败: "+err.Error())
		return
	}

	publicURL := strings.TrimRight(config.Config.Server.PublicURL, "/")
	if publicURL == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "服务端未配置 publicUrl，请在 config.yaml 中设置 server.publicUrl")
		return
	}
	callbackURL := publicURL + "/api/v1/tool/feishu/oauth/callback"
	state := feishu.EncodeState(keyManageID, frontendRedirect)
	oauthURL := feishu.BuildOAuthURL(appID, callbackURL, state)

	result.Success(c, gin.H{
		"oauthUrl":    oauthURL,
		"callbackUrl": callbackURL,
	})
}

// HandleOAuthCallback 处理飞书 OAuth 回调
func (s *feishuServiceImpl) HandleOAuthCallback(c *gin.Context, code, state string) {
	oauthState, err := feishu.DecodeState(state)
	if err != nil {
		c.Redirect(http.StatusFound, "/?feishuError=invalid_state")
		return
	}

	appToken, err := s.getAppAccessToken(oauthState.KeyManageID)
	if err != nil {
		log.Log().Errorf("[Feishu OAuth] 获取 app token 失败: %v", err)
		s.redirectWithError(c, oauthState.RedirectURI, "app_token_failed")
		return
	}

	resp, err := feishu.ExchangeCodeForToken(appToken, code)
	if err != nil {
		log.Log().Errorf("[Feishu OAuth] 换取 user token 失败: %v", err)
		s.redirectWithError(c, oauthState.RedirectURI, "token_exchange_failed")
		return
	}

	encAccess, err := util.AESEncrypt(resp.Data.AccessToken)
	if err != nil {
		s.redirectWithError(c, oauthState.RedirectURI, "encrypt_failed")
		return
	}
	encRefresh, err := util.AESEncrypt(resp.Data.RefreshToken)
	if err != nil {
		s.redirectWithError(c, oauthState.RedirectURI, "encrypt_failed")
		return
	}

	record := &tokenModel.FeishuUserToken{
		KeyManageID:      oauthState.KeyManageID,
		OpenID:           resp.Data.OpenID,
		UnionID:          resp.Data.UnionID,
		UserName:         resp.Data.Name,
		AvatarURL:        resp.Data.AvatarURL,
		AccessToken:      encAccess,
		RefreshToken:     encRefresh,
		ExpiresAt:        time.Now().Add(time.Duration(resp.Data.ExpiresIn) * time.Second),
		RefreshExpiresAt: time.Now().Add(time.Duration(resp.Data.RefreshExpiresIn) * time.Second),
	}
	if err := tokenDao.UpsertFeishuUserToken(record); err != nil {
		log.Log().Errorf("[Feishu OAuth] 存储 token 失败: %v", err)
		s.redirectWithError(c, oauthState.RedirectURI, "db_save_failed")
		return
	}

	log.Log().Infof("[Feishu OAuth] 用户授权成功 openId=%s userName=%s", resp.Data.OpenID, resp.Data.Name)

	redirectURL := oauthState.RedirectURI
	sep := "?"
	if strings.Contains(redirectURL, "?") {
		sep = "&"
	}
	redirectURL += sep + "feishuOpenId=" + resp.Data.OpenID + "&feishuUser=" + resp.Data.Name
	c.Redirect(http.StatusFound, redirectURL)
}

func (s *feishuServiceImpl) redirectWithError(c *gin.Context, frontendURI, errCode string) {
	if frontendURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errCode})
		return
	}
	sep := "?"
	if strings.Contains(frontendURI, "?") {
		sep = "&"
	}
	c.Redirect(http.StatusFound, frontendURI+sep+"feishuError="+errCode)
}

// ListUserTokens 列出已授权的用户（脱敏）
func (s *feishuServiceImpl) ListUserTokens(c *gin.Context, keyManageID uint) {
	list, err := tokenDao.ListFeishuUserTokens(keyManageID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "获取授权用户列表失败: "+err.Error())
		return
	}
	result.Success(c, list)
}

// RevokeUserToken 撤销用户授权
func (s *feishuServiceImpl) RevokeUserToken(c *gin.Context, keyManageID uint, openID string) {
	if err := tokenDao.DeleteFeishuUserToken(keyManageID, openID); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "撤销失败: "+err.Error())
		return
	}
	result.Success(c, nil)
}

// -------- 单例 --------

var feishuService = &feishuServiceImpl{
	keyManageDao: ccDao.NewKeyManageDao(),
}

// FeishuService 获取飞书服务实例
func FeishuService() IFeishuService {
	return feishuService
}
