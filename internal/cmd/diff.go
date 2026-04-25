package cmd

import (
	"fmt"
	"os"

	"github.com/user/envoy-sync/internal/diff"
	"github.com/user/envoy-sync/internal/store"
)

// DiffOptions holds configuration for the diff command.
type DiffOptions struct {
	BaseFile     string
	TargetFile   string
	Verbose      bool
	OutputFormat string // "text" or "summary"
}

// RunDiff compares two .env files and prints the differences.
func RunDiff(opts DiffOptions) error {
	base, err := store.NewFileStore(opts.BaseFile)
	if err != nil {
		return fmt.Errorf("loading base file %q: %w", opts.BaseFile, err)
	}

	target, err := store.NewFileStore(opts.TargetFile)
	if err != nil {
		return fmt.Errorf("loading target file %q: %w", opts.TargetFile, err)
	}

	baseVars := storeToMap(base)
	targetVars := storeToMap(target)

	changes := diff.Compare(baseVars, targetVars)

	switch opts.OutputFormat {
	case "summary":
		summary := diff.Summarize(changes)
		fmt.Println(summary.String())
	default:
		diff.Fprint(os.Stdout, changes, opts.Verbose)
	}

	return nil
}

// storeToMap extracts all key-value pairs from a store into a plain map.
func storeToMap(s interface {
	Keys() []string
	Get(string) (string, bool)
}) map[string]string {
	result := make(map[string]string)
	for _, k := range s.Keys() {
		if v, ok := s.Get(k); ok {
			result[k] = v
		}
	}
	return result
}
