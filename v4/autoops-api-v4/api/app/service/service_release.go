package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"dodevops-api/api/app/dao"
	"dodevops-api/api/app/model"
	ccdao "dodevops-api/api/configcenter/dao"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IServiceReleaseService 服务上线服务接口
type IServiceReleaseService interface {
	// 创建服务上线工单
	CreateServiceRelease(c *gin.Context, req *model.CreateServiceReleaseRequest)
	// 获取服务上线工单列表
	GetServiceReleaseList(c *gin.Context, req *model.ServiceReleaseListRequest)
	// 获取服务上线工单详情
	GetServiceReleaseDetail(c *gin.Context, id uint)
	// 删除服务上线工单
	DeleteServiceRelease(c *gin.Context, id uint)

	// 审批操作
	ApproveRelease(c *gin.Context, req *model.ApprovalRequest)

	// 运维发布
	DeployRelease(c *gin.Context, req *model.DeployReleaseRequest)

	// 回归测试
	RegressionTest(c *gin.Context, req *model.RegressionTestRequest)

	// 获取单个服务的Jenkins日志
	GetServiceItemLog(c *gin.Context, itemID uint, start int)

	// 检查用户操作权限
	CheckUserPermission(c *gin.Context, releaseID uint, operationType string)
}

// ServiceReleaseService 服务上线服务实现
type ServiceReleaseService struct {
	dao    dao.IServiceReleaseDao
	appDao dao.IApplicationDao
	db     *gorm.DB
}

// NewServiceReleaseService 创建服务上线服务
func NewServiceReleaseService(db *gorm.DB) IServiceReleaseService {
	return &ServiceReleaseService{
		dao:    dao.NewServiceReleaseDao(db),
		appDao: dao.NewApplicationDao(db),
		db:     db,
	}
}

// CreateServiceRelease 创建服务上线工单
func (s *ServiceReleaseService) CreateServiceRelease(c *gin.Context, req *model.CreateServiceReleaseRequest) {
	// TODO: 从JWT获取当前用户信息
	applicantID := uint(1)
	applicantName := "管理员"

	// 获取审批人姓名
	ownerApproverName := s.getApproverName(req.OwnerApproverID)
	securityApproverName := s.getApproverName(req.SecurityApproverID)
	testApproverName := s.getApproverName(req.TestApproverID)

	// 构建服务上线工单对象
	release := &model.ServiceRelease{
		Title:                  req.Title,
		BusinessGroupID:        req.BusinessGroupID,
		ImpactFeature:          req.ImpactFeature,
		ApplicantID:            applicantID,
		ApplicantName:          applicantName,
		OwnerApproverID:        req.OwnerApproverID,
		OwnerApproverName:      ownerApproverName,
		SecurityApproverID:     req.SecurityApproverID,
		SecurityApproverName:   securityApproverName,
		TestApproverID:         req.TestApproverID,
		TestApproverName:       testApproverName,
		OwnerApprovalStatus:    1, // 未审批
		SecurityApprovalStatus: 1,
		TestApprovalStatus:     1,
		DeployStatus:           1, // 未发布
		RegressionTestStatus:   1, // 未测试
		Status:                 1, // 流程进行中
		ServiceCount:           len(req.Services),
	}

	// 使用事务创建工单和服务明细
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建工单
	if err := tx.Create(release).Error; err != nil {
		tx.Rollback()
		result.Failed(c, 500, "创建服务上线工单失败: "+err.Error())
		return
	}

	// 创建服务明细
	for index, serviceReq := range req.Services {
		// 获取应用信息
		app, err := s.appDao.GetApplicationByID(serviceReq.AppID)
		if err != nil {
			tx.Rollback()
			result.Failed(c, 400, fmt.Sprintf("应用ID %d 不存在", serviceReq.AppID))
			return
		}

		// 获取Jenkins环境配置
		envs, err := s.appDao.GetJenkinsEnvsByAppID(serviceReq.AppID)
		if err != nil {
			tx.Rollback()
			result.Failed(c, 500, "获取Jenkins环境配置失败: "+err.Error())
			return
		}

		// 查找目标环境配置
		var targetEnv *model.JenkinsEnv
		for _, env := range envs {
			if env.EnvName == serviceReq.Environment {
				targetEnv = &env
				break
			}
		}

		var jenkinsEnvID *uint
		if targetEnv != nil {
			jenkinsEnvID = &targetEnv.ID
		}

		// 创建服务明细
		// 将参数map转换为JSON字符串
		var parametersJSON string
		if len(serviceReq.Parameters) > 0 {
			if jsonBytes, err := json.Marshal(serviceReq.Parameters); err == nil {
				parametersJSON = string(jsonBytes)
			}
		}

		item := model.ServiceReleaseItem{
			ReleaseID:      release.ID,
			AppID:          serviceReq.AppID,
			AppName:        app.Name,
			AppCode:        app.Code,
			ProjectName:    serviceReq.ProjectName,
			RepoURL:        serviceReq.RepoURL,
			Branch:         serviceReq.Branch,
			CommitID:       serviceReq.CommitID,
			ImpactFeature:  serviceReq.ImpactFeature,
			FunctionModule: serviceReq.FunctionModule,
			DBChange:       serviceReq.DBChange,
			ConfigChange:   serviceReq.ConfigChange,
			Remark:         serviceReq.Remark,
			JenkinsEnvID:   jenkinsEnvID,
			Parameters:     parametersJSON, // 保存参数JSON字符串
			Status:         1,               // 未发布
			ExecuteOrder:   index + 1,
		}

		if serviceReq.Branch == "" {
			item.Branch = "master"
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			result.Failed(c, 500, "创建服务明细失败: "+err.Error())
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		result.Failed(c, 500, "提交事务失败: "+err.Error())
		return
	}

	// 返回创建的工单（包含服务明细）
	createdRelease, _ := s.dao.GetServiceReleaseByID(release.ID)
	result.Success(c, createdRelease)
}

// GetServiceReleaseList 获取服务上线工单列表
func (s *ServiceReleaseService) GetServiceReleaseList(c *gin.Context, req *model.ServiceReleaseListRequest) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	releases, total, err := s.dao.GetServiceReleaseList(req)
	if err != nil {
		result.Failed(c, 500, "获取服务上线列表失败: "+err.Error())
		return
	}

	response := &model.ServiceReleaseListResponse{
		Total: int(total),
		List:  releases,
	}

	result.Success(c, response)
}

// GetServiceReleaseDetail 获取服务上线工单详情
func (s *ServiceReleaseService) GetServiceReleaseDetail(c *gin.Context, id uint) {
	release, err := s.dao.GetServiceReleaseByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "服务上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取服务上线工单详情失败: "+err.Error())
		return
	}

	result.Success(c, release)
}

// DeleteServiceRelease 删除服务上线工单
func (s *ServiceReleaseService) DeleteServiceRelease(c *gin.Context, id uint) {
	// 检查工单是否存在
	release, err := s.dao.GetServiceReleaseByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "服务上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取服务上线工单失败: "+err.Error())
		return
	}

	// 检查工单状态，发布中的不允许删除
	if release.DeployStatus == 2 {
		result.Failed(c, 400, "发布中的工单不能删除")
		return
	}

	// 删除工单（级联删除服务明细）
	if err := s.dao.DeleteServiceRelease(id); err != nil {
		result.Failed(c, 500, "删除服务上线工单失败: "+err.Error())
		return
	}

	result.Success(c, "删除成功")
}

// ApproveRelease 审批操作
func (s *ServiceReleaseService) ApproveRelease(c *gin.Context, req *model.ApprovalRequest) {
	// 获取工单
	release, err := s.dao.GetServiceReleaseByID(req.ReleaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "服务上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取服务上线工单失败: "+err.Error())
		return
	}

	// 检查工单状态
	if release.Status != 1 {
		result.Failed(c, 400, "该工单当前状态不允许审批")
		return
	}

	// 验证审批动作
	if req.Action != 2 && req.Action != 3 {
		result.Failed(c, 400, "无效的审批动作")
		return
	}

	now := time.Now()
	updates := make(map[string]interface{})

	// 根据审批类型更新对应字段
	switch req.ApprovalType {
	case "owner":
		if release.OwnerApprovalStatus != 1 {
			result.Failed(c, 400, "负责人已审批过")
			return
		}
		updates["owner_approval_status"] = req.Action
		updates["owner_approval_time"] = &now
		updates["owner_approval_remark"] = req.Remark

	case "security":
		if release.SecurityApprovalStatus != 1 {
			result.Failed(c, 400, "安全已审批过")
			return
		}
		// 检查负责人是否已通过
		if release.OwnerApprovalStatus != 2 {
			result.Failed(c, 400, "需要负责人审批通过后才能进行安全审批")
			return
		}
		updates["security_approval_status"] = req.Action
		updates["security_approval_time"] = &now
		updates["security_approval_remark"] = req.Remark

	case "test":
		if release.TestApprovalStatus != 1 {
			result.Failed(c, 400, "测试已审批过")
			return
		}
		// 检查安全是否已通过
		if release.SecurityApprovalStatus != 2 {
			result.Failed(c, 400, "需要安全审批通过后才能进行测试审批")
			return
		}
		updates["test_approval_status"] = req.Action
		updates["test_approval_time"] = &now
		updates["test_approval_remark"] = req.Remark

	default:
		result.Failed(c, 400, "无效的审批类型")
		return
	}

	// 如果是驳回，立即结束流程并标记为异常状态
	if req.Action == 3 {
		updates["status"] = 4              // 已驳回
		updates["deploy_status"] = 4       // 发布失败（异常）
		updates["regression_test_status"] = 4 // 测试失败（异常）
		now := time.Now()
		updates["end_time"] = &now
	}

	// 更新工单
	if err := s.dao.UpdateServiceRelease(req.ReleaseID, updates); err != nil {
		result.Failed(c, 500, "更新审批状态失败: "+err.Error())
		return
	}

	result.Success(c, "审批成功")
}

// DeployRelease 运维发布
func (s *ServiceReleaseService) DeployRelease(c *gin.Context, req *model.DeployReleaseRequest) {
	// 获取工单
	release, err := s.dao.GetServiceReleaseByID(req.ReleaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "服务上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取服务上线工单失败: "+err.Error())
		return
	}

	// 检查工单状态
	if release.Status != 1 {
		result.Failed(c, 400, "该工单当前状态不允许发布")
		return
	}

	// 检查三方审批是否都通过
	if release.OwnerApprovalStatus != 2 {
		result.Failed(c, 400, "负责人审批未通过，不能发布")
		return
	}
	if release.SecurityApprovalStatus != 2 {
		result.Failed(c, 400, "安全审批未通过，不能发布")
		return
	}
	if release.TestApprovalStatus != 2 {
		result.Failed(c, 400, "测试审批未通过，不能发布")
		return
	}

	// 检查发布状态
	if release.DeployStatus != 1 {
		result.Failed(c, 400, "该工单已发布或正在发布中")
		return
	}

	// 更新发布状态
	now := time.Now()
	updates := map[string]interface{}{
		"deploy_status": 2, // 发布中
		"start_time":    &now,
	}

	if err := s.dao.UpdateServiceRelease(req.ReleaseID, updates); err != nil {
		result.Failed(c, 500, "更新发布状态失败: "+err.Error())
		return
	}

	// 启动异步发布流程
	executionMode := req.ExecutionMode
	if executionMode == 0 {
		executionMode = 1 // 默认并行
	}

	go s.executeReleaseDeploy(req.ReleaseID, executionMode)

	result.Success(c, "发布流程已启动")
}

// executeReleaseDeploy 执行发布流程
func (s *ServiceReleaseService) executeReleaseDeploy(releaseID uint, executionMode int) {
	// 获取服务明细列表
	items, err := s.dao.GetServiceReleaseItems(releaseID)
	if err != nil {
		return
	}

	// 根据执行模式选择发布策略
	if executionMode == 2 {
		// 串行发布
		s.deployItemsSerially(items)
	} else {
		// 并行发布
		s.deployItemsParallel(items)
	}
}

// deployItemsParallel 并行发布
func (s *ServiceReleaseService) deployItemsParallel(items []model.ServiceReleaseItem) {
	for _, item := range items {
		go s.deployItem(&item)
	}
}

// deployItemsSerially 串行发布
func (s *ServiceReleaseService) deployItemsSerially(items []model.ServiceReleaseItem) {
	for _, item := range items {
		if !s.deployItem(&item) {
			// 如果发布失败，取消后续任务
			break
		}
		// 等待任务完成
		s.waitForItemCompletion(item.ID)
	}
}

// deployItem 发布单个服务
func (s *ServiceReleaseService) deployItem(item *model.ServiceReleaseItem) bool {
	startTime := time.Now()

	// 更新状态为发布中
	s.dao.UpdateServiceReleaseItem(item.ID, map[string]interface{}{
		"status":     2, // 发布中
		"start_time": &startTime,
	})

	// 检查Jenkins配置
	if item.JenkinsEnvID == nil {
		s.updateItemFailure(item.ID, startTime, "未配置Jenkins环境")
		return false
	}

	// 获取Jenkins环境配置
	var jenkinsEnv model.JenkinsEnv
	if err := s.db.First(&jenkinsEnv, item.JenkinsEnvID).Error; err != nil {
		s.updateItemFailure(item.ID, startTime, "获取Jenkins配置失败: "+err.Error())
		return false
	}

	if jenkinsEnv.JenkinsServerID == nil || jenkinsEnv.JobName == "" {
		s.updateItemFailure(item.ID, startTime, "Jenkins配置不完整")
		return false
	}

	// 使用存储的Jenkins参数
	// 如果item.Parameters为空，则使用默认参数
	var parametersJSON string
	if item.Parameters != "" {
		// 使用存储的参数
		parametersJSON = item.Parameters
	} else {
		// 如果没有存储参数，使用默认参数（Branch和CommitID）
		defaultParams := map[string]string{
			"BRANCH":    item.Branch,
			"COMMIT_ID": item.CommitID,
		}
		if jsonBytes, err := json.Marshal(defaultParams); err == nil {
			parametersJSON = string(jsonBytes)
		}
	}

	// 启动Jenkins任务
	started, buildNumber, logURL, err := executeJenkinsJob(
		*jenkinsEnv.JenkinsServerID,
		jenkinsEnv.JobName,
		parametersJSON,
	)

	if err != nil || !started {
		errorMsg := "启动Jenkins任务失败"
		if err != nil {
			errorMsg += ": " + err.Error()
		}
		s.updateItemFailure(item.ID, startTime, errorMsg)
		return false
	}

	// 更新构建信息
	s.dao.UpdateServiceReleaseItem(item.ID, map[string]interface{}{
		"build_number": buildNumber,
		"log_url":      logURL,
	})

	// 启动异步监控
	go s.monitorItemDeploy(item.ID, *jenkinsEnv.JenkinsServerID, jenkinsEnv.JobName, buildNumber, startTime)

	return true
}

// monitorItemDeploy 监控服务发布状态
func (s *ServiceReleaseService) monitorItemDeploy(itemID uint, serverID uint, jobName string, buildNumber int, startTime time.Time) {
	// 获取Jenkins配置
	accountDao := ccdao.NewAccountAuthDao()
	account, err := accountDao.GetByID(serverID)
	if err != nil {
		s.updateItemFailure(itemID, startTime, "获取Jenkins配置失败: "+err.Error())
		return
	}

	password, err := account.DecryptPassword()
	if err != nil {
		s.updateItemFailure(itemID, startTime, "解密Jenkins密码失败: "+err.Error())
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
			s.updateItemFailure(itemID, startTime, "发布超时")
			return
		}

		isBuilding, result, _, err := getBuildStatusDetailed(baseURL, account.Name, password, jobName, buildNumber)
		if err != nil {
			time.Sleep(checkInterval)
			continue
		}

		if !isBuilding {
			s.updateItemCompletion(itemID, startTime, result)
			return
		}

		time.Sleep(checkInterval)
	}
}

// updateItemFailure 更新服务明细失败状态
func (s *ServiceReleaseService) updateItemFailure(itemID uint, startTime time.Time, errorMessage string) {
	endTime := time.Now()
	duration := int(endTime.Sub(startTime).Seconds())

	s.dao.UpdateServiceReleaseItem(itemID, map[string]interface{}{
		"status":        4, // 失败
		"end_time":      &endTime,
		"duration":      duration,
		"error_message": errorMessage,
	})

	s.checkAndUpdateReleaseStatus(itemID)
}

// updateItemCompletion 更新服务明细完成状态
func (s *ServiceReleaseService) updateItemCompletion(itemID uint, startTime time.Time, jenkinsResult string) {
	endTime := time.Now()
	duration := int(endTime.Sub(startTime).Seconds())

	var status int
	var errorMessage string

	switch jenkinsResult {
	case "SUCCESS":
		status = 3 // 成功
	case "FAILURE":
		status = 4 // 失败
		errorMessage = "Jenkins构建失败"
	case "ABORTED":
		status = 4
		errorMessage = "Jenkins构建被中止"
	default:
		status = 4
		errorMessage = fmt.Sprintf("Jenkins构建结果未知: %s", jenkinsResult)
	}

	s.dao.UpdateServiceReleaseItem(itemID, map[string]interface{}{
		"status":        status,
		"end_time":      &endTime,
		"duration":      duration,
		"error_message": errorMessage,
	})

	s.checkAndUpdateReleaseStatus(itemID)
}

// waitForItemCompletion 等待服务明细完成
func (s *ServiceReleaseService) waitForItemCompletion(itemID uint) {
	for {
		var item model.ServiceReleaseItem
		if err := s.db.First(&item, itemID).Error; err != nil {
			return
		}

		if item.Status == 3 || item.Status == 4 {
			return
		}

		time.Sleep(1 * time.Second)
	}
}

// checkAndUpdateReleaseStatus 检查并更新工单发布状态
func (s *ServiceReleaseService) checkAndUpdateReleaseStatus(itemID uint) {
	// 获取服务明细
	var item model.ServiceReleaseItem
	if err := s.db.First(&item, itemID).Error; err != nil {
		return
	}

	// 获取所有服务明细
	items, err := s.dao.GetServiceReleaseItems(item.ReleaseID)
	if err != nil {
		return
	}

	// 检查是否所有服务都已完成
	allCompleted := true
	allSuccess := true
	for _, it := range items {
		if it.Status != 3 && it.Status != 4 {
			allCompleted = false
			break
		}
		if it.Status != 3 {
			allSuccess = false
		}
	}

	if allCompleted {
		endTime := time.Now()
		totalDuration := 0
		for _, it := range items {
			totalDuration += it.Duration
		}

		var deployStatus int
		if allSuccess {
			deployStatus = 3 // 发布成功
		} else {
			deployStatus = 4 // 发布失败
		}

		updates := map[string]interface{}{
			"deploy_status": deployStatus,
			"end_time":      &endTime,
			"duration":      totalDuration,
		}

		// 如果发布成功，进入回归测试阶段
		if allSuccess {
			updates["regression_test_status"] = 2 // 测试中
		} else {
			// 发布失败，整体流程失败
			updates["status"] = 3 // 发布失败
		}

		s.dao.UpdateServiceRelease(item.ReleaseID, updates)
	}
}

// RegressionTest 回归测试
func (s *ServiceReleaseService) RegressionTest(c *gin.Context, req *model.RegressionTestRequest) {
	// 获取工单
	release, err := s.dao.GetServiceReleaseByID(req.ReleaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "服务上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取服务上线工单失败: "+err.Error())
		return
	}

	// 检查发布状态
	if release.DeployStatus != 3 {
		result.Failed(c, 400, "发布未成功，不能进行回归测试")
		return
	}

	// 检查回归测试状态
	if release.RegressionTestStatus != 2 {
		result.Failed(c, 400, "当前状态不允许更新回归测试结果")
		return
	}

	updates := map[string]interface{}{
		"regression_test_status": req.Status,
	}

	// 根据测试结果更新整体流程状态
	if req.Status == 3 {
		updates["status"] = 2 // 发布成功
	} else {
		updates["status"] = 3 // 发布失败
	}

	if err := s.dao.UpdateServiceRelease(req.ReleaseID, updates); err != nil {
		result.Failed(c, 500, "更新回归测试状态失败: "+err.Error())
		return
	}

	result.Success(c, "回归测试结果已更新")
}

// getApproverName 获取审批人姓名
func (s *ServiceReleaseService) getApproverName(approverID *uint) string {
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

// executeJenkinsJob 启动Jenkins任务（共享函数）
func executeJenkinsJob(serverID uint, jobName string, parametersJSON string) (bool, int, string, error) {
	// 获取Jenkins客户端
	accountDao := ccdao.NewAccountAuthDao()
	account, err := accountDao.GetByID(serverID)
	if err != nil {
		return false, 0, "", fmt.Errorf("获取Jenkins服务器配置失败: %v", err)
	}

	if account.Type != 4 { // Jenkins account type
		return false, 0, "", fmt.Errorf("账号类型不是Jenkins")
	}

	password, err := account.DecryptPassword()
	if err != nil {
		return false, 0, "", fmt.Errorf("解密Jenkins密码失败: %v", err)
	}

	// 创建Jenkins客户端
	jenkinsClient := NewJenkinsClient(account.Host, account.Port, account.Name, password)
	baseURL := jenkinsClient.BaseURL

	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	// 1. 获取CSRF crumb令牌
	crumbURL := fmt.Sprintf("%s/crumbIssuer/api/json", baseURL)
	crumbReq, err := http.NewRequest("GET", crumbURL, nil)
	if err != nil {
		return false, 0, "", fmt.Errorf("创建crumb请求失败: %v", err)
	}

	crumbReq.SetBasicAuth(account.Name, password)
	crumbResp, err := httpClient.Do(crumbReq)
	if err != nil {
		return false, 0, "", fmt.Errorf("获取crumb失败: %v", err)
	}
	defer crumbResp.Body.Close()

	var crumbData struct {
		Crumb            string `json:"crumb"`
		CrumbRequestField string `json:"crumbRequestField"`
	}

	if err := json.NewDecoder(crumbResp.Body).Decode(&crumbData); err != nil {
		return false, 0, "", fmt.Errorf("解析crumb响应失败: %v", err)
	}

	// 2. 解析参数JSON字符串
	var parameters map[string]string
	if parametersJSON != "" {
		if err := json.Unmarshal([]byte(parametersJSON), &parameters); err != nil {
			return false, 0, "", fmt.Errorf("解析参数JSON失败: %v", err)
		}
	}

	// 3. 构建请求URL和表单数据
	var buildURL string
	var formData url.Values

	if len(parameters) > 0 {
		// 带参数构建
		buildURL = fmt.Sprintf("%s/job/%s/buildWithParameters", baseURL, url.PathEscape(jobName))
		formData = url.Values{}
		for key, value := range parameters {
			formData.Set(key, value)
		}
	} else {
		// 无参数构建
		buildURL = fmt.Sprintf("%s/job/%s/build", baseURL, url.PathEscape(jobName))
	}

	// 4. 启动构建（带crumb令牌）
	var req *http.Request
	if len(parameters) > 0 {
		req, err = http.NewRequest("POST", buildURL+"?"+formData.Encode(), nil)
	} else {
		req, err = http.NewRequest("POST", buildURL, nil)
	}
	if err != nil {
		return false, 0, "", fmt.Errorf("创建构建请求失败: %v", err)
	}

	req.SetBasicAuth(account.Name, password)
	req.Header.Set(crumbData.CrumbRequestField, crumbData.Crumb)

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, 0, "", fmt.Errorf("启动Jenkins任务失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, 0, "", fmt.Errorf("Jenkins响应错误，状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	// 获取队列位置
	location := resp.Header.Get("Location")
	if location == "" {
		return false, 0, "", fmt.Errorf("未获取到Jenkins队列位置")
	}

	// 等待新构建开始并获取构建编号
	buildNumber, err := waitForNewBuildStart(baseURL, account.Name, password, jobName)
	if err != nil {
		return false, 0, "", fmt.Errorf("等待新构建开始失败: %v", err)
	}

	// 获取构建URL
	logURL := fmt.Sprintf("%s/job/%s/%d/console", baseURL, url.PathEscape(jobName), buildNumber)

	return true, buildNumber, logURL, nil
}

// waitForNewBuildStart 等待新构建开始并返回新的构建号
func waitForNewBuildStart(baseURL, username, password, jobName string) (int, error) {
	// 先获取当前最新构建号作为基准
	currentBuildNumber, err := getLatestBuildNumber(baseURL, username, password, jobName)
	if err != nil {
		// 如果获取失败，可能是第一次构建，设为0
		currentBuildNumber = 0
	}

	maxWaitTime := 60 * time.Second
	checkInterval := 2 * time.Second
	startTime := time.Now()

	for {
		if time.Since(startTime) > maxWaitTime {
			return 0, fmt.Errorf("等待新构建开始超时")
		}

		latestBuildNumber, err := getLatestBuildNumber(baseURL, username, password, jobName)
		if err != nil {
			time.Sleep(checkInterval)
			continue
		}

		if latestBuildNumber > currentBuildNumber {
			return latestBuildNumber, nil
		}

		time.Sleep(checkInterval)
	}
}

// getLatestBuildNumber 获取最新构建编号
func getLatestBuildNumber(baseURL, username, password, jobName string) (int, error) {
	jobURL := fmt.Sprintf("%s/job/%s/api/json?tree=lastBuild[number]", baseURL, url.PathEscape(jobName))

	req, err := http.NewRequest("GET", jobURL, nil)
	if err != nil {
		return 0, err
	}

	req.SetBasicAuth(username, password)

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("获取任务信息失败，状态码: %d", resp.StatusCode)
	}

	var jobData struct {
		LastBuild struct {
			Number int `json:"number"`
		} `json:"lastBuild"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jobData); err != nil {
		return 0, err
	}

	return jobData.LastBuild.Number, nil
}

// getBuildStatusDetailed 获取详细的构建状态
func getBuildStatusDetailed(baseURL, username, password, jobName string, buildNumber int) (bool, string, string, error) {
	buildURL := fmt.Sprintf("%s/job/%s/%d/api/json", baseURL, url.PathEscape(jobName), buildNumber)

	req, err := http.NewRequest("GET", buildURL, nil)
	if err != nil {
		return false, "", "", err
	}

	req.SetBasicAuth(username, password)

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", "", fmt.Errorf("获取构建状态失败，状态码: %d", resp.StatusCode)
	}

	var buildData struct {
		Building bool   `json:"building"`
		Result   string `json:"result"`
		URL      string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&buildData); err != nil {
		return false, "", "", err
	}

	logURL := buildData.URL + "console"

	if buildData.Building {
		return true, "", logURL, nil
	}

	return false, buildData.Result, logURL, nil
}

// GetServiceItemLog 获取单个服务的Jenkins日志
func (s *ServiceReleaseService) GetServiceItemLog(c *gin.Context, itemID uint, start int) {
	// 获取服务明细
	var item model.ServiceReleaseItem
	if err := s.db.First(&item, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "服务明细不存在")
			return
		}
		result.Failed(c, 500, "获取服务明细失败: "+err.Error())
		return
	}

	// 检查是否已启动Jenkins任务
	if item.BuildNumber == 0 {
		result.Failed(c, 400, "Jenkins任务尚未启动")
		return
	}

	// 获取Jenkins环境配置
	if item.JenkinsEnvID == nil {
		result.Failed(c, 400, "Jenkins环境配置不存在")
		return
	}

	var jenkinsEnv model.JenkinsEnv
	if err := s.db.First(&jenkinsEnv, item.JenkinsEnvID).Error; err != nil {
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
		baseURL, url.PathEscape(jenkinsEnv.JobName), item.BuildNumber, start)

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
		"item_id":       itemID,
		"build_number":  item.BuildNumber,
		"log_content":   string(logContent),
		"next_start":    nextStart,
		"more_data":     moreData == "true",
		"has_more":      moreData == "true",
		"log_url":       item.LogURL,
		"job_name":      jenkinsEnv.JobName,
		"jenkins_url":   baseURL,
		"start":         start,
		"status":        item.Status,
		"error_message": item.ErrorMessage,
	})
}

// CheckUserPermission 检查用户操作权限
func (s *ServiceReleaseService) CheckUserPermission(c *gin.Context, releaseID uint, operationType string) {
	// 从JWT获取当前用户ID
	currentUserID, err := jwt.GetAdminId(c)
	if err != nil {
		result.Failed(c, 401, "无法获取当前用户信息")
		return
	}

	// 获取工单信息
	release, err := s.dao.GetServiceReleaseByID(releaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "服务上线工单不存在")
			return
		}
		result.Failed(c, 500, "获取服务上线工单失败: "+err.Error())
		return
	}

	var hasPermission bool
	var permissionMsg string
	var currentStep string
	var allowedUsers []uint
	var allowedUserNames []string

	switch operationType {
	case "owner_approval":
		// 负责人审批权限
		currentStep = "负责人审批"
		if release.OwnerApproverID != nil && *release.OwnerApproverID == currentUserID {
			hasPermission = true
		}
		if release.OwnerApproverID != nil {
			allowedUsers = []uint{*release.OwnerApproverID}
			approverName := s.getApproverName(release.OwnerApproverID)
			allowedUserNames = []string{approverName}
			permissionMsg = fmt.Sprintf("只有负责人 %s 才能进行审批", approverName)
		}

	case "security_approval":
		// 安全审批权限
		currentStep = "安全审批"
		if release.SecurityApproverID != nil && *release.SecurityApproverID == currentUserID {
			hasPermission = true
		}
		if release.SecurityApproverID != nil {
			allowedUsers = []uint{*release.SecurityApproverID}
			approverName := s.getApproverName(release.SecurityApproverID)
			allowedUserNames = []string{approverName}
			permissionMsg = fmt.Sprintf("只有安全审批人 %s 才能进行审批", approverName)
		}

	case "test_approval":
		// 测试审批权限
		currentStep = "测试审批"
		if release.TestApproverID != nil && *release.TestApproverID == currentUserID {
			hasPermission = true
		}
		if release.TestApproverID != nil {
			allowedUsers = []uint{*release.TestApproverID}
			approverName := s.getApproverName(release.TestApproverID)
			allowedUserNames = []string{approverName}
			permissionMsg = fmt.Sprintf("只有测试审批人 %s 才能进行审批", approverName)
		}

	case "deploy":
		// 运维发布权限 - 从应用详情中的运维负责人获取
		currentStep = "运维发布"

		// 获取工单关联的所有应用的运维负责人
		items, err := s.dao.GetServiceReleaseItems(releaseID)
		if err != nil {
			result.Failed(c, 500, "获取服务明细失败: "+err.Error())
			return
		}

		// 收集所有应用的运维负责人
		opsOwnerSet := make(map[uint]bool)
		for _, item := range items {
			app, err := s.appDao.GetApplicationByID(item.AppID)
			if err != nil {
				continue
			}
			// 解析运维负责人
			for _, ownerID := range app.OpsOwners {
				opsOwnerSet[ownerID] = true
				allowedUsers = append(allowedUsers, ownerID)
			}
		}

		// 检查当前用户是否在运维负责人列表中
		if opsOwnerSet[currentUserID] {
			hasPermission = true
		}

		// 获取运维负责人名称
		for ownerID := range opsOwnerSet {
			allowedUserNames = append(allowedUserNames, s.getApproverName(&ownerID))
		}

		if len(allowedUserNames) > 0 {
			permissionMsg = fmt.Sprintf("只有运维负责人 [%s] 才能执行发布", joinStrings(allowedUserNames))
		} else {
			permissionMsg = "未配置运维负责人，请先在应用详情中配置"
		}

	case "regression_test":
		// 回归测试权限 - 从应用详情中的测试负责人获取
		currentStep = "回归测试"

		// 获取工单关联的所有应用的测试负责人
		items, err := s.dao.GetServiceReleaseItems(releaseID)
		if err != nil {
			result.Failed(c, 500, "获取服务明细失败: "+err.Error())
			return
		}

		// 收集所有应用的测试负责人
		testOwnerSet := make(map[uint]bool)
		for _, item := range items {
			app, err := s.appDao.GetApplicationByID(item.AppID)
			if err != nil {
				continue
			}
			// 解析测试负责人
			for _, ownerID := range app.TestOwners {
				testOwnerSet[ownerID] = true
				allowedUsers = append(allowedUsers, ownerID)
			}
		}

		// 检查当前用户是否在测试负责人列表中
		if testOwnerSet[currentUserID] {
			hasPermission = true
		}

		// 获取测试负责人名称
		for ownerID := range testOwnerSet {
			allowedUserNames = append(allowedUserNames, s.getApproverName(&ownerID))
		}

		if len(allowedUserNames) > 0 {
			permissionMsg = fmt.Sprintf("只有测试负责人 [%s] 才能执行回归测试", joinStrings(allowedUserNames))
		} else {
			permissionMsg = "未配置测试负责人，请先在应用详情中配置"
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

// joinStrings 将字符串数组连接成逗号分隔的字符串
func joinStrings(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += ", " + strs[i]
	}
	return result
}
