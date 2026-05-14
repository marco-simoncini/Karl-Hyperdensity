package installadmission

import "fmt"

func validateAdmissionSections(doc map[string]interface{}) error {
	rbac, _ := doc["rbacAdmissionResult"].(map[string]interface{})
	if rbac == nil {
		return fmt.Errorf("rbacAdmissionResult required")
	}
	if boolv(rbac["clusterAdmin"]) || !boolv(rbac["namespaceScoped"]) || boolv(rbac["podsExecAllowed"]) ||
		boolv(rbac["nodesWriteAllowed"]) || boolv(rbac["wildcardResourceAllowed"]) || boolv(rbac["wildcardVerbAllowed"]) ||
		boolv(rbac["rawRuntimeControlsAllowed"]) || boolv(rbac["directLibvirtAllowed"]) || boolv(rbac["directCgroupAllowed"]) {
		return fmt.Errorf("unsafe RBAC admission")
	}
	if strv(rbac["claimBoundary"]) == "" || !hasNonEmptyStringList(rbac["evidenceRefs"]) {
		return fmt.Errorf("rbacAdmissionResult missing evidenceRefs or claimBoundary")
	}

	manifest, _ := doc["manifestAdmissionResult"].(map[string]interface{})
	if manifest == nil {
		return fmt.Errorf("manifestAdmissionResult required")
	}
	if strv(manifest["claimBoundary"]) == "" || !hasNonEmptyStringList(manifest["evidenceRefs"]) {
		return fmt.Errorf("manifestAdmissionResult missing evidenceRefs or claimBoundary")
	}
	if boolv(manifest["referenceOnly"]) || boolv(manifest["notProductionInstall"]) {
		if boolv(doc["productionInstallEligible"]) {
			return fmt.Errorf("reference-only manifest cannot be eligible")
		}
	}

	for _, key := range []string{
		"namespacePreflight",
		"crdAdmissionResult",
		"durableStateAdmissionResult",
		"leaderElectionAdmissionResult",
		"probeAdmissionResult",
		"metricsAdmissionResult",
		"rollbackAdmissionResult",
		"clusterAdmissionGate",
		"productionInstallEligibility",
	} {
		m, _ := doc[key].(map[string]interface{})
		if m == nil {
			return fmt.Errorf("%s required", key)
		}
		if strv(m["claimBoundary"]) == "" || !hasNonEmptyStringList(m["evidenceRefs"]) {
			return fmt.Errorf("%s missing evidenceRefs or claimBoundary", key)
		}
	}

	leader := doc["leaderElectionAdmissionResult"].(map[string]interface{})
	if boolv(leader["haProductionProven"]) && !hasEvidenceRef(leader["evidenceRefs"], "ha") {
		return fmt.Errorf("haProductionProven without evidence")
	}

	metrics := doc["metricsAdmissionResult"].(map[string]interface{})
	if num(metrics["generalProductionAutoEnabledGauge"]) > 0 || num(metrics["productionAutoWithPolicyEnabledGauge"]) > 0 ||
		boolv(metrics["forbiddenAutoGaugesPresent"]) {
		return fmt.Errorf("forbidden auto gauge present")
	}

	eligibility := doc["productionInstallEligibility"].(map[string]interface{})
	if strv(eligibility["nextAllowedState"]) == "production_auto_with_policy" {
		return fmt.Errorf("nextAllowedState cannot be production_auto_with_policy")
	}
	if boolv(eligibility["productionInstallApplied"]) && !boolv(doc["realApplyEvidencePresent"]) {
		return fmt.Errorf("productionInstallApplied without evidence")
	}
	return nil
}
