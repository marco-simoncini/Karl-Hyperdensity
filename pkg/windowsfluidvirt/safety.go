package windowsfluidvirt

// Planning-lane safety (Sprint 15): read-only view over The-Father-Windows types.
//
// TF-W uses WindowsFluidActionSlate with top-level MutationAllowed / ApplyAllowed
// (envelope) plus per-step WindowsFluidAction.MutationAllowed. Sprint 15 lane
// policy: the default slate built by BuildWindowsFluidActionSlate must present
// as planning-only — no envelope apply, no executed runtime mutation, and no
// per-step mutation flags (steps remain typed for future gated phases only).
//
// FVI “product model” fields (PlanningOnly, ApplyEnabled, AutonomousApplyEnabled,
// RawRuntimeControlsExposed, ProductionReadyClaim, …) are NOT on TF-W structs;
// this file maps TF-W fields into that conceptual model for assertions and docs.

// WindowsFluidVirtPlanningSafety is a derived, read-only summary for tests and audits.
type WindowsFluidVirtPlanningSafety struct {
	PlanningOnly                  bool
	ApplyEnabled                  bool
	RuntimeMutationEnabled        bool
	AutonomousApplyEnabled        bool
	RawRuntimeControlsExposed     bool
	ProductionReadyClaim          bool
	VCPUHotplugClaim              bool
	LogicalCPUScalingClaim        bool
	PoolScalingClaim              bool
	AnyActionDeclaresStepMutation bool
}

// DerivePlanningSafety maps TF-W slate fields into the planning-only safety model.
//
// Mapping (TF-W → model):
//
//	ApplyEnabled                  == slate.ApplyAllowed
//	RuntimeMutationEnabled        == slate.RuntimeMutationExecuted || slate.MutationAllowed
//	AutonomousApplyEnabled        == false (no distinct TF-W field; envelope forbids apply)
//	RawRuntimeControlsExposed     == false (no TF-W field; reserved — must stay false in lane)
//	ProductionReadyClaim          == false (lane never asserts GA; callers use evidence gates)
//	VCPUHotplug / LogicalCPU / PoolScaling claims are target/lease-driven — pass separately.
//	PlanningOnly                  == no violations from EvaluatePlanningSafetyViolations(slate)
//	  AND no per-step MutationAllowed (strict planning lane).
func DerivePlanningSafety(slate WindowsFluidActionSlate) WindowsFluidVirtPlanningSafety {
	violations := EvaluatePlanningSafetyViolations(slate)
	anyStep := false
	for _, a := range slate.Actions {
		if a.MutationAllowed {
			anyStep = true
			break
		}
	}
	return WindowsFluidVirtPlanningSafety{
		PlanningOnly:                  len(violations) == 0,
		ApplyEnabled:                  slate.ApplyAllowed,
		RuntimeMutationEnabled:        slate.RuntimeMutationExecuted || slate.MutationAllowed,
		AutonomousApplyEnabled:        false,
		RawRuntimeControlsExposed:     false,
		ProductionReadyClaim:          false,
		VCPUHotplugClaim:              false,
		LogicalCPUScalingClaim:        false,
		PoolScalingClaim:              false,
		AnyActionDeclaresStepMutation: anyStep,
	}
}

// DerivePlanningSafetyWithTargetLease overlays forbidden Windows *claims* from TF-W target/lease.
func DerivePlanningSafetyWithTargetLease(slate WindowsFluidActionSlate, target WindowsHyperdensityTarget, lease WindowsFluidResourceLease) WindowsFluidVirtPlanningSafety {
	s := DerivePlanningSafety(slate)
	s.VCPUHotplugClaim = target.VCPUHotplugRequested || lease.RequestsVCPUHotplug
	s.LogicalCPUScalingClaim = target.LogicalCPUScalingClaimed || lease.LogicalCPUScalingClaim
	s.PoolScalingClaim = target.PoolScalingRequested || lease.RequestedMechanism == "pool-scaling"
	s.PlanningOnly = s.PlanningOnly && !s.VCPUHotplugClaim && !s.LogicalCPUScalingClaim && !s.PoolScalingClaim
	return s
}

// EvaluatePlanningSafetyViolations returns human-readable violations of the planning lane.
func EvaluatePlanningSafetyViolations(slate WindowsFluidActionSlate) []string {
	var v []string
	if slate.ApplyAllowed {
		v = append(v, "planning_lane: envelope ApplyAllowed must be false")
	}
	if slate.MutationAllowed {
		v = append(v, "planning_lane: envelope MutationAllowed must be false")
	}
	if slate.RuntimeMutationExecuted {
		v = append(v, "planning_lane: RuntimeMutationExecuted must be false")
	}
	for _, a := range slate.Actions {
		if a.MutationAllowed {
			v = append(v, "planning_lane: action "+a.ActionID+" must not declare MutationAllowed in planning lane")
		}
	}
	return v
}

// IsCanonicalWindowsFluidBlocker reports whether id is registered in CanonicalBlockers.
func IsCanonicalWindowsFluidBlocker(id string) bool {
	_, ok := LookupBlocker(id)
	return ok
}
