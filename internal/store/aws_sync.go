package store

import (
	"fmt"
	"io"
)

// SyncToAWS pushes all keys from a local Store into an AWSStore,
// deleting remote keys that are absent locally when deleteStale is true.
// Progress is written to w. Returns the number of changes applied.
func SyncToAWS(local Store, remote *AWSStore, deleteStale bool, w io.Writer) (int, error) {
	localMap := ToMap(local)
	remoteMap := remote.ToMap()

	changes := 0

	// Upsert local keys into remote.
	for k, v := range localMap {
		remoteVal, err := remote.Get(k)
		if err != nil || remoteVal != v {
			if err := remote.Set(k, v); err != nil {
				return changes, fmt.Errorf("aws sync: set %q: %w", k, err)
			}
			fmt.Fprintf(w, "  set %s\n", k)
			changes++
		}
	}

	// Optionally remove stale remote keys.
	if deleteStale {
		for k := range remoteMap {
			if _, exists := localMap[k]; !exists {
				if err := remote.Delete(k); err != nil {
					return changes, fmt.Errorf("aws sync: delete %q: %w", k, err)
				}
				fmt.Fprintf(w, "  deleted %s\n", k)
				changes++
			}
		}
	}

	if changes == 0 {
		fmt.Fprintln(w, "already in sync")
	}
	return changes, nil
}
