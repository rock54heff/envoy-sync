package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeGitHubEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunGitHubDiff_NoChanges(t *testing.T) {
	path := writeGitHubEnvFile(t, "TOKEN=abc\nSECRET=xyz\n")
	var buf bytes.Buffer
	err := RunGitHubDiff(GitHubDiffOptions{
		LocalFile:  path,
		Repo:       "org/repo",
		Namespace:  "prod",
		RemoteVars: map[string]string{"TOKEN": "abc", "SECRET": "xyz"},
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "+") || strings.Contains(buf.String(), "-") {
		t.Errorf("expected no changes, got: %s", buf.String())
	}
}

func TestRunGitHubDiff_WithChanges(t *testing.T) {
	path := writeGitHubEnvFile(t, "TOKEN=newval\nADDED=yes\n")
	var buf bytes.Buffer
	err := RunGitHubDiff(GitHubDiffOptions{
		LocalFile:  path,
		Repo:       "org/repo",
		Namespace:  "prod",
		RemoteVars: map[string]string{"TOKEN": "oldval"},
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "TOKEN") {
		t.Errorf("expected TOKEN in diff output, got: %s", output)
	}
}

func TestRunGitHubDiff_SummaryFormat(t *testing.T) {
	path := writeGitHubEnvFile(t, "A=1\nB=2\n")
	var buf bytes.Buffer
	err := RunGitHubDiff(GitHubDiffOptions{
		LocalFile:  path,
		Repo:       "org/repo",
		Namespace:  "staging",
		RemoteVars: map[string]string{"A": "1"},
		Summary:    true,
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "added") && !strings.Contains(buf.String(), "1") {
		t.Errorf("expected summary output, got: %s", buf.String())
	}
}

func TestRunGitHubDiff_InvalidLocalFile(t *testing.T) {
	err := RunGitHubDiff(GitHubDiffOptions{
		LocalFile: "/nonexistent/.env",
		Repo:      "org/repo",
		Namespace: "prod",
	})
	if err == nil {
		t.Error("expected error for invalid local file")
	}
}

func TestRunGitHubDiff_EmptyNamespace(t *testing.T) {
	path := writeGitHubEnvFile(t, "KEY=val\n")
	err := RunGitHubDiff(GitHubDiffOptions{
		LocalFile: path,
		Repo:      "org/repo",
		Namespace: "",
	})
	if err == nil {
		t.Error("expected error for empty namespace")
	}
}
