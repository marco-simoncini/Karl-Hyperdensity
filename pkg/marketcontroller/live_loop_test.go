package marketcontroller

import (
	"testing"
	"time"
)

func testLiveInput() LiveLoopInput {
	return LiveLoopInput{
		IdempotencyKey:          "tick-idem-sprint13-001",
		TickSequence:            1,
		RateLimitRemaining:      5,
		ProductionCanaryEnabled: true,
		CurrentCompressionRate:    0.04941176470588235,
		ProjectedCompressionRate: 0.22352941176470587,
		EligibleIdleValue:       0.085,
		GuaranteedEligibleTotal: 0.00336,
		Snapshot:                testSnapshot(),
		Observed: ObservedMarketState{
			ObservedSnapshotID:      "observed-sprint13-001",
			CollectedAt:             time.Now().UTC(),
			ShellCount:              12,
			ProductionCanaryEnabled: true,
		},
		Desired: DesiredMarketState{
			DesiredStateID:           "desired-sprint13",
			ForbiddenExecutionScopes: []string{"general_production_auto", "production_auto_with_policy"},
		},
	}
}

func TestStateStoreLoadSave(t *testing.T) {
	store := NewInMemoryStateStore()
	state, err := store.LoadControllerState()
	if err != nil {
		t.Fatal(err)
	}
	state.StateVersion = 2
	if err := store.SaveControllerState(state); err != nil {
		t.Fatal(err)
	}
	loaded, _ := store.LoadControllerState()
	if loaded.StateVersion != 2 {
		t.Fatal("state not saved")
	}
}

func TestIdempotencySuppressesDuplicate(t *testing.T) {
	store := NewInMemoryStateStore()
	loop := NewLiveControllerLoop(store, "live-loop-sprint13", "ctrl-sprint13")
	input := testLiveInput()
	r1, err := loop.RunScheduledTick(input, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if r1.TickStatus != "completed" {
		t.Fatalf("expected completed got %s", r1.TickStatus)
	}
	r2, err := loop.RunScheduledTick(input, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if !r2.IdempotentReplay {
		t.Fatal("expected idempotent replay")
	}
	if r2.TickStatus != "completed_idempotent_replay" {
		t.Fatal("expected completed_idempotent_replay")
	}
}

func TestLeaseTTLExpires(t *testing.T) {
	state := ControllerState{}
	now := time.Now().UTC()
	lease := newLeaseLifecycle("lease-ttl-1", "action-1", "d", "r", "cpu", "100m", now.Add(-2*time.Hour), 60)
	state.LeaseLifecycles = []map[string]interface{}{lease}
	exp, _, _ := RefreshLifecycles(&state, now)
	if exp != 1 {
		t.Fatal("lease should expire")
	}
	if !boolOr(state.LeaseLifecycles[0]["expired"]) {
		t.Fatal("lease expired flag")
	}
}

func TestFutureExpires(t *testing.T) {
	state := ControllerState{}
	now := time.Now().UTC()
	flc := map[string]interface{}{
		"futureId":  "f-1",
		"expiresAt": now.Add(-time.Hour).Format(time.RFC3339),
		"expired":   false,
		"state":     "refreshed",
	}
	state.FutureLifecycles = []map[string]interface{}{flc}
	_, exp, _ := RefreshLifecycles(&state, now)
	if exp != 1 {
		t.Fatal("future should expire")
	}
}

func TestKillSwitchBlocksSelection(t *testing.T) {
	input := testLiveInput()
	input.KillSwitchActive = true
	_, _, post := SelectExecutableActions(nil, []map[string]interface{}{
		{"actionId": "a1", "donorShellId": "shell-container-donor-a", "executionScopeRecommendation": "operator_controlled"},
	}, input, time.Now().UTC())
	for _, p := range post {
		if boolOr(p["realizedMovementKept"]) {
			t.Fatal("kill switch must block realized movement")
		}
	}
}

func TestProjectedNotCountedAsRealized(t *testing.T) {
	tracker := TrackRealizedCompression(testLiveInput(), []map[string]interface{}{
		{"realizedMovementKept": false, "realizedIdleValue": 0.0},
	}, time.Now().UTC())
	if boolOr(tracker["projectedCompressionCountedAsRealized"]) {
		t.Fatal("projected must not count as realized")
	}
}

func TestRealizedRequiresMutationEvidence(t *testing.T) {
	p := buildPostExecution("a1", "h1", "shell-container-donor-a", map[string]interface{}{"expectedMovedIdleValue": 0.0042}, true)
	if !boolOr(p["mutationObserved"]) || !boolOr(p["postVerifyPassed"]) || !boolOr(p["realizedMovementKept"]) {
		t.Fatal("evidence-backed movement should be realized kept")
	}
	p2 := buildPostExecution("a2", "h2", "shell-rollback-test", map[string]interface{}{}, true)
	if boolOr(p2["realizedMovementKept"]) {
		t.Fatal("rollback-required must not count as realized kept")
	}
}

func TestProductionCanaryRequiresScope(t *testing.T) {
	input := testLiveInput()
	input.ProductionCanaryEnabled = false
	sel, _, _ := SelectExecutableActions(nil, []map[string]interface{}{
		{"actionId": "a1", "donorShellId": "shell-container-donor-a", "executionScopeRecommendation": "production_canary_eligible"},
	}, input, time.Now().UTC())
	for _, s := range sel {
		if strOr(s["executionScope"]) == "production_canary_auto" && boolOr(s["selected"]) {
			t.Fatal("canary without productionCanaryEnabled must not select")
		}
	}
}

func TestWindowsRemediationNotAutoSelected(t *testing.T) {
	input := testLiveInput()
	sel, _, _ := SelectExecutableActions(nil, []map[string]interface{}{
		{"actionId": "w1", "donorShellId": "shell-windows-hyper", "executionScopeRecommendation": "remediation_only"},
	}, input, time.Now().UTC())
	if boolOr(sel[0]["selected"]) {
		t.Fatal("windows must not be auto selected")
	}
}

func TestRunScheduledTickAuditTrail(t *testing.T) {
	store := NewInMemoryStateStore()
	loop := NewLiveControllerLoop(store, "live-loop-sprint13", "ctrl-sprint13")
	input := testLiveInput()
	input.IdempotencyKey = "audit-tick-unique"
	res, err := loop.RunScheduledTick(input, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	types := map[string]bool{}
	for _, ev := range res.State.AuditEvents {
		types[strOr(ev["eventType"])] = true
	}
	required := []string{"state_loaded", "observed_snapshot_collected", "desired_state_computed",
		"reconciliation_diff_computed", "lease_created", "action_lifecycle_updated", "future_refreshed",
		"execution_selected", "execution_handed_off", "compression_tracker_updated", "state_saved"}
	for _, r := range required {
		if !types[r] {
			t.Fatalf("missing audit event %s", r)
		}
	}
}

func TestForbiddenGeneralProductionAuto(t *testing.T) {
	if isPermittedScope("general_production_auto", testLiveInput()) {
		t.Fatal("general production auto forbidden")
	}
	if isPermittedScope("production_auto_with_policy", testLiveInput()) {
		t.Fatal("production_auto_with_policy forbidden")
	}
}
