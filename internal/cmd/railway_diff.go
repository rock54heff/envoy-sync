package cmd

import (
	"fmt"
	"io"

	"envoy-sync/internal/diff"
	"envoy-sync/internal/envfile"
	"envoy-sync/internal/store"
)

// RailwayDiffOptions holds parameters for the railway diff command.
type RailwayDiffOptions struct {
	LocalFile   string
	ProjectID   string
	Environment string
	RemoteVars  map[string]string
	Verbose     bool
	Format      string
}

// RunRailwayDiff compares a local .env file against a Railway environment store.
func RunRailwayDiff(opts RailwayDiffOptions, w io.Writer) error {
	if opts.LocalFile == "" {
		return fmt.Errorf("local file path must not be empty")
	}
	if opts.ProjectID == "" {
		return fmt.Errorf("railway project ID must not be empty")
	}

	local, err := envfile.Parse(opts.LocalFile)
	if err != nil {
		return fmt.Errorf("failed to parse local file: %w", err)
	}

	remote := store.NewRailwayStore(opts.RemoteVars, opts.ProjectID, opts.Environment)
	remoteMap := remote.All()

	changes := diff.Compare(local, remoteMap)

	if opts.Format == "summary" {
		summary := diff.Summarize(changes)
		fmt.Fprintln(w, summary.String())
		return nil
	}

	return diff.Fprint(w, changes, opts.Verbose)
}
