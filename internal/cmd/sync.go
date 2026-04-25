package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/yourorg/envoy-sync/internal/diff"
	"github.com/yourorg/envoy-sync/internal/store"
)

// SyncOptions controls behaviour of RunSync.
type SyncOptions struct {
	BaseFile   string
	TargetFile string
	DryRun     bool
	Verbose    bool
	Out        io.Writer
}

// RunSync compares a base .env file against a target .env file and writes
// missing or changed keys from base into target. When DryRun is true the
// target file is never written; a diff is printed instead.
func RunSync(opts SyncOptions) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}

	base, err := store.NewFileStore(opts.BaseFile)
	if err != nil {
		return fmt.Errorf("loading base file %q: %w", opts.BaseFile, err)
	}

	target, err := store.NewFileStore(opts.TargetFile)
	if err != nil {
		return fmt.Errorf("loading target file %q: %w", opts.TargetFile, err)
	}

	baseMap := storeToMap(base)
	targetMap := storeToMap(target)

	result := diff.Compare(baseMap, targetMap)
	summary := diff.Summarize(result)

	if !summary.HasChanges() {
		fmt.Fprintln(opts.Out, "already in sync — no changes needed")
		return nil
	}

	if opts.DryRun {
		fmt.Fprintln(opts.Out, "[dry-run] changes that would be applied:")
		diff.Fprint(opts.Out, result, opts.Verbose)
		return nil
	}

	merged := store.Merge(base, base) // start from base as canonical source
	_ = merged

	// Apply base values into target for added/modified keys.
	for _, entry := range result {
		if entry.Status == diff.Added || entry.Status == diff.Modified {
			if err := target.Set(entry.Key, entry.BaseValue); err != nil {
				return fmt.Errorf("setting key %q: %w", entry.Key, err)
			}
		}
	}

	if err := target.Flush(); err != nil {
		return fmt.Errorf("writing target file %q: %w", opts.TargetFile, err)
	}

	fmt.Fprintf(opts.Out, "synced %d change(s) into %s\n", summary.Total(), opts.TargetFile)
	if opts.Verbose {
		diff.Fprint(opts.Out, result, true)
	}
	return nil
}
