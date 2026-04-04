// 运维知识库 DAO层
package dao

import (
	"dodevops-api/api/tool/model"
	"dodevops-api/common"
	"time"
)

// CreateKnowledge 创建知识
func CreateKnowledge(dto model.AddKnowledgeDto, author string) error {
	category := dto.Category
	if category == "" {
		category = "其他"
	}
	status := dto.Status
	if status == 0 {
		status = 1
	}

	data := map[string]interface{}{
		"title":    dto.Title,
		"category": category,
		"content":  dto.Content,
		"tags":     dto.Tags,
		"status":   status,
		"author":   author,
	}

	return common.GetDB().Table("ops_knowledge").Create(&data).Error
}

// GetKnowledgeByID 根据ID获取知识
func GetKnowledgeByID(id uint) (*model.Knowledge, error) {
	var knowledge model.Knowledge
	err := common.GetDB().Where("id = ?", id).First(&knowledge).Error
	if err != nil {
		return nil, err
	}
	return &knowledge, nil
}

// UpdateKnowledge 更新知识
func UpdateKnowledge(dto model.UpdateKnowledgeDto) error {
	updates := map[string]interface{}{
		"title":       dto.Title,
		"category":    dto.Category,
		"content":     dto.Content,
		"tags":        dto.Tags,
		"status":      dto.Status,
		"update_time": time.Now(),
	}
	if dto.Category == "" {
		updates["category"] = "其他"
	}
	return common.GetDB().Model(&model.Knowledge{}).Where("id = ?", dto.ID).Updates(updates).Error
}

// DeleteKnowledge 删除知识
func DeleteKnowledge(id uint) error {
	return common.GetDB().Where("id = ?", id).Delete(&model.Knowledge{}).Error
}

// GetKnowledgeList 获取知识列表(分页)
func GetKnowledgeList(dto model.KnowledgeQueryDto) ([]model.KnowledgeVo, int64, error) {
	var list []model.Knowledge
	var total int64

	db := common.GetDB().Model(&model.Knowledge{})

	if dto.Title != "" {
		db = db.Where("title LIKE ?", "%"+dto.Title+"%")
	}
	if dto.Category != "" {
		db = db.Where("category = ?", dto.Category)
	}
	if dto.Status != nil {
		db = db.Where("status = ?", *dto.Status)
	}

	db.Count(&total)

	offset := (dto.PageNum - 1) * dto.PageSize
	if offset < 0 {
		offset = 0
	}
	if dto.PageSize <= 0 {
		dto.PageSize = 10
	}

	err := db.Order("create_time DESC").Offset(offset).Limit(dto.PageSize).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	var voList []model.KnowledgeVo
	for _, item := range list {
		vo := model.KnowledgeVo{
			ID:         item.ID,
			Title:      item.Title,
			Category:   item.Category,
			Content:    item.Content,
			Tags:       item.Tags,
			Status:     item.Status,
			StatusText: getStatusText(item.Status),
			Author:     item.Author,
			CreateTime: item.CreateTime,
			UpdateTime: item.UpdateTime,
		}
		voList = append(voList, vo)
	}

	return voList, total, nil
}

// GetAllCategories 获取所有分类
func GetAllCategories() ([]model.KnowledgeCategory, error) {
	var categories []model.KnowledgeCategory
	err := common.GetDB().Table("ops_knowledge_category").
		Order("sort ASC, id ASC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// ==================== 知识分类管理 ====================

// CreateCategory 创建分类
func CreateCategory(dto model.AddCategoryDto) error {
	category := model.KnowledgeCategory{
		Name:        dto.Name,
		Sort:        dto.Sort,
		Description: dto.Description,
		CreateTime:  time.Now(),
	}
	return common.GetDB().Create(&category).Error
}

// GetCategoryByID 根据ID获取分类
func GetCategoryByID(id uint) (*model.KnowledgeCategory, error) {
	var category model.KnowledgeCategory
	err := common.GetDB().Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory 更新分类
func UpdateCategory(dto model.UpdateCategoryDto) error {
	updates := map[string]interface{}{
		"name":        dto.Name,
		"sort":        dto.Sort,
		"description": dto.Description,
	}
	return common.GetDB().Model(&model.KnowledgeCategory{}).Where("id = ?", dto.ID).Updates(updates).Error
}

// DeleteCategory 删除分类
func DeleteCategory(id uint) error {
	return common.GetDB().Where("id = ?", id).Delete(&model.KnowledgeCategory{}).Error
}

// GetCategoryList 获取分类列表(含文档数量)
func GetCategoryList() ([]model.CategoryVo, error) {
	var categories []model.KnowledgeCategory
	err := common.GetDB().Order("sort ASC, id ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	var voList []model.CategoryVo
	for _, cat := range categories {
		var docCount int64
		common.GetDB().Model(&model.Knowledge{}).Where("category = ?", cat.Name).Count(&docCount)

		voList = append(voList, model.CategoryVo{
			ID:          cat.ID,
			Name:        cat.Name,
			Sort:        cat.Sort,
			Description: cat.Description,
			CreateTime:  cat.CreateTime,
			DocCount:    int(docCount),
		})
	}
	return voList, nil
}

// GetAllCategoryNames 获取所有分类名称
func GetAllCategoryNames() ([]string, error) {
	var names []string
	err := common.GetDB().Model(&model.KnowledgeCategory{}).
		Order("sort ASC, id ASC").
		Pluck("name", &names).Error
	return names, err
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
