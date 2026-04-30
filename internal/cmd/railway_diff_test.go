package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeRailwayEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return p
}

func TestRunRailwayDiff_NoChanges(t *testing.T) {
	path := writeRailwayEnvFile(t, "FOO=bar\nBAZ=qux\n")
	var buf bytes.Buffer
	err := RunRailwayDiff(RailwayDiffOptions{
		LocalFile:  path,
		ProjectID:  "proj-1",
		Environment: "production",
		RemoteVars: map[string]string{"FOO": "bar", "BAZ": "qux"},
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "+") || strings.Contains(buf.String(), "-") {
		t.Errorf("expected no diff output, got: %s", buf.String())
	}
}

func TestRunRailwayDiff_WithChanges(t *testing.T) {
	path := writeRailwayEnvFile(t, "FOO=newval\nNEWKEY=hello\n")
	var buf bytes.Buffer
	err := RunRailwayDiff(RailwayDiffOptions{
		LocalFile:  path,
		ProjectID:  "proj-1",
		Environment: "staging",
		RemoteVars: map[string]string{"FOO": "oldval", "GONE": "bye"},
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in diff output, got: %s", out)
	}
}

func TestRunRailwayDiff_SummaryFormat(t *testing.T) {
	path := writeRailwayEnvFile(t, "FOO=bar\n")
	var buf bytes.Buffer
	err := RunRailwayDiff(RailwayDiffOptions{
		LocalFile:   path,
		ProjectID:   "proj-1",
		Environment: "production",
		RemoteVars:  map[string]string{"FOO": "different"},
		Format:      "summary",
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected summary output, got empty")
	}
}

func TestRunRailwayDiff_InvalidLocalFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunRailwayDiff(RailwayDiffOptions{
		LocalFile: "/nonexistent/.env",
		ProjectID: "proj-1",
	}, &buf)
	if err == nil {
		t.Error("expected error for invalid local file")
	}
}

func TestRunRailwayDiff_EmptyProjectID(t *testing.T) {
	path := writeRailwayEnvFile(t, "FOO=bar\n")
	var buf bytes.Buffer
	err := RunRailwayDiff(RailwayDiffOptions{
		LocalFile: path,
		ProjectID: "",
	}, &buf)
	if err == nil {
		t.Error("expected error for empty project ID")
	}
}
