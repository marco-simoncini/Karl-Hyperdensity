package windowsfluidvirt

import "testing"

// Fundamental canonical blockers the planning / execution-safety lane relies on.
func TestCanonicalBlockerCatalogIncludesSafetyFundamentals(t *testing.T) {
	required := []string{
		BlockerPoolScalingAsMechanism,
		BlockerLeaseRequestsVCPUHotplug,
		BlockerLeaseRequestsVMSpecPatch,
		BlockerActuatorArbitraryWrite,
		BlockerRollbackNotReady,
		BlockerReturnToFloorNotReady,
		BlockerFutureApplyExecutorDisabled,
		BlockerNodeFluidActuatorUnavailable,
		BlockerRAMBalloonUnavailable,
		BlockerActuatorReplayDetected,
	}
	for _, id := range required {
		if !IsCanonicalWindowsFluidBlocker(id) {
			t.Fatalf("missing canonical blocker %q", id)
		}
	}
}
