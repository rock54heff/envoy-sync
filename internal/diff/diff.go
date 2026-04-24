package diff

import "sort"

// Result holds the diff between two env maps.
type Result struct {
	Added    map[string]string // keys present in b but not in a
	Removed  map[string]string // keys present in a but not in b
	Modified map[string]string // keys present in both but with different values (value is b's value)
	Unchanged map[string]string // keys present in both with same values
}

// Compare computes the diff between env map a (source) and env map b (target).
func Compare(a, b map[string]string) Result {
	res := Result{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Modified:  make(map[string]string),
		Unchanged: make(map[string]string),
	}

	for k, va := range a {
		if vb, ok := b[k]; ok {
			if va == vb {
				res.Unchanged[k] = va
			} else {
				res.Modified[k] = vb
			}
		} else {
			res.Removed[k] = va
		}
	}

	for k, vb := range b {
		if _, ok := a[k]; !ok {
			res.Added[k] = vb
		}
	}

	return res
}

// HasChanges returns true if there are any added, removed, or modified keys.
func (r Result) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Modified) > 0
}

// SortedKeys returns a sorted slice of keys from the given map.
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
