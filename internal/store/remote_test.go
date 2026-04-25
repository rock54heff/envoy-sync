package store

import (
	"testing"
)

func TestNewRemoteStore_CopiesVars(t *testing.T) {
	orig := map[string]string{"KEY": "val"}
	rs := NewRemoteStore("prod", orig)
	orig["KEY"] = "mutated"
	v, err := rs.Get("KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "val" {
		t.Errorf("expected 'val', got %q", v)
	}
}

func TestRemoteStore_Namespace(t *testing.T) {
	rs := NewRemoteStore("staging", nil)
	if rs.Namespace() != "staging" {
		t.Errorf("expected 'staging', got %q", rs.Namespace())
	}
}

func TestRemoteStore_GetMissingKey(t *testing.T) {
	rs := NewRemoteStore("prod", nil)
	_, err := rs.Get("MISSING")
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestRemoteStore_SetAndGet(t *testing.T) {
	rs := NewRemoteStore("prod", nil)
	if err := rs.Set("DB_URL", "postgres://localhost"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	v, err := rs.Get("DB_URL")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if v != "postgres://localhost" {
		t.Errorf("expected 'postgres://localhost', got %q", v)
	}
}

func TestRemoteStore_Delete(t *testing.T) {
	rs := NewRemoteStore("prod", map[string]string{"TOKEN": "abc"})
	if err := rs.Delete("TOKEN"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err := rs.Get("TOKEN")
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestRemoteStore_DeleteMissingKey(t *testing.T) {
	rs := NewRemoteStore("prod", nil)
	if err := rs.Delete("GHOST"); err == nil {
		t.Error("expected error deleting missing key, got nil")
	}
}

func TestRemoteStore_All(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	rs := NewRemoteStore("dev", vars)
	all := rs.All()
	if len(all) != 2 {
		t.Errorf("expected 2 keys, got %d", len(all))
	}
	all["A"] = "mutated"
	v, _ := rs.Get("A")
	if v != "1" {
		t.Error("All() should return a copy, not a reference")
	}
}

func TestRemoteStore_Keys_Sorted(t *testing.T) {
	vars := map[string]string{"Z": "1", "A": "2", "M": "3"}
	rs := NewRemoteStore("dev", vars)
	keys := rs.Keys()
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("keys[%d]: expected %q, got %q", i, expected[i], k)
		}
	}
}
