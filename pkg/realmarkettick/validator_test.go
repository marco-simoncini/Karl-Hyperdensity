package realmarkettick

import "testing"

func validSurface() map[string]interface{} {
	return referenceSurfaceDoc()
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/runtime-market-tick-real-shell-inputs-no-apply-reference.json"); err != nil {
		t.Fatal(err)
	}
}

func TestValidatorNegativeCases(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]interface{})
	}{
		{"tick_disabled", func(d map[string]interface{}) { d["realInputMarketTickEnabled"] = false }},
		{"not_observation", func(d map[string]interface{}) { d["productionObservationMode"] = false }},
		{"not_no_apply", func(d map[string]interface{}) { d["noApplyMode"] = false }},
		{"not_executed", func(d map[string]interface{}) { d["realInputMarketTickExecuted"] = false }},
		{"passed_no_shell", func(d map[string]interface{}) { d["realShellInputsObserved"] = false }},
		{"passed_no_idle", func(d map[string]interface{}) { d["realIdleSignalsObserved"] = false }},
		{"passed_no_pressure", func(d map[string]interface{}) { d["realPressureSignalsObserved"] = false }},
		{"passed_no_slo", func(d map[string]interface{}) { d["realSloReadinessObserved"] = false }},
		{"passed_no_rollback", func(d map[string]interface{}) { d["realRollbackReadinessObserved"] = false }},
		{"passed_no_freshness", func(d map[string]interface{}) { d["inputFreshnessValidated"] = false }},
		{"passed_no_invalidation", func(d map[string]interface{}) { d["staleInputsInvalidated"] = false }},
		{"passed_no_donor_index", func(d map[string]interface{}) { d["realInputDonorIndexGenerated"] = false }},
		{"passed_no_receiver_index", func(d map[string]interface{}) { d["realInputReceiverIndexGenerated"] = false }},
		{"passed_no_pairing", func(d map[string]interface{}) { d["boundedPairingVerified"] = false }},
		{"passed_no_actions", func(d map[string]interface{}) { d["noApplyActionSlateGenerated"] = false }},
		{"passed_no_futures", func(d map[string]interface{}) { d["noApplyFuturesGenerated"] = false }},
		{"passed_no_metrics", func(d map[string]interface{}) { d["observationMetricsEmitted"] = false }},
		{"passed_no_events", func(d map[string]interface{}) { d["observationEventsEmitted"] = false }},
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
		{"action_selected", func(d map[string]interface{}) {
			d["generatedActions"].([]interface{})[0].(map[string]interface{})["selectedForExecution"] = true
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
		{"fluidvirt_tick", func(d map[string]interface{}) { d["fluidvirtTickAuthority"] = true }},
		{"inventory_executor", func(d map[string]interface{}) { d["inventoryRuntimeExecutor"] = true }},
		{"forbidden_gauge", func(d map[string]interface{}) {
			d["observationMetrics"].(map[string]interface{})["generalProductionAutoEnabledGauge"] = float64(1)
		}},
		{"passed_with_blockers", func(d map[string]interface{}) { d["blockers"] = []interface{}{map[string]interface{}{"blockerCode": "x"}} }},
		{"missing_evidence", func(d map[string]interface{}) { d["donorIndex"].(map[string]interface{})["evidenceRefs"] = []interface{}{} }},
		{"missing_claim", func(d map[string]interface{}) { d["donorIndex"].(map[string]interface{})["claimBoundary"] = "" }},
		{"unsafe_safety_boundary", func(d map[string]interface{}) {
			d["noApplySafetyBoundary"].(map[string]interface{})["productionMovementExecuted"] = true
		}},
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
