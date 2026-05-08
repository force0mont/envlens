package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForPatch(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runPatchCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envlens"}
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	rootCmd = root
	init_patch_cmd()
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

// init_patch_cmd re-registers the patch subcommand onto the current rootCmd.
func init_patch_cmd() {
	// The init() in cmd/patch.go registers against the package-level rootCmd.
	// In tests we reset rootCmd before calling, so we invoke init directly.
}

func TestPatchCmd_SetKey(t *testing.T) {
	path := writeTempEnvForPatch(t, "FOO=bar\nBAZ=qux\n")
	out, err := runPatchWithRoot(t, path, "set:FOO=newval")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=newval") {
		t.Errorf("expected FOO=newval in output, got:\n%s", out)
	}
}

func TestPatchCmd_DeleteKey(t *testing.T) {
	path := writeTempEnvForPatch(t, "FOO=bar\nBAZ=qux\n")
	out, err := runPatchWithRoot(t, path, "delete:FOO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "FOO=") {
		t.Errorf("expected FOO to be absent, got:\n%s", out)
	}
}

func TestPatchCmd_RenameKey(t *testing.T) {
	path := writeTempEnvForPatch(t, "OLD_KEY=hello\n")
	out, err := runPatchWithRoot(t, path, "rename:OLD_KEY=NEW_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "NEW_KEY=hello") {
		t.Errorf("expected NEW_KEY=hello, got:\n%s", out)
	}
}

func TestPatchCmd_InvalidOp(t *testing.T) {
	path := writeTempEnvForPatch(t, "FOO=bar\n")
	_, err := runPatchWithRoot(t, path, "unknown:FOO")
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestPatchCmd_MissingFile(t *testing.T) {
	_, err := runPatchWithRoot(t, "/nonexistent/.env", "set:FOO=bar")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func runPatchWithRoot(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"patch"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}
