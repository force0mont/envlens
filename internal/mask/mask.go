package mask

import (
	"fmt"
	"strings"
)

// Result holds the masked version of an env map.
type Result struct {
	Masked map[string]string
	MaskCount int
}

// defaultSensitivePrefixes are key substrings that trigger masking.
var defaultSensitivePrefixes = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "CERT", "SIGNING",
}

// Mask returns a copy of env where sensitive values are partially obscured.
// extraPatterns allows callers to supply additional key substrings to mask.
// visibleChars controls how many leading characters of the value remain visible.
func Mask(env map[string]string, extraPatterns []string, visibleChars int) Result {
	if visibleChars < 0 {
		visibleChars = 0
	}

	patterns := append(defaultSensitivePrefixes, extraPatterns...)
	masked := make(map[string]string, len(env))
	count := 0

	for k, v := range env {
		if isSensitive(k, patterns) && v != "" {
			masked[k] = maskValue(v, visibleChars)
			count++
		} else {
			masked[k] = v
		}
	}

	return Result{Masked: masked, MaskCount: count}
}

// Format renders the masked env map as KEY=VALUE lines, sorted.
func Format(r Result) string {
	keys := make([]string, 0, len(r.Masked))
	for k := range r.Masked {
		keys = append(keys, k)
	}
	sortStrings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, r.Masked[k])
	}
	return sb.String()
}

func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func maskValue(v string, visible int) string {
	if visible >= len(v) {
		return v
	}
	return v[:visible] + strings.Repeat("*", len(v)-visible)
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
