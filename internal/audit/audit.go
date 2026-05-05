package audit

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envlens/internal/parser"
)

// Severity represents the importance level of an audit finding.
type Severity string

const (
	SeverityInfo    Severity = "INFO"
	SeverityWarning Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"
)

// Finding represents a single audit result for a key/value pair.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

// Report holds all findings from an audit run.
type Report struct {
	File     string
	Findings []Finding
}

// HasIssues returns true if the report contains any non-info findings.
func (r *Report) HasIssues() bool {
	for _, f := range r.Findings {
		if f.Severity != SeverityInfo {
			return true
		}
	}
	return false
}

// Audit inspects an EnvFile and returns a Report of findings.
func Audit(filename string, env parser.EnvFile) Report {
	report := Report{File: filename}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := env[key]

		if val == "" {
			report.Findings = append(report.Findings, Finding{
				Key:      key,
				Message:  "value is empty",
				Severity: SeverityWarning,
			})
		}

		if looksLikeSecret(key) && val != "" {
			report.Findings = append(report.Findings, Finding{
				Key:      key,
				Message:  fmt.Sprintf("potential secret detected (key matches pattern): %q", key),
				Severity: SeverityCritical,
			})
		}

		if strings.Contains(val, "localhost") || strings.Contains(val, "127.0.0.1") {
			report.Findings = append(report.Findings, Finding{
				Key:      key,
				Message:  "value references localhost — may not be suitable for production",
				Severity: SeverityWarning,
			})
		}
	}

	return report
}

func looksLikeSecret(key string) bool {
	upper := strings.ToUpper(key)
	secretPatterns := []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE_KEY", "CREDENTIAL"}
	for _, p := range secretPatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
