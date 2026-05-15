package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envlens/cmd"
)

func writeTempEnvForSearch(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func TestSearchCmd_MatchesKey(t *testing.T) {
	path := writeTempEnvForSearch(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=myapp\n")
	out, err := cmd.ExecuteWithArgs("search", "^DB_", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "DB_PORT") {
		t.Errorf("expected DB keys in output, got: %s", out)
	}
	if strings.Contains(out, "APP_NAME") {
		t.Errorf("APP_NAME should not appear in output")
	}
}

func TestSearchCmd_MatchesValue(t *testing.T) {
	path := writeTempEnvForSearch(t, "HOST=localhost\nENV=production\n")
	out, err := cmd.ExecuteWithArgs("search", "--field=value", "localhost", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST in output, got: %s", out)
	}
}

func TestSearchCmd_NoMatches(t *testing.T) {
	path := writeTempEnvForSearch(t, "FOO=bar\nBAZ=qux\n")
	out, err := cmd.ExecuteWithArgs("search", "NOTFOUND", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No matches") {
		t.Errorf("expected no-match message, got: %s", out)
	}
}

func TestSearchCmd_MissingFile(t *testing.T) {
	_, err := cmd.ExecuteWithArgs("search", "FOO", "/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSearchCmd_InvalidPattern(t *testing.T) {
	path := writeTempEnvForSearch(t, "FOO=bar\n")
	_, err := cmd.ExecuteWithArgs("search", "[", path)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}
