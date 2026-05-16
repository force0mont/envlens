package envnorm

import (
	"strings"

	"github.com/user/envlens/internal/parser"
)

// Rule describes a single normalisation transformation.
type Rule string

const (
	RuleUppercase   Rule = "uppercase"   // convert keys to UPPER_CASE
	RuleLowercase   Rule = "lowercase"   // convert keys to lower_case
	RuleTrimValues  Rule = "trim_values" // strip leading/trailing whitespace from values
	RuleSnakeCase   Rule = "snake_case"  // replace hyphens in keys with underscores
)

// Result holds an original entry alongside the normalised version.
type Result struct {
	Original  parser.Entry
	Normalised parser.Entry
	Changed   bool
	Rules     []Rule
}

// Normalise applies the given rules to every entry and returns a Result per entry.
func Normalise(entries []parser.Entry, rules []Rule) []Result {
	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		norm, applied := applyRules(e, rules)
		results = append(results, Result{
			Original:   e,
			Normalised: norm,
			Changed:    e.Key != norm.Key || e.Value != norm.Value,
			Rules:      applied,
		})
	}
	return results
}

func applyRules(e parser.Entry, rules []Rule) (parser.Entry, []Rule) {
	applied := []Rule{}
	out := parser.Entry{Key: e.Key, Value: e.Value}
	for _, r := range rules {
		switch r {
		case RuleUppercase:
			out.Key = strings.ToUpper(out.Key)
			applied = append(applied, r)
		case RuleLowercase:
			out.Key = strings.ToLower(out.Key)
			applied = append(applied, r)
		case RuleTrimValues:
			out.Value = strings.TrimSpace(out.Value)
			applied = append(applied, r)
		case RuleSnakeCase:
			out.Key = strings.ReplaceAll(out.Key, "-", "_")
			applied = append(applied, r)
		}
	}
	return out, applied
}

// ToEntries extracts only the normalised entries from a slice of Results.
func ToEntries(results []Result) []parser.Entry {
	out := make([]parser.Entry, len(results))
	for i, r := range results {
		out[i] = r.Normalised
	}
	return out
}
