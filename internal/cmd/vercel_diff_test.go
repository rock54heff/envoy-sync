package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeVercelEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return p
}

func TestRunVercelDiff_NoChanges(t *testing.T) {
	p := writeVercelEnvFile(t, "KEY=value\n")
	var buf bytes.Buffer
	err := RunVercelDiff(VercelDiffOptions{
		LocalFile: p,
		Namespace: "my-project",
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunVercelDiff_WithChanges(t *testing.T) {
	p := writeVercelEnvFile(t, "KEY=value\nEXTRA=only_local\n")
	var buf bytes.Buffer
	err := RunVercelDiff(VercelDiffOptions{
		LocalFile: p,
		Namespace: "my-project",
		Verbose:   true,
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunVercelDiff_SummaryFormat(t *testing.T) {
	p := writeVercelEnvFile(t, "FOO=bar\n")
	var buf bytes.Buffer
	err := RunVercelDiff(VercelDiffOptions{
		LocalFile: p,
		Namespace: "my-project",
		Summary:   true,
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected summary output, got empty")
	}
}

func TestRunVercelDiff_InvalidLocalFile(t *testing.T) {
	err := RunVercelDiff(VercelDiffOptions{
		LocalFile: "/nonexistent/.env",
		Namespace: "my-project",
	}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for invalid local file, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse local file") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunVercelDiff_EmptyNamespace(t *testing.T) {
	p := writeVercelEnvFile(t, "KEY=value\n")
	err := RunVercelDiff(VercelDiffOptions{
		LocalFile: p,
		Namespace: "",
	}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for empty namespace, got nil")
	}
	if !strings.Contains(err.Error(), "namespace") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunVercelDiff_EmptyLocalFile(t *testing.T) {
	err := RunVercelDiff(VercelDiffOptions{
		LocalFile: "",
		Namespace: "my-project",
	}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for empty local file path, got nil")
	}
	if !strings.Contains(err.Error(), "local file path is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}
