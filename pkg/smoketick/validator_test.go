package smoketick

import "testing"

func validSurface() map[string]interface{} {
	return map[string]interface{}{
		"milestone": Milestone, "surfaceVersion": "v1", "smokeTickId": "tick-1",
		"controllerMode": "production_canary_only",
		"installedControllerSmokeTickEnabled": true, "smokeTickRequested": true, "smokeTickExecuted": true, "smokeTickPassed": true,
		"leaderHeldTickObserved": true, "durableStateReadVerified": true, "durableStateWriteVerified": true,
		"marketSnapshotCollected": true, "indicesRefreshed": true, "boundedPairingVerified": true,
		"actionSlateGenerated": true, "resourceFuturesGenerated": true, "tickMetricsEmitted": true,
		"tickEventsEmitted": true, "dashboardProjectionUpdated": true,
		"productionMovementExecuted": false, "broadProductionMutationExecuted": false, "noFullNxNPairing": true,
		"fullPairSpace": float64(100), "evaluatedPairCount": float64(25), "avoidedPairCount": float64(75),
		"topKDonors": float64(5), "topKReceivers": float64(5),
		"generatedActionCount": float64(3), "generatedFutureCount": float64(2),
		"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
		"universalGuaranteedSavingsAllowed": false, "universalGuaranteedSavingsClaimed": false,
		"estimatedIdleCountedAsMoved": false, "projectedCompressionCountedAsRealized": false,
		"syntheticFleetCountedAsProduction": false, "referenceFleetCountedAsProduction": false,
		"dashboardExecutor": false, "fluidvirtPolicyAuthority": false, "fluidvirtReconcilerAuthority": false,
		"inventoryRuntimeExecutor": false,
		"smokeTickRequest": map[string]interface{}{
			"smokeTickRequestId": "req-1", "productionMovementAllowed": false,
			"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
			"evidenceRefs": []interface{}{"req"}, "claimBoundary": "request",
		},
		"smokeTickResult": map[string]interface{}{
			"smokeTickResultId": "res-1", "status": "passed", "smokeTickPassed": true,
			"evidenceRefs": []interface{}{"res"}, "claimBoundary": "result",
		},
		"stateAccess": map[string]interface{}{
			"durableStateReadVerified": true, "durableStateWriteVerified": true,
			"evidenceRefs": []interface{}{"state"}, "claimBoundary": "state",
		},
		"marketSnapshot": map[string]interface{}{"marketSnapshotCollected": true, "evidenceRefs": []interface{}{"snap"}, "claimBoundary": "snapshot"},
		"indexRefresh": map[string]interface{}{"indicesRefreshed": true, "evidenceRefs": []interface{}{"idx"}, "claimBoundary": "index"},
		"pairingWindow": map[string]interface{}{
			"noFullNxNPairing": true, "fullPairSpace": float64(100), "evaluatedPairCount": float64(25), "avoidedPairCount": float64(75),
			"evidenceRefs": []interface{}{"pair"}, "claimBoundary": "pairing",
		},
		"generatedActions": []interface{}{
			map[string]interface{}{
				"actionId": "act-1", "executionScopeRecommendation": "operator_controlled",
				"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false, "selectedForExecution": false,
				"evidenceRefs": []interface{}{"act"}, "claimBoundary": "action",
			},
		},
		"generatedFutures": []interface{}{
			map[string]interface{}{
				"futureId": "fut-1", "executionScopeRecommendation": "production_canary_eligible",
				"evidenceRefs": []interface{}{"fut"}, "claimBoundary": "future",
			},
		},
		"tickMetrics": map[string]interface{}{
			"generalProductionAutoEnabledGauge": float64(0), "productionAutoWithPolicyEnabledGauge": float64(0),
			"evidenceRefs": []interface{}{"metrics"}, "claimBoundary": "metrics",
		},
		"tickEvents": []interface{}{
			map[string]interface{}{"eventId": "evt-1", "eventType": "smoke_tick_passed", "evidenceRefs": []interface{}{"evt"}, "claimBoundary": "event"},
		},
		"dashboardProjection": map[string]interface{}{
			"projectionUpdated": true, "projectionReadOnly": true,
			"evidenceRefs": []interface{}{"dash"}, "claimBoundary": "dash",
		},
		"safetyBoundary": map[string]interface{}{
			"productionMovementExecuted": false, "broadProductionMutationExecuted": false,
			"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
			"dashboardExecutor": false, "fluidvirtReconcilerAuthority": false,
			"rawRuntimeControlsExposed": false, "clusterAdmin": false, "podsExecAllowed": false, "nodesWriteAllowed": false,
			"evidenceRefs": []interface{}{"safety"}, "claimBoundary": "safety",
		},
		"healthDecision": map[string]interface{}{
			"decision": "passed", "smokeTickPassed": true,
			"evidenceRefs": []interface{}{"health"}, "claimBoundary": "health",
		},
		"auditEvents": []interface{}{
			map[string]interface{}{"auditEventId": "aud-1", "eventType": "smoke_tick_passed", "evidenceRefs": []interface{}{"aud"}, "claimBoundary": "audit"},
		},
		"blockers": []interface{}{},
		"claimBoundaries": []interface{}{"smoke tick generates candidates only"},
	}
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/controller-runtime-reconciliation-smoke-tick-reference.json"); err != nil {
		t.Fatal(err)
	}
}

func TestValidatorNegativeCases(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]interface{})
	}{
		{"tick_disabled", func(d map[string]interface{}) { d["installedControllerSmokeTickEnabled"] = false }},
		{"not_requested", func(d map[string]interface{}) { d["smokeTickRequested"] = false }},
		{"not_executed", func(d map[string]interface{}) { d["smokeTickExecuted"] = false }},
		{"passed_no_leader", func(d map[string]interface{}) { d["leaderHeldTickObserved"] = false }},
		{"passed_no_read", func(d map[string]interface{}) { d["durableStateReadVerified"] = false }},
		{"passed_no_write", func(d map[string]interface{}) { d["durableStateWriteVerified"] = false }},
		{"passed_no_snapshot", func(d map[string]interface{}) { d["marketSnapshotCollected"] = false }},
		{"passed_no_indices", func(d map[string]interface{}) { d["indicesRefreshed"] = false }},
		{"passed_no_pairing", func(d map[string]interface{}) { d["boundedPairingVerified"] = false }},
		{"passed_no_actions", func(d map[string]interface{}) { d["actionSlateGenerated"] = false }},
		{"passed_no_futures", func(d map[string]interface{}) { d["resourceFuturesGenerated"] = false }},
		{"passed_no_metrics", func(d map[string]interface{}) { d["tickMetricsEmitted"] = false }},
		{"passed_no_events", func(d map[string]interface{}) { d["tickEventsEmitted"] = false }},
		{"passed_no_dashboard", func(d map[string]interface{}) { d["dashboardProjectionUpdated"] = false }},
		{"no_full_nxn_false", func(d map[string]interface{}) {
			d["noFullNxNPairing"] = false
			d["pairingWindow"].(map[string]interface{})["noFullNxNPairing"] = false
		}},
		{"pair_count_exceeded", func(d map[string]interface{}) {
			d["pairingWindow"].(map[string]interface{})["evaluatedPairCount"] = float64(100)
		}},
		{"avoided_mismatch", func(d map[string]interface{}) {
			d["pairingWindow"].(map[string]interface{})["avoidedPairCount"] = float64(50)
		}},
		{"forbidden_action_auto", func(d map[string]interface{}) {
			d["generatedActions"].([]interface{})[0].(map[string]interface{})["generalProductionAutoAllowed"] = true
		}},
		{"forbidden_action_scope", func(d map[string]interface{}) {
			d["generatedActions"].([]interface{})[0].(map[string]interface{})["executionScopeRecommendation"] = "general_production_auto"
		}},
		{"forbidden_future_scope", func(d map[string]interface{}) {
			d["generatedFutures"].([]interface{})[0].(map[string]interface{})["executionScopeRecommendation"] = "production_auto_with_policy"
		}},
		{"production_movement", func(d map[string]interface{}) { d["productionMovementExecuted"] = true }},
		{"broad_mutation", func(d map[string]interface{}) { d["broadProductionMutationExecuted"] = true }},
		{"general_auto", func(d map[string]interface{}) { d["generalProductionAutoAllowed"] = true }},
		{"prod_auto_policy", func(d map[string]interface{}) { d["productionAutoWithPolicy"] = true }},
		{"universal_allowed", func(d map[string]interface{}) { d["universalGuaranteedSavingsAllowed"] = true }},
		{"universal_claimed", func(d map[string]interface{}) { d["universalGuaranteedSavingsClaimed"] = true }},
		{"projected_realized", func(d map[string]interface{}) { d["projectedCompressionCountedAsRealized"] = true }},
		{"estimated_moved", func(d map[string]interface{}) { d["estimatedIdleCountedAsMoved"] = true }},
		{"synthetic_prod", func(d map[string]interface{}) { d["syntheticFleetCountedAsProduction"] = true }},
		{"reference_prod", func(d map[string]interface{}) { d["referenceFleetCountedAsProduction"] = true }},
		{"dashboard_executor", func(d map[string]interface{}) { d["dashboardExecutor"] = true }},
		{"fluidvirt_policy", func(d map[string]interface{}) { d["fluidvirtPolicyAuthority"] = true }},
		{"fluidvirt_reconciler", func(d map[string]interface{}) { d["fluidvirtReconcilerAuthority"] = true }},
		{"inventory_executor", func(d map[string]interface{}) { d["inventoryRuntimeExecutor"] = true }},
		{"forbidden_gauge", func(d map[string]interface{}) { d["tickMetrics"].(map[string]interface{})["generalProductionAutoEnabledGauge"] = float64(1) }},
		{"dashboard_not_readonly", func(d map[string]interface{}) { d["dashboardProjection"].(map[string]interface{})["projectionReadOnly"] = false }},
		{"passed_with_blockers", func(d map[string]interface{}) { d["blockers"] = []interface{}{map[string]interface{}{"blockerCode": "x"}} }},
		{"missing_evidence", func(d map[string]interface{}) { d["stateAccess"].(map[string]interface{})["evidenceRefs"] = []interface{}{} }},
		{"missing_claim", func(d map[string]interface{}) { d["stateAccess"].(map[string]interface{})["claimBoundary"] = "" }},
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
