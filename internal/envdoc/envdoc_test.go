package envdoc_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/envdoc"
)

func TestGenerate_BasicEntry(t *testing.T) {
	lines := []string{
		"# The application port",
		"PORT=8080",
	}
	doc := envdoc.Generate(lines)
	if len(doc.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(doc.Entries))
	}
	e := doc.Entries[0]
	if e.Key != "PORT" {
		t.Errorf("expected key PORT, got %s", e.Key)
	}
	if e.Value != "8080" {
		t.Errorf("expected value 8080, got %s", e.Value)
	}
	if e.Comment != "The application port" {
		t.Errorf("unexpected comment: %q", e.Comment)
	}
}

func TestGenerate_RequiredFlag(t *testing.T) {
	lines := []string{
		"# required: database connection string",
		"DATABASE_URL=postgres://localhost/db",
	}
	doc := envdoc.Generate(lines)
	if !doc.Entries[0].Required {
		t.Error("expected entry to be marked required")
	}
}

func TestGenerate_NoComment(t *testing.T) {
	lines := []string{
		"SECRET_KEY=abc123",
	}
	doc := envdoc.Generate(lines)
	if len(doc.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(doc.Entries))
	}
	if doc.Entries[0].Comment != "" {
		t.Errorf("expected empty comment, got %q", doc.Entries[0].Comment)
	}
}

func TestGenerate_BlankLineResetsComment(t *testing.T) {
	lines := []string{
		"# stale comment",
		"",
		"PORT=9000",
	}
	doc := envdoc.Generate(lines)
	if doc.Entries[0].Comment != "" {
		t.Errorf("expected comment to be cleared by blank line, got %q", doc.Entries[0].Comment)
	}
}

func TestFormat_NoEntries(t *testing.T) {
	out := envdoc.Format(envdoc.Doc{})
	if !strings.Contains(out, "No documented") {
		t.Errorf("expected fallback message, got: %s", out)
	}
}

func TestFormat_WithEntries(t *testing.T) {
	doc := envdoc.Doc{
		Entries: []envdoc.Entry{
			{Key: "PORT", Value: "8080", Comment: "HTTP port", Required: false},
			{Key: "DB_URL", Value: "", Comment: "required db url", Required: true},
		},
	}
	out := envdoc.Format(doc)
	if !strings.Contains(out, "PORT") {
		t.Error("expected PORT in output")
	}
	if !strings.Contains(out, "yes") {
		t.Error("expected 'yes' for required field")
	}
	if !strings.Contains(out, "HTTP port") {
		t.Error("expected comment in output")
	}
}
