package skill

import "testing"

func TestValidateManagedFileOperation(t *testing.T) {
	if err := validateManagedFileOperation("clear", "/var/log/app.log"); err != nil {
		t.Fatalf("expected clear log to be allowed, got %v", err)
	}
	if err := validateManagedFileOperation("delete", "/tmp/debug.txt"); err != nil {
		t.Fatalf("expected delete txt to be allowed, got %v", err)
	}
	if err := validateManagedFileOperation("delete", "/etc/passwd"); err == nil {
		t.Fatal("expected non-whitelist suffix to be rejected")
	}
	if err := validateManagedFileOperation("delete", "/tmp/archive.zip"); err == nil {
		t.Fatal("expected zip suffix to be rejected")
	}
}

func TestBuildManagedFileOperationCommand(t *testing.T) {
	clearCmd := buildManagedFileOperationCommand("clear", "/var/log/app.log")
	if clearCmd != "if test -f /var/log/app.log; then : > /var/log/app.log && printf 'file=/var/log/app.log\\naction=clear\\nstatus=cleared\\n'; else printf 'file=/var/log/app.log\\naction=clear\\nstatus=not_found\\n'; fi" {
		t.Fatalf("unexpected clear command: %s", clearCmd)
	}

	deleteCmd := buildManagedFileOperationCommand("delete", "/var/log/app.log")
	if deleteCmd != "if test -f /var/log/app.log; then rm -f /var/log/app.log && printf 'file=/var/log/app.log\\naction=delete\\nstatus=deleted\\nremaining_dir=/var/log\\n'; ls -1A /var/log 2>/dev/null | sed 's/^/remaining_file=/'; else printf 'file=/var/log/app.log\\naction=delete\\nstatus=not_found\\n'; fi" {
		t.Fatalf("unexpected delete command: %s", deleteCmd)
	}
}

func TestValidateDeleteCommand(t *testing.T) {
	if err := validateDeleteCommand("rm -f /tmp/a.log /tmp/b.txt /tmp/c.gz"); err != nil {
		t.Fatalf("expected allowed delete command, got %v", err)
	}
	if err := validateDeleteCommand("rm -f /root/gin-api.sql"); err == nil {
		t.Fatal("expected sql delete to be rejected")
	}
	if err := validateDeleteCommand("rm -rf /tmp/testdir"); err == nil {
		t.Fatal("expected recursive delete to be rejected")
	}
	if err := validateDeleteCommand("rm -f /tmp/*.log"); err == nil {
		t.Fatal("expected wildcard delete to be rejected")
	}
}
