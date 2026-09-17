package dao

import (
	"dodevops-api/api/app/model"
	"gorm.io/gorm"
)

// IServiceReleaseDao 服务上线数据访问接口
type IServiceReleaseDao interface {
	// 创建服务上线工单
	CreateServiceRelease(release *model.ServiceRelease) error
	// 获取服务上线工单详情
	GetServiceReleaseByID(id uint) (*model.ServiceRelease, error)
	// 获取服务上线工单列表
	GetServiceReleaseList(req *model.ServiceReleaseListRequest) ([]model.ServiceRelease, int64, error)
	// 更新服务上线工单
	UpdateServiceRelease(id uint, updates map[string]interface{}) error
	// 删除服务上线工单
	DeleteServiceRelease(id uint) error

	// 创建服务明细
	CreateServiceReleaseItem(item *model.ServiceReleaseItem) error
	// 获取服务明细列表
	GetServiceReleaseItems(releaseID uint) ([]model.ServiceReleaseItem, error)
	// 更新服务明细
	UpdateServiceReleaseItem(id uint, updates map[string]interface{}) error
	// 删除服务明细
	DeleteServiceReleaseItem(id uint) error
}

// ServiceReleaseDao 服务上线数据访问实现
type ServiceReleaseDao struct {
	db *gorm.DB
}

// NewServiceReleaseDao 创建服务上线数据访问对象
func NewServiceReleaseDao(db *gorm.DB) IServiceReleaseDao {
	return &ServiceReleaseDao{db: db}
}

// CreateServiceRelease 创建服务上线工单
func (d *ServiceReleaseDao) CreateServiceRelease(release *model.ServiceRelease) error {
	return d.db.Create(release).Error
}

// GetServiceReleaseByID 获取服务上线工单详情
func (d *ServiceReleaseDao) GetServiceReleaseByID(id uint) (*model.ServiceRelease, error) {
	var release model.ServiceRelease
	err := d.db.Preload("Services.Application").
		Preload("Services.JenkinsEnv").
		First(&release, id).Error
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// GetServiceReleaseList 获取服务上线工单列表
func (d *ServiceReleaseDao) GetServiceReleaseList(req *model.ServiceReleaseListRequest) ([]model.ServiceRelease, int64, error) {
	var releases []model.ServiceRelease
	var total int64

	query := d.db.Model(&model.ServiceRelease{})

	// 构建查询条件
	if req.BusinessGroupID != nil {
		query = query.Where("business_group_id = ?", *req.BusinessGroupID)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.ApplicantID != nil {
		query = query.Where("applicant_id = ?", *req.ApplicantID)
	}
	if req.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Services").
		Offset(offset).
		Limit(req.PageSize).
		Order("created_at DESC").
		Find(&releases).Error

	return releases, total, err
}

// UpdateServiceRelease 更新服务上线工单
func (d *ServiceReleaseDao) UpdateServiceRelease(id uint, updates map[string]interface{}) error {
	return d.db.Model(&model.ServiceRelease{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteServiceRelease 删除服务上线工单
func (d *ServiceReleaseDao) DeleteServiceRelease(id uint) error {
	return d.db.Delete(&model.ServiceRelease{}, id).Error
}

// CreateServiceReleaseItem 创建服务明细
func (d *ServiceReleaseDao) CreateServiceReleaseItem(item *model.ServiceReleaseItem) error {
	return d.db.Create(item).Error
}

// GetServiceReleaseItems 获取服务明细列表
func (d *ServiceReleaseDao) GetServiceReleaseItems(releaseID uint) ([]model.ServiceReleaseItem, error) {
	var items []model.ServiceReleaseItem
	err := d.db.Where("release_id = ?", releaseID).
		Preload("Application").
		Preload("JenkinsEnv").
		Order("execute_order ASC").
		Find(&items).Error
	return items, err
}

// UpdateServiceReleaseItem 更新服务明细
func (d *ServiceReleaseDao) UpdateServiceReleaseItem(id uint, updates map[string]interface{}) error {
	return d.db.Model(&model.ServiceReleaseItem{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteServiceReleaseItem 删除服务明细
func (d *ServiceReleaseDao) DeleteServiceReleaseItem(id uint) error {
	return d.db.Delete(&model.ServiceReleaseItem{}, id).Error
}
