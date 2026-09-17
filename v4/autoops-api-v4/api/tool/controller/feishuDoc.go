// 飞书文档知识库 Controller 层（方案A：纯代理模式 + OAuth 用户授权）
package controller

import (
	"strconv"

	"dodevops-api/api/tool/service"

	"github.com/gin-gonic/gin"
)

// -------- 应用级接口 --------

// GetFeishuApps 获取飞书应用列表
// @Tags Feishu飞书知识库
// @Summary 获取飞书应用列表
// @Description 获取密钥管理中所有 KeyType=7 的飞书应用（脱敏，不返回 AppSecret）
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/tool/feishu/apps [get]
// @Security ApiKeyAuth
func GetFeishuApps(c *gin.Context) {
	service.FeishuService().GetFeishuApps(c)
}

// GetFeishuWikiSpaces 获取 Wiki 空间列表
// @Tags Feishu飞书知识库
// @Summary 获取 Wiki 空间列表
// @Produce json
// @Param keyManageId query int true "飞书应用密钥ID（KeyType=7）"
// @Param openId query string false "用户 openId（不传则使用 app token 查团队空间，传入则查个人空间）"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/feishu/spaces [get]
// @Security ApiKeyAuth
func GetFeishuWikiSpaces(c *gin.Context) {
	keyManageID, err := parseUintQuery(c, "keyManageId")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "keyManageId 参数错误"})
		return
	}
	openID := c.Query("openId")
	service.FeishuService().GetWikiSpaces(c, keyManageID, openID)
}

// GetFeishuWikiNodes 获取 Wiki 节点列表
// @Tags Feishu飞书知识库
// @Summary 获取 Wiki 节点列表（目录树）
// @Produce json
// @Param keyManageId query int true "飞书应用密钥ID"
// @Param spaceId path string true "Wiki 空间ID"
// @Param parentNodeToken query string false "父节点Token（不传则获取根节点）"
// @Param openId query string false "用户 openId（个人知识库时传入）"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/feishu/spaces/{spaceId}/nodes [get]
// @Security ApiKeyAuth
func GetFeishuWikiNodes(c *gin.Context) {
	keyManageID, err := parseUintQuery(c, "keyManageId")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "keyManageId 参数错误"})
		return
	}
	spaceID := c.Param("spaceId")
	if spaceID == "" {
		c.JSON(400, gin.H{"code": 400, "message": "spaceId 不能为空"})
		return
	}
	parentNodeToken := c.Query("parentNodeToken")
	openID := c.Query("openId")
	service.FeishuService().GetWikiNodes(c, keyManageID, spaceID, parentNodeToken, openID)
}

// GetFeishuDocContent 获取飞书文档内容
// @Tags Feishu飞书知识库
// @Summary 获取飞书文档 Block 内容
// @Produce json
// @Param keyManageId query int true "飞书应用密钥ID"
// @Param docToken path string true "飞书文档 Token（document_id）"
// @Param openId query string false "用户 openId（个人知识库时传入）"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/feishu/docs/{docToken}/content [get]
// @Security ApiKeyAuth
func GetFeishuDocContent(c *gin.Context) {
	keyManageID, err := parseUintQuery(c, "keyManageId")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "keyManageId 参数错误"})
		return
	}
	docToken := c.Param("docToken")
	if docToken == "" {
		c.JSON(400, gin.H{"code": 400, "message": "docToken 不能为空"})
		return
	}
	openID := c.Query("openId")
	service.FeishuService().GetDocContent(c, keyManageID, docToken, openID)
}

// SearchFeishuDocs 搜索飞书文档
// @Tags Feishu飞书知识库
// @Summary 全文搜索飞书文档
// @Produce json
// @Param keyManageId query int true "飞书应用密钥ID"
// @Param q query string true "搜索关键词"
// @Param openId query string false "用户 openId（个人知识库时传入）"
// @Param pageNum query int false "页码，默认1"
// @Param pageSize query int false "每页数量，默认20"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/feishu/search [get]
// @Security ApiKeyAuth
func SearchFeishuDocs(c *gin.Context) {
	keyManageID, err := parseUintQuery(c, "keyManageId")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "keyManageId 参数错误"})
		return
	}
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(400, gin.H{"code": 400, "message": "搜索关键词 q 不能为空"})
		return
	}
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	openID := c.Query("openId")
	service.FeishuService().SearchDocs(c, keyManageID, keyword, pageNum, pageSize, openID)
}

// ValidateFeishu 诊断接口
// @Tags Feishu飞书知识库
// @Summary 验证飞书 token 并返回原始响应（排查权限用）
// @Produce json
// @Param keyManageId query int true "飞书应用密钥ID"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/feishu/validate [get]
// @Security ApiKeyAuth
func ValidateFeishu(c *gin.Context) {
	keyManageID, err := parseUintQuery(c, "keyManageId")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "keyManageId 参数错误"})
		return
	}
	service.FeishuService().Validate(c, keyManageID)
}

// -------- OAuth 用户授权接口 --------

// GetFeishuOAuthURL 获取飞书 OAuth 授权链接
// @Tags Feishu飞书知识库
// @Summary 获取 OAuth 授权链接（用于访问个人知识库）
// @Produce json
// @Param keyManageId query int true "飞书应用密钥ID"
// @Param frontendRedirect query string true "授权完成后前端接收结果的页面地址"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/feishu/oauth/url [get]
// @Security ApiKeyAuth
func GetFeishuOAuthURL(c *gin.Context) {
	keyManageID, err := parseUintQuery(c, "keyManageId")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "keyManageId 参数错误"})
		return
	}
	frontendRedirect := c.Query("frontendRedirect")
	if frontendRedirect == "" {
		c.JSON(400, gin.H{"code": 400, "message": "frontendRedirect 不能为空"})
		return
	}
	service.FeishuService().GetOAuthURL(c, keyManageID, frontendRedirect)
}

// HandleFeishuOAuthCallback 飞书 OAuth 回调（飞书授权后重定向到此地址，需在飞书开放平台注册）
// @Tags Feishu飞书知识库
// @Summary 飞书 OAuth 回调地址（由飞书服务器调用）
// @Router /api/v1/tool/feishu/oauth/callback [get]
func HandleFeishuOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		c.JSON(400, gin.H{"error": "缺少 code 或 state 参数"})
		return
	}
	service.FeishuService().HandleOAuthCallback(c, code, state)
}

// ListFeishuUserTokens 列出已授权的飞书用户
// @Tags Feishu飞书知识库
// @Summary 列出已授权的飞书用户（脱敏，不含 token）
// @Produce json
// @Param keyManageId query int true "飞书应用密钥ID"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/feishu/oauth/users [get]
// @Security ApiKeyAuth
func ListFeishuUserTokens(c *gin.Context) {
	keyManageID, err := parseUintQuery(c, "keyManageId")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "keyManageId 参数错误"})
		return
	}
	service.FeishuService().ListUserTokens(c, keyManageID)
}

// RevokeFeishuUserToken 撤销飞书用户授权
// @Tags Feishu飞书知识库
// @Summary 撤销指定用户的飞书授权
// @Produce json
// @Param keyManageId query int true "飞书应用密钥ID"
// @Param openId query string true "用户 openId"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/feishu/oauth/users [delete]
// @Security ApiKeyAuth
func RevokeFeishuUserToken(c *gin.Context) {
	keyManageID, err := parseUintQuery(c, "keyManageId")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "keyManageId 参数错误"})
		return
	}
	openID := c.Query("openId")
	if openID == "" {
		c.JSON(400, gin.H{"code": 400, "message": "openId 不能为空"})
		return
	}
	service.FeishuService().RevokeUserToken(c, keyManageID, openID)
}

// -------- 内部工具 --------

func parseUintQuery(c *gin.Context, key string) (uint, error) {
	val, err := strconv.ParseUint(c.Query(key), 10, 32)
	return uint(val), err
}
