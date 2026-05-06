package compare

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/parser"
)

func makeEnv(pairs ...string) parser.EnvFile {
	env := parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		env[pairs[i]] = pairs[i+1]
	}
	return env
}

func TestCompare_IdenticalValues(t *testing.T) {
	envs := map[string]parser.EnvFile{
		"a": makeEnv("HOST", "localhost"),
		"b": makeEnv("HOST", "localhost"),
	}
	results := Compare(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Unique {
		t.Error("expected Unique=true for identical values")
	}
}

func TestCompare_DifferentValues(t *testing.T) {
	envs := map[string]parser.EnvFile{
		"prod": makeEnv("DB_HOST", "db.prod.example.com"),
		"dev":  makeEnv("DB_HOST", "localhost"),
	}
	results := Compare(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Unique {
		t.Error("expected Unique=false for differing values")
	}
}

func TestCompare_MissingKeyInOneFile(t *testing.T) {
	envs := map[string]parser.EnvFile{
		"a": makeEnv("ONLY_A", "yes"),
		"b": makeEnv("ONLY_B", "yes"),
	}
	results := Compare(envs)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Unique {
			t.Errorf("key %s should not be Unique (missing in one file)", r.Key)
		}
	}
}

func TestCompare_MultipleFiles(t *testing.T) {
	envs := map[string]parser.EnvFile{
		"a": makeEnv("PORT", "8080", "HOST", "localhost"),
		"b": makeEnv("PORT", "9090", "HOST", "localhost"),
		"c": makeEnv("PORT", "8080", "HOST", "localhost"),
	}
	results := Compare(envs)
	for _, r := range results {
		switch r.Key {
		case "PORT":
			if r.Unique {
				t.Error("PORT has differing values; Unique should be false")
			}
		case "HOST":
			if !r.Unique {
				t.Error("HOST has identical values; Unique should be true")
			}
		}
	}
}

func TestFormat_ContainsMissingMarker(t *testing.T) {
	envs := map[string]parser.EnvFile{
		"a": makeEnv("SECRET", "abc"),
		"b": {},
	}
	results := Compare(envs)
	out := Format(results, []string{"a", "b"})
	if !strings.Contains(out, "<missing>") {
		t.Error("expected <missing> marker in output")
	}
}

func TestFormat_EmptyResults(t *testing.T) {
	out := Format([]Result{}, []string{"a", "b"})
	if !strings.Contains(out, "no keys found") {
		t.Errorf("unexpected output for empty results: %q", out)
	}
}
