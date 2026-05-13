package defaults

import (
	"strings"
	"testing"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestApply_FillsMissingKeys(t *testing.T) {
	env := makeEnv("EXISTING", "hello")
	defs := makeEnv("MISSING", "default_val", "EXISTING", "ignored")

	r := Apply(env, defs)

	if r.Env["MISSING"] != "default_val" {
		t.Errorf("expected MISSING=default_val, got %q", r.Env["MISSING"])
	}
	if r.Env["EXISTING"] != "hello" {
		t.Errorf("expected EXISTING=hello, got %q", r.Env["EXISTING"])
	}
}

func TestApply_FillsEmptyValue(t *testing.T) {
	env := makeEnv("KEY", "")
	defs := makeEnv("KEY", "fallback")

	r := Apply(env, defs)

	if r.Env["KEY"] != "fallback" {
		t.Errorf("expected KEY=fallback, got %q", r.Env["KEY"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := makeEnv("A", "")
	defs := makeEnv("A", "default")

	Apply(env, defs)

	if env["A"] != "" {
		t.Error("Apply must not mutate the original env map")
	}
}

func TestApply_EntryAppliedFlag(t *testing.T) {
	env := makeEnv("SET", "value")
	defs := makeEnv("SET", "ignored", "UNSET", "used")

	r := Apply(env, defs)

	for _, e := range r.Entries {
		if e.Key == "SET" && e.Applied {
			t.Error("SET should not be marked applied")
		}
		if e.Key == "UNSET" && !e.Applied {
			t.Error("UNSET should be marked applied")
		}
	}
}

func TestApply_NoDefaults(t *testing.T) {
	env := makeEnv("X", "1")
	r := Apply(env, map[string]string{})

	if len(r.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries))
	}
	if r.Env["X"] != "1" {
		t.Errorf("expected X=1, got %q", r.Env["X"])
	}
}

func TestFormat_ShowsAppliedAndSkipped(t *testing.T) {
	env := makeEnv("PRESENT", "yes")
	defs := makeEnv("PRESENT", "no", "ABSENT", "default")

	r := Apply(env, defs)
	out := Format(r)

	if !strings.Contains(out, "applied") {
		t.Error("expected 'applied' in output")
	}
	if !strings.Contains(out, "skipped") {
		t.Error("expected 'skipped' in output")
	}
	if !strings.Contains(out, "ABSENT=default") {
		t.Error("expected ABSENT=default in output")
	}
}

func TestFormat_NoDefaultsDefined(t *testing.T) {
	r := Result{Env: makeEnv(), Entries: nil}
	out := Format(r)

	if !strings.Contains(out, "no defaults defined") {
		t.Errorf("unexpected output: %q", out)
	}
}
