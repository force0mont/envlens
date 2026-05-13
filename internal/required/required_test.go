package required

import (
	"strings"
	"testing"

	"github.com/envlens/envlens/internal/parser"
)

func makeEnv(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestCheck_AllPresent(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "DB_PORT", "5432")
	r := Check(env, []string{"DB_HOST", "DB_PORT"})
	if len(r.Findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(r.Findings))
	}
	if r.Checked != 2 {
		t.Fatalf("expected Checked=2, got %d", r.Checked)
	}
}

func TestCheck_MissingKey(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost")
	r := Check(env, []string{"DB_HOST", "DB_PORT"})
	if len(r.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(r.Findings))
	}
	if r.Findings[0].Key != "DB_PORT" || r.Findings[0].Reason != "missing" {
		t.Errorf("unexpected finding: %+v", r.Findings[0])
	}
}

func TestCheck_EmptyValue(t *testing.T) {
	env := makeEnv("SECRET_KEY", "")
	r := Check(env, []string{"SECRET_KEY"})
	if len(r.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(r.Findings))
	}
	if r.Findings[0].Reason != "empty value" {
		t.Errorf("expected 'empty value', got %q", r.Findings[0].Reason)
	}
}

func TestCheck_EmptyRequired(t *testing.T) {
	env := makeEnv("A", "1")
	r := Check(env, []string{})
	if len(r.Findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(r.Findings))
	}
	if r.Checked != 0 {
		t.Fatalf("expected Checked=0, got %d", r.Checked)
	}
}

func TestFormat_NoFindings(t *testing.T) {
	r := Result{Findings: nil, Checked: 3}
	out := Format(r)
	if !strings.Contains(out, "OK") {
		t.Errorf("expected OK in output, got: %s", out)
	}
}

func TestFormat_WithFindings(t *testing.T) {
	r := Result{
		Findings: []Finding{{Key: "DB_URL", Reason: "missing"}},
		Checked:  2,
	}
	out := Format(r)
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected DB_URL in output, got: %s", out)
	}
	if !strings.Contains(out, "missing") {
		t.Errorf("expected 'missing' in output, got: %s", out)
	}
}
