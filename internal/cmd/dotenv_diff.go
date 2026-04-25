package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/yourorg/envoy-sync/internal/diff"
	"github.com/yourorg/envoy-sync/internal/envfile"
	"github.com/yourorg/envoy-sync/internal/store"
)

// DotenvDiffOptions configures the RunDotenvDiff command.
type DotenvDiffOptions struct {
	BaseFile   string
	TargetFile string
	Namespace  string
	Summary    bool
	Verbose    bool
	Out        io.Writer
}

// RunDotenvDiff compares two .env files (optionally scoped by namespace)
// and prints the diff to opts.Out.
func RunDotenvDiff(opts DotenvDiffOptions) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}

	baseVars, err := envfile.Parse(opts.BaseFile)
	if err != nil {
		return fmt.Errorf("reading base file %q: %w", opts.BaseFile, err)
	}

	targetVars, err := envfile.Parse(opts.TargetFile)
	if err != nil {
		return fmt.Errorf("reading target file %q: %w", opts.TargetFile, err)
	}

	baseStore := store.NewDotenvStore(opts.BaseFile, opts.Namespace, baseVars)
	targetStore := store.NewDotenvStore(opts.TargetFile, opts.Namespace, targetVars)

	baseMap := baseStore.All()
	targetMap := targetStore.All()

	changes := diff.Compare(baseMap, targetMap)

	if opts.Summary {
		summary := diff.Summarize(changes)
		fmt.Fprintln(opts.Out, summary.String())
		return nil
	}

	diff.Fprint(opts.Out, changes, opts.Verbose)
	return nil
}
