package windowsfluidvirt

type WindowsFluidVirtActionStatus string

const (
	ActionStatusCandidate WindowsFluidVirtActionStatus = "candidate"
	ActionStatusGated     WindowsFluidVirtActionStatus = "gated"
	ActionStatusBlocked   WindowsFluidVirtActionStatus = "blocked"
)

type WindowsFluidVirtActionKind string

const (
	ActionObserveRuntime       WindowsFluidVirtActionKind = "observe_runtime"
	ActionAssessCompliance     WindowsFluidVirtActionKind = "assess_compliance"
	ActionPrepareLeasePlan     WindowsFluidVirtActionKind = "prepare_lease_plan"
	ActionBuildActionSlate     WindowsFluidVirtActionKind = "build_action_slate"
	ActionDryRunEvidence       WindowsFluidVirtActionKind = "dry_run_evidence"
	ActionVerifyGuestWitness   WindowsFluidVirtActionKind = "verify_guest_witness"
	ActionValidateRollbackPlan WindowsFluidVirtActionKind = "validate_rollback_plan"
	ActionValidateReturnToFloor WindowsFluidVirtActionKind = "validate_return_to_floor_plan"
	ActionCheckAuditChain      WindowsFluidVirtActionKind = "check_audit_chain"
)

type WindowsFluidVirtAction struct {
	ActionID                  string                     `json:"actionId"`
	Kind                      WindowsFluidVirtActionKind `json:"kind"`
	Status                    WindowsFluidVirtActionStatus `json:"status"`
	PlanningOnly              bool                       `json:"planningOnly"`
	RuntimeMutationEnabled    bool                       `json:"runtimeMutationEnabled"`
	AutonomousApplyEnabled    bool                       `json:"autonomousApplyEnabled"`
	RawRuntimeControlsExposed bool                       `json:"rawRuntimeControlsExposed"`
	EvidenceRefs              []string                   `json:"evidenceRefs,omitempty"`
	Blockers                  []WindowsFluidVirtBlocker  `json:"blockers,omitempty"`
}

type WindowsFluidVirtActionSlate struct {
	SlateID                   string                    `json:"slateId"`
	SlateVersion              string                    `json:"slateVersion"`
	PlanningOnly              bool                      `json:"planningOnly"`
	ApplyEnabled              bool                      `json:"applyEnabled"`
	RuntimeMutationEnabled    bool                      `json:"runtimeMutationEnabled"`
	AutonomousApplyEnabled    bool                      `json:"autonomousApplyEnabled"`
	RawRuntimeControlsExposed bool                      `json:"rawRuntimeControlsExposed"`
	Actions                   []WindowsFluidVirtAction  `json:"actions"`
	Blockers                  []WindowsFluidVirtBlocker `json:"blockers,omitempty"`
}

func NewDefaultWindowsFluidVirtActionSlate(version string) WindowsFluidVirtActionSlate {
	actions := []WindowsFluidVirtAction{
		{ActionID: "a1", Kind: ActionObserveRuntime, Status: ActionStatusCandidate, PlanningOnly: true},
		{ActionID: "a2", Kind: ActionAssessCompliance, Status: ActionStatusCandidate, PlanningOnly: true},
		{ActionID: "a3", Kind: ActionPrepareLeasePlan, Status: ActionStatusCandidate, PlanningOnly: true},
		{ActionID: "a4", Kind: ActionBuildActionSlate, Status: ActionStatusCandidate, PlanningOnly: true},
		{ActionID: "a5", Kind: ActionDryRunEvidence, Status: ActionStatusGated, PlanningOnly: true},
		{ActionID: "a6", Kind: ActionVerifyGuestWitness, Status: ActionStatusGated, PlanningOnly: true},
		{ActionID: "a7", Kind: ActionValidateRollbackPlan, Status: ActionStatusGated, PlanningOnly: true},
		{ActionID: "a8", Kind: ActionValidateReturnToFloor, Status: ActionStatusGated, PlanningOnly: true},
		{ActionID: "a9", Kind: ActionCheckAuditChain, Status: ActionStatusGated, PlanningOnly: true},
	}
	for i := range actions {
		actions[i].RuntimeMutationEnabled = false
		actions[i].AutonomousApplyEnabled = false
		actions[i].RawRuntimeControlsExposed = false
	}
	return WindowsFluidVirtActionSlate{
		SlateID:                   "windows-fluidvirt-action-slate-core-v1",
		SlateVersion:              version,
		PlanningOnly:              true,
		ApplyEnabled:              false,
		RuntimeMutationEnabled:    false,
		AutonomousApplyEnabled:    false,
		RawRuntimeControlsExposed: false,
		Actions:                   actions,
	}
}

func (s WindowsFluidVirtActionSlate) PlanningOnlyBlockers() []WindowsFluidVirtBlocker {
	blockers := make([]WindowsFluidVirtBlocker, 0, 4)
	if s.RuntimeMutationEnabled || s.ApplyEnabled {
		blockers = append(blockers, BlockerAutonomousApplyForbidden)
	}
	if s.AutonomousApplyEnabled {
		blockers = append(blockers, BlockerAutonomousApplyForbidden)
	}
	if s.RawRuntimeControlsExposed {
		blockers = append(blockers, BlockerRawRuntimeControlForbidden)
	}
	for _, action := range s.Actions {
		if !action.PlanningOnly || action.RuntimeMutationEnabled || action.AutonomousApplyEnabled {
			blockers = append(blockers, BlockerAutonomousApplyForbidden)
		}
		if action.RawRuntimeControlsExposed {
			blockers = append(blockers, BlockerRawRuntimeControlForbidden)
		}
	}
	return dedupeBlockers(blockers)
}
