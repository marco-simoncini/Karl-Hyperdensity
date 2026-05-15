// Package integrity provides local-only canonical serialization and digests for KHR evidence bundles (Sprint 10).
package integrity

import (
	"bytes"
	"encoding/json"
)

// CanonicalJSON returns compact, deterministic JSON bytes suitable for hashing.
// Rules: no HTML escaping in strings; no trailing newline; struct field order follows Go struct definitions.
func CanonicalJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return b, nil
}
