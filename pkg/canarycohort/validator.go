package canarycohort

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const Milestone = "hyperdensity_production_canary_cohort_expansion_guarded_policy_graduation_evidence_v1"

func ValidateSurface(doc map[string]interface{}) error {
	if strv(doc["milestone"]) != Milestone {
		return fmt.Errorf("milestone must be %s", Milestone)
	}
	for _, key := range []string{
		"generalProductionAutoAllowed", "productionAutoWithPolicy", "productionAutoWithPolicyEnabled",
		"productionAutoWithPolicyActivated", "universalGuaranteedSavingsAllowed",
		"universalGuaranteedSavingsClaimed", "estimatedIdleCountedAsMoved",
		"projectedCompressionCountedAsRealized", "projectedValueCountedAsRealized",
		"syntheticFleetCountedAsProduction", "referenceFleetCountedAsProduction",
		"dashboardExecutor", "fluidvirtPolicyAuthority", "fluidvirtControllerAuthority",
		"inventoryRuntimeExecutor", "broadProductionMutationExecuted",
	} {
		if boolv(doc[key]) {
			return fmt.Errorf("%s must be false", key)
		}
	}
	if err := validateCohortSections(doc); err != nil {
		return err
	}
	if err := evaluateCohortPass(doc); err != nil {
		return err
	}
	if err := validateAccountingAggregation(doc); err != nil {
		return err
	}
	if err := validateGraduationEvidence(doc); err != nil {
		return err
	}
	if len(sliceMap(doc["blockers"])) > 0 && boolv(doc["productionCanaryCohortExecuted"]) {
		return fmt.Errorf("productionCanaryCohortExecuted true while blockers exist")
	}
	if !hasNonEmptyStringList(doc["claimBoundaries"]) {
		return fmt.Errorf("claimBoundaries required")
	}
	if list, ok := doc["auditEvents"].([]interface{}); !ok || len(list) == 0 {
		return fmt.Errorf("auditEvents required")
	}
	return nil
}

func ValidateReferenceFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return err
	}
	return ValidateSurface(doc)
}

func strv(v interface{}) string {
	s, _ := v.(string)
	return s
}
func boolv(v interface{}) bool {
	b, _ := v.(bool)
	return b
}
func num(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case int64:
		return float64(t)
	default:
		return 0
	}
}
func hasNonEmptyStringList(v interface{}) bool {
	items, ok := v.([]interface{})
	if !ok || len(items) == 0 {
		return false
	}
	for _, item := range items {
		if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
			return true
		}
	}
	return false
}
func sliceMap(v interface{}) []map[string]interface{} {
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}
