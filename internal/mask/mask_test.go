package mask

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

func TestMask_SensitiveKeyIsObscured(t *testing.T) {
	env := makeEnv("DB_PASSWORD", "supersecret")
	r := Mask(env, nil, 2)
	if r.Masked["DB_PASSWORD"] != "su*********" {
		t.Errorf("unexpected masked value: %q", r.Masked["DB_PASSWORD"])
	}
	if r.MaskCount != 1 {
		t.Errorf("expected MaskCount 1, got %d", r.MaskCount)
	}
}

func TestMask_NonSensitiveKeyPreserved(t *testing.T) {
	env := makeEnv("APP_ENV", "production")
	r := Mask(env, nil, 2)
	if r.Masked["APP_ENV"] != "production" {
		t.Errorf("expected value preserved, got %q", r.Masked["APP_ENV"])
	}
	if r.MaskCount != 0 {
		t.Errorf("expected MaskCount 0, got %d", r.MaskCount)
	}
}

func TestMask_EmptyValueNotMasked(t *testing.T) {
	env := makeEnv("API_KEY", "")
	r := Mask(env, nil, 2)
	if r.Masked["API_KEY"] != "" {
		t.Errorf("expected empty value unchanged, got %q", r.Masked["API_KEY"])
	}
	if r.MaskCount != 0 {
		t.Errorf("expected MaskCount 0 for empty value, got %d", r.MaskCount)
	}
}

func TestMask_ExtraPattern(t *testing.T) {
	env := makeEnv("DEPLOY_PASSPHRASE", "abc123")
	r := Mask(env, []string{"PASSPHRASE"}, 0)
	if r.Masked["DEPLOY_PASSPHRASE"] != "******" {
		t.Errorf("unexpected masked value: %q", r.Masked["DEPLOY_PASSPHRASE"])
	}
}

func TestMask_ZeroVisibleChars(t *testing.T) {
	env := makeEnv("SECRET_KEY", "hello")
	r := Mask(env, nil, 0)
	if r.Masked["SECRET_KEY"] != "*****" {
		t.Errorf("expected all stars, got %q", r.Masked["SECRET_KEY"])
	}
}

func TestFormat_OutputContainsKeys(t *testing.T) {
	env := makeEnv("APP_NAME", "envlens", "API_TOKEN", "tok123")
	r := Mask(env, nil, 0)
	out := Format(r)
	if !strings.Contains(out, "APP_NAME=envlens") {
		t.Errorf("expected APP_NAME in output, got:\n%s", out)
	}
	if !strings.Contains(out, "API_TOKEN=") {
		t.Errorf("expected API_TOKEN in output, got:\n%s", out)
	}
}
