package store

import "fmt"

// NetlifyStore holds environment variables scoped to a Netlify site.
type NetlifyStore struct {
	vars      map[string]string
	siteID    string
	namespace string
}

// NewNetlifyStore creates a NetlifyStore pre-populated with the given vars.
func NewNetlifyStore(vars map[string]string, siteID, namespace string) *NetlifyStore {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &NetlifyStore{vars: copy, siteID: siteID, namespace: namespace}
}

// Name returns a human-readable identifier for this store.
func (s *NetlifyStore) Name() string {
	if s.siteID != "" {
		return fmt.Sprintf("netlify:%s", s.siteID)
	}
	return "netlify"
}

// Namespace returns the scoping namespace (e.g. "production", "deploy-preview").
func (s *NetlifyStore) Namespace() string {
	return s.namespace
}

// Get returns the value for key, and whether it was found.
func (s *NetlifyStore) Get(key string) (string, bool) {
	v, ok := s.vars[key]
	return v, ok
}

// Set adds or updates a key in the store.
func (s *NetlifyStore) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("netlify: key must not be empty")
	}
	s.vars[key] = value
	return nil
}

// Delete removes a key from the store.
func (s *NetlifyStore) Delete(key string) {
	delete(s.vars, key)
}

// All returns a copy of all stored variables.
func (s *NetlifyStore) All() map[string]string {
	copy := make(map[string]string, len(s.vars))
	for k, v := range s.vars {
		copy[k] = v
	}
	return copy
}
