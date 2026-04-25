package store

import (
	"testing"
)

func TestNewVaultStore_CopiesVars(t *testing.T) {
	orig := map[string]string{"ns/KEY": "val"}
	vs := NewVaultStore("ns", orig)
	orig["ns/KEY"] = "mutated"
	if v, _ := vs.Get("KEY"); v != "val" {
		t.Fatalf("expected original value, got %q", v)
	}
}

func TestVaultStore_Namespace(t *testing.T) {
	vs := NewVaultStore("prod/app", nil)
	if vs.Namespace() != "prod/app" {
		t.Fatalf("unexpected namespace %q", vs.Namespace())
	}
}

func TestVaultStore_GetMissingKey(t *testing.T) {
	vs := NewVaultStore("ns", nil)
	if _, ok := vs.Get("MISSING"); ok {
		t.Fatal("expected missing key to return ok=false")
	}
}

func TestVaultStore_SetAndGet(t *testing.T) {
	vs := NewVaultStore("ns", nil)
	if err := vs.Set("FOO", "bar"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := vs.Get("FOO"); !ok || v != "bar" {
		t.Fatalf("expected bar, got %q ok=%v", v, ok)
	}
}

func TestVaultStore_SetEmptyKey(t *testing.T) {
	vs := NewVaultStore("ns", nil)
	if err := vs.Set("", "val"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestVaultStore_Delete(t *testing.T) {
	vs := NewVaultStore("ns", map[string]string{"ns/KEY": "val"})
	if err := vs.Delete("KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := vs.Get("KEY"); ok {
		t.Fatal("key should be deleted")
	}
}

func TestVaultStore_DeleteMissing(t *testing.T) {
	vs := NewVaultStore("ns", nil)
	if err := vs.Delete("GHOST"); err == nil {
		t.Fatal("expected error deleting missing key")
	}
}

func TestVaultStore_All_StripsNamespace(t *testing.T) {
	vs := NewVaultStore("ns", map[string]string{"ns/A": "1", "ns/B": "2"})
	all := vs.All()
	if all["A"] != "1" || all["B"] != "2" {
		t.Fatalf("unexpected All() result: %v", all)
	}
	if _, ok := all["ns/A"]; ok {
		t.Fatal("All() should strip namespace prefix")
	}
}

func TestVaultStore_SortedKeys(t *testing.T) {
	vs := NewVaultStore("ns", map[string]string{"ns/Z": "z", "ns/A": "a", "ns/M": "m"})
	keys := vs.SortedKeys()
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Fatalf("index %d: want %q got %q", i, expected[i], k)
		}
	}
}
