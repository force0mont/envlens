package copy

import (
	"fmt"
	"sort"
)

// Result holds the outcome of copying a key from one env map to another.
type Result struct {
	Key       string
	Value     string
	Overwrote bool
	Missing   bool
}

// Copy copies the specified keys from src into dst.
// If overwrite is false, existing keys in dst are left untouched.
// If a key does not exist in src, a Result with Missing=true is recorded.
func Copy(src, dst map[string]string, keys []string, overwrite bool) (map[string]string, []Result) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var results []Result
	for _, key := range keys {
		val, ok := src[key]
		if !ok {
			results = append(results, Result{Key: key, Missing: true})
			continue
		}
		_, exists := out[key]
		if exists && !overwrite {
			results = append(results, Result{Key: key, Value: val, Overwrote: false})
			continue
		}
		out[key] = val
		results = append(results, Result{Key: key, Value: val, Overwrote: exists})
	}
	return out, results
}

// Format returns a human-readable summary of copy results.
func Format(results []Result, srcLabel, dstLabel string) string {
	if len(results) == 0 {
		return "no keys specified\n"
	}

	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	out := fmt.Sprintf("copying keys from %s → %s\n\n", srcLabel, dstLabel)
	for _, r := range sorted {
		switch {
		case r.Missing:
			out += fmt.Sprintf("  ✗ %-30s (not found in source)\n", r.Key)
		case r.Overwrote:
			out += fmt.Sprintf("  ↺ %-30s (overwritten)\n", r.Key)
		default:
			out += fmt.Sprintf("  ✓ %-30s\n", r.Key)
		}
	}
	return out
}
