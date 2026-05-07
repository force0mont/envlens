package promote

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

func TestPromote_AllKeys(t *testing.T) {
	src := makeEnv("FOO", "bar", "BAZ", "qux")
	dst := makeEnv("EXISTING", "yes")
	out, results := Promote(src, dst, nil)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Error("expected src keys to appear in output")
	}
	if out["EXISTING"] != "yes" {
		t.Error("existing dst key should be preserved")
	}
}

func TestPromote_SelectedKeys(t *testing.T) {
	src := makeEnv("FOO", "1", "BAR", "2", "BAZ", "3")
	dst := makeEnv()
	_, results := Promote(src, dst, []string{"FOO", "BAZ"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Key == "BAR" {
			t.Error("BAR should not have been promoted")
		}
	}
}

func TestPromote_OverwriteFlag(t *testing.T) {
	src := makeEnv("KEY", "new")
	dst := makeEnv("KEY", "old")
	out, results := Promote(src, dst, nil)
	if out["KEY"] != "new" {
		t.Errorf("expected overwritten value 'new', got %q", out["KEY"])
	}
	if !results[0].Overwrote {
		t.Error("expected Overwrote to be true")
	}
}

func TestPromote_MissingKeyInSrc(t *testing.T) {
	src := makeEnv("FOO", "bar")
	dst := makeEnv()
	_, results := Promote(src, dst, []string{"MISSING"})
	if len(results) != 0 {
		t.Errorf("expected 0 results for missing key, got %d", len(results))
	}
}

func TestFormat_NoResults(t *testing.T) {
	out := Format(nil, "staging", "production")
	if !strings.Contains(out, "No keys promoted") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_WithResults(t *testing.T) {
	results := []Result{
		{Key: "FOO", Value: "bar", Overwrote: false},
		{Key: "BAZ", Value: "qux", Overwrote: true},
	}
	out := Format(results, "staging", "production")
	if !strings.Contains(out, "+ FOO") {
		t.Errorf("expected '+ FOO' in output: %q", out)
	}
	if !strings.Contains(out, "~ BAZ (overwritten)") {
		t.Errorf("expected '~ BAZ (overwritten)' in output: %q", out)
	}
}
