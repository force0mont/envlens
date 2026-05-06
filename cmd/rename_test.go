package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvForRename(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestRenameCmd_BasicRename(t *testing.T) {
	path := writeTempEnvForRename(t, "OLD_KEY=hello\nOTHER=world\n")

	out, err := executeCommand(rootCmd, "rename", path, "--pair", "OLD_KEY=NEW_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "RENAMED") {
		t.Errorf("expected RENAMED in output, got: %q", out)
	}
}

func TestRenameCmd_MissingKey(t *testing.T) {
	path := writeTempEnvForRename(t, "EXISTING=val\n")

	out, err := executeCommand(rootCmd, "rename", path, "--pair", "GHOST=NEW_GHOST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING in output, got: %q", out)
	}
}

func TestRenameCmd_InvalidPairFormat(t *testing.T) {
	path := writeTempEnvForRename(t, "KEY=val\n")

	_, err := executeCommand(rootCmd, "rename", path, "--pair", "BADFORMAT")
	if err == nil {
		t.Error("expected error for invalid pair format")
	}
}

func TestRenameCmd_WritesOutputFile(t *testing.T) {
	path := writeTempEnvForRename(t, "DB_HOST=localhost\n")
	outFile := filepath.Join(t.TempDir(), "out.env")

	_, err := executeCommand(rootCmd, "rename", path, "--pair", "DB_HOST=DATABASE_HOST", "--output", outFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if !strings.Contains(string(data), "DATABASE_HOST=localhost") {
		t.Errorf("expected DATABASE_HOST=localhost in output file, got: %q", string(data))
	}
}

func TestRenameCmd_MissingFile(t *testing.T) {
	_, err := executeCommand(rootCmd, "rename", "/nonexistent/.env", "--pair", "A=B")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
