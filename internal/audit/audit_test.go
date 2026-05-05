package audit_test

import (
	"testing"

	"github.com/user/envlens/internal/audit"
	"github.com/user/envlens/internal/parser"
)

func makeEnv(pairs ...string) parser.EnvFile {
	env := make(parser.EnvFile)
	for i := 0; i+1 < len(pairs); i += 2 {
		env[pairs[i]] = pairs[i+1]
	}
	return env
}

func TestAudit_EmptyValue(t *testing.T) {
	env := makeEnv("DB_HOST", "")
	report := audit.Audit("test.env", env)
	if len(report.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(report.Findings))
	}
	if report.Findings[0].Severity != audit.SeverityWarning {
		t.Errorf("expected WARNING, got %s", report.Findings[0].Severity)
	}
}

func TestAudit_SecretKey(t *testing.T) {
	env := makeEnv("API_KEY", "abc123")
	report := audit.Audit("test.env", env)
	if len(report.Findings) == 0 {
		t.Fatal("expected at least one finding for secret key")
	}
	found := false
	for _, f := range report.Findings {
		if f.Severity == audit.SeverityCritical {
			found = true
		}
	}
	if !found {
		t.Error("expected a CRITICAL finding for API_KEY")
	}
}

func TestAudit_LocalhostWarning(t *testing.T) {
	env := makeEnv("DB_URL", "postgres://localhost:5432/mydb")
	report := audit.Audit("test.env", env)
	if len(report.Findings) == 0 {
		t.Fatal("expected a finding for localhost")
	}
	if report.Findings[0].Severity != audit.SeverityWarning {
		t.Errorf("expected WARNING, got %s", report.Findings[0].Severity)
	}
}

func TestAudit_CleanEnv(t *testing.T) {
	env := makeEnv("APP_ENV", "production", "PORT", "8080")
	report := audit.Audit("test.env", env)
	if report.HasIssues() {
		t.Errorf("expected no issues, got %+v", report.Findings)
	}
}

func TestReport_HasIssues(t *testing.T) {
	r := audit.Report{
		Findings: []audit.Finding{
			{Key: "X", Message: "ok", Severity: audit.SeverityInfo},
		},
	}
	if r.HasIssues() {
		t.Error("INFO-only report should not HasIssues")
	}
	r.Findings = append(r.Findings, audit.Finding{Key: "Y", Message: "bad", Severity: audit.SeverityWarning})
	if !r.HasIssues() {
		t.Error("report with WARNING should HasIssues")
	}
}
