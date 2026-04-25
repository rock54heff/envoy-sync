package cmd

import (
	"fmt"
	"io"

	"github.com/yourorg/envoy-sync/internal/diff"
	"github.com/yourorg/envoy-sync/internal/envfile"
	"github.com/yourorg/envoy-sync/internal/store"
)

// AWSdiffOptions holds parameters for the aws-diff command.
type AWSDiffOptions struct {
	LocalFile string
	Namespace string
	// RemoteVars simulates the remote SSM state (injected for testing).
	RemoteVars map[string]string
	Summary    bool
	Verbose    bool
}

// RunAWSDiff compares a local .env file against a simulated AWS SSM namespace
// and prints the diff to w. Returns an error on I/O or parse failure.
func RunAWSDiff(opts AWSDiffOptions, w io.Writer) error {
	if opts.Namespace == "" {
		return fmt.Errorf("aws-diff: namespace must not be empty")
	}

	localVars, err := envfile.Parse(opts.LocalFile)
	if err != nil {
		return fmt.Errorf("aws-diff: parse local file: %w", err)
	}

	localStore := store.New(localVars)

	remoteStore, err := store.NewAWSStore(opts.Namespace, opts.RemoteVars)
	if err != nil {
		return fmt.Errorf("aws-diff: create remote store: %w", err)
	}

	baseMap := store.ToMap(localStore)
	targetMap := remoteStore.ToMap()

	result := diff.Compare(baseMap, targetMap)

	if opts.Summary {
		summary := diff.Summarize(result)
		fmt.Fprintln(w, summary.String())
		return nil
	}

	diff.Fprint(w, result, opts.Verbose)
	return nil
}
