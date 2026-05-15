package envsearch

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/envlens/internal/parser"
)

// MatchField controls which part of an entry to search.
type MatchField string

const (
	MatchKey   MatchField = "key"
	MatchValue MatchField = "value"
	MatchBoth  MatchField = "both"
)

// Result holds a single search hit.
type Result struct {
	Key     string
	Value   string
	Matched MatchField
}

// Options configures a Search call.
type Options struct {
	Pattern     string
	Field       MatchField
	CaseSensitive bool
}

// Search scans env entries for keys or values matching the given pattern.
func Search(entries []parser.Entry, opts Options) ([]Result, error) {
	if opts.Field == "" {
		opts.Field = MatchBoth
	}

	pattern := opts.Pattern
	if !opts.CaseSensitive {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern %q: %w", opts.Pattern, err)
	}

	var results []Result
	for _, e := range entries {
		var matched MatchField
		switch opts.Field {
		case MatchKey:
			if re.MatchString(e.Key) {
				matched = MatchKey
			}
		case MatchValue:
			if re.MatchString(e.Value) {
				matched = MatchValue
			}
		default: // MatchBoth
			if re.MatchString(e.Key) {
				matched = MatchKey
			} else if re.MatchString(e.Value) {
				matched = MatchValue
			}
		}
		if matched != "" {
			results = append(results, Result{Key: e.Key, Value: e.Value, Matched: matched})
		}
	}
	return results, nil
}

// Format renders search results as a human-readable string.
func Format(results []Result, pattern string) string {
	if len(results) == 0 {
		return fmt.Sprintf("No matches found for %q.\n", pattern)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d match(es) for %q:\n", len(results), pattern)
	for _, r := range results {
		fmt.Fprintf(&sb, "  [%s] %s=%s\n", r.Matched, r.Key, r.Value)
	}
	return sb.String()
}
