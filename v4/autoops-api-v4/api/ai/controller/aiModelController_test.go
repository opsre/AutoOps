package controller

import (
	"testing"

	agentModel "dodevops-api/api/agent/model"
	"dodevops-api/api/ai/model"
	aiService "dodevops-api/api/ai/service"
)

func TestShouldTreatAsFollowUp(t *testing.T) {
	ctx := aiService.SessionStructuredContext{
		LastTargetType: "host",
		LastTargetName: "180.76.231.65",
	}

	if !shouldTreatAsFollowUp("分析下 最大占用目录**: /var (15G) 这个目录里面哪些文件可以删除，哪些文件占用的空间比较大", ctx) {
		t.Fatal("expected message to be treated as follow-up")
	}
}

func TestIsAgentIntentForClearLogFile(t *testing.T) {
	req := model.AiChatRequest{
		Messages: []model.AiChatMessage{
			{Role: "user", Content: "180.76.231.65 检查下这台机器磁盘使用情况，那个目录占用空间最大"},
			{Role: "assistant", Content: "最大占用目录: /var"},
			{Role: "user", Content: "/var/log/dodevops-agent.log 请清空这个日志文件"},
		},
	}

	if !isAgentIntent(req) {
		t.Fatal("expected clear log file request to be treated as agent intent")
	}
}

func TestHasHostNamePatternIgnoresDirectoryAndSizeNoise(t *testing.T) {
	if hasHostNamePattern("**最大占用目录**: /var (15G) 分析下这个目录下都是什么文件，哪些可以删除") {
		t.Fatal("expected directory path and size token not to be treated as host pattern")
	}
}

func TestHistoryQueryIntentShouldIgnoreFileFollowUp(t *testing.T) {
	req := model.AiChatRequest{
		Messages: []model.AiChatMessage{
			{Role: "user", Content: "180.76.231.65 检查这台服务器的磁盘，找到那个目录占用空间最大，并分析给出优化建议"},
			{Role: "assistant", Content: "最大占用目录: /var"},
			{Role: "user", Content: "这个日志文件清空后多大"},
		},
	}

	if isHistoryQueryIntent(req) {
		t.Fatal("expected file follow-up not to be treated as history query")
	}
}

func TestIsAgentIntentForDirectoryFollowUpQuestion(t *testing.T) {
	req := model.AiChatRequest{
		Messages: []model.AiChatMessage{
			{Role: "user", Content: "180.76.231.65 检查这台服务器的磁盘，找到那个目录占用空间最大，并分析给出优化建议"},
			{Role: "assistant", Content: "最大占用目录: /var"},
			{Role: "user", Content: "这个目录还能删什么"},
		},
	}

	if !isAgentIntent(req) {
		t.Fatal("expected directory follow-up to be treated as agent intent")
	}
}

func TestShouldBackfillContextFromRecentHistory(t *testing.T) {
	if !shouldBackfillContextFromRecentHistory("检查 /var 这个目录下都有哪些文件，哪些文件可以删除") {
		t.Fatal("expected directory follow-up to trigger recent-history backfill")
	}
}

func TestConvertAgentHistoryToAiMessages(t *testing.T) {
	histories := []agentModel.AgentChatHistory{
		{Role: "user", Message: "first"},
		{Role: "assistant", Message: "second"},
	}

	msgs := convertAgentHistoriesToAiMessages(histories)
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Role != "user" || msgs[0].Content != "first" {
		t.Fatalf("unexpected first message: %+v", msgs[0])
	}
}

func TestIsAgentIntentForManagedFileOperationWithoutExplicitTarget(t *testing.T) {
	req := model.AiChatRequest{
		Messages: []model.AiChatMessage{
			{Role: "user", Content: "请清空 /var/log/dodevops-agent.log 这个日志目录"},
		},
	}

	if !isAgentIntent(req) {
		t.Fatal("expected managed log operation to be treated as agent intent")
	}
}

func TestIsAgentIntentForAbsoluteDirectoryFollowUpWithoutExplicitTarget(t *testing.T) {
	req := model.AiChatRequest{
		Messages: []model.AiChatMessage{
			{Role: "user", Content: "180.76.231.65 检查下这台机器磁盘使用情况，那个目录占用空间最大"},
			{Role: "assistant", Content: "最大占用目录: /var"},
			{Role: "user", Content: "请继续检查下/root/ 有哪些文件"},
		},
	}

	if !isAgentIntent(req) {
		t.Fatal("expected absolute directory follow-up to be treated as agent intent")
	}
}

func TestIsAgentIntentForContinuationCommandUsesSessionTarget(t *testing.T) {
	sessionID := "test-session-continuation-command"
	aiService.SaveSessionStructuredContext(sessionID, aiService.SessionStructuredContext{
		LastTargetType: "host",
		LastTargetName: "180.76.231.65",
		LastCheckType:  "disk",
	})

	req := model.AiChatRequest{
		SessionID: sessionID,
		Messages: []model.AiChatMessage{
			{Role: "user", Content: "继续执行命令 df -h"},
		},
	}

	if !isAgentIntent(req) {
		t.Fatal("expected continuation command to reuse session target and be treated as agent intent")
	}
}

func TestExtractDirectoryPathWithoutWhitespace(t *testing.T) {
	got := extractDirectoryPath("请继续检查下/root/下面有哪些文件")
	if got != "/root/" && got != "/root" {
		t.Fatalf("expected /root or /root/, got %q", got)
	}
}

func TestIsAgentIntentForInstallService(t *testing.T) {
	req := model.AiChatRequest{
		Messages: []model.AiChatMessage{
			{Role: "user", Content: "180.76.231.65 请帮我安装一个nginx服务"},
		},
	}

	if !isAgentIntent(req) {
		t.Fatal("expected install service request to be treated as agent intent")
	}
}

func TestIsAgentIntentForRestartService(t *testing.T) {
	sessionID := "test-session-restart-service"
	aiService.SaveSessionStructuredContext(sessionID, aiService.SessionStructuredContext{
		LastTargetType: "host",
		LastTargetName: "180.76.231.65",
		LastCheckType:  "service",
	})

	req := model.AiChatRequest{
		SessionID: sessionID,
		Messages: []model.AiChatMessage{
			{Role: "user", Content: "请 重启 nginx这个服务"},
		},
	}

	if !isAgentIntent(req) {
		t.Fatal("expected restart service request to be treated as agent intent")
	}
}

func TestShouldTreatAsFollowUpForServiceOperation(t *testing.T) {
	ctx := aiService.SessionStructuredContext{
		LastTargetType: "host",
		LastTargetName: "180.76.231.65",
		LastCheckType:  "service",
	}

	if !shouldTreatAsFollowUp("请检查 nginx 这个服务状态", ctx) {
		t.Fatal("expected service operation to be treated as follow-up")
	}
}

func TestIsAgentIntentForInstallServiceTypoAlias(t *testing.T) {
	req := model.AiChatRequest{
		Messages: []model.AiChatMessage{
			{Role: "user", Content: "180.76.231.65 请在这个服务器上安装一个ningx"},
		},
	}

	if !isAgentIntent(req) {
		t.Fatal("expected install service request with typo alias to be treated as agent intent")
	}
}
