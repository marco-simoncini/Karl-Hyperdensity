package resourcefuture

import (
	"testing"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/nativelive"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/policygates"
)

func testCertifiedRegistry() *certregistry.Registry {
	summary := nativelive.CertificationSummary{
		CertificationID: nativelive.CertificationID,
		Lane:            nativelive.LaneNativeLive,
		Status:          nativelive.CertificationCertified,
		Invariants: nativelive.Invariants{
			NoRestart: true, NoRollout: true, NoRecreate: true,
		},
		Metrics: nativelive.CertificationMetrics{
			RollbackLatencyMs: nativelive.RollbackLatencyMs{CPU: 1},
		},
		Scores: nativelive.CertificationScores{ContinuityScore: 1},
		ContinuityProof: nativelive.ContinuityCertificationProof{
			ShellContinuityPreserved: true,
		},
	}
	reg := certregistry.GenerateFromSummary("KHR-V", summary, "evidence/test", 86400, time.Now().UTC())
	return &reg
}

func TestEvaluateLiveInPlace_EligibleCertifiedLane(t *testing.T) {
	policy := PolicyContext{
		Registry: testCertifiedRegistry(),
		Gates:    policygates.DefaultNativeLiveGates(),
		Now:      time.Now().UTC(),
	}
	ok, _, state, _, stale, uncert := evaluateLiveInPlace(
		lanediscovery.LaneNativeLive, lanediscovery.ClassificationNativeLive, true, policy,
	)
	if !ok || state != policygates.EligibilityEligible || stale || uncert {
		t.Fatalf("ok=%v state=%s stale=%v uncert=%v", ok, state, stale, uncert)
	}
}

func TestEvaluateLiveInPlace_StaleEvidenceBlocks(t *testing.T) {
	reg := testCertifiedRegistry()
	reg.Entries[0].LastCertifiedAt = "2019-01-01T00:00:00Z"
	policy := PolicyContext{Registry: reg, Gates: policygates.DefaultNativeLiveGates(), Now: time.Now().UTC()}
	ok, _, _, _, stale, _ := evaluateLiveInPlace(lanediscovery.LaneNativeLive, lanediscovery.ClassificationNativeLive, true, policy)
	if ok || !stale {
		t.Fatal("expected stale block")
	}
}

func TestEvaluateLiveInPlace_UncertifiedCompatibilityBlocks(t *testing.T) {
	policy := PolicyContext{
		Registry: testCertifiedRegistry(),
		Gates:    policygates.DefaultNativeLiveGates(),
		Now:      time.Now().UTC(),
	}
	ok, _, _, _, _, uncert := evaluateLiveInPlace(
		lanediscovery.LaneLinuxVMCompatibility, lanediscovery.ClassificationCompatibilityFallback, true, policy,
	)
	if ok || !uncert {
		t.Fatal("expected uncertified compatibility block")
	}
}

func TestEvaluateLiveInPlace_FailedGateBlocks(t *testing.T) {
	reg := testCertifiedRegistry()
	reg.Entries[0].Attestation.ShellContinuityPreserved = false
	policy := PolicyContext{Registry: reg, Gates: policygates.DefaultNativeLiveGates(), Now: time.Now().UTC()}
	ok, reason, state, _, _, _ := evaluateLiveInPlace(lanediscovery.LaneNativeLive, lanediscovery.ClassificationNativeLive, true, policy)
	if ok || state != policygates.EligibilityBlocked || reason == "" {
		t.Fatalf("ok=%v reason=%q state=%s", ok, reason, state)
	}
}
