package sort_test

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/parser"
	"github.com/user/envlens/internal/sort"
)

func makeEnv(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestSort_Ascending(t *testing.T) {
	env := makeEnv("ZEBRA", "1", "APPLE", "2", "MANGO", "3")
	r := sort.Sort(env, sort.Ascending)
	if r.Entries[0].Key != "APPLE" || r.Entries[1].Key != "MANGO" || r.Entries[2].Key != "ZEBRA" {
		t.Errorf("unexpected order: %v", r.Entries)
	}
}

func TestSort_Descending(t *testing.T) {
	env := makeEnv("ZEBRA", "1", "APPLE", "2", "MANGO", "3")
	r := sort.Sort(env, sort.Descending)
	if r.Entries[0].Key != "ZEBRA" || r.Entries[1].Key != "MANGO" || r.Entries[2].Key != "APPLE" {
		t.Errorf("unexpected order: %v", r.Entries)
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	env := makeEnv("ZEBRA", "1", "APPLE", "2")
	origFirst := env[0].Key
	sort.Sort(env, sort.Ascending)
	if env[0].Key != origFirst {
		t.Errorf("original slice was mutated")
	}
}

func TestSort_EmptyEnv(t *testing.T) {
	r := sort.Sort([]parser.Entry{}, sort.Ascending)
	out := sort.Format(r)
	if !strings.Contains(out, "no entries") {
		t.Errorf("expected empty message, got %q", out)
	}
}

func TestSort_CaseInsensitive(t *testing.T) {
	env := makeEnv("zebra", "1", "APPLE", "2", "Mango", "3")
	r := sort.Sort(env, sort.Ascending)
	if r.Entries[0].Key != "APPLE" {
		t.Errorf("expected APPLE first, got %s", r.Entries[0].Key)
	}
}

func TestFormat_ContainsAllKeys(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	r := sort.Sort(env, sort.Ascending)
	out := sort.Format(r)
	if !strings.Contains(out, "BAZ=qux") || !strings.Contains(out, "FOO=bar") {
		t.Errorf("missing keys in output: %q", out)
	}
}
