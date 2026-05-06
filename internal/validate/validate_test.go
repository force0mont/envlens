package validate

import (
	"strings"
	"testing"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestValidate_RequiredPresent(t *testing.T) {
	env := makeEnv("DATABASE_URL", "postgres://localhost/db")
	rules := []Rule{{Key: "DATABASE_URL", Required: true}}
	findings := Validate(env, rules)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %v", findings)
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	env := makeEnv()
	rules := []Rule{{Key: "DATABASE_URL", Required: true}}
	findings := Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !strings.Contains(findings[0].Message, "missing") {
		t.Errorf("unexpected message: %s", findings[0].Message)
	}
}

func TestValidate_RequiredEmptyValue(t *testing.T) {
	env := makeEnv("PORT", "   ")
	rules := []Rule{{Key: "PORT", Required: true}}
	findings := Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !strings.Contains(findings[0].Message, "empty") {
		t.Errorf("unexpected message: %s", findings[0].Message)
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	env := makeEnv("PORT", "8080")
	rules := []Rule{{Key: "PORT", Pattern: `^\d+$`}}
	findings := Validate(env, rules)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %v", findings)
	}
}

func TestValidate_PatternNoMatch(t *testing.T) {
	env := makeEnv("PORT", "not-a-port")
	rules := []Rule{{Key: "PORT", Pattern: `^\d+$`}}
	findings := Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !strings.Contains(findings[0].Message, "does not match") {
		t.Errorf("unexpected message: %s", findings[0].Message)
	}
}

func TestValidate_OptionalMissingSkipped(t *testing.T) {
	env := makeEnv()
	rules := []Rule{{Key: "OPTIONAL_KEY", Pattern: `^\d+$`}}
	findings := Validate(env, rules)
	if len(findings) != 0 {
		t.Errorf("expected no findings for missing optional key, got %v", findings)
	}
}

func TestFormat_AllPassed(t *testing.T) {
	out := Format(nil, "staging")
	if !strings.Contains(out, "all rules passed") {
		t.Errorf("expected pass message, got: %s", out)
	}
}

func TestFormat_WithFindings(t *testing.T) {
	findings := []Finding{{Key: "DATABASE_URL", Message: "required key is missing"}}
	out := Format(findings, "prod")
	if !strings.Contains(out, "DATABASE_URL") || !strings.Contains(out, "1 validation") {
		t.Errorf("unexpected output: %s", out)
	}
}
