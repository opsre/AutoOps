package model

import "dodevops-api/common/util"

// EsInstance Elasticsearch 实例模型
type EsInstance struct {
    ID                 uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
    Code               string     `gorm:"column:code;type:varchar(36);comment:'实例编码';uniqueIndex;NOT NULL" json:"code"`
    Name               string     `gorm:"column:name;type:varchar(100);comment:'实例名称';NOT NULL" json:"name"`
    Protocol           string     `gorm:"column:protocol;type:varchar(10);comment:'协议:http|https';default:'http'" json:"protocol"`
    Host               string     `gorm:"column:host;type:varchar(255);comment:'主机';NOT NULL" json:"host"`
    Port               int        `gorm:"column:port;comment:'端口';default:9200" json:"port"`
    Username           string     `gorm:"column:username;type:varchar(100);comment:'用户名(可选)'" json:"username"`
    Password           string     `gorm:"column:password;type:varchar(500);comment:'密码(加密)';default:''" json:"password"`
    Remark             string     `gorm:"column:remark;type:varchar(500);comment:'备注'" json:"remark"`
    SshTunnelMachineId uint       `gorm:"column:ssh_tunnel_machine_id;comment:'SSH隧道机器ID';default:0" json:"sshTunnelMachineId"`
    Status             int        `gorm:"column:status;comment:'状态:1->启用,2->禁用';default:1" json:"status"`
    CreateTime         util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
    UpdateTime         util.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
    Creator            string     `gorm:"column:creator;type:varchar(64);comment:'创建人'" json:"creator"`
    CreatorId          uint       `gorm:"column:creator_id;comment:'创建人ID'" json:"creatorId"`
    Modifier           string     `gorm:"column:modifier;type:varchar(64);comment:'修改人'" json:"modifier"`
    ModifierId         uint       `gorm:"column:modifier_id;comment:'修改人ID'" json:"modifierId"`
}

func (EsInstance) TableName() string { return "db_es_instance" }

// 查询条件
type EsInstanceQuery struct {
    Name   string `json:"name"`
    Host   string `json:"host"`
    Status int    `json:"status"`

    PageNum  int `json:"pageNum"`
    PageSize int `json:"pageSize"`
}

// DTO 用于创建/更新
type EsInstanceDto struct {
    ID                 uint   `json:"id"`
    Code               string `json:"code"`
    Name               string `json:"name"`
    Protocol           string `json:"protocol"`
    Host               string `json:"host"`
    Port               int    `json:"port"`
    Username           string `json:"username"`
    Password           string `json:"password"`
    Remark             string `json:"remark"`
    SshTunnelMachineId uint   `json:"sshTunnelMachineId"`
}

