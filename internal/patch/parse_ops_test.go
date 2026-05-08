package patch

import (
	"testing"
)

func TestParseInstruction_SetValid(t *testing.T) {
	ins := Instruction{Op: OpSet, Key: "FOO", Value: "bar"}
	if ins.Op != OpSet {
		t.Errorf("expected OpSet, got %q", ins.Op)
	}
	if ins.Key != "FOO" || ins.Value != "bar" {
		t.Errorf("unexpected instruction: %+v", ins)
	}
}

func TestParseInstruction_DeleteValid(t *testing.T) {
	ins := Instruction{Op: OpDelete, Key: "BAR"}
	if ins.Op != OpDelete {
		t.Errorf("expected OpDelete, got %q", ins.Op)
	}
	if ins.Key != "BAR" {
		t.Errorf("unexpected key: %q", ins.Key)
	}
}

func TestParseInstruction_RenameValid(t *testing.T) {
	ins := Instruction{Op: OpRename, Key: "OLD", NewKey: "NEW"}
	if ins.NewKey != "NEW" {
		t.Errorf("expected NewKey=NEW, got %q", ins.NewKey)
	}
}

func TestApply_MultipleOps(t *testing.T) {
	env := makeEnv("A", "1", "B", "2", "C", "3")
	instructions := []Instruction{
		{Op: OpSet, Key: "A", Value: "updated"},
		{Op: OpDelete, Key: "B"},
		{Op: OpRename, Key: "C", NewKey: "D"},
	}
	result := Apply(env, instructions)
	if result.Env["A"] != "updated" {
		t.Errorf("expected A=updated, got %q", result.Env["A"])
	}
	if _, ok := result.Env["B"]; ok {
		t.Error("expected B to be deleted")
	}
	if result.Env["D"] != "3" {
		t.Errorf("expected D=3, got %q", result.Env["D"])
	}
	if len(result.Applied) != 3 {
		t.Errorf("expected 3 applied, got %d", len(result.Applied))
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(result.Skipped))
	}
}

func TestApply_EmptyInstructions(t *testing.T) {
	env := makeEnv("FOO", "bar")
	result := Apply(env, nil)
	if result.Env["FOO"] != "bar" {
		t.Error("expected env to be unchanged")
	}
	if len(result.Applied) != 0 {
		t.Errorf("expected 0 applied, got %d", len(result.Applied))
	}
}

func TestFormat_EmptyEnv(t *testing.T) {
	result := Apply(map[string]string{}, nil)
	out := Format(result)
	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}
