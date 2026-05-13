package required

import (
	"fmt"
	"sort"

	"github.com/envlens/envlens/internal/parser"
)

// Finding represents a missing or empty required key.
type Finding struct {
	Key    string
	Reason string
}

// Result holds the outcome of a required-key check.
type Result struct {
	Findings []Finding
	Checked  int
}

// Check verifies that every key in required exists and is non-empty in env.
func Check(env []parser.Entry, required []string) Result {
	lookup := make(map[string]string, len(env))
	for _, e := range env {
		lookup[e.Key] = e.Value
	}

	var findings []Finding
	for _, key := range required {
		val, exists := lookup[key]
		if !exists {
			findings = append(findings, Finding{Key: key, Reason: "missing"})
		} else if val == "" {
			findings = append(findings, Finding{Key: key, Reason: "empty value"})
		}
	}

	return Result{
		Findings: findings,
		Checked:  len(required),
	}
}

// Format renders the Result as a human-readable string.
func Format(r Result) string {
	if len(r.Findings) == 0 {
		return fmt.Sprintf("OK — all %d required key(s) present and non-empty.\n", r.Checked)
	}

	sorted := make([]Finding, len(r.Findings))
	copy(sorted, r.Findings)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	out := fmt.Sprintf("%d of %d required key(s) failed:\n", len(sorted), r.Checked)
	for _, f := range sorted {
		out += fmt.Sprintf("  ✗ %-30s %s\n", f.Key, f.Reason)
	}
	return out
}
