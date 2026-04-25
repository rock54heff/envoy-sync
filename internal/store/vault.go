package store

import (
	"fmt"
	"sort"
	"strings"
)

// VaultStore simulates a HashiCorp Vault-like secret store backed by a
// namespace-prefixed in-memory map. In a real implementation the HTTP
// client calls would live here.
type VaultStore struct {
	namespace string
	vars      map[string]string
}

// NewVaultStore creates a VaultStore pre-loaded with the given vars under
// the provided namespace prefix.
func NewVaultStore(namespace string, vars map[string]string) *VaultStore {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &VaultStore{namespace: namespace, vars: copy}
}

// Namespace returns the configured namespace / secret path prefix.
func (v *VaultStore) Namespace() string {
	return v.namespace
}

// Get returns the value for key, with ok=false when absent.
func (v *VaultStore) Get(key string) (string, bool) {
	val, ok := v.vars[v.qualifiedKey(key)]
	return val, ok
}

// Set stores a value under the namespaced key.
func (v *VaultStore) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("vault: key must not be empty")
	}
	v.vars[v.qualifiedKey(key)] = value
	return nil
}

// Delete removes a key from the store, returning an error if not found.
func (v *VaultStore) Delete(key string) error {
	qk := v.qualifiedKey(key)
	if _, ok := v.vars[qk]; !ok {
		return fmt.Errorf("vault: key %q not found", key)
	}
	delete(v.vars, qk)
	return nil
}

// All returns a flat map with namespace prefixes stripped from keys.
func (v *VaultStore) All() map[string]string {
	prefix := v.namespace + "/"
	out := make(map[string]string, len(v.vars))
	for k, val := range v.vars {
		bare := strings.TrimPrefix(k, prefix)
		out[bare] = val
	}
	return out
}

// SortedKeys returns bare (un-prefixed) keys in alphabetical order.
func (v *VaultStore) SortedKeys() []string {
	all := v.All()
	keys := make([]string, 0, len(all))
	for k := range all {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (v *VaultStore) qualifiedKey(key string) string {
	if v.namespace == "" {
		return key
	}
	return v.namespace + "/" + key
}
