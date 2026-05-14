package productioncanary_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/productioncanary"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestSprint9SchemaFilesExist(t *testing.T) {
	root := repoRoot(t)
	for _, name := range productioncanary.SchemaFilesRequiredSprint9() {
		if _, err := os.Stat(filepath.Join(root, "schemas", name)); err != nil {
			t.Fatalf("missing schema %s: %v", name, err)
		}
	}
}

func TestSprint9ExamplesValidate(t *testing.T) {
	if err := productioncanary.ValidateSprint9Examples(repoRoot(t)); err != nil {
		t.Fatalf("ValidateSprint9Examples: %v", err)
	}
}

func TestGeneralProductionAutoRejected(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": productioncanary.MilestoneProductionCanaryAutoApply,
		"productionCanaryAutoApplyAllowed": true, "productionCanaryScope": true,
		"generalProductionAutoAllowed": true, "productionAutoWithPolicy": false,
		"productionScope": true, "rawRuntimeControlsExposed": false,
		"safetyInvariants": map[string]interface{}{"generalProductionAutoAllowed": true},
	}
	if err := productioncanary.ValidateProductionCanaryAutoApplySurface(doc); err == nil {
		t.Fatal("expected generalProductionAutoAllowed rejection")
	}
}

func TestProductionScopeWithoutCanaryRejected(t *testing.T) {
	s := map[string]interface{}{
		"selected": true, "candidateState": "production_canary_ready",
		"productionScope": true, "productionCanaryScope": false,
		"allowlisted": true, "broadProductionScope": false,
		"claimBoundary": []interface{}{"x"},
	}
	if err := productioncanary.ValidateCanarySelection(s); err == nil {
		t.Fatal("expected productionCanaryScope rejection")
	}
}

func TestBroadProductionSelectedRejected(t *testing.T) {
	s := map[string]interface{}{
		"selected": true, "candidateState": "production_canary_ready",
		"productionScope": true, "productionCanaryScope": true,
		"allowlisted": true, "broadProductionScope": true,
		"claimBoundary": []interface{}{"x"},
	}
	if err := productioncanary.ValidateCanarySelection(s); err == nil {
		t.Fatal("expected broadProductionScope rejection")
	}
}

func TestMutationWithoutAllowlistRejected(t *testing.T) {
	inv := map[string]interface{}{
		"actuator": "FluidVirt", "invocationMode": "guarded_production_canary",
		"mutationExecuted": true, "productionScope": true, "productionCanaryScope": true,
		"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
		"rawRuntimeControlsExposed": false,
		"claimBoundary": []interface{}{"x"},
	}
	if err := productioncanary.ValidateCanaryFluidvirtInvocation(inv, true, true, false, true, true, true); err == nil {
		t.Fatal("expected allowlist rejection")
	}
}

func TestPromotedToGeneralProductionRejected(t *testing.T) {
	c := map[string]interface{}{
		"promotedToGeneralProduction": true, "nextAllowedState": "production_canary_only",
		"claimBoundary": []interface{}{"x"},
	}
	if err := productioncanary.ValidateCanaryCloseout(c); err == nil {
		t.Fatal("expected promotedToGeneralProduction rejection")
	}
}

func TestImmutableAuditNotWrittenRejected(t *testing.T) {
	a := map[string]interface{}{
		"immutableAuditRequired": true, "immutableAuditWritten": false,
		"events": []interface{}{},
		"claimBoundary": []interface{}{"x"},
	}
	if err := productioncanary.ValidateImmutableAuditTrail(a); err == nil {
		t.Fatal("expected immutableAuditWritten rejection")
	}
}
