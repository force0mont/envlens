package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=production\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(ef.Entries))
	}
	if ef.Index["DB_HOST"].Value != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", ef.Index["DB_HOST"].Value)
	}
}

func TestParse_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=value\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(ef.Entries))
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"` + "\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Index["SECRET"].Value != "my secret value" {
		t.Errorf("expected unquoted value, got %q", ef.Index["SECRET"].Value)
	}
}

func TestParse_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080 # http port\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := ef.Index["PORT"]
	if entry.Value != "8080" {
		t.Errorf("expected value 8080, got %q", entry.Value)
	}
	if entry.Comment != "http port" {
		t.Errorf("expected comment 'http port', got %q", entry.Comment)
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "NOEQUALSIGN\n")
	_, err := Parse(path)
	if err == nil {
		t.Error("expected error for invalid line, got nil")
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
