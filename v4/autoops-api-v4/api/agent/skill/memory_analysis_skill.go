// 内存分析 Skill
package skill

import (
	"context"
	"fmt"
	"strings"

	baseSkill "dodevops-api/pkg/agent/skill"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
)

// MemoryAnalysisSkill 内存分析技能
type MemoryAnalysisSkill struct {
	baseSkill.BaseSkill
	hostDao cmdbDao.CmdbHostDao
}

// NewMemoryAnalysisSkill 创建内存分析技能实例
func NewMemoryAnalysisSkill() *MemoryAnalysisSkill {
	return &MemoryAnalysisSkill{
		hostDao: cmdbDao.NewCmdbHostDao(),
	}
}

// Name 技能名称
func (s *MemoryAnalysisSkill) Name() string {
	return "memory_analysis"
}

// Description 技能描述
func (s *MemoryAnalysisSkill) Description() string {
	return "分析服务器内存使用情况，识别内存问题"
}

// Category 技能分类
func (s *MemoryAnalysisSkill) Category() baseSkill.SkillCategory {
	return baseSkill.CategoryAnalysis
}

// CanHandle 判断是否能处理该意图
func (s *MemoryAnalysisSkill) CanHandle(intent *baseSkill.Intent) bool {
	return intent.Type == "incident" && intent.Entities["checkType"] == "memory"
}

// Validate 验证参数
func (s *MemoryAnalysisSkill) Validate(req *baseSkill.SkillRequest) error {
	targetName := req.Intent.Entities["targetName"]
	if targetName == "" {
		return fmt.Errorf("请指定要分析的主机IP或名称")
	}
	return nil
}

// Execute 执行内存分析
func (s *MemoryAnalysisSkill) Execute(ctx context.Context, req *baseSkill.SkillRequest) (*baseSkill.SkillResponse, error) {
	targetName := req.Intent.Entities["targetName"]

	// 获取主机信息
	var host *cmdbModel.CmdbHost
	hostByIP, err := s.hostDao.GetCmdbHostByIP(targetName)
	if err == nil {
		host = &hostByIP
	} else {
		hostByName, err := s.hostDao.GetCmdbHostByName(targetName)
		if err == nil {
			host = &hostByName
		} else {
			if strings.Contains(err.Error(), "record not found") {
				return nil, fmt.Errorf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", targetName)
			}
			return nil, fmt.Errorf("分析失败: %v", err)
		}
	}

	return &baseSkill.SkillResponse{
		Message: fmt.Sprintf("正在分析主机 %s (%s) 的内存情况...", host.HostName, host.SSHIP),
		Type:    "text",
		TaskInfo: map[string]interface{}{
			"taskType": "memory_analysis",
			"hostId":   host.ID,
			"hostName": host.HostName,
			"hostIP":   host.SSHIP,
		},
	}, nil
}
