package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/config"
	"dodevops-api/common/result"
	"dodevops-api/common/util"

	"github.com/gin-gonic/gin"
)

type APIMonitorService interface {
	// CRUD操作
	CreateMonitor(c *gin.Context, dto *model.CreateAPIMonitorDto)
	SimpleCreateMonitor(c *gin.Context, dto *model.SimpleCreateAPIMonitorDto)
	UpdateMonitor(c *gin.Context, dto *model.UpdateAPIMonitorDto)
	DeleteMonitor(c *gin.Context, id uint)
	GetMonitor(c *gin.Context, id uint)
	GetMonitorList(c *gin.Context, query *model.APIMonitorListQuery)
	UpdateMonitorStatus(c *gin.Context, id uint, status int)

	// 统计接口
	GetStats(c *gin.Context)
	GetGroupStats(c *gin.Context)
	GetServiceStats(c *gin.Context, groupID uint)

	// 测试接口
	TestAPI(c *gin.Context, dto *model.TestAPIDto)

	// 核心检测逻辑（供调度器调用）
	ExecuteCheck(monitor *model.APIMonitor) *model.APIMonitorResult
	PushMetricToPrometheus(monitor *model.APIMonitor, result *model.APIMonitorResult) error

	// 分层关系管理接口
	AddChildMonitor(c *gin.Context, parentID uint, childMonitorID uint)
	ConvertToChildMonitor(c *gin.Context, monitorID uint, parentID uint)
	GetMonitorHierarchy(c *gin.Context, monitorID uint, depth int)
	DeleteMonitorWithChildren(c *gin.Context, monitorID uint, operation string)
	CalculateParentHealthStatus(parentID uint) (string, error)
	UpdateAllParentAggregateStats() error
	RemoveParentRelationship(c *gin.Context, childMonitorID uint)
}

type APIMonitorServiceImpl struct {
	monitorDao dao.APIMonitorDao
	resultDao  dao.APIMonitorResultDao
}

func NewAPIMonitorService() APIMonitorService {
	return &APIMonitorServiceImpl{
		monitorDao: dao.NewAPIMonitorDao(),
		resultDao:  dao.NewAPIMonitorResultDao(),
	}
}

// CreateMonitor 创建API监控
func (s *APIMonitorServiceImpl) CreateMonitor(c *gin.Context, dto *model.CreateAPIMonitorDto) {
	monitor := &model.APIMonitor{
		Name:            dto.Name,
		Description:     dto.Description,
		GroupID:         dto.GroupID,
		GroupName:       dto.GroupName,
		ServiceName:     dto.ServiceName,
		URL:             dto.URL,
		Method:          dto.Method,
		Headers:         dto.Headers,
		Body:            dto.Body,
		Timeout:         dto.Timeout,
		ExpectedStatus:  dto.ExpectedStatus,
		ExpectedKeyword: dto.ExpectedKeyword,
		ValidateJSON:    dto.ValidateJSON,
		CheckInterval:   dto.CheckInterval,
		RetryCount:      dto.RetryCount,
		AlertEnabled:    dto.AlertEnabled,
		PolicyID:        dto.PolicyID,
		Status:          dto.Status,
		Tags:            dto.Tags,
		HealthStatus:    "unknown",
		ParentID:        dto.ParentID,        // 设置父监控ID
		HierarchyLevel:  dto.HierarchyLevel,  // 设置层级深度
	}

	// 设置默认值
	if monitor.Method == "" {
		monitor.Method = "GET"
	}
	if monitor.Timeout == 0 {
		monitor.Timeout = 5000
	}
	if monitor.ExpectedStatus == 0 {
		monitor.ExpectedStatus = 200
	}
	if monitor.CheckInterval == 0 {
		monitor.CheckInterval = 60
	}
	if monitor.RetryCount == 0 {
		monitor.RetryCount = 3
	}
	if monitor.Status == 0 {
		monitor.Status = 1
	}

	if err := s.monitorDao.Create(monitor); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建API监控失败: %v", err))
		return
	}

	log.Printf("✅ [CreateMonitor] 创建成功: ID=%d, Name=%s, ParentID=%v", monitor.ID, monitor.Name, monitor.ParentID)
	result.Success(c, monitor)
}

// SimpleCreateMonitor 简化版创建API监控
func (s *APIMonitorServiceImpl) SimpleCreateMonitor(c *gin.Context, dto *model.SimpleCreateAPIMonitorDto) {
	log.Printf("📝 [SimpleCreateMonitor] 开始创建监控: Name=%s, ParentID=%v", dto.Name, dto.ParentID)

	var businessFunction *model.APIBusinessFunction
	var err error
	var hierarchyLevel int = 0

	// 如果是创建子监控，逻辑不同
	if dto.ParentID != nil && *dto.ParentID > 0 {
		log.Printf("📝 [SimpleCreateMonitor] 创建子监控，ParentID=%d", *dto.ParentID)

		// 验证父监控是否存在
		parentMonitor, err := s.monitorDao.GetByID(*dto.ParentID)
		if err != nil {
			result.Failed(c, http.StatusBadRequest, "父监控不存在")
			return
		}

		// 计算层级深度
		hierarchyLevel = parentMonitor.HierarchyLevel + 1
		log.Printf("📝 [SimpleCreateMonitor] 子监控层级: %d", hierarchyLevel)

		// 子监控复用父监控的业务功能
		businessFunctionService := NewAPIBusinessFunctionService()

		// 优先通过ID获取业务功能
		if parentMonitor.BusinessFunctionID > 0 {
			businessFunction, err = businessFunctionService.GetByID(parentMonitor.BusinessFunctionID)
			if err != nil {
				log.Printf("⚠️ [SimpleCreateMonitor] 通过ID获取业务功能失败: %v", err)
			}
		}

		// 如果通过ID获取失败，则通过名称查找业务功能
		if businessFunction == nil {
			log.Printf("📝 [SimpleCreateMonitor] 尝试通过名称查找业务功能: %s", dto.BusinessFunctionName)
			businessFunction, err = businessFunctionService.GetByGroupIDAndName(dto.GroupID, dto.BusinessFunctionName)
			if err != nil {
				// 如果也找不到，再创建新的
				log.Printf("⚠️ [SimpleCreateMonitor] 业务功能不存在，创建新的")
				businessFunction, err = businessFunctionService.CreateIfNotExists(dto.GroupID, dto.BusinessFunctionName)
				if err != nil {
					result.Failed(c, http.StatusBadRequest, "业务功能创建失败: "+err.Error())
					return
				}
			} else {
				log.Printf("✅ [SimpleCreateMonitor] 找到已存在的业务功能: ID=%d, Name=%s", businessFunction.ID, businessFunction.Name)
			}
		}

	} else {
		log.Printf("📝 [SimpleCreateMonitor] 创建独立监控")

		// 1. 检查业务功能名称是否存在，不存在则创建（存在则返回错误）
		businessFunctionService := NewAPIBusinessFunctionService()
		businessFunction, err = businessFunctionService.CreateIfNotExists(dto.GroupID, dto.BusinessFunctionName)
		if err != nil {
			result.Failed(c, http.StatusBadRequest, "业务功能名称已存在: "+dto.BusinessFunctionName)
			return
		}

		// 2. 检查API监控名称在业务组内的唯一性
		exists, err := s.monitorDao.CheckNameExistsInGroup(dto.GroupID, dto.Name, 0)
		if err != nil {
			result.Failed(c, http.StatusInternalServerError, "检查名称失败")
			return
		}
		if exists {
			result.Failed(c, http.StatusBadRequest, "该业务组内已存在同名API监控: "+dto.Name)
			return
		}
	}

	// 3. 构建完整的CreateAPIMonitorDto，填充默认值
	method := dto.Method
	if method == "" {
		method = "GET"
	}

	fullDto := &model.CreateAPIMonitorDto{
		Name:                 dto.Name,
		Description:          dto.Description,
		URL:                  dto.URL,
		Method:               method,
		GroupID:              dto.GroupID,
		BusinessFunctionID:   businessFunction.ID,
		BusinessFunctionName: businessFunction.Name,
		GroupName:            businessFunction.GroupName,
		ParentID:             dto.ParentID,        // 设置父监控ID
		HierarchyLevel:       hierarchyLevel,      // 设置层级深度

		// 默认参数
		Timeout:         5000,      // 5秒超时
		CheckInterval:   30,        // 30秒检测
		ExpectedStatus:  200,       // 期望200
		ExpectedKeyword: "success", // 或包含success
		RetryCount:      3,         // 重试3次
		AlertEnabled:    true,      // 启用告警
		PolicyID:        nil,       // 不关联告警策略，使用默认告警
		Status:          1,         // 启用状态
	}

	// 4. 调用原有的CreateMonitor
	s.CreateMonitor(c, fullDto)
}

// UpdateMonitor 更新API监控
func (s *APIMonitorServiceImpl) UpdateMonitor(c *gin.Context, dto *model.UpdateAPIMonitorDto) {
	// 检查是否存在
	existing, err := s.monitorDao.GetByID(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("API监控不存在: %v", err))
		return
	}

	// 更新字段
	if dto.Name != "" {
		existing.Name = dto.Name
	}
	if dto.Description != "" {
		existing.Description = dto.Description
	}
	if dto.GroupID > 0 {
		existing.GroupID = dto.GroupID
	}
	if dto.GroupName != "" {
		existing.GroupName = dto.GroupName
	}
	if dto.BusinessFunctionID > 0 {
		existing.BusinessFunctionID = dto.BusinessFunctionID
	}
	if dto.BusinessFunctionName != "" {
		existing.BusinessFunctionName = dto.BusinessFunctionName
	}
	if dto.ServiceName != "" {
		existing.ServiceName = dto.ServiceName
	}
	if dto.URL != "" {
		existing.URL = dto.URL
	}
	if dto.Method != "" {
		existing.Method = dto.Method
	}
	if dto.Headers != "" {
		existing.Headers = dto.Headers
	}
	if dto.Body != "" {
		existing.Body = dto.Body
	}
	if dto.Timeout > 0 {
		existing.Timeout = dto.Timeout
	}
	if dto.ExpectedStatus > 0 {
		existing.ExpectedStatus = dto.ExpectedStatus
	}
	if dto.ExpectedKeyword != "" {
		existing.ExpectedKeyword = dto.ExpectedKeyword
	}
	if dto.ValidateJSON != "" {
		existing.ValidateJSON = dto.ValidateJSON
	}
	if dto.CheckInterval > 0 {
		existing.CheckInterval = dto.CheckInterval
	}
	if dto.RetryCount > 0 {
		existing.RetryCount = dto.RetryCount
	}
	existing.AlertEnabled = dto.AlertEnabled
	if dto.PolicyID != nil {
		existing.PolicyID = dto.PolicyID
	}
	if dto.Status >= 0 {
		existing.Status = dto.Status
	}
	if dto.Tags != "" {
		existing.Tags = dto.Tags
	}

	if err := s.monitorDao.Update(existing); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新API监控失败: %v", err))
		return
	}

	result.Success(c, existing)
}

// DeleteMonitor 删除API监控
func (s *APIMonitorServiceImpl) DeleteMonitor(c *gin.Context, id uint) {
	log.Printf("🗑️ [DeleteMonitor] 开始删除监控, ID=%d", id)

	// 检查是否存在
	monitor, err := s.monitorDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "监控不存在")
		return
	}

	// 检查是否有子监控
	children, err := s.monitorDao.GetChildMonitors(id)
	if err != nil {
		log.Printf("⚠️  [DeleteMonitor] 查询子监控失败: %v", err)
	}

	if len(children) > 0 {
		log.Printf("🗑️ [DeleteMonitor] 检测到 %d 个子监控，将级联删除", len(children))

		// 先删除所有子监控（包括Pushgateway数据）
		for _, child := range children {
			// 清理Pushgateway数据
			if err := s.cleanupPushgatewayMetrics(monitor.Name, child.Name); err != nil {
				log.Printf("⚠️  清理子监控Pushgateway数据失败[%s/%s]: %v", monitor.Name, child.Name, err)
			}

			if err := s.monitorDao.Delete(child.ID); err != nil {
				result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("删除子监控失败[ID=%d]: %v", child.ID, err))
				return
			}
			log.Printf("🗑️ [DeleteMonitor] 已删除子监控: ID=%d, Name=%s", child.ID, child.Name)
		}
	}

	// 清理父监控的Pushgateway数据
	// 如果是顶级监控，清理整个job；如果是子监控，只清理自己的instance
	if monitor.ParentID == nil {
		// 顶级监控：清理整个job
		if err := s.cleanupPushgatewayJob(monitor.Name); err != nil {
			log.Printf("⚠️  清理父监控Pushgateway数据失败[job=%s]: %v", monitor.Name, err)
		}
	} else {
		// 子监控：只清理自己的instance
		parentMonitor, err := s.monitorDao.GetByID(*monitor.ParentID)
		if err == nil {
			if err := s.cleanupPushgatewayMetrics(parentMonitor.Name, monitor.Name); err != nil {
				log.Printf("⚠️  清理子监控Pushgateway数据失败[%s/%s]: %v", parentMonitor.Name, monitor.Name, err)
			}
		}
	}

	// 删除父监控
	if err := s.monitorDao.Delete(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("删除API监控失败: %v", err))
		return
	}

	log.Printf("✅ [DeleteMonitor] 删除成功: ID=%d, Name=%s, 子监控数量=%d", id, monitor.Name, len(children))
	result.Success(c, gin.H{
		"message": "删除成功",
		"deletedCount": 1 + len(children),
		"childCount": len(children),
	})
}

// GetMonitor 获取API监控详情
func (s *APIMonitorServiceImpl) GetMonitor(c *gin.Context, id uint) {
	log.Printf("📖 [GetMonitor] 开始获取监控详情, ID=%d", id)

	monitor, err := s.monitorDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取API监控失败: %v", err))
		return
	}
	log.Printf("📖 [GetMonitor] 成功获取基本信息: Name=%s, ParentID=%v", monitor.Name, monitor.ParentID)

	// 计算健康率
	healthRate, err := s.resultDao.GetHealthRate24h(id)
	if err != nil {
		log.Printf("⚠️  获取健康率失败: %v", err)
	} else {
		monitor.HealthRate = &healthRate
	}

	// ======================== 新增：加载完整的分层信息 ========================
	// 1. 如果是子监控，加载父监控名称
	if monitor.ParentID != nil && *monitor.ParentID > 0 {
		log.Printf("📖 [GetMonitor] 这是子监控，正在加载父监控信息, ParentID=%d", *monitor.ParentID)
		parentMonitor, err := s.monitorDao.GetByID(*monitor.ParentID)
		if err == nil {
			monitor.ParentName = parentMonitor.Name
			monitor.ParentMonitor = parentMonitor
			log.Printf("📖 [GetMonitor] 成功加载父监控: ParentName=%s", parentMonitor.Name)
		} else {
			log.Printf("⚠️  [GetMonitor] 加载父监控失败: %v", err)
		}
	}

	// 2. 如果是父监控，加载所有子监控及其统计信息
	log.Printf("📖 [GetMonitor] 正在查询子监控...")
	children, err := s.monitorDao.GetChildMonitors(id)
	if err != nil {
		log.Printf("⚠️  [GetMonitor] 查询子监控失败: %v", err)
	} else {
		log.Printf("📖 [GetMonitor] 查询到 %d 个子监控", len(children))
		if len(children) > 0 {
			monitor.ChildMonitors = children
			monitor.ChildCount = len(children)

			// 统计健康的子监控
			healthyCount := 0
			for _, child := range children {
				log.Printf("📖 [GetMonitor] 子监控: ID=%d, Name=%s, URL=%s, HealthStatus=%s",
					child.ID, child.Name, child.URL, child.HealthStatus)
				if child.HealthStatus == "healthy" {
					healthyCount++
				}
			}
			monitor.HealthyChildren = healthyCount
			log.Printf("📖 [GetMonitor] 子监控统计: Total=%d, Healthy=%d", len(children), healthyCount)
		}
	}

	log.Printf("✅ [GetMonitor] 返回完整监控信息")
	result.Success(c, monitor)
}

// GetMonitorList 获取API监控列表
func (s *APIMonitorServiceImpl) GetMonitorList(c *gin.Context, query *model.APIMonitorListQuery) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	monitors, total, err := s.monitorDao.GetList(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取API监控列表失败: %v", err))
		return
	}

	// 批量计算健康率
	if len(monitors) > 0 {
		monitorIDs := make([]uint, len(monitors))
		for i, monitor := range monitors {
			monitorIDs[i] = monitor.ID
		}

		healthRates, err := s.resultDao.GetHealthRateBatch24h(monitorIDs)
		if err != nil {
			log.Printf("⚠️  获取健康率失败: %v", err)
		} else {
			// 将健康率赋值给每个监控项
			for _, monitor := range monitors {
				if rate, exists := healthRates[monitor.ID]; exists {
					monitor.HealthRate = &rate
				}
			}
		}

		// ======================== 新增：批量加载子监控统计信息 ========================
		// 收集所有有子监控的父监控ID
		parentIDs := make([]uint, 0)
		for _, monitor := range monitors {
			if monitor.ParentID == nil { // 只有顶级监控才可能有子监控
				parentIDs = append(parentIDs, monitor.ID)
			}
		}

		if len(parentIDs) > 0 {
			// 批量获取子监控统计信息
			childrenStats, err := s.monitorDao.GetChildrenStatsInBatch(parentIDs)
			if err != nil {
				log.Printf("⚠️  获取子监控统计失败: %v", err)
			} else {
				// 将统计信息赋值给对应的监控项
				for _, monitor := range monitors {
					if stats, exists := childrenStats[monitor.ID]; exists {
						monitor.ChildCount = stats.ChildCount
						monitor.HealthyChildren = stats.HealthyChildren
						// 如果该监控是聚合监控，更新其聚合健康状态
						if monitor.ChildCount > 0 {
							monitor.IsAggregate = true
						}
					}
				}
			}
		}
	}

	result.Success(c, gin.H{
		"list":     monitors,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

// UpdateMonitorStatus 更新API监控状态
func (s *APIMonitorServiceImpl) UpdateMonitorStatus(c *gin.Context, id uint, status int) {
	// 先获取监控配置
	monitor, err := s.monitorDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取监控配置失败: %v", err))
		return
	}

	// 更新状态
	if err := s.monitorDao.UpdateStatus(id, status); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新状态失败: %v", err))
		return
	}

	statusText := "禁用"
	if status == 1 {
		statusText = "启用"
		// 启用时，立即执行一次检测
		go func() {
			log.Printf("🔍 启用监控[%s]，立即执行检测...", monitor.Name)
			checkResult := s.ExecuteCheck(monitor)

			// 保存检测结果
			if err := s.resultDao.Create(checkResult); err != nil {
				log.Printf("⚠️  保存检测结果失败: %v", err)
			}

			// 更新监控的最新检测结果
			monitor.LastCheckTime = &checkResult.CheckTime
			monitor.LastResponseTime = checkResult.ResponseTime
			monitor.LastStatusCode = checkResult.StatusCode
			monitor.LastErrorMessage = checkResult.ErrorMessage

			if checkResult.Success {
				monitor.HealthStatus = "healthy"
				monitor.ContinuousFailures = 0
			} else {
				monitor.HealthStatus = "unhealthy"
				monitor.ContinuousFailures++
			}

			if err := s.monitorDao.UpdateLastCheckResult(monitor); err != nil {
				log.Printf("⚠️  更新监控状态失败: %v", err)
			}

			log.Printf("✅ 监控[%s]检测完成: %v", monitor.Name, checkResult.Success)
		}()
	}

	result.Success(c, fmt.Sprintf("API监控已%s", statusText))
}

// GetStats 获取统计信息
func (s *APIMonitorServiceImpl) GetStats(c *gin.Context) {
	stats, err := s.monitorDao.GetStats()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取统计信息失败: %v", err))
		return
	}
	result.Success(c, stats)
}

// GetGroupStats 获取按业务组聚合的统计（树形结构）- 优化版，避免N+1查询
func (s *APIMonitorServiceImpl) GetGroupStats(c *gin.Context) {
	stats, err := s.monitorDao.GetGroupStats()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取业务组统计失败: %v", err))
		return
	}

	// 收集所有子组的ID（ParentID != 0）
	var childGroupIDs []uint
	for _, groupStat := range stats {
		if groupStat.ParentID != 0 {
			childGroupIDs = append(childGroupIDs, groupStat.GroupID)
		}
	}

	// 一次性查询所有子组的API监控数据，避免N+1查询
	var allBusinessFunctions map[uint][]model.APIMonitorBusinessFunctionStats
	if len(childGroupIDs) > 0 {
		allBusinessFunctions, err = s.getBusinessFunctionStatsByGroupIDs(childGroupIDs)
		if err != nil {
			log.Printf("⚠️  获取业务功能统计失败: %v", err)
			allBusinessFunctions = make(map[uint][]model.APIMonitorBusinessFunctionStats)
		}
	}

	// 为每个子组分配业务功能统计
	for _, groupStat := range stats {
		if groupStat.ParentID != 0 {
			if functions, exists := allBusinessFunctions[groupStat.GroupID]; exists {
				groupStat.BusinessFunctions = functions
			}
		}
	}

	// 构建树形结构
	tree := model.BuildAPIMonitorGroupTree(stats)

	result.Success(c, tree)
}

// GetServiceStats 获取指定业务组的服务统计
func (s *APIMonitorServiceImpl) GetServiceStats(c *gin.Context, groupID uint) {
	serviceStats, err := s.getBusinessFunctionStatsByGroupID(groupID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取服务统计失败: %v", err))
		return
	}
	result.Success(c, serviceStats)
}

// getBusinessFunctionStatsByGroupID 内部方法：获取业务组下的业务功能统计
func (s *APIMonitorServiceImpl) getBusinessFunctionStatsByGroupID(groupID uint) ([]model.APIMonitorBusinessFunctionStats, error) {
	monitors, err := s.monitorDao.GetByGroupID(groupID)
	if err != nil {
		return nil, err
	}

	return s.buildBusinessFunctionStats(monitors), nil
}

// getBusinessFunctionStatsByGroupIDs 内部方法：批量获取多个业务组的业务功能统计（优化版，避免N+1查询）
func (s *APIMonitorServiceImpl) getBusinessFunctionStatsByGroupIDs(groupIDs []uint) (map[uint][]model.APIMonitorBusinessFunctionStats, error) {
	if len(groupIDs) == 0 {
		return make(map[uint][]model.APIMonitorBusinessFunctionStats), nil
	}

	// 一次性查询所有业务组的监控数据
	monitors, err := s.monitorDao.GetByGroupIDs(groupIDs)
	if err != nil {
		return nil, err
	}

	// 按 group_id 分组
	monitorsByGroup := make(map[uint][]*model.APIMonitor)
	for _, monitor := range monitors {
		monitorsByGroup[monitor.GroupID] = append(monitorsByGroup[monitor.GroupID], monitor)
	}

	// 为每个业务组构建业务功能统计
	result := make(map[uint][]model.APIMonitorBusinessFunctionStats, len(groupIDs))
	for _, groupID := range groupIDs {
		if groupMonitors, exists := monitorsByGroup[groupID]; exists {
			result[groupID] = s.buildBusinessFunctionStats(groupMonitors)
		} else {
			result[groupID] = []model.APIMonitorBusinessFunctionStats{}
		}
	}

	return result, nil
}

// buildBusinessFunctionStats 内部方法：从监控列表构建业务功能统计
func (s *APIMonitorServiceImpl) buildBusinessFunctionStats(monitors []*model.APIMonitor) []model.APIMonitorBusinessFunctionStats {
	// 按业务功能ID分组
	functionMap := make(map[uint]*model.APIMonitorBusinessFunctionStats)
	for _, monitor := range monitors {
		functionID := monitor.BusinessFunctionID
		if functionID == 0 {
			continue // 跳过未分类的
		}

		if _, exists := functionMap[functionID]; !exists {
			functionMap[functionID] = &model.APIMonitorBusinessFunctionStats{
				BusinessFunctionID:   functionID,
				BusinessFunctionName: monitor.BusinessFunctionName,
				APIs:                 []model.APIMonitor{},
			}
		}

		stat := functionMap[functionID]
		stat.TotalAPIs++
		stat.APIs = append(stat.APIs, *monitor)

		if monitor.HealthStatus == "healthy" {
			stat.HealthyAPIs++
		} else if monitor.HealthStatus == "unhealthy" {
			stat.UnhealthyAPIs++
		}

		if monitor.LastResponseTime > 0 {
			stat.AvgResponseTime = (stat.AvgResponseTime*float64(stat.TotalAPIs-1) +
				float64(monitor.LastResponseTime)) / float64(stat.TotalAPIs)
		}
	}

	// 转换为数组
	var stats []model.APIMonitorBusinessFunctionStats
	for _, stat := range functionMap {
		stats = append(stats, *stat)
	}

	return stats
}

// TestAPI 测试API接口
func (s *APIMonitorServiceImpl) TestAPI(c *gin.Context, dto *model.TestAPIDto) {
	// 构建临时监控对象
	tempMonitor := &model.APIMonitor{
		URL:            dto.URL,
		Method:         dto.Method,
		Headers:        dto.Headers,
		Body:           dto.Body,
		Timeout:        dto.Timeout,
		ExpectedStatus: dto.ExpectedStatus,
	}

	// 设置默认值
	if tempMonitor.Method == "" {
		tempMonitor.Method = "GET"
	}
	if tempMonitor.Timeout == 0 {
		tempMonitor.Timeout = 5000
	}
	if tempMonitor.ExpectedStatus == 0 {
		tempMonitor.ExpectedStatus = 200
	}

	// 执行检测
	testResult := s.ExecuteCheck(tempMonitor)

	result.Success(c, gin.H{
		"success":      testResult.Success,
		"statusCode":   testResult.StatusCode,
		"responseTime": testResult.ResponseTime,
		"errorMessage": testResult.ErrorMessage,
		"responseBody": testResult.ResponseBody,
	})
}

// ExecuteCheck 执行API检测（核心逻辑）
func (s *APIMonitorServiceImpl) ExecuteCheck(monitor *model.APIMonitor) *model.APIMonitorResult {
	startTime := time.Now()

	// 构建HTTP请求
	var reqBody io.Reader
	if monitor.Body != "" {
		reqBody = bytes.NewBufferString(monitor.Body)
	}

	req, err := http.NewRequest(monitor.Method, monitor.URL, reqBody)
	if err != nil {
		return &model.APIMonitorResult{
			MonitorID:    monitor.ID,
			CheckTime:    util.HTime{Time: time.Now()},
			Success:      false,
			ErrorMessage: fmt.Sprintf("构建请求失败: %v", err),
		}
	}

	// 添加请求头
	if monitor.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(monitor.Headers), &headers); err == nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
	}

	// 执行HTTP请求
	client := &http.Client{
		Timeout: time.Duration(monitor.Timeout) * time.Millisecond,
	}
	resp, err := client.Do(req)
	responseTime := int(time.Since(startTime).Milliseconds())

	checkResult := &model.APIMonitorResult{
		MonitorID:    monitor.ID,
		CheckTime:    util.HTime{Time: time.Now()},
		ResponseTime: responseTime,
	}

	if err != nil {
		checkResult.Success = false
		checkResult.ErrorMessage = fmt.Sprintf("请求失败: %v", err)
		return checkResult
	}
	defer resp.Body.Close()

	checkResult.StatusCode = resp.StatusCode

	// 读取响应体
	body, _ := io.ReadAll(resp.Body)
	checkResult.ResponseBody = string(body)

	// 校验状态码
	if resp.StatusCode != monitor.ExpectedStatus {
		checkResult.Success = false
		checkResult.ErrorMessage = fmt.Sprintf("期望状态码%d, 实际%d", monitor.ExpectedStatus, resp.StatusCode)
		return checkResult
	}

	// 校验关键字
	if monitor.ExpectedKeyword != "" {
		if !strings.Contains(string(body), monitor.ExpectedKeyword) {
			checkResult.Success = false
			checkResult.ErrorMessage = fmt.Sprintf("响应体不包含期望关键字: %s", monitor.ExpectedKeyword)
			return checkResult
		}
	}

	// 校验JSON（如果配置了校验规则）
	if monitor.ValidateJSON != "" {
		// TODO: 实现JSON Path校验逻辑
	}

	checkResult.Success = true
	return checkResult
}

// PushMetricToPrometheus 推送指标到Prometheus（通过Pushgateway）
func (s *APIMonitorServiceImpl) PushMetricToPrometheus(monitor *model.APIMonitor, checkResult *model.APIMonitorResult) error {
	pushgatewayURL := config.Config.Monitor.Pushgateway.URL
	if pushgatewayURL == "" {
		return fmt.Errorf("Pushgateway地址未配置")
	}

	// 确定job名称：如果有父监控，使用父监控的名称；否则使用自己的名称
	var jobName string
	var parentMonitorName string
	if monitor.ParentID != nil && *monitor.ParentID > 0 {
		// 子监控：获取父监控名称作为job
		parentMonitor, err := s.monitorDao.GetByID(*monitor.ParentID)
		if err != nil {
			log.Printf("⚠️  获取父监控失败(ID=%d): %v，使用当前监控名称作为job", *monitor.ParentID, err)
			jobName = monitor.Name
			parentMonitorName = monitor.Name
		} else {
			jobName = parentMonitor.Name
			parentMonitorName = parentMonitor.Name
		}
	} else {
		// 顶级监控：使用自己的名称作为job
		jobName = monitor.Name
		parentMonitorName = monitor.Name
	}

	// 构建Prometheus指标
	successValue := 0
	if checkResult.Success {
		successValue = 1
	}

	metrics := fmt.Sprintf(`# HELP api_monitor_success API监控成功状态(1=成功,0=失败)
# TYPE api_monitor_success gauge
api_monitor_success{monitor_id="%d",name="%s",parent_monitor="%s",service="%s",group_id="%d",group_name="%s",url="%s",method="%s"} %d

# HELP api_monitor_response_time API响应时间(毫秒)
# TYPE api_monitor_response_time gauge
api_monitor_response_time{monitor_id="%d",name="%s",parent_monitor="%s",service="%s",group_id="%d",group_name="%s",url="%s",method="%s"} %d

# HELP api_monitor_status_code HTTP状态码
# TYPE api_monitor_status_code gauge
api_monitor_status_code{monitor_id="%d",name="%s",parent_monitor="%s",service="%s",group_id="%d",group_name="%s",url="%s",method="%s"} %d
`,
		monitor.ID, monitor.Name, parentMonitorName, monitor.ServiceName, monitor.GroupID, monitor.GroupName, monitor.URL, monitor.Method, successValue,
		monitor.ID, monitor.Name, parentMonitorName, monitor.ServiceName, monitor.GroupID, monitor.GroupName, monitor.URL, monitor.Method, checkResult.ResponseTime,
		monitor.ID, monitor.Name, parentMonitorName, monitor.ServiceName, monitor.GroupID, monitor.GroupName, monitor.URL, monitor.Method, checkResult.StatusCode,
	)

	// 推送到Pushgateway：使用父监控名称作为job，当前监控名称作为instance
	pushURL := fmt.Sprintf("%s/metrics/job/%s/instance/%s", pushgatewayURL, jobName, monitor.Name)

	req, err := http.NewRequest("POST", pushURL, strings.NewReader(metrics))
	if err != nil {
		return fmt.Errorf("构建推送请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("推送到Pushgateway失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("推送失败，状态码%d: %s", resp.StatusCode, string(body))
	}

	log.Printf("✅ 推送指标到Pushgateway成功 [监控:%s, URL:%s]", monitor.Name, monitor.URL)
	return nil
}

// ======================== 分层关系管理方法 ========================

const (
	MaxHierarchyDepth = 3 // 最大层级深度限制（0=顶级, 1=一级子, 2=二级子, 3=三级子）
)

// AddChildMonitor 为父监控添加子监控
func (s *APIMonitorServiceImpl) AddChildMonitor(c *gin.Context, parentID uint, childMonitorID uint) {
	// 验证父监控是否存在
	parentMonitor, err := s.monitorDao.GetByID(parentID)
	if err != nil {
		result.Failed(c, 400, "父监控不存在")
		return
	}

	// 验证子监控是否存在
	childMonitor, err := s.monitorDao.GetByID(childMonitorID)
	if err != nil {
		result.Failed(c, 400, "子监控不存在")
		return
	}

	// 检查是否会形成循环依赖
	hasCircular, err := s.monitorDao.CheckCircularDependency(childMonitorID, parentID)
	if err != nil || hasCircular {
		result.Failed(c, 400, "会形成循环依赖，操作被禁止")
		return
	}

	// 检查层级深度限制（获取childMonitorID的最大深度，然后加上1）
	childMaxDepth, err := s.monitorDao.GetMaxHierarchyDepth(childMonitorID)
	if err == nil && childMaxDepth >= MaxHierarchyDepth {
		result.Failed(c, 400, fmt.Sprintf("已达到最大层级深度限制(%d层)，无法添加子监控", MaxHierarchyDepth))
		return
	}

	// 计算新的层级深度
	newHierarchyLevel := parentMonitor.HierarchyLevel + 1

	// 更新子监控的parent_id和hierarchy_level
	childMonitor.ParentID = &parentID
	childMonitor.HierarchyLevel = int(newHierarchyLevel)

	if err := s.monitorDao.Update(childMonitor); err != nil {
		result.Failed(c, 500, "更新子监控失败: "+err.Error())
		return
	}

	// 更新父监控的聚合数据
	if err := s.monitorDao.UpdateParentAggregateStats(parentID); err != nil {
		result.Failed(c, 500, "更新父监控聚合数据失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"message": "子监控添加成功",
		"childMonitorId": childMonitorID,
		"parentMonitorId": parentID,
	})
}

// ConvertToChildMonitor 将独立监控转换为子监控
func (s *APIMonitorServiceImpl) ConvertToChildMonitor(c *gin.Context, monitorID uint, parentID uint) {
	// 验证父监控是否存在
	parentMonitor, err := s.monitorDao.GetByID(parentID)
	if err != nil {
		result.Failed(c, 400, "父监控不存在")
		return
	}

	// 验证子监控是否存在
	monitor, err := s.monitorDao.GetByID(monitorID)
	if err != nil {
		result.Failed(c, 400, "监控不存在")
		return
	}

	// 检查循环依赖
	hasCircular, err := s.monitorDao.CheckCircularDependency(monitorID, parentID)
	if err != nil || hasCircular {
		result.Failed(c, 400, "会形成循环依赖，操作被禁止")
		return
	}

	// 更新监控
	monitor.ParentID = &parentID
	monitor.HierarchyLevel = int(parentMonitor.HierarchyLevel) + 1

	if err := s.monitorDao.Update(monitor); err != nil {
		result.Failed(c, 500, "转换失败: "+err.Error())
		return
	}

	// 更新父监控的聚合数据
	if err := s.monitorDao.UpdateParentAggregateStats(parentID); err != nil {
		result.Failed(c, 500, "更新父监控聚合数据失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"message": "监控转换成功",
		"monitorId": monitorID,
		"parentMonitorId": parentID,
	})
}

// GetMonitorHierarchy 获取监控的树形结构（包含子监控）
func (s *APIMonitorServiceImpl) GetMonitorHierarchy(c *gin.Context, monitorID uint, depth int) {
	monitor, err := s.monitorDao.GetByID(monitorID)
	if err != nil {
		result.Failed(c, 400, "监控不存在")
		return
	}

	if depth <= 0 {
		depth = 10 // 默认深度限制
	}

	hierarchyDto := s.buildHierarchyDto(monitor, depth)
	result.Success(c, hierarchyDto)
}

// buildHierarchyDto 递归构建树形结构DTO
func (s *APIMonitorServiceImpl) buildHierarchyDto(monitor *model.APIMonitor, remainDepth int) model.HierarchyDto {
	dto := model.HierarchyDto{
		ID:             monitor.ID,
		Name:           monitor.Name,
		URL:            monitor.URL,
		HealthStatus:   monitor.HealthStatus,
		HierarchyLevel: monitor.HierarchyLevel,
		ChildCount:     0,
		HealthyChildren: 0,
	}

	if remainDepth <= 0 {
		return dto
	}

	// 获取子监控
	children, _ := s.monitorDao.GetChildMonitors(monitor.ID)
	if len(children) > 0 {
		dto.ChildCount = len(children)
		dto.HealthyChildren = 0
		for _, child := range children {
			if child.HealthStatus == "healthy" {
				dto.HealthyChildren++
			}
			dto.Children = append(dto.Children, s.buildHierarchyDto(child, remainDepth-1))
		}
	}

	return dto
}

// DeleteMonitorWithChildren 删除监控及其子监控
func (s *APIMonitorServiceImpl) DeleteMonitorWithChildren(c *gin.Context, monitorID uint, operation string) {
	monitor, err := s.monitorDao.GetByID(monitorID)
	if err != nil {
		result.Failed(c, 400, "监控不存在")
		return
	}

	// 获取所有后代监控
	descendants, err := s.monitorDao.GetAllDescendants(monitorID)
	if err != nil {
		result.Failed(c, 500, "获取子监控失败: "+err.Error())
		return
	}

	switch operation {
	case "cascade":
		// 级联删除所有子监控及其检测结果
		allToDelete := append([]*model.APIMonitor{monitor}, descendants...)
		for _, m := range allToDelete {
			if err := s.monitorDao.Delete(m.ID); err != nil {
				result.Failed(c, 500, "删除失败: "+err.Error())
				return
			}
		}
		result.Success(c, map[string]interface{}{
			"message": "监控及其子监控已删除",
			"deletedCount": len(allToDelete),
		})

	case "orphan":
		// 把子监控变为独立监控
		for _, child := range descendants {
			child.ParentID = nil
			child.HierarchyLevel = 0
			if err := s.monitorDao.Update(child); err != nil {
				result.Failed(c, 500, "更新子监控失败: "+err.Error())
				return
			}
		}
		if err := s.monitorDao.Delete(monitorID); err != nil {
			result.Failed(c, 500, "删除监控失败: "+err.Error())
			return
		}
		result.Success(c, map[string]interface{}{
			"message": "监控已删除，子监控已变为独立监控",
			"orphanedCount": len(descendants),
		})

	default:
		// 防止删除（如果有子监控）
		if len(descendants) > 0 {
			result.Failed(c, 400, fmt.Sprintf("该监控有%d个子监控，无法删除。请选择其他操作", len(descendants)))
			return
		}
		// 无子监控，可以删除
		if err := s.monitorDao.Delete(monitorID); err != nil {
			result.Failed(c, 500, "删除失败: "+err.Error())
			return
		}
		result.Success(c, map[string]interface{}{
			"message": "监控已删除",
		})
	}
}

// CalculateParentHealthStatus 计算父监控的聚合健康状态
func (s *APIMonitorServiceImpl) CalculateParentHealthStatus(parentID uint) (string, error) {
	children, err := s.monitorDao.GetChildMonitors(parentID)
	if err != nil {
		return "unknown", err
	}

	if len(children) == 0 {
		return "healthy", nil // 没有子监控时默认为healthy
	}

	// 所有子监控都健康则为healthy，否则为unhealthy
	for _, child := range children {
		if child.HealthStatus != "healthy" {
			return "unhealthy", nil
		}
	}

	return "healthy", nil
}

// UpdateAllParentAggregateStats 定时任务：更新所有父监控的聚合数据
func (s *APIMonitorServiceImpl) UpdateAllParentAggregateStats() error {
	// 获取所有有子监控的父监控
	monitors, err := s.monitorDao.GetAll()
	if err != nil {
		return err
	}

	// 收集所有有子监控的parentID
	parentIDs := make([]uint, 0)
	parentMap := make(map[uint]bool)

	for _, m := range monitors {
		if m.ParentID != nil && *m.ParentID > 0 {
			if !parentMap[*m.ParentID] {
				parentIDs = append(parentIDs, *m.ParentID)
				parentMap[*m.ParentID] = true
			}
		}
	}

	if len(parentIDs) == 0 {
		return nil
	}

	// 批量更新父监控的聚合数据
	for _, parentID := range parentIDs {
		if err := s.monitorDao.UpdateParentAggregateStats(parentID); err != nil {
			log.Printf("⚠️ 更新父监控%d的聚合数据失败: %v", parentID, err)
		}
	}

	log.Printf("✅ 更新了%d个父监控的聚合数据", len(parentIDs))
	return nil
}

// RemoveParentRelationship 移除子监控与父监控的关系，使其变为独立监控
func (s *APIMonitorServiceImpl) RemoveParentRelationship(c *gin.Context, childMonitorID uint) {
	monitor, err := s.monitorDao.GetByID(childMonitorID)
	if err != nil {
		result.Failed(c, 400, "监控不存在")
		return
	}

	if monitor.ParentID == nil {
		result.Failed(c, 400, "该监控没有父监控关系")
		return
	}

	oldParentID := *monitor.ParentID

	// 清除父监控关系
	monitor.ParentID = nil
	monitor.HierarchyLevel = 0

	if err := s.monitorDao.Update(monitor); err != nil {
		result.Failed(c, 500, "更新失败: "+err.Error())
		return
	}

	// 更新原父监控的聚合数据
	if err := s.monitorDao.UpdateParentAggregateStats(oldParentID); err != nil {
		result.Failed(c, 500, "更新父监控聚合数据失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"message": "父监控关系已移除，监控现为独立监控",
		"monitorId": childMonitorID,
	})
}

// cleanupPushgatewayJob 清理Pushgateway中的整个job
func (s *APIMonitorServiceImpl) cleanupPushgatewayJob(jobName string) error {
	pushgatewayURL := config.Config.Monitor.Pushgateway.URL
	if pushgatewayURL == "" {
		return fmt.Errorf("Pushgateway地址未配置")
	}

	// DELETE /metrics/job/<job_name>
	deleteURL := fmt.Sprintf("%s/metrics/job/%s", pushgatewayURL, jobName)
	log.Printf("🧹 清理Pushgateway job: %s", deleteURL)

	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return fmt.Errorf("构建删除请求失败: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("清理失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("清理失败，状态码%d: %s", resp.StatusCode, string(body))
	}

	log.Printf("✅ 清理Pushgateway job成功: %s", jobName)
	return nil
}

// cleanupPushgatewayMetrics 清理Pushgateway中的特定instance
func (s *APIMonitorServiceImpl) cleanupPushgatewayMetrics(jobName string, instanceName string) error {
	pushgatewayURL := config.Config.Monitor.Pushgateway.URL
	if pushgatewayURL == "" {
		return fmt.Errorf("Pushgateway地址未配置")
	}

	// DELETE /metrics/job/<job_name>/instance/<instance_name>
	deleteURL := fmt.Sprintf("%s/metrics/job/%s/instance/%s", pushgatewayURL, jobName, instanceName)
	log.Printf("🧹 清理Pushgateway指标: %s", deleteURL)

	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return fmt.Errorf("构建删除请求失败: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("清理失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("清理失败，状态码%d: %s", resp.StatusCode, string(body))
	}

	log.Printf("✅ 清理Pushgateway指标成功: job=%s, instance=%s", jobName, instanceName)
	return nil
}

