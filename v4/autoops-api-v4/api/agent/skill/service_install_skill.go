package skill

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	cmdbDao "dodevops-api/api/cmdb/dao"
	cmdbModel "dodevops-api/api/cmdb/model"
	monitorService "dodevops-api/api/monitor/service"
	"dodevops-api/common/util"
	baseSkill "dodevops-api/pkg/agent/skill"
)

// ServiceInstallSkill 服务安装技能
type ServiceInstallSkill struct {
	baseSkill.BaseSkill
	hostDao    cmdbDao.CmdbHostDao
	ecsAuthDao monitorService.EcsAuthDao
	sshUtil    *util.SSHUtil
}

func NewServiceInstallSkill() *ServiceInstallSkill {
	return &ServiceInstallSkill{
		hostDao:    cmdbDao.NewCmdbHostDao(),
		ecsAuthDao: monitorService.NewEcsAuthDao(),
		sshUtil:    util.NewSSHUtil(),
	}
}

func (s *ServiceInstallSkill) Name() string {
	return "service_install"
}

func (s *ServiceInstallSkill) Description() string {
	return "在指定服务器上安装并启动系统服务"
}

func (s *ServiceInstallSkill) Category() baseSkill.SkillCategory {
	return baseSkill.CategoryOperation
}

func (s *ServiceInstallSkill) CanHandle(intent *baseSkill.Intent) bool {
	return intent.Type == "incident" && intent.Entities["actionType"] == "install_service"
}

func (s *ServiceInstallSkill) Validate(req *baseSkill.SkillRequest) error {
	targetName := req.Intent.Entities["targetName"]
	if targetName == "" {
		return fmt.Errorf("请指定要安装服务的主机IP或名称")
	}
	serviceName := strings.TrimSpace(req.Intent.Entities["serviceName"])
	if serviceName == "" {
		return fmt.Errorf("请指定要安装的服务名称，例如 nginx")
	}
	return nil
}

func (s *ServiceInstallSkill) Execute(ctx context.Context, req *baseSkill.SkillRequest) (*baseSkill.SkillResponse, error) {
	serviceName := strings.TrimSpace(req.Intent.Entities["serviceName"])
	targetName := strings.TrimSpace(req.Intent.Entities["targetName"])

	host, err := s.getHostInfo(targetName)
	if err != nil {
		return nil, err
	}

	sshConfig, err := s.getSSHConfig(host)
	if err != nil {
		return nil, fmt.Errorf("获取SSH认证信息失败: %v", err)
	}

	command := buildServiceInstallCommand(serviceName)
	result, err := s.sshUtil.ExecuteRemoteCommand(sshConfig, command)
	if err != nil {
		return nil, fmt.Errorf("服务安装失败: %v", err)
	}

	report := formatServiceInstallResult(host, serviceName, result)

	return &baseSkill.SkillResponse{
		Message: report,
		Type:    "text",
		TaskInfo: map[string]interface{}{
			"taskType":    "service_install",
			"hostId":      host.ID,
			"hostName":    host.HostName,
			"hostIP":      host.SSHIP,
			"serviceName": serviceName,
			"result":      result,
		},
	}, nil
}

func (s *ServiceInstallSkill) getHostInfo(target string) (*cmdbModel.CmdbHost, error) {
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

func (s *ServiceInstallSkill) getSSHConfig(host *cmdbModel.CmdbHost) (*util.SSHConfig, error) {
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

func buildServiceInstallCommand(serviceName string) string {
	return fmt.Sprintf(`SERVICE_NAME=%s
if [ -f /etc/os-release ]; then
  . /etc/os-release
  OS_ID="${ID:-unknown}"
  OS_VERSION="${VERSION_ID:-unknown}"
else
  OS_ID=unknown
  OS_VERSION=unknown
fi
echo "service=%s"
echo "os_id=$OS_ID"
echo "os_version=$OS_VERSION"
if systemctl list-unit-files 2>/dev/null | grep -q "^${SERVICE_NAME}\.service"; then
  echo "service_exists=installed"
elif command -v "$SERVICE_NAME" >/dev/null 2>&1; then
  echo "service_exists=binary_only"
else
  echo "service_exists=missing"
fi
PKG_MGR=unknown
if command -v apt-get >/dev/null 2>&1; then
  PKG_MGR=apt
elif command -v dnf >/dev/null 2>&1; then
  PKG_MGR=dnf
elif command -v yum >/dev/null 2>&1; then
  PKG_MGR=yum
fi
echo "package_manager=$PKG_MGR"
CURRENT_STATUS=inactive
if systemctl is-active "$SERVICE_NAME" >/dev/null 2>&1; then
  CURRENT_STATUS=active
fi
echo "current_status=$CURRENT_STATUS"
INSTALL_RESULT=skipped
if ! systemctl list-unit-files 2>/dev/null | grep -q "^${SERVICE_NAME}\.service" && ! command -v "$SERVICE_NAME" >/dev/null 2>&1; then
  if [ "$PKG_MGR" = "apt" ]; then
    DEBIAN_FRONTEND=noninteractive apt-get update -y >/tmp/dodevops-install.log 2>&1 && DEBIAN_FRONTEND=noninteractive apt-get install -y "$SERVICE_NAME" >>/tmp/dodevops-install.log 2>&1 && INSTALL_RESULT=installed || INSTALL_RESULT=install_failed
  elif [ "$PKG_MGR" = "dnf" ]; then
    dnf install -y "$SERVICE_NAME" >/tmp/dodevops-install.log 2>&1 && INSTALL_RESULT=installed || INSTALL_RESULT=install_failed
  elif [ "$PKG_MGR" = "yum" ]; then
    yum install -y "$SERVICE_NAME" >/tmp/dodevops-install.log 2>&1 && INSTALL_RESULT=installed || INSTALL_RESULT=install_failed
  else
    INSTALL_RESULT=unsupported_os
  fi
fi
echo "install_result=$INSTALL_RESULT"
START_RESULT=skipped
if [ "$INSTALL_RESULT" != "install_failed" ] && [ "$INSTALL_RESULT" != "unsupported_os" ]; then
  if [ "$CURRENT_STATUS" = "active" ]; then
    START_RESULT=already_running
  else
    if systemctl enable --now "$SERVICE_NAME" >/tmp/dodevops-service.log 2>&1; then
      START_RESULT=started
    elif systemctl start "$SERVICE_NAME" >>/tmp/dodevops-service.log 2>&1; then
      START_RESULT=started
    else
      START_RESULT=start_failed
    fi
  fi
fi
echo "start_result=$START_RESULT"
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

func formatServiceInstallResult(host *cmdbModel.CmdbHost, serviceName, raw string) string {
	values := parseKeyValueLines(raw)
	statusOutput := extractSection(raw, "STATUS")
	journalOutput := extractSection(raw, "JOURNAL")
	analysis := analyzeServiceStartFailure(serviceName, statusOutput, journalOutput)
	serviceExists := firstNonEmpty(values["service_exists"], "unknown")
	installResult := firstNonEmpty(values["install_result"], "unknown")
	startResult := firstNonEmpty(values["start_result"], "unknown")
	serviceStatus := firstNonEmpty(values["service_status"], "unknown")
	var report strings.Builder
	report.WriteString("# 服务安装结果\n\n")
	report.WriteString(fmt.Sprintf("- **主机**: `%s (%s)`\n", host.HostName, host.SSHIP))
	report.WriteString(fmt.Sprintf("- **服务**: `%s`\n", serviceName))
	report.WriteString(fmt.Sprintf("- **系统**: `%s %s`\n", firstNonEmpty(values["os_id"], "unknown"), firstNonEmpty(values["os_version"], "unknown")))
	report.WriteString(fmt.Sprintf("- **现有状态**: `%s`\n", serviceExists))
	report.WriteString(fmt.Sprintf("- **包管理器**: `%s`\n", firstNonEmpty(values["package_manager"], "unknown")))
	report.WriteString(fmt.Sprintf("- **安装结果**: %s\n", renderInstallOutcome(serviceExists, installResult)))
	report.WriteString(fmt.Sprintf("- **启动结果**: %s\n", renderStartOutcome(startResult)))
	report.WriteString(fmt.Sprintf("- **服务状态**: %s `%s`\n\n", renderStatusText(serviceStatus, true), serviceStatus))
	report.WriteString("## 检查结论\n\n")
	if serviceStatus == "active" {
		switch {
		case installResult == "skipped" && startResult == "already_running":
			report.WriteString(fmt.Sprintf("%s 服务 `%s` 已安装且当前正在运行，无需重复安装。\n", renderStatusText("success", false), serviceName))
		case installResult == "skipped":
			report.WriteString(fmt.Sprintf("%s 服务 `%s` 已安装，本次跳过安装并已完成状态检查。\n", renderStatusText("success", false), serviceName))
		default:
			report.WriteString(fmt.Sprintf("%s 服务 `%s` 安装并启动成功。\n", renderStatusText("success", false), serviceName))
		}
	} else {
		report.WriteString(fmt.Sprintf("%s 服务 `%s` 未成功启动。\n", renderStatusText("failed", false), serviceName))
		if analysis.Reason != "" {
			report.WriteString("\n")
			report.WriteString("### 失败原因\n\n")
			report.WriteString("- ")
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

func renderStatusText(status string, activeIsSuccess bool) string {
	normalized := strings.ToLower(strings.TrimSpace(status))
	label := "未知"

	switch normalized {
	case "active":
		if activeIsSuccess {
			label = "成功"
		} else {
			label = "运行中"
		}
	case "inactive", "failed", "start_failed", "restart_failed":
		label = "失败"
	case "installed", "installed_now":
		label = "成功"
	case "missing", "not_installed":
		label = "失败"
	case "skipped":
		label = "跳过"
	case "already_running":
		label = "已在运行"
	case "restarted":
		label = "成功"
	default:
		if strings.Contains(normalized, "success") || strings.Contains(normalized, "active") || strings.Contains(normalized, "installed") {
			label = "成功"
		} else if strings.Contains(normalized, "fail") || strings.Contains(normalized, "inactive") || strings.Contains(normalized, "error") {
			label = "失败"
		} else if normalized != "" {
			label = normalized
		}
	}
	return fmt.Sprintf("**%s**", label)
}

func renderInstallOutcome(serviceExists, installResult string) string {
	switch {
	case installResult == "installed":
		return fmt.Sprintf("%s `已安装`", renderStatusText("installed", false))
	case installResult == "skipped" && (serviceExists == "installed" || serviceExists == "binary_only"):
		return fmt.Sprintf("%s `已安装，跳过重复安装`", renderStatusText("skipped", false))
	case installResult == "unsupported_os":
		return fmt.Sprintf("%s `当前系统暂不支持自动安装`", renderStatusText("failed", false))
	default:
		return fmt.Sprintf("%s `%s`", renderStatusText(installResult, false), firstNonEmpty(installResult, "unknown"))
	}
}

func renderStartOutcome(startResult string) string {
	switch startResult {
	case "already_running":
		return fmt.Sprintf("%s `服务已在运行，无需重复启动`", renderStatusText(startResult, false))
	case "started":
		return fmt.Sprintf("%s `已启动`", renderStatusText(startResult, false))
	default:
		return fmt.Sprintf("%s `%s`", renderStatusText(startResult, false), firstNonEmpty(startResult, "unknown"))
	}
}

func parseKeyValueLines(raw string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "__SECTION_") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			result[key] = val
		}
	}
	return result
}

type serviceFailureAnalysis struct {
	Reason      string
	Suggestions []string
}

func analyzeServiceStartFailure(serviceName, statusOutput, journalOutput string) serviceFailureAnalysis {
	combined := strings.ToLower(statusOutput + "\n" + journalOutput)

	switch {
	case strings.Contains(combined, "address already in use") ||
		strings.Contains(combined, "failed to bind") ||
		strings.Contains(combined, "bind() to") ||
		strings.Contains(combined, "port is already allocated"):
		port := extractConflictingPort(statusOutput + "\n" + journalOutput)
		reason := "端口被占用"
		suggestions := []string{
			"执行 `ss -lntp | grep :" + port + "` 查看占用该端口的进程",
			"确认占用进程后，再决定停掉旧服务或修改 " + serviceName + " 的监听端口",
		}
		if port == "" {
			suggestions = []string{
				"执行 `ss -lntp` 查看端口占用情况",
				"确认占用进程后，再决定停掉旧服务或修改 " + serviceName + " 的监听端口",
			}
		}
		if port != "" {
			reason = fmt.Sprintf("%s 端口 `%s` 被占用", serviceName, port)
		}
		return serviceFailureAnalysis{Reason: reason, Suggestions: suggestions}
	case strings.Contains(combined, "syntax is invalid") ||
		strings.Contains(combined, "syntax error") ||
		strings.Contains(combined, "unexpected") ||
		strings.Contains(combined, "invalid number of arguments"):
		return serviceFailureAnalysis{
			Reason: "配置文件语法错误",
			Suggestions: []string{
				"先执行配置检查命令，例如 `nginx -t`",
				"修正配置文件语法后，再重新启动服务",
			},
		}
	case strings.Contains(combined, "permission denied") ||
		strings.Contains(combined, "operation not permitted"):
		return serviceFailureAnalysis{
			Reason: "权限不足",
			Suggestions: []string{
				"检查服务运行用户对配置、日志、数据目录是否有读写权限",
				"必要时修正目录属主属组与文件权限后重试",
			},
		}
	case strings.Contains(combined, "no such file or directory") ||
		strings.Contains(combined, "cannot load certificate") ||
		strings.Contains(combined, "failed to open") ||
		strings.Contains(combined, "unable to load certificate"):
		return serviceFailureAnalysis{
			Reason: "缺少目录、文件或证书",
			Suggestions: []string{
				"检查配置中引用的目录、日志文件、证书和私钥文件是否存在",
				"确认路径正确且服务运行用户有访问权限",
			},
		}
	case strings.Contains(combined, "selinux") ||
		strings.Contains(combined, "avc denied") ||
		strings.Contains(combined, "firewalld") ||
		strings.Contains(combined, "firewall-cmd"):
		return serviceFailureAnalysis{
			Reason: "SELinux 或防火墙策略限制",
			Suggestions: []string{
				"检查 SELinux 审计日志和防火墙规则是否拦截了服务行为",
				"必要时调整 SELinux 上下文、布尔值或开放对应端口",
			},
		}
	default:
		return serviceFailureAnalysis{
			Reason: "服务启动失败，但未匹配到已知原因",
			Suggestions: []string{
				"优先查看下方 `systemctl status` 和 `journalctl` 输出中的关键报错",
				"根据报错进一步检查端口、配置、权限和依赖服务",
			},
		}
	}
}

func extractSection(raw, name string) string {
	startMarker := "__SECTION_" + name + "_START__"
	endMarker := "__SECTION_" + name + "_END__"
	start := strings.Index(raw, startMarker)
	if start == -1 {
		return ""
	}
	start += len(startMarker)
	end := strings.Index(raw[start:], endMarker)
	if end == -1 {
		return strings.TrimSpace(raw[start:])
	}
	return strings.TrimSpace(raw[start : start+end])
}

func extractConflictingPort(raw string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?:0\.0\.0\.0|\[::\]|localhost):([0-9]{1,5})`),
		regexp.MustCompile(`port[[:space:]]+([0-9]{1,5})`),
	}
	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(raw)
		if len(matches) == 2 && isAllDigits(matches[1]) {
			return matches[1]
		}
	}
	return ""
}

func isAllDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func filterNonKeyValueLines(raw string) string {
	var lines []string
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "__SECTION_") {
			continue
		}
		if strings.Contains(trimmed, "=") {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 && !strings.Contains(parts[0], " ") {
				continue
			}
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func shellEscape(value string) string {
	value = strings.ReplaceAll(value, `'`, `'\''`)
	return "'" + value + "'"
}
