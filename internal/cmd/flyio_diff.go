package cmd

import (
	"fmt"
	"io"

	"github.com/your-org/envoy-sync/internal/diff"
	"github.com/your-org/envoy-sync/internal/envfile"
	"github.com/your-org/envoy-sync/internal/store"
)

// FlyioDiffOptions holds parameters for the fly.io diff command.
type FlyioDiffOptions struct {
	LocalFile string
	AppName   string
	Namespace string
	Verbose   bool
	Summary   bool
}

// RunFlyioDiff compares a local .env file against a FlyioStore and prints the diff.
func RunFlyioDiff(opts FlyioDiffOptions, out io.Writer) error {
	if opts.LocalFile == "" {
		return fmt.Errorf("flyio diff: local file path must not be empty")
	}
	if opts.AppName == "" {
		return fmt.Errorf("flyio diff: app name must not be empty")
	}

	localVars, err := envfile.Parse(opts.LocalFile)
	if err != nil {
		return fmt.Errorf("flyio diff: parse local file: %w", err)
	}

	remoteVars := make(map[string]string)
	if opts.Namespace != "" {
		// In a real implementation this would fetch from Fly.io API.
		// For now we start with an empty remote to allow testing.
		remoteVars = map[string]string{}
	}

	remote := store.NewFlyioStore(opts.AppName, remoteVars)
	result := diff.Compare(localVars, remote.ToMap())

	if opts.Summary {
		summary := diff.Summarize(result)
		_, err = fmt.Fprintln(out, summary.String())
		return err
	}

	return diff.Fprint(out, result, opts.Verbose)
}
