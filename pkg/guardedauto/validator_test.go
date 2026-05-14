package guardedauto_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/guardedauto"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestSprint7SchemaFilesExist(t *testing.T) {
	root := repoRoot(t)
	for _, name := range guardedauto.SchemaFilesRequiredSprint7() {
		if _, err := os.Stat(filepath.Join(root, "schemas", name)); err != nil {
			t.Fatalf("missing schema %s: %v", name, err)
		}
	}
}

func TestSprint7ExamplesValidate(t *testing.T) {
	if err := guardedauto.ValidateSprint7Examples(repoRoot(t)); err != nil {
		t.Fatalf("ValidateSprint7Examples: %v", err)
	}
}

func baseSurface() map[string]interface{} {
	return map[string]interface{}{
		"milestone": guardedauto.MilestoneGuardedAutoPolicyEngine,
		"autoApplyExecutionEnabled": false, "productionAutonomousApplyAllowed": false, "productionScope": false,
		"policyMode": "guarded_auto_candidate",
		"safetyInvariants": map[string]interface{}{
			"autoApplyExecutionEnabled": false, "productionAutonomousApplyAllowed": false, "productionScope": false,
			"dashboardPolicySourceOfTruth": false, "fluidvirtPolicyAuthority": false, "inventoryRuntimeApply": false,
		},
		"sourceOfTruth": map[string]interface{}{"runtimeConstraintEvidence": "FluidVirt"},
	}
}

func TestAutoApplyExecutionEnabledRejected(t *testing.T) {
	doc := baseSurface()
	doc["autoApplyExecutionEnabled"] = true
	if err := guardedauto.ValidateGuardedAutoPolicyEngineSurface(doc); err == nil {
		t.Fatal("expected autoApplyExecutionEnabled rejection")
	}
}

func TestProductionAutonomousApplyRejected(t *testing.T) {
	doc := baseSurface()
	doc["productionAutonomousApplyAllowed"] = true
	if err := guardedauto.ValidateGuardedAutoPolicyEngineSurface(doc); err == nil {
		t.Fatal("expected productionAutonomousApplyAllowed rejection")
	}
}

func TestProductionScopeRejected(t *testing.T) {
	doc := baseSurface()
	doc["productionScope"] = true
	if err := guardedauto.ValidateGuardedAutoPolicyEngineSurface(doc); err == nil {
		t.Fatal("expected productionScope rejection")
	}
}

func TestProductionAutoWithPolicyModeRejected(t *testing.T) {
	doc := baseSurface()
	doc["policyMode"] = "production_auto_with_policy"
	if err := guardedauto.ValidateGuardedAutoPolicyEngineSurface(doc); err == nil {
		t.Fatal("expected production_auto_with_policy rejection")
	}
}

func baseCandidate() map[string]interface{} {
	return map[string]interface{}{
		"candidateState": "candidate_only", "rollbackReady": true, "sloGuardPassed": true,
		"noRegressionCertified": true, "donorHealthPreserved": true, "receiverHealthPreserved": true,
		"killSwitchClear": true, "circuitBreakerClosed": true, "rateLimitAvailable": true, "cooldownExpired": true,
		"ledgerRecordRef": "l", "performanceProofRef": "p",
		"autoApplyExecutionEnabled": false, "productionAutonomousApplyAllowed": false, "productionScope": false,
		"evidenceRefs": []interface{}{"e"}, "claimBoundary": []interface{}{"scoped"},
	}
}

func TestCandidateAutoApplyRejected(t *testing.T) {
	c := baseCandidate()
	c["autoApplyExecutionEnabled"] = true
	if err := guardedauto.ValidateGuardedAutoCandidate(c); err == nil {
		t.Fatal("expected candidate autoApplyExecutionEnabled rejection")
	}
}

func TestCandidateMissingRollbackRejected(t *testing.T) {
	c := baseCandidate()
	c["rollbackReady"] = false
	if err := guardedauto.ValidateGuardedAutoCandidate(c); err == nil {
		t.Fatal("expected rollbackReady rejection")
	}
}

func TestCandidateMissingSloGuardRejected(t *testing.T) {
	c := baseCandidate()
	c["sloGuardPassed"] = false
	if err := guardedauto.ValidateGuardedAutoCandidate(c); err == nil {
		t.Fatal("expected sloGuardPassed rejection")
	}
}

func TestCandidateKillSwitchRejected(t *testing.T) {
	c := baseCandidate()
	c["killSwitchClear"] = false
	if err := guardedauto.ValidateGuardedAutoCandidate(c); err == nil {
		t.Fatal("expected killSwitchClear rejection")
	}
}

func baseDecision() map[string]interface{} {
	return map[string]interface{}{
		"decision": "guarded_auto_candidate", "eligibleAsGuardedAutoCandidate": true, "eligibleAsProductionCanary": false,
		"autoApplyExecutionEnabled": false, "productionAutonomousApplyAllowed": false, "productionScope": false,
		"dryRunStatus": "valid", "rollbackStatus": "ready", "sloGuardStatus": "passed",
		"noRegressionStatus": "certified", "donorHealthStatus": "preserved", "receiverHealthStatus": "preserved",
		"ledgerStatus": "record_present", "blastRadiusStatus": "available", "killSwitchStatus": "clear",
		"circuitBreakerStatus": "closed", "rateLimitStatus": "available", "cooldownStatus": "expired",
		"windowsEvidenceGateStatus": "not_applicable", "syntheticProductionStatus": "not_synthetic",
		"claimBoundary": []interface{}{"x"},
	}
}

func TestEligibleAsProductionCanaryRejected(t *testing.T) {
	d := baseDecision()
	d["eligibleAsProductionCanary"] = true
	if err := guardedauto.ValidateEligibilityDecision(d); err == nil {
		t.Fatal("expected eligibleAsProductionCanary rejection")
	}
}

func TestWindowsGatedCandidateRejected(t *testing.T) {
	d := baseDecision()
	d["windowsEvidenceGateStatus"] = "gated"
	if err := guardedauto.ValidateEligibilityDecision(d); err == nil {
		t.Fatal("expected Windows gated candidate rejection")
	}
}

func TestSyntheticProductionCandidateRejected(t *testing.T) {
	d := baseDecision()
	d["syntheticProductionStatus"] = "synthetic_shadow"
	if err := guardedauto.ValidateEligibilityDecision(d); err == nil {
		t.Fatal("expected synthetic production rejection")
	}
}

func TestDeniedRegressionNotCandidate(t *testing.T) {
	d := baseDecision()
	d["decision"] = "denied_regression"
	d["eligibleAsGuardedAutoCandidate"] = true
	if err := guardedauto.ValidateEligibilityDecision(d); err == nil {
		t.Fatal("expected denied_regression rejection")
	}
}

func TestPolicyWindowsGatedAutoRejected(t *testing.T) {
	p := map[string]interface{}{
		"autoApplyExecutionEnabled": false, "productionAutonomousApplyAllowed": false, "productionScope": false,
		"allowWindowsEvidenceGatedAuto": true, "allowSyntheticProductionProof": false,
		"claimBoundary": []interface{}{"x"},
	}
	if err := guardedauto.ValidateGuardedAutoPolicy(p); err == nil {
		t.Fatal("expected allowWindowsEvidenceGatedAuto rejection")
	}
}

func TestPolicySyntheticProofRejected(t *testing.T) {
	p := map[string]interface{}{
		"autoApplyExecutionEnabled": false, "productionAutonomousApplyAllowed": false, "productionScope": false,
		"allowWindowsEvidenceGatedAuto": false, "allowSyntheticProductionProof": true,
		"claimBoundary": []interface{}{"x"},
	}
	if err := guardedauto.ValidateGuardedAutoPolicy(p); err == nil {
		t.Fatal("expected allowSyntheticProductionProof rejection")
	}
}
