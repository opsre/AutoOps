package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"dodevops-api/api/configcenter/dao"
	"dodevops-api/api/configcenter/model"
)

type AccountAuthService struct {
	dao *dao.AccountAuthDao
}

func NewAccountAuthService() *AccountAuthService {
	return &AccountAuthService{
		dao: dao.NewAccountAuthDao(),
	}
}

// Create 创建账号
func (s *AccountAuthService) Create(account *model.AccountAuth) error {
	// 检查alias是否已存在
	existing, err := s.GetByAlias(account.Alias)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("alias '%s' 已存在", account.Alias)
	}
	return s.dao.Create(account)
}

// Update 更新账号
func (s *AccountAuthService) Update(account *model.AccountAuth) error {
	return s.dao.Update(account)
}

// Delete 删除账号
func (s *AccountAuthService) Delete(id uint) error {
	return s.dao.Delete(id)
}

// GetByID 根据ID查询账号
func (s *AccountAuthService) GetByID(id uint) (*model.AccountAuth, error) {
	return s.dao.GetByID(id)
}

// List 获取账号列表
func (s *AccountAuthService) List() ([]model.AccountAuth, error) {
	return s.dao.List()
}

// ListWithPage 获取账号列表（分页）
func (s *AccountAuthService) ListWithPage(page, pageSize int) ([]model.AccountAuth, int64, error) {
	return s.dao.ListWithPage(page, pageSize)
}

// DecryptPassword 解密密码
func (s *AccountAuthService) DecryptPassword(id uint) (string, error) {
	account, err := s.dao.GetByID(id)
	if err != nil {
		return "", err
	}
	return account.DecryptPassword()
}

// GetByType 根据类型查询账号
func (s *AccountAuthService) GetByType(accountType string) ([]model.AccountAuth, error) {
	return s.dao.GetByType(accountType)
}

// GetByAlias 根据别名查询账号
func (s *AccountAuthService) GetByAlias(alias string) (*model.AccountAuth, error) {
	return s.dao.GetByAlias(alias)
}

// TestConnection 测试账号连接
func (s *AccountAuthService) TestConnection(id uint) (bool, string, error) {
	account, err := s.dao.GetByID(id)
	if err != nil {
		return false, "", fmt.Errorf("获取账号失败: %w", err)
	}

	password, err := account.DecryptPassword()
	if err != nil {
		return false, "", fmt.Errorf("解密密码失败: %w", err)
	}

	switch account.Type {
	case 4: // Jenkins
		return testJenkinsConnection(account.Host, account.Port, account.Name, password)
	case 5: // Zabbix
		return testZabbixConnection(account.Host, account.Port, account.Name, password)
	default:
		return false, "不支持的账号类型", nil
	}
}

// testJenkinsConnection 测试Jenkins连接
func testJenkinsConnection(host string, port int, username, password string) (bool, string, error) {
	var baseURL string
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		baseURL = strings.TrimRight(host, "/")
	} else if port == 443 {
		baseURL = fmt.Sprintf("https://%s", host)
	} else if port > 0 {
		baseURL = fmt.Sprintf("http://%s:%d", host, port)
	} else {
		baseURL = fmt.Sprintf("http://%s:8080", host)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", baseURL+"/api/json", nil)
	if err != nil {
		return false, "创建请求失败", err
	}

	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return false, "连接失败: " + err.Error(), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, "连接成功", nil
	} else if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return false, "认证失败: 用户名或密码错误", nil
	} else {
		return false, fmt.Sprintf("连接失败: HTTP状态码 %d", resp.StatusCode), nil
	}
}

// testZabbixConnection 测试Zabbix连接
func testZabbixConnection(host string, port int, username, password string) (bool, string, error) {
	var baseURL string
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		baseURL = strings.TrimRight(host, "/")
		// 移除可能的 /zabbix.php 后缀
		baseURL = strings.TrimSuffix(baseURL, "/zabbix.php")
	} else if port == 443 {
		baseURL = fmt.Sprintf("https://%s", host)
	} else if port > 0 {
		baseURL = fmt.Sprintf("http://%s:%d", host, port)
	} else {
		baseURL = fmt.Sprintf("http://%s", host)
	}

	// Zabbix API登录请求
	loginData := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "user.login",
		"params": map[string]string{
			"user":     username,
			"password": password,
		},
		"id":   1,
		"auth": nil,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return false, "构建请求数据失败", err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", baseURL+"/api_jsonrpc.php", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, "创建请求失败", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return false, "连接失败: " + err.Error(), nil
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, "解析响应失败", err
	}

	if authToken, ok := result["result"].(string); ok && authToken != "" {
		return true, "连接成功", nil
	} else if errData, ok := result["error"].(map[string]interface{}); ok {
		if message, ok := errData["data"].(string); ok {
			return false, "认证失败: " + message, nil
		}
		return false, "认证失败: 用户名或密码错误", nil
	}

	return false, "未知错误", nil
}
