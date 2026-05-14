package installruntime

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const Milestone = "hyperdensity_controller_install_apply_runtime_health_gate_v1"

func ValidateSurface(doc map[string]interface{}) error {
	if strv(doc["milestone"]) != Milestone {
		return fmt.Errorf("milestone must be %s", Milestone)
	}
	for _, key := range []string{
		"runtimeHealthGateEnabled",
		"controlledInstallScope",
	} {
		if !boolv(doc[key]) {
			return fmt.Errorf("%s must be true", key)
		}
	}
	for _, key := range []string{
		"generalProductionAutoAllowed",
		"productionAutoWithPolicy",
		"universalGuaranteedSavingsAllowed",
		"universalGuaranteedSavingsClaimed",
		"estimatedIdleCountedAsMoved",
		"projectedCompressionCountedAsRealized",
		"syntheticFleetCountedAsProduction",
		"referenceFleetCountedAsProduction",
		"dashboardExecutor",
		"fluidvirtPolicyAuthority",
		"fluidvirtInstallAuthority",
		"inventoryRuntimeExecutor",
	} {
		if boolv(doc[key]) {
			return fmt.Errorf("%s must be false", key)
		}
	}

	if err := validateApplyRequestResult(doc); err != nil {
		return err
	}
	if err := validateHealthSections(doc); err != nil {
		return err
	}
	if err := evaluateRuntimeHealthGate(doc); err != nil {
		return err
	}

	if boolv(doc["productionInstallApplied"]) {
		if !boolv(doc["realApplyEvidencePresent"]) {
			return fmt.Errorf("productionInstallApplied without evidence")
		}
		if strv(doc["installApplyMode"]) != "real_cluster_apply" {
			return fmt.Errorf("productionInstallApplied requires real_cluster_apply mode")
		}
	}
	if len(sliceMap(doc["blockers"])) > 0 && boolv(doc["runtimeHealthGatePassed"]) {
		return fmt.Errorf("runtimeHealthGatePassed true while blockers exist")
	}
	if !boolv(doc["productionInstallEligible"]) && boolv(doc["productionInstallApplied"]) {
		return fmt.Errorf("install applied without prior admission eligibility")
	}
	if !hasNonEmptyStringList(doc["claimBoundaries"]) {
		return fmt.Errorf("claimBoundaries required")
	}
	if list, ok := doc["postInstallAuditEvents"].([]interface{}); !ok || len(list) == 0 {
		return fmt.Errorf("postInstallAuditEvents required")
	}
	return nil
}

func validateApplyRequestResult(doc map[string]interface{}) error {
	req, _ := doc["installApplyRequest"].(map[string]interface{})
	if req == nil {
		return fmt.Errorf("installApplyRequest required")
	}
	if strv(req["claimBoundary"]) == "" || !hasNonEmptyStringList(req["evidenceRefs"]) {
		return fmt.Errorf("installApplyRequest missing evidenceRefs or claimBoundary")
	}
	if !boolv(req["admissionGatePassed"]) && boolv(doc["productionInstallApplied"]) {
		return fmt.Errorf("install applied when admissionGatePassed=false")
	}
	if boolv(req["generalProductionAutoAllowed"]) || boolv(req["productionAutoWithPolicy"]) {
		return fmt.Errorf("forbidden auto flags in apply request")
	}

	res, _ := doc["installApplyResult"].(map[string]interface{})
	if res == nil {
		return fmt.Errorf("installApplyResult required")
	}
	if strv(res["claimBoundary"]) == "" || !hasNonEmptyStringList(res["evidenceRefs"]) {
		return fmt.Errorf("installApplyResult missing evidenceRefs or claimBoundary")
	}
	if boolv(res["productionInstallApplied"]) && !boolv(res["realApplyEvidencePresent"]) {
		return fmt.Errorf("apply result claims applied without evidence")
	}

	inv, _ := doc["appliedManifestInventory"].(map[string]interface{})
	if inv == nil {
		return fmt.Errorf("appliedManifestInventory required")
	}
	if strv(inv["claimBoundary"]) == "" || !hasNonEmptyStringList(inv["evidenceRefs"]) {
		return fmt.Errorf("appliedManifestInventory missing evidenceRefs or claimBoundary")
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
func hasEvidenceRef(v interface{}, token string) bool {
	items, ok := v.([]interface{})
	if !ok || len(items) == 0 {
		return false
	}
	for _, item := range items {
		if s, ok := item.(string); ok && strings.Contains(strings.ToLower(s), strings.ToLower(token)) {
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
