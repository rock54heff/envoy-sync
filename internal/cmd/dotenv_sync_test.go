package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeDotenvSyncFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeDotenvSyncFile: %v", err)
	}
	return p
}

func TestRunDotenvSync_AlreadyInSync(t *testing.T) {
	dir := t.TempDir()
	base := writeDotenvSyncFile(t, dir, "base.env", "KEY=value\n")
	target := writeDotenvSyncFile(t, dir, "target.env", "KEY=value\n")

	var buf bytes.Buffer
	err := RunDotenvSync(DotenvSyncOptions{
		BaseFile:   base,
		TargetFile: target,
		Namespace:  "test",
		DryRun:     false,
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Already in sync") {
		t.Errorf("expected 'Already in sync', got: %s", buf.String())
	}
}

func TestRunDotenvSync_DryRun(t *testing.T) {
	dir := t.TempDir()
	base := writeDotenvSyncFile(t, dir, "base.env", "KEY=newval\n")
	target := writeDotenvSyncFile(t, dir, "target.env", "KEY=oldval\n")

	var buf bytes.Buffer
	err := RunDotenvSync(DotenvSyncOptions{
		BaseFile:   base,
		TargetFile: target,
		Namespace:  "test",
		DryRun:     true,
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run notice, got: %s", out)
	}
	// target file must be unchanged
	raw, _ := os.ReadFile(target)
	if !strings.Contains(string(raw), "oldval") {
		t.Errorf("target file should not be modified in dry-run")
	}
}

func TestRunDotenvSync_InvalidBaseFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunDotenvSync(DotenvSyncOptions{
		BaseFile:  "/nonexistent/base.env",
		Namespace: "test",
		Out:       &buf,
	})
	if err == nil {
		t.Fatal("expected error for missing base file")
	}
	if !strings.Contains(err.Error(), "loading base file") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunDotenvSync_WritesChanges(t *testing.T) {
	dir := t.TempDir()
	base := writeDotenvSyncFile(t, dir, "base.env", "KEY=newval\nEXTRA=added\n")
	target := writeDotenvSyncFile(t, dir, "target.env", "KEY=oldval\n")

	var buf bytes.Buffer
	err := RunDotenvSync(DotenvSyncOptions{
		BaseFile:   base,
		TargetFile: target,
		Namespace:  "prod",
		DryRun:     false,
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw, _ := os.ReadFile(target)
	contents := string(raw)
	if !strings.Contains(contents, "newval") {
		t.Errorf("expected updated value in target, got: %s", contents)
	}
	if !strings.Contains(contents, "EXTRA") {
		t.Errorf("expected new key in target, got: %s", contents)
	}
}
