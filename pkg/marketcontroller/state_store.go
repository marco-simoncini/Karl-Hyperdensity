package marketcontroller

import (
	"fmt"
	"sync"
	"time"
)

const MilestoneLiveControllerReconciliation = "hyperdensity_live_controller_reconciliation_execution_loop_v1"

// ControllerState is the persisted live-loop state.
type ControllerState struct {
	StateVersion       int                      `json:"stateVersion"`
	OptimisticLock     int                      `json:"optimisticLockVersion"`
	Actions            []map[string]interface{} `json:"actions"`
	Futures            []map[string]interface{} `json:"futures"`
	Leases             []map[string]interface{} `json:"leases"`
	ActionLifecycles   []map[string]interface{} `json:"actionLifecycles"`
	FutureLifecycles   []map[string]interface{} `json:"futureLifecycles"`
	LeaseLifecycles    []map[string]interface{} `json:"leaseLifecycles"`
	Invalidations      []map[string]interface{} `json:"invalidations"`
	ExecutionSelections []map[string]interface{} `json:"executionSelections"`
	ExecutionHandoffs  []map[string]interface{} `json:"executionHandoffs"`
	PostExecutions     []map[string]interface{} `json:"postExecutions"`
	IdempotencyRecords []map[string]interface{} `json:"idempotencyRecords"`
	AuditEvents        []map[string]interface{} `json:"auditEvents"`
	LoadedAt           time.Time                `json:"loadedAt"`
	SavedAt            time.Time                `json:"savedAt"`
}

// StateStore abstracts controller persistence.
type StateStore interface {
	LoadControllerState() (ControllerState, error)
	SaveControllerState(state ControllerState) error
	AppendAuditEvent(event map[string]interface{}) error
	SaveLease(lease map[string]interface{}) error
	SaveActionLifecycle(lc map[string]interface{}) error
	SaveFutureLifecycle(lc map[string]interface{}) error
	SaveInvalidation(inv map[string]interface{}) error
}

// InMemoryStateStore is a deterministic test/reference store.
type InMemoryStateStore struct {
	mu    sync.Mutex
	state ControllerState
}

func NewInMemoryStateStore() *InMemoryStateStore {
	return &InMemoryStateStore{state: ControllerState{StateVersion: 1}}
}

func (s *InMemoryStateStore) LoadControllerState() (ControllerState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state.LoadedAt.IsZero() {
		s.state.LoadedAt = time.Now().UTC()
	}
	return s.state, nil
}

func (s *InMemoryStateStore) SaveControllerState(state ControllerState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	state.SavedAt = time.Now().UTC()
	state.OptimisticLock++
	s.state = state
	return nil
}

func (s *InMemoryStateStore) AppendAuditEvent(event map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.AuditEvents = append(s.state.AuditEvents, event)
	return nil
}

func (s *InMemoryStateStore) SaveLease(lease map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Leases = appendOrReplace(s.state.Leases, lease, "leaseId")
	s.state.LeaseLifecycles = appendOrReplace(s.state.LeaseLifecycles, lease, "leaseId")
	return nil
}

func (s *InMemoryStateStore) SaveActionLifecycle(lc map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.ActionLifecycles = appendOrReplace(s.state.ActionLifecycles, lc, "actionId")
	return nil
}

func (s *InMemoryStateStore) SaveFutureLifecycle(lc map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.FutureLifecycles = appendOrReplace(s.state.FutureLifecycles, lc, "futureId")
	return nil
}

func (s *InMemoryStateStore) SaveInvalidation(inv map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Invalidations = append(s.state.Invalidations, inv)
	return nil
}

func appendOrReplace(items []map[string]interface{}, item map[string]interface{}, key string) []map[string]interface{} {
	id, _ := item[key].(string)
	if id == "" {
		return append(items, item)
	}
	for i, existing := range items {
		if strOr(existing[key]) == id {
			items[i] = item
			return items
		}
	}
	return append(items, item)
}

// BuildStateStoreRecord returns a reference state-store surface fragment.
func BuildStateStoreRecord() map[string]interface{} {
	now := time.Now().UTC().Format(time.RFC3339)
	return map[string]interface{}{
		"stateStoreId":          "state-store-sprint13",
		"storeType":             "in_memory_reference",
		"persistenceMode":       "ephemeral_reference",
		"loadedAt":              now,
		"savedAt":               now,
		"previousStateRef":      "state-v0",
		"currentStateRef":       "state-v1",
		"stateVersion":          1,
		"optimisticLockVersion": 1,
		"auditAppendOnly":       true,
		"idempotencyEnabled":    true,
		"leaseTtlEnforced":      true,
		"futureTtlEnforced":     true,
		"actionTtlEnforced":     true,
		"evidenceRefs":          []interface{}{"hyperdensity-live-controller-reconciliation-v1"},
		"claimBoundary":         "in_memory_reference only; not production persistence",
	}
}

func auditEvent(eventType, liveLoopID string, refs ...string) map[string]interface{} {
	ev := map[string]interface{}{
		"eventType":   eventType,
		"occurredAt":  time.Now().UTC().Format(time.RFC3339),
		"liveLoopId":  liveLoopID,
		"evidenceRefs": []interface{}{"hyperdensity-live-controller-reconciliation-v1"},
		"claimBoundary": "immutable audit append-only",
	}
	if len(refs) > 0 {
		ev["objectRef"] = refs[0]
	}
	return ev
}

func requireStore(store StateStore) error {
	if store == nil {
		return fmt.Errorf("state store required")
	}
	return nil
}
