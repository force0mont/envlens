package envpin_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/envpin"
)

func TestFormat_ShowsAddedKey(t *testing.T) {
	p := envpin.Create(makeEnv())
	drift := envpin.Detect(p, makeEnv("NEW_KEY", "hello"))
	out := envpin.Format(drift)
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "+") {
		t.Errorf("expected '+' prefix for added key, got: %s", out)
	}
}

func TestFormat_ShowsRemovedKey(t *testing.T) {
	p := envpin.Create(makeEnv("OLD_KEY", "bye"))
	drift := envpin.Detect(p, makeEnv())
	out := envpin.Format(drift)
	if !strings.Contains(out, "OLD_KEY") {
		t.Errorf("expected OLD_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "-") {
		t.Errorf("expected '-' prefix for removed key, got: %s", out)
	}
}

func TestFormat_ShowsChangedKey(t *testing.T) {
	p := envpin.Create(makeEnv("TOKEN", "old-secret"))
	drift := envpin.Detect(p, makeEnv("TOKEN", "new-secret"))
	out := envpin.Format(drift)
	if !strings.Contains(out, "TOKEN") {
		t.Errorf("expected TOKEN in output, got: %s", out)
	}
	if !strings.Contains(out, "~") {
		t.Errorf("expected '~' prefix for changed key, got: %s", out)
	}
	if !strings.Contains(out, "old-secret") {
		t.Errorf("expected old value in output, got: %s", out)
	}
}

func TestFormat_SummaryCount(t *testing.T) {
	p := envpin.Create(makeEnv("A", "1", "B", "2"))
	current := makeEnv("A", "changed", "C", "new")
	drift := envpin.Detect(p, current)
	out := envpin.Format(drift)
	if !strings.Contains(out, "3 drift(s)") {
		t.Errorf("expected summary with count, got: %s", out)
	}
}

func TestFormat_SortedOutput(t *testing.T) {
	p := envpin.Create(makeEnv("Z_KEY", "1", "A_KEY", "2"))
	drift := envpin.Detect(p, makeEnv("Z_KEY", "changed", "A_KEY", "changed"))
	out := envpin.Format(drift)
	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA > idxZ {
		t.Error("expected output sorted alphabetically (A before Z)")
	}
}
