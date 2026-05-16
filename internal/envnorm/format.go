package envnorm

import (
	"fmt"
	"strings"
)

// Format renders a human-readable summary of normalisation results.
func Format(results []Result) string {
	if len(results) == 0 {
		return "No entries to normalise.\n"
	}

	var sb strings.Builder
	changed := 0

	for _, r := range results {
		if !r.Changed {
			continue
		}
		changed++
		rulesStr := formatRules(r.Rules)
		if r.Original.Key != r.Normalised.Key {
			fmt.Fprintf(&sb, "  key   %q → %q  [%s]\n", r.Original.Key, r.Normalised.Key, rulesStr)
		}
		if r.Original.Value != r.Normalised.Value {
			fmt.Fprintf(&sb, "  value %s: %q → %q  [%s]\n",
				r.Normalised.Key, r.Original.Value, r.Normalised.Value, rulesStr)
		}
	}

	if changed == 0 {
		sb.WriteString("All entries already conform to the requested rules.\n")
	} else {
		fmt.Fprintf(&sb, "\n%d of %d %s normalised.\n",
			changed, len(results), plural(len(results), "entry", "entries"))
	}
	return sb.String()
}

func formatRules(rules []Rule) string {
	parts := make([]string, len(rules))
	for i, r := range rules {
		parts[i] = string(r)
	}
	return strings.Join(parts, ",")
}

func plural(n int, singular, pluralForm string) string {
	if n == 1 {
		return singular
	}
	return pluralForm
}
