package store

import "fmt"

// SyncToNetlify applies changes from a local env map to the target NetlifyStore.
// It adds/updates keys present in local and removes keys absent from local
// that exist in the remote store. Returns the number of changes applied.
func SyncToNetlify(local map[string]string, remote *NetlifyStore, dryRun bool) (int, error) {
	if remote == nil {
		return 0, fmt.Errorf("netlify: remote store must not be nil")
	}

	changes := 0
	remoteAll := remote.All()

	// Add or update keys from local.
	for k, v := range local {
		existing, ok := remoteAll[k]
		if !ok || existing != v {
			if !dryRun {
				if err := remote.Set(k, v); err != nil {
					return changes, fmt.Errorf("netlify: set %q: %w", k, err)
				}
			}
			changes++
		}
	}

	// Remove keys not present in local.
	for k := range remoteAll {
		if _, ok := local[k]; !ok {
			if !dryRun {
				remote.Delete(k)
			}
			changes++
		}
	}

	return changes, nil
}
