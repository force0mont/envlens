package envclone

import (
	"fmt"
	"strings"

	"github.com/envlens/internal/parser"
)

// Result holds the outcome of a clone operation for a single key.
type Result struct {
	SourceKey string
	DestKey   string
	Value     string
	Skipped   bool // true when dest key already exists and overwrite is false
}

// Clone duplicates values from source keys to destination keys within an env.
// pairs maps srcKey -> destKey. If overwrite is false, existing dest keys are
// left untouched and the result is marked Skipped.
func Clone(entries []parser.Entry, pairs map[string]string, overwrite bool) ([]parser.Entry, []Result) {
	index := make(map[string]int, len(entries))
	for i, e := range entries {
		index[e.Key] = i
	}

	out := make([]parser.Entry, len(entries))
	copy(out, entries)

	var results []Result

	for src, dst := range pairs {
		srcIdx, srcFound := index[src]
		if !srcFound {
			results = append(results, Result{SourceKey: src, DestKey: dst, Skipped: true})
			continue
		}
		val := out[srcIdx].Value

		if dstIdx, exists := index[dst]; exists {
			if !overwrite {
				results = append(results, Result{SourceKey: src, DestKey: dst, Value: val, Skipped: true})
				continue
			}
			out[dstIdx].Value = val
			results = append(results, Result{SourceKey: src, DestKey: dst, Value: val})
		} else {
			out = append(out, parser.Entry{Key: dst, Value: val})
			index[dst] = len(out) - 1
			results = append(results, Result{SourceKey: src, DestKey: dst, Value: val})
		}
	}

	return out, results
}

// Format renders a human-readable summary of clone results.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no keys cloned\n"
	}
	var sb strings.Builder
	cloned, skipped := 0, 0
	for _, r := range results {
		if r.Skipped {
			skipped++
			sb.WriteString(fmt.Sprintf("  SKIP  %s -> %s\n", r.SourceKey, r.DestKey))
		} else {
			cloned++
			sb.WriteString(fmt.Sprintf("  OK    %s -> %s = %q\n", r.SourceKey, r.DestKey, r.Value))
		}
	}
	sb.WriteString(fmt.Sprintf("\n%d cloned, %d skipped\n", cloned, skipped))
	return sb.String()
}
