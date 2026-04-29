package store

import "fmt"

// GitLabStore simulates a GitLab CI/CD variable store (namespace-scoped).
type GitLabStore struct {
	vars      map[string]string
	namespace string
}

// NewGitLabStore creates a GitLabStore pre-populated with the given vars.
func NewGitLabStore(namespace string, vars map[string]string) *GitLabStore {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &GitLabStore{vars: copy, namespace: namespace}
}

func (s *GitLabStore) Name() string {
	return "gitlab"
}

func (s *GitLabStore) Namespace() string {
	return s.namespace
}

func (s *GitLabStore) Get(key string) (string, bool) {
	v, ok := s.vars[key]
	return v, ok
}

func (s *GitLabStore) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("gitlab: key must not be empty")
	}
	s.vars[key] = value
	return nil
}

func (s *GitLabStore) Delete(key string) error {
	delete(s.vars, key)
	return nil
}

func (s *GitLabStore) All() map[string]string {
	copy := make(map[string]string, len(s.vars))
	for k, v := range s.vars {
		copy[k] = v
	}
	return copy
}
