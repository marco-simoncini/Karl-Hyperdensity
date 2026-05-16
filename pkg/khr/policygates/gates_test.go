package policygates

import (
	"testing"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
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

func freshRegistry() *certregistry.Registry {
	reg := certregistry.GenerateFromSummary("KHR-V", certifiedSummary(), "evidence/x", 86400, time.Now().UTC())
	return &reg
}

func TestEvaluateEligibleCertifiedLane(t *testing.T) {
	reg := freshRegistry()
	out := Evaluate(lanediscovery.LaneNativeLive, reg, DefaultNativeLiveGates(), time.Now().UTC())
	if !out.Eligible || out.Blocked || out.UncertifiedLane || out.StaleEvidence {
		t.Fatalf("out=%+v", out)
	}
}

func TestEvaluateUncertifiedLaneBlocks(t *testing.T) {
	reg := freshRegistry()
	out := Evaluate(lanediscovery.LaneLinuxVMCompatibility, reg, DefaultNativeLiveGates(), time.Now().UTC())
	if !out.UncertifiedLane || out.Eligible {
		t.Fatalf("out=%+v", out)
	}
}

func TestEvaluateStaleEvidenceBlocks(t *testing.T) {
	reg := freshRegistry()
	reg.Entries[0].LastCertifiedAt = "2020-01-01T00:00:00Z"
	out := Evaluate(lanediscovery.LaneNativeLive, reg, DefaultNativeLiveGates(), time.Now().UTC())
	if !out.StaleEvidence || out.Eligible {
		t.Fatalf("out=%+v", out)
	}
}

func TestEvaluateFailedGateBlocks(t *testing.T) {
	reg := freshRegistry()
	reg.Entries[0].Attestation.NoRestart = false
	out := Evaluate(lanediscovery.LaneNativeLive, reg, DefaultNativeLiveGates(), time.Now().UTC())
	if out.Eligible || len(out.FailedGates) == 0 {
		t.Fatalf("out=%+v", out)
	}
}
