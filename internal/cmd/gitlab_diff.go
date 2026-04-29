package cmd

import (
	"fmt"
	"io"

	"envoy-sync/internal/diff"
	"envoy-sync/internal/envfile"
	"envoy-sync/internal/store"
)

// GitLabDiffOptions holds parameters for the gitlab-diff command.
type GitLabDiffOptions struct {
	LocalFile string
	Namespace string
	RemoteVars map[string]string // injected in tests; represents remote GitLab state
	Format     string            // "text" or "summary"
	Verbose    bool
}

// RunGitLabDiff compares a local .env file against a simulated GitLab variable store
// and writes the diff result to out.
func RunGitLabDiff(opts GitLabDiffOptions, out io.Writer) error {
	if opts.LocalFile == "" {
		return fmt.Errorf("gitlab-diff: local file path must not be empty")
	}
	if opts.Namespace == "" {
		return fmt.Errorf("gitlab-diff: namespace must not be empty")
	}

	localVars, err := envfile.Parse(opts.LocalFile)
	if err != nil {
		return fmt.Errorf("gitlab-diff: parse local file: %w", err)
	}

	remote := store.NewGitLabStore(opts.Namespace, opts.RemoteVars)
	remoteMap := remote.All()

	result := diff.Compare(localVars, remoteMap)

	switch opts.Format {
	case "summary":
		summary := diff.Summarize(result)
		_, err = fmt.Fprintln(out, summary.String())
	default:
		err = diff.Fprint(out, result, opts.Verbose)
	}
	return err
}
