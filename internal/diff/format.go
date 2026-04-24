package diff

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// FormatOptions controls how the diff is rendered.
type FormatOptions struct {
	Color   bool
	Verbose bool // include unchanged keys
}

// Fprint writes a human-readable diff to w.
func Fprint(w io.Writer, r Result, opts FormatOptions) {
	if !r.HasChanges() && !opts.Verbose {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	for _, k := range SortedKeys(r.Added) {
		line := fmt.Sprintf("+ %s=%s", k, r.Added[k])
		if opts.Color {
			line = colorGreen + line + colorReset
		}
		fmt.Fprintln(w, line)
	}

	for _, k := range SortedKeys(r.Removed) {
		line := fmt.Sprintf("- %s=%s", k, r.Removed[k])
		if opts.Color {
			line = colorRed + line + colorReset
		}
		fmt.Fprintln(w, line)
	}

	for _, k := range SortedKeys(r.Modified) {
		line := fmt.Sprintf("~ %s=%s", k, r.Modified[k])
		if opts.Color {
			line = colorYellow + line + colorReset
		}
		fmt.Fprintln(w, line)
	}

	if opts.Verbose {
		for _, k := range SortedKeys(r.Unchanged) {
			fmt.Fprintf(w, "  %s=%s\n", k, r.Unchanged[k])
		}
	}
}

// Sprint returns the formatted diff as a string.
func Sprint(r Result, opts FormatOptions) string {
	var sb strings.Builder
	Fprint(&sb, r, opts)
	return sb.String()
}
