package store

import (
	"errors"
	"fmt"
	"strings"
)

// AWSStore simulates an AWS SSM Parameter Store-backed secret store.
// In production this would use the AWS SDK; here we use an in-memory
// map to keep the implementation self-contained and testable.
type AWSStore struct {
	namespace string
	vars      map[string]string
}

// NewAWSStore creates an AWSStore pre-populated with the given vars
// under the given namespace prefix (e.g. "/myapp/prod").
func NewAWSStore(namespace string, vars map[string]string) (*AWSStore, error) {
	if namespace == "" {
		return nil, errors.New("aws store: namespace must not be empty")
	}
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &AWSStore{namespace: namespace, vars: copy}, nil
}

// Namespace returns the SSM path prefix used by this store.
func (a *AWSStore) Namespace() string { return a.namespace }

func (a *AWSStore) fullKey(key string) string {
	return strings.TrimRight(a.namespace, "/") + "/" + key
}

// Get retrieves a parameter by key. Returns an error when not found.
func (a *AWSStore) Get(key string) (string, error) {
	v, ok := a.vars[a.fullKey(key)]
	if !ok {
		return "", fmt.Errorf("aws store: key %q not found in namespace %q", key, a.namespace)
	}
	return v, nil
}

// Set writes a parameter value under the namespaced key.
func (a *AWSStore) Set(key, value string) error {
	if key == "" {
		return errors.New("aws store: key must not be empty")
	}
	a.vars[a.fullKey(key)] = value
	return nil
}

// Delete removes a parameter from the store.
func (a *AWSStore) Delete(key string) error {
	delete(a.vars, a.fullKey(key))
	return nil
}

// ToMap returns a flat map of bare keys (namespace prefix stripped) to values.
func (a *AWSStore) ToMap() map[string]string {
	prefix := strings.TrimRight(a.namespace, "/") + "/"
	out := make(map[string]string, len(a.vars))
	for k, v := range a.vars {
		bare := strings.TrimPrefix(k, prefix)
		out[bare] = v
	}
	return out
}
