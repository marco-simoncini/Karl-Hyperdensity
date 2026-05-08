package windowsfluidvirt

import "testing"

func TestNewWindowsFluidVirtProductModelSafetyDefaults(t *testing.T) {
	model := NewWindowsFluidVirtProductModel("v1")

	if model.ReleaseTrack != ReleaseTrackTechnicalPreview {
		t.Fatalf("release track = %q, want technical preview", model.ReleaseTrack)
	}
	if model.LaneStatus != LaneStatusTechnicalPreviewCandidate {
		t.Fatalf("lane status = %q, want technical preview candidate", model.LaneStatus)
	}
	if model.SupportBoundary.ProductionMutationAllowed {
		t.Fatalf("production mutation must be disabled")
	}
	if model.SupportBoundary.AutonomousApplyAllowed {
		t.Fatalf("autonomous apply must be disabled")
	}
	if model.SupportBoundary.EnforcementMode != "disabled" {
		t.Fatalf("enforcement mode = %q, want disabled", model.SupportBoundary.EnforcementMode)
	}
	if model.SupportBoundary.WindowsGaClaimAllowed || model.SupportBoundary.WindowsProductionReadyClaimAllowed {
		t.Fatalf("GA and production-ready claims must be disabled")
	}
	if model.SupportBoundary.WindowsExecutionReadyByDefault {
		t.Fatalf("execution-ready by default must be disabled")
	}
	if model.SupportBoundary.VCPUHotplugClaimAllowed || model.SupportBoundary.LogicalCPUScalingClaimAllowed || model.SupportBoundary.PoolScalingClaimAllowed {
		t.Fatalf("hotplug/logical/pool claims must be disabled")
	}
	if model.SupportBoundary.LiveMigrationClaimAllowed || model.SupportBoundary.RebootRecreateRolloutMechanismAllowed {
		t.Fatalf("migration and reboot/recreate/rollout claims must be disabled")
	}
	if model.SupportBoundary.RawRuntimeControlsExposed {
		t.Fatalf("raw controls must be disabled")
	}
	if model.ActionSlate.ApplyEnabled || model.ActionSlate.RuntimeMutationEnabled || model.ActionSlate.AutonomousApplyEnabled {
		t.Fatalf("action slate must stay planning-only")
	}
	if model.CPULiquidityModel.ApplyEnabled || model.RAMLiquidityModel.ApplyEnabled || model.LeaseModel.ApplyEnabled {
		t.Fatalf("apply must be disabled in all core models")
	}
	if blockers := model.ValidateSafety(); len(blockers) != 0 {
		t.Fatalf("unexpected safety blockers for default model: %v", blockers)
	}
}

func TestValidateSafetyFindsForbiddenClaims(t *testing.T) {
	model := NewWindowsFluidVirtProductModel("v1")
	model.SupportBoundary.WindowsProductionReadyClaimAllowed = true
	model.SupportBoundary.VCPUHotplugClaimAllowed = true
	model.SupportBoundary.RawRuntimeControlsExposed = true
	model.LeaseModel.ReturnToFloorRequired = false

	blockers := model.ValidateSafety()
	if !containsBlocker(blockers, BlockerProductionReadyClaimForbidden) {
		t.Fatalf("expected production ready blocker, got %v", blockers)
	}
	if !containsBlocker(blockers, BlockerVCPUHotplugClaimForbidden) {
		t.Fatalf("expected vcpu blocker, got %v", blockers)
	}
	if !containsBlocker(blockers, BlockerRawRuntimeControlForbidden) {
		t.Fatalf("expected raw control blocker, got %v", blockers)
	}
	if !containsBlocker(blockers, BlockerMissingReturnToFloorPlan) {
		t.Fatalf("expected return-to-floor blocker, got %v", blockers)
	}
}

func containsBlocker(values []WindowsFluidVirtBlocker, target WindowsFluidVirtBlocker) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
