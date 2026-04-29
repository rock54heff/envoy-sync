package cmd

import (
	"fmt"
	"io"

	"github.com/your-org/envoy-sync/internal/diff"
	"github.com/your-org/envoy-sync/internal/envfile"
	"github.com/your-org/envoy-sync/internal/store"
)

// VercelDiffOptions holds parameters for the vercel diff command.
type VercelDiffOptions struct {
	LocalFile string
	Namespace string
	TeamID    string
	Verbose   bool
	Summary   bool
}

// RunVercelDiff compares a local .env file against a Vercel environment store.
func RunVercelDiff(opts VercelDiffOptions, out io.Writer) error {
	if opts.LocalFile == "" {
		return fmt.Errorf("local file path is required")
	}
	if opts.Namespace == "" {
		return fmt.Errorf("namespace (Vercel project name) is required")
	}

	localVars, err := envfile.Parse(opts.LocalFile)
	if err != nil {
		return fmt.Errorf("failed to parse local file %q: %w", opts.LocalFile, err)
	}

	localStore := store.New(localVars)

	remoteStore := store.NewVercelStore(map[string]string{}, opts.Namespace, opts.TeamID)

	localMap := store.ToMap(localStore)
	remoteMap := store.ToMap(remoteStore)

	results := diff.Compare(localMap, remoteMap)

	if opts.Summary {
		summary := diff.Summarize(results)
		_, err = fmt.Fprintln(out, summary.String())
		return err
	}

	return diff.Fprint(out, results, opts.Verbose)
}
