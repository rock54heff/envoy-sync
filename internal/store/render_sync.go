package store

import "fmt"

// SyncToRender pushes changes from a local env map into a RenderStore.
// If dryRun is true, changes are computed but not applied.
// Returns the number of keys added/updated, deleted, and any error.
func SyncToRender(local map[string]string, remote *RenderStore, dryRun bool) (added, deleted int, err error) {
	if remote == nil {
		return 0, 0, fmt.Errorf("render: remote store must not be nil")
	}

	remoteAll := remote.All()

	// Determine keys to set (new or changed)
	for k, v := range local {
		remoteVal, exists := remoteAll[k]
		if !exists || remoteVal != v {
			if !dryRun {
				if setErr := remote.Set(k, v); setErr != nil {
					return added, deleted, fmt.Errorf("render: failed to set %q: %w", k, setErr)
				}
			}
			added++
		}
	}

	// Determine keys to delete (present in remote but not in local)
	for k := range remoteAll {
		if _, exists := local[k]; !exists {
			if !dryRun {
				if delErr := remote.Delete(k); delErr != nil {
					return added, deleted, fmt.Errorf("render: failed to delete %q: %w", k, delErr)
				}
			}
			deleted++
		}
	}

	return added, deleted, nil
}
