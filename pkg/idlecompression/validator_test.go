package idlecompression

import "testing"

func coverageDoc(overrides map[string]interface{}) map[string]interface{} {
	doc := map[string]interface{}{
		"milestone":                           MilestoneIdleTimeCompressionFleetValue,
		"idleCompressionMeasured":             true,
		"fleetValueCoverageMeasured":          true,
		"guaranteedCoverageMeasured":          true,
		"highSavingsOpportunityIdentified":      true,
		"universalGuaranteedSavingsAllowed":     false,
		"universalGuaranteedSavingsClaimed":     false,
		"generalProductionAutoAllowed":          false,
		"productionAutoWithPolicy":              false,
		"universalPerformanceImprovementClaimed": false,
		"estimatedIdleCountedAsMoved":           false,
		"syntheticFleetCountedAsProduction":     false,
		"referenceFleetCountedAsProduction":     false,
		"totalIdleValue":                        0.125,
		"eligibleIdleValue":                     0.085,
		"movedIdleValue":                        0.0042,
		"unmovedEligibleIdleValue":              0.0808,
		"realizedMovedIdleValue":                0.0042,
		"guaranteedEligibleSavingsTotal":      0.00336,
		"idleCompressionRate":                   0.0042 / 0.085,
		"guaranteeCoveragePercent":              0.00336 / 0.0042,
		"fleetLiquidityRate":                    0.085 / 0.125,
		"productionEvidenceSeparation": map[string]interface{}{
			"referenceCountedAsProduction": false, "syntheticCountedAsProduction": false,
			"canaryCountedAsGeneralProduction": false,
			"evidenceRefs": []interface{}{"ref-1"}, "claimBoundary": "separated",
		},
		"claimBoundaries": []interface{}{"scoped idle compression only; universal guarantee remains disabled"},
	}
	for k, v := range overrides {
		doc[k] = v
	}
	return doc
}

func TestValidateCoverageSurfaceRejectsUniversalGuarantee(t *testing.T) {
	doc := coverageDoc(map[string]interface{}{"universalGuaranteedSavingsAllowed": true})
	if err := ValidateCoverageSurface(doc); err == nil {
		t.Fatal("universal guarantee must be rejected")
	}
}

func TestValidateCoverageSurfaceRejectsEstimatedAsMoved(t *testing.T) {
	doc := coverageDoc(map[string]interface{}{"estimatedIdleCountedAsMoved": true})
	if err := ValidateCoverageSurface(doc); err == nil {
		t.Fatal("estimated as moved must be rejected")
	}
}

func TestValidateCoverageSurfaceRejectsSyntheticAsProduction(t *testing.T) {
	doc := coverageDoc(map[string]interface{}{"syntheticFleetCountedAsProduction": true})
	if err := ValidateCoverageSurface(doc); err == nil {
		t.Fatal("synthetic as production must be rejected")
	}
}

func TestValidateCoverageSurfaceRejectsCompressionRateMismatch(t *testing.T) {
	doc := coverageDoc(map[string]interface{}{"idleCompressionRate": 0.99})
	if err := ValidateCoverageSurface(doc); err == nil {
		t.Fatal("compression rate mismatch must be rejected")
	}
}

func TestValidateCoverageSurfaceRejectsGuaranteeCoverageMismatch(t *testing.T) {
	doc := coverageDoc(map[string]interface{}{"guaranteeCoveragePercent": 0.1})
	if err := ValidateCoverageSurface(doc); err == nil {
		t.Fatal("guarantee coverage mismatch must be rejected")
	}
}

func TestValidateCoverageSurfaceRejectsLiquidityMismatch(t *testing.T) {
	doc := coverageDoc(map[string]interface{}{"fleetLiquidityRate": 0.1})
	if err := ValidateCoverageSurface(doc); err == nil {
		t.Fatal("liquidity rate mismatch must be rejected")
	}
}

func TestValidateCoverageSurfaceRejectsEligibleLessThanMoved(t *testing.T) {
	doc := coverageDoc(map[string]interface{}{"eligibleIdleValue": 0.001, "movedIdleValue": 0.0042})
	if err := ValidateCoverageSurface(doc); err == nil {
		t.Fatal("eligible < moved must be rejected")
	}
}

func TestValidateFleetIdleObservationRejectsReferenceProduction(t *testing.T) {
	doc := map[string]interface{}{
		"referenceOnly": true, "productionEvidence": true,
		"evidenceRefs": []interface{}{"ref-1"}, "claimBoundary": "test",
	}
	if err := ValidateFleetIdleObservation(doc); err == nil {
		t.Fatal("reference production evidence must be rejected")
	}
}

func TestValidateMovedIdleResourceRejectsNoMutation(t *testing.T) {
	doc := map[string]interface{}{
		"mutationObserved": false, "movementSource": "operator_controlled",
		"evidenceRefs": []interface{}{"ref-1"}, "claimBoundary": "test",
	}
	if err := ValidateMovedIdleResource(doc); err == nil {
		t.Fatal("mutationObserved=false must be rejected")
	}
}

func TestValidateMovedIdleResourceRejectsSandboxAsProduction(t *testing.T) {
	doc := map[string]interface{}{
		"mutationObserved": true, "movementSource": "sandbox_auto",
		"productionEvidence": true, "environmentClass": "sandbox",
		"evidenceRefs": []interface{}{"ref-1"}, "claimBoundary": "test",
	}
	if err := ValidateMovedIdleResource(doc); err == nil {
		t.Fatal("sandbox as production must be rejected")
	}
}

func TestValidateProductionEvidenceSeparation(t *testing.T) {
	doc := map[string]interface{}{
		"referenceCountedAsProduction": true,
		"evidenceRefs": []interface{}{"ref-1"}, "claimBoundary": "test",
	}
	if err := ValidateProductionEvidenceSeparation(doc); err == nil {
		t.Fatal("reference as production must be rejected")
	}
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/idle-time-compression-fleet-value-coverage-reference.json"); err != nil {
		t.Fatalf("reference file invalid: %v", err)
	}
}
