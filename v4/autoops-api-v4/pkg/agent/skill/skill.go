// Skill 核心接口定义
package skill

import (
	"context"
)

// Skill 技能接口 - 所有Agent能力的统一抽象
type Skill interface {
	// 技能元信息
	Name() string                    // 技能名称，如 "inspection", "pressure_analysis"
	Description() string             // 技能描述
	Category() SkillCategory         // 技能分类

	// 意图匹配
	CanHandle(intent *Intent) bool   // 判断是否能处理该意图

	// 执行
	Execute(ctx context.Context, req *SkillRequest) (*SkillResponse, error)

	// 可选：参数验证
	Validate(req *SkillRequest) error
}

// SkillCategory 技能分类
type SkillCategory string

const (
	CategoryInspection  SkillCategory = "inspection"  // 巡检类
	CategoryAnalysis    SkillCategory = "analysis"    // 分析类
	CategoryQuery       SkillCategory = "query"       // 查询类
	CategoryOperation   SkillCategory = "operation"   // 操作类
)

// Intent 意图
type Intent struct {
	Type       string            // inspection, incident, query
	Action     string            // execute, analyze, query
	Confidence float64
	Entities   map[string]string // 提取的实体
}

// SkillRequest 统一请求
type SkillRequest struct {
	SessionID  string
	UserID     uint
	Intent     *Intent
	Message    string                 // 原始消息
	Context    map[string]interface{} // 上下文数据
	TargetList []string               // 目标列表（支持多个主机）
}

// SkillResponse 统一响应
type SkillResponse struct {
	Message     string
	Type        string  // text, task, error, batch（批量）
	TaskInfo    map[string]interface{}
	Suggestions []string
	Metadata    map[string]interface{} // 元数据（执行时间、记录数等）
	BatchResult []map[string]interface{} // 批量结果（用于多主机场景）
}
