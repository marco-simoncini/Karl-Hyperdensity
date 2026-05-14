package guaranteedsavings

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const MilestoneGuaranteedEligibleSavingsActivation = "hyperdensity_guaranteed_eligible_savings_activation_v1"

var forbiddenSourceTypes = map[string]bool{
	"placeholder_non_guarantee": true, "estimated_market_reference": true,
}

var forbiddenPositiveClaims = []string{
	"universal guaranteed savings",
	"guaranteed savings for all workloads",
	"guaranteed savings active without eligibility",
	"estimated value guaranteed",
	"synthetic value guaranteed",
	"reference payload value guaranteed",
	"windows evidence-gated value guaranteed",
	"placeholder price guaranteed",
	"universal performance improvement",
	"general production auto",
	"production auto with policy",
	"production_auto_with_policy",
	"windows total ram hotplug supported",
	"logical vcpu hotplug supported",
	"1000 production workloads proven",
	"dashboard accounting authority",
	"dashboard guarantee authority",
	"fluidvirt accounting authority",
	"fluidvirt guarantee authority",
	"inventory accounting authority",
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

func scanForbiddenClaims(text string) error {
	lower := strings.ToLower(text)
	for _, phrase := range forbiddenPositiveClaims {
		if strings.Contains(lower, phrase) {
			if strings.Contains(lower, "remains disabled") || strings.Contains(lower, "not guaranteed") ||
				strings.Contains(lower, "excluded") || strings.Contains(lower, "forbidden") ||
				strings.Contains(lower, "never counted") || strings.Contains(lower, "not universal") {
				continue
			}
			return fmt.Errorf("forbidden claim phrase: %s", phrase)
		}
	}
	return nil
}

// ValidateActivationSurface checks Sprint 11A activation invariants.
func ValidateActivationSurface(doc map[string]interface{}) error {
	if m, _ := doc["milestone"].(string); m != MilestoneGuaranteedEligibleSavingsActivation {
		return fmt.Errorf("milestone must be %s", MilestoneGuaranteedEligibleSavingsActivation)
	}
	for _, key := range []string{
		"universalGuaranteedSavingsAllowed", "universalGuaranteedSavingsClaimed",
		"estimatedValueCountedAsGuaranteed", "syntheticValueCountedAsGuaranteed",
		"referencePayloadCountedAsGuaranteed", "windowsEvidenceGatedValueGuaranteed",
		"generalProductionAutoAllowed", "productionAutoWithPolicy", "universalPerformanceImprovementClaimed",
	} {
		if err := requireFalse(doc, key); err != nil {
			return err
		}
	}
	if !boolOr(doc["guaranteedEligibleSavingsAllowed"]) || !boolOr(doc["guaranteedEligibleSavingsClaimed"]) {
		return fmt.Errorf("guaranteedEligibleSavingsAllowed and guaranteedEligibleSavingsClaimed must be true in Sprint 11A")
	}
	if scope, _ := doc["guaranteedSavingsScope"].(string); scope != "scoped_eligible_records" {
		return fmt.Errorf("guaranteedSavingsScope must be scoped_eligible_records")
	}
	if inv, ok := doc["safetyInvariants"].(map[string]interface{}); ok {
		for _, key := range []string{
			"universalGuaranteedSavingsAllowed", "dashboardAccountingAuthority", "dashboardGuaranteeAuthority",
			"fluidvirtAccountingAuthority", "fluidvirtGuaranteeAuthority", "inventoryAccountingAuthority",
		} {
			if boolOr(inv[key]) {
				return fmt.Errorf("safetyInvariants.%s must be false", key)
			}
		}
	}
	if records, ok := doc["guaranteedRecords"].([]interface{}); ok {
		for _, item := range records {
			rec, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateGuaranteedRecord(rec); err != nil {
				return fmt.Errorf("guaranteed record %v: %w", rec["guaranteedSavingsRecordId"], err)
			}
		}
	}
	if policies, ok := doc["guaranteePolicies"].([]interface{}); ok {
		for _, item := range policies {
			pol, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateGuaranteePolicy(pol); err != nil {
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

// ValidateGuaranteePolicy checks guarantee policy invariants.
func ValidateGuaranteePolicy(doc map[string]interface{}) error {
	for _, key := range []string{
		"allowEstimatedOpportunity", "allowSyntheticShadow", "allowReferencePayload",
		"allowWindowsEvidenceGated", "allowPlaceholderPrice", "universalGuaranteeAllowed",
		"generalProductionAutoAllowed", "productionAutoWithPolicy",
	} {
		if err := requireFalse(doc, key); err != nil {
			return err
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateGuaranteedRecord checks a guaranteed savings record.
func ValidateGuaranteedRecord(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	status, _ := doc["guaranteeClaimStatus"].(string)
	if status != "guaranteed_eligible" {
		return nil
	}
	net := floatOr(doc["netRealizedValue"])
	gv := floatOr(doc["guaranteedValue"])
	conf := floatOr(doc["confidence"])
	if net <= 0 || gv <= 0 {
		return fmt.Errorf("guaranteed_eligible requires netRealizedValue>0 and guaranteedValue>0")
	}
	if gv > net {
		return fmt.Errorf("guaranteedValue cannot exceed netRealizedValue")
	}
	if !boolOr(doc["sloPreserved"]) || !boolOr(doc["donorHealthPreserved"]) || !boolOr(doc["rollbackReady"]) {
		return fmt.Errorf("guaranteed_eligible requires slo, donor health and rollback ready")
	}
	if boolOr(doc["rollbackFailed"]) {
		return fmt.Errorf("guaranteed_eligible cannot have rollbackFailed=true")
	}
	if floatOr(doc["durationSeconds"]) <= 0 {
		return fmt.Errorf("guaranteed_eligible requires durationSeconds>0")
	}
	for _, key := range []string{"estimatedOnly", "synthetic", "referenceOnly", "windowsEvidenceGated", "placeholderPrice"} {
		if boolOr(doc[key]) {
			return fmt.Errorf("guaranteed_eligible cannot have %s=true", key)
		}
	}
	if conf < 0.75 {
		return fmt.Errorf("guaranteed_eligible confidence below default floor")
	}
	return nil
}

// ValidateEligibilityCheck checks eligibility check document.
func ValidateEligibilityCheck(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	status, _ := doc["eligibilityStatus"].(string)
	if status == "eligible" && !boolOr(doc["eligible"]) {
		return fmt.Errorf("eligible status requires eligible=true")
	}
	if boolOr(doc["estimatedOnly"]) && boolOr(doc["eligible"]) {
		return fmt.Errorf("estimatedOnly cannot be eligible")
	}
	if boolOr(doc["synthetic"]) && boolOr(doc["eligible"]) {
		return fmt.Errorf("synthetic cannot be eligible")
	}
	if boolOr(doc["windowsEvidenceGated"]) && boolOr(doc["eligible"]) {
		return fmt.Errorf("windowsEvidenceGated cannot be eligible")
	}
	if boolOr(doc["placeholderPrice"]) && boolOr(doc["eligible"]) {
		return fmt.Errorf("placeholderPrice cannot be eligible")
	}
	return nil
}

// ValidateValueCalculation checks value calculation invariants.
func ValidateValueCalculation(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	net := floatOr(doc["netRealizedValue"])
	gv := floatOr(doc["guaranteedValue"])
	if gv > net && net > 0 {
		return fmt.Errorf("guaranteedValue cannot exceed netRealizedValue")
	}
	if gv < 0 {
		return fmt.Errorf("negative guaranteedValue forbidden")
	}
	return nil
}

// ValidateUnitPriceAuthority checks unit price authority.
func ValidateUnitPriceAuthority(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	st, _ := doc["sourceType"].(string)
	if forbiddenSourceTypes[st] && doc["authorityStatus"] == "active" {
		return fmt.Errorf("forbidden sourceType %s cannot be active for guarantee", st)
	}
	return nil
}

// ValidateGuaranteeReleaseGate checks release gate.
func ValidateGuaranteeReleaseGate(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	return requireFalse(doc, "universalGuaranteedSavingsAllowed")
}

// ValidateReferenceFile validates main activation reference.
func ValidateReferenceFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return err
	}
	if err := ValidateActivationSurface(doc); err != nil {
		return err
	}
	if grg, ok := doc["guaranteeReleaseGate"].(map[string]interface{}); ok {
		if err := ValidateGuaranteeReleaseGate(grg); err != nil {
			return err
		}
	}
	if checks, ok := doc["eligibilityChecks"].([]interface{}); ok {
		for _, item := range checks {
			ch, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateEligibilityCheck(ch); err != nil {
				return err
			}
		}
	}
	if calcs, ok := doc["valueCalculations"].([]interface{}); ok {
		for _, item := range calcs {
			c, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateValueCalculation(c); err != nil {
				return err
			}
		}
	}
	if prices, ok := doc["unitPriceAuthorities"].([]interface{}); ok {
		for _, item := range prices {
			p, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateUnitPriceAuthority(p); err != nil {
				return err
			}
		}
	}
	return nil
}
