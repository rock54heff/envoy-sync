package store

import (
	"testing"
)

func TestNewFlyioStore_CopiesVars(t *testing.T) {
	orig := map[string]string{"A": "1", "B": "2"}
	s := NewFlyioStore("my-app", orig)
	orig["A"] = "mutated"
	v, err := s.Get("A")
	if err != nil || v != "1" {
		t.Fatalf("expected '1', got %q (err=%v)", v, err)
	}
}

func TestFlyioStore_Name(t *testing.T) {
	s := NewFlyioStore("my-app", nil)
	if s.Name() != "fly.io(my-app)" {
		t.Errorf("unexpected name: %s", s.Name())
	}
}

func TestFlyioStore_NameNoApp(t *testing.T) {
	s := NewFlyioStore("", nil)
	if s.Name() != "fly.io" {
		t.Errorf("unexpected name: %s", s.Name())
	}
}

func TestFlyioStore_Namespace(t *testing.T) {
	s := NewFlyioStore("prod", nil)
	if s.Namespace() != "flyio/prod" {
		t.Errorf("unexpected namespace: %s", s.Namespace())
	}
}

func TestFlyioStore_GetMissingKey(t *testing.T) {
	s := NewFlyioStore("app", nil)
	_, err := s.Get("MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestFlyioStore_SetAndGet(t *testing.T) {
	s := NewFlyioStore("app", nil)
	if err := s.Set("KEY", "val"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := s.Get("KEY")
	if err != nil || v != "val" {
		t.Fatalf("expected 'val', got %q (err=%v)", v, err)
	}
}

func TestFlyioStore_SetEmptyKey(t *testing.T) {
	s := NewFlyioStore("app", nil)
	if err := s.Set("", "val"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestFlyioStore_Delete(t *testing.T) {
	s := NewFlyioStore("app", map[string]string{"X": "1"})
	_ = s.Delete("X")
	_, err := s.Get("X")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestFlyioStore_ToMap(t *testing.T) {
	s := NewFlyioStore("app", map[string]string{"A": "1"})
	m := s.ToMap()
	if m["A"] != "1" {
		t.Errorf("expected '1', got %q", m["A"])
	}
	m["A"] = "mutated"
	v, _ := s.Get("A")
	if v != "1" {
		t.Error("ToMap should return a copy")
	}
}
