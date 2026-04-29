package cmd

import (
	"fmt"
	"io"

	"github.com/yourusername/envoy-sync/internal/diff"
	"github.com/yourusername/envoy-sync/internal/envfile"
	"github.com/yourusername/envoy-sync/internal/store"
)

// GitLabSyncOptions holds parameters for the gitlab-sync command.
type GitLabSyncOptions struct {
	BaseFile  string
	Token     string
	ProjectID string
	Namespace string
	DryRun    bool
	Out       io.Writer
}

// RunGitLabSync syncs a local .env file to GitLab CI/CD variables.
func RunGitLabSync(opts GitLabSyncOptions) error {
	vars, err := envfile.Parse(opts.BaseFile)
	if err != nil {
		return fmt.Errorf("parsing base file: %w", err)
	}

	local := store.New(vars)

	remoteVars := map[string]string{} // in production: fetched via GitLab API
	remote := store.NewGitLabStore(remoteVars, opts.Token, opts.ProjectID, opts.Namespace)

	results := diff.Compare(store.ToMap(local), store.ToMap(remote))
	summary := diff.Summarize(results)

	if !summary.HasChanges() {
		fmt.Fprintln(opts.Out, "gitlab: already in sync, nothing to do")
		return nil
	}

	if opts.DryRun {
		fmt.Fprintln(opts.Out, "gitlab: dry-run mode, no changes written")
		diff.Fprint(opts.Out, results, false)
		return nil
	}

	if err := store.SyncToGitLab(local, remote); err != nil {
		return fmt.Errorf("syncing to gitlab: %w", err)
	}

	fmt.Fprintln(opts.Out, summary.String())
	return nil
}
