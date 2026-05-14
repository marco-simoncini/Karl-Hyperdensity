package leaseslate_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/leaseslate"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestSprint3SchemaFilesExist(t *testing.T) {
	root := repoRoot(t)
	for _, name := range leaseslate.SchemaFilesRequiredSprint3() {
		if _, err := os.Stat(filepath.Join(root, "schemas", name)); err != nil {
			t.Fatalf("missing schema %s: %v", name, err)
		}
	}
}

func TestSprint3ExamplesValidate(t *testing.T) {
	if err := leaseslate.ValidateSprint3Examples(repoRoot(t)); err != nil {
		t.Fatalf("ValidateSprint3Examples: %v", err)
	}
}

func TestLeaseCandidateAutoApplyRejected(t *testing.T) {
	doc := map[string]interface{}{
		"autoApplyAllowed": true, "productionMutationAllowed": false,
		"claimBoundary": []interface{}{"x"}, "evidenceRefs": []interface{}{"y"},
	}
	if err := leaseslate.ValidateResourceLeaseCandidate(doc); err == nil {
		t.Fatal("expected autoApply rejection")
	}
}

func TestActionSlateNoFullNxNRejected(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": leaseslate.MilestoneResourceLeaseActionSlateReadiness,
		"noFullNxNPairing": false,
		"evaluatedPairCount": 2, "maxEvaluatedPairs": 4, "fullPairSpace": 6,
		"avoidedPairCount": 4,
		"safetyInvariants": map[string]interface{}{
			"autoApplyAllowed": false, "productionMutationAllowed": false,
			"guaranteedSavingsClaimed": false, "universalPerformanceImprovementClaimed": false,
			"windowsTotalRamHotplugClaimed": false, "logicalVcpuHotplugClaimed": false,
			"dashboardAppliesRuntimeChanges": false,
		},
		"sourceOfTruth": map[string]interface{}{"dryRunReadinessEvidence": "FluidVirt"},
		"actionEntries": []interface{}{},
		"protectedExcludedCount": 1,
	}
	if err := leaseslate.ValidateActionSlateReadiness(doc); err == nil {
		t.Fatal("expected noFullNxNPairing rejection")
	}
}

func TestBlockedActionNotOperatorReady(t *testing.T) {
	doc := map[string]interface{}{
		"autoApplyAllowed": false, "productionMutationAllowed": false,
		"dryRunStatus": "blocked", "rollbackStatus": "unknown",
		"claimBoundary": []interface{}{"x"}, "evidenceRefs": []interface{}{"y"},
		"readyForApply": "operator_controlled_ready",
		"blockers":      []interface{}{"cpu_hotplug_headroom_missing"},
	}
	if err := leaseslate.ValidateActionSlateEntry(doc); err == nil {
		t.Fatal("expected blocked action rejection")
	}
}

func TestDryRunReadinessFluidVirtOnly(t *testing.T) {
	doc := map[string]interface{}{
		"source": "Dashboard", "actuator": "Dashboard",
		"mutationExecuted": false, "dryRunStatus": "valid",
	}
	if err := leaseslate.ValidateActionDryrunReadiness(doc); err == nil {
		t.Fatal("expected non-FluidVirt rejection")
	}
}

func TestEvaluatedPairsExceedsMax(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": leaseslate.MilestoneResourceLeaseActionSlateReadiness,
		"noFullNxNPairing": true, "evaluatedPairCount": 10, "maxEvaluatedPairs": 4,
		"fullPairSpace": 20, "avoidedPairCount": 10, "protectedExcludedCount": 1,
		"safetyInvariants": map[string]interface{}{
			"autoApplyAllowed": false, "productionMutationAllowed": false,
			"guaranteedSavingsClaimed": false, "universalPerformanceImprovementClaimed": false,
			"windowsTotalRamHotplugClaimed": false, "logicalVcpuHotplugClaimed": false,
			"dashboardAppliesRuntimeChanges": false,
		},
		"sourceOfTruth": map[string]interface{}{"dryRunReadinessEvidence": "FluidVirt"},
		"actionEntries": []interface{}{},
	}
	if err := leaseslate.ValidateActionSlateReadiness(doc); err == nil {
		t.Fatal("expected maxEvaluatedPairs rejection")
	}
}
