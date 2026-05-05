package diff

import (
	"testing"

	"github.com/yourusername/envlens/internal/parser"
)

func makeEnvFile(path string, pairs map[string]string) *parser.EnvFile {
	ef := &parser.EnvFile{Path: path, Index: make(map[string]parser.EnvEntry)}
	for k, v := range pairs {
		entry := parser.EnvEntry{Key: k, Value: v}
		ef.Entries = append(ef.Entries, entry)
		ef.Index[k] = entry
	}
	return ef
}

func TestDiff_Added(t *testing.T) {
	left := makeEnvFile("a.env", map[string]string{"A": "1"})
	right := makeEnvFile("b.env", map[string]string{"A": "1", "B": "2"})

	result := Diff(left, right)
	summary := result.Summary()

	if summary[KindAdded] != 1 {
		t.Errorf("expected 1 added, got %d", summary[KindAdded])
	}
}

func TestDiff_Removed(t *testing.T) {
	left := makeEnvFile("a.env", map[string]string{"A": "1", "B": "2"})
	right := makeEnvFile("b.env", map[string]string{"A": "1"})

	result := Diff(left, right)
	summary := result.Summary()

	if summary[KindRemoved] != 1 {
		t.Errorf("expected 1 removed, got %d", summary[KindRemoved])
	}
}

func TestDiff_Changed(t *testing.T) {
	left := makeEnvFile("a.env", map[string]string{"DB": "localhost"})
	right := makeEnvFile("b.env", map[string]string{"DB": "prod-host"})

	result := Diff(left, right)
	summary := result.Summary()

	if summary[KindChanged] != 1 {
		t.Errorf("expected 1 changed, got %d", summary[KindChanged])
	}
}

func TestDiff_Identical(t *testing.T) {
	left := makeEnvFile("a.env", map[string]string{"X": "same"})
	right := makeEnvFile("b.env", map[string]string{"X": "same"})

	result := Diff(left, right)
	summary := result.Summary()

	if summary[KindIdentical] != 1 {
		t.Errorf("expected 1 identical, got %d", summary[KindIdentical])
	}
}

func TestFormat_ContainsDiffMarkers(t *testing.T) {
	left := makeEnvFile("a.env", map[string]string{"OLD": "x"})
	right := makeEnvFile("b.env", map[string]string{"NEW": "y"})

	result := Diff(left, right)
	output := Format(result)

	if len(output) == 0 {
		t.Error("expected non-empty format output")
	}
	if output[:3] != "---" {
		t.Errorf("expected output to start with ---, got: %q", output[:3])
	}
}
