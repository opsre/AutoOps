package model

// JumpServer同步主机DTO
type SyncJumpServerHostDto struct {
	URL       string `validate:"required" json:"url"`       // JumpServer URL
	AccessKey string `validate:"required" json:"accessKey"` // Access Key ID
	SecretKey string `validate:"required" json:"secretKey"` // Secret Key
}

// JumpServer API响应结构
type JumpServerAssetsResponse struct {
	Count    int                  `json:"count"`
	Next     *string              `json:"next"`
	Previous *string              `json:"previous"`
	Results  []JumpServerAssetDto `json:"results"`
}

// JumpServer资产数据
type JumpServerAssetDto struct {
	ID       string                   `json:"id"`
	Name     string                   `json:"name"`
	Address  string                   `json:"address"`
	Comment  string                   `json:"comment"`
	Platform JumpServerPlatformDto    `json:"platform"`
	Nodes    []JumpServerNodeDto      `json:"nodes"`
	Protocols []JumpServerProtocolDto `json:"protocols"`
}

// JumpServer平台信息
type JumpServerPlatformDto struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// JumpServer节点(业务分组)信息
type JumpServerNodeDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// JumpServer协议信息
type JumpServerProtocolDto struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

// JumpServer同步结果
type SyncJumpServerResult struct {
	Total       int      `json:"total"`       // 总数
	Success     int      `json:"success"`     // 成功数
	Failed      int      `json:"failed"`      // 失败数
	Skipped     int      `json:"skipped"`     // 跳过数(已存在)
	FailedHosts []string `json:"failedHosts"` // 失败主机列表
	Message     string   `json:"message"`     // 消息
}
