package store

import "fmt"

// SyncToGitLab applies changes from a local env map into a GitLabStore.
// It adds/updates keys present in local and removes keys absent from local
// that exist in the remote store. Returns the number of changes applied.
func SyncToGitLab(local map[string]string, remote *GitLabStore, dryRun bool) (int, error) {
	remoteAll := remote.All()
	changes := 0

	// Set or update keys from local.
	for k, v := range local {
		remoteVal, exists := remoteAll[k]
		if !exists || remoteVal != v {
			if !dryRun {
				if err := remote.Set(k, v); err != nil {
					return changes, fmt.Errorf("gitlab sync: set %q: %w", k, err)
				}
			}
			changes++
		}
	}

	// Remove keys not present in local.
	for k := range remoteAll {
		if _, ok := local[k]; !ok {
			if !dryRun {
				if err := remote.Delete(k); err != nil {
					return changes, fmt.Errorf("gitlab sync: delete %q: %w", k, err)
				}
			}
			changes++
		}
	}

	return changes, nil
}
