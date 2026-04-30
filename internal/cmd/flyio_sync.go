package cmd

import (
	"fmt"
	"io"

	"github.com/your-org/envoy-sync/internal/store"
)

// FlyioSyncOptions holds the parameters for syncing a local .env file to Fly.io.
type FlyioSyncOptions struct {
	BaseFile  string
	AppName   string
	Token     string
	DryRun    bool
	Verbose   bool
	Out       io.Writer
}

// RunFlyioSync compares a local .env file against a Fly.io app's secrets and
// pushes any additions or modifications. Removals are never applied automatically.
func RunFlyioSync(opts FlyioSyncOptions) error {
	if opts.BaseFile == "" {
		return fmt.Errorf("base file path is required")
	}
	if opts.AppName == "" {
		return fmt.Errorf("fly.io app name is required")
	}

	local, err := store.NewFileStore(opts.BaseFile)
	if err != nil {
		return fmt.Errorf("loading base file: %w", err)
	}

	remote := store.NewFlyioStore(map[string]string{}, opts.AppName, opts.Token)

	result, err := store.SyncToFlyio(local, remote)
	if err != nil {
		return fmt.Errorf("syncing to fly.io: %w", err)
	}

	if result.Added == 0 && result.Updated == 0 {
		fmt.Fprintln(opts.Out, "fly.io secrets already in sync")
		return nil
	}

	if opts.DryRun {
		fmt.Fprintf(opts.Out, "[dry-run] would set %d secret(s) on fly.io app %q\n",
			result.Added+result.Updated, opts.AppName)
		if opts.Verbose {
			for _, k := range result.Keys {
				fmt.Fprintf(opts.Out, "  ~ %s\n", k)
			}
		}
		return nil
	}

	fmt.Fprintf(opts.Out, "synced %d secret(s) to fly.io app %q\n",
		result.Added+result.Updated, opts.AppName)
	if opts.Verbose {
		for _, k := range result.Keys {
			fmt.Fprintf(opts.Out, "  ~ %s\n", k)
		}
	}
	return nil
}
