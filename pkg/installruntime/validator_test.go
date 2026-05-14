package installruntime

import "testing"

func validSurface() map[string]interface{} {
	return map[string]interface{}{
		"milestone": Milestone, "surfaceVersion": "v1", "runtimeHealthGateId": "health-gate-1",
		"controlledInstallScope": true, "installApplyEnabled": true, "installApplyMode": "real_cluster_apply",
		"productionInstallEligible": true, "productionInstallApplied": true, "realApplyEvidencePresent": true,
		"installApplySucceeded": true, "runtimeHealthGateEnabled": true, "runtimeHealthGatePassed": true,
		"deploymentAvailable": true, "podsReady": true, "probesHealthy": true, "metricsReachable": true,
		"leaderElectionLeaseObserved": true, "durableStateReadWriteVerified": true, "statusConditionsObserved": true,
		"eventsObserved": true, "rbacRuntimeVerified": true, "rollbackReady": true, "failClosedVerified": true,
		"degradedModeActive": false, "largeFleetProductionProven": false,
		"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
		"universalGuaranteedSavingsAllowed": false, "universalGuaranteedSavingsClaimed": false,
		"estimatedIdleCountedAsMoved": false, "projectedCompressionCountedAsRealized": false,
		"syntheticFleetCountedAsProduction": false, "referenceFleetCountedAsProduction": false,
		"dashboardExecutor": false, "fluidvirtPolicyAuthority": false, "fluidvirtInstallAuthority": false,
		"inventoryRuntimeExecutor": false,
		"installApplyRequest": map[string]interface{}{
			"applyRequestId": "apply-req-1", "admissionGatePassed": true, "controlledInstallScope": true,
			"applyMode": "real_cluster_apply", "productionInstallApplied": true,
			"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
			"evidenceRefs": []interface{}{"apply-req"}, "claimBoundary": "apply request",
		},
		"installApplyResult": map[string]interface{}{
			"applyResultId": "apply-res-1", "status": "applied", "productionInstallApplied": true,
			"realApplyEvidencePresent": true, "evidenceRefs": []interface{}{"apply-res"}, "claimBoundary": "apply result",
		},
		"appliedManifestInventory": map[string]interface{}{
			"inventoryId": "inv-1", "namespace": "karl-system",
			"evidenceRefs": []interface{}{"inventory"}, "claimBoundary": "inventory",
		},
		"runtimeHealthGate": map[string]interface{}{
			"healthGateId": "hg-1", "decision": "healthy", "runtimeHealthGatePassed": true,
			"evidenceRefs": []interface{}{"health-gate"}, "claimBoundary": "health gate",
		},
		"deploymentHealth": map[string]interface{}{"deploymentAvailable": true, "evidenceRefs": []interface{}{"dep"}, "claimBoundary": "deployment"},
		"podReadiness": map[string]interface{}{"podsReady": true, "evidenceRefs": []interface{}{"pod"}, "claimBoundary": "pod"},
		"probeVerification": map[string]interface{}{"probesHealthy": true, "evidenceRefs": []interface{}{"probe"}, "claimBoundary": "probe"},
		"metricsReachability": map[string]interface{}{
			"reachable": true, "generalProductionAutoEnabledGauge": float64(0), "productionAutoWithPolicyEnabledGauge": float64(0),
			"evidenceRefs": []interface{}{"metrics"}, "claimBoundary": "metrics",
		},
		"leaderElectionRuntime": map[string]interface{}{
			"leaderElectionLeaseObserved": true, "haProductionProven": false,
			"evidenceRefs": []interface{}{"leader"}, "claimBoundary": "leader",
		},
		"durableStateRuntime": map[string]interface{}{
			"durableStateReadWriteVerified": true, "evidenceRefs": []interface{}{"state"}, "claimBoundary": "state",
		},
		"statusConditionRuntime": map[string]interface{}{"statusConditionsObserved": true, "evidenceRefs": []interface{}{"status"}, "claimBoundary": "status"},
		"eventRuntime": map[string]interface{}{"eventsObserved": true, "evidenceRefs": []interface{}{"event"}, "claimBoundary": "event"},
		"rbacRuntimeVerification": map[string]interface{}{
			"clusterAdmin": false, "podsExecAllowed": false, "nodesWriteAllowed": false,
			"rawRuntimeControlsAllowed": false, "directLibvirtAllowed": false, "directCgroupAllowed": false,
			"runtimeRbacSafe": true, "evidenceRefs": []interface{}{"rbac"}, "claimBoundary": "rbac",
		},
		"rollbackReadiness": map[string]interface{}{
			"rollbackReady": true, "evidenceRefs": []interface{}{"rollback"}, "claimBoundary": "rollback",
		},
		"postInstallAuditEvents": []interface{}{
			map[string]interface{}{"auditEventId": "evt-1", "eventType": "install_apply_succeeded", "evidenceRefs": []interface{}{"evt"}, "claimBoundary": "audit"},
		},
		"blockers": []interface{}{},
		"claimBoundaries": []interface{}{"controlled apply and runtime health only"},
	}
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/controller-install-apply-runtime-health-gate-reference.json"); err != nil {
		t.Fatal(err)
	}
}

func TestValidatorNegativeCases(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]interface{})
	}{
		{"health_gate_disabled", func(d map[string]interface{}) { d["runtimeHealthGateEnabled"] = false }},
		{"controlled_scope_disabled", func(d map[string]interface{}) { d["controlledInstallScope"] = false }},
		{"health_passed_without_apply", func(d map[string]interface{}) { d["productionInstallApplied"] = false }},
		{"applied_without_evidence", func(d map[string]interface{}) { d["realApplyEvidencePresent"] = false }},
		{"applied_wrong_mode", func(d map[string]interface{}) { d["installApplyMode"] = "server_dryrun_only" }},
		{"deployment_unavailable", func(d map[string]interface{}) { d["deploymentAvailable"] = false }},
		{"pods_not_ready", func(d map[string]interface{}) { d["podsReady"] = false }},
		{"probes_unhealthy", func(d map[string]interface{}) { d["probesHealthy"] = false }},
		{"metrics_unreachable", func(d map[string]interface{}) { d["metricsReachable"] = false }},
		{"leader_missing", func(d map[string]interface{}) { d["leaderElectionLeaseObserved"] = false }},
		{"state_unverified", func(d map[string]interface{}) { d["durableStateReadWriteVerified"] = false }},
		{"rollback_not_ready", func(d map[string]interface{}) { d["rollbackReady"] = false }},
		{"fail_closed_unverified", func(d map[string]interface{}) { d["failClosedVerified"] = false }},
		{"admission_not_passed", func(d map[string]interface{}) { d["installApplyRequest"].(map[string]interface{})["admissionGatePassed"] = false }},
		{"cluster_admin", func(d map[string]interface{}) { d["rbacRuntimeVerification"].(map[string]interface{})["clusterAdmin"] = true }},
		{"pods_exec_allowed", func(d map[string]interface{}) { d["rbacRuntimeVerification"].(map[string]interface{})["podsExecAllowed"] = true }},
		{"nodes_write_allowed", func(d map[string]interface{}) { d["rbacRuntimeVerification"].(map[string]interface{})["nodesWriteAllowed"] = true }},
		{"raw_runtime_controls", func(d map[string]interface{}) { d["rbacRuntimeVerification"].(map[string]interface{})["rawRuntimeControlsAllowed"] = true }},
		{"direct_libvirt", func(d map[string]interface{}) { d["rbacRuntimeVerification"].(map[string]interface{})["directLibvirtAllowed"] = true }},
		{"direct_cgroup", func(d map[string]interface{}) { d["rbacRuntimeVerification"].(map[string]interface{})["directCgroupAllowed"] = true }},
		{"general_auto_enabled", func(d map[string]interface{}) { d["generalProductionAutoAllowed"] = true }},
		{"prod_auto_policy", func(d map[string]interface{}) { d["productionAutoWithPolicy"] = true }},
		{"universal_savings_allowed", func(d map[string]interface{}) { d["universalGuaranteedSavingsAllowed"] = true }},
		{"universal_savings_claimed", func(d map[string]interface{}) { d["universalGuaranteedSavingsClaimed"] = true }},
		{"projected_realized", func(d map[string]interface{}) { d["projectedCompressionCountedAsRealized"] = true }},
		{"estimated_idle_moved", func(d map[string]interface{}) { d["estimatedIdleCountedAsMoved"] = true }},
		{"synthetic_prod", func(d map[string]interface{}) { d["syntheticFleetCountedAsProduction"] = true }},
		{"reference_prod", func(d map[string]interface{}) { d["referenceFleetCountedAsProduction"] = true }},
		{"dashboard_executor", func(d map[string]interface{}) { d["dashboardExecutor"] = true }},
		{"fluidvirt_policy", func(d map[string]interface{}) { d["fluidvirtPolicyAuthority"] = true }},
		{"fluidvirt_install", func(d map[string]interface{}) { d["fluidvirtInstallAuthority"] = true }},
		{"inventory_executor", func(d map[string]interface{}) { d["inventoryRuntimeExecutor"] = true }},
		{"forbidden_gauge", func(d map[string]interface{}) { d["metricsReachability"].(map[string]interface{})["generalProductionAutoEnabledGauge"] = float64(1) }},
		{"ha_without_evidence", func(d map[string]interface{}) { d["leaderElectionRuntime"].(map[string]interface{})["haProductionProven"] = true }},
		{"large_fleet_without_evidence", func(d map[string]interface{}) { d["largeFleetProductionProven"] = true }},
		{"missing_evidence_refs", func(d map[string]interface{}) { d["rbacRuntimeVerification"].(map[string]interface{})["evidenceRefs"] = []interface{}{} }},
		{"missing_claim_boundary", func(d map[string]interface{}) { d["rbacRuntimeVerification"].(map[string]interface{})["claimBoundary"] = "" }},
		{"health_with_blockers", func(d map[string]interface{}) { d["blockers"] = []interface{}{map[string]interface{}{"blockerCode": "x"}} }},
		{"ineligible_applied", func(d map[string]interface{}) { d["productionInstallEligible"] = false }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			doc := validSurface()
			tc.mutate(doc)
			if err := ValidateSurface(doc); err == nil {
				t.Fatalf("expected rejection for case %s", tc.name)
			}
		})
	}
}
