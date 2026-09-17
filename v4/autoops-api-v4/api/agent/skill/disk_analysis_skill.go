// 磁盘分析 Skill
package skill

import (
	"context"
	"fmt"
	"strings"

	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	baseSkill "dodevops-api/pkg/agent/skill"
)

// DiskAnalysisSkill 磁盘分析技能
type DiskAnalysisSkill struct {
	baseSkill.BaseSkill
	hostDao cmdbDao.CmdbHostDao
}

// NewDiskAnalysisSkill 创建磁盘分析技能实例
func NewDiskAnalysisSkill() *DiskAnalysisSkill {
	return &DiskAnalysisSkill{
		hostDao: cmdbDao.NewCmdbHostDao(),
	}
}

// Name 技能名称
func (s *DiskAnalysisSkill) Name() string {
	return "disk_analysis"
}

// Description 技能描述
func (s *DiskAnalysisSkill) Description() string {
	return "分析服务器磁盘使用情况"
}

// Category 技能分类
func (s *DiskAnalysisSkill) Category() baseSkill.SkillCategory {
	return baseSkill.CategoryAnalysis
}

// CanHandle 判断是否能处理该意图
func (s *DiskAnalysisSkill) CanHandle(intent *baseSkill.Intent) bool {
	queryType := intent.Entities["queryType"]
	return intent.Type == "incident" &&
		intent.Entities["checkType"] == "disk" &&
		(queryType == "" || queryType == "directories") &&
		intent.Entities["actionType"] == "" &&
		intent.Entities["fileOp"] == ""
}

// Validate 验证参数
func (s *DiskAnalysisSkill) Validate(req *baseSkill.SkillRequest) error {
	targetName := req.Intent.Entities["targetName"]
	if targetName == "" {
		return fmt.Errorf("请指定要分析的主机IP或名称")
	}
	return nil
}

// Execute 执行磁盘分析
func (s *DiskAnalysisSkill) Execute(ctx context.Context, req *baseSkill.SkillRequest) (*baseSkill.SkillResponse, error) {
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

	// 检查是否需要查看目录详情
	showDirectories := req.Intent.Entities["queryType"] == "directories"

	return &baseSkill.SkillResponse{
		Message: fmt.Sprintf("正在分析主机 %s (%s) 的磁盘使用情况...", host.HostName, host.SSHIP),
		Type:    "text",
		TaskInfo: map[string]interface{}{
			"taskType":        "disk_analysis",
			"hostId":          host.ID,
			"hostName":        host.HostName,
			"hostIP":          host.SSHIP,
			"showDirectories": showDirectories,
			"queryType":       req.Intent.Entities["queryType"],
		},
	}, nil
}
