package extract

import (
	"fmt"
	"sort"
	"strings"
)

// EnvFile represents a parsed environment file as a key-value map.
type EnvFile = map[string]string

// Result holds the extracted subset of an environment file.
type Result struct {
	Key   string
	Value string
	Found bool
}

// Extract returns the values for the given keys from the env map.
// Keys not found in the map are included with Found=false.
func Extract(env EnvFile, keys []string) []Result {
	results := make([]Result, 0, len(keys))
	for _, k := range keys {
		v, ok := env[k]
		results = append(results, Result{Key: k, Value: v, Found: ok})
	}
	return results
}

// Format renders extracted results as a human-readable string.
// Missing keys are clearly flagged. A summary line is appended.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no keys requested\n"
	}

	// stable output: sort by key
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	var sb strings.Builder
	missing := 0
	for _, r := range sorted {
		if r.Found {
			fmt.Fprintf(&sb, "  %s=%s\n", r.Key, r.Value)
		} else {
			fmt.Fprintf(&sb, "  %s=<missing>\n", r.Key)
			missing++
		}
	}

	found := len(sorted) - missing
	fmt.Fprintf(&sb, "\n%d extracted, %d missing\n", found, missing)
	return sb.String()
}
