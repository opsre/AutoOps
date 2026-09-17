package service

import (
	"strings"
	"testing"

	"dodevops-api/api/agent/model"
)

func TestExtractEntitiesDirectoryDeepAnalysis(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	entities := svc.extractEntities("分析下 /var 里面都是什么文件，哪些可以删除")

	if entities["directoryPath"] != "/var" {
		t.Fatalf("expected directoryPath=/var, got %q", entities["directoryPath"])
	}
	if entities["checkType"] != "disk" {
		t.Fatalf("expected checkType=disk, got %q", entities["checkType"])
	}
	if entities["queryType"] != "directory_deep_analysis" {
		t.Fatalf("expected queryType=directory_deep_analysis, got %q", entities["queryType"])
	}
}

func TestExtractEntitiesLargestDirectoryQuery(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	entities := svc.extractEntities("180.76.231.65 检查下这台机器磁盘使用情况，那个目录占用空间最大")

	if entities["targetName"] != "180.76.231.65" {
		t.Fatalf("expected targetName=180.76.231.65, got %q", entities["targetName"])
	}
	if entities["checkType"] != "disk" {
		t.Fatalf("expected checkType=disk, got %q", entities["checkType"])
	}
	if entities["queryType"] != "directories" {
		t.Fatalf("expected queryType=directories, got %q", entities["queryType"])
	}
}

func TestShouldUseDirectoryDeepAnalysisFollowUp(t *testing.T) {
	if !shouldUseDirectoryDeepAnalysis("继续分析这个目录，哪些文件可以删除", "") {
		t.Fatal("expected follow-up directory deep analysis to be detected")
	}
}

func TestExtractEntitiesDirectoryDeepAnalysisIgnoresSizeAsTarget(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	entities := svc.extractEntities("分析下 最大占用目录**: /var (15G) 这个目录里面哪些文件可以删除，哪些文件占用的空间比较大")

	if entities["directoryPath"] != "/var" {
		t.Fatalf("expected directoryPath=/var, got %q", entities["directoryPath"])
	}
	if entities["targetName"] != "" {
		t.Fatalf("expected empty targetName, got %q", entities["targetName"])
	}
	if entities["queryType"] != "directory_deep_analysis" {
		t.Fatalf("expected queryType=directory_deep_analysis, got %q", entities["queryType"])
	}
}

func TestExtractEntitiesDirectoryPathWithoutWhitespace(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	entities := svc.extractEntities("请继续检查下/root/下面有哪些文件")

	if entities["directoryPath"] != "/root/" && entities["directoryPath"] != "/root" {
		t.Fatalf("expected directoryPath=/root or /root/, got %q", entities["directoryPath"])
	}
	if entities["checkType"] != "disk" {
		t.Fatalf("expected checkType=disk, got %q", entities["checkType"])
	}
}

func TestExtractEntitiesListDirectoryFilesQuery(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	entities := svc.extractEntities("请继续检查下 /root 目录下有哪些文件")

	if entities["queryType"] != "list_files" {
		t.Fatalf("expected queryType=list_files, got %q", entities["queryType"])
	}
}

func TestExtractEntitiesInstallService(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	entities := svc.extractEntities("180.76.231.65 请帮我安装一个nginx服务")

	if entities["actionType"] != "install_service" {
		t.Fatalf("expected actionType=install_service, got %q", entities["actionType"])
	}
	if entities["serviceName"] != "nginx" {
		t.Fatalf("expected serviceName=nginx, got %q", entities["serviceName"])
	}
	if entities["checkType"] != "service" {
		t.Fatalf("expected checkType=service, got %q", entities["checkType"])
	}
}

func TestExtractEntitiesInstallServiceWithTypoAlias(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	entities := svc.extractEntities("180.76.231.65 请在这个服务器上安装一个ningx")

	if entities["actionType"] != "install_service" {
		t.Fatalf("expected actionType=install_service, got %q", entities["actionType"])
	}
	if entities["serviceName"] != "nginx" {
		t.Fatalf("expected serviceName=nginx, got %q", entities["serviceName"])
	}
}

func TestExtractEntitiesRestartService(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	entities := svc.extractEntities("请重启 nginx这个服务")

	if entities["actionType"] != "restart_service" {
		t.Fatalf("expected actionType=restart_service, got %q", entities["actionType"])
	}
	if entities["serviceName"] != "nginx" {
		t.Fatalf("expected serviceName=nginx, got %q", entities["serviceName"])
	}
}

func TestExtractEntitiesClearLogFileCommand(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	entities := svc.extractEntities("/var/log/dodevops-agent.log 请清空这个日志文件")

	if entities["actionType"] != "execute_command" {
		t.Fatalf("expected actionType=execute_command, got %q", entities["actionType"])
	}
	if entities["filePath"] != "/var/log/dodevops-agent.log" {
		t.Fatalf("expected filePath=/var/log/dodevops-agent.log, got %q", entities["filePath"])
	}
	if entities["fileOp"] != "clear" {
		t.Fatalf("expected fileOp=clear, got %q", entities["fileOp"])
	}
	if entities["queryType"] != "" {
		t.Fatalf("expected empty queryType, got %q", entities["queryType"])
	}
}

func TestExtractEntitiesDeleteLogFileCommand(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	entities := svc.extractEntities("/var/log/dodevops-agent.log 请删除这个日志文件")

	if entities["actionType"] != "execute_command" {
		t.Fatalf("expected actionType=execute_command, got %q", entities["actionType"])
	}
	if entities["filePath"] != "/var/log/dodevops-agent.log" {
		t.Fatalf("expected filePath=/var/log/dodevops-agent.log, got %q", entities["filePath"])
	}
	if entities["fileOp"] != "delete" {
		t.Fatalf("expected fileOp=delete, got %q", entities["fileOp"])
	}
}

func TestExtractEntitiesRecognizeDeleteNonAllowedFile(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	entities := svc.extractEntities("/etc/passwd 请删除这个文件")

	if entities["fileOp"] != "delete" {
		t.Fatalf("expected fileOp=delete, got %q", entities["fileOp"])
	}
	if entities["filePath"] != "/etc/passwd" {
		t.Fatalf("expected filePath=/etc/passwd, got %q", entities["filePath"])
	}
	if entities["queryType"] != "" {
		t.Fatalf("expected empty queryType, got %q", entities["queryType"])
	}
}

func TestFallbackIntentRecognitionForClearLogFile(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	intent := svc.fallbackIntentRecognition("/var/log/dodevops-agent.log 请帮我 清空这个日志")

	if intent.Type != "incident" {
		t.Fatalf("expected incident intent, got %q", intent.Type)
	}
	if intent.Action != "execute" {
		t.Fatalf("expected execute action, got %q", intent.Action)
	}
	if intent.Entities["fileOp"] != "clear" {
		t.Fatalf("expected fileOp=clear, got %q", intent.Entities["fileOp"])
	}
	if intent.Entities["queryType"] != "" {
		t.Fatalf("expected empty queryType, got %q", intent.Entities["queryType"])
	}
}

func TestFallbackIntentRecognitionForDirectoryInspectionFollowUp(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	intent := svc.fallbackIntentRecognition("请继续检查 /root/ 目录下有哪些文件")

	if intent.Type != "incident" {
		t.Fatalf("expected incident intent, got %q", intent.Type)
	}
	if intent.Action != "analyze" {
		t.Fatalf("expected analyze action, got %q", intent.Action)
	}
	if intent.Entities["directoryPath"] != "/root/" && intent.Entities["directoryPath"] != "/root" {
		t.Fatalf("expected directoryPath to be /root or /root/, got %q", intent.Entities["directoryPath"])
	}
	if intent.Entities["checkType"] != "disk" {
		t.Fatalf("expected checkType=disk, got %q", intent.Entities["checkType"])
	}
}

func TestFallbackIntentRecognitionForDeleteAnyFilePath(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	intent := svc.fallbackIntentRecognition("请继续删除 /root/a.sh")

	if intent.Type != "incident" {
		t.Fatalf("expected incident intent, got %q", intent.Type)
	}
	if intent.Action != "execute" {
		t.Fatalf("expected execute action, got %q", intent.Action)
	}
	if intent.Entities["fileOp"] != "delete" {
		t.Fatalf("expected fileOp=delete, got %q", intent.Entities["fileOp"])
	}
	if intent.Entities["filePath"] != "/root/a.sh" {
		t.Fatalf("expected filePath=/root/a.sh, got %q", intent.Entities["filePath"])
	}
	if intent.Entities["queryType"] != "" {
		t.Fatalf("expected empty queryType, got %q", intent.Entities["queryType"])
	}
}

func TestFallbackIntentRecognitionForInstallService(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	intent := svc.fallbackIntentRecognition("180.76.231.65 请帮我安装一个nginx服务")

	if intent.Type != "incident" {
		t.Fatalf("expected incident intent, got %q", intent.Type)
	}
	if intent.Action != "execute" {
		t.Fatalf("expected execute action, got %q", intent.Action)
	}
	if intent.Entities["actionType"] != "install_service" {
		t.Fatalf("expected actionType=install_service, got %q", intent.Entities["actionType"])
	}
	if intent.Entities["serviceName"] != "nginx" {
		t.Fatalf("expected serviceName=nginx, got %q", intent.Entities["serviceName"])
	}
}

func TestFallbackIntentRecognitionForInstallServiceTypoAlias(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	intent := svc.fallbackIntentRecognition("180.76.231.65 请在这个服务器上安装一个ningx")

	if intent.Type != "incident" {
		t.Fatalf("expected incident intent, got %q", intent.Type)
	}
	if intent.Entities["actionType"] != "install_service" {
		t.Fatalf("expected actionType=install_service, got %q", intent.Entities["actionType"])
	}
	if intent.Entities["serviceName"] != "nginx" {
		t.Fatalf("expected serviceName=nginx, got %q", intent.Entities["serviceName"])
	}
}

func TestFallbackIntentRecognitionForRestartService(t *testing.T) {
	svc := &AgentChatServiceImpl{}

	intent := svc.fallbackIntentRecognition("请重启 nginx这个服务")

	if intent.Type != "incident" {
		t.Fatalf("expected incident intent, got %q", intent.Type)
	}
	if intent.Action != "execute" {
		t.Fatalf("expected execute action, got %q", intent.Action)
	}
	if intent.Entities["actionType"] != "restart_service" {
		t.Fatalf("expected actionType=restart_service, got %q", intent.Entities["actionType"])
	}
}

func TestManagedFileFollowUpDirectory(t *testing.T) {
	if got := managedFileFollowUpDirectory("/var/log/dodevops-agent.log"); got != "/var/log" {
		t.Fatalf("expected /var/log, got %q", got)
	}
	if got := managedFileFollowUpDirectory(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
	if got := managedFileFollowUpDirectory("relative.log"); got != "" {
		t.Fatalf("expected empty string for relative path, got %q", got)
	}
}

func TestMergeResolvedContextIntoIntentSkipsQueryTypeForManagedFileOp(t *testing.T) {
	entities := map[string]string{
		"actionType": "execute_command",
		"fileOp":     "clear",
		"filePath":   "/var/log/dodevops-agent.log",
		"checkType":  "disk",
	}

	req := &model.AgentChatRequest{
		ResolvedTargetType:    "host",
		ResolvedTargetName:    "180.76.231.65",
		ResolvedCheckType:     "disk",
		ResolvedQueryType:     "directory_deep_analysis",
		ResolvedDirectoryPath: "/var",
	}

	merged := mergeResolvedContextIntoIntent(entities, req)

	if merged["targetName"] != "180.76.231.65" {
		t.Fatalf("expected targetName to be backfilled, got %q", merged["targetName"])
	}
	if merged["queryType"] != "" {
		t.Fatalf("expected queryType to stay empty, got %q", merged["queryType"])
	}
	if merged["directoryPath"] != "" {
		t.Fatalf("expected directoryPath to stay empty, got %q", merged["directoryPath"])
	}
}

func TestFormatManagedFileStatus(t *testing.T) {
	response := &model.AgentChatResponse{
		Message: "file=/var/log/app.log\naction=delete\nstatus=deleted\n",
	}

	if got := formatManagedFileStatus(response, "delete"); got != "删除成功" {
		t.Fatalf("expected 删除成功, got %q", got)
	}

	response.Message = "file=/var/log/app.log\naction=clear\nstatus=cleared\n"
	if got := formatManagedFileStatus(response, "clear"); got != "清空成功" {
		t.Fatalf("expected 清空成功, got %q", got)
	}
}

func TestBuildTopFilesSection(t *testing.T) {
	section := buildTopFilesSection([]model.DirectoryInfo{
		{Path: "/var/log/app.log", Size: "1M"},
	})

	if !strings.Contains(section, "## 最大文件 TOP 5") {
		t.Fatalf("expected top files heading, got %q", section)
	}
	if !strings.Contains(section, "/var/log/app.log") {
		t.Fatalf("expected file path in section, got %q", section)
	}
}
