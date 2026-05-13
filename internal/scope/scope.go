package scope

import (
	"fmt"
	"sort"
	"strings"
)

// Entry holds a key/value pair with its originating scope label.
type Entry struct {
	Key   string
	Value string
	Scope string
}

// Result is the output of a Scope operation.
type Result struct {
	Entries []Entry
	Scopes  []string
}

// Tag annotates every key in each env map with its scope label.
// Later scopes override earlier ones when keys collide (last-write-wins).
func Tag(envs []map[string]string, labels []string) (Result, error) {
	if len(envs) != len(labels) {
		return Result{}, fmt.Errorf("scope: %d envs but %d labels", len(envs), len(labels))
	}

	// track winning entry per key
	type slot struct {
		entry Entry
		order int
	}
	index := make(map[string]slot)

	for i, env := range envs {
		for k, v := range env {
			index[k] = slot{
				entry: Entry{Key: k, Value: v, Scope: labels[i]},
				order: i,
			}
		}
	}

	keys := make([]string, 0, len(index))
	for k := range index {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, index[k].entry)
	}

	return Result{Entries: entries, Scopes: labels}, nil
}

// Format renders the result as a human-readable string.
func Format(r Result) string {
	if len(r.Entries) == 0 {
		return "no entries\n"
	}
	var sb strings.Builder
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "[%s] %s=%s\n", e.Scope, e.Key, e.Value)
	}
	return sb.String()
}
