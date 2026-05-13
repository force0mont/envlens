package typecheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Rule defines an expected type constraint for an environment variable key.
type Rule struct {
	Key  string
	Type string // "int", "bool", "url", "email", "nonempty"
}

// Finding represents a type violation.
type Finding struct {
	Key      string
	Value    string
	Expected string
	Message  string
}

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
var urlRe = regexp.MustCompile(`^https?://[^\s]+$`)

// Check validates the env map against the provided rules.
func Check(env map[string]string, rules []Rule) []Finding {
	var findings []Finding
	for _, rule := range rules {
		val, exists := env[rule.Key]
		if !exists {
			findings = append(findings, Finding{
				Key:      rule.Key,
				Value:    "",
				Expected: rule.Type,
				Message:  "key not found",
			})
			continue
		}
		if msg := validate(val, rule.Type); msg != "" {
			findings = append(findings, Finding{
				Key:      rule.Key,
				Value:    val,
				Expected: rule.Type,
				Message:  msg,
			})
		}
	}
	return findings
}

func validate(val, typ string) string {
	switch strings.ToLower(typ) {
	case "int":
		if _, err := strconv.Atoi(val); err != nil {
			return fmt.Sprintf("%q is not a valid integer", val)
		}
	case "bool":
		lower := strings.ToLower(val)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return fmt.Sprintf("%q is not a valid boolean", val)
		}
	case "url":
		if !urlRe.MatchString(val) {
			return fmt.Sprintf("%q is not a valid URL", val)
		}
	case "email":
		if !emailRe.MatchString(val) {
			return fmt.Sprintf("%q is not a valid email", val)
		}
	case "nonempty":
		if strings.TrimSpace(val) == "" {
			return "value must not be empty"
		}
	}
	return ""
}

// Format renders findings as a human-readable string.
func Format(findings []Finding) string {
	if len(findings) == 0 {
		return "typecheck: all values pass type constraints\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "typecheck: %d violation(s) found\n", len(findings))
	for _, f := range findings {
		fmt.Fprintf(&sb, "  %-24s expected=%-10s %s\n", f.Key, f.Expected, f.Message)
	}
	return sb.String()
}
