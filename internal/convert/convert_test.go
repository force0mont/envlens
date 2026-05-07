package convert_test

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/convert"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestConvert_EnvFormat(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	res, err := convert.Convert(env, convert.FormatEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got:\n%s", res.Output)
	}
	if !strings.Contains(res.Output, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got:\n%s", res.Output)
	}
}

func TestConvert_ExportFormat(t *testing.T) {
	env := makeEnv("PORT", "8080")
	res, err := convert.Convert(env, convert.FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "export PORT=8080") {
		t.Errorf("expected export prefix, got:\n%s", res.Output)
	}
}

func TestConvert_JSONFormat(t *testing.T) {
	env := makeEnv("APP", "envlens")
	res, err := convert.Convert(env, convert.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, `"APP"`) {
		t.Errorf("expected JSON key APP, got:\n%s", res.Output)
	}
	if !strings.Contains(res.Output, `"envlens"`) {
		t.Errorf("expected JSON value envlens, got:\n%s", res.Output)
	}
}

func TestConvert_YAMLFormat(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost")
	res, err := convert.Convert(env, convert.FormatYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "DB_HOST:") {
		t.Errorf("expected YAML key DB_HOST, got:\n%s", res.Output)
	}
}

func TestConvert_UnsupportedFormat(t *testing.T) {
	env := makeEnv("X", "1")
	_, err := convert.Convert(env, convert.Format("toml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestConvert_EmptyEnv(t *testing.T) {
	res, err := convert.Convert(map[string]string{}, convert.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "{") {
		t.Errorf("expected empty JSON object, got:\n%s", res.Output)
	}
}
