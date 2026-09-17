package controller

import (
	"dodevops-api/api/system/model"
	"dodevops-api/api/system/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// @Tags 激活管理
// @Summary 获取部署机器码
// @Description 获取当前部署实例的机器码，用于生成激活码
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/license/machine-code [get]
// @Security ApiKeyAuth
func GetLicenseMachineCode(c *gin.Context) {
	service.NewLicenseService().GetMachineCode(c)
}

// @Tags 激活管理
// @Summary 获取激活状态
// @Description 获取当前系统是否已激活、到期时间和剩余天数
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/license/status [get]
// @Security ApiKeyAuth
func GetLicenseStatus(c *gin.Context) {
	service.NewLicenseService().GetLicenseStatus(c)
}

// @Tags 激活管理
// @Summary 提交激活码
// @Description 提交离线签名的激活码并完成系统激活
// @Accept json
// @Produce json
// @Param data body model.ActivateLicenseRequest true "激活请求"
// @Success 200 {object} result.Result
// @Router /api/v1/license/activate [post]
// @Security ApiKeyAuth
func ActivateLicense(c *gin.Context) {
	var req model.ActivateLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}
	service.NewLicenseService().Activate(c, &req)
}
