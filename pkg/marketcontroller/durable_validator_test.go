package marketcontroller

import "testing"

func TestValidateDurableReferenceFile(t *testing.T) {
	if err := ValidateDurableReferenceFile("../../examples/durable-controller-state-kubernetes-reconciler-reference.json"); err != nil {
		t.Fatal(err)
	}
}

func TestValidateDurableRejectsForbiddenAuto(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": DurableMilestone,
		"durableStateStoreEnabled": true, "kubernetesReconcilerEnabled": true,
		"fakeClientTestsEnabled": true, "controllerStatusConditionsEnabled": true,
		"kubernetesEventsEnabled": true, "leaderElectionReady": true,
		"rbacManifestsDefined": true, "metricsExportDefined": true, "recoverySemanticsDefined": true,
		"generalProductionAutoAllowed": true,
		"stateStore": map[string]interface{}{"storeType": "kubernetes_configmap", "optimisticLockEnabled": true, "idempotencyPersists": true},
		"statusConditions": []interface{}{map[string]interface{}{"type": ConditionReady}},
		"kubernetesEvents": []interface{}{map[string]interface{}{"reason": "ControllerTickCompleted"}},
	}
	if err := ValidateDurableSurface(doc); err == nil {
		t.Fatal("expected rejection")
	}
}
