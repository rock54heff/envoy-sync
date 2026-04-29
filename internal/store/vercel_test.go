package store

import (
	"testing"
)

func TestNewVercelStore_CopiesVars(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	s := NewVercelStore(src, "my-project", "")
	src["FOO"] = "mutated"
	if v, _ := s.Get("FOO"); v != "bar" {
		t.Errorf("expected original value, got %q", v)
	}
}

func TestVercelStore_Name(t *testing.T) {
	s := NewVercelStore(nil, "my-project", "")
	if s.Name() != "vercel:my-project" {
		t.Errorf("unexpected name: %s", s.Name())
	}
}

func TestVercelStore_NameWithTeam(t *testing.T) {
	s := NewVercelStore(nil, "my-project", "acme")
	if s.Name() != "vercel:acme/my-project" {
		t.Errorf("unexpected name: %s", s.Name())
	}
}

func TestVercelStore_Namespace(t *testing.T) {
	s := NewVercelStore(nil, "proj", "team")
	if s.Namespace() != "proj" {
		t.Errorf("expected 'proj', got %q", s.Namespace())
	}
}

func TestVercelStore_GetMissingKey(t *testing.T) {
	s := NewVercelStore(nil, "proj", "")
	_, ok := s.Get("MISSING")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestVercelStore_SetAndGet(t *testing.T) {
	s := NewVercelStore(nil, "proj", "")
	if err := s.Set("KEY", "value"); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	v, ok := s.Get("KEY")
	if !ok || v != "value" {
		t.Errorf("expected 'value', got %q ok=%v", v, ok)
	}
}

func TestVercelStore_SetEmptyKey(t *testing.T) {
	s := NewVercelStore(nil, "proj", "")
	if err := s.Set("", "val"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestVercelStore_Delete(t *testing.T) {
	s := NewVercelStore(map[string]string{"X": "1"}, "proj", "")
	if err := s.Delete("X"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	_, ok := s.Get("X")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestVercelStore_All(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	s := NewVercelStore(vars, "proj", "")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
	all["A"] = "mutated"
	if v, _ := s.Get("A"); v != "1" {
		t.Error("All() should return a copy")
	}
}
