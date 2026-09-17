package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/api/monitor/model"
	"dodevops-api/api/monitor/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type APIBusinessFunctionController struct {
	service service.APIBusinessFunctionService
}

func NewAPIBusinessFunctionController() *APIBusinessFunctionController {
	return &APIBusinessFunctionController{
		service: service.NewAPIBusinessFunctionService(),
	}
}

// Create 创建业务功能
// @Summary 创建业务功能
// @Description 创建新的业务功能分类
// @Tags API业务功能
// @Accept json
// @Produce json
// @Param function body model.CreateAPIBusinessFunctionDto true "业务功能信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/business-functions [post]
// @Security ApiKeyAuth
func (ctrl *APIBusinessFunctionController) Create(c *gin.Context) {
	var dto model.CreateAPIBusinessFunctionDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	function, err := ctrl.service.Create(&dto)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	result.Success(c, function)
}

// Update 更新业务功能
// @Summary 更新业务功能
// @Description 更新已有的业务功能信息
// @Tags API业务功能
// @Accept json
// @Produce json
// @Param id path uint true "业务功能ID"
// @Param function body model.UpdateAPIBusinessFunctionDto true "业务功能信息"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/business-functions/{id} [put]
// @Security ApiKeyAuth
func (ctrl *APIBusinessFunctionController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的业务功能ID")
		return
	}

	var dto model.UpdateAPIBusinessFunctionDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	dto.ID = uint(id)
	function, err := ctrl.service.Update(&dto)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	result.Success(c, function)
}

// Delete 删除业务功能
// @Summary 删除业务功能
// @Description 删除指定的业务功能
// @Tags API业务功能
// @Accept json
// @Produce json
// @Param id path uint true "业务功能ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/business-functions/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *APIBusinessFunctionController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的业务功能ID")
		return
	}

	if err := ctrl.service.Delete(uint(id)); err != nil {
		result.Failed(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	result.Success(c, nil)
}

// GetByID 获取业务功能详情
// @Summary 获取业务功能详情
// @Description 获取指定ID的业务功能详细信息
// @Tags API业务功能
// @Accept json
// @Produce json
// @Param id path uint true "业务功能ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/business-functions/{id} [get]
// @Security ApiKeyAuth
func (ctrl *APIBusinessFunctionController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的业务功能ID")
		return
	}

	function, err := ctrl.service.GetByID(uint(id))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	if function == nil {
		result.Failed(c, http.StatusNotFound, "业务功能不存在")
		return
	}

	result.Success(c, function)
}

// GetList 获取业务功能列表
// @Summary 获取业务功能列表
// @Description 获取业务功能列表，支持按业务组筛选和分页
// @Tags API业务功能
// @Accept json
// @Produce json
// @Param groupId query uint false "业务组ID"
// @Param keyword query string false "关键字（功能名称）"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/business-functions [get]
// @Security ApiKeyAuth
func (ctrl *APIBusinessFunctionController) GetList(c *gin.Context) {
	// 解析分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 解析筛选参数
	groupIDStr := c.Query("groupId")
	var groupID *uint
	if groupIDStr != "" {
		id, err := strconv.ParseUint(groupIDStr, 10, 32)
		if err == nil {
			u := uint(id)
			groupID = &u
		}
	}

	keyword := c.Query("keyword")

	// 查询
	functions, total, err := ctrl.service.GetList(groupID, keyword, page, pageSize)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	result.SuccessWithPage(c, functions, total, page, pageSize)
}

// GetFunctionsByGroupID 获取指定业务组的所有业务功能（用于级联选择）
// @Summary 获取业务组的所有功能
// @Description 获取指定业务组下的所有业务功能，用于前端级联选择
// @Tags API业务功能
// @Accept json
// @Produce json
// @Param groupId query uint true "业务组ID"
// @Success 200 {object} result.Result
// @Router /api/v1/monitor/business-functions-by-group [get]
// @Security ApiKeyAuth
func (ctrl *APIBusinessFunctionController) GetFunctionsByGroupID(c *gin.Context) {
	groupIDStr := c.Query("groupId")
	if groupIDStr == "" {
		result.Failed(c, http.StatusBadRequest, "业务组ID不能为空")
		return
	}

	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的业务组ID")
		return
	}

	functions, err := ctrl.service.GetByGroupID(uint(groupID))
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	// 返回包含分页信息的结构化数据（兼容前端期望的格式）
	result.Success(c, gin.H{
		"list":     functions,
		"total":    len(functions),
		"page":     1,
		"pageSize": len(functions),
	})
}
