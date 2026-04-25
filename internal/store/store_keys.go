package store

// Keyer is implemented by any store that can enumerate its keys.
type Keyer interface {
	Keys() []string
}

// Getter is implemented by any store that supports key lookup.
type Getter interface {
	Get(key string) (string, bool)
}

// ReadableStore combines Keyer and Getter for read-only store access.
type ReadableStore interface {
	Keyer
	Getter
}

// ToMap converts any ReadableStore into a plain map[string]string.
func ToMap(s ReadableStore) map[string]string {
	result := make(map[string]string)
	for _, k := range s.Keys() {
		if v, ok := s.Get(k); ok {
			result[k] = v
		}
	}
	return result
}

// Equal returns true if two readable stores contain identical key-value pairs.
func Equal(a, b ReadableStore) bool {
	am := ToMap(a)
	bm := ToMap(b)
	if len(am) != len(bm) {
		return false
	}
	for k, av := range am {
		if bv, ok := bm[k]; !ok || av != bv {
			return false
		}
	}
	return true
}
