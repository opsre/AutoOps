package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"dodevops-api/api/app/dao"
	"dodevops-api/api/app/model"
	ccdao "dodevops-api/api/configcenter/dao"
	ccmodel "dodevops-api/api/configcenter/model"
	monitormodel "dodevops-api/api/monitor/model"
	monitordao "dodevops-api/api/monitor/dao"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApplicationService 应用服务接口
type IApplicationService interface {
	// 应用管理
	CreateApplication(c *gin.Context, req *model.CreateApplicationRequest)
	GetApplicationList(c *gin.Context, req *model.ApplicationListRequest)
	GetApplicationDetail(c *gin.Context, id uint)
	UpdateApplication(c *gin.Context, id uint, req *model.UpdateApplicationRequest)
	DeleteApplication(c *gin.Context, id uint)

	// Jenkins环境配置管理（与应用绑定）
	CreateJenkinsEnv(c *gin.Context, req *model.CreateJenkinsEnvRequest)
	GetAppJenkinsEnvs(c *gin.Context, appID uint)
	UpdateAppJenkinsEnv(c *gin.Context, appID uint, envID uint, req *model.UpdateJenkinsEnvRequest)
	DeleteAppJenkinsEnv(c *gin.Context, appID uint, envID uint)


	// Jenkins服务器管理
	GetJenkinsServers(c *gin.Context)
	ValidateJenkinsJob(c *gin.Context, req *model.ValidateJenkinsJobRequest)

	// 快速发布管理
	GetApplicationsForDeployment(c *gin.Context, req *model.GetApplicationsForDeploymentRequest)
	CreateQuickDeployment(c *gin.Context, req *model.CreateQuickDeploymentRequest)
	ExecuteQuickDeployment(c *gin.Context, req *model.ExecuteQuickDeploymentRequest)
	GetQuickDeploymentList(c *gin.Context, req *model.QuickDeploymentListRequest)
	GetQuickDeploymentDetail(c *gin.Context, id uint)
	DeleteQuickDeployment(c *gin.Context, id uint)

	// 业务线服务树管理
	GetServiceTree(c *gin.Context, req *model.GetServiceTreeRequest)

	// 应用视图树（四层结构）
	GetApplicationViewTree(c *gin.Context, req *model.ApplicationViewTreeRequest)

	// 获取应用完整详情（用于详情面板）
	GetApplicationDetailForView(c *gin.Context, id uint)

	// 获取应用发布历史
	GetApplicationReleaseHistory(c *gin.Context, req *model.ApplicationReleaseHistoryRequest)

	// 任务日志监控
	GetTaskBuildLog(c *gin.Context, taskID uint, start int)

	// 连级选择器数据源
	GetBusinessGroupOptions(c *gin.Context)

	// 获取单个应用环境配置
	GetAppEnvironment(c *gin.Context, req *model.GetAppEnvironmentRequest)

	// 获取任务状态
	GetTaskStatus(c *gin.Context, taskID uint)
}

// ApplicationService 应用服务实现
type ApplicationService struct {
	appDao    dao.IApplicationDao
	robotDao  monitordao.NotifyRobotDao
	webhookDao monitordao.WebhookDao
	db        *gorm.DB
}

// NewApplicationService 创建应用服务
func NewApplicationService(db *gorm.DB) IApplicationService {
	return &ApplicationService{
		appDao:    dao.NewApplicationDao(db),
		robotDao:  monitordao.NewNotifyRobotDao(),
		webhookDao: monitordao.NewWebhookDao(),
		db:        db,
	}
}

// CreateApplication 创建应用
func (s *ApplicationService) CreateApplication(c *gin.Context, req *model.CreateApplicationRequest) {
	// 生成或验证应用编码
	var appCode string
	if req.Code != "" {
		// 用户提供了编码，验证格式
		if !util.ValidateAppCode(req.Code) {
			result.Failed(c, 400, "应用编码格式无效，只允许小写字母、数字和连字符，且必须以字母或数字开头和结尾")
			return
		}

		// 检查编码是否已存在
		existingApp, _ := s.appDao.GetApplicationByCode(req.Code)
		if existingApp != nil {
			result.Failed(c, 400, "应用编码已存在")
			return
		}
		appCode = req.Code
	} else {
		// 自动生成编码
		appCode = util.GenerateUniqueAppCode(req.Name, func(code string) bool {
			existingApp, _ := s.appDao.GetApplicationByCode(code)
			return existingApp != nil
		})
	}

	// 构建应用对象
	app := &model.Application{
		Name:            req.Name,
		Code:            appCode,
		Description:     req.Description,
		RepoURL:         req.RepoURL,
		BusinessGroupID: req.BusinessGroupID,
		BusinessDeptID:  req.BusinessDeptID,
		Status:          1, // 默认未激活，需要配置Jenkins环境后自动激活

		DevOwners:  model.UserIDs(req.DevOwners),
		TestOwners: model.UserIDs(req.TestOwners),
		OpsOwners:  model.UserIDs(req.OpsOwners),

		ProgrammingLang: req.ProgrammingLang,
		StartCommand:    req.StartCommand,
		StopCommand:     req.StopCommand,
		HealthAPI:       req.HealthAPI,

		Domains:   model.DomainsJSON(req.Domains),
		Hosts:     model.ResourceIDs(req.Hosts),
		Databases: model.ResourceIDs(req.Databases),
		OtherRes:  req.OtherRes,
	}

	// 创建应用
	if err := s.appDao.CreateApplication(app); err != nil {
		result.Failed(c, 500, "创建应用失败: "+err.Error())
		return
	}

	// 创建Jenkins环境配置
	if len(req.JenkinsEnvs) > 0 {
		// 用户提供了自定义Jenkins环境配置
		for _, envReq := range req.JenkinsEnvs {
			envReq.AppID = app.ID // 设置应用ID
			if err := s.createJenkinsEnvFromRequest(&envReq); err != nil {
				// 记录错误但不影响应用创建成功
			}
		}
	} else {
		// 自动创建默认的3套Jenkins环境配置(prod, test, dev)
		if err := s.createDefaultJenkinsEnvs(app.ID); err != nil {
			// 记录错误但不影响应用创建成功
		}
	}

	result.Success(c, app)
}

// GetApplicationList 获取应用列表
func (s *ApplicationService) GetApplicationList(c *gin.Context, req *model.ApplicationListRequest) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	apps, total, err := s.appDao.GetApplicationList(req)
	if err != nil {
		result.Failed(c, 500, "获取应用列表失败: "+err.Error())
		return
	}

	response := &model.ApplicationListResponse{
		Total: int(total),
		List:  apps,
	}

	result.Success(c, response)
}

// GetApplicationDetail 获取应用详情
func (s *ApplicationService) GetApplicationDetail(c *gin.Context, id uint) {
	app, err := s.appDao.GetApplicationByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "应用不存在")
			return
		}
		result.Failed(c, 500, "获取应用详情失败: "+err.Error())
		return
	}

	result.Success(c, app)
}

// UpdateApplication 更新应用
func (s *ApplicationService) UpdateApplication(c *gin.Context, id uint, req *model.UpdateApplicationRequest) {
	// 检查应用是否存在
	_, err := s.appDao.GetApplicationByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "应用不存在")
			return
		}
		result.Failed(c, 500, "获取应用信息失败: "+err.Error())
		return
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.RepoURL != nil {
		updates["repo_url"] = *req.RepoURL
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.BusinessGroupID != nil {
		updates["business_group_id"] = *req.BusinessGroupID
	}
	if req.BusinessDeptID != nil {
		updates["business_dept_id"] = *req.BusinessDeptID
	}
	if req.DevOwners != nil {
		updates["dev_owners"] = model.UserIDs(*req.DevOwners)
	}
	if req.TestOwners != nil {
		updates["test_owners"] = model.UserIDs(*req.TestOwners)
	}
	if req.OpsOwners != nil {
		updates["ops_owners"] = model.UserIDs(*req.OpsOwners)
	}
	if req.ProgrammingLang != nil {
		updates["programming_lang"] = *req.ProgrammingLang
	}
	if req.StartCommand != nil {
		updates["start_command"] = *req.StartCommand
	}
	if req.StopCommand != nil {
		updates["stop_command"] = *req.StopCommand
	}
	if req.HealthAPI != nil {
		updates["health_api"] = *req.HealthAPI
	}
	if req.Domains != nil {
		updates["domains"] = model.DomainsJSON(*req.Domains)
	}
	if req.Hosts != nil {
		updates["hosts"] = model.ResourceIDs(*req.Hosts)
	}
	if req.Databases != nil {
		updates["databases"] = model.ResourceIDs(*req.Databases)
	}
	if req.OtherRes != nil {
		updates["other_res"] = *req.OtherRes
	}

	// 更新应用
	if err := s.appDao.UpdateApplication(id, updates); err != nil {
		result.Failed(c, 500, "更新应用失败: "+err.Error())
		return
	}

	// 处理Jenkins环境配置更新
	if req.JenkinsEnvs != nil {
		if err := s.updateJenkinsEnvs(id, *req.JenkinsEnvs); err != nil {
			// 记录错误但不影响应用更新成功
			// TODO: 添加日志记录
		}
	}

	// 返回更新后的应用信息
	updatedApp, _ := s.appDao.GetApplicationByID(id)
	result.Success(c, updatedApp)
}

// DeleteApplication 删除应用
func (s *ApplicationService) DeleteApplication(c *gin.Context, id uint) {
	// 检查应用是否存在
	app, err := s.appDao.GetApplicationByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "应用不存在")
			return
		}
		result.Failed(c, 500, "获取应用信息失败: "+err.Error())
		return
	}

	// 手动删除关联的Jenkins环境配置（双重保障）
	// 虽然设置了级联删除，但手动删除可以确保数据一致性并记录操作日志
	if err := s.deleteAppJenkinsEnvs(id); err != nil {
		result.Failed(c, 500, "删除Jenkins环境配置失败: "+err.Error())
		return
	}

	// 删除应用
	if err := s.appDao.DeleteApplication(id); err != nil {
		result.Failed(c, 500, "删除应用失败: "+err.Error())
		return
	}

	result.Success(c, fmt.Sprintf("应用 '%s' 及其关联的Jenkins环境配置删除成功", app.Name))
}

// CreateJenkinsEnv 创建Jenkins环境配置
func (s *ApplicationService) CreateJenkinsEnv(c *gin.Context, req *model.CreateJenkinsEnvRequest) {
	// 检查应用是否存在
	_, err := s.appDao.GetApplicationByID(req.AppID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "应用不存在")
			return
		}
		result.Failed(c, 500, "获取应用信息失败: "+err.Error())
		return
	}

	// 检查同一应用下环境名是否已存在
	existingEnv, _ := s.appDao.GetJenkinsEnvByAppAndEnv(req.AppID, req.EnvName)
	if existingEnv != nil {
		result.Failed(c, 400, "该应用下环境名称已存在")
		return
	}

	// 构建Jenkins环境配置对象
	env := &model.JenkinsEnv{
		AppID:           req.AppID,
		EnvName:         req.EnvName,
		JenkinsServerID: req.JenkinsServerID,
		JobName:         req.JobName,
	}

	// 创建Jenkins环境配置
	if err := s.appDao.CreateJenkinsEnv(env); err != nil {
		result.Failed(c, 500, "创建Jenkins环境配置失败: "+err.Error())
		return
	}

	// 检查并更新应用状态
	if err := s.checkAndUpdateAppStatus(req.AppID); err != nil {
		// 记录错误但不影响创建成功的响应
		// TODO: 添加日志记录
	}

	result.Success(c, env)
}

// UpdateAppJenkinsEnv 更新应用的Jenkins环境配置
func (s *ApplicationService) UpdateAppJenkinsEnv(c *gin.Context, appID uint, envID uint, req *model.UpdateJenkinsEnvRequest) {
	// 检查应用是否存在
	_, err := s.appDao.GetApplicationByID(appID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "应用不存在")
			return
		}
		result.Failed(c, 500, "获取应用信息失败: "+err.Error())
		return
	}

	// 检查Jenkins环境配置是否存在且属于该应用
	env, err := s.appDao.GetJenkinsEnvByID(envID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "Jenkins环境配置不存在")
			return
		}
		result.Failed(c, 500, "获取Jenkins环境配置信息失败: "+err.Error())
		return
	}

	if env.AppID != appID {
		result.Failed(c, 403, "Jenkins环境配置不属于该应用")
		return
	}

	// 如果要更新环境名称，检查是否重复
	if req.EnvName != nil && *req.EnvName != env.EnvName {
		existingEnv, _ := s.appDao.GetJenkinsEnvByAppAndEnv(appID, *req.EnvName)
		if existingEnv != nil {
			result.Failed(c, 400, "该应用下环境名称已存在")
			return
		}
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.EnvName != nil {
		updates["env_name"] = *req.EnvName
	}
	if req.JenkinsServerID != nil {
		updates["jenkins_server_id"] = *req.JenkinsServerID
	}
	if req.JobName != nil {
		updates["job_name"] = *req.JobName
	}

	// 更新Jenkins环境配置
	if err := s.appDao.UpdateJenkinsEnv(envID, updates); err != nil {
		result.Failed(c, 500, "更新Jenkins环境配置失败: "+err.Error())
		return
	}

	// 检查并更新应用状态
	if err := s.checkAndUpdateAppStatus(appID); err != nil {
		// 记录错误但不影响更新成功的响应
		// TODO: 添加日志记录
	}

	// 返回更新后的环境配置信息
	updatedEnv, _ := s.appDao.GetJenkinsEnvByID(envID)
	result.Success(c, updatedEnv)
}

// DeleteAppJenkinsEnv 删除应用的Jenkins环境配置
func (s *ApplicationService) DeleteAppJenkinsEnv(c *gin.Context, appID uint, envID uint) {
	// 检查应用是否存在
	_, err := s.appDao.GetApplicationByID(appID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "应用不存在")
			return
		}
		result.Failed(c, 500, "获取应用信息失败: "+err.Error())
		return
	}

	// 检查Jenkins环境配置是否存在且属于该应用
	env, err := s.appDao.GetJenkinsEnvByID(envID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "Jenkins环境配置不存在")
			return
		}
		result.Failed(c, 500, "获取Jenkins环境配置信息失败: "+err.Error())
		return
	}

	if env.AppID != appID {
		result.Failed(c, 403, "Jenkins环境配置不属于该应用")
		return
	}

	// 删除Jenkins环境配置
	if err := s.appDao.DeleteJenkinsEnv(envID); err != nil {
		result.Failed(c, 500, "删除Jenkins环境配置失败: "+err.Error())
		return
	}

	// 检查并更新应用状态
	if err := s.checkAndUpdateAppStatus(appID); err != nil {
		// 记录错误但不影响删除成功的响应
		// TODO: 添加日志记录
	}

	result.Success(c, "删除成功")
}

// GetAppJenkinsEnvs 获取应用的所有Jenkins环境配置
func (s *ApplicationService) GetAppJenkinsEnvs(c *gin.Context, appID uint) {
	// 检查应用是否存在
	_, err := s.appDao.GetApplicationByID(appID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "应用不存在")
			return
		}
		result.Failed(c, 500, "获取应用信息失败: "+err.Error())
		return
	}

	// 获取应用的所有Jenkins环境配置
	envs, err := s.appDao.GetJenkinsEnvsByAppID(appID)
	if err != nil {
		result.Failed(c, 500, "获取应用Jenkins环境配置失败: "+err.Error())
		return
	}

	result.Success(c, envs)
}

// createDefaultJenkinsEnvs 创建默认的3套Jenkins环境配置
func (s *ApplicationService) createDefaultJenkinsEnvs(appID uint) error {
	// 默认环境配置
	defaultEnvs := []string{"prod", "test", "dev"}

	for _, envName := range defaultEnvs {
		jenkinsEnv := &model.JenkinsEnv{
			AppID:           appID,
			EnvName:         envName,
			JenkinsServerID: nil, // 默认为空，需要用户手动配置
			JobName:         "",  // 默认为空，需要用户手动配置
		}

		if err := s.appDao.CreateJenkinsEnv(jenkinsEnv); err != nil {
			return err
		}
	}

	return nil
}

// GetApplicationsForDeployment 获取可发布的应用列表
func (s *ApplicationService) GetApplicationsForDeployment(c *gin.Context, req *model.GetApplicationsForDeploymentRequest) {
	// 验证环境不能是prod
	if req.Environment == "prod" {
		result.Failed(c, 400, "生产环境禁止快速发布")
		return
	}

	// 获取指定业务组和部门下的所有已激活应用
	apps, _, err := s.appDao.GetApplicationList(&model.ApplicationListRequest{
		Page:            1,
		PageSize:        1000, // 获取所有应用
		BusinessGroupID: &req.BusinessGroupID,
		BusinessDeptID:  &req.BusinessDeptID,
		Status:          func() *int { status := 2; return &status }(), // 只获取已激活的应用
	})
	if err != nil {
		result.Failed(c, 500, "获取应用列表失败: "+err.Error())
		return
	}

	var deployableApps []model.ApplicationForDeployment
	for _, app := range apps {
		// 获取应用的Jenkins环境配置
		envs, err := s.appDao.GetJenkinsEnvsByAppID(app.ID)
		if err != nil {
			continue
		}

		// 查找指定环境的配置
		var targetEnv *model.JenkinsEnv
		for _, env := range envs {
			if env.EnvName == req.Environment {
				targetEnv = &env
				break
			}
		}

		deployApp := model.ApplicationForDeployment{
			ID:          app.ID,
			Name:        app.Name,
			Code:        app.Code,
			Environment: req.Environment,
		}

		if targetEnv == nil {
			deployApp.CanDeploy = false
			deployApp.Reason = "未配置该环境"
		} else if targetEnv.JenkinsServerID == nil || *targetEnv.JenkinsServerID == 0 {
			deployApp.CanDeploy = false
			deployApp.Reason = "未配置Jenkins服务器"
		} else if targetEnv.JobName == "" {
			deployApp.CanDeploy = false
			deployApp.Reason = "未配置Jenkins任务名称"
		} else {
			deployApp.CanDeploy = true
			deployApp.JenkinsEnvID = targetEnv.ID
			deployApp.JobName = targetEnv.JobName
			deployApp.Reason = "可以发布"
		}

		deployableApps = append(deployableApps, deployApp)
	}

	result.Success(c, deployableApps)
}

// CreateQuickDeployment 创建快速发布
func (s *ApplicationService) CreateQuickDeployment(c *gin.Context, req *model.CreateQuickDeploymentRequest) {
	// 验证所有应用的环境不能是prod
	for _, app := range req.Applications {
		if app.Environment == "prod" {
			result.Failed(c, 400, "生产环境禁止快速发布")
			return
		}
	}

	// 获取当前用户信息（这里需要从JWT token或session中获取）
	// 暂时使用硬编码，实际应该从认证中间件获取
	creatorID := uint(1)      // TODO: 从认证信息获取
	creatorName := "管理员"    // TODO: 从认证信息获取

	// 创建快速发布记录
	deployment := &model.QuickDeployment{
		Title:           req.Title,
		BusinessGroupID: req.BusinessGroupID,
		BusinessDeptID:  req.BusinessDeptID,
		Description:     req.Description,
		Status:          1, // 待发布
		TaskCount:       len(req.Applications), // 记录任务数量
		CreatorID:       creatorID,
		CreatorName:     creatorName,
		Notice:          req.Notice != nil && *req.Notice, // 是否开启通知
		RobotID:         req.RobotID, // 通知机器人ID
	}

	// 使用事务创建发布记录和任务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建发布记录
	if err := tx.Create(deployment).Error; err != nil {
		tx.Rollback()
		result.Failed(c, 500, "创建发布记录失败: "+err.Error())
		return
	}

	// 创建发布任务，按数组顺序设置执行顺序
	var tasks []model.QuickDeploymentTask
	for index, appReq := range req.Applications {
		// 获取应用信息
		app, err := s.appDao.GetApplicationByID(appReq.AppID)
		if err != nil {
			tx.Rollback()
			result.Failed(c, 400, fmt.Sprintf("应用ID %d 不存在", appReq.AppID))
			return
		}

		// 获取Jenkins环境配置
		envs, err := s.appDao.GetJenkinsEnvsByAppID(appReq.AppID)
		if err != nil {
			tx.Rollback()
			result.Failed(c, 500, "获取Jenkins环境配置失败: "+err.Error())
			return
		}

		// 查找目标环境配置
		var targetEnv *model.JenkinsEnv
		for _, env := range envs {
			if env.EnvName == appReq.Environment {
				targetEnv = &env
				break
			}
		}

		if targetEnv == nil {
			tx.Rollback()
			result.Failed(c, 400, fmt.Sprintf("应用 %s 未配置 %s 环境", app.Name, appReq.Environment))
			return
		}

		if targetEnv.JenkinsServerID == nil || *targetEnv.JenkinsServerID == 0 || targetEnv.JobName == "" {
			tx.Rollback()
			result.Failed(c, 400, fmt.Sprintf("应用 %s 的 %s 环境Jenkins配置不完整", app.Name, appReq.Environment))
			return
		}

		// 生成Jenkins任务URL
		jobURL := "" // TODO: 根据Jenkins服务器配置生成URL

		// 将参数map转换为JSON字符串
		var parametersJSON string
		if len(appReq.Parameters) > 0 {
			if jsonBytes, err := json.Marshal(appReq.Parameters); err == nil {
				parametersJSON = string(jsonBytes)
			}
		}

		task := model.QuickDeploymentTask{
			DeploymentID:  deployment.ID,
			AppID:         appReq.AppID,
			AppName:       app.Name,
			AppCode:       app.Code,
			Environment:   appReq.Environment,
			JenkinsEnvID:  targetEnv.ID,
			JenkinsJobURL: jobURL,
			Parameters:    parametersJSON, // 保存参数JSON字符串
			Status:        1, // 未部署
			ExecuteOrder:  index + 1, // 使用数组索引+1作为执行顺序
		}

		if err := tx.Create(&task).Error; err != nil {
			tx.Rollback()
			result.Failed(c, 500, "创建发布任务失败: "+err.Error())
			return
		}

		tasks = append(tasks, task)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		result.Failed(c, 500, "提交事务失败: "+err.Error())
		return
	}

	// 设置关联数据
	deployment.Tasks = tasks

	result.Success(c, deployment)
}

// ExecuteQuickDeployment 执行快速发布
func (s *ApplicationService) ExecuteQuickDeployment(c *gin.Context, req *model.ExecuteQuickDeploymentRequest) {
	// 获取发布记录
	var deployment model.QuickDeployment
	if err := s.db.Preload("Tasks").First(&deployment, req.DeploymentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "发布记录不存在")
			return
		}
		result.Failed(c, 500, "获取发布记录失败: "+err.Error())
		return
	}

	// 检查发布状态
	if deployment.Status != 1 {
		result.Failed(c, 400, "该发布流程状态不允许执行")
		return
	}

	// 更新发布状态为发布中，并设置执行模式
	now := time.Now()
	deployment.Status = 2 // 发布中
	deployment.StartTime = &now

	// 设置执行模式，如果用户没有指定则使用默认值(1=并行)
	if req.ExecutionMode != nil {
		deployment.ExecutionMode = *req.ExecutionMode
	} else {
		deployment.ExecutionMode = 1 // 默认并行
	}

	if err := s.db.Save(&deployment).Error; err != nil {
		result.Failed(c, 500, "更新发布状态失败: "+err.Error())
		return
	}

	// 启动异步发布流程
	go s.executeDeploymentTasks(deployment.ID)

	result.Success(c, "发布流程已启动")
}

// GetQuickDeploymentList 获取快速发布列表
func (s *ApplicationService) GetQuickDeploymentList(c *gin.Context, req *model.QuickDeploymentListRequest) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var deployments []model.QuickDeployment
	var total int64

	query := s.db.Model(&model.QuickDeployment{})

	// 构建查询条件
	if req.BusinessGroupID != nil {
		query = query.Where("business_group_id = ?", *req.BusinessGroupID)
	}
	if req.BusinessDeptID != nil {
		query = query.Where("business_dept_id = ?", *req.BusinessDeptID)
	}
	if req.Environment != "" {
		query = query.Where("environment = ?", req.Environment)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.CreatorID != nil {
		query = query.Where("creator_id = ?", *req.CreatorID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		result.Failed(c, 500, "获取发布列表总数失败: "+err.Error())
		return
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Tasks").Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&deployments).Error; err != nil {
		result.Failed(c, 500, "获取发布列表失败: "+err.Error())
		return
	}

	response := &model.QuickDeploymentListResponse{
		Total: int(total),
		List:  deployments,
	}

	result.Success(c, response)
}

// GetQuickDeploymentDetail 获取快速发布详情
func (s *ApplicationService) GetQuickDeploymentDetail(c *gin.Context, id uint) {
	var deployment model.QuickDeployment
	if err := s.db.Preload("Tasks.Application").Preload("Tasks.JenkinsEnv").First(&deployment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "发布记录不存在")
			return
		}
		result.Failed(c, 500, "获取发布详情失败: "+err.Error())
		return
	}

	result.Success(c, deployment)
}

// executeDeploymentTasks 异步执行发布任务（串行）
func (s *ApplicationService) executeDeploymentTasks(deploymentID uint) {
	// 获取发布记录和任务
	var deployment model.QuickDeployment
	if err := s.db.Preload("Tasks").First(&deployment, deploymentID).Error; err != nil {
		return
	}

	// 按执行顺序排序任务
	tasks := deployment.Tasks
	for i := 0; i < len(tasks)-1; i++ {
		for j := i + 1; j < len(tasks); j++ {
			if tasks[i].ExecuteOrder > tasks[j].ExecuteOrder {
				tasks[i], tasks[j] = tasks[j], tasks[i]
			}
		}
	}

	// 根据执行模式决定串行或并行执行
	if deployment.ExecutionMode == 2 { // 串行模式
		s.executeTasksSerially(tasks)
	} else { // 并行模式 (默认)
		s.executeTasksParallel(tasks)
	}

	// 注意：不再在这里立即更新父任务状态
	// 父任务状态将在所有子任务完成时通过 checkAndUpdateDeploymentStatus 更新
}

// executeTasksParallel 并行执行任务
func (s *ApplicationService) executeTasksParallel(tasks []model.QuickDeploymentTask) {
	// 并行启动所有任务
	for _, task := range tasks {
		if !s.executeTask(&task) {
			// 如果任务启动失败，直接标记为失败状态
			s.updateTaskFailure(&task, time.Now(), "任务启动失败")
		}
	}
}

// executeTasksSerially 串行执行任务
func (s *ApplicationService) executeTasksSerially(tasks []model.QuickDeploymentTask) {
	// 串行执行任务，每个任务完成后才执行下一个
	for i, task := range tasks {
		if !s.executeTask(&task) {
			// 如果任务启动失败，直接标记为失败状态
			s.updateTaskFailure(&task, time.Now(), "任务启动失败")
			// 串行模式下，如果某个任务启动失败，取消后续任务
			s.cancelRemainingTasks(tasks[i+1:], "前置任务失败，取消执行")
			return
		}

		// 等待当前任务完成
		s.waitForTaskCompletion(task.ID)
	}
}

// waitForTaskCompletion 等待任务完成
func (s *ApplicationService) waitForTaskCompletion(taskID uint) {
	// 轮询检查任务状态，直到完成
	for {
		var task model.QuickDeploymentTask
		if err := s.db.First(&task, taskID).Error; err != nil {
			return
		}

		// 检查任务是否完成 (3=成功, 4=失败)
		if task.Status == 3 || task.Status == 4 {
			return
		}

		// 等待1秒后再次检查
		time.Sleep(1 * time.Second)
	}
}

// cancelRemainingTasks 取消剩余任务
func (s *ApplicationService) cancelRemainingTasks(tasks []model.QuickDeploymentTask, reason string) {
	for _, task := range tasks {
		s.db.Model(&task).Updates(map[string]interface{}{
			"status":        5, // 已取消
			"error_message": reason,
			"end_time":      time.Now(),
			"duration":      0,
		})
	}
}

// executeTask 执行单个任务
func (s *ApplicationService) executeTask(task *model.QuickDeploymentTask) bool {
	// 更新任务状态为部署中
	startTime := time.Now()
	s.db.Model(task).Updates(map[string]interface{}{
		"status":     2, // 部署中
		"start_time": &startTime,
	})

	// 获取Jenkins环境配置
	var jenkinsEnv model.JenkinsEnv
	if err := s.db.First(&jenkinsEnv, task.JenkinsEnvID).Error; err != nil {
		s.updateTaskFailure(task, startTime, "获取Jenkins环境配置失败: "+err.Error())
		return false
	}

	if jenkinsEnv.JenkinsServerID == nil || *jenkinsEnv.JenkinsServerID == 0 {
		s.updateTaskFailure(task, startTime, "Jenkins服务器配置缺失")
		return false
	}

	if jenkinsEnv.JobName == "" {
		s.updateTaskFailure(task, startTime, "Jenkins任务名称未配置")
		return false
	}

	// 启动Jenkins任务（不等待完成）
	started, buildNumber, logURL, err := s.executeJenkinsJob(*jenkinsEnv.JenkinsServerID, jenkinsEnv.JobName, task.Parameters)
	if err != nil {
		s.updateTaskFailure(task, startTime, "执行Jenkins任务失败: "+err.Error())
		return false
	}

	if !started {
		s.updateTaskFailure(task, startTime, "Jenkins任务启动失败")
		return false
	}

	// 验证构建编号有效性
	if buildNumber <= 0 {
		s.updateTaskFailure(task, startTime, fmt.Sprintf("获取到无效的构建编号: %d", buildNumber))
		return false
	}

	// 更新任务为运行中状态，记录构建编号和日志URL
	s.db.Model(task).Updates(map[string]interface{}{
		"status":       2, // 部署中
		"build_number": buildNumber,
		"log_url":      logURL,
	})

	// 启动异步监控goroutine来监控任务完成状态
	go s.monitorJenkinsTask(task.ID, *jenkinsEnv.JenkinsServerID, jenkinsEnv.JobName, buildNumber, startTime)

	return true
}

// monitorJenkinsTask 异步监控Jenkins任务状态
func (s *ApplicationService) monitorJenkinsTask(taskID uint, serverID uint, jobName string, buildNumber int, startTime time.Time) {
	fmt.Printf("INFO: 开始监控Jenkins任务, taskID=%d, buildNumber=%d, jobName=%s\n", taskID, buildNumber, jobName)

	// 获取Jenkins配置
	accountDao := ccdao.NewAccountAuthDao()
	account, err := accountDao.GetByID(serverID)
	if err != nil {
		fmt.Printf("ERROR: 获取Jenkins配置失败: %v\n", err)
		s.updateTaskFailureByID(taskID, startTime, "获取Jenkins配置失败: "+err.Error())
		return
	}

	password, err := account.DecryptPassword()
	if err != nil {
		fmt.Printf("ERROR: 解密Jenkins密码失败: %v\n", err)
		s.updateTaskFailureByID(taskID, startTime, "解密Jenkins密码失败: "+err.Error())
		return
	}

	jenkinsClient := NewJenkinsClient(account.Host, account.Port, account.Name, password)
	baseURL := jenkinsClient.BaseURL

	// 循环监控直到任务完成
	checkInterval := 5 * time.Second  // 改为5秒检查一次，更快响应
	maxWaitTime := 30 * time.Minute
	monitorStartTime := time.Now()
	checkCount := 0

	for {
		checkCount++

		// 检查监控超时
		if time.Since(monitorStartTime) > maxWaitTime {
			fmt.Printf("ERROR: Jenkins任务监控超时, taskID=%d, buildNumber=%d\n", taskID, buildNumber)
			s.updateTaskFailureByID(taskID, startTime, "Jenkins任务监控超时")
			return
		}

		// 检查构建状态
		isBuilding, result, _, err := s.getBuildStatusDetailed(baseURL, account.Name, password, jobName, buildNumber)
		if err != nil {
			// 网络错误等，等待后重试
			fmt.Printf("WARN: 获取构建状态失败(第%d次检查): %v，等待后重试\n", checkCount, err)
			time.Sleep(checkInterval)
			continue
		}

		fmt.Printf("DEBUG: 构建状态检查(第%d次) - taskID=%d, buildNumber=%d, isBuilding=%v, result=%s\n",
			checkCount, taskID, buildNumber, isBuilding, result)

		// 如果任务完成，更新最终状态
		if !isBuilding {
			fmt.Printf("INFO: Jenkins任务完成, taskID=%d, buildNumber=%d, result=%s\n", taskID, buildNumber, result)
			s.updateTaskCompletion(taskID, startTime, result)
			return
		}

		// 任务还在运行，等待后继续检查
		time.Sleep(checkInterval)
	}
}

// updateTaskFailureByID 通过任务ID更新失败状态
func (s *ApplicationService) updateTaskFailureByID(taskID uint, startTime time.Time, errorMessage string) {
	endTime := time.Now()
	duration := int(endTime.Sub(startTime).Seconds())

	s.db.Model(&model.QuickDeploymentTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":        4, // 异常
		"end_time":      &endTime,
		"duration":      duration,
		"error_message": errorMessage,
	})

	// 检查并更新父任务状态
	s.checkAndUpdateDeploymentStatus(taskID)
}

// updateTaskCompletion 更新任务完成状态
func (s *ApplicationService) updateTaskCompletion(taskID uint, startTime time.Time, jenkinsResult string) {
	endTime := time.Now()
	duration := int(endTime.Sub(startTime).Seconds())

	var status int
	var errorMessage string

	switch jenkinsResult {
	case "SUCCESS":
		status = 3 // 成功
	case "FAILURE":
		status = 4 // 异常
		errorMessage = "Jenkins构建失败"
	case "ABORTED":
		status = 4 // 异常
		errorMessage = "Jenkins构建被中止"
	default:
		status = 4 // 异常
		errorMessage = fmt.Sprintf("Jenkins构建结果未知: %s", jenkinsResult)
	}

	s.db.Model(&model.QuickDeploymentTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":        status,
		"end_time":      &endTime,
		"duration":      duration,
		"error_message": errorMessage,
	})

	// 检查并更新父任务状态
	s.checkAndUpdateDeploymentStatus(taskID)
}

// updateTaskFailure 更新任务失败状态
func (s *ApplicationService) updateTaskFailure(task *model.QuickDeploymentTask, startTime time.Time, errorMessage string) {
	endTime := time.Now()
	duration := int(endTime.Sub(startTime).Seconds())

	s.db.Model(task).Updates(map[string]interface{}{
		"status":        4, // 异常
		"end_time":      &endTime,
		"duration":      duration,
		"error_message": errorMessage,
	})
}

// createJenkinsEnvFromRequest 从请求创建Jenkins环境配置
func (s *ApplicationService) createJenkinsEnvFromRequest(req *model.CreateJenkinsEnvRequest) error {
	env := &model.JenkinsEnv{
		AppID:           req.AppID,
		EnvName:         req.EnvName,
		JenkinsServerID: req.JenkinsServerID,
		JobName:         req.JobName,
	}

	return s.appDao.CreateJenkinsEnv(env)
}



// updateJenkinsEnvs 更新应用的Jenkins环境配置
func (s *ApplicationService) updateJenkinsEnvs(appID uint, envs []model.UpdateJenkinsEnvRequest) error {
	// 获取现有的Jenkins环境配置
	existingEnvs, err := s.appDao.GetJenkinsEnvsByAppID(appID)
	if err != nil {
		return err
	}

	// 构建现有环境配置映射（按ID）
	existingEnvMap := make(map[uint]*model.JenkinsEnv)
	for i := range existingEnvs {
		existingEnvMap[existingEnvs[i].ID] = &existingEnvs[i]
	}

	// 构建要保留的环境配置ID集合
	keepEnvIDs := make(map[uint]bool)

	// 处理每个环境配置更新请求
	for _, envReq := range envs {
		if envReq.ID != nil {
			// 更新现有环境配置
			envID := *envReq.ID
			keepEnvIDs[envID] = true

			if existingEnv, exists := existingEnvMap[envID]; exists {
				// 验证环境配置属于该应用
				if existingEnv.AppID != appID {
					continue // 跳过不属于该应用的环境配置
				}

				// 构建更新字段
				updates := make(map[string]interface{})
				if envReq.EnvName != nil {
					updates["env_name"] = *envReq.EnvName
				}
				if envReq.JenkinsServerID != nil {
					updates["jenkins_server_id"] = *envReq.JenkinsServerID
				}
				if envReq.JobName != nil {
					updates["job_name"] = *envReq.JobName
				}

				// 执行更新
				s.appDao.UpdateJenkinsEnv(envID, updates)
			}
		} else {
			// 创建新的环境配置
			createReq := &model.CreateJenkinsEnvRequest{
				AppID:           appID,
				EnvName:         *envReq.EnvName,
				JenkinsServerID: envReq.JenkinsServerID,
				JobName:         *envReq.JobName,
			}

			s.createJenkinsEnvFromRequest(createReq)
		}
	}

	// 删除不在更新列表中的环境配置
	for envID := range existingEnvMap {
		if !keepEnvIDs[envID] {
			s.appDao.DeleteJenkinsEnv(envID)
		}
	}

	return nil
}


// GetJenkinsServers 获取Jenkins服务器列表
func (s *ApplicationService) GetJenkinsServers(c *gin.Context) {
	var servers []ccmodel.AccountAuth

	// 查询类型为4(Jenkins)的服务器配置
	if err := s.db.Where("type = ?", 4).Find(&servers).Error; err != nil {
		result.Failed(c, 500, "获取Jenkins服务器列表失败: "+err.Error())
		return
	}

	// 转换为前端所需的格式，只返回ID和名称
	var options []model.JenkinsServerOption
	for _, server := range servers {
		options = append(options, model.JenkinsServerOption{
			ID:   server.ID,
			Name: server.Alias, // 使用别名作为显示名称
		})
	}

	result.Success(c, options)
}

// ValidateJenkinsJob 验证Jenkins任务是否存在
func (s *ApplicationService) ValidateJenkinsJob(c *gin.Context, req *model.ValidateJenkinsJobRequest) {
	// 1. 获取Jenkins服务器配置
	var jenkinsServer ccmodel.AccountAuth
	if err := s.db.Where("id = ? AND type = 4", req.JenkinsServerID).First(&jenkinsServer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "Jenkins服务器配置不存在")
			return
		}
		result.Failed(c, 500, "获取Jenkins服务器配置失败: "+err.Error())
		return
	}

	// 2. 解密密码
	password, err := jenkinsServer.DecryptPassword()
	if err != nil {
		result.Failed(c, 500, "解密Jenkins密码失败: "+err.Error())
		return
	}

	// 3. 使用 NewJenkinsClient 构建正确的 baseURL
	jenkinsClient := NewJenkinsClient(jenkinsServer.Host, jenkinsServer.Port, jenkinsServer.Name, password)
	jobURL := fmt.Sprintf("%s/job/%s", jenkinsClient.BaseURL, req.JobName)
	apiURL := fmt.Sprintf("%s/api/json", jobURL)

	// 4. 创建HTTP客户端并验证任务是否存在
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	httpReq, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		result.Failed(c, 500, "创建HTTP请求失败: "+err.Error())
		return
	}

	// 设置基础认证
	httpReq.SetBasicAuth(jenkinsServer.Name, password)

	// 5. 发送请求
	resp, err := client.Do(httpReq)
	if err != nil {
		response := &model.ValidateJenkinsJobResponse{
			Exists:   false,
			JobName:  req.JobName,
			JobURL:   "",
			Message:  fmt.Sprintf("连接Jenkins服务器失败: %v", err),
			ServerID: req.JenkinsServerID,
		}
		result.Success(c, response)
		return
	}
	defer resp.Body.Close()

	// 6. 根据HTTP状态码判断任务是否存在
	var response *model.ValidateJenkinsJobResponse

	switch resp.StatusCode {
	case 200:
		// 任务存在
		response = &model.ValidateJenkinsJobResponse{
			Exists:   true,
			JobName:  req.JobName,
			JobURL:   jobURL,
			Message:  "任务存在",
			ServerID: req.JenkinsServerID,
		}
	case 404:
		// 任务不存在
		response = &model.ValidateJenkinsJobResponse{
			Exists:   false,
			JobName:  req.JobName,
			JobURL:   "",
			Message:  "任务不存在",
			ServerID: req.JenkinsServerID,
		}
	case 401, 403:
		// 认证失败
		response = &model.ValidateJenkinsJobResponse{
			Exists:   false,
			JobName:  req.JobName,
			JobURL:   "",
			Message:  "Jenkins认证失败，请检查用户名和密码",
			ServerID: req.JenkinsServerID,
		}
	default:
		// 其他错误
		response = &model.ValidateJenkinsJobResponse{
			Exists:   false,
			JobName:  req.JobName,
			JobURL:   "",
			Message:  fmt.Sprintf("验证失败，HTTP状态码: %d", resp.StatusCode),
			ServerID: req.JenkinsServerID,
		}
	}

	result.Success(c, response)
}

// checkAndUpdateAppStatus 检查并更新应用状态
// 当应用至少有一个完整配置的Jenkins环境时，将应用状态设为已激活(status=2)
func (s *ApplicationService) checkAndUpdateAppStatus(appID uint) error {
	// 获取应用的所有Jenkins环境配置
	envs, err := s.appDao.GetJenkinsEnvsByAppID(appID)
	if err != nil {
		return fmt.Errorf("获取应用Jenkins环境配置失败: %v", err)
	}

	// 检查是否有至少一个完整配置的环境
	hasCompleteEnv := false
	for _, env := range envs {
		if env.JenkinsServerID != nil && *env.JenkinsServerID > 0 && env.JobName != "" {
			hasCompleteEnv = true
			break
		}
	}

	// 根据检查结果更新应用状态
	var newStatus int
	if hasCompleteEnv {
		newStatus = 2 // 已激活
	} else {
		newStatus = 1 // 未激活
	}

	// 更新应用状态
	updates := map[string]interface{}{
		"status": newStatus,
	}

	if err := s.appDao.UpdateApplication(appID, updates); err != nil {
		return fmt.Errorf("更新应用状态失败: %v", err)
	}

	return nil
}

// deleteAppJenkinsEnvs 删除应用的所有Jenkins环境配置
func (s *ApplicationService) deleteAppJenkinsEnvs(appID uint) error {
	// 获取要删除的Jenkins环境配置列表（用于日志记录）
	envs, err := s.appDao.GetJenkinsEnvsByAppID(appID)
	if err != nil {
		return fmt.Errorf("获取Jenkins环境配置列表失败: %v", err)
	}

	// 如果没有环境配置，直接返回
	if len(envs) == 0 {
		return nil
	}

	// 删除所有关联的Jenkins环境配置
	for _, env := range envs {
		if err := s.appDao.DeleteJenkinsEnv(env.ID); err != nil {
			return fmt.Errorf("删除Jenkins环境配置 '%s' 失败: %v", env.EnvName, err)
		}
	}

	return nil
}

// GetServiceTree 获取业务线服务树
func (s *ApplicationService) GetServiceTree(c *gin.Context, req *model.GetServiceTreeRequest) {
	// 调试：打印接收到的参数
	fmt.Printf("DEBUG: BusinessGroupIDs = %v, len = %d\n", req.BusinessGroupIDs, len(req.BusinessGroupIDs))
	fmt.Printf("DEBUG: Status = %v\n", req.Status)
	fmt.Printf("DEBUG: Environment = %s\n", req.Environment)

	// 构建查询条件
	query := s.db.Model(&model.Application{}).
		Preload("JenkinsEnvs").
		Select("app_application.*, COUNT(app_jenkins_env.id) as jenkins_env_count").
		Joins("LEFT JOIN app_jenkins_env ON app_application.id = app_jenkins_env.app_id").
		Group("app_application.id")

	// 按业务组筛选
	if len(req.BusinessGroupIDs) > 0 {
		query = query.Where("app_application.business_group_id IN ?", req.BusinessGroupIDs)
	}

	// 按状态筛选
	if req.Status != nil {
		query = query.Where("app_application.status = ?", *req.Status)
	}

	// 按环境筛选Jenkins配置
	if req.Environment != "" {
		query = query.Where("app_jenkins_env.env_name = ? OR app_jenkins_env.env_name IS NULL", req.Environment)
	}

	// 查询应用列表
	var apps []model.Application
	if err := query.Order("app_application.business_group_id, app_application.name").Find(&apps).Error; err != nil {
		result.Failed(c, 500, "获取应用列表失败: "+err.Error())
		return
	}

	// 按业务组分组并构建服务树
	businessLineMap := make(map[uint]*model.BusinessLineServiceTree)
	statusTextMap := map[int]string{
		0: "待配置", 1: "未激活", 2: "已激活", 3: "已停用",
	}

	for _, app := range apps {
		// 如果业务组不存在，创建新的业务线服务树
		if _, exists := businessLineMap[app.BusinessGroupID]; !exists {
			businessLineMap[app.BusinessGroupID] = &model.BusinessLineServiceTree{
				BusinessGroupID:   app.BusinessGroupID,
				BusinessGroupName: s.getBusinessGroupName(app.BusinessGroupID), // 需要从业务组表获取名称
				ServiceCount:      0,
				Services:          []model.ServiceTreeNode{},
			}
		}

		// 构建Jenkins环境配置列表
		jenkinsEnvs := make([]model.ServiceJenkinsEnv, 0, len(app.JenkinsEnvs))
		for _, env := range app.JenkinsEnvs {
			// 如果指定了环境筛选，只包含匹配的环境
			if req.Environment != "" && env.EnvName != req.Environment {
				continue
			}

			isConfigured := env.JenkinsServerID != nil && *env.JenkinsServerID > 0 && env.JobName != ""
			jenkinsEnvs = append(jenkinsEnvs, model.ServiceJenkinsEnv{
				ID:              env.ID,
				EnvName:         env.EnvName,
				JenkinsServerID: env.JenkinsServerID,
				JobName:         env.JobName,
				IsConfigured:    isConfigured,
			})
		}

		// 构建服务树节点
		serviceNode := model.ServiceTreeNode{
			ID:                app.ID,
			Name:              app.Name,
			Code:              app.Code,
			Status:            app.Status,
			StatusText:        statusTextMap[app.Status],
			ProgrammingLang:   app.ProgrammingLang,
			BusinessDeptID:    app.BusinessDeptID,
			BusinessDeptName:  s.getBusinessDeptName(app.BusinessDeptID), // 需要从业务部门表获取名称
			CreatedAt:         app.CreatedAt.Format("2006-01-02 15:04:05"),
			JenkinsEnvs:       jenkinsEnvs,
		}

		// 添加到对应业务线
		businessLine := businessLineMap[app.BusinessGroupID]
		businessLine.Services = append(businessLine.Services, serviceNode)
		businessLine.ServiceCount++
	}

	// 转换为数组格式
	var serviceTree []model.BusinessLineServiceTree
	for _, businessLine := range businessLineMap {
		serviceTree = append(serviceTree, *businessLine)
	}

	// 按业务组ID排序
	for i := 0; i < len(serviceTree)-1; i++ {
		for j := i + 1; j < len(serviceTree); j++ {
			if serviceTree[i].BusinessGroupID > serviceTree[j].BusinessGroupID {
				serviceTree[i], serviceTree[j] = serviceTree[j], serviceTree[i]
			}
		}
	}

	result.Success(c, serviceTree)
}

// getBusinessGroupName 获取业务组名称（占位符方法，需要根据实际业务组表实现）
func (s *ApplicationService) getBusinessGroupName(groupID uint) string {
	// TODO: 从业务组表获取名称，这里暂时返回格式化字符串
	return fmt.Sprintf("业务组_%d", groupID)
}

// getBusinessDeptName 获取业务部门名称（占位符方法，需要根据实际业务部门表实现）
func (s *ApplicationService) getBusinessDeptName(deptID uint) string {
	// TODO: 从业务部门表获取名称，这里暂时返回格式化字符串
	return fmt.Sprintf("业务部门_%d", deptID)
}

// BuildResult 构建启动结果
type BuildResult struct {
	QueueItemLocation string
	QueueID           int
}

// executeJenkinsJob 启动Jenkins任务（不等待完成）
func (s *ApplicationService) executeJenkinsJob(serverID uint, jobName string, parametersJSON string) (bool, int, string, error) {
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

	// 创建Jenkins客户端（添加cookie jar支持）
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
		fmt.Printf("DEBUG: 解析到的参数: %+v\n", parameters)
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
		fmt.Printf("DEBUG: 构建URL: %s?%s\n", buildURL, formData.Encode())
	} else {
		// 无参数构建
		buildURL = fmt.Sprintf("%s/job/%s/build", baseURL, url.PathEscape(jobName))
		fmt.Printf("DEBUG: 构建URL(无参数): %s\n", buildURL)
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

	fmt.Printf("DEBUG: Jenkins响应状态码: %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// 读取响应体以获取详细错误信息
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("DEBUG: Jenkins响应内容: %s\n", string(bodyBytes))
		return false, 0, "", fmt.Errorf("Jenkins响应错误，状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	// 获取队列位置
	location := resp.Header.Get("Location")
	if location == "" {
		fmt.Printf("WARN: 未获取到Jenkins队列位置，尝试直接获取最新构建\n")
	} else {
		fmt.Printf("INFO: Jenkins任务已加入队列，Location: %s\n", location)
	}

	// 等待新构建开始并获取构建编号
	buildNumber, err := s.waitForNewBuildStart(baseURL, account.Name, password, jobName, location)
	if err != nil {
		fmt.Printf("WARN: 等待新构建开始失败: %v，但任务可能已成功加入队列\n", err)
		// 即使获取buildNumber失败，任务可能已经成功加入队列
		// 尝试再次获取最新的构建编号
		latestBuild, getErr := s.getLatestBuildNumber(baseURL, account.Name, password, jobName)
		if getErr == nil && latestBuild > 0 {
			buildNumber = latestBuild
			fmt.Printf("INFO: 获取到最新构建编号: %d\n", buildNumber)
		} else {
			// 如果还是获取不到，返回错误
			return false, 0, "", fmt.Errorf("等待新构建开始失败: %v", err)
		}
	}

	fmt.Printf("INFO: Jenkins任务启动成功，构建编号: %d\n", buildNumber)

	// 获取构建URL
	logURL := fmt.Sprintf("%s/job/%s/%d/console", baseURL, url.PathEscape(jobName), buildNumber)

	// 任务启动成功，返回构建编号和日志URL
	// 注意：这里不等待构建完成，返回true表示启动成功
	return true, buildNumber, logURL, nil
}

// getBuildNumberFromQueue 从队列API获取构建编号
func (s *ApplicationService) getBuildNumberFromQueue(baseURL, username, password, queueLocation string) (int, error) {
	// queueLocation 格式如: http://jenkins/queue/item/123/
	// 需要转换为 API URL
	queueAPIURL := strings.TrimRight(queueLocation, "/") + "/api/json"

	fmt.Printf("DEBUG: 从队列获取构建号，URL: %s\n", queueAPIURL)

	maxWaitTime := 2 * time.Minute // 最大等待2分钟
	checkInterval := 2 * time.Second
	startTime := time.Now()

	for {
		if time.Since(startTime) > maxWaitTime {
			return 0, fmt.Errorf("从队列获取构建号超时")
		}

		req, err := http.NewRequest("GET", queueAPIURL, nil)
		if err != nil {
			return 0, err
		}

		req.SetBasicAuth(username, password)

		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		if err != nil {
			time.Sleep(checkInterval)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			time.Sleep(checkInterval)
			continue
		}

		var queueData struct {
			Executable struct {
				Number int `json:"number"`
			} `json:"executable"`
			Cancelled bool `json:"cancelled"`
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(bodyBytes, &queueData); err != nil {
			time.Sleep(checkInterval)
			continue
		}

		// 如果任务被取消
		if queueData.Cancelled {
			return 0, fmt.Errorf("队列任务被取消")
		}

		// 如果已经有构建编号，说明任务已经开始执行
		if queueData.Executable.Number > 0 {
			fmt.Printf("INFO: 从队列获取到构建号: %d\n", queueData.Executable.Number)
			return queueData.Executable.Number, nil
		}

		// 继续等待
		time.Sleep(checkInterval)
	}
}

// getLatestBuildNumber 获取最新构建编号
func (s *ApplicationService) getLatestBuildNumber(baseURL, username, password, jobName string) (int, error) {
	jobURL := fmt.Sprintf("%s/job/%s/api/json?tree=lastBuild[number]", baseURL, url.PathEscape(jobName))

	fmt.Printf("DEBUG: 获取最新构建号，URL: %s\n", jobURL)

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

	// 读取响应体用于调试
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %v", err)
	}

	fmt.Printf("DEBUG: Jenkins响应状态码: %d, 响应内容前200字符: %s\n",
		resp.StatusCode, string(bodyBytes[:min(200, len(bodyBytes))]))

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("获取任务信息失败，状态码: %d", resp.StatusCode)
	}

	var jobData struct {
		LastBuild *struct {
			Number int `json:"number"`
		} `json:"lastBuild"`
	}

	if err := json.Unmarshal(bodyBytes, &jobData); err != nil {
		return 0, fmt.Errorf("解析响应失败: %v", err)
	}

	// 如果 lastBuild 为 null，说明还没有任何构建
	if jobData.LastBuild == nil {
		fmt.Printf("DEBUG: lastBuild 为 null，任务尚无构建历史\n")
		return 0, nil
	}

	fmt.Printf("DEBUG: 获取到最新构建号: %d\n", jobData.LastBuild.Number)
	return jobData.LastBuild.Number, nil
}

// waitForNewBuildStart 等待新构建开始并返回新的构建号
func (s *ApplicationService) waitForNewBuildStart(baseURL, username, password, jobName, queueLocation string) (int, error) {
	// 如果有队列Location，先尝试通过队列API获取
	if queueLocation != "" {
		buildNumber, err := s.getBuildNumberFromQueue(baseURL, username, password, queueLocation)
		if err == nil && buildNumber > 0 {
			fmt.Printf("INFO: 通过队列API获取到构建号: %d\n", buildNumber)
			return buildNumber, nil
		}
		fmt.Printf("WARN: 通过队列API获取构建号失败: %v，切换到轮询方式\n", err)
	}

	// 先获取当前最新构建号作为基准
	currentBuildNumber, err := s.getLatestBuildNumber(baseURL, username, password, jobName)
	if err != nil {
		// 如果获取失败，可能是第一次构建，设为0
		fmt.Printf("WARN: 获取当前构建号失败: %v，假设为第一次构建\n", err)
		currentBuildNumber = 0
	} else {
		fmt.Printf("INFO: 当前最新构建号: %d，等待新构建开始...\n", currentBuildNumber)
	}

	maxWaitTime := 5 * time.Minute  // 最大等待5分钟
	checkInterval := 3 * time.Second // 每3秒检查一次
	startTime := time.Now()
	checkCount := 0

	for {
		checkCount++

		// 检查是否超时
		if time.Since(startTime) > maxWaitTime {
			fmt.Printf("ERROR: 等待新构建开始超时，检查了%d次\n", checkCount)
			return 0, fmt.Errorf("等待新构建开始超时")
		}

		// 获取最新构建号
		latestBuildNumber, err := s.getLatestBuildNumber(baseURL, username, password, jobName)
		if err != nil {
			// 如果是400错误，可能是Jenkins配置问题，尝试其他方式
			if strings.Contains(err.Error(), "状态码: 400") {
				fmt.Printf("WARN: Jenkins API返回400错误，可能是权限或配置问题(第%d次)\n", checkCount)
				// 继续重试，也许任务完成后API就可以访问了
			} else {
				fmt.Printf("WARN: 获取最新构建号失败(第%d次): %v\n", checkCount, err)
			}
			time.Sleep(checkInterval)
			continue
		}

		fmt.Printf("DEBUG: 检查构建号(第%d次): current=%d, latest=%d\n", checkCount, currentBuildNumber, latestBuildNumber)

		// 如果构建号增加了，说明新构建已经开始
		if latestBuildNumber > currentBuildNumber {
			fmt.Printf("INFO: 检测到新构建开始，构建号: %d\n", latestBuildNumber)
			return latestBuildNumber, nil
		}

		// 等待后继续检查
		time.Sleep(checkInterval)
	}
}

// waitForBuildCompletion 等待构建完成并获取最终状态
func (s *ApplicationService) waitForBuildCompletion(baseURL, username, password, jobName string, buildNumber int) (bool, string, error) {
	maxWaitTime := 20 * time.Minute // 最大等待20分钟
	checkInterval := 10 * time.Second // 每10秒检查一次
	startTime := time.Now()

	for {
		// 检查是否超时
		if time.Since(startTime) > maxWaitTime {
			return false, "", fmt.Errorf("构建等待超时")
		}

		// 获取构建状态
		isBuilding, result, logURL, err := s.getBuildStatusDetailed(baseURL, username, password, jobName, buildNumber)
		if err != nil {
			return false, "", err
		}

		// 如果构建完成，返回结果
		if !isBuilding {
			success := result == "SUCCESS"
			return success, logURL, nil
		}

		// 等待后继续检查
		time.Sleep(checkInterval)
	}
}

// getBuildStatusDetailed 获取详细的构建状态
func (s *ApplicationService) getBuildStatusDetailed(baseURL, username, password, jobName string, buildNumber int) (bool, string, string, error) {
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
		Result   string `json:"result"`   // 可能是 "SUCCESS", "FAILURE", "ABORTED" 或 null
		URL      string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&buildData); err != nil {
		return false, "", "", err
	}

	logURL := buildData.URL + "console"

	// 如果还在构建中，result可能为空字符串
	if buildData.Building {
		return true, "", logURL, nil
	}

	// 构建完成，返回最终结果
	return false, buildData.Result, logURL, nil
}

// getBuildStatus 获取构建状态（保留原方法用于兼容）
func (s *ApplicationService) getBuildStatus(baseURL, username, password, jobName string, buildNumber int) (bool, string, error) {
	_, result, logURL, err := s.getBuildStatusDetailed(baseURL, username, password, jobName, buildNumber)
	if err != nil {
		return false, "", err
	}

	success := result == "SUCCESS"
	return success, logURL, nil
}

// GetTaskBuildLog 获取任务的Jenkins构建日志
func (s *ApplicationService) GetTaskBuildLog(c *gin.Context, taskID uint, start int) {
	// 获取任务信息
	var task model.QuickDeploymentTask
	if err := s.db.Preload("JenkinsEnv").First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "任务不存在")
			return
		}
		result.Failed(c, 500, "获取任务信息失败: "+err.Error())
		return
	}

	// 检查任务是否有构建编号
	if task.BuildNumber == 0 {
		result.Failed(c, 400, "任务尚未执行或未获取到构建编号")
		return
	}

	// 获取Jenkins环境配置
	if task.JenkinsEnv.JenkinsServerID == nil || *task.JenkinsEnv.JenkinsServerID == 0 {
		result.Failed(c, 400, "Jenkins服务器配置缺失")
		return
	}

	// 获取Jenkins服务器配置
	accountDao := ccdao.NewAccountAuthDao()
	account, err := accountDao.GetByID(*task.JenkinsEnv.JenkinsServerID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins服务器配置失败: "+err.Error())
		return
	}

	password, err := account.DecryptPassword()
	if err != nil {
		result.Failed(c, 500, "解密Jenkins密码失败: "+err.Error())
		return
	}

	// 获取构建日志
	jenkinsClient := NewJenkinsClient(account.Host, account.Port, account.Name, password)
	baseURL := jenkinsClient.BaseURL

	// 使用 consoleText 端点获取纯文本日志（不含ANSI转义序列）
	logURL := fmt.Sprintf("%s/job/%s/%d/consoleText",
		baseURL, url.PathEscape(task.JenkinsEnv.JobName), task.BuildNumber)

	fmt.Printf("DEBUG: 获取日志URL: %s\n", logURL)

	req, err := http.NewRequest("GET", logURL, nil)
	if err != nil {
		result.Failed(c, 500, "创建日志请求失败: "+err.Error())
		return
	}

	req.SetBasicAuth(account.Name, password)
	// 明确请求纯文本格式
	req.Header.Set("Accept", "text/plain")

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		result.Failed(c, 500, "获取构建日志失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Failed(c, 500, fmt.Sprintf("Jenkins响应错误，状态码: %d", resp.StatusCode))
		return
	}

	// 读取所有日志内容
	logBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Failed(c, 500, "读取日志内容失败: "+err.Error())
		return
	}

	logContent := string(logBytes)

	// 如果需要分页，可以根据 start 参数截取
	if start > 0 && start < len(logContent) {
		logContent = logContent[start:]
	}

	logResponse := map[string]interface{}{
		"log":          logContent,
		"has_more":     false, // consoleText 返回完整日志
		"text_size":    len(logBytes),
		"task_id":      taskID,
		"build_number": task.BuildNumber,
		"job_name":     task.JenkinsEnv.JobName,
	}

	result.Success(c, logResponse)
}

// GetBusinessGroupOptions 获取业务组选项（连级选择器数据源）
func (s *ApplicationService) GetBusinessGroupOptions(c *gin.Context) {
	// 查询所有有应用的业务组
	var businessGroups []struct {
		BusinessGroupID   uint   `json:"business_group_id"`
		BusinessGroupName string `json:"business_group_name"`
		ServiceCount      int    `json:"service_count"`
	}

	// 执行SQL查询，按业务组分组统计应用数量
	query := `
		SELECT
			business_group_id,
			CONCAT('业务组_', business_group_id) as business_group_name,
			COUNT(*) as service_count
		FROM app_application
		WHERE deleted_at IS NULL
		GROUP BY business_group_id
		HAVING service_count > 0
		ORDER BY business_group_id
	`

	if err := s.db.Raw(query).Scan(&businessGroups).Error; err != nil {
		result.Failed(c, 500, "获取业务组选项失败: "+err.Error())
		return
	}

	// 构建连级选择器数据结构
	var options []map[string]interface{}
	for _, group := range businessGroups {
		// 查询该业务组下的业务部门
		var depts []struct {
			BusinessDeptID   uint   `json:"business_dept_id"`
			BusinessDeptName string `json:"business_dept_name"`
			ServiceCount     int    `json:"service_count"`
		}

		deptQuery := `
			SELECT
				business_dept_id,
				CONCAT('业务部门_', business_dept_id) as business_dept_name,
				COUNT(*) as service_count
			FROM app_application
			WHERE deleted_at IS NULL AND business_group_id = ?
			GROUP BY business_dept_id
			HAVING service_count > 0
			ORDER BY business_dept_id
		`

		s.db.Raw(deptQuery, group.BusinessGroupID).Scan(&depts)

		// 构建部门选项
		var children []map[string]interface{}
		for _, dept := range depts {
			children = append(children, map[string]interface{}{
				"value": dept.BusinessDeptID,
				"label": fmt.Sprintf("%s (%d个服务)", dept.BusinessDeptName, dept.ServiceCount),
				"dept_id": dept.BusinessDeptID,
				"dept_name": dept.BusinessDeptName,
			})
		}

		// 构建业务组选项
		option := map[string]interface{}{
			"value": group.BusinessGroupID,
			"label": fmt.Sprintf("%s (%d个服务)", group.BusinessGroupName, group.ServiceCount),
			"group_id": group.BusinessGroupID,
			"group_name": group.BusinessGroupName,
			"children": children,
		}

		options = append(options, option)
	}

	result.Success(c, options)
}

// GetAppEnvironment 获取单个应用的环境配置
func (s *ApplicationService) GetAppEnvironment(c *gin.Context, req *model.GetAppEnvironmentRequest) {
	// 获取应用信息
	app, err := s.appDao.GetApplicationByID(req.AppID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "应用不存在")
			return
		}
		result.Failed(c, 500, "获取应用信息失败: "+err.Error())
		return
	}

	// 获取应用的Jenkins环境配置
	envs, err := s.appDao.GetJenkinsEnvsByAppID(req.AppID)
	if err != nil {
		result.Failed(c, 500, "获取Jenkins环境配置失败: "+err.Error())
		return
	}

	// 查找指定环境配置
	var targetEnv *model.JenkinsEnv
	for _, env := range envs {
		if env.EnvName == req.Environment {
			targetEnv = &env
			break
		}
	}

	// 构建响应数据
	statusTextMap := map[int]string{
		0: "待配置", 1: "未激活", 2: "已激活", 3: "已停用",
	}

	response := &model.AppEnvironmentResponse{
		AppID:           app.ID,
		AppName:         app.Name,
		AppCode:         app.Code,
		Environment:     req.Environment,
		Status:          app.Status,
		StatusText:      statusTextMap[app.Status],
		BusinessGroupID: app.BusinessGroupID,
		BusinessDeptID:  app.BusinessDeptID,
		ProgrammingLang: app.ProgrammingLang,
	}

	if targetEnv != nil {
		// 环境已配置
		response.IsConfigured = targetEnv.JenkinsServerID != nil && *targetEnv.JenkinsServerID > 0 && targetEnv.JobName != ""
		response.JenkinsServerID = targetEnv.JenkinsServerID
		response.JobName = targetEnv.JobName

		// 获取Jenkins服务器名称
		if targetEnv.JenkinsServerID != nil && *targetEnv.JenkinsServerID > 0 {
			accountDao := ccdao.NewAccountAuthDao()
			if account, err := accountDao.GetByID(*targetEnv.JenkinsServerID); err == nil {
				response.JenkinsServerName = account.Alias
				// 构建Jenkins任务URL
				if targetEnv.JobName != "" {
					jenkinsClient := NewJenkinsClient(account.Host, account.Port, "", "")
					response.JenkinsJobURL = fmt.Sprintf("%s/job/%s", jenkinsClient.BaseURL, targetEnv.JobName)
				}
			}
		}
	} else {
		// 环境未配置
		response.IsConfigured = false
	}

	result.Success(c, response)
}

// GetTaskStatus 获取任务状态
func (s *ApplicationService) GetTaskStatus(c *gin.Context, taskID uint) {
	// 参数验证
	if taskID == 0 {
		result.Failed(c, 400, "无效的任务ID")
		return
	}

	// 获取任务信息
	var task model.QuickDeploymentTask
	if err := s.db.First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "任务不存在")
			return
		}
		result.Failed(c, 500, "获取任务信息失败: "+err.Error())
		return
	}

	// 构建状态文本映射
	statusTextMap := map[int]string{
		1: "未部署",
		2: "部署中",
		3: "成功",
		4: "异常",
	}

	// 安全获取状态文本
	statusText := statusTextMap[task.Status]
	if statusText == "" {
		statusText = "未知状态"
	}

	// 计算进度百分比（添加异常处理）
	var progress int
	func() {
		defer func() {
			if r := recover(); r != nil {
				progress = 0 // 如果计算出错，设为0
			}
		}()
		progress = s.calculateTaskProgress(&task)
	}()

	// 构建响应数据
	response := &model.TaskStatusResponse{
		TaskID:       task.ID,
		Status:       task.Status,
		StatusText:   statusText,
		AppName:      task.AppName,
		AppCode:      task.AppCode,
		Environment:  task.Environment,
		BuildNumber:  task.BuildNumber,
		StartTime:    task.StartTime,
		EndTime:      task.EndTime,
		Duration:     task.Duration,
		ErrorMessage: task.ErrorMessage,
		LogURL:       task.LogURL,
		Progress:     progress,
	}

	result.Success(c, response)
}

// calculateTaskProgress 计算任务进度百分比
func (s *ApplicationService) calculateTaskProgress(task *model.QuickDeploymentTask) int {
	// 防御性检查
	if task == nil {
		return 0
	}

	switch task.Status {
	case 1: // 未部署
		return 0
	case 2: // 部署中
		if task.StartTime == nil {
			return 10 // 刚开始
		}

		// 防止时间异常
		now := time.Now()
		if task.StartTime.After(now) {
			return 10 // 开始时间异常，返回初始进度
		}

		// 根据运行时间估算进度，假设一般任务需要5分钟
		runningTime := now.Sub(*task.StartTime).Minutes()

		// 防止负数和异常值
		if runningTime < 0 {
			return 10
		}

		estimatedTime := 5.0 // 5分钟
		progressRatio := runningTime / estimatedTime
		progress := int(progressRatio * 80) + 10 // 10%-90%

		// 限制进度范围
		if progress < 10 {
			progress = 10
		}
		if progress > 90 {
			progress = 90
		}

		return progress
	case 3: // 成功
		return 100
	case 4: // 异常
		return 100
	default:
		return 0
	}
}

// DeleteQuickDeployment 删除快速发布（级联删除所有子任务）
func (s *ApplicationService) DeleteQuickDeployment(c *gin.Context, id uint) {
	// 参数验证
	if id == 0 {
		result.Failed(c, 400, "无效的发布ID")
		return
	}

	// 检查发布是否存在
	var deployment model.QuickDeployment
	if err := s.db.First(&deployment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "发布记录不存在")
			return
		}
		result.Failed(c, 500, "查询发布记录失败: "+err.Error())
		return
	}

	// 检查发布状态，正在执行中的发布不允许删除
	if deployment.Status == 2 { // 2=发布中
		result.Failed(c, 400, "正在发布中的任务不能删除，请等待发布完成或先取消发布")
		return
	}

	// 使用事务删除
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 删除所有子任务（由于配置了CASCADE，这一步应该自动执行，但为了确保，我们手动删除）
	if err := tx.Where("deployment_id = ?", id).Delete(&model.QuickDeploymentTask{}).Error; err != nil {
		tx.Rollback()
		result.Failed(c, 500, "删除发布任务失败: "+err.Error())
		return
	}

	// 2. 删除父发布记录
	if err := tx.Delete(&deployment).Error; err != nil {
		tx.Rollback()
		result.Failed(c, 500, "删除发布记录失败: "+err.Error())
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		result.Failed(c, 500, "提交删除操作失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{
		"message": "删除成功",
		"deleted_deployment_id": id,
	})
}

// calculateTotalTaskDuration 计算发布下所有子任务的总耗时
func (s *ApplicationService) calculateTotalTaskDuration(deploymentID uint) int {
	var tasks []model.QuickDeploymentTask
	if err := s.db.Where("deployment_id = ?", deploymentID).Find(&tasks).Error; err != nil {
		return 0
	}

	totalDuration := 0
	for _, task := range tasks {
		totalDuration += task.Duration
	}

	return totalDuration
}

// checkAndUpdateDeploymentStatus 检查并更新发布状态
func (s *ApplicationService) checkAndUpdateDeploymentStatus(taskID uint) {
	// 获取任务信息
	var task model.QuickDeploymentTask
	if err := s.db.First(&task, taskID).Error; err != nil {
		return
	}

	deploymentID := task.DeploymentID

	// 检查该发布下的所有任务状态
	var tasks []model.QuickDeploymentTask
	if err := s.db.Where("deployment_id = ?", deploymentID).Find(&tasks).Error; err != nil {
		return
	}

	// 检查是否所有任务都已完成（状态为3成功或4失败）
	allCompleted := true
	allSuccess := true
	for _, t := range tasks {
		if t.Status != 3 && t.Status != 4 { // 3=成功, 4=失败
			allCompleted = false
			break
		}
		if t.Status != 3 {
			allSuccess = false
		}
	}

	// 如果所有任务都完成了，更新父任务状态和总时间
	if allCompleted {
		endTime := time.Now()
		totalDuration := s.calculateTotalTaskDuration(deploymentID)

		var finalStatus int
		if allSuccess {
			finalStatus = 3 // 发布成功
		} else {
			finalStatus = 4 // 发布失败
		}

		s.db.Model(&model.QuickDeployment{}).Where("id = ?", deploymentID).Updates(map[string]interface{}{
			"status":   finalStatus,
			"end_time": &endTime,
			"duration": totalDuration,
		})

		// 发送通知
		s.sendDeploymentNotification(deploymentID, finalStatus, allSuccess)
	}
}

// GetApplicationViewTree 获取应用视图树（四层结构：一级业务组->二级业务组->应用->环境）
// @Summary 获取应用视图树
// @Description 获取四层结构的应用视图树：一级业务组->二级业务组->应用->环境
// @Tags Application
// @Accept json
// @Produce json
// @Param business_group_ids[] query []uint false "一级业务组ID列表，为空则查询所有"
// @Param status query int false "应用状态筛选，为空则查询所有状态"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=[]model.ApplicationViewTree}
// @Router /apps/view-tree [get]
func (s *ApplicationService) GetApplicationViewTree(c *gin.Context, req *model.ApplicationViewTreeRequest) {
	// 1. 查询所有应用及其Jenkins环境配置
	query := s.db.Model(&model.Application{}).Preload("JenkinsEnvs")

	// 按状态筛选
	if req.Status != nil {
		query = query.Where("app_application.status = ?", *req.Status)
	}

	var apps []model.Application
	if err := query.Order("app_application.business_group_id, app_application.name").Find(&apps).Error; err != nil {
		result.Failed(c, 500, "获取应用列表失败: "+err.Error())
		return
	}

	// 2. 查询所有业务组（cmdb_group表）
	type GroupInfo struct {
		ID       uint   `gorm:"column:id"`
		ParentID uint   `gorm:"column:parent_id"`
		Name     string `gorm:"column:name"`
	}
	var allGroups []GroupInfo
	if err := s.db.Table("cmdb_group").Find(&allGroups).Error; err != nil {
		result.Failed(c, 500, "获取业务组列表失败: "+err.Error())
		return
	}

	// 3. 构建业务组映射（ID -> 业务组信息）
	groupMap := make(map[uint]GroupInfo)
	for _, g := range allGroups {
		groupMap[g.ID] = g
	}

	// 4. 查询部门信息映射
	type DeptInfo struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:dept_name"`
	}
	var depts []DeptInfo
	s.db.Table("sys_dept").Find(&depts)
	deptMap := make(map[uint]string)
	for _, d := range depts {
		deptMap[d.ID] = d.Name
	}

	// 5. 查询Jenkins服务器映射
	var jenkinsServers []ccmodel.AccountAuth
	s.db.Where("type = ?", 4).Find(&jenkinsServers)
	jenkinsServerMap := make(map[uint]string)
	for _, js := range jenkinsServers {
		jenkinsServerMap[js.ID] = js.Alias
	}

	// 6. 应用状态映射
	statusTextMap := map[int]string{
		0: "待配置", 1: "未激活", 2: "已激活", 3: "已停用",
	}

	// 7. 按一级业务组分组处理
	// 先找出所有一级业务组（parent_id=0）
	firstLevelGroupMap := make(map[uint]*model.ApplicationViewTree)
	for _, g := range allGroups {
		if g.ParentID == 0 {
			// 如果指定了业务组筛选，只处理指定的业务组
			if len(req.BusinessGroupIDs) > 0 {
				found := false
				for _, id := range req.BusinessGroupIDs {
					if id == g.ID {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			firstLevelGroupMap[g.ID] = &model.ApplicationViewTree{
				GroupID:   g.ID,
				GroupName: g.Name,
				SubGroups: []model.ApplicationViewSubGroup{},
			}
		}
	}

	// 8. 处理每个应用，组织成四层结构
	for _, app := range apps {
		// 找到应用所属的业务组
		businessGroup, exists := groupMap[app.BusinessGroupID]
		if !exists {
			continue
		}

		// 判断业务组是一级还是二级
		var firstLevelGroupID uint
		var secondLevelGroupID uint
		var secondLevelGroupName string

		if businessGroup.ParentID == 0 {
			// 业务组本身是一级，没有二级
			firstLevelGroupID = businessGroup.ID
			secondLevelGroupID = 0 // 没有二级，设为0
			secondLevelGroupName = "默认分组"
		} else {
			// 业务组是二级，需要找到父级（一级）
			secondLevelGroupID = businessGroup.ID
			secondLevelGroupName = businessGroup.Name
			firstLevelGroupID = businessGroup.ParentID
		}

		// 检查一级业务组是否存在于映射中
		firstLevelGroup, exists := firstLevelGroupMap[firstLevelGroupID]
		if !exists {
			// 如果不存在，可能是因为筛选条件，跳过此应用
			continue
		}

		// 查找或创建二级业务组
		var targetSubGroup *model.ApplicationViewSubGroup
		for i := range firstLevelGroup.SubGroups {
			if firstLevelGroup.SubGroups[i].SubGroupID == secondLevelGroupID {
				targetSubGroup = &firstLevelGroup.SubGroups[i]
				break
			}
		}
		if targetSubGroup == nil {
			// 创建新的二级业务组
			newSubGroup := model.ApplicationViewSubGroup{
				SubGroupID:   secondLevelGroupID,
				SubGroupName: secondLevelGroupName,
				Applications: []model.ApplicationViewService{},
			}
			firstLevelGroup.SubGroups = append(firstLevelGroup.SubGroups, newSubGroup)
			targetSubGroup = &firstLevelGroup.SubGroups[len(firstLevelGroup.SubGroups)-1]
		}

		// 构建应用的环境列表
		environments := []model.ApplicationViewEnvironment{}
		for _, env := range app.JenkinsEnvs {
			isConfigured := env.JenkinsServerID != nil && *env.JenkinsServerID > 0 && env.JobName != ""
			jenkinsServerName := ""
			if env.JenkinsServerID != nil {
				jenkinsServerName = jenkinsServerMap[*env.JenkinsServerID]
			}

			environments = append(environments, model.ApplicationViewEnvironment{
				EnvID:            env.ID,
				EnvName:          env.EnvName,
				JenkinsServerID:  env.JenkinsServerID,
				JenkinsServerName: jenkinsServerName,
				JobName:          env.JobName,
				IsConfigured:     isConfigured,
				AppID:            app.ID, // 关联应用ID，用于点击时获取详情
			})
		}

		// 构建应用节点
		appNode := model.ApplicationViewService{
			AppID:           app.ID,
			AppName:         app.Name,
			AppCode:         app.Code,
			Status:          app.Status,
			StatusText:      statusTextMap[app.Status],
			ProgrammingLang: app.ProgrammingLang,
			Description:     app.Description,
			RepoURL:         app.RepoURL,
			DeptID:          app.BusinessDeptID,
			DeptName:        deptMap[app.BusinessDeptID],
			CreatedAt:       app.CreatedAt.Format("2006-01-02 15:04:05"),
			Environments:    environments,
		}

		// 添加到二级业务组
		targetSubGroup.Applications = append(targetSubGroup.Applications, appNode)
	}

	// 9. 转换为数组返回
	var viewTree []model.ApplicationViewTree
	for _, tree := range firstLevelGroupMap {
		viewTree = append(viewTree, *tree)
	}

	result.Success(c, viewTree)
}

// GetApplicationDetailForView 获取应用完整详情（用于详情面板）
// @Summary 获取应用完整详情
// @Description 获取应用的完整详情信息，包括负责人姓名、环境配置等，用于右侧详情面板展示
// @Tags Application
// @Accept json
// @Produce json
// @Param id path int true "应用ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.ApplicationDetailResponse}
// @Router /apps/{id}/detail [get]
func (s *ApplicationService) GetApplicationDetailForView(c *gin.Context, id uint) {
	// 1. 获取应用基本信息及关联数据
	var app model.Application
	if err := s.db.Preload("JenkinsEnvs").First(&app, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "应用不存在")
			return
		}
		result.Failed(c, 500, "获取应用信息失败: "+err.Error())
		return
	}

	// 2. 获取业务组名称
	var businessGroup struct {
		Name string `gorm:"column:name"`
	}
	s.db.Table("cmdb_group").Where("id = ?", app.BusinessGroupID).First(&businessGroup)

	// 3. 获取业务部门名称
	var businessDept struct {
		Name string `gorm:"column:dept_name"`
	}
	s.db.Table("sys_dept").Where("id = ?", app.BusinessDeptID).First(&businessDept)

	// 4. 获取负责人姓名
	type AdminInfo struct {
		ID       uint   `gorm:"column:id"`
		Nickname string `gorm:"column:nickname"` // 使用nickname字段
	}
	var allAdmins []AdminInfo
	s.db.Table("sys_admin").Find(&allAdmins)
	adminNameMap := make(map[uint]string)
	for _, admin := range allAdmins {
		adminNameMap[admin.ID] = admin.Nickname
	}

	// 转换负责人ID为姓名
	devOwnerNames := []string{}
	for _, id := range app.DevOwners {
		if name, ok := adminNameMap[id]; ok {
			devOwnerNames = append(devOwnerNames, name)
		}
	}
	testOwnerNames := []string{}
	for _, id := range app.TestOwners {
		if name, ok := adminNameMap[id]; ok {
			testOwnerNames = append(testOwnerNames, name)
		}
	}
	opsOwnerNames := []string{}
	for _, id := range app.OpsOwners {
		if name, ok := adminNameMap[id]; ok {
			opsOwnerNames = append(opsOwnerNames, name)
		}
	}

	// 5. 查询Jenkins服务器映射
	var jenkinsServers []ccmodel.AccountAuth
	s.db.Where("type = ?", 4).Find(&jenkinsServers)
	jenkinsServerMap := make(map[uint]ccmodel.AccountAuth)
	for _, js := range jenkinsServers {
		jenkinsServerMap[js.ID] = js
	}

	// 6. 构建环境配置列表（带完整信息）
	environments := []model.ApplicationDetailEnvironment{}
	for _, env := range app.JenkinsEnvs {
		isConfigured := env.JenkinsServerID != nil && *env.JenkinsServerID > 0 && env.JobName != ""
		jenkinsServerName := ""
		jobURL := ""

		if env.JenkinsServerID != nil {
			if server, ok := jenkinsServerMap[*env.JenkinsServerID]; ok {
				jenkinsServerName = server.Alias
				// 构建Jenkins任务URL
				if env.JobName != "" {
					jenkinsClient := NewJenkinsClient(server.Host, server.Port, "", "")
					jobURL = fmt.Sprintf("%s/job/%s", jenkinsClient.BaseURL, env.JobName)
				}
			}
		}

		environments = append(environments, model.ApplicationDetailEnvironment{
			EnvID:            env.ID,
			EnvName:          env.EnvName,
			JenkinsServerID:  env.JenkinsServerID,
			JenkinsServerName: jenkinsServerName,
			JobName:          env.JobName,
			JobURL:           jobURL,
			IsConfigured:     isConfigured,
		})
	}

	// 7. 应用状态映射
	statusTextMap := map[int]string{
		0: "待配置", 1: "未激活", 2: "已激活", 3: "已停用",
	}

	// 8. 构建响应
	response := model.ApplicationDetailResponse{
		// 基本信息
		AppID:       app.ID,
		AppName:     app.Name,
		AppCode:     app.Code,
		Status:      app.Status,
		StatusText:  statusTextMap[app.Status],
		Description: app.Description,
		RepoURL:     app.RepoURL,

		// 业务信息
		BusinessGroupID:   app.BusinessGroupID,
		BusinessGroupName: businessGroup.Name,
		BusinessDeptID:    app.BusinessDeptID,
		BusinessDeptName:  businessDept.Name,

		// 负责人信息
		DevOwners:      app.DevOwners,
		TestOwners:     app.TestOwners,
		OpsOwners:      app.OpsOwners,
		DevOwnerNames:  devOwnerNames,
		TestOwnerNames: testOwnerNames,
		OpsOwnerNames:  opsOwnerNames,

		// 技术信息
		ProgrammingLang: app.ProgrammingLang,
		StartCommand:    app.StartCommand,
		StopCommand:     app.StopCommand,
		HealthAPI:       app.HealthAPI,

		// 关联资源
		Domains:   app.Domains,
		Hosts:     app.Hosts,
		Databases: app.Databases,
		OtherRes:  app.OtherRes,

		// 环境配置
		Environments: environments,

		// 时间信息
		CreatedAt: app.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: app.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	result.Success(c, response)
}

// GetApplicationReleaseHistory 获取应用发布历史
// @Summary 获取应用发布历史
// @Description 获取指定应用在指定环境的发布历史记录
// @Tags Application
// @Accept json
// @Produce json
// @Param app_id query int true "应用ID"
// @Param environment query string false "环境名称（可选，为空则查询所有环境）"
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.ApplicationReleaseHistoryResponse}
// @Router /apps/release-history [get]
func (s *ApplicationService) GetApplicationReleaseHistory(c *gin.Context, req *model.ApplicationReleaseHistoryRequest) {
	// 1. 获取应用基本信息
	var app model.Application
	if err := s.db.First(&app, req.AppID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed(c, 404, "应用不存在")
			return
		}
		result.Failed(c, 500, "获取应用信息失败: "+err.Error())
		return
	}

	// 2. 构建查询
	query := s.db.Table("app_service_release_item AS item").
		Select(`
			item.id AS item_id,
			item.release_id,
			sr.title,
			je.env_name AS environment,
			sr.applicant_id,
			sr.applicant_name,
			item.status,
			item.branch,
			item.commit_id,
			item.build_number,
			item.jenkins_job_url,
			item.start_time,
			item.end_time,
			item.duration,
			item.impact_feature,
			item.function_module,
			item.created_at
		`).
		Joins("LEFT JOIN app_service_release AS sr ON item.release_id = sr.id").
		Joins("LEFT JOIN app_jenkins_env AS je ON item.jenkins_env_id = je.id").
		Where("item.app_id = ? AND item.deleted_at IS NULL", req.AppID)

	// 按环境筛选
	if req.Environment != "" {
		query = query.Where("je.env_name = ?", req.Environment)
	}

	// 3. 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		result.Failed(c, 500, "查询发布历史总数失败: "+err.Error())
		return
	}

	// 4. 分页查询
	var records []struct {
		ItemID          uint       `gorm:"column:item_id"`
		ReleaseID       uint       `gorm:"column:release_id"`
		Title           string     `gorm:"column:title"`
		Environment     string     `gorm:"column:environment"`
		ApplicantID     uint       `gorm:"column:applicant_id"`
		ApplicantName   string     `gorm:"column:applicant_name"`
		Status          int        `gorm:"column:status"`
		Branch          string     `gorm:"column:branch"`
		CommitID        string     `gorm:"column:commit_id"`
		BuildNumber     int        `gorm:"column:build_number"`
		JenkinsJobURL   string     `gorm:"column:jenkins_job_url"`
		StartTime       *time.Time `gorm:"column:start_time"`
		EndTime         *time.Time `gorm:"column:end_time"`
		Duration        int        `gorm:"column:duration"`
		ImpactFeature   string     `gorm:"column:impact_feature"`
		FunctionModule  string     `gorm:"column:function_module"`
		CreatedAt       time.Time  `gorm:"column:created_at"`
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("item.created_at DESC").
		Limit(req.PageSize).
		Offset(offset).
		Scan(&records).Error; err != nil {
		result.Failed(c, 500, "查询发布历史失败: "+err.Error())
		return
	}

	// 5. 状态文本映射
	statusTextMap := map[int]string{
		1: "未发布",
		2: "发布中",
		3: "发布成功",
		4: "发布失败",
	}

	// 6. 构建响应
	historyList := []model.ApplicationReleaseHistory{}
	for _, record := range records {
		historyList = append(historyList, model.ApplicationReleaseHistory{
			ReleaseID:      record.ReleaseID,
			ItemID:         record.ItemID,
			Title:          record.Title,
			Environment:    record.Environment,
			ApplicantID:    record.ApplicantID,
			ApplicantName:  record.ApplicantName,
			Status:         record.Status,
			StatusText:     statusTextMap[record.Status],
			Branch:         record.Branch,
			CommitID:       record.CommitID,
			BuildNumber:    record.BuildNumber,
			JenkinsJobURL:  record.JenkinsJobURL,
			StartTime:      record.StartTime,
			EndTime:        record.EndTime,
			Duration:       record.Duration,
			ImpactFeature:  record.ImpactFeature,
			FunctionModule: record.FunctionModule,
			CreatedAt:      record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response := model.ApplicationReleaseHistoryResponse{
		AppID:   app.ID,
		AppName: app.Name,
		AppCode: app.Code,
		Total:   int(total),
		List:    historyList,
	}

	result.Success(c, response)
}

// sendDeploymentNotification 发送发布完成通知
func (s *ApplicationService) sendDeploymentNotification(deploymentID uint, finalStatus int, allSuccess bool) {
	fmt.Printf("DEBUG: 开始发送发布通知, deploymentID=%d, finalStatus=%d, allSuccess=%v\n", deploymentID, finalStatus, allSuccess)

	// 获取发布信息
	var deployment model.QuickDeployment
	if err := s.db.Preload("Tasks").First(&deployment, deploymentID).Error; err != nil {
		fmt.Printf("ERROR: 获取发布信息失败: %v\n", err)
		return
	}

	fmt.Printf("DEBUG: 发布信息 - Notice=%v, RobotID=%v\n", deployment.Notice, deployment.RobotID)

	// 检查是否开启通知
	if !deployment.Notice || deployment.RobotID == nil {
		fmt.Printf("DEBUG: 通知未开启或机器人ID为空，跳过发送\n")
		return
	}

	// 获取机器人信息
	robot, err := s.robotDao.GetByID(*deployment.RobotID)
	if err != nil || robot == nil {
		fmt.Printf("ERROR: 获取机器人信息失败: %v\n", err)
		return
	}

	fmt.Printf("DEBUG: 获取到机器人 - ID=%d, Name=%s, Type=%s\n", robot.ID, robot.Name, robot.Type)

	// 构建通知内容
	statusText := "成功"
	level := "info"
	if !allSuccess {
		statusText = "失败"
		level = "error"
	}

	// 收集所有任务名称
	var appNames []string
	for _, task := range deployment.Tasks {
		appNames = append(appNames, task.AppName)
	}
	appNamesStr := strings.Join(appNames, ", ")

	// 格式化时间
	createdAtStr := deployment.CreatedAt.Format("2006-01-02 15:04:05")
	durationStr := fmt.Sprintf("%d秒", deployment.Duration)

	// 构建通知消息
	title := deployment.Title
	content := fmt.Sprintf(
		"发布标题: %s\n发布环境: %s\n服务名称: %s\n发布人员: %s\n发布时间: %s\n运行时长: %s\n发布状态: %s",
		deployment.Title,
		deployment.Tasks[0].Environment, // 使用第一个任务的环境
		appNamesStr,
		deployment.CreatorName,
		createdAtStr,
		durationStr,
		statusText,
	)

	fmt.Printf("DEBUG: 构建通知消息 - Title=%s, Level=%s\n", title, level)

	// 创建webhook请求
	webhookReq := &monitormodel.WebhookRequest{
		Source:         "deployment",
		Title:          title,
		Content:        content,
		Level:          level,
		NotifyRobotIDs: []uint{*deployment.RobotID},
	}

	// 调用webhook服务发送通知
	s.sendWebhookNotification(webhookReq)
}

// sendWebhookNotification 发送webhook通知
func (s *ApplicationService) sendWebhookNotification(req *monitormodel.WebhookRequest) {
	fmt.Printf("DEBUG: 开始发送webhook通知, Source=%s, RobotIDs=%v\n", req.Source, req.NotifyRobotIDs)

	// 创建WebhookLog记录
	webhookLog := &monitormodel.WebhookLog{
		Source:         req.Source,
		Title:          req.Title,
		Content:        req.Content,
		Level:          req.Level,
		NotifyRobotIDs: fmt.Sprintf("%v", req.NotifyRobotIDs),
		Status:         "success",
	}

	// 保存WebhookLog
	if err := s.webhookDao.CreateLog(webhookLog); err != nil {
		fmt.Printf("ERROR: 保存WebhookLog失败: %v\n", err)
		return
	}

	fmt.Printf("DEBUG: WebhookLog已保存, ID=%d\n", webhookLog.ID)

	// 获取目标机器人
	var robots []monitormodel.NotifyRobot
	if len(req.NotifyRobotIDs) > 0 {
		s.db.Where("id IN ? AND status = ?", req.NotifyRobotIDs, 1).Find(&robots)
	} else {
		s.db.Where("status = ?", 1).Find(&robots)
	}

	fmt.Printf("DEBUG: 获取到%d个机器人\n", len(robots))

	// 为每个机器人发送通知
	successCount := 0
	failedCount := 0
	for _, robot := range robots {
		fmt.Printf("DEBUG: 开始发送通知到机器人 - ID=%d, Type=%s\n", robot.ID, robot.Type)

		notifyLog := &monitormodel.WebhookNotifyLog{
			WebhookLogID: webhookLog.ID,
			RobotID:      robot.ID,
			RobotName:    robot.Name,
			RobotType:    robot.Type,
			Status:       "success",
		}

		// 根据机器人类型发送通知
		if s.sendNotificationByRobotType(&robot, req.Title, req.Content, req.Level) {
			successCount++
			notifyLog.Status = "success"
			fmt.Printf("DEBUG: 通知发送成功 - RobotID=%d\n", robot.ID)
		} else {
			failedCount++
			notifyLog.Status = "failed"
			notifyLog.ErrorMsg = "发送失败"
			fmt.Printf("DEBUG: 通知发送失败 - RobotID=%d\n", robot.ID)
		}

		// 保存通知日志
		s.webhookDao.CreateNotifyLog(notifyLog)
	}

	// 更新WebhookLog统计
	s.db.Model(webhookLog).Updates(map[string]interface{}{
		"notify_count":  len(robots),
		"success_count": successCount,
		"failed_count":  failedCount,
	})

	fmt.Printf("DEBUG: webhook通知完成 - 总数=%d, 成功=%d, 失败=%d\n", len(robots), successCount, failedCount)
}

// sendNotificationByRobotType 根据机器人类型发送通知
func (s *ApplicationService) sendNotificationByRobotType(robot *monitormodel.NotifyRobot, title, content, level string) bool {
	switch robot.Type {
	case "dingtalk":
		return s.sendDingTalkNotification(robot, title, content, level)
	case "feishu":
		return s.sendFeishuNotification(robot, title, content, level)
	case "wechat":
		return s.sendWechatNotification(robot, title, content, level)
	case "webhook":
		return s.sendCustomWebhookNotification(robot, title, content, level)
	case "teams":
		return s.sendTeamsNotification(robot, title, content, level)
	default:
		return false
	}
}

// sendDingTalkNotification 发送钉钉通知
func (s *ApplicationService) sendDingTalkNotification(robot *monitormodel.NotifyRobot, title, content, level string) bool {
	if robot.Webhook == "" {
		fmt.Printf("ERROR: 钉钉机器人webhook为空\n")
		return false
	}

	fmt.Printf("DEBUG: 发送钉钉通知 - Webhook=%s\n", robot.Webhook)

	// 构建钉钉消息
	message := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("%s\n%s", title, content),
		},
	}

	jsonData, _ := json.Marshal(message)
	resp, err := http.Post(robot.Webhook, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("ERROR: 钉钉通知发送失败: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	success := resp.StatusCode == http.StatusOK
	fmt.Printf("DEBUG: 钉钉通知响应 - StatusCode=%d, Success=%v\n", resp.StatusCode, success)
	return success
}

// sendFeishuNotification 发送飞书通知
func (s *ApplicationService) sendFeishuNotification(robot *monitormodel.NotifyRobot, title, content, level string) bool {
	if robot.Webhook == "" {
		return false
	}

	// 构建飞书消息
	message := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": fmt.Sprintf("%s\n%s", title, content),
		},
	}

	jsonData, _ := json.Marshal(message)
	resp, err := http.Post(robot.Webhook, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// sendWechatNotification 发送企业微信通知
func (s *ApplicationService) sendWechatNotification(robot *monitormodel.NotifyRobot, title, content, level string) bool {
	if robot.Webhook == "" {
		return false
	}

	// 构建企业微信消息
	message := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("%s\n%s", title, content),
		},
	}

	jsonData, _ := json.Marshal(message)
	resp, err := http.Post(robot.Webhook, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// sendCustomWebhookNotification 发送自定义webhook通知
func (s *ApplicationService) sendCustomWebhookNotification(robot *monitormodel.NotifyRobot, title, content, level string) bool {
	if robot.Webhook == "" {
		return false
	}

	// 构建自定义消息
	message := map[string]interface{}{
		"title":   title,
		"content": content,
		"level":   level,
	}

	jsonData, _ := json.Marshal(message)
	resp, err := http.Post(robot.Webhook, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// sendTeamsNotification 发送Teams通知
func (s *ApplicationService) sendTeamsNotification(robot *monitormodel.NotifyRobot, title, content, level string) bool {
	if robot.Webhook == "" {
		return false
	}

	// 构建Teams MessageCard
	message := map[string]interface{}{
		"@type":       "MessageCard",
		"@context":    "https://schema.org/extensions",
		"summary":     title,
		"themeColor":  "0078D4",
		"sections": []map[string]interface{}{
			{
				"activityTitle": title,
				"text":          content,
			},
		},
	}

	jsonData, _ := json.Marshal(message)
	resp, err := http.Post(robot.Webhook, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

