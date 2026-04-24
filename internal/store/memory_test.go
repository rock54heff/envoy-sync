package store

import (
	"testing"
)

func TestNew_CopiesVars(t *testing.T) {
	orig := map[string]string{"A": "1", "B": "2"}
	s := New("test", orig)
	orig["A"] = "mutated"
	if s.Vars["A"] != "1" {
		t.Errorf("expected '1', got %q", s.Vars["A"])
	}
}

func TestGet_ExistingKey(t *testing.T) {
	s := New("test", map[string]string{"FOO": "bar"})
	v, ok := s.Get("FOO")
	if !ok || v != "bar" {
		t.Errorf("expected ('bar', true), got (%q, %v)", v, ok)
	}
}

func TestGet_MissingKey(t *testing.T) {
	s := New("test", map[string]string{})
	_, ok := s.Get("MISSING")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestSet_AddsKey(t *testing.T) {
	s := New("test", map[string]string{})
	s.Set("NEW", "value")
	v, ok := s.Get("NEW")
	if !ok || v != "value" {
		t.Errorf("expected ('value', true), got (%q, %v)", v, ok)
	}
}

func TestDelete_RemovesKey(t *testing.T) {
	s := New("test", map[string]string{"DEL": "bye"})
	s.Delete("DEL")
	_, ok := s.Get("DEL")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestKeys_SortedOrder(t *testing.T) {
	s := New("test", map[string]string{"Z": "1", "A": "2", "M": "3"})
	keys := s.Keys()
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestMerge_OtherTakesPrecedence(t *testing.T) {
	a := New("a", map[string]string{"X": "from_a", "Y": "only_a"})
	b := New("b", map[string]string{"X": "from_b", "Z": "only_b"})
	merged := a.Merge(b)

	if merged.Vars["X"] != "from_b" {
		t.Errorf("expected 'from_b', got %q", merged.Vars["X"])
	}
	if merged.Vars["Y"] != "only_a" {
		t.Errorf("expected 'only_a', got %q", merged.Vars["Y"])
	}
	if merged.Vars["Z"] != "only_b" {
		t.Errorf("expected 'only_b', got %q", merged.Vars["Z"])
	}
	if merged.Name != "a+b" {
		t.Errorf("expected name 'a+b', got %q", merged.Name)
	}
}

func TestLen(t *testing.T) {
	s := New("test", map[string]string{"A": "1", "B": "2"})
	if s.Len() != 2 {
		t.Errorf("expected 2, got %d", s.Len())
	}
}
