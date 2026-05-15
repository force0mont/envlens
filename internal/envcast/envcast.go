package envcast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/envlens/envlens/internal/parser"
)

// TargetType represents the type to cast an env value to.
type TargetType string

const (
	TypeString TargetType = "string"
	TypeInt    TargetType = "int"
	TypeFloat  TargetType = "float"
	TypeBool   TargetType = "bool"
)

// Rule defines a cast rule for a single key.
type Rule struct {
	Key  string
	Type TargetType
}

// Result holds the outcome of casting a single entry.
type Result struct {
	Key     string
	RawValue string
	CastValue string
	Type    TargetType
	Err     error
}

// Cast applies type-cast rules to the given env entries and returns results.
func Cast(entries []parser.Entry, rules []Rule) []Result {
	index := make(map[string]string, len(entries))
	for _, e := range entries {
		index[e.Key] = e.Value
	}

	results := make([]Result, 0, len(rules))
	for _, r := range rules {
		raw, ok := index[r.Key]
		if !ok {
			results = append(results, Result{
				Key:  r.Key,
				Type: r.Type,
				Err:  fmt.Errorf("key %q not found", r.Key),
			})
			continue
		}
		cast, err := castValue(raw, r.Type)
		results = append(results, Result{
			Key:       r.Key,
			RawValue:  raw,
			CastValue: cast,
			Type:      r.Type,
			Err:       err,
		})
	}
	return results
}

func castValue(raw string, t TargetType) (string, error) {
	switch t {
	case TypeString:
		return raw, nil
	case TypeInt:
		v, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to int: %w", raw, err)
		}
		return strconv.FormatInt(v, 10), nil
	case TypeFloat:
		v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to float: %w", raw, err)
		}
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case TypeBool:
		v, err := strconv.ParseBool(strings.TrimSpace(raw))
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to bool: %w", raw, err)
		}
		return strconv.FormatBool(v), nil
	default:
		return "", fmt.Errorf("unknown type %q", t)
	}
}

// Format renders cast results as a human-readable string.
func Format(results []Result) string {
	if len(results) == 0 {
		return "No cast rules applied.\n"
	}
	var sb strings.Builder
	ok, errs := 0, 0
	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(&sb, "  ERROR  %-24s %s\n", r.Key, r.Err)
			errs++
		} else {
			fmt.Fprintf(&sb, "  OK     %-24s %-10s => %s\n", r.Key, "("+string(r.Type)+")", r.CastValue)
			ok++
		}
	}
	fmt.Fprintf(&sb, "\n%d cast(s) succeeded, %d failed.\n", ok, errs)
	return sb.String()
}
