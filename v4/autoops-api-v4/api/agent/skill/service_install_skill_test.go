package skill

import (
	"strings"
	"testing"

	cmdbModel "dodevops-api/api/cmdb/model"
	baseSkill "dodevops-api/pkg/agent/skill"
)

func TestServiceInstallSkillCanHandle(t *testing.T) {
	skill := &ServiceManageSkill{}
	intent := &baseSkill.Intent{
		Type: "incident",
		Entities: map[string]string{
			"actionType":  "install_service",
			"checkType":   "service",
			"serviceName": "nginx",
		},
	}

	if !skill.CanHandle(intent) {
		t.Fatal("expected service manage skill to handle install_service intent")
	}
}

func TestServiceManageSkillCanHandleRestartServiceIntent(t *testing.T) {
	skill := &ServiceManageSkill{}
	intent := &baseSkill.Intent{
		Type: "incident",
		Entities: map[string]string{
			"checkType":   "service",
			"actionType":  "restart_service",
			"serviceName": "nginx",
		},
	}

	if !skill.CanHandle(intent) {
		t.Fatal("expected service manage skill to handle restart_service intent")
	}
}

func TestBuildServiceInstallCommandContainsCoreSteps(t *testing.T) {
	command := buildServiceInstallCommand("nginx")
	required := []string{
		"service=nginx",
		"package_manager=",
		"install_result=",
		"start_result=",
		"service_status=",
		"systemctl enable --now",
		"journalctl -u \"$SERVICE_NAME\" -n 50",
		"__SECTION_STATUS_START__",
		"__SECTION_JOURNAL_START__",
	}
	for _, item := range required {
		if !strings.Contains(command, item) {
			t.Fatalf("expected command to contain %q, got %s", item, command)
		}
	}
}

func TestAnalyzeServiceStartFailurePortOccupied(t *testing.T) {
	analysis := analyzeServiceStartFailure(
		"nginx",
		"bind() to 0.0.0.0:80 failed (98: Address already in use)",
		"",
	)

	if !strings.Contains(analysis.Reason, "80") {
		t.Fatalf("expected occupied port in reason, got %q", analysis.Reason)
	}
	if len(analysis.Suggestions) == 0 || !strings.Contains(analysis.Suggestions[0], "ss -lntp | grep :80") {
		t.Fatalf("expected port suggestion, got %+v", analysis.Suggestions)
	}
}

func TestExtractConflictingPort(t *testing.T) {
	raw := "Mar 24 19:55:26 nginx[3218265]: nginx: [emerg] bind() to 0.0.0.0:80 failed (98: Unknown error)"
	if got := extractConflictingPort(raw); got != "80" {
		t.Fatalf("expected port 80, got %q", got)
	}
}

func TestAnalyzeServiceStartFailureConfigError(t *testing.T) {
	analysis := analyzeServiceStartFailure(
		"nginx",
		"nginx: [emerg] unexpected end of file, expecting \";\"",
		"",
	)

	if analysis.Reason != "配置文件语法错误" {
		t.Fatalf("expected config syntax error, got %q", analysis.Reason)
	}
}

func TestFormatServiceInstallResultIncludesFailureAnalysis(t *testing.T) {
	raw := strings.Join([]string{
		"service=nginx",
		"os_id=ubuntu",
		"os_version=22.04",
		"service_exists=missing",
		"package_manager=apt",
		"install_result=installed",
		"start_result=start_failed",
		"service_status=inactive",
		"__SECTION_STATUS_START__",
		"bind() to 0.0.0.0:80 failed (98: Address already in use)",
		"__SECTION_STATUS_END__",
		"__SECTION_JOURNAL_START__",
		"Address already in use",
		"__SECTION_JOURNAL_END__",
	}, "\n")

	host := &cmdbModel.CmdbHost{HostName: "test-host", SSHIP: "127.0.0.1"}
	report := formatServiceInstallResult(host, "nginx", raw)

	if !strings.Contains(report, "### 失败原因") {
		t.Fatalf("expected failure reason section, got %s", report)
	}
	if !strings.Contains(report, "ss -lntp | grep :80") {
		t.Fatalf("expected port check suggestion, got %s", report)
	}
	if strings.Contains(report, "## 最近日志") || strings.Contains(report, "## 服务状态") {
		t.Fatalf("expected compact report without raw logs, got %s", report)
	}
}

func TestFormatServiceInstallResultAlreadyInstalledAndRunning(t *testing.T) {
	raw := strings.Join([]string{
		"service=nginx",
		"os_id=ubuntu",
		"os_version=22.04",
		"service_exists=installed",
		"package_manager=apt",
		"current_status=active",
		"install_result=skipped",
		"start_result=already_running",
		"service_status=active",
	}, "\n")

	host := &cmdbModel.CmdbHost{HostName: "test-host", SSHIP: "127.0.0.1"}
	report := formatServiceInstallResult(host, "nginx", raw)

	if !strings.Contains(report, "已安装，跳过重复安装") {
		t.Fatalf("expected duplicate install skip message, got %s", report)
	}
	if !strings.Contains(report, "服务已在运行，无需重复启动") {
		t.Fatalf("expected already running message, got %s", report)
	}
	if !strings.Contains(report, "已安装且当前正在运行，无需重复安装") {
		t.Fatalf("expected clear conclusion, got %s", report)
	}
}

func TestFormatServiceStatusResultReturnsCompactFailureAnalysis(t *testing.T) {
	raw := strings.Join([]string{
		"service=nginx",
		"os_id=ubuntu",
		"os_version=22.04",
		"service_exists=installed",
		"service_status=inactive",
		"__SECTION_STATUS_START__",
		"bind() to 0.0.0.0:80 failed (98: Address already in use)",
		"__SECTION_STATUS_END__",
		"__SECTION_JOURNAL_START__",
		"Address already in use",
		"__SECTION_JOURNAL_END__",
	}, "\n")

	host := &cmdbModel.CmdbHost{HostName: "test-host", SSHIP: "127.0.0.1"}
	report := formatServiceStatusResult(host, "nginx", raw)

	if !strings.Contains(report, "检查结论") {
		t.Fatalf("expected conclusion section, got %s", report)
	}
	if !strings.Contains(report, "80") {
		t.Fatalf("expected port analysis, got %s", report)
	}
	if !strings.Contains(report, "ss -lntp | grep :80") {
		t.Fatalf("expected suggestion for occupied port, got %s", report)
	}
	if strings.Contains(report, "systemctl status") || strings.Contains(report, "```plaintext") {
		t.Fatalf("expected compact report without raw status output, got %s", report)
	}
}

func TestFormatServiceStatusResultReturnsHealthyConclusion(t *testing.T) {
	raw := strings.Join([]string{
		"service=nginx",
		"os_id=ubuntu",
		"os_version=22.04",
		"service_exists=installed",
		"service_status=active",
	}, "\n")

	host := &cmdbModel.CmdbHost{HostName: "test-host", SSHIP: "127.0.0.1"}
	report := formatServiceStatusResult(host, "nginx", raw)

	if !strings.Contains(report, "服务 `nginx` 正常运行") {
		t.Fatalf("expected healthy conclusion, got %s", report)
	}
	if strings.Contains(report, "失败原因") || strings.Contains(report, "建议处理") {
		t.Fatalf("expected no failure analysis for healthy service, got %s", report)
	}
}
