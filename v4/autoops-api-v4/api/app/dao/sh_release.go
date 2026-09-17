package dao

import (
	"dodevops-api/api/app/model"
	"gorm.io/gorm"
)

// IShReleaseDao 脚本上线数据访问接口
type IShReleaseDao interface {
	// 创建脚本上线工单
	CreateShRelease(release *model.ShRelease) error
	// 获取脚本上线工单详情
	GetShReleaseByID(id uint) (*model.ShRelease, error)
	// 获取脚本上线工单列表
	GetShReleaseList(req *model.ShReleaseListRequest) ([]model.ShRelease, int64, error)
	// 更新脚本上线工单
	UpdateShRelease(id uint, updates map[string]interface{}) error
	// 删除脚本上线工单
	DeleteShRelease(id uint) error
}

// ShReleaseDao 脚本上线数据访问实现
type ShReleaseDao struct {
	db *gorm.DB
}

// NewShReleaseDao 创建脚本上线数据访问对象
func NewShReleaseDao(db *gorm.DB) IShReleaseDao {
	return &ShReleaseDao{db: db}
}

// CreateShRelease 创建脚本上线工单
func (d *ShReleaseDao) CreateShRelease(release *model.ShRelease) error {
	return d.db.Create(release).Error
}

// GetShReleaseByID 获取脚本上线工单详情
func (d *ShReleaseDao) GetShReleaseByID(id uint) (*model.ShRelease, error) {
	var release model.ShRelease
	err := d.db.Preload("Application").
		Preload("JenkinsEnv").
		First(&release, id).Error
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// GetShReleaseList 获取脚本上线工单列表
func (d *ShReleaseDao) GetShReleaseList(req *model.ShReleaseListRequest) ([]model.ShRelease, int64, error) {
	var releases []model.ShRelease
	var total int64

	query := d.db.Model(&model.ShRelease{})

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
	err := query.Preload("Application").
		Preload("JenkinsEnv").
		Offset(offset).
		Limit(req.PageSize).
		Order("created_at DESC").
		Find(&releases).Error

	return releases, total, err
}

// UpdateShRelease 更新脚本上线工单
func (d *ShReleaseDao) UpdateShRelease(id uint, updates map[string]interface{}) error {
	return d.db.Model(&model.ShRelease{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteShRelease 删除脚本上线工单
func (d *ShReleaseDao) DeleteShRelease(id uint) error {
	return d.db.Delete(&model.ShRelease{}, id).Error
}
