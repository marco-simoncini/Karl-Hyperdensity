package performanceproof_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/performanceproof"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestSprint6SchemaFilesExist(t *testing.T) {
	root := repoRoot(t)
	for _, name := range performanceproof.SchemaFilesRequiredSprint6() {
		if _, err := os.Stat(filepath.Join(root, "schemas", name)); err != nil {
			t.Fatalf("missing schema %s: %v", name, err)
		}
	}
}

func TestSprint6ExamplesValidate(t *testing.T) {
	if err := performanceproof.ValidateSprint6Examples(repoRoot(t)); err != nil {
		t.Fatalf("ValidateSprint6Examples: %v", err)
	}
}

func TestUniversalPerformanceImprovementRejected(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": performanceproof.MilestoneUniversalSloGuard,
		"universalPerformanceImprovementClaimed": true,
		"guaranteedSavingsClaimed": false, "autoApplyAllowed": false,
		"productionAutonomousApplyAllowed": false,
		"safetyInvariants": map[string]interface{}{
			"universalPerformanceImprovementClaimed": true,
			"guaranteedSavingsClaimed": false, "autoApplyAllowed": false,
			"productionAutonomousApplyAllowed": false,
			"dashboardPerformanceSourceOfTruth": false, "fluidvirtPerformanceClaimAuthority": false,
		},
		"sourceOfTruth": map[string]interface{}{"performanceMeasurements": "FluidVirt"},
	}
	if err := performanceproof.ValidateUniversalSloGuardSurface(doc); err == nil {
		t.Fatal("expected universalPerformanceImprovementClaimed rejection")
	}
}

func TestCertifiedUpliftIoBoundRejected(t *testing.T) {
	doc := map[string]interface{}{
		"certified": true, "confidence": 0.9,
		"baselineSampleRef": "b", "postMutationSampleRef": "p",
		"excludedFromUniversalClaim": true,
		"claimBoundary": []interface{}{"scoped only"},
		"evidenceRefs": []interface{}{"e"},
	}
	if err := performanceproof.ValidateCertifiedPerformanceUplift(doc, "io_bound", true, true, true); err == nil {
		t.Fatal("expected io_bound certified uplift rejection")
	}
}

func TestCertifiedUpliftLowConfidenceRejected(t *testing.T) {
	doc := map[string]interface{}{
		"certified": true, "confidence": 0.5,
		"baselineSampleRef": "b", "postMutationSampleRef": "p",
		"excludedFromUniversalClaim": true,
		"claimBoundary": []interface{}{"scoped only"},
		"evidenceRefs": []interface{}{"e"},
	}
	if err := performanceproof.ValidateCertifiedPerformanceUplift(doc, "cpu_bound", true, true, true); err == nil {
		t.Fatal("expected low confidence rejection")
	}
}

func TestSloGuardPassWithRegressionRejected(t *testing.T) {
	doc := map[string]interface{}{
		"sloGuardStatus": "passed", "regressionDetected": true, "rollbackRequired": false,
		"claimBoundary": []interface{}{"x"}, "evidenceRefs": []interface{}{"y"},
	}
	if err := performanceproof.ValidateSloGuardEvaluation(doc); err == nil {
		t.Fatal("expected regression rejection")
	}
}

func TestNoRegressionWithRollbackRejected(t *testing.T) {
	doc := map[string]interface{}{
		"noRegressionCertified": true, "rollbackRequired": true, "regressionBlocked": false,
		"claimBoundary": []interface{}{"x"},
	}
	if err := performanceproof.ValidateNoRegressionResult(doc); err == nil {
		t.Fatal("expected rollbackRequired rejection")
	}
}

func TestNeutralNoClaimNotCertified(t *testing.T) {
	doc := map[string]interface{}{
		"proofStatus": "neutral_no_claim", "certifiedUplift": true,
		"claimBoundary": []interface{}{"x"},
	}
	if err := performanceproof.ValidatePerformanceProofClassification(doc); err == nil {
		t.Fatal("expected neutral certifiedUplift rejection")
	}
}

func TestInsufficientEvidenceNotCertified(t *testing.T) {
	doc := map[string]interface{}{
		"proofStatus": "insufficient_evidence", "certifiedUplift": true,
		"claimBoundary": []interface{}{"x"},
	}
	if err := performanceproof.ValidatePerformanceProofClassification(doc); err == nil {
		t.Fatal("expected insufficient_evidence certifiedUplift rejection")
	}
}

func TestAutoApplyRejected(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": performanceproof.MilestoneUniversalSloGuard,
		"universalPerformanceImprovementClaimed": false,
		"guaranteedSavingsClaimed": false, "autoApplyAllowed": true,
		"productionAutonomousApplyAllowed": false,
		"safetyInvariants": map[string]interface{}{
			"universalPerformanceImprovementClaimed": false,
			"guaranteedSavingsClaimed": false, "autoApplyAllowed": true,
			"productionAutonomousApplyAllowed": false,
			"dashboardPerformanceSourceOfTruth": false, "fluidvirtPerformanceClaimAuthority": false,
		},
		"sourceOfTruth": map[string]interface{}{"performanceMeasurements": "FluidVirt"},
	}
	if err := performanceproof.ValidateUniversalSloGuardSurface(doc); err == nil {
		t.Fatal("expected autoApply rejection")
	}
}
