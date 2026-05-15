package safety

import "testing"

func TestMutationsForbidden(t *testing.T) {
	if !MutationsForbidden(false) {
		t.Fatal("expected forbidden when unsafe false")
	}
	if MutationsForbidden(true) {
		t.Fatal("expected not forbidden when unsafe true")
	}
}
