package audit_test

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/audit"
)

func TestFormat_NoFindings(t *testing.T) {
	r := audit.Report{File: "prod.env", Findings: nil}
	out := audit.Format(r)
	if !strings.Contains(out, "No issues found") {
		t.Errorf("expected 'No issues found' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "prod.env") {
		t.Errorf("expected filename in output")
	}
}

func TestFormat_WithFindings(t *testing.T) {
	r := audit.Report{
		File: "staging.env",
		Findings: []audit.Finding{
			{Key: "DB_PASSWORD", Message: "potential secret detected", Severity: audit.SeverityCritical},
			{Key: "DB_HOST", Message: "value is empty", Severity: audit.SeverityWarning},
		},
	}
	out := audit.Format(r)

	if !strings.Contains(out, "CRITICAL") {
		t.Error("expected CRITICAL in output")
	}
	if !strings.Contains(out, "WARNING") {
		t.Error("expected WARNING in output")
	}
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD in output")
	}
	if !strings.Contains(out, "Summary:") {
		t.Error("expected Summary line in output")
	}
}

func TestFormat_SummaryCount(t *testing.T) {
	r := audit.Report{
		File: "dev.env",
		Findings: []audit.Finding{
			{Key: "A", Message: "m1", Severity: audit.SeverityCritical},
			{Key: "B", Message: "m2", Severity: audit.SeverityCritical},
			{Key: "C", Message: "m3", Severity: audit.SeverityWarning},
		},
	}
	out := audit.Format(r)
	if !strings.Contains(out, "2 critical") {
		t.Errorf("expected '2 critical' in summary, got:\n%s", out)
	}
	if !strings.Contains(out, "1 warning") {
		t.Errorf("expected '1 warning' in summary, got:\n%s", out)
	}
}
