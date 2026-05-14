package applygate_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/applygate"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestSprint4SchemaFilesExist(t *testing.T) {
	root := repoRoot(t)
	for _, name := range applygate.SchemaFilesRequiredSprint4() {
		if _, err := os.Stat(filepath.Join(root, "schemas", name)); err != nil {
			t.Fatalf("missing schema %s: %v", name, err)
		}
	}
}

func TestSprint4ExamplesValidate(t *testing.T) {
	if err := applygate.ValidateSprint4Examples(repoRoot(t)); err != nil {
		t.Fatalf("ValidateSprint4Examples: %v", err)
	}
}

func TestAutoApplyRejected(t *testing.T) {
	doc := map[string]interface{}{
		"actuator": "FluidVirt", "productionScope": false,
		"autoApplyAllowed": true, "productionAutonomousApplyAllowed": false,
		"rawRuntimeControlsExposed": false,
		"approvalRef": "a", "dryRunRef": "b", "rollbackReadinessRef": "c",
		"sloPrecheckRef": "d", "riskAssessmentRef": "e",
		"actionId": "x", "leaseCandidateId": "y", "donorShellId": "d", "receiverShellId": "r",
		"claimBoundary": []interface{}{"ok"}, "evidenceRefs": []interface{}{"e"},
	}
	if err := applygate.ValidateApplyRequest(doc); err == nil {
		t.Fatal("expected autoApply rejection")
	}
}

func TestProductionScopeRejected(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": applygate.MilestoneOperatorControlledApplyGate,
		"operatorControlledApplyAllowed": true,
		"autoApplyAllowed": false, "productionAutonomousApplyAllowed": false,
		"rawRuntimeControlsExposed": false, "productionScope": true,
		"mutationScope": "technical_preview_operator_controlled",
		"safetyInvariants": map[string]interface{}{
			"autoApplyAllowed": false, "productionAutonomousApplyAllowed": false,
			"rawRuntimeControlsExposed": false, "dashboardAppliesRuntimeChanges": false,
			"inventoryRuntimeApply": false,
		},
		"sourceOfTruth": map[string]interface{}{"runtimeMutation": "FluidVirt"},
		"scenarioSummaries": map[string]interface{}{
			"successfulOperatorControlled": map[string]interface{}{},
			"blockedApply":                 map[string]interface{}{"mutationExecuted": false},
			"rollbackRequired":             map[string]interface{}{},
			"windowsEvidenceGated":         map[string]interface{}{"guestDeltaVerified": false},
		},
	}
	if err := applygate.ValidateOperatorApplyGateSurface(doc); err == nil {
		t.Fatal("expected productionScope rejection")
	}
}

func TestMutationWithoutApprovalRejected(t *testing.T) {
	doc := map[string]interface{}{
		"actuator": "FluidVirt", "mutationExecuted": true, "productionScope": false,
		"rawRuntimeControlsExposed": false, "autoApplyAllowed": false,
		"productionAutonomousApplyAllowed": false,
		"claimBoundary": []interface{}{"x"}, "evidenceRefs": []interface{}{"y"},
	}
	if err := applygate.ValidateFluidvirtInvocationRecord(doc, false); err == nil {
		t.Fatal("expected missing approval rejection")
	}
}

func TestPostVerifyPassWithoutPlanMatchRejected(t *testing.T) {
	doc := map[string]interface{}{
		"verifyStatus": "passed", "runtimeDeltaVerified": true,
		"sloGuardStatus": "passed", "donorHealthStatus": "preserved",
		"rollbackReady": true, "mutationRolledBack": false, "guestDeltaVerified": true,
		"claimBoundary": []interface{}{"x"}, "evidenceRefs": []interface{}{"y"},
	}
	if err := applygate.ValidatePostVerifyResult(doc, false); err == nil {
		t.Fatal("expected mutationMatchesPlan rejection")
	}
}

func TestDisruptiveMutationRejected(t *testing.T) {
	doc := map[string]interface{}{
		"identityPreserved": true, "rebootUsed": true, "recreateUsed": false,
		"rolloutUsed": false, "migrationUsed": false,
		"claimBoundary": []interface{}{"x"}, "evidenceRefs": []interface{}{"y"},
	}
	if err := applygate.ValidateRuntimeMutationObservation(doc); err == nil {
		t.Fatal("expected rebootUsed rejection")
	}
}

func TestNonFluidVirtActuatorRejected(t *testing.T) {
	doc := map[string]interface{}{
		"actuator": "Dashboard", "productionScope": false,
		"autoApplyAllowed": false, "productionAutonomousApplyAllowed": false,
		"rawRuntimeControlsExposed": false,
		"approvalRef": "a", "dryRunRef": "b", "rollbackReadinessRef": "c",
		"sloPrecheckRef": "d", "riskAssessmentRef": "e",
		"actionId": "x", "leaseCandidateId": "y", "donorShellId": "d", "receiverShellId": "r",
		"claimBoundary": []interface{}{"ok"}, "evidenceRefs": []interface{}{"e"},
	}
	if err := applygate.ValidateApplyRequest(doc); err == nil {
		t.Fatal("expected non-FluidVirt rejection")
	}
}

func TestMissingApprovalRefRejected(t *testing.T) {
	doc := map[string]interface{}{
		"actuator": "FluidVirt", "productionScope": false,
		"autoApplyAllowed": false, "productionAutonomousApplyAllowed": false,
		"rawRuntimeControlsExposed": false,
		"dryRunRef": "b", "rollbackReadinessRef": "c",
		"sloPrecheckRef": "d", "riskAssessmentRef": "e",
		"actionId": "x", "leaseCandidateId": "y", "donorShellId": "d", "receiverShellId": "r",
		"claimBoundary": []interface{}{"ok"}, "evidenceRefs": []interface{}{"e"},
	}
	if err := applygate.ValidateApplyRequest(doc); err == nil {
		t.Fatal("expected missing approvalRef rejection")
	}
}

func TestPostVerifyPassWithoutRollbackReadyRejected(t *testing.T) {
	doc := map[string]interface{}{
		"verifyStatus": "passed", "runtimeDeltaVerified": true,
		"sloGuardStatus": "passed", "donorHealthStatus": "preserved",
		"rollbackReady": false, "mutationRolledBack": false, "guestDeltaVerified": true,
		"claimBoundary": []interface{}{"x"}, "evidenceRefs": []interface{}{"y"},
	}
	if err := applygate.ValidatePostVerifyResult(doc, true); err == nil {
		t.Fatal("expected rollbackReady rejection")
	}
}
