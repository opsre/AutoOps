// 运维知识库 Controller层
package controller

import (
	"dodevops-api/api/tool/model"
	"dodevops-api/api/tool/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateKnowledge 创建知识
// @Tags Knowledge运维知识库
// @Summary 创建知识
// @Description 创建知识接口
// @Accept json
// @Produce json
// @Param data body model.AddKnowledgeDto true "知识信息"
// @Success 200 {object} result.Result
// @Router /api/v1/knowledge [post]
// @Security ApiKeyAuth
func CreateKnowledge(c *gin.Context) {
	var dto model.AddKnowledgeDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	service.KnowledgeService().CreateKnowledge(c, dto)
}

// GetKnowledgeByID 根据ID获取知识
// @Tags Knowledge运维知识库
// @Summary 根据ID获取知识
// @Description 根据ID获取知识详情
// @Produce json
// @Param id path int true "知识ID"
// @Success 200 {object} result.Result
// @Router /api/v1/knowledge/{id} [get]
// @Security ApiKeyAuth
func GetKnowledgeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "ID格式错误"})
		return
	}
	service.KnowledgeService().GetKnowledgeByID(c, uint(id))
}

// UpdateKnowledge 更新知识
// @Tags Knowledge运维知识库
// @Summary 更新知识
// @Description 更新知识接口
// @Accept json
// @Produce json
// @Param data body model.UpdateKnowledgeDto true "知识信息"
// @Success 200 {object} result.Result
// @Router /api/v1/knowledge [put]
// @Security ApiKeyAuth
func UpdateKnowledge(c *gin.Context) {
	var dto model.UpdateKnowledgeDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	service.KnowledgeService().UpdateKnowledge(c, dto)
}

// DeleteKnowledge 删除知识
// @Tags Knowledge运维知识库
// @Summary 删除知识
// @Description 删除知识接口
// @Produce json
// @Param id path int true "知识ID"
// @Success 200 {object} result.Result
// @Router /api/v1/knowledge/{id} [delete]
// @Security ApiKeyAuth
func DeleteKnowledge(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "ID格式错误"})
		return
	}
	service.KnowledgeService().DeleteKnowledge(c, uint(id))
}

// GetKnowledgeList 获取知识列表(分页)
// @Tags Knowledge运维知识库
// @Summary 获取知识列表
// @Description 获取知识列表(分页)
// @Produce json
// @Param title query string false "标题(模糊查询)"
// @Param category query string false "分类"
// @Param status query int false "状态"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/knowledge/list [get]
// @Security ApiKeyAuth
func GetKnowledgeList(c *gin.Context) {
	var dto model.KnowledgeQueryDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	service.KnowledgeService().GetKnowledgeList(c, dto)
}

// GetAllKnowledgeCategories 获取所有分类
// @Tags Knowledge运维知识库
// @Summary 获取所有分类
// @Description 获取所有知识分类
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/knowledge/categories [get]
// @Security ApiKeyAuth
func GetAllKnowledgeCategories(c *gin.Context) {
	service.KnowledgeService().GetAllCategories(c)
}

// ==================== 知识分类管理 ====================

// CreateCategory 创建分类
// @Tags Knowledge运维知识库
// @Summary 创建分类
// @Description 创建知识分类
// @Accept json
// @Produce json
// @Param data body model.AddCategoryDto true "分类信息"
// @Success 200 {object} result.Result
// @Router /api/v1/knowledge/category [post]
// @Security ApiKeyAuth
func CreateCategory(c *gin.Context) {
	var dto model.AddCategoryDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	service.CategoryService().CreateCategory(c, dto)
}

// GetCategoryByID 根据ID获取分类
// @Tags Knowledge运维知识库
// @Summary 根据ID获取分类
// @Description 根据ID获取分类详情
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} result.Result
// @Router /api/v1/knowledge/category/{id} [get]
// @Security ApiKeyAuth
func GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "ID格式错误"})
		return
	}
	service.CategoryService().GetCategoryByID(c, uint(id))
}

// UpdateCategory 更新分类
// @Tags Knowledge运维知识库
// @Summary 更新分类
// @Description 更新知识分类
// @Accept json
// @Produce json
// @Param data body model.UpdateCategoryDto true "分类信息"
// @Success 200 {object} result.Result
// @Router /api/v1/knowledge/category [put]
// @Security ApiKeyAuth
func UpdateCategory(c *gin.Context) {
	var dto model.UpdateCategoryDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	service.CategoryService().UpdateCategory(c, dto)
}

// DeleteCategory 删除分类
// @Tags Knowledge运维知识库
// @Summary 删除分类
// @Description 删除知识分类
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} result.Result
// @Router /api/v1/knowledge/category/{id} [delete]
// @Security ApiKeyAuth
func DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "ID格式错误"})
		return
	}
	service.CategoryService().DeleteCategory(c, uint(id))
}

// GetCategoryList 获取分类列表
// @Tags Knowledge运维知识库
// @Summary 获取分类列表
// @Description 获取知识分类列表(含文档数量)
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/knowledge/category/list [get]
// @Security ApiKeyAuth
func GetCategoryList(c *gin.Context) {
	service.CategoryService().GetCategoryList(c)
}
