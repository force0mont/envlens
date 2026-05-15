package envcount_test

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/envcount"
	"github.com/user/envlens/internal/parser"
)

func makeEnv(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestCount_BasicTotals(t *testing.T) {
	env := makeEnv("APP_HOST", "localhost", "APP_PORT", "8080", "DB_URL", "postgres://")
	r := envcount.Count("test", env)
	if r.Total != 3 {
		t.Errorf("expected Total=3, got %d", r.Total)
	}
	if r.NonEmpty != 3 {
		t.Errorf("expected NonEmpty=3, got %d", r.NonEmpty)
	}
	if r.Empty != 0 {
		t.Errorf("expected Empty=0, got %d", r.Empty)
	}
}

func TestCount_EmptyValues(t *testing.T) {
	env := makeEnv("KEY_A", "", "KEY_B", "value", "KEY_C", "")
	r := envcount.Count("test", env)
	if r.Empty != 2 {
		t.Errorf("expected Empty=2, got %d", r.Empty)
	}
	if r.NonEmpty != 1 {
		t.Errorf("expected NonEmpty=1, got %d", r.NonEmpty)
	}
}

func TestCount_PrefixGrouping(t *testing.T) {
	env := makeEnv(
		"APP_HOST", "localhost",
		"APP_PORT", "8080",
		"DB_URL", "postgres://",
		"NOPREFIX", "val",
	)
	r := envcount.Count("test", env)
	if r.Prefixes["APP"] != 2 {
		t.Errorf("expected APP prefix count=2, got %d", r.Prefixes["APP"])
	}
	if r.Prefixes["DB"] != 1 {
		t.Errorf("expected DB prefix count=1, got %d", r.Prefixes["DB"])
	}
	if _, ok := r.Prefixes["NOPREFIX"]; ok {
		t.Error("NOPREFIX should not appear in prefix map")
	}
}

func TestCount_EmptyEnv(t *testing.T) {
	r := envcount.Count("empty", []parser.Entry{})
	if r.Total != 0 || r.Empty != 0 || r.NonEmpty != 0 {
		t.Errorf("expected all zeros for empty env, got %+v", r)
	}
}

func TestFormat_NoFiles(t *testing.T) {
	out := envcount.Format(nil)
	if !strings.Contains(out, "no files") {
		t.Errorf("expected 'no files' message, got: %s", out)
	}
}

func TestFormat_WithResults(t *testing.T) {
	env := makeEnv("APP_HOST", "localhost", "APP_PORT", "", "DB_URL", "postgres://")
	results := []envcount.Result{envcount.Count("prod", env)}
	out := envcount.Format(results)
	if !strings.Contains(out, "[prod]") {
		t.Errorf("expected label in output, got: %s", out)
	}
	if !strings.Contains(out, "total   : 3") {
		t.Errorf("expected total count in output, got: %s", out)
	}
	if !strings.Contains(out, "APP: 2") {
		t.Errorf("expected APP prefix count in output, got: %s", out)
	}
}
