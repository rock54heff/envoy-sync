package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFlyioSyncFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeFlyioSyncFile: %v", err)
	}
	return p
}

func TestRunFlyioSync_AlreadyInSync(t *testing.T) {
	path := writeFlyioSyncFile(t, "KEY=value\n")
	var buf bytes.Buffer
	err := RunFlyioSync(FlyioSyncOptions{
		BaseFile: path,
		AppName:  "my-app",
		Token:    "tok",
		Out:      &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "already in sync") {
		t.Errorf("expected 'already in sync', got: %q", buf.String())
	}
}

func TestRunFlyioSync_DryRun(t *testing.T) {
	path := writeFlyioSyncFile(t, "NEW_KEY=hello\nANOTHER=world\n")
	var buf bytes.Buffer
	err := RunFlyioSync(FlyioSyncOptions{
		BaseFile: path,
		AppName:  "my-app",
		Token:    "tok",
		DryRun:   true,
		Out:      &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "dry-run") {
		t.Errorf("expected dry-run output, got: %q", buf.String())
	}
}

func TestRunFlyioSync_InvalidBaseFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunFlyioSync(FlyioSyncOptions{
		BaseFile: "/nonexistent/.env",
		AppName:  "my-app",
		Token:    "tok",
		Out:      &buf,
	})
	if err == nil {
		t.Fatal("expected error for invalid base file")
	}
	if !strings.Contains(err.Error(), "loading base file") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunFlyioSync_EmptyAppName(t *testing.T) {
	path := writeFlyioSyncFile(t, "KEY=value\n")
	var buf bytes.Buffer
	err := RunFlyioSync(FlyioSyncOptions{
		BaseFile: path,
		AppName:  "",
		Token:    "tok",
		Out:      &buf,
	})
	if err == nil {
		t.Fatal("expected error for empty app name")
	}
	if !strings.Contains(err.Error(), "app name is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunFlyioSync_EmptyBaseFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunFlyioSync(FlyioSyncOptions{
		BaseFile: "",
		AppName:  "my-app",
		Token:    "tok",
		Out:      &buf,
	})
	if err == nil {
		t.Fatal("expected error for empty base file")
	}
	if !strings.Contains(err.Error(), "base file path is required") {
		t.Errorf("unexpected error: %v", err)
	}
}
