package template

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

func TestRender_NoPlaceholders(t *testing.T) {
	src := makeEnv("HOST", "localhost", "PORT", "5432")
	vars := makeEnv()
	r := Render(src, vars)
	if len(r.Missing) != 0 {
		t.Fatalf("expected no missing, got %v", r.Missing)
	}
	for _, e := range r.Entries {
		if e.Key == "HOST" && e.Value != "localhost" {
			t.Errorf("HOST: want localhost, got %s", e.Value)
		}
	}
}

func TestRender_SubstitutesKnownVar(t *testing.T) {
	src := makeEnv("DSN", "postgres://${DB_HOST}:${DB_PORT}/mydb")
	vars := makeEnv("DB_HOST", "db.internal", "DB_PORT", "5432")
	r := Render(src, vars)
	if len(r.Missing) != 0 {
		t.Fatalf("unexpected missing: %v", r.Missing)
	}
	if r.Entries[0].Value != "postgres://db.internal:5432/mydb" {
		t.Errorf("unexpected value: %s", r.Entries[0].Value)
	}
}

func TestRender_RecordsMissingVar(t *testing.T) {
	src := makeEnv("URL", "https://${API_HOST}/v1")
	vars := makeEnv()
	r := Render(src, vars)
	if len(r.Missing) != 1 || r.Missing[0] != "API_HOST" {
		t.Fatalf("expected [API_HOST] missing, got %v", r.Missing)
	}
	if !strings.Contains(r.Entries[0].Value, "${API_HOST}") {
		t.Errorf("placeholder should remain intact: %s", r.Entries[0].Value)
	}
}

func TestRender_DeduplicatesMissing(t *testing.T) {
	src := makeEnv("A", "${X}-${X}", "B", "${X}")
	vars := makeEnv()
	r := Render(src, vars)
	if len(r.Missing) != 1 {
		t.Errorf("expected 1 unique missing var, got %d: %v", len(r.Missing), r.Missing)
	}
}

func TestRender_PartialSubstitution(t *testing.T) {
	src := makeEnv("CONN", "${KNOWN}:${UNKNOWN}")
	vars := makeEnv("KNOWN", "resolved")
	r := Render(src, vars)
	if !strings.HasPrefix(r.Entries[0].Value, "resolved:") {
		t.Errorf("expected partial substitution, got: %s", r.Entries[0].Value)
	}
	if len(r.Missing) != 1 || r.Missing[0] != "UNKNOWN" {
		t.Errorf("expected UNKNOWN in missing, got %v", r.Missing)
	}
}

func TestFormat_NoEntries(t *testing.T) {
	r := Result{}
	out := Format(r)
	if out != "no entries\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_WithMissing(t *testing.T) {
	src := makeEnv("KEY", "${MISSING_VAR}")
	r := Render(src, makeEnv())
	out := Format(r)
	if !strings.Contains(out, "unresolved") {
		t.Errorf("expected 'unresolved' in output: %s", out)
	}
	if !strings.Contains(out, "MISSING_VAR") {
		t.Errorf("expected variable name in output: %s", out)
	}
}
