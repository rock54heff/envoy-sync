package store

import "fmt"

// SyncToVercel applies the diff between a local env map and a VercelStore,
// writing additions and modifications and removing deleted keys.
// When dryRun is true no mutations are applied and the planned changes
// are returned as a human-readable summary.
func SyncToVercel(local map[string]string, remote *VercelStore, dryRun bool) (string, error) {
	remoteMap := remote.All()

	var added, updated, removed []string

	// Keys to add or update
	for k, v := range local {
		rv, exists := remoteMap[k]
		if !exists {
			added = append(added, k)
			if !dryRun {
				if err := remote.Set(k, v); err != nil {
					return "", fmt.Errorf("SyncToVercel: set %q: %w", k, err)
				}
			}
		} else if rv != v {
			updated = append(updated, k)
			if !dryRun {
				if err := remote.Set(k, v); err != nil {
					return "", fmt.Errorf("SyncToVercel: update %q: %w", k, err)
				}
			}
		}
	}

	// Keys to remove (present in remote but not in local)
	for k := range remoteMap {
		if _, exists := local[k]; !exists {
			removed = append(removed, k)
			if !dryRun {
				if err := remote.Delete(k); err != nil {
					return "", fmt.Errorf("SyncToVercel: delete %q: %w", k, err)
				}
			}
		}
	}

	if len(added)+len(updated)+len(removed) == 0 {
		return fmt.Sprintf("vercel[%s]: already in sync", remote.Name()), nil
	}

	summary := fmt.Sprintf("vercel[%s]: +%d ~%d -%d",
		remote.Name(), len(added), len(updated), len(removed))
	if dryRun {
		summary += " (dry-run)"
	}
	return summary, nil
}
