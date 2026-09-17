package model

// CloudCost 云厂商费用结构
type CloudCost struct {
	BusinessLine     string  `json:"business_line"`      // 业务线/服务名称
	Amount           float64 `json:"amount"`             // 费用金额（人民币）
	OriginalAmount   float64 `json:"original_amount"`    // 原始金额
	OriginalCurrency string  `json:"original_currency"`  // 原始货币
}

// TencentBillingRequest 腾讯云账单请求
type TencentBillingRequest struct {
	KeyID uint   `json:"keyId" form:"keyId" validate:"required"` // 密钥ID
	Month string `json:"month" form:"month"`                     // 账期(格式: 2024-11, 为空则查询上月)
}

// TencentBillingResponse 腾讯云账单响应
type TencentBillingResponse struct {
	Month      string      `json:"month"`       // 账期
	Costs      []CloudCost `json:"costs"`       // 费用明细
	TotalCost  float64     `json:"total_cost"`  // 总费用
}

// AliyunProductBillingRequest 阿里云账单（按产品）请求
type AliyunProductBillingRequest struct {
	KeyID uint   `json:"keyId" form:"keyId" validate:"required"` // 密钥ID
	Month string `json:"month" form:"month"`                     // 账期(格式: 2024-11, 为空则查询上月)
}

// AliyunProductBillingResponse 阿里云账单（按产品）响应
type AliyunProductBillingResponse struct {
	Month     string      `json:"month"`      // 账期
	Costs     []CloudCost `json:"costs"`      // 费用明细（按产品）
	TotalCost float64     `json:"total_cost"` // 总费用
}

// AliyunTagBillingRequest 阿里云账单（按标签）请求
type AliyunTagBillingRequest struct {
	KeyID uint   `json:"keyId" form:"keyId" validate:"required"` // 密钥ID
	Month string `json:"month" form:"month"`                     // 账期(格式: 2024-11, 为空则查询上月)
}

// AliyunTagBillingResponse 阿里云账单（按标签）响应
type AliyunTagBillingResponse struct {
	Month     string      `json:"month"`      // 账期
	Costs     []CloudCost `json:"costs"`      // 费用明细（按标签bline）
	TotalCost float64     `json:"total_cost"` // 总费用
}

// AWSBillingRequest AWS账单请求
type AWSBillingRequest struct {
	KeyID uint   `json:"keyId" form:"keyId" validate:"required"` // 密钥ID
	Month string `json:"month" form:"month"`                     // 账期(格式: 2024-11, 为空则查询上月)
}

// AWSBillingResponse AWS账单响应
type AWSBillingResponse struct {
	Month        string      `json:"month"`         // 账期
	Costs        []CloudCost `json:"costs"`         // 费用明细
	TotalCost    float64     `json:"total_cost"`    // 总费用（人民币）
	TotalCostUSD float64     `json:"total_cost_usd"` // 总费用（美元）
	ExchangeRate float64     `json:"exchange_rate"` // 美元汇率
}

// CloudSummaryRequest 云费用汇总请求
type CloudSummaryRequest struct {
	TencentKeyID uint `json:"tencentKeyId" form:"tencentKeyId"` // 腾讯云密钥ID（可选）
	AliyunKeyID  uint `json:"aliyunKeyId" form:"aliyunKeyId"`  // 阿里云密钥ID（可选）
	AWSKeyID     uint `json:"awsKeyId" form:"awsKeyId"`     // AWS密钥ID（可选）
}

// CloudCostSummary 单个云厂商费用汇总
type CloudCostSummary struct {
	CloudProvider string  `json:"cloud_provider"` // 云厂商名称
	CurrentMonth  float64 `json:"current_month"`  // 上月费用（CNY）
	PreviousMonth float64 `json:"previous_month"` // 上上月费用（CNY）
	MonthOnMonth  float64 `json:"month_on_month"` // 环比（%）
	Currency      string  `json:"currency"`       // 货币单位
	CostUSD       float64 `json:"cost_usd,omitempty"` // 美元金额（仅AWS）
}

// LastMonthSummaryResponse 上月费用汇总响应
type LastMonthSummaryResponse struct {
	CurrentMonth  string             `json:"current_month"`  // 当前统计月份
	PreviousMonth string             `json:"previous_month"` // 对比月份
	Summary       []CloudCostSummary `json:"summary"`        // 各云厂商汇总
}

// MonthlyTrend 单月费用趋势
type MonthlyTrend struct {
	Month    string  `json:"month"`     // 月份 (2024-01)
	Amount   float64 `json:"amount"`    // 金额 (CNY)
	AmountUSD float64 `json:"amount_usd,omitempty"` // 美元金额（仅AWS）
}

// CloudTrendData 单个云厂商趋势数据
type CloudTrendData struct {
	CloudProvider string         `json:"cloud_provider"` // 云厂商名称
	Trends        []MonthlyTrend `json:"trends"`         // 月度趋势
	Currency      string         `json:"currency"`       // 货币单位
}

// Last12MonthsTrendResponse 最近12个月趋势响应
type Last12MonthsTrendResponse struct {
	StartMonth string           `json:"start_month"` // 开始月份
	EndMonth   string           `json:"end_month"`   // 结束月份
	Data       []CloudTrendData `json:"data"`        // 各云厂商趋势数据
}
