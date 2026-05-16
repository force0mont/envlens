package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvForReplace(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestEnvreplaceCmd_SubstringReplacement(t *testing.T) {
	p := writeTempEnvForReplace(t, "DB_HOST=localhost:5432\nCACHE_HOST=localhost:6379\n")
	out, err := runRootCmd("envreplace", p, "localhost", "prod.internal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "2 replacements made") {
		t.Errorf("expected summary; got: %s", out)
	}
}

func TestEnvreplaceCmd_LiteralFlag(t *testing.T) {
	p := writeTempEnvForReplace(t, "ENV=staging\nAPP_ENV=staging-eu\n")
	out, err := runRootCmd("envreplace", p, "staging", "production", "--literal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "1 replacement made") {
		t.Errorf("expected 1 replacement; got: %s", out)
	}
}

func TestEnvreplaceCmd_NoMatch(t *testing.T) {
	p := writeTempEnvForReplace(t, "FOO=bar\n")
	out, err := runRootCmd("envreplace", p, "notfound", "x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No replacements made") {
		t.Errorf("expected no-replacement message; got: %s", out)
	}
}

func TestEnvreplaceCmd_WritesOutputFile(t *testing.T) {
	p := writeTempEnvForReplace(t, "HOST=localhost\n")
	outPath := filepath.Join(t.TempDir(), "out.env")
	_, err := runRootCmd("envreplace", p, "localhost", "remote", "--output", outPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	if !strings.Contains(string(data), "HOST=remote") {
		t.Errorf("expected updated value in output file; got: %s", string(data))
	}
}

func TestEnvreplaceCmd_MissingFile(t *testing.T) {
	_, err := runRootCmd("envreplace", "/nonexistent/.env", "a", "b")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
