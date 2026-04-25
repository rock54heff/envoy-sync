package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeMergeEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeMergeEnvFile: %v", err)
	}
	return p
}

func TestRunMerge_DryRun(t *testing.T) {
	dir := t.TempDir()
	base := writeMergeEnvFile(t, dir, "base.env", "FOO=bar\nSHARED=base\n")
	override := writeMergeEnvFile(t, dir, "override.env", "BAZ=qux\nSHARED=override\n")

	var buf bytes.Buffer
	err := RunMerge(MergeOptions{
		BaseFile:     base,
		OverrideFile: override,
		DryRun:       true,
		Out:          &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "SHARED=override") {
		t.Errorf("expected override to win; got:\n%s", out)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected base key FOO; got:\n%s", out)
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected override key BAZ; got:\n%s", out)
	}
}

func TestRunMerge_WritesOutputFile(t *testing.T) {
	dir := t.TempDir()
	base := writeMergeEnvFile(t, dir, "base.env", "KEY=base\n")
	override := writeMergeEnvFile(t, dir, "override.env", "KEY=override\nEXTRA=yes\n")
	out := filepath.Join(dir, "merged.env")

	var buf bytes.Buffer
	err := RunMerge(MergeOptions{
		BaseFile:     base,
		OverrideFile: override,
		OutputFile:   out,
		Out:          &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	contents := string(data)
	if !strings.Contains(contents, "KEY=override") {
		t.Errorf("expected KEY=override in output; got:\n%s", contents)
	}
	if !strings.Contains(contents, "EXTRA=yes") {
		t.Errorf("expected EXTRA=yes in output; got:\n%s", contents)
	}
}

func TestRunMerge_InvalidBaseFile(t *testing.T) {
	dir := t.TempDir()
	override := writeMergeEnvFile(t, dir, "override.env", "K=v\n")
	err := RunMerge(MergeOptions{
		BaseFile:     filepath.Join(dir, "nonexistent.env"),
		OverrideFile: override,
		DryRun:       true,
		Out:          &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for missing base file")
	}
}

func TestRunMerge_EmptyOutputPath(t *testing.T) {
	dir := t.TempDir()
	base := writeMergeEnvFile(t, dir, "base.env", "K=v\n")
	override := writeMergeEnvFile(t, dir, "override.env", "K=v2\n")
	err := RunMerge(MergeOptions{
		BaseFile:     base,
		OverrideFile: override,
		OutputFile:   "",
		Out:          &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error when output file path is empty")
	}
}
