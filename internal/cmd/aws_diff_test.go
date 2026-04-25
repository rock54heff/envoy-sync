package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeAWSEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestRunAWSDiff_NoChanges(t *testing.T) {
	path := writeAWSEnvFile(t, "KEY=value\n")
	var buf bytes.Buffer
	err := RunAWSDiff(AWSDiffOptions{
		LocalFile:  path,
		Namespace:  "/app/prod",
		RemoteVars: map[string]string{"/app/prod/KEY": "value"},
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "+") || strings.Contains(buf.String(), "-") {
		t.Errorf("expected no diff, got: %s", buf.String())
	}
}

func TestRunAWSDiff_WithChanges(t *testing.T) {
	path := writeAWSEnvFile(t, "KEY=new\nADDED=yes\n")
	var buf bytes.Buffer
	err := RunAWSDiff(AWSDiffOptions{
		LocalFile:  path,
		Namespace:  "/app/prod",
		RemoteVars: map[string]string{"/app/prod/KEY": "old"},
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ADDED") {
		t.Errorf("expected ADDED in diff output, got: %s", out)
	}
}

func TestRunAWSDiff_SummaryFormat(t *testing.T) {
	path := writeAWSEnvFile(t, "A=1\nB=2\n")
	var buf bytes.Buffer
	err := RunAWSDiff(AWSDiffOptions{
		LocalFile:  path,
		Namespace:  "/ns",
		RemoteVars: map[string]string{"/ns/A": "1"},
		Summary:    true,
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "added") && !strings.Contains(buf.String(), "1") {
		t.Errorf("summary output unexpected: %s", buf.String())
	}
}

func TestRunAWSDiff_InvalidLocalFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunAWSDiff(AWSDiffOptions{
		LocalFile: "/nonexistent/.env",
		Namespace: "/ns",
	}, &buf)
	if err == nil {
		t.Fatal("expected error for invalid local file")
	}
}

func TestRunAWSDiff_EmptyNamespace(t *testing.T) {
	path := writeAWSEnvFile(t, "KEY=val\n")
	var buf bytes.Buffer
	err := RunAWSDiff(AWSDiffOptions{
		LocalFile: path,
		Namespace: "",
	}, &buf)
	if err == nil {
		t.Fatal("expected error for empty namespace")
	}
}
