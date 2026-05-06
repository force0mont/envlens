package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Key     string
	Pattern string // optional regex pattern the value must match
	Required bool
}

// Finding represents a validation issue found for a key.
type Finding struct {
	Key     string
	Message string
}

// Validate checks an env map against a set of rules and returns any findings.
func Validate(env map[string]string, rules []Rule) []Finding {
	var findings []Finding

	for _, rule := range rules {
		val, exists := env[rule.Key]

		if rule.Required && !exists {
			findings = append(findings, Finding{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		if rule.Required && strings.TrimSpace(val) == "" {
			findings = append(findings, Finding{
				Key:     rule.Key,
				Message: "required key has empty value",
			})
			continue
		}

		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				findings = append(findings, Finding{
					Key:     rule.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				findings = append(findings, Finding{
					Key:     rule.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
				})
			}
		}
	}

	return findings
}

// Format renders validation findings as a human-readable string.
func Format(findings []Finding, label string) string {
	if len(findings) == 0 {
		return fmt.Sprintf("[%s] all rules passed\n", label)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "[%s] %d validation issue(s):\n", label, len(findings))
	for _, f := range findings {
		fmt.Fprintf(&sb, "  %-30s %s\n", f.Key, f.Message)
	}
	return sb.String()
}
