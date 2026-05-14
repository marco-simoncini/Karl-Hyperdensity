package guaranteedsavings

import "testing"

func activationDoc(overrides map[string]interface{}) map[string]interface{} {
	doc := map[string]interface{}{
		"milestone":                               MilestoneGuaranteedEligibleSavingsActivation,
		"guaranteedEligibleSavingsAllowed":        true,
		"guaranteedEligibleSavingsClaimed":        true,
		"guaranteedSavingsScope":                  "scoped_eligible_records",
		"universalGuaranteedSavingsAllowed":       false,
		"universalGuaranteedSavingsClaimed":       false,
		"estimatedValueCountedAsGuaranteed":       false,
		"syntheticValueCountedAsGuaranteed":       false,
		"referencePayloadCountedAsGuaranteed":     false,
		"windowsEvidenceGatedValueGuaranteed":   false,
		"generalProductionAutoAllowed":            false,
		"productionAutoWithPolicy":                false,
		"universalPerformanceImprovementClaimed":  false,
		"guaranteePolicies": []interface{}{
			map[string]interface{}{
				"allowEstimatedOpportunity": false, "allowSyntheticShadow": false,
				"allowReferencePayload": false, "allowWindowsEvidenceGated": false,
				"allowPlaceholderPrice": false, "universalGuaranteeAllowed": false,
				"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
				"claimBoundary": "scoped only",
			},
		},
		"guaranteedRecords": []interface{}{},
		"claimBoundaries":   []interface{}{"scoped eligible records only; universal guaranteed savings remains disabled"},
	}
	for k, v := range overrides {
		doc[k] = v
	}
	return doc
}

func guaranteedRecord(overrides map[string]interface{}) map[string]interface{} {
	doc := map[string]interface{}{
		"guaranteeClaimStatus": "guaranteed_eligible",
		"netRealizedValue":     0.0042,
		"guaranteedValue":      0.00336,
		"confidence":           0.87,
		"durationSeconds":      300.0,
		"sloPreserved":         true,
		"donorHealthPreserved": true,
		"rollbackReady":        true,
		"rollbackFailed":       false,
		"evidenceRefs":         []interface{}{"ref-1"},
		"claimBoundary":        "test",
	}
	for k, v := range overrides {
		doc[k] = v
	}
	return doc
}

func TestValidateActivationSurfaceRejectsUniversalGuarantee(t *testing.T) {
	for _, key := range []string{"universalGuaranteedSavingsAllowed", "universalGuaranteedSavingsClaimed"} {
		doc := activationDoc(map[string]interface{}{key: true})
		if err := ValidateActivationSurface(doc); err == nil {
			t.Fatalf("expected rejection for %s=true", key)
		}
	}
}

func TestValidateActivationSurfaceRejectsEstimatedAsGuaranteed(t *testing.T) {
	doc := activationDoc(map[string]interface{}{"estimatedValueCountedAsGuaranteed": true})
	if err := ValidateActivationSurface(doc); err == nil {
		t.Fatal("estimated as guaranteed must be rejected")
	}
}

func TestValidateActivationSurfaceRejectsSyntheticAsGuaranteed(t *testing.T) {
	doc := activationDoc(map[string]interface{}{"syntheticValueCountedAsGuaranteed": true})
	if err := ValidateActivationSurface(doc); err == nil {
		t.Fatal("synthetic as guaranteed must be rejected")
	}
}

func TestValidateGuaranteedRecordRejectsOverGuarantee(t *testing.T) {
	rec := guaranteedRecord(map[string]interface{}{"guaranteedValue": 0.01, "netRealizedValue": 0.0042})
	doc := activationDoc(map[string]interface{}{"guaranteedRecords": []interface{}{rec}})
	if err := ValidateActivationSurface(doc); err == nil {
		t.Fatal("guaranteedValue > netRealizedValue must be rejected")
	}
}

func TestValidateGuaranteedRecordRejectsRollbackFailed(t *testing.T) {
	rec := guaranteedRecord(map[string]interface{}{"rollbackFailed": true})
	doc := activationDoc(map[string]interface{}{"guaranteedRecords": []interface{}{rec}})
	if err := ValidateActivationSurface(doc); err == nil {
		t.Fatal("rollbackFailed guaranteed record must be rejected")
	}
}

func TestValidateGuaranteedRecordRejectsWindowsEvidenceGated(t *testing.T) {
	rec := guaranteedRecord(map[string]interface{}{"windowsEvidenceGated": true})
	doc := activationDoc(map[string]interface{}{"guaranteedRecords": []interface{}{rec}})
	if err := ValidateActivationSurface(doc); err == nil {
		t.Fatal("windows evidence gated guaranteed record must be rejected")
	}
}

func TestValidateEligibilityCheckRejectsEstimatedEligible(t *testing.T) {
	doc := map[string]interface{}{
		"eligible": true, "eligibilityStatus": "eligible", "estimatedOnly": true,
		"evidenceRefs": []interface{}{"ref-1"}, "claimBoundary": "test",
	}
	if err := ValidateEligibilityCheck(doc); err == nil {
		t.Fatal("estimatedOnly eligible must be rejected")
	}
}

func TestValidateUnitPriceAuthorityRejectsPlaceholderActive(t *testing.T) {
	doc := map[string]interface{}{
		"sourceType": "placeholder_non_guarantee", "authorityStatus": "active",
		"evidenceRefs": []interface{}{"ref-1"}, "claimBoundary": "test",
	}
	if err := ValidateUnitPriceAuthority(doc); err == nil {
		t.Fatal("placeholder active authority must be rejected")
	}
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/guaranteed-eligible-savings-activation-reference.json"); err != nil {
		t.Fatalf("reference file invalid: %v", err)
	}
}

func TestScanForbiddenClaimsUniversalGuarantee(t *testing.T) {
	if err := scanForbiddenClaims("universal guaranteed savings for all workloads"); err == nil {
		t.Fatal("universal guarantee claim must be detected")
	}
}
