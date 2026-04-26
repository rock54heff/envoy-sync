package cmd

import (
	"fmt"
	"io"

	"github.com/yourorg/envoy-sync/internal/store"
)

// DotenvSyncOptions configures the behaviour of RunDotenvSync.
type DotenvSyncOptions struct {
	BaseFile   string
	TargetFile string
	Namespace  string
	DryRun     bool
	Verbose    bool
	Out        io.Writer
}

// RunDotenvSync syncs a local .env file into a dotenv-namespaced store,
// optionally writing the merged result back to the target file.
func RunDotenvSync(opts DotenvSyncOptions) error {
	base, err := store.NewFileStore(opts.BaseFile)
	if err != nil {
		return fmt.Errorf("loading base file %q: %w", opts.BaseFile, err)
	}

	var target *store.DotenvStore
	if opts.TargetFile != "" {
		target, err = store.NewDotenvStore(opts.TargetFile, opts.Namespace)
		if err != nil {
			return fmt.Errorf("loading target file %q: %w", opts.TargetFile, err)
		}
	} else {
		target = store.NewEmptyDotenvStore(opts.Namespace)
	}

	result, err := store.SyncToDotenv(base, target)
	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	if result.Summary.HasChanges() {
		fmt.Fprintf(opts.Out, "Changes detected (%s):\n", opts.Namespace)
		result.Summary.Fprint(opts.Out, opts.Verbose)
	} else {
		fmt.Fprintln(opts.Out, "Already in sync.")
	}

	if opts.DryRun {
		fmt.Fprintln(opts.Out, "[dry-run] no changes written.")
		return nil
	}

	outPath := opts.TargetFile
	if outPath == "" {
		outPath = opts.BaseFile
	}

	if err := store.WriteEnvFile(outPath, result.Store); err != nil {
		return fmt.Errorf("writing synced file %q: %w", outPath, err)
	}

	fmt.Fprintf(opts.Out, "Synced → %s\n", outPath)
	return nil
}
