package interpolate

import (
	"fmt"
	"regexp"
	"strings"
)

// varPattern matches ${VAR} and $VAR style references.
var varPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Result holds the interpolated environment map and any warnings.
type Result struct {
	Env      map[string]string
	Warnings []string
}

// Interpolate resolves variable references within env values using the same
// map as the source of substitutions. References to undefined variables are
// left as-is and a warning is recorded.
func Interpolate(env map[string]string) Result {
	resolved := make(map[string]string, len(env))
	var warnings []string

	for key, value := range env {
		resolved[key] = expand(value, env, &warnings)
	}

	return Result{Env: resolved, Warnings: warnings}
}

// InterpolateWithOverrides resolves variable references within env values,
// preferring values from overrides when a variable name exists in both maps.
// This is useful for substituting variables using a separate set of
// definitions (e.g. OS environment variables) without merging the maps.
func InterpolateWithOverrides(env map[string]string, overrides map[string]string) Result {
	resolved := make(map[string]string, len(env))
	var warnings []string

	// Build a combined lookup: overrides take precedence over env.
	lookup := make(map[string]string, len(env)+len(overrides))
	for k, v := range env {
		lookup[k] = v
	}
	for k, v := range overrides {
		lookup[k] = v
	}

	for key, value := range env {
		resolved[key] = expand(value, lookup, &warnings)
	}

	return Result{Env: resolved, Warnings: warnings}
}

// expand replaces all variable references in s using lookup.
func expand(s string, lookup map[string]string, warnings *[]string) string {
	return varPattern.ReplaceAllStringFunc(s, func(match string) string {
		name := extractName(match)
		if val, ok := lookup[name]; ok {
			return val
		}
		*warnings = append(*warnings, fmt.Sprintf("undefined variable: %s", name))
		return match
	})
}

// extractName strips the sigil and braces from a matched token.
func extractName(token string) string {
	token = strings.TrimPrefix(token, "$")
	token = strings.TrimPrefix(token, "{")
	token = strings.TrimSuffix(token, "}")
	return token
}
