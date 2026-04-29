package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeGitLabSyncFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeGitLabSyncFile: %v", err)
	}
	return p
}

func TestRunGitLabSync_AlreadyInSync(t *testing.T) {
	p := writeGitLabSyncFile(t, "KEY=value\n")
	var buf bytes.Buffer
	err := RunGitLabSync(GitLabSyncOptions{
		BaseFile:  p,
		Token:     "tok",
		ProjectID: "123",
		Namespace: "prod",
		DryRun:    false,
		Out:       &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "already in sync") {
		t.Errorf("expected 'already in sync', got: %s", buf.String())
	}
}

func TestRunGitLabSync_DryRun(t *testing.T) {
	p := writeGitLabSyncFile(t, "NEW_KEY=hello\n")
	var buf bytes.Buffer
	err := RunGitLabSync(GitLabSyncOptions{
		BaseFile:  p,
		Token:     "tok",
		ProjectID: "456",
		Namespace: "staging",
		DryRun:    true,
		Out:       &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "dry-run") {
		t.Errorf("expected dry-run message, got: %s", buf.String())
	}
}

func TestRunGitLabSync_InvalidBaseFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunGitLabSync(GitLabSyncOptions{
		BaseFile:  "/nonexistent/.env",
		Token:     "tok",
		ProjectID: "789",
		Namespace: "prod",
		DryRun:    false,
		Out:       &buf,
	})
	if err == nil {
		t.Fatal("expected error for invalid base file, got nil")
	}
	if !strings.Contains(err.Error(), "parsing base file") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunGitLabSync_EmptyNamespace(t *testing.T) {
	p := writeGitLabSyncFile(t, "FOO=bar\n")
	var buf bytes.Buffer
	err := RunGitLabSync(GitLabSyncOptions{
		BaseFile:  p,
		Token:     "tok",
		ProjectID: "321",
		Namespace: "",
		DryRun:    true,
		Out:       &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error with empty namespace: %v", err)
	}
}
