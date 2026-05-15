package canarymovement

import "fmt"

func validateCanarySections(doc map[string]interface{}) error {
	for _, key := range []string{
		"candidateSelection", "allowlistGate", "blastRadiusGate", "executionLease",
		"applyRequest", "fluidvirtInvocation", "mutationObservation", "postVerifyResult",
		"rollbackWindow", "realizedCompressionRecord", "realizedValueRecord",
		"projectedRealizedSeparation", "canaryCloseout",
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
	if boolv(doc["productionCanaryMovementExecuted"]) {
		for _, gate := range []string{
			"stableObservationWindowObserved", "allowlistGatePassed", "blastRadiusGatePassed",
			"sloGatePassed", "rollbackGatePassed", "killSwitchGatePassed",
			"circuitBreakerGatePassed", "rateLimitGatePassed",
		} {
			if !boolv(doc[gate]) {
				return fmt.Errorf("canary movement executed with failed gate: %s", gate)
			}
		}
		inv := doc["fluidvirtInvocation"].(map[string]interface{})
		if strv(inv["actuator"]) != "FluidVirt" {
			return fmt.Errorf("FluidVirt invocation required")
		}
		mut := doc["mutationObservation"].(map[string]interface{})
		if !boolv(mut["mutationObserved"]) || !boolv(mut["mutationMatchesPlan"]) || !boolv(mut["identityPreserved"]) {
			return fmt.Errorf("mutation observation invalid")
		}
		for _, key := range []string{"rebootUsed", "recreateUsed", "rolloutUsed", "migrationUsed"} {
			if boolv(mut[key]) {
				return fmt.Errorf("%s must be false for canary movement", key)
			}
		}
		pv := doc["postVerifyResult"].(map[string]interface{})
		if strv(pv["verifyStatus"]) != "passed" || !boolv(pv["mutationKept"]) || boolv(pv["mutationRolledBack"]) {
			return fmt.Errorf("post-verify must pass with mutation kept")
		}
		if strv(pv["sloGuardStatus"]) != "passed" {
			return fmt.Errorf("SLO guard must pass")
		}
	}
	closeout := doc["canaryCloseout"].(map[string]interface{})
	if boolv(closeout["promotedToGeneralProduction"]) {
		return fmt.Errorf("promotedToGeneralProduction must be false")
	}
	if strv(closeout["nextAllowedState"]) == "production_auto_with_policy" {
		return fmt.Errorf("nextAllowedState production_auto_with_policy forbidden")
	}
	sep := doc["projectedRealizedSeparation"].(map[string]interface{})
	if boolv(sep["projectedCompressionCountedAsRealized"]) || boolv(sep["projectedValueCountedAsRealized"]) || boolv(sep["estimatedIdleCountedAsMoved"]) {
		return fmt.Errorf("projected/estimated must not count as realized")
	}
	if num(doc["realizedMovementCount"]) > 0 {
		rcr := doc["realizedCompressionRecord"].(map[string]interface{})
		if !boolv(rcr["movementEvidencePresent"]) {
			return fmt.Errorf("realizedMovementCount > 0 without movement evidence")
		}
	}
	if num(doc["realizedMovedIdleValue"]) > 0 && !boolv(doc["realizedCompressionRecord"].(map[string]interface{})["movementEvidencePresent"]) {
		return fmt.Errorf("realizedMovedIdleValue > 0 without movement evidence")
	}
	apply := doc["applyRequest"].(map[string]interface{})
	if strv(apply["actuator"]) != "FluidVirt" || !boolv(apply["productionCanaryScope"]) || boolv(apply["generalProductionScope"]) {
		return fmt.Errorf("apply request must be FluidVirt canary-scoped")
	}
	return nil
}

func evaluateCanaryExpansionPass(doc map[string]interface{}) error {
	if !boolv(doc["productionCanaryMovementExpansionEnabled"]) || !boolv(doc["stableObservationWindowRequired"]) {
		return fmt.Errorf("canary expansion flags not enabled")
	}
	if !boolv(doc["productionCanaryMovementExecuted"]) {
		return nil
	}
	if !boolv(doc["productionCanaryScope"]) || !boolv(doc["productionCanaryMovementAllowed"]) {
		return fmt.Errorf("production canary scope required")
	}
	for _, key := range []string{
		"candidateSelectionPassed", "executionLeaseCreated", "fluidvirtInvocationRecorded",
		"mutationObserved", "postVerifyPassed", "rollbackWindowOpen",
		"realizedCompressionRecorded", "realizedValueRecorded", "projectedRealizedSeparated",
		"canaryCloseoutRecorded",
	} {
		if !boolv(doc[key]) {
			return fmt.Errorf("canary movement executed with failed check: %s", key)
		}
	}
	if num(doc["realizedMovementCount"]) <= 0 || num(doc["realizedMovedIdleValue"]) <= 0 {
		return fmt.Errorf("successful canary movement requires realized counts > 0")
	}
	return nil
}
