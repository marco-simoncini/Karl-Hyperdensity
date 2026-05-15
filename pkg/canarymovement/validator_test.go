package canarymovement

import "testing"

func validSurface() map[string]interface{} {
	return referenceSurfaceDoc()
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/production-canary-movement-expansion-reference.json"); err != nil {
		t.Fatal(err)
	}
}

func TestValidatorNegativeCases(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]interface{})
	}{
		{"expansion_disabled", func(d map[string]interface{}) { d["productionCanaryMovementExpansionEnabled"] = false }},
		{"stable_window_not_required", func(d map[string]interface{}) { d["stableObservationWindowRequired"] = false }},
		{"executed_no_stable_window", func(d map[string]interface{}) { d["stableObservationWindowObserved"] = false }},
		{"executed_no_allowlist", func(d map[string]interface{}) { d["allowlistGatePassed"] = false }},
		{"executed_no_blast_radius", func(d map[string]interface{}) { d["blastRadiusGatePassed"] = false }},
		{"executed_no_slo", func(d map[string]interface{}) { d["sloGatePassed"] = false }},
		{"executed_no_rollback", func(d map[string]interface{}) { d["rollbackGatePassed"] = false }},
		{"executed_kill_switch", func(d map[string]interface{}) { d["killSwitchGatePassed"] = false }},
		{"executed_circuit_breaker", func(d map[string]interface{}) { d["circuitBreakerGatePassed"] = false }},
		{"executed_rate_limit", func(d map[string]interface{}) { d["rateLimitGatePassed"] = false }},
		{"selected_not_canary", func(d map[string]interface{}) {
			d["candidateSelection"].(map[string]interface{})["candidates"].([]interface{})[0].(map[string]interface{})["productionCanaryScope"] = false
		}},
		{"selected_general_scope", func(d map[string]interface{}) {
			d["candidateSelection"].(map[string]interface{})["candidates"].([]interface{})[0].(map[string]interface{})["generalProductionScope"] = true
		}},
		{"selected_not_allowlisted", func(d map[string]interface{}) {
			d["candidateSelection"].(map[string]interface{})["candidates"].([]interface{})[0].(map[string]interface{})["allowlisted"] = false
		}},
		{"general_auto", func(d map[string]interface{}) { d["generalProductionAutoAllowed"] = true }},
		{"prod_auto_policy", func(d map[string]interface{}) { d["productionAutoWithPolicy"] = true }},
		{"broad_mutation", func(d map[string]interface{}) { d["broadProductionMutationExecuted"] = true }},
		{"fluidvirt_missing", func(d map[string]interface{}) {
			d["fluidvirtInvocation"].(map[string]interface{})["actuator"] = "Dashboard"
		}},
		{"mutation_not_observed", func(d map[string]interface{}) {
			d["mutationObservation"].(map[string]interface{})["mutationObserved"] = false
		}},
		{"mutation_mismatch", func(d map[string]interface{}) {
			d["mutationObservation"].(map[string]interface{})["mutationMatchesPlan"] = false
		}},
		{"reboot_used", func(d map[string]interface{}) { d["mutationObservation"].(map[string]interface{})["rebootUsed"] = true }},
		{"recreate_used", func(d map[string]interface{}) { d["mutationObservation"].(map[string]interface{})["recreateUsed"] = true }},
		{"rollout_used", func(d map[string]interface{}) { d["mutationObservation"].(map[string]interface{})["rolloutUsed"] = true }},
		{"migration_used", func(d map[string]interface{}) { d["mutationObservation"].(map[string]interface{})["migrationUsed"] = true }},
		{"post_verify_slo_fail", func(d map[string]interface{}) {
			d["postVerifyResult"].(map[string]interface{})["sloGuardStatus"] = "failed"
		}},
		{"realized_no_evidence", func(d map[string]interface{}) {
			d["realizedCompressionRecord"].(map[string]interface{})["movementEvidencePresent"] = false
		}},
		{"projected_compression_realized", func(d map[string]interface{}) { d["projectedCompressionCountedAsRealized"] = true }},
		{"projected_value_realized", func(d map[string]interface{}) { d["projectedValueCountedAsRealized"] = true }},
		{"estimated_moved", func(d map[string]interface{}) { d["estimatedIdleCountedAsMoved"] = true }},
		{"synthetic_prod", func(d map[string]interface{}) { d["syntheticFleetCountedAsProduction"] = true }},
		{"reference_prod", func(d map[string]interface{}) { d["referenceFleetCountedAsProduction"] = true }},
		{"dashboard_executor", func(d map[string]interface{}) { d["dashboardExecutor"] = true }},
		{"fluidvirt_policy", func(d map[string]interface{}) { d["fluidvirtPolicyAuthority"] = true }},
		{"fluidvirt_controller", func(d map[string]interface{}) { d["fluidvirtControllerAuthority"] = true }},
		{"inventory_executor", func(d map[string]interface{}) { d["inventoryRuntimeExecutor"] = true }},
		{"promoted_general", func(d map[string]interface{}) {
			d["canaryCloseout"].(map[string]interface{})["promotedToGeneralProduction"] = true
		}},
		{"next_auto_policy", func(d map[string]interface{}) {
			d["canaryCloseout"].(map[string]interface{})["nextAllowedState"] = "production_auto_with_policy"
		}},
		{"passed_with_blockers", func(d map[string]interface{}) { d["blockers"] = []interface{}{map[string]interface{}{"blockerCode": "x"}} }},
		{"missing_evidence", func(d map[string]interface{}) {
			d["allowlistGate"].(map[string]interface{})["evidenceRefs"] = []interface{}{}
		}},
		{"missing_claim", func(d map[string]interface{}) {
			d["allowlistGate"].(map[string]interface{})["claimBoundary"] = ""
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
