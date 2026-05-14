package marketcontroller

import "testing"

func TestValidateControllerSurfaceRejectsNoFullNxN(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": Milestone, "noFullNxNPairing": false,
		"continuousControllerEnabled": true, "controllerTickExecuted": true,
		"actionSlateGeneratedByController": true, "resourceFuturesGeneratedByController": true,
		"controllerCoverageExpanded": true, "idleCompressionTargetTracked": true,
		"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
		"universalGuaranteedSavingsAllowed": false, "universalGuaranteedSavingsClaimed": false,
		"estimatedIdleCountedAsMoved": false, "syntheticFleetCountedAsProduction": false,
		"referenceFleetCountedAsProduction": false, "dashboardExecutor": false,
		"fluidvirtPolicyAuthority": false, "inventoryRuntimeExecutor": false,
		"fullPairSpace": 96.0, "evaluatedPairCount": 96.0, "avoidedPairCount": 0.0,
		"topKDonors": 3.0, "topKReceivers": 3.0,
		"currentIdleCompressionRate": 0.04, "projectedIdleCompressionRate": 0.12,
		"backpressure": map[string]interface{}{}, "invalidationRules": []interface{}{map[string]interface{}{}},
		"claimBoundaries": []interface{}{"test"},
	}
	if err := ValidateControllerSurface(doc); err == nil {
		t.Fatal("noFullNxNPairing=false must be rejected")
	}
}

func TestValidateGeneratedActionRejectsGeneralAuto(t *testing.T) {
	doc := map[string]interface{}{
		"actionId": "a", "donorShellId": "d", "receiverShellId": "r", "resource": "cpu", "amount": "1",
		"generalProductionAutoAllowed": true, "productionAutoWithPolicy": false,
		"evidenceRefs": []interface{}{"ref"}, "claimBoundary": "test",
	}
	if err := ValidateGeneratedAction(doc); err == nil {
		t.Fatal("generalProductionAutoAllowed must be rejected")
	}
}

func TestValidateGeneratedFutureMissingExpiration(t *testing.T) {
	doc := map[string]interface{}{
		"evidenceRefs": []interface{}{"ref"}, "claimBoundary": "test",
		"invalidationReasons": []interface{}{"expired"},
	}
	if err := ValidateGeneratedFuture(doc); err == nil {
		t.Fatal("missing expiration must be rejected")
	}
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/continuous-resource-market-controller-reference.json"); err != nil {
		t.Fatalf("reference invalid: %v", err)
	}
}
