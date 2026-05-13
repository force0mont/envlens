package defaults

import (
	"fmt"
	"sort"
	"strings"
)

// Entry represents a key with an optional default value and whether it was applied.
type Entry struct {
	Key     string
	Default string
	Applied bool
}

// Result holds the merged env map and the list of defaults that were applied.
type Result struct {
	Env     map[string]string
	Entries []Entry
}

// Apply fills in missing or empty keys in env using the provided defaults map.
// Keys already present with non-empty values are left untouched.
func Apply(env map[string]string, defaults map[string]string) Result {
	merged := make(map[string]string, len(env))
	for k, v := range env {
		merged[k] = v
	}

	keys := make([]string, 0, len(defaults))
	for k := range defaults {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var entries []Entry
	for _, k := range keys {
		defVal := defaults[k]
		current, exists := merged[k]
		applied := !exists || current == ""
		if applied {
			merged[k] = defVal
		}
		entries = append(entries, Entry{
			Key:     k,
			Default: defVal,
			Applied: applied,
		})
	}

	return Result{Env: merged, Entries: entries}
}

// Format returns a human-readable summary of which defaults were applied.
func Format(r Result) string {
	if len(r.Entries) == 0 {
		return "no defaults defined\n"
	}

	var sb strings.Builder
	applied := 0
	for _, e := range r.Entries {
		if e.Applied {
			applied++
			sb.WriteString(fmt.Sprintf("  applied  %s=%s\n", e.Key, e.Default))
		} else {
			sb.WriteString(fmt.Sprintf("  skipped  %s (already set)\n", e.Key))
		}
	}

	header := fmt.Sprintf("defaults: %d applied, %d skipped\n", applied, len(r.Entries)-applied)
	return header + sb.String()
}
