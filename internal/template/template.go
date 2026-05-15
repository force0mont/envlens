package template

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Entry holds a key and its rendered value after template substitution.
type Entry struct {
	Key      string
	Value    string
	Missing  []string // variable references that could not be resolved
}

// Result is the output of a Render call.
type Result struct {
	Entries  []Entry
	Missing  []string // deduplicated list of all unresolved references
}

var refRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// Render substitutes ${VAR} placeholders in the values of src using the
// variables provided in vars. Keys in src that reference unknown variables are
// recorded in Result.Missing rather than causing an error.
func Render(src map[string]string, vars map[string]string) Result {
	missingSet := map[string]struct{}{}
	var entries []Entry

	keys := sortedKeys(src)
	for _, k := range keys {
		raw := src[k]
		var keyMissing []string

		rendered := refRe.ReplaceAllStringFunc(raw, func(match string) string {
			name := refRe.FindStringSubmatch(match)[1]
			if v, ok := vars[name]; ok {
				return v
			}
			keyMissing = append(keyMissing, name)
			missingSet[name] = struct{}{}
			return match // leave placeholder intact
		})

		entries = append(entries, Entry{
			Key:     k,
			Value:   rendered,
			Missing: keyMissing,
		})
	}

	var missing []string
	for k := range missingSet {
		missing = append(missing, k)
	}
	sort.Strings(missing)

	return Result{Entries: entries, Missing: missing}
}

// Format returns a human-readable summary of the render result.
func Format(r Result) string {
	if len(r.Entries) == 0 {
		return "no entries\n"
	}

	var sb strings.Builder
	for _, e := range r.Entries {
		if len(e.Missing) > 0 {
			sb.WriteString(fmt.Sprintf("%-30s = %s  [unresolved: %s]\n",
				e.Key, e.Value, strings.Join(e.Missing, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("%-30s = %s\n", e.Key, e.Value))
		}
	}

	if len(r.Missing) > 0 {
		sb.WriteString(fmt.Sprintf("\n%d unresolved reference(s): %s\n",
			len(r.Missing), strings.Join(r.Missing, ", ")))
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
