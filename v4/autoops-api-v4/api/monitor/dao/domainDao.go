package dao

import (
	"time"

	"dodevops-api/api/monitor/model"
	"dodevops-api/common"

	"gorm.io/gorm"
)

// DomainDao 域名监控数据访问对象接口
type DomainDao interface {
	Create(domain *model.Domain) error
	BatchCreate(domains []*model.Domain) error
	GetByDomain(domain string) (*model.Domain, error)
	GetByDomains(domains []string) (map[string]*model.Domain, error)
	GetByID(id uint) (*model.Domain, error)
	Update(domain *model.Domain) error
	UpdateCheckResult(id uint, result *model.DomainCheckResult) error
	Delete(id uint) error
	GetList(query *model.DomainQuery) ([]model.Domain, int64, error)
	GetAll() ([]model.Domain, error)
	GetEnabledDomains() ([]model.Domain, error)
	GetSSLExpiringSoon(days int) ([]model.Domain, error)
	GetStatistics() (*model.DomainStatistics, error)
}

// DomainDaoImpl 域名监控数据访问对象实现
type DomainDaoImpl struct {
	db *gorm.DB
}

// NewDomainDao 创建Domain DAO实例
func NewDomainDao() DomainDao {
	return &DomainDaoImpl{
		db: common.GetDB(),
	}
}

// Create 创建域名监控记录
func (d *DomainDaoImpl) Create(domain *model.Domain) error {
	domain.CreateTime.Time = time.Now()
	domain.CreateTime.Valid = true
	domain.UpdateTime.Time = time.Now()
	domain.UpdateTime.Valid = true
	return d.db.Create(domain).Error
}

// BatchCreate 批量创建域名监控记录
func (d *DomainDaoImpl) BatchCreate(domains []*model.Domain) error {
	if len(domains) == 0 {
		return nil
	}

	now := time.Now()
	for _, domain := range domains {
		domain.CreateTime.Time = now
		domain.CreateTime.Valid = true
		domain.UpdateTime.Time = now
		domain.UpdateTime.Valid = true
	}

	// 使用批量插入，每批最多500条
	batchSize := 500
	return d.db.CreateInBatches(domains, batchSize).Error
}

// GetByDomains 批量查询域名是否存在
func (d *DomainDaoImpl) GetByDomains(domainNames []string) (map[string]*model.Domain, error) {
	if len(domainNames) == 0 {
		return make(map[string]*model.Domain), nil
	}

	var domains []model.Domain
	err := d.db.Where("domain IN ?", domainNames).Find(&domains).Error
	if err != nil {
		return nil, err
	}

	// 转换为map方便查找
	result := make(map[string]*model.Domain)
	for i := range domains {
		result[domains[i].Domain] = &domains[i]
	}

	return result, nil
}

// GetByID 根据ID获取域名监控记录
func (d *DomainDaoImpl) GetByID(id uint) (*model.Domain, error) {
	var domain model.Domain
	err := d.db.Where("id = ?", id).First(&domain).Error
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

// GetByDomain 根据域名获取域名监控记录
func (d *DomainDaoImpl) GetByDomain(domainName string) (*model.Domain, error) {
	var domain model.Domain
	err := d.db.Where("domain = ?", domainName).First(&domain).Error
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

// Update 更新域名监控记录
func (d *DomainDaoImpl) Update(domain *model.Domain) error {
	domain.UpdateTime.Time = time.Now()
	domain.UpdateTime.Valid = true
	return d.db.Save(domain).Error
}

// UpdateCheckResult 更新检查结果
func (d *DomainDaoImpl) UpdateCheckResult(id uint, result *model.DomainCheckResult) error {
	updates := map[string]interface{}{
		"is_alive":      result.IsAlive,
		"status_code":   result.StatusCode,
		"response_time": result.ResponseTime,
		"error_msg":     result.ErrorMsg,
		"last_check_at": time.Now(),
		"update_time":   time.Now(),
	}

	// 如果有SSL信息
	if result.SSLExpireAt != "" {
		updates["ssl_expire_at"] = result.SSLExpireAt
		updates["ssl_days_left"] = result.SSLDaysLeft
		updates["ssl_issuer"] = result.SSLIssuer
	}

	return d.db.Model(&model.Domain{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除域名监控记录
func (d *DomainDaoImpl) Delete(id uint) error {
	return d.db.Delete(&model.Domain{}, id).Error
}

// GetList 获取域名监控列表
func (d *DomainDaoImpl) GetList(query *model.DomainQuery) ([]model.Domain, int64, error) {
	var domains []model.Domain
	var total int64

	db := d.db.Model(&model.Domain{})

	// 域名模糊查询
	if query.Domain != "" {
		db = db.Where("domain LIKE ?", "%"+query.Domain+"%")
	}

	// 标签模糊查询
	if query.Tags != "" {
		db = db.Where("tags LIKE ?", "%"+query.Tags+"%")
	}

	// 状态查询
	if query.Status != -1 {
		db = db.Where("status = ?", query.Status)
	}

	// 存活状态查询
	if query.IsAlive != -1 {
		db = db.Where("is_alive = ?", query.IsAlive)
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if query.PageNum > 0 && query.PageSize > 0 {
		offset := (query.PageNum - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 按更新时间倒序
	err := db.Order("update_time DESC").Find(&domains).Error
	return domains, total, err
}

// GetAll 获取所有域名监控记录
func (d *DomainDaoImpl) GetAll() ([]model.Domain, error) {
	var domains []model.Domain
	err := d.db.Find(&domains).Error
	return domains, err
}

// GetEnabledDomains 获取启用的域名监控记录
func (d *DomainDaoImpl) GetEnabledDomains() ([]model.Domain, error) {
	var domains []model.Domain
	err := d.db.Where("status = ?", 1).Find(&domains).Error
	return domains, err
}

// GetSSLExpiringSoon 获取SSL证书即将过期的域名
func (d *DomainDaoImpl) GetSSLExpiringSoon(days int) ([]model.Domain, error) {
	var domains []model.Domain

	// 查询条件：
	// 1. SSL证书过期时间不为空
	// 2. SSL证书剩余天数 <= 指定天数
	// 3. SSL证书剩余天数 > 0 (未过期)
	err := d.db.Where("ssl_expire_at IS NOT NULL AND ssl_days_left <= ? AND ssl_days_left > 0", days).
		Order("ssl_days_left ASC").
		Find(&domains).Error

	return domains, err
}

// GetStatistics 获取域名监控统计
func (d *DomainDaoImpl) GetStatistics() (*model.DomainStatistics, error) {
	stats := &model.DomainStatistics{}

	// 总数
	d.db.Model(&model.Domain{}).Count(&stats.Total)

	// 正常数量
	d.db.Model(&model.Domain{}).Where("is_alive = ?", 1).Count(&stats.Alive)

	// 异常数量
	d.db.Model(&model.Domain{}).Where("is_alive = ?", 0).Count(&stats.Dead)

	// SSL证书已过期
	d.db.Model(&model.Domain{}).Where("ssl_expire_at IS NOT NULL AND ssl_expire_at < ?", time.Now()).Count(&stats.SSLExpired)

	// SSL证书即将过期(30天内)
	thirtyDaysLater := time.Now().AddDate(0, 0, 30)
	d.db.Model(&model.Domain{}).
		Where("ssl_expire_at IS NOT NULL AND ssl_expire_at BETWEEN ? AND ?", time.Now(), thirtyDaysLater).
		Count(&stats.SSLSoon)

	return stats, nil
}
