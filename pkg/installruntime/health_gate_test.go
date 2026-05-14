package installruntime

import "testing"

func TestValidateHealthSectionsRejectsUnsafeRBAC(t *testing.T) {
	doc := minimalHealthDoc()
	doc["rbacRuntimeVerification"].(map[string]interface{})["clusterAdmin"] = true
	if err := validateHealthSections(doc); err == nil {
		t.Fatal("expected unsafe RBAC rejection")
	}
}

func TestEvaluateRuntimeHealthGateRejectsHealthWithoutApply(t *testing.T) {
	doc := validSurface()
	doc["runtimeHealthGatePassed"] = true
	doc["productionInstallApplied"] = false
	if err := evaluateRuntimeHealthGate(doc); err == nil {
		t.Fatal("expected health without apply rejection")
	}
}

func minimalHealthDoc() map[string]interface{} {
	return map[string]interface{}{
		"runtimeHealthGate": map[string]interface{}{"evidenceRefs": []interface{}{"gate"}, "claimBoundary": "gate"},
		"deploymentHealth": map[string]interface{}{"evidenceRefs": []interface{}{"dep"}, "claimBoundary": "dep"},
		"podReadiness": map[string]interface{}{"evidenceRefs": []interface{}{"pod"}, "claimBoundary": "pod"},
		"probeVerification": map[string]interface{}{"evidenceRefs": []interface{}{"probe"}, "claimBoundary": "probe"},
		"metricsReachability": map[string]interface{}{
			"generalProductionAutoEnabledGauge": float64(0), "productionAutoWithPolicyEnabledGauge": float64(0),
			"evidenceRefs": []interface{}{"metrics"}, "claimBoundary": "metrics",
		},
		"leaderElectionRuntime": map[string]interface{}{"haProductionProven": false, "evidenceRefs": []interface{}{"leader"}, "claimBoundary": "leader"},
		"durableStateRuntime": map[string]interface{}{"evidenceRefs": []interface{}{"state"}, "claimBoundary": "state"},
		"statusConditionRuntime": map[string]interface{}{"evidenceRefs": []interface{}{"status"}, "claimBoundary": "status"},
		"eventRuntime": map[string]interface{}{"evidenceRefs": []interface{}{"event"}, "claimBoundary": "event"},
		"rbacRuntimeVerification": map[string]interface{}{
			"clusterAdmin": false, "podsExecAllowed": false, "nodesWriteAllowed": false,
			"rawRuntimeControlsAllowed": false, "directLibvirtAllowed": false, "directCgroupAllowed": false,
			"evidenceRefs": []interface{}{"rbac"}, "claimBoundary": "rbac",
		},
		"rollbackReadiness": map[string]interface{}{"evidenceRefs": []interface{}{"rollback"}, "claimBoundary": "rollback"},
	}
}
