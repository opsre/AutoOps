package model

import (
	"database/sql/driver"
	"encoding/json"
	"dodevops-api/common/util"
)

// UintArray 自定义类型用于存储uint数组到JSON
type UintArray []uint

// Scan 实现sql.Scanner接口
func (u *UintArray) Scan(value interface{}) error {
	if value == nil {
		*u = []uint{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, u)
}

// Value 实现driver.Valuer接口
func (u UintArray) Value() (driver.Value, error) {
	if u == nil {
		return "[]", nil
	}
	return json.Marshal(u)
}

// 堡垒机用户组模型
type BastionUserGroup struct {
	ID          uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Name        string     `gorm:"column:name;type:varchar(64);comment:'用户组名称';NOT NULL;uniqueIndex" json:"name"`
	Description string     `gorm:"column:description;type:varchar(500);comment:'用户组描述'" json:"description"`
	MemberIds   UintArray  `gorm:"column:member_ids;type:json;comment:'成员ID列表(JSON数组)'" json:"memberIds"` // 成员ID列表
	Status      int        `gorm:"column:status;default:1;comment:'状态:1->启用,2->禁用';NOT NULL" json:"status"`
	CreateTime  util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime  util.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
}

func (BastionUserGroup) TableName() string {
	return "bastion_user_group"
}

// 创建用户组DTO
type CreateBastionUserGroupDto struct {
	Name        string `validate:"required" json:"name"`   // 用户组名称
	Description string `json:"description"`                // 用户组描述
	MemberIds   []uint `json:"memberIds"`                  // 成员ID列表
	Status      int    `validate:"required" json:"status"` // 状态:1->启用,2->禁用
}

// 更新用户组DTO
type UpdateBastionUserGroupDto struct {
	ID          uint   `json:"id"`                              // 用户组ID
	Name        string `validate:"required" json:"name"`        // 用户组名称
	Description string `json:"description"`                     // 用户组描述
	MemberIds   []uint `json:"memberIds"`                       // 成员ID列表
	Status      int    `validate:"required" json:"status"`      // 状态:1->启用,2->禁用
}

// 用户组ID DTO
type BastionUserGroupIdDto struct {
	ID uint `json:"id"` // 用户组ID
}

// 添加用户组成员DTO
type AddGroupMembersDto struct {
	GroupID uint   `validate:"required" json:"groupId"` // 用户组ID
	UserIDs []uint `validate:"required" json:"userIds"` // 要添加的用户ID列表
}

// 删除用户组成员DTO
type RemoveGroupMembersDto struct {
	GroupID uint   `validate:"required" json:"groupId"` // 用户组ID
	UserIDs []uint `validate:"required" json:"userIds"` // 要删除的用户ID列表
}

// 用户组VO
type BastionUserGroupVo struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	MemberIds   []uint     `json:"memberIds"`   // 成员ID列表
	MemberCount int        `json:"memberCount"` // 用户组成员数量
	HostCount   int        `json:"hostCount"`   // 已授权主机数量
	Status      int        `json:"status"`
	CreateTime  util.HTime `json:"createTime"`
	UpdateTime  util.HTime `json:"updateTime"`
}

// 用户组详情VO（包含成员详细信息）
type BastionUserGroupDetailVo struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	MemberIds   []uint       `json:"memberIds"`
	Members     []MemberInfo `json:"members"` // 成员详细信息
	HostCount   int          `json:"hostCount"`
	Status      int          `json:"status"`
	CreateTime  util.HTime   `json:"createTime"`
	UpdateTime  util.HTime   `json:"updateTime"`
}

// 成员信息
type MemberInfo struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
