package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-sync/internal/store"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnvFile: %v", err)
	}
	return path
}

func TestNewFileStore_LoadsVars(t *testing.T) {
	path := writeTempEnvFile(t, "FOO=bar\nBAZ=qux\n")
	fs, err := store.NewFileStore(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := fs.Get("FOO"); !ok || v != "bar" {
		t.Errorf("expected FOO=bar, got %q (ok=%v)", v, ok)
	}
}

func TestNewFileStore_InvalidPath(t *testing.T) {
	_, err := store.NewFileStore("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestFileStore_GetMissingKey(t *testing.T) {
	path := writeTempEnvFile(t, "FOO=bar\n")
	fs, _ := store.NewFileStore(path)
	_, ok := fs.Get("MISSING")
	if ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestFileStore_SetAndGet(t *testing.T) {
	path := writeTempEnvFile(t, "FOO=bar\n")
	fs, _ := store.NewFileStore(path)
	fs.Set("NEW_KEY", "new_value")
	if v, ok := fs.Get("NEW_KEY"); !ok || v != "new_value" {
		t.Errorf("expected NEW_KEY=new_value, got %q (ok=%v)", v, ok)
	}
}

func TestFileStore_Delete(t *testing.T) {
	path := writeTempEnvFile(t, "FOO=bar\nBAZ=qux\n")
	fs, _ := store.NewFileStore(path)
	fs.Delete("FOO")
	_, ok := fs.Get("FOO")
	if ok {
		t.Error("expected FOO to be deleted")
	}
}

func TestFileStore_AllReturnsCopy(t *testing.T) {
	path := writeTempEnvFile(t, "FOO=bar\nBAZ=qux\n")
	fs, _ := store.NewFileStore(path)
	all := fs.All()
	all["FOO"] = "mutated"
	if v, _ := fs.Get("FOO"); v != "bar" {
		t.Errorf("All() should return a copy; original mutated to %q", v)
	}
}

func TestFileStore_Path(t *testing.T) {
	path := writeTempEnvFile(t, "")
	fs, _ := store.NewFileStore(path)
	if fs.Path() != path {
		t.Errorf("expected path %q, got %q", path, fs.Path())
	}
}
