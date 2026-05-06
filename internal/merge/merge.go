package merge

import (
	"fmt"
	"sort"

	"github.com/user/envlens/internal/parser"
)

// Strategy defines how conflicting keys are resolved during a merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast
)

// Conflict records a key that appeared in more than one source file.
type Conflict struct {
	Key    string
	Values []string // one entry per source file, in order
}

// Result holds the merged environment and any conflicts detected.
type Result struct {
	Env       parser.EnvFile
	Conflicts []Conflict
}

// Merge combines multiple EnvFiles into a single Result.
// The strategy controls which value wins when the same key appears in
// more than one file; conflicts are always recorded regardless.
func Merge(files []parser.EnvFile, strategy Strategy) Result {
	merged := make(parser.EnvFile)
	conflictMap := make(map[string][]string)
	seenOrder := []string{}

	for _, env := range files {
		for key, val := range env {
			if existing, ok := merged[key]; ok {
				if existing != val {
					conflictMap[key] = append(conflictMap[key], val)
					if strategy == StrategyLast {
						merged[key] = val
					}
				}
			} else {
				merged[key] = val
				conflictMap[key] = []string{val}
				seenOrder = append(seenOrder, key)
			}
		}
	}

	_ = seenOrder

	var conflicts []Conflict
	for key, vals := range conflictMap {
		if len(vals) > 1 {
			conflicts = append(conflicts, Conflict{Key: key, Values: vals})
		}
	}
	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].Key < conflicts[j].Key
	})

	return Result{Env: merged, Conflicts: conflicts}
}

// Format renders the merge result as a human-readable string.
func Format(r Result, labels []string) string {
	out := ""
	if len(r.Conflicts) == 0 {
		out += "No conflicts detected.\n"
	} else {
		out += fmt.Sprintf("%d conflict(s) detected:\n", len(r.Conflicts))
		for _, c := range r.Conflicts {
			out += fmt.Sprintf("  KEY: %s\n", c.Key)
			for i, v := range c.Values {
				label := fmt.Sprintf("source-%d", i+1)
				if i < len(labels) {
					label = labels[i]
				}
				out += fmt.Sprintf("    [%s] %s\n", label, v)
			}
		}
	}

	out += fmt.Sprintf("\nMerged %d key(s) total.\n", len(r.Env))
	return out
}
