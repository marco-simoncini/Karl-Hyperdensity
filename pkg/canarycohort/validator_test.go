package canarycohort

import "testing"

func validSurface() map[string]interface{} {
	return referenceSurfaceDoc()
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/production-canary-cohort-expansion-reference.json"); err != nil {
		t.Fatal(err)
	}
}

func TestValidatorNegativeCases(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]interface{})
	}{
		{"cohort_disabled", func(d map[string]interface{}) { d["productionCanaryCohortExpansionEnabled"] = false }},
		{"executed_no_stable_window", func(d map[string]interface{}) { d["stableObservationWindowObserved"] = false }},
		{"executed_no_prior_canary", func(d map[string]interface{}) { d["priorCanaryMovementEvidenceObserved"] = false }},
		{"executed_no_allowlist", func(d map[string]interface{}) { d["allowlistGatePassed"] = false }},
		{"executed_no_blast_radius", func(d map[string]interface{}) { d["blastRadiusGatePassed"] = false }},
		{"executed_no_slo", func(d map[string]interface{}) { d["sloGatePassed"] = false }},
		{"executed_no_rollback", func(d map[string]interface{}) { d["rollbackGatePassed"] = false }},
		{"executed_kill_switch", func(d map[string]interface{}) { d["killSwitchGatePassed"] = false }},
		{"executed_circuit_breaker", func(d map[string]interface{}) { d["circuitBreakerGatePassed"] = false }},
		{"executed_rate_limit", func(d map[string]interface{}) { d["rateLimitGatePassed"] = false }},
		{"executed_failure_budget", func(d map[string]interface{}) { d["failureBudgetPassed"] = false }},
		{"failure_budget_exceeded", func(d map[string]interface{}) {
			d["failureBudget"].(map[string]interface{})["failureBudgetExceeded"] = true
		}},
		{"cohort_too_small", func(d map[string]interface{}) { d["cohortMovementCount"] = float64(1) }},
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
		{"prod_auto_policy_enabled", func(d map[string]interface{}) { d["productionAutoWithPolicyEnabled"] = true }},
		{"prod_auto_policy_activated", func(d map[string]interface{}) { d["productionAutoWithPolicyActivated"] = true }},
		{"broad_mutation", func(d map[string]interface{}) { d["broadProductionMutationExecuted"] = true }},
		{"fluidvirt_missing", func(d map[string]interface{}) {
			d["fluidvirtInvocations"].([]interface{})[0].(map[string]interface{})["actuator"] = "Dashboard"
		}},
		{"mutation_mismatch", func(d map[string]interface{}) {
			d["mutationObservations"].([]interface{})[0].(map[string]interface{})["mutationMatchesPlan"] = false
		}},
		{"reboot_used", func(d map[string]interface{}) {
			d["mutationObservations"].([]interface{})[0].(map[string]interface{})["rebootUsed"] = true
		}},
		{"aggregate_no_evidence", func(d map[string]interface{}) {
			d["realizedCompressionAggregation"].(map[string]interface{})["movementEvidencePresent"] = false
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
			d["canaryCohortCloseout"].(map[string]interface{})["promotedToGeneralProduction"] = true
		}},
		{"graduation_auto_enabled", func(d map[string]interface{}) {
			d["guardedPolicyGraduationEvidence"].(map[string]interface{})["graduationDecision"] = "production_auto_with_policy_enabled"
		}},
		{"readiness_incomplete", func(d map[string]interface{}) {
			d["productionAutoWithPolicyReadinessCandidate"] = true
			d["guardedPolicyGraduationEvidence"].(map[string]interface{})["graduationEvidenceComplete"] = false
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
