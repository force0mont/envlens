package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envlens/internal/snapshot"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestNew_CopiesEnv(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	s := snapshot.New("test", env)

	env["FOO"] = "mutated"
	if s.Env["FOO"] != "bar" {
		t.Errorf("expected snapshot to be independent of original map")
	}
	if s.Label != "test" {
		t.Errorf("expected label 'test', got %q", s.Label)
	}
	if s.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "PORT", "5432")
	s := snapshot.New("production", env)
	s.CreatedAt = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(s, tmp); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(tmp)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Label != "production" {
		t.Errorf("expected label 'production', got %q", loaded.Label)
	}
	if loaded.Env["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", loaded.Env["DB_HOST"])
	}
	if loaded.Env["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", loaded.Env["PORT"])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	s := snapshot.New("x", makeEnv())
	err := snapshot.Save(s, "/nonexistent/dir/snap.json")
	if err == nil {
		t.Error("expected error for invalid path")
	}
	_ = os.Remove("/nonexistent/dir/snap.json")
}
