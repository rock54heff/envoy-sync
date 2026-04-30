package store

import "fmt"

// RailwayStore is an in-memory representation of Railway environment variables.
type RailwayStore struct {
	vars      map[string]string
	projectID string
	environment string
}

// NewRailwayStore creates a new RailwayStore with the given vars, projectID, and environment.
func NewRailwayStore(vars map[string]string, projectID, environment string) *RailwayStore {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &RailwayStore{
		vars:        copy,
		projectID:   projectID,
		environment: environment,
	}
}

// Name returns a human-readable identifier for this store.
func (r *RailwayStore) Name() string {
	if r.projectID != "" {
		return fmt.Sprintf("railway:%s", r.projectID)
	}
	return "railway"
}

// Namespace returns a scoped identifier including the environment.
func (r *RailwayStore) Namespace() string {
	if r.environment != "" {
		return fmt.Sprintf("%s/%s", r.Name(), r.environment)
	}
	return r.Name()
}

// Get returns the value for the given key, or an error if not found.
func (r *RailwayStore) Get(key string) (string, error) {
	v, ok := r.vars[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in railway store", key)
	}
	return v, nil
}

// Set sets or updates a key-value pair in the store.
func (r *RailwayStore) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}
	r.vars[key] = value
	return nil
}

// Delete removes a key from the store.
func (r *RailwayStore) Delete(key string) error {
	delete(r.vars, key)
	return nil
}

// All returns a copy of all key-value pairs in the store.
func (r *RailwayStore) All() map[string]string {
	copy := make(map[string]string, len(r.vars))
	for k, v := range r.vars {
		copy[k] = v
	}
	return copy
}
