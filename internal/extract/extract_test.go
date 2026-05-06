package extract

import (
	"strings"
	"testing"
)

func makeEnv(pairs ...string) EnvFile {
	m := make(EnvFile)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestExtract_AllFound(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "DB_PORT", "5432", "APP_ENV", "prod")
	results := Extract(env, []string{"DB_HOST", "DB_PORT"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Found {
			t.Errorf("expected key %q to be found", r.Key)
		}
	}
}

func TestExtract_SomeMissing(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost")
	results := Extract(env, []string{"DB_HOST", "DB_PASS"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	foundCount := 0
	for _, r := range results {
		if r.Found {
			foundCount++
		}
	}
	if foundCount != 1 {
		t.Errorf("expected 1 found, got %d", foundCount)
	}
}

func TestExtract_EmptyKeys(t *testing.T) {
	env := makeEnv("A", "1")
	results := Extract(env, []string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestExtract_EmptyValue(t *testing.T) {
	env := makeEnv("EMPTY_KEY", "")
	results := Extract(env, []string{"EMPTY_KEY"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if !results[0].Found {
		t.Error("expected key with empty value to be found")
	}
	if results[0].Value != "" {
		t.Errorf("expected empty value, got %q", results[0].Value)
	}
}

func TestFormat_NoFindings(t *testing.T) {
	out := Format([]Result{})
	if !strings.Contains(out, "no keys requested") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_WithMissing(t *testing.T) {
	results := []Result{
		{Key: "DB_HOST", Value: "localhost", Found: true},
		{Key: "DB_PASS", Value: "", Found: false},
	}
	out := Format(results)
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST in output, got: %q", out)
	}
	if !strings.Contains(out, "DB_PASS=<missing>") {
		t.Errorf("expected DB_PASS missing marker, got: %q", out)
	}
	if !strings.Contains(out, "1 extracted, 1 missing") {
		t.Errorf("expected summary line, got: %q", out)
	}
}

func TestFormat_SortedOutput(t *testing.T) {
	results := []Result{
		{Key: "Z_KEY", Value: "z", Found: true},
		{Key: "A_KEY", Value: "a", Found: true},
	}
	out := Format(results)
	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA > idxZ {
		t.Error("expected output to be sorted alphabetically")
	}
}
