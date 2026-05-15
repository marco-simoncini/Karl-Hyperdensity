package windowsfluidvirt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type WindowsFluidVirtControlledApplyPlanRecord struct {
	PlanRecordID                 string `json:"planRecordId"`
	PlanRecordVersion            string `json:"planRecordVersion"`
	PlanState                    string `json:"planState"`
	SourceComplianceReplayRef    string `json:"sourceComplianceReplayRef"`
	SourceReadonlyReplayRef      string `json:"sourceReadonlyReplayRef"`
	SourceActuatorBoundaryRef    string `json:"sourceActuatorBoundaryRef"`
	ControlledApplyPlanDefined   bool   `json:"controlledApplyPlanDefined"`
	ControlledApplyEnabled       bool   `json:"controlledApplyEnabled"`
	ControlledApplyExecuted      bool   `json:"controlledApplyExecuted"`
	ControlledApplyReady         bool   `json:"controlledApplyReady"`
	RuntimeMutationAllowed       bool   `json:"runtimeMutationAllowed"`
	RawControlExposureAllowed    bool   `json:"rawControlExposureAllowed"`
	EvidenceConfidence           string `json:"evidenceConfidence"`
	NextRequiredState            string `json:"nextRequiredState"`
	ClaimBoundary                string `json:"claimBoundary"`
}

type WindowsFluidVirtCandidateActionBoundary struct {
	CandidateActionState                      string `json:"candidateActionState"`
	CPUEntitlementCandidateAllowedForPlanning bool   `json:"cpuEntitlementCandidateAllowedForPlanning"`
	MemoryBalloonCandidateAllowedForPlanning  bool   `json:"memoryBalloonCandidateAllowedForPlanning"`
	CandidateExecutionAllowed                 bool   `json:"candidateExecutionAllowed"`
	CandidateRequiresManualApproval           bool   `json:"candidateRequiresManualApproval"`
	CandidateRequiresLease                    bool   `json:"candidateRequiresLease"`
	CandidateRequiresAuditChain               bool   `json:"candidateRequiresAuditChain"`
	CandidateRequiresRollback                 bool   `json:"candidateRequiresRollback"`
	CandidateRequiresReturnToFloor            bool   `json:"candidateRequiresReturnToFloor"`
	CandidateRequiresGuestWitness             bool   `json:"candidateRequiresGuestWitness"`
	CandidateRequiresSameBootProof            bool   `json:"candidateRequiresSameBootProof"`
	CandidateRequiresSameQEMUProof            bool   `json:"candidateRequiresSameQemuProof"`
	ClaimBoundary                             string `json:"claimBoundary"`
}

type WindowsFluidVirtApprovalBoundary struct {
	ManualApprovalRequired         bool   `json:"manualApprovalRequired"`
	ApprovalRecordRequired         bool   `json:"approvalRecordRequired"`
	ApprovalBeforeExecutionRequired bool  `json:"approvalBeforeExecutionRequired"`
	OperatorIdentityRequired       bool   `json:"operatorIdentityRequired"`
	AutonomousApprovalAllowed      bool   `json:"autonomousApprovalAllowed"`
	ApprovalBypassAllowed          bool   `json:"approvalBypassAllowed"`
	ClaimBoundary                  string `json:"claimBoundary"`
}

type WindowsFluidVirtLeaseBoundary struct {
	LeaseRequired                    bool   `json:"leaseRequired"`
	LeaseTTLRequired                 bool   `json:"leaseTtlRequired"`
	ActiveLeaseRequiredBeforeApply   bool   `json:"activeLeaseRequiredBeforeApply"`
	LeaseExpiryRequired              bool   `json:"leaseExpiryRequired"`
	ReturnToFloorOnExpiryRequired    bool   `json:"returnToFloorOnExpiryRequired"`
	LeaseWithoutAuditAllowed         bool   `json:"leaseWithoutAuditAllowed"`
	ClaimBoundary                    string `json:"claimBoundary"`
}

type WindowsFluidVirtRollbackBoundary struct {
	RollbackPlanRequired                 bool   `json:"rollbackPlanRequired"`
	RollbackBaselineRequired             bool   `json:"rollbackBaselineRequired"`
	RollbackVerificationRequired         bool   `json:"rollbackVerificationRequired"`
	RollbackBeforeApplyReadinessRequired bool   `json:"rollbackBeforeApplyReadinessRequired"`
	RollbackPlanProven                   bool   `json:"rollbackPlanProven"`
	RollbackExecutionAllowed             bool   `json:"rollbackExecutionAllowed"`
	ClaimBoundary                        string `json:"claimBoundary"`
}

type WindowsFluidVirtReturnToFloorBoundary struct {
	ReturnToFloorPlanRequired                 bool   `json:"returnToFloorPlanRequired"`
	ReturnToFloorVerificationRequired         bool   `json:"returnToFloorVerificationRequired"`
	ReturnToFloorBeforeApplyReadinessRequired bool   `json:"returnToFloorBeforeApplyReadinessRequired"`
	ReturnToFloorPlanProven                   bool   `json:"returnToFloorPlanProven"`
	ReturnToFloorExecutionAllowed             bool   `json:"returnToFloorExecutionAllowed"`
	ClaimBoundary                             string `json:"claimBoundary"`
}

type WindowsFluidVirtGuestWitnessBoundary struct {
	GuestWitnessRequired         bool   `json:"guestWitnessRequired"`
	FluidShellOrQGARequired      bool   `json:"fluidShellOrQgaRequired"`
	GuestAckRequired             bool   `json:"guestAckRequired"`
	SameBootProofRequired        bool   `json:"sameBootProofRequired"`
	SameQEMUProofRequired        bool   `json:"sameQemuProofRequired"`
	PendingRebootCheckRequired   bool   `json:"pendingRebootCheckRequired"`
	WitnessIntegrated            bool   `json:"witnessIntegrated"`
	WitnessBypassAllowed         bool   `json:"witnessBypassAllowed"`
	ClaimBoundary                string `json:"claimBoundary"`
}

type WindowsFluidVirtAuditChainBoundary struct {
	AuditHashChainRequired           bool   `json:"auditHashChainRequired"`
	ComplianceReplayRequired         bool   `json:"complianceReplayRequired"`
	AuditHashChainVerifiedRequired   bool   `json:"auditHashChainVerifiedRequired"`
	EveryCandidateMustBeReplayable   bool   `json:"everyCandidateMustBeReplayable"`
	EveryTransitionMustBeAudited     bool   `json:"everyTransitionMustBeAudited"`
	AuditBypassAllowed               bool   `json:"auditBypassAllowed"`
	ClaimBoundary                    string `json:"claimBoundary"`
}

type WindowsFluidVirtKillSwitchBoundary struct {
	KillSwitchRequired              bool   `json:"killSwitchRequired"`
	KillSwitchBeforeRuntimeRequired bool   `json:"killSwitchBeforeRuntimeRequired"`
	KillSwitchVerified              bool   `json:"killSwitchVerified"`
	ApplyWithoutKillSwitchAllowed   bool   `json:"applyWithoutKillSwitchAllowed"`
	ClaimBoundary                   string `json:"claimBoundary"`
}

type WindowsFluidVirtRuntimeExecutionBoundary struct {
	RuntimeExecutionAllowed     bool   `json:"runtimeExecutionAllowed"`
	ExecutorEnabled             bool   `json:"executorEnabled"`
	ActuatorRuntimeEnabled      bool   `json:"actuatorRuntimeEnabled"`
	CgroupWriteEnabled          bool   `json:"cgroupWriteEnabled"`
	QMPCommandExecutionAllowed  bool   `json:"qmpCommandExecutionAllowed"`
	QGACommandExecutionAllowed  bool   `json:"qgaCommandExecutionAllowed"`
	RawRuntimeControlsExposed   bool   `json:"rawRuntimeControlsExposed"`
	ProductionMutationAllowed   bool   `json:"productionMutationAllowed"`
	AutonomousApplyAllowed      bool   `json:"autonomousApplyAllowed"`
	ClaimBoundary               string `json:"claimBoundary"`
}

type WindowsFluidVirtControlledApplyReadinessScorecard struct {
	ScorecardState             string `json:"scorecardState"`
	ControlledApplyPlanDefined bool   `json:"controlledApplyPlanDefined"`
	ControlledApplyEnabled     bool   `json:"controlledApplyEnabled"`
	ControlledApplyExecuted    bool   `json:"controlledApplyExecuted"`
	ControlledApplyReady       bool   `json:"controlledApplyReady"`
	RuntimeMutationEnabled     bool   `json:"runtimeMutationEnabled"`
	ExecutorEnabled            bool   `json:"executorEnabled"`
	ActuatorRuntimeEnabled     bool   `json:"actuatorRuntimeEnabled"`
	ComplianceReplayAvailable  bool   `json:"complianceReplayAvailable"`
	AuditHashChainVerified     bool   `json:"auditHashChainVerified"`
	BlockingReasonCount        int    `json:"blockingReasonCount"`
	Confidence                 string `json:"confidence"`
	Decision                   string `json:"decision"`
	DecisionReason             string `json:"decisionReason"`
}

type WindowsFluidVirtControlledApplyPlanBoundary struct {
	ControlledApplyPlanID                   string                                      `json:"controlledApplyPlanId"`
	ControlledApplyPlanVersion              string                                      `json:"controlledApplyPlanVersion"`
	ReleaseTrack                            string                                      `json:"releaseTrack"`
	LaneStatus                              string                                      `json:"laneStatus"`
	PlanMode                                string                                      `json:"planMode"`
	SourceProductModelRef                   string                                      `json:"sourceProductModelRef"`
	SourceNodeActuatorBoundaryRef           string                                      `json:"sourceNodeActuatorBoundaryRef"`
	SourceNodeActuatorReadonlyReplayRef     string                                      `json:"sourceNodeActuatorReadonlyReplayRef"`
	SourceComplianceReplayRef               string                                      `json:"sourceComplianceReplayRef"`
	ControlledApplyPlanDefined              bool                                        `json:"controlledApplyPlanDefined"`
	ControlledApplyEnabled                  bool                                        `json:"controlledApplyEnabled"`
	ControlledApplyExecuted                 bool                                        `json:"controlledApplyExecuted"`
	ControlledApplyReady                    bool                                        `json:"controlledApplyReady"`
	RuntimeMutationEnabled                  bool                                        `json:"runtimeMutationEnabled"`
	ActuatorRuntimeEnabled                  bool                                        `json:"actuatorRuntimeEnabled"`
	ExecutorEnabled                         bool                                        `json:"executorEnabled"`
	CgroupWriteEnabled                      bool                                        `json:"cgroupWriteEnabled"`
	QMPCommandExecutionAllowed              bool                                        `json:"qmpCommandExecutionAllowed"`
	QGACommandExecutionAllowed              bool                                        `json:"qgaCommandExecutionAllowed"`
	RawRuntimeControlsExposed               bool                                        `json:"rawRuntimeControlsExposed"`
	ProductionMutationAllowed               bool                                        `json:"productionMutationAllowed"`
	AutonomousApplyAllowed                  bool                                        `json:"autonomousApplyAllowed"`
	EnforcementMode                         string                                      `json:"enforcementMode"`
	WindowsGaClaimAllowed                   bool                                        `json:"windowsGaClaimAllowed"`
	WindowsProductionReadyClaimAllowed      bool                                        `json:"windowsProductionReadyClaimAllowed"`
	WindowsExecutionReadyByDefault          bool                                        `json:"windowsExecutionReadyByDefault"`
	VCPUHotplugClaimAllowed                 bool                                        `json:"vcpuHotplugClaimAllowed"`
	LogicalCPUScalingClaimAllowed           bool                                        `json:"logicalCpuScalingClaimAllowed"`
	PoolScalingClaimAllowed                 bool                                        `json:"poolScalingClaimAllowed"`
	LiveMigrationClaimAllowed               bool                                        `json:"liveMigrationClaimAllowed"`
	RebootRecreateRolloutMechanismAllowed   bool                                        `json:"rebootRecreateRolloutMechanismAllowed"`
	ControlledApplyPlanRecord               WindowsFluidVirtControlledApplyPlanRecord   `json:"controlledApplyPlanRecord"`
	CandidateActionBoundary                 WindowsFluidVirtCandidateActionBoundary     `json:"candidateActionBoundary"`
	ApprovalBoundary                        WindowsFluidVirtApprovalBoundary            `json:"approvalBoundary"`
	LeaseBoundary                           WindowsFluidVirtLeaseBoundary               `json:"leaseBoundary"`
	RollbackBoundary                        WindowsFluidVirtRollbackBoundary            `json:"rollbackBoundary"`
	ReturnToFloorBoundary                   WindowsFluidVirtReturnToFloorBoundary       `json:"returnToFloorBoundary"`
	GuestWitnessBoundary                    WindowsFluidVirtGuestWitnessBoundary        `json:"guestWitnessBoundary"`
	AuditChainBoundary                      WindowsFluidVirtAuditChainBoundary          `json:"auditChainBoundary"`
	KillSwitchBoundary                      WindowsFluidVirtKillSwitchBoundary          `json:"killSwitchBoundary"`
	RuntimeExecutionBoundary                WindowsFluidVirtRuntimeExecutionBoundary    `json:"runtimeExecutionBoundary"`
	RequiredBeforeControlledApply           []string                                    `json:"requiredBeforeControlledApply"`
	ControlledApplyBlockingReasons          []string                                    `json:"controlledApplyBlockingReasons"`
	ControlledApplyReadinessScorecard       WindowsFluidVirtControlledApplyReadinessScorecard `json:"controlledApplyReadinessScorecard"`
	AllowedControlledApplyPlanActions       []string                                    `json:"allowedControlledApplyPlanActions"`
	ForbiddenControlledApplyPlanActions     []string                                    `json:"forbiddenControlledApplyPlanActions"`
	ClaimBoundaries                         []string                                    `json:"claimBoundaries"`
	EvidenceRefs                            []string                                    `json:"evidenceRefs"`
	AuditTrail                              []string                                    `json:"auditTrail"`
}

func NewWindowsFluidVirtControlledApplyPlanBoundaryMinimal() WindowsFluidVirtControlledApplyPlanBoundary {
	blockingReasons := []string{
		"controlled_apply_boundary_only",
		"executor_not_ported",
		"actuator_runtime_not_enabled",
		"cgroup_write_not_enabled",
		"manual_approval_flow_missing",
		"active_lease_model_missing",
		"guest_witness_missing",
		"rollback_plan_not_proven",
		"return_to_floor_plan_not_proven",
		"kill_switch_not_verified",
		"node_allowlist_missing",
		"cgroup_path_validation_missing",
		"production_ready_claim_forbidden",
		"autonomous_apply_forbidden",
		"raw_runtime_control_forbidden",
	}
	return WindowsFluidVirtControlledApplyPlanBoundary{
		ControlledApplyPlanID:                 "windows_fluidvirt_controlled_apply_plan_boundary_v1",
		ControlledApplyPlanVersion:            "v1",
		ReleaseTrack:                          "technical_preview",
		LaneStatus:                            "gated_preview",
		PlanMode:                              "controlled_apply_plan_boundary_only",
		SourceProductModelRef:                 "windows-fluidvirt-product-model-v1",
		SourceNodeActuatorBoundaryRef:         "windows_fluidvirt_node_actuator_contract_boundary_v1",
		SourceNodeActuatorReadonlyReplayRef:   "windows_fluidvirt_node_actuator_mvp_readonly_replay_v1",
		SourceComplianceReplayRef:             "windows_fluidvirt_compliance_replay_audit_chain_v1",
		ControlledApplyPlanDefined:            true,
		ControlledApplyEnabled:                false,
		ControlledApplyExecuted:               false,
		ControlledApplyReady:                  false,
		RuntimeMutationEnabled:                false,
		ActuatorRuntimeEnabled:                false,
		ExecutorEnabled:                       false,
		CgroupWriteEnabled:                    false,
		QMPCommandExecutionAllowed:            false,
		QGACommandExecutionAllowed:            false,
		RawRuntimeControlsExposed:             false,
		ProductionMutationAllowed:             false,
		AutonomousApplyAllowed:                false,
		EnforcementMode:                       "disabled",
		WindowsGaClaimAllowed:                 false,
		WindowsProductionReadyClaimAllowed:    false,
		WindowsExecutionReadyByDefault:        false,
		VCPUHotplugClaimAllowed:               false,
		LogicalCPUScalingClaimAllowed:         false,
		PoolScalingClaimAllowed:               false,
		LiveMigrationClaimAllowed:             false,
		RebootRecreateRolloutMechanismAllowed: false,
		ControlledApplyPlanRecord: WindowsFluidVirtControlledApplyPlanRecord{
			PlanRecordID:               "controlled_apply_plan_record_v1",
			PlanRecordVersion:          "v1",
			PlanState:                  "defined_not_enabled",
			SourceComplianceReplayRef:  "windows_fluidvirt_compliance_replay_audit_chain_v1",
			SourceReadonlyReplayRef:    "windows_fluidvirt_node_actuator_mvp_readonly_replay_v1",
			SourceActuatorBoundaryRef:  "windows_fluidvirt_node_actuator_contract_boundary_v1",
			ControlledApplyPlanDefined: true,
			ControlledApplyEnabled:     false,
			ControlledApplyExecuted:    false,
			ControlledApplyReady:       false,
			RuntimeMutationAllowed:     false,
			RawControlExposureAllowed:  false,
			EvidenceConfidence:         "low",
			NextRequiredState:          "guarded_fake_runtime_or_inventory_witness_integration",
			ClaimBoundary:              "boundary_only",
		},
		CandidateActionBoundary: WindowsFluidVirtCandidateActionBoundary{
			CandidateActionState:                      "planning_only",
			CPUEntitlementCandidateAllowedForPlanning: true,
			MemoryBalloonCandidateAllowedForPlanning:  true,
			CandidateExecutionAllowed:                 false,
			CandidateRequiresManualApproval:           true,
			CandidateRequiresLease:                    true,
			CandidateRequiresAuditChain:               true,
			CandidateRequiresRollback:                 true,
			CandidateRequiresReturnToFloor:            true,
			CandidateRequiresGuestWitness:             true,
			CandidateRequiresSameBootProof:            true,
			CandidateRequiresSameQEMUProof:            true,
			ClaimBoundary:                             "planning_only_no_execution",
		},
		ApprovalBoundary: WindowsFluidVirtApprovalBoundary{
			ManualApprovalRequired:          true,
			ApprovalRecordRequired:          true,
			ApprovalBeforeExecutionRequired: true,
			OperatorIdentityRequired:        true,
			AutonomousApprovalAllowed:       false,
			ApprovalBypassAllowed:           false,
			ClaimBoundary:                   "manual_approval_required",
		},
		LeaseBoundary: WindowsFluidVirtLeaseBoundary{
			LeaseRequired:                  true,
			LeaseTTLRequired:               true,
			ActiveLeaseRequiredBeforeApply: true,
			LeaseExpiryRequired:            true,
			ReturnToFloorOnExpiryRequired:  true,
			LeaseWithoutAuditAllowed:       false,
			ClaimBoundary:                  "lease_required_before_apply",
		},
		RollbackBoundary: WindowsFluidVirtRollbackBoundary{
			RollbackPlanRequired:                 true,
			RollbackBaselineRequired:             true,
			RollbackVerificationRequired:         true,
			RollbackBeforeApplyReadinessRequired: true,
			RollbackPlanProven:                   false,
			RollbackExecutionAllowed:             false,
			ClaimBoundary:                        "rollback_required_not_proven",
		},
		ReturnToFloorBoundary: WindowsFluidVirtReturnToFloorBoundary{
			ReturnToFloorPlanRequired:                 true,
			ReturnToFloorVerificationRequired:         true,
			ReturnToFloorBeforeApplyReadinessRequired: true,
			ReturnToFloorPlanProven:                   false,
			ReturnToFloorExecutionAllowed:             false,
			ClaimBoundary:                             "return_to_floor_required_not_proven",
		},
		GuestWitnessBoundary: WindowsFluidVirtGuestWitnessBoundary{
			GuestWitnessRequired:       true,
			FluidShellOrQGARequired:    true,
			GuestAckRequired:           true,
			SameBootProofRequired:      true,
			SameQEMUProofRequired:      true,
			PendingRebootCheckRequired: true,
			WitnessIntegrated:          false,
			WitnessBypassAllowed:       false,
			ClaimBoundary:              "guest_witness_required_not_integrated",
		},
		AuditChainBoundary: WindowsFluidVirtAuditChainBoundary{
			AuditHashChainRequired:         true,
			ComplianceReplayRequired:       true,
			AuditHashChainVerifiedRequired: true,
			EveryCandidateMustBeReplayable: true,
			EveryTransitionMustBeAudited:   true,
			AuditBypassAllowed:             false,
			ClaimBoundary:                  "audit_chain_required",
		},
		KillSwitchBoundary: WindowsFluidVirtKillSwitchBoundary{
			KillSwitchRequired:              true,
			KillSwitchBeforeRuntimeRequired: true,
			KillSwitchVerified:              false,
			ApplyWithoutKillSwitchAllowed:   false,
			ClaimBoundary:                   "kill_switch_required_not_verified",
		},
		RuntimeExecutionBoundary: WindowsFluidVirtRuntimeExecutionBoundary{
			RuntimeExecutionAllowed:   false,
			ExecutorEnabled:           false,
			ActuatorRuntimeEnabled:    false,
			CgroupWriteEnabled:        false,
			QMPCommandExecutionAllowed: false,
			QGACommandExecutionAllowed: false,
			RawRuntimeControlsExposed: false,
			ProductionMutationAllowed: false,
			AutonomousApplyAllowed:    false,
			ClaimBoundary:             "runtime_execution_forbidden",
		},
		RequiredBeforeControlledApply: []string{
			"manual_approval_flow",
			"active_lease_model",
			"lease_ttl",
			"lease_expiry_return_to_floor",
			"guest_witness_integration",
			"same_boot_proof",
			"same_qemu_proof",
			"rollback_plan",
			"rollback_baseline",
			"return_to_floor_plan",
			"audit_hash_chain_verified",
			"compliance_replay_for_every_candidate",
			"kill_switch",
			"node_allowlist",
			"cgroup_path_validation_strategy",
			"no_raw_runtime_control_surface",
			"no_autonomous_apply",
			"no_production_apply",
			"operator_identity_audit",
		},
		ControlledApplyBlockingReasons: blockingReasons,
		ControlledApplyReadinessScorecard: WindowsFluidVirtControlledApplyReadinessScorecard{
			ScorecardState:             "controlled_apply_plan_defined_not_ready",
			ControlledApplyPlanDefined: true,
			ControlledApplyEnabled:     false,
			ControlledApplyExecuted:    false,
			ControlledApplyReady:       false,
			RuntimeMutationEnabled:     false,
			ExecutorEnabled:            false,
			ActuatorRuntimeEnabled:     false,
			ComplianceReplayAvailable:  true,
			AuditHashChainVerified:     true,
			BlockingReasonCount:        len(blockingReasons),
			Confidence:                 "low",
			Decision:                   "define_plan_boundary_without_apply_enablement",
			DecisionReason:             "boundary defined while runtime, witness, and rollback proofs remain incomplete",
		},
		AllowedControlledApplyPlanActions: []string{
			"define_controlled_apply_plan_boundary",
			"record_candidate_action_boundary",
			"record_approval_boundary",
			"record_lease_boundary",
			"record_rollback_boundary",
			"record_return_to_floor_boundary",
			"record_guest_witness_boundary",
			"record_audit_chain_boundary",
			"record_kill_switch_boundary",
			"record_runtime_execution_boundary",
			"record_required_before_controlled_apply",
			"record_controlled_apply_blocking_reasons",
		},
		ForbiddenControlledApplyPlanActions: []string{
			"enable_controlled_apply",
			"execute_controlled_apply",
			"enable_executor",
			"execute_windows_executor",
			"execute_node_actuator_runtime",
			"write_cgroup_cpu_max",
			"mutate_qmp_balloon",
			"execute_qmp_command",
			"execute_qga_command",
			"touch_real_cgroup",
			"touch_real_qmp",
			"touch_real_qga",
			"enable_runtime_actuator",
			"enable_autonomous_apply",
			"enable_production_auto",
			"mark_windows_execution_ready",
			"claim_windows_ga",
			"claim_windows_production_ready",
			"expose_raw_runtime_controls",
			"create_secret",
			"commit_secret",
			"log_token",
			"package_os_iso",
			"touch_dashboard",
			"touch_inventory",
		},
		ClaimBoundaries: []string{
			"controlled apply plan defined does not imply controlled apply readiness",
			"no runtime execution enablement",
			"not windows ga",
			"not windows production ready",
		},
		EvidenceRefs: []string{
			"contract://windows_fluidvirt_node_actuator_contract_boundary_v1",
			"replay://windows_fluidvirt_node_actuator_mvp_readonly_replay_v1",
			"compliance://windows_fluidvirt_compliance_replay_audit_chain_v1",
			"plan://windows_fluidvirt_controlled_apply_plan_boundary_v1",
		},
		AuditTrail: []string{
			"audit://controlled-apply-plan-boundary-defined",
			"audit://blocking-reasons-recorded",
			"audit://apply-enablement-forbidden",
		},
	}
}

func ValidateWindowsFluidVirtControlledApplyPlanBoundary(plan WindowsFluidVirtControlledApplyPlanBoundary) error {
	if plan.PlanMode != "controlled_apply_plan_boundary_only" || !plan.ControlledApplyPlanDefined {
		return fmt.Errorf("plan must remain boundary-only and defined")
	}
	if plan.ControlledApplyEnabled || plan.ControlledApplyExecuted || plan.ControlledApplyReady {
		return fmt.Errorf("controlled apply must remain disabled and not ready")
	}
	if plan.ExecutorEnabled || plan.ActuatorRuntimeEnabled || plan.RuntimeMutationEnabled || plan.CgroupWriteEnabled {
		return fmt.Errorf("runtime execution and mutation surfaces must remain disabled")
	}
	if plan.QMPCommandExecutionAllowed || plan.QGACommandExecutionAllowed || plan.RawRuntimeControlsExposed {
		return fmt.Errorf("raw runtime controls and qmp/qga execution must remain disabled")
	}
	if plan.AutonomousApplyAllowed || plan.ProductionMutationAllowed {
		return fmt.Errorf("autonomous and production mutation must remain disabled")
	}
	if plan.WindowsGaClaimAllowed || plan.WindowsProductionReadyClaimAllowed || plan.WindowsExecutionReadyByDefault {
		return fmt.Errorf("windows ga/production/execution-ready claims must remain disabled")
	}
	if plan.ControlledApplyReadinessScorecard.ScorecardState != "controlled_apply_plan_defined_not_ready" {
		return fmt.Errorf("readiness scorecard must remain not-ready")
	}
	if len(plan.RequiredBeforeControlledApply) == 0 || len(plan.ControlledApplyBlockingReasons) == 0 {
		return fmt.Errorf("required-before and blocking reasons must be populated")
	}
	return nil
}

func LoadControlledApplyPlanBoundaryFixtureFromTemporaryFile(path string) (WindowsFluidVirtControlledApplyPlanBoundary, error) {
	clean := filepath.Clean(path)
	if strings.HasPrefix(clean, "/sys/fs/cgroup") {
		return WindowsFluidVirtControlledApplyPlanBoundary{}, fmt.Errorf("real cgroup path forbidden")
	}
	if !strings.Contains(clean, os.TempDir()) {
		return WindowsFluidVirtControlledApplyPlanBoundary{}, fmt.Errorf("fixture must be loaded from temporary path")
	}
	raw, err := os.ReadFile(clean)
	if err != nil {
		return WindowsFluidVirtControlledApplyPlanBoundary{}, err
	}
	var plan WindowsFluidVirtControlledApplyPlanBoundary
	if err := json.Unmarshal(raw, &plan); err != nil {
		return WindowsFluidVirtControlledApplyPlanBoundary{}, err
	}
	if err := ValidateWindowsFluidVirtControlledApplyPlanBoundary(plan); err != nil {
		return WindowsFluidVirtControlledApplyPlanBoundary{}, err
	}
	return plan, nil
}
