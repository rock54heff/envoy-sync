package store

import (
	"fmt"
	"strings"
)

// Storer is implemented by any store that supports Get, Set, Delete, and All.
type Storer interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	All() map[string]string
}

// SyncResult summarises the changes applied during a Sync operation.
type SyncResult struct {
	Added   []string
	Updated []string
	Deleted []string
}

// String returns a human-readable summary of the sync result.
func (r SyncResult) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("added: %d, updated: %d, deleted: %d",
		len(r.Added), len(r.Updated), len(r.Deleted)))
	return sb.String()
}

// HasChanges returns true if any keys were added, updated, or deleted.
func (r SyncResult) HasChanges() bool {
	return len(r.Added)+len(r.Updated)+len(r.Deleted) > 0
}

// Sync copies all key/value pairs from src into dst.
// Keys present in dst but absent in src are deleted when deleteOrphans is true.
func Sync(dst, src Storer, deleteOrphans bool) (SyncResult, error) {
	var result SyncResult

	srcVars := src.All()
	dstVars := dst.All()

	for k, v := range srcVars {
		if existing, ok := dstVars[k]; !ok {
			if err := dst.Set(k, v); err != nil {
				return result, fmt.Errorf("sync: set %q: %w", k, err)
			}
			result.Added = append(result.Added, k)
		} else if existing != v {
			if err := dst.Set(k, v); err != nil {
				return result, fmt.Errorf("sync: update %q: %w", k, err)
			}
			result.Updated = append(result.Updated, k)
		}
	}

	if deleteOrphans {
		for k := range dstVars {
			if _, ok := srcVars[k]; !ok {
				if err := dst.Delete(k); err != nil {
					return result, fmt.Errorf("sync: delete %q: %w", k, err)
				}
				result.Deleted = append(result.Deleted, k)
			}
		}
	}

	return result, nil
}
