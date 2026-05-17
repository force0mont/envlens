package envexport_test

import (
	"strings"
	"testing"

	"github.com/envlens/envlens/internal/envexport"
	"github.com/envlens/envlens/internal/parser"
)

func makeEnv(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestExport_DotenvFormat(t *testing.T) {
	entries := makeEnv("APP_ENV", "production", "PORT", "8080")
	res, err := envexport.Export(entries, envexport.Options{OutputFormat: envexport.FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Content, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", res.Content)
	}
	if res.Count != 2 {
		t.Errorf("expected count 2, got %d", res.Count)
	}
}

func TestExport_ShellFormat(t *testing.T) {
	entries := makeEnv("DB_PASS", "secret")
	res, err := envexport.Export(entries, envexport.Options{OutputFormat: envexport.FormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(res.Content, "export ") {
		t.Errorf("expected shell export prefix, got: %s", res.Content)
	}
}

func TestExport_DockerFormat(t *testing.T) {
	entries := makeEnv("APP_ENV", "staging")
	res, err := envexport.Export(entries, envexport.Options{OutputFormat: envexport.FormatDocker})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Content, "--env APP_ENV=staging") {
		t.Errorf("expected docker flag format, got: %s", res.Content)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	entries := makeEnv("KEY", "val")
	res, err := envexport.Export(entries, envexport.Options{OutputFormat: envexport.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Content, `"KEY": "val"`) {
		t.Errorf("expected JSON output, got: %s", res.Content)
	}
}

func TestExport_OmitEmpty(t *testing.T) {
	entries := makeEnv("A", "hello", "B", "", "C", "world")
	res, err := envexport.Export(entries, envexport.Options{OutputFormat: envexport.FormatDotenv, OmitEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Count != 2 {
		t.Errorf("expected 2 entries after omitting empty, got %d", res.Count)
	}
	if strings.Contains(res.Content, "B=") {
		t.Errorf("empty key B should have been omitted")
	}
}

func TestExport_WithPrefix(t *testing.T) {
	entries := makeEnv("HOST", "localhost")
	res, err := envexport.Export(entries, envexport.Options{OutputFormat: envexport.FormatDotenv, Prefix: "MY_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Content, "MY_HOST=localhost") {
		t.Errorf("expected prefixed key, got: %s", res.Content)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	entries := makeEnv("A", "1")
	_, err := envexport.Export(entries, envexport.Options{OutputFormat: "xml"})
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
