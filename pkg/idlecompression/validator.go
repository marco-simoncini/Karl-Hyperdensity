package idlecompression

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
)

const MilestoneIdleTimeCompressionFleetValue = "hyperdensity_idle_time_compression_fleet_value_coverage_v1"

const highSavingsOpportunityThreshold = 0.05

var forbiddenPositiveClaims = []string{
	"universal guaranteed savings",
	"guaranteed savings for all workloads",
	"estimated idle value realized",
	"estimated idle counted as moved",
	"universal performance improvement",
	"general production auto",
	"production auto with policy",
	"production_auto_with_policy",
	"windows total ram hotplug supported",
	"logical vcpu hotplug supported",
	"1000 production workloads proven",
	"dashboard measurement authority",
	"dashboard accounting authority",
	"dashboard guarantee authority",
	"fluidvirt accounting authority",
	"fluidvirt guarantee authority",
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

func strOr(v interface{}) string {
	s, _ := v.(string)
	return s
}

func requireFalse(doc map[string]interface{}, key string) error {
	if v, ok := doc[key].(bool); !ok || v {
		return fmt.Errorf("%s must be false", key)
	}
	return nil
}

func requireClaimBoundary(doc map[string]interface{}) error {
	if cbs, ok := doc["claimBoundaries"].([]interface{}); ok && len(cbs) > 0 {
		return nil
	}
	switch cb := doc["claimBoundary"].(type) {
	case string:
		if strings.TrimSpace(cb) == "" {
			return fmt.Errorf("claimBoundary required")
		}
	case []interface{}:
		if len(cb) == 0 {
			return fmt.Errorf("claimBoundary required")
		}
	default:
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
	return math.Abs(a-b) < 1e-9
}

func scanForbiddenClaims(text string) error {
	lower := strings.ToLower(text)
	for _, phrase := range forbiddenPositiveClaims {
		if strings.Contains(lower, phrase) {
			if strings.Contains(lower, "remains disabled") || strings.Contains(lower, "not production") ||
				strings.Contains(lower, "not guaranteed") || strings.Contains(lower, "excluded") ||
				strings.Contains(lower, "never counted") || strings.Contains(lower, "not universal") ||
				strings.Contains(lower, "projection only") || strings.Contains(lower, "not general production") {
				continue
			}
			return fmt.Errorf("forbidden claim phrase: %s", phrase)
		}
	}
	return nil
}

// ValidateCoverageSurface checks Sprint 11B coverage invariants.
func ValidateCoverageSurface(doc map[string]interface{}) error {
	if m, _ := doc["milestone"].(string); m != MilestoneIdleTimeCompressionFleetValue {
		return fmt.Errorf("milestone must be %s", MilestoneIdleTimeCompressionFleetValue)
	}
	for _, key := range []string{
		"universalGuaranteedSavingsAllowed", "universalGuaranteedSavingsClaimed",
		"generalProductionAutoAllowed", "productionAutoWithPolicy", "universalPerformanceImprovementClaimed",
		"estimatedIdleCountedAsMoved", "syntheticFleetCountedAsProduction", "referenceFleetCountedAsProduction",
	} {
		if err := requireFalse(doc, key); err != nil {
			return err
		}
	}
	if !boolOr(doc["idleCompressionMeasured"]) || !boolOr(doc["fleetValueCoverageMeasured"]) {
		return fmt.Errorf("idleCompressionMeasured and fleetValueCoverageMeasured must be true")
	}
	eligible := floatOr(doc["eligibleIdleValue"])
	moved := floatOr(doc["movedIdleValue"])
	if eligible < moved {
		return fmt.Errorf("eligibleIdleValue cannot be less than movedIdleValue")
	}
	if floatOr(doc["unmovedEligibleIdleValue"]) < 0 {
		return fmt.Errorf("unmovedEligibleIdleValue cannot be negative")
	}
	if eligible > 0 {
		expectedRate := moved / eligible
		if rate := floatOr(doc["idleCompressionRate"]); !approxEqual(rate, expectedRate) {
			return fmt.Errorf("idleCompressionRate inconsistent with moved/eligible")
		}
	}
	realized := floatOr(doc["realizedMovedIdleValue"])
	guaranteed := floatOr(doc["guaranteedEligibleSavingsTotal"])
	if realized > 0 {
		expectedCov := guaranteed / realized
		if cov := floatOr(doc["guaranteeCoveragePercent"]); !approxEqual(cov, expectedCov) {
			return fmt.Errorf("guaranteeCoveragePercent inconsistent with guaranteed/realized")
		}
	}
	if guaranteed > realized && realized > 0 {
		return fmt.Errorf("guaranteedEligibleSavingsTotal cannot exceed realizedMovedIdleValue")
	}
	total := floatOr(doc["totalIdleValue"])
	if total > 0 {
		expectedLiq := eligible / total
		if liq := floatOr(doc["fleetLiquidityRate"]); !approxEqual(liq, expectedLiq) {
			return fmt.Errorf("fleetLiquidityRate inconsistent with eligible/total")
		}
	}
	unmoved := floatOr(doc["unmovedEligibleIdleValue"])
	if boolOr(doc["highSavingsOpportunityIdentified"]) == false && unmoved > highSavingsOpportunityThreshold {
		return fmt.Errorf("highSavingsOpportunityIdentified must be true when unmovedEligibleIdleValue exceeds threshold")
	}
	if inv, ok := doc["safetyInvariants"].(map[string]interface{}); ok {
		for _, key := range []string{"dashboardMeasurementAuthority", "fluidvirtAccountingAuthority", "fluidvirtGuaranteeAuthority"} {
			if boolOr(inv[key]) {
				return fmt.Errorf("safetyInvariants.%s must be false", key)
			}
		}
	}
	if sep, ok := doc["productionEvidenceSeparation"].(map[string]interface{}); ok {
		if err := ValidateProductionEvidenceSeparation(sep); err != nil {
			return err
		}
	}
	if movedResources, ok := doc["movedIdleResources"].([]interface{}); ok {
		for _, item := range movedResources {
			mr, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateMovedIdleResource(mr); err != nil {
				return fmt.Errorf("moved idle resource %v: %w", mr["movedIdleResourceId"], err)
			}
		}
	}
	if obs, ok := doc["fleetIdleObservations"].([]interface{}); ok {
		for _, item := range obs {
			o, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateFleetIdleObservation(o); err != nil {
				return err
			}
		}
	}
	if cbs, ok := doc["claimBoundaries"].([]interface{}); ok {
		for _, item := range cbs {
			if s, ok := item.(string); ok {
				if err := scanForbiddenClaims(s); err != nil {
					return err
				}
			}
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateFleetIdleObservation checks observation rules.
func ValidateFleetIdleObservation(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	if boolOr(doc["referenceOnly"]) && boolOr(doc["productionEvidence"]) {
		return fmt.Errorf("referenceOnly cannot have productionEvidence=true")
	}
	if boolOr(doc["syntheticShadow"]) && boolOr(doc["productionEvidence"]) {
		return fmt.Errorf("syntheticShadow cannot have productionEvidence=true")
	}
	if ec, _ := doc["environmentClass"].(string); ec == "production_canary" && boolOr(doc["productionEvidence"]) {
		// canary may have canaryEvidence but not general production - productionEvidence on observation should be false for canary class unless explicitly canaryEvidence
		if boolOr(doc["productionEvidence"]) && !boolOr(doc["canaryEvidence"]) {
			return fmt.Errorf("production_canary observation cannot claim general productionEvidence")
		}
	}
	return nil
}

// ValidateMovedIdleResource checks moved resource rules.
func ValidateMovedIdleResource(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	if !boolOr(doc["mutationObserved"]) {
		return fmt.Errorf("moved idle resource requires mutationObserved=true")
	}
	src, _ := doc["movementSource"].(string)
	ec, _ := doc["environmentClass"].(string)
	if (src == "sandbox_auto" || src == "nonprod_auto" || ec == "sandbox" || ec == "nonprod") && boolOr(doc["productionEvidence"]) {
		return fmt.Errorf("sandbox/nonprod movement cannot count as production proof")
	}
	if src == "production_canary_auto" && !boolOr(doc["productionCanaryScope"]) {
		return fmt.Errorf("production_canary_auto requires productionCanaryScope=true")
	}
	if boolOr(doc["postVerifyPassed"]) == false && floatOr(doc["realizedValue"]) > 0 {
		return fmt.Errorf("postVerifyPassed=false cannot count toward realizedMovedIdleValue")
	}
	return nil
}

// ValidateProductionEvidenceSeparation checks evidence separation.
func ValidateProductionEvidenceSeparation(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	for _, key := range []string{"referenceCountedAsProduction", "syntheticCountedAsProduction", "canaryCountedAsGeneralProduction"} {
		if err := requireFalse(doc, key); err != nil {
			return err
		}
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
	return ValidateCoverageSurface(doc)
}
