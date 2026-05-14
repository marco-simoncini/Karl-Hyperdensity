package installadmission

import "testing"

func TestValidateAdmissionSectionsRejectsUnsafeRBAC(t *testing.T) {
	doc := minimalAdmissionDoc()
	doc["rbacAdmissionResult"].(map[string]interface{})["clusterAdmin"] = true
	if err := validateAdmissionSections(doc); err == nil {
		t.Fatal("expected unsafe RBAC rejection")
	}
}

func TestValidateAdmissionSectionsRejectsForbiddenMetricsGauge(t *testing.T) {
	doc := minimalAdmissionDoc()
	doc["metricsAdmissionResult"].(map[string]interface{})["generalProductionAutoEnabledGauge"] = float64(1)
	if err := validateAdmissionSections(doc); err == nil {
		t.Fatal("expected forbidden metrics gauge rejection")
	}
}

func TestValidateAdmissionSectionsRejectsReferenceEligible(t *testing.T) {
	doc := minimalAdmissionDoc()
	doc["productionInstallEligible"] = true
	doc["manifestAdmissionResult"].(map[string]interface{})["referenceOnly"] = true
	if err := validateAdmissionSections(doc); err == nil {
		t.Fatal("expected reference-only eligible rejection")
	}
}

func minimalAdmissionDoc() map[string]interface{} {
	return map[string]interface{}{
		"productionInstallEligible": false,
		"rbacAdmissionResult": map[string]interface{}{
			"clusterAdmin": false, "namespaceScoped": true, "podsExecAllowed": false, "nodesWriteAllowed": false,
			"wildcardResourceAllowed": false, "wildcardVerbAllowed": false, "rawRuntimeControlsAllowed": false, "directLibvirtAllowed": false, "directCgroupAllowed": false,
			"evidenceRefs": []interface{}{"rbac"}, "claimBoundary": "rbac admission",
		},
		"manifestAdmissionResult": map[string]interface{}{
			"referenceOnly": false, "notProductionInstall": false,
			"evidenceRefs": []interface{}{"manifest"}, "claimBoundary": "manifest admission",
		},
		"namespacePreflight": map[string]interface{}{"evidenceRefs": []interface{}{"ns"}, "claimBoundary": "namespace preflight"},
		"crdAdmissionResult": map[string]interface{}{"evidenceRefs": []interface{}{"crd"}, "claimBoundary": "crd admission"},
		"durableStateAdmissionResult": map[string]interface{}{"evidenceRefs": []interface{}{"state"}, "claimBoundary": "state admission"},
		"leaderElectionAdmissionResult": map[string]interface{}{"haProductionProven": false, "evidenceRefs": []interface{}{"leader"}, "claimBoundary": "leader admission"},
		"probeAdmissionResult": map[string]interface{}{"evidenceRefs": []interface{}{"probe"}, "claimBoundary": "probe admission"},
		"metricsAdmissionResult": map[string]interface{}{
			"generalProductionAutoEnabledGauge": float64(0), "productionAutoWithPolicyEnabledGauge": float64(0), "forbiddenAutoGaugesPresent": false,
			"evidenceRefs": []interface{}{"metrics"}, "claimBoundary": "metrics admission",
		},
		"rollbackAdmissionResult": map[string]interface{}{"evidenceRefs": []interface{}{"rollback"}, "claimBoundary": "rollback admission"},
		"clusterAdmissionGate": map[string]interface{}{"evidenceRefs": []interface{}{"gate"}, "claimBoundary": "cluster admission gate"},
		"productionInstallEligibility": map[string]interface{}{
			"productionInstallApplied": false, "nextAllowedState": "production_canary_only",
			"evidenceRefs": []interface{}{"eligibility"}, "claimBoundary": "eligibility",
		},
	}
}
