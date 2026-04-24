package diff_test

import (
	"testing"

	"github.com/yourorg/envoy-sync/internal/diff"
)

func TestCompare_Added(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "NEW": "val"}

	res := diff.Compare(a, b)
	if len(res.Added) != 1 || res.Added["NEW"] != "val" {
		t.Errorf("expected NEW=val in Added, got %v", res.Added)
	}
	if len(res.Removed) != 0 || len(res.Modified) != 0 {
		t.Errorf("unexpected changes: removed=%v modified=%v", res.Removed, res.Modified)
	}
}

func TestCompare_Removed(t *testing.T) {
	a := map[string]string{"FOO": "bar", "OLD": "gone"}
	b := map[string]string{"FOO": "bar"}

	res := diff.Compare(a, b)
	if len(res.Removed) != 1 || res.Removed["OLD"] != "gone" {
		t.Errorf("expected OLD=gone in Removed, got %v", res.Removed)
	}
}

func TestCompare_Modified(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}

	res := diff.Compare(a, b)
	if len(res.Modified) != 1 || res.Modified["FOO"] != "new" {
		t.Errorf("expected FOO=new in Modified, got %v", res.Modified)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}

	res := diff.Compare(a, b)
	if res.HasChanges() {
		t.Errorf("expected no changes, got added=%v removed=%v modified=%v", res.Added, res.Removed, res.Modified)
	}
	if len(res.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged keys, got %d", len(res.Unchanged))
	}
}

func TestCompare_Empty(t *testing.T) {
	res := diff.Compare(map[string]string{}, map[string]string{})
	if res.HasChanges() {
		t.Error("expected no changes for empty maps")
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]string{"Z": "1", "A": "2", "M": "3"}
	keys := diff.SortedKeys(m)
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], k)
		}
	}
}
