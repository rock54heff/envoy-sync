package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteEnvFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.env")

	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := WriteEnvFile(p, vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	content := string(raw)
	if !strings.Contains(content, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got:\n%s", content)
	}
	if !strings.Contains(content, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got:\n%s", content)
	}
}

func TestWriteEnvFile_QuotesSpaces(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.env")

	vars := map[string]string{"MSG": "hello world"}
	if err := WriteEnvFile(p, vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, _ := os.ReadFile(p)
	if !strings.Contains(string(raw), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", string(raw))
	}
}

func TestWriteEnvFile_EmptyPath(t *testing.T) {
	err := WriteEnvFile("", map[string]string{"K": "v"})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestWriteEnvFile_SortedOutput(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "sorted.env")

	vars := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"}
	if err := WriteEnvFile(p, vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, _ := os.ReadFile(p)
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected first line ALPHA, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected last line ZEBRA, got %s", lines[2])
	}
}

func TestNeedsQuoting(t *testing.T) {
	cases := []struct {
		v    string
		want bool
	}{
		{"simple", false},
		{"with space", true},
		{"with#hash", true},
		{"with\ttab", true},
		{"normal123", false},
	}
	for _, c := range cases {
		if got := needsQuoting(c.v); got != c.want {
			t.Errorf("needsQuoting(%q) = %v, want %v", c.v, got, c.want)
		}
	}
}
