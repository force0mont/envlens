package flatten

import (
	"fmt"
	"sort"
	"strings"

	"github.com/envlens/envlens/internal/parser"
)

// Result holds a single flattened entry with its origin file label.
type Result struct {
	Key    string
	Value  string
	Source string
}

// Options controls how Flatten behaves when the same key appears in
// multiple files.
type Options struct {
	// Strategy is either "first" (keep first occurrence) or "last" (keep last).
	Strategy string
}

// Flatten merges multiple env files into a single ordered list of entries.
// Duplicate keys are resolved according to opts.Strategy.
// labels must be the same length as envs; if a label is empty the index is used.
func Flatten(envs [][]parser.Entry, labels []string, opts Options) []Result {
	if opts.Strategy == "" {
		opts.Strategy = "first"
	}

	seen := make(map[string]bool)
	order := []string{}
	best := make(map[string]Result)

	for i, entries := range envs {
		label := fmt.Sprintf("file%d", i+1)
		if i < len(labels) && labels[i] != "" {
			label = labels[i]
		}
		for _, e := range entries {
			r := Result{Key: e.Key, Value: e.Value, Source: label}
			if !seen[e.Key] {
				seen[e.Key] = true
				order = append(order, e.Key)
				best[e.Key] = r
			} else if opts.Strategy == "last" {
				best[e.Key] = r
			}
		}
	}

	out := make([]Result, 0, len(order))
	for _, k := range order {
		out = append(out, best[k])
	}
	return out
}

// Format renders flattened results as KEY=VALUE lines with an optional
// source annotation comment above each entry that changed source.
func Format(results []Result, annotate bool) string {
	if len(results) == 0 {
		return "(empty)\n"
	}

	keys := make([]string, len(results))
	for i, r := range results {
		keys[i] = r.Key
	}
	_ = keys // preserve order already determined by Flatten

	var sb strings.Builder
	for _, r := range results {
		if annotate {
			fmt.Fprintf(&sb, "# source: %s\n", r.Source)
		}
		fmt.Fprintf(&sb, "%s=%s\n", r.Key, r.Value)
	}
	return sb.String()
}

// sortedKeys returns a sorted copy of the keys slice (used internally for tests).
func sortedKeys(results []Result) []string {
	ks := make([]string, len(results))
	for i, r := range results {
		ks[i] = r.Key
	}
	sort.Strings(ks)
	return ks
}
