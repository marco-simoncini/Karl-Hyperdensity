package windowsfluidvirt

import (
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestGovernanceFixtureMatrix(t *testing.T) {
	fixtures := []string{
		"governance-master-win11-cpu.contract-prepared.json",
		"governance-master-win11-stale.needs-revalidation.json",
		"governance-master-win11-missing-rollback.blocked.json",
		"governance-master-win11-missing-return.blocked.json",
		"governance-master-win11-identity-change.quarantined.json",
		"governance-win11-pool.denied.json",
		"governance-generic-windows.denied.json",
		"governance-ram-memory-safety-missing.blocked.json",
	}
	for _, fixtureName := range fixtures {
		t.Run(fixtureName, func(t *testing.T) {
			fixturePath := governanceFixtureAbsPath(t, fixtureName)
			fixture, err := LoadGovernanceReplayFixture(fixturePath)
			if err != nil {
				t.Fatalf("load governance fixture: %v", err)
			}
			result := EvaluateWindowsFluidApplyGovernance(ApplyGovernanceEvaluationInput{
				AdmissionDecision: fixture.AdmissionDecision,
				Bundle:            fixture.Bundle,
				PolicyPack:        fixture.PolicyPack,
				RequestedAction:   fixture.RequestedAction,
				EvaluationTime:    time.Date(2026, 5, 7, 14, 45, 0, 0, time.UTC),
			})
			if result.FinalGovernancePhase != fixture.ExpectedGovernancePhase {
				t.Fatalf("phase mismatch expected=%s got=%s blockers=%v denial=%v",
					fixture.ExpectedGovernancePhase, result.FinalGovernancePhase, result.GovernanceContract.Blockers, result.GovernanceContract.DenialReasons)
			}
			for _, expected := range fixture.ExpectedBlockers {
				assertHas(t, result.GovernanceContract.Blockers, expected)
			}
			if result.GovernanceContract.MutationAllowed {
				t.Fatal("mutationAllowed must always be false")
			}
			if result.GovernanceContract.ApplyAllowed {
				t.Fatal("applyAllowed must always be false")
			}
		})
	}
}

func TestGovernancePreparedNeverMeansApply(t *testing.T) {
	fixture, err := LoadGovernanceReplayFixture(governanceFixtureAbsPath(t, "governance-master-win11-cpu.contract-prepared.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	result := EvaluateWindowsFluidApplyGovernance(ApplyGovernanceEvaluationInput{
		AdmissionDecision: fixture.AdmissionDecision,
		Bundle:            fixture.Bundle,
		RequestedAction:   fixture.RequestedAction,
		EvaluationTime:    time.Date(2026, 5, 7, 14, 45, 0, 0, time.UTC),
	})
	if result.FinalGovernancePhase != GovernanceContractPrepared {
		t.Fatalf("expected CONTRACT_PREPARED got %s", result.FinalGovernancePhase)
	}
	if result.GovernanceContract.ApplyAllowed {
		t.Fatal("contract prepared must not enable apply")
	}
	if result.GovernanceContract.RuntimeMode != "in-place-qmp" {
		t.Fatalf("unexpected runtime mode %s", result.GovernanceContract.RuntimeMode)
	}
}

func TestGovernanceP0AndP1PriorityBehavior(t *testing.T) {
	identityFixture, err := LoadGovernanceReplayFixture(governanceFixtureAbsPath(t, "governance-master-win11-identity-change.quarantined.json"))
	if err != nil {
		t.Fatalf("load identity fixture: %v", err)
	}
	identityResult := EvaluateWindowsFluidApplyGovernance(ApplyGovernanceEvaluationInput{
		AdmissionDecision: identityFixture.AdmissionDecision,
		Bundle:            identityFixture.Bundle,
		RequestedAction:   identityFixture.RequestedAction,
		EvaluationTime:    time.Date(2026, 5, 7, 14, 45, 0, 0, time.UTC),
	})
	if identityResult.FinalGovernancePhase != GovernanceContractQuarantined {
		t.Fatalf("expected quarantine with P0 blockers, got %s", identityResult.FinalGovernancePhase)
	}

	rollbackFixture, err := LoadGovernanceReplayFixture(governanceFixtureAbsPath(t, "governance-master-win11-missing-rollback.blocked.json"))
	if err != nil {
		t.Fatalf("load rollback fixture: %v", err)
	}
	rollbackResult := EvaluateWindowsFluidApplyGovernance(ApplyGovernanceEvaluationInput{
		AdmissionDecision: rollbackFixture.AdmissionDecision,
		Bundle:            rollbackFixture.Bundle,
		RequestedAction:   rollbackFixture.RequestedAction,
		EvaluationTime:    time.Date(2026, 5, 7, 14, 45, 0, 0, time.UTC),
	})
	if rollbackResult.FinalGovernancePhase != GovernanceContractBlocked {
		t.Fatalf("expected blocked with P1 blockers, got %s", rollbackResult.FinalGovernancePhase)
	}
}

func TestGovernanceInvariantsAndRevalidationContract(t *testing.T) {
	fixture, err := LoadGovernanceReplayFixture(governanceFixtureAbsPath(t, "governance-master-win11-cpu.contract-prepared.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	result := EvaluateWindowsFluidApplyGovernance(ApplyGovernanceEvaluationInput{
		AdmissionDecision: fixture.AdmissionDecision,
		Bundle:            fixture.Bundle,
		RequestedAction:   fixture.RequestedAction,
		EvaluationTime:    time.Date(2026, 5, 7, 14, 45, 0, 0, time.UTC),
	})
	if !result.RuntimeInvariantSet.QMPReadOnlyUntilApplyPhaseInvariant.Passed {
		t.Fatal("qmpReadOnlyUntilApplyPhaseInvariant must stay true in this phase")
	}
	if !result.PreApplyRevalidation.RequiredFreshEvidence.KubeVirtIdentity ||
		!result.PreApplyRevalidation.RequiredFreshEvidence.QMPEvidence ||
		!result.PreApplyRevalidation.RequiredFreshEvidence.GuestEvidence {
		t.Fatal("pre-apply revalidation must require fresh identity/qmp/guest evidence")
	}
}

func TestPolicyAttestationModeAndSignatureSafety(t *testing.T) {
	fixture, err := LoadGovernanceReplayFixture(governanceFixtureAbsPath(t, "governance-master-win11-cpu.contract-prepared.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	result := EvaluateWindowsFluidApplyGovernance(ApplyGovernanceEvaluationInput{
		AdmissionDecision: fixture.AdmissionDecision,
		Bundle:            fixture.Bundle,
		RequestedAction:   fixture.RequestedAction,
		EvaluationTime:    time.Date(2026, 5, 7, 14, 45, 0, 0, time.UTC),
	})
	mode := result.PolicyAttestation.Signature.Mode
	if mode != "unsigned-dev" && mode != "future-signable" {
		t.Fatalf("unexpected signature mode %s", mode)
	}
	if result.PolicyAttestation.Signature.Value != "" {
		t.Fatal("signature value must stay empty in this phase")
	}
}

func TestGovernanceCLIOutputDeterministic(t *testing.T) {
	fixturePath := governanceFixtureAbsPath(t, "governance-master-win11-cpu.contract-prepared.json")
	cmd := exec.Command("go", "run", "./cmd/karl-fluid-governance", "-fixture", fixturePath, "-evaluation-time", "2026-05-07T14:45:00Z")
	cmd.Dir = governanceRepoRoot(t)
	out1, err := cmd.Output()
	if err != nil {
		t.Fatalf("first governance cli run failed: %v", err)
	}
	cmd = exec.Command("go", "run", "./cmd/karl-fluid-governance", "-fixture", fixturePath, "-evaluation-time", "2026-05-07T14:45:00Z")
	cmd.Dir = governanceRepoRoot(t)
	out2, err := cmd.Output()
	if err != nil {
		t.Fatalf("second governance cli run failed: %v", err)
	}
	if string(out1) != string(out2) {
		t.Fatal("governance cli output must be deterministic with fixed evaluation-time")
	}
}

func governanceFixtureAbsPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(governanceRepoRoot(t), "examples", "windows-fluid-governance-fixtures", name)
}

func governanceRepoRoot(t *testing.T) string {
	t.Helper()
	return admissionRepoRoot(t)
}
