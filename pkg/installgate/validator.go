package installgate

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	Milestone = "hyperdensity_controller_deployment_hardening_production_install_gate_v1"
)

func ValidateSurface(doc map[string]interface{}) error {
	if str(doc, "milestone") != Milestone {
		return fmt.Errorf("milestone must be %s", Milestone)
	}
	for _, key := range []string{
		"productionInstallGateEnabled",
		"hardenedManifestsDefined",
		"installFailClosed",
		"manifestLintPassed",
		"rbacHardeningPassed",
		"serviceAccountBoundaryPassed",
		"deploymentSpecBoundaryPassed",
		"probesConfigured",
		"metricsEndpointConfigured",
		"leaderElectionWired",
		"durableStateWired",
		"upgradeSafetyDefined",
		"rollbackSafetyDefined",
		"installAuditEnabled",
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
		"inventoryRuntimeExecutor",
	} {
		if boolv(doc[key]) {
			return fmt.Errorf("%s must be false", key)
		}
	}

	candidate, _ := doc["installCandidate"].(map[string]interface{})
	if candidate == nil {
		return fmt.Errorf("installCandidate is required")
	}
	if boolv(candidate["referenceOnly"]) || boolv(candidate["notProductionInstall"]) {
		if boolv(doc["productionInstallAllowed"]) {
			return fmt.Errorf("productionInstallAllowed true while candidate is reference-only")
		}
	}

	if list, ok := doc["installBlockers"].([]interface{}); ok && len(list) > 0 && boolv(doc["productionInstallAllowed"]) {
		return fmt.Errorf("productionInstallAllowed true with blockers present")
	}

	rbac, _ := doc["rbacHardening"].(map[string]interface{})
	if rbac == nil {
		return fmt.Errorf("rbacHardening is required")
	}
	for _, unsafe := range []string{
		"clusterAdmin",
		"wildcardResourceAllowed",
		"wildcardVerbAllowed",
		"podsExecAllowed",
		"nodesWriteAllowed",
		"rawRuntimeControlsAllowed",
		"directLibvirtAllowed",
		"directCgroupAllowed",
	} {
		if boolv(rbac[unsafe]) {
			return fmt.Errorf("unsafe RBAC: %s", unsafe)
		}
	}
	if !boolv(rbac["namespaceScoped"]) {
		return fmt.Errorf("namespaceScoped must be true")
	}

	dep, _ := doc["deploymentSpecBoundary"].(map[string]interface{})
	if dep == nil {
		return fmt.Errorf("deploymentSpecBoundary is required")
	}
	for _, mustTrue := range []string{"livenessProbe", "readinessProbe", "startupProbe", "resourceRequests", "resourceLimits", "runAsNonRoot"} {
		if !boolv(dep[mustTrue]) {
			return fmt.Errorf("deployment boundary missing: %s", mustTrue)
		}
	}
	for _, mustFalse := range []string{"privileged", "allowPrivilegeEscalation", "hostPID", "hostNetwork", "hostIPC", "generalProductionAutoAllowedEnv", "productionAutoWithPolicyEnv"} {
		if boolv(dep[mustFalse]) {
			return fmt.Errorf("deployment boundary unsafe: %s", mustFalse)
		}
	}

	metrics, _ := doc["metricsEndpoint"].(map[string]interface{})
	if metrics == nil {
		return fmt.Errorf("metricsEndpoint is required")
	}
	if num(metrics["generalProductionAutoEnabledGauge"]) > 0 || num(metrics["productionAutoWithPolicyEnabledGauge"]) > 0 {
		return fmt.Errorf("forbidden auto metrics gauge > 0")
	}

	leader, _ := doc["leaderElectionWiring"].(map[string]interface{})
	if leader == nil {
		return fmt.Errorf("leaderElectionWiring is required")
	}
	if !boolv(leader["enabled"]) || !boolv(leader["rbacConfigured"]) {
		return fmt.Errorf("leader election not fully wired")
	}
	if boolv(leader["haProductionProven"]) {
		if !hasEvidenceRef(leader["evidenceRefs"], "ha") {
			return fmt.Errorf("HA production proven without evidence")
		}
	}

	decision, _ := doc["productionInstallDecision"].(map[string]interface{})
	if decision == nil {
		return fmt.Errorf("productionInstallDecision is required")
	}
	if str(decision, "nextAllowedState") == "production_auto_with_policy" {
		return fmt.Errorf("nextAllowedState cannot be production_auto_with_policy")
	}
	if str(decision, "decision") == "install_allowed" && !boolv(doc["productionInstallAllowed"]) {
		return fmt.Errorf("install_allowed decision requires productionInstallAllowed=true")
	}

	if !hasNonEmptyStringList(doc["claimBoundaries"]) {
		return fmt.Errorf("claimBoundaries required")
	}
	if !hasNonEmptyStringList(candidate["evidenceRefs"]) {
		return fmt.Errorf("installCandidate.evidenceRefs required")
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

func hasNonEmptyStringList(v interface{}) bool {
	items, ok := v.([]interface{})
	if !ok || len(items) == 0 {
		return false
	}
	for _, item := range items {
		if s, ok := item.(string); ok && s != "" {
			return true
		}
	}
	return false
}

func hasEvidenceRef(v interface{}, token string) bool {
	items, ok := v.([]interface{})
	if !ok {
		return false
	}
	for _, item := range items {
		if s, ok := item.(string); ok && s != "" {
			if token == "" {
				return true
			}
			if containsFold(s, token) {
				return true
			}
		}
	}
	return false
}

func containsFold(s, token string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(token))
}

func str(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
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
