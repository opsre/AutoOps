// 飞书 OAuth 2.0 授权工具
// 用于获取 user_access_token，访问个人知识库
package feishu

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// OAuthState OAuth state 结构（base64 编码后放入 state 参数，防 CSRF，携带回调信息）
type OAuthState struct {
	KeyManageID uint   `json:"k"` // 使用哪个飞书应用
	RedirectURI string `json:"r"` // 授权完成后前端跳转的地址
	Timestamp   int64  `json:"t"` // 防重放
}

// OAuthTokenResp 飞书 OAuth token 响应
type OAuthTokenResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		TokenType        string `json:"token_type"`
		ExpiresIn        int    `json:"expires_in"`         // 秒，通常 7200
		RefreshExpiresIn int    `json:"refresh_expires_in"` // 秒，通常 2592000（30天）
		Name             string `json:"name"`
		EnName           string `json:"en_name"`
		AvatarURL        string `json:"avatar_url"`
		OpenID           string `json:"open_id"`
		UnionID          string `json:"union_id"`
		UserID           string `json:"user_id"`
	} `json:"data"`
}

// EncodeState 将 state 结构编码为 base64 字符串
func EncodeState(keyManageID uint, frontendRedirectURI string) string {
	s := OAuthState{
		KeyManageID: keyManageID,
		RedirectURI: frontendRedirectURI,
		Timestamp:   time.Now().Unix(),
	}
	b, _ := json.Marshal(s)
	return base64.URLEncoding.EncodeToString(b)
}

// DecodeState 解码 state 字符串
func DecodeState(state string) (*OAuthState, error) {
	b, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		return nil, fmt.Errorf("state 解码失败: %w", err)
	}
	var s OAuthState
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("state 解析失败: %w", err)
	}
	return &s, nil
}

// BuildOAuthURL 构建飞书 OAuth 授权 URL
// appID: 飞书应用 App ID
// callbackURL: 后端接收授权码的回调地址（需在飞书开放平台注册）
// state: 由 EncodeState 生成
func BuildOAuthURL(appID, callbackURL, state string) string {
	params := url.Values{}
	params.Set("redirect_uri", callbackURL)
	params.Set("app_id", appID)
	params.Set("state", state)
	return feishuBaseURL + "/authen/v1/index?" + params.Encode()
}

// ExchangeCodeForToken 使用授权码换取 user_access_token
// appAccessToken: 应用级 tenant_access_token，用于请求头鉴权
func ExchangeCodeForToken(appAccessToken, code string) (*OAuthTokenResp, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"grant_type": "authorization_code",
		"code":       code,
	})
	body, err := doRequest("POST", feishuBaseURL+"/authen/v1/access_token", appAccessToken, reqBody)
	if err != nil {
		return nil, fmt.Errorf("换取 user_access_token 失败: %w", err)
	}
	var resp OAuthTokenResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析 token 响应失败: %w", err)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("换取 token 失败(code=%d): %s", resp.Code, resp.Msg)
	}
	return &resp, nil
}

// RefreshUserToken 用 refresh_token 刷新 user_access_token
// appAccessToken: 应用级 tenant_access_token，用于请求头鉴权
func RefreshUserToken(appAccessToken, refreshToken string) (*OAuthTokenResp, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	})
	body, err := doRequest("POST", feishuBaseURL+"/authen/v1/refresh_access_token", appAccessToken, reqBody)
	if err != nil {
		return nil, fmt.Errorf("刷新 token 失败: %w", err)
	}
	var resp OAuthTokenResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析刷新 token 响应失败: %w", err)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("刷新 token 失败(code=%d): %s", resp.Code, resp.Msg)
	}
	return &resp, nil
}
