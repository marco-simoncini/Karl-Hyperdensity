package actionapproval

import (
	"testing"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/nativelive"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/policygates"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/provenance"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcefuture"
)

func certifiedRegistryWithProvenance() certregistry.Registry {
	evidence := []byte(`{"status":"certified"}`)
	now := time.Now().UTC()
	return certregistry.GenerateFromSummaryWithEvidence("KHR-Y", nativelive.CertificationSummary{
		CertificationID: nativelive.CertificationID,
		Lane:            nativelive.LaneNativeLive,
		Status:          nativelive.CertificationCertified,
		Invariants:      nativelive.Invariants{NoRestart: true, NoRollout: true, NoRecreate: true},
		Metrics:         nativelive.CertificationMetrics{RollbackLatencyMs: nativelive.RollbackLatencyMs{CPU: 1}},
		Scores:          nativelive.CertificationScores{ContinuityScore: 1},
		ContinuityProof: nativelive.ContinuityCertificationProof{ShellContinuityPreserved: true},
	}, "evidence/cert", 86400, now, evidence, "karl-metal-01@ovh")
}

func certifiedRegistry() certregistry.Registry {
	now := time.Now().UTC()
	return certregistry.GenerateFromSummary("KHR-W", nativelive.CertificationSummary{
		CertificationID: nativelive.CertificationID,
		Lane:            nativelive.LaneNativeLive,
		Status:          nativelive.CertificationCertified,
		Invariants:      nativelive.Invariants{NoRestart: true, NoRollout: true, NoRecreate: true},
		Metrics:         nativelive.CertificationMetrics{RollbackLatencyMs: nativelive.RollbackLatencyMs{CPU: 1}},
		Scores:          nativelive.CertificationScores{ContinuityScore: 1},
		ContinuityProof: nativelive.ContinuityCertificationProof{ShellContinuityPreserved: true},
	}, "evidence/cert", 86400, now)
}

func eligibleSimulation() resourcefuture.SimulationResult {
	return resourcefuture.SimulationResult{
		LiveInPlaceEligibility: []resourcefuture.LiveInPlaceEligibility{{
			TargetRef: "khr-runtime-sandbox/Cell/khr-native-live-target",
			Lane:      lanediscovery.LaneNativeLive,
			Eligible:  true,
		}},
	}
}

func TestGeneratePendingFromEligible(t *testing.T) {
	now := time.Now().UTC()
	pending, err := GeneratePending(GenerateInput{
		Simulation: eligibleSimulation(), Registry: certifiedRegistry(),
		Gates: policygates.DefaultNativeLiveGates(), CertificationRef: "evidence/cert",
		TTLSeconds: 3600, Now: now,
	})
	if err != nil || len(pending) != 1 || pending[0].ApprovalState != StatePending {
		t.Fatalf("pending=%+v err=%v", pending, err)
	}
}

func TestApproveCertifiedOnly(t *testing.T) {
	reg := certifiedRegistry()
	now := time.Now().UTC()
	pending, _ := GeneratePending(GenerateInput{
		Simulation: eligibleSimulation(), Registry: reg,
		Gates: policygates.DefaultNativeLiveGates(), CertificationRef: "evidence/cert", Now: now,
	})
	approved, err := Approve(pending[0], &reg, policygates.DefaultNativeLiveGates(), "operator-a", now)
	if err != nil || approved.ApprovalState != StateApproved {
		t.Fatalf("approved=%+v err=%v", approved, err)
	}
}

func TestRejectPath(t *testing.T) {
	reg := certifiedRegistry()
	now := time.Now().UTC()
	pending, _ := GeneratePending(GenerateInput{
		Simulation: eligibleSimulation(), Registry: reg,
		Gates: policygates.DefaultNativeLiveGates(), Now: now,
	})
	rejected, err := Reject(pending[0], "operator-b", "sandbox reject test", now)
	if err != nil || rejected.ApprovalState != StateRejected {
		t.Fatalf("rejected=%+v err=%v", rejected, err)
	}
}

func TestExpirePath(t *testing.T) {
	reg := certifiedRegistry()
	now := time.Now().UTC()
	pending, _ := GeneratePending(GenerateInput{
		Simulation: eligibleSimulation(), Registry: reg,
		Gates: policygates.DefaultNativeLiveGates(), Now: now,
	})
	expired := SimulateExpire(pending[0], now)
	if expired.ApprovalState != StateExpired {
		t.Fatalf("state=%q", expired.ApprovalState)
	}
}

func TestStaleCertificationBlocksApproval(t *testing.T) {
	reg := certifiedRegistry()
	reg.Entries[0].LastCertifiedAt = "2019-01-01T00:00:00Z"
	now := time.Now().UTC()
	pending, _ := GeneratePending(GenerateInput{
		Simulation: eligibleSimulation(), Registry: reg,
		Gates: policygates.DefaultNativeLiveGates(), Now: now,
	})
	if len(pending) > 0 {
		t.Fatal("stale registry should not generate eligible pending")
	}
	// manual pending with stale registry at approve time
	p := ActionApproval{
		ActionID: "test", LaneID: lanediscovery.LaneNativeLive,
		ApprovalState: StatePending,
		ExpiresAt:     now.Add(time.Hour).UTC().Format(time.RFC3339),
		PolicyGateResult: policygates.EligibilityOutcome{Eligible: true},
	}
	_, err := Approve(p, &reg, policygates.DefaultNativeLiveGates(), "op", now)
	if err == nil {
		t.Fatal("expected stale block")
	}
}

func TestInvalidApprovalProvenance(t *testing.T) {
	reg := certifiedRegistryWithProvenance()
	now := time.Now().UTC()
	pending, _ := GeneratePending(GenerateInput{
		Simulation: eligibleSimulation(), Registry: reg,
		Gates: policygates.DefaultNativeLiveGates(), Now: now,
	})
	pending[0].Provenance = provenance.NewRecord("tamper", provenance.SourceContext{
		Cluster: "other-cluster", Lane: lanediscovery.LaneNativeLive,
	}, "x", []byte("wrong"), now)
	_, err := Approve(pending[0], &reg, policygates.DefaultNativeLiveGates(), "op", now)
	if err == nil {
		t.Fatal("expected provenance mismatch on approve")
	}
}

func TestNoApplyFlags(t *testing.T) {
	reg := certifiedRegistry()
	pending, _ := GeneratePending(GenerateInput{
		Simulation: eligibleSimulation(), Registry: reg,
		Gates: policygates.DefaultNativeLiveGates(), Now: time.Now().UTC(),
	})
	if !pending[0].NoApply || !pending[0].NoMutation || !pending[0].NoAutonomousOrchestration {
		t.Fatalf("safety flags=%+v", pending[0])
	}
}
