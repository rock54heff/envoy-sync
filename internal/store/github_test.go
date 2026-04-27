package store

import (
	"testing"
)

func TestNewGitHubStore_CopiesVars(t *testing.T) {
	orig := map[string]string{"TOKEN": "abc", "SECRET": "xyz"}
	s := NewGitHubStore("org/repo", "production", orig)
	orig["TOKEN"] = "mutated"
	if v, _ := s.Get("TOKEN"); v != "abc" {
		t.Errorf("expected abc, got %s", v)
	}
}

func TestGitHubStore_Name(t *testing.T) {
	s := NewGitHubStore("org/repo", "prod", nil)
	if s.Name() != "github:org/repo" {
		t.Errorf("unexpected name: %s", s.Name())
	}
}

func TestGitHubStore_Namespace(t *testing.T) {
	s := NewGitHubStore("org/repo", "staging", nil)
	if s.Namespace() != "staging" {
		t.Errorf("expected staging, got %s", s.Namespace())
	}
}

func TestGitHubStore_GetMissingKey(t *testing.T) {
	s := NewGitHubStore("org/repo", "", nil)
	_, ok := s.Get("MISSING")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestGitHubStore_SetAndGet(t *testing.T) {
	s := NewGitHubStore("org/repo", "", nil)
	if err := s.Set("API_KEY", "secret"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get("API_KEY")
	if !ok || v != "secret" {
		t.Errorf("expected secret, got %s", v)
	}
}

func TestGitHubStore_SetEmptyKey(t *testing.T) {
	s := NewGitHubStore("org/repo", "", nil)
	if err := s.Set("", "value"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestGitHubStore_Delete(t *testing.T) {
	s := NewGitHubStore("org/repo", "", map[string]string{"KEY": "val"})
	_ = s.Delete("KEY")
	_, ok := s.Get("KEY")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestGitHubStore_Keys_Sorted(t *testing.T) {
	s := NewGitHubStore("org/repo", "", map[string]string{"Z": "1", "A": "2", "M": "3"})
	keys := s.Keys()
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("keys not sorted: %v", keys)
	}
}
