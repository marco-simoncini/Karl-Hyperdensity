package observationwindow

import "fmt"

func validateWindowSections(doc map[string]interface{}) error {
	for _, key := range []string{
		"tickSequence", "inputStabilityAnalysis", "donorStabilityAnalysis",
		"receiverPressurePersistence", "idleOpportunityPersistence",
		"blockerDecayAnalysis", "staleInputDecayAnalysis",
		"actionSlateRefreshEvidence", "resourceFutureRefreshEvidence",
		"projectedCompressionTrend", "projectedValueTrend",
		"realizedValueSeparation", "noApplySafetyWindow", "observationMetrics",
	} {
		m, _ := doc[key].(map[string]interface{})
		if m == nil {
			return fmt.Errorf("%s required", key)
		}
		if strv(m["claimBoundary"]) == "" || !hasNonEmptyStringList(m["evidenceRefs"]) {
			return fmt.Errorf("%s missing evidenceRefs or claimBoundary", key)
		}
	}
	seq := doc["tickSequence"].(map[string]interface{})
	ticks := sliceMap(seq["ticks"])
	if num(doc["tickCount"]) < 3 && boolv(doc["observationWindowPassed"]) {
		return fmt.Errorf("observationWindowPassed requires tickCount >= 3")
	}
	if len(ticks) < 3 && boolv(doc["observationWindowPassed"]) {
		return fmt.Errorf("observationWindowPassed requires at least 3 tick records")
	}
	if !boolv(seq["allTicksNoApply"]) && boolv(doc["observationWindowPassed"]) {
		return fmt.Errorf("not all ticks are no-apply")
	}
	for _, tick := range ticks {
		if boolv(tick["productionMovementExecuted"]) {
			return fmt.Errorf("tick %s has production movement", strv(tick["tickId"]))
		}
		if num(tick["selectedForExecutionCount"]) > 0 {
			return fmt.Errorf("tick %s has selectedForExecutionCount > 0", strv(tick["tickId"]))
		}
	}
	refresh := doc["actionSlateRefreshEvidence"].(map[string]interface{})
	if num(refresh["selectedForExecutionCount"]) > 0 || !boolv(refresh["allActionsNoApply"]) {
		return fmt.Errorf("action slate refresh must be no-apply")
	}
	pct := doc["projectedCompressionTrend"].(map[string]interface{})
	if boolv(pct["projectedCompressionCountedAsRealized"]) {
		return fmt.Errorf("projectedCompressionCountedAsRealized must be false")
	}
	pvt := doc["projectedValueTrend"].(map[string]interface{})
	if boolv(pvt["projectedValueCountedAsRealized"]) {
		return fmt.Errorf("projectedValueCountedAsRealized must be false")
	}
	rvs := doc["realizedValueSeparation"].(map[string]interface{})
	if !boolv(rvs["projectedValueExcludedFromRealized"]) {
		return fmt.Errorf("projectedValueExcludedFromRealized must be true")
	}
	if num(doc["realizedMovementCount"]) > 0 && !boolv(rvs["movementEvidencePresent"]) {
		return fmt.Errorf("realizedMovementCount > 0 without movement evidence")
	}
	if num(doc["realizedMovedIdleValue"]) > 0 && !boolv(rvs["movementEvidencePresent"]) {
		return fmt.Errorf("realizedMovedIdleValue > 0 without movement evidence")
	}
	if num(doc["realizedCompressionDelta"]) > 0 && !boolv(rvs["movementEvidencePresent"]) {
		return fmt.Errorf("realizedCompressionDelta > 0 without movement evidence")
	}
	metrics := doc["observationMetrics"].(map[string]interface{})
	if num(metrics["generalProductionAutoEnabledGauge"]) > 0 || num(metrics["productionAutoWithPolicyEnabledGauge"]) > 0 {
		return fmt.Errorf("forbidden auto gauge present")
	}
	if num(metrics["realizedMovedIdleValueGauge"]) > 0 && !boolv(rvs["movementEvidencePresent"]) {
		return fmt.Errorf("realizedMovedIdleValueGauge > 0 without movement evidence")
	}
	safety := doc["noApplySafetyWindow"].(map[string]interface{})
	for _, key := range []string{
		"productionMovementExecuted", "broadProductionMutationExecuted",
		"generalProductionAutoAllowed", "productionAutoWithPolicy",
		"dashboardExecutor", "fluidvirtObservationWindowAuthority",
		"rawRuntimeControlsExposed", "inventoryRuntimeExecutor",
	} {
		if boolv(safety[key]) {
			return fmt.Errorf("unsafe safety window: %s", key)
		}
	}
	if num(safety["selectedForExecutionCount"]) > 0 {
		return fmt.Errorf("selectedForExecutionCount must be 0 in safety window")
	}
	return nil
}

func evaluateWindowPass(doc map[string]interface{}) error {
	if !boolv(doc["productionObservationWindowEnabled"]) || !boolv(doc["productionObservationMode"]) || !boolv(doc["noApplyMode"]) {
		return fmt.Errorf("observation window flags not enabled")
	}
	if !boolv(doc["multiTickWindowObserved"]) {
		return fmt.Errorf("multiTickWindowObserved must be true")
	}
	if !boolv(doc["observationWindowPassed"]) {
		return nil
	}
	if num(doc["tickCount"]) < 3 {
		return fmt.Errorf("observationWindowPassed with tickCount < 3")
	}
	for _, key := range []string{
		"realInputTicksObserved", "inputStabilityAnalyzed", "donorStabilityAnalyzed",
		"receiverPressurePersistenceAnalyzed", "idleOpportunityPersistenceAnalyzed",
		"blockerDecayAnalyzed", "staleInputDecayAnalyzed", "actionSlateRefreshed",
		"resourceFuturesRefreshed", "projectedCompressionTrendComputed",
		"projectedValueTrendComputed", "realizedValueSeparated",
		"noApplySafetyWindowVerified", "dashboardProjectionUpdated",
	} {
		if !boolv(doc[key]) {
			return fmt.Errorf("observationWindowPassed with failed check: %s", key)
		}
	}
	if boolv(doc["productionMovementExecuted"]) || boolv(doc["broadProductionMutationExecuted"]) {
		return fmt.Errorf("observationWindowPassed with production movement")
	}
	return nil
}
