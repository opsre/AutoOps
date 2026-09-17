package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"dodevops-api/api/app/dao"
	"dodevops-api/api/app/model"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	ccdao "dodevops-api/api/configcenter/dao"
	ccModel "dodevops-api/api/configcenter/model"
	"dodevops-api/common"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

// IShReleaseService 脚本上线服务接口
type IShReleaseService interface {
	// 创建脚本上线工单
	CreateShRelease(c *gin.Context, req *model.CreateShReleaseRequest)
	// 获取脚本上线工单列表
	GetShReleaseList(c *gin.Context, req *model.ShReleaseListRequest)
	// 获取脚本上线工单详情
	GetShReleaseDetail(c *gin.Context, id uint)
	// 删除脚本上线工单
	DeleteShRelease(c *gin.Context, id uint)

	// 审核操作
	ApproveShRelease(c *gin.Context, req *model.ShApprovalRequest)

	// 拉取代码（通过Jenkins）
	PullCode(c *gin.Context, req *model.ShExecuteRequest)

	// 执行脚本（通过SSH登录服务器执行）
	ExecuteScript(c *gin.Context, req *model.ShExecuteScriptRequest)

	// 获取Jenkins日志（代码拉取日志）
	GetShReleaseLog(c *gin.Context, releaseID uint, start int)

	// 获取脚本执行输出
	GetScriptOutput(c *gin.Context, releaseID uint)

	// 列出脚本上线服务器上的文件
	ListReleaseFiles(c *gin.Context, releaseID uint, remotePath string)

	// 下载脚本上线服务器上的文件
	DownloadReleaseFile(c *gin.Context, releaseID uint, remotePath string)

	// 检查用户操作权限
	CheckUserPermission(c *gin.Context, releaseID uint, operationType string)
}

// ShReleaseService 脚本上线服务实现
type ShReleaseService struct {
	dao    dao.IShReleaseDao
	appDao dao.IApplicationDao
	db     *gorm.DB
}

// NewShReleaseService 创建脚本上线服务
func NewShReleaseService(db *gorm.DB) IShReleaseService {
	return &ShReleaseService{
		dao:    dao.NewShReleaseDao(db),
		appDao: dao.NewApplicationDao(db),
		db:     db,
	}
}

// CreateShRelease 创建脚本上线工单
func (s *ShReleaseService) CreateShRelease(c *gin.Context, req *model.CreateShReleaseRequest) {
	// TODO: 从JWT获取当前用户信息
	applicantID := uint(1)
	applicantName := "管理员"

	// 获取应用信息
	app, err := s.appDao.GetApplicationByID(req.AppID)
	if err != nil {
		result.Failed(c, 400, fmt.Sprintf("应用ID %d 不存在", req.AppID))
		return
	}

	// 如果前端未指定环境，默认使用 prod
	environment := req.Environment
	if environment == "" {
		environment = "prod"
	}

	// 根据环境获取 Jenkins 配置
	envs, err := s.appDao.GetJenkinsEnvsByAppID(req.AppID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins环境配置失败: "+err.Error())
		return
	}

	var targetEnv *model.JenkinsEnv
	for _, env := range envs {
		if env.EnvName == environment {
			targetEnv = &env
			break
		}
	}

	if targetEnv == nil {
		result.Failed(c, 400, fmt.Sprintf("该应用未配置 %s 环境，无法创建脚本上线工单", environment))
		return
	}

	// 验证Jenkins环境配置完整性
	if targetEnv.JenkinsServerID == nil || targetEnv.JobName == "" {
		result.Failed(c, 400, fmt.Sprintf("该应用的 %s 环境Jenkins配置不完整", environment))
		return
	}

	// 获取审核人姓名
	approverName := s.getApproverName(req.ApproverID)
	executorName := s.getApproverName(req.ExecutorID)

	// 将参数map转换为JSON字符串
	var parametersJSON string
	if len(req.Parameters) > 0 {
		if jsonBytes, err := json.Marshal(req.Parameters); err == nil {
			parametersJSON = string(jsonBytes)
		}
	}

	// 构建脚本上线工单对象
	release := &model.ShRelease{
		Title:           req.Title,
		Reason:          req.Reason,
		BusinessGroupID: req.BusinessGroupID,
		AppID:           req.AppID,
		AppName:         app.Name,
		AppCode:         app.Code,
		ApplicantID:     applicantID,
		ApplicantName:   applicantName,
		ApproverID:      req.ApproverID,
		ApproverName:    approverName,
		ExecutorID:      req.ExecutorID,
		ExecutorName:    executorName,
		ServerHostID:    req.ServerHostID, // 执行服务器ID
		ExecuteDir:      req.ExecuteDir,
		ScriptContent:   req.ScriptContent,
		Parameters:      parametersJSON, // 保存参数JSON字符串
		ApprovalStatus:  1,              // 未审核
		ExecuteStatus:   1,              // 未执行
		Status:          1,              // 流程进行中
		JenkinsEnvID:    &targetEnv.ID,
	}

	// 创建工单
	if err := s.dao.CreateShRelease(release); err != nil {
		result.Failed(c, 500, "创建脚本上线工单失败: "+err.Error())
		return
	}

	// 返回创建的工单
	createdRelease, _ := s.dao.GetShReleaseByID(release.ID)
	result.Success(c, createdRelease)
}

// GetShReleaseList 获取脚本上线工单列表
func (s *ShReleaseService) GetShReleaseList(c *gin.Context, req *model.ShReleaseListRequest) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	releases, total, err := s.dao.GetShReleaseList(req)
	if err != nil {
		result.Failed(c, 500, "获取脚本上线列表失败: "+err.Error())
		return
	}

	response := &model.ShReleaseListResponse{
		Total: int(total),
		List:  releases,
	}

	result.Success(c, response)
}

// GetShReleaseDetail 获取脚本上线工单详情
func (s *ShReleaseService) GetShReleaseDetail(c *gin.Context, id uint) {
	release, err := s.dao.GetShReleaseByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "脚本上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取脚本上线工单详情失败: "+err.Error())
		return
	}

	result.Success(c, release)
}

// DeleteShRelease 删除脚本上线工单
func (s *ShReleaseService) DeleteShRelease(c *gin.Context, id uint) {
	// 检查工单是否存在
	release, err := s.dao.GetShReleaseByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "脚本上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取脚本上线工单失败: "+err.Error())
		return
	}

	// 检查工单状态，执行中的不允许删除
	if release.ExecuteStatus == 2 {
		result.Failed(c, 400, "执行中的工单不能删除")
		return
	}

	// 删除工单
	if err := s.dao.DeleteShRelease(id); err != nil {
		result.Failed(c, 500, "删除脚本上线工单失败: "+err.Error())
		return
	}

	result.Success(c, "删除成功")
}

// ApproveShRelease 审核操作
func (s *ShReleaseService) ApproveShRelease(c *gin.Context, req *model.ShApprovalRequest) {
	// 获取工单
	release, err := s.dao.GetShReleaseByID(req.ReleaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "脚本上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取脚本上线工单失败: "+err.Error())
		return
	}

	// 检查工单状态
	if release.Status != 1 {
		result.Failed(c, 400, "该工单当前状态不允许审核")
		return
	}

	// 检查审核状态
	if release.ApprovalStatus != 1 {
		result.Failed(c, 400, "该工单已审核过")
		return
	}

	// 验证审核动作
	if req.Action != 2 && req.Action != 3 {
		result.Failed(c, 400, "无效的审核动作")
		return
	}

	now := time.Now()
	updates := make(map[string]interface{})

	updates["approval_status"] = req.Action
	updates["approval_time"] = &now
	updates["approval_remark"] = req.Remark

	// 如果是驳回，立即结束流程
	if req.Action == 3 {
		updates["status"] = 4         // 已驳回
		updates["execute_status"] = 4 // 执行失败（异常）
		updates["end_time"] = &now
	}

	// 更新工单
	if err := s.dao.UpdateShRelease(req.ReleaseID, updates); err != nil {
		result.Failed(c, 500, "更新审核状态失败: "+err.Error())
		return
	}

	result.Success(c, "审核成功")
}

// PullCode 拉取代码（通过Jenkins）
func (s *ShReleaseService) PullCode(c *gin.Context, req *model.ShExecuteRequest) {
	// 获取工单
	release, err := s.dao.GetShReleaseByID(req.ReleaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "脚本上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取脚本上线工单失败: "+err.Error())
		return
	}

	// 检查工单状态
	if release.Status != 1 {
		result.Failed(c, 400, "该工单当前状态不允许执行")
		return
	}

	// 检查审核是否通过
	if release.ApprovalStatus != 2 {
		result.Failed(c, 400, "审核未通过，不能拉取代码")
		return
	}

	// 检查执行状态
	if release.ExecuteStatus != 1 {
		result.Failed(c, 400, "该工单已开始执行")
		return
	}

	// 更新执行状态为拉取代码中
	now := time.Now()
	updates := map[string]interface{}{
		"execute_status":      2, // 拉取代码中
		"pull_code_start_time": &now,
	}

	if err := s.dao.UpdateShRelease(req.ReleaseID, updates); err != nil{
		result.Failed(c, 500, "更新状态失败: "+err.Error())
		return
	}

	// 启动异步拉取代码流程
	go s.pullCodeFromJenkins(req.ReleaseID)

	result.Success(c, "代码拉取流程已启动")
}

// pullCodeFromJenkins 通过Jenkins拉取代码
func (s *ShReleaseService) pullCodeFromJenkins(releaseID uint) {
	// 获取工单详情
	release, err := s.dao.GetShReleaseByID(releaseID)
	if err != nil {
		return
	}

	startTime := time.Now()

	// 检查Jenkins配置
	if release.JenkinsEnvID == nil {
		s.updateExecuteFailure(releaseID, startTime, "未配置Jenkins环境")
		return
	}

	// 获取Jenkins环境配置
	var jenkinsEnv model.JenkinsEnv
	if err := s.db.First(&jenkinsEnv, release.JenkinsEnvID).Error; err != nil {
		s.updateExecuteFailure(releaseID, startTime, "获取Jenkins配置失败: "+err.Error())
		return
	}

	if jenkinsEnv.JenkinsServerID == nil || jenkinsEnv.JobName == "" {
		s.updateExecuteFailure(releaseID, startTime, "Jenkins配置不完整")
		return
	}

	// 构建Jenkins参数：合并固定参数和自定义参数
	parameters := map[string]string{
		"EXECUTE_DIR":    release.ExecuteDir,
		"SCRIPT_CONTENT": release.ScriptContent,
	}

	// 如果有存储的自定义参数，合并进去
	if release.Parameters != "" {
		var customParams map[string]string
		if err := json.Unmarshal([]byte(release.Parameters), &customParams); err == nil {
			// 合并自定义参数（自定义参数优先级更高，可以覆盖默认参数）
			for key, value := range customParams {
				parameters[key] = value
			}
		}
	}

	parametersJSON, err := json.Marshal(parameters)
	if err != nil {
		s.updateExecuteFailure(releaseID, startTime, "构建参数失败: "+err.Error())
		return
	}

	// 启动Jenkins任务
	started, buildNumber, logURL, err := executeJenkinsJob(
		*jenkinsEnv.JenkinsServerID,
		jenkinsEnv.JobName,
		string(parametersJSON),
	)

	if err != nil || !started {
		errorMsg := "启动Jenkins任务失败"
		if err != nil {
			errorMsg += ": " + err.Error()
		}
		s.updateExecuteFailure(releaseID, startTime, errorMsg)
		return
	}

	// 更新构建信息
	s.dao.UpdateShRelease(releaseID, map[string]interface{}{
		"build_number": buildNumber,
		"log_url":      logURL,
	})

	// 启动异步监控
	go s.monitorExecution(releaseID, *jenkinsEnv.JenkinsServerID, jenkinsEnv.JobName, buildNumber, startTime)
}

// monitorExecution 监控脚本执行状态
func (s *ShReleaseService) monitorExecution(releaseID uint, serverID uint, jobName string, buildNumber int, startTime time.Time) {
	// 获取Jenkins配置
	accountDao := ccdao.NewAccountAuthDao()
	account, err := accountDao.GetByID(serverID)
	if err != nil {
		s.updateExecuteFailure(releaseID, startTime, "获取Jenkins配置失败: "+err.Error())
		return
	}

	password, err := account.DecryptPassword()
	if err != nil {
		s.updateExecuteFailure(releaseID, startTime, "解密Jenkins密码失败: "+err.Error())
		return
	}

	jenkinsClient := NewJenkinsClient(account.Host, account.Port, account.Name, password)
	baseURL := jenkinsClient.BaseURL

	// 循环监控
	checkInterval := 10 * time.Second
	maxWaitTime := 30 * time.Minute
	monitorStartTime := time.Now()

	for {
		if time.Since(monitorStartTime) > maxWaitTime {
			s.updateExecuteFailure(releaseID, startTime, "执行超时")
			return
		}

		isBuilding, result, _, err := getBuildStatusDetailed(baseURL, account.Name, password, jobName, buildNumber)
		if err != nil {
			time.Sleep(checkInterval)
			continue
		}

		if !isBuilding {
			s.updateExecuteCompletion(releaseID, startTime, result)
			return
		}

		time.Sleep(checkInterval)
	}
}

// updateExecuteFailure 更新拉取代码失败状态
func (s *ShReleaseService) updateExecuteFailure(releaseID uint, startTime time.Time, errorMessage string) {
	endTime := time.Now()

	s.dao.UpdateShRelease(releaseID, map[string]interface{}{
		"execute_status":     4, // 代码拉取失败
		"status":             3, // 流程失败
		"pull_code_end_time": &endTime,
		"error_message":      errorMessage,
	})
}

// updateExecuteCompletion 更新拉取代码完成状态
func (s *ShReleaseService) updateExecuteCompletion(releaseID uint, startTime time.Time, jenkinsResult string) {
	endTime := time.Now()

	var executeStatus int
	var status int
	var errorMessage string

	switch jenkinsResult {
	case "SUCCESS":
		executeStatus = 3 // 代码拉取成功
		status = 1        // 流程继续（等待执行脚本）
	case "FAILURE":
		executeStatus = 4 // 代码拉取失败
		status = 3        // 流程失败
		errorMessage = "Jenkins构建失败"
	case "ABORTED":
		executeStatus = 4
		status = 3
		errorMessage = "Jenkins构建被中止"
	default:
		executeStatus = 4
		status = 3
		errorMessage = fmt.Sprintf("Jenkins构建结果未知: %s", jenkinsResult)
	}

	s.dao.UpdateShRelease(releaseID, map[string]interface{}{
		"execute_status":     executeStatus,
		"status":             status,
		"pull_code_end_time": &endTime,
		"error_message":      errorMessage,
	})
}

// GetShReleaseLog 获取Jenkins日志
func (s *ShReleaseService) GetShReleaseLog(c *gin.Context, releaseID uint, start int) {
	// 获取工单详情
	release, err := s.dao.GetShReleaseByID(releaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "脚本上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取脚本上线工单失败: "+err.Error())
		return
	}

	// 检查是否已启动Jenkins任务
	if release.BuildNumber == 0 {
		result.Failed(c, 400, "Jenkins任务尚未启动")
		return
	}

	// 获取Jenkins环境配置
	if release.JenkinsEnvID == nil {
		result.Failed(c, 400, "Jenkins环境配置不存在")
		return
	}

	var jenkinsEnv model.JenkinsEnv
	if err := s.db.First(&jenkinsEnv, release.JenkinsEnvID).Error; err != nil {
		result.Failed(c, 500, "获取Jenkins配置失败: "+err.Error())
		return
	}

	if jenkinsEnv.JenkinsServerID == nil {
		result.Failed(c, 400, "Jenkins服务器配置不完整")
		return
	}

	// 获取Jenkins服务器信息
	accountDao := ccdao.NewAccountAuthDao()
	account, err := accountDao.GetByID(*jenkinsEnv.JenkinsServerID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins服务器信息失败: "+err.Error())
		return
	}

	password, err := account.DecryptPassword()
	if err != nil {
		result.Failed(c, 500, "解密Jenkins密码失败: "+err.Error())
		return
	}

	// 构建Jenkins客户端
	jenkinsClient := NewJenkinsClient(account.Host, account.Port, account.Name, password)
	baseURL := jenkinsClient.BaseURL

	// 获取Jenkins日志
	logURL := fmt.Sprintf("%s/job/%s/%d/logText/progressiveText?start=%d",
		baseURL, jenkinsEnv.JobName, release.BuildNumber, start)

	req, err := http.NewRequest("GET", logURL, nil)
	if err != nil {
		result.Failed(c, 500, "创建请求失败: "+err.Error())
		return
	}

	req.SetBasicAuth(account.Name, password)

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		result.Failed(c, 500, "请求Jenkins日志失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("获取Jenkins日志失败，状态码: %d", resp.StatusCode))
		return
	}

	// 读取日志内容
	logContent, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Failed(c, 500, "读取日志内容失败: "+err.Error())
		return
	}

	// 获取下一次请求的起始位置
	nextStart := resp.Header.Get("X-Text-Size")
	moreData := resp.Header.Get("X-More-Data")

	// 返回日志数据
	result.Success(c, gin.H{
		"release_id":    releaseID,
		"build_number":  release.BuildNumber,
		"log_content":   string(logContent),
		"next_start":    nextStart,
		"more_data":     moreData == "true",
		"has_more":      moreData == "true",
		"log_url":       release.LogURL,
		"job_name":      jenkinsEnv.JobName,
		"jenkins_url":   baseURL,
		"start":         start,
		"status":        release.ExecuteStatus,
		"error_message": release.ErrorMessage,
	})
}

// CheckUserPermission 检查用户操作权限
func (s *ShReleaseService) CheckUserPermission(c *gin.Context, releaseID uint, operationType string) {
	// 从JWT获取当前用户ID
	currentUserID, err := jwt.GetAdminId(c)
	if err != nil {
		result.Failed(c, 401, "无法获取当前用户信息")
		return
	}

	// 获取工单信息
	release, err := s.dao.GetShReleaseByID(releaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "脚本上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取脚本上线工单失败: "+err.Error())
		return
	}

	var hasPermission bool
	var permissionMsg string
	var currentStep string
	var allowedUsers []uint
	var allowedUserNames []string

	switch operationType {
	case "approval":
		// 审核权限
		currentStep = "审核"
		if release.ApproverID != nil && *release.ApproverID == currentUserID {
			hasPermission = true
		}
		if release.ApproverID != nil {
			allowedUsers = []uint{*release.ApproverID}
			approverName := s.getApproverName(release.ApproverID)
			allowedUserNames = []string{approverName}
			permissionMsg = fmt.Sprintf("只有审核人 %s 才能进行审核", approverName)
		}

	case "execute":
		// 执行权限
		currentStep = "执行"
		if release.ExecutorID != nil && *release.ExecutorID == currentUserID {
			hasPermission = true
		}
		if release.ExecutorID != nil {
			allowedUsers = []uint{*release.ExecutorID}
			executorName := s.getApproverName(release.ExecutorID)
			allowedUserNames = []string{executorName}
			permissionMsg = fmt.Sprintf("只有执行人 %s 才能执行脚本", executorName)
		}

	default:
		result.Failed(c, 400, "无效的操作类型")
		return
	}

	// 返回权限检查结果
	result.Success(c, gin.H{
		"has_permission":     hasPermission,
		"current_step":       currentStep,
		"operation_type":     operationType,
		"release_id":         releaseID,
		"current_user_id":    currentUserID,
		"allowed_user_ids":   allowedUsers,
		"allowed_user_names": allowedUserNames,
		"message":            permissionMsg,
	})
}

// getApproverName 获取审批人姓名
func (s *ShReleaseService) getApproverName(approverID *uint) string {
	if approverID == nil {
		return ""
	}

	// 从sys_admin表获取用户信息
	var user struct {
		Username string
		Nickname string
	}
	if err := s.db.Table("sys_admin").
		Select("username, nickname").
		Where("id = ?", *approverID).
		First(&user).Error; err != nil {
		// 如果查询失败，返回用户ID作为后备
		return fmt.Sprintf("用户_%d", *approverID)
	}

	// 优先使用昵称，如果没有昵称则使用用户名
	if user.Nickname != "" {
		return user.Nickname
	}
	if user.Username != "" {
		return user.Username
	}
	return fmt.Sprintf("用户_%d", *approverID)
}

// ExecuteScript 执行脚本（通过SSH登录服务器执行）
func (s *ShReleaseService) ExecuteScript(c *gin.Context, req *model.ShExecuteScriptRequest) {
	// 获取工单
	release, err := s.dao.GetShReleaseByID(req.ReleaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "脚本上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取脚本上线工单失败: "+err.Error())
		return
	}

	// 检查工单状态
	if release.Status != 1 {
		result.Failed(c, 400, "该工单当前状态不允许执行脚本")
		return
	}

	// 检查审核是否通过
	if release.ApprovalStatus != 2 {
		result.Failed(c, 400, "审核未通过，不能执行脚本")
		return
	}

	// 检查代码是否已拉取成功
	if release.ExecuteStatus != 3 {
		result.Failed(c, 400, "请先拉取代码，代码拉取成功后才能执行脚本")
		return
	}

	// 更新执行状态为脚本执行中
	now := time.Now()
	updates := map[string]interface{}{
		"execute_status": 5, // 脚本执行中
		"start_time":     &now,
	}

	if err := s.dao.UpdateShRelease(req.ReleaseID, updates); err != nil {
		result.Failed(c, 500, "更新状态失败: "+err.Error())
		return
	}

	// 启动异步执行脚本流程
	go s.executeScriptOnServer(req.ReleaseID)

	result.Success(c, "脚本执行已启动")
}

// executeScriptOnServer 在服务器上执行脚本
func (s *ShReleaseService) executeScriptOnServer(releaseID uint) {
	// 获取工单详情
	release, err := s.dao.GetShReleaseByID(releaseID)
	if err != nil {
		return
	}

	startTime := time.Now()

	// 获取主机SSH信息
	var host cmdbModel.CmdbHost
	if err := s.db.First(&host, release.ServerHostID).Error; err != nil {
		s.updateScriptExecuteFailure(releaseID, startTime, "获取主机信息失败: "+err.Error())
		return
	}

	if host.SSHIP == "" || host.SSHPort == 0 || host.SSHName == "" {
		s.updateScriptExecuteFailure(releaseID, startTime, "主机SSH信息不完整")
		return
	}

	// 获取SSH认证信息
	var authMethod ssh.AuthMethod
	if host.SSHKeyID != 0 {
		credential, err := cmdbDao.NewCmdbHostSSHDao().GetSSHCredentials(host.SSHKeyID)
		if err != nil {
			s.updateScriptExecuteFailure(releaseID, startTime, "获取SSH凭据失败: "+err.Error())
			return
		}

		var auth ccModel.EcsAuth
		if err := common.GetDB().Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&auth).Error; err != nil {
			s.updateScriptExecuteFailure(releaseID, startTime, "获取认证类型失败: "+err.Error())
			return
		}

		switch auth.Type {
		case 1: // 密码认证
			authMethod = ssh.Password(credential)
		case 2: // 密钥认证
			authMethod, err = util.NewSSHUtil().PublicKeyAuth(credential)
			if err != nil {
				s.updateScriptExecuteFailure(releaseID, startTime, "创建密钥认证失败: "+err.Error())
				return
			}
		case 3: // 公钥免认证
			if userKeyAuth, err := util.NewSSHUtil().UserKeyAuth(); err == nil {
				authMethod = userKeyAuth
			} else if defaultKeyAuth, err := util.NewSSHUtil().DefaultKeyAuth(); err == nil {
				authMethod = defaultKeyAuth
			} else {
				s.updateScriptExecuteFailure(releaseID, startTime, "公钥免认证失败: 本地未找到私钥文件")
				return
			}
		default:
			s.updateScriptExecuteFailure(releaseID, startTime, "不支持的认证方式")
			return
		}
	} else {
		s.updateScriptExecuteFailure(releaseID, startTime, "主机未配置SSH凭据")
		return
	}

	// 创建SSH连接
	config := &ssh.ClientConfig{
		User:            host.SSHName,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.SSHIP, host.SSHPort), config)
	if err != nil {
		s.updateScriptExecuteFailure(releaseID, startTime, "创建SSH连接失败: "+err.Error())
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		s.updateScriptExecuteFailure(releaseID, startTime, "创建SSH会话失败: "+err.Error())
		return
	}
	defer session.Close()

	// 构建执行命令：先cd到目录，然后执行脚本
	command := fmt.Sprintf("cd %s && %s", release.ExecuteDir, release.ScriptContent)

	// 执行命令并获取输出
	output, err := session.CombinedOutput(command)
	outputStr := string(output)

	// 更新脚本执行结果
	s.updateScriptExecuteCompletion(releaseID, startTime, outputStr, err)
}

// updateScriptExecuteFailure 更新脚本执行失败状态
func (s *ShReleaseService) updateScriptExecuteFailure(releaseID uint, startTime time.Time, errorMessage string) {
	endTime := time.Now()
	duration := int(endTime.Sub(startTime).Seconds())

	s.dao.UpdateShRelease(releaseID, map[string]interface{}{
		"execute_status": 7, // 脚本执行失败
		"status":         3, // 流程失败
		"end_time":       &endTime,
		"duration":       duration,
		"error_message":  errorMessage,
	})
}

// updateScriptExecuteCompletion 更新脚本执行完成状态
func (s *ShReleaseService) updateScriptExecuteCompletion(releaseID uint, startTime time.Time, output string, execErr error) {
	endTime := time.Now()
	duration := int(endTime.Sub(startTime).Seconds())

	var executeStatus int
	var status int
	var errorMessage string

	if execErr != nil {
		executeStatus = 7 // 脚本执行失败
		status = 3        // 流程失败
		errorMessage = fmt.Sprintf("脚本执行失败: %v", execErr)
	} else {
		executeStatus = 6 // 脚本执行成功
		status = 2        // 流程成功
	}

	s.dao.UpdateShRelease(releaseID, map[string]interface{}{
		"execute_status": executeStatus,
		"status":         status,
		"end_time":       &endTime,
		"duration":       duration,
		"script_output":  output,
		"error_message":  errorMessage,
	})
}

// GetScriptOutput 获取脚本执行输出
func (s *ShReleaseService) GetScriptOutput(c *gin.Context, releaseID uint) {
	// 获取工单详情
	release, err := s.dao.GetShReleaseByID(releaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "脚本上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取脚本上线工单失败: "+err.Error())
		return
	}

	// 返回脚本执行信息
	result.Success(c, gin.H{
		"release_id":      releaseID,
		"execute_status":  release.ExecuteStatus,
		"status_text":     getExecuteStatusText(release.ExecuteStatus),
		"script_output":   release.ScriptOutput,
		"error_message":   release.ErrorMessage,
		"start_time":      release.StartTime,
		"end_time":        release.EndTime,
		"duration":        release.Duration,
		"execute_dir":     release.ExecuteDir,
		"script_content":  release.ScriptContent,
		"server_host_id":  release.ServerHostID,
	})
}

// getExecuteStatusText 获取执行状态文本
func getExecuteStatusText(status int) string {
	statusMap := map[int]string{
		1: "未执行",
		2: "拉取代码中",
		3: "代码拉取成功",
		4: "代码拉取失败",
		5: "脚本执行中",
		6: "脚本执行成功",
		7: "脚本执行失败",
	}
	if text, ok := statusMap[status]; ok {
		return text
	}
	return "未知状态"
}

// ListReleaseFiles 列出脚本上线服务器上的文件
func (s *ShReleaseService) ListReleaseFiles(c *gin.Context, releaseID uint, remotePath string) {
	// 获取工单详情
	release, err := s.dao.GetShReleaseByID(releaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "脚本上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取脚本上线工单失败: "+err.Error())
		return
	}

	// 如果未指定路径，使用默认目录
	if remotePath == "" {
		remotePath = "/home/dingding/saas3-data-export/xlsx"
	}

	// 获取主机SSH信息
	var host cmdbModel.CmdbHost
	if err := s.db.First(&host, release.ServerHostID).Error; err != nil {
		result.Failed(c, 500, "获取主机信息失败: "+err.Error())
		return
	}

	if host.SSHIP == "" || host.SSHPort == 0 || host.SSHName == "" {
		result.Failed(c, 500, "主机SSH信息不完整")
		return
	}

	// 获取SSH认证信息
	var auth ccModel.EcsAuth
	if host.SSHKeyID != 0 {
		if err := s.db.Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&auth).Error; err != nil {
			result.Failed(c, 500, "获取认证信息失败: "+err.Error())
			return
		}
	} else {
		result.Failed(c, 500, "主机未配置SSH凭据")
		return
	}

	// 构建SSH配置
	sshConfig := util.SSHConfig{
		IP:        host.SSHIP,
		Port:      host.SSHPort,
		Type:      auth.Type,
		Username:  host.SSHName,
		Password:  auth.Password,
		PublicKey: auth.PublicKey,
	}

	// 使用SSHUtil列出文件
	sshUtil := util.NewSSHUtil()
	fileInfos, err := sshUtil.ListFiles(&sshConfig, remotePath)
	if err != nil {
		result.Failed(c, 500, "列出文件失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{
		"release_id":  releaseID,
		"host_id":     release.ServerHostID,
		"host_ip":     host.SSHIP,
		"remotePath":  remotePath,
		"files":       fileInfos,
		"total":       len(fileInfos),
	})
}

// DownloadReleaseFile 下载脚本上线服务器上的文件
func (s *ShReleaseService) DownloadReleaseFile(c *gin.Context, releaseID uint, remotePath string) {
	// 获取工单详情
	release, err := s.dao.GetShReleaseByID(releaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "脚本上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取脚本上线工单失败: "+err.Error())
		return
	}

	if remotePath == "" {
		result.Failed(c, 400, "远程文件路径不能为空")
		return
	}

	// 获取主机SSH信息
	var host cmdbModel.CmdbHost
	if err := s.db.First(&host, release.ServerHostID).Error; err != nil {
		result.Failed(c, 500, "获取主机信息失败: "+err.Error())
		return
	}

	if host.SSHIP == "" || host.SSHPort == 0 || host.SSHName == "" {
		result.Failed(c, 500, "主机SSH信息不完整")
		return
	}

	// 获取SSH认证信息
	var auth ccModel.EcsAuth
	if host.SSHKeyID != 0 {
		if err := s.db.Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&auth).Error; err != nil {
			result.Failed(c, 500, "获取认证信息失败: "+err.Error())
			return
		}
	} else {
		result.Failed(c, 500, "主机未配置SSH凭据")
		return
	}

	// 构建SSH配置
	sshConfig := util.SSHConfig{
		IP:        host.SSHIP,
		Port:      host.SSHPort,
		Type:      auth.Type,
		Username:  host.SSHName,
		Password:  auth.Password,
		PublicKey: auth.PublicKey,
	}

	// 使用SSHUtil下载文件
	sshUtil := util.NewSSHUtil()

	// 创建临时目录和文件路径
	tempDir := fmt.Sprintf("/tmp/sh_release_downloads/%d", releaseID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		result.Failed(c, 500, "创建临时目录失败: "+err.Error())
		return
	}

	// 生成本地临时文件路径
	fileName := filepath.Base(remotePath)
	localPath := filepath.Join(tempDir, fmt.Sprintf("%s_%d", fileName, time.Now().Unix()))

	// 下载文件
	if err := sshUtil.DownloadFile(&sshConfig, remotePath, localPath); err != nil {
		result.Failed(c, 500, "文件下载失败: "+err.Error())
		return
	}

	// 设置响应头并返回文件
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(remotePath)))
	c.Header("Content-Type", "application/octet-stream")
	c.File(localPath)

	// 异步删除临时文件
	go func() {
		time.Sleep(10 * time.Second)
		os.Remove(localPath)
	}()
}
