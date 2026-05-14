package smoketick

import "testing"

func TestValidateTickSectionsRejectsProductionMovement(t *testing.T) {
	doc := minimalTickDoc()
	doc["smokeTickRequest"].(map[string]interface{})["productionMovementAllowed"] = true
	if err := validateTickSections(doc); err == nil {
		t.Fatal("expected productionMovementAllowed rejection")
	}
}

func TestEvaluateSmokeTickPassRejectsLeaderMissing(t *testing.T) {
	doc := validSurface()
	doc["leaderHeldTickObserved"] = false
	if err := evaluateSmokeTickPass(doc); err == nil {
		t.Fatal("expected leader missing rejection")
	}
}

func minimalTickDoc() map[string]interface{} {
	return map[string]interface{}{
		"topKDonors": float64(5), "topKReceivers": float64(5),
		"smokeTickRequest": map[string]interface{}{
			"productionMovementAllowed": false, "generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
			"evidenceRefs": []interface{}{"req"}, "claimBoundary": "request",
		},
		"smokeTickResult": map[string]interface{}{"evidenceRefs": []interface{}{"res"}, "claimBoundary": "result"},
		"stateAccess": map[string]interface{}{"evidenceRefs": []interface{}{"state"}, "claimBoundary": "state"},
		"marketSnapshot": map[string]interface{}{"evidenceRefs": []interface{}{"snap"}, "claimBoundary": "snapshot"},
		"indexRefresh": map[string]interface{}{"evidenceRefs": []interface{}{"idx"}, "claimBoundary": "index"},
		"pairingWindow": map[string]interface{}{
			"noFullNxNPairing": true, "fullPairSpace": float64(100), "evaluatedPairCount": float64(25), "avoidedPairCount": float64(75),
			"evidenceRefs": []interface{}{"pair"}, "claimBoundary": "pairing",
		},
		"tickMetrics": map[string]interface{}{
			"generalProductionAutoEnabledGauge": float64(0), "productionAutoWithPolicyEnabledGauge": float64(0),
			"evidenceRefs": []interface{}{"metrics"}, "claimBoundary": "metrics",
		},
		"dashboardProjection": map[string]interface{}{"projectionReadOnly": true, "evidenceRefs": []interface{}{"dash"}, "claimBoundary": "dash"},
		"safetyBoundary": map[string]interface{}{
			"productionMovementExecuted": false, "broadProductionMutationExecuted": false,
			"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
			"dashboardExecutor": false, "fluidvirtReconcilerAuthority": false,
			"rawRuntimeControlsExposed": false, "clusterAdmin": false, "podsExecAllowed": false, "nodesWriteAllowed": false,
			"evidenceRefs": []interface{}{"safety"}, "claimBoundary": "safety",
		},
		"healthDecision": map[string]interface{}{"evidenceRefs": []interface{}{"health"}, "claimBoundary": "health"},
		"generatedActions": []interface{}{},
		"generatedFutures": []interface{}{},
	}
}
