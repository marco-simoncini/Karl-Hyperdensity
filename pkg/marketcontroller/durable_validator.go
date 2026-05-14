package marketcontroller

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const DurableMilestone = MilestoneDurableControllerKubernetesReconciler

var requiredConditionTypes = []string{
	ConditionReady, ConditionDegraded, ConditionGeneralProductionAutoDisabled,
}

var requiredEventReasons = []string{"ControllerTickCompleted", "StatePersisted"}

// ValidateDurableSurface validates Sprint 14 durable reconciler surface.
func ValidateDurableSurface(doc map[string]interface{}) error {
	if strOr(doc["milestone"]) != DurableMilestone {
		return fmt.Errorf("milestone must be %s", DurableMilestone)
	}
	for _, key := range []string{
		"durableStateStoreEnabled", "kubernetesReconcilerEnabled", "fakeClientTestsEnabled",
		"controllerStatusConditionsEnabled", "kubernetesEventsEnabled", "leaderElectionReady",
		"rbacManifestsDefined", "metricsExportDefined", "recoverySemanticsDefined",
	} {
		if err := requireTrue(doc, key); err != nil {
			return err
		}
	}
	for _, key := range []string{
		"generalProductionAutoAllowed", "productionAutoWithPolicy",
		"universalGuaranteedSavingsAllowed", "estimatedIdleCountedAsMoved",
		"projectedCompressionCountedAsRealized", "syntheticFleetCountedAsProduction",
		"referenceFleetCountedAsProduction", "dashboardExecutor", "fluidvirtPolicyAuthority",
		"inventoryRuntimeExecutor",
	} {
		if err := requireFalse(doc, key); err != nil {
			return err
		}
	}
	if store, ok := doc["stateStore"].(map[string]interface{}); ok {
		if strOr(store["storeType"]) == "in_memory_reference" {
			return fmt.Errorf("cannot claim production persistence with in_memory_reference")
		}
		if !boolOr(store["optimisticLockEnabled"]) || !boolOr(store["idempotencyPersists"]) {
			return fmt.Errorf("state store must enforce optimistic lock and idempotency")
		}
	}
	if results, ok := doc["reconcileResults"].([]interface{}); ok {
		for _, item := range results {
			r, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if boolOr(r["staleWriteRejected"]) && !boolOr(r["staleWriteRejected"]) {
				return fmt.Errorf("stale write scenario must set staleWriteRejected")
			}
			if strOr(r["status"]) == "idempotent_replay" && !boolOr(r["duplicateActionsSuppressed"]) {
				return fmt.Errorf("idempotent replay must suppress duplicates")
			}
		}
	}
	if rbac, ok := doc["rbacBoundary"].(map[string]interface{}); ok {
		if boolOr(rbac["clusterAdmin"]) || boolOr(rbac["rawRuntimeControlsAllowed"]) {
			return fmt.Errorf("unsafe RBAC")
		}
	}
	if metrics, ok := doc["metrics"].(map[string]interface{}); ok {
		if floatOr(metrics["generalProductionAutoEnabledGauge"]) > 0 || floatOr(metrics["productionAutoWithPolicyEnabledGauge"]) > 0 {
			return fmt.Errorf("forbidden auto metrics gauge")
		}
	}
	if dm, ok := doc["degradedMode"].(map[string]interface{}); ok {
		if boolOr(dm["active"]) && (!boolOr(dm["failClosed"]) || !boolOr(dm["executionSelectionDisabled"])) {
			return fmt.Errorf("degraded must be fail-closed")
		}
	}
	if conds, ok := doc["statusConditions"].([]interface{}); ok {
		found := map[string]bool{}
		for _, item := range conds {
			c, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			found[strOr(c["type"])] = true
		}
		for _, want := range requiredConditionTypes {
			if !found[want] {
				return fmt.Errorf("missing condition %s", want)
			}
		}
	}
	if events, ok := doc["kubernetesEvents"].([]interface{}); ok {
		found := map[string]bool{}
		for _, item := range events {
			e, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			found[strOr(e["reason"])] = true
		}
		for _, want := range requiredEventReasons {
			if !found[want] {
				return fmt.Errorf("missing event %s", want)
			}
		}
	}
	claim := strings.ToLower(strOr(doc["claimBoundary"]))
	if strings.Contains(claim, "ha production") && !strings.Contains(claim, "not ha production proven") {
		return fmt.Errorf("HA production claim without evidence")
	}
	return nil
}

func ValidateDurableReferenceFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	return ValidateDurableSurface(doc)
}
