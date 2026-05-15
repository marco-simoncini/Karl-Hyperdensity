package windowsfluidvirt

import "testing"

func TestWindowsFluidVirtBlockerCatalogContainsRequiredBlockers(t *testing.T) {
	required := []WindowsFluidVirtBlocker{
		BlockerMissingGuestWitness,
		BlockerMissingSameBootProof,
		BlockerMissingReturnToFloorPlan,
		BlockerMissingRollbackPlan,
		BlockerMissingAuditChain,
		BlockerNoManualApproval,
		BlockerRawRuntimeControlForbidden,
		BlockerAutonomousApplyForbidden,
		BlockerProductionReadyClaimForbidden,
		BlockerVCPUHotplugClaimForbidden,
		BlockerLogicalCPUScalingClaimForbidden,
		BlockerPoolScalingClaimForbidden,
	}
	for _, blocker := range required {
		if !IsKnownWindowsFluidVirtBlocker(blocker) {
			t.Fatalf("required blocker missing from catalog: %s", blocker)
		}
	}
}
