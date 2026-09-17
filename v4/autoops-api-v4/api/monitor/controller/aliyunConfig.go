package controller

import (
	"strconv"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// AliyunConfigController 阿里云配置控制器
type AliyunConfigController struct {
	aliyunConfigService service.AliyunConfigService
}

// NewAliyunConfigController 创建阿里云配置控制器实例
func NewAliyunConfigController() *AliyunConfigController {
	return &AliyunConfigController{
		aliyunConfigService: service.NewAliyunConfigService(),
	}
}

// @Tags 阿里云配置管理
// @Summary 创建阿里云配置
// @Description 创建阿里云配置
// @Accept json
// @Produce json
// @Description 添加新的阿里云AK/SK配置
// @Param data body model.CreateAliyunConfigDto true "配置信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/aliyun-config [post]
// @Security ApiKeyAuth
func (ctrl *AliyunConfigController) CreateAliyunConfig(c *gin.Context) {
	var dto model.CreateAliyunConfigDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	ctrl.aliyunConfigService.CreateAliyunConfig(c, &dto)
}

// @Tags 阿里云配置管理
// @Summary 更新阿里云配置
// @Description 更新阿里云配置
// @Accept json
// @Produce json
// @Description 更新阿里云配置信息
// @Param data body model.UpdateAliyunConfigDto true "配置信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/aliyun-config [put]
// @Security ApiKeyAuth
func (ctrl *AliyunConfigController) UpdateAliyunConfig(c *gin.Context) {
	var dto model.UpdateAliyunConfigDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	if dto.ID == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "ID不能为空")
		return
	}

	ctrl.aliyunConfigService.UpdateAliyunConfig(c, &dto)
}

// @Tags 阿里云配置管理
// @Summary 获取阿里云配置详情
// @Description 获取阿里云配置详情
// @Produce json
// @Description 根据ID获取阿里云配置详情
// @Param id path int true "配置ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/aliyun-config/{id} [get]
// @Security ApiKeyAuth
func (ctrl *AliyunConfigController) GetAliyunConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.aliyunConfigService.GetAliyunConfig(c, uint(id))
}

// @Tags 阿里云配置管理
// @Summary 获取阿里云配置列表
// @Description 获取阿里云配置列表
// @Produce json
// @Description 查询所有阿里云配置列表
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/aliyun-config/list [get]
// @Security ApiKeyAuth
func (ctrl *AliyunConfigController) ListAliyunConfigs(c *gin.Context) {
	ctrl.aliyunConfigService.ListAliyunConfigs(c)
}

// @Tags 阿里云配置管理
// @Summary 删除阿里云配置
// @Description 删除阿里云配置
// @Produce json
// @Description 删除阿里云配置记录
// @Param id path int true "配置ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/aliyun-config/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *AliyunConfigController) DeleteAliyunConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.aliyunConfigService.DeleteAliyunConfig(c, uint(id))
}
