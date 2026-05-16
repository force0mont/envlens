package envreplace

import (
	"fmt"
	"strings"

	"github.com/envlens/internal/parser"
)

// Replacement describes a find/replace operation on env values.
type Replacement struct {
	Key     string
	OldVal  string
	NewVal  string
	Changed bool
}

// Result holds the updated entries and the list of replacements made.
type Result struct {
	Entries      []parser.Entry
	Replacements []Replacement
}

// Replace performs a find-and-replace on values across all entries.
// If keys is non-empty, only those keys are considered.
// If literal is false, old is treated as a substring match.
func Replace(entries []parser.Entry, old, new string, keys []string, literal bool) Result {
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}

	updated := make([]parser.Entry, len(entries))
	var replacements []Replacement

	for i, e := range entries {
		updated[i] = e
		if len(keySet) > 0 && !keySet[e.Key] {
			continue
		}
		var newVal string
		if literal {
			if e.Value == old {
				newVal = new
			} else {
				continue
			}
		} else {
			if !strings.Contains(e.Value, old) {
				continue
			}
			newVal = strings.ReplaceAll(e.Value, old, new)
		}
		replacements = append(replacements, Replacement{
			Key:     e.Key,
			OldVal:  e.Value,
			NewVal:  newVal,
			Changed: true,
		})
		updated[i].Value = newVal
	}

	return Result{Entries: updated, Replacements: replacements}
}

// Format returns a human-readable summary of the replacement result.
func Format(r Result) string {
	if len(r.Replacements) == 0 {
		return "No replacements made.\n"
	}
	var sb strings.Builder
	for _, rep := range r.Replacements {
		fmt.Fprintf(&sb, "  %s: %q -> %q\n", rep.Key, rep.OldVal, rep.NewVal)
	}
	plural := "replacement"
	if len(r.Replacements) != 1 {
		plural = "replacements"
	}
	fmt.Fprintf(&sb, "\n%d %s made.\n", len(r.Replacements), plural)
	return sb.String()
}
