package marketcontroller

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const LiveMilestone = MilestoneLiveControllerReconciliation

var liveForbiddenScopes = []string{"general_production_auto", "production_auto_with_policy"}

var requiredAuditEventTypes = []string{
	"state_loaded", "observed_snapshot_collected", "desired_state_computed",
	"reconciliation_diff_computed", "lease_created", "action_lifecycle_updated",
	"future_refreshed", "invalidation_recorded", "execution_selected",
	"execution_handed_off", "post_execution_reconciled", "compression_tracker_updated", "state_saved",
}

// ValidateLiveSurface validates Sprint 13 live reconciliation surface.
func ValidateLiveSurface(doc map[string]interface{}) error {
	if strOr(doc["milestone"]) != LiveMilestone {
		return fmt.Errorf("milestone must be %s", LiveMilestone)
	}
	for _, key := range []string{
		"liveReconciliationEnabled", "stateStoreEnabled", "scheduledTicksEnabled",
		"leaseLifecycleEnabled", "actionLifecycleEnabled", "futuresRefreshEnabled",
		"executionSelectionEnabled", "realizedCompressionTrackingEnabled",
	} {
		if err := requireTrue(doc, key); err != nil {
			return err
		}
	}
	for _, key := range []string{
		"generalProductionAutoAllowed", "productionAutoWithPolicy",
		"universalGuaranteedSavingsAllowed", "universalGuaranteedSavingsClaimed",
		"estimatedIdleCountedAsMoved", "projectedCompressionCountedAsRealized",
		"syntheticFleetCountedAsProduction", "referenceFleetCountedAsProduction",
		"dashboardExecutor", "fluidvirtPolicyAuthority", "inventoryRuntimeExecutor",
	} {
		if err := requireFalse(doc, key); err != nil {
			return err
		}
	}
	mode := strOr(doc["controllerMode"])
	if mode == "general_production_auto" || mode == "production_auto_with_policy" {
		return fmt.Errorf("forbidden controllerMode %s", mode)
	}
	if tracker, ok := doc["realizedCompressionTracker"].(map[string]interface{}); ok {
		if err := requireFalse(tracker, "projectedCompressionCountedAsRealized"); err != nil {
			return err
		}
	}
	if selections, ok := doc["executionSelections"].([]interface{}); ok {
		for _, item := range selections {
			sel, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			scope := strOr(sel["executionScope"])
			for _, f := range liveForbiddenScopes {
				if scope == f {
					return fmt.Errorf("forbidden executionScope %s", scope)
				}
			}
			if boolOr(sel["productionScope"]) && !boolOr(sel["productionCanaryScope"]) {
				return fmt.Errorf("productionScope requires productionCanaryScope")
			}
			if boolOr(sel["generalProductionAutoAllowed"]) || boolOr(sel["productionAutoWithPolicy"]) {
				return fmt.Errorf("forbidden production auto flags on selection")
			}
		}
	}
	if handoffs, ok := doc["executionHandoffs"].([]interface{}); ok {
		for _, item := range handoffs {
			h, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if boolOr(h["accepted"]) && strOr(h["actuator"]) != "FluidVirt" {
				return fmt.Errorf("accepted handoff requires FluidVirt actuator")
			}
		}
	}
	if posts, ok := doc["postExecutionReconciliations"].([]interface{}); ok {
		for _, item := range posts {
			p, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if boolOr(p["realizedMovementKept"]) {
				if !boolOr(p["mutationObserved"]) || !boolOr(p["postVerifyPassed"]) {
					return fmt.Errorf("realizedMovementKept requires mutation and post-verify")
				}
				if boolOr(p["rollbackRequired"]) {
					return fmt.Errorf("rollback-required cannot be realized kept")
				}
			}
		}
	}
	if leases, ok := doc["leaseLifecycleStates"].([]interface{}); ok {
		for _, item := range leases {
			lc, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if boolOr(lc["expired"]) && strOr(lc["state"]) == "active" {
				return fmt.Errorf("expired lease cannot remain active")
			}
		}
	}
	if futures, ok := doc["futureLifecycleStates"].([]interface{}); ok {
		for _, item := range futures {
			flc, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if boolOr(flc["expired"]) && strOr(flc["state"]) == "active" {
				return fmt.Errorf("expired future cannot remain active")
			}
		}
	}
	if audit, ok := doc["auditTrail"].(map[string]interface{}); ok {
		events, _ := audit["events"].([]interface{})
		found := map[string]bool{}
		for _, ev := range events {
			e, ok := ev.(map[string]interface{})
			if !ok {
				continue
			}
			found[strOr(e["eventType"])] = true
		}
		for _, want := range requiredAuditEventTypes {
			if !found[want] {
				return fmt.Errorf("missing audit event type %s", want)
			}
		}
	}
	if idem, ok := doc["idempotencyRecords"].([]interface{}); ok {
		for _, item := range idem {
			rec, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if intOr(rec["replayCount"]) > 0 && !boolOr(rec["duplicateSuppressed"]) {
				return fmt.Errorf("idempotent replay must set duplicateSuppressed")
			}
		}
	}
	claim := strings.ToLower(strOr(doc["claimBoundary"]))
	if strings.Contains(claim, "dashboard executor") {
		return fmt.Errorf("dashboard cannot be executor")
	}
	return nil
}

// ValidateLiveReferenceFile validates reference JSON file.
func ValidateLiveReferenceFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	return ValidateLiveSurface(doc)
}
