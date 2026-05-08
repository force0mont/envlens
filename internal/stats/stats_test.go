package stats_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/stats"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestCompute_BasicCounts(t *testing.T) {
	files := []stats.EnvFile{
		{Label: "prod", Vars: makeEnv("HOST", "example.com", "PORT", "8080", "SECRET", "")},
	}
	result := stats.Compute(files)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	s := result[0]
	if s.Total != 3 {
		t.Errorf("Total: want 3, got %d", s.Total)
	}
	if s.Empty != 1 {
		t.Errorf("Empty: want 1, got %d", s.Empty)
	}
}

func TestCompute_UniqueKeys(t *testing.T) {
	files := []stats.EnvFile{
		{Label: "a", Vars: makeEnv("SHARED", "val", "ONLY_A", "foo")},
		{Label: "b", Vars: makeEnv("SHARED", "val", "ONLY_B", "bar")},
	}
	result := stats.Compute(files)
	byLabel := make(map[string]stats.Stats)
	for _, s := range result {
		byLabel[s.Label] = s
	}
	if byLabel["a"].Unique != 1 {
		t.Errorf("a Unique: want 1, got %d", byLabel["a"].Unique)
	}
	if byLabel["b"].Unique != 1 {
		t.Errorf("b Unique: want 1, got %d", byLabel["b"].Unique)
	}
}

func TestCompute_DuplicatedValues(t *testing.T) {
	files := []stats.EnvFile{
		{Label: "a", Vars: makeEnv("KEY1", "same", "KEY2", "different")},
		{Label: "b", Vars: makeEnv("KEY3", "same", "KEY4", "other")},
	}
	result := stats.Compute(files)
	byLabel := make(map[string]stats.Stats)
	for _, s := range result {
		byLabel[s.Label] = s
	}
	// "same" appears in both files, so KEY1 and KEY3 should be counted as duplicated
	if byLabel["a"].Duplicated != 1 {
		t.Errorf("a Duplicated: want 1, got %d", byLabel["a"].Duplicated)
	}
	if byLabel["b"].Duplicated != 1 {
		t.Errorf("b Duplicated: want 1, got %d", byLabel["b"].Duplicated)
	}
}

func TestCompute_AvgValueLen(t *testing.T) {
	files := []stats.EnvFile{
		{Label: "x", Vars: makeEnv("A", "ab", "B", "abcd")}, // avg = (2+4)/2 = 3.0
	}
	result := stats.Compute(files)
	if result[0].AvgValueLen != 3.0 {
		t.Errorf("AvgValueLen: want 3.0, got %.2f", result[0].AvgValueLen)
	}
}

func TestCompute_EmptyFiles(t *testing.T) {
	result := stats.Compute([]stats.EnvFile{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestFormat_ContainsHeaders(t *testing.T) {
	files := []stats.EnvFile{
		{Label: "dev", Vars: makeEnv("FOO", "bar")},
	}
	out := stats.Format(stats.Compute(files))
	for _, hdr := range []string{"FILE", "TOTAL", "EMPTY", "UNIQUE"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("Format output missing header %q", hdr)
		}
	}
	if !strings.Contains(out, "dev") {
		t.Errorf("Format output missing label 'dev'")
	}
}

func TestFormat_NoFiles(t *testing.T) {
	out := stats.Format(nil)
	if !strings.Contains(out, "no files") {
		t.Errorf("expected 'no files' message, got: %s", out)
	}
}
