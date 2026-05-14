// Package envdoc generates documentation for environment variables
// by extracting inline comments from .env files.
package envdoc

import (
	"fmt"
	"strings"
)

// Entry holds a documented environment variable.
type Entry struct {
	Key     string
	Value   string
	Comment string
	Required bool
}

// Doc holds the full documentation result for an env file.
type Doc struct {
	Entries []Entry
}

// Generate produces documentation entries from parsed env lines.
// rawLines should be the original lines of the .env file (including comments).
func Generate(rawLines []string) Doc {
	var entries []Entry
	var pendingComment string

	for _, line := range rawLines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			pendingComment = ""
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			comment := strings.TrimSpace(strings.TrimPrefix(trimmed, "#"))
			pendingComment = comment
			continue
		}

		key, value, ok := parsePair(trimmed)
		if !ok {
			pendingComment = ""
			continue
		}

		required := strings.Contains(strings.ToLower(pendingComment), "required")

		entries = append(entries, Entry{
			Key:      key,
			Value:    value,
			Comment:  pendingComment,
			Required: required,
		})
		pendingComment = ""
	}

	return Doc{Entries: entries}
}

// Format renders the documentation as a human-readable Markdown table.
func Format(d Doc) string {
	if len(d.Entries) == 0 {
		return "No documented variables found.\n"
	}

	var sb strings.Builder
	sb.WriteString("| Variable | Default | Required | Description |\n")
	sb.WriteString("|----------|---------|----------|-------------|\n")

	for _, e := range d.Entries {
		req := "no"
		if e.Required {
			req = "yes"
		}
		desc := e.Comment
		if desc == "" {
			desc = "-"
		}
		sb.WriteString(fmt.Sprintf("| `%s` | `%s` | %s | %s |\n", e.Key, e.Value, req, desc))
	}

	return sb.String()
}

func parsePair(line string) (key, value string, ok bool) {
	idx := strings.IndexByte(line, '=')
	if idx < 1 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:idx])
	value = strings.Trim(strings.TrimSpace(line[idx+1:]), `"'`)
	return key, value, true
}
