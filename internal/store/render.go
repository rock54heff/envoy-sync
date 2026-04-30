package store

import (
	"fmt"
	"sort"
)

// RenderStore is a store backed by Render environment groups.
type RenderStore struct {
	vars      map[string]string
	serviceID string
	apiKey    string
}

// NewRenderStore creates a RenderStore pre-populated with the given vars.
func NewRenderStore(serviceID, apiKey string, vars map[string]string) *RenderStore {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &RenderStore{
		vars:      copy,
		serviceID: serviceID,
		apiKey:    apiKey,
	}
}

// Name returns a human-readable identifier for this store.
func (r *RenderStore) Name() string {
	if r.serviceID != "" {
		return fmt.Sprintf("render(%s)", r.serviceID)
	}
	return "render"
}

// Namespace returns a unique key scoped to the service.
func (r *RenderStore) Namespace() string {
	return fmt.Sprintf("render/%s", r.serviceID)
}

// Get returns the value for key, and whether it was found.
func (r *RenderStore) Get(key string) (string, bool) {
	v, ok := r.vars[key]
	return v, ok
}

// Set sets a key-value pair in the store.
func (r *RenderStore) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("render: key must not be empty")
	}
	r.vars[key] = value
	return nil
}

// Delete removes a key from the store.
func (r *RenderStore) Delete(key string) error {
	delete(r.vars, key)
	return nil
}

// All returns a sorted copy of all key-value pairs.
func (r *RenderStore) All() map[string]string {
	out := make(map[string]string, len(r.vars))
	for k, v := range r.vars {
		out[k] = v
	}
	return out
}

// Keys returns sorted keys from the store.
func (r *RenderStore) SortedKeys() []string {
	keys := make([]string, 0, len(r.vars))
	for k := range r.vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
