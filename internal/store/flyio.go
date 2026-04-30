package store

import "fmt"

// FlyioStore represents a Fly.io secrets store for a given app.
type FlyioStore struct {
	vars    map[string]string
	appName string
}

// NewFlyioStore creates a FlyioStore pre-loaded with the provided vars.
func NewFlyioStore(appName string, vars map[string]string) *FlyioStore {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &FlyioStore{vars: copy, appName: appName}
}

// Name returns a human-readable identifier for the store.
func (s *FlyioStore) Name() string {
	if s.appName == "" {
		return "fly.io"
	}
	return fmt.Sprintf("fly.io(%s)", s.appName)
}

// Namespace returns the app-scoped key prefix.
func (s *FlyioStore) Namespace() string {
	if s.appName == "" {
		return "flyio"
	}
	return fmt.Sprintf("flyio/%s", s.appName)
}

// Get returns the value for key, or an error if not found.
func (s *FlyioStore) Get(key string) (string, error) {
	v, ok := s.vars[key]
	if !ok {
		return "", fmt.Errorf("flyio: key %q not found", key)
	}
	return v, nil
}

// Set stores a key-value pair.
func (s *FlyioStore) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("flyio: key must not be empty")
	}
	s.vars[key] = value
	return nil
}

// Delete removes a key from the store.
func (s *FlyioStore) Delete(key string) error {
	delete(s.vars, key)
	return nil
}

// ToMap returns a shallow copy of all stored vars.
func (s *FlyioStore) ToMap() map[string]string {
	out := make(map[string]string, len(s.vars))
	for k, v := range s.vars {
		out[k] = v
	}
	return out
}
