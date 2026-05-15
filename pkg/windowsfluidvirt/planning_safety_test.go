package windowsfluidvirt

import (
	"testing"
	"time"
)

func TestDefaultWindowsFluidActionSlateIsPlanningOnly(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := PrepareWindowsFluidResourceLease(target, baseCombinedLease())
	slate := BuildWindowsFluidActionSlate(target, lease)

	safe := DerivePlanningSafety(slate)
	if !safe.PlanningOnly {
		t.Fatalf("expected planning-only slate, violations=%v", EvaluatePlanningSafetyViolations(slate))
	}
	if safe.ApplyEnabled || safe.RuntimeMutationEnabled {
		t.Fatalf("envelope must not allow apply/runtime mutation: %+v", safe)
	}
	if safe.AutonomousApplyEnabled || safe.RawRuntimeControlsExposed || safe.ProductionReadyClaim {
		t.Fatalf("lane must not assert autonomous/raw/production: %+v", safe)
	}
	if safe.AnyActionDeclaresStepMutation {
		t.Fatalf("no action step should declare MutationAllowed in default slate")
	}
	for _, a := range slate.Actions {
		if a.MutationAllowed {
			t.Fatalf("action %s should not declare mutation in planning lane", a.ActionID)
		}
	}
}

func TestDryRunEnvelopeSlateIsPlanningSafe(t *testing.T) {
	shell := WindowsFluidShell{
		Spec: WindowsFluidShellSpec{VMRef: "vm-1"},
		Status: WindowsFluidShellStatus{
			EvidenceRef:        "e1",
			LastTransitionTime: time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC),
		},
	}
	slate := BuildDryRunActionSlate("s1", shell, StateReady, nil, map[string]bool{
		"rollbackReady": true, "returnToFloorReady": true, "qmpReady": true,
		"guestAckReady": true, "noRebootProof": true, "sameQemuProcess": true,
		"sameNode": true, "sameVirtLauncherPod": true,
	}, nil)
	v := EvaluatePlanningSafetyViolations(slate)
	if len(v) != 0 {
		t.Fatalf("dry-run envelope slate: %v", v)
	}
}

func TestUnsafeSlateViolationsEnvelopeApply(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := PrepareWindowsFluidResourceLease(target, baseCombinedLease())
	slate := BuildWindowsFluidActionSlate(target, lease)
	slate.ApplyAllowed = true
	if len(EvaluatePlanningSafetyViolations(slate)) == 0 {
		t.Fatal("expected violations when ApplyAllowed true")
	}
}

func TestUnsafeSlateViolationsRuntimeMutationExecuted(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := PrepareWindowsFluidResourceLease(target, baseCombinedLease())
	slate := BuildWindowsFluidActionSlate(target, lease)
	slate.RuntimeMutationExecuted = true
	if len(EvaluatePlanningSafetyViolations(slate)) == 0 {
		t.Fatal("expected violations when RuntimeMutationExecuted true")
	}
}

func TestUnsafeSlateViolationsStepMutation(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := PrepareWindowsFluidResourceLease(target, baseCombinedLease())
	slate := BuildWindowsFluidActionSlate(target, lease)
	if len(slate.Actions) == 0 {
		t.Fatal("expected actions")
	}
	slate.Actions[0].MutationAllowed = true
	if len(EvaluatePlanningSafetyViolations(slate)) == 0 {
		t.Fatal("expected violations when a step declares mutation")
	}
}

func TestDerivePlanningSafetyWithTargetLeaseCleanOverlay(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := PrepareWindowsFluidResourceLease(target, baseCombinedLease())
	slate := BuildWindowsFluidActionSlate(target, lease)
	s := DerivePlanningSafetyWithTargetLease(slate, target, lease)
	if !s.PlanningOnly || s.VCPUHotplugClaim || s.LogicalCPUScalingClaim || s.PoolScalingClaim {
		t.Fatalf("expected clean planning overlay, got %+v", s)
	}
}
