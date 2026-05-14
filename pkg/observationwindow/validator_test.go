package observationwindow

import "testing"

func validSurface() map[string]interface{} {
	return referenceSurfaceDoc()
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/production-observation-window-multitick-compression-evidence-reference.json"); err != nil {
		t.Fatal(err)
	}
}

func TestValidatorNegativeCases(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]interface{})
	}{
		{"window_disabled", func(d map[string]interface{}) { d["productionObservationWindowEnabled"] = false }},
		{"not_observation", func(d map[string]interface{}) { d["productionObservationMode"] = false }},
		{"not_no_apply", func(d map[string]interface{}) { d["noApplyMode"] = false }},
		{"not_multi_tick", func(d map[string]interface{}) { d["multiTickWindowObserved"] = false }},
		{"passed_too_few_ticks", func(d map[string]interface{}) { d["tickCount"] = float64(2) }},
		{"passed_no_real_ticks", func(d map[string]interface{}) { d["realInputTicksObserved"] = false }},
		{"passed_no_input_stability", func(d map[string]interface{}) { d["inputStabilityAnalyzed"] = false }},
		{"passed_no_donor_stability", func(d map[string]interface{}) { d["donorStabilityAnalyzed"] = false }},
		{"passed_no_pressure_persistence", func(d map[string]interface{}) { d["receiverPressurePersistenceAnalyzed"] = false }},
		{"passed_no_idle_persistence", func(d map[string]interface{}) { d["idleOpportunityPersistenceAnalyzed"] = false }},
		{"passed_no_blocker_decay", func(d map[string]interface{}) { d["blockerDecayAnalyzed"] = false }},
		{"passed_no_stale_decay", func(d map[string]interface{}) { d["staleInputDecayAnalyzed"] = false }},
		{"passed_no_action_refresh", func(d map[string]interface{}) { d["actionSlateRefreshed"] = false }},
		{"passed_no_future_refresh", func(d map[string]interface{}) { d["resourceFuturesRefreshed"] = false }},
		{"passed_no_compression_trend", func(d map[string]interface{}) { d["projectedCompressionTrendComputed"] = false }},
		{"passed_no_value_trend", func(d map[string]interface{}) { d["projectedValueTrendComputed"] = false }},
		{"passed_no_realized_separation", func(d map[string]interface{}) { d["realizedValueSeparated"] = false }},
		{"passed_no_safety_window", func(d map[string]interface{}) { d["noApplySafetyWindowVerified"] = false }},
		{"passed_no_dashboard", func(d map[string]interface{}) { d["dashboardProjectionUpdated"] = false }},
		{"tick_production_movement", func(d map[string]interface{}) {
			d["tickSequence"].(map[string]interface{})["ticks"].([]interface{})[0].(map[string]interface{})["productionMovementExecuted"] = true
		}},
		{"tick_selected_execution", func(d map[string]interface{}) {
			d["tickSequence"].(map[string]interface{})["ticks"].([]interface{})[0].(map[string]interface{})["selectedForExecutionCount"] = float64(1)
		}},
		{"action_refresh_selected", func(d map[string]interface{}) {
			d["actionSlateRefreshEvidence"].(map[string]interface{})["selectedForExecutionCount"] = float64(1)
		}},
		{"production_movement", func(d map[string]interface{}) { d["productionMovementExecuted"] = true }},
		{"broad_mutation", func(d map[string]interface{}) { d["broadProductionMutationExecuted"] = true }},
		{"realized_movement_no_evidence", func(d map[string]interface{}) { d["realizedMovementCount"] = float64(1) }},
		{"realized_value_no_evidence", func(d map[string]interface{}) { d["realizedMovedIdleValue"] = float64(100) }},
		{"realized_compression_no_evidence", func(d map[string]interface{}) { d["realizedCompressionDelta"] = 0.01 }},
		{"general_auto", func(d map[string]interface{}) { d["generalProductionAutoAllowed"] = true }},
		{"prod_auto_policy", func(d map[string]interface{}) { d["productionAutoWithPolicy"] = true }},
		{"universal_allowed", func(d map[string]interface{}) { d["universalGuaranteedSavingsAllowed"] = true }},
		{"universal_claimed", func(d map[string]interface{}) { d["universalGuaranteedSavingsClaimed"] = true }},
		{"projected_compression_realized", func(d map[string]interface{}) { d["projectedCompressionCountedAsRealized"] = true }},
		{"projected_value_realized", func(d map[string]interface{}) { d["projectedValueCountedAsRealized"] = true }},
		{"estimated_moved", func(d map[string]interface{}) { d["estimatedIdleCountedAsMoved"] = true }},
		{"synthetic_prod", func(d map[string]interface{}) { d["syntheticFleetCountedAsProduction"] = true }},
		{"reference_prod", func(d map[string]interface{}) { d["referenceFleetCountedAsProduction"] = true }},
		{"dashboard_executor", func(d map[string]interface{}) { d["dashboardExecutor"] = true }},
		{"fluidvirt_policy", func(d map[string]interface{}) { d["fluidvirtPolicyAuthority"] = true }},
		{"fluidvirt_window_authority", func(d map[string]interface{}) { d["fluidvirtObservationWindowAuthority"] = true }},
		{"inventory_executor", func(d map[string]interface{}) { d["inventoryRuntimeExecutor"] = true }},
		{"forbidden_gauge", func(d map[string]interface{}) {
			d["observationMetrics"].(map[string]interface{})["generalProductionAutoEnabledGauge"] = float64(1)
		}},
		{"trend_projected_realized", func(d map[string]interface{}) {
			d["projectedCompressionTrend"].(map[string]interface{})["projectedCompressionCountedAsRealized"] = true
		}},
		{"value_trend_projected_realized", func(d map[string]interface{}) {
			d["projectedValueTrend"].(map[string]interface{})["projectedValueCountedAsRealized"] = true
		}},
		{"unsafe_safety_window", func(d map[string]interface{}) {
			d["noApplySafetyWindow"].(map[string]interface{})["productionMovementExecuted"] = true
		}},
		{"passed_with_blockers", func(d map[string]interface{}) { d["blockers"] = []interface{}{map[string]interface{}{"blockerCode": "x"}} }},
		{"missing_evidence", func(d map[string]interface{}) {
			d["donorStabilityAnalysis"].(map[string]interface{})["evidenceRefs"] = []interface{}{}
		}},
		{"missing_claim", func(d map[string]interface{}) {
			d["donorStabilityAnalysis"].(map[string]interface{})["claimBoundary"] = ""
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
