package diff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yourusername/envlens/internal/parser"
)

// DiffKind classifies the type of difference between two env files.
type DiffKind string

const (
	KindAdded    DiffKind = "added"    // present in right, missing in left
	KindRemoved  DiffKind = "removed"  // present in left, missing in right
	KindChanged  DiffKind = "changed"  // present in both, different value
	KindIdentical DiffKind = "identical" // present in both, same value
)

// DiffEntry represents a single key comparison result.
type DiffEntry struct {
	Key        string
	Kind       DiffKind
	LeftValue  string
	RightValue string
}

// Result holds the full diff between two env files.
type Result struct {
	Left    string
	Right   string
	Entries []DiffEntry
}

// Summary returns counts of each diff kind.
func (r *Result) Summary() map[DiffKind]int {
	counts := map[DiffKind]int{}
	for _, e := range r.Entries {
		counts[e.Kind]++
	}
	return counts
}

// HasChanges reports whether the result contains any non-identical entries.
func (r *Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Kind != KindIdentical {
			return true
		}
	}
	return false
}

// Diff compares two parsed EnvFiles and returns a Result.
func Diff(left, right *parser.EnvFile) *Result {
	result := &Result{Left: left.Path, Right: right.Path}

	keys := allKeys(left, right)
	for _, key := range keys {
		lEntry, lOk := left.Index[key]
		rEntry, rOk := right.Index[key]

		var entry DiffEntry
		entry.Key = key

		switch {
		case lOk && !rOk:
			entry.Kind = KindRemoved
			entry.LeftValue = lEntry.Value
		case !lOk && rOk:
			entry.Kind = KindAdded
			entry.RightValue = rEntry.Value
		case lEntry.Value != rEntry.Value:
			entry.Kind = KindChanged
			entry.LeftValue = lEntry.Value
			entry.RightValue = rEntry.Value
		default:
			entry.Kind = KindIdentical
			entry.LeftValue = lEntry.Value
			entry.RightValue = rEntry.Value
		}

		result.Entries = append(result.Entries, entry)
	}

	return result
}

// Format returns a human-readable diff string.
func Format(r *Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "--- %s\n+++ %s\n\n", r.Left, r.Right)
	for _, e := range r.Entries {
		switch e.Kind {
		case KindAdded:
			fmt.Fprintf(&sb, "+ %s=%s\n", e.Key, e.RightValue)
		case KindRemoved:
			fmt.Fprintf(&sb, "- %s=%s\n", e.Key, e.LeftValue)
		case KindChanged:
			fmt.Fprintf(&sb, "~ %s: %q -> %q\n", e.Key, e.LeftValue, e.RightValue)
		}
	}
	return sb.String()
}

func allKeys(left, right *parser.EnvFile) []string {
	seen := map[string]struct{}{}
	for k := range left.Index {
		seen[k] = struct{}{}
	}
	for k := range right.Index {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
