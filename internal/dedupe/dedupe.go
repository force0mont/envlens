package dedupe

import (
	"fmt"
	"strings"

	"github.com/envlens/internal/parser"
)

// Result holds the outcome of a deduplication operation.
type Result struct {
	Key      string
	Kept     string
	Dupes    []string
	Occurred int
}

// Dedupe removes duplicate keys from an env map, keeping the first occurrence
// by default. If keepLast is true, the last occurrence is kept instead.
func Dedupe(entries []parser.Entry, keepLast bool) ([]parser.Entry, []Result) {
	seen := make(map[string][]parser.Entry)
	order := make([]string, 0)

	for _, e := range entries {
		if _, exists := seen[e.Key]; !exists {
			order = append(order, e.Key)
		}
		seen[e.Key] = append(seen[e.Key], e)
	}

	results := make([]Result, 0)
	deduped := make([]parser.Entry, 0, len(order))

	for _, key := range order {
		group := seen[key]
		if len(group) == 1 {
			deduped = append(deduped, group[0])
			continue
		}

		var kept parser.Entry
		if keepLast {
			kept = group[len(group)-1]
		} else {
			kept = group[0]
		}

		dupeValues := make([]string, 0, len(group)-1)
		for _, e := range group {
			if e.Value != kept.Value {
				dupeValues = append(dupeValues, e.Value)
			}
		}

		results = append(results, Result{
			Key:      key,
			Kept:     kept.Value,
			Dupes:    dupeValues,
			Occurred: len(group),
		})
		deduped = append(deduped, kept)
	}

	return deduped, results
}

// Format renders the deduplication results as a human-readable string.
func Format(results []Result) string {
	if len(results) == 0 {
		return "No duplicate keys found.\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d duplicate key(s):\n", len(results)))
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("  %s: %d occurrences, kept %q\n", r.Key, r.Occurred, r.Kept))
	}
	return sb.String()
}
