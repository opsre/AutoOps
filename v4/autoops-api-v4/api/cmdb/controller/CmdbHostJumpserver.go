package controller

import (
	"dodevops-api/api/cmdb/model"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type CmdbHostJumpServerController struct {
	service service.CmdbHostJumpServerServiceInterface
}

func NewCmdbHostJumpServerController() *CmdbHostJumpServerController {
	return &CmdbHostJumpServerController{
		service: service.GetCmdbHostJumpServerService(),
	}
}

// SyncJumpServerHosts 同步JumpServer主机
// @Summary 同步JumpServer主机
// @Description 从JumpServer同步主机资产到CMDB,根据业务分组自动匹配二级分组
// @Tags CMDB资产管理
// @Accept json
// @Produce json
// @Param data body model.SyncJumpServerHostDto true "JumpServer配置信息"
// @Success 200 {object} result.Result
// @Router /api/v1/cmdb/hostjumpserversync [post]
// @Security ApiKeyAuth
func (ctrl *CmdbHostJumpServerController) SyncJumpServerHosts(ctx *gin.Context) {
	var dto model.SyncJumpServerHostDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}
	ctrl.service.SyncJumpServerHosts(ctx, &dto)
}
