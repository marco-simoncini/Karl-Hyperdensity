package realmarkettick

import "fmt"

func validateInputSections(doc map[string]interface{}) error {
	for _, key := range []string{
		"realShellInputSnapshot", "shellInputSourceMap", "observedIdleSignals",
		"observedPressureSignals", "observedSloReadiness", "observedRollbackReadiness",
		"observedRiskSignals", "inputFreshnessValidation", "donorIndex", "receiverIndex",
		"pairingWindow", "observationMetrics", "noApplySafetyBoundary",
	} {
		m, _ := doc[key].(map[string]interface{})
		if m == nil {
			return fmt.Errorf("%s required", key)
		}
		if strv(m["claimBoundary"]) == "" || !hasNonEmptyStringList(m["evidenceRefs"]) {
			return fmt.Errorf("%s missing evidenceRefs or claimBoundary", key)
		}
	}
	fresh := doc["inputFreshnessValidation"].(map[string]interface{})
	if !boolv(fresh["inputFreshnessValidated"]) && boolv(doc["realInputMarketTickPassed"]) {
		return fmt.Errorf("tick passed without freshness validation")
	}
	if !boolv(doc["staleInputsInvalidated"]) && boolv(doc["realInputMarketTickPassed"]) {
		return fmt.Errorf("tick passed without stale input invalidation")
	}
	return validatePairingAndActions(doc)
}

func validatePairingAndActions(doc map[string]interface{}) error {
	pw := doc["pairingWindow"].(map[string]interface{})
	if !boolv(pw["noFullNxNPairing"]) || !boolv(doc["noFullNxNPairing"]) {
		return fmt.Errorf("noFullNxNPairing must be true")
	}
	full := num(pw["fullPairSpace"])
	eval := num(pw["evaluatedPairCount"])
	avoided := num(pw["avoidedPairCount"])
	topD := num(doc["topKDonors"])
	topR := num(doc["topKReceivers"])
	if eval > topD*topR {
		return fmt.Errorf("evaluatedPairCount exceeds top-K window")
	}
	if avoided != full-eval {
		return fmt.Errorf("avoidedPairCount mismatch")
	}
	metrics := doc["observationMetrics"].(map[string]interface{})
	if num(metrics["generalProductionAutoEnabledGauge"]) > 0 || num(metrics["productionAutoWithPolicyEnabledGauge"]) > 0 {
		return fmt.Errorf("forbidden auto gauge present")
	}
	safety := doc["noApplySafetyBoundary"].(map[string]interface{})
	for _, key := range []string{
		"productionMovementExecuted", "broadProductionMutationExecuted",
		"generalProductionAutoAllowed", "productionAutoWithPolicy", "selectedProductionExecution",
		"dashboardExecutor", "fluidvirtTickAuthority", "rawRuntimeControlsExposed", "inventoryRuntimeExecutor",
	} {
		if boolv(safety[key]) {
			return fmt.Errorf("unsafe safety boundary: %s", key)
		}
	}
	for _, item := range sliceMap(doc["generatedActions"]) {
		if boolv(item["selectedForExecution"]) {
			return fmt.Errorf("action selectedForExecution in no-apply mode")
		}
		scope := strv(item["executionScopeRecommendation"])
		if scope == "general_production_auto" || scope == "production_auto_with_policy" {
			return fmt.Errorf("forbidden action scope: %s", scope)
		}
		if boolv(item["generalProductionAutoAllowed"]) || boolv(item["productionAutoWithPolicy"]) {
			return fmt.Errorf("forbidden auto in generated action")
		}
		if strv(item["claimBoundary"]) == "" || !hasNonEmptyStringList(item["evidenceRefs"]) {
			return fmt.Errorf("generated action missing evidenceRefs or claimBoundary")
		}
	}
	for _, item := range sliceMap(doc["generatedFutures"]) {
		scope := strv(item["executionScopeRecommendation"])
		if scope == "general_production_auto" || scope == "production_auto_with_policy" {
			return fmt.Errorf("forbidden future scope: %s", scope)
		}
		if strv(item["claimBoundary"]) == "" || !hasNonEmptyStringList(item["evidenceRefs"]) {
			return fmt.Errorf("generated future missing evidenceRefs or claimBoundary")
		}
	}
	return nil
}

func evaluateRealInputTickPass(doc map[string]interface{}) error {
	if !boolv(doc["realInputMarketTickEnabled"]) || !boolv(doc["productionObservationMode"]) || !boolv(doc["noApplyMode"]) {
		return fmt.Errorf("real input tick flags not enabled")
	}
	if !boolv(doc["realInputMarketTickExecuted"]) {
		return fmt.Errorf("realInputMarketTickExecuted must be true")
	}
	if !boolv(doc["realInputMarketTickPassed"]) {
		return nil
	}
	for _, key := range []string{
		"realShellInputsObserved", "realIdleSignalsObserved", "realPressureSignalsObserved",
		"realSloReadinessObserved", "realRollbackReadinessObserved", "inputFreshnessValidated",
		"staleInputsInvalidated", "realInputDonorIndexGenerated", "realInputReceiverIndexGenerated",
		"boundedPairingVerified", "noApplyActionSlateGenerated", "noApplyFuturesGenerated",
		"observationMetricsEmitted", "observationEventsEmitted", "dashboardProjectionUpdated",
	} {
		if !boolv(doc[key]) {
			return fmt.Errorf("realInputMarketTickPassed with failed check: %s", key)
		}
	}
	if boolv(doc["productionMovementExecuted"]) || boolv(doc["broadProductionMutationExecuted"]) {
		return fmt.Errorf("tick passed with production movement")
	}
	return nil
}
