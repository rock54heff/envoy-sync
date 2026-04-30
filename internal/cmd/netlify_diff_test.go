package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeNetlifyEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeNetlifyEnvFile: %v", err)
	}
	return p
}

func TestRunNetlifyDiff_NoChanges(t *testing.T) {
	p := writeNetlifyEnvFile(t, "FOO=bar\nBAZ=qux\n")
	var buf bytes.Buffer
	err := RunNetlifyDiff(NetlifyDiffOptions{
		LocalFile:  p,
		SiteID:     "site-1",
		RemoteVars: map[string]string{"FOO": "bar", "BAZ": "qux"},
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "+") || strings.Contains(buf.String(), "-") {
		t.Errorf("expected no diff output, got: %s", buf.String())
	}
}

func TestRunNetlifyDiff_WithChanges(t *testing.T) {
	p := writeNetlifyEnvFile(t, "FOO=newval\nNEW=key\n")
	var buf bytes.Buffer
	err := RunNetlifyDiff(NetlifyDiffOptions{
		LocalFile:  p,
		SiteID:     "site-1",
		RemoteVars: map[string]string{"FOO": "oldval", "OLD": "gone"},
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in diff, got: %s", out)
	}
}

func TestRunNetlifyDiff_SummaryFormat(t *testing.T) {
	p := writeNetlifyEnvFile(t, "FOO=bar\n")
	var buf bytes.Buffer
	err := RunNetlifyDiff(NetlifyDiffOptions{
		LocalFile:  p,
		SiteID:     "site-1",
		RemoteVars: map[string]string{"FOO": "baz"},
		Summary:    true,
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected summary output")
	}
}

func TestRunNetlifyDiff_InvalidLocalFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunNetlifyDiff(NetlifyDiffOptions{
		LocalFile:  "/nonexistent/.env",
		SiteID:     "site-1",
		RemoteVars: map[string]string{},
		Out:        &buf,
	})
	if err == nil {
		t.Error("expected error for invalid file")
	}
}

func TestRunNetlifyDiff_EmptySiteID(t *testing.T) {
	p := writeNetlifyEnvFile(t, "FOO=bar\n")
	var buf bytes.Buffer
	err := RunNetlifyDiff(NetlifyDiffOptions{
		LocalFile:  p,
		SiteID:     "",
		RemoteVars: map[string]string{},
		Out:        &buf,
	})
	if err == nil {
		t.Error("expected error for empty site ID")
	}
}
