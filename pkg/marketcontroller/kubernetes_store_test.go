package marketcontroller

import (
	"context"
	"testing"
)

func TestKubernetesStoreStaleWrite(t *testing.T) {
	client := NewFakeKubernetesClient()
	store := NewKubernetesStateStore(client, DefaultDurableStateKey())
	ctx := context.Background()
	key := DefaultDurableStateKey()
	ds := newEmptyDurableState()
	rv1, err := store.Save(ctx, key, ds, "")
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.Save(ctx, key, ds, "wrong-rv")
	if err != ErrStaleWrite {
		t.Fatal("stale write must be rejected")
	}
	ds.StateVersion = 2
	rv2, err := store.Save(ctx, key, ds, rv1)
	if err != nil || rv2 == rv1 {
		t.Fatal("valid write must succeed with new rv")
	}
}

func TestKubernetesStoreIdempotencyPersist(t *testing.T) {
	client := NewFakeKubernetesClient()
	store := NewKubernetesStateStore(client, DefaultDurableStateKey())
	rec := map[string]interface{}{"idempotencyKey": "k1", "replayCount": 1, "duplicateSuppressed": true}
	_ = store.SaveIdempotencyRecord(context.Background(), rec)
	loaded, _ := store.LoadIdempotencyRecord(context.Background(), "k1")
	if loaded == nil || !boolOr(loaded["duplicateSuppressed"]) {
		t.Fatal("idempotency must persist")
	}
}
