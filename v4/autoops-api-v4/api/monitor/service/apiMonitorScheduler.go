package service

import (
	"log"
	"sync"
	"time"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/util"
)

// APIMonitorScheduler API监控调度器（全自动执行检测）
type APIMonitorScheduler struct {
	monitorDao    dao.APIMonitorDao
	resultDao     dao.APIMonitorResultDao
	monitorService APIMonitorService

	// 调度器状态
	running       bool
	stopChan      chan struct{}
	mu            sync.RWMutex

	// 下次执行时间映射 (monitorID -> nextCheckTime)
	nextCheckTimes map[uint]time.Time
}

// NewAPIMonitorScheduler 创建调度器
func NewAPIMonitorScheduler() *APIMonitorScheduler {
	return &APIMonitorScheduler{
		monitorDao:     dao.NewAPIMonitorDao(),
		resultDao:      dao.NewAPIMonitorResultDao(),
		monitorService: NewAPIMonitorService(),
		stopChan:       make(chan struct{}),
		nextCheckTimes: make(map[uint]time.Time),
	}
}

// Start 启动调度器
func (s *APIMonitorScheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		log.Println("⚠️  API监控调度器已经在运行中")
		return
	}
	s.running = true
	s.mu.Unlock()

	log.Println("🚀 API监控调度器启动中...")

	// 启动主循环（每10秒检查一次）
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.checkAndExecute()
			case <-s.stopChan:
				ticker.Stop()
				log.Println("⏹  API监控调度器已停止")
				return
			}
		}
	}()

	// 启动清理旧数据任务（每天凌晨2点执行）
	go s.startCleanupTask()

	// 启动父监控聚合数据更新任务（每分钟执行一次）
	go s.startAggregateUpdateTask()

	log.Println("✅ API监控调度器已启动")
}

// Stop 停止调度器
func (s *APIMonitorScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	close(s.stopChan)
}

// checkAndExecute 检查并执行需要检测的API
func (s *APIMonitorScheduler) checkAndExecute() {
	// 获取所有启用的API监控
	monitors, err := s.monitorDao.GetAllEnabled()
	if err != nil {
		log.Printf("❌ 获取API监控列表失败: %v", err)
		return
	}

	now := time.Now()
	var wg sync.WaitGroup

	for _, monitor := range monitors {
		// 检查是否到达执行时间
		nextCheckTime, exists := s.nextCheckTimes[monitor.ID]
		if !exists || now.After(nextCheckTime) {
			// 需要执行检测
			wg.Add(1)
			go func(m *model.APIMonitor) {
				defer wg.Done()
				s.executeMonitor(m)
			}(monitor)

			// 更新下次检测时间
			s.nextCheckTimes[monitor.ID] = now.Add(time.Duration(monitor.CheckInterval) * time.Second)
		}
	}

	wg.Wait()
}

// executeMonitor 执行单个API监控
func (s *APIMonitorScheduler) executeMonitor(monitor *model.APIMonitor) {
	log.Printf("🔍 开始检测API [监控:%s, URL:%s]", monitor.Name, monitor.URL)

	// 执行检测（支持重试）
	var checkResult *model.APIMonitorResult
	for i := 0; i <= monitor.RetryCount; i++ {
		checkResult = s.monitorService.ExecuteCheck(monitor)

		if checkResult.Success {
			break
		}

		// 失败后重试
		if i < monitor.RetryCount {
			log.Printf("⚠️  检测失败，%d秒后重试 [监控:%s, 第%d次重试]", 2, monitor.Name, i+1)
			time.Sleep(2 * time.Second)
		}
	}

	// 保存检测结果到数据库
	if err := s.resultDao.Create(checkResult); err != nil {
		log.Printf("❌ 保存检测结果失败: %v", err)
	}

	// 更新监控配置的最近检测结果
	now := time.Now()
	htime := util.HTime{Time: now}
	monitor.LastCheckTime = &htime
	monitor.LastResponseTime = checkResult.ResponseTime
	monitor.LastStatusCode = checkResult.StatusCode
	monitor.LastErrorMessage = checkResult.ErrorMessage

	if checkResult.Success {
		monitor.ContinuousFailures = 0
		monitor.HealthStatus = "healthy"
		log.Printf("✅ API检测成功 [监控:%s, 响应时间:%dms]", monitor.Name, checkResult.ResponseTime)
	} else {
		monitor.ContinuousFailures++
		monitor.HealthStatus = "unhealthy"
		log.Printf("❌ API检测失败 [监控:%s, 连续失败:%d次, 错误:%s]",
			monitor.Name, monitor.ContinuousFailures, checkResult.ErrorMessage)
	}

	// 更新到数据库
	if err := s.monitorDao.UpdateLastCheckResult(monitor); err != nil {
		log.Printf("❌ 更新检测结果失败: %v", err)
	}

	// 推送指标到Prometheus
	if err := s.monitorService.PushMetricToPrometheus(monitor, checkResult); err != nil {
		log.Printf("⚠️  推送指标到Prometheus失败: %v", err)
	}

	// 触发告警（如果启用且检测失败）
	if !checkResult.Success && monitor.AlertEnabled && monitor.PolicyID != nil {
		s.triggerAlert(monitor, checkResult)
	}
}

// triggerAlert 触发告警
func (s *APIMonitorScheduler) triggerAlert(monitor *model.APIMonitor, checkResult *model.APIMonitorResult) {
	// 判断是否达到告警阈值（例如：连续失败3次才告警）
	alertThreshold := 3
	if monitor.ContinuousFailures < alertThreshold {
		log.Printf("⏳ 未达到告警阈值 [监控:%s, 连续失败:%d/%d]", monitor.Name, monitor.ContinuousFailures, alertThreshold)
		return
	}

	log.Printf("🚨 触发告警 [监控:%s, 连续失败:%d次]", monitor.Name, monitor.ContinuousFailures)

	// TODO: 集成现有的告警通知系统
	// 可以调用 NotifyRobotService 发送告警通知
	// 或者创建 AlertStats 记录
}

// startCleanupTask 启动清理旧数据任务
func (s *APIMonitorScheduler) startCleanupTask() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 每天凌晨2点执行
			now := time.Now()
			if now.Hour() == 2 {
				log.Println("🧹 开始清理API监控旧数据...")
				// 保留最近30天的数据
				if err := s.resultDao.DeleteOldRecords(30); err != nil {
					log.Printf("❌ 清理旧数据失败: %v", err)
				} else {
					log.Println("✅ 清理旧数据完成")
				}
			}
		case <-s.stopChan:
			return
		}
	}
}

// GetStatus 获取调度器状态
func (s *APIMonitorScheduler) GetStatus() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"running":        s.running,
		"monitored_apis": len(s.nextCheckTimes),
	}
}

// startAggregateUpdateTask 启动父监控聚合数据更新任务（每分钟执行一次）
func (s *APIMonitorScheduler) startAggregateUpdateTask() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.monitorService.UpdateAllParentAggregateStats(); err != nil {
				log.Printf("❌ 更新父监控聚合数据失败: %v", err)
			}
		case <-s.stopChan:
			return
		}
	}
}
