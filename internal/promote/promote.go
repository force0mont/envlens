package promote

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the outcome of promoting keys from a source env to a target env.
type Result struct {
	Key      string
	Value    string
	Overwrote bool
}

// Promote copies the given keys from src into dst.
// If a key already exists in dst, it is overwritten and Overwrote is set to true.
// If keys is empty, all keys from src are promoted.
func Promote(src, dst map[string]string, keys []string) (map[string]string, []Result) {
	if len(keys) == 0 {
		keys = sortedKeys(src)
	}

	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var results []Result
	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			continue
		}
		_, existed := out[k]
		out[k] = v
		results = append(results, Result{Key: k, Value: v, Overwrote: existed})
	}
	return out, results
}

// Format renders a human-readable summary of promotion results.
func Format(results []Result, srcLabel, dstLabel string) string {
	if len(results) == 0 {
		return fmt.Sprintf("No keys promoted from %s to %s.\n", srcLabel, dstLabel)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Promoted %d key(s) from %s → %s:\n", len(results), srcLabel, dstLabel))
	for _, r := range results {
		if r.Overwrote {
			sb.WriteString(fmt.Sprintf("  ~ %s (overwritten)\n", r.Key))
		} else {
			sb.WriteString(fmt.Sprintf("  + %s\n", r.Key))
		}
	}
	return sb.String()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
