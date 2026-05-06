package rename

import (
	"fmt"
	"strings"
)

// Result holds the outcome of a single rename operation.
type Result struct {
	OldKey  string
	NewKey  string
	Value   string
	Missing bool // true when OldKey was not found in the env map
}

// Rename applies a set of key renames to an env map, returning a new map and
// a slice of Results describing what happened. The original map is not mutated.
func Rename(env map[string]string, pairs map[string]string) (map[string]string, []Result) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	results := make([]Result, 0, len(pairs))
	for oldKey, newKey := range pairs {
		val, ok := out[oldKey]
		if !ok {
			results = append(results, Result{OldKey: oldKey, NewKey: newKey, Missing: true})
			continue
		}
		delete(out, oldKey)
		out[newKey] = val
		results = append(results, Result{OldKey: oldKey, NewKey: newKey, Value: val})
	}
	return out, results
}

// Format returns a human-readable summary of rename results.
func Format(results []Result) string {
	if len(results) == 0 {
		return "No renames performed.\n"
	}

	var sb strings.Builder
	renamed, missing := 0, 0
	for _, r := range results {
		if r.Missing {
			fmt.Fprintf(&sb, "  MISSING  %s (not found)\n", r.OldKey)
			missing++
		} else {
			fmt.Fprintf(&sb, "  RENAMED  %s -> %s\n", r.OldKey, r.NewKey)
			renamed++
		}
	}
	fmt.Fprintf(&sb, "\n%d renamed, %d missing.\n", renamed, missing)
	return sb.String()
}
