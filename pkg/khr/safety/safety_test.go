package safety

import "testing"

func TestMutationsForbiddenAlwaysTrue(t *testing.T) {
	if !MutationsForbidden(false) {
		t.Fatal("expected forbidden")
	}
	if !MutationsForbidden(true) {
		t.Fatal("Sprint 6: expected forbidden even with unsafe flag")
	}
}

func TestUnsafeApplyRequested(t *testing.T) {
	if UnsafeApplyRequested(false) {
		t.Fatal("expected false")
	}
	if !UnsafeApplyRequested(true) {
		t.Fatal("expected true")
	}
}
