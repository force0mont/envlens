package filter_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/filter"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestFilter_ByPrefix(t *testing.T) {
	env := makeEnv("APP_HOST", "localhost", "APP_PORT", "8080", "DB_URL", "postgres://")
	r, err := filter.Filter(env, filter.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
	if r.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", r.Skipped)
	}
}

func TestFilter_ByPattern(t *testing.T) {
	env := makeEnv("SECRET_KEY", "abc", "API_KEY", "xyz", "HOST", "localhost")
	r, err := filter.Filter(env, filter.Options{Pattern: "KEY$"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
}

func TestFilter_ByKeys(t *testing.T) {
	env := makeEnv("A", "1", "B", "2", "C", "3")
	r, err := filter.Filter(env, filter.Options{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
	if _, ok := r.Matched["B"]; ok {
		t.Error("B should not be in matched")
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	env := makeEnv("FOO", "bar")
	_, err := filter.Filter(env, filter.Options{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	env := makeEnv("X", "1", "Y", "2")
	r, err := filter.Filter(env, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
}

func TestFilter_NoMatches(t *testing.T) {
	env := makeEnv("FOO", "bar")
	r, err := filter.Filter(env, filter.Options{Prefix: "NOPE_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := filter.Format(r)
	if !strings.Contains(out, "No matching") {
		t.Errorf("expected 'No matching' in output, got: %s", out)
	}
}

func TestFormat_WithMatches(t *testing.T) {
	env := makeEnv("APP_HOST", "localhost", "APP_PORT", "8080")
	r, _ := filter.Filter(env, filter.Options{Prefix: "APP_"})
	out := filter.Format(r)
	if !strings.Contains(out, "APP_HOST") {
		t.Errorf("expected APP_HOST in output")
	}
	if !strings.Contains(out, "Matched 2") {
		t.Errorf("expected match count in output")
	}
}
