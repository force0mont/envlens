package trim

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the outcome of trimming a single key.
type Result struct {
	Key      string
	Original string
	Trimmed  string
	Changed  bool
}

// Trim removes leading and trailing whitespace from all values in the env map.
// It returns a new map and a slice of Results describing what changed.
func Trim(env map[string]string) (map[string]string, []Result) {
	out := make(map[string]string, len(env))
	results := make([]Result, 0)

	for k, v := range env {
		trimmed := strings.TrimSpace(v)
		out[k] = trimmed
		results = append(results, Result{
			Key:      k,
			Original: v,
			Trimmed:  trimmed,
			Changed:  trimmed != v,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return out, results
}

// Format renders the trim results as a human-readable string.
func Format(results []Result) string {
	var sb strings.Builder

	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
		}
	}

	if changed == 0 {
		sb.WriteString("No values required trimming.\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("%d value(s) trimmed:\n", changed))
	for _, r := range results {
		if r.Changed {
			sb.WriteString(fmt.Sprintf("  %-24s %q -> %q\n", r.Key, r.Original, r.Trimmed))
		}
	}

	return sb.String()
}
