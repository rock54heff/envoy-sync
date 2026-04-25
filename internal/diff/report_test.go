package diff

import (
	"testing"
)

func TestSummarize_AllTypes(t *testing.T) {
	results := []Result{
		{Key: "A", Status: Added},
		{Key: "B", Status: Removed},
		{Key: "C", Status: Modified},
		{Key: "D", Status: Unchanged},
		{Key: "E", Status: Unchanged},
	}
	s := Summarize(results)
	if s.Added != 1 {
		t.Errorf("expected Added=1, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("expected Modified=1, got %d", s.Modified)
	}
	if s.Unchanged != 2 {
		t.Errorf("expected Unchanged=2, got %d", s.Unchanged)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize([]Result{})
	if s.HasChanges() {
		t.Error("expected no changes for empty results")
	}
}

func TestSummary_HasChanges_True(t *testing.T) {
	s := Summary{Added: 1}
	if !s.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestSummary_HasChanges_False(t *testing.T) {
	s := Summary{Unchanged: 3}
	if s.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestSummary_String_NoChanges(t *testing.T) {
	s := Summary{}
	if s.String() != "no differences" {
		t.Errorf("unexpected string: %q", s.String())
	}
}

func TestSummary_String_WithChanges(t *testing.T) {
	s := Summary{Added: 2, Removed: 1, Modified: 3, Unchanged: 4}
	got := s.String()
	for _, want := range []string{"2 added", "1 removed", "3 modified", "4 unchanged"} {
		if !contains(got, want) {
			t.Errorf("expected %q in summary string %q", want, got)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
