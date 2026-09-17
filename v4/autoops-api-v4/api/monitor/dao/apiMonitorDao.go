package dao

import (
	"log"

	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
)

type APIMonitorDao interface {
	Create(monitor *model.APIMonitor) error
	Update(monitor *model.APIMonitor) error
	Delete(id uint) error
	GetByID(id uint) (*model.APIMonitor, error)
	GetAll() ([]*model.APIMonitor, error)
	GetAllEnabled() ([]*model.APIMonitor, error)
	GetList(query *model.APIMonitorListQuery) ([]*model.APIMonitor, int64, error)
	GetByGroupID(groupID uint) ([]*model.APIMonitor, error)
	GetByGroupIDs(groupIDs []uint) ([]*model.APIMonitor, error) // 新增：批量查询
	GetByServiceName(serviceName string) ([]*model.APIMonitor, error)
	UpdateStatus(id uint, status int) error
	UpdateHealthStatus(id uint, healthStatus string) error
	UpdateLastCheckResult(monitor *model.APIMonitor) error
	GetStats() (*model.APIMonitorStats, error)
	GetGroupStats() ([]*model.APIMonitorGroupStats, error)
	CheckNameExistsInGroup(groupID uint, name string, excludeID uint) (bool, error)

	// 分层关系查询方法
	GetMonitorTree(parentID uint) (*model.APIMonitor, error)              // 获取父监控及其所有子监控
	GetTopLevelMonitors(query *model.APIMonitorListQuery) ([]*model.APIMonitor, int64, error) // 获取所有顶级监控
	GetChildMonitors(parentID uint) ([]*model.APIMonitor, error)          // 获取指定监控的所有直接子监控
	GetAllDescendants(parentID uint) ([]*model.APIMonitor, error)         // 获取指定监控的所有后代监控（递归）
	CheckCircularDependency(childID uint, parentID uint) (bool, error)    // 检查是否会形成循环依赖
	GetChildrenStatsInBatch(parentIDs []uint) (map[uint]*model.ChildrenStats, error) // 批量获取子监控的统计信息
	UpdateParentAggregateStats(parentID uint) error                       // 更新父监控的聚合字段
	GetMaxHierarchyDepth(monitorID uint) (int, error)                     // 获取指定监控的最大层级深度（用于检查是否超过限制）
}

type APIMonitorDaoImpl struct{}

func NewAPIMonitorDao() APIMonitorDao {
	return &APIMonitorDaoImpl{}
}

// Create 创建API监控
func (d *APIMonitorDaoImpl) Create(monitor *model.APIMonitor) error {
	return common.GetDB().Create(monitor).Error
}

// Update 更新API监控
func (d *APIMonitorDaoImpl) Update(monitor *model.APIMonitor) error {
	return common.GetDB().Model(&model.APIMonitor{}).
		Where("id = ?", monitor.ID).
		Updates(monitor).Error
}

// Delete 删除API监控
func (d *APIMonitorDaoImpl) Delete(id uint) error {
	return common.GetDB().Delete(&model.APIMonitor{}, id).Error
}

// GetByID 根据ID获取API监控
func (d *APIMonitorDaoImpl) GetByID(id uint) (*model.APIMonitor, error) {
	var monitor model.APIMonitor
	err := common.GetDB().Where("id = ?", id).First(&monitor).Error
	return &monitor, err
}

// GetAll 获取所有API监控
func (d *APIMonitorDaoImpl) GetAll() ([]*model.APIMonitor, error) {
	var monitors []*model.APIMonitor
	err := common.GetDB().Find(&monitors).Error
	return monitors, err
}

// GetAllEnabled 获取所有启用的API监控（会过滤掉父监控被禁用的子监控）
func (d *APIMonitorDaoImpl) GetAllEnabled() ([]*model.APIMonitor, error) {
	var monitors []*model.APIMonitor

	// 查询所有启用的监控，但要确保：
	// 1. 监控本身是启用的（status = 1）
	// 2. 如果有父监控（parent_id IS NOT NULL），则父监控也必须是启用的
	err := common.GetDB().
		Where("status = ?", 1).
		Where(
			common.GetDB().Where("parent_id IS NULL"). // 顶级监控（没有父监控）
				Or("parent_id IN (?)", // 或者父监控是启用状态
					common.GetDB().Model(&model.APIMonitor{}).
						Select("id").
						Where("status = ?", 1),
				),
		).
		Find(&monitors).Error

	return monitors, err
}

// GetList 分页查询API监控列表
func (d *APIMonitorDaoImpl) GetList(query *model.APIMonitorListQuery) ([]*model.APIMonitor, int64, error) {
	var monitors []*model.APIMonitor
	var total int64

	db := common.GetDB().Model(&model.APIMonitor{})

	// 分层关系筛选：如果指定了 parentId，则查询该父监控的直接子监控
	if query.ParentID != nil {
		db = db.Where("parent_id = ?", *query.ParentID)
	} else if query.OnlyTopLevel {
		// 只查询顶级监控（parentId=null）- 默认向下兼容行为
		db = db.Where("parent_id IS NULL")
	}

	// 业务组筛选
	if query.GroupID != nil {
		db = db.Where("group_id = ?", *query.GroupID)
	}

	// 业务功能筛选
	if query.BusinessFunctionID != nil {
		db = db.Where("business_function_id = ?", *query.BusinessFunctionID)
	}
	// 服务名称筛选
	if query.ServiceName != "" {
		db = db.Where("service_name = ?", query.ServiceName)
	}

	// 状态筛选
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	// 关键字搜索
	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR url LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 统计总数
	db.Count(&total)

	// 分页
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 是否加载子监控
	if query.IncludeChildren {
		err := db.Order("id DESC").Preload("ChildMonitors").Find(&monitors).Error
		return monitors, total, err
	}

	err := db.Order("id DESC").Find(&monitors).Error
	return monitors, total, err
}

// GetByGroupID 根据业务组ID获取API监控列表
func (d *APIMonitorDaoImpl) GetByGroupID(groupID uint) ([]*model.APIMonitor, error) {
	var monitors []*model.APIMonitor
	err := common.GetDB().Where("group_id = ?", groupID).Find(&monitors).Error
	return monitors, err
}

// GetByGroupIDs 根据多个业务组ID批量获取API监控列表（优化版，避免N+1查询）
func (d *APIMonitorDaoImpl) GetByGroupIDs(groupIDs []uint) ([]*model.APIMonitor, error) {
	if len(groupIDs) == 0 {
		return []*model.APIMonitor{}, nil
	}

	var monitors []*model.APIMonitor
	err := common.GetDB().Where("group_id IN ?", groupIDs).Find(&monitors).Error
	return monitors, err
}

// GetByServiceName 根据服务名称获取API监控列表
func (d *APIMonitorDaoImpl) GetByServiceName(serviceName string) ([]*model.APIMonitor, error) {
	var monitors []*model.APIMonitor
	err := common.GetDB().Where("service_name = ?", serviceName).Find(&monitors).Error
	return monitors, err
}

// UpdateStatus 更新状态
func (d *APIMonitorDaoImpl) UpdateStatus(id uint, status int) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// 禁用时，重置健康状态为healthy，清空错误信息和连续失败次数
	if status == 0 {
		updates["health_status"] = "healthy"
		updates["last_error_message"] = ""
		updates["continuous_failures"] = 0
	} else if status == 1 {
		// 启用时，设置健康状态为unknown，等待下次检测
		updates["health_status"] = "unknown"
	}

	return common.GetDB().Model(&model.APIMonitor{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// UpdateHealthStatus 更新健康状态
func (d *APIMonitorDaoImpl) UpdateHealthStatus(id uint, healthStatus string) error {
	return common.GetDB().Model(&model.APIMonitor{}).
		Where("id = ?", id).
		Update("health_status", healthStatus).Error
}

// UpdateLastCheckResult 更新最近检测结果
func (d *APIMonitorDaoImpl) UpdateLastCheckResult(monitor *model.APIMonitor) error {
	return common.GetDB().Model(&model.APIMonitor{}).
		Where("id = ?", monitor.ID).
		Updates(map[string]interface{}{
			"last_check_time":      monitor.LastCheckTime,
			"last_response_time":   monitor.LastResponseTime,
			"last_status_code":     monitor.LastStatusCode,
			"last_error_message":   monitor.LastErrorMessage,
			"continuous_failures":  monitor.ContinuousFailures,
			"health_status":        monitor.HealthStatus,
		}).Error
}

// GetStats 获取统计信息
func (d *APIMonitorDaoImpl) GetStats() (*model.APIMonitorStats, error) {
	var stats model.APIMonitorStats

	db := common.GetDB().Model(&model.APIMonitor{})

	// 总数
	var totalCount int64
	db.Count(&totalCount)
	stats.TotalCount = int(totalCount)

	// 启用数
	var enabledCount int64
	db.Where("status = ?", 1).Count(&enabledCount)
	stats.EnabledCount = int(enabledCount)

	// 健康数
	var healthyCount int64
	db.Where("health_status = ?", "healthy").Count(&healthyCount)
	stats.HealthyCount = int(healthyCount)

	// 异常数
	var unhealthyCount int64
	db.Where("health_status = ?", "unhealthy").Count(&unhealthyCount)
	stats.UnhealthyCount = int(unhealthyCount)

	// 平均响应时间
	var avgResponseTime float64
	common.GetDB().Model(&model.APIMonitor{}).
		Select("AVG(last_response_time)").
		Where("last_response_time > 0").
		Scan(&avgResponseTime)
	stats.AvgResponseTime = avgResponseTime

	return &stats, nil
}

// GetGroupStats 获取按业务组聚合的统计信息（包含父子关系，显示所有业务组）
func (d *APIMonitorDaoImpl) GetGroupStats() ([]*model.APIMonitorGroupStats, error) {
	type groupResult struct {
		GroupID         uint
		GroupName       string
		ParentID        uint
		TotalAPIs       int
		HealthyAPIs     int
		UnhealthyAPIs   int
		AvgResponseTime float64
	}

	var results []groupResult
	// 从cmdb_group表获取所有业务组，左联monitor_api_monitor表统计数据
	err := common.GetDB().Table("cmdb_group cg").
		Select(`
			cg.id as group_id,
			cg.name as group_name,
			COALESCE(cg.parent_id, 0) as parent_id,
			COUNT(am.id) as total_apis,
			SUM(CASE WHEN am.health_status = 'healthy' THEN 1 ELSE 0 END) as healthy_apis,
			SUM(CASE WHEN am.health_status = 'unhealthy' THEN 1 ELSE 0 END) as unhealthy_apis,
			AVG(CASE WHEN am.last_response_time > 0 THEN am.last_response_time ELSE NULL END) as avg_response_time
		`).
		Joins("LEFT JOIN monitor_api_monitor am ON cg.id = am.group_id").
		Group("cg.id, cg.name, cg.parent_id").
		Order("cg.id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为返回格式
	var stats []*model.APIMonitorGroupStats
	for _, r := range results {
		stats = append(stats, &model.APIMonitorGroupStats{
			GroupID:         r.GroupID,
			GroupName:       r.GroupName,
			ParentID:        r.ParentID,
			TotalAPIs:       r.TotalAPIs,
			HealthyAPIs:     r.HealthyAPIs,
			UnhealthyAPIs:   r.UnhealthyAPIs,
			AvgResponseTime: r.AvgResponseTime,
		})
	}

	return stats, nil
}

// CheckNameExistsInGroup 检查同一业务组内是否存在同名监控
func (d *APIMonitorDaoImpl) CheckNameExistsInGroup(groupID uint, name string, excludeID uint) (bool, error) {
	var count int64
	query := common.GetDB().Model(&model.APIMonitor{}).
		Where("group_id = ? AND name = ?", groupID, name)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// GetMonitorTree 获取父监控及其所有子监控（树形结构）
func (d *APIMonitorDaoImpl) GetMonitorTree(parentID uint) (*model.APIMonitor, error) {
	var monitor model.APIMonitor
	err := common.GetDB().Where("id = ?", parentID).
		Preload("ChildMonitors").
		First(&monitor).Error
	return &monitor, err
}

// GetTopLevelMonitors 获取所有顶级监控（parentID=null）
func (d *APIMonitorDaoImpl) GetTopLevelMonitors(query *model.APIMonitorListQuery) ([]*model.APIMonitor, int64, error) {
	query.OnlyTopLevel = true
	return d.GetList(query)
}

// GetChildMonitors 获取指定监控的所有直接子监控
func (d *APIMonitorDaoImpl) GetChildMonitors(parentID uint) ([]*model.APIMonitor, error) {
	var monitors []*model.APIMonitor
	log.Printf("🔍 [GetChildMonitors] 查询子监控, parentID=%d", parentID)

	db := common.GetDB().Where("parent_id = ?", parentID).Order("id ASC")

	// 先执行 count 查看是否有记录
	var count int64
	db.Model(&model.APIMonitor{}).Count(&count)
	log.Printf("🔍 [GetChildMonitors] Count查询结果: %d 条", count)

	// 执行实际查询
	err := db.Find(&monitors).Error
	log.Printf("🔍 [GetChildMonitors] Find查询结果: %d 条, error=%v", len(monitors), err)

	if len(monitors) > 0 {
		for _, m := range monitors {
			log.Printf("🔍 [GetChildMonitors] 子监控: ID=%d, Name=%s, ParentID=%v", m.ID, m.Name, m.ParentID)
		}
	}

	return monitors, err
}

// GetAllDescendants 获取指定监控的所有后代监控（递归）
func (d *APIMonitorDaoImpl) GetAllDescendants(parentID uint) ([]*model.APIMonitor, error) {
	var descendants []*model.APIMonitor

	// 获取直接子监控
	var children []*model.APIMonitor
	err := common.GetDB().Where("parent_id = ?", parentID).Find(&children).Error
	if err != nil {
		return nil, err
	}

	descendants = append(descendants, children...)

	// 递归获取孙监控及更深层级
	for _, child := range children {
		grandChildren, err := d.GetAllDescendants(child.ID)
		if err != nil {
			return nil, err
		}
		descendants = append(descendants, grandChildren...)
	}

	return descendants, nil
}

// CheckCircularDependency 检查是否会形成循环依赖（childID作为parentID时，parentID不应该是childID的后代）
func (d *APIMonitorDaoImpl) CheckCircularDependency(childID uint, parentID uint) (bool, error) {
	if childID == parentID {
		return true, nil // 自引用
	}

	// 获取parentID的所有后代，检查是否包含childID
	descendants, err := d.GetAllDescendants(parentID)
	if err != nil {
		return false, err
	}

	for _, desc := range descendants {
		if desc.ID == childID {
			return true, nil // 存在循环依赖
		}
	}

	return false, nil
}

// GetChildrenStatsInBatch 批量获取子监控的统计信息
func (d *APIMonitorDaoImpl) GetChildrenStatsInBatch(parentIDs []uint) (map[uint]*model.ChildrenStats, error) {
	if len(parentIDs) == 0 {
		return make(map[uint]*model.ChildrenStats), nil
	}

	type statsResult struct {
		ParentID        uint
		ChildCount      int
		HealthyChildren int
		HealthStatus    string
	}

	var results []statsResult

	// 查询所有子监控，按parentID分组统计
	err := common.GetDB().Model(&model.APIMonitor{}).
		Select(`
			parent_id as parent_id,
			COUNT(*) as child_count,
			SUM(CASE WHEN health_status = 'healthy' THEN 1 ELSE 0 END) as healthy_children,
			CASE
				WHEN COUNT(*) > 0 AND SUM(CASE WHEN health_status = 'healthy' THEN 1 ELSE 0 END) = COUNT(*)
				THEN 'healthy'
				ELSE 'unhealthy'
			END as health_status
		`).
		Where("parent_id IN ?", parentIDs).
		Group("parent_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为map
	statsMap := make(map[uint]*model.ChildrenStats)
	for _, r := range results {
		statsMap[r.ParentID] = &model.ChildrenStats{
			ParentID:        r.ParentID,
			ChildCount:      r.ChildCount,
			HealthyChildren: r.HealthyChildren,
			HealthStatus:    r.HealthStatus,
		}
	}

	// 确保所有parentID都在结果中（即使没有子监控）
	for _, parentID := range parentIDs {
		if _, exists := statsMap[parentID]; !exists {
			statsMap[parentID] = &model.ChildrenStats{
				ParentID:        parentID,
				ChildCount:      0,
				HealthyChildren: 0,
				HealthStatus:    "healthy", // 没有子监控时默认为healthy
			}
		}
	}

	return statsMap, nil
}

// UpdateParentAggregateStats 更新父监控的聚合字段（ChildCount、HealthyChildren等）
func (d *APIMonitorDaoImpl) UpdateParentAggregateStats(parentID uint) error {
	// 获取所有子监控的统计
	children, err := d.GetChildMonitors(parentID)
	if err != nil {
		return err
	}

	childCount := len(children)
	healthyChildren := 0
	for _, child := range children {
		if child.HealthStatus == "healthy" {
			healthyChildren++
		}
	}

	// 计算聚合健康状态
	aggregateHealthStatus := "healthy"
	if childCount > 0 && healthyChildren < childCount {
		aggregateHealthStatus = "unhealthy"
	}

	// 更新父监控记录
	return common.GetDB().Model(&model.APIMonitor{}).
		Where("id = ?", parentID).
		Updates(map[string]interface{}{
			"health_status": aggregateHealthStatus,
		}).Error
}

// GetMaxHierarchyDepth 获取指定监控的最大层级深度（0为根级，1为一级子，2为二级子...）
func (d *APIMonitorDaoImpl) GetMaxHierarchyDepth(monitorID uint) (int, error) {
	var maxDepth int
	err := common.GetDB().Model(&model.APIMonitor{}).
		Where("id IN (?)",
			common.GetDB().Model(&model.APIMonitor{}).
				Select("id").
				Where("parent_id IN (?)",
					common.GetDB().Model(&model.APIMonitor{}).
						Select("id").
						Where("id = ? OR parent_id = ?", monitorID, monitorID),
				),
		).
		Select("MAX(hierarchy_level)").
		Scan(&maxDepth).Error

	// 简化实现：直接查询整条链的最大深度
	var monitor model.APIMonitor
	if err := common.GetDB().Where("id = ?", monitorID).First(&monitor).Error; err != nil {
		return 0, err
	}

	maxDepth = monitor.HierarchyLevel

	// 获取所有后代中的最大深度
	descendants, err := d.GetAllDescendants(monitorID)
	if err != nil {
		return maxDepth, nil
	}

	for _, desc := range descendants {
		if desc.HierarchyLevel > maxDepth {
			maxDepth = desc.HierarchyLevel
		}
	}

	return maxDepth, nil
}
