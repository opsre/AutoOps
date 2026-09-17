package controller

import (
	"strconv"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

// SslCertController SSL证书控制器
type SslCertController struct {
	sslCertService service.SslCertService
}

// NewSslCertController 创建SSL证书控制器实例
func NewSslCertController() *SslCertController {
	return &SslCertController{
		sslCertService: service.NewSslCertService(),
	}
}

// @Tags SSL证书管理
// @Summary 申请SSL证书
// @Description 申请SSL证书
// @Accept json
// @Produce json
// @Description 申请阿里云免费SSL证书
// @Param data body model.CreateSslCertDto true "证书申请信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/ssl-cert/apply [post]
// @Security ApiKeyAuth
func (ctrl *SslCertController) ApplySslCert(c *gin.Context) {
	var dto model.CreateSslCertDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	ctrl.sslCertService.ApplySslCert(c, &dto)
}

// @Tags SSL证书管理
// @Summary 查询证书状态
// @Description 查询证书状态
// @Produce json
// @Description 查询证书申请状态和验证信息
// @Param id path int true "证书ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/ssl-cert/{id}/state [get]
// @Security ApiKeyAuth
func (ctrl *SslCertController) QueryCertState(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.sslCertService.QueryCertState(c, uint(id))
}

// @Tags SSL证书管理
// @Summary 下载证书
// @Description 下载证书
// @Produce json
// @Description 下载已签发的SSL证书
// @Param id path int true "证书ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/ssl-cert/{id}/download [get]
// @Security ApiKeyAuth
func (ctrl *SslCertController) DownloadCert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.sslCertService.DownloadCert(c, uint(id))
}

// @Tags SSL证书管理
// @Summary 获取SSL证书列表
// @Description 获取SSL证书列表
// @Produce json
// @Description 查询SSL证书列表（分页）
// @Param domain query string false "域名（模糊查询）"
// @Param aliyunConfigId query int false "阿里云配置ID"
// @Param status query int false "状态"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/ssl-cert/list [get]
// @Security ApiKeyAuth
func (ctrl *SslCertController) ListSslCerts(c *gin.Context) {
	var dto model.SslCertListDto

	dto.Domain = c.Query("domain")
	dto.Status, _ = strconv.Atoi(c.Query("status"))
	dto.Page, _ = strconv.Atoi(c.Query("page"))
	dto.PageSize, _ = strconv.Atoi(c.Query("pageSize"))

	if c.Query("aliyunConfigId") != "" {
		configID, _ := strconv.ParseUint(c.Query("aliyunConfigId"), 10, 64)
		dto.AliyunConfigID = uint(configID)
	}

	ctrl.sslCertService.ListSslCerts(c, &dto)
}

// @Tags SSL证书管理
// @Summary 删除SSL证书
// @Description 删除SSL证书
// @Produce json
// @Description 删除SSL证书记录
// @Param id path int true "证书ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/ssl-cert/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *SslCertController) DeleteSslCert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.sslCertService.DeleteSslCert(c, uint(id))
}

// @Tags SSL证书管理
// @Summary 添加DNS验证记录
// @Description 添加DNS验证记录
// @Produce json
// @Description 自动添加DNS TXT验证记录到阿里云DNS
// @Param id path int true "证书ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/ssl-cert/{id}/dns-validation [post]
// @Security ApiKeyAuth
func (ctrl *SslCertController) AddDnsValidation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.sslCertService.AddDnsValidation(c, uint(id))
}

// @Tags SSL证书管理
// @Summary 检查DNS验证状态
// @Description 检查DNS验证状态
// @Produce json
// @Description 检查DNS验证是否通过
// @Param id path int true "证书ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/ssl-cert/{id}/check-validation [get]
// @Security ApiKeyAuth
func (ctrl *SslCertController) CheckDnsValidation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "ID参数错误")
		return
	}

	ctrl.sslCertService.CheckDnsValidation(c, uint(id))
}

// @Tags SSL证书管理
// @Summary 部署SSL证书（异步）
// @Description 部署SSL证书（异步）
// @Accept json
// @Produce json
// @Description 异步部署SSL证书到目标主机(备份、拷贝、重载Nginx)，立即返回部署日志ID，前端可通过查询部署日志接口获取实时进度
// @Param data body model.DeploySslCertRequest true "部署请求"
// @Success 200 {object} result.Result{data=model.DeploySslCertResponse}
// @Router /api/v1/monitor/ssl-cert/deploy [post]
// @Security ApiKeyAuth
func (ctrl *SslCertController) DeploySslCert(c *gin.Context) {
	var req model.DeploySslCertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "参数错误: "+err.Error())
		return
	}

	ctrl.sslCertService.DeploySslCert(c, &req)
}

// @Tags SSL证书管理
// @Summary 获取SSL证书部署日志
// @Description 获取SSL证书部署日志
// @Produce json
// @Description 获取SSL证书部署日志列表，支持多条件筛选和分页
// @Param certId query int false "证书ID"
// @Param hostId query int false "主机ID"
// @Param domain query string false "域名(模糊查询)"
// @Param status query int false "部署状态(1=部署中 2=成功 3=失败)"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} result.Result{data=model.GetDeployLogsResponse}
// @Router /api/v1/monitor/ssl-cert/deploy-logs [get]
// @Security ApiKeyAuth
func (ctrl *SslCertController) GetDeployLogs(c *gin.Context) {
	var req model.GetDeployLogsRequest

	// 解析查询参数
	if certID := c.Query("certId"); certID != "" {
		id, _ := strconv.ParseUint(certID, 10, 64)
		req.CertID = uint(id)
	}

	if hostID := c.Query("hostId"); hostID != "" {
		id, _ := strconv.ParseUint(hostID, 10, 64)
		req.HostID = uint(id)
	}

	req.Domain = c.Query("domain")

	if status := c.Query("status"); status != "" {
		s, _ := strconv.Atoi(status)
		req.Status = &s
	}

	req.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	req.PageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	ctrl.sslCertService.GetDeployLogs(c, &req)
}

// @Tags SSL证书管理
// @Summary 上传SSL证书
// @Description 上传SSL证书
// @Accept multipart/form-data
// @Produce json
// @Description 上传SSL证书压缩包(仅支持zip格式),自动解析证书和私钥并录入数据库。压缩包中应包含xxx.key和xxx.pem/crt/cer文件,域名从文件名自动提取
// @Param file formData file true "证书压缩包(zip格式)"
// @Success 200 {object} result.Result{data=model.UploadSslCertResponse}
// @Router /api/v1/monitor/ssl-cert/upload [post]
// @Security ApiKeyAuth
func (ctrl *SslCertController) UploadSslCert(c *gin.Context) {
	ctrl.sslCertService.UploadSslCert(c)
}

