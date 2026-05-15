package cgroup

import "testing"

func TestPlanEnvelopeNeverWrites(t *testing.T) {
	p := PlanEnvelope(false, "-10us", "-64Mi")
	if p.WouldWrite {
		t.Fatal("expected no write")
	}
	if len(p.WritePaths) != 0 {
		t.Fatal("expected no write paths")
	}
}

func TestPlanEnvelopeUnsafeFlagStillNoWrite(t *testing.T) {
	p := PlanEnvelope(true, "-10us", "-64Mi")
	if p.WouldWrite {
		t.Fatal("Sprint 6: unsafe flag must not enable write hints")
	}
}
