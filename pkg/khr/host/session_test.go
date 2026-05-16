package host

import "testing"

func TestRuntimeSessionStableIDs(t *testing.T) {
	ResetRuntimeSession()
	cfg := testHostConfig()
	s1 := InitRuntimeSession(cfg)
	s2 := InitRuntimeSession(cfg)
	if s1.RuntimeSessionID != s2.RuntimeSessionID {
		t.Fatalf("session id not stable: %q vs %q", s1.RuntimeSessionID, s2.RuntimeSessionID)
	}
	if s1.HostRuntimeInstanceID != s2.HostRuntimeInstanceID {
		t.Fatal("instance id not stable")
	}
	if s1.CorrelationID == s2.CorrelationID {
		t.Fatal("correlation id must advance per call")
	}
}

func TestSetCorrelationID(t *testing.T) {
	ResetRuntimeSession()
	InitRuntimeSession(testHostConfig())
	SetCorrelationID("khr-corr-apply-test")
	s := CurrentRuntimeSession()
	if s.CorrelationID != "khr-corr-apply-test" {
		t.Fatalf("correlation=%q", s.CorrelationID)
	}
}
