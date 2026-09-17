package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"dodevops-api/api/configcenter/model"
	"dodevops-api/api/configcenter/service"
	sysModel "dodevops-api/api/system/model"
	sysService "dodevops-api/api/system/service"
	"dodevops-api/common/util"
	"github.com/robfig/cron/v3"
)

// SyncScheduler 定时同步调度器
type SyncScheduler struct {
	cron                *cron.Cron
	syncScheduleService *service.SyncScheduleService
	keyManageService    *service.KeyManageService
	activityLogService  *sysService.SysActivityLogService
	scheduleJobs        map[uint]cron.EntryID // 记录配置ID到任务ID的映射
	mutex               sync.RWMutex
	ctx                 context.Context
	cancel              context.CancelFunc
}

// NewSyncScheduler 创建新的定时同步调度器
func NewSyncScheduler() *SyncScheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &SyncScheduler{
		cron:                cron.New(),
		syncScheduleService: service.NewSyncScheduleService(),
		keyManageService:    service.NewKeyManageService(),
		activityLogService:  sysService.NewSysActivityLogService(),
		scheduleJobs:        make(map[uint]cron.EntryID),
		ctx:                 ctx,
		cancel:              cancel,
	}
}

// Start 启动调度器
func (s *SyncScheduler) Start() error {
	log.Println("启动定时同步调度器...")

	// 加载所有启用的定时同步配置
	if err := s.LoadActiveSchedules(); err != nil {
		return fmt.Errorf("failed to load active schedules: %v", err)
	}

	// 启动cron调度器
	s.cron.Start()

	log.Println("定时同步调度器启动成功")
	return nil
}

// Stop 停止调度器
func (s *SyncScheduler) Stop() {
	log.Println("停止定时同步调度器...")

	s.cron.Stop()
	s.cancel()

	log.Println("定时同步调度器已停止")
}

// LoadActiveSchedules 加载所有启用的定时同步配置
func (s *SyncScheduler) LoadActiveSchedules() error {
	schedules, err := s.syncScheduleService.GetActiveSchedules()
	if err != nil {
		return fmt.Errorf("failed to get active schedules: %v", err)
	}

	log.Printf("加载到 %d 个启用的定时同步配置", len(schedules))

	for _, schedule := range schedules {
		if err := s.AddSchedule(&schedule); err != nil {
			log.Printf("Failed to add schedule %d: %v", schedule.ID, err)
		}
	}

	return nil
}

// AddSchedule 添加定时同步配置到调度器
func (s *SyncScheduler) AddSchedule(schedule *model.SyncSchedule) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 如果已经存在，先移除
	if entryID, exists := s.scheduleJobs[schedule.ID]; exists {
		s.cron.Remove(entryID)
		delete(s.scheduleJobs, schedule.ID)
	}

	// 添加新的定时任务（使用标准5位cron格式）
	entryID, err := s.cron.AddFunc(schedule.CronExpr, func() {
		s.executeSyncJob(schedule)
	})

	if err != nil {
		return fmt.Errorf("failed to add cron job for schedule %d: %v", schedule.ID, err)
	}

	s.scheduleJobs[schedule.ID] = entryID
	log.Printf("添加定时同步配置: ID=%d, Name=%s, CronExpr=%s", schedule.ID, schedule.Name, schedule.CronExpr)

	return nil
}

// RemoveSchedule 从调度器中移除定时同步配置
func (s *SyncScheduler) RemoveSchedule(scheduleID uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if entryID, exists := s.scheduleJobs[scheduleID]; exists {
		s.cron.Remove(entryID)
		delete(s.scheduleJobs, scheduleID)
		log.Printf("移除定时同步配置: ID=%d", scheduleID)
	}
}

// UpdateSchedule 更新定时同步配置
func (s *SyncScheduler) UpdateSchedule(schedule *model.SyncSchedule) error {
	// 先移除旧的配置
	s.RemoveSchedule(schedule.ID)

	// 如果配置是启用状态，则重新添加
	if schedule.Status == 1 {
		return s.AddSchedule(schedule)
	}

	return nil
}

// executeSyncJob 执行同步任务
func (s *SyncScheduler) executeSyncJob(schedule *model.SyncSchedule) {
	log.Printf("开始执行定时同步任务: %s (ID: %d)", schedule.Name, schedule.ID)
	startTime := time.Now()

	// 解析密钥ID列表
	var keyIDs []int
	if err := json.Unmarshal([]byte(schedule.KeyTypes), &keyIDs); err != nil {
		log.Printf("解析密钥ID列表失败: %v", err)
		return
	}

	// 收集同步日志和统计
	var syncLogBuilder string
	syncLogBuilder += fmt.Sprintf("[%s] 开始同步\n", startTime.Format("2006-01-02 15:04:05"))

	totalNewAdded := 0
	totalFetched := 0
	totalFailed := 0
	keyCount := 0

	// 按云厂商类型分组，同时缓存密钥对象
	type KeyGroup struct {
		KeyType int
		Keys    []*model.KeyManage
	}
	cloudTypeGroups := make(map[int][]*model.KeyManage) // keyType -> KeyManages

	for _, keyID := range keyIDs {
		keyManage, err := s.keyManageService.GetByID(uint(keyID))
		if err != nil {
			log.Printf("获取密钥 %d 失败: %v", keyID, err)
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
			keyLog, result := s.syncSingleKeyWithStats(keyManage, keyType)
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
	//syncLogBuilder += fmt.Sprintf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	syncLogBuilder += fmt.Sprintf("本次同步 %d 台新增主机，%d 个密钥，共计耗时 %v\n", totalNewAdded, keyCount, duration)
	syncLogBuilder += fmt.Sprintf("[%s] 同步完成\n", time.Now().Format("2006-01-02 15:04:05"))

	// 更新上次执行时间和同步日志
	if err := s.syncScheduleService.UpdateLastRunTimeAndLog(schedule.ID, startTime, syncLogBuilder); err != nil {
		log.Printf("更新执行时间和日志失败: %v", err)
	}

	// 记录到系统动态日志
	activityLog := &sysModel.SysActivityLog{
		ActivityType: sysModel.ActivityTypeSyncKey,
		Title:        schedule.Name, // 直接使用配置名称
		Content:      syncLogBuilder,
		Status:       1, // 成功
		RelatedID:    schedule.ID,
		Summary:      fmt.Sprintf("同步%d个密钥，共计新增%d台主机，耗时%.1f秒", keyCount, totalNewAdded, duration.Seconds()),
		Duration:     int64(duration.Seconds()), // 转换为秒
		CreateTime:   util.HTime{Time: time.Now()},
	}
	if err := s.activityLogService.Create(activityLog); err != nil {
		log.Printf("记录系统动态日志失败: %v", err)
	}

	log.Printf("定时同步任务完成: %s (ID: %d), 耗时: %v, 新增: %d台", schedule.Name, schedule.ID, duration, totalNewAdded)
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

// syncCloudHostsByType 根据云厂商类型同步云主机
func (s *SyncScheduler) syncCloudHostsByType(keyType int) ([]string, error) {
	// 获取该类型的所有密钥
	keyManages, err := s.keyManageService.GetByType(keyType)
	if err != nil {
		return nil, fmt.Errorf("failed to get keys for type %d: %v", keyType, err)
	}

	if len(keyManages) == 0 {
		log.Printf("云厂商类型 %d 没有配置密钥，跳过同步", keyType)
		return []string{"未配置密钥"}, nil
	}

	var logs []string
	// 为每个密钥执行同步
	for _, keyManage := range keyManages {
		syncLog, err := s.syncSingleKey(&keyManage, keyType)
		if err != nil {
			log.Printf("同步密钥 %d 失败: %v", keyManage.ID, err)
			logs = append(logs, syncLog)
		} else {
			log.Printf("成功同步密钥 %d", keyManage.ID)
			logs = append(logs, syncLog)
		}
	}

	return logs, nil
}

// syncSingleKeyWithStats 同步单个密钥的云主机并返回统计信息
func (s *SyncScheduler) syncSingleKeyWithStats(keyManage *model.KeyManage, keyType int) (string, *service.SyncResult) {
	// 固定使用默认分组ID为1
	groupID := uint(1)
	// 固定使用全区域扫描
	region := "all"

	// 根据keyType调用不同的后台同步方法
	switch keyType {
	case 1: // 阿里云
		result, err := s.keyManageService.SyncAliyunHostsBackground(keyManage.ID, groupID, region)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), result
		}
		// 构建区域统计信息
		regionInfo := ""
		if len(result.RegionStats) > 0 {
			regionInfo = " ["
			for region, count := range result.RegionStats {
				regionInfo += fmt.Sprintf("%s:%d台 ", region, count)
			}
			regionInfo += "]"
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台%s",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount, regionInfo), result
	case 2: // 腾讯云
		result, err := s.keyManageService.SyncTencentHostsBackground(keyManage.ID, groupID)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), result
		}
		// 构建区域统计信息
		regionInfo := ""
		if len(result.RegionStats) > 0 {
			regionInfo = " ["
			for region, count := range result.RegionStats {
				regionInfo += fmt.Sprintf("%s:%d台 ", region, count)
			}
			regionInfo += "]"
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台%s",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount, regionInfo), result
	case 3: // 百度云
		result, err := s.keyManageService.SyncBaiduHostsBackground(keyManage.ID, groupID)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), result
		}
		// 构建区域统计信息
		regionInfo := ""
		if len(result.RegionStats) > 0 {
			regionInfo = " ["
			for region, count := range result.RegionStats {
				regionInfo += fmt.Sprintf("%s:%d台 ", region, count)
			}
			regionInfo += "]"
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台%s",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount, regionInfo), result
	case 4: // 华为云
		result, err := s.keyManageService.SyncHuaweiHostsBackground(keyManage.ID, groupID)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), result
		}
		// 构建区域统计信息
		regionInfo := ""
		if len(result.RegionStats) > 0 {
			regionInfo = " ["
			for region, count := range result.RegionStats {
				regionInfo += fmt.Sprintf("%s:%d台 ", region, count)
			}
			regionInfo += "]"
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台%s",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount, regionInfo), result
	case 5: // AWS云
		result, err := s.keyManageService.SyncAWSHostsBackground(keyManage.ID, groupID)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), result
		}
		// 构建区域统计信息
		regionInfo := ""
		if len(result.RegionStats) > 0 {
			regionInfo = " ["
			for region, count := range result.RegionStats {
				regionInfo += fmt.Sprintf("%s:%d台 ", region, count)
			}
			regionInfo += "]"
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台%s",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount, regionInfo), result
	default:
		return fmt.Sprintf("  - 未知云厂商(%d): 不支持", keyType), nil
	}
}

// syncSingleKey 同步单个密钥的云主机
func (s *SyncScheduler) syncSingleKey(keyManage *model.KeyManage, keyType int) (string, error) {
	// 固定使用默认分组ID为1
	groupID := uint(1)
	// 固定使用全区域扫描
	region := "all"

	// 根据keyType调用不同的后台同步方法
	switch keyType {
	case 1: // 阿里云
		result, err := s.keyManageService.SyncAliyunHostsBackground(keyManage.ID, groupID, region)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), err
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount), nil
	case 2: // 腾讯云
		result, err := s.keyManageService.SyncTencentHostsBackground(keyManage.ID, groupID)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), err
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount), nil
	case 3: // 百度云
		log.Printf("百度云同步暂未实现")
		return "  - 百度云: 暂未实现", nil
	case 4: // 华为云
		log.Printf("华为云同步暂未实现")
		return "  - 华为云: 暂未实现", nil
	case 5: // AWS云
		result, err := s.keyManageService.SyncAWSHostsBackground(keyManage.ID, groupID)
		if err != nil {
			return fmt.Sprintf("  - %s[%s]: 同步失败 - %s", result.KeyName, result.CloudType, result.Error), err
		}
		return fmt.Sprintf("  - %s[%s]: 获取%d台，新增%d台，跳过%d台，失败%d台",
			result.KeyName, result.CloudType, result.TotalFetched, result.NewAdded, result.SkippedCount, result.FailedCount), nil
	default:
		return fmt.Sprintf("  - 未知云厂商(%d): 不支持", keyType), fmt.Errorf("unsupported cloud type: %d", keyType)
	}
}

// GetJobStats 获取调度器状态信息
func (s *SyncScheduler) GetJobStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return map[string]interface{}{
		"total_jobs":   len(s.scheduleJobs),
		"cron_entries": len(s.cron.Entries()),
		"active_jobs":  s.scheduleJobs,
	}
}