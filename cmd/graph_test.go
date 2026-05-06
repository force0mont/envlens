package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvForGraph(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func TestGraphCmd_BasicOutput(t *testing.T) {
	file1 := writeTempEnvForGraph(t, "DB_HOST=localhost\nAPP_PORT=8080\n")
	file2 := writeTempEnvForGraph(t, "DB_HOST=db.prod\nSECRET=abc\n")

	out, err := executeCommand(rootCmd, "graph", file1, file2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected shared key DB_HOST in output, got: %q", out)
	}
}

func TestGraphCmd_CustomLabels(t *testing.T) {
	file1 := writeTempEnvForGraph(t, "KEY=1\n")
	file2 := writeTempEnvForGraph(t, "KEY=2\n")

	out, err := executeCommand(rootCmd, "graph", "--labels", "dev,prod", file1, file2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "dev") || !strings.Contains(out, "prod") {
		t.Errorf("expected custom labels in output, got: %q", out)
	}
}

func TestGraphCmd_MismatchedLabels(t *testing.T) {
	file1 := writeTempEnvForGraph(t, "A=1\n")
	file2 := writeTempEnvForGraph(t, "B=2\n")

	_, err := executeCommand(rootCmd, "graph", "--labels", "only-one", file1, file2)
	if err == nil {
		t.Error("expected error for mismatched label count, got nil")
	}
}

func TestGraphCmd_MissingFile(t *testing.T) {
	_, err := executeCommand(rootCmd, "graph", "/nonexistent/.env", "/also/missing/.env")
	if err == nil {
		t.Error("expected error for missing files, got nil")
	}
}
