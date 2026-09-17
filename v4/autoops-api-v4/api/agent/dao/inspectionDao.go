// 巡检任务数据访问层
package dao

import (
	"encoding/json"
	"time"

	"dodevops-api/api/agent/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"
)

// InspectionDao 巡检任务DAO接口
type InspectionDao interface {
	Create(inspection *model.AgentInspection) error
	GetByID(id uint) (*model.AgentInspection, error)
	Update(inspection *model.AgentInspection) error
	Delete(id uint) error
	GetList(query *model.InspectionQuery) ([]model.AgentInspection, int64, error)
	GetEnabledInspections() ([]model.AgentInspection, error)
	GetScheduledInspections() ([]model.AgentInspection, error)
}

// InspectionDaoImpl 巡检任务DAO实现
type InspectionDaoImpl struct{}

// NewInspectionDao 创建巡检任务DAO实例
func NewInspectionDao() InspectionDao {
	return &InspectionDaoImpl{}
}

// Create 创建巡检任务
func (d *InspectionDaoImpl) Create(inspection *model.AgentInspection) error {
	now := util.HTime{Time: time.Now()}
	inspection.CreateTime = now
	inspection.UpdateTime = now
	return Db.Create(inspection).Error
}

// GetByID 根据ID获取巡检任务
func (d *InspectionDaoImpl) GetByID(id uint) (*model.AgentInspection, error) {
	var inspection model.AgentInspection
	err := Db.Where("id = ?", id).First(&inspection).Error
	if err != nil {
		return nil, err
	}
	return &inspection, nil
}

// Update 更新巡检任务
func (d *InspectionDaoImpl) Update(inspection *model.AgentInspection) error {
	inspection.UpdateTime = util.HTime{Time: time.Now()}
	return Db.Model(&model.AgentInspection{}).Where("id = ?", inspection.ID).Updates(inspection).Error
}

// Delete 删除巡检任务
func (d *InspectionDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.AgentInspection{}, id).Error
}

// GetList 获取巡检任务列表
func (d *InspectionDaoImpl) GetList(query *model.InspectionQuery) ([]model.AgentInspection, int64, error) {
	var inspections []model.AgentInspection
	var total int64

	db := Db.Model(&model.AgentInspection{})

	// 条件筛选
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.TargetType != "" {
		db = db.Where("target_type = ?", query.TargetType)
	}
	if query.ScheduleType != "" {
		db = db.Where("schedule_type = ?", query.ScheduleType)
	}
	if query.Enabled > 0 {
		db = db.Where("enabled = ?", query.Enabled)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if query.PageNum > 0 && query.PageSize > 0 {
		offset := (query.PageNum - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 按创建时间倒序
	db = db.Order("create_time DESC")

	if err := db.Find(&inspections).Error; err != nil {
		return nil, 0, err
	}

	return inspections, total, nil
}

// GetEnabledInspections 获取所有启用的巡检任务
func (d *InspectionDaoImpl) GetEnabledInspections() ([]model.AgentInspection, error) {
	var inspections []model.AgentInspection
	err := Db.Where("enabled = ?", 1).Find(&inspections).Error
	return inspections, err
}

// GetScheduledInspections 获取所有定时巡检任务
func (d *InspectionDaoImpl) GetScheduledInspections() ([]model.AgentInspection, error) {
	var inspections []model.AgentInspection
	err := Db.Where("enabled = ? AND schedule_type = ?", 1, "cron").Find(&inspections).Error
	return inspections, err
}

// InspectionRecordDao 巡检记录DAO接口
type InspectionRecordDao interface {
	Create(record *model.InspectionRecord) error
	GetByID(id uint) (*model.InspectionRecord, error)
	Update(record *model.InspectionRecord) error
	Delete(id uint) error
	GetList(query *model.InspectionRecordQuery) ([]model.InspectionRecord, int64, error)
	GetLatestByHost(hostID uint, limit int) ([]model.InspectionRecord, error)
	GetStatistics(inspectionID uint) (map[string]interface{}, error)
}

// InspectionRecordDaoImpl 巡检记录DAO实现
type InspectionRecordDaoImpl struct{}

// NewInspectionRecordDao 创建巡检记录DAO实例
func NewInspectionRecordDao() InspectionRecordDao {
	return &InspectionRecordDaoImpl{}
}

// Create 创建巡检记录
func (d *InspectionRecordDaoImpl) Create(record *model.InspectionRecord) error {
	record.CreateTime = util.HTime{Time: time.Now()}
	return Db.Create(record).Error
}

// GetByID 根据ID获取巡检记录
func (d *InspectionRecordDaoImpl) GetByID(id uint) (*model.InspectionRecord, error) {
	var record model.InspectionRecord
	err := Db.Where("id = ?", id).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// Update 更新巡检记录
func (d *InspectionRecordDaoImpl) Update(record *model.InspectionRecord) error {
	return Db.Model(&model.InspectionRecord{}).Where("id = ?", record.ID).Updates(record).Error
}

// Delete 删除巡检记录
func (d *InspectionRecordDaoImpl) Delete(id uint) error {
	// 同时删除关联的巡检项
	Db.Where("record_id = ?", id).Delete(&model.InspectionItem{})
	return Db.Delete(&model.InspectionRecord{}, id).Error
}

// GetList 获取巡检记录列表
func (d *InspectionRecordDaoImpl) GetList(query *model.InspectionRecordQuery) ([]model.InspectionRecord, int64, error) {
	var records []model.InspectionRecord
	var total int64

	db := Db.Model(&model.InspectionRecord{})

	// 条件筛选
	if query.InspectionID > 0 {
		db = db.Where("inspection_id = ?", query.InspectionID)
	}
	if query.HostID > 0 {
		db = db.Where("host_id = ?", query.HostID)
	}
	if query.HostName != "" {
		db = db.Where("host_name LIKE ?", "%"+query.HostName+"%")
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.RiskLevel >= 0 {
		db = db.Where("risk_level = ?", query.RiskLevel)
	}
	if query.StartTime != "" {
		db = db.Where("start_time >= ?", query.StartTime)
	}
	if query.EndTime != "" {
		db = db.Where("start_time <= ?", query.EndTime)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if query.PageNum > 0 && query.PageSize > 0 {
		offset := (query.PageNum - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 按开始时间倒序
	db = db.Order("start_time DESC")

	if err := db.Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetLatestByHost 获取主机最近的巡检记录
func (d *InspectionRecordDaoImpl) GetLatestByHost(hostID uint, limit int) ([]model.InspectionRecord, error) {
	var records []model.InspectionRecord
	err := Db.Where("host_id = ?", hostID).
		Order("start_time DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

// GetStatistics 获取巡检统计
func (d *InspectionRecordDaoImpl) GetStatistics(inspectionID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 按状态统计
	var statusStats []struct {
		Status int
		Count  int64
	}
	db := Db.Model(&model.InspectionRecord{})
	if inspectionID > 0 {
		db = db.Where("inspection_id = ?", inspectionID)
	}
	err := db.Select("status, count(*) as count").
		Group("status").
		Scan(&statusStats).Error
	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]int64)
	for _, stat := range statusStats {
		switch stat.Status {
		case 1:
			statusMap["running"] = stat.Count
		case 2:
			statusMap["success"] = stat.Count
		case 3:
			statusMap["failed"] = stat.Count
		case 4:
			statusMap["partial"] = stat.Count
		}
	}
	stats["by_status"] = statusMap

	// 按风险等级统计
	var riskStats []struct {
		RiskLevel int
		Count     int64
	}
	db = Db.Model(&model.InspectionRecord{})
	if inspectionID > 0 {
		db = db.Where("inspection_id = ?", inspectionID)
	}
	err = db.Select("risk_level, count(*) as count").
		Group("risk_level").
		Scan(&riskStats).Error
	if err != nil {
		return nil, err
	}

	riskMap := make(map[string]int64)
	for _, stat := range riskStats {
		switch stat.RiskLevel {
		case 0:
			riskMap["normal"] = stat.Count
		case 1:
			riskMap["low"] = stat.Count
		case 2:
			riskMap["medium"] = stat.Count
		case 3:
			riskMap["high"] = stat.Count
		}
	}
	stats["by_risk"] = riskMap

	return stats, nil
}

// InspectionItemDao 巡检项DAO接口
type InspectionItemDao interface {
	Create(item *model.InspectionItem) error
	BatchCreate(items []model.InspectionItem) error
	GetByRecordID(recordID uint) ([]model.InspectionItem, error)
	Delete(id uint) error
	DeleteByRecordID(recordID uint) error
}

// InspectionItemDaoImpl 巡检项DAO实现
type InspectionItemDaoImpl struct{}

// NewInspectionItemDao 创建巡检项DAO实例
func NewInspectionItemDao() InspectionItemDao {
	return &InspectionItemDaoImpl{}
}

// Create 创建巡检项
func (d *InspectionItemDaoImpl) Create(item *model.InspectionItem) error {
	item.CreateTime = util.HTime{Time: time.Now()}
	return Db.Create(item).Error
}

// BatchCreate 批量创建巡检项
func (d *InspectionItemDaoImpl) BatchCreate(items []model.InspectionItem) error {
	now := util.HTime{Time: time.Now()}
	for i := range items {
		items[i].CreateTime = now
	}
	return Db.Create(&items).Error
}

// GetByRecordID 根据记录ID获取巡检项列表
func (d *InspectionItemDaoImpl) GetByRecordID(recordID uint) ([]model.InspectionItem, error) {
	var items []model.InspectionItem
	err := Db.Where("record_id = ?", recordID).
		Order("category ASC, id ASC").
		Find(&items).Error
	return items, err
}

// Delete 删除巡检项
func (d *InspectionItemDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.InspectionItem{}, id).Error
}

// DeleteByRecordID 根据记录ID删除所有巡检项
func (d *InspectionItemDaoImpl) DeleteByRecordID(recordID uint) error {
	return Db.Where("record_id = ?", recordID).Delete(&model.InspectionItem{}).Error
}

// ParseTargetIDs 解析目标ID列表
func ParseTargetIDs(targetIDsJSON string) ([]uint, error) {
	var ids []uint
	if err := json.Unmarshal([]byte(targetIDsJSON), &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// EncodeTargetIDs 编码目标ID列表
func EncodeTargetIDs(ids []uint) (string, error) {
	data, err := json.Marshal(ids)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
