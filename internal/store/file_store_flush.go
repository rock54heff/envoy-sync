package store

// Flush persists the current in-memory state of a FileStore back to the
// file it was loaded from. It satisfies the Flusher interface.
func (fs *FileStore) Flush() error {
	vars := ToMap(fs)
	return WriteEnvFile(fs.path, vars)
}
