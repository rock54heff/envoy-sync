package store

import (
	"fmt"

	"github.com/user/envoy-sync/internal/envfile"
)

// FileStore is a Store backed by a .env file on disk.
type FileStore struct {
	path string
	vars map[string]string
}

// NewFileStore loads a .env file from the given path and returns a FileStore.
func NewFileStore(path string) (*FileStore, error) {
	vars, err := envfile.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("file store: parsing %q: %w", path, err)
	}
	return &FileStore{path: path, vars: vars}, nil
}

// Get returns the value for key, and whether it was found.
func (f *FileStore) Get(key string) (string, bool) {
	v, ok := f.vars[key]
	return v, ok
}

// Set sets key to value in the in-memory map.
// Note: does not persist to disk; use Save to write changes.
func (f *FileStore) Set(key, value string) {
	f.vars[key] = value
}

// Delete removes key from the in-memory map.
func (f *FileStore) Delete(key string) {
	delete(f.vars, key)
}

// All returns a copy of all key-value pairs.
func (f *FileStore) All() map[string]string {
	out := make(map[string]string, len(f.vars))
	for k, v := range f.vars {
		out[k] = v
	}
	return out
}

// Path returns the file path backing this store.
func (f *FileStore) Path() string {
	return f.path
}
