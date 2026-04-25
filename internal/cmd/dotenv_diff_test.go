package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeDotenvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeDotenvFile: %v", err)
	}
	return p
}

func TestRunDotenvDiff_NoChanges(t *testing.T) {
	dir := t.TempDir()
	content := "KEY=value\nFOO=bar\n"
	base := writeDotenvFile(t, dir, "base.env", content)
	target := writeDotenvFile(t, dir, "target.env", content)

	var buf bytes.Buffer
	err := RunDotenvDiff(DotenvDiffOptions{BaseFile: base, TargetFile: target, Out: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "+") || strings.Contains(buf.String(), "-") {
		t.Errorf("expected no diff output, got: %s", buf.String())
	}
}

func TestRunDotenvDiff_WithChanges(t *testing.T) {
	dir := t.TempDir()
	base := writeDotenvFile(t, dir, "base.env", "KEY=old\n")
	target := writeDotenvFile(t, dir, "target.env", "KEY=new\nEXTRA=yes\n")

	var buf bytes.Buffer
	err := RunDotenvDiff(DotenvDiffOptions{BaseFile: base, TargetFile: target, Out: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected KEY in diff output, got: %s", out)
	}
}

func TestRunDotenvDiff_SummaryFormat(t *testing.T) {
	dir := t.TempDir()
	base := writeDotenvFile(t, dir, "base.env", "A=1\n")
	target := writeDotenvFile(t, dir, "target.env", "A=1\nB=2\n")

	var buf bytes.Buffer
	err := RunDotenvDiff(DotenvDiffOptions{BaseFile: base, TargetFile: target, Summary: true, Out: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "added") && !strings.Contains(buf.String(), "1") {
		t.Errorf("expected summary output, got: %s", buf.String())
	}
}

func TestRunDotenvDiff_InvalidBaseFile(t *testing.T) {
	dir := t.TempDir()
	target := writeDotenvFile(t, dir, "target.env", "KEY=val\n")
	err := RunDotenvDiff(DotenvDiffOptions{BaseFile: "/nonexistent/base.env", TargetFile: target})
	if err == nil {
		t.Error("expected error for invalid base file")
	}
}

func TestRunDotenvDiff_InvalidTargetFile(t *testing.T) {
	dir := t.TempDir()
	base := writeDotenvFile(t, dir, "base.env", "KEY=val\n")
	err := RunDotenvDiff(DotenvDiffOptions{BaseFile: base, TargetFile: "/nonexistent/target.env"})
	if err == nil {
		t.Error("expected error for invalid target file")
	}
}
