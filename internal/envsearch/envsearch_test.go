package envsearch_test

import (
	"testing"

	"github.com/envlens/internal/envsearch"
	"github.com/envlens/internal/parser"
)

func makeEnv(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestSearch_MatchKey(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "myapp")
	res, err := envsearch.Search(env, envsearch.Options{Pattern: "^DB_", Field: envsearch.MatchKey})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	for _, r := range res {
		if r.Matched != envsearch.MatchKey {
			t.Errorf("expected MatchKey, got %s", r.Matched)
		}
	}
}

func TestSearch_MatchValue(t *testing.T) {
	env := makeEnv("HOST", "localhost", "PORT", "5432", "ENV", "production")
	res, err := envsearch.Search(env, envsearch.Options{Pattern: "localhost", Field: envsearch.MatchValue})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].Key != "HOST" {
		t.Errorf("expected HOST, got %+v", res)
	}
}

func TestSearch_MatchBoth(t *testing.T) {
	env := makeEnv("SECRET_KEY", "abc123", "APP_HOST", "secret.internal")
	res, err := envsearch.Search(env, envsearch.Options{Pattern: "secret", Field: envsearch.MatchBoth})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	env := makeEnv("api_key", "Value123", "OTHER", "nope")
	res, err := envsearch.Search(env, envsearch.Options{Pattern: "API", Field: envsearch.MatchKey, CaseSensitive: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("expected 1 result, got %d", len(res))
	}
}

func TestSearch_NoMatches(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	res, err := envsearch.Search(env, envsearch.Options{Pattern: "NOTFOUND", Field: envsearch.MatchBoth})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected 0 results, got %d", len(res))
	}
}

func TestSearch_InvalidPattern(t *testing.T) {
	env := makeEnv("FOO", "bar")
	_, err := envsearch.Search(env, envsearch.Options{Pattern: "[", Field: envsearch.MatchKey})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestFormat_NoResults(t *testing.T) {
	out := envsearch.Format(nil, "missing")
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_WithResults(t *testing.T) {
	res := []envsearch.Result{{Key: "DB_HOST", Value: "localhost", Matched: envsearch.MatchKey}}
	out := envsearch.Format(res, "DB")
	if out == "" {
		t.Error("expected non-empty output")
	}
}
