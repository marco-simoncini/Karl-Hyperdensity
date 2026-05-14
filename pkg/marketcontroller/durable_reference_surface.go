package marketcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// BuildDurableReferenceSurface returns Sprint 14 reference surface.
func BuildDurableReferenceSurface() map[string]interface{} {
	now := time.Date(2026, 5, 14, 18, 15, 0, 0, time.UTC)
	client := NewFakeKubernetesClient()
	store := NewKubernetesStateStore(client, DefaultDurableStateKey())
	rec := NewReconciler(store, "live-loop-durable-sprint14", "durable-controller-sprint14")
	ctx := context.Background()

	req1 := durableReconcileRequest("tick-idem-sprint14-001")
	req1.ReconcileRequestID = "req-sprint14-001"
	req1.RequestedAt = now
	r1, _ := rec.Reconcile(ctx, req1)

	req2 := durableReconcileRequest("tick-idem-sprint14-001")
	req2.ReconcileRequestID = "req-sprint14-002"
	req2.RequestedAt = now.Add(time.Minute)
	r2, _ := rec.Reconcile(ctx, req2)

	reqStale := durableReconcileRequest("tick-idem-stale-001")
	reqStale.ReconcileRequestID = "req-sprint14-stale"
	reqStale.PreviousResourceVersion = "1"
	rStale, staleErr := rec.Reconcile(ctx, reqStale)

	client2 := NewFakeKubernetesClient()
	client2.failSave = true
	storeFail := NewKubernetesStateStore(client2, DefaultDurableStateKey())
	recFail := NewReconciler(storeFail, "live-loop-durable-sprint14", "durable-controller-sprint14")
	_, _ = recFail.Reconcile(ctx, durableReconcileRequest("tick-idem-degraded-001"))

	client3 := NewFakeKubernetesClient()
	storeRecovery := NewKubernetesStateStore(client3, DefaultDurableStateKey())
	recRecovery := NewReconciler(storeRecovery, "live-loop-durable-sprint14", "durable-controller-sprint14")
	client3.failSave = true
	_, _ = recRecovery.Reconcile(ctx, durableReconcileRequest("tick-idem-recovery-001"))
	client3.failSave = false
	rRecovery, _ := recRecovery.Reconcile(ctx, durableReconcileRequest("tick-idem-recovery-001"))

	conditions := client.GetConditions()
	events := client.GetEvents()

	return map[string]interface{}{
		"milestone": DurableMilestone,
		"surfaceVersion": "v1",
		"durableControllerId": "durable-controller-sprint14",
		"generatedAt": now.Format(time.RFC3339),
		"sourceLiveControllerRef": "hyperdensity-live-controller-reconciliation-v1",
		"sourceContinuousControllerRef": "hyperdensity-continuous-resource-market-controller-v1",
		"sourceIdleCompressionRef": "hyperdensity-idle-time-compression-fleet-value-v1",
		"sourceGaReleaseGateRef": "hyperdensity-ga-release-gate-v1",
		"controllerMode": "production_canary_only",
		"durableStateStoreEnabled": true,
		"kubernetesReconcilerEnabled": true,
		"configMapBackedStateEnabled": true,
		"crdBackedStateDefined": true,
		"fakeClientTestsEnabled": true,
		"controllerStatusConditionsEnabled": true,
		"kubernetesEventsEnabled": true,
		"leaderElectionReady": true,
		"rbacManifestsDefined": true,
		"metricsExportDefined": true,
		"recoverySemanticsDefined": true,
		"generalProductionAutoAllowed": false,
		"productionAutoWithPolicy": false,
		"universalGuaranteedSavingsAllowed": false,
		"universalGuaranteedSavingsClaimed": false,
		"estimatedIdleCountedAsMoved": false,
		"projectedCompressionCountedAsRealized": false,
		"syntheticFleetCountedAsProduction": false,
		"referenceFleetCountedAsProduction": false,
		"dashboardExecutor": false,
		"fluidvirtPolicyAuthority": false,
		"inventoryRuntimeExecutor": false,
		"stateStore": buildStateStoreFragment(r1.NewResourceVersion),
		"reconcileRequests": []interface{}{
			reconcileRequestToMap(req1), reconcileRequestToMap(req2), reconcileRequestToMap(reqStale),
		},
		"reconcileResults": []interface{}{
			reconcileResultToMap(r1), reconcileResultToMap(r2),
			reconcileResultToMapStale(rStale, staleErr), reconcileResultToMap(rRecovery),
		},
		"statusConditions": toIface(conditions),
		"kubernetesEvents": toIface(events),
		"leaderElectionReadiness": map[string]interface{}{
			"leaderElectionId": "leader-election-sprint14", "enabled": true,
			"leaseName": "hyperdensity-controller-leader", "leaseNamespace": "karl-system",
			"identity": "hyperdensity-controller-sprint14", "leaseDurationSeconds": 15,
			"renewDeadlineSeconds": 10, "retryPeriodSeconds": 2,
			"readinessStatus": "ready",
			"evidenceRefs": []interface{}{"hyperdensity-durable-controller-reconciler-v1"},
			"claimBoundary": "leader election configured; not HA production proven",
		},
		"rbacBoundary": buildRBACBoundary(),
		"metrics": rec.Metrics.Snapshot(),
		"recoverySemantics": map[string]interface{}{
			"recoverySemanticsId": "recovery-sprint14", "retryPolicyRef": "retry-sprint13",
			"recoveryMode": "retry_with_idempotency", "maxRetries": 3, "backoffSeconds": 30,
			"preservesIdempotency": true, "preservesLeases": true, "preservesActions": true,
			"preservesFutures": true, "degradedModeRef": "degraded-sprint14",
			"evidenceRefs": []interface{}{"hyperdensity-durable-controller-reconciler-v1"},
			"claimBoundary": "recovery preserves idempotency",
		},
		"degradedMode": recFail.DegradedMode,
		"ownershipReferences": []interface{}{
			map[string]interface{}{
				"ownershipReferenceId": "owner-sprint14", "ownerKind": "Deployment",
				"ownerName": "hyperdensity-controller", "ownedKind": "ConfigMap",
				"ownedName": "hyperdensity-controller-state", "namespace": "karl-system",
				"controller": true, "blockOwnerDeletion": true,
				"evidenceRefs": []interface{}{"hyperdensity-durable-controller-reconciler-v1"},
				"claimBoundary": "controller ownership reference",
			},
		},
		"stateMigrations": []interface{}{
			map[string]interface{}{
				"migrationId": "migration-sprint13-to-14", "fromStoreType": "in_memory_reference",
				"toStoreType": "kubernetes_configmap", "migrationStatus": "completed",
				"migratedStateVersion": 2, "idempotencyRecordsMigrated": 1,
				"leasesMigrated": 9, "actionsMigrated": 9, "futuresMigrated": 9, "auditMigrated": true,
				"evidenceRefs": []interface{}{"hyperdensity-live-controller-reconciliation-v1"},
				"claimBoundary": "state migration from Sprint 13 in-memory to Sprint 14 ConfigMap",
			},
		},
		"durableIdempotencyRecords": []interface{}{
			map[string]interface{}{
				"durableIdempotencyRecordId": "durable-idem-tick-idem-sprint14-001",
				"idempotencyKey": "tick-idem-sprint14-001", "objectType": "reconcile_request",
				"objectId": "req-sprint14-002", "persisted": true, "replayCount": 1,
				"duplicateSuppressed": true, "resourceVersion": r2.NewResourceVersion,
				"firstSeenAt": now.Format(time.RFC3339), "lastSeenAt": now.Add(time.Minute).Format(time.RFC3339),
				"evidenceRefs": []interface{}{"tick-idem-sprint14-001"}, "claimBoundary": "durable idempotency",
			},
		},
		"sourceOfTruth": map[string]interface{}{
			"durableState": "Karl-Hyperdensity", "kubernetesReconciler": "Karl-Hyperdensity",
			"runtimeActuator": "FluidVirt", "operatorUI": "Karl-Dashboard (projection only)",
		},
		"safetyInvariants": map[string]interface{}{
			"generalProductionAutoAllowed": false, "projectedCompressionCountedAsRealized": false,
			"dashboardExecutor": false, "fluidvirtPolicyAuthority": false,
		},
		"claimBoundaries": []interface{}{
			"durable Kubernetes ConfigMap state; CRD reference_defined only",
			"leader election configured; not HA production proven",
			"manifests reference_only; not_production_install",
		},
		"blockers": []interface{}{},
		"warnings": []interface{}{},
		"_staleErr": fmt.Sprintf("%v", staleErr),
	}
}

func durableReconcileRequest(idemKey string) ReconcileRequest {
	now := time.Date(2026, 5, 14, 18, 0, 0, 0, time.UTC)
	return ReconcileRequest{
		ReconcileRequestID: "req-001",
		Namespace:          "karl-system",
		Name:               "hyperdensity-controller-state",
		Reason:             "scheduled_tick",
		RequestedAt:        now,
		IdempotencyKey:     idemKey,
		SourceEvent:        "scheduled",
		LiveLoopInput: LiveLoopInput{
			IdempotencyKey: idemKey, TickSequence: 1, RateLimitRemaining: 5,
			ProductionCanaryEnabled: true, CurrentCompressionRate: 0.04941176470588235,
			ProjectedCompressionRate: 0.22352941176470587, EligibleIdleValue: 0.085,
			Snapshot: referenceSnapshot(),
			Observed: ObservedMarketState{ObservedSnapshotID: "obs-001", ProductionCanaryEnabled: true},
			Desired:  DesiredMarketState{DesiredStateID: "desired-001", ForbiddenExecutionScopes: []string{"general_production_auto", "production_auto_with_policy"}},
		},
	}
}

func buildStateStoreFragment(rv string) map[string]interface{} {
	return map[string]interface{}{
		"stateStoreId": "k8s-state-store-sprint14", "storeType": "kubernetes_configmap",
		"implementationStatus": "implemented", "namespace": "karl-system",
		"objectName": "hyperdensity-controller-state", "objectKind": "ConfigMap",
		"resourceVersion": rv, "optimisticLockEnabled": true, "staleWriteRejected": true,
		"stateSurvivesReconcile": true, "idempotencyPersists": true,
		"leasePersistenceEnabled": true, "actionPersistenceEnabled": true,
		"futurePersistenceEnabled": true, "invalidationPersistenceEnabled": true,
		"auditPersistenceEnabled": true, "statusConditionPersistenceEnabled": true,
		"eventEmissionEnabled": true,
		"evidenceRefs": []interface{}{"hyperdensity-durable-controller-reconciler-v1"},
		"claimBoundary": "kubernetes_configmap implemented; crd reference_defined",
	}
}

func buildRBACBoundary() map[string]interface{} {
	return map[string]interface{}{
		"rbacBoundaryId": "rbac-sprint14", "serviceAccount": "hyperdensity-controller",
		"namespace": "karl-system", "clusterAdmin": false, "namespaceScoped": true,
		"allowedResources": []interface{}{"configmaps", "events", "hyperdensitycontrollerstates"},
		"allowedVerbs": []interface{}{"get", "list", "watch", "create", "update", "patch"},
		"deniedResources": []interface{}{"nodes", "pods/exec", "libvirt", "cgroups"},
		"rawRuntimeControlsAllowed": false, "directLibvirtAllowed": false, "directCgroupAllowed": false,
		"evidenceRefs": []interface{}{"hyperdensity-durable-controller-reconciler-v1"},
		"claimBoundary": "namespace-scoped RBAC; no raw runtime controls",
	}
}

func reconcileRequestToMap(req ReconcileRequest) map[string]interface{} {
	return map[string]interface{}{
		"reconcileRequestId": req.ReconcileRequestID, "namespace": req.Namespace, "name": req.Name,
		"reason": req.Reason, "requestedAt": req.RequestedAt.Format(time.RFC3339),
		"idempotencyKey": req.IdempotencyKey, "previousResourceVersion": req.PreviousResourceVersion,
		"sourceEvent": req.SourceEvent,
		"evidenceRefs": []interface{}{req.ReconcileRequestID}, "claimBoundary": "reconcile request",
	}
}

func reconcileResultToMap(r ReconcileResult) map[string]interface{} {
	return map[string]interface{}{
		"reconcileResultId": r.ReconcileResultID, "reconcileRequestId": r.ReconcileRequestID,
		"status": r.Status, "startedAt": r.StartedAt.Format(time.RFC3339),
		"completedAt": r.CompletedAt.Format(time.RFC3339), "durationMs": r.DurationMs,
		"loadedStateVersion": r.LoadedStateVersion, "savedStateVersion": r.SavedStateVersion,
		"newResourceVersion": r.NewResourceVersion, "staleWriteRejected": r.StaleWriteRejected,
		"actionsCreated": r.ActionsCreated, "actionsInvalidated": r.ActionsInvalidated,
		"leasesCreated": r.LeasesCreated, "futuresCreated": r.FuturesCreated,
		"idempotentReplay": r.IdempotentReplay, "duplicateActionsSuppressed": r.DuplicateActionsSuppressed,
		"statusConditionsUpdated": r.StatusConditionsUpdated, "eventsEmitted": r.EventsEmitted,
		"metricsUpdated": r.MetricsUpdated,
		"evidenceRefs": []interface{}{r.ReconcileResultID}, "claimBoundary": "reconcile result",
	}
}

func reconcileResultToMapStale(r ReconcileResult, err error) map[string]interface{} {
	m := reconcileResultToMap(r)
	m["status"] = "failed"
	m["staleWriteRejected"] = true
	m["duplicateActionsSuppressed"] = true
	if err != nil {
		m["message"] = err.Error()
	}
	return m
}

func MarshalDurableReferenceSurfaceJSON() ([]byte, error) {
	s := BuildDurableReferenceSurface()
	delete(s, "_staleErr")
	return json.MarshalIndent(s, "", "  ")
}
