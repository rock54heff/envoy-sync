package store

import (
	"fmt"
	"os"
	"strings"
)

// Flusher is implemented by stores that can persist their state to a file.
type Flusher interface {
	Flush() error
}

// WriteEnvFile serialises vars into standard KEY=VALUE format and writes
// them to path, creating or truncating the file as needed.
func WriteEnvFile(path string, vars map[string]string) error {
	if path == "" {
		return fmt.Errorf("WriteEnvFile: path must not be empty")
	}

	keys := sortedStringKeys(vars)
	var sb strings.Builder
	for _, k := range keys {
		v := vars[k]
		// Quote values that contain spaces or special characters.
		if needsQuoting(v) {
			v = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(v)
		sb.WriteByte('\n')
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func needsQuoting(v string) bool {
	return strings.ContainsAny(v, " \t\n#")
}

func sortedStringKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// simple insertion sort — maps are small
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
