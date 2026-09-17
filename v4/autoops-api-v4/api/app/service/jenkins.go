package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"dodevops-api/api/configcenter/dao"
	"dodevops-api/api/app/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// IJenkinsService Jenkins服务接口
type IJenkinsService interface {
	// 服务器管理
	GetJenkinsServers(c *gin.Context)
	GetJenkinsServerDetail(c *gin.Context, serverID uint)
	TestJenkinsConnection(c *gin.Context, req *model.TestJenkinsConnectionRequest)

	// 任务管理
	GetJobs(c *gin.Context, serverID uint)
	SearchJobs(c *gin.Context, serverID uint, keyword string)
	GetJobDetail(c *gin.Context, serverID uint, jobName string)
	GetJobParameters(c *gin.Context, serverID uint, jobName string)
	StartJob(c *gin.Context, serverID uint, jobName string, req *model.StartJobRequest)
	StopBuild(c *gin.Context, serverID uint, jobName string, buildNumber int, req *model.StopBuildRequest)

	// 构建管理
	GetBuildDetail(c *gin.Context, serverID uint, jobName string, buildNumber int)
	GetBuildLog(c *gin.Context, serverID uint, jobName string, buildNumber int, req *model.GetBuildLogRequest)

	// 系统信息
	GetSystemInfo(c *gin.Context, serverID uint)
	GetQueueInfo(c *gin.Context, serverID uint)
}

// JenkinsServiceImpl Jenkins服务实现
type JenkinsServiceImpl struct {
	db         *gorm.DB
	accountDao *dao.AccountAuthDao
}

// NewJenkinsService 创建Jenkins服务实例
func NewJenkinsService(db *gorm.DB) IJenkinsService {
	return &JenkinsServiceImpl{
		db:         db,
		accountDao: dao.NewAccountAuthDao(),
	}
}

// JenkinsClient Jenkins客户端
type JenkinsClient struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// NewJenkinsClient 创建Jenkins客户端
func NewJenkinsClient(host string, port int, username, password string) *JenkinsClient {
	var baseURL string

	// 检查 host 是否已经包含协议
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		// host 已包含完整URL，直接使用，但要去除尾随斜杠
		baseURL = strings.TrimRight(host, "/")
	} else if strings.Contains(host, ":") {
		// host 包含端口但没有协议，根据端口判断协议
		if strings.Contains(host, ":443") {
			baseURL = fmt.Sprintf("https://%s", host)
		} else {
			baseURL = fmt.Sprintf("http://%s", host)
		}
	} else {
		// host 不包含端口，使用 port 参数
		if port == 443 {
			baseURL = fmt.Sprintf("https://%s", host)
		} else if port > 0 {
			baseURL = fmt.Sprintf("http://%s:%d", host, port)
		} else {
			baseURL = fmt.Sprintf("http://%s:8080", host)
		}
	}

	return &JenkinsClient{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest 发送HTTP请求
func (jc *JenkinsClient) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := jc.BaseURL + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(jc.Username, jc.Password)
	req.Header.Set("Content-Type", "application/json")

	return jc.HTTPClient.Do(req)
}

// getJenkinsClient 获取Jenkins客户端
func (s *JenkinsServiceImpl) getJenkinsClient(serverID uint) (*JenkinsClient, error) {
	account, err := s.accountDao.GetByID(serverID)
	if err != nil {
		return nil, err
	}

	if account.Type != model.JenkinsAccountType {
		return nil, fmt.Errorf("账号类型不是Jenkins")
	}

	password, err := account.DecryptPassword()
	if err != nil {
		return nil, err
	}

	return NewJenkinsClient(account.Host, account.Port, account.Name, password), nil
}

// GetJenkinsServers 获取Jenkins服务器列表
func (s *JenkinsServiceImpl) GetJenkinsServers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	accounts, total, err := s.accountDao.GetAccountsByType(model.JenkinsAccountType, page, pageSize)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins服务器列表失败: "+err.Error())
		return
	}

	var serverList []model.JenkinsServerInfo
	for _, account := range accounts {
		serverInfo := model.JenkinsServerInfo{
			ID:          account.ID,
			Alias:       account.Alias,
			Host:        account.Host,
			Port:        account.Port,
			Username:    account.Name,
			Description: account.Remark,
			CreatedAt:   account.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   account.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		serverList = append(serverList, serverInfo)
	}

	response := model.JenkinsServerListResponse{
		List:  serverList,
		Total: total,
	}

	result.Success(c, response)
}

// GetJenkinsServerDetail 获取Jenkins服务器详情
func (s *JenkinsServiceImpl) GetJenkinsServerDetail(c *gin.Context, serverID uint) {
	account, err := s.accountDao.GetByID(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins服务器详情失败: "+err.Error())
		return
	}

	if account.Type != model.JenkinsAccountType {
		result.Failed(c, 400, "账号类型不是Jenkins")
		return
	}

	serverInfo := model.JenkinsServerInfo{
		ID:          account.ID,
		Alias:       account.Alias,
		Host:        account.Host,
		Port:        account.Port,
		Username:    account.Name,
		Description: account.Remark,
		CreatedAt:   account.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   account.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	result.Success(c, serverInfo)
}

// TestJenkinsConnection 测试Jenkins连接
func (s *JenkinsServiceImpl) TestJenkinsConnection(c *gin.Context, req *model.TestJenkinsConnectionRequest) {
	// 解析URL获取host和port
	u, err := url.Parse(req.URL)
	if err != nil {
		result.Failed(c, 500, "无效的Jenkins URL: "+err.Error())
		return
	}

	host := u.Hostname()
	port := 80
	if u.Port() != "" {
		port, _ = strconv.Atoi(u.Port())
	} else if u.Scheme == "https" {
		port = 443
	}

	client := NewJenkinsClient(host, port, req.Username, req.Password)

	// 获取系统信息来测试连接
	resp, err := client.makeRequest("GET", "/api/json", nil)
	if err != nil {
		result.Failed(c, 500, "连接Jenkins失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var systemInfo model.JenkinsSystemInfo
	if err := json.NewDecoder(resp.Body).Decode(&systemInfo); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	response := model.TestJenkinsConnectionResponse{
		Success:    true,
		Message:    "连接成功",
		SystemInfo: &systemInfo,
	}

	result.Success(c, response)
}

// GetJobs 获取任务列表
func (s *JenkinsServiceImpl) GetJobs(c *gin.Context, serverID uint) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	resp, err := client.makeRequest("GET", "/api/json?tree=jobs[name,displayName,description,url,buildable,color,_class,lastBuild[number,url,displayName,result,building,duration,timestamp],lastStableBuild[number,url,displayName,result],lastSuccessfulBuild[number,url,displayName,result],lastFailedBuild[number,url,displayName,result]]", nil)
	if err != nil {
		result.Failed(c, 500, "获取任务列表失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var jenkinsData struct {
		Jobs []model.JenkinsJob `json:"jobs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jenkinsData); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.JenkinsJobListResponse{
		Jobs:   jenkinsData.Jobs,
		Total:  len(jenkinsData.Jobs),
		Server: serverName,
	}

	result.Success(c, response)
}

// SearchJobs 搜索Jenkins任务
func (s *JenkinsServiceImpl) SearchJobs(c *gin.Context, serverID uint, keyword string) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	// 获取所有任务
	resp, err := client.makeRequest("GET", "/api/json?tree=jobs[name,displayName,description,url,buildable,color,_class,lastBuild[number,url,displayName,result,building,duration,timestamp],lastStableBuild[number,url,displayName,result],lastSuccessfulBuild[number,url,displayName,result],lastFailedBuild[number,url,displayName,result]]", nil)
	if err != nil {
		result.Failed(c, 500, "获取任务列表失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var jenkinsData struct {
		Jobs []model.JenkinsJob `json:"jobs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jenkinsData); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	// 过滤任务 - 支持名称、显示名称、描述的模糊匹配
	var filteredJobs []model.JenkinsJob
	keyword = strings.ToLower(keyword)

	for _, job := range jenkinsData.Jobs {
		// 检查任务名称
		if strings.Contains(strings.ToLower(job.Name), keyword) {
			filteredJobs = append(filteredJobs, job)
			continue
		}
		// 检查显示名称
		if strings.Contains(strings.ToLower(job.DisplayName), keyword) {
			filteredJobs = append(filteredJobs, job)
			continue
		}
		// 检查描述
		if strings.Contains(strings.ToLower(job.Description), keyword) {
			filteredJobs = append(filteredJobs, job)
			continue
		}
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.JenkinsJobListResponse{
		Jobs:   filteredJobs,
		Total:  len(filteredJobs),
		Server: serverName,
	}

	result.Success(c, response)
}

// GetJobDetail 获取任务详情
func (s *JenkinsServiceImpl) GetJobDetail(c *gin.Context, serverID uint, jobName string) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	// 获取任务详情
	jobURL := fmt.Sprintf("/job/%s/api/json", url.PathEscape(jobName))
	resp, err := client.makeRequest("GET", jobURL, nil)
	if err != nil {
		result.Failed(c, 500, "获取任务详情失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var job model.JenkinsJob
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	// 获取构建历史
	buildsURL := fmt.Sprintf("/job/%s/api/json?tree=builds[number,url,displayName,result,building,duration,timestamp,keepLog,queueId]", url.PathEscape(jobName))
	resp, err = client.makeRequest("GET", buildsURL, nil)
	if err != nil {
		result.Failed(c, 500, "获取构建历史失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var buildsData struct {
		Builds []model.JenkinsBuild `json:"builds"`
	}

	if resp.StatusCode == http.StatusOK {
		json.NewDecoder(resp.Body).Decode(&buildsData)
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.JenkinsJobDetailResponse{
		Job:    job,
		Builds: buildsData.Builds,
		Server: serverName,
	}

	result.Success(c, response)
}

// GetJobParameters 获取任务参数列表
func (s *JenkinsServiceImpl) GetJobParameters(c *gin.Context, serverID uint, jobName string) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	// 获取任务配置，包含参数定义
	// 先获取基本的参数信息
	jobURL := fmt.Sprintf("/job/%s/api/json", url.PathEscape(jobName))
	resp, err := client.makeRequest("GET", jobURL, nil)
	if err != nil {
		result.Failed(c, 500, "获取任务参数失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d, 响应内容: %s", resp.StatusCode, string(bodyBytes[:min(200, len(bodyBytes))])))
		return
	}

	// 先读取响应体用于调试
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Failed(c, 500, "读取Jenkins响应失败: "+err.Error())
		return
	}

	// 检查是否是HTML响应
	bodyStr := string(bodyBytes)
	if strings.HasPrefix(strings.TrimSpace(bodyStr), "<") {
		result.Failed(c, 500, fmt.Sprintf("Jenkins返回了HTML而不是JSON，可能是任务不存在或认证失败。响应前200字符: %s", bodyStr[:min(200, len(bodyStr))]))
		return
	}

	var jobData struct {
		Actions []struct {
			Class                string                   `json:"_class"`
			ParameterDefinitions []map[string]interface{} `json:"parameterDefinitions"`
		} `json:"actions"`
		Property []struct {
			Class                string                   `json:"_class"`
			ParameterDefinitions []map[string]interface{} `json:"parameterDefinitions"`
		} `json:"property"`
	}

	if err := json.Unmarshal(bodyBytes, &jobData); err != nil {
		result.Failed(c, 500, fmt.Sprintf("解析Jenkins响应失败: %s, 响应前200字符: %s", err.Error(), bodyStr[:min(200, len(bodyStr))]))
		return
	}

	// 获取 XML 配置以获取扩展参数的完整信息
	configURL := fmt.Sprintf("/job/%s/config.xml", url.PathEscape(jobName))
	configResp, err := client.makeRequest("GET", configURL, nil)
	var extendedParamConfig map[string][]string // 存储参数名 -> choices 的映射
	extendedParamConfig = make(map[string][]string)

	if err == nil && configResp.StatusCode == http.StatusOK {
		defer configResp.Body.Close()
		bodyBytes, _ := io.ReadAll(configResp.Body)
		configXML := string(bodyBytes)

		// 解析 XML 中的 Extended Choice Parameter 配置
		parseExtendedChoiceParams(configXML, extendedParamConfig)
	}

	// 提取参数定义（从 actions 或 property 中）
	var parameters []model.BuildParameter

	// 先从 actions 中查找（Freestyle 项目）
	for _, action := range jobData.Actions {
		if action.Class == "hudson.model.ParametersDefinitionProperty" {
			for _, paramDef := range action.ParameterDefinitions {
				param := buildParameterFromMap(paramDef, extendedParamConfig)
				parameters = append(parameters, param)
			}
			break
		}
	}

	// 如果 actions 中没找到，从 property 中查找（Pipeline 项目）
	if len(parameters) == 0 {
		for _, prop := range jobData.Property {
			if prop.Class == "hudson.model.ParametersDefinitionProperty" {
				for _, paramDef := range prop.ParameterDefinitions {
					param := buildParameterFromMap(paramDef, extendedParamConfig)
					parameters = append(parameters, param)
				}
				break
			}
		}
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.JobParametersResponse{
		JobName:    jobName,
		Parameters: parameters,
		Server:     serverName,
		Total:      len(parameters),
	}

	result.Success(c, response)
}

// extractParameterType 从完整的类名中提取参数类型
func extractParameterType(fullType string) string {
	// 提取类名的最后一部分
	// 例如: "hudson.model.StringParameterDefinition" -> "String"
	parts := strings.Split(fullType, ".")
	if len(parts) > 0 {
		className := parts[len(parts)-1]
		// 移除 "ParameterDefinition" 后缀
		className = strings.TrimSuffix(className, "ParameterDefinition")
		return className
	}
	return fullType
}

// parseExtendedChoiceParams 解析 XML 配置中的 Extended Choice Parameter
func parseExtendedChoiceParams(configXML string, result map[string][]string) {
	// 查找所有的 ExtendedChoiceParameterDefinition 块
	// XML 格式: <com.cwctravel.hudson.plugins.extended__choice__parameter.ExtendedChoiceParameterDefinition>

	searchPattern := "extended__choice__parameter.ExtendedChoiceParameterDefinition"
	xmlCopy := configXML

	for {
		// 查找参数定义块的开始位置
		startIdx := strings.Index(xmlCopy, searchPattern)
		if startIdx == -1 {
			break
		}

		// 往回找到完整的开标签 <com.cwctravel...>
		tagStart := strings.LastIndex(xmlCopy[:startIdx], "<com.cwctravel")
		if tagStart == -1 {
			xmlCopy = xmlCopy[startIdx+len(searchPattern):]
			continue
		}

		// 找到对应的结束标签
		endTagPattern := "</com.cwctravel.hudson.plugins.extended__choice__parameter.ExtendedChoiceParameterDefinition>"
		endIdx := strings.Index(xmlCopy[tagStart:], endTagPattern)
		if endIdx == -1 {
			xmlCopy = xmlCopy[startIdx+len(searchPattern):]
			continue
		}

		// 提取整个参数定义块
		paramBlock := xmlCopy[tagStart : tagStart+endIdx+len(endTagPattern)]

		// 提取参数名
		nameStart := strings.Index(paramBlock, "<name>")
		nameEnd := strings.Index(paramBlock, "</name>")
		if nameStart != -1 && nameEnd != -1 {
			paramName := paramBlock[nameStart+6 : nameEnd]

			// 提取 value 字段（包含逗号分隔的选项）
			valueStart := strings.Index(paramBlock, "<value>")
			valueEnd := strings.Index(paramBlock, "</value>")
			if valueStart != -1 && valueEnd != -1 {
				valueStr := paramBlock[valueStart+7 : valueEnd]
				if valueStr != "" {
					// 提取分隔符（默认为逗号）
					delimiter := ","
					delimiterStart := strings.Index(paramBlock, "<multiSelectDelimiter>")
					delimiterEnd := strings.Index(paramBlock, "</multiSelectDelimiter>")
					if delimiterStart != -1 && delimiterEnd != -1 {
						customDelimiter := paramBlock[delimiterStart+22 : delimiterEnd]
						if customDelimiter != "" {
							delimiter = customDelimiter
						}
					}

					// 智能检测：如果配置的分隔符在valueStr中不存在，尝试使用英文逗号
					if !strings.Contains(valueStr, delimiter) {
						if strings.Contains(valueStr, ",") {
							delimiter = ","
						}
					}

					// 按分隔符分隔选项
					choices := strings.Split(valueStr, delimiter)
					for i, choice := range choices {
						choices[i] = strings.TrimSpace(choice)
					}

					result[paramName] = choices
				}
			}
		}

		// 继续查找下一个
		xmlCopy = xmlCopy[tagStart+endIdx+len(endTagPattern):]
	}
}

// buildParameterFromMap 从参数定义 map 构建参数对象
func buildParameterFromMap(paramDef map[string]interface{}, extendedParamConfig map[string][]string) model.BuildParameter {
	// 提取基本字段
	name, _ := paramDef["name"].(string)
	paramTypeRaw, _ := paramDef["type"].(string)
	description, _ := paramDef["description"].(string)

	// 提取参数类型
	paramType := extractParameterType(paramTypeRaw)

	param := model.BuildParameter{
		Name:        name,
		Type:        paramType,
		Description: description,
		Required:    isParameterRequired(paramType),
	}

	// 提取 choices（用于 Choice、Checkbox 等参数）
	// 尝试多个可能的字段名
	choicesFound := false

	// 1. 从 XML 配置中查找（针对 Extended Choice Parameter）
	if xmlChoices, ok := extendedParamConfig[name]; ok && len(xmlChoices) > 0 {
		param.Choices = xmlChoices
		choicesFound = true
	}

	// 2. 标准的 choices 字段
	if !choicesFound {
		if choicesRaw, ok := paramDef["choices"]; ok && choicesRaw != nil {
			if choicesArray, ok := choicesRaw.([]interface{}); ok {
				choices := make([]string, 0, len(choicesArray))
				for _, choice := range choicesArray {
					if choiceStr, ok := choice.(string); ok {
						// PT_CHECKBOX 或 ExtendedChoiceParameter 类型的 choices 可能是逗号分隔的字符串，需要进一步分割
						if (paramType == "PT_CHECKBOX" || strings.Contains(paramTypeRaw, "ExtendedChoice")) && strings.Contains(choiceStr, ",") {
							subChoices := strings.Split(choiceStr, ",")
							for _, subChoice := range subChoices {
								trimmed := strings.TrimSpace(subChoice)
								if trimmed != "" {
									choices = append(choices, trimmed)
								}
							}
						} else {
							choices = append(choices, choiceStr)
						}
					}
				}
				param.Choices = choices
				choicesFound = true
			}
		}
	}

	// 3. PT_CHECKBOX 插件的 format 字段（格式: "opt1,opt2,opt3"）
	if !choicesFound {
		if formatRaw, ok := paramDef["format"]; ok {
			if formatStr, ok := formatRaw.(string); ok && formatStr != "" {
				choices := strings.Split(formatStr, ",")
				// 去除空格
				for i, choice := range choices {
					choices[i] = strings.TrimSpace(choice)
				}
				param.Choices = choices
				choicesFound = true
			}
		}
	}

	// 提取 defaultParameterValue
	if defaultValueRaw, ok := paramDef["defaultParameterValue"]; ok {
		if defaultValueMap, ok := defaultValueRaw.(map[string]interface{}); ok {
			if value, ok := defaultValueMap["value"]; ok {
				param.DefaultValue = normalizeParameterValue(paramType, value)
			}
		}
	}

	// 如果没有默认值，根据类型设置合理的空值
	if param.DefaultValue == nil {
		param.DefaultValue = getDefaultValueForType(paramType)
	}

	return param
}

// buildParameterFromDefinition 从参数定义构建参数对象（保留用于向后兼容）
func buildParameterFromDefinition(paramDef struct {
	Name                 string `json:"name"`
	Type                 string `json:"type"`
	Description          string `json:"description"`
	DefaultParameterValue struct {
		Value interface{} `json:"value"`
	} `json:"defaultParameterValue"`
	Choices []string `json:"choices"`
}) model.BuildParameter {
	paramType := extractParameterType(paramDef.Type)

	param := model.BuildParameter{
		Name:        paramDef.Name,
		Type:        paramType,
		Description: paramDef.Description,
		Choices:     paramDef.Choices,
		Required:    isParameterRequired(paramType),
	}

	// 处理默认值 - 根据参数类型进行类型转换
	if paramDef.DefaultParameterValue.Value != nil {
		param.DefaultValue = normalizeParameterValue(paramType, paramDef.DefaultParameterValue.Value)
	} else {
		// 如果没有默认值，根据类型设置合理的空值
		param.DefaultValue = getDefaultValueForType(paramType)
	}

	return param
}

// isParameterRequired 判断参数是否必填
func isParameterRequired(paramType string) bool {
	// 文件参数和密码参数通常不是必填的
	switch paramType {
	case model.ParamTypeFile, model.ParamTypePassword:
		return false
	default:
		return false // Jenkins 参数默认都不是必填，除非在 Pipeline 中明确指定
	}
}

// normalizeParameterValue 标准化参数值
func normalizeParameterValue(paramType string, value interface{}) interface{} {
	switch paramType {
	case model.ParamTypeBoolean:
		// 确保布尔值类型正确
		if v, ok := value.(bool); ok {
			return v
		}
		// 尝试从字符串转换
		if v, ok := value.(string); ok {
			return v == "true" || v == "True" || v == "TRUE"
		}
		return false

	case model.ParamTypeString, model.ParamTypeChoice, model.ParamTypeText, model.ParamTypePassword:
		// 字符串类型
		if v, ok := value.(string); ok {
			return v
		}
		return fmt.Sprintf("%v", value)

	default:
		return value
	}
}

// getDefaultValueForType 根据参数类型获取默认值
func getDefaultValueForType(paramType string) interface{} {
	switch paramType {
	case model.ParamTypeBoolean:
		return false
	case model.ParamTypeString, model.ParamTypeChoice, model.ParamTypeText, model.ParamTypePassword:
		return ""
	default:
		return nil
	}
}

// StartJob 启动任务
func (s *JenkinsServiceImpl) StartJob(c *gin.Context, serverID uint, jobName string, req *model.StartJobRequest) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	var buildURL string
	if req.Parameters != nil && len(req.Parameters) > 0 {
		// 带参数构建
		buildURL = fmt.Sprintf("/job/%s/buildWithParameters", url.PathEscape(jobName))

		// 构建表单数据
		formData := url.Values{}
		for key, value := range req.Parameters {
			formData.Set(key, value)
		}

		resp, err := client.makeRequest("POST", buildURL+"?"+formData.Encode(), nil)
		if err != nil {
			result.Failed(c, 500, "启动任务失败: "+err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
			return
		}
	} else {
		// 无参数构建
		buildURL = fmt.Sprintf("/job/%s/build", url.PathEscape(jobName))

		resp, err := client.makeRequest("POST", buildURL, nil)
		if err != nil {
			result.Failed(c, 500, "启动任务失败: "+err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
			return
		}
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.StartJobResponse{
		Success: true,
		Message: "任务启动成功",
		JobName: jobName,
		Server:  serverName,
	}

	result.Success(c, response)
}

// StopBuild 停止构建
func (s *JenkinsServiceImpl) StopBuild(c *gin.Context, serverID uint, jobName string, buildNumber int, req *model.StopBuildRequest) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	stopURL := fmt.Sprintf("/job/%s/%d/stop", url.PathEscape(jobName), buildNumber)
	resp, err := client.makeRequest("POST", stopURL, nil)
	if err != nil {
		result.Failed(c, 500, "停止构建失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.StopBuildResponse{
		Success:     true,
		Message:     "构建停止成功",
		JobName:     jobName,
		BuildNumber: buildNumber,
		Server:      serverName,
	}

	result.Success(c, response)
}

// GetBuildDetail 获取构建详情
func (s *JenkinsServiceImpl) GetBuildDetail(c *gin.Context, serverID uint, jobName string, buildNumber int) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	buildURL := fmt.Sprintf("/job/%s/%d/api/json", url.PathEscape(jobName), buildNumber)
	resp, err := client.makeRequest("GET", buildURL, nil)
	if err != nil {
		result.Failed(c, 500, "获取构建详情失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var build model.JenkinsBuild
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.JenkinsBuildDetailResponse{
		Build:  build,
		Server: serverName,
	}

	result.Success(c, response)
}

// GetBuildLog 获取构建日志
func (s *JenkinsServiceImpl) GetBuildLog(c *gin.Context, serverID uint, jobName string, buildNumber int, req *model.GetBuildLogRequest) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	logURL := fmt.Sprintf("/job/%s/%d/logText/progressiveText?start=%d", url.PathEscape(jobName), buildNumber, req.Start)
	resp, err := client.makeRequest("GET", logURL, nil)
	if err != nil {
		result.Failed(c, 500, "获取构建日志失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	logBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Failed(c, 500, "读取日志内容失败: "+err.Error())
		return
	}

	// 检查是否有更多日志
	moreData := resp.Header.Get("X-More-Data") == "true"
	textSize, _ := strconv.Atoi(resp.Header.Get("X-Text-Size"))

	// 获取服务器信息
	account, _ := s.accountDao.GetByID(serverID)
	serverName := account.Alias

	response := model.GetBuildLogResponse{
		Log:         string(logBytes),
		HasMore:     moreData,
		TextSize:    textSize,
		MoreData:    moreData,
		JobName:     jobName,
		BuildNumber: buildNumber,
		Server:      serverName,
	}

	result.Success(c, response)
}

// GetSystemInfo 获取系统信息
func (s *JenkinsServiceImpl) GetSystemInfo(c *gin.Context, serverID uint) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	resp, err := client.makeRequest("GET", "/api/json", nil)
	if err != nil {
		result.Failed(c, 500, "获取系统信息失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var systemInfo model.JenkinsSystemInfo
	if err := json.NewDecoder(resp.Body).Decode(&systemInfo); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	result.Success(c, systemInfo)
}

// GetQueueInfo 获取队列信息
func (s *JenkinsServiceImpl) GetQueueInfo(c *gin.Context, serverID uint) {
	client, err := s.getJenkinsClient(serverID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins客户端失败: "+err.Error())
		return
	}

	resp, err := client.makeRequest("GET", "/queue/api/json", nil)
	if err != nil {
		result.Failed(c, 500, "获取队列信息失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误: %d", resp.StatusCode))
		return
	}

	var queueInfo model.JenkinsQueue
	if err := json.NewDecoder(resp.Body).Decode(&queueInfo); err != nil {
		result.Failed(c, 500, "解析Jenkins响应失败: "+err.Error())
		return
	}

	result.Success(c, queueInfo)
}