package marketcontroller

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
)

const Milestone = MilestoneContinuousResourceMarketController

var forbiddenScopes = map[string]bool{
	"general_production_auto": true, "production_auto_with_policy": true,
}

var forbiddenPositiveClaims = []string{
	"general production auto",
	"production auto with policy",
	"production_auto_with_policy",
	"universal guaranteed savings",
	"universal performance improvement",
	"windows total ram hotplug supported",
	"logical vcpu hotplug supported",
	"1000 production workloads proven",
	"dashboard executor",
	"fluidvirt policy authority",
	"inventory runtime executor",
}

func boolOr(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

func floatOr(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	default:
		return 0
	}
}

func intOr(v interface{}) int {
	switch t := v.(type) {
	case int:
		return t
	case float64:
		return int(t)
	default:
		return 0
	}
}

func requireFalse(doc map[string]interface{}, key string) error {
	if v, ok := doc[key].(bool); !ok || v {
		return fmt.Errorf("%s must be false", key)
	}
	return nil
}

func requireTrue(doc map[string]interface{}, key string) error {
	if v, ok := doc[key].(bool); !ok || !v {
		return fmt.Errorf("%s must be true", key)
	}
	return nil
}

func requireClaimBoundary(doc map[string]interface{}) error {
	if cbs, ok := doc["claimBoundaries"].([]interface{}); ok && len(cbs) > 0 {
		return nil
	}
	cb, ok := doc["claimBoundary"].(string)
	if !ok || strings.TrimSpace(cb) == "" {
		return fmt.Errorf("claimBoundary required")
	}
	return nil
}

func requireEvidenceRefs(doc map[string]interface{}) error {
	refs, ok := doc["evidenceRefs"].([]interface{})
	if !ok || len(refs) == 0 {
		return fmt.Errorf("evidenceRefs required")
	}
	return nil
}

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}

// ValidateControllerSurface checks Sprint 12 surface invariants.
func ValidateControllerSurface(doc map[string]interface{}) error {
	if m, _ := doc["milestone"].(string); m != Milestone {
		return fmt.Errorf("milestone must be %s", Milestone)
	}
	for _, key := range []string{
		"continuousControllerEnabled", "controllerTickExecuted", "noFullNxNPairing",
		"actionSlateGeneratedByController", "resourceFuturesGeneratedByController",
		"controllerCoverageExpanded", "idleCompressionTargetTracked",
	} {
		if err := requireTrue(doc, key); err != nil {
			return err
		}
	}
	for _, key := range []string{
		"generalProductionAutoAllowed", "productionAutoWithPolicy",
		"universalGuaranteedSavingsAllowed", "universalGuaranteedSavingsClaimed",
		"estimatedIdleCountedAsMoved", "syntheticFleetCountedAsProduction",
		"referenceFleetCountedAsProduction", "dashboardExecutor",
		"fluidvirtPolicyAuthority", "inventoryRuntimeExecutor",
	} {
		if err := requireFalse(doc, key); err != nil {
			return err
		}
	}
	full := int(floatOr(doc["fullPairSpace"]))
	evaluated := int(floatOr(doc["evaluatedPairCount"]))
	avoided := int(floatOr(doc["avoidedPairCount"]))
	topKDonors := int(floatOr(doc["topKDonors"]))
	topKReceivers := int(floatOr(doc["topKReceivers"]))
	if evaluated > topKDonors*topKReceivers {
		return fmt.Errorf("evaluatedPairCount exceeds top-K bound")
	}
	if avoided != full-evaluated {
		return fmt.Errorf("avoidedPairCount must equal fullPairSpace - evaluatedPairCount")
	}
	if projected := floatOr(doc["projectedIdleCompressionRate"]); projected <= floatOr(doc["currentIdleCompressionRate"]) {
		return fmt.Errorf("projectedIdleCompressionRate must exceed currentIdleCompressionRate")
	}
	if actions, ok := doc["generatedActionSlate"].(map[string]interface{}); ok {
		if acts, ok := actions["actions"].([]interface{}); ok {
			for _, item := range acts {
				a, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				if err := ValidateGeneratedAction(a); err != nil {
					return err
				}
			}
		}
	}
	if futures, ok := doc["generatedResourceFutures"].([]interface{}); ok {
		for _, item := range futures {
			f, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateGeneratedFuture(f); err != nil {
				return err
			}
		}
	}
	if doc["backpressure"] == nil {
		return fmt.Errorf("backpressure required")
	}
	if rules, ok := doc["invalidationRules"].([]interface{}); !ok || len(rules) == 0 {
		return fmt.Errorf("invalidationRules required")
	}
	return requireClaimBoundary(doc)
}

// ValidateGeneratedAction checks a generated action.
func ValidateGeneratedAction(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	for _, key := range []string{"generalProductionAutoAllowed", "productionAutoWithPolicy"} {
		if err := requireFalse(doc, key); err != nil {
			return err
		}
	}
	scope, _ := doc["executionScopeRecommendation"].(string)
	if forbiddenScopes[scope] {
		return fmt.Errorf("forbidden execution scope: %s", scope)
	}
	if boolOr(doc["productionScope"]) && !boolOr(doc["productionCanaryScope"]) {
		return fmt.Errorf("productionScope requires productionCanaryScope")
	}
	for _, key := range []string{"actionId", "donorShellId", "receiverShellId", "resource", "amount"} {
		if doc[key] == nil || doc[key] == "" {
			return fmt.Errorf("action missing %s", key)
		}
	}
	return nil
}

// ValidateGeneratedFuture checks a generated future.
func ValidateGeneratedFuture(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	if doc["expiration"] == nil || doc["expiration"] == "" {
		return fmt.Errorf("future missing expiration")
	}
	inv, ok := doc["invalidationReasons"].([]interface{})
	if !ok || len(inv) == 0 {
		return fmt.Errorf("future missing invalidationReasons")
	}
	return nil
}

// ValidateReferenceFile validates main reference payload.
func ValidateReferenceFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return err
	}
	return ValidateControllerSurface(doc)
}
