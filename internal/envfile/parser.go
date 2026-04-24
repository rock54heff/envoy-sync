package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a set of environment variable key-value pairs.
type EnvMap map[string]string

// Parse reads a .env file and returns an EnvMap of its contents.
// It supports KEY=VALUE pairs, ignores blank lines and comments (#).
func Parse(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envfile: open %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("envfile: %q line %d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envfile: scanning %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a single "KEY=VALUE" line into its components.
// Inline comments and surrounding quotes are stripped from values.
func parseLine(line string) (string, string, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 1 {
		return "", "", fmt.Errorf("invalid line %q: expected KEY=VALUE", line)
	}

	key := strings.TrimSpace(line[:idx])
	raw := strings.TrimSpace(line[idx+1:])

	// Strip inline comment
	if ci := strings.Index(raw, " #"); ci != -1 {
		raw = strings.TrimSpace(raw[:ci])
	}

	// Strip surrounding quotes (single or double)
	value := stripQuotes(raw)

	if key == "" {
		return "", "", fmt.Errorf("empty key in line %q", line)
	}

	return key, value, nil
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
