package trim

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

func TestTrim_NoChanges(t *testing.T) {
	env := makeEnv("KEY", "value", "HOST", "localhost")
	out, results := Trim(env)

	for _, r := range results {
		if r.Changed {
			t.Errorf("expected no changes, but %q was marked changed", r.Key)
		}
	}
	if out["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", out["KEY"])
	}
}

func TestTrim_LeadingAndTrailingSpaces(t *testing.T) {
	env := makeEnv("KEY", "  hello world  ")
	out, results := Trim(env)

	if out["KEY"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", out["KEY"])
	}
	if len(results) != 1 || !results[0].Changed {
		t.Error("expected result to be marked as changed")
	}
}

func TestTrim_TabsAndNewlines(t *testing.T) {
	env := makeEnv("TOKEN", "\t\nmy-token\n\t")
	out, _ := Trim(env)

	if out["TOKEN"] != "my-token" {
		t.Errorf("expected 'my-token', got %q", out["TOKEN"])
	}
}

func TestTrim_EmptyValue(t *testing.T) {
	env := makeEnv("EMPTY", "")
	out, results := Trim(env)

	if out["EMPTY"] != "" {
		t.Errorf("expected empty string, got %q", out["EMPTY"])
	}
	if results[0].Changed {
		t.Error("empty value should not be marked as changed")
	}
}

func TestTrim_DoesNotMutateOriginal(t *testing.T) {
	env := makeEnv("KEY", "  padded  ")
	Trim(env)

	if env["KEY"] != "  padded  " {
		t.Error("original map should not be mutated")
	}
}

func TestFormat_NoChanges(t *testing.T) {
	env := makeEnv("A", "clean")
	_, results := Trim(env)
	out := Format(results)

	if !strings.Contains(out, "No values required trimming") {
		t.Errorf("expected no-change message, got: %s", out)
	}
}

func TestFormat_WithChanges(t *testing.T) {
	env := makeEnv("API_KEY", " secret ")
	_, results := Trim(env)
	out := Format(results)

	if !strings.Contains(out, "1 value(s) trimmed") {
		t.Errorf("expected trimmed count, got: %s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected key in output, got: %s", out)
	}
}
