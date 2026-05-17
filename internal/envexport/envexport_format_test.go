package envexport_test

import (
	"strings"
	"testing"

	"github.com/envlens/envlens/internal/envexport"
)

func TestExport_SortedOutput(t *testing.T) {
	entries := makeEnv("ZEBRA", "z", "ALPHA", "a", "MIDDLE", "m")
	res, err := envexport.Export(entries, envexport.Options{OutputFormat: envexport.FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(res.Content), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected ALPHA first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected ZEBRA last, got: %s", lines[2])
	}
}

func TestExport_JSONStructure(t *testing.T) {
	entries := makeEnv("HOST", "db", "PORT", "5432")
	res, err := envexport.Export(entries, envexport.Options{OutputFormat: envexport.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(res.Content), "{") {
		t.Errorf("expected JSON to start with '{', got: %s", res.Content)
	}
	if !strings.HasSuffix(strings.TrimSpace(res.Content), "}") {
		t.Errorf("expected JSON to end with '}', got: %s", res.Content)
	}
}

func TestExport_EmptyEntries(t *testing.T) {
	res, err := envexport.Export(nil, envexport.Options{OutputFormat: envexport.FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Count != 0 {
		t.Errorf("expected count 0, got %d", res.Count)
	}
	if strings.TrimSpace(res.Content) != "" {
		t.Errorf("expected empty output, got: %s", res.Content)
	}
}

func TestExport_FormatPreservedInResult(t *testing.T) {
	entries := makeEnv("K", "v")
	res, err := envexport.Export(entries, envexport.Options{OutputFormat: envexport.FormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Format != envexport.FormatShell {
		t.Errorf("expected format shell in result, got: %s", res.Format)
	}
}
