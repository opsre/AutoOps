package service

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AliyunConfigService 阿里云配置服务接口
type AliyunConfigService interface {
	CreateAliyunConfig(c *gin.Context, dto *model.CreateAliyunConfigDto)
	UpdateAliyunConfig(c *gin.Context, dto *model.UpdateAliyunConfigDto)
	GetAliyunConfig(c *gin.Context, id uint)
	ListAliyunConfigs(c *gin.Context)
	DeleteAliyunConfig(c *gin.Context, id uint)
}

// AliyunConfigServiceImpl 阿里云配置服务实现
type AliyunConfigServiceImpl struct{}

// NewAliyunConfigService 创建阿里云配置服务实例
func NewAliyunConfigService() AliyunConfigService {
	return &AliyunConfigServiceImpl{}
}

// CreateAliyunConfig 创建阿里云配置
func (s *AliyunConfigServiceImpl) CreateAliyunConfig(c *gin.Context, dto *model.CreateAliyunConfigDto) {
	db := common.GetDB()

	// 检查配置名称是否重复
	var count int64
	db.Model(&model.AliyunConfig{}).Where("name = ?", dto.Name).Count(&count)
	if count > 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "配置名称已存在")
		return
	}

	// 设置默认区域
	region := dto.Region
	if region == "" {
		region = "cn-hangzhou"
	}

	now := time.Now()
	aliyunConfig := model.AliyunConfig{
		Name:         dto.Name,
		AccessKey:    dto.AccessKey,
		AccessSecret: dto.AccessSecret,
		Region:       region,
		Email:        dto.Email,
		Phone:        dto.Phone,
		Username:     dto.Username,
		Status:       1,
		CreateTime:   util.HTime{Time: now},
		UpdateTime:   util.HTime{Time: now},
	}

	if err := db.Create(&aliyunConfig).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, aliyunConfig.ToVO(true))
}

// UpdateAliyunConfig 更新阿里云配置
func (s *AliyunConfigServiceImpl) UpdateAliyunConfig(c *gin.Context, dto *model.UpdateAliyunConfigDto) {
	db := common.GetDB()

	var aliyunConfig model.AliyunConfig
	if err := db.First(&aliyunConfig, dto.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "配置不存在")
			return
		}
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 检查名称是否与其他记录重复
	var count int64
	db.Model(&model.AliyunConfig{}).Where("name = ? AND id != ?", dto.Name, dto.ID).Count(&count)
	if count > 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "配置名称已存在")
		return
	}

	updates := map[string]interface{}{
		"name":          dto.Name,
		"access_key":    dto.AccessKey,
		"access_secret": dto.AccessSecret,
		"region":        dto.Region,
		"email":         dto.Email,
		"phone":         dto.Phone,
		"username":      dto.Username,
		"status":        dto.Status,
		"update_time":   time.Now(),
	}

	if err := db.Model(&aliyunConfig).Updates(updates).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	db.First(&aliyunConfig, dto.ID)
	result.Success(c, aliyunConfig.ToVO(true))
}

// GetAliyunConfig 获取阿里云配置详情
func (s *AliyunConfigServiceImpl) GetAliyunConfig(c *gin.Context, id uint) {
	db := common.GetDB()

	var aliyunConfig model.AliyunConfig
	if err := db.First(&aliyunConfig, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "配置不存在")
			return
		}
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, aliyunConfig.ToVO(false)) // 详情不脱敏
}

// ListAliyunConfigs 查询阿里云配置列表
func (s *AliyunConfigServiceImpl) ListAliyunConfigs(c *gin.Context) {
	db := common.GetDB()

	var configs []model.AliyunConfig
	if err := db.Order("create_time DESC").Find(&configs).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	vos := make([]*model.AliyunConfigVO, 0, len(configs))
	for i := range configs {
		vos = append(vos, configs[i].ToVO(true)) // 列表脱敏
	}

	result.Success(c, vos)
}

// DeleteAliyunConfig 删除阿里云配置
func (s *AliyunConfigServiceImpl) DeleteAliyunConfig(c *gin.Context, id uint) {
	db := common.GetDB()

	// 检查是否有关联的证书
	var count int64
	db.Model(&model.SslCert{}).Where("aliyun_config_id = ?", id).Count(&count)
	if count > 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "该配置下还有关联的证书,无法删除")
		return
	}

	if err := db.Delete(&model.AliyunConfig{}, id).Error; err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, "删除成功")
}
