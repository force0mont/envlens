package setop_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/setop"
)

func makeEnv(pairs ...string) setop.EnvFile {
	m := make(setop.EnvFile)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestIntersect_CommonKeys(t *testing.T) {
	a := makeEnv("A", "1", "B", "2", "C", "3")
	b := makeEnv("B", "20", "C", "30", "D", "40")
	r := setop.Intersect(a, b)
	if len(r.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Keys))
	}
	if r.Keys[0] != "B" || r.Keys[1] != "C" {
		t.Errorf("unexpected keys: %v", r.Keys)
	}
	// values come from first file
	if r.Values["B"] != "2" {
		t.Errorf("expected value '2', got %q", r.Values["B"])
	}
}

func TestIntersect_NoCommonKeys(t *testing.T) {
	a := makeEnv("A", "1")
	b := makeEnv("B", "2")
	r := setop.Intersect(a, b)
	if len(r.Keys) != 0 {
		t.Errorf("expected empty intersection, got %v", r.Keys)
	}
}

func TestUnion_AllKeys(t *testing.T) {
	a := makeEnv("A", "1", "B", "2")
	b := makeEnv("B", "99", "C", "3")
	r := setop.Union(a, b)
	if len(r.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(r.Keys))
	}
	// B value should come from first file
	if r.Values["B"] != "2" {
		t.Errorf("expected B=2, got %q", r.Values["B"])
	}
}

func TestDifference_UniqueToBase(t *testing.T) {
	base := makeEnv("A", "1", "B", "2", "C", "3")
	other := makeEnv("B", "2", "D", "4")
	r := setop.Difference(base, other)
	if len(r.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Keys))
	}
	if r.Keys[0] != "A" || r.Keys[1] != "C" {
		t.Errorf("unexpected keys: %v", r.Keys)
	}
}

func TestDifference_EmptyResult(t *testing.T) {
	base := makeEnv("A", "1")
	other := makeEnv("A", "99")
	r := setop.Difference(base, other)
	if len(r.Keys) != 0 {
		t.Errorf("expected empty difference, got %v", r.Keys)
	}
}

func TestFormat_WithKeys(t *testing.T) {
	r := setop.Result{
		Keys:   []string{"FOO", "BAR"},
		Values: map[string]string{"FOO": "hello", "BAR": "world"},
	}
	out := setop.Format(r)
	if !strings.Contains(out, "FOO=hello") {
		t.Errorf("missing FOO=hello in output: %q", out)
	}
}

func TestFormat_Empty(t *testing.T) {
	r := setop.Result{Values: map[string]string{}}
	out := setop.Format(r)
	if !strings.Contains(out, "(no keys)") {
		t.Errorf("expected '(no keys)', got %q", out)
	}
}
