// 运维知识库模型
package model

import (
	"time"
)

// Knowledge 运维知识库模型
type Knowledge struct {
	ID         uint      `gorm:"column:id;comment:'主键';primaryKey;autoIncrement;NOT NULL" json:"id"`
	Title      string    `gorm:"column:title;varchar(200);comment:'标题';NOT NULL" json:"title"`
	Category   string    `gorm:"column:category;varchar(50);default:'其他';comment:'分类'" json:"category"`
	Content    string    `gorm:"column:content;type:longtext;comment:'Markdown内容'" json:"content"`
	Tags       string    `gorm:"column:tags;varchar(500);default:'';comment:'标签(JSON数组)'" json:"tags"`
	Status     int       `gorm:"column:status;default:1;comment:'状态:1->已发布,2->草稿'" json:"status"`
	Author     string    `gorm:"column:author;varchar(50);default:'';comment:'作者'" json:"author"`
	CreateTime time.Time `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
}

func (Knowledge) TableName() string {
	return "ops_knowledge"
}

// AddKnowledgeDto 新增知识DTO
type AddKnowledgeDto struct {
	Title    string `json:"title" validate:"required,min=1,max=200"` // 标题
	Category string `json:"category"`                                // 分类
	Content  string `json:"content"`                                 // Markdown内容
	Tags     string `json:"tags"`                                    // 标签(JSON数组)
	Status   int    `json:"status"`                                  // 状态:1->已发布,2->草稿
}

// UpdateKnowledgeDto 更新知识DTO
type UpdateKnowledgeDto struct {
	ID       uint   `json:"id" validate:"required"`                  // ID
	Title    string `json:"title" validate:"required,min=1,max=200"` // 标题
	Category string `json:"category"`                                // 分类
	Content  string `json:"content"`                                 // Markdown内容
	Tags     string `json:"tags"`                                    // 标签(JSON数组)
	Status   int    `json:"status"`                                  // 状态:1->已发布,2->草稿
}

// KnowledgeQueryDto 查询DTO
type KnowledgeQueryDto struct {
	Title    string `form:"title"`    // 标题(模糊查询)
	Category string `form:"category"` // 分类
	Status   *int   `form:"status"`   // 状态
	PageNum  int    `form:"pageNum"`  // 页码
	PageSize int    `form:"pageSize"` // 每页数量
}

// KnowledgeVo 知识VO
type KnowledgeVo struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Category   string    `json:"category"`
	Content    string    `json:"content"`
	Tags       string    `json:"tags"`
	Status     int       `json:"status"`
	StatusText string    `json:"statusText"`
	Author     string    `json:"author"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

// KnowledgeCategory 知识分类
type KnowledgeCategory struct {
	ID          uint      `gorm:"column:id;comment:'主键';primaryKey;autoIncrement;NOT NULL" json:"id"`
	Name        string    `gorm:"column:name;varchar(50);comment:'分类名称';NOT NULL;unique" json:"name"`
	Sort        int       `gorm:"column:sort;default:0;comment:'排序'" json:"sort"`
	CreateTime  time.Time `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	Description string    `gorm:"column:description;varchar(200);comment:'分类描述'" json:"description"`
}

func (KnowledgeCategory) TableName() string {
	return "ops_knowledge_category"
}

// AddCategoryDto 新增分类DTO
type AddCategoryDto struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Sort        int    `json:"sort"`
	Description string `json:"description"`
}

// UpdateCategoryDto 更新分类DTO
type UpdateCategoryDto struct {
	ID          uint   `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Sort        int    `json:"sort"`
	Description string `json:"description"`
}

// CategoryVo 分类VO
type CategoryVo struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Sort        int       `json:"sort"`
	Description string    `json:"description"`
	CreateTime  time.Time `json:"createTime"`
	DocCount    int       `json:"docCount"`
}
