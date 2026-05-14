package installadmission

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const Milestone = "hyperdensity_controller_install_dryrun_cluster_admission_gate_v1"

func ValidateSurface(doc map[string]interface{}) error {
	if strv(doc["milestone"]) != Milestone {
		return fmt.Errorf("milestone must be %s", Milestone)
	}
	for _, key := range []string{
		"productionInstallDryRunEnabled",
		"clientSideDryRunModeled",
		"serverSideDryRunModeled",
		"clusterAdmissionGateEnabled",
		"installFailClosed",
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
		"fluidvirtAdmissionAuthority",
		"inventoryRuntimeExecutor",
	} {
		if boolv(doc[key]) {
			return fmt.Errorf("%s must be false", key)
		}
	}

	if err := validateDryRunRequests(doc); err != nil {
		return err
	}
	serverPass, clientPass, err := evaluateDryRunResults(doc)
	if err != nil {
		return err
	}
	if err := validateAdmissionSections(doc); err != nil {
		return err
	}

	if boolv(doc["productionInstallApplied"]) && !boolv(doc["realApplyEvidencePresent"]) {
		return fmt.Errorf("productionInstallApplied without evidence")
	}
	if len(sliceMap(doc["admissionBlockers"])) > 0 && boolv(doc["productionInstallEligible"]) {
		return fmt.Errorf("productionInstallEligible true while blockers exist")
	}
	if boolv(doc["productionInstallEligible"]) {
		if !serverPass {
			return fmt.Errorf("eligible with server dry-run failed")
		}
		if !clientPass {
			return fmt.Errorf("eligible with client dry-run failed")
		}
		for _, key := range []string{
			"namespacePreflightPassed",
			"rbacPreflightPassed",
			"manifestAdmissionPassed",
			"crdPreflightPassed",
			"durableStatePreflightPassed",
			"leaderElectionPreflightPassed",
			"probesPreflightPassed",
			"metricsPreflightPassed",
			"rollbackPreflightPassed",
			"productionInstallDryRunPassed",
			"clusterAdmissionGatePassed",
		} {
			if !boolv(doc[key]) {
				return fmt.Errorf("eligible with failed preflight: %s", key)
			}
		}
	}
	if !hasNonEmptyStringList(doc["claimBoundaries"]) {
		return fmt.Errorf("claimBoundaries required")
	}
	if !hasNonEmptyStringList(doc["auditEvents"]) {
		// allow object list; verify non-empty separately
		if list, ok := doc["auditEvents"].([]interface{}); !ok || len(list) == 0 {
			return fmt.Errorf("auditEvents required")
		}
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
