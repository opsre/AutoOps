// CPU分析 Skill
package skill

import (
	"context"
	"fmt"
	"strings"

	baseSkill "dodevops-api/pkg/agent/skill"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
)

// CPUAnalysisSkill CPU分析技能
type CPUAnalysisSkill struct {
	baseSkill.BaseSkill
	hostDao cmdbDao.CmdbHostDao
}

// NewCPUAnalysisSkill 创建CPU分析技能实例
func NewCPUAnalysisSkill() *CPUAnalysisSkill {
	return &CPUAnalysisSkill{
		hostDao: cmdbDao.NewCmdbHostDao(),
	}
}

// Name 技能名称
func (s *CPUAnalysisSkill) Name() string {
	return "cpu_analysis"
}

// Description 技能描述
func (s *CPUAnalysisSkill) Description() string {
	return "分析服务器CPU使用情况"
}

// Category 技能分类
func (s *CPUAnalysisSkill) Category() baseSkill.SkillCategory {
	return baseSkill.CategoryAnalysis
}

// CanHandle 判断是否能处理该意图
func (s *CPUAnalysisSkill) CanHandle(intent *baseSkill.Intent) bool {
	return intent.Type == "incident" && intent.Entities["checkType"] == "cpu"
}

// Validate 验证参数
func (s *CPUAnalysisSkill) Validate(req *baseSkill.SkillRequest) error {
	// 检查 TargetList（新架构）
	if len(req.TargetList) > 0 {
		return nil
	}

	// 兼容旧架构：检查 targetName
	targetName := req.Intent.Entities["targetName"]
	if targetName == "" {
		return fmt.Errorf("请指定要分析的主机IP或名称")
	}
	return nil
}

// Execute 执行CPU分析，支持批量处理多台主机
func (s *CPUAnalysisSkill) Execute(ctx context.Context, req *baseSkill.SkillRequest) (*baseSkill.SkillResponse, error) {
	// 检查是否需要深度分析（"为什么"）
	needAnalysis := req.Intent.Entities["needAnalysis"] == "true"

	// 支持批量处理：如果 TargetList 有多个主机，批量处理
	targetList := req.TargetList
	if len(targetList) == 0 {
		// 兼容旧架构：从 entities 中提取
		targetName := req.Intent.Entities["targetName"]
		if targetName != "" {
			targetList = []string{targetName}
		}
	}

	if len(targetList) == 0 {
		return nil, fmt.Errorf("请指定要分析的主机IP或名称")
	}

	// 单主机处理（快速路径）
	if len(targetList) == 1 {
		host, err := s.getHostInfo(targetList[0])
		if err != nil {
			return nil, err
		}

		return &baseSkill.SkillResponse{
			Message: fmt.Sprintf("正在分析主机 %s (%s) 的CPU使用情况...", host.HostName, host.SSHIP),
			Type:    "text",
			TaskInfo: map[string]interface{}{
				"taskType":     "cpu_analysis",
				"hostId":       host.ID,
				"hostName":     host.HostName,
				"hostIP":       host.SSHIP,
				"needAnalysis": needAnalysis,
			},
		}, nil
	}

	// 多主机批量处理
	var batchResults []map[string]interface{}
	var hostNames []string

	for _, target := range targetList {
		host, err := s.getHostInfo(target)
		if err != nil {
			// 记录错误但不中断
			batchResults = append(batchResults, map[string]interface{}{
				"target": target,
				"error":  err.Error(),
			})
			continue
		}

		batchResults = append(batchResults, map[string]interface{}{
			"taskType":     "cpu_analysis",
			"hostId":       host.ID,
			"hostName":     host.HostName,
			"hostIP":       host.SSHIP,
			"needAnalysis": needAnalysis,
		})
		hostNames = append(hostNames, fmt.Sprintf("%s(%s)", host.HostName, host.SSHIP))
	}

	return &baseSkill.SkillResponse{
		Message:     fmt.Sprintf("正在批量分析 %d 台主机的CPU使用情况：%s", len(batchResults), strings.Join(hostNames, ", ")),
		Type:        "batch",
		BatchResult: batchResults,
		Metadata: map[string]interface{}{
			"totalHosts": len(targetList),
			"successfulHosts": len(batchResults),
		},
	}, nil
}

// getHostInfo 根据IP或名称获取主机信息
func (s *CPUAnalysisSkill) getHostInfo(target string) (*cmdbModel.CmdbHost, error) {
	// 先尝试按IP查询
	hostByIP, err := s.hostDao.GetCmdbHostByIP(target)
	if err == nil {
		return &hostByIP, nil
	}

	// 再尝试按名称查询
	hostByName, err := s.hostDao.GetCmdbHostByName(target)
	if err == nil {
		return &hostByName, nil
	}

	// 都失败，返回友好错误
	if strings.Contains(err.Error(), "record not found") {
		return nil, fmt.Errorf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", target)
	}
	return nil, fmt.Errorf("查找主机 %s 失败: %v", target, err)
}
