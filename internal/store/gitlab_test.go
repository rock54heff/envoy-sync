package store

import "testing"

func TestNewGitLabStore_CopiesVars(t *testing.T) {
	orig := map[string]string{"A": "1", "B": "2"}
	s := NewGitLabStore("myproject", orig)
	orig["A"] = "mutated"
	if v, _ := s.Get("A"); v != "1" {
		t.Errorf("expected '1', got %q", v)
	}
}

func TestGitLabStore_Name(t *testing.T) {
	s := NewGitLabStore("proj", nil)
	if s.Name() != "gitlab" {
		t.Errorf("expected 'gitlab', got %q", s.Name())
	}
}

func TestGitLabStore_Namespace(t *testing.T) {
	s := NewGitLabStore("my-group/my-project", nil)
	if s.Namespace() != "my-group/my-project" {
		t.Errorf("unexpected namespace: %q", s.Namespace())
	}
}

func TestGitLabStore_GetMissingKey(t *testing.T) {
	s := NewGitLabStore("proj", nil)
	_, ok := s.Get("MISSING")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestGitLabStore_SetAndGet(t *testing.T) {
	s := NewGitLabStore("proj", nil)
	if err := s.Set("KEY", "value"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get("KEY")
	if !ok || v != "value" {
		t.Errorf("expected 'value', got %q (ok=%v)", v, ok)
	}
}

func TestGitLabStore_SetEmptyKey(t *testing.T) {
	s := NewGitLabStore("proj", nil)
	if err := s.Set("", "val"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestGitLabStore_Delete(t *testing.T) {
	s := NewGitLabStore("proj", map[string]string{"X": "1"})
	_ = s.Delete("X")
	_, ok := s.Get("X")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestGitLabStore_All(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	s := NewGitLabStore("proj", vars)
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 keys, got %d", len(all))
	}
	all["A"] = "mutated"
	if v, _ := s.Get("A"); v != "1" {
		t.Error("All() should return a copy")
	}
}
