package store

import (
	"fmt"
	"sort"
	"sync"
)

// RemoteStore simulates a remote secret store (e.g. AWS SSM, Vault).
// It holds key/value pairs in memory with a namespace prefix.
type RemoteStore struct {
	mu        sync.RWMutex
	namespace string
	vars      map[string]string
}

// NewRemoteStore creates a RemoteStore with the given namespace and initial vars.
func NewRemoteStore(namespace string, vars map[string]string) *RemoteStore {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &RemoteStore{
		namespace: namespace,
		vars:      copy,
	}
}

// Namespace returns the store's namespace.
func (r *RemoteStore) Namespace() string {
	return r.namespace
}

// Get retrieves a value by key. Returns an error if the key does not exist.
func (r *RemoteStore) Get(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.vars[key]
	if !ok {
		return "", fmt.Errorf("remote store %q: key %q not found", r.namespace, key)
	}
	return v, nil
}

// Set stores a key/value pair.
func (r *RemoteStore) Set(key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.vars[key] = value
	return nil
}

// Delete removes a key from the store.
func (r *RemoteStore) Delete(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.vars[key]; !ok {
		return fmt.Errorf("remote store %q: key %q not found", r.namespace, key)
	}
	delete(r.vars, key)
	return nil
}

// All returns a snapshot of all key/value pairs.
func (r *RemoteStore) All() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copy := make(map[string]string, len(r.vars))
	for k, v := range r.vars {
		copy[k] = v
	}
	return copy
}

// Keys returns sorted keys in the store.
func (r *RemoteStore) Keys() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	keys := make([]string, 0, len(r.vars))
	for k := range r.vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
