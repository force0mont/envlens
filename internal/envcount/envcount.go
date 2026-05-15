package envcount

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envlens/internal/parser"
)

// Result holds the count summary for a single env file.
type Result struct {
	Label    string
	Total    int
	Empty    int
	NonEmpty int
	Prefixes map[string]int // count of keys per top-level prefix
}

// Count analyses a parsed env file and returns a Result.
func Count(label string, entries []parser.Entry) Result {
	r := Result{
		Label:    label,
		Prefixes: make(map[string]int),
	}
	for _, e := range entries {
		r.Total++
		if e.Value == "" {
			r.Empty++
		} else {
			r.NonEmpty++
		}
		if idx := strings.Index(e.Key, "_"); idx > 0 {
			prefix := e.Key[:idx]
			r.Prefixes[prefix]++
		}
	}
	return r
}

// Format renders a slice of Results as a human-readable table.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no files to count\n"
	}
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("[%s]\n", r.Label))
		sb.WriteString(fmt.Sprintf("  total   : %d\n", r.Total))
		sb.WriteString(fmt.Sprintf("  non-empty: %d\n", r.NonEmpty))
		sb.WriteString(fmt.Sprintf("  empty   : %d\n", r.Empty))
		if len(r.Prefixes) > 0 {
			sb.WriteString("  prefixes:\n")
			keys := make([]string, 0, len(r.Prefixes))
			for k := range r.Prefixes {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				sb.WriteString(fmt.Sprintf("    %s: %d\n", k, r.Prefixes[k]))
			}
		}
	}
	return sb.String()
}
