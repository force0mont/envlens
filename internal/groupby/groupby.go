package groupby

import (
	"fmt"
	"sort"
	"strings"

	"github.com/envlens/internal/parser"
)

// Group holds all entries that share the same prefix.
type Group struct {
	Prefix  string
	Entries []parser.Entry
}

// Result is the output of a GroupBy operation.
type Result struct {
	Groups   []Group
	Ungrouped []parser.Entry
}

// ByPrefix groups environment entries by their key prefix, where the prefix is
// determined by splitting on sep (e.g. "_") and taking the first N parts.
func ByPrefix(env []parser.Entry, sep string, depth int) Result {
	if sep == "" {
		sep = "_"
	}
	if depth < 1 {
		depth = 1
	}

	buckets := map[string][]parser.Entry{}
	var ungrouped []parser.Entry

	for _, e := range env {
		parts := strings.SplitN(e.Key, sep, depth+1)
		if len(parts) <= depth {
			ungrouped = append(ungrouped, e)
			continue
		}
		prefix := strings.Join(parts[:depth], sep)
		buckets[prefix] = append(buckets[prefix], e)
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	groups := make([]Group, 0, len(keys))
	for _, k := range keys {
		groups = append(groups, Group{Prefix: k, Entries: buckets[k]})
	}

	return Result{Groups: groups, Ungrouped: ungrouped}
}

// Format renders a Result as a human-readable string.
func Format(r Result) string {
	var sb strings.Builder

	for _, g := range r.Groups {
		sb.WriteString(fmt.Sprintf("[%s] (%d keys)\n", g.Prefix, len(g.Entries)))
		for _, e := range g.Entries {
			sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, e.Value))
		}
	}

	if len(r.Ungrouped) > 0 {
		sb.WriteString(fmt.Sprintf("[ungrouped] (%d keys)\n", len(r.Ungrouped)))
		for _, e := range r.Ungrouped {
			sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, e.Value))
		}
	}

	return sb.String()
}
