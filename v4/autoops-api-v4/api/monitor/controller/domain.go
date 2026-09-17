package controller

import (
	"strconv"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// DomainController 域名监控控制器
type DomainController struct {
	domainService service.DomainService
}

// NewDomainController 创建域名监控控制器实例
func NewDomainController() *DomainController {
	return &DomainController{
		domainService: service.NewDomainService(),
	}
}

// @Tags 域名监控
// @Summary 创建域名监控
// @Description 创建域名监控
// @Accept json
// @Produce json
// @Description 添加新的域名监控
// @Param data body model.DomainDto true "域名信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain [post]
// @Security ApiKeyAuth
func (ctrl *DomainController) CreateDomain(c *gin.Context) {
	var dto model.DomainDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	ctrl.domainService.Create(c, &dto)
}

// @Tags 域名监控
// @Summary 获取域名监控详情
// @Description 获取域名监控详情
// @Produce json
// @Description 根据ID获取域名监控详情
// @Param id path int true "域名ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/{id} [get]
// @Security ApiKeyAuth
func (ctrl *DomainController) GetDomain(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.domainService.GetByID(c, uint(id))
}

// @Tags 域名监控
// @Summary 更新域名监控
// @Description 更新域名监控
// @Accept json
// @Produce json
// @Description 更新域名监控信息
// @Param data body model.DomainDto true "域名信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain [put]
// @Security ApiKeyAuth
func (ctrl *DomainController) UpdateDomain(c *gin.Context) {
	var dto model.DomainDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	if dto.ID == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "ID不能为空")
		return
	}

	ctrl.domainService.Update(c, &dto)
}

// @Tags 域名监控
// @Summary 删除域名监控
// @Description 删除域名监控
// @Produce json
// @Description 根据ID删除域名监控
// @Param id path int true "域名ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *DomainController) DeleteDomain(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.domainService.Delete(c, uint(id))
}

// @Tags 域名监控
// @Summary 获取域名监控列表
// @Description 获取域名监控列表
// @Produce json
// @Description 获取域名监控列表（支持分页和筛选）
// @Param domain query string false "域名(模糊查询)"
// @Param tags query string false "标签(模糊查询)"
// @Param status query int false "状态(-1:全部,1:启用,0:禁用)"
// @Param isAlive query int false "存活状态(-1:全部,1:正常,0:异常)"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/list [get]
// @Security ApiKeyAuth
func (ctrl *DomainController) GetDomainList(c *gin.Context) {
	query := &model.DomainQuery{
		Domain:   c.Query("domain"),
		Tags:     c.Query("tags"),
		Status:   -1, // 默认查询全部
		IsAlive:  -1, // 默认查询全部
		PageNum:  1,
		PageSize: 10,
	}

	if pageNum, err := strconv.Atoi(c.Query("pageNum")); err == nil && pageNum > 0 {
		query.PageNum = pageNum
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize")); err == nil && pageSize > 0 {
		query.PageSize = pageSize
	}
	if status, err := strconv.Atoi(c.Query("status")); err == nil {
		query.Status = status
	}
	if isAlive, err := strconv.Atoi(c.Query("isAlive")); err == nil {
		query.IsAlive = isAlive
	}

	ctrl.domainService.GetList(c, query)
}

// @Tags 域名监控
// @Summary 检查单个域名
// @Description 检查单个域名
// @Produce json
// @Description 立即检查指定域名的存活状态和SSL证书
// @Param id path int true "域名ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/check/{id} [post]
// @Security ApiKeyAuth
func (ctrl *DomainController) CheckDomain(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.domainService.CheckDomain(c, uint(id))
}

// @Tags 域名监控
// @Summary 检查所有域名
// @Description 检查所有域名
// @Produce json
// @Description 立即检查所有启用的域名
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/check-all [post]
// @Security ApiKeyAuth
func (ctrl *DomainController) CheckAllDomains(c *gin.Context) {
	ctrl.domainService.CheckAllDomains(c)
}

// @Tags 域名监控
// @Summary 获取SSL证书即将过期的域名
// @Description 获取SSL证书即将过期的域名
// @Produce json
// @Description 获取SSL证书即将过期的域名列表
// @Param days query int false "天数，默认30天"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/ssl-expiring [get]
// @Security ApiKeyAuth
func (ctrl *DomainController) GetSSLExpiringSoon(c *gin.Context) {
	days := 30 // 默认30天
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 {
			days = d
		}
	}

	ctrl.domainService.GetSSLExpiringSoon(c, days)
}

// @Tags 域名监控
// @Summary 获取统计信息
// @Description 获取统计信息
// @Produce json
// @Description 获取域名监控统计信息（总数、正常、异常、SSL过期等）
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/statistics [get]
// @Security ApiKeyAuth
func (ctrl *DomainController) GetStatistics(c *gin.Context) {
	ctrl.domainService.GetStatistics(c)
}

// @Tags 域名监控
// @Summary 批量上传域名文件
// @Description 批量上传域名文件
// @Accept multipart/form-data
// @Produce json
// @Description 上传包含域名的文本文件，每行一个域名，自动去重并异步创建
// @Param file formData file true "域名文件(txt)"
// @Param tags formData string false "统一标签"
// @Param remark formData string false "统一备注"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/batch-upload [post]
// @Security ApiKeyAuth
func (ctrl *DomainController) BatchUploadDomains(c *gin.Context) {
	ctrl.domainService.BatchUploadDomains(c)
}

// @Tags 域名监控
// @Summary 从阿里云解析导入域名
// @Description 从阿里云解析导入域名
// @Accept json
// @Produce json
// @Description 通过阿里云密钥ID从阿里云解析自动导入所有域名,支持导入子域名
// @Param data body model.ImportFromAliyunDNSDto true "导入配置"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/domain/import-aliyun-dns [post]
// @Security ApiKeyAuth
func (ctrl *DomainController) ImportFromAliyunDNS(c *gin.Context) {
	var dto model.ImportFromAliyunDNSDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	ctrl.domainService.ImportFromAliyunDNS(c, &dto)
}
