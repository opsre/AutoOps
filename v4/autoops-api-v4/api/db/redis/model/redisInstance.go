package model

import "dodevops-api/common/util"

// RedisInstance Redis 实例模型
type RedisInstance struct {
    ID                 uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
    Code               string     `gorm:"column:code;type:varchar(36);comment:'实例编码';uniqueIndex;NOT NULL" json:"code"`
    Name               string     `gorm:"column:name;type:varchar(100);comment:'实例名称';NOT NULL" json:"name"`
    Mode               string     `gorm:"column:mode;type:varchar(20);comment:'模式:standalone,cluster,sentinel';NOT NULL" json:"mode"`
    Host               string     `gorm:"column:host;type:varchar(300);comment:'主机: standalone为host:port, cluster为逗号分隔, sentinel为master=hosts';NOT NULL" json:"host"`
    Port               int        `gorm:"column:port;comment:'端口(standalone可用)';default:0" json:"port"`
    DB                 int        `gorm:"column:db;comment:'默认库号';default:0" json:"db"`
    Username           string     `gorm:"column:username;type:varchar(100);comment:'用户名(可选)'" json:"username"`
    Password           string     `gorm:"column:password;type:varchar(500);comment:'密码(加密)';NOT NULL" json:"password"`
    RedisNodePassword  string     `gorm:"column:redis_node_password;type:varchar(500);comment:'节点密码(仅sentinel)';default:''" json:"redisNodePassword"`
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

func (RedisInstance) TableName() string { return "db_redis_instance" }

// RedisInstanceQuery 查询条件
type RedisInstanceQuery struct {
    Name   string `json:"name"`
    Host   string `json:"host"`
    Mode   string `json:"mode"`
    Status int    `json:"status"`

    PageNum  int `json:"pageNum"`
    PageSize int `json:"pageSize"`
}

// RedisInstanceDto 创建/更新 DTO
type RedisInstanceDto struct {
    ID                 uint   `json:"id"`
    Code               string `json:"code"`
    Name               string `json:"name"`
    Mode               string `json:"mode"`
    Host               string `json:"host"`
    Port               int    `json:"port"`
    DB                 int    `json:"db"`
    Username           string `json:"username"`
    Password           string `json:"password"`
    RedisNodePassword  string `json:"redisNodePassword"`
    Remark             string `json:"remark"`
    SshTunnelMachineId uint   `json:"sshTunnelMachineId"`
}

