// 故障分析数据访问层
package dao

import (
	"time"

	"dodevops-api/api/agent/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"
)

// IncidentDao 故障分析DAO接口
type IncidentDao interface {
	Create(incident *model.IncidentAnalysis) error
	GetByID(id uint) (*model.IncidentAnalysis, error)
	GetByIncidentID(incidentID string) (*model.IncidentAnalysis, error)
	Update(incident *model.IncidentAnalysis) error
	Delete(id uint) error
	GetList(query *model.IncidentQuery) ([]model.IncidentAnalysis, int64, error)
	GetOpenIncidents() ([]model.IncidentAnalysis, error)
	GetStatistics() (map[string]interface{}, error)
}

// IncidentDaoImpl 故障分析DAO实现
type IncidentDaoImpl struct{}

// NewIncidentDao 创建故障分析DAO实例
func NewIncidentDao() IncidentDao {
	return &IncidentDaoImpl{}
}

// Create 创建故障分析
func (d *IncidentDaoImpl) Create(incident *model.IncidentAnalysis) error {
	now := util.HTime{Time: time.Now()}
	incident.CreateTime = now
	incident.UpdateTime = now
	return Db.Create(incident).Error
}

// GetByID 根据ID获取故障分析
func (d *IncidentDaoImpl) GetByID(id uint) (*model.IncidentAnalysis, error) {
	var incident model.IncidentAnalysis
	err := Db.Where("id = ?", id).First(&incident).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

// GetByIncidentID 根据故障ID获取故障分析
func (d *IncidentDaoImpl) GetByIncidentID(incidentID string) (*model.IncidentAnalysis, error) {
	var incident model.IncidentAnalysis
	err := Db.Where("incident_id = ?", incidentID).First(&incident).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

// Update 更新故障分析
func (d *IncidentDaoImpl) Update(incident *model.IncidentAnalysis) error {
	incident.UpdateTime = util.HTime{Time: time.Now()}
	return Db.Model(&model.IncidentAnalysis{}).Where("id = ?", incident.ID).Updates(incident).Error
}

// Delete 删除故障分析
func (d *IncidentDaoImpl) Delete(id uint) error {
	// 同时删除关联的证据
	var incident model.IncidentAnalysis
	if err := Db.Where("id = ?", id).First(&incident).Error; err != nil {
		return err
	}
	Db.Where("incident_id = ?", incident.IncidentID).Delete(&model.IncidentEvidence{})
	return Db.Delete(&model.IncidentAnalysis{}, id).Error
}

// GetList 获取故障分析列表
func (d *IncidentDaoImpl) GetList(query *model.IncidentQuery) ([]model.IncidentAnalysis, int64, error) {
	var incidents []model.IncidentAnalysis
	var total int64

	db := Db.Model(&model.IncidentAnalysis{})

	// 条件筛选
	if query.IncidentID != "" {
		db = db.Where("incident_id = ?", query.IncidentID)
	}
	if query.HostID > 0 {
		db = db.Where("host_id = ?", query.HostID)
	}
	if query.HostName != "" {
		db = db.Where("host_name LIKE ?", "%"+query.HostName+"%")
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if query.Severity > 0 {
		db = db.Where("severity = ?", query.Severity)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.StartTime != "" {
		db = db.Where("detect_time >= ?", query.StartTime)
	}
	if query.EndTime != "" {
		db = db.Where("detect_time <= ?", query.EndTime)
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

	// 按检测时间倒序
	db = db.Order("detect_time DESC")

	if err := db.Find(&incidents).Error; err != nil {
		return nil, 0, err
	}

	return incidents, total, nil
}

// GetOpenIncidents 获取未关闭的故障
func (d *IncidentDaoImpl) GetOpenIncidents() ([]model.IncidentAnalysis, error) {
	var incidents []model.IncidentAnalysis
	err := Db.Where("status IN (?, ?, ?)", 1, 2, 3). // 分析中、已定位、修复中
							Order("severity DESC, detect_time ASC").
							Find(&incidents).Error
	return incidents, err
}

// GetStatistics 获取故障统计
func (d *IncidentDaoImpl) GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 按状态统计
	var statusStats []struct {
		Status int
		Count  int64
	}
	err := Db.Model(&model.IncidentAnalysis{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusStats).Error
	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]int64)
	for _, stat := range statusStats {
		switch stat.Status {
		case 1:
			statusMap["analyzing"] = stat.Count
		case 2:
			statusMap["located"] = stat.Count
		case 3:
			statusMap["fixing"] = stat.Count
		case 4:
			statusMap["fixed"] = stat.Count
		case 5:
			statusMap["closed"] = stat.Count
		}
	}
	stats["by_status"] = statusMap

	// 按类别统计
	var categoryStats []struct {
		Category string
		Count    int64
	}
	err = Db.Model(&model.IncidentAnalysis{}).
		Select("category, count(*) as count").
		Group("category").
		Scan(&categoryStats).Error
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string]int64)
	for _, stat := range categoryStats {
		categoryMap[stat.Category] = stat.Count
	}
	stats["by_category"] = categoryMap

	// 按严重程度统计
	var severityStats []struct {
		Severity int
		Count    int64
	}
	err = Db.Model(&model.IncidentAnalysis{}).
		Select("severity, count(*) as count").
		Group("severity").
		Scan(&severityStats).Error
	if err != nil {
		return nil, err
	}

	severityMap := make(map[string]int64)
	for _, stat := range severityStats {
		switch stat.Severity {
		case 1:
			severityMap["low"] = stat.Count
		case 2:
			severityMap["medium"] = stat.Count
		case 3:
			severityMap["high"] = stat.Count
		case 4:
			severityMap["critical"] = stat.Count
		}
	}
	stats["by_severity"] = severityMap

	// 今日故障数
	var todayCount int64
	Db.Model(&model.IncidentAnalysis{}).
		Where("DATE(detect_time) = CURDATE()").
		Count(&todayCount)
	stats["today_count"] = todayCount

	// 自动修复率
	var totalCount, autoFixedCount int64
	Db.Model(&model.IncidentAnalysis{}).Count(&totalCount)
	Db.Model(&model.IncidentAnalysis{}).Where("is_auto_fixed = ?", true).Count(&autoFixedCount)
	if totalCount > 0 {
		stats["auto_fix_rate"] = float64(autoFixedCount) / float64(totalCount) * 100
	} else {
		stats["auto_fix_rate"] = 0.0
	}

	return stats, nil
}

// IncidentEvidenceDao 故障证据DAO接口
type IncidentEvidenceDao interface {
	Create(evidence *model.IncidentEvidence) error
	BatchCreate(evidences []model.IncidentEvidence) error
	GetByIncidentID(incidentID string) ([]model.IncidentEvidence, error)
	Delete(id uint) error
	DeleteByIncidentID(incidentID string) error
}

// IncidentEvidenceDaoImpl 故障证据DAO实现
type IncidentEvidenceDaoImpl struct{}

// NewIncidentEvidenceDao 创建故障证据DAO实例
func NewIncidentEvidenceDao() IncidentEvidenceDao {
	return &IncidentEvidenceDaoImpl{}
}

// Create 创建故障证据
func (d *IncidentEvidenceDaoImpl) Create(evidence *model.IncidentEvidence) error {
	evidence.CreateTime = util.HTime{Time: time.Now()}
	return Db.Create(evidence).Error
}

// BatchCreate 批量创建故障证据
func (d *IncidentEvidenceDaoImpl) BatchCreate(evidences []model.IncidentEvidence) error {
	now := util.HTime{Time: time.Now()}
	for i := range evidences {
		evidences[i].CreateTime = now
	}
	return Db.Create(&evidences).Error
}

// GetByIncidentID 根据故障ID获取证据列表
func (d *IncidentEvidenceDaoImpl) GetByIncidentID(incidentID string) ([]model.IncidentEvidence, error) {
	var evidences []model.IncidentEvidence
	err := Db.Where("incident_id = ?", incidentID).
		Order("relevance DESC, collect_time ASC").
		Find(&evidences).Error
	return evidences, err
}

// Delete 删除故障证据
func (d *IncidentEvidenceDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.IncidentEvidence{}, id).Error
}

// DeleteByIncidentID 根据故障ID删除所有证据
func (d *IncidentEvidenceDaoImpl) DeleteByIncidentID(incidentID string) error {
	return Db.Where("incident_id = ?", incidentID).Delete(&model.IncidentEvidence{}).Error
}

// IncidentKnowledgeDao 故障知识库DAO接口
type IncidentKnowledgeDao interface {
	Create(knowledge *model.IncidentKnowledge) error
	GetByID(id uint) (*model.IncidentKnowledge, error)
	Update(knowledge *model.IncidentKnowledge) error
	Delete(id uint) error
	GetList(query *model.IncidentKnowledgeQuery) ([]model.IncidentKnowledge, int64, error)
	Search(keyword string, category string, limit int) ([]model.IncidentKnowledge, error)
	IncrementMatchCount(id uint) error
	UpdateSuccessRate(id uint, isSuccess bool) error
}

// IncidentKnowledgeDaoImpl 故障知识库DAO实现
type IncidentKnowledgeDaoImpl struct{}

// NewIncidentKnowledgeDao 创建故障知识库DAO实例
func NewIncidentKnowledgeDao() IncidentKnowledgeDao {
	return &IncidentKnowledgeDaoImpl{}
}

// Create 创建故障知识
func (d *IncidentKnowledgeDaoImpl) Create(knowledge *model.IncidentKnowledge) error {
	now := util.HTime{Time: time.Now()}
	knowledge.CreateTime = now
	knowledge.UpdateTime = now
	return Db.Create(knowledge).Error
}

// GetByID 根据ID获取故障知识
func (d *IncidentKnowledgeDaoImpl) GetByID(id uint) (*model.IncidentKnowledge, error) {
	var knowledge model.IncidentKnowledge
	err := Db.Where("id = ?", id).First(&knowledge).Error
	if err != nil {
		return nil, err
	}
	return &knowledge, nil
}

// Update 更新故障知识
func (d *IncidentKnowledgeDaoImpl) Update(knowledge *model.IncidentKnowledge) error {
	knowledge.UpdateTime = util.HTime{Time: time.Now()}
	return Db.Model(&model.IncidentKnowledge{}).Where("id = ?", knowledge.ID).Updates(knowledge).Error
}

// Delete 删除故障知识
func (d *IncidentKnowledgeDaoImpl) Delete(id uint) error {
	return Db.Delete(&model.IncidentKnowledge{}, id).Error
}

// GetList 获取故障知识列表
func (d *IncidentKnowledgeDaoImpl) GetList(query *model.IncidentKnowledgeQuery) ([]model.IncidentKnowledge, int64, error) {
	var knowledges []model.IncidentKnowledge
	var total int64

	db := Db.Model(&model.IncidentKnowledge{})

	// 条件筛选
	if query.Title != "" {
		db = db.Where("title LIKE ?", "%"+query.Title+"%")
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if query.Keyword != "" {
		db = db.Where("keywords LIKE ?", "%"+query.Keyword+"%")
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

	// 按匹配次数和成功率排序
	db = db.Order("match_count DESC, success_rate DESC")

	if err := db.Find(&knowledges).Error; err != nil {
		return nil, 0, err
	}

	return knowledges, total, nil
}

// Search 搜索故障知识
func (d *IncidentKnowledgeDaoImpl) Search(keyword string, category string, limit int) ([]model.IncidentKnowledge, error) {
	var knowledges []model.IncidentKnowledge

	db := Db.Where("enabled = ?", 1)

	if category != "" {
		db = db.Where("category = ?", category)
	}

	if keyword != "" {
		db = db.Where("title LIKE ? OR keywords LIKE ? OR root_cause LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := db.Order("match_count DESC, success_rate DESC").
		Limit(limit).
		Find(&knowledges).Error

	return knowledges, err
}

// IncrementMatchCount 增加匹配次数
func (d *IncidentKnowledgeDaoImpl) IncrementMatchCount(id uint) error {
	return Db.Model(&model.IncidentKnowledge{}).
		Where("id = ?", id).
		UpdateColumn("match_count", Db.Raw("match_count + 1")).Error
}

// UpdateSuccessRate 更新成功率
func (d *IncidentKnowledgeDaoImpl) UpdateSuccessRate(id uint, isSuccess bool) error {
	knowledge, err := d.GetByID(id)
	if err != nil {
		return err
	}

	// 简单的移动平均算法
	totalAttempts := float64(knowledge.MatchCount)
	if totalAttempts == 0 {
		totalAttempts = 1
	}

	currentSuccessCount := knowledge.SuccessRate * float64(knowledge.MatchCount) / 100
	if isSuccess {
		currentSuccessCount++
	}

	newSuccessRate := (currentSuccessCount / totalAttempts) * 100

	return Db.Model(&model.IncidentKnowledge{}).
		Where("id = ?", id).
		Update("success_rate", newSuccessRate).Error
}
