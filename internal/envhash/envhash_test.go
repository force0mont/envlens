package envhash

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/parser"
)

func makeEnv(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestHash_BasicResult(t *testing.T) {
	env := makeEnv("HOST", "localhost", "PORT", "5432")
	r := Hash("test", env)
	if r.Hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if r.KeyCount != 2 {
		t.Fatalf("expected 2 keys, got %d", r.KeyCount)
	}
	if r.Label != "test" {
		t.Fatalf("expected label 'test', got %s", r.Label)
	}
}

func TestHash_DeterministicAcrossInsertionOrder(t *testing.T) {
	env1 := makeEnv("A", "1", "B", "2")
	env2 := makeEnv("B", "2", "A", "1")
	r1 := Hash("e1", env1)
	r2 := Hash("e2", env2)
	if r1.Hash != r2.Hash {
		t.Fatalf("hashes differ despite same content: %s vs %s", r1.Hash, r2.Hash)
	}
}

func TestHash_DiffersWhenValueChanges(t *testing.T) {
	env1 := makeEnv("HOST", "prod.example.com")
	env2 := makeEnv("HOST", "staging.example.com")
	r1 := Hash("prod", env1)
	r2 := Hash("staging", env2)
	if r1.Hash == r2.Hash {
		t.Fatal("expected hashes to differ for different values")
	}
}

func TestHash_EmptyEnv(t *testing.T) {
	r := Hash("empty", []parser.Entry{})
	if r.Hash == "" {
		t.Fatal("expected non-empty hash even for empty env")
	}
	if r.KeyCount != 0 {
		t.Fatalf("expected 0 keys, got %d", r.KeyCount)
	}
}

func TestFormat_AllSame(t *testing.T) {
	env := makeEnv("X", "1")
	r1 := Hash("a", env)
	r2 := Hash("b", env)
	out := Format([]Result{r1, r2})
	if !strings.Contains(out, "all files produce the same hash") {
		t.Errorf("expected same-hash message, got:\n%s", out)
	}
}

func TestFormat_Differ(t *testing.T) {
	r1 := Hash("a", makeEnv("K", "1"))
	r2 := Hash("b", makeEnv("K", "2"))
	out := Format([]Result{r1, r2})
	if !strings.Contains(out, "files differ") {
		t.Errorf("expected differ message, got:\n%s", out)
	}
}

func TestFormat_NoFiles(t *testing.T) {
	out := Format(nil)
	if !strings.Contains(out, "no files to hash") {
		t.Errorf("expected no-files message, got: %s", out)
	}
}
