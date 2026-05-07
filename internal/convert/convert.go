package convert

import (
	"fmt"
	"sort"
	"strings"
)

// Format represents a supported output serialisation format.
type Format string

const (
	FormatEnv    Format = "env"
	FormatJSON   Format = "json"
	FormatYAML   Format = "yaml"
	FormatExport Format = "export"
)

// Result holds the converted output and the format it was rendered in.
type Result struct {
	Format Format
	Output string
}

// Convert serialises env vars from the given map into the requested format.
func Convert(env map[string]string, format Format) (Result, error) {
	keys := sortedKeys(env)

	var sb strings.Builder

	switch format {
	case FormatEnv:
		for _, k := range keys {
			fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
		}

	case FormatExport:
		for _, k := range keys {
			fmt.Fprintf(&sb, "export %s=%s\n", k, env[k])
		}

	case FormatJSON:
		sb.WriteString("{\n")
		for i, k := range keys {
			comma := ","
			if i == len(keys)-1 {
				comma = ""
			}
			fmt.Fprintf(&sb, "  %q: %q%s\n", k, env[k], comma)
		}
		sb.WriteString("}\n")

	case FormatYAML:
		for _, k := range keys {
			fmt.Fprintf(&sb, "%s: %q\n", k, env[k])
		}

	default:
		return Result{}, fmt.Errorf("unsupported format %q: choose one of env, export, json, yaml", format)
	}

	return Result{Format: format, Output: sb.String()}, nil
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
