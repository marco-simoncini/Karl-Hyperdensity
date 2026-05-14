package installruntime

import "fmt"

func validateHealthSections(doc map[string]interface{}) error {
	for _, key := range []string{
		"runtimeHealthGate",
		"deploymentHealth",
		"podReadiness",
		"probeVerification",
		"metricsReachability",
		"leaderElectionRuntime",
		"durableStateRuntime",
		"statusConditionRuntime",
		"eventRuntime",
		"rbacRuntimeVerification",
		"rollbackReadiness",
	} {
		m, _ := doc[key].(map[string]interface{})
		if m == nil {
			return fmt.Errorf("%s required", key)
		}
		if strv(m["claimBoundary"]) == "" || !hasNonEmptyStringList(m["evidenceRefs"]) {
			return fmt.Errorf("%s missing evidenceRefs or claimBoundary", key)
		}
	}

	rbac := doc["rbacRuntimeVerification"].(map[string]interface{})
	if boolv(rbac["clusterAdmin"]) || boolv(rbac["podsExecAllowed"]) || boolv(rbac["nodesWriteAllowed"]) ||
		boolv(rbac["rawRuntimeControlsAllowed"]) || boolv(rbac["directLibvirtAllowed"]) || boolv(rbac["directCgroupAllowed"]) {
		return fmt.Errorf("unsafe RBAC runtime")
	}

	metrics := doc["metricsReachability"].(map[string]interface{})
	if num(metrics["generalProductionAutoEnabledGauge"]) > 0 || num(metrics["productionAutoWithPolicyEnabledGauge"]) > 0 {
		return fmt.Errorf("forbidden auto gauge present")
	}

	leader := doc["leaderElectionRuntime"].(map[string]interface{})
	if boolv(leader["haProductionProven"]) && !hasEvidenceRef(leader["evidenceRefs"], "ha") {
		return fmt.Errorf("haProductionProven without evidence")
	}

	if boolv(doc["largeFleetProductionProven"]) && !hasEvidenceRef(doc["fleetEvidenceRefs"], "fleet") {
		return fmt.Errorf("largeFleetProductionProven without evidence")
	}
	return nil
}

func evaluateRuntimeHealthGate(doc map[string]interface{}) error {
	if !boolv(doc["runtimeHealthGateEnabled"]) {
		return fmt.Errorf("runtimeHealthGateEnabled must be true")
	}
	if !boolv(doc["controlledInstallScope"]) {
		return fmt.Errorf("controlledInstallScope must be true")
	}

	mode := strv(doc["installApplyMode"])
	if mode != "real_cluster_apply" {
		if boolv(doc["runtimeHealthGatePassed"]) {
			return fmt.Errorf("runtime health passed with non-real apply mode")
		}
		if boolv(doc["productionInstallApplied"]) {
			return fmt.Errorf("productionInstallApplied with non-real apply mode")
		}
		return nil
	}

	if boolv(doc["productionInstallApplied"]) {
		if !boolv(doc["realApplyEvidencePresent"]) {
			return fmt.Errorf("productionInstallApplied without evidence")
		}
	}
	if boolv(doc["runtimeHealthGatePassed"]) {
		if !boolv(doc["productionInstallApplied"]) {
			return fmt.Errorf("runtime health passed without install applied")
		}
		for _, key := range []string{
			"deploymentAvailable",
			"podsReady",
			"probesHealthy",
			"metricsReachable",
			"leaderElectionLeaseObserved",
			"durableStateReadWriteVerified",
			"rollbackReady",
			"failClosedVerified",
		} {
			if !boolv(doc[key]) {
				return fmt.Errorf("runtime health passed with failed check: %s", key)
			}
		}
	}
	return nil
}
