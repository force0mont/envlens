package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/snapshot"
)

func writeEnvForSnapshot(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestSnapshotSave_CreatesFile(t *testing.T) {
	envFile := writeEnvForSnapshot(t, "APP_ENV=production\nPORT=8080\n")
	outFile := filepath.Join(t.TempDir(), "snap.json")

	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	s := snapshot.New("test-label", env)
	if err := snapshot.Save(s, outFile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if !strings.Contains(string(data), "test-label") {
		t.Errorf("expected label in output, got: %s", string(data))
	}
	_ = envFile
}

func TestSnapshotLoad_ParsesJSON(t *testing.T) {
	tmp := t.TempDir()
	env := map[string]string{"KEY": "value", "OTHER": "123"}
	s := snapshot.New("my-snap", env)
	path := filepath.Join(tmp, "s.json")
	if err := snapshot.Save(s, path); err != nil {
		t.Fatal(err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if loaded.Label != "my-snap" {
		t.Errorf("wrong label: %s", loaded.Label)
	}
	if loaded.Env["KEY"] != "value" {
		t.Errorf("wrong KEY value: %s", loaded.Env["KEY"])
	}
}

func TestSnapshotSave_ValidJSON(t *testing.T) {
	env := map[string]string{"X": "1"}
	s := snapshot.New("valid", env)
	path := filepath.Join(t.TempDir(), "out.json")
	if err := snapshot.Save(s, path); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}
}
