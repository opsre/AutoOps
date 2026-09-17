package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	cmdbDao "dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	configDao "dodevops-api/api/configcenter/dao"
	"dodevops-api/common/result"
	"dodevops-api/common/util"

	"github.com/gin-gonic/gin"
)

type CmdbHostJumpServerServiceInterface interface {
	SyncJumpServerHosts(c *gin.Context, dto *model.SyncJumpServerHostDto) // 同步JumpServer主机
}

type CmdbHostJumpServerServiceImpl struct {
	dao         cmdbDao.CmdbHostDao
	groupDao    cmdbDao.CmdbGroupDao
	ecsAuthDao  configDao.EcsAuthDao
}

var jumpServerServiceInstance *CmdbHostJumpServerServiceImpl

func GetCmdbHostJumpServerService() CmdbHostJumpServerServiceInterface {
	if jumpServerServiceInstance == nil {
		jumpServerServiceInstance = &CmdbHostJumpServerServiceImpl{
			dao:        cmdbDao.NewCmdbHostDao(),
			groupDao:   cmdbDao.NewCmdbGroupDao(),
			ecsAuthDao: configDao.NewEcsAuthDao(),
		}
	}
	return jumpServerServiceInstance
}

// SyncJumpServerHosts 同步JumpServer主机
func (s *CmdbHostJumpServerServiceImpl) SyncJumpServerHosts(c *gin.Context, dto *model.SyncJumpServerHostDto) {
	// 异步处理同步任务
	go func() {
		syncResult := s.syncJumpServerHostsBackground(dto)
		fmt.Printf("JumpServer同步完成: %+v\n", syncResult)
	}()

	// 立即返回响应
	result.Success(c, gin.H{
		"message": "JumpServer主机同步任务已启动，正在后台处理...",
	})
}

// syncJumpServerHostsBackground 后台同步JumpServer主机
func (s *CmdbHostJumpServerServiceImpl) syncJumpServerHostsBackground(dto *model.SyncJumpServerHostDto) *model.SyncJumpServerResult {
	syncResult := &model.SyncJumpServerResult{
		FailedHosts: make([]string, 0),
	}

	// 1. 获取默认业务组(ID=1)
	defaultGroup, err := s.groupDao.GetCmdbGroupById(1)
	if err != nil {
		syncResult.Message = "获取默认业务组失败: " + err.Error()
		return syncResult
	}

	// 2. 查找"免密码认证"的SSH凭据
	var sshKeyID uint = 0
	if auth, err := s.ecsAuthDao.GetEcsAuthByName("免密码认证"); err == nil {
		sshKeyID = auth.ID
		fmt.Printf("找到免密码认证凭据: ID=%d, Name=%s\n", auth.ID, auth.Name)
	} else {
		fmt.Printf("未找到免密码认证凭据，将不设置sshKeyId\n")
	}

	// 3. 调用JumpServer API获取资产列表
	assets, err := s.fetchJumpServerAssets(dto)
	if err != nil {
		syncResult.Message = "获取JumpServer资产失败: " + err.Error()
		return syncResult
	}

	syncResult.Total = len(assets)
	fmt.Printf("从JumpServer获取到 %d 台主机\n", len(assets))

	// 4. 遍历资产，逐个导入
	for _, asset := range assets {
		// 提取节点名称(业务分组)
		var groupName string
		if len(asset.Nodes) > 0 {
			// 取第一个节点的名称作为业务分组
			groupName = asset.Nodes[0].Name
		}

		// 提取SSH端口
		sshPort := 22 // 默认端口
		for _, protocol := range asset.Protocols {
			if protocol.Name == "ssh" {
				sshPort = protocol.Port
				break
			}
		}

		// 5. 检查SSH IP是否已存在
		if existHost := s.dao.GetCmdbHostBySSHIP(asset.Address); existHost != nil {
			fmt.Printf("跳过已存在的主机: %s (%s)\n", asset.Name, asset.Address)
			syncResult.Skipped++
			continue
		}

		// 6. 根据业务分组名称匹配二级分组
		targetGroupID := defaultGroup.ID // 默认使用默认业务组
		if groupName != "" {
			// 尝试按精确名称查找二级分组
			if group, err := s.groupDao.GetCmdbGroupByExactName(groupName); err == nil {
				// 检查是否为二级分组
				if isSecondLevel, _ := s.groupDao.IsSecondLevelGroup(group.ID); isSecondLevel {
					targetGroupID = group.ID
					fmt.Printf("匹配到二级分组: %s (ID=%d)\n", groupName, group.ID)
				} else {
					fmt.Printf("分组 %s 不是二级分组，使用默认分组\n", groupName)
				}
			} else {
				fmt.Printf("未找到匹配的分组: %s，使用默认分组\n", groupName)
			}
		}

		// 7. 创建主机记录
		host := model.CmdbHost{
			HostName:   asset.Name,
			Name:       asset.Name,
			GroupID:    targetGroupID,
			SSHIP:      asset.Address,
			SSHName:    "root", // 默认root用户
			SSHPort:    sshPort,
			SSHKeyID:   sshKeyID, // 关联免密码认证凭据
			Remark:     fmt.Sprintf("JumpServer: %s", asset.Comment),
			CreateTime: util.HTime{Time: time.Now()},
			Vendor:     1,      // 默认自建主机
			Status:     2,      // 未认证
		}

		if err := s.dao.CreateCmdbHost(&host); err != nil {
			fmt.Printf("创建主机失败: %s, 错误: %v\n", asset.Name, err)
			syncResult.Failed++
			syncResult.FailedHosts = append(syncResult.FailedHosts, asset.Name)
			continue
		}

		syncResult.Success++
		fmt.Printf("成功导入主机: %s (%s) -> 分组ID=%d, SSHKeyID=%d\n", asset.Name, asset.Address, targetGroupID, sshKeyID)
	}

	// 生成同步结果消息
	syncResult.Message = fmt.Sprintf("同步完成: 总数=%d, 成功=%d, 失败=%d, 跳过=%d",
		syncResult.Total, syncResult.Success, syncResult.Failed, syncResult.Skipped)

	return syncResult
}

// fetchJumpServerAssets 从JumpServer API获取资产列表
func (s *CmdbHostJumpServerServiceImpl) fetchJumpServerAssets(dto *model.SyncJumpServerHostDto) ([]model.JumpServerAssetDto, error) {
	assets := make([]model.JumpServerAssetDto, 0)

	// 获取认证Token
	token, err := s.getJumpServerToken(dto)
	if err != nil {
		return nil, fmt.Errorf("获取认证Token失败: %v", err)
	}

	// 构造API URL - 注意这里使用 /api/v1/assets/hosts/ 而不是 /api/v1/assets/assets/
	baseURL := strings.TrimSuffix(dto.URL, "/")
	apiURL := fmt.Sprintf("%s/api/v1/assets/hosts/?display=1", baseURL)

	fmt.Printf("请求JumpServer API: %s\n", apiURL)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置认证头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API返回错误状态码 %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// JumpServer可能返回数组或分页对象,需要判断
	var rawResp interface{}
	if err := json.Unmarshal(body, &rawResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 判断是数组还是对象
	switch data := rawResp.(type) {
	case []interface{}:
		// 返回的是数组格式,直接解析
		if err := json.Unmarshal(body, &assets); err != nil {
			return nil, fmt.Errorf("解析资产数组失败: %v", err)
		}
		fmt.Printf("获取到 %d 台主机(数组格式)\n", len(assets))
	case map[string]interface{}:
		// 返回的是分页对象格式
		if results, ok := data["results"].([]interface{}); ok {
			resultsBytes, _ := json.Marshal(results)
			if err := json.Unmarshal(resultsBytes, &assets); err != nil {
				return nil, fmt.Errorf("解析资产分页结果失败: %v", err)
			}
			count := int(data["count"].(float64))
			fmt.Printf("获取到 %d/%d 台主机(分页格式)\n", len(assets), count)
		} else {
			return nil, fmt.Errorf("无法解析分页结果")
		}
	default:
		return nil, fmt.Errorf("未知的响应格式")
	}

	return assets, nil
}

// getJumpServerToken 获取JumpServer认证Token
// 支持两种认证方式:
// 1. AccessKey/SecretKey (UUID格式,长度>30) - 直接使用SecretKey作为Token
// 2. 用户名/密码 (普通字符串) - 先调用登录接口获取Token
func (s *CmdbHostJumpServerServiceImpl) getJumpServerToken(dto *model.SyncJumpServerHostDto) (string, error) {
	// 判断认证方式: 如果AccessKey或SecretKey包含'-'且长度>30,则认为是AccessKey认证
	isAccessKey := (len(dto.AccessKey) > 30 && strings.Contains(dto.AccessKey, "-")) ||
		(len(dto.SecretKey) > 30 && strings.Contains(dto.SecretKey, "-"))

	if isAccessKey {
		// AccessKey认证: 直接使用SecretKey作为Token
		fmt.Printf("使用AccessKey认证\n")
		return dto.SecretKey, nil
	}

	// 用户名/密码认证: 调用登录接口获取Token
	fmt.Printf("使用用户名/密码认证\n")
	baseURL := strings.TrimSuffix(dto.URL, "/")
	loginURL := fmt.Sprintf("%s/api/v1/authentication/tokens/", baseURL)

	// 构造登录请求体
	loginData := map[string]string{
		"username": dto.AccessKey,
		"password": dto.SecretKey,
	}
	loginBody, _ := json.Marshal(loginData)

	// 创建登录请求
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(string(loginBody)))
	if err != nil {
		return "", fmt.Errorf("创建登录请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送登录请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("登录请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("登录失败(状态码 %d): %s", resp.StatusCode, string(body))
	}

	// 解析Token
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取登录响应失败: %v", err)
	}

	var tokenResp map[string]interface{}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("解析Token响应失败: %v", err)
	}

	token, ok := tokenResp["token"].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("无法获取Token")
	}

	fmt.Printf("登录成功,获取到Token: %s\n", token)
	return token, nil
}
