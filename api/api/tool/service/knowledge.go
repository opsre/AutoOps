// 运维知识库 Service层
package service

import (
	sysModel "dodevops-api/api/system/model"
	dao "dodevops-api/api/tool/dao1"
	"dodevops-api/api/tool/model"
	"dodevops-api/common"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)



// IKnowledgeService 接口定义
type IKnowledgeService interface {
	CreateKnowledge(c *gin.Context, dto model.AddKnowledgeDto)
	GetKnowledgeByID(c *gin.Context, id uint)
	UpdateKnowledge(c *gin.Context, dto model.UpdateKnowledgeDto)
	DeleteKnowledge(c *gin.Context, id uint)
	GetKnowledgeList(c *gin.Context, dto model.KnowledgeQueryDto)
	GetAllCategories(c *gin.Context)
}

type KnowledgeServiceImpl struct{}

// CreateKnowledge 创建知识
func (s KnowledgeServiceImpl) CreateKnowledge(c *gin.Context, dto model.AddKnowledgeDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), result.ApiCode.GetMessage(result.ApiCode.ValidationParameterError))
		return
	}

	author := getCurrentUsername(c)
	err = dao.CreateKnowledge(dto, author)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "创建失败: "+err.Error())
		return
	}

	result.Success(c, nil)
}

// GetKnowledgeByID 根据ID获取知识
func (s KnowledgeServiceImpl) GetKnowledgeByID(c *gin.Context, id uint) {
	if id == 0 {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), "ID不能为空")
		return
	}

	knowledge, err := dao.GetKnowledgeByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询失败: "+err.Error())
		return
	}

	vo := model.KnowledgeVo{
		ID:         knowledge.ID,
		Title:      knowledge.Title,
		Category:   knowledge.Category,
		Content:    knowledge.Content,
		Tags:       knowledge.Tags,
		Status:     knowledge.Status,
		StatusText: getStatusText(knowledge.Status),
		Author:     knowledge.Author,
		CreateTime: knowledge.CreateTime,
		UpdateTime: knowledge.UpdateTime,
	}

	result.Success(c, vo)
}

// UpdateKnowledge 更新知识
func (s KnowledgeServiceImpl) UpdateKnowledge(c *gin.Context, dto model.UpdateKnowledgeDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), result.ApiCode.GetMessage(result.ApiCode.ValidationParameterError))
		return
	}

	_, err = dao.GetKnowledgeByID(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "知识不存在")
		return
	}

	err = dao.UpdateKnowledge(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "更新失败: "+err.Error())
		return
	}

	result.Success(c, nil)
}

// DeleteKnowledge 删除知识
func (s KnowledgeServiceImpl) DeleteKnowledge(c *gin.Context, id uint) {
	if id == 0 {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), "ID不能为空")
		return
	}

	_, err := dao.GetKnowledgeByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "知识不存在")
		return
	}

	err = dao.DeleteKnowledge(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "删除失败: "+err.Error())
		return
	}

	result.Success(c, nil)
}

// GetKnowledgeList 获取知识列表
func (s KnowledgeServiceImpl) GetKnowledgeList(c *gin.Context, dto model.KnowledgeQueryDto) {
	if dto.PageNum <= 0 {
		dto.PageNum = 1
	}
	if dto.PageSize <= 0 {
		dto.PageSize = 10
	}

	list, total, err := dao.GetKnowledgeList(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询失败: "+err.Error())
		return
	}

	data := map[string]interface{}{
		"list":     list,
		"total":    total,
		"pageNum":  dto.PageNum,
		"pageSize": dto.PageSize,
	}

	result.Success(c, data)
}

// GetAllCategories 获取所有分类
func (s KnowledgeServiceImpl) GetAllCategories(c *gin.Context) {
	categories, err := dao.GetAllCategories()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询失败: "+err.Error())
		return
	}
	result.Success(c, categories)
}

func getCurrentUsername(c *gin.Context) string {
	sysAdmin, exists := c.Get("sysAdmin")
	if !exists {
		return "admin"
	}
	if admin, ok := sysAdmin.(*sysModel.JwtAdmin); ok {
		return admin.Username
	}
	return "admin"
}

func getStatusText(status int) string {
	switch status {
	case 1:
		return "已发布"
	case 2:
		return "草稿"
	default:
		return "未知"
	}
}

var knowledgeService = KnowledgeServiceImpl{}

func KnowledgeService() IKnowledgeService {
	return &knowledgeService
}

// ==================== 知识分类管理 ====================

// ICategoryService 分类接口定义
type ICategoryService interface {
	CreateCategory(c *gin.Context, dto model.AddCategoryDto)
	GetCategoryByID(c *gin.Context, id uint)
	UpdateCategory(c *gin.Context, dto model.UpdateCategoryDto)
	DeleteCategory(c *gin.Context, id uint)
	GetCategoryList(c *gin.Context)
}

type CategoryServiceImpl struct{}

// CreateCategory 创建分类
func (s CategoryServiceImpl) CreateCategory(c *gin.Context, dto model.AddCategoryDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), result.ApiCode.GetMessage(result.ApiCode.ValidationParameterError))
		return
	}

	err = dao.CreateCategory(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "创建失败: "+err.Error())
		return
	}

	result.Success(c, nil)
}

// GetCategoryByID 根据ID获取分类
func (s CategoryServiceImpl) GetCategoryByID(c *gin.Context, id uint) {
	if id == 0 {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), "ID不能为空")
		return
	}

	category, err := dao.GetCategoryByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询失败: "+err.Error())
		return
	}

	result.Success(c, category)
}

// UpdateCategory 更新分类
func (s CategoryServiceImpl) UpdateCategory(c *gin.Context, dto model.UpdateCategoryDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), result.ApiCode.GetMessage(result.ApiCode.ValidationParameterError))
		return
	}

	_, err = dao.GetCategoryByID(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "分类不存在")
		return
	}

	oldCategory, _ := dao.GetCategoryByID(dto.ID)
	err = dao.UpdateCategory(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "更新失败: "+err.Error())
		return
	}

	// 如果分类名称变更，同步更新知识库中的分类名称
	if oldCategory.Name != dto.Name {
		common.GetDB().Model(&model.Knowledge{}).Where("category = ?", oldCategory.Name).Update("category", dto.Name)
	}

	result.Success(c, nil)
}

// DeleteCategory 删除分类
func (s CategoryServiceImpl) DeleteCategory(c *gin.Context, id uint) {
	if id == 0 {
		result.Failed(c, int(result.ApiCode.ValidationParameterError), "ID不能为空")
		return
	}

	category, err := dao.GetCategoryByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "分类不存在")
		return
	}

	err = dao.DeleteCategory(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "删除失败: "+err.Error())
		return
	}

	// 删除分类后，将该分类下的知识库文档归类为"未分类"
	common.GetDB().Model(&model.Knowledge{}).Where("category = ?", category.Name).Update("category", "未分类")

	result.Success(c, nil)
}

// GetCategoryList 获取分类列表
func (s CategoryServiceImpl) GetCategoryList(c *gin.Context) {
	list, err := dao.GetCategoryList()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询失败: "+err.Error())
		return
	}

	result.Success(c, list)
}

var categoryService = CategoryServiceImpl{}

func CategoryService() ICategoryService {
	return &categoryService
}
