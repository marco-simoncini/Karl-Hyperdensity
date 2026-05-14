package smoketick

import "fmt"

func validateTickSections(doc map[string]interface{}) error {
	for _, key := range []string{
		"smokeTickRequest", "smokeTickResult", "stateAccess", "marketSnapshot",
		"indexRefresh", "pairingWindow", "tickMetrics", "dashboardProjection", "safetyBoundary", "healthDecision",
	} {
		m, _ := doc[key].(map[string]interface{})
		if m == nil {
			return fmt.Errorf("%s required", key)
		}
		if strv(m["claimBoundary"]) == "" || !hasNonEmptyStringList(m["evidenceRefs"]) {
			return fmt.Errorf("%s missing evidenceRefs or claimBoundary", key)
		}
	}

	req := doc["smokeTickRequest"].(map[string]interface{})
	if boolv(req["productionMovementAllowed"]) {
		return fmt.Errorf("productionMovementAllowed must be false in Sprint 18")
	}
	if boolv(req["generalProductionAutoAllowed"]) || boolv(req["productionAutoWithPolicy"]) {
		return fmt.Errorf("forbidden auto flags in smoke tick request")
	}

	pw := doc["pairingWindow"].(map[string]interface{})
	if !boolv(pw["noFullNxNPairing"]) {
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

	metrics := doc["tickMetrics"].(map[string]interface{})
	if num(metrics["generalProductionAutoEnabledGauge"]) > 0 || num(metrics["productionAutoWithPolicyEnabledGauge"]) > 0 {
		return fmt.Errorf("forbidden auto gauge present")
	}

	proj := doc["dashboardProjection"].(map[string]interface{})
	if !boolv(proj["projectionReadOnly"]) {
		return fmt.Errorf("dashboard projection must be read-only")
	}

	safety := doc["safetyBoundary"].(map[string]interface{})
	for _, key := range []string{
		"productionMovementExecuted", "broadProductionMutationExecuted",
		"generalProductionAutoAllowed", "productionAutoWithPolicy",
		"dashboardExecutor", "fluidvirtReconcilerAuthority",
		"rawRuntimeControlsExposed", "clusterAdmin", "podsExecAllowed", "nodesWriteAllowed",
	} {
		if boolv(safety[key]) {
			return fmt.Errorf("unsafe safety boundary: %s", key)
		}
	}
	return validateGeneratedActionsFutures(doc)
}

func validateGeneratedActionsFutures(doc map[string]interface{}) error {
	for _, item := range sliceMap(doc["generatedActions"]) {
		if boolv(item["generalProductionAutoAllowed"]) || boolv(item["productionAutoWithPolicy"]) {
			return fmt.Errorf("forbidden auto in generated action")
		}
		scope := strv(item["executionScopeRecommendation"])
		if scope == "general_production_auto" || scope == "production_auto_with_policy" {
			return fmt.Errorf("forbidden action scope: %s", scope)
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

func evaluateSmokeTickPass(doc map[string]interface{}) error {
	if !boolv(doc["installedControllerSmokeTickEnabled"]) || !boolv(doc["smokeTickRequested"]) || !boolv(doc["smokeTickExecuted"]) {
		return fmt.Errorf("smoke tick not enabled/requested/executed")
	}
	if !boolv(doc["smokeTickPassed"]) {
		return nil
	}
	for _, key := range []string{
		"leaderHeldTickObserved", "durableStateReadVerified", "durableStateWriteVerified",
		"marketSnapshotCollected", "indicesRefreshed", "boundedPairingVerified",
		"actionSlateGenerated", "resourceFuturesGenerated", "tickMetricsEmitted",
		"tickEventsEmitted", "dashboardProjectionUpdated", "noFullNxNPairing",
	} {
		if !boolv(doc[key]) {
			return fmt.Errorf("smokeTickPassed with failed check: %s", key)
		}
	}
	if boolv(doc["productionMovementExecuted"]) || boolv(doc["broadProductionMutationExecuted"]) {
		return fmt.Errorf("smokeTickPassed with production movement")
	}
	return nil
}
