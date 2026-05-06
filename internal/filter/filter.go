package filter

import (
	"fmt"
	"regexp"
	"strings"
)

// Result holds the filtered environment variables and metadata.
type Result struct {
	Matched map[string]string
	Skipped int
}

// Options controls how filtering is applied.
type Options struct {
	Prefix  string
	Pattern string
	Keys    []string
}

// Filter returns only the entries from env that match the given options.
// If multiple options are set, entries must satisfy all of them.
func Filter(env map[string]string, opts Options) (Result, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return Result{}, fmt.Errorf("invalid pattern %q: %w", opts.Pattern, err)
		}
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	matched := make(map[string]string)
	for k, v := range env {
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if re != nil && !re.MatchString(k) {
			continue
		}
		if len(keySet) > 0 {
			if _, ok := keySet[k]; !ok {
				continue
			}
		}
		matched[k] = v
	}

	return Result{
		Matched: matched,
		Skipped: len(env) - len(matched),
	}, nil
}

// Format renders the filter result as a human-readable string.
func Format(r Result) string {
	if len(r.Matched) == 0 {
		return fmt.Sprintf("No matching keys found (%d skipped).\n", r.Skipped)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Matched %d key(s) (%d skipped):\n", len(r.Matched), r.Skipped)

	keys := make([]string, 0, len(r.Matched))
	for k := range r.Matched {
		keys = append(keys, k)
	}
	sortStrings(keys)

	for _, k := range keys {
		fmt.Fprintf(&sb, "  %s=%s\n", k, r.Matched[k])
	}
	return sb.String()
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
