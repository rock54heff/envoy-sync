// Package store defines the Store interface and common store utilities
// used across envoy-sync backends (memory, file, remote).
package store

// Store is the common interface for any key-value environment variable backend.
type Store interface {
	// Get returns the value associated with key and whether the key exists.
	Get(key string) (string, bool)

	// Set inserts or updates key with the given value.
	Set(key, value string)

	// Delete removes key from the store.
	Delete(key string)

	// All returns a snapshot of all key-value pairs in the store.
	All() map[string]string
}

// Merge copies all entries from src into dst, overwriting existing keys.
func Merge(dst, src Store) {
	for k, v := range src.All() {
		dst.Set(k, v)
	}
}

// Keys returns all keys present in the store in an unordered slice.
func Keys(s Store) []string {
	all := s.All()
	keys := make([]string, 0, len(all))
	for k := range all {
		keys = append(keys, k)
	}
	return keys
}
