package integrity

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

const (
	// SignatureAlgorithmLocalDev labels dev-only Ed25519 signatures (not production PKI).
	SignatureAlgorithmLocalDev = "ED25519-local-dev"
)

// SignLocalDev signs canonical bundle bytes with an Ed25519 private key from PEM (PKCS#8).
// Intended for developer workflows only; does not replace production code signing or admission.
func SignLocalDev(canonical []byte, pemPath string) ([]byte, error) {
	raw, err := os.ReadFile(pemPath)
	if err != nil {
		return nil, fmt.Errorf("read signing key: %w", err)
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("signing key PEM decode failed")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse PKCS8 private key: %w", err)
	}
	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not Ed25519 (got %T)", key)
	}
	if len(priv) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid Ed25519 private key length")
	}
	return ed25519.Sign(priv, canonical), nil
}

// ErrLocalDevRequiresKeyFile is returned when signing-mode is local-dev without -signing-key-file.
var ErrLocalDevRequiresKeyFile = fmt.Errorf("collect-evidence: signing-mode=local-dev requires -signing-key-file (dev-only; not production security)")

// RequireLocalDevKey returns ErrLocalDevRequiresKeyFile when mode is local-dev and path is empty.
func RequireLocalDevKey(signingMode, keyFile string) error {
	if NormalizeSigningMode(signingMode) != "local-dev" {
		return nil
	}
	if strings.TrimSpace(keyFile) == "" {
		return ErrLocalDevRequiresKeyFile
	}
	return nil
}
