package skill

import (
	"testing"

	baseSkill "dodevops-api/pkg/agent/skill"
)

func TestDiskAnalysisSkillCanHandleRejectsManagedFileOperation(t *testing.T) {
	skill := &DiskAnalysisSkill{}
	intent := &baseSkill.Intent{
		Type: "incident",
		Entities: map[string]string{
			"checkType":  "disk",
			"actionType": "execute_command",
			"fileOp":     "clear",
			"filePath":   "/var/log/dodevops-agent.log",
		},
	}

	if skill.CanHandle(intent) {
		t.Fatal("expected disk analysis skill to reject managed file operations")
	}
}

func TestDirectoryDeepAnalysisSkillCanHandleListFilesQuery(t *testing.T) {
	skill := &DirectoryDeepAnalysisSkill{}
	intent := &baseSkill.Intent{
		Type: "incident",
		Entities: map[string]string{
			"checkType":     "disk",
			"queryType":     "list_files",
			"directoryPath": "/root",
		},
	}

	if !skill.CanHandle(intent) {
		t.Fatal("expected directory deep analysis skill to handle list_files query")
	}
}

func TestDiskAnalysisSkillRejectsListFilesQuery(t *testing.T) {
	skill := &DiskAnalysisSkill{}
	intent := &baseSkill.Intent{
		Type: "incident",
		Entities: map[string]string{
			"checkType":     "disk",
			"queryType":     "list_files",
			"directoryPath": "/root",
		},
	}

	if skill.CanHandle(intent) {
		t.Fatal("expected disk analysis skill to reject list_files query")
	}
}
