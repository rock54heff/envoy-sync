package store

import (
	"fmt"
	"sort"
)

// Store represents a named collection of environment variables.
type Store struct {
	Name string
	Vars map[string]string
}

// New creates a new Store with the given name and variables.
func New(name string, vars map[string]string) *Store {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &Store{Name: name, Vars: copy}
}

// Get returns the value for a key and whether it exists.
func (s *Store) Get(key string) (string, bool) {
	v, ok := s.Vars[key]
	return v, ok
}

// Set sets a key-value pair in the store.
func (s *Store) Set(key, value string) {
	if s.Vars == nil {
		s.Vars = make(map[string]string)
	}
	s.Vars[key] = value
}

// Delete removes a key from the store.
func (s *Store) Delete(key string) {
	delete(s.Vars, key)
}

// Keys returns all keys in sorted order.
func (s *Store) Keys() []string {
	keys := make([]string, 0, len(s.Vars))
	for k := range s.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Len returns the number of variables in the store.
func (s *Store) Len() int {
	return len(s.Vars)
}

// Merge returns a new Store combining s with other.
// Keys in other take precedence over keys in s.
func (s *Store) Merge(other *Store) *Store {
	merged := make(map[string]string, len(s.Vars)+len(other.Vars))
	for k, v := range s.Vars {
		merged[k] = v
	}
	for k, v := range other.Vars {
		merged[k] = v
	}
	name := fmt.Sprintf("%s+%s", s.Name, other.Name)
	return New(name, merged)
}
