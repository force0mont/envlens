package rename

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

func TestRename_BasicRename(t *testing.T) {
	env := makeEnv("OLD_KEY", "value1")
	out, results := Rename(env, map[string]string{"OLD_KEY": "NEW_KEY"})

	if _, ok := out["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if out["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", out["NEW_KEY"])
	}
	if len(results) != 1 || results[0].Missing {
		t.Error("expected one successful rename result")
	}
}

func TestRename_MissingKey(t *testing.T) {
	env := makeEnv("EXISTING", "val")
	_, results := Rename(env, map[string]string{"GHOST": "NEW_GHOST"})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Missing {
		t.Error("expected result to be marked as missing")
	}
}

func TestRename_DoesNotMutateOriginal(t *testing.T) {
	env := makeEnv("A", "1")
	Rename(env, map[string]string{"A": "B"})
	if _, ok := env["A"]; !ok {
		t.Error("original map should not be mutated")
	}
}

func TestRename_MultipleKeys(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "myapp")
	out, results := Rename(env, map[string]string{
		"DB_HOST": "DATABASE_HOST",
		"DB_PORT": "DATABASE_PORT",
	})
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("unexpected DATABASE_HOST value: %q", out["DATABASE_HOST"])
	}
	if out["DATABASE_PORT"] != "5432" {
		t.Errorf("unexpected DATABASE_PORT value: %q", out["DATABASE_PORT"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Error("APP_NAME should be preserved")
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestFormat_SummaryLine(t *testing.T) {
	results := []Result{
		{OldKey: "A", NewKey: "B", Value: "v"},
		{OldKey: "C", NewKey: "D", Missing: true},
	}
	out := Format(results)
	if !strings.Contains(out, "1 renamed, 1 missing") {
		t.Errorf("unexpected summary in output: %q", out)
	}
}

func TestFormat_Empty(t *testing.T) {
	out := Format(nil)
	if !strings.Contains(out, "No renames") {
		t.Errorf("expected empty message, got: %q", out)
	}
}
