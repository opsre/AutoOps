package skill

import (
	"context"
	"fmt"
	"strings"

	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	baseSkill "dodevops-api/pkg/agent/skill"
)

// DirectoryDeepAnalysisSkill 目录深度分析技能
type DirectoryDeepAnalysisSkill struct {
	baseSkill.BaseSkill
	hostDao cmdbDao.CmdbHostDao
}

func NewDirectoryDeepAnalysisSkill() *DirectoryDeepAnalysisSkill {
	return &DirectoryDeepAnalysisSkill{
		hostDao: cmdbDao.NewCmdbHostDao(),
	}
}

func (s *DirectoryDeepAnalysisSkill) Name() string {
	return "directory_deep_analysis"
}

func (s *DirectoryDeepAnalysisSkill) Description() string {
	return "深度分析指定目录的子目录、大文件和可清理项"
}

func (s *DirectoryDeepAnalysisSkill) Category() baseSkill.SkillCategory {
	return baseSkill.CategoryAnalysis
}

func (s *DirectoryDeepAnalysisSkill) CanHandle(intent *baseSkill.Intent) bool {
	return intent.Type == "incident" &&
		intent.Entities["checkType"] == "disk" &&
		(intent.Entities["queryType"] == "directory_deep_analysis" || intent.Entities["queryType"] == "list_files") &&
		intent.Entities["fileOp"] == ""
}

func (s *DirectoryDeepAnalysisSkill) Validate(req *baseSkill.SkillRequest) error {
	targetName := req.Intent.Entities["targetName"]
	if targetName == "" {
		return fmt.Errorf("请先指定要分析的主机。你当前这条消息只有目录，没有主机信息；请直接带上主机IP/主机名，或保持同一个会话继续提问")
	}
	if strings.TrimSpace(req.Intent.Entities["directoryPath"]) == "" {
		return fmt.Errorf("请指定要深度分析的目录路径")
	}
	return nil
}

func (s *DirectoryDeepAnalysisSkill) Execute(ctx context.Context, req *baseSkill.SkillRequest) (*baseSkill.SkillResponse, error) {
	targetName := req.Intent.Entities["targetName"]
	directoryPath := strings.TrimSpace(req.Intent.Entities["directoryPath"])

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
		Message: fmt.Sprintf("正在深度分析主机 %s (%s) 的目录 %s ...", host.HostName, host.SSHIP, directoryPath),
		Type:    "text",
		TaskInfo: map[string]interface{}{
			"taskType":      "directory_deep_analysis",
			"hostId":        host.ID,
			"hostName":      host.HostName,
			"hostIP":        host.SSHIP,
			"directoryPath": directoryPath,
			"queryType":     req.Intent.Entities["queryType"],
		},
	}, nil
}
