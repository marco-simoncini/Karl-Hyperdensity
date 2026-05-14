package marketcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

// FakeConfigMap models a Kubernetes ConfigMap with resourceVersion.
type FakeConfigMap struct {
	Namespace       string
	Name            string
	Data            map[string]string
	ResourceVersion string
}

// FakeKubernetesClient is an in-memory fake Kubernetes API for tests.
type FakeKubernetesClient struct {
	mu         sync.Mutex
	configMaps map[string]*FakeConfigMap
	events     []map[string]interface{}
	conditions []map[string]interface{}
	leases     []map[string]interface{}
	actions    []map[string]interface{}
	futures    []map[string]interface{}
	invalids   []map[string]interface{}
	idempotency map[string]map[string]interface{}
	nextRV     int
	failSave   bool
	failLoad   bool
}

func NewFakeKubernetesClient() *FakeKubernetesClient {
	return &FakeKubernetesClient{
		configMaps:  map[string]*FakeConfigMap{},
		idempotency: map[string]map[string]interface{}{},
		nextRV:      1,
	}
}

func (c *FakeKubernetesClient) cmKey(ns, name string) string {
	return ns + "/" + name
}

// KubernetesStateStore is a ConfigMap-backed durable store using FakeKubernetesClient.
type KubernetesStateStore struct {
	Client *FakeKubernetesClient
	Key    DurableStateKey
}

func NewKubernetesStateStore(client *FakeKubernetesClient, key DurableStateKey) *KubernetesStateStore {
	return &KubernetesStateStore{Client: client, Key: key}
}

func (s *KubernetesStateStore) Load(ctx context.Context, key DurableStateKey) (DurablePersistedState, string, error) {
	_ = ctx
	if s.Client.failLoad {
		return DurablePersistedState{}, "", fmt.Errorf("fake load failure")
	}
	s.Client.mu.Lock()
	defer s.Client.mu.Unlock()
	cm, ok := s.Client.configMaps[s.Client.cmKey(key.Namespace, key.Name)]
	if !ok || cm.Data == nil || cm.Data[key.Key] == "" {
		return newEmptyDurableState(), "0", nil
	}
	var ds DurablePersistedState
	if err := json.Unmarshal([]byte(cm.Data[key.Key]), &ds); err != nil {
		return DurablePersistedState{}, "", err
	}
	return ds, cm.ResourceVersion, nil
}

func (s *KubernetesStateStore) Save(ctx context.Context, key DurableStateKey, state DurablePersistedState, resourceVersion string) (string, error) {
	_ = ctx
	if s.Client.failSave {
		return "", fmt.Errorf("fake save failure")
	}
	s.Client.mu.Lock()
	defer s.Client.mu.Unlock()
	ck := s.Client.cmKey(key.Namespace, key.Name)
	cm, ok := s.Client.configMaps[ck]
	if ok && resourceVersion != "" && cm.ResourceVersion != resourceVersion {
		return "", ErrStaleWrite
	}
	raw, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	s.Client.nextRV++
	newRV := strconv.Itoa(s.Client.nextRV)
	if !ok {
		cm = &FakeConfigMap{Namespace: key.Namespace, Name: key.Name, Data: map[string]string{}}
		s.Client.configMaps[ck] = cm
	}
	if cm.Data == nil {
		cm.Data = map[string]string{}
	}
	cm.Data[key.Key] = string(raw)
	cm.ResourceVersion = newRV
	return newRV, nil
}

func (s *KubernetesStateStore) AppendAudit(ctx context.Context, key DurableStateKey, event map[string]interface{}) error {
	ds, rv, err := s.Load(ctx, key)
	if err != nil {
		return err
	}
	ds.AuditEvents = append(ds.AuditEvents, event)
	_, err = s.Save(ctx, key, ds, rv)
	return err
}

func (s *KubernetesStateStore) SaveLease(ctx context.Context, lease map[string]interface{}) error {
	s.Client.mu.Lock()
	defer s.Client.mu.Unlock()
	s.Client.leases = appendOrReplace(s.Client.leases, lease, "leaseId")
	return nil
}

func (s *KubernetesStateStore) SaveAction(ctx context.Context, action map[string]interface{}) error {
	s.Client.mu.Lock()
	defer s.Client.mu.Unlock()
	s.Client.actions = appendOrReplace(s.Client.actions, action, "actionId")
	return nil
}

func (s *KubernetesStateStore) SaveFuture(ctx context.Context, future map[string]interface{}) error {
	s.Client.mu.Lock()
	defer s.Client.mu.Unlock()
	s.Client.futures = appendOrReplace(s.Client.futures, future, "futureId")
	return nil
}

func (s *KubernetesStateStore) SaveInvalidation(ctx context.Context, inv map[string]interface{}) error {
	s.Client.mu.Lock()
	defer s.Client.mu.Unlock()
	s.Client.invalids = append(s.Client.invalids, inv)
	return nil
}

func (s *KubernetesStateStore) SaveIdempotencyRecord(ctx context.Context, record map[string]interface{}) error {
	s.Client.mu.Lock()
	defer s.Client.mu.Unlock()
	key := strOr(record["idempotencyKey"])
	if key == "" {
		return fmt.Errorf("idempotencyKey required")
	}
	s.Client.idempotency[key] = record
	return nil
}

func (s *KubernetesStateStore) LoadIdempotencyRecord(ctx context.Context, idempotencyKey string) (map[string]interface{}, error) {
	_ = ctx
	s.Client.mu.Lock()
	defer s.Client.mu.Unlock()
	rec, ok := s.Client.idempotency[idempotencyKey]
	if !ok {
		return nil, nil
	}
	return rec, nil
}

func (s *KubernetesStateStore) SetStatusCondition(ctx context.Context, condition map[string]interface{}) error {
	s.Client.mu.Lock()
	defer s.Client.mu.Unlock()
	s.Client.conditions = appendOrReplace(s.Client.conditions, condition, "type")
	return nil
}

func (s *KubernetesStateStore) EmitEvent(ctx context.Context, event map[string]interface{}) error {
	s.Client.mu.Lock()
	defer s.Client.mu.Unlock()
	s.Client.events = append(s.Client.events, event)
	return nil
}

func (c *FakeKubernetesClient) GetEvents() []map[string]interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]map[string]interface{}, len(c.events))
	copy(out, c.events)
	return out
}

func (c *FakeKubernetesClient) GetConditions() []map[string]interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]map[string]interface{}, len(c.conditions))
	copy(out, c.conditions)
	return out
}
