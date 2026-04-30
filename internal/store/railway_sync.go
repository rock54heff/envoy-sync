package store

import (
	"fmt"
	"io"
)

// SyncToRailway applies changes from a local env map to the target RailwayStore.
// It adds/updates keys present in local and removes keys absent from local.
// If dryRun is true, changes are described to w but not applied.
func SyncToRailway(local map[string]string, target *RailwayStore, dryRun bool, w io.Writer) error {
	remote := target.All()

	for k, v := range local {
		if existing, ok := remote[k]; !ok {
			if dryRun {
				fmt.Fprintf(w, "[dry-run] would add %s=%s to %s\n", k, v, target.Namespace())
			} else {
				if err := target.Set(k, v); err != nil {
					return fmt.Errorf("railway sync: set %q: %w", k, err)
				}
			}
		} else if existing != v {
			if dryRun {
				fmt.Fprintf(w, "[dry-run] would update %s in %s\n", k, target.Namespace())
			} else {
				if err := target.Set(k, v); err != nil {
					return fmt.Errorf("railway sync: update %q: %w", k, err)
				}
			}
		}
	}

	for k := range remote {
		if _, ok := local[k]; !ok {
			if dryRun {
				fmt.Fprintf(w, "[dry-run] would delete %s from %s\n", k, target.Namespace())
			} else {
				if err := target.Delete(k); err != nil {
					return fmt.Errorf("railway sync: delete %q: %w", k, err)
				}
			}
		}
	}

	return nil
}
