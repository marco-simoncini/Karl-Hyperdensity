package observationwindow

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const Milestone = "hyperdensity_production_observation_window_multitick_compression_evidence_v1"

func ValidateSurface(doc map[string]interface{}) error {
	if strv(doc["milestone"]) != Milestone {
		return fmt.Errorf("milestone must be %s", Milestone)
	}
	for _, key := range []string{
		"generalProductionAutoAllowed", "productionAutoWithPolicy",
		"universalGuaranteedSavingsAllowed", "universalGuaranteedSavingsClaimed",
		"estimatedIdleCountedAsMoved", "projectedCompressionCountedAsRealized",
		"projectedValueCountedAsRealized", "syntheticFleetCountedAsProduction",
		"referenceFleetCountedAsProduction", "dashboardExecutor",
		"fluidvirtPolicyAuthority", "fluidvirtObservationWindowAuthority",
		"inventoryRuntimeExecutor", "productionMovementExecuted", "broadProductionMutationExecuted",
	} {
		if boolv(doc[key]) {
			return fmt.Errorf("%s must be false", key)
		}
	}
	if num(doc["realizedMovementCount"]) > 0 {
		rvs, _ := doc["realizedValueSeparation"].(map[string]interface{})
		if rvs == nil || !boolv(rvs["movementEvidencePresent"]) {
			return fmt.Errorf("realizedMovementCount > 0 without movement evidence")
		}
	}
	if err := validateWindowSections(doc); err != nil {
		return err
	}
	if err := evaluateWindowPass(doc); err != nil {
		return err
	}
	if len(sliceMap(doc["blockers"])) > 0 && boolv(doc["observationWindowPassed"]) {
		return fmt.Errorf("observationWindowPassed true while blockers exist")
	}
	if !hasNonEmptyStringList(doc["claimBoundaries"]) {
		return fmt.Errorf("claimBoundaries required")
	}
	if list, ok := doc["observationEvents"].([]interface{}); !ok || len(list) == 0 {
		return fmt.Errorf("observationEvents required")
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
