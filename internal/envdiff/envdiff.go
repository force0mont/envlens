package envdiff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/envlens/internal/parser"
)

// ChangeKind describes the type of change between two snapshots.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
	Same    ChangeKind = "same"
)

// Change represents a single key-level difference between two env states.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds the full set of changes between two env maps.
type Result struct {
	Changes []Change
	Added   int
	Removed int
	Changed int
}

// Compare produces a Result describing all differences between before and after.
func Compare(before, after []parser.Entry) Result {
	bMap := toMap(before)
	aMap := toMap(after)

	keys := unionKeys(bMap, aMap)
	sort.Strings(keys)

	var result Result
	for _, k := range keys {
		bVal, inB := bMap[k]
		aVal, inA := aMap[k]

		switch {
		case inB && !inA:
			result.Changes = append(result.Changes, Change{Key: k, Kind: Removed, OldValue: bVal})
			result.Removed++
		case !inB && inA:
			result.Changes = append(result.Changes, Change{Key: k, Kind: Added, NewValue: aVal})
			result.Added++
		case bVal != aVal:
			result.Changes = append(result.Changes, Change{Key: k, Kind: Changed, OldValue: bVal, NewValue: aVal})
			result.Changed++
		default:
			result.Changes = append(result.Changes, Change{Key: k, Kind: Same})
		}
	}
	return result
}

// Format renders a Result as a human-readable diff string.
func Format(r Result, beforeLabel, afterLabel string) string {
	if beforeLabel == "" {
		beforeLabel = "before"
	}
	if afterLabel == "" {
		afterLabel = "after"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "envdiff: %s → %s\n", beforeLabel, afterLabel)
	fmt.Fprintf(&sb, "  added: %d  removed: %d  changed: %d\n\n", r.Added, r.Removed, r.Changed)

	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			fmt.Fprintf(&sb, "+ %s=%s\n", c.Key, c.NewValue)
		case Removed:
			fmt.Fprintf(&sb, "- %s=%s\n", c.Key, c.OldValue)
		case Changed:
			fmt.Fprintf(&sb, "~ %s: %q → %q\n", c.Key, c.OldValue, c.NewValue)
		}
	}
	return sb.String()
}

func toMap(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
