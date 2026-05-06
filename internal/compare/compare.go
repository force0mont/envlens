package compare

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envlens/internal/parser"
)

// Result holds the outcome of comparing a key across multiple env files.
type Result struct {
	Key    string
	Values map[string]string // label -> value
	Unique bool              // true if all present values are identical
}

// Compare checks every key that appears in at least one env file and reports
// whether its value is consistent across all files that define it.
func Compare(envs map[string]parser.EnvFile) []Result {
	// collect all keys
	keySet := map[string]struct{}{}
	for _, env := range envs {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	results := make([]Result, 0, len(keys))
	for _, key := range keys {
		values := map[string]string{}
		for label, env := range envs {
			if v, ok := env[key]; ok {
				values[label] = v
			}
		}

		unique := allSame(values)
		results = append(results, Result{
			Key:    key,
			Values: values,
			Unique: unique,
		})
	}
	return results
}

// Format renders a human-readable comparison table.
func Format(results []Result, labels []string) string {
	if len(results) == 0 {
		return "no keys found\n"
	}

	var sb strings.Builder
	header := fmt.Sprintf("%-30s  %s\n", "KEY", strings.Join(labels, "  |  "))
	sb.WriteString(header)
	sb.WriteString(strings.Repeat("-", len(header)) + "\n")

	for _, r := range results {
		parts := make([]string, len(labels))
		for i, lbl := range labels {
			if v, ok := r.Values[lbl]; ok {
				parts[i] = v
			} else {
				parts[i] = "<missing>"
			}
		}
		marker := " "
		if !r.Unique {
			marker = "!"
		}
		sb.WriteString(fmt.Sprintf("%s %-29s  %s\n", marker, r.Key, strings.Join(parts, "  |  ")))
	}
	return sb.String()
}

func allSame(values map[string]string) bool {
	var ref string
	first := true
	for _, v := range values {
		if first {
			ref = v
			first = false
			continue
		}
		if v != ref {
			return false
		}
	}
	return true
}
