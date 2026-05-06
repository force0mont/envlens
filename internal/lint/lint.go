package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// Severity represents the level of a lint finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Finding represents a single lint issue found in an env file.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

var (
	keyPattern     = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	doubleUnder    = regexp.MustCompile(`__`)
	spaceInValue   = regexp.MustCompile(`\s{2,}`)
)

// Lint analyses a parsed env map and returns a list of findings.
func Lint(env map[string]string) []Finding {
	var findings []Finding

	for k, v := range env {
		// Keys should be uppercase with underscores only.
		if !keyPattern.MatchString(k) {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "key should be UPPER_SNAKE_CASE and start with a letter",
				Severity: SeverityError,
			})
		}

		// Warn on consecutive underscores.
		if doubleUnder.MatchString(k) {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "key contains consecutive underscores",
				Severity: SeverityWarning,
			})
		}

		// Warn on values with leading/trailing whitespace.
		if v != strings.TrimSpace(v) {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "value has leading or trailing whitespace",
				Severity: SeverityWarning,
			})
		}

		// Info on values with multiple consecutive spaces.
		if spaceInValue.MatchString(v) {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "value contains multiple consecutive spaces",
				Severity: SeverityInfo,
			})
		}
	}

	return findings
}

// Format renders lint findings as a human-readable string.
func Format(findings []Finding) string {
	if len(findings) == 0 {
		return "No lint issues found.\n"
	}
	var sb strings.Builder
	for _, f := range findings {
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", strings.ToUpper(string(f.Severity)), f.Key, f.Message))
	}
	sb.WriteString(fmt.Sprintf("\n%d issue(s) found.\n", len(findings)))
	return sb.String()
}
