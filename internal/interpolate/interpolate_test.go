package interpolate

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

func TestInterpolate_NoBraces(t *testing.T) {
	env := makeEnv("HOME", "/home/user", "PATH", "$HOME/bin")
	res := Interpolate(env)
	if got := res.Env["PATH"]; got != "/home/user/bin" {
		t.Errorf("PATH: got %q, want %q", got, "/home/user/bin")
	}
	if len(res.Warnings) != 0 {
		t.Errorf("unexpected warnings: %v", res.Warnings)
	}
}

func TestInterpolate_WithBraces(t *testing.T) {
	env := makeEnv("BASE", "/app", "LOG_DIR", "${BASE}/logs")
	res := Interpolate(env)
	if got := res.Env["LOG_DIR"]; got != "/app/logs" {
		t.Errorf("LOG_DIR: got %q, want %q", got, "/app/logs")
	}
}

func TestInterpolate_UndefinedVariable(t *testing.T) {
	env := makeEnv("URL", "http://${HOST}:8080")
	res := Interpolate(env)
	if got := res.Env["URL"]; got != "http://${HOST}:8080" {
		t.Errorf("URL: got %q, want original", got)
	}
	if len(res.Warnings) == 0 {
		t.Error("expected a warning for undefined variable HOST")
	}
	if !strings.Contains(res.Warnings[0], "HOST") {
		t.Errorf("warning should mention HOST, got: %s", res.Warnings[0])
	}
}

func TestInterpolate_NoReferences(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	res := Interpolate(env)
	if res.Env["FOO"] != "bar" || res.Env["BAZ"] != "qux" {
		t.Error("plain values should be unchanged")
	}
	if len(res.Warnings) != 0 {
		t.Errorf("unexpected warnings: %v", res.Warnings)
	}
}

func TestInterpolate_MultipleReferences(t *testing.T) {
	env := makeEnv(
		"PROTO", "https",
		"HOST", "example.com",
		"PORT", "443",
		"URL", "${PROTO}://${HOST}:${PORT}",
	)
	res := Interpolate(env)
	want := "https://example.com:443"
	if got := res.Env["URL"]; got != want {
		t.Errorf("URL: got %q, want %q", got, want)
	}
}

func TestInterpolate_EmptyValue(t *testing.T) {
	env := makeEnv("EMPTY", "", "REF", "${EMPTY}-suffix")
	res := Interpolate(env)
	if got := res.Env["REF"]; got != "-suffix" {
		t.Errorf("REF: got %q, want %q", got, "-suffix")
	}
}
