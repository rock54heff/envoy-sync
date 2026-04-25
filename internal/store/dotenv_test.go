package store

import (
	"testing"
)

func TestNewDotenvStore_CopiesVars(t *testing.T) {
	orig := map[string]string{"KEY": "val"}
	s := NewDotenvStore("test", "", orig)
	orig["KEY"] = "mutated"
	if v, _ := s.Get("KEY"); v != "val" {
		t.Errorf("expected 'val', got %q", v)
	}
}

func TestDotenvStore_Name(t *testing.T) {
	s := NewDotenvStore("myfile", "", nil)
	if s.Name() != "myfile" {
		t.Errorf("expected 'myfile', got %q", s.Name())
	}
}

func TestDotenvStore_Namespace(t *testing.T) {
	s := NewDotenvStore("x", "APP", map[string]string{"APP_PORT": "8080"})
	if s.Namespace() != "APP" {
		t.Errorf("expected 'APP', got %q", s.Namespace())
	}
	v, ok := s.Get("PORT")
	if !ok || v != "8080" {
		t.Errorf("expected '8080', got %q ok=%v", v, ok)
	}
}

func TestDotenvStore_GetMissingKey(t *testing.T) {
	s := NewDotenvStore("x", "", nil)
	_, ok := s.Get("MISSING")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestDotenvStore_SetAndGet(t *testing.T) {
	s := NewDotenvStore("x", "", nil)
	s.Set("FOO", "bar")
	v, ok := s.Get("FOO")
	if !ok || v != "bar" {
		t.Errorf("expected 'bar', got %q ok=%v", v, ok)
	}
}

func TestDotenvStore_SetEmptyKey(t *testing.T) {
	s := NewDotenvStore("x", "", nil)
	s.Set("  ", "val")
	if len(s.All()) != 0 {
		t.Error("expected empty key to be ignored")
	}
}

func TestDotenvStore_Delete(t *testing.T) {
	s := NewDotenvStore("x", "", map[string]string{"KEY": "v"})
	s.Delete("KEY")
	_, ok := s.Get("KEY")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestDotenvStore_All(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	s := NewDotenvStore("x", "", vars)
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
	// mutating All() result must not affect store
	all["A"] = "mutated"
	if v, _ := s.Get("A"); v != "1" {
		t.Errorf("All() returned non-copy, store mutated to %q", v)
	}
}
