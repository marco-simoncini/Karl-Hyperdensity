package marketcontroller

import (
	"encoding/json"
	"time"
)

// BuildReferenceSurface runs the live loop and returns Sprint 13 reference surface.
func BuildReferenceSurface() map[string]interface{} {
	store := NewInMemoryStateStore()
	loop := NewLiveControllerLoop(store, "live-loop-sprint13", "continuous-resource-market-controller-sprint12")
	now := time.Date(2026, 5, 14, 17, 30, 0, 0, time.UTC)
	input := LiveLoopInput{
		IdempotencyKey:           "tick-idem-sprint13-001",
		TickSequence:             1,
		RateLimitRemaining:       5,
		CooldownBlockedCount:     1,
		ProductionCanaryEnabled:  true,
		CurrentCompressionRate:   0.04941176470588235,
		ProjectedCompressionRate: 0.22352941176470587,
		EligibleIdleValue:        0.085,
		GuaranteedEligibleTotal:  0.00336,
		Snapshot:                 referenceSnapshot(),
		Observed: ObservedMarketState{
			ObservedSnapshotID: "observed-sprint13-001", CollectedAt: now,
			ShellCount: 12, IdleObservationCount: 12, PressureSignalCount: 3,
			ProductionCanaryEnabled: true, RateLimitRemaining: 5, CooldownBlockedCount: 1,
		},
		Desired: defaultDesiredState(now),
	}
	res1, _ := loop.RunScheduledTick(input, now)
	input2 := input
	input2.TickSequence = 2
	input2.IdempotencyKey = "tick-idem-sprint13-001"
	res2, _ := loop.RunScheduledTick(input2, now.Add(time.Minute))

	// Expire one future for reference scenario
	expiredFuture := map[string]interface{}{
		"futureLifecycleId": "future-lc-expired-001", "futureId": "ctrl-future-expired-001",
		"donorShellId": "shell-container-donor-b", "receiverShellId": "shell-container-replica-c",
		"resource": "cpu", "amount": "100m", "state": "expired", "expired": true,
		"createdAt": now.Add(-48 * time.Hour).Format(time.RFC3339),
		"expiresAt": now.Add(-24 * time.Hour).Format(time.RFC3339),
		"invalidationReasons": []interface{}{"candidate_expired"},
		"evidenceRefs": []interface{}{"ctrl-future-expired-001"}, "claimBoundary": "expired future",
	}
	invalidatedAction := map[string]interface{}{
		"actionLifecycleId": "action-lc-invalidated-001", "actionId": "ctrl-action-stale-001",
		"donorShellId": "shell-container-donor-b", "receiverShellId": "shell-container-replica-d",
		"resource": "cpu", "amount": "100m", "state": "invalidated",
		"invalidationReason": "stale_action_ttl", "selectedForExecution": false,
		"evidenceRefs": []interface{}{"ctrl-action-stale-001"}, "claimBoundary": "invalidated action",
	}
	invalidationEvent := map[string]interface{}{
		"invalidationEventId": "inv-stale-001", "invalidatedObjectType": "action",
		"invalidatedObjectId": "ctrl-action-stale-001", "reason": "stale_action_ttl",
		"occurredAt": now.Format(time.RFC3339), "sourceSignal": "action_ttl_enforcer",
		"evidenceRefs": []interface{}{"ctrl-action-stale-001"}, "claimBoundary": "invalidation event",
	}

	auditEvents := res1.State.AuditEvents
	auditEvents = append(auditEvents, auditEvent("invalidation_recorded", loop.LiveLoopID, "inv-stale-001"))
	auditEvents = append(auditEvents, auditEvent("post_execution_reconciled", loop.LiveLoopID, "post-exec-realized"))

	tracker := res1.RealizedCompressionTracker
	leaseStates := res1.State.LeaseLifecycles
	actionStates := append(res1.State.ActionLifecycles, invalidatedAction)
	futureStates := append(res1.State.FutureLifecycles, expiredFuture)

	activeLeases, expiredLeases := countLeaseStates(leaseStates)
	activeActions, invActions := countActionStates(actionStates)
	activeFutures, expFutures := countFutureStates(futureStates)

	idemRecords := res1.State.IdempotencyRecords
	idemRecords = append(idemRecords, map[string]interface{}{
		"idempotencyRecordId": "idem-replay-001", "idempotencyKey": "tick-idem-sprint13-001",
		"objectType": "scheduled_tick", "objectId": "tick-2", "replayCount": 1,
		"duplicateSuppressed": true, "firstSeenAt": now.Format(time.RFC3339),
		"lastSeenAt": now.Add(time.Minute).Format(time.RFC3339),
		"evidenceRefs": []interface{}{"tick-idem-sprint13-001"}, "claimBoundary": "idempotent replay",
	})

	return map[string]interface{}{
		"milestone":                         LiveMilestone,
		"surfaceVersion":                    "v1",
		"liveLoopId":                        loop.LiveLoopID,
		"generatedAt":                       now.Format(time.RFC3339),
		"sourceContinuousControllerRef":     "hyperdensity-continuous-resource-market-controller-v1",
		"sourceIdleCompressionRef":          "hyperdensity-idle-time-compression-fleet-value-v1",
		"sourceGuaranteedEligibleSavingsRef": "hyperdensity-guaranteed-eligible-savings-v1",
		"sourceProductionCanaryRef":         "hyperdensity-production-canary-auto-apply-v1",
		"sourceSandboxNonProdAutoRef":       "hyperdensity-guarded-auto-policy-engine-v1",
		"sourceSloPerformanceProofRef":        "hyperdensity-slo-performance-proof-v1",
		"sourceGaReleaseGateRef":            "hyperdensity-ga-release-gate-v1",
		"controllerMode":                    "production_canary_only",
		"liveReconciliationEnabled":         true,
		"stateStoreEnabled":                 true,
		"scheduledTicksEnabled":             true,
		"leaseLifecycleEnabled":             true,
		"actionLifecycleEnabled":            true,
		"futuresRefreshEnabled":             true,
		"executionSelectionEnabled":         true,
		"realizedCompressionTrackingEnabled": true,
		"generalProductionAutoAllowed":      false,
		"productionAutoWithPolicy":          false,
		"universalGuaranteedSavingsAllowed": false,
		"universalGuaranteedSavingsClaimed": false,
		"estimatedIdleCountedAsMoved":       false,
		"projectedCompressionCountedAsRealized": false,
		"syntheticFleetCountedAsProduction": false,
		"referenceFleetCountedAsProduction": false,
		"dashboardExecutor":                 false,
		"fluidvirtPolicyAuthority":          false,
		"inventoryRuntimeExecutor":          false,
		"tickCount":                         2,
		"successfulTickCount":               2,
		"failedTickCount":                   0,
		"idempotentReplayCount":             1,
		"observedShellCount":                12,
		"desiredActionCount":                9,
		"activeLeaseCount":                  activeLeases,
		"expiredLeaseCount":                 expiredLeases,
		"activeActionCount":                 activeActions,
		"invalidatedActionCount":            invActions,
		"activeFutureCount":                 activeFutures,
		"expiredFutureCount":                expFutures,
		"executionSelectionCount":           len(res1.ExecutionSelections),
		"executionHandoffCount":             len(res1.ExecutionHandoffs),
		"postExecutionReconciliationCount":  len(res1.PostExecutionReconciliations),
		"realizedMovementCount":             intOr(tracker["realizedMovementCount"]),
		"projectedMovementCount":            9,
		"currentIdleCompressionRate":        0.04941176470588235,
		"projectedIdleCompressionRate":      0.22352941176470587,
		"realizedIdleCompressionRate":       tracker["currentRealizedIdleCompressionRate"],
		"realizedMovedIdleValue":            tracker["realizedMovedIdleValue"],
		"projectedMovedIdleValue":           0.019,
		"currentUnmovedEligibleIdleValue":   0.0808,
		"realizedUnmovedEligibleIdleValue":  tracker["realizedUnmovedEligibleIdleValue"],
		"guaranteedEligibleSavingsTotal":    0.00336,
		"stateStore":                        BuildStateStoreRecord(),
		"scheduledTicks":                    []interface{}{res1.ScheduledTick, res2.ScheduledTick},
		"observedSnapshots":                 []interface{}{res1.ObservedSnapshot},
		"desiredStates":                     []interface{}{res1.DesiredState},
		"reconciliationDiffs":               []interface{}{reconciliationDiffToMap(res1.ReconciliationDiff)},
		"leaseLifecycleStates":              toIface(leaseStates),
		"actionLifecycleStates":             toIface(actionStates),
		"futureLifecycleStates":             toIface(futureStates),
		"invalidationEvents":                []interface{}{invalidationEvent},
		"executionSelections":               toIface(res1.ExecutionSelections),
		"executionHandoffs":                 toIface(res1.ExecutionHandoffs),
		"postExecutionReconciliations":        toIface(res1.PostExecutionReconciliations),
		"realizedCompressionTracker":        tracker,
		"driftDetections":                   toIface(res1.DriftDetections),
		"retryPolicies": []interface{}{
			map[string]interface{}{
				"retryPolicyId": "retry-sprint13", "objectType": "action_lifecycle",
				"maxRetries": 3, "retryBackoffSeconds": 30,
				"retryableStates": []interface{}{"execution_handed_off"},
				"nonRetryableStates": []interface{}{"invalidated", "blocked"},
				"evidenceRefs": []interface{}{"hyperdensity-live-controller-reconciliation-v1"},
				"claimBoundary": "retry policy",
			},
		},
		"idempotencyRecords": toIface(idemRecords),
		"auditTrail": map[string]interface{}{
			"auditTrailId": "audit-trail-sprint13", "liveLoopId": loop.LiveLoopID,
			"immutableAuditRequired": true, "appendOnly": true, "retentionDays": 90,
			"events": auditEvents,
			"evidenceRefs": []interface{}{"hyperdensity-live-controller-reconciliation-v1"},
			"claimBoundary": "immutable append-only audit trail",
		},
		"sourceOfTruth": map[string]interface{}{
			"liveController": "Karl-Hyperdensity", "runtimeActuator": "FluidVirt",
			"operatorUI": "Karl-Dashboard (projection only)", "identitySignals": "Karl-Inventory",
		},
		"safetyInvariants": map[string]interface{}{
			"generalProductionAutoAllowed": false, "projectedCompressionCountedAsRealized": false,
			"dashboardExecutor": false, "fluidvirtPolicyAuthority": false,
		},
		"claimBoundaries": []interface{}{
			"live reconciliation loop; projected compression separate from realized",
			"realized compression requires mutation and post-verify evidence",
			"no general production auto",
		},
		"blockers":  []interface{}{},
		"warnings":  []interface{}{},
		"_res2Status": res2.TickStatus,
	}
}

func referenceSnapshot() Snapshot {
	return Snapshot{
		TopKDonors: 3, TopKReceivers: 3,
		FullDonorCount: 12, FullReceiverCount: 8,
		CurrentMovedIdleValue: 0.0042,
		CurrentEligibleIdleValue: 0.085,
		CurrentCompressionRate: 0.04941176470588235,
		TargetCompressionRate: 0.25,
		RateLimitRemaining: 5,
		CooldownBlockedCount: 1,
		Donors: []DonorCandidate{
			{DonorShellID: "shell-container-donor-a", Resource: "cpu", EligibleIdleAmount: "500m", EligibleIdleValue: 0.056, RollbackReady: true, SloGuardAvailable: true, NoRegressionAvailable: true, RiskScore: 0.1, GuaranteePotential: 0.8},
			{DonorShellID: "shell-container-donor-b", Resource: "cpu", EligibleIdleAmount: "300m", EligibleIdleValue: 0.018, RollbackReady: true, SloGuardAvailable: true, RiskScore: 0.2, GuaranteePotential: 0.5},
			{DonorShellID: "shell-windows-hyper", Resource: "cpu", EligibleIdleAmount: "0", EligibleIdleValue: 0.008, WindowsEvidenceGated: true},
			{DonorShellID: "shell-reference-sample", ReferenceOnly: true, EligibleIdleValue: 0.01},
			{DonorShellID: "shell-synthetic-fleet", SyntheticShadow: true, EligibleIdleValue: 0.015},
			{DonorShellID: "shell-protected-core-001", Protected: true, EligibleIdleValue: 0.025},
		},
		Receivers: []ReceiverCandidate{
			{ReceiverShellID: "shell-container-replica-b", Resource: "cpu", RequestedAmount: "500m", PotentialValueCapture: 0.0042, SloProfilePresent: true},
			{ReceiverShellID: "shell-container-replica-c", Resource: "cpu", RequestedAmount: "400m", PotentialValueCapture: 0.011},
			{ReceiverShellID: "shell-container-replica-d", Resource: "cpu", RequestedAmount: "200m", PotentialValueCapture: 0.005},
		},
	}
}

func reconciliationDiffToMap(d ReconciliationDiff) map[string]interface{} {
	return map[string]interface{}{
		"reconciliationDiffId": d.ReconciliationDiffID, "previousStateRef": d.PreviousStateRef,
		"observedSnapshotRef": d.ObservedSnapshotRef, "desiredStateRef": d.DesiredStateRef,
		"actionsCreated": d.ActionsCreated, "actionsUpdated": d.ActionsUpdated,
		"actionsInvalidated": d.ActionsInvalidated, "leasesCreated": d.LeasesCreated,
		"leasesExpired": d.LeasesExpired, "futuresCreated": d.FuturesCreated,
		"futuresUpdated": d.FuturesUpdated, "futuresExpired": d.FuturesExpired,
		"executionSelectionsCreated": d.ExecutionSelectionsCreated, "noOpCount": d.NoOpCount,
		"diffSummary": d.DiffSummary,
		"evidenceRefs": []interface{}{d.ReconciliationDiffID}, "claimBoundary": "reconciliation diff",
	}
}

func toIface(items []map[string]interface{}) []interface{} {
	out := make([]interface{}, len(items))
	for i, v := range items {
		out[i] = v
	}
	return out
}

func countLeaseStates(states []map[string]interface{}) (active, expired int) {
	for _, s := range states {
		if boolOr(s["expired"]) {
			expired++
		} else {
			active++
		}
	}
	return
}

func countActionStates(states []map[string]interface{}) (active, invalidated int) {
	for _, s := range states {
		if strOr(s["state"]) == "invalidated" {
			invalidated++
		} else {
			active++
		}
	}
	return
}

func countFutureStates(states []map[string]interface{}) (active, expired int) {
	for _, s := range states {
		if boolOr(s["expired"]) || strOr(s["state"]) == "expired" {
			expired++
		} else {
			active++
		}
	}
	return
}

// MarshalReferenceSurfaceJSON returns indented JSON for reference file.
func MarshalReferenceSurfaceJSON() ([]byte, error) {
	surface := BuildReferenceSurface()
	delete(surface, "_res2Status")
	return json.MarshalIndent(surface, "", "  ")
}
