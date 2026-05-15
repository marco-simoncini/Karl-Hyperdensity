package canarycohort

import "fmt"

func validateCohortSections(doc map[string]interface{}) error {
	for _, key := range []string{
		"candidateSelection", "allowlistGate", "blastRadiusGate", "sloGate", "rollbackGate",
		"failureBudget", "realizedCompressionAggregation", "realizedValueAggregation",
		"projectedRealizedSeparation", "guardedPolicyGraduationEvidence", "canaryCohortCloseout",
	} {
		m, _ := doc[key].(map[string]interface{})
		if m == nil {
			return fmt.Errorf("%s required", key)
		}
		if strv(m["claimBoundary"]) == "" || !hasNonEmptyStringList(m["evidenceRefs"]) {
			return fmt.Errorf("%s missing evidenceRefs or claimBoundary", key)
		}
	}
	sel := doc["candidateSelection"].(map[string]interface{})
	for _, c := range sliceMap(sel["candidates"]) {
		if boolv(c["selected"]) {
			if !boolv(c["productionCanaryScope"]) || boolv(c["generalProductionScope"]) || !boolv(c["allowlisted"]) {
				return fmt.Errorf("selected candidate must be canary-scoped and allowlisted")
			}
		}
	}
	if num(doc["cohortMovementCount"]) < 2 && boolv(doc["productionCanaryCohortExecuted"]) {
		return fmt.Errorf("cohortMovementCount must be >= 2")
	}
	fb := doc["failureBudget"].(map[string]interface{})
	if boolv(fb["failureBudgetExceeded"]) {
		return fmt.Errorf("failure budget exceeded")
	}
	if boolv(doc["productionCanaryCohortExecuted"]) {
		for _, gate := range []string{
			"stableObservationWindowObserved", "priorCanaryMovementEvidenceObserved",
			"allowlistGatePassed", "blastRadiusGatePassed", "sloGatePassed",
			"rollbackGatePassed", "killSwitchGatePassed", "circuitBreakerGatePassed",
			"rateLimitGatePassed", "failureBudgetPassed",
		} {
			if !boolv(doc[gate]) {
				return fmt.Errorf("cohort executed with failed gate: %s", gate)
			}
		}
	}
	for _, inv := range sliceMap(doc["fluidvirtInvocations"]) {
		if strv(inv["actuator"]) != "FluidVirt" {
			return fmt.Errorf("FluidVirt invocation required")
		}
	}
	for _, mut := range sliceMap(doc["mutationObservations"]) {
		if boolv(mut["mutationObserved"]) && !boolv(mut["mutationMatchesPlan"]) {
			return fmt.Errorf("mutation must match plan")
		}
		for _, key := range []string{"rebootUsed", "recreateUsed", "rolloutUsed", "migrationUsed"} {
			if boolv(mut[key]) {
				return fmt.Errorf("%s must be false", key)
			}
		}
	}
	closeout := doc["canaryCohortCloseout"].(map[string]interface{})
	if boolv(closeout["promotedToGeneralProduction"]) || boolv(closeout["productionAutoWithPolicyEnabled"]) {
		return fmt.Errorf("closeout must not promote or enable production_auto_with_policy")
	}
	return nil
}

func evaluateCohortPass(doc map[string]interface{}) error {
	if !boolv(doc["productionCanaryCohortExpansionEnabled"]) {
		return fmt.Errorf("productionCanaryCohortExpansionEnabled must be true")
	}
	if !boolv(doc["productionCanaryCohortExecuted"]) {
		return nil
	}
	if !boolv(doc["productionCanaryScope"]) {
		return fmt.Errorf("productionCanaryScope required")
	}
	for _, key := range []string{
		"executionLeasesCreated", "fluidvirtInvocationsRecorded", "mutationsObserved",
		"postVerifiesPassed", "rollbackWindowsOpen", "realizedCompressionAggregated",
		"realizedValueAggregated", "projectedRealizedSeparated",
		"policyGraduationEvidenceProduced", "canaryCohortCloseoutRecorded",
	} {
		if !boolv(doc[key]) {
			return fmt.Errorf("cohort executed with failed check: %s", key)
		}
	}
	if num(doc["realizedMovementCount"]) <= 0 || num(doc["aggregateRealizedMovedIdleValue"]) <= 0 {
		return fmt.Errorf("successful cohort requires aggregate realized counts > 0")
	}
	return nil
}
