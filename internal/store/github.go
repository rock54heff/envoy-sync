package store

import (
	"fmt"
	"sort"
)

// GitHubStore represents a store backed by GitHub Actions secrets (simulated in-memory for testing).
type GitHubStore struct {
	vars      map[string]string
	namespace string
	repo      string
}

// NewGitHubStore creates a new GitHubStore with the given repo, namespace, and initial vars.
func NewGitHubStore(repo, namespace string, vars map[string]string) *GitHubStore {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &GitHubStore{
		vars:      copy,
		namespace: namespace,
		repo:      repo,
	}
}

func (g *GitHubStore) Name() string {
	return fmt.Sprintf("github:%s", g.repo)
}

func (g *GitHubStore) Namespace() string {
	return g.namespace
}

func (g *GitHubStore) Get(key string) (string, bool) {
	v, ok := g.vars[key]
	return v, ok
}

func (g *GitHubStore) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}
	g.vars[key] = value
	return nil
}

func (g *GitHubStore) Delete(key string) error {
	delete(g.vars, key)
	return nil
}

func (g *GitHubStore) Keys() []string {
	keys := make([]string, 0, len(g.vars))
	for k := range g.vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (g *GitHubStore) ToMap() map[string]string {
	copy := make(map[string]string, len(g.vars))
	for k, v := range g.vars {
		copy[k] = v
	}
	return copy
}
