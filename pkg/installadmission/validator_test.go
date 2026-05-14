package installadmission

import "testing"

func validSurface() map[string]interface{} {
	return map[string]interface{}{
		"milestone": Milestone, "surfaceVersion": "v1", "admissionGateId": "gate-1",
		"productionInstallDryRunEnabled": true, "clientSideDryRunModeled": true, "serverSideDryRunModeled": true,
		"clusterAdmissionGateEnabled": true, "productionInstallDryRunPassed": true, "clusterAdmissionGatePassed": true,
		"productionInstallEligible": true, "productionInstallApplied": false, "realApplyEvidencePresent": false, "installFailClosed": true,
		"namespacePreflightPassed": true, "rbacPreflightPassed": true, "manifestAdmissionPassed": true, "crdPreflightPassed": true,
		"durableStatePreflightPassed": true, "leaderElectionPreflightPassed": true, "probesPreflightPassed": true, "metricsPreflightPassed": true, "rollbackPreflightPassed": true,
		"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false, "universalGuaranteedSavingsAllowed": false, "universalGuaranteedSavingsClaimed": false,
		"estimatedIdleCountedAsMoved": false, "projectedCompressionCountedAsRealized": false, "syntheticFleetCountedAsProduction": false, "referenceFleetCountedAsProduction": false,
		"dashboardExecutor": false, "fluidvirtPolicyAuthority": false, "fluidvirtAdmissionAuthority": false, "inventoryRuntimeExecutor": false,
		"dryRunRequests": []interface{}{
			map[string]interface{}{"dryRunRequestId": "req-client", "productionInstallApplied": false, "evidenceRefs": []interface{}{"req-client"}, "claimBoundary": "dry-run request"},
			map[string]interface{}{"dryRunRequestId": "req-server", "productionInstallApplied": false, "evidenceRefs": []interface{}{"req-server"}, "claimBoundary": "dry-run request"},
		},
		"dryRunResults": []interface{}{
			map[string]interface{}{"dryRunResultId": "res-client", "dryRunMode": "client", "clientSideDryRunPassed": true, "evidenceRefs": []interface{}{"res-client"}, "claimBoundary": "dry-run result"},
			map[string]interface{}{"dryRunResultId": "res-server", "dryRunMode": "server", "serverSideDryRunPassed": true, "evidenceRefs": []interface{}{"res-server"}, "claimBoundary": "dry-run result"},
		},
		"clusterAdmissionGate": map[string]interface{}{"claimBoundary": "cluster admission gate", "evidenceRefs": []interface{}{"gate"}},
		"namespacePreflight":   map[string]interface{}{"claimBoundary": "namespace preflight", "evidenceRefs": []interface{}{"ns"}},
		"rbacAdmissionResult": map[string]interface{}{
			"clusterAdmin": false, "namespaceScoped": true, "podsExecAllowed": false, "nodesWriteAllowed": false,
			"wildcardResourceAllowed": false, "wildcardVerbAllowed": false, "rawRuntimeControlsAllowed": false,
			"directLibvirtAllowed": false, "directCgroupAllowed": false,
			"claimBoundary": "rbac admission", "evidenceRefs": []interface{}{"rbac"},
		},
		"manifestAdmissionResult": map[string]interface{}{"referenceOnly": false, "notProductionInstall": false, "claimBoundary": "manifest admission", "evidenceRefs": []interface{}{"manifest"}},
		"crdAdmissionResult":      map[string]interface{}{"claimBoundary": "crd admission", "evidenceRefs": []interface{}{"crd"}},
		"durableStateAdmissionResult": map[string]interface{}{"claimBoundary": "durable state admission", "evidenceRefs": []interface{}{"durable"}},
		"leaderElectionAdmissionResult": map[string]interface{}{"haProductionProven": false, "claimBoundary": "leader election admission", "evidenceRefs": []interface{}{"leader"}},
		"probeAdmissionResult": map[string]interface{}{"claimBoundary": "probe admission", "evidenceRefs": []interface{}{"probe"}},
		"metricsAdmissionResult": map[string]interface{}{
			"forbiddenAutoGaugesPresent": false, "generalProductionAutoEnabledGauge": float64(0), "productionAutoWithPolicyEnabledGauge": float64(0),
			"claimBoundary": "metrics admission", "evidenceRefs": []interface{}{"metrics"},
		},
		"rollbackAdmissionResult": map[string]interface{}{"claimBoundary": "rollback admission", "evidenceRefs": []interface{}{"rollback"}},
		"productionInstallEligibility": map[string]interface{}{
			"productionInstallApplied": false, "nextAllowedState": "production_canary_only",
			"claimBoundary": "eligibility", "evidenceRefs": []interface{}{"eligibility"},
		},
		"auditEvents": []interface{}{"dryrun_requested"},
		"admissionBlockers": []interface{}{},
		"claimBoundaries":   []interface{}{"dry-run validates eligibility only"},
	}
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/controller-install-dryrun-cluster-admission-gate-reference.json"); err != nil {
		t.Fatal(err)
	}
}

func TestValidatorNegativeCases(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]interface{})
	}{
		{"dryrun_disabled", func(d map[string]interface{}) { d["productionInstallDryRunEnabled"] = false }},
		{"admission_disabled", func(d map[string]interface{}) { d["clusterAdmissionGateEnabled"] = false }},
		{"fail_closed_disabled", func(d map[string]interface{}) { d["installFailClosed"] = false }},
		{"eligible_with_blockers", func(d map[string]interface{}) { d["admissionBlockers"] = []interface{}{map[string]interface{}{"blockerCode": "x"}} }},
		{"eligible_server_failed", func(d map[string]interface{}) { d["dryRunResults"].([]interface{})[1].(map[string]interface{})["serverSideDryRunPassed"] = false }},
		{"eligible_client_failed", func(d map[string]interface{}) { d["dryRunResults"].([]interface{})[0].(map[string]interface{})["clientSideDryRunPassed"] = false }},
		{"eligible_namespace_failed", func(d map[string]interface{}) { d["namespacePreflightPassed"] = false }},
		{"eligible_rbac_failed", func(d map[string]interface{}) { d["rbacPreflightPassed"] = false }},
		{"eligible_manifest_failed", func(d map[string]interface{}) { d["manifestAdmissionPassed"] = false }},
		{"eligible_crd_failed", func(d map[string]interface{}) { d["crdPreflightPassed"] = false }},
		{"eligible_state_failed", func(d map[string]interface{}) { d["durableStatePreflightPassed"] = false }},
		{"eligible_leader_failed", func(d map[string]interface{}) { d["leaderElectionPreflightPassed"] = false }},
		{"eligible_probe_failed", func(d map[string]interface{}) { d["probesPreflightPassed"] = false }},
		{"eligible_metrics_failed", func(d map[string]interface{}) { d["metricsPreflightPassed"] = false }},
		{"eligible_rollback_failed", func(d map[string]interface{}) { d["rollbackPreflightPassed"] = false }},
		{"eligible_reference_only", func(d map[string]interface{}) { d["manifestAdmissionResult"].(map[string]interface{})["referenceOnly"] = true }},
		{"eligible_not_production_install", func(d map[string]interface{}) { d["manifestAdmissionResult"].(map[string]interface{})["notProductionInstall"] = true }},
		{"applied_without_evidence", func(d map[string]interface{}) { d["productionInstallApplied"] = true }},
		{"cluster_admin", func(d map[string]interface{}) { d["rbacAdmissionResult"].(map[string]interface{})["clusterAdmin"] = true }},
		{"namespace_not_scoped", func(d map[string]interface{}) { d["rbacAdmissionResult"].(map[string]interface{})["namespaceScoped"] = false }},
		{"pods_exec_allowed", func(d map[string]interface{}) { d["rbacAdmissionResult"].(map[string]interface{})["podsExecAllowed"] = true }},
		{"nodes_write_allowed", func(d map[string]interface{}) { d["rbacAdmissionResult"].(map[string]interface{})["nodesWriteAllowed"] = true }},
		{"wildcard_resource_allowed", func(d map[string]interface{}) { d["rbacAdmissionResult"].(map[string]interface{})["wildcardResourceAllowed"] = true }},
		{"wildcard_verb_allowed", func(d map[string]interface{}) { d["rbacAdmissionResult"].(map[string]interface{})["wildcardVerbAllowed"] = true }},
		{"raw_runtime_controls_allowed", func(d map[string]interface{}) { d["rbacAdmissionResult"].(map[string]interface{})["rawRuntimeControlsAllowed"] = true }},
		{"direct_libvirt_allowed", func(d map[string]interface{}) { d["rbacAdmissionResult"].(map[string]interface{})["directLibvirtAllowed"] = true }},
		{"direct_cgroup_allowed", func(d map[string]interface{}) { d["rbacAdmissionResult"].(map[string]interface{})["directCgroupAllowed"] = true }},
		{"general_auto_enabled", func(d map[string]interface{}) { d["generalProductionAutoAllowed"] = true }},
		{"prod_auto_policy_enabled", func(d map[string]interface{}) { d["productionAutoWithPolicy"] = true }},
		{"universal_savings_allowed", func(d map[string]interface{}) { d["universalGuaranteedSavingsAllowed"] = true }},
		{"universal_savings_claimed", func(d map[string]interface{}) { d["universalGuaranteedSavingsClaimed"] = true }},
		{"estimated_idle_moved", func(d map[string]interface{}) { d["estimatedIdleCountedAsMoved"] = true }},
		{"projected_realized", func(d map[string]interface{}) { d["projectedCompressionCountedAsRealized"] = true }},
		{"synthetic_prod", func(d map[string]interface{}) { d["syntheticFleetCountedAsProduction"] = true }},
		{"reference_prod", func(d map[string]interface{}) { d["referenceFleetCountedAsProduction"] = true }},
		{"dashboard_executor", func(d map[string]interface{}) { d["dashboardExecutor"] = true }},
		{"fluidvirt_policy_authority", func(d map[string]interface{}) { d["fluidvirtPolicyAuthority"] = true }},
		{"fluidvirt_admission_authority", func(d map[string]interface{}) { d["fluidvirtAdmissionAuthority"] = true }},
		{"inventory_runtime_executor", func(d map[string]interface{}) { d["inventoryRuntimeExecutor"] = true }},
		{"forbidden_metrics_gauge", func(d map[string]interface{}) { d["metricsAdmissionResult"].(map[string]interface{})["generalProductionAutoEnabledGauge"] = float64(1) }},
		{"next_state_policy_auto", func(d map[string]interface{}) { d["productionInstallEligibility"].(map[string]interface{})["nextAllowedState"] = "production_auto_with_policy" }},
		{"ha_without_evidence", func(d map[string]interface{}) { d["leaderElectionAdmissionResult"].(map[string]interface{})["haProductionProven"] = true }},
		{"missing_evidence_refs", func(d map[string]interface{}) { d["rbacAdmissionResult"].(map[string]interface{})["evidenceRefs"] = []interface{}{} }},
		{"missing_claim_boundary", func(d map[string]interface{}) { d["rbacAdmissionResult"].(map[string]interface{})["claimBoundary"] = "" }},
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
