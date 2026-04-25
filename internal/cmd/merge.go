package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/yourorg/envoy-sync/internal/store"
)

// MergeOptions configures the behaviour of RunMerge.
type MergeOptions struct {
	BaseFile     string
	OverrideFile string
	OutputFile   string
	DryRun       bool
	Out          io.Writer
}

// RunMerge loads two .env files, merges them (override wins on conflict),
// and writes the result to OutputFile (or stdout when DryRun is true).
func RunMerge(opts MergeOptions) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}

	base, err := store.NewFileStore(opts.BaseFile)
	if err != nil {
		return fmt.Errorf("loading base file %q: %w", opts.BaseFile, err)
	}

	override, err := store.NewFileStore(opts.OverrideFile)
	if err != nil {
		return fmt.Errorf("loading override file %q: %w", opts.OverrideFile, err)
	}

	merged := store.Merge(base, override)

	if opts.DryRun {
		fmt.Fprintln(opts.Out, "# dry-run: merged output")
		for _, k := range store.Keys(merged) {
			v, _ := merged.Get(k)
			fmt.Fprintf(opts.Out, "%s=%s\n", k, v)
		}
		return nil
	}

	if opts.OutputFile == "" {
		return fmt.Errorf("output file path must not be empty")
	}

	if err := store.WriteEnvFile(opts.OutputFile, merged); err != nil {
		return fmt.Errorf("writing merged output to %q: %w", opts.OutputFile, err)
	}

	fmt.Fprintf(opts.Out, "merged %q + %q -> %q\n", opts.BaseFile, opts.OverrideFile, opts.OutputFile)
	return nil
}
