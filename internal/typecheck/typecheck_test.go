package typecheck_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/typecheck"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestCheck_ValidInt(t *testing.T) {
	env := makeEnv("PORT", "8080")
	rules := []typecheck.Rule{{Key: "PORT", Type: "int"}}
	findings := typecheck.Check(env, rules)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestCheck_InvalidInt(t *testing.T) {
	env := makeEnv("PORT", "not-a-number")
	rules := []typecheck.Rule{{Key: "PORT", Type: "int"}}
	findings := typecheck.Check(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !strings.Contains(findings[0].Message, "not a valid integer") {
		t.Errorf("unexpected message: %s", findings[0].Message)
	}
}

func TestCheck_ValidBool(t *testing.T) {
	for _, v := range []string{"true", "false", "1", "0", "True", "FALSE"} {
		env := makeEnv("DEBUG", v)
		rules := []typecheck.Rule{{Key: "DEBUG", Type: "bool"}}
		if findings := typecheck.Check(env, rules); len(findings) != 0 {
			t.Errorf("expected no findings for value %q", v)
		}
	}
}

func TestCheck_InvalidBool(t *testing.T) {
	env := makeEnv("DEBUG", "yes")
	rules := []typecheck.Rule{{Key: "DEBUG", Type: "bool"}}
	findings := typecheck.Check(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestCheck_ValidURL(t *testing.T) {
	env := makeEnv("API_URL", "https://example.com/api")
	rules := []typecheck.Rule{{Key: "API_URL", Type: "url"}}
	if findings := typecheck.Check(env, rules); len(findings) != 0 {
		t.Errorf("expected no findings")
	}
}

func TestCheck_InvalidURL(t *testing.T) {
	env := makeEnv("API_URL", "not-a-url")
	rules := []typecheck.Rule{{Key: "API_URL", Type: "url"}}
	findings := typecheck.Check(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestCheck_ValidEmail(t *testing.T) {
	env := makeEnv("ADMIN_EMAIL", "admin@example.com")
	rules := []typecheck.Rule{{Key: "ADMIN_EMAIL", Type: "email"}}
	if findings := typecheck.Check(env, rules); len(findings) != 0 {
		t.Errorf("expected no findings")
	}
}

func TestCheck_MissingKey(t *testing.T) {
	env := makeEnv()
	rules := []typecheck.Rule{{Key: "PORT", Type: "int"}}
	findings := typecheck.Check(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for missing key")
	}
	if findings[0].Message != "key not found" {
		t.Errorf("unexpected message: %s", findings[0].Message)
	}
}

func TestCheck_Nonempty(t *testing.T) {
	env := makeEnv("APP_NAME", "   ")
	rules := []typecheck.Rule{{Key: "APP_NAME", Type: "nonempty"}}
	findings := typecheck.Check(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for blank value")
	}
}

func TestFormat_NoFindings(t *testing.T) {
	out := typecheck.Format(nil)
	if !strings.Contains(out, "all values pass") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormat_WithFindings(t *testing.T) {
	findings := []typecheck.Finding{
		{Key: "PORT", Value: "abc", Expected: "int", Message: `"abc" is not a valid integer`},
	}
	out := typecheck.Format(findings)
	if !strings.Contains(out, "1 violation") {
		t.Errorf("expected violation count in output: %s", out)
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected key in output: %s", out)
	}
}
