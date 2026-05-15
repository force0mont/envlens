package envcast_test

import (
	"testing"

	"github.com/envlens/envlens/internal/envcast"
	"github.com/envlens/envlens/internal/parser"
)

func makeEnv(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestCast_StringType(t *testing.T) {
	env := makeEnv("NAME", "alice")
	results := envcast.Cast(env, []envcast.Rule{{Key: "NAME", Type: envcast.TypeString}})
	if len(results) != 1 || results[0].Err != nil {
		t.Fatalf("expected successful string cast, got %+v", results)
	}
	if results[0].CastValue != "alice" {
		t.Errorf("expected 'alice', got %q", results[0].CastValue)
	}
}

func TestCast_IntValid(t *testing.T) {
	env := makeEnv("PORT", "8080")
	results := envcast.Cast(env, []envcast.Rule{{Key: "PORT", Type: envcast.TypeInt}})
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].CastValue != "8080" {
		t.Errorf("expected '8080', got %q", results[0].CastValue)
	}
}

func TestCast_IntInvalid(t *testing.T) {
	env := makeEnv("PORT", "not-a-number")
	results := envcast.Cast(env, []envcast.Rule{{Key: "PORT", Type: envcast.TypeInt}})
	if results[0].Err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestCast_FloatValid(t *testing.T) {
	env := makeEnv("RATIO", "3.14")
	results := envcast.Cast(env, []envcast.Rule{{Key: "RATIO", Type: envcast.TypeFloat}})
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].CastValue != "3.14" {
		t.Errorf("expected '3.14', got %q", results[0].CastValue)
	}
}

func TestCast_BoolValid(t *testing.T) {
	env := makeEnv("DEBUG", "true")
	results := envcast.Cast(env, []envcast.Rule{{Key: "DEBUG", Type: envcast.TypeBool}})
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].CastValue != "true" {
		t.Errorf("expected 'true', got %q", results[0].CastValue)
	}
}

func TestCast_BoolInvalid(t *testing.T) {
	env := makeEnv("DEBUG", "yes-please")
	results := envcast.Cast(env, []envcast.Rule{{Key: "DEBUG", Type: envcast.TypeBool}})
	if results[0].Err == nil {
		t.Fatal("expected error for invalid bool")
	}
}

func TestCast_MissingKey(t *testing.T) {
	env := makeEnv("PORT", "9000")
	results := envcast.Cast(env, []envcast.Rule{{Key: "MISSING", Type: envcast.TypeInt}})
	if results[0].Err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestFormat_WithErrors(t *testing.T) {
	env := makeEnv("PORT", "bad")
	results := envcast.Cast(env, []envcast.Rule{{Key: "PORT", Type: envcast.TypeInt}})
	out := envcast.Format(results)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}

func TestFormat_Empty(t *testing.T) {
	out := envcast.Format(nil)
	if out != "No cast rules applied.\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
