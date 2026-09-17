package controller

import (
	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type CloudBillingController struct {
	cloudBillingService service.CloudBillingService
}

func NewCloudBillingController() *CloudBillingController {
	return &CloudBillingController{
		cloudBillingService: service.NewCloudBillingService(),
	}
}

// @Tags 云账单管理
// @Summary 获取腾讯云账单
// @Description 获取腾讯云账单
// @Accept json
// @Produce json
// @Description 根据密钥ID获取腾讯云账单（按业务线统计）
// @Param keyId query uint true "密钥ID"
// @Param month query string false "账期(格式: 2024-11, 为空则查询上月)"
// @Success 200 {object} result.Result{data=model.TencentBillingResponse}
// @Router /api/v1/monitor/billing/tencent [get]
// @Security ApiKeyAuth
func (ctrl *CloudBillingController) GetTencentBilling(c *gin.Context) {
	var req model.TencentBillingRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误")
		return
	}
	ctrl.cloudBillingService.GetTencentBilling(c, &req)
}

// @Tags 云账单管理
// @Summary 获取阿里云账单（按产品）
// @Description 获取阿里云账单（按产品）
// @Accept json
// @Produce json
// @Description 根据密钥ID获取阿里云账单（按产品统计）
// @Param keyId query uint true "密钥ID"
// @Param month query string false "账期(格式: 2024-11, 为空则查询上月)"
// @Success 200 {object} result.Result{data=model.AliyunProductBillingResponse}
// @Router /api/v1/monitor/billing/aliyun-product [get]
// @Security ApiKeyAuth
func (ctrl *CloudBillingController) GetAliyunProductBilling(c *gin.Context) {
	var req model.AliyunProductBillingRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误")
		return
	}
	ctrl.cloudBillingService.GetAliyunProductBilling(c, &req)
}

// @Tags 云账单管理
// @Summary 获取阿里云账单（按标签）
// @Description 获取阿里云账单（按标签）
// @Accept json
// @Produce json
// @Description 根据密钥ID获取阿里云账单（按标签bline统计）
// @Param keyId query uint true "密钥ID"
// @Param month query string false "账期(格式: 2024-11, 为空则查询上月)"
// @Success 200 {object} result.Result{data=model.AliyunTagBillingResponse}
// @Router /api/v1/monitor/billing/aliyun-tag [get]
// @Security ApiKeyAuth
func (ctrl *CloudBillingController) GetAliyunTagBilling(c *gin.Context) {
	var req model.AliyunTagBillingRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误")
		return
	}
	ctrl.cloudBillingService.GetAliyunTagBilling(c, &req)
}

// @Tags 云账单管理
// @Summary 获取AWS账单
// @Description 获取AWS账单
// @Accept json
// @Produce json
// @Description 根据密钥ID获取AWS账单（美元转换为人民币，按SERVICE统计）
// @Param keyId query uint true "密钥ID"
// @Param month query string false "账期(格式: 2024-11, 为空则查询上月)"
// @Success 200 {object} result.Result{data=model.AWSBillingResponse}
// @Router /api/v1/monitor/billing/aws [get]
// @Security ApiKeyAuth
func (ctrl *CloudBillingController) GetAWSBilling(c *gin.Context) {
	var req model.AWSBillingRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误")
		return
	}
	ctrl.cloudBillingService.GetAWSBilling(c, &req)
}

// @Tags 云账单管理
// @Summary 获取上月费用汇总及同环比
// @Description 获取上月费用汇总及同环比
// @Accept json
// @Produce json
// @Description 获取各云厂商上月费用汇总以及与上上月的环比
// @Param tencentKeyId query uint false "腾讯云密钥ID"
// @Param aliyunKeyId query uint false "阿里云密钥ID"
// @Param awsKeyId query uint false "AWS密钥ID"
// @Success 200 {object} result.Result{data=model.LastMonthSummaryResponse}
// @Router /api/v1/monitor/billing/last-month-summary [get]
// @Security ApiKeyAuth
func (ctrl *CloudBillingController) GetLastMonthSummary(c *gin.Context) {
	var req model.CloudSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误")
		return
	}
	ctrl.cloudBillingService.GetLastMonthSummary(c, &req)
}

// @Tags 云账单管理
// @Summary 获取最近12个月费用趋势
// @Description 获取最近12个月费用趋势
// @Accept json
// @Produce json
// @Description 获取各云厂商最近12个月的费用趋势数据（折线图）
// @Param tencentKeyId query uint false "腾讯云密钥ID"
// @Param aliyunKeyId query uint false "阿里云密钥ID"
// @Param awsKeyId query uint false "AWS密钥ID"
// @Success 200 {object} result.Result{data=model.Last12MonthsTrendResponse}
// @Router /api/v1/monitor/billing/last-12-months-trend [get]
// @Security ApiKeyAuth
func (ctrl *CloudBillingController) GetLast12MonthsTrend(c *gin.Context) {
	var req model.CloudSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误")
		return
	}
	ctrl.cloudBillingService.GetLast12MonthsTrend(c, &req)
}
