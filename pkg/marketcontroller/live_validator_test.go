package marketcontroller

import "testing"

func TestValidateLiveReferenceFile(t *testing.T) {
	if err := ValidateLiveReferenceFile("../../examples/live-controller-reconciliation-execution-loop-reference.json"); err != nil {
		t.Fatal(err)
	}
}

func TestValidateLiveSurfaceRejectsProjectedAsRealized(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": LiveMilestone,
		"liveReconciliationEnabled": true, "stateStoreEnabled": true,
		"scheduledTicksEnabled": true, "leaseLifecycleEnabled": true,
		"actionLifecycleEnabled": true, "futuresRefreshEnabled": true,
		"executionSelectionEnabled": true, "realizedCompressionTrackingEnabled": true,
		"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
		"universalGuaranteedSavingsAllowed": false, "universalGuaranteedSavingsClaimed": false,
		"estimatedIdleCountedAsMoved": false, "projectedCompressionCountedAsRealized": true,
		"syntheticFleetCountedAsProduction": false, "referenceFleetCountedAsProduction": false,
		"dashboardExecutor": false, "fluidvirtPolicyAuthority": false, "inventoryRuntimeExecutor": false,
		"controllerMode": "production_canary_only",
		"realizedCompressionTracker": map[string]interface{}{"projectedCompressionCountedAsRealized": true},
		"auditTrail": map[string]interface{}{"events": []interface{}{}},
	}
	if err := ValidateLiveSurface(doc); err == nil {
		t.Fatal("expected projected-as-realized rejection")
	}
}

func TestValidateLiveSurfaceRejectsForbiddenScope(t *testing.T) {
	doc := baseValidLiveDoc()
	doc["executionSelections"] = []interface{}{
		map[string]interface{}{"executionScope": "general_production_auto"},
	}
	if err := ValidateLiveSurface(doc); err == nil {
		t.Fatal("expected forbidden scope rejection")
	}
}

func TestValidateLiveSurfaceRejectsRealizedWithoutEvidence(t *testing.T) {
	doc := baseValidLiveDoc()
	doc["postExecutionReconciliations"] = []interface{}{
		map[string]interface{}{"realizedMovementKept": true, "mutationObserved": false, "postVerifyPassed": false},
	}
	if err := ValidateLiveSurface(doc); err == nil {
		t.Fatal("expected realized without evidence rejection")
	}
}

func baseValidLiveDoc() map[string]interface{} {
	events := make([]interface{}, 0, len(requiredAuditEventTypes))
	for _, et := range requiredAuditEventTypes {
		events = append(events, map[string]interface{}{"eventType": et})
	}
	return map[string]interface{}{
		"milestone": LiveMilestone,
		"liveReconciliationEnabled": true, "stateStoreEnabled": true,
		"scheduledTicksEnabled": true, "leaseLifecycleEnabled": true,
		"actionLifecycleEnabled": true, "futuresRefreshEnabled": true,
		"executionSelectionEnabled": true, "realizedCompressionTrackingEnabled": true,
		"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
		"universalGuaranteedSavingsAllowed": false, "universalGuaranteedSavingsClaimed": false,
		"estimatedIdleCountedAsMoved": false, "projectedCompressionCountedAsRealized": false,
		"syntheticFleetCountedAsProduction": false, "referenceFleetCountedAsProduction": false,
		"dashboardExecutor": false, "fluidvirtPolicyAuthority": false, "inventoryRuntimeExecutor": false,
		"controllerMode": "production_canary_only",
		"realizedCompressionTracker": map[string]interface{}{"projectedCompressionCountedAsRealized": false},
		"auditTrail": map[string]interface{}{"events": events},
	}
}
