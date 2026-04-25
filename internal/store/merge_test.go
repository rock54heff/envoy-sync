package store_test

import (
	"testing"

	"github.com/yourorg/envoy-sync/internal/store"
)

func TestMerge_BaseOnly(t *testing.T) {
	base := store.New(map[string]string{"A": "1", "B": "2"})
	override := store.New(map[string]string{})

	result := store.Merge(base, override)

	if v, ok := result.Get("A"); !ok || v != "1" {
		t.Errorf("expected A=1, got %q ok=%v", v, ok)
	}
	if v, ok := result.Get("B"); !ok || v != "2" {
		t.Errorf("expected B=2, got %q ok=%v", v, ok)
	}
}

func TestMerge_OverrideWins(t *testing.T) {
	base := store.New(map[string]string{"A": "base", "B": "base"})
	override := store.New(map[string]string{"A": "override"})

	result := store.Merge(base, override)

	if v, ok := result.Get("A"); !ok || v != "override" {
		t.Errorf("expected A=override, got %q ok=%v", v, ok)
	}
	if v, ok := result.Get("B"); !ok || v != "base" {
		t.Errorf("expected B=base, got %q ok=%v", v, ok)
	}
}

func TestMerge_OverrideAddsKeys(t *testing.T) {
	base := store.New(map[string]string{"A": "1"})
	override := store.New(map[string]string{"B": "2"})

	result := store.Merge(base, override)

	keys := store.Keys(result)
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestMerge_DoesNotMutateBase(t *testing.T) {
	base := store.New(map[string]string{"A": "original"})
	override := store.New(map[string]string{"A": "changed"})

	store.Merge(base, override)

	if v, ok := base.Get("A"); !ok || v != "original" {
		t.Errorf("base mutated: expected A=original, got %q ok=%v", v, ok)
	}
}

func TestKeys_ReturnsSortedKeys(t *testing.T) {
	s := store.New(map[string]string{"C": "3", "A": "1", "B": "2"})
	keys := store.Keys(s)

	expected := []string{"A", "B", "C"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("keys[%d]: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestKeys_Empty(t *testing.T) {
	s := store.New(map[string]string{})
	keys := store.Keys(s)
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}
