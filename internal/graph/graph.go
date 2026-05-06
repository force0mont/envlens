package graph

import (
	"fmt"
	"sort"
	"strings"

	"github.com/envlens/internal/parser"
)

// Node represents a single .env file in the graph.
type Node struct {
	Label string
	Env   parser.EnvFile
}

// Edge describes a relationship between two nodes and the keys they share.
type Edge struct {
	From    string
	To      string
	Shared  []string
	Unique  map[string][]string // keys unique to each side
}

// Graph holds all nodes and computed edges.
type Graph struct {
	Nodes []Node
	Edges []Edge
}

// Build constructs a Graph from a map of label -> EnvFile.
func Build(files map[string]parser.EnvFile) *Graph {
	g := &Graph{}

	labels := make([]string, 0, len(files))
	for label, env := range files {
		g.Nodes = append(g.Nodes, Node{Label: label, Env: env})
		labels = append(labels, label)
	}
	sort.Strings(labels)

	for i := 0; i < len(labels); i++ {
		for j := i + 1; j < len(labels); j++ {
			a, b := labels[i], labels[j]
			edge := computeEdge(a, files[a], b, files[b])
			g.Edges = append(g.Edges, edge)
		}
	}
	return g
}

func computeEdge(aLabel string, a parser.EnvFile, bLabel string, b parser.EnvFile) Edge {
	shared := []string{}
	uniqueA := []string{}
	uniqueB := []string{}

	for k := range a {
		if _, ok := b[k]; ok {
			shared = append(shared, k)
		} else {
			uniqueA = append(uniqueA, k)
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			uniqueB = append(uniqueB, k)
		}
	}
	sort.Strings(shared)
	sort.Strings(uniqueA)
	sort.Strings(uniqueB)

	return Edge{
		From:   aLabel,
		To:     bLabel,
		Shared: shared,
		Unique: map[string][]string{aLabel: uniqueA, bLabel: uniqueB},
	}
}

// Format renders the graph as a human-readable summary.
func Format(g *Graph) string {
	var sb strings.Builder
	for _, e := range g.Edges {
		sb.WriteString(fmt.Sprintf("[%s] <-> [%s]\n", e.From, e.To))
		sb.WriteString(fmt.Sprintf("  shared (%d): %s\n", len(e.Shared), strings.Join(e.Shared, ", ")))
		for label, keys := range e.Unique {
			if len(keys) > 0 {
				sb.WriteString(fmt.Sprintf("  only in %s (%d): %s\n", label, len(keys), strings.Join(keys, ", ")))
			}
		}
	}
	return sb.String()
}
