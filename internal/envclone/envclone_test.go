package envclone_test

import (
	"strings"
	"testing"

	"github.com/envlens/internal/envclone"
	"github.com/envlens/internal/parser"
)

func makeEnv(kv ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(kv); i += 2 {
		entries = append(entries, parser.Entry{Key: kv[i], Value: kv[i+1]})
	}
	return entries
}

func findKey(entries []parser.Entry, key string) (string, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}

func TestClone_NewDestKey(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost")
	out, results := envclone.Clone(env, map[string]string{"DB_HOST": "DATABASE_HOST"}, false)
	if len(results) != 1 || results[0].Skipped {
		t.Fatal("expected one non-skipped result")
	}
	val, ok := findKey(out, "DATABASE_HOST")
	if !ok || val != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q ok=%v", val, ok)
	}
}

func TestClone_NoOverwriteByDefault(t *testing.T) {
	env := makeEnv("SRC", "new", "DST", "old")
	_, results := envclone.Clone(env, map[string]string{"SRC": "DST"}, false)
	if len(results) != 1 || !results[0].Skipped {
		t.Fatal("expected result to be skipped")
	}
}

func TestClone_OverwriteExisting(t *testing.T) {
	env := makeEnv("SRC", "new", "DST", "old")
	out, results := envclone.Clone(env, map[string]string{"SRC": "DST"}, true)
	if len(results) != 1 || results[0].Skipped {
		t.Fatal("expected non-skipped result")
	}
	val, _ := findKey(out, "DST")
	if val != "new" {
		t.Errorf("expected DST=new, got %q", val)
	}
}

func TestClone_MissingSourceSkipped(t *testing.T) {
	env := makeEnv("A", "1")
	_, results := envclone.Clone(env, map[string]string{"MISSING": "DEST"}, false)
	if len(results) != 1 || !results[0].Skipped {
		t.Fatal("expected skip for missing source")
	}
}

func TestClone_DoesNotMutateOriginal(t *testing.T) {
	env := makeEnv("X", "val")
	orig := make([]parser.Entry, len(env))
	copy(orig, env)
	envclone.Clone(env, map[string]string{"X": "Y"}, false)
	if len(env) != len(orig) || env[0].Key != orig[0].Key {
		t.Error("original slice was mutated")
	}
}

func TestFormat_SummaryLine(t *testing.T) {
	results := []envclone.Result{
		{SourceKey: "A", DestKey: "B", Value: "v"},
		{SourceKey: "C", DestKey: "D", Skipped: true},
	}
	out := envclone.Format(results)
	if !strings.Contains(out, "1 cloned, 1 skipped") {
		t.Errorf("unexpected summary: %q", out)
	}
}

func TestFormat_NoResults(t *testing.T) {
	out := envclone.Format(nil)
	if !strings.Contains(out, "no keys cloned") {
		t.Errorf("expected empty message, got %q", out)
	}
}
