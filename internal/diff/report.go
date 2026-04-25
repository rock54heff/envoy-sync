package diff

import (
	"fmt"
	"strings"
)

// Summary holds aggregate counts from a diff result.
type Summary struct {
	Added     int
	Removed   int
	Modified  int
	Unchanged int
}

// Summarize counts the number of each change type in a slice of Results.
func Summarize(results []Result) Summary {
	var s Summary
	for _, r := range results {
		switch r.Status {
		case Added:
			s.Added++
		case Removed:
			s.Removed++
		case Modified:
			s.Modified++
		case Unchanged:
			s.Unchanged++
		}
	}
	return s
}

// HasChanges returns true if any keys were added, removed, or modified.
func (s Summary) HasChanges() bool {
	return s.Added > 0 || s.Removed > 0 || s.Modified > 0
}

// String returns a human-readable one-line summary of the diff.
func (s Summary) String() string {
	parts := []string{}
	if s.Added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", s.Added))
	}
	if s.Removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", s.Removed))
	}
	if s.Modified > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", s.Modified))
	}
	if s.Unchanged > 0 {
		parts = append(parts, fmt.Sprintf("%d unchanged", s.Unchanged))
	}
	if len(parts) == 0 {
		return "no differences"
	}
	return strings.Join(parts, ", ")
}
