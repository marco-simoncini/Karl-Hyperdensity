package autoapply_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/autoapply"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestSprint8SchemaFilesExist(t *testing.T) {
	root := repoRoot(t)
	for _, name := range autoapply.SchemaFilesRequiredSprint8() {
		if _, err := os.Stat(filepath.Join(root, "schemas", name)); err != nil {
			t.Fatalf("missing schema %s: %v", name, err)
		}
	}
}

func TestSprint8ExamplesValidate(t *testing.T) {
	if err := autoapply.ValidateSprint8Examples(repoRoot(t)); err != nil {
		t.Fatalf("ValidateSprint8Examples: %v", err)
	}
}

func TestProductionScopeRejected(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": autoapply.MilestoneGuardedAutoApplySandboxNonprod,
		"sandboxAutoApplyAllowed": true, "nonProdAutoApplyAllowed": true,
		"autoApplyExecutionEnabled": true, "productionAutonomousApplyAllowed": false,
		"productionScope": true, "productionMutationAllowed": false, "rawRuntimeControlsExposed": false,
		"autoExecutionScope": "sandbox",
	}
	if err := autoapply.ValidateGuardedAutoApplySandboxSurface(doc); err == nil {
		t.Fatal("expected productionScope rejection")
	}
}

func TestCandidateOnlySelectedRejected(t *testing.T) {
	s := map[string]interface{}{
		"selected": true, "candidateState": "candidate_only",
		"productionScope": false, "sandboxOrNonProd": false,
		"claimBoundary": []interface{}{"x"},
	}
	if err := autoapply.ValidateAutoApplySelection(s); err == nil {
		t.Fatal("expected candidate_only selection rejection")
	}
}

func TestMutationWithoutPreflightRejected(t *testing.T) {
	inv := map[string]interface{}{
		"actuator": "FluidVirt", "invocationMode": "guarded_auto_sandbox_or_nonprod",
		"mutationExecuted": true, "productionScope": false,
		"productionAutonomousApplyAllowed": false, "rawRuntimeControlsExposed": false,
		"claimBoundary": []interface{}{"x"},
	}
	if err := autoapply.ValidateAutoFluidvirtInvocation(inv, false, true, true, true); err == nil {
		t.Fatal("expected preflight rejection")
	}
}

func TestRebootUsedRejected(t *testing.T) {
	m := map[string]interface{}{
		"identityPreserved": true, "rebootUsed": true,
		"recreateUsed": false, "rolloutUsed": false, "migrationUsed": false,
		"claimBoundary": []interface{}{"x"},
	}
	if err := autoapply.ValidateAutoMutationObservation(m); err == nil {
		t.Fatal("expected rebootUsed rejection")
	}
}

func TestPostVerifyPassedWithoutSloRejected(t *testing.T) {
	pv := map[string]interface{}{
		"verifyStatus": "passed", "runtimeDeltaVerified": true,
		"sloGuardStatus": "failed", "donorHealthStatus": "preserved",
		"noRegressionStatus": "certified", "rollbackReady": true, "rollbackRequired": false,
		"claimBoundary": []interface{}{"x"},
	}
	if err := autoapply.ValidateAutoPostVerifyResult(pv); err == nil {
		t.Fatal("expected sloGuardStatus rejection")
	}
}

func TestRollbackFailedHealthRestoredRejected(t *testing.T) {
	re := map[string]interface{}{
		"executed": true, "rollbackPassed": false, "healthRestored": true,
		"claimBoundary": []interface{}{"x"},
	}
	if err := autoapply.ValidateAutoRollbackExecution(re); err == nil {
		t.Fatal("expected failed rollback healthRestored rejection")
	}
}
