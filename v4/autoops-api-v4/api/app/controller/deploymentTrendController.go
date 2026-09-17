package controller

import (
	"dodevops-api/api/app/service"
	"github.com/gin-gonic/gin"
)

// GetDeploymentTrend 获取上线趋势统计
// @Summary 获取上线趋势统计
// @Description 获取服务上线、脚本上线、快速发布的趋势统计数据，支持近一周、近一个月、近一年
// @Tags 应用管理
// @Accept json
// @Produce json
// @Param period query string true "时间周期: week=近一周, month=近一个月, year=近一年"
// @Success 200 {object} result.Result{data=model.DeploymentTrendResponse}
// @Router /api/v1/apps/deployment-trend [get]
func GetDeploymentTrend(c *gin.Context) {
	service.GetDeploymentTrend(c)
}
