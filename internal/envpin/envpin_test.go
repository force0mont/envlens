package envpin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlens/internal/envpin"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestCreate_CopiesEntries(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	p := envpin.Create(env)
	if p.Entries["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", p.Entries["FOO"])
	}
	env["FOO"] = "mutated"
	if p.Entries["FOO"] != "bar" {
		t.Error("Create should deep-copy the env map")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := makeEnv("DB_URL", "postgres://localhost/dev")
	p := envpin.Create(env)
	dir := t.TempDir()
	path := filepath.Join(dir, "pin.json")
	if err := envpin.Save(p, path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := envpin.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Entries["DB_URL"] != "postgres://localhost/dev" {
		t.Errorf("unexpected value: %s", loaded.Entries["DB_URL"])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := envpin.Load("/nonexistent/pin.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDetect_NoDrift(t *testing.T) {
	env := makeEnv("A", "1", "B", "2")
	p := envpin.Create(env)
	drift := envpin.Detect(p, env)
	if len(drift) != 0 {
		t.Errorf("expected no drift, got %d", len(drift))
	}
}

func TestDetect_Changed(t *testing.T) {
	p := envpin.Create(makeEnv("API_KEY", "old"))
	drift := envpin.Detect(p, makeEnv("API_KEY", "new"))
	if len(drift) != 1 || drift[0].Status != "changed" {
		t.Errorf("expected 1 changed entry, got %+v", drift)
	}
}

func TestDetect_Added(t *testing.T) {
	p := envpin.Create(makeEnv("A", "1"))
	drift := envpin.Detect(p, makeEnv("A", "1", "B", "2"))
	if len(drift) != 1 || drift[0].Status != "added" || drift[0].Key != "B" {
		t.Errorf("expected 1 added entry, got %+v", drift)
	}
}

func TestDetect_Removed(t *testing.T) {
	p := envpin.Create(makeEnv("A", "1", "B", "2"))
	drift := envpin.Detect(p, makeEnv("A", "1"))
	if len(drift) != 1 || drift[0].Status != "removed" || drift[0].Key != "B" {
		t.Errorf("expected 1 removed entry, got %+v", drift)
	}
}

func TestFormat_NoDrift(t *testing.T) {
	out := envpin.Format(nil)
	if out != "no drift detected\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_WithDrift(t *testing.T) {
	p := envpin.Create(makeEnv("KEY", "old"))
	drift := envpin.Detect(p, makeEnv("KEY", "new"))
	out := envpin.Format(drift)
	if out == "" {
		t.Error("expected non-empty format output")
	}
	_ = os.Stdout
}
