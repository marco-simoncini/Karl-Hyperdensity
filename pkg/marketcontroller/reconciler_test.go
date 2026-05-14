package marketcontroller

import (
	"context"
	"testing"
	"time"
)

func testReconcileRequest(idemKey string) ReconcileRequest {
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

func TestStateSurvivesReconcile(t *testing.T) {
	client := NewFakeKubernetesClient()
	store := NewKubernetesStateStore(client, DefaultDurableStateKey())
	rec := NewReconciler(store, "live-loop-durable", "ctrl-durable")
	ctx := context.Background()
	r1, err := rec.Reconcile(ctx, testReconcileRequest("idem-survive-001"))
	if err != nil {
		t.Fatal(err)
	}
	if r1.NewResourceVersion == "" {
		t.Fatal("expected new resourceVersion")
	}
	ds, rv, _ := store.Load(ctx, DefaultDurableStateKey())
	if ds.StateVersion < 1 || rv == "" {
		t.Fatal("state must survive")
	}
}

func TestIdempotencyPersistsAcrossReconcile(t *testing.T) {
	client := NewFakeKubernetesClient()
	store := NewKubernetesStateStore(client, DefaultDurableStateKey())
	rec := NewReconciler(store, "live-loop-durable", "ctrl-durable")
	ctx := context.Background()
	key := "idem-persist-001"
	_, _ = rec.Reconcile(ctx, testReconcileRequest(key))
	r2, _ := rec.Reconcile(ctx, testReconcileRequest(key))
	if !r2.IdempotentReplay && r2.Status != "idempotent_replay" {
		t.Fatal("expected idempotent replay on second reconcile")
	}
	rec2, _ := store.LoadIdempotencyRecord(ctx, key)
	if rec2 == nil {
		t.Fatal("idempotency record must persist")
	}
}

func TestStaleResourceVersionRejected(t *testing.T) {
	client := NewFakeKubernetesClient()
	store := NewKubernetesStateStore(client, DefaultDurableStateKey())
	rec := NewReconciler(store, "live-loop-durable", "ctrl-durable")
	ctx := context.Background()
	req := testReconcileRequest("idem-stale-first")
	_, _ = rec.Reconcile(ctx, req)
	req2 := testReconcileRequest("idem-stale-second")
	req2.PreviousResourceVersion = "1"
	_, err := rec.Reconcile(ctx, req2)
	if err != ErrStaleWrite {
		t.Fatalf("expected ErrStaleWrite got %v", err)
	}
}

func TestReadyConditionOnSuccess(t *testing.T) {
	client := NewFakeKubernetesClient()
	store := NewKubernetesStateStore(client, DefaultDurableStateKey())
	rec := NewReconciler(store, "live-loop-durable", "ctrl-durable")
	_, _ = rec.Reconcile(context.Background(), testReconcileRequest("idem-ready-001"))
	found := false
	for _, c := range client.GetConditions() {
		if strOr(c["type"]) == ConditionReady && strOr(c["status"]) == "True" {
			found = true
		}
	}
	if !found {
		t.Fatal("Ready=True expected")
	}
}

func TestDegradedOnSaveFailure(t *testing.T) {
	client := NewFakeKubernetesClient()
	client.failSave = true
	store := NewKubernetesStateStore(client, DefaultDurableStateKey())
	rec := NewReconciler(store, "live-loop-durable", "ctrl-durable")
	_, err := rec.Reconcile(context.Background(), testReconcileRequest("idem-fail-001"))
	if err == nil {
		t.Fatal("expected save failure")
	}
	if !boolOr(rec.DegradedMode["active"]) || !boolOr(rec.DegradedMode["failClosed"]) {
		t.Fatal("degraded fail-closed expected")
	}
}

func TestEventsEmitted(t *testing.T) {
	client := NewFakeKubernetesClient()
	store := NewKubernetesStateStore(client, DefaultDurableStateKey())
	rec := NewReconciler(store, "live-loop-durable", "ctrl-durable")
	_, _ = rec.Reconcile(context.Background(), testReconcileRequest("idem-events-001"))
	reasons := map[string]bool{}
	for _, e := range client.GetEvents() {
		reasons[strOr(e["reason"])] = true
	}
	for _, want := range []string{"ControllerTickCompleted", "StatePersisted"} {
		if !reasons[want] {
			t.Fatalf("missing event %s", want)
		}
	}
}

func TestMetricsForbiddenAutoZero(t *testing.T) {
	m := NewMetricsCollector()
	snap := m.Snapshot()
	if floatOr(snap["generalProductionAutoEnabledGauge"]) != 0 || floatOr(snap["productionAutoWithPolicyEnabledGauge"]) != 0 {
		t.Fatal("forbidden auto gauges must be 0")
	}
}

func TestRecoveryRetryPreservesIdempotency(t *testing.T) {
	client := NewFakeKubernetesClient()
	store := NewKubernetesStateStore(client, DefaultDurableStateKey())
	rec := NewReconciler(store, "live-loop-durable", "ctrl-durable")
	ctx := context.Background()
	key := "idem-recovery-001"
	client.failSave = true
	_, _ = rec.Reconcile(ctx, testReconcileRequest(key))
	client.failSave = false
	r2, err := rec.Reconcile(ctx, testReconcileRequest(key))
	if err != nil {
		t.Fatal(err)
	}
	if r2.Status != "success" && r2.Status != "idempotent_replay" {
		t.Fatalf("recovery expected success got %s", r2.Status)
	}
}
