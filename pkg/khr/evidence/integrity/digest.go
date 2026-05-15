package integrity

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256Hex returns lowercase hex of SHA-256 over b.
func SHA256Hex(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
