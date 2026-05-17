package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envlens/envlens/cmd"
)

func writeTempEnvForExport(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestEnvexportCmd_DotenvDefault(t *testing.T) {
	file := writeTempEnvForExport(t, "APP=prod\nPORT=9000\n")
	out, err := cmd.ExecuteWithArgs("envexport", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP=prod") {
		t.Errorf("expected APP=prod in output, got: %s", out)
	}
}

func TestEnvexportCmd_ShellFormat(t *testing.T) {
	file := writeTempEnvForExport(t, "TOKEN=abc123\n")
	out, err := cmd.ExecuteWithArgs("envexport", "--format", "shell", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export TOKEN=") {
		t.Errorf("expected shell export, got: %s", out)
	}
}

func TestEnvexportCmd_OmitEmpty(t *testing.T) {
	file := writeTempEnvForExport(t, "A=hello\nB=\nC=world\n")
	out, err := cmd.ExecuteWithArgs("envexport", "--omit-empty", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "B=") {
		t.Errorf("empty key B should be omitted, got: %s", out)
	}
}

func TestEnvexportCmd_WritesOutputFile(t *testing.T) {
	file := writeTempEnvForExport(t, "X=1\n")
	outPath := filepath.Join(t.TempDir(), "out.env")
	_, err := cmd.ExecuteWithArgs("envexport", "--output", outPath, file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if !strings.Contains(string(data), "X=1") {
		t.Errorf("expected X=1 in output file, got: %s", string(data))
	}
}

func TestEnvexportCmd_MissingFile(t *testing.T) {
	_, err := cmd.ExecuteWithArgs("envexport", "/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
