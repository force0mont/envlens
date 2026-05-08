package patch

import (
	"testing"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestApply_SetNewKey(t *testing.T) {
	env := makeEnv("FOO", "bar")
	result := Apply(env, []Instruction{{Op: OpSet, Key: "BAZ", Value: "qux"}})
	if result.Env["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", result.Env["BAZ"])
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(result.Applied))
	}
}

func TestApply_SetOverwritesExisting(t *testing.T) {
	env := makeEnv("FOO", "old")
	result := Apply(env, []Instruction{{Op: OpSet, Key: "FOO", Value: "new"}})
	if result.Env["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %q", result.Env["FOO"])
	}
}

func TestApply_DeleteExistingKey(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	result := Apply(env, []Instruction{{Op: OpDelete, Key: "FOO"}})
	if _, ok := result.Env["FOO"]; ok {
		t.Error("expected FOO to be deleted")
	}
	if result.Env["BAZ"] != "qux" {
		t.Error("expected BAZ to be preserved")
	}
}

func TestApply_DeleteMissingKey_Skipped(t *testing.T) {
	env := makeEnv("FOO", "bar")
	result := Apply(env, []Instruction{{Op: OpDelete, Key: "MISSING"}})
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if result.SkipMsgs["MISSING"] == "" {
		t.Error("expected skip message for MISSING")
	}
}

func TestApply_RenameKey(t *testing.T) {
	env := makeEnv("OLD_KEY", "value")
	result := Apply(env, []Instruction{{Op: OpRename, Key: "OLD_KEY", NewKey: "NEW_KEY"}})
	if _, ok := result.Env["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if result.Env["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", result.Env["NEW_KEY"])
	}
}

func TestApply_RenameMissingKey_Skipped(t *testing.T) {
	env := makeEnv("FOO", "bar")
	result := Apply(env, []Instruction{{Op: OpRename, Key: "GHOST", NewKey: "SPIRIT"}})
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	env := makeEnv("FOO", "bar")
	Apply(env, []Instruction{{Op: OpSet, Key: "FOO", Value: "changed"}})
	if env["FOO"] != "bar" {
		t.Error("original env was mutated")
	}
}

func TestFormat_IncludesSkipped(t *testing.T) {
	env := makeEnv("FOO", "bar")
	result := Apply(env, []Instruction{{Op: OpDelete, Key: "MISSING"}})
	out := Format(result)
	if out == "" {
		t.Error("expected non-empty format output")
	}
	if len(result.Skipped) == 0 {
		t.Error("expected skipped entries in result")
	}
}
