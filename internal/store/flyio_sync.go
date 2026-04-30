package store

import "fmt"

// SyncToFlyio applies changes from a local env map into a FlyioStore.
// It adds/updates keys present in local and deletes keys absent from local
// that exist in the remote store. Returns the number of changes applied.
func SyncToFlyio(local map[string]string, remote *FlyioStore, dryRun bool) (int, error) {
	remoteMap := remote.ToMap()
	changes := 0

	// Set or update keys from local.
	for k, v := range local {
		remoteVal, err := remote.Get(k)
		if err != nil || remoteVal != v {
			if !dryRun {
				if setErr := remote.Set(k, v); setErr != nil {
					return changes, fmt.Errorf("flyio sync: set %q: %w", k, setErr)
				}
			}
			changes++
		}
	}

	// Delete keys present in remote but absent from local.
	for k := range remoteMap {
		if _, exists := local[k]; !exists {
			if !dryRun {
				if delErr := remote.Delete(k); delErr != nil {
					return changes, fmt.Errorf("flyio sync: delete %q: %w", k, delErr)
				}
			}
			changes++
		}
	}

	return changes, nil
}
