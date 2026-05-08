package windowsfluidvirt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type WindowsFluidVirtExecutorBoundaryRecord struct {
	ExecutorBoundaryRecordID             string `json:"executorBoundaryRecordId"`
	ExecutorBoundaryRecordVersion        string `json:"executorBoundaryRecordVersion"`
	BoundaryState                        string `json:"boundaryState"`
	SourceControlledApplyPlanBoundaryRef string `json:"sourceControlledApplyPlanBoundaryRef"`
	SourceComplianceReplayRef            string `json:"sourceComplianceReplayRef"`
	SourceReadonlyReplayRef              string `json:"sourceReadonlyReplayRef"`
	ExecutorBoundaryDefined              bool   `json:"executorBoundaryDefined"`
	ExecutorEnabled                      bool   `json:"executorEnabled"`
	ExecutorRuntimeAvailable             bool   `json:"executorRuntimeAvailable"`
	ExecutorExecuted                     bool   `json:"executorExecuted"`
	RuntimeMutationAllowed               bool   `json:"runtimeMutationAllowed"`
	RawControlExposureAllowed            bool   `json:"rawControlExposureAllowed"`
	EvidenceConfidence                   string `json:"evidenceConfidence"`
	NextRequiredState                    string `json:"nextRequiredState"`
	ClaimBoundary                        string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorInputBoundary struct {
	InputBoundaryState                    string   `json:"inputBoundaryState"`
	AllowedInputTypes                     []string `json:"allowedInputTypes"`
	ForbiddenInputTypes                   []string `json:"forbiddenInputTypes"`
	CandidateActionInputAllowedForPlanning bool    `json:"candidateActionInputAllowedForPlanning"`
	RealRuntimeCommandInputAllowed        bool     `json:"realRuntimeCommandInputAllowed"`
	RawCgroupPathInputAllowed             bool     `json:"rawCgroupPathInputAllowed"`
	RawQMPCommandInputAllowed             bool     `json:"rawQmpCommandInputAllowed"`
	RawQGACommandInputAllowed             bool     `json:"rawQgaCommandInputAllowed"`
	SecretMaterialInputAllowed            bool     `json:"secretMaterialInputAllowed"`
	RequiresComplianceReplay              bool     `json:"requiresComplianceReplay"`
	RequiresAuditHashChain                bool     `json:"requiresAuditHashChain"`
	RequiresManualApproval                bool     `json:"requiresManualApproval"`
	RequiresActiveLease                   bool     `json:"requiresActiveLease"`
	ClaimBoundary                         string   `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorOutputBoundary struct {
	OutputBoundaryState                 string   `json:"outputBoundaryState"`
	AllowedOutputTypes                  []string `json:"allowedOutputTypes"`
	ForbiddenOutputTypes                []string `json:"forbiddenOutputTypes"`
	RuntimeMutationOutputAllowed        bool     `json:"runtimeMutationOutputAllowed"`
	RawRuntimeControlOutputAllowed      bool     `json:"rawRuntimeControlOutputAllowed"`
	SecretOutputAllowed                 bool     `json:"secretOutputAllowed"`
	AuditEventOutputRequired            bool     `json:"auditEventOutputRequired"`
	DeterministicReplayOutputRequired   bool     `json:"deterministicReplayOutputRequired"`
	ClaimBoundary                       string   `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorAuthorizationBoundary struct {
	ManualApprovalRequired          bool   `json:"manualApprovalRequired"`
	OperatorIdentityRequired        bool   `json:"operatorIdentityRequired"`
	ApprovalRecordRequired          bool   `json:"approvalRecordRequired"`
	ApprovalBeforeExecutionRequired bool   `json:"approvalBeforeExecutionRequired"`
	AutonomousApprovalAllowed       bool   `json:"autonomousApprovalAllowed"`
	ApprovalBypassAllowed           bool   `json:"approvalBypassAllowed"`
	ClaimBoundary                   string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorLeaseBoundary struct {
	ActiveLeaseRequired               bool   `json:"activeLeaseRequired"`
	LeaseTTLRequired                  bool   `json:"leaseTtlRequired"`
	LeaseExpiryReturnToFloorRequired  bool   `json:"leaseExpiryReturnToFloorRequired"`
	LeaseWithoutAuditAllowed          bool   `json:"leaseWithoutAuditAllowed"`
	LeaseBypassAllowed                bool   `json:"leaseBypassAllowed"`
	ClaimBoundary                     string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorWitnessBoundary struct {
	GuestWitnessRequired       bool   `json:"guestWitnessRequired"`
	FluidShellOrQGARequired    bool   `json:"fluidShellOrQgaRequired"`
	GuestAckRequired           bool   `json:"guestAckRequired"`
	SameBootProofRequired      bool   `json:"sameBootProofRequired"`
	SameQEMUProofRequired      bool   `json:"sameQemuProofRequired"`
	PendingRebootCheckRequired bool   `json:"pendingRebootCheckRequired"`
	WitnessIntegrated          bool   `json:"witnessIntegrated"`
	WitnessBypassAllowed       bool   `json:"witnessBypassAllowed"`
	ClaimBoundary              string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorRollbackBoundary struct {
	RollbackPlanRequired         bool   `json:"rollbackPlanRequired"`
	RollbackBaselineRequired     bool   `json:"rollbackBaselineRequired"`
	RollbackVerificationRequired bool   `json:"rollbackVerificationRequired"`
	RollbackExecutionAllowed     bool   `json:"rollbackExecutionAllowed"`
	RollbackBypassAllowed        bool   `json:"rollbackBypassAllowed"`
	ClaimBoundary                string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorReturnToFloorBoundary struct {
	ReturnToFloorPlanRequired         bool   `json:"returnToFloorPlanRequired"`
	ReturnToFloorVerificationRequired bool   `json:"returnToFloorVerificationRequired"`
	ReturnToFloorExecutionAllowed     bool   `json:"returnToFloorExecutionAllowed"`
	ReturnToFloorBypassAllowed        bool   `json:"returnToFloorBypassAllowed"`
	ClaimBoundary                     string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorAuditBoundary struct {
	AuditHashChainRequired         bool   `json:"auditHashChainRequired"`
	ComplianceReplayRequired       bool   `json:"complianceReplayRequired"`
	AuditHashChainVerifiedRequired bool   `json:"auditHashChainVerifiedRequired"`
	EveryInputMustBeReplayable     bool   `json:"everyInputMustBeReplayable"`
	EveryOutputMustBeAudited       bool   `json:"everyOutputMustBeAudited"`
	AuditBypassAllowed             bool   `json:"auditBypassAllowed"`
	ClaimBoundary                  string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorKillSwitchBoundary struct {
	KillSwitchRequired               bool   `json:"killSwitchRequired"`
	KillSwitchBeforeRuntimeRequired  bool   `json:"killSwitchBeforeRuntimeRequired"`
	KillSwitchVerified               bool   `json:"killSwitchVerified"`
	ExecutionWithoutKillSwitchAllowed bool  `json:"executionWithoutKillSwitchAllowed"`
	ClaimBoundary                    string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorRuntimeBoundary struct {
	RuntimeExecutionAllowed     bool   `json:"runtimeExecutionAllowed"`
	ExecutorEnabled             bool   `json:"executorEnabled"`
	ActuatorRuntimeEnabled      bool   `json:"actuatorRuntimeEnabled"`
	ControlledApplyEnabled      bool   `json:"controlledApplyEnabled"`
	CgroupWriteEnabled          bool   `json:"cgroupWriteEnabled"`
	QMPCommandExecutionAllowed  bool   `json:"qmpCommandExecutionAllowed"`
	QGACommandExecutionAllowed  bool   `json:"qgaCommandExecutionAllowed"`
	RawRuntimeControlsExposed   bool   `json:"rawRuntimeControlsExposed"`
	ProductionMutationAllowed   bool   `json:"productionMutationAllowed"`
	AutonomousApplyAllowed      bool   `json:"autonomousApplyAllowed"`
	ClaimBoundary               string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorReadinessScorecard struct {
	ScorecardState             string `json:"scorecardState"`
	ExecutorBoundaryDefined    bool   `json:"executorBoundaryDefined"`
	ExecutorEnabled            bool   `json:"executorEnabled"`
	ExecutorRuntimeAvailable   bool   `json:"executorRuntimeAvailable"`
	ExecutorExecuted           bool   `json:"executorExecuted"`
	ControlledApplyEnabled     bool   `json:"controlledApplyEnabled"`
	ControlledApplyReady       bool   `json:"controlledApplyReady"`
	RuntimeMutationEnabled     bool   `json:"runtimeMutationEnabled"`
	ActuatorRuntimeEnabled     bool   `json:"actuatorRuntimeEnabled"`
	ComplianceReplayAvailable  bool   `json:"complianceReplayAvailable"`
	AuditHashChainVerified     bool   `json:"auditHashChainVerified"`
	BlockingReasonCount        int    `json:"blockingReasonCount"`
	Confidence                 string `json:"confidence"`
	Decision                   string `json:"decision"`
	DecisionReason             string `json:"decisionReason"`
}

type WindowsFluidVirtGuardedExecutorBoundary struct {
	GuardedExecutorBoundaryID             string                                       `json:"guardedExecutorBoundaryId"`
	GuardedExecutorBoundaryVersion        string                                       `json:"guardedExecutorBoundaryVersion"`
	ReleaseTrack                          string                                       `json:"releaseTrack"`
	LaneStatus                            string                                       `json:"laneStatus"`
	BoundaryMode                          string                                       `json:"boundaryMode"`
	SourceProductModelRef                 string                                       `json:"sourceProductModelRef"`
	SourceNodeActuatorBoundaryRef         string                                       `json:"sourceNodeActuatorBoundaryRef"`
	SourceNodeActuatorReadonlyReplayRef   string                                       `json:"sourceNodeActuatorReadonlyReplayRef"`
	SourceComplianceReplayRef             string                                       `json:"sourceComplianceReplayRef"`
	SourceControlledApplyPlanBoundaryRef  string                                       `json:"sourceControlledApplyPlanBoundaryRef"`
	ExecutorBoundaryDefined               bool                                         `json:"executorBoundaryDefined"`
	ExecutorEnabled                       bool                                         `json:"executorEnabled"`
	ExecutorRuntimeAvailable              bool                                         `json:"executorRuntimeAvailable"`
	ExecutorExecuted                      bool                                         `json:"executorExecuted"`
	ControlledApplyEnabled                bool                                         `json:"controlledApplyEnabled"`
	ControlledApplyExecuted               bool                                         `json:"controlledApplyExecuted"`
	ControlledApplyReady                  bool                                         `json:"controlledApplyReady"`
	RuntimeMutationEnabled                bool                                         `json:"runtimeMutationEnabled"`
	ActuatorRuntimeEnabled                bool                                         `json:"actuatorRuntimeEnabled"`
	CgroupWriteEnabled                    bool                                         `json:"cgroupWriteEnabled"`
	QMPCommandExecutionAllowed            bool                                         `json:"qmpCommandExecutionAllowed"`
	QGACommandExecutionAllowed            bool                                         `json:"qgaCommandExecutionAllowed"`
	RawRuntimeControlsExposed             bool                                         `json:"rawRuntimeControlsExposed"`
	ProductionMutationAllowed             bool                                         `json:"productionMutationAllowed"`
	AutonomousApplyAllowed                bool                                         `json:"autonomousApplyAllowed"`
	EnforcementMode                       string                                       `json:"enforcementMode"`
	WindowsGaClaimAllowed                 bool                                         `json:"windowsGaClaimAllowed"`
	WindowsProductionReadyClaimAllowed    bool                                         `json:"windowsProductionReadyClaimAllowed"`
	WindowsExecutionReadyByDefault        bool                                         `json:"windowsExecutionReadyByDefault"`
	VCPUHotplugClaimAllowed               bool                                         `json:"vcpuHotplugClaimAllowed"`
	LogicalCPUScalingClaimAllowed         bool                                         `json:"logicalCpuScalingClaimAllowed"`
	PoolScalingClaimAllowed               bool                                         `json:"poolScalingClaimAllowed"`
	LiveMigrationClaimAllowed             bool                                         `json:"liveMigrationClaimAllowed"`
	RebootRecreateRolloutMechanismAllowed bool                                         `json:"rebootRecreateRolloutMechanismAllowed"`
	ExecutorBoundaryRecord                WindowsFluidVirtExecutorBoundaryRecord       `json:"executorBoundaryRecord"`
	ExecutorInputBoundary                 WindowsFluidVirtExecutorInputBoundary        `json:"executorInputBoundary"`
	ExecutorOutputBoundary                WindowsFluidVirtExecutorOutputBoundary       `json:"executorOutputBoundary"`
	ExecutorAuthorizationBoundary         WindowsFluidVirtExecutorAuthorizationBoundary `json:"executorAuthorizationBoundary"`
	ExecutorLeaseBoundary                 WindowsFluidVirtExecutorLeaseBoundary        `json:"executorLeaseBoundary"`
	ExecutorWitnessBoundary               WindowsFluidVirtExecutorWitnessBoundary      `json:"executorWitnessBoundary"`
	ExecutorRollbackBoundary              WindowsFluidVirtExecutorRollbackBoundary     `json:"executorRollbackBoundary"`
	ExecutorReturnToFloorBoundary         WindowsFluidVirtExecutorReturnToFloorBoundary `json:"executorReturnToFloorBoundary"`
	ExecutorAuditBoundary                 WindowsFluidVirtExecutorAuditBoundary        `json:"executorAuditBoundary"`
	ExecutorKillSwitchBoundary            WindowsFluidVirtExecutorKillSwitchBoundary   `json:"executorKillSwitchBoundary"`
	ExecutorRuntimeBoundary               WindowsFluidVirtExecutorRuntimeBoundary      `json:"executorRuntimeBoundary"`
	RequiredBeforeExecutorEnablement      []string                                     `json:"requiredBeforeExecutorEnablement"`
	ExecutorBlockingReasons               []string                                     `json:"executorBlockingReasons"`
	ExecutorReadinessScorecard            WindowsFluidVirtExecutorReadinessScorecard   `json:"executorReadinessScorecard"`
	AllowedExecutorBoundaryActions        []string                                     `json:"allowedExecutorBoundaryActions"`
	ForbiddenExecutorBoundaryActions      []string                                     `json:"forbiddenExecutorBoundaryActions"`
	ClaimBoundaries                       []string                                     `json:"claimBoundaries"`
	EvidenceRefs                          []string                                     `json:"evidenceRefs"`
	AuditTrail                            []string                                     `json:"auditTrail"`
}

func NewWindowsFluidVirtGuardedExecutorBoundaryMinimal() WindowsFluidVirtGuardedExecutorBoundary {
	allowedInput := []string{
		"controlled_apply_candidate_model",
		"compliance_replay_verified_candidate",
		"audit_hash_chain_verified_candidate",
		"operator_approved_candidate_reference",
		"lease_bound_candidate_reference",
	}
	forbiddenInput := []string{
		"raw_cgroup_path",
		"raw_qmp_command",
		"raw_qga_command",
		"raw_shell_command",
		"raw_secret",
		"kubeconfig",
		"token",
		"unapproved_candidate",
		"unaudited_candidate",
		"production_auto_candidate",
	}
	allowedOutput := []string{
		"boundary_validation_result",
		"executor_blocking_reasons",
		"required_before_executor_enablement",
		"audit_event_model",
		"no_runtime_mutation_attestation",
	}
	forbiddenOutput := []string{
		"runtime_apply_result",
		"cgroup_write_result",
		"qmp_command_result",
		"qga_command_result",
		"raw_secret_output",
		"production_auto_result",
	}
	blockingReasons := []string{
		"executor_boundary_only",
		"executor_runtime_not_ported",
		"executor_enabled_false",
		"controlled_apply_not_ready",
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
	return WindowsFluidVirtGuardedExecutorBoundary{
		GuardedExecutorBoundaryID:             "windows_fluidvirt_guarded_executor_boundary_v1",
		GuardedExecutorBoundaryVersion:        "v1",
		ReleaseTrack:                          "technical_preview",
		LaneStatus:                            "gated_preview",
		BoundaryMode:                          "guarded_executor_boundary_only",
		SourceProductModelRef:                 "windows-fluidvirt-product-model-v1",
		SourceNodeActuatorBoundaryRef:         "windows_fluidvirt_node_actuator_contract_boundary_v1",
		SourceNodeActuatorReadonlyReplayRef:   "windows_fluidvirt_node_actuator_mvp_readonly_replay_v1",
		SourceComplianceReplayRef:             "windows_fluidvirt_compliance_replay_audit_chain_v1",
		SourceControlledApplyPlanBoundaryRef:  "windows_fluidvirt_controlled_apply_plan_boundary_v1",
		ExecutorBoundaryDefined:               true,
		ExecutorEnabled:                       false,
		ExecutorRuntimeAvailable:              false,
		ExecutorExecuted:                      false,
		ControlledApplyEnabled:                false,
		ControlledApplyExecuted:               false,
		ControlledApplyReady:                  false,
		RuntimeMutationEnabled:                false,
		ActuatorRuntimeEnabled:                false,
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
		ExecutorBoundaryRecord: WindowsFluidVirtExecutorBoundaryRecord{
			ExecutorBoundaryRecordID:             "executor_boundary_record_v1",
			ExecutorBoundaryRecordVersion:        "v1",
			BoundaryState:                        "defined_disabled",
			SourceControlledApplyPlanBoundaryRef: "windows_fluidvirt_controlled_apply_plan_boundary_v1",
			SourceComplianceReplayRef:            "windows_fluidvirt_compliance_replay_audit_chain_v1",
			SourceReadonlyReplayRef:              "windows_fluidvirt_node_actuator_mvp_readonly_replay_v1",
			ExecutorBoundaryDefined:              true,
			ExecutorEnabled:                      false,
			ExecutorRuntimeAvailable:             false,
			ExecutorExecuted:                     false,
			RuntimeMutationAllowed:               false,
			RawControlExposureAllowed:            false,
			EvidenceConfidence:                   "low",
			NextRequiredState:                    "guarded_fake_runtime_mvp_or_inventory_witness_integration",
			ClaimBoundary:                        "executor_boundary_only",
		},
		ExecutorInputBoundary: WindowsFluidVirtExecutorInputBoundary{
			InputBoundaryState:                     "planning_only",
			AllowedInputTypes:                      allowedInput,
			ForbiddenInputTypes:                    forbiddenInput,
			CandidateActionInputAllowedForPlanning: true,
			RealRuntimeCommandInputAllowed:         false,
			RawCgroupPathInputAllowed:              false,
			RawQMPCommandInputAllowed:              false,
			RawQGACommandInputAllowed:              false,
			SecretMaterialInputAllowed:             false,
			RequiresComplianceReplay:               true,
			RequiresAuditHashChain:                 true,
			RequiresManualApproval:                 true,
			RequiresActiveLease:                    true,
			ClaimBoundary:                          "input_planning_only",
		},
		ExecutorOutputBoundary: WindowsFluidVirtExecutorOutputBoundary{
			OutputBoundaryState:               "planning_only",
			AllowedOutputTypes:                allowedOutput,
			ForbiddenOutputTypes:              forbiddenOutput,
			RuntimeMutationOutputAllowed:      false,
			RawRuntimeControlOutputAllowed:    false,
			SecretOutputAllowed:               false,
			AuditEventOutputRequired:          true,
			DeterministicReplayOutputRequired: true,
			ClaimBoundary:                     "output_planning_only",
		},
		ExecutorAuthorizationBoundary: WindowsFluidVirtExecutorAuthorizationBoundary{
			ManualApprovalRequired:          true,
			OperatorIdentityRequired:        true,
			ApprovalRecordRequired:          true,
			ApprovalBeforeExecutionRequired: true,
			AutonomousApprovalAllowed:       false,
			ApprovalBypassAllowed:           false,
			ClaimBoundary:                   "manual_approval_required",
		},
		ExecutorLeaseBoundary: WindowsFluidVirtExecutorLeaseBoundary{
			ActiveLeaseRequired:              true,
			LeaseTTLRequired:                 true,
			LeaseExpiryReturnToFloorRequired: true,
			LeaseWithoutAuditAllowed:         false,
			LeaseBypassAllowed:               false,
			ClaimBoundary:                    "lease_required",
		},
		ExecutorWitnessBoundary: WindowsFluidVirtExecutorWitnessBoundary{
			GuestWitnessRequired:       true,
			FluidShellOrQGARequired:    true,
			GuestAckRequired:           true,
			SameBootProofRequired:      true,
			SameQEMUProofRequired:      true,
			PendingRebootCheckRequired: true,
			WitnessIntegrated:          false,
			WitnessBypassAllowed:       false,
			ClaimBoundary:              "witness_required_not_integrated",
		},
		ExecutorRollbackBoundary: WindowsFluidVirtExecutorRollbackBoundary{
			RollbackPlanRequired:         true,
			RollbackBaselineRequired:     true,
			RollbackVerificationRequired: true,
			RollbackExecutionAllowed:     false,
			RollbackBypassAllowed:        false,
			ClaimBoundary:                "rollback_required_not_proven",
		},
		ExecutorReturnToFloorBoundary: WindowsFluidVirtExecutorReturnToFloorBoundary{
			ReturnToFloorPlanRequired:         true,
			ReturnToFloorVerificationRequired: true,
			ReturnToFloorExecutionAllowed:     false,
			ReturnToFloorBypassAllowed:        false,
			ClaimBoundary:                     "return_to_floor_required_not_proven",
		},
		ExecutorAuditBoundary: WindowsFluidVirtExecutorAuditBoundary{
			AuditHashChainRequired:         true,
			ComplianceReplayRequired:       true,
			AuditHashChainVerifiedRequired: true,
			EveryInputMustBeReplayable:     true,
			EveryOutputMustBeAudited:       true,
			AuditBypassAllowed:             false,
			ClaimBoundary:                  "audit_required",
		},
		ExecutorKillSwitchBoundary: WindowsFluidVirtExecutorKillSwitchBoundary{
			KillSwitchRequired:                true,
			KillSwitchBeforeRuntimeRequired:   true,
			KillSwitchVerified:                false,
			ExecutionWithoutKillSwitchAllowed: false,
			ClaimBoundary:                     "kill_switch_required_not_verified",
		},
		ExecutorRuntimeBoundary: WindowsFluidVirtExecutorRuntimeBoundary{
			RuntimeExecutionAllowed:    false,
			ExecutorEnabled:            false,
			ActuatorRuntimeEnabled:     false,
			ControlledApplyEnabled:     false,
			CgroupWriteEnabled:         false,
			QMPCommandExecutionAllowed: false,
			QGACommandExecutionAllowed: false,
			RawRuntimeControlsExposed:  false,
			ProductionMutationAllowed:  false,
			AutonomousApplyAllowed:     false,
			ClaimBoundary:              "runtime_execution_forbidden",
		},
		RequiredBeforeExecutorEnablement: []string{
			"explicit_executor_runtime_milestone",
			"controlled_apply_plan_ready",
			"manual_approval_flow",
			"operator_identity_audit",
			"active_lease_model",
			"lease_ttl",
			"guest_witness_integration",
			"same_boot_proof",
			"same_qemu_proof",
			"rollback_plan",
			"rollback_baseline",
			"return_to_floor_plan",
			"audit_hash_chain_verified",
			"compliance_replay_for_every_candidate",
			"kill_switch_verified",
			"node_allowlist",
			"cgroup_path_validation_strategy",
			"no_raw_runtime_control_surface",
			"no_autonomous_apply",
			"no_production_apply",
		},
		ExecutorBlockingReasons: blockingReasons,
		ExecutorReadinessScorecard: WindowsFluidVirtExecutorReadinessScorecard{
			ScorecardState:            "executor_boundary_defined_not_ready",
			ExecutorBoundaryDefined:   true,
			ExecutorEnabled:           false,
			ExecutorRuntimeAvailable:  false,
			ExecutorExecuted:          false,
			ControlledApplyEnabled:    false,
			ControlledApplyReady:      false,
			RuntimeMutationEnabled:    false,
			ActuatorRuntimeEnabled:    false,
			ComplianceReplayAvailable: true,
			AuditHashChainVerified:    true,
			BlockingReasonCount:       len(blockingReasons),
			Confidence:                "low",
			Decision:                  "define_executor_boundary_without_runtime_enablement",
			DecisionReason:            "boundary defined while runtime prerequisites remain unmet",
		},
		AllowedExecutorBoundaryActions: []string{
			"define_guarded_executor_boundary",
			"record_executor_input_boundary",
			"record_executor_output_boundary",
			"record_executor_authorization_boundary",
			"record_executor_lease_boundary",
			"record_executor_witness_boundary",
			"record_executor_rollback_boundary",
			"record_executor_return_to_floor_boundary",
			"record_executor_audit_boundary",
			"record_executor_kill_switch_boundary",
			"record_executor_runtime_boundary",
			"record_required_before_executor_enablement",
			"record_executor_blocking_reasons",
		},
		ForbiddenExecutorBoundaryActions: []string{
			"enable_executor",
			"execute_windows_executor",
			"execute_controlled_apply",
			"enable_controlled_apply",
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
			"executor boundary defined does not imply executor readiness",
			"controlled apply remains disabled",
			"not windows ga",
			"not windows production ready",
		},
		EvidenceRefs: []string{
			"plan://windows_fluidvirt_controlled_apply_plan_boundary_v1",
			"executor-boundary://windows_fluidvirt_guarded_executor_boundary_v1",
		},
		AuditTrail: []string{
			"audit://executor-boundary-defined",
			"audit://executor-blocking-reasons-recorded",
			"audit://executor-runtime-disabled-attested",
		},
	}
}

func ValidateWindowsFluidVirtGuardedExecutorBoundary(boundary WindowsFluidVirtGuardedExecutorBoundary) error {
	if boundary.BoundaryMode != "guarded_executor_boundary_only" || !boundary.ExecutorBoundaryDefined {
		return fmt.Errorf("boundary must remain guarded-executor-boundary-only and defined")
	}
	if boundary.ExecutorEnabled || boundary.ExecutorRuntimeAvailable || boundary.ExecutorExecuted {
		return fmt.Errorf("executor must remain disabled and unavailable")
	}
	if boundary.ControlledApplyEnabled || boundary.ControlledApplyExecuted || boundary.ControlledApplyReady {
		return fmt.Errorf("controlled apply must remain disabled and not ready")
	}
	if boundary.RuntimeMutationEnabled || boundary.ActuatorRuntimeEnabled || boundary.CgroupWriteEnabled {
		return fmt.Errorf("runtime mutation and cgroup write must remain disabled")
	}
	if boundary.QMPCommandExecutionAllowed || boundary.QGACommandExecutionAllowed || boundary.RawRuntimeControlsExposed {
		return fmt.Errorf("qmp/qga execution and raw controls must remain disabled")
	}
	if boundary.ProductionMutationAllowed || boundary.AutonomousApplyAllowed {
		return fmt.Errorf("production/autonomous mutation must remain disabled")
	}
	if boundary.WindowsGaClaimAllowed || boundary.WindowsProductionReadyClaimAllowed || boundary.WindowsExecutionReadyByDefault {
		return fmt.Errorf("windows ga/production/execution-ready claims must remain disabled")
	}
	if boundary.ExecutorReadinessScorecard.ScorecardState != "executor_boundary_defined_not_ready" {
		return fmt.Errorf("executor scorecard must remain not ready")
	}
	if len(boundary.RequiredBeforeExecutorEnablement) == 0 || len(boundary.ExecutorBlockingReasons) == 0 {
		return fmt.Errorf("required-before-executor and blocking reasons must be populated")
	}
	return nil
}

func LoadGuardedExecutorBoundaryFixtureFromTemporaryFile(path string) (WindowsFluidVirtGuardedExecutorBoundary, error) {
	clean := filepath.Clean(path)
	if strings.HasPrefix(clean, "/sys/fs/cgroup") {
		return WindowsFluidVirtGuardedExecutorBoundary{}, fmt.Errorf("real cgroup path forbidden")
	}
	if !strings.Contains(clean, os.TempDir()) {
		return WindowsFluidVirtGuardedExecutorBoundary{}, fmt.Errorf("fixture must be loaded from temporary path")
	}
	raw, err := os.ReadFile(clean)
	if err != nil {
		return WindowsFluidVirtGuardedExecutorBoundary{}, err
	}
	var boundary WindowsFluidVirtGuardedExecutorBoundary
	if err := json.Unmarshal(raw, &boundary); err != nil {
		return WindowsFluidVirtGuardedExecutorBoundary{}, err
	}
	if err := ValidateWindowsFluidVirtGuardedExecutorBoundary(boundary); err != nil {
		return WindowsFluidVirtGuardedExecutorBoundary{}, err
	}
	return boundary, nil
}
