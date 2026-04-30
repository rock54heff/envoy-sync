package cmd

import (
	"fmt"
	"io"

	"github.com/your-org/envoy-sync/internal/diff"
	"github.com/your-org/envoy-sync/internal/envfile"
	"github.com/your-org/envoy-sync/internal/store"
)

// NetlifyDiffOptions configures the netlify diff command.
type NetlifyDiffOptions struct {
	LocalFile string
	SiteID    string
	Namespace string
	RemoteVars map[string]string
	Summary   bool
	Verbose   bool
	Out       io.Writer
}

// RunNetlifyDiff compares a local .env file against a simulated Netlify store
// and writes the diff to opts.Out.
func RunNetlifyDiff(opts NetlifyDiffOptions) error {
	if opts.LocalFile == "" {
		return fmt.Errorf("netlify diff: local file path is required")
	}
	if opts.SiteID == "" {
		return fmt.Errorf("netlify diff: site ID is required")
	}

	local, err := envfile.Parse(opts.LocalFile)
	if err != nil {
		return fmt.Errorf("netlify diff: parse local file: %w", err)
	}

	remote := store.NewNetlifyStore(opts.RemoteVars, opts.SiteID, opts.Namespace)
	remoteMap := remote.All()

	results := diff.Compare(local, remoteMap)

	if opts.Summary {
		summary := diff.Summarize(results)
		_, err = fmt.Fprintln(opts.Out, summary.String())
		return err
	}

	return diff.Fprint(opts.Out, results, opts.Verbose)
}
