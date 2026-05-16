package provenance

import (
	"testing"
	"time"
)

func TestFingerprintMismatch(t *testing.T) {
	data := []byte(`{"status":"certified"}`)
	rec := NewRecord("khr-cert-registry", SourceContext{
		Cluster: "karl-metal-01@ovh", Namespace: "khr-runtime-sandbox", Lane: "native-live",
	}, "cert-1", data, time.Now().UTC())
	if err := VerifyFingerprint(rec, []byte(`{"status":"failed"}`)); err == nil {
		t.Fatal("expected fingerprint mismatch")
	}
}

func TestStaleProvenance(t *testing.T) {
	rec := NewRecord("test", SourceContext{Cluster: "c"}, "s", []byte("x"), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	if !IsStaleProvenance(rec, time.Now().UTC(), 3600) {
		t.Fatal("expected stale")
	}
}

func TestLineageMismatch(t *testing.T) {
	rec := NewRecord("test", SourceContext{Cluster: "c", Lane: "native-live"}, "s", []byte("x"), time.Now().UTC())
	if err := VerifyLineage(rec, "lineage:deadbeef"); err == nil {
		t.Fatal("expected lineage mismatch")
	}
}

func TestInvalidApprovalProvenance(t *testing.T) {
	cert := NewRecord("cert", SourceContext{Cluster: "c", Lane: "native-live"}, "c", []byte("cert"), time.Now().UTC())
	approval := NewRecord("approval", SourceContext{Cluster: "other", Lane: "native-live"}, "a", []byte("approval"), time.Now().UTC())
	if err := VerifyApprovalProvenance(approval, cert); err == nil {
		t.Fatal("expected approval provenance mismatch")
	}
}

func TestMatchSuccess(t *testing.T) {
	data := []byte("same")
	src := SourceContext{Cluster: "karl-metal-01@ovh", Namespace: "khr-runtime-sandbox", Lane: "native-live"}
	a := NewRecord("a", src, "seed", data, time.Now().UTC())
	b := NewRecord("b", src, "seed", data, time.Now().UTC())
	if !Match(a, b) {
		t.Fatal("expected match on same evidence and lineage")
	}
}
