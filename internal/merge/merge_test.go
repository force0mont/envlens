package merge

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/parser"
)

func makeEnv(pairs ...string) parser.EnvFile {
	env := make(parser.EnvFile)
	for i := 0; i+1 < len(pairs); i += 2 {
		env[pairs[i]] = pairs[i+1]
	}
	return env
}

func TestMerge_NoConflicts(t *testing.T) {
	a := makeEnv("FOO", "bar", "BAZ", "qux")
	b := makeEnv("HELLO", "world")
	r := Merge([]parser.EnvFile{a, b}, StrategyFirst)
	if len(r.Conflicts) != 0 {
		t.Errorf("expected 0 conflicts, got %d", len(r.Conflicts))
	}
	if r.Env["FOO"] != "bar" || r.Env["HELLO"] != "world" {
		t.Errorf("unexpected merged env: %v", r.Env)
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	a := makeEnv("KEY", "original")
	b := makeEnv("KEY", "override")
	r := Merge([]parser.EnvFile{a, b}, StrategyFirst)
	if r.Env["KEY"] != "original" {
		t.Errorf("expected 'original', got %q", r.Env["KEY"])
	}
	if len(r.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(r.Conflicts))
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	a := makeEnv("KEY", "original")
	b := makeEnv("KEY", "override")
	r := Merge([]parser.EnvFile{a, b}, StrategyLast)
	if r.Env["KEY"] != "override" {
		t.Errorf("expected 'override', got %q", r.Env["KEY"])
	}
}

func TestMerge_IdenticalValueNoConflict(t *testing.T) {
	a := makeEnv("KEY", "same")
	b := makeEnv("KEY", "same")
	r := Merge([]parser.EnvFile{a, b}, StrategyFirst)
	if len(r.Conflicts) != 0 {
		t.Errorf("identical values should not produce a conflict")
	}
}

func TestFormat_NoConflicts(t *testing.T) {
	r := Result{Env: makeEnv("A", "1", "B", "2"), Conflicts: nil}
	out := Format(r, nil)
	if !strings.Contains(out, "No conflicts") {
		t.Errorf("expected no-conflict message, got: %s", out)
	}
	if !strings.Contains(out, "2 key(s)") {
		t.Errorf("expected key count in output, got: %s", out)
	}
}

func TestFormat_WithLabels(t *testing.T) {
	r := Result{
		Env: makeEnv("KEY", "override"),
		Conflicts: []Conflict{
			{Key: "KEY", Values: []string{"original", "override"}},
		},
	}
	out := Format(r, []string{"prod", "staging"})
	if !strings.Contains(out, "[prod]") || !strings.Contains(out, "[staging]") {
		t.Errorf("expected labels in output, got: %s", out)
	}
}
