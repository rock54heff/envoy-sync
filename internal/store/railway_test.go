package store

import (
	"testing"
)

func TestNewRailwayStore_CopiesVars(t *testing.T) {
	orig := map[string]string{"A": "1", "B": "2"}
	s := NewRailwayStore(orig, "proj-123", "production")
	orig["A"] = "mutated"
	v, _ := s.Get("A")
	if v != "1" {
		t.Errorf("expected '1', got %q", v)
	}
}

func TestRailwayStore_Name(t *testing.T) {
	s := NewRailwayStore(nil, "proj-abc", "")
	if s.Name() != "railway:proj-abc" {
		t.Errorf("unexpected name: %s", s.Name())
	}
}

func TestRailwayStore_NameNoProject(t *testing.T) {
	s := NewRailwayStore(nil, "", "")
	if s.Name() != "railway" {
		t.Errorf("unexpected name: %s", s.Name())
	}
}

func TestRailwayStore_Namespace(t *testing.T) {
	s := NewRailwayStore(nil, "proj-abc", "staging")
	expected := "railway:proj-abc/staging"
	if s.Namespace() != expected {
		t.Errorf("expected %q, got %q", expected, s.Namespace())
	}
}

func TestRailwayStore_GetMissingKey(t *testing.T) {
	s := NewRailwayStore(map[string]string{}, "p", "e")
	_, err := s.Get("MISSING")
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestRailwayStore_SetAndGet(t *testing.T) {
	s := NewRailwayStore(map[string]string{}, "p", "e")
	if err := s.Set("FOO", "bar"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := s.Get("FOO")
	if err != nil || v != "bar" {
		t.Errorf("expected 'bar', got %q (err: %v)", v, err)
	}
}

func TestRailwayStore_SetEmptyKey(t *testing.T) {
	s := NewRailwayStore(map[string]string{}, "p", "e")
	if err := s.Set("", "value"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestRailwayStore_Delete(t *testing.T) {
	s := NewRailwayStore(map[string]string{"X": "1"}, "p", "e")
	_ = s.Delete("X")
	_, err := s.Get("X")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestRailwayStore_All(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	s := NewRailwayStore(vars, "p", "e")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 keys, got %d", len(all))
	}
}
