package marketcontroller

import (
	"context"
	"fmt"
	"time"
)

const MilestoneDurableControllerKubernetesReconciler = "hyperdensity_durable_controller_state_kubernetes_reconciler_v1"

// DurableStateKey identifies a persisted controller state object.
type DurableStateKey struct {
	Namespace string
	Name      string
	Key       string
}

// DurablePersistedState is the serializable durable controller blob.
type DurablePersistedState struct {
	StateVersion        int                      `json:"stateVersion"`
	ControllerState     ControllerState          `json:"controllerState"`
	IdempotencyRecords  []map[string]interface{} `json:"idempotencyRecords"`
	LeaseLifecycles     []map[string]interface{} `json:"leaseLifecycles"`
	ActionLifecycles    []map[string]interface{} `json:"actionLifecycles"`
	FutureLifecycles    []map[string]interface{} `json:"futureLifecycles"`
	Invalidations       []map[string]interface{} `json:"invalidations"`
	AuditEvents         []map[string]interface{} `json:"auditEvents"`
	StatusConditions    []map[string]interface{} `json:"statusConditions"`
}

// DurableStateStore abstracts Kubernetes-backed persistence.
type DurableStateStore interface {
	Load(ctx context.Context, key DurableStateKey) (DurablePersistedState, string, error)
	Save(ctx context.Context, key DurableStateKey, state DurablePersistedState, resourceVersion string) (string, error)
	AppendAudit(ctx context.Context, key DurableStateKey, event map[string]interface{}) error
	SaveLease(ctx context.Context, lease map[string]interface{}) error
	SaveAction(ctx context.Context, action map[string]interface{}) error
	SaveFuture(ctx context.Context, future map[string]interface{}) error
	SaveInvalidation(ctx context.Context, invalidation map[string]interface{}) error
	SaveIdempotencyRecord(ctx context.Context, record map[string]interface{}) error
	LoadIdempotencyRecord(ctx context.Context, idempotencyKey string) (map[string]interface{}, error)
	SetStatusCondition(ctx context.Context, condition map[string]interface{}) error
	EmitEvent(ctx context.Context, event map[string]interface{}) error
}

// ErrStaleWrite indicates optimistic-lock / resourceVersion conflict.
var ErrStaleWrite = fmt.Errorf("stale resourceVersion write rejected")

// DefaultDurableStateKey is the Sprint 14 reference key.
func DefaultDurableStateKey() DurableStateKey {
	return DurableStateKey{Namespace: "karl-system", Name: "hyperdensity-controller-state", Key: "state.json"}
}

func newEmptyDurableState() DurablePersistedState {
	return DurablePersistedState{
		StateVersion: 1,
		ControllerState: ControllerState{StateVersion: 1},
	}
}

func durableStateFromController(state ControllerState) DurablePersistedState {
	return DurablePersistedState{
		StateVersion:       state.StateVersion,
		ControllerState:    state,
		IdempotencyRecords: state.IdempotencyRecords,
		LeaseLifecycles:    state.LeaseLifecycles,
		ActionLifecycles:   state.ActionLifecycles,
		FutureLifecycles:   state.FutureLifecycles,
		Invalidations:      state.Invalidations,
		AuditEvents:        state.AuditEvents,
	}
}

func controllerFromDurable(ds DurablePersistedState) ControllerState {
	cs := ds.ControllerState
	cs.StateVersion = ds.StateVersion
	cs.IdempotencyRecords = ds.IdempotencyRecords
	cs.LeaseLifecycles = ds.LeaseLifecycles
	cs.ActionLifecycles = ds.ActionLifecycles
	cs.FutureLifecycles = ds.FutureLifecycles
	cs.Invalidations = ds.Invalidations
	cs.AuditEvents = ds.AuditEvents
	return cs
}

func nowRFC3339() string { return time.Now().UTC().Format(time.RFC3339) }
