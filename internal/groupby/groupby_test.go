package groupby_test

import (
	"strings"
	"testing"

	"github.com/envlens/internal/groupby"
	"github.com/envlens/internal/parser"
)

func makeEnv(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestByPrefix_BasicGrouping(t *testing.T) {
	env := makeEnv(
		"DB_HOST", "localhost",
		"DB_PORT", "5432",
		"APP_NAME", "envlens",
		"APP_ENV", "production",
	)
	r := groupby.ByPrefix(env, "_", 1)
	if len(r.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(r.Groups))
	}
	if r.Groups[0].Prefix != "APP" {
		t.Errorf("expected first group APP, got %s", r.Groups[0].Prefix)
	}
	if r.Groups[1].Prefix != "DB" {
		t.Errorf("expected second group DB, got %s", r.Groups[1].Prefix)
	}
	if len(r.Ungrouped) != 0 {
		t.Errorf("expected no ungrouped entries, got %d", len(r.Ungrouped))
	}
}

func TestByPrefix_UngroupedKeys(t *testing.T) {
	env := makeEnv(
		"DB_HOST", "localhost",
		"PORT", "8080",
	)
	r := groupby.ByPrefix(env, "_", 1)
	if len(r.Ungrouped) != 1 {
		t.Fatalf("expected 1 ungrouped entry, got %d", len(r.Ungrouped))
	}
	if r.Ungrouped[0].Key != "PORT" {
		t.Errorf("expected ungrouped key PORT, got %s", r.Ungrouped[0].Key)
	}
}

func TestByPrefix_DepthTwo(t *testing.T) {
	env := makeEnv(
		"AWS_S3_BUCKET", "my-bucket",
		"AWS_S3_REGION", "us-east-1",
		"AWS_EC2_AMI", "ami-123",
	)
	r := groupby.ByPrefix(env, "_", 2)
	if len(r.Groups) != 2 {
		t.Fatalf("expected 2 groups at depth 2, got %d", len(r.Groups))
	}
	if r.Groups[0].Prefix != "AWS_EC2" {
		t.Errorf("expected AWS_EC2, got %s", r.Groups[0].Prefix)
	}
}

func TestByPrefix_EmptyEnv(t *testing.T) {
	r := groupby.ByPrefix(nil, "_", 1)
	if len(r.Groups) != 0 {
		t.Errorf("expected 0 groups for empty env")
	}
	if len(r.Ungrouped) != 0 {
		t.Errorf("expected 0 ungrouped for empty env")
	}
}

func TestFormat_ContainsPrefixHeader(t *testing.T) {
	env := makeEnv(
		"DB_HOST", "localhost",
		"DB_PORT", "5432",
	)
	r := groupby.ByPrefix(env, "_", 1)
	out := groupby.Format(r)
	if !strings.Contains(out, "[DB]") {
		t.Errorf("expected output to contain [DB], got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected output to contain DB_HOST=localhost")
	}
}

func TestFormat_UngroupedSection(t *testing.T) {
	env := makeEnv("SOLO", "value")
	r := groupby.ByPrefix(env, "_", 1)
	out := groupby.Format(r)
	if !strings.Contains(out, "[ungrouped]") {
		t.Errorf("expected [ungrouped] section, got:\n%s", out)
	}
}
