package graph

import (
	"strings"
	"testing"

	"github.com/envlens/internal/parser"
)

func makeEnv(pairs ...string) parser.EnvFile {
	env := parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		env[pairs[i]] = pairs[i+1]
	}
	return env
}

func TestBuild_SharedKeys(t *testing.T) {
	files := map[string]parser.EnvFile{
		"dev":  makeEnv("DB_HOST", "localhost", "APP_PORT", "8080"),
		"prod": makeEnv("DB_HOST", "db.prod", "APP_PORT", "443", "SECRET", "xyz"),
	}
	g := Build(files)
	if len(g.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(g.Edges))
	}
	e := g.Edges[0]
	if len(e.Shared) != 2 {
		t.Errorf("expected 2 shared keys, got %d: %v", len(e.Shared), e.Shared)
	}
}

func TestBuild_UniqueKeys(t *testing.T) {
	files := map[string]parser.EnvFile{
		"a": makeEnv("ONLY_A", "1", "COMMON", "x"),
		"b": makeEnv("ONLY_B", "2", "COMMON", "y"),
	}
	g := Build(files)
	e := g.Edges[0]
	if len(e.Unique["a"]) != 1 || e.Unique["a"][0] != "ONLY_A" {
		t.Errorf("unexpected unique keys for a: %v", e.Unique["a"])
	}
	if len(e.Unique["b"]) != 1 || e.Unique["b"][0] != "ONLY_B" {
		t.Errorf("unexpected unique keys for b: %v", e.Unique["b"])
	}
}

func TestBuild_MultipleEdges(t *testing.T) {
	files := map[string]parser.EnvFile{
		"dev":     makeEnv("A", "1"),
		"staging": makeEnv("A", "2", "B", "3"),
		"prod":    makeEnv("B", "4", "C", "5"),
	}
	g := Build(files)
	// 3 files => 3 edges
	if len(g.Edges) != 3 {
		t.Errorf("expected 3 edges, got %d", len(g.Edges))
	}
}

func TestBuild_SingleFile(t *testing.T) {
	// A single file produces no edges since there are no pairs to compare.
	files := map[string]parser.EnvFile{
		"only": makeEnv("KEY", "value"),
	}
	g := Build(files)
	if len(g.Edges) != 0 {
		t.Errorf("expected 0 edges for single file, got %d", len(g.Edges))
	}
}

func TestFormat_ContainsLabels(t *testing.T) {
	files := map[string]parser.EnvFile{
		"alpha": makeEnv("KEY", "1"),
		"beta":  makeEnv("KEY", "2", "EXTRA", "3"),
	}
	g := Build(files)
	out := Format(g)
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "beta") {
		t.Errorf("format output missing labels: %q", out)
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("format output missing shared key: %q", out)
	}
}

func TestFormat_EmptyGraph(t *testing.T) {
	g := &Graph{}
	out := Format(g)
	if out != "" {
		t.Errorf("expected empty output for empty graph, got %q", out)
	}
}
