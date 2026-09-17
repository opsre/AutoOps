package skill

import (
	"context"
	"fmt"
	"strings"

	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	monitorService "dodevops-api/api/monitor/service"
	"dodevops-api/common/util"
	baseSkill "dodevops-api/pkg/agent/skill"
)

type ServiceManageSkill struct {
	baseSkill.BaseSkill
	hostDao    cmdbDao.CmdbHostDao
	ecsAuthDao monitorService.EcsAuthDao
	sshUtil    *util.SSHUtil
}

func NewServiceManageSkill() *ServiceManageSkill {
	return &ServiceManageSkill{
		hostDao:    cmdbDao.NewCmdbHostDao(),
		ecsAuthDao: monitorService.NewEcsAuthDao(),
		sshUtil:    util.NewSSHUtil(),
	}
}

func (s *ServiceManageSkill) Name() string {
	return "service_manage"
}

func (s *ServiceManageSkill) Description() string {
	return "在指定服务器上检查、安装和重启系统服务"
}

func (s *ServiceManageSkill) Category() baseSkill.SkillCategory {
	return baseSkill.CategoryOperation
}

func (s *ServiceManageSkill) CanHandle(intent *baseSkill.Intent) bool {
	if intent.Type != "incident" || intent.Entities["checkType"] != "service" {
		return false
	}
	actionType := intent.Entities["actionType"]
	return actionType == "install_service" || actionType == "restart_service" || actionType == "query_service"
}

func (s *ServiceManageSkill) Validate(req *baseSkill.SkillRequest) error {
	if strings.TrimSpace(req.Intent.Entities["targetName"]) == "" {
		return fmt.Errorf("请指定要操作服务的主机IP或名称")
	}
	if strings.TrimSpace(req.Intent.Entities["serviceName"]) == "" {
		return fmt.Errorf("请指定服务名称，例如 nginx")
	}
	return nil
}

func (s *ServiceManageSkill) Execute(ctx context.Context, req *baseSkill.SkillRequest) (*baseSkill.SkillResponse, error) {
	host, err := s.getHostInfo(strings.TrimSpace(req.Intent.Entities["targetName"]))
	if err != nil {
		return nil, err
	}
	sshConfig, err := s.getSSHConfig(host)
	if err != nil {
		return nil, fmt.Errorf("获取SSH认证信息失败: %v", err)
	}

	serviceName := strings.TrimSpace(req.Intent.Entities["serviceName"])
	actionType := strings.TrimSpace(req.Intent.Entities["actionType"])

	var command string
	switch actionType {
	case "install_service":
		command = buildServiceInstallCommand(serviceName)
	case "restart_service":
		command = buildServiceRestartCommand(serviceName)
	case "query_service":
		command = buildServiceStatusCommand(serviceName)
	default:
		return nil, fmt.Errorf("不支持的服务操作: %s", actionType)
	}

	result, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, command)
	if err != nil {
		return nil, fmt.Errorf("服务操作失败: %v", err)
	}

	var report string
	switch actionType {
	case "install_service":
		report = formatServiceInstallResult(host, serviceName, result)
	case "restart_service":
		report = formatServiceRestartResult(host, serviceName, result)
	case "query_service":
		report = formatServiceStatusResult(host, serviceName, result)
	}

	return &baseSkill.SkillResponse{
		Message: report,
		Type:    "text",
		TaskInfo: map[string]interface{}{
			"taskType":    "service_manage",
			"hostId":      host.ID,
			"hostName":    host.HostName,
			"hostIP":      host.SSHIP,
			"serviceName": serviceName,
			"actionType":  actionType,
		},
	}, nil
}

func (s *ServiceManageSkill) getHostInfo(target string) (*cmdbModel.CmdbHost, error) {
	hostByIP, err := s.hostDao.GetCmdbHostByIP(target)
	if err == nil {
		return &hostByIP, nil
	}
	hostByName, err := s.hostDao.GetCmdbHostByName(target)
	if err == nil {
		return &hostByName, nil
	}
	if strings.Contains(err.Error(), "record not found") {
		return nil, fmt.Errorf("该机器 %s 没有配置密钥无法访问，请增加密钥后重新提问", target)
	}
	return nil, fmt.Errorf("查找主机 %s 失败: %v", target, err)
}

func (s *ServiceManageSkill) getSSHConfig(host *cmdbModel.CmdbHost) (*util.SSHConfig, error) {
	ecsAuth, err := s.ecsAuthDao.GetEcsAuthById(host.SSHKeyID)
	if err != nil {
		return nil, err
	}
	sshConfig := &util.SSHConfig{
		IP:       host.SSHIP,
		Port:     host.SSHPort,
		Username: host.SSHName,
		Type:     ecsAuth.Type,
	}
	switch ecsAuth.Type {
	case 1:
		sshConfig.Password = ecsAuth.Password
	case 2:
		sshConfig.PublicKey = ecsAuth.PublicKey
	default:
		return nil, fmt.Errorf("不支持的认证类型: %d", ecsAuth.Type)
	}
	return sshConfig, nil
}

func buildServiceStatusCommand(serviceName string) string {
	return fmt.Sprintf(`SERVICE_NAME=%s
echo "service=%s"
if [ -f /etc/os-release ]; then
  . /etc/os-release
  echo "os_id=${ID:-unknown}"
  echo "os_version=${VERSION_ID:-unknown}"
fi
if systemctl list-unit-files 2>/dev/null | grep -q "^${SERVICE_NAME}\.service"; then
  echo "service_exists=installed"
else
  echo "service_exists=missing"
fi
if systemctl is-active "$SERVICE_NAME" >/dev/null 2>&1; then
  echo "service_status=active"
else
  echo "service_status=inactive"
fi
echo "__SECTION_STATUS_START__"
systemctl status "$SERVICE_NAME" --no-pager -l 2>/dev/null | head -20
echo "__SECTION_STATUS_END__"
echo "__SECTION_JOURNAL_START__"
journalctl -u "$SERVICE_NAME" -n 50 --no-pager 2>/dev/null
echo "__SECTION_JOURNAL_END__"`, shellEscape(serviceName), serviceName)
}

func buildServiceRestartCommand(serviceName string) string {
	return fmt.Sprintf(`SERVICE_NAME=%s
echo "service=%s"
if [ -f /etc/os-release ]; then
  . /etc/os-release
  echo "os_id=${ID:-unknown}"
  echo "os_version=${VERSION_ID:-unknown}"
fi
if systemctl list-unit-files 2>/dev/null | grep -q "^${SERVICE_NAME}\.service"; then
  echo "service_exists=installed"
else
  echo "service_exists=missing"
fi
RESTART_RESULT=skipped
if systemctl list-unit-files 2>/dev/null | grep -q "^${SERVICE_NAME}\.service"; then
  if systemctl restart "$SERVICE_NAME" >/tmp/dodevops-service.log 2>&1; then
    RESTART_RESULT=restarted
  else
    RESTART_RESULT=restart_failed
  fi
else
  RESTART_RESULT=not_installed
fi
echo "restart_result=$RESTART_RESULT"
if systemctl is-active "$SERVICE_NAME" >/dev/null 2>&1; then
  echo "service_status=active"
else
  echo "service_status=inactive"
fi
echo "__SECTION_STATUS_START__"
systemctl status "$SERVICE_NAME" --no-pager -l 2>/dev/null | head -20
echo "__SECTION_STATUS_END__"
echo "__SECTION_JOURNAL_START__"
journalctl -u "$SERVICE_NAME" -n 50 --no-pager 2>/dev/null
echo "__SECTION_JOURNAL_END__"`, shellEscape(serviceName), serviceName)
}

func formatServiceStatusResult(host *cmdbModel.CmdbHost, serviceName, raw string) string {
	values := parseKeyValueLines(raw)
	statusOutput := extractSection(raw, "STATUS")
	journalOutput := extractSection(raw, "JOURNAL")
	analysis := analyzeServiceStartFailure(serviceName, statusOutput, journalOutput)
	var report strings.Builder
	report.WriteString("# 服务检查结果\n\n")
	report.WriteString(fmt.Sprintf("- **主机**: `%s (%s)`\n", host.HostName, host.SSHIP))
	report.WriteString(fmt.Sprintf("- **服务**: `%s`\n", serviceName))
	report.WriteString(fmt.Sprintf("- **系统**: `%s %s`\n", firstNonEmpty(values["os_id"], "unknown"), firstNonEmpty(values["os_version"], "unknown")))
	report.WriteString(fmt.Sprintf("- **服务是否存在**: `%s`\n", firstNonEmpty(values["service_exists"], "unknown")))
	report.WriteString(fmt.Sprintf("- **服务状态**: %s `%s`\n\n", renderStatusText(values["service_status"], true), firstNonEmpty(values["service_status"], "unknown")))

	if values["service_status"] == "active" {
		report.WriteString("## 检查结论\n\n")
		report.WriteString(fmt.Sprintf("%s 服务 `%s` 正常运行。\n", renderStatusText("success", false), serviceName))
		return report.String()
	}

	report.WriteString("## 检查结论\n\n")
	report.WriteString(fmt.Sprintf("%s 服务 `%s` 异常，未正常运行。\n", renderStatusText("failed", false), serviceName))
	if analysis.Reason != "" {
		report.WriteString("\n## 失败原因\n\n")
		report.WriteString(analysis.Reason)
		report.WriteString("\n")
	}
	if len(analysis.Suggestions) > 0 {
		report.WriteString("\n## 建议处理\n\n")
		for _, suggestion := range analysis.Suggestions {
			report.WriteString("- ")
			report.WriteString(suggestion)
			report.WriteString("\n")
		}
	}
	return report.String()
}

func formatServiceRestartResult(host *cmdbModel.CmdbHost, serviceName, raw string) string {
	values := parseKeyValueLines(raw)
	statusOutput := extractSection(raw, "STATUS")
	journalOutput := extractSection(raw, "JOURNAL")
	analysis := analyzeServiceStartFailure(serviceName, statusOutput, journalOutput)
	var report strings.Builder
	report.WriteString("# 服务重启结果\n\n")
	report.WriteString(fmt.Sprintf("- **主机**: `%s (%s)`\n", host.HostName, host.SSHIP))
	report.WriteString(fmt.Sprintf("- **服务**: `%s`\n", serviceName))
	report.WriteString(fmt.Sprintf("- **系统**: `%s %s`\n", firstNonEmpty(values["os_id"], "unknown"), firstNonEmpty(values["os_version"], "unknown")))
	report.WriteString(fmt.Sprintf("- **现有状态**: `%s`\n", firstNonEmpty(values["service_exists"], "unknown")))
	report.WriteString(fmt.Sprintf("- **重启结果**: %s `%s`\n", renderStatusText(values["restart_result"], false), firstNonEmpty(values["restart_result"], "unknown")))
	report.WriteString(fmt.Sprintf("- **服务状态**: %s `%s`\n\n", renderStatusText(values["service_status"], true), firstNonEmpty(values["service_status"], "unknown")))
	report.WriteString("## 检查结论\n\n")
	if values["service_status"] == "active" {
		report.WriteString(fmt.Sprintf("%s 服务 `%s` 已成功重启。\n", renderStatusText("success", false), serviceName))
	} else {
		report.WriteString(fmt.Sprintf("%s 服务 `%s` 重启失败。\n", renderStatusText("failed", false), serviceName))
		if analysis.Reason != "" {
			report.WriteString("\n### 失败原因\n\n- ")
			report.WriteString(analysis.Reason)
			report.WriteString("\n")
		}
		if len(analysis.Suggestions) > 0 {
			report.WriteString("\n### 建议处理\n\n")
			for _, suggestion := range analysis.Suggestions {
				report.WriteString("- ")
				report.WriteString(suggestion)
				report.WriteString("\n")
			}
		}
	}
	return report.String()
}
