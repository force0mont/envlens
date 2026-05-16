package envreplace_test

import (
	"testing"

	"github.com/envlens/internal/envreplace"
	"github.com/envlens/internal/parser"
)

func makeEnv(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestReplace_SubstringMatch(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost:5432", "REDIS_HOST", "localhost:6379")
	r := envreplace.Replace(env, "localhost", "prod.internal", nil, false)
	if len(r.Replacements) != 2 {
		t.Fatalf("expected 2 replacements, got %d", len(r.Replacements))
	}
	if r.Entries[0].Value != "prod.internal:5432" {
		t.Errorf("unexpected value: %s", r.Entries[0].Value)
	}
}

func TestReplace_LiteralMatch(t *testing.T) {
	env := makeEnv("ENV", "staging", "APP_ENV", "staging-eu")
	r := envreplace.Replace(env, "staging", "production", nil, true)
	if len(r.Replacements) != 1 {
		t.Fatalf("expected 1 replacement, got %d", len(r.Replacements))
	}
	if r.Replacements[0].Key != "ENV" {
		t.Errorf("expected ENV to be replaced, got %s", r.Replacements[0].Key)
	}
}

func TestReplace_ScopedToKeys(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "CACHE_HOST", "localhost")
	r := envreplace.Replace(env, "localhost", "remote", []string{"DB_HOST"}, false)
	if len(r.Replacements) != 1 {
		t.Fatalf("expected 1 replacement, got %d", len(r.Replacements))
	}
	if r.Entries[1].Value != "localhost" {
		t.Errorf("CACHE_HOST should be unchanged")
	}
}

func TestReplace_NoMatch(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	r := envreplace.Replace(env, "notfound", "x", nil, false)
	if len(r.Replacements) != 0 {
		t.Errorf("expected 0 replacements, got %d", len(r.Replacements))
	}
}

func TestReplace_DoesNotMutateInput(t *testing.T) {
	env := makeEnv("KEY", "old_value")
	envreplace.Replace(env, "old_value", "new_value", nil, true)
	if env[0].Value != "old_value" {
		t.Errorf("original entries should not be mutated")
	}
}

func TestFormat_NoReplacements(t *testing.T) {
	r := envreplace.Result{}
	out := envreplace.Format(r)
	if out != "No replacements made.\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_WithReplacements(t *testing.T) {
	r := envreplace.Result{
		Replacements: []envreplace.Replacement{
			{Key: "DB_HOST", OldVal: "localhost", NewVal: "prod.db", Changed: true},
		},
	}
	out := envreplace.Format(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
	if !contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output")
	}
	if !contains(out, "1 replacement made") {
		t.Errorf("expected summary line in output")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
