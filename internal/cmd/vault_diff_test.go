package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunVaultDiff_NoChanges(t *testing.T) {
	f := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	var buf bytes.Buffer
	err := RunVaultDiff(VaultDiffOptions{
		LocalFile:  f,
		Namespace:  "ns",
		RemoteVars: map[string]string{"ns/KEY": "value", "ns/FOO": "bar"},
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "+") || strings.Contains(buf.String(), "-") {
		t.Fatalf("expected no diff output, got:\n%s", buf.String())
	}
}

func TestRunVaultDiff_WithChanges(t *testing.T) {
	f := writeTempEnv(t, "KEY=local\n")
	var buf bytes.Buffer
	err := RunVaultDiff(VaultDiffOptions{
		LocalFile:  f,
		Namespace:  "ns",
		RemoteVars: map[string]string{"ns/KEY": "remote", "ns/EXTRA": "only_vault"},
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "KEY") {
		t.Errorf("expected KEY in diff output")
	}
	if !strings.Contains(output, "EXTRA") {
		t.Errorf("expected EXTRA in diff output")
	}
}

func TestRunVaultDiff_SummaryFormat(t *testing.T) {
	f := writeTempEnv(t, "A=1\nB=2\n")
	var buf bytes.Buffer
	err := RunVaultDiff(VaultDiffOptions{
		LocalFile:  f,
		Namespace:  "ns",
		RemoteVars: map[string]string{"ns/A": "1", "ns/B": "changed"},
		Summary:    true,
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "modified") && !strings.Contains(buf.String(), "1") {
		t.Errorf("summary should mention modified count, got: %s", buf.String())
	}
}

func TestRunVaultDiff_InvalidLocalFile(t *testing.T) {
	err := RunVaultDiff(VaultDiffOptions{
		LocalFile:  "/nonexistent/path/.env",
		Namespace:  "ns",
		RemoteVars: map[string]string{},
	})
	if err == nil {
		t.Fatal("expected error for invalid local file")
	}
}

func TestRunVaultDiff_EmptyNamespace(t *testing.T) {
	f := writeTempEnv(t, "KEY=val\n")
	var buf bytes.Buffer
	err := RunVaultDiff(VaultDiffOptions{
		LocalFile:  f,
		Namespace:  "",
		RemoteVars: map[string]string{"KEY": "val"},
		Out:        &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
