package controller

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"dodevops-api/api/configcenter/model"
	"dodevops-api/api/configcenter/service"
	sysModel "dodevops-api/api/system/model"
	sysService "dodevops-api/api/system/service"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/scheduler"
	"github.com/gin-gonic/gin"
)

type SyncScheduleController struct {
	service *service.SyncScheduleService
}

func NewSyncScheduleController() *SyncScheduleController {
	return &SyncScheduleController{
		service: service.NewSyncScheduleService(),
	}
}

// Create 创建定时同步配置
// @Summary 创建定时同步配置
// @Description 创建定时同步配置
// @Tags Config配置中心
// @Param syncSchedule body model.CreateSyncScheduleDto true "定时同步配置信息"
// @Success 200 {object} result.Result{data=model.SyncSchedule}
// @Router /api/v1/config/sync-schedule [post]
// @Security ApiKeyAuth
func (c *SyncScheduleController) Create(ctx *gin.Context) {
	var dto model.CreateSyncScheduleDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	// 将DTO转换为实体
	syncSchedule := model.SyncSchedule{
		Name:     dto.Name,
		CronExpr: dto.CronExpr,
		KeyTypes: dto.KeyTypes,
		Status:   dto.Status,
		Remark:   dto.Remark,
	}

	if err := c.service.Create(&syncSchedule); err != nil {
		result.Failed(ctx, 500, "创建失败: "+err.Error())
		return
	}

	// 如果配置是启用状态，则添加到调度器
	if syncSchedule.Status == 1 {
		if err := scheduler.GetManager().AddSyncSchedule(&syncSchedule); err != nil {
			// 调度器添加失败不影响数据库创建，只记录日志
			// log.Printf("Failed to add schedule to scheduler: %v", err)
		}
	}

	result.Success(ctx, syncSchedule)
}

// Update 更新定时同步配置
// @Summary 更新定时同步配置
// @Description 更新定时同步配置
// @Tags Config配置中心
// @Param syncSchedule body model.UpdateSyncScheduleDto true "定时同步配置信息"
// @Success 200 {object} result.Result{data=model.SyncSchedule}
// @Router /api/v1/config/sync-schedule [put]
// @Security ApiKeyAuth
func (c *SyncScheduleController) Update(ctx *gin.Context) {
	var dto model.UpdateSyncScheduleDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	// 将DTO转换为实体
	syncSchedule := model.SyncSchedule{
		ID:       dto.ID,
		Name:     dto.Name,
		CronExpr: dto.CronExpr,
		KeyTypes: dto.KeyTypes,
		Status:   dto.Status,
		Remark:   dto.Remark,
	}

	if err := c.service.Update(&syncSchedule); err != nil {
		result.Failed(ctx, 500, "更新失败: "+err.Error())
		return
	}

	// 更新调度器中的配置
	if err := scheduler.GetManager().UpdateSyncSchedule(&syncSchedule); err != nil {
		// 调度器更新失败不影响数据库更新，只记录日志
		// log.Printf("Failed to update schedule in scheduler: %v", err)
	}

	result.Success(ctx, syncSchedule)
}

// Delete 删除定时同步配置
// @Summary 删除定时同步配置
// @Description 删除定时同步配置
// @Tags Config配置中心
// @Param id query uint true "配置ID"
// @Success 200 {object} result.Result
// @Router /api/v1/config/sync-schedule [delete]
// @Security ApiKeyAuth
func (c *SyncScheduleController) Delete(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	if err := c.service.Delete(req.ID); err != nil {
		result.Failed(ctx, 500, "删除失败: "+err.Error())
		return
	}

	// 从调度器中移除配置
	scheduler.GetManager().RemoveSyncSchedule(req.ID)

	result.Success(ctx, nil)
}

// GetByID 根据ID获取定时同步配置详情
// @Summary 根据ID获取定时同步配置详情
// @Description 根据ID获取定时同步配置详情
// @Tags Config配置中心
// @Param id query uint true "配置ID"
// @Success 200 {object} result.Result{data=model.SyncSchedule}
// @Router /api/v1/config/sync-schedule [get]
// @Security ApiKeyAuth
func (c *SyncScheduleController) GetByID(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	syncSchedule, err := c.service.GetByID(req.ID)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, syncSchedule)
}

// List 获取定时同步配置列表（分页）
// @Summary 获取定时同步配置列表（分页）
// @Description 获取定时同步配置列表（分页）
// @Tags Config配置中心
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Success 200 {object} result.Result{data=result.PageResult{list=[]model.SyncSchedule}}
// @Router /api/v1/config/sync-schedule/list [get]
// @Security ApiKeyAuth
func (c *SyncScheduleController) List(ctx *gin.Context) {
	var params PageParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	syncSchedules, total, err := c.service.ListWithPage(params.Page, params.PageSize)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	pageResult := result.PageResult{
		List:     syncSchedules,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
	result.Success(ctx, pageResult)
}

// ToggleStatus 切换配置状态（启用/禁用）
// @Summary 切换配置状态（启用/禁用）
// @Description 切换配置状态（启用/禁用）
// @Tags Config配置中心
// @Param request body object{id=uint,status=int} true "配置ID和状态"
// @Success 200 {object} result.Result
// @Router /api/v1/config/sync-schedule/toggle-status [post]
// @Security ApiKeyAuth
func (c *SyncScheduleController) ToggleStatus(ctx *gin.Context) {
	var req struct {
		ID     uint `json:"id" binding:"required"`
		Status int  `json:"status" binding:"min=0,max=1"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	if err := c.service.ToggleStatus(req.ID, req.Status); err != nil {
		result.Failed(ctx, 500, "操作失败: "+err.Error())
		return
	}

	// 根据状态更新调度器
	if req.Status == 1 {
		// 启用：获取配置并添加到调度器
		if schedule, err := c.service.GetByID(req.ID); err == nil {
			scheduler.GetManager().AddSyncSchedule(schedule)
		}
	} else {
		// 禁用：从调度器中移除
		scheduler.GetManager().RemoveSyncSchedule(req.ID)
	}

	result.Success(ctx, nil)
}

// GetActiveSchedules 获取所有启用的定时同步配置
// @Summary 获取所有启用的定时同步配置
// @Description 获取所有启用的定时同步配置
// @Tags Config配置中心
// @Success 200 {object} result.Result{data=[]model.SyncSchedule}
// @Router /api/v1/config/sync-schedule/active [get]
// @Security ApiKeyAuth
func (c *SyncScheduleController) GetActiveSchedules(ctx *gin.Context) {
	schedules, err := c.service.GetActiveSchedules()
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, schedules)
}

// TriggerManualSync 手动触发定时同步（测试用）
// @Summary 手动触发定时同步（测试用）
// @Description 手动触发定时同步（测试用）
// @Tags Config配置中心
// @Param id query uint true "配置ID"
// @Success 200 {object} result.Result
// @Router /api/v1/config/sync-schedule/trigger [post]
// @Security ApiKeyAuth
// TriggerManualSync 手动触发同步任务
// @Summary 手动触发同步任务
// @Description 手动触发同步任务
// @Tags Config配置中心
// @Param id query uint true "配置ID"
// @Success 200 {object} result.Result
// @Router /api/v1/config/sync-schedule/trigger [post]
// @Security ApiKeyAuth
func (c *SyncScheduleController) TriggerManualSync(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		result.Failed(ctx, 400, "配置ID不能为空")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(ctx, 400, "配置ID格式错误")
		return
	}

	// 获取配置信息
	schedule, err := c.service.GetByID(uint(id))
	if err != nil {
		result.Failed(ctx, 500, "查询配置失败: "+err.Error())
		return
	}

	if schedule.Status != 1 {
		result.Failed(ctx, 400, "配置未启用，无法执行同步")
		return
	}

	// 先返回响应，避免超时
	result.Success(ctx, gin.H{
		"message":  "同步任务已触发，正在后台执行",
		"schedule": schedule.Name,
	})

	// 异步执行同步任务
	go c.executeManualSync(schedule)
}

// executeManualSync 执行手动同步任务
func (c *SyncScheduleController) executeManualSync(schedule *model.SyncSchedule) {
	startTime := time.Now()
	keyManageService := service.NewKeyManageService()
	activityLogService := sysService.NewSysActivityLogService()

	// 解析密钥ID列表
	var keyIDs []int
	if err := json.Unmarshal([]byte(schedule.KeyTypes), &keyIDs); err != nil {
		fmt.Printf("解析密钥ID列表失败: %v\n", err)
		return
	}

	// 收集同步日志和统计
	var syncLogBuilder string
	syncLogBuilder += fmt.Sprintf("[%s] 手动触发同步\n", startTime.Format("2006-01-02 15:04:05"))

	totalNewAdded := 0
	totalFetched := 0
	totalFailed := 0
	keyCount := 0

	// 按云厂商类型分组，同时缓存密钥对象
	cloudTypeGroups := make(map[int][]*model.KeyManage) // keyType -> KeyManages

	for _, keyID := range keyIDs {
		keyManage, err := keyManageService.GetByID(uint(keyID))
		if err != nil {
			fmt.Printf("获取密钥 %d 失败: %v\n", keyID, err)
			continue
		}
		cloudTypeGroups[keyManage.KeyType] = append(cloudTypeGroups[keyManage.KeyType], keyManage)
	}

	// 为每种云厂商类型执行同步
	for keyType, keyManages := range cloudTypeGroups {
		cloudTypeName := getCloudTypeName(keyType)
		syncLogBuilder += fmt.Sprintf("\n%s:\n", cloudTypeName)

		for _, keyManage := range keyManages {
			keyCount++
			keyLog, result := syncSingleKeyWithStats(keyManageService, keyManage, keyType)
			syncLogBuilder += keyLog + "\n"

			if result != nil {
				totalNewAdded += result.NewAdded
				totalFetched += result.TotalFetched
				totalFailed += result.FailedCount
			}
		}
	}

	duration := time.Since(startTime)

	// 添加总体统计
	syncLogBuilder += fmt.Sprintf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	syncLogBuilder += fmt.Sprintf("本次同步 %d 台新增主机，%d 个密钥，共计耗时 %v\n", totalNewAdded, keyCount, duration)
	syncLogBuilder += fmt.Sprintf("[%s] 同步完成\n", time.Now().Format("2006-01-02 15:04:05"))

	// 更新同步日志到配置表
	c.service.UpdateLastRunTimeAndLog(schedule.ID, startTime, syncLogBuilder)

	// 记录到系统动态日志
	activityLog := &sysModel.SysActivityLog{
		ActivityType: sysModel.ActivityTypeSyncKey,
		Title:        fmt.Sprintf("手动同步任务: %s", schedule.Name),
		Content:      syncLogBuilder,
		Status:       1,
		RelatedID:    schedule.ID,
		Summary:      fmt.Sprintf("手动触发同步%d个密钥，新增%d台主机", keyCount, totalNewAdded),
		Duration:     duration.Milliseconds(),
		CreateTime:   util.HTime{Time: time.Now()},
	}
	activityLogService.Create(activityLog)

	fmt.Printf("手动同步任务完成: %s (ID: %d), 耗时: %v, 新增: %d台\n", schedule.Name, schedule.ID, duration, totalNewAdded)
}

// syncSingleKeyWithStats 同步单个密钥并返回统计信息
func syncSingleKeyWithStats(keyManageService *service.KeyManageService, keyManage *model.KeyManage, keyType int) (string, *service.SyncResult) {
	groupID := uint(1)
	region := "all"

	switch keyType {
	case 1: // 阿里云
		result, err := keyManageService.SyncAliyunHostsBackground(keyManage.ID, groupID, region)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), result
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount), result
	case 2: // 腾讯云
		result, err := keyManageService.SyncTencentHostsBackground(keyManage.ID, groupID)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), result
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount), result
	case 3: // 百度云
		result, err := keyManageService.SyncBaiduHostsBackground(keyManage.ID, groupID)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), result
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount), result
	case 4: // 华为云
		result, err := keyManageService.SyncHuaweiHostsBackground(keyManage.ID, groupID)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), result
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount), result
	case 5: // AWS云
		result, err := keyManageService.SyncAWSHostsBackground(keyManage.ID, groupID)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), result
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount), result
	default:
		return fmt.Sprintf("  - 未知云厂商(%d): 不支持", keyType), nil
	}
}

// syncSingleKey 同步单个密钥的云主机
func syncSingleKey(keyManageService *service.KeyManageService, keyManage *model.KeyManage, keyType int) string {
	groupID := uint(1)
	region := "all"

	switch keyType {
	case 1: // 阿里云
		result, err := keyManageService.SyncAliyunHostsBackground(keyManage.ID, groupID, region)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error)
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount)
	case 2: // 腾讯云
		result, err := keyManageService.SyncTencentHostsBackground(keyManage.ID, groupID)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error)
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount)
	case 3: // 百度云
		return "  - 百度云: 暂未实现"
	case 4: // 华为云
		return "  - 华为云: 暂未实现"
	case 5: // AWS云
		return "  - AWS云: 暂未实现"
	default:
		return fmt.Sprintf("  - 未知云厂商(%d): 不支持", keyType)
	}
}

// getCloudTypeName 获取云厂商类型名称
func getCloudTypeName(keyType int) string {
	switch keyType {
	case 1:
		return "阿里云"
	case 2:
		return "腾讯云"
	case 3:
		return "百度云"
	case 4:
		return "华为云"
	case 5:
		return "AWS云"
	default:
		return fmt.Sprintf("未知云厂商(%d)", keyType)
	}
}

// GetSchedulerStats 获取调度器状态信息
// @Summary 获取调度器状态信息
// @Description 获取调度器状态信息
// @Tags Config配置中心
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Router /api/v1/config/sync-schedule/scheduler-stats [get]
// @Security ApiKeyAuth
func (c *SyncScheduleController) GetSchedulerStats(ctx *gin.Context) {
	stats := scheduler.GetManager().GetSyncSchedulerStats()
	result.Success(ctx, stats)
}

// GetSyncLog 获取同步日志
// @Summary 获取同步日志
// @Description 获取同步日志
// @Tags Config配置中心
// @Param id query uint true "配置ID"
// @Success 200 {object} result.Result{data=object{syncLog=string,lastRunTime=string}}
// @Router /api/v1/config/sync-schedule/log [get]
// @Security ApiKeyAuth
func (c *SyncScheduleController) GetSyncLog(ctx *gin.Context) {
	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}

	syncSchedule, err := c.service.GetByID(req.ID)
	if err != nil {
		result.Failed(ctx, 500, "查询失败: "+err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"id":          syncSchedule.ID,
		"name":        syncSchedule.Name,
		"syncLog":     syncSchedule.SyncLog,
		"lastRunTime": syncSchedule.LastRunTime,
		"nextRunTime": syncSchedule.NextRunTime,
	})
}