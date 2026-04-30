package store

import (
	"testing"
)

func TestNewRenderStore_CopiesVars(t *testing.T) {
	orig := map[string]string{"FOO": "bar"}
	s := NewRenderStore("svc-123", "key-abc", orig)
	orig["FOO"] = "mutated"
	v, ok := s.Get("FOO")
	if !ok || v != "bar" {
		t.Errorf("expected 'bar', got %q (ok=%v)", v, ok)
	}
}

func TestRenderStore_Name(t *testing.T) {
	s := NewRenderStore("svc-123", "key-abc", nil)
	if s.Name() != "render(svc-123)" {
		t.Errorf("unexpected name: %s", s.Name())
	}
}

func TestRenderStore_NameNoService(t *testing.T) {
	s := NewRenderStore("", "key-abc", nil)
	if s.Name() != "render" {
		t.Errorf("unexpected name: %s", s.Name())
	}
}

func TestRenderStore_Namespace(t *testing.T) {
	s := NewRenderStore("svc-123", "key-abc", nil)
	if s.Namespace() != "render/svc-123" {
		t.Errorf("unexpected namespace: %s", s.Namespace())
	}
}

func TestRenderStore_GetMissingKey(t *testing.T) {
	s := NewRenderStore("svc-123", "key-abc", nil)
	_, ok := s.Get("MISSING")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestRenderStore_SetAndGet(t *testing.T) {
	s := NewRenderStore("svc-123", "key-abc", nil)
	if err := s.Set("NEW_KEY", "new_val"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get("NEW_KEY")
	if !ok || v != "new_val" {
		t.Errorf("expected 'new_val', got %q", v)
	}
}

func TestRenderStore_SetEmptyKey(t *testing.T) {
	s := NewRenderStore("svc-123", "key-abc", nil)
	if err := s.Set("", "val"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestRenderStore_Delete(t *testing.T) {
	s := NewRenderStore("svc-123", "key-abc", map[string]string{"TO_DEL": "yes"})
	if err := s.Delete("TO_DEL"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("TO_DEL")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestRenderStore_SortedKeys(t *testing.T) {
	s := NewRenderStore("svc-123", "key-abc", map[string]string{
		"ZZZ": "1",
		"AAA": "2",
		"MMM": "3",
	})
	keys := s.SortedKeys()
	expected := []string{"AAA", "MMM", "ZZZ"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], k)
		}
	}
}
