package lint

import (
	"strings"
	"testing"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestLint_CleanEnv(t *testing.T) {
	env := makeEnv("APP_HOST", "localhost", "DB_PORT", "5432")
	findings := Lint(env)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := makeEnv("app_host", "localhost")
	findings := Lint(env)
	if !hasFindings(findings, "app_host", SeverityError) {
		t.Error("expected error finding for lowercase key")
	}
}

func TestLint_MixedCaseKey(t *testing.T) {
	env := makeEnv("AppHost", "localhost")
	findings := Lint(env)
	if !hasFindings(findings, "AppHost", SeverityError) {
		t.Error("expected error finding for mixed-case key")
	}
}

func TestLint_DoubleUnderscore(t *testing.T) {
	env := makeEnv("APP__HOST", "localhost")
	findings := Lint(env)
	if !hasFindings(findings, "APP__HOST", SeverityWarning) {
		t.Error("expected warning for double underscore in key")
	}
}

func TestLint_LeadingWhitespaceInValue(t *testing.T) {
	env := makeEnv("APP_HOST", " localhost")
	findings := Lint(env)
	if !hasFindings(findings, "APP_HOST", SeverityWarning) {
		t.Error("expected warning for leading whitespace in value")
	}
}

func TestLint_MultipleSpacesInValue(t *testing.T) {
	env := makeEnv("APP_DESC", "hello  world")
	findings := Lint(env)
	if !hasFindings(findings, "APP_DESC", SeverityInfo) {
		t.Error("expected info finding for multiple spaces in value")
	}
}

func TestFormat_NoFindings(t *testing.T) {
	out := Format(nil)
	if !strings.Contains(out, "No lint issues") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_WithFindings(t *testing.T) {
	findings := []Finding{
		{Key: "bad_key", Message: "key should be UPPER_SNAKE_CASE", Severity: SeverityError},
	}
	out := Format(findings)
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected ERROR in output, got: %q", out)
	}
	if !strings.Contains(out, "1 issue(s)") {
		t.Errorf("expected issue count in output, got: %q", out)
	}
}

func hasFindings(findings []Finding, key string, sev Severity) bool {
	for _, f := range findings {
		if f.Key == key && f.Severity == sev {
			return true
		}
	}
	return false
}
