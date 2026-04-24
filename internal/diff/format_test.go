package diff_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envoy-sync/internal/diff"
)

func TestFprint_NoChanges(t *testing.T) {
	r := diff.Compare(
		map[string]string{"A": "1"},
		map[string]string{"A": "1"},
	)
	out := diff.Sprint(r, diff.FormatOptions{})
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no-diff message, got: %q", out)
	}
}

func TestFprint_ShowsAdded(t *testing.T) {
	r := diff.Compare(
		map[string]string{},
		map[string]string{"NEW": "val"},
	)
	out := diff.Sprint(r, diff.FormatOptions{})
	if !strings.Contains(out, "+ NEW=val") {
		t.Errorf("expected '+ NEW=val' in output, got: %q", out)
	}
}

func TestFprint_ShowsRemoved(t *testing.T) {
	r := diff.Compare(
		map[string]string{"OLD": "gone"},
		map[string]string{},
	)
	out := diff.Sprint(r, diff.FormatOptions{})
	if !strings.Contains(out, "- OLD=gone") {
		t.Errorf("expected '- OLD=gone' in output, got: %q", out)
	}
}

func TestFprint_ShowsModified(t *testing.T) {
	r := diff.Compare(
		map[string]string{"KEY": "old"},
		map[string]string{"KEY": "new"},
	)
	out := diff.Sprint(r, diff.FormatOptions{})
	if !strings.Contains(out, "~ KEY=new") {
		t.Errorf("expected '~ KEY=new' in output, got: %q", out)
	}
}

func TestFprint_VerboseShowsUnchanged(t *testing.T) {
	r := diff.Compare(
		map[string]string{"SAME": "val", "DIFF": "a"},
		map[string]string{"SAME": "val", "DIFF": "b"},
	)
	out := diff.Sprint(r, diff.FormatOptions{Verbose: true})
	if !strings.Contains(out, "SAME=val") {
		t.Errorf("expected unchanged key in verbose output, got: %q", out)
	}
}
