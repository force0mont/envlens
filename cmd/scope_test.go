package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvForScope(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnvForScope: %v", err)
	}
	return p
}

func TestScopeCmd_BasicOutput(t *testing.T) {
	f1 := writeTempEnvForScope(t, "APP_ENV=development\nPORT=3000\n")
	f2 := writeTempEnvForScope(t, "APP_ENV=production\nPORT=8080\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"scope", "--labels", "dev,prod", f1, f2})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[prod]") {
		t.Errorf("expected [prod] in output, got: %s", out)
	}
}

func TestScopeCmd_DefaultLabels(t *testing.T) {
	f1 := writeTempEnvForScope(t, "KEY=val\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"scope", f1})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "KEY=val") {
		t.Errorf("expected KEY=val in output, got: %s", out)
	}
}

func TestScopeCmd_MissingFile(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"scope", "--labels", "x", "/nonexistent/.env"})

	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}
