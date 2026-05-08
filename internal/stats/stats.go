package stats

import (
	"fmt"
	"sort"
	"strings"
)

// EnvFile represents a named environment variable set.
type EnvFile struct {
	Label string
	Vars  map[string]string
}

// Stats holds summary statistics for an env file.
type Stats struct {
	Label       string
	Total       int
	Empty       int
	Unique      int
	Duplicated  int
	AvgValueLen float64
}

// Compute calculates statistics for each provided EnvFile.
// A key is considered duplicated if its value appears in more than one file.
func Compute(files []EnvFile) []Stats {
	// Build a map of value -> count across all files to detect duplicates.
	valueCount := make(map[string]int)
	for _, f := range files {
		seen := make(map[string]bool)
		for _, v := range f.Vars {
			if v != "" && !seen[v] {
				valueCount[v]++
				seen[v] = true
			}
		}
	}

	// Build a set of all keys across all files.
	allKeys := make(map[string]int) // key -> number of files containing it
	for _, f := range files {
		for k := range f.Vars {
			allKeys[k]++
		}
	}

	result := make([]Stats, 0, len(files))
	for _, f := range files {
		s := Stats{Label: f.Label, Total: len(f.Vars)}
		totalLen := 0
		for k, v := range f.Vars {
			if v == "" {
				s.Empty++
			} else {
				totalLen += len(v)
				if valueCount[v] > 1 {
					s.Duplicated++
				}
			}
			if allKeys[k] == 1 {
				s.Unique++
			}
		}
		nonEmpty := s.Total - s.Empty
		if nonEmpty > 0 {
			s.AvgValueLen = float64(totalLen) / float64(nonEmpty)
		}
		result = append(result, s)
	}
	return result
}

// Format renders statistics as a human-readable table.
func Format(stats []Stats) string {
	if len(stats) == 0 {
		return "no files to report\n"
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Label < stats[j].Label
	})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-20s %6s %6s %7s %11s %12s\n",
		"FILE", "TOTAL", "EMPTY", "UNIQUE", "DUPLICATED", "AVG_VAL_LEN"))
	sb.WriteString(strings.Repeat("-", 68) + "\n")
	for _, s := range stats {
		sb.WriteString(fmt.Sprintf("%-20s %6d %6d %7d %11d %12.1f\n",
			s.Label, s.Total, s.Empty, s.Unique, s.Duplicated, s.AvgValueLen))
	}
	return sb.String()
}
