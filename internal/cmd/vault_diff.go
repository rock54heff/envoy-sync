package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envoy-sync/internal/diff"
	"github.com/user/envoy-sync/internal/store"
)

// VaultDiffOptions configures the RunVaultDiff command.
type VaultDiffOptions struct {
	LocalFile  string
	Namespace  string
	RemoteVars map[string]string // injected in tests; nil means real Vault call
	Verbose    bool
	Summary    bool
	Out        io.Writer
}

// RunVaultDiff compares a local .env file against a VaultStore and prints
// the diff to opts.Out (defaults to os.Stdout).
func RunVaultDiff(opts VaultDiffOptions) error {
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}

	fs, err := store.NewFileStore(opts.LocalFile)
	if err != nil {
		return fmt.Errorf("vault-diff: loading local file: %w", err)
	}

	remoteVars := opts.RemoteVars
	if remoteVars == nil {
		// Real implementation would fetch from Vault here.
		remoteVars = map[string]string{}
	}

	vs := store.NewVaultStore(opts.Namespace, remoteVars)

	localMap := store.ToMap(fs)
	vaultMap := vs.All()

	results := diff.Compare(localMap, vaultMap)

	if opts.Summary {
		summary := diff.Summarize(results)
		fmt.Fprintln(out, summary.String())
		return nil
	}

	diff.Fprint(out, results, opts.Verbose)
	return nil
}
