// 服务查询 Skill
package skill

import (
	"context"
	"fmt"
	"strings"

	baseSkill "dodevops-api/pkg/agent/skill"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
)

// ServiceQuerySkill 服务查询技能
type ServiceQuerySkill struct {
	baseSkill.BaseSkill
	hostDao cmdbDao.CmdbHostDao
}

// NewServiceQuerySkill 创建服务查询技能实例
func NewServiceQuerySkill() *ServiceQuerySkill {
	return &ServiceQuerySkill{
		hostDao: cmdbDao.NewCmdbHostDao(),
	}
}

// Name 技能名称
func (s *ServiceQuerySkill) Name() string {
	return "service_query"
}

// Description 技能描述
func (s *ServiceQuerySkill) Description() string {
	return "查询服务器上运行的服务"
}

// Category 技能分类
func (s *ServiceQuerySkill) Category() baseSkill.SkillCategory {
	return baseSkill.CategoryQuery
}

// CanHandle 判断是否能处理该意图
func (s *ServiceQuerySkill) CanHandle(intent *baseSkill.Intent) bool {
	return intent.Type == "incident" &&
		intent.Entities["checkType"] == "service" &&
		intent.Entities["actionType"] == ""
}

// Validate 验证参数
func (s *ServiceQuerySkill) Validate(req *baseSkill.SkillRequest) error {
	targetName := req.Intent.Entities["targetName"]
	if targetName == "" {
		return fmt.Errorf("请指定要查询的主机IP或名称")
	}
	return nil
}

// Execute 执行服务查询
func (s *ServiceQuerySkill) Execute(ctx context.Context, req *baseSkill.SkillRequest) (*baseSkill.SkillResponse, error) {
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
			return nil, fmt.Errorf("查询失败: %v", err)
		}
	}

	return &baseSkill.SkillResponse{
		Message: fmt.Sprintf("正在查询主机 %s (%s) 上运行的服务...", host.HostName, host.SSHIP),
		Type:    "text",
		TaskInfo: map[string]interface{}{
			"taskType": "service_query",
			"hostId":   host.ID,
			"hostName": host.HostName,
			"hostIP":   host.SSHIP,
		},
	}, nil
}
