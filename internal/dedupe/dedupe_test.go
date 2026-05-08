package dedupe

import (
	"strings"
	"testing"

	"github.com/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestDedupe_NoDuplicates(t *testing.T) {
	entries := makeEntries("FOO", "bar", "BAZ", "qux")
	out, results := Dedupe(entries, false)
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestDedupe_KeepFirst(t *testing.T) {
	entries := makeEntries("FOO", "first", "BAR", "other", "FOO", "second")
	out, results := Dedupe(entries, false)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Kept != "first" {
		t.Errorf("expected kept=first, got %q", results[0].Kept)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 deduped entries, got %d", len(out))
	}
}

func TestDedupe_KeepLast(t *testing.T) {
	entries := makeEntries("FOO", "first", "FOO", "second", "FOO", "third")
	out, results := Dedupe(entries, true)
	if results[0].Kept != "third" {
		t.Errorf("expected kept=third, got %q", results[0].Kept)
	}
	if results[0].Occurred != 3 {
		t.Errorf("expected 3 occurrences, got %d", results[0].Occurred)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestDedupe_PreservesOrder(t *testing.T) {
	entries := makeEntries("B", "1", "A", "2", "C", "3", "A", "4")
	out, _ := Dedupe(entries, false)
	if out[0].Key != "B" || out[1].Key != "A" || out[2].Key != "C" {
		t.Errorf("order not preserved: got %v", out)
	}
}

func TestFormat_NoDuplicates(t *testing.T) {
	out := Format([]Result{})
	if !strings.Contains(out, "No duplicate") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_WithDuplicates(t *testing.T) {
	results := []Result{
		{Key: "SECRET", Kept: "abc", Dupes: []string{"xyz"}, Occurred: 2},
	}
	out := Format(results)
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET in output, got %q", out)
	}
	if !strings.Contains(out, "2 occurrences") {
		t.Errorf("expected occurrence count in output, got %q", out)
	}
}
