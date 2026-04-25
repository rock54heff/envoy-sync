package envfile

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "bar" || env["BAZ"] != "qux" {
		t.Errorf("got %v, want FOO=bar BAZ=qux", env)
	}
}

func TestParse_IgnoresCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=value\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 || env["KEY"] != "value" {
		t.Errorf("unexpected env: %v", env)
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SINGLE='hello world'
DOUBLE="goodbye world"
`)
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SINGLE"] != "hello world" {
		t.Errorf("SINGLE: got %q, want %q", env["SINGLE"], "hello world")
	}
	if env["DOUBLE"] != "goodbye world" {
		t.Errorf("DOUBLE: got %q, want %q", env["DOUBLE"], "goodbye world")
	}
}

func TestParse_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080 # http port\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["PORT"] != "8080" {
		t.Errorf("PORT: got %q, want %q", env["PORT"], "8080")
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "NOTAVALIDLINE\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParse_MissingFile(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParse_EmptyFile(t *testing.T) {
	path := writeTempEnv(t, "")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error for empty file: %v", err)
	}
	if len(env) != 0 {
		t.Errorf("expected empty map, got %v", env)
	}
}
