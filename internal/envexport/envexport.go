package envexport

import (
	"fmt"
	"sort"
	"strings"

	"github.com/envlens/envlens/internal/parser"
)

// Format represents the export output format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatShell  Format = "shell"
	FormatDocker Format = "docker"
	FormatJSON   Format = "json"
)

// Options controls export behaviour.
type Options struct {
	OutputFormat Format
	Prefix       string
	OmitEmpty    bool
}

// Result holds the exported content and metadata.
type Result struct {
	Content string
	Count   int
	Format  Format
}

// Export converts a slice of env entries into the requested output format.
func Export(entries []parser.Entry, opts Options) (Result, error) {
	filtered := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if opts.OmitEmpty && e.Value == "" {
			continue
		}
		filtered = append(filtered, e)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Key < filtered[j].Key
	})

	var content string
	var err error

	switch opts.OutputFormat {
	case FormatShell:
		content = renderShell(filtered, opts.Prefix)
	case FormatDocker:
		content = renderDocker(filtered, opts.Prefix)
	case FormatJSON:
		content = renderJSON(filtered, opts.Prefix)
	case FormatDotenv, "":
		content = renderDotenv(filtered, opts.Prefix)
	default:
		err = fmt.Errorf("unknown format: %q", opts.OutputFormat)
	}

	if err != nil {
		return Result{}, err
	}
	return Result{Content: content, Count: len(filtered), Format: opts.OutputFormat}, nil
}

func renderDotenv(entries []parser.Entry, prefix string) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s%s=%s\n", prefix, e.Key, e.Value)
	}
	return sb.String()
}

func renderShell(entries []parser.Entry, prefix string) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "export %s%s=%q\n", prefix, e.Key, e.Value)
	}
	return sb.String()
}

func renderDocker(entries []parser.Entry, prefix string) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "--env %s%s=%s\n", prefix, e.Key, e.Value)
	}
	return sb.String()
}

func renderJSON(entries []parser.Entry, prefix string) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, e := range entries {
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", prefix+e.Key, e.Value, comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}
