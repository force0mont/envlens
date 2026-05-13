package scope

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

func TestTag_BasicScoping(t *testing.T) {
	dev := makeEnv("APP_ENV", "development", "PORT", "3000")
	prod := makeEnv("APP_ENV", "production", "PORT", "8080")

	r, err := Tag([]map[string]string{dev, prod}, []string{"dev", "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
	for _, e := range r.Entries {
		if e.Scope != "prod" {
			t.Errorf("expected scope prod, got %s for key %s", e.Scope, e.Key)
		}
	}
}

func TestTag_UniqueKeys(t *testing.T) {
	a := makeEnv("ONLY_A", "1")
	b := makeEnv("ONLY_B", "2")

	r, err := Tag([]map[string]string{a, b}, []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
}

func TestTag_MismatchedLengths(t *testing.T) {
	_, err := Tag([]map[string]string{makeEnv("K", "V")}, []string{"a", "b"})
	if err == nil {
		t.Fatal("expected error for mismatched lengths")
	}
}

func TestTag_EmptyEnvs(t *testing.T) {
	r, err := Tag([]map[string]string{}, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(r.Entries))
	}
}

func TestFormat_WithEntries(t *testing.T) {
	r := Result{
		Entries: []Entry{
			{Key: "FOO", Value: "bar", Scope: "dev"},
		},
	}
	out := Format(r)
	if !strings.Contains(out, "[dev] FOO=bar") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormat_Empty(t *testing.T) {
	out := Format(Result{})
	if !strings.Contains(out, "no entries") {
		t.Errorf("expected 'no entries', got: %s", out)
	}
}
