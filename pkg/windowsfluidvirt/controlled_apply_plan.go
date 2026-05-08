package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"time"
)

const (
	ControlledApplyBlockerManualApprovalRequired = "controlled_apply_manual_approval_required"
	ControlledApplyBlockerApprovalRejected       = "controlled_apply_manual_approval_rejected"
	ControlledApplyBlockerApprovalExpired        = "controlled_apply_manual_approval_expired"
	ControlledApplyBlockerAutonomousApplyDenied  = "controlled_apply_autonomous_apply_denied"
	ControlledApplyBlockerDryRunRequired         = "controlled_apply_dry_run_required"
	ControlledApplyBlockerWorkloadVerifyMissing  = "controlled_apply_workload_verify_missing"
	ControlledApplyBlockerAuditBundleMissing     = "controlled_apply_audit_bundle_missing"
	ControlledApplyBlockerKillSwitchBlocked      = "controlled_apply_kill_switch_blocked"
	ControlledApplyBlockerNamespaceNotAllowlisted = "controlled_apply_namespace_not_allowlisted"
	ControlledApplyBlockerTargetNotAllowlisted   = "controlled_apply_target_not_allowlisted"
	ControlledApplyBlockerLeaseKindNotAllowlisted = "controlled_apply_lease_kind_not_allowlisted"
	ControlledApplyBlockerBlastRadiusInvalid     = "controlled_apply_blast_radius_invalid"
)

type WindowsFluidControlledApplyGate struct {
	GateID                  string                 `json:"gateId"`
	WindowsFluidVirtEnabled bool                   `json:"windowsFluidVirtEnabled"`
	DryRunRequired          bool                   `json:"dryRunRequired"`
	ManualApprovalRequired  bool                   `json:"manualApprovalRequired"`
	AutonomousApplyEnabled  bool                   `json:"autonomousApplyEnabled"`
	CPUApplyEnabled         bool                   `json:"cpuApplyEnabled"`
	RAMApplyEnabled         bool                   `json:"ramApplyEnabled"`
	NodeActuatorApplyEnabled bool                  `json:"nodeActuatorApplyEnabled"`
	QMPBalloonApplyEnabled  bool                   `json:"qmpBalloonApplyEnabled"`
	GuestVerifyRequired     bool                   `json:"guestVerifyRequired"`
	WorkloadVerifyRequired  bool                   `json:"workloadVerifyRequired"`
	RollbackRequired        bool                   `json:"rollbackRequired"`
	ReturnToFloorRequired   bool                   `json:"returnToFloorRequired"`
	AuditBundleRequired     bool                   `json:"auditBundleRequired"`
	KillSwitchRequired      bool                   `json:"killSwitchRequired"`
	MaxBlastRadius          int                    `json:"maxBlastRadius"`
	AllowedNamespaces       []string               `json:"allowedNamespaces"`
	AllowedTargets          []string               `json:"allowedTargets"`
	AllowedLeaseKinds       []WindowsFluidLeaseKind `json:"allowedLeaseKinds"`
	CreatedAt               string                 `json:"createdAt"`
}

type WindowsFluidManualApproval struct {
	ApprovalID  string `json:"approvalId"`
	ApprovedBy  string `json:"approvedBy,omitempty"`
	Reason      string `json:"reason,omitempty"`
	State       string `json:"state"`
	CreatedAt   string `json:"createdAt"`
	ExpiresAt   string `json:"expiresAt,omitempty"`
}

type WindowsFluidControlledApplyApprovalState string

const (
	ApprovalNotRequired WindowsFluidControlledApplyApprovalState = "not_required"
	ApprovalRequired    WindowsFluidControlledApplyApprovalState = "required"
	ApprovalApproved    WindowsFluidControlledApplyApprovalState = "approved"
	ApprovalRejected    WindowsFluidControlledApplyApprovalState = "rejected"
	ApprovalExpired     WindowsFluidControlledApplyApprovalState = "expired"
)

type WindowsFluidControlledPlanPhase string

const (
	ControlledPlanPrepared              WindowsFluidControlledPlanPhase = "prepared"
	ControlledPlanDryRunReady           WindowsFluidControlledPlanPhase = "dry_run_ready"
	ControlledPlanAwaitingApproval      WindowsFluidControlledPlanPhase = "awaiting_approval"
	ControlledPlanApplyReady            WindowsFluidControlledPlanPhase = "apply_ready"
	ControlledPlanApplyBlocked          WindowsFluidControlledPlanPhase = "apply_blocked"
	ControlledPlanVerifyRequired        WindowsFluidControlledPlanPhase = "verify_required"
	ControlledPlanRollbackRequired      WindowsFluidControlledPlanPhase = "rollback_required"
	ControlledPlanReturnToFloorRequired WindowsFluidControlledPlanPhase = "return_to_floor_required"
	ControlledPlanCompleted             WindowsFluidControlledPlanPhase = "completed"
	ControlledPlanBlocked               WindowsFluidControlledPlanPhase = "blocked"
	ControlledPlanQuarantined           WindowsFluidControlledPlanPhase = "quarantined"
)

type WindowsFluidControlledApplyPlan struct {
	PlanID              string                                   `json:"planId"`
	TargetRef           string                                   `json:"targetRef"`
	LeaseRef            string                                   `json:"leaseRef"`
	GateSnapshot        WindowsFluidControlledApplyGate          `json:"gateSnapshot"`
	ComplianceSnapshot  EvaluateWindowsHyperdensityReadyComplianceOutput `json:"complianceSnapshot"`
	ActionSlate         WindowsFluidActionSlate                  `json:"actionSlate"`
	ActuatorRequest     map[string]any                           `json:"actuatorRequest,omitempty"`
	QMPBalloonRequest   map[string]any                           `json:"qmpBalloonRequest,omitempty"`
	GuestVerifyPlan     map[string]any                           `json:"guestVerifyPlan"`
	WorkloadVerifyPlan  map[string]any                           `json:"workloadVerifyPlan"`
	RollbackPlan        map[string]any                           `json:"rollbackPlan"`
	ReturnToFloorPlan   map[string]any                           `json:"returnToFloorPlan"`
	AuditBundlePlan     map[string]any                           `json:"auditBundlePlan"`
	ApprovalState       WindowsFluidControlledApplyApprovalState `json:"approvalState"`
	PlanPhase           WindowsFluidControlledPlanPhase          `json:"planPhase"`
	MutationAllowed     bool                                     `json:"mutationAllowed"`
	ApplyAllowed        bool                                     `json:"applyAllowed"`
	Blockers            []string                                 `json:"blockers"`
	EvidenceRefs        []string                                 `json:"evidenceRefs"`
	CreatedAt           string                                   `json:"createdAt"`
}

type WindowsFluidControlledApplyFixture struct {
	Name              string                          `json:"name"`
	Mode              string                          `json:"mode"`
	Target            WindowsHyperdensityTarget       `json:"target"`
	Lease             WindowsFluidResourceLease       `json:"lease"`
	Gate              WindowsFluidControlledApplyGate `json:"gate"`
	Approval          *WindowsFluidManualApproval     `json:"approval,omitempty"`
	ExpectedPlanPhase WindowsFluidControlledPlanPhase `json:"expectedPlanPhase"`
	ExpectedApplyAllowed bool                         `json:"expectedApplyAllowed"`
	ExpectedBlockers  []string                        `json:"expectedBlockers"`
}

func DefaultWindowsFluidControlledApplyGate() WindowsFluidControlledApplyGate {
	return WindowsFluidControlledApplyGate{
		GateID:                   "windows-fluid-controlled-apply-gate-v1",
		WindowsFluidVirtEnabled:  true,
		DryRunRequired:           true,
		ManualApprovalRequired:   true,
		AutonomousApplyEnabled:   false,
		CPUApplyEnabled:          false,
		RAMApplyEnabled:          false,
		NodeActuatorApplyEnabled: false,
		QMPBalloonApplyEnabled:   false,
		GuestVerifyRequired:      true,
		WorkloadVerifyRequired:   true,
		RollbackRequired:         true,
		ReturnToFloorRequired:    true,
		AuditBundleRequired:      true,
		KillSwitchRequired:       true,
		MaxBlastRadius:           1,
		AllowedNamespaces:        []string{"karl"},
		AllowedTargets:           []string{"master-win11"},
		AllowedLeaseKinds:        []WindowsFluidLeaseKind{LeaseKindCPUEntitlement, LeaseKindRAMBalloon, LeaseKindCombinedEnvelope},
		CreatedAt:                "2026-05-08T00:00:00Z",
	}
}

func EvaluateWindowsFluidControlledApplyGate(
	target WindowsHyperdensityTarget,
	lease WindowsFluidResourceLease,
	gate WindowsFluidControlledApplyGate,
) []string {
	blockers := make([]string, 0, 20)
	if !gate.WindowsFluidVirtEnabled {
		blockers = append(blockers, BlockerFutureApplyExecutorDisabled)
	}
	if gate.AutonomousApplyEnabled {
		blockers = append(blockers, ControlledApplyBlockerAutonomousApplyDenied)
	}
	if target.PoolScalingRequested || lease.RequestedMechanism == "pool-scaling" {
		blockers = append(blockers, BlockerPoolScalingAsMechanism)
	}
	if lease.RequestsVCPUHotplug || target.VCPUHotplugRequested {
		blockers = append(blockers, BlockerLeaseRequestsVCPUHotplug)
	}
	if lease.LogicalCPUScalingClaim || target.LogicalCPUScalingClaimed {
		blockers = append(blockers, BlockerLeaseRequestsVCPUHotplug)
	}
	if lease.RequestsVMSpecPatch {
		blockers = append(blockers, BlockerLeaseRequestsVMSpecPatch)
	}
	if gate.MaxBlastRadius <= 0 {
		blockers = append(blockers, ControlledApplyBlockerBlastRadiusInvalid)
	}
	if len(gate.AllowedNamespaces) > 0 && !contains(gate.AllowedNamespaces, target.Namespace) {
		blockers = append(blockers, ControlledApplyBlockerNamespaceNotAllowlisted)
	}
	if len(gate.AllowedTargets) > 0 && !contains(gate.AllowedTargets, target.VMRef) {
		blockers = append(blockers, ControlledApplyBlockerTargetNotAllowlisted)
	}
	if len(gate.AllowedLeaseKinds) > 0 && !containsLeaseKind(gate.AllowedLeaseKinds, lease.LeaseKind) {
		blockers = append(blockers, ControlledApplyBlockerLeaseKindNotAllowlisted)
	}
	if (lease.LeaseKind == LeaseKindCPUEntitlement || lease.LeaseKind == LeaseKindCombinedEnvelope) && !target.CPU.ActuatorRequired {
		blockers = append(blockers, BlockerNodeFluidActuatorUnavailable)
	}
	if (lease.LeaseKind == LeaseKindRAMBalloon || lease.LeaseKind == LeaseKindCombinedEnvelope) && !target.Memory.QMPRequired {
		blockers = append(blockers, BlockerRAMBalloonUnavailable)
	}
	return dedupe(blockers)
}

func EvaluateWindowsFluidManualApproval(
	gate WindowsFluidControlledApplyGate,
	approval *WindowsFluidManualApproval,
	evaluationTime time.Time,
) (WindowsFluidControlledApplyApprovalState, []string) {
	if !gate.ManualApprovalRequired {
		return ApprovalNotRequired, nil
	}
	if approval == nil {
		return ApprovalRequired, []string{ControlledApplyBlockerManualApprovalRequired}
	}
	state := WindowsFluidControlledApplyApprovalState(approval.State)
	switch state {
	case ApprovalApproved:
		if approval.ExpiresAt != "" {
			expiresAt, err := time.Parse(time.RFC3339, approval.ExpiresAt)
			if err != nil || normalizeControlledTime(evaluationTime).After(expiresAt.UTC()) {
				return ApprovalExpired, []string{ControlledApplyBlockerApprovalExpired}
			}
		}
		return ApprovalApproved, nil
	case ApprovalRejected:
		return ApprovalRejected, []string{ControlledApplyBlockerApprovalRejected}
	case ApprovalExpired:
		return ApprovalExpired, []string{ControlledApplyBlockerApprovalExpired}
	default:
		return ApprovalRequired, []string{ControlledApplyBlockerManualApprovalRequired}
	}
}

func EvaluateWindowsFluidDryRunGate(
	gate WindowsFluidControlledApplyGate,
	lease WindowsFluidResourceLease,
	slate WindowsFluidActionSlate,
) []string {
	blockers := make([]string, 0, 4)
	if gate.DryRunRequired && !hasActionType(slate, ActionActuatorDryRun) {
		blockers = append(blockers, ControlledApplyBlockerDryRunRequired)
	}
	if gate.DryRunRequired && len(slate.Blockers) > 0 {
		blockers = append(blockers, ControlledApplyBlockerDryRunRequired)
	}
	if gate.DryRunRequired && lease.Status == LeaseStatusBlocked {
		blockers = append(blockers, ControlledApplyBlockerDryRunRequired)
	}
	return dedupe(blockers)
}

func EvaluateWindowsFluidVerificationPlan(
	gate WindowsFluidControlledApplyGate,
	lease WindowsFluidResourceLease,
) (map[string]any, []string) {
	blockers := make([]string, 0, 2)
	workloadRef := policySnapshotString(lease.PolicySnapshot, "workloadEvidenceRef")
	guestReady := !gate.GuestVerifyRequired || lease.GuestEvidenceRef != ""
	workloadReady := !gate.WorkloadVerifyRequired || workloadRef != ""
	if gate.GuestVerifyRequired && !guestReady {
		blockers = append(blockers, BlockerGuestAckMissing)
	}
	if gate.WorkloadVerifyRequired && !workloadReady {
		blockers = append(blockers, ControlledApplyBlockerWorkloadVerifyMissing)
	}
	return map[string]any{
		"guestVerifyRequired":    gate.GuestVerifyRequired,
		"guestEvidenceRef":       lease.GuestEvidenceRef,
		"guestVerifyReady":       guestReady,
		"workloadVerifyRequired": gate.WorkloadVerifyRequired,
		"workloadEvidenceRef":    workloadRef,
		"workloadVerifyReady":    workloadReady,
	}, dedupe(blockers)
}

func EvaluateWindowsFluidRollbackPlan(
	gate WindowsFluidControlledApplyGate,
	lease WindowsFluidResourceLease,
) (map[string]any, []string) {
	blockers := make([]string, 0, 2)
	ready := len(EvaluateWindowsFluidRollbackReadiness(lease)) == 0
	if gate.RollbackRequired && !ready {
		blockers = append(blockers, BlockerRollbackNotReady)
	}
	return map[string]any{
		"required": gate.RollbackRequired,
		"ready":    ready,
		"target":   lease.RollbackTarget,
	}, dedupe(blockers)
}

func EvaluateWindowsFluidReturnToFloorPlan(
	gate WindowsFluidControlledApplyGate,
	lease WindowsFluidResourceLease,
) (map[string]any, []string) {
	blockers := make([]string, 0, 2)
	ready := len(EvaluateWindowsFluidReturnToFloorReadiness(lease)) == 0
	if gate.ReturnToFloorRequired && !ready {
		blockers = append(blockers, BlockerReturnToFloorNotReady)
	}
	return map[string]any{
		"required": gate.ReturnToFloorRequired,
		"ready":    ready,
		"target":   lease.ReturnToFloorTarget,
	}, dedupe(blockers)
}

func EvaluateWindowsFluidAuditBundlePlan(
	gate WindowsFluidControlledApplyGate,
	lease WindowsFluidResourceLease,
) (map[string]any, []string) {
	blockers := make([]string, 0, 2)
	ready := !gate.AuditBundleRequired || lease.AuditBundleRef != ""
	if gate.AuditBundleRequired && !ready {
		blockers = append(blockers, ControlledApplyBlockerAuditBundleMissing)
	}
	return map[string]any{
		"required":  gate.AuditBundleRequired,
		"ready":     ready,
		"bundleRef": lease.AuditBundleRef,
	}, dedupe(blockers)
}

func EvaluateWindowsFluidApplyReadiness(
	target WindowsHyperdensityTarget,
	lease WindowsFluidResourceLease,
	gate WindowsFluidControlledApplyGate,
	approvalState WindowsFluidControlledApplyApprovalState,
	dryRunBlockers []string,
	verifyBlockers []string,
	rollbackBlockers []string,
	returnBlockers []string,
	auditBlockers []string,
) []string {
	blockers := make([]string, 0, 20)
	if !target.HyperdensityReady || target.CompliancePhase != ComplianceHyperdensityReadyWindowsShell {
		blockers = append(blockers, BlockerCPUTopologyNotConfirmed)
	}
	if gate.AutonomousApplyEnabled {
		blockers = append(blockers, ControlledApplyBlockerAutonomousApplyDenied)
	}
	if gate.ManualApprovalRequired && approvalState != ApprovalApproved {
		blockers = append(blockers, ControlledApplyBlockerManualApprovalRequired)
	}
	if len(dryRunBlockers) > 0 {
		blockers = append(blockers, dryRunBlockers...)
	}
	if len(verifyBlockers) > 0 {
		blockers = append(blockers, verifyBlockers...)
	}
	if len(rollbackBlockers) > 0 {
		blockers = append(blockers, rollbackBlockers...)
	}
	if len(returnBlockers) > 0 {
		blockers = append(blockers, returnBlockers...)
	}
	if len(auditBlockers) > 0 {
		blockers = append(blockers, auditBlockers...)
	}
	killSwitchState := policySnapshotString(lease.PolicySnapshot, "killSwitchState")
	if gate.KillSwitchRequired && killSwitchState != "allow" {
		blockers = append(blockers, ControlledApplyBlockerKillSwitchBlocked)
	}
	if (lease.LeaseKind == LeaseKindCPUEntitlement || lease.LeaseKind == LeaseKindCombinedEnvelope) &&
		(!gate.CPUApplyEnabled || !gate.NodeActuatorApplyEnabled) {
		blockers = append(blockers, BlockerNodeFluidActuatorUnavailable)
	}
	if (lease.LeaseKind == LeaseKindRAMBalloon || lease.LeaseKind == LeaseKindCombinedEnvelope) &&
		(!gate.RAMApplyEnabled || !gate.QMPBalloonApplyEnabled) {
		blockers = append(blockers, BlockerRAMBalloonUnavailable)
	}
	return dedupe(blockers)
}

func BuildWindowsFluidControlledApplyPlan(
	target WindowsHyperdensityTarget,
	lease WindowsFluidResourceLease,
	gate WindowsFluidControlledApplyGate,
	approval *WindowsFluidManualApproval,
	evaluationTime time.Time,
) WindowsFluidControlledApplyPlan {
	ts := normalizeControlledTime(evaluationTime)
	target = EvaluateWindowsHyperdensityTarget(target)
	lease = PrepareWindowsFluidResourceLease(target, lease)
	slate := BuildWindowsFluidActionSlate(target, lease)

	compliance := EvaluateWindowsHyperdensityReadyComplianceOutput{
		CompliancePhase: target.CompliancePhase,
		Blockers:        dedupe(append([]string{}, target.Blockers...)),
		Risk:            RiskLow,
		EvidenceSummary: map[string]any{
			"vmRef":             target.VMRef,
			"namespace":         target.Namespace,
			"sameQemu":          target.Guarantees.SameQEMU,
			"sameBoot":          target.Guarantees.SameBoot,
			"cpuActuatorReady":  target.CPU.ActuatorRequired,
			"ramBalloonReady":   target.Memory.QMPRequired,
			"guestAckRequired":  target.Guest.GuestAckRequired,
		},
	}

	gateBlockers := EvaluateWindowsFluidControlledApplyGate(target, lease, gate)
	approvalState, approvalBlockers := EvaluateWindowsFluidManualApproval(gate, approval, ts)
	dryRunBlockers := EvaluateWindowsFluidDryRunGate(gate, lease, slate)
	verifyPlan, verifyBlockers := EvaluateWindowsFluidVerificationPlan(gate, lease)
	rollbackPlan, rollbackBlockers := EvaluateWindowsFluidRollbackPlan(gate, lease)
	returnPlan, returnBlockers := EvaluateWindowsFluidReturnToFloorPlan(gate, lease)
	auditPlan, auditBlockers := EvaluateWindowsFluidAuditBundlePlan(gate, lease)

	applyBlockers := EvaluateWindowsFluidApplyReadiness(
		target,
		lease,
		gate,
		approvalState,
		dryRunBlockers,
		verifyBlockers,
		rollbackBlockers,
		returnBlockers,
		auditBlockers,
	)

	allBlockers := dedupe(append(
		append(
			append(
				append(
					append(
						append(
							append(
								append([]string{}, compliance.Blockers...),
								lease.Blockers...,
							),
							gateBlockers...,
						),
						approvalBlockers...,
					),
					dryRunBlockers...,
				),
				verifyBlockers...,
			),
			rollbackBlockers...,
		),
		append(returnBlockers, append(auditBlockers, applyBlockers...)...)...,
	))

	phase := ControlledPlanPrepared
	if len(dryRunBlockers) == 0 {
		phase = ControlledPlanDryRunReady
	}
	if gate.ManualApprovalRequired && approvalState == ApprovalRequired && len(allBlockers) == 1 && allBlockers[0] == ControlledApplyBlockerManualApprovalRequired {
		phase = ControlledPlanAwaitingApproval
	}
	if contains(allBlockers, ControlledApplyBlockerApprovalRejected) || contains(allBlockers, ControlledApplyBlockerApprovalExpired) {
		phase = ControlledPlanBlocked
	}

	applyAllowed := len(allBlockers) == 0
	mutationAllowed := applyAllowed
	if applyAllowed {
		phase = ControlledPlanApplyReady
	} else if phase != ControlledPlanAwaitingApproval {
		phase = ControlledPlanApplyBlocked
	}

	if contains(allBlockers, BlockerQemuPIDChanged) || contains(allBlockers, BlockerCgroupPathMismatch) {
		phase = ControlledPlanQuarantined
	}

	actuatorReq := map[string]any{}
	if lease.LeaseKind == LeaseKindCPUEntitlement || lease.LeaseKind == LeaseKindCombinedEnvelope {
		actuatorReq = map[string]any{
			"ref":        lease.ActuatorRequestRef,
			"requested":  lease.Requested.CPUMax,
			"previous":   lease.Previous.CPUMax,
			"rollback":   lease.RollbackTarget.CPUMax,
		}
	}

	qmpReq := map[string]any{}
	if lease.LeaseKind == LeaseKindRAMBalloon || lease.LeaseKind == LeaseKindCombinedEnvelope {
		qmpReq = map[string]any{
			"ref":         lease.QMPRequestRef,
			"requested":   lease.Requested.MemoryBytes,
			"previous":    lease.Previous.MemoryBytes,
			"returnFloor": lease.ReturnToFloorTarget.MemoryBytes,
		}
	}

	return WindowsFluidControlledApplyPlan{
		PlanID:             "windows-fluid-controlled-plan-" + shortHash(target.TargetID+"|"+lease.LeaseID+"|"+ts.Format(time.RFC3339)),
		TargetRef:          target.TargetID,
		LeaseRef:           lease.LeaseID,
		GateSnapshot:       gate,
		ComplianceSnapshot: compliance,
		ActionSlate:        slate,
		ActuatorRequest:    actuatorReq,
		QMPBalloonRequest:  qmpReq,
		GuestVerifyPlan: map[string]any{
			"required":    gate.GuestVerifyRequired,
			"evidenceRef": lease.GuestEvidenceRef,
			"ready":       lease.GuestEvidenceRef != "",
		},
		WorkloadVerifyPlan: verifyPlan,
		RollbackPlan:       rollbackPlan,
		ReturnToFloorPlan:  returnPlan,
		AuditBundlePlan:    auditPlan,
		ApprovalState:      approvalState,
		PlanPhase:          phase,
		MutationAllowed:    mutationAllowed,
		ApplyAllowed:       applyAllowed,
		Blockers:           allBlockers,
		EvidenceRefs:       dedupe(append(append([]string{}, target.EvidenceRefs...), lease.EvidenceRefs...)),
		CreatedAt:          ts.Format(time.RFC3339),
	}
}

func LoadWindowsFluidControlledApplyFixture(path string) (WindowsFluidControlledApplyFixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return WindowsFluidControlledApplyFixture{}, err
	}
	var fixture WindowsFluidControlledApplyFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return WindowsFluidControlledApplyFixture{}, err
	}
	return fixture, nil
}

func containsLeaseKind(values []WindowsFluidLeaseKind, value WindowsFluidLeaseKind) bool {
	for _, current := range values {
		if current == value {
			return true
		}
	}
	return false
}

func hasActionType(slate WindowsFluidActionSlate, actionType WindowsFluidActionType) bool {
	for _, action := range slate.Actions {
		if action.ActionType == actionType {
			return true
		}
	}
	return false
}

func policySnapshotString(policy map[string]any, key string) string {
	if len(policy) == 0 {
		return ""
	}
	value, ok := policy[key]
	if !ok {
		return ""
	}
	raw, ok := value.(string)
	if !ok {
		return ""
	}
	return raw
}

func normalizeControlledTime(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now().UTC()
	}
	return value.UTC()
}
