// 压力分析 Skill
package skill

import (
	"context"
	"fmt"
	"strings"

	baseSkill "dodevops-api/pkg/agent/skill"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
)

// PressureAnalysisSkill 服务器压力分析技能
type PressureAnalysisSkill struct {
	baseSkill.BaseSkill
	hostDao cmdbDao.CmdbHostDao
}

// NewPressureAnalysisSkill 创建压力分析技能实例
func NewPressureAnalysisSkill() *PressureAnalysisSkill {
	return &PressureAnalysisSkill{
		hostDao: cmdbDao.NewCmdbHostDao(),
	}
}

// Name 技能名称
func (s *PressureAnalysisSkill) Name() string {
	return "pressure_analysis"
}

// Description 技能描述
func (s *PressureAnalysisSkill) Description() string {
	return "分析服务器压力，包括CPU、内存、负载等指标"
}

// Category 技能分类
func (s *PressureAnalysisSkill) Category() baseSkill.SkillCategory {
	return baseSkill.CategoryAnalysis
}

// CanHandle 判断是否能处理该意图
func (s *PressureAnalysisSkill) CanHandle(intent *baseSkill.Intent) bool {
	return intent.Type == "incident" && intent.Entities["checkType"] == "performance"
}

// Validate 验证参数
func (s *PressureAnalysisSkill) Validate(req *baseSkill.SkillRequest) error {
	targetName := req.Intent.Entities["targetName"]
	if targetName == "" {
		return fmt.Errorf("请指定要分析的主机IP或名称，例如：\n- 检查 ip=192.168.1.100 的服务器压力\n- 分析 name=web-server-01 的内存问题")
	}
	return nil
}

// Execute 执行压力分析
func (s *PressureAnalysisSkill) Execute(ctx context.Context, req *baseSkill.SkillRequest) (*baseSkill.SkillResponse, error) {
	targetName := req.Intent.Entities["targetName"]

	// 获取主机信息
	var host *cmdbModel.CmdbHost
	var err error

	// 尝试按IP查询
	hostByIP, err := s.hostDao.GetCmdbHostByIP(targetName)
	if err == nil {
		host = &hostByIP
	} else {
		// 尝试按名称查询
		hostByName, err := s.hostDao.GetCmdbHostByName(targetName)
		if err == nil {
			host = &hostByName
		} else {
			// 优化错误提示
			if strings.Contains(err.Error(), "record not found") {
				return nil, fmt.Errorf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", targetName)
			}
			return nil, fmt.Errorf("分析失败: %v", err)
		}
	}

	// 调用分析服务（通过 Context 传递，避免循环依赖）
	// 这里返回一个占位响应，实际分析逻辑在 service 层完成
	return &baseSkill.SkillResponse{
		Message: fmt.Sprintf("正在分析主机 %s (%s) 的压力情况...", host.HostName, host.SSHIP),
		Type:    "text",
		TaskInfo: map[string]interface{}{
			"taskType": "pressure_analysis",
			"hostId":   host.ID,
			"hostName": host.HostName,
			"hostIP":   host.SSHIP,
		},
		Metadata: map[string]interface{}{
			"needAnalysis": true,
		},
	}, nil
}
