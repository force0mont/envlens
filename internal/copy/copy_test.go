package copy

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

func TestCopy_AllFound(t *testing.T) {
	src := makeEnv("DB_HOST", "prod-db", "API_KEY", "secret")
	dst := makeEnv("APP_NAME", "myapp")
	out, results := Copy(src, dst, []string{"DB_HOST", "API_KEY"}, false)
	if out["DB_HOST"] != "prod-db" {
		t.Errorf("expected DB_HOST=prod-db, got %s", out["DB_HOST"])
	}
	if out["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %s", out["API_KEY"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Error("expected APP_NAME to be preserved")
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestCopy_MissingKey(t *testing.T) {
	src := makeEnv("DB_HOST", "prod-db")
	dst := makeEnv()
	_, results := Copy(src, dst, []string{"DB_HOST", "MISSING_KEY"}, false)
	var missing []Result
	for _, r := range results {
		if r.Missing {
			missing = append(missing, r)
		}
	}
	if len(missing) != 1 || missing[0].Key != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY to be reported missing, got %+v", missing)
	}
}

func TestCopy_NoOverwriteWhenFlagFalse(t *testing.T) {
	src := makeEnv("DB_HOST", "prod-db")
	dst := makeEnv("DB_HOST", "local-db")
	out, results := Copy(src, dst, []string{"DB_HOST"}, false)
	if out["DB_HOST"] != "local-db" {
		t.Errorf("expected DB_HOST to remain local-db, got %s", out["DB_HOST"])
	}
	if results[0].Overwrote {
		t.Error("expected Overwrote=false")
	}
}

func TestCopy_OverwriteWhenFlagTrue(t *testing.T) {
	src := makeEnv("DB_HOST", "prod-db")
	dst := makeEnv("DB_HOST", "local-db")
	out, results := Copy(src, dst, []string{"DB_HOST"}, true)
	if out["DB_HOST"] != "prod-db" {
		t.Errorf("expected DB_HOST=prod-db, got %s", out["DB_HOST"])
	}
	if !results[0].Overwrote {
		t.Error("expected Overwrote=true")
	}
}

func TestCopy_DoesNotMutateDst(t *testing.T) {
	src := makeEnv("DB_HOST", "prod-db")
	dst := makeEnv("APP_NAME", "myapp")
	original := makeEnv("APP_NAME", "myapp")
	Copy(src, dst, []string{"DB_HOST"}, false)
	for k, v := range original {
		if dst[k] != v {
			t.Errorf("dst was mutated: key %s changed", k)
		}
	}
}

func TestFormat_WithMissing(t *testing.T) {
	results := []Result{
		{Key: "DB_HOST", Value: "prod-db"},
		{Key: "MISSING", Missing: true},
		{Key: "API_KEY", Value: "x", Overwrote: true},
	}
	out := Format(results, "prod", "staging")
	if !strings.Contains(out, "prod → staging") {
		t.Error("expected label header in output")
	}
	if !strings.Contains(out, "not found in source") {
		t.Error("expected missing key note")
	}
	if !strings.Contains(out, "overwritten") {
		t.Error("expected overwritten note")
	}
}
