package integrity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildManifestWithoutSignature(t *testing.T) {
	t.Setenv("KHR_TEST_INTEGRITY_NOW", "2026-07-01T12:00:00Z")
	t.Setenv("KHR_TEST_INTEGRITY_CHAIN_STUB", "1")
	canonical := []byte(`{"k":"v"}`)
	m := BuildManifest("agent-1", "artifact-1", "none", canonical, SHA256Hex(canonical), "", "")
	if m.SignaturePresent {
		t.Fatal("expected no signature")
	}
	if m.SigningMode != "none" {
		t.Fatalf("signingMode %q", m.SigningMode)
	}
	if m.SourceMode != "collect-evidence" || !m.MutationsForbidden {
		t.Fatalf("manifest contract fields: %+v", m)
	}
	if m.BundleBytes != len(canonical) {
		t.Fatalf("bundleBytes")
	}
	if m.BundleSha256 != SHA256Hex(canonical) {
		t.Fatal("sha mismatch")
	}
}

func TestWriteManifestAndDigestFiles(t *testing.T) {
	t.Setenv("KHR_TEST_INTEGRITY_NOW", "2026-07-01T12:00:00Z")
	t.Setenv("KHR_TEST_INTEGRITY_CHAIN_STUB", "1")
	dir := t.TempDir()
	canonical := []byte(`{"x":1}`)
	sha := SHA256Hex(canonical)
	m := BuildManifest("a", "", "none", canonical, sha, "", "")
	if err := WriteDigestFile(filepath.Join(dir, "d.txt"), sha); err != nil {
		t.Fatal(err)
	}
	if err := WriteManifestFile(filepath.Join(dir, "m.json"), m); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(dir, "d.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != sha+"\n" {
		t.Fatalf("digest file: %q", raw)
	}
	var got ArtifactManifest
	if err := json.Unmarshal(mustRead(t, filepath.Join(dir, "m.json")), &got); err != nil {
		t.Fatal(err)
	}
	if got.BundleSha256 != sha {
		t.Fatal(got.BundleSha256)
	}
}

func mustRead(t *testing.T, p string) []byte {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
