package flatten

import (
	"strings"
	"testing"

	"github.com/envlens/envlens/internal/parser"
)

func makeEnv(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestFlatten_NoOverlap(t *testing.T) {
	a := makeEnv("HOST", "localhost", "PORT", "5432")
	b := makeEnv("DEBUG", "true")
	res := Flatten([][]parser.Entry{a, b}, []string{"a", "b"}, Options{})
	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}
}

func TestFlatten_StrategyFirst(t *testing.T) {
	a := makeEnv("KEY", "from-a")
	b := makeEnv("KEY", "from-b")
	res := Flatten([][]parser.Entry{a, b}, []string{"a", "b"}, Options{Strategy: "first"})
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Value != "from-a" {
		t.Errorf("expected from-a, got %s", res[0].Value)
	}
	if res[0].Source != "a" {
		t.Errorf("expected source a, got %s", res[0].Source)
	}
}

func TestFlatten_StrategyLast(t *testing.T) {
	a := makeEnv("KEY", "from-a")
	b := makeEnv("KEY", "from-b")
	res := Flatten([][]parser.Entry{a, b}, []string{"a", "b"}, Options{Strategy: "last"})
	if res[0].Value != "from-b" {
		t.Errorf("expected from-b, got %s", res[0].Value)
	}
	if res[0].Source != "b" {
		t.Errorf("expected source b, got %s", res[0].Source)
	}
}

func TestFlatten_DefaultLabelFallback(t *testing.T) {
	a := makeEnv("X", "1")
	res := Flatten([][]parser.Entry{a}, []string{}, Options{})
	if res[0].Source != "file1" {
		t.Errorf("expected default label file1, got %s", res[0].Source)
	}
}

func TestFlatten_PreservesInsertionOrder(t *testing.T) {
	a := makeEnv("B", "2", "A", "1")
	res := Flatten([][]parser.Entry{a}, []string{"src"}, Options{})
	if res[0].Key != "B" || res[1].Key != "A" {
		t.Errorf("order not preserved: got %s, %s", res[0].Key, res[1].Key)
	}
}

func TestFormat_Annotated(t *testing.T) {
	res := []Result{{Key: "FOO", Value: "bar", Source: "prod"}}
	out := Format(res, true)
	if !strings.Contains(out, "# source: prod") {
		t.Errorf("expected annotation, got: %s", out)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output")
	}
}

func TestFormat_NoAnnotation(t *testing.T) {
	res := []Result{{Key: "FOO", Value: "bar", Source: "prod"}}
	out := Format(res, false)
	if strings.Contains(out, "#") {
		t.Errorf("unexpected annotation in output: %s", out)
	}
}

func TestFormat_Empty(t *testing.T) {
	out := Format(nil, false)
	if out != "(empty)\n" {
		t.Errorf("expected (empty), got %s", out)
	}
}
