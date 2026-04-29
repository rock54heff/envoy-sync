package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/envoy-sync/internal/diff"
	"github.com/yourusername/envoy-sync/internal/envfile"
	"github.com/yourusername/envoy-sync/internal/store"
)

// RunGitLabDiff compares a local .env file against a GitLab CI/CD variable namespace.
func RunGitLabDiff(localPath, namespace string, summary bool, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	if namespace == "" {
		return fmt.Errorf("namespace must not be empty")
	}

	localVars, err := envfile.Parse(localPath)
	if err != nil {
		return fmt.Errorf("failed to parse local file %q: %w", localPath, err)
	}

	glStore := store.NewGitLabStore(localVars, namespace)
	remoteVars := store.ToMap(glStore)

	changes := diff.Compare(localVars, remoteVars)

	if summary {
		rep := diff.Summarize(changes)
		_, err = fmt.Fprintln(w, rep.String())
		return err
	}

	diff.Fprint(w, changes, false)
	return nil
}
