package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeEnvFile: %v", err)
	}
	return p
}

func TestRunSync_AlreadyInSync(t *testing.T) {
	dir := t.TempDir()
	base := writeEnvFile(t, dir, "base.env", "FOO=bar\nBAZ=qux\n")
	target := writeEnvFile(t, dir, "target.env", "FOO=bar\nBAZ=qux\n")

	var buf bytes.Buffer
	err := RunSync(SyncOptions{BaseFile: base, TargetFile: target, Out: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "already in sync") {
		t.Errorf("expected 'already in sync', got: %s", buf.String())
	}
}

func TestRunSync_DryRun(t *testing.T) {
	dir := t.TempDir()
	base := writeEnvFile(t, dir, "base.env", "FOO=bar\nNEW=value\n")
	target := writeEnvFile(t, dir, "target.env", "FOO=bar\n")

	var buf bytes.Buffer
	err := RunSync(SyncOptions{BaseFile: base, TargetFile: target, DryRun: true, Out: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run label, got: %s", out)
	}
	// target file must remain unchanged
	raw, _ := os.ReadFile(target)
	if strings.Contains(string(raw), "NEW") {
		t.Error("dry-run must not write to target file")
	}
}

func TestRunSync_InvalidBaseFile(t *testing.T) {
	dir := t.TempDir()
	target := writeEnvFile(t, dir, "target.env", "FOO=bar\n")
	err := RunSync(SyncOptions{BaseFile: "/nonexistent/base.env", TargetFile: target})
	if err == nil {
		t.Fatal("expected error for missing base file")
	}
}

func TestRunSync_InvalidTargetFile(t *testing.T) {
	dir := t.TempDir()
	base := writeEnvFile(t, dir, "base.env", "FOO=bar\n")
	err := RunSync(SyncOptions{BaseFile: base, TargetFile: "/nonexistent/target.env"})
	if err == nil {
		t.Fatal("expected error for missing target file")
	}
}

func TestRunSync_VerboseOutput(t *testing.T) {
	dir := t.TempDir()
	base := writeEnvFile(t, dir, "base.env", "FOO=bar\nNEW=hello\n")
	target := writeEnvFile(t, dir, "target.env", "FOO=bar\n")

	var buf bytes.Buffer
	err := RunSync(SyncOptions{BaseFile: base, TargetFile: target, DryRun: true, Verbose: true, Out: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "NEW") {
		t.Errorf("verbose output should mention NEW key, got: %s", buf.String())
	}
}
