package service

import (
	"strings"
	"testing"
	"time"

	agentModel "dodevops-api/api/agent/model"
)

func TestClassifyFileDeletionPolicy(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{path: "/var/log/app.log", want: "可删除"},
		{path: "/tmp/debug.txt", want: "可删除"},
		{path: "/backup/archive.gz", want: "可删除"},
		{path: "/backup/archive.zip", want: "可删除"},
		{path: "/data/mysql/ibdata1", want: "不可自动删除，需用户手动处理"},
		{path: "/app/bin/server", want: "不可自动删除，需用户手动处理"},
	}

	for _, tc := range cases {
		got := classifyFileDeletionPolicy(tc.path)
		if got != tc.want {
			t.Fatalf("path=%s want=%s got=%s", tc.path, tc.want, got)
		}
	}
}

func TestGenerateDirectoryDeepAnalysisReportKeepsOnlyEssentialSections(t *testing.T) {
	svc := &PerformanceAnalysisService{}
	result := &agentModel.DiskAnalysisResult{
		HostName:     "阿里云-ops",
		HostIP:       "8.130.14.34",
		AnalyzedPath: "/var",
		CollectTime:  time.Date(2026, 3, 24, 17, 55, 21, 0, time.Local),
		DiskUsage:    "/dev/vda3 40G 23G 16G 60% /",
		Directories: []agentModel.DirectoryInfo{
			{Path: "/var/lib", Size: "7.1G"},
			{Path: "/var/log", Size: "5.0G"},
		},
		TopFiles: []agentModel.DirectoryInfo{
			{Path: "/var/log/dodevops-agent.log", Size: "517M"},
		},
		CleanupHints: []string{"优先排查 /var/log、/var/lib/docker"},
	}

	report := svc.GenerateDirectoryDeepAnalysisReport(result)

	required := []string{
		"**分析目录**: /var",
		"## 子项占用 TOP 5",
		"## 最大文件 TOP 5",
		"## 可清理建议",
		"## 删除规则",
	}
	for _, item := range required {
		if !strings.Contains(report, item) {
			t.Fatalf("expected report to contain %q, got:\n%s", item, report)
		}
	}

	unwanted := []string{
		"**主机信息**",
		"## 当前磁盘使用情况",
		"## 这个目录里通常是什么",
		"## 建议动作",
	}
	for _, item := range unwanted {
		if strings.Contains(report, item) {
			t.Fatalf("expected report to exclude %q, got:\n%s", item, report)
		}
	}
}
