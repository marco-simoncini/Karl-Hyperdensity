package kernelboundary_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/kernelboundary"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestSchemaFilesExist(t *testing.T) {
	root := repoRoot(t)
	for _, name := range kernelboundary.SchemaFilesRequired() {
		if _, err := os.Stat(filepath.Join(root, "schemas", name)); err != nil {
			t.Fatalf("missing schema %s: %v", name, err)
		}
	}
}

func TestSprint1ExamplesValidate(t *testing.T) {
	if err := kernelboundary.ValidateSprint1Examples(repoRoot(t)); err != nil {
		t.Fatalf("ValidateSprint1Examples: %v", err)
	}
}

func TestForbiddenPhraseRejected(t *testing.T) {
	doc := map[string]interface{}{
		"claimPolicyId":                          kernelboundary.ClaimPolicyV2ID,
		"guaranteedSavingsAllowed":               false,
		"universalPerformanceImprovementAllowed": false,
		"logicalVcpuHotplugClaimAllowed":         false,
		"windowsTotalRamHotplugClaimAllowed":     false,
		"ramAboveOriginalClaimAllowed":           false,
		"productionAutonomousApplyAllowed":       false,
		"syntheticFleetProductionClaimAllowed":   false,
		"allowedPhrases":                         []interface{}{"guaranteed savings active"},
	}
	if err := kernelboundary.ValidateClaimPolicyV2(doc); err == nil {
		t.Fatal("expected forbidden phrase rejection")
	}
}

func TestProductionAutonomousApplyMustBeFalse(t *testing.T) {
	doc := map[string]interface{}{
		"claimPolicyId":                    kernelboundary.ClaimPolicyV2ID,
		"guaranteedSavingsAllowed":           false,
		"universalPerformanceImprovementAllowed": false,
		"logicalVcpuHotplugClaimAllowed":   false,
		"windowsTotalRamHotplugClaimAllowed": false,
		"ramAboveOriginalClaimAllowed":     false,
		"productionAutonomousApplyAllowed": true,
		"syntheticFleetProductionClaimAllowed": false,
	}
	if err := kernelboundary.ValidateClaimPolicyV2(doc); err == nil {
		t.Fatal("expected productionAutonomousApplyAllowed=false enforcement")
	}
}

func TestFluidVirtIsOnlyRuntimeActuator(t *testing.T) {
	root := repoRoot(t)
	if err := kernelboundary.ValidateExampleFile(root, "production-kernel-boundary-reference.json"); err != nil {
		t.Fatalf("boundary reference: %v", err)
	}
}

func TestDashboardProjectionOnly(t *testing.T) {
	root := repoRoot(t)
	b, err := os.ReadFile(filepath.Join(root, "examples", "production-kernel-boundary-reference.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"projectionOnly": true`) {
		t.Fatal("Dashboard must be projectionOnly in boundary reference")
	}
	if strings.Contains(strings.ToLower(string(b)), `"dashboard is source of truth"`) {
		t.Fatal("forbidden Dashboard source-of-truth claim in reference positive fields")
	}
}
