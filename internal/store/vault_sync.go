package store

import "fmt"

// VaultSyncResult records what happened to each key during a vault sync.
type VaultSyncResult struct {
	Added    []string
	Updated  []string
	Deleted  []string
	Skipped  []string
}

// SyncToVault pushes keys from src (e.g. a FileStore) into dst (a VaultStore).
// When dryRun is true no mutations are applied but the result is still populated.
func SyncToVault(src EnvReader, dst *VaultStore, dryRun bool) (VaultSyncResult, error) {
	var result VaultSyncResult

	srcMap := ToMap(src)
	dstMap := dst.All()

	// Add / update keys present in src.
	for _, k := range sortedStringKeys(srcMap) {
		v := srcMap[k]
		existing, exists := dstMap[k]
		if !exists {
			result.Added = append(result.Added, k)
			if !dryRun {
				if err := dst.Set(k, v); err != nil {
					return result, fmt.Errorf("vault sync: set %q: %w", k, err)
				}
			}
		} else if existing != v {
			result.Updated = append(result.Updated, k)
			if !dryRun {
				if err := dst.Set(k, v); err != nil {
					return result, fmt.Errorf("vault sync: update %q: %w", k, err)
				}
			}
		} else {
			result.Skipped = append(result.Skipped, k)
		}
	}

	// Collect keys in dst not present in src (informational only — not deleted).
	for _, k := range sortedStringKeys(dstMap) {
		if _, ok := srcMap[k]; !ok {
			result.Deleted = append(result.Deleted, k)
		}
	}

	return result, nil
}

// EnvReader is satisfied by any store that exposes Keys and Get.
type EnvReader interface {
	Keys() []string
	Get(key string) (string, bool)
}
