package store

import "sort"

// Store is the interface representing a readable and writable key-value env store.
type Store interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Delete(key string)
	All() map[string]string
}

// Merge returns a new in-memory Store containing all keys from base,
// with any keys in override taking precedence. Neither base nor override
// is mutated.
func Merge(base, override Store) Store {
	merged := make(map[string]string)

	for k, v := range base.All() {
		merged[k] = v
	}

	for k, v := range override.All() {
		merged[k] = v
	}

	return New(merged)
}

// Keys returns a sorted slice of all keys present in the given Store.
func Keys(s Store) []string {
	all := s.All()
	keys := make([]string, 0, len(all))
	for k := range all {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
