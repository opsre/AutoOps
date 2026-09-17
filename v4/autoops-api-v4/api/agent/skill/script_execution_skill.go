// 脚本执行 Skill
package skill

import (
	"context"
	"fmt"
	"strings"

	baseSkill "dodevops-api/pkg/agent/skill"
	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	monitorService "dodevops-api/api/monitor/service"
	"dodevops-api/common/util"
)

// ScriptExecutionSkill 脚本执行技能
type ScriptExecutionSkill struct {
	baseSkill.BaseSkill
	hostDao    cmdbDao.CmdbHostDao
	ecsAuthDao monitorService.EcsAuthDao
	sshUtil    *util.SSHUtil
}

// NewScriptExecutionSkill 创建脚本执行技能实例
func NewScriptExecutionSkill() *ScriptExecutionSkill {
	return &ScriptExecutionSkill{
		hostDao:    cmdbDao.NewCmdbHostDao(),
		ecsAuthDao: monitorService.NewEcsAuthDao(),
		sshUtil:    util.NewSSHUtil(),
	}
}

// Name 技能名称
func (s *ScriptExecutionSkill) Name() string {
	return "script_execution"
}

// Description 技能描述
func (s *ScriptExecutionSkill) Description() string {
	return "在指定服务器上执行Shell脚本"
}

// Category 技能分类
func (s *ScriptExecutionSkill) Category() baseSkill.SkillCategory {
	return baseSkill.CategoryOperation
}

// CanHandle 判断是否能处理该意图
func (s *ScriptExecutionSkill) CanHandle(intent *baseSkill.Intent) bool {
	// 检测"执行脚本"、"运行脚本"等关键词
	actionType := intent.Entities["actionType"]
	return (intent.Type == "incident" || intent.Type == "query") &&
		(actionType == "execute_script" || actionType == "run_script")
}

// Validate 验证参数
func (s *ScriptExecutionSkill) Validate(req *baseSkill.SkillRequest) error {
	// 检查目标列表
	if len(req.TargetList) == 0 {
		targetName := req.Intent.Entities["targetName"]
		if targetName == "" {
			return fmt.Errorf("请指定要执行脚本的主机IP或名称")
		}
	}

	// 检查脚本
	script := req.Intent.Entities["script"]
	if script == "" {
		return fmt.Errorf("请指定要执行的脚本内容或脚本路径")
	}

	return nil
}

// Execute 执行脚本
func (s *ScriptExecutionSkill) Execute(ctx context.Context, req *baseSkill.SkillRequest) (*baseSkill.SkillResponse, error) {
	script := req.Intent.Entities["script"]
	scriptType := req.Intent.Entities["scriptType"] // "inline" 或 "file"

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
		return nil, fmt.Errorf("请指定要执行脚本的主机IP或名称")
	}

	// 单主机处理（快速路径）
	if len(targetList) == 1 {
		host, err := s.getHostInfo(targetList[0])
		if err != nil {
			return nil, err
		}

		// 执行脚本
		result, err := s.executeScript(host, script, scriptType)
		if err != nil {
			return nil, err
		}

		return &baseSkill.SkillResponse{
			Message: fmt.Sprintf("✅ 脚本在主机 %s (%s) 上执行完成\n\n**输出**:\n```\n%s\n```",
				host.HostName, host.SSHIP, result),
			Type: "text",
			TaskInfo: map[string]interface{}{
				"taskType": "script_execution",
				"hostId":   host.ID,
				"hostName": host.HostName,
				"hostIP":   host.SSHIP,
				"script":   script,
				"result":   result,
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

		// 执行脚本
		result, err := s.executeScript(host, script, scriptType)
		if err != nil {
			batchResults = append(batchResults, map[string]interface{}{
				"hostName": host.HostName,
				"hostIP":   host.SSHIP,
				"error":    err.Error(),
			})
			continue
		}

		batchResults = append(batchResults, map[string]interface{}{
			"taskType": "script_execution",
			"hostId":   host.ID,
			"hostName": host.HostName,
			"hostIP":   host.SSHIP,
			"script":   script,
			"result":   result,
		})
		hostNames = append(hostNames, fmt.Sprintf("%s(%s)", host.HostName, host.SSHIP))
	}

	return &baseSkill.SkillResponse{
		Message:     fmt.Sprintf("正在在 %d 台主机上执行脚本：%s", len(batchResults), strings.Join(hostNames, ", ")),
		Type:        "batch",
		BatchResult: batchResults,
		Metadata: map[string]interface{}{
			"totalHosts": len(targetList),
		},
	}, nil
}

// getHostInfo 根据IP或名称获取主机信息
func (s *ScriptExecutionSkill) getHostInfo(target string) (*cmdbModel.CmdbHost, error) {
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

// executeScript 在主机上执行脚本
func (s *ScriptExecutionSkill) executeScript(host *cmdbModel.CmdbHost, script string, scriptType string) (string, error) {
	// 获取SSH认证信息
	ecsAuth, err := s.ecsAuthDao.GetEcsAuthById(host.SSHKeyID)
	if err != nil {
		return "", fmt.Errorf("获取SSH认证信息失败: %v", err)
	}

	// 构建SSH配置
	sshConfig := &util.SSHConfig{
		IP:       host.SSHIP,
		Port:     host.SSHPort,
		Username: host.SSHName,
		Type:     ecsAuth.Type,
	}

	// 根据认证类型设置认证信息
	switch ecsAuth.Type {
	case 1: // 密码认证
		sshConfig.Password = ecsAuth.Password
	case 2: // 私钥认证
		sshConfig.PublicKey = ecsAuth.PublicKey
	default:
		return "", fmt.Errorf("不支持的认证类型: %d", ecsAuth.Type)
	}

	var result string

	if scriptType == "file" {
		// 执行脚本文件
		result, err = s.sshUtil.ExecuteRemoteCommand(sshConfig, fmt.Sprintf("bash %s", script))
	} else {
		// 执行内联脚本
		result, err = s.sshUtil.ExecuteScript(sshConfig, script)
	}

	if err != nil {
		return "", fmt.Errorf("脚本执行失败: %v", err)
	}

	return result, nil
}
