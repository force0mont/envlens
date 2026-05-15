package envdiff_test

import (
	"strings"
	"testing"

	"github.com/envlens/internal/envdiff"
	"github.com/envlens/internal/parser"
)

func makeEnv(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestCompare_NoChanges(t *testing.T) {
	before := makeEnv("HOST", "localhost", "PORT", "8080")
	after := makeEnv("HOST", "localhost", "PORT", "8080")

	r := envdiff.Compare(before, after)
	if r.Added != 0 || r.Removed != 0 || r.Changed != 0 {
		t.Errorf("expected no changes, got added=%d removed=%d changed=%d", r.Added, r.Removed, r.Changed)
	}
}

func TestCompare_Added(t *testing.T) {
	before := makeEnv("HOST", "localhost")
	after := makeEnv("HOST", "localhost", "PORT", "9090")

	r := envdiff.Compare(before, after)
	if r.Added != 1 {
		t.Fatalf("expected 1 added, got %d", r.Added)
	}
	if r.Changes[0].Key != "PORT" || r.Changes[0].Kind != envdiff.Added {
		t.Errorf("unexpected change: %+v", r.Changes[0])
	}
}

func TestCompare_Removed(t *testing.T) {
	before := makeEnv("HOST", "localhost", "DEBUG", "true")
	after := makeEnv("HOST", "localhost")

	r := envdiff.Compare(before, after)
	if r.Removed != 1 {
		t.Fatalf("expected 1 removed, got %d", r.Removed)
	}
	found := false
	for _, c := range r.Changes {
		if c.Key == "DEBUG" && c.Kind == envdiff.Removed {
			found = true
		}
	}
	if !found {
		t.Error("expected DEBUG to be marked as removed")
	}
}

func TestCompare_Changed(t *testing.T) {
	before := makeEnv("DB_URL", "postgres://old")
	after := makeEnv("DB_URL", "postgres://new")

	r := envdiff.Compare(before, after)
	if r.Changed != 1 {
		t.Fatalf("expected 1 changed, got %d", r.Changed)
	}
	c := r.Changes[0]
	if c.OldValue != "postgres://old" || c.NewValue != "postgres://new" {
		t.Errorf("unexpected values: old=%q new=%q", c.OldValue, c.NewValue)
	}
}

func TestFormat_ContainsSummary(t *testing.T) {
	before := makeEnv("A", "1", "B", "2")
	after := makeEnv("A", "99", "C", "3")

	r := envdiff.Compare(before, after)
	out := envdiff.Format(r, "staging", "production")

	if !strings.Contains(out, "staging → production") {
		t.Error("expected label header in output")
	}
	if !strings.Contains(out, "added: 1") {
		t.Error("expected added count in summary")
	}
	if !strings.Contains(out, "removed: 1") {
		t.Error("expected removed count in summary")
	}
	if !strings.Contains(out, "+ C=3") {
		t.Error("expected added key line")
	}
	if !strings.Contains(out, "- B=2") {
		t.Error("expected removed key line")
	}
}

func TestFormat_DefaultLabels(t *testing.T) {
	r := envdiff.Compare(nil, nil)
	out := envdiff.Format(r, "", "")
	if !strings.Contains(out, "before → after") {
		t.Errorf("expected default labels, got: %s", out)
	}
}
