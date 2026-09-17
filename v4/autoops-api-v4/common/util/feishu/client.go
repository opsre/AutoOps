// 飞书 Open API 客户端
// 封装 tenant_access_token 获取（Redis缓存）、Wiki、文档、搜索等接口
package feishu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"dodevops-api/common"
	"dodevops-api/pkg/log"
)

const (
	feishuBaseURL = "https://open.feishu.cn/open-apis"
	tokenCacheTTL = 6480 * time.Second // 1.8 小时，飞书 token 有效期 2 小时
)

// -------- 响应结构体 --------

// AppAccessTokenResp 获取 tenant_access_token 响应
type AppAccessTokenResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

// WikiSpace Wiki 空间信息
type WikiSpace struct {
	SpaceID     string `json:"space_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SpaceType   string `json:"space_type"`
}

// WikiNode Wiki 节点信息
type WikiNode struct {
	SpaceID   string `json:"space_id"`
	NodeToken string `json:"node_token"`
	ObjToken  string `json:"obj_token"`
	ObjType   string `json:"obj_type"` // docx / wiki
	Title     string `json:"title"`
	HasChild  bool   `json:"has_child"`
}

// SearchResult 文档搜索结果
type SearchResult struct {
	DocToken string `json:"docs_token"`
	DocType  string `json:"docs_type"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	OwnerID  string `json:"owner_id"`
}

// -------- Token 管理 --------

// GetAppAccessToken 获取飞书 tenant_access_token，优先读 Redis 缓存
// keyManageID 用于构造唯一缓存 key，支持多飞书应用
func GetAppAccessToken(keyManageID uint, appID, appSecret string) (string, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("feishu:token:%d", keyManageID)

	// 先查 Redis
	rc := common.GetRedisClient()
	if rc != nil {
		token, err := rc.Get(ctx, cacheKey).Result()
		if err == nil && token != "" {
			log.Log().Infof("[Feishu] token hit cache, key=%s prefix=%s...", cacheKey, token[:10])
			return token, nil
		}
	}

	// 调飞书接口获取 token
	reqBody, _ := json.Marshal(map[string]string{
		"app_id":     appID,
		"app_secret": appSecret,
	})
	resp, err := http.Post(
		feishuBaseURL+"/auth/v3/tenant_access_token/internal",
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return "", fmt.Errorf("请求飞书 token 接口失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result AppAccessTokenResp
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析飞书 token 响应失败: %w", err)
	}
	if result.Code != 0 {
		return "", fmt.Errorf("飞书 token 获取失败(code=%d): %s", result.Code, result.Msg)
	}

	// 写入 Redis，TTL = 1.8h
	if rc != nil {
		_ = rc.Set(ctx, cacheKey, result.TenantAccessToken, tokenCacheTTL).Err()
	}

	log.Log().Infof("[Feishu] token fetched from feishu, key=%s prefix=%s...", cacheKey, result.TenantAccessToken[:10])
	return result.TenantAccessToken, nil
}

// -------- 内部工具 --------

// doRequest 统一执行带 Authorization 头的 HTTP 请求
func doRequest(method, url, accessToken string, reqBody []byte) ([]byte, error) {
	var bodyReader io.Reader
	if reqBody != nil {
		bodyReader = bytes.NewReader(reqBody)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// -------- Wiki 接口 --------

// GetWikiSpaces 获取 Wiki 空间列表（最多50个）
func GetWikiSpaces(accessToken string) ([]WikiSpace, error) {
	body, err := doRequest("GET", feishuBaseURL+"/wiki/v2/spaces?page_size=50", accessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("请求 Wiki 空间列表失败: %w", err)
	}

	log.Log().Infof("[Feishu] GetWikiSpaces raw response: %s", string(body))

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Items []WikiSpace `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 Wiki 空间响应失败: %w", err)
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("获取 Wiki 空间列表失败(code=%d): %s", result.Code, result.Msg)
	}
	return result.Data.Items, nil
}

// GetWikiNodes 获取 Wiki 空间节点列表
// parentNodeToken 为空时获取根节点，非空时获取该节点的子节点
func GetWikiNodes(accessToken, spaceID, parentNodeToken string) ([]WikiNode, error) {
	url := fmt.Sprintf("%s/wiki/v2/spaces/%s/nodes?page_size=50", feishuBaseURL, spaceID)
	if parentNodeToken != "" {
		url += "&parent_node_token=" + parentNodeToken
	}

	body, err := doRequest("GET", url, accessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("请求 Wiki 节点列表失败: %w", err)
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Items []WikiNode `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 Wiki 节点响应失败: %w", err)
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("获取 Wiki 节点列表失败(code=%d): %s", result.Code, result.Msg)
	}
	return result.Data.Items, nil
}

// -------- 文档接口 --------

// CreateDocx 在飞书创建一个新文档，返回 document_id（即 docToken）
// folderToken 为空则创建在根目录
func CreateDocx(accessToken, title, folderToken string) (string, error) {
	payload := map[string]string{"title": title}
	if folderToken != "" {
		payload["folder_token"] = folderToken
	}
	reqBody, _ := json.Marshal(payload)

	body, err := doRequest("POST", feishuBaseURL+"/docx/v1/documents", accessToken, reqBody)
	if err != nil {
		return "", fmt.Errorf("创建飞书文档失败: %w", err)
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Document struct {
				DocumentID string `json:"document_id"`
				Title      string `json:"title"`
			} `json:"document"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析创建文档响应失败: %w", err)
	}
	if result.Code != 0 {
		return "", fmt.Errorf("创建文档失败(code=%d): %s", result.Code, result.Msg)
	}
	return result.Data.Document.DocumentID, nil
}

// GetDocumentContent 获取飞书文档 Block 内容（原始 JSON）
func GetDocumentContent(accessToken, docToken string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/docx/v1/documents/%s/blocks", feishuBaseURL, docToken)
	body, err := doRequest("GET", url, accessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("获取文档内容失败: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析文档内容失败: %w", err)
	}
	return result, nil
}

// -------- 搜索接口 --------

// SearchDocuments 通过飞书搜索 API 全文搜索文档
func SearchDocuments(accessToken, keyword string, page, pageSize int) ([]SearchResult, int, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}
	reqBody, _ := json.Marshal(map[string]interface{}{
		"search_key": keyword,
		"count":      pageSize,
		"offset":     (page - 1) * pageSize,
		"docs_types": []string{"docx", "wiki"},
	})

	body, err := doRequest("POST", feishuBaseURL+"/suite/docs-api/search/object", accessToken, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("搜索文档失败: %w", err)
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			DocsEntitys []SearchResult `json:"docs_entitys"`
			Total       int            `json:"total"`
			HasMore     bool           `json:"has_more"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, 0, fmt.Errorf("解析搜索结果失败: %w", err)
	}
	if result.Code != 0 {
		return nil, 0, fmt.Errorf("搜索文档失败(code=%d): %s", result.Code, result.Msg)
	}
	return result.Data.DocsEntitys, result.Data.Total, nil
}

// -------- 诊断工具 --------

// ValidateToken 验证 token 是否有效，并返回飞书原始响应（用于排查权限问题）
// 调用 /auth/v3/app_info 获取应用信息
func ValidateToken(accessToken string) (map[string]interface{}, error) {
	body, err := doRequest("GET", feishuBaseURL+"/application/v6/applications/underauditlist?lang=zh_cn&page_size=1", accessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("验证请求失败: %w", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	return raw, nil
}

// GetWikiSpacesRaw 返回飞书 Wiki 空间接口的原始响应（用于排查权限问题）
func GetWikiSpacesRaw(accessToken string) (map[string]interface{}, error) {
	body, err := doRequest("GET", feishuBaseURL+"/wiki/v2/spaces?page_size=50", accessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	var raw map[string]interface{}
	_ = json.Unmarshal(body, &raw)
	return raw, nil
}

// -------- URL 工具 --------

// BuildDocURL 构建飞书文档访问 URL（用于前端 iframe 嵌入）
// docType: docx / wiki
func BuildDocURL(docType, docToken string) string {
	if docType == "wiki" {
		return fmt.Sprintf("https://feishu.cn/wiki/%s", docToken)
	}
	return fmt.Sprintf("https://feishu.cn/docx/%s", docToken)
}
