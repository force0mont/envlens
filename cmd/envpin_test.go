package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/envpin"
)

func writeTempEnvForPin(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestEnvpinSave_CreatesFile(t *testing.T) {
	envFile := writeTempEnvForPin(t, "FOO=bar\nBAZ=qux\n")
	pinPath := filepath.Join(t.TempDir(), "pin.json")

	out, err := executeCommand(rootCmd, "envpin", "save", envFile, "--output", pinPath)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if _, err := os.Stat(pinPath); os.IsNotExist(err) {
		t.Error("pin file was not created")
	}
}

func TestEnvpinSave_ValidJSON(t *testing.T) {
	envFile := writeTempEnvForPin(t, "KEY=value\n")
	pinPath := filepath.Join(t.TempDir(), "pin.json")

	_, err := executeCommand(rootCmd, "envpin", "save", envFile, "--output", pinPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(pinPath)
	var p envpin.Pin
	if err := json.Unmarshal(data, &p); err != nil {
		t.Errorf("pin file is not valid JSON: %v", err)
	}
	if p.Entries["KEY"] != "value" {
		t.Errorf("expected KEY=value in pin, got %v", p.Entries)
	}
}

func TestEnvpinCheck_NoDrift(t *testing.T) {
	envFile := writeTempEnvForPin(t, "FOO=bar\n")
	pinPath := filepath.Join(t.TempDir(), "pin.json")

	_, err := executeCommand(rootCmd, "envpin", "save", envFile, "--output", pinPath)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}
	out, err := executeCommand(rootCmd, "envpin", "check", envFile, "--pin", pinPath)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if !strings.Contains(out, "no drift") {
		t.Errorf("expected no drift message, got: %s", out)
	}
}

func TestEnvpinCheck_MissingPin(t *testing.T) {
	envFile := writeTempEnvForPin(t, "FOO=bar\n")
	_, err := executeCommand(rootCmd, "envpin", "check", envFile, "--pin", "/nonexistent/pin.json")
	if err == nil {
		t.Error("expected error for missing pin file")
	}
}
