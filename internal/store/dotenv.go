package store

import (
	"fmt"
	"strings"
)

// DotenvStore is an in-memory store backed by a named .env source.
// It supports namespacing keys with a prefix.
type DotenvStore struct {
	name      string
	namespace string
	vars      map[string]string
}

// NewDotenvStore creates a DotenvStore with the given name, namespace prefix,
// and initial variables. The vars map is copied defensively.
func NewDotenvStore(name, namespace string, vars map[string]string) *DotenvStore {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &DotenvStore{name: name, namespace: namespace, vars: copy}
}

// Name returns the human-readable name of this store.
func (d *DotenvStore) Name() string { return d.name }

// Namespace returns the key prefix used by this store.
func (d *DotenvStore) Namespace() string { return d.namespace }

// prefixed prepends the namespace (if any) to the key.
func (d *DotenvStore) prefixed(key string) string {
	if d.namespace == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", d.namespace, key)
}

// Get retrieves a value by key, applying the namespace prefix.
func (d *DotenvStore) Get(key string) (string, bool) {
	v, ok := d.vars[d.prefixed(key)]
	return v, ok
}

// Set stores a value under the namespaced key.
func (d *DotenvStore) Set(key, value string) {
	if strings.TrimSpace(key) == "" {
		return
	}
	d.vars[d.prefixed(key)] = value
}

// Delete removes a key from the store.
func (d *DotenvStore) Delete(key string) {
	delete(d.vars, d.prefixed(key))
}

// All returns a copy of all key-value pairs in the store.
func (d *DotenvStore) All() map[string]string {
	out := make(map[string]string, len(d.vars))
	for k, v := range d.vars {
		out[k] = v
	}
	return out
}
