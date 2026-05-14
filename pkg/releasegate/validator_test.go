package releasegate

import (
	"testing"
)

func gateDoc(overrides map[string]interface{}) map[string]interface{} {
	doc := map[string]interface{}{
		"milestone":                               MilestoneDashboardEnterpriseCleanupGA,
		"releaseDecision":                         "canary_only",
		"gaReady":                                 false,
		"canaryReady":                             true,
		"generalProductionAutoAllowed":            false,
		"productionAutoWithPolicy":                false,
		"guaranteedSavingsClaimed":                false,
		"universalPerformanceImprovementClaimed":  false,
		"referencePayloadCountedAsProduction":     false,
		"syntheticProofCountedAsProduction":       false,
		"dashboardRuntimeControlsExposed":         false,
		"dashboardExecutor":                       false,
		"fluidvirtPolicyAuthority":                false,
		"inventoryRuntimeExecutor":                false,
		"claimBoundaries":                         []interface{}{"Sprint 10 allowed claim only"},
		"checks": []interface{}{
			map[string]interface{}{"checkId": "general-production-auto-boundary", "status": "failed", "requiredForGA": true},
		},
	}
	for k, v := range overrides {
		doc[k] = v
	}
	return doc
}

func surfaceDoc(overrides map[string]interface{}) map[string]interface{} {
	doc := map[string]interface{}{
		"surfaceId": "test-surface", "surfaceName": "Test", "route": "/hyperdensity/kernel",
		"releaseTrack": "ga", "visibleInExecutive": true, "referenceOnly": false,
		"syntheticOrShadow": false, "canBeUsedForProductionProof": true,
		"claimBoundary": "test boundary",
	}
	for k, v := range overrides {
		doc[k] = v
	}
	return doc
}

func TestValidateGAReleaseGateRejectsUnsafeFlags(t *testing.T) {
	keys := []string{
		"generalProductionAutoAllowed", "productionAutoWithPolicy", "guaranteedSavingsClaimed",
		"universalPerformanceImprovementClaimed", "referencePayloadCountedAsProduction",
		"syntheticProofCountedAsProduction", "dashboardRuntimeControlsExposed", "dashboardExecutor",
		"fluidvirtPolicyAuthority", "inventoryRuntimeExecutor",
	}
	for _, key := range keys {
		doc := gateDoc(map[string]interface{}{key: true})
		if err := ValidateGAReleaseGate(doc); err == nil {
			t.Fatalf("expected rejection for %s=true", key)
		}
	}
}

func TestValidateGAReleaseGateRejectsGaPassWithFailedCheck(t *testing.T) {
	doc := gateDoc(map[string]interface{}{
		"releaseDecision": "ga_pass",
		"gaReady":         true,
		"checks": []interface{}{
			map[string]interface{}{"checkId": "windows-claim-boundary", "status": "failed", "requiredForGA": true},
		},
	})
	if err := ValidateGAReleaseGate(doc); err == nil {
		t.Fatal("ga_pass with failed required check must be rejected")
	}
}

func TestValidateSurfaceClassificationRejectsReferenceProductionProof(t *testing.T) {
	doc := surfaceDoc(map[string]interface{}{
		"referenceOnly": true, "canBeUsedForProductionProof": true, "releaseTrack": "reference_only",
		"route": "/lab/reference/sample", "visibleInExecutive": false,
	})
	if err := ValidateSurfaceClassification(doc); err == nil {
		t.Fatal("referenceOnly production proof must be rejected")
	}
}

func TestValidateSurfaceClassificationRejectsSyntheticProductionProof(t *testing.T) {
	doc := surfaceDoc(map[string]interface{}{
		"syntheticOrShadow": true, "canBeUsedForProductionProof": true, "releaseTrack": "synthetic_shadow",
		"route": "/lab/evidence/synthetic-shadow", "visibleInExecutive": false,
	})
	if err := ValidateSurfaceClassification(doc); err == nil {
		t.Fatal("synthetic production proof must be rejected")
	}
}

func TestValidateSurfaceClassificationRejectsLabInExecutive(t *testing.T) {
	doc := surfaceDoc(map[string]interface{}{
		"releaseTrack": "lab", "visibleInExecutive": true, "route": "/lab/hyperdensity",
		"canBeUsedForProductionProof": false,
	})
	if err := ValidateSurfaceClassification(doc); err == nil {
		t.Fatal("lab surface in executive must be rejected")
	}
}

func TestValidateSurfaceClassificationRejectsDemoRouteInExecutive(t *testing.T) {
	doc := surfaceDoc(map[string]interface{}{
		"releaseTrack": "preview", "visibleInExecutive": true, "route": "/demo/hyperdensity/mock-test",
		"canBeUsedForProductionProof": false,
	})
	if err := ValidateSurfaceClassification(doc); err == nil {
		t.Fatal("demo route in executive must be rejected")
	}
}

func TestValidateSurfaceClassificationRejectsCanaryAsGeneralAuto(t *testing.T) {
	doc := surfaceDoc(map[string]interface{}{
		"releaseTrack": "production_canary", "claimBoundary": "general production auto enabled",
		"route": "/hyperdensity/production-canary", "canBeUsedForProductionProof": true,
	})
	if err := ValidateSurfaceClassification(doc); err == nil {
		t.Fatal("canary labeled general production auto must be rejected")
	}
}

func TestValidateReleaseGateDecisionRejectsProductionAutoWithPolicy(t *testing.T) {
	doc := map[string]interface{}{
		"decidedByComponent": "Karl-Hyperdensity",
		"nextAllowedState":   "production_auto_with_policy",
		"claimBoundary":      "test",
	}
	if err := ValidateReleaseGateDecision(doc); err == nil {
		t.Fatal("production_auto_with_policy next state must be rejected")
	}
}

func TestValidateClaimGateRejectsForbiddenClaimsPresent(t *testing.T) {
	doc := map[string]interface{}{
		"forbiddenClaimsPresent": true,
		"claimBoundary":          "test",
		"evidenceRefs":           []interface{}{"ref-1"},
	}
	if err := ValidateClaimGate(doc); err == nil {
		t.Fatal("forbiddenClaimsPresent must be rejected")
	}
}

func TestValidateEvidenceGateRejectsReferenceInProduction(t *testing.T) {
	doc := map[string]interface{}{
		"claimBoundary":          "test",
		"evidenceRefs":           []interface{}{"ref-1"},
		"productionEvidenceRefs": []interface{}{"hyperdensity-reference-configmap-sample-v1"},
		"referenceOnlyRefs":      []interface{}{"hyperdensity-reference-configmap-sample-v1"},
		"syntheticShadowRefs":    []interface{}{},
	}
	if err := ValidateEvidenceGate(doc); err == nil {
		t.Fatal("reference in production evidence must be rejected")
	}
}

func TestValidateReferencePayloadGateRequiresFalseCountedAsProduction(t *testing.T) {
	doc := map[string]interface{}{"countedAsProduction": true, "claimBoundary": "test"}
	if err := ValidateReferencePayloadGate(doc); err == nil {
		t.Fatal("countedAsProduction=true must be rejected")
	}
}

func TestValidateSyntheticProofGateRequiresFalseCountedAsProduction(t *testing.T) {
	doc := map[string]interface{}{"countedAsProduction": true, "claimBoundary": "test"}
	if err := ValidateSyntheticProofGate(doc); err == nil {
		t.Fatal("countedAsProduction=true must be rejected")
	}
}

func TestValidateCountConsistencyGateMismatchCount(t *testing.T) {
	doc := map[string]interface{}{
		"claimBoundary": "test",
		"evidenceRefs":  []interface{}{"ref-1"},
		"mismatchCount": 2,
		"mismatches": []interface{}{
			map[string]interface{}{"surfaceId": "hyperdensityGuardedAutoApplySandbox"},
		},
	}
	if err := ValidateCountConsistencyGate(doc); err == nil {
		t.Fatal("mismatchCount mismatch must be rejected")
	}
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/ga-release-gate-reference.json"); err != nil {
		t.Fatalf("reference file invalid: %v", err)
	}
}

func TestScanForbiddenClaimsInAllowedClaim(t *testing.T) {
	if err := scanForbiddenClaims("guaranteed savings active"); err == nil {
		t.Fatal("forbidden claim must be detected")
	}
	if err := scanForbiddenClaims("guaranteed savings remains disabled"); err != nil {
		t.Fatalf("negated claim should pass: %v", err)
	}
}

func TestMissingClaimBoundaryRejected(t *testing.T) {
	doc := gateDoc(map[string]interface{}{"claimBoundaries": []interface{}{}})
	if err := ValidateGAReleaseGate(doc); err == nil {
		t.Fatal("empty claimBoundaries must be rejected")
	}
}
