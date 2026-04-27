package cmd

import (
	"fmt"
	"io"
	"os"

	"envoy-sync/internal/diff"
	"envoy-sync/internal/envfile"
	"envoy-sync/internal/store"
)

// GitHubDiffOptions holds configuration for the github-diff command.
type GitHubDiffOptions struct {
	LocalFile string
	Repo      string
	Namespace string
	RemoteVars map[string]string
	Summary   bool
	Verbose   bool
	Out       io.Writer
}

// RunGitHubDiff compares a local .env file against a simulated GitHub secrets store.
func RunGitHubDiff(opts GitHubDiffOptions) error {
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}

	local, err := envfile.Parse(opts.LocalFile)
	if err != nil {
		return fmt.Errorf("parse local file: %w", err)
	}

	if opts.Namespace == "" {
		return fmt.Errorf("namespace must not be empty")
	}

	remote := store.NewGitHubStore(opts.Repo, opts.Namespace, opts.RemoteVars)
	remoteMap := remote.ToMap()

	result := diff.Compare(local, remoteMap)

	if opts.Summary {
		summary := diff.Summarize(result)
		fmt.Fprintln(out, summary.String())
		return nil
	}

	diff.Fprint(out, result, opts.Verbose)
	return nil
}
