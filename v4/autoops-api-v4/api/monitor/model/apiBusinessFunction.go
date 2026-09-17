package model

import "dodevops-api/common/util"

// APIBusinessFunction API业务功能分类
type APIBusinessFunction struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"type:varchar(200);not null;comment:业务功能名称"`
	Description string     `json:"description" gorm:"type:varchar(500);comment:描述"`

	// 关联业务组
	GroupID     uint       `json:"groupId" gorm:"type:int;not null;index;comment:业务组ID"`
	GroupName   string     `json:"groupName" gorm:"type:varchar(100);comment:业务组名称"`

	// 排序和状态
	Sort        int        `json:"sort" gorm:"type:int;default:0;comment:排序"`
	Status      int        `json:"status" gorm:"type:tinyint;default:1;index;comment:启用状态(1:启用,0:禁用)"`

	// 统计字段（不存数据库，查询时计算）
	APICount    int        `json:"apiCount" gorm:"-"`
	HealthyAPIs int        `json:"healthyApis" gorm:"-"`
	UnhealthyAPIs int      `json:"unhealthyApis" gorm:"-"`

	CreatedBy   string     `json:"createdBy" gorm:"type:varchar(100);comment:创建人"`
	UpdatedBy   string     `json:"updatedBy" gorm:"type:varchar(100);comment:更新人"`
	CreatedAt   util.HTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   util.HTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (APIBusinessFunction) TableName() string {
	return "monitor_api_business_function"
}

// DTO - 创建业务功能
type CreateAPIBusinessFunctionDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	GroupID     uint   `json:"groupId" binding:"required"`
	GroupName   string `json:"groupName"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
}

// DTO - 更新业务功能
type UpdateAPIBusinessFunctionDto struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupID     uint   `json:"groupId"`
	GroupName   string `json:"groupName"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
}

// DTO - 业务功能统计（含API列表）
type APIBusinessFunctionWithAPIs struct {
	APIBusinessFunction
	APIs []APIMonitor `json:"apis"` // API列表
}
