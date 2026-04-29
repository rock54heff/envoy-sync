package store

import "fmt"

// VercelStore simulates a Vercel environment variable store.
// In production this would call the Vercel API; here it holds
// an in-memory map so the rest of the toolchain can work with it.
type VercelStore struct {
	vars      map[string]string
	namespace string
	teamID    string
}

// NewVercelStore creates a VercelStore pre-populated with vars.
// namespace is typically the Vercel project name.
// teamID is the Vercel team slug/id (may be empty for personal accounts).
func NewVercelStore(vars map[string]string, namespace, teamID string) *VercelStore {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &VercelStore{
		vars:      copy,
		namespace: namespace,
		teamID:    teamID,
	}
}

// Name returns a human-readable identifier for the store.
func (s *VercelStore) Name() string {
	if s.teamID != "" {
		return fmt.Sprintf("vercel:%s/%s", s.teamID, s.namespace)
	}
	return fmt.Sprintf("vercel:%s", s.namespace)
}

// Namespace returns the project name this store is scoped to.
func (s *VercelStore) Namespace() string { return s.namespace }

// Get returns the value for key and whether it was found.
func (s *VercelStore) Get(key string) (string, bool) {
	v, ok := s.vars[key]
	return v, ok
}

// Set stores a key/value pair.
func (s *VercelStore) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("vercel store: key must not be empty")
	}
	s.vars[key] = value
	return nil
}

// Delete removes a key from the store.
func (s *VercelStore) Delete(key string) error {
	delete(s.vars, key)
	return nil
}

// All returns a shallow copy of all stored variables.
func (s *VercelStore) All() map[string]string {
	out := make(map[string]string, len(s.vars))
	for k, v := range s.vars {
		out[k] = v
	}
	return out
}
