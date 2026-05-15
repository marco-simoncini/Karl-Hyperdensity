package cgroup

import "testing"

func TestPlanEnvelopeNoWriteByDefault(t *testing.T) {
	p := PlanEnvelope(false, "-10us", "-64Mi")
	if p.WouldWrite {
		t.Fatal("expected no write")
	}
	if len(p.WritePaths) != 0 {
		t.Fatal("expected no write paths")
	}
}

func TestPlanEnvelopeWouldWriteOnlyWithUnsafe(t *testing.T) {
	p := PlanEnvelope(true, "-10us", "")
	if !p.WouldWrite {
		t.Fatal("expected wouldWrite when unsafe enabled with delta")
	}
}
