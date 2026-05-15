package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvForCast(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestEnvcastCmd_ValidIntRule(t *testing.T) {
	file := writeTempEnvForCast(t, "PORT=8080\n")
	out, err := executeCommand(rootCmd, "envcast", file, "--rule", "PORT:int")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "OK") {
		t.Errorf("expected OK in output, got: %s", out)
	}
}

func TestEnvcastCmd_InvalidIntValue(t *testing.T) {
	file := writeTempEnvForCast(t, "PORT=not-a-number\n")
	// command exits 1 on cast failure; RunE returns nil but os.Exit is called
	// We just check output contains ERROR
	out, _ := executeCommand(rootCmd, "envcast", file, "--rule", "PORT:int")
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in output, got: %s", out)
	}
}

func TestEnvcastCmd_BoolRule(t *testing.T) {
	file := writeTempEnvForCast(t, "DEBUG=true\n")
	out, err := executeCommand(rootCmd, "envcast", file, "--rule", "DEBUG:bool")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "true") {
		t.Errorf("expected cast value in output, got: %s", out)
	}
}

func TestEnvcastCmd_InvalidRuleFormat(t *testing.T) {
	file := writeTempEnvForCast(t, "PORT=8080\n")
	_, err := executeCommand(rootCmd, "envcast", file, "--rule", "PORTONLY")
	if err == nil {
		t.Fatal("expected error for malformed rule")
	}
}

func TestEnvcastCmd_UnknownType(t *testing.T) {
	file := writeTempEnvForCast(t, "PORT=8080\n")
	_, err := executeCommand(rootCmd, "envcast", file, "--rule", "PORT:uuid")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestEnvcastCmd_MissingFile(t *testing.T) {
	_, err := executeCommand(rootCmd, "envcast", "/no/such/file.env", "--rule", "PORT:int")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
