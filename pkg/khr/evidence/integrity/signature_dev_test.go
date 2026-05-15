package integrity

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
)

func TestRequireLocalDevKeyError(t *testing.T) {
	if err := RequireLocalDevKey("local-dev", ""); err == nil {
		t.Fatal("expected error")
	}
	if err := RequireLocalDevKey("none", ""); err != nil {
		t.Fatal(err)
	}
}

func TestSignLocalDevRoundTrip(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "k.pem")
	if err := os.WriteFile(keyPath, pemBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	msg := []byte("canonical-bundle-bytes")
	sig, err := SignLocalDev(msg, keyPath)
	if err != nil {
		t.Fatal(err)
	}
	if !ed25519.Verify(pub, msg, sig) {
		t.Fatal("verify failed")
	}
}

func TestEmitLocalDevRequiresManifestPath(t *testing.T) {
	t.Setenv("KHR_TEST_INTEGRITY_NOW", "2026-07-01T12:00:00Z")
	dir := t.TempDir()
	err := EmitEvidenceSidecars(
		map[string]int{"a": 1},
		"ag",
		"",
		"",
		filepath.Join(dir, "d.txt"),
		"local-dev",
		filepath.Join(dir, "k.pem"),
	)
	if err == nil {
		t.Fatal("expected error for local-dev without manifest path")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("evidence-manifest-output")) {
		t.Fatalf("unexpected err: %v", err)
	}
}
