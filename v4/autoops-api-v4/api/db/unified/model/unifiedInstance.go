package model

import (
	"dodevops-api/common/util"
	"encoding/json"
	"strings"
)

// UnifiedInstance 统一数据库实例模型（对应视图 view_db_instance_all）
type UnifiedInstance struct {
	ID                 uint            `gorm:"column:id" json:"id"`
	Code               string          `gorm:"column:code" json:"code"`
	Name               string          `gorm:"column:name" json:"name"`
	DbType             string          `gorm:"column:db_type" json:"dbType"`                         // mysql, postgres, redis, mongodb, elasticsearch
	SubType            *string         `gorm:"column:sub_type" json:"subType"`                       // 子类型：mysql/postgres, standalone/cluster, mongodb/atlas, http/https
	Host               *string         `gorm:"column:host" json:"host"`                              // 主机地址（MongoDB 为空）
	Port               *int            `gorm:"column:port" json:"port"`                              // 端口（MongoDB 为空）
	Username           *string         `gorm:"column:username" json:"username"`                      // 用户名
	Password           *string         `gorm:"column:password" json:"password"`                      // 密码（脱敏）
	Remark             *string         `gorm:"column:remark" json:"remark"`                          // 备注
	SshTunnelMachineId uint            `gorm:"column:ssh_tunnel_machine_id" json:"sshTunnelMachineId"` // SSH隧道机器ID
	Status             int             `gorm:"column:status" json:"status"`                          // 状态
	CreateTime         util.HTime      `gorm:"column:create_time" json:"createTime"`                 // 创建时间
	UpdateTime         *util.HTime     `gorm:"column:update_time" json:"updateTime"`                 // 更新时间
	Creator            *string         `gorm:"column:creator" json:"creator"`                        // 创建人
	CreatorId          *uint           `gorm:"column:creator_id" json:"creatorId"`                   // 创建人ID
	Modifier           *string         `gorm:"column:modifier" json:"modifier"`                      // 修改人
	ModifierId         *uint           `gorm:"column:modifier_id" json:"modifierId"`                 // 修改人ID
	ConnectionConfig   json.RawMessage `gorm:"column:connection_config" json:"connectionConfig"`     // 连接配置（JSON）
}

// TableName 指定表名（实际是视图名）
func (UnifiedInstance) TableName() string {
	return "db_instance_all"
}

// UnifiedInstanceQuery 统一查询条件
type UnifiedInstanceQuery struct {
	DbType   string `json:"dbType"`   // 数据库类型过滤
	Name     string `json:"name"`     // 名称模糊搜索
	Status   *int   `json:"status"`   // 状态过滤
	PageNum  int    `json:"pageNum"`  // 页码
	PageSize int    `json:"pageSize"` // 每页数量
}

// UnifiedInstanceVO 统一实例视图对象（返回给前端）
type UnifiedInstanceVO struct {
	ID                 uint                   `json:"id"`
	Code               string                 `json:"code"`
	Name               string                 `json:"name"`
	DbType             string                 `json:"dbType"`
	SubType            string                 `json:"subType,omitempty"`
	Host               string                 `json:"host,omitempty"`
	Port               int                    `json:"port,omitempty"`
	Username           string                 `json:"username,omitempty"`
	Password           string                 `json:"password,omitempty"` // 始终脱敏为 ******
	Remark             string                 `json:"remark,omitempty"`
	SshTunnelMachineId uint                   `json:"sshTunnelMachineId"`
	Status             int                    `json:"status"`
	CreateTime         util.HTime             `json:"createTime"`
	UpdateTime         util.HTime             `json:"updateTime,omitempty"`
	Creator            string                 `json:"creator,omitempty"`
	CreatorId          uint                   `json:"creatorId,omitempty"`
	Modifier           string                 `json:"modifier,omitempty"`
	ModifierId         uint                   `json:"modifierId,omitempty"`
	ConnectionConfig   map[string]interface{} `json:"connectionConfig,omitempty"` // 解析后的 JSON

	// 原始表 ID 信息（用于跳转到详情页面）
	SourceTable        string                 `json:"sourceTable"` // 原始表名
	SourceId           uint                   `json:"sourceId"`    // 原始表 ID
}

// ToVO 将 UnifiedInstance 转换为 VO
func (u *UnifiedInstance) ToVO() UnifiedInstanceVO {
	vo := UnifiedInstanceVO{
		ID:                 u.ID,
		Code:               u.Code,
		Name:               u.Name,
		DbType:             u.DbType,
		SshTunnelMachineId: u.SshTunnelMachineId,
		Status:             u.Status,
		CreateTime:         u.CreateTime,
		Password:           "******", // 始终脱敏
		SourceId:           u.ID,
	}

	// 处理可空字段
	if u.SubType != nil {
		vo.SubType = *u.SubType
	}
	if u.Host != nil {
		vo.Host = *u.Host
	}
	if u.Port != nil {
		vo.Port = *u.Port
	}
	if u.Username != nil {
		vo.Username = *u.Username
	}
	if u.Remark != nil {
		vo.Remark = *u.Remark
	}
	if u.UpdateTime != nil {
		vo.UpdateTime = *u.UpdateTime
	}
	if u.Creator != nil {
		vo.Creator = *u.Creator
	}
	if u.CreatorId != nil {
		vo.CreatorId = *u.CreatorId
	}
	if u.Modifier != nil {
		vo.Modifier = *u.Modifier
	}
	if u.ModifierId != nil {
		vo.ModifierId = *u.ModifierId
	}

	// 解析 ConnectionConfig JSON 并脱敏处理
	if len(u.ConnectionConfig) > 0 {
		var config map[string]interface{}
		if err := json.Unmarshal(u.ConnectionConfig, &config); err == nil {
			// 特殊处理 MongoDB URI - 脱敏并提取 username
			if u.DbType == "mongodb" {
				if uri, ok := config["uri"].(string); ok {
					// 解析 MongoDB URI 提取 username 和 host:port
					username, hostPort := parseMongoURI(uri)
					if username != "" {
						vo.Username = username
					}
					if hostPort != "" {
						// 将脱敏后的 host:port 放回 config
						config["uri"] = hostPort
					}
				}
			}
			vo.ConnectionConfig = config
		}
	}

	// 确定原始表名
	switch u.DbType {
	case "mysql", "postgres":
		vo.SourceTable = "db_instance"
	case "redis":
		vo.SourceTable = "db_redis_instance"
	case "mongodb":
		vo.SourceTable = "db_mongo_instance"
	case "elasticsearch":
		vo.SourceTable = "db_es_instance"
	}

	return vo
}

// parseMongoURI 解析 MongoDB URI，提取 username 和 host:port，去除密码
// 输入示例: mongodb://admin:password@192.168.2.123:27018
// 返回: username="admin", hostPort="192.168.2.123:27018"
func parseMongoURI(uri string) (username, hostPort string) {
	// 移除协议前缀
	uri = strings.TrimPrefix(uri, "mongodb://")
	uri = strings.TrimPrefix(uri, "mongodb+srv://")

	// 查找 @ 符号，分离认证信息和主机信息
	atIndex := strings.Index(uri, "@")
	if atIndex > 0 {
		// 有认证信息
		authPart := uri[:atIndex]
		hostPart := uri[atIndex+1:]

		// 从认证部分提取 username（冒号前的部分）
		colonIndex := strings.Index(authPart, ":")
		if colonIndex > 0 {
			username = authPart[:colonIndex]
		} else {
			username = authPart
		}

		// 从 host 部分提取 host:port（去掉数据库名和参数）
		slashIndex := strings.Index(hostPart, "/")
		questionIndex := strings.Index(hostPart, "?")

		if slashIndex > 0 {
			hostPort = hostPart[:slashIndex]
		} else if questionIndex > 0 {
			hostPort = hostPart[:questionIndex]
		} else {
			hostPort = hostPart
		}
	} else {
		// 无认证信息，直接提取 host:port
		slashIndex := strings.Index(uri, "/")
		questionIndex := strings.Index(uri, "?")

		if slashIndex > 0 {
			hostPort = uri[:slashIndex]
		} else if questionIndex > 0 {
			hostPort = uri[:questionIndex]
		} else {
			hostPort = uri
		}
	}

	return
}
