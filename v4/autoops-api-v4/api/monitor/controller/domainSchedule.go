package controller

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// DomainScheduleController 域名监控定时任务控制器
type DomainScheduleController struct {
	scheduleService service.DomainScheduleService
}

// NewDomainScheduleController 创建定时任务控制器实例
func NewDomainScheduleController() *DomainScheduleController {
	return &DomainScheduleController{
		scheduleService: service.NewDomainScheduleService(),
	}
}

// @Tags 域名监控-定时任务
// @Summary 获取定时任务配置
// @Description 获取定时任务配置
// @Produce json
// @Description 获取域名监控定时任务配置信息
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/schedule/config [get]
// @Security ApiKeyAuth
func (ctrl *DomainScheduleController) GetScheduleConfig(c *gin.Context) {
	ctrl.scheduleService.GetScheduleConfig(c)
}

// @Tags 域名监控-定时任务
// @Summary 更新定时任务配置
// @Description 更新定时任务配置
// @Accept json
// @Produce json
// @Description 更新域名监控定时任务配置（开启/关闭，设置检查间隔）
// @Param data body model.DomainScheduleDto true "定时任务配置"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/schedule/config [post]
// @Security ApiKeyAuth
func (ctrl *DomainScheduleController) UpdateScheduleConfig(c *gin.Context) {
	var dto model.DomainScheduleDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	ctrl.scheduleService.UpdateScheduleConfig(c, &dto)
}

// @Tags 域名监控-定时任务
// @Summary 暂停定时任务
// @Description 暂停定时任务
// @Produce json
// @Description 暂停域名监控定时任务
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/schedule/pause [post]
// @Security ApiKeyAuth
func (ctrl *DomainScheduleController) PauseSchedule(c *gin.Context) {
	ctrl.scheduleService.PauseSchedule(c)
}

// @Tags 域名监控-定时任务
// @Summary 恢复定时任务
// @Description 恢复定时任务
// @Produce json
// @Description 恢复域名监控定时任务
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/schedule/resume [post]
// @Security ApiKeyAuth
func (ctrl *DomainScheduleController) ResumeSchedule(c *gin.Context) {
	ctrl.scheduleService.ResumeSchedule(c)
}

// @Tags 域名监控-定时任务
// @Summary 获取定时任务状态
// @Description 获取定时任务状态
// @Produce json
// @Description 获取域名监控定时任务运行状态
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/schedule/status [get]
// @Security ApiKeyAuth
func (ctrl *DomainScheduleController) GetScheduleStatus(c *gin.Context) {
	ctrl.scheduleService.GetScheduleStatus(c)
}
