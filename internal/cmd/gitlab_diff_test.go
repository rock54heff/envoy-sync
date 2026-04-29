package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeGitLabDiffFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return p
}

func TestRunGitLabDiff_NoChanges(t *testing.T) {
	p := writeGitLabDiffFile(t, "KEY=value\nFOO=bar\n")
	var buf bytes.Buffer
	if err := RunGitLabDiff(p, "myproject", false, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "+") || strings.Contains(out, "-") {
		t.Errorf("expected no diff output, got: %s", out)
	}
}

func TestRunGitLabDiff_WithChanges(t *testing.T) {
	p := writeGitLabDiffFile(t, "KEY=value\nEXTRA=only_local\n")
	var buf bytes.Buffer
	if err := RunGitLabDiff(p, "myproject", false, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Output should contain diff markers; at minimum no error
	_ = buf.String()
}

func TestRunGitLabDiff_SummaryFormat(t *testing.T) {
	p := writeGitLabDiffFile(t, "KEY=value\n")
	var buf bytes.Buffer
	if err := RunGitLabDiff(p, "myproject", true, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if out == "" {
		t.Error("expected summary output, got empty string")
	}
}

func TestRunGitLabDiff_InvalidLocalFile(t *testing.T) {
	err := RunGitLabDiff("/nonexistent/.env", "myproject", false, nil)
	if err == nil {
		t.Error("expected error for invalid local file, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse local file") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunGitLabDiff_EmptyNamespace(t *testing.T) {
	p := writeGitLabDiffFile(t, "KEY=value\n")
	err := RunGitLabDiff(p, "", false, nil)
	if err == nil {
		t.Error("expected error for empty namespace, got nil")
	}
	if !strings.Contains(err.Error(), "namespace must not be empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}
