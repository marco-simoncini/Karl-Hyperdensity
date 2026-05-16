package certregistry

import (
	"testing"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/nativelive"
)

func certifiedSummary() nativelive.CertificationSummary {
	return nativelive.CertificationSummary{
		CertificationID: nativelive.CertificationID,
		Lane:            nativelive.LaneNativeLive,
		Status:          nativelive.CertificationCertified,
		Invariants: nativelive.Invariants{
			NoRestart: true, NoRollout: true, NoRecreate: true,
		},
		Metrics: nativelive.CertificationMetrics{
			RollbackLatencyMs: nativelive.RollbackLatencyMs{CPU: 5},
		},
		Scores: nativelive.CertificationScores{
			ContinuityScore: 1, LiveScaleConfidence: nativelive.ConfidenceHigh,
		},
		ContinuityProof: nativelive.ContinuityCertificationProof{
			ShellContinuityPreserved: true,
		},
	}
}

func TestGenerateFromSummary(t *testing.T) {
	now := time.Date(2026, 5, 16, 17, 0, 0, 0, time.UTC)
	reg := GenerateFromSummary("KHR-V", certifiedSummary(), "evidence/cert", 3600, now)
	if len(reg.Entries) != 1 || reg.Entries[0].CertificationState != CertStateCertified {
		t.Fatalf("reg=%+v", reg)
	}
}

func TestIsFreshStale(t *testing.T) {
	entry := LaneEntry{
		LastCertifiedAt: "2026-05-01T00:00:00Z",
		ValidForSeconds: 3600,
	}
	now := time.Date(2026, 5, 16, 0, 0, 0, 0, time.UTC)
	if IsFresh(entry, now) {
		t.Fatal("expected stale")
	}
}
