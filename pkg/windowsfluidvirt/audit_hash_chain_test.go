package windowsfluidvirt

import "testing"

func TestAuditHashChainBuildsDeterministically(t *testing.T) {
	replay := NewWindowsFluidVirtComplianceReplayMinimal()
	chainA, eventsA, err := BuildWindowsFluidVirtAuditHashChain(replay.ReplayedEvents)
	if err != nil {
		t.Fatalf("build chain A: %v", err)
	}
	chainB, eventsB, err := BuildWindowsFluidVirtAuditHashChain(replay.ReplayedEvents)
	if err != nil {
		t.Fatalf("build chain B: %v", err)
	}
	if chainA.RootHash != chainB.RootHash || chainA.TerminalHash != chainB.TerminalHash {
		t.Fatalf("deterministic hash chain mismatch")
	}
	if len(eventsA) != len(eventsB) {
		t.Fatalf("event lengths mismatch")
	}
	if !VerifyWindowsFluidVirtAuditHashChain(chainA, eventsA) {
		t.Fatalf("verification should succeed for deterministic chain")
	}
}

func TestAuditHashChainDetectsEventTampering(t *testing.T) {
	replay := NewWindowsFluidVirtComplianceReplayMinimal()
	chain, events, err := BuildWindowsFluidVirtAuditHashChain(replay.ReplayedEvents)
	if err != nil {
		t.Fatalf("build chain: %v", err)
	}
	events[0].EventType = "tampered_event_type"
	if VerifyWindowsFluidVirtAuditHashChain(chain, events) {
		t.Fatalf("verification must fail after tampering")
	}
}
