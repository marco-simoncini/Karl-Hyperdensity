package windowsfluidvirt

import "testing"

func TestNewDefaultWindowsFluidVirtActionSlatePlanningOnly(t *testing.T) {
	slate := NewDefaultWindowsFluidVirtActionSlate("v1")

	if !slate.PlanningOnly {
		t.Fatalf("planningOnly must be true")
	}
	if slate.ApplyEnabled || slate.RuntimeMutationEnabled || slate.AutonomousApplyEnabled {
		t.Fatalf("apply/mutation/autonomous flags must be false")
	}
	if slate.RawRuntimeControlsExposed {
		t.Fatalf("raw controls must be false")
	}
	if len(slate.Actions) == 0 {
		t.Fatalf("expected at least one action")
	}
	for _, action := range slate.Actions {
		if !action.PlanningOnly {
			t.Fatalf("action %s must be planning-only", action.ActionID)
		}
		if action.RuntimeMutationEnabled || action.AutonomousApplyEnabled || action.RawRuntimeControlsExposed {
			t.Fatalf("action %s exposes forbidden runtime toggles", action.ActionID)
		}
	}
	if blockers := slate.PlanningOnlyBlockers(); len(blockers) != 0 {
		t.Fatalf("unexpected blockers for default slate: %v", blockers)
	}
}

func TestPlanningOnlyBlockersDetectsUnsafeSlate(t *testing.T) {
	slate := NewDefaultWindowsFluidVirtActionSlate("v1")
	slate.ApplyEnabled = true
	slate.Actions[0].RawRuntimeControlsExposed = true

	blockers := slate.PlanningOnlyBlockers()
	if !containsBlocker(blockers, BlockerAutonomousApplyForbidden) {
		t.Fatalf("expected autonomous apply blocker, got %v", blockers)
	}
	if !containsBlocker(blockers, BlockerRawRuntimeControlForbidden) {
		t.Fatalf("expected raw control blocker, got %v", blockers)
	}
}
