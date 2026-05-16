package envnorm

import (
	"testing"

	"github.com/user/envlens/internal/parser"
)

func makeEnv(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestNormalise_Uppercase(t *testing.T) {
	entries := makeEnv("db_host", "localhost", "api_key", "secret")
	results := Normalise(entries, []Rule{RuleUppercase})
	if results[0].Normalised.Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", results[0].Normalised.Key)
	}
	if !results[0].Changed {
		t.Error("expected Changed=true")
	}
}

func TestNormalise_Lowercase(t *testing.T) {
	entries := makeEnv("DB_HOST", "localhost")
	results := Normalise(entries, []Rule{RuleLowercase})
	if results[0].Normalised.Key != "db_host" {
		t.Errorf("expected db_host, got %s", results[0].Normalised.Key)
	}
}

func TestNormalise_TrimValues(t *testing.T) {
	entries := makeEnv("HOST", "  localhost  ")
	results := Normalise(entries, []Rule{RuleTrimValues})
	if results[0].Normalised.Value != "localhost" {
		t.Errorf("expected 'localhost', got %q", results[0].Normalised.Value)
	}
	if !results[0].Changed {
		t.Error("expected Changed=true")
	}
}

func TestNormalise_SnakeCase(t *testing.T) {
	entries := makeEnv("my-key", "val", "another-env-var", "42")
	results := Normalise(entries, []Rule{RuleSnakeCase})
	if results[0].Normalised.Key != "my_key" {
		t.Errorf("expected my_key, got %s", results[0].Normalised.Key)
	}
	if results[1].Normalised.Key != "another_env_var" {
		t.Errorf("expected another_env_var, got %s", results[1].Normalised.Key)
	}
}

func TestNormalise_NoChange(t *testing.T) {
	entries := makeEnv("DB_HOST", "localhost")
	results := Normalise(entries, []Rule{RuleUppercase})
	if results[0].Changed {
		t.Error("expected Changed=false for already-uppercase key")
	}
}

func TestNormalise_MultipleRules(t *testing.T) {
	entries := makeEnv("my-db-host", "  prod  ")
	results := Normalise(entries, []Rule{RuleSnakeCase, RuleUppercase, RuleTrimValues})
	if results[0].Normalised.Key != "MY_DB_HOST" {
		t.Errorf("expected MY_DB_HOST, got %s", results[0].Normalised.Key)
	}
	if results[0].Normalised.Value != "prod" {
		t.Errorf("expected 'prod', got %q", results[0].Normalised.Value)
	}
}

func TestToEntries_ExtractsNormalised(t *testing.T) {
	entries := makeEnv("my-key", "val")
	results := Normalise(entries, []Rule{RuleSnakeCase})
	out := ToEntries(results)
	if len(out) != 1 || out[0].Key != "my_key" {
		t.Errorf("unexpected entries: %+v", out)
	}
}
