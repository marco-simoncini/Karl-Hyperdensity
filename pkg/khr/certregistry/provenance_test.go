package certregistry

import (
	"testing"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/nativelive"
)

func TestVerifyIntegrityFingerprintMismatch(t *testing.T) {
	evidence := []byte(`{"status":"certified"}`)
	summary := nativelive.CertificationSummary{
		Lane: nativelive.LaneNativeLive, Status: nativelive.CertificationCertified,
		Invariants: nativelive.Invariants{NoRestart: true, NoRollout: true, NoRecreate: true},
	}
	reg := GenerateFromSummaryWithEvidence("KHR-Y", summary, "evidence/cert.json", 3600,
		time.Now().UTC(), evidence, "karl-metal-01@ovh")
	if err := VerifyIntegrity(reg, []byte(`{"tampered":true}`)); err == nil {
		t.Fatal("expected fingerprint mismatch")
	}
}

func TestVerifyIntegrityPass(t *testing.T) {
	evidence := []byte(`{"status":"certified"}`)
	summary := nativelive.CertificationSummary{
		Lane: lanediscovery.LaneNativeLive, Status: nativelive.CertificationCertified,
		Invariants: nativelive.Invariants{NoRestart: true, NoRollout: true, NoRecreate: true},
	}
	reg := GenerateFromSummaryWithEvidence("KHR-Y", summary, "evidence/cert.json", 3600,
		time.Now().UTC(), evidence, "karl-metal-01@ovh")
	if err := VerifyIntegrity(reg, evidence); err != nil {
		t.Fatal(err)
	}
}
