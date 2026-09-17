package model

// Conn 连接参数
// Mode: standalone | cluster | sentinel
// - standalone: Host=host:port
// - cluster: Host=host1:port1,host2:port2
// - sentinel: Host=masterName=host1:port1,host2:port2, RedisNodePassword 用于节点密码
type Conn struct {
    Host              string `json:"host" binding:"required"`  // 主机地址，格式见上
    Port              int    `json:"port"`                      // 仅 standalone 使用
    Password          string `json:"password"`                  // 密码（standalone/哨兵密码/集群总密码）
    DB                int    `json:"db"`
    Mode              string `json:"mode,omitempty"`            // standalone | cluster | sentinel
    RedisNodePassword string `json:"redisNodePassword,omitempty"`
}

// ScanReq 扫描请求
type ScanReq struct {
    Conn
    // 支持 test/api 的字段命名 match；同时兼容旧字段 pattern
    Match   string            `json:"match"`
    Pattern string            `json:"pattern"` // 兼容旧字段
    Count   int               `json:"count"`
    // 集群下为每个节点维护一个 cursor；单机仍可使用旧的单值 cursor
    CursorMap map[string]uint64 `json:"cursor,omitempty"`
    Cursor    uint64            `json:"cursorSingle,omitempty"`
}

// CmdReq 通用命令
type CmdReq struct {
    Conn
    Command string   `json:"command" binding:"required"`
    Args    []string `json:"args"`
}

