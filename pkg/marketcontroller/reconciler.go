package marketcontroller

import (
	"context"
	"fmt"
	"time"
)

// ReconcileRequest is a Kubernetes-style reconcile request.
type ReconcileRequest struct {
	ReconcileRequestID      string
	Namespace               string
	Name                    string
	Reason                  string
	RequestedAt             time.Time
	IdempotencyKey          string
	PreviousResourceVersion string
	SourceEvent             string
	LiveLoopInput           LiveLoopInput
}

// ReconcileResult is the outcome of one reconcile.
type ReconcileResult struct {
	ReconcileResultID        string
	ReconcileRequestID       string
	Status                   string
	StartedAt                time.Time
	CompletedAt              time.Time
	DurationMs               int64
	LoadedStateVersion       int
	SavedStateVersion        int
	NewResourceVersion       string
	StaleWriteRejected       bool
	ActionsCreated           int
	ActionsUpdated           int
	ActionsInvalidated       int
	LeasesCreated            int
	LeasesExpired            int
	FuturesCreated           int
	FuturesExpired           int
	IdempotentReplay         bool
	DuplicateActionsSuppressed bool
	StatusConditionsUpdated  int
	EventsEmitted            int
	MetricsUpdated           bool
	Requeue                  bool
	RequeueAfterSeconds      int
	RealizedCompressionRate  float64
	ProjectedCompressionRate float64
	Degraded                 bool
	DegradedReason           string
}

// Reconciler runs durable Kubernetes reconciliation.
type Reconciler struct {
	Store         DurableStateStore
	StateKey      DurableStateKey
	LiveLoopID    string
	ControllerID  string
	Metrics       *MetricsCollector
	DegradedMode  map[string]interface{}
}

func NewReconciler(store DurableStateStore, liveLoopID, controllerID string) *Reconciler {
	return &Reconciler{
		Store:        store,
		StateKey:     DefaultDurableStateKey(),
		LiveLoopID:   liveLoopID,
		ControllerID: controllerID,
		Metrics:      NewMetricsCollector(),
		DegradedMode: map[string]interface{}{"active": false, "failClosed": false, "executionSelectionDisabled": false},
	}
}

// Reconcile executes one durable reconcile cycle.
func (r *Reconciler) Reconcile(ctx context.Context, req ReconcileRequest) (ReconcileResult, error) {
	started := time.Now().UTC()
	result := ReconcileResult{
		ReconcileResultID:  fmt.Sprintf("reconcile-result-%s", req.ReconcileRequestID),
		ReconcileRequestID: req.ReconcileRequestID,
		StartedAt:          started,
	}

	_ = r.Store.SetStatusCondition(ctx, newStatusCondition(ConditionReconciling, "True", "Reconciling", "reconcile in progress"))

	ds, rv, err := r.Store.Load(ctx, r.StateKey)
	if err != nil {
		return r.failDegraded(ctx, result, started, "StateLoadFailed", err)
	}
	result.LoadedStateVersion = ds.StateVersion

	// Check durable idempotency — record already persisted for this key
	if rec, _ := r.Store.LoadIdempotencyRecord(ctx, req.IdempotencyKey); rec != nil {
		replay := intOr(rec["replayCount"]) + 1
		rec["replayCount"] = replay
		rec["lastSeenAt"] = time.Now().UTC().Format(time.RFC3339)
		rec["duplicateSuppressed"] = true
		_ = r.Store.SaveIdempotencyRecord(ctx, rec)
		result.Status = "idempotent_replay"
		result.IdempotentReplay = true
		result.DuplicateActionsSuppressed = true
		result.CompletedAt = time.Now().UTC()
		result.DurationMs = result.CompletedAt.Sub(started).Milliseconds()
		_ = r.emitEvent(ctx, "Normal", "IdempotentReplaySuppressed", "duplicate reconcile suppressed", req.Name)
		_ = r.Store.SetStatusCondition(ctx, newStatusCondition(ConditionReady, "True", "IdempotentReplay", "idempotent replay suppressed"))
		r.Metrics.RecordReconcile(true, int(result.DurationMs), result)
		result.MetricsUpdated = true
		return result, nil
	}

	// Bridge to Sprint 13 live loop via in-memory adapter
	mem := controllerFromDurable(ds)
	memStore := &memoryBridgeStore{state: mem}
	loop := NewLiveControllerLoop(memStore, r.LiveLoopID, r.ControllerID)
	tickTime := req.RequestedAt
	if tickTime.IsZero() {
		tickTime = started
	}
	liveRes, err := loop.RunScheduledTick(req.LiveLoopInput, tickTime)
	if err != nil {
		return r.failDegraded(ctx, result, started, "LiveLoopFailed", err)
	}

	if liveRes.IdempotentReplay {
		result.Status = "idempotent_replay"
		result.IdempotentReplay = true
		result.DuplicateActionsSuppressed = true
	} else {
		result.Status = "success"
		result.ActionsCreated = liveRes.ReconciliationDiff.ActionsCreated
		result.ActionsInvalidated = liveRes.ReconciliationDiff.ActionsInvalidated
		result.LeasesCreated = liveRes.ReconciliationDiff.LeasesCreated
		result.LeasesExpired = liveRes.ReconciliationDiff.LeasesExpired
		result.FuturesCreated = liveRes.ReconciliationDiff.FuturesCreated
		result.FuturesExpired = liveRes.ReconciliationDiff.FuturesExpired
	}

	newDS := durableStateFromController(liveRes.State)
	newDS.StateVersion = ds.StateVersion + 1
	newDS.StatusConditions = setReadyConditions(baseStatusConditions())
	result.SavedStateVersion = newDS.StateVersion

	useRV := rv
	if req.PreviousResourceVersion != "" {
		useRV = req.PreviousResourceVersion
	}
	newRV, err := r.Store.Save(ctx, r.StateKey, newDS, useRV)
	if err == ErrStaleWrite {
		result.StaleWriteRejected = true
		result.Status = "failed"
		_ = r.emitEvent(ctx, "Warning", "ReconcileDegraded", "stale write rejected", req.Name)
		result.CompletedAt = time.Now().UTC()
		result.DurationMs = result.CompletedAt.Sub(started).Milliseconds()
		r.Metrics.RecordReconcile(false, int(result.DurationMs), result)
		result.MetricsUpdated = true
		return result, ErrStaleWrite
	}
	if err != nil {
		return r.failDegraded(ctx, result, started, "StateSaveFailed", err)
	}
	result.NewResourceVersion = newRV

	idemRec := map[string]interface{}{
		"durableIdempotencyRecordId": fmt.Sprintf("durable-idem-%s", req.IdempotencyKey),
		"idempotencyKey": req.IdempotencyKey, "objectType": "reconcile_request",
		"objectId": req.ReconcileRequestID, "persisted": true,
		"firstSeenAt": started.Format(time.RFC3339), "lastSeenAt": time.Now().UTC().Format(time.RFC3339),
		"replayCount": 0, "duplicateSuppressed": liveRes.IdempotentReplay,
		"resourceVersion": newRV,
		"evidenceRefs": []interface{}{req.IdempotencyKey}, "claimBoundary": "durable idempotency",
	}
	_ = r.Store.SaveIdempotencyRecord(ctx, idemRec)

	for _, lc := range liveRes.State.LeaseLifecycles {
		_ = r.Store.SaveLease(ctx, lc)
	}
	for _, lc := range liveRes.State.ActionLifecycles {
		_ = r.Store.SaveAction(ctx, lc)
	}
	for _, lc := range liveRes.State.FutureLifecycles {
		_ = r.Store.SaveFuture(ctx, lc)
	}

	for _, cond := range newDS.StatusConditions {
		_ = r.Store.SetStatusCondition(ctx, cond)
		result.StatusConditionsUpdated++
	}

	_ = r.emitEvent(ctx, "Normal", "ControllerTickCompleted", "controller tick completed", req.Name)
	_ = r.emitEvent(ctx, "Normal", "StatePersisted", "durable state persisted", req.Name)
	if result.ActionsInvalidated > 0 {
		_ = r.emitEvent(ctx, "Normal", "ActionInvalidated", "stale action invalidated", req.Name)
	}
	if result.FuturesExpired > 0 {
		_ = r.emitEvent(ctx, "Normal", "FutureExpired", "future expired", req.Name)
	}
	if result.LeasesExpired > 0 {
		_ = r.emitEvent(ctx, "Normal", "LeaseExpired", "lease expired", req.Name)
	}

	if tracker := liveRes.RealizedCompressionTracker; tracker != nil {
		result.RealizedCompressionRate = floatOr(tracker["currentRealizedIdleCompressionRate"])
		result.ProjectedCompressionRate = floatOr(tracker["projectedIdleCompressionRate"])
	}

	result.CompletedAt = time.Now().UTC()
	result.DurationMs = result.CompletedAt.Sub(started).Milliseconds()
	r.Metrics.RecordReconcile(result.Status == "success" || result.Status == "idempotent_replay", int(result.DurationMs), result)
	result.MetricsUpdated = true
	return result, nil
}

func (r *Reconciler) failDegraded(ctx context.Context, result ReconcileResult, started time.Time, reason string, err error) (ReconcileResult, error) {
	result.Status = "degraded"
	result.Degraded = true
	result.DegradedReason = reason
	result.CompletedAt = time.Now().UTC()
	result.DurationMs = result.CompletedAt.Sub(started).Milliseconds()
	r.DegradedMode = map[string]interface{}{
		"degradedModeId": "degraded-sprint14", "active": true, "reason": reason,
		"failClosed": true, "executionSelectionDisabled": true,
		"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
		"statusConditionRef": ConditionDegraded, "eventRef": "ReconcileDegraded",
		"evidenceRefs": []interface{}{"hyperdensity-durable-controller-reconciler-v1"},
		"claimBoundary": "degraded fail-closed mode",
	}
	for _, cond := range setDegradedConditions(baseStatusConditions(), reason, err.Error()) {
		_ = r.Store.SetStatusCondition(ctx, cond)
		result.StatusConditionsUpdated++
	}
	_ = r.emitEvent(ctx, "Warning", "ReconcileDegraded", err.Error(), result.ReconcileRequestID)
	r.Metrics.RecordReconcile(false, int(result.DurationMs), result)
	result.MetricsUpdated = true
	return result, err
}

func (r *Reconciler) emitEvent(ctx context.Context, eventType, reason, message, objectName string) error {
	ev := map[string]interface{}{
		"eventId": fmt.Sprintf("evt-%s-%d", reason, time.Now().UnixNano()),
		"type": eventType, "reason": reason, "message": message,
		"involvedObjectKind": "ConfigMap", "involvedObjectName": objectName,
		"namespace": r.StateKey.Namespace, "emittedAt": nowRFC3339(),
		"sourceComponent": "hyperdensity-controller",
		"evidenceRefs": []interface{}{"hyperdensity-durable-controller-reconciler-v1"},
		"claimBoundary": "kubernetes event",
	}
	return r.Store.EmitEvent(ctx, ev)
}

// memoryBridgeStore adapts ControllerState to StateStore for Sprint 13 loop.
type memoryBridgeStore struct {
	state ControllerState
}

func (m *memoryBridgeStore) LoadControllerState() (ControllerState, error) { return m.state, nil }
func (m *memoryBridgeStore) SaveControllerState(state ControllerState) error { m.state = state; return nil }
func (m *memoryBridgeStore) AppendAuditEvent(event map[string]interface{}) error {
	m.state.AuditEvents = append(m.state.AuditEvents, event)
	return nil
}
func (m *memoryBridgeStore) SaveLease(lease map[string]interface{}) error {
	m.state.LeaseLifecycles = appendOrReplace(m.state.LeaseLifecycles, lease, "leaseId")
	return nil
}
func (m *memoryBridgeStore) SaveActionLifecycle(lc map[string]interface{}) error {
	m.state.ActionLifecycles = appendOrReplace(m.state.ActionLifecycles, lc, "actionId")
	return nil
}
func (m *memoryBridgeStore) SaveFutureLifecycle(lc map[string]interface{}) error {
	m.state.FutureLifecycles = appendOrReplace(m.state.FutureLifecycles, lc, "futureId")
	return nil
}
func (m *memoryBridgeStore) SaveInvalidation(inv map[string]interface{}) error {
	m.state.Invalidations = append(m.state.Invalidations, inv)
	return nil
}
