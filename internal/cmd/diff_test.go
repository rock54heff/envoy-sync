package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-sync/internal/cmd"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestRunDiff_NoChanges(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	target := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	opts := cmd.DiffOptions{
		BaseFile:     base,
		TargetFile:   target,
		Verbose:      false,
		OutputFormat: "text",
	}
	if err := cmd.RunDiff(opts); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRunDiff_WithChanges(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nOLD=value\n")
	target := writeTempEnv(t, "FOO=changed\nNEW=value\n")

	opts := cmd.DiffOptions{
		BaseFile:     base,
		TargetFile:   target,
		Verbose:      true,
		OutputFormat: "text",
	}
	if err := cmd.RunDiff(opts); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRunDiff_SummaryFormat(t *testing.T) {
	base := writeTempEnv(t, "A=1\nB=2\n")
	target := writeTempEnv(t, "A=1\nC=3\n")

	opts := cmd.DiffOptions{
		BaseFile:     base,
		TargetFile:   target,
		OutputFormat: "summary",
	}
	if err := cmd.RunDiff(opts); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRunDiff_InvalidBaseFile(t *testing.T) {
	target := writeTempEnv(t, "FOO=bar\n")

	opts := cmd.DiffOptions{
		BaseFile:   "/nonexistent/.env",
		TargetFile: target,
	}
	if err := cmd.RunDiff(opts); err == nil {
		t.Error("expected error for invalid base file, got nil")
	}
}

func TestRunDiff_InvalidTargetFile(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\n")

	opts := cmd.DiffOptions{
		BaseFile:   base,
		TargetFile: "/nonexistent/.env",
	}
	if err := cmd.RunDiff(opts); err == nil {
		t.Error("expected error for invalid target file, got nil")
	}
}
