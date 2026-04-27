package store

import "fmt"

// SyncToGitHub syncs vars from a local map into a GitHubStore.
// In dry-run mode it reports what would change without applying.
func SyncToGitHub(local map[string]string, remote *GitHubStore, dryRun bool) ([]string, error) {
	var actions []string

	for k, localVal := range local {
		remoteVal, exists := remote.Get(k)
		if !exists {
			actions = append(actions, fmt.Sprintf("+ %s", k))
			if !dryRun {
				if err := remote.Set(k, localVal); err != nil {
					return actions, fmt.Errorf("set %s: %w", k, err)
				}
			}
		} else if remoteVal != localVal {
			actions = append(actions, fmt.Sprintf("~ %s", k))
			if !dryRun {
				if err := remote.Set(k, localVal); err != nil {
					return actions, fmt.Errorf("update %s: %w", k, err)
				}
			}
		}
	}

	for _, k := range remote.Keys() {
		if _, exists := local[k]; !exists {
			actions = append(actions, fmt.Sprintf("- %s", k))
			if !dryRun {
				if err := remote.Delete(k); err != nil {
					return actions, fmt.Errorf("delete %s: %w", k, err)
				}
			}
		}
	}

	return actions, nil
}
