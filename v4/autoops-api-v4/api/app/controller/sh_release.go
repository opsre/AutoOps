package controller

import (
	"strconv"

	"dodevops-api/api/app/model"
	"dodevops-api/api/app/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ShReleaseController 脚本上线控制器
type ShReleaseController struct {
	service service.IShReleaseService
}

// NewShReleaseController 创建脚本上线控制器
func NewShReleaseController(db *gorm.DB) *ShReleaseController {
	return &ShReleaseController{
		service: service.NewShReleaseService(db),
	}
}

// CreateShRelease 创建脚本上线工单
// @Summary 创建脚本上线工单
// @Description 创建脚本上线工单，用于执行生产环境脚本
// @Tags ShRelease
// @Accept json
// @Produce json
// @Param request body model.CreateShReleaseRequest true "创建脚本上线请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.ShRelease}
// @Router /apps/sh-release [post]
func (ctrl *ShReleaseController) CreateShRelease(c *gin.Context) {
	var req model.CreateShReleaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.CreateShRelease(c, &req)
}

// GetShReleaseList 获取脚本上线工单列表
// @Summary 获取脚本上线工单列表
// @Description 获取脚本上线工单列表，支持分页和筛选
// @Tags ShRelease
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param business_group_id query int false "业务线ID"
// @Param status query int false "流程状态"
// @Param applicant_id query int false "申请人ID"
// @Param keyword query string false "搜索关键词（标题）"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.ShReleaseListResponse}
// @Router /apps/sh-release/list [get]
func (ctrl *ShReleaseController) GetShReleaseList(c *gin.Context) {
	var req model.ShReleaseListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.GetShReleaseList(c, &req)
}

// GetShReleaseDetail 获取脚本上线工单详情
// @Summary 获取脚本上线工单详情
// @Description 获取脚本上线工单详情
// @Tags ShRelease
// @Accept json
// @Produce json
// @Param id path int true "工单ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=model.ShRelease}
// @Router /apps/sh-release/{id} [get]
func (ctrl *ShReleaseController) GetShReleaseDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的工单ID")
		return
	}

	ctrl.service.GetShReleaseDetail(c, uint(id))
}

// DeleteShRelease 删除脚本上线工单
// @Summary 删除脚本上线工单
// @Description 删除脚本上线工单
// @Tags ShRelease
// @Accept json
// @Produce json
// @Param id path int true "工单ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=string}
// @Router /apps/sh-release/{id} [delete]
func (ctrl *ShReleaseController) DeleteShRelease(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的工单ID")
		return
	}

	ctrl.service.DeleteShRelease(c, uint(id))
}

// ApproveShRelease 审核脚本上线工单
// @Summary 审核脚本上线工单
// @Description 审核人审核脚本上线工单
// @Tags ShRelease
// @Accept json
// @Produce json
// @Param request body model.ShApprovalRequest true "审核请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=string}
// @Router /apps/sh-release/approve [post]
func (ctrl *ShReleaseController) ApproveShRelease(c *gin.Context) {
	var req model.ShApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.ApproveShRelease(c, &req)
}

// PullCode 拉取代码
// @Summary 拉取代码
// @Description 审核通过后，通过Jenkins拉取代码
// @Tags ShRelease
// @Accept json
// @Produce json
// @Param request body model.ShExecuteRequest true "拉取代码请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=string}
// @Router /apps/sh-release/pull-code [post]
func (ctrl *ShReleaseController) PullCode(c *gin.Context) {
	var req model.ShExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.PullCode(c, &req)
}

// ExecuteScript 执行脚本
// @Summary 执行脚本
// @Description 代码拉取成功后，通过SSH在服务器上执行脚本
// @Tags ShRelease
// @Accept json
// @Produce json
// @Param request body model.ShExecuteScriptRequest true "执行脚本请求"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=string}
// @Router /apps/sh-release/execute-script [post]
func (ctrl *ShReleaseController) ExecuteScript(c *gin.Context) {
	var req model.ShExecuteScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.ExecuteScript(c, &req)
}

// GetShReleaseLog 获取Jenkins日志
// @Summary 获取Jenkins日志
// @Description 获取脚本上线工单的Jenkins代码拉取日志
// @Tags ShRelease
// @Accept json
// @Produce json
// @Param id path int true "工单ID"
// @Param start query int false "日志起始行数" default(0)
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=object}
// @Router /apps/sh-release/{id}/jenkins-log [get]
func (ctrl *ShReleaseController) GetShReleaseLog(c *gin.Context) {
	releaseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的工单ID")
		return
	}

	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))

	ctrl.service.GetShReleaseLog(c, uint(releaseID), start)
}

// GetScriptOutput 获取脚本执行输出
// @Summary 获取脚本执行输出
// @Description 获取脚本通过SSH执行后的输出结果
// @Tags ShRelease
// @Accept json
// @Produce json
// @Param id path int true "工单ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=object}
// @Router /apps/sh-release/{id}/script-output [get]
func (ctrl *ShReleaseController) GetScriptOutput(c *gin.Context) {
	releaseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的工单ID")
		return
	}

	ctrl.service.GetScriptOutput(c, uint(releaseID))
}

// CheckUserPermission 检查用户操作权限
// @Summary 检查用户操作权限
// @Description 检查当前用户是否有权限执行指定的操作
// @Tags ShRelease
// @Accept json
// @Produce json
// @Param id path int true "工单ID"
// @Param operation_type query string true "操作类型" Enums(approval, execute)
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=object}
// @Router /apps/sh-release/{id}/permission [get]
func (ctrl *ShReleaseController) CheckUserPermission(c *gin.Context) {
	releaseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的工单ID")
		return
	}

	operationType := c.Query("operation_type")
	if operationType == "" {
		result.Failed(c, 400, "请指定操作类型")
		return
	}

	ctrl.service.CheckUserPermission(c, uint(releaseID), operationType)
}

// ListReleaseFiles 列出脚本上线服务器上的文件
// @Summary 列出脚本上线服务器上的文件
// @Description 列出脚本上线工单关联的服务器上指定目录的文件列表，默认目录为/home/dingding/saas3-data-export/xlsx
// @Tags ShRelease
// @Accept json
// @Produce json
// @Param id path int true "工单ID"
// @Param remotePath query string false "远程目录路径(默认为/home/dingding/saas3-data-export/xlsx)"
// @Security ApiKeyAuth
// @Success 200 {object} result.Result{data=object}
// @Router /apps/sh-release/{id}/files [get]
func (ctrl *ShReleaseController) ListReleaseFiles(c *gin.Context) {
	releaseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的工单ID")
		return
	}

	remotePath := c.Query("remotePath")
	ctrl.service.ListReleaseFiles(c, uint(releaseID), remotePath)
}

// DownloadReleaseFile 下载脚本上线服务器上的文件
// @Summary 下载脚本上线服务器上的文件
// @Description 从脚本上线工单关联的服务器上下载指定文件
// @Tags ShRelease
// @Accept json
// @Produce octet-stream
// @Param id path int true "工单ID"
// @Param remotePath query string true "远程文件路径"
// @Security ApiKeyAuth
// @Success 200 {file} file "下载的文件"
// @Router /apps/sh-release/{id}/download [get]
func (ctrl *ShReleaseController) DownloadReleaseFile(c *gin.Context) {
	releaseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		result.Failed(c, 400, "无效的工单ID")
		return
	}

	remotePath := c.Query("remotePath")
	if remotePath == "" {
		result.Failed(c, 400, "远程文件路径不能为空")
		return
	}

	ctrl.service.DownloadReleaseFile(c, uint(releaseID), remotePath)
}
