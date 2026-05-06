package redact

import (
	"testing"

	"github.com/user/envlens/internal/parser"
)

func makeEnv(pairs ...string) parser.EnvFile {
	env := make(parser.EnvFile)
	for i := 0; i+1 < len(pairs); i += 2 {
		env[pairs[i]] = pairs[i+1]
	}
	return env
}

func TestRedact_SensitiveKeyIsRedacted(t *testing.T) {
	env := makeEnv("DB_PASSWORD", "s3cr3t", "APP_NAME", "myapp")
	out := Redact(env, Options{})
	if out["DB_PASSWORD"] != redactedValue {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", out["APP_NAME"])
	}
}

func TestRedact_TokenKeyIsRedacted(t *testing.T) {
	env := makeEnv("GITHUB_TOKEN", "ghp_abc123")
	out := Redact(env, Options{})
	if out["GITHUB_TOKEN"] != redactedValue {
		t.Errorf("expected GITHUB_TOKEN to be redacted")
	}
}

func TestRedact_ExtraPatternIsRedacted(t *testing.T) {
	env := makeEnv("STRIPE_WEBHOOK", "wh_live_xyz")
	out := Redact(env, Options{ExtraPatterns: []string{"webhook"}})
	if out["STRIPE_WEBHOOK"] != redactedValue {
		t.Errorf("expected STRIPE_WEBHOOK to be redacted with extra pattern")
	}
}

func TestRedact_AllowlistPreservesValue(t *testing.T) {
	env := makeEnv("API_KEY", "public-key-ok")
	out := Redact(env, Options{Allowlist: []string{"API_KEY"}})
	if out["API_KEY"] != "public-key-ok" {
		t.Errorf("expected allowlisted API_KEY to be preserved, got %q", out["API_KEY"])
	}
}

func TestRedact_OriginalEnvUnmodified(t *testing.T) {
	env := makeEnv("SECRET_KEY", "original")
	_ = Redact(env, Options{})
	if env["SECRET_KEY"] != "original" {
		t.Errorf("Redact must not mutate the original env map")
	}
}

func TestRedact_CaseInsensitiveMatch(t *testing.T) {
	env := makeEnv("db_Password", "hunter2")
	out := Redact(env, Options{})
	if out["db_Password"] != redactedValue {
		t.Errorf("expected case-insensitive match to redact db_Password")
	}
}
