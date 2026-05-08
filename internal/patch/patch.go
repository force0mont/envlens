package patch

import (
	"fmt"
	"sort"
	"strings"
)

// Op represents the type of patch operation.
type Op string

const (
	OpSet    Op = "set"
	OpDelete Op = "delete"
	OpRename Op = "rename"
)

// Instruction describes a single patch operation.
type Instruction struct {
	Op      Op
	Key     string
	Value   string // used by set
	NewKey  string // used by rename
}

// Result holds the outcome of applying a patch.
type Result struct {
	Env      map[string]string
	Applied  []Instruction
	Skipped  []Instruction
	SkipMsgs map[string]string
}

// Apply applies a slice of Instructions to the given env map.
// It returns a new map and a Result describing what happened.
func Apply(env map[string]string, instructions []Instruction) Result {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	result := Result{
		Env:      out,
		SkipMsgs: make(map[string]string),
	}

	for _, ins := range instructions {
		switch ins.Op {
		case OpSet:
			out[ins.Key] = ins.Value
			result.Applied = append(result.Applied, ins)
		case OpDelete:
			if _, ok := out[ins.Key]; !ok {
				result.Skipped = append(result.Skipped, ins)
				result.SkipMsgs[ins.Key] = fmt.Sprintf("key %q not found", ins.Key)
				continue
			}
			delete(out, ins.Key)
			result.Applied = append(result.Applied, ins)
		case OpRename:
			val, ok := out[ins.Key]
			if !ok {
				result.Skipped = append(result.Skipped, ins)
				result.SkipMsgs[ins.Key] = fmt.Sprintf("key %q not found", ins.Key)
				continue
			}
			delete(out, ins.Key)
			out[ins.NewKey] = val
			result.Applied = append(result.Applied, ins)
		}
	}

	return result
}

// Format returns a human-readable summary of the patch result.
func Format(r Result) string {
	var sb strings.Builder

	keys := make([]string, 0, len(r.Env))
	for k := range r.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, r.Env[k])
	}

	if len(r.Skipped) > 0 {
		sb.WriteString("\n# Skipped:\n")
		for _, ins := range r.Skipped {
			fmt.Fprintf(&sb, "#  [%s] %s: %s\n", ins.Op, ins.Key, r.SkipMsgs[ins.Key])
		}
	}

	return sb.String()
}
