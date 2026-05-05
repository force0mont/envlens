package audit

import (
	"fmt"
	"strings"
)

// Format renders a Report as a human-readable string.
func Format(report Report) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Audit Report: %s\n", report.File))
	sb.WriteString(strings.Repeat("-", 50) + "\n")

	if len(report.Findings) == 0 {
		sb.WriteString("  No issues found.\n")
		return sb.String()
	}

	counts := map[Severity]int{}
	for _, f := range report.Findings {
		counts[f.Severity]++
	}

	for _, f := range report.Findings {
		sb.WriteString(fmt.Sprintf("  [%-8s] %s: %s\n", f.Severity, f.Key, f.Message))
	}

	sb.WriteString(strings.Repeat("-", 50) + "\n")
	sb.WriteString(fmt.Sprintf("  Summary: %d critical, %d warning, %d info\n",
		counts[SeverityCritical],
		counts[SeverityWarning],
		counts[SeverityInfo],
	))

	return sb.String()
}
