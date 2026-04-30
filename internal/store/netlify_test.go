package store

import "testing"

func TestNewNetlifyStore_CopiesVars(t *testing.T) {
	orig := map[string]string{"A": "1"}
	s := NewNetlifyStore(orig, "site-abc", "production")
	orig["A"] = "mutated"
	if v, _ := s.Get("A"); v != "1" {
		t.Errorf("expected '1', got %q", v)
	}
}

func TestNetlifyStore_Name(t *testing.T) {
	s := NewNetlifyStore(nil, "site-xyz", "")
	if s.Name() != "netlify:site-xyz" {
		t.Errorf("unexpected name: %s", s.Name())
	}
}

func TestNetlifyStore_NameNoSite(t *testing.T) {
	s := NewNetlifyStore(nil, "", "")
	if s.Name() != "netlify" {
		t.Errorf("unexpected name: %s", s.Name())
	}
}

func TestNetlifyStore_Namespace(t *testing.T) {
	s := NewNetlifyStore(nil, "s", "deploy-preview")
	if s.Namespace() != "deploy-preview" {
		t.Errorf("unexpected namespace: %s", s.Namespace())
	}
}

func TestNetlifyStore_GetMissingKey(t *testing.T) {
	s := NewNetlifyStore(nil, "", "")
	_, ok := s.Get("MISSING")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestNetlifyStore_SetAndGet(t *testing.T) {
	s := NewNetlifyStore(nil, "", "")
	if err := s.Set("FOO", "bar"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := s.Get("FOO"); !ok || v != "bar" {
		t.Errorf("expected 'bar', got %q", v)
	}
}

func TestNetlifyStore_SetEmptyKey(t *testing.T) {
	s := NewNetlifyStore(nil, "", "")
	if err := s.Set("", "val"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestNetlifyStore_Delete(t *testing.T) {
	s := NewNetlifyStore(map[string]string{"X": "1"}, "", "")
	s.Delete("X")
	if _, ok := s.Get("X"); ok {
		t.Error("expected key to be deleted")
	}
}

func TestNetlifyStore_All(t *testing.T) {
	s := NewNetlifyStore(map[string]string{"A": "1", "B": "2"}, "", "")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 keys, got %d", len(all))
	}
	all["A"] = "mutated"
	if v, _ := s.Get("A"); v != "1" {
		t.Error("All() should return a copy")
	}
}
