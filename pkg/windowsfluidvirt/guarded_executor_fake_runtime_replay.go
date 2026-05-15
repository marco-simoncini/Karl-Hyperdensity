package windowsfluidvirt

import (
	"fmt"
	"strings"
)

type WindowsFluidVirtExecutorFakeRuntimeBoundary struct {
	FakeRuntimeOnly         bool   `json:"fakeRuntimeOnly"`
	DeterministicReplay     bool   `json:"deterministicReplay"`
	SafeForCI               bool   `json:"safeForCi"`
	UsesTemporaryFilesOnly  bool   `json:"usesTemporaryFilesOnly"`
	TouchesRealCgroup       bool   `json:"touchesRealCgroup"`
	TouchesRealQMP          bool   `json:"touchesRealQmp"`
	TouchesRealQGA          bool   `json:"touchesRealQga"`
	TouchesHostRuntime      bool   `json:"touchesHostRuntime"`
	RequiresNoPrivileges    bool   `json:"requiresNoPrivileges"`
	RejectsSysFsCgroupPath  bool   `json:"rejectsSysFsCgroupPath"`
	RejectsRawQMPInput      bool   `json:"rejectsRawQmpInput"`
	RejectsRawQGAInput      bool   `json:"rejectsRawQgaInput"`
	RejectsSecretMaterial   bool   `json:"rejectsSecretMaterial"`
	ClaimBoundary           string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorReplayInput struct {
	InputID                  string `json:"inputId"`
	InputVersion             string `json:"inputVersion"`
	SourceFixtureRef         string `json:"sourceFixtureRef"`
	SourceExecutorBoundaryRef string `json:"sourceExecutorBoundaryRef"`
	SourceControlledApplyPlanBoundaryRef string `json:"sourceControlledApplyPlanBoundaryRef"`
	Deterministic            bool   `json:"deterministic"`
	ContainsSecretMaterial   bool   `json:"containsSecretMaterial"`
	TouchesRuntime           bool   `json:"touchesRuntime"`
	TouchesRealCgroup        bool   `json:"touchesRealCgroup"`
	TouchesRealQMP           bool   `json:"touchesRealQmp"`
	TouchesRealQGA           bool   `json:"touchesRealQga"`
	CandidateActionInputType string `json:"candidateActionInputType"`
	InputBoundaryRespected   bool   `json:"inputBoundaryRespected"`
	ClaimBoundary            string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorReplayScenario struct {
	ScenarioID                string `json:"scenarioId"`
	ScenarioVersion           string `json:"scenarioVersion"`
	ScenarioState             string `json:"scenarioState"`
	ScenarioMode              string `json:"scenarioMode"`
	CandidateActionState      string `json:"candidateActionState"`
	ManualApprovalPresent     bool   `json:"manualApprovalPresent"`
	ActiveLeasePresent        bool   `json:"activeLeasePresent"`
	GuestWitnessPresent       bool   `json:"guestWitnessPresent"`
	RollbackPlanProven        bool   `json:"rollbackPlanProven"`
	ReturnToFloorPlanProven   bool   `json:"returnToFloorPlanProven"`
	KillSwitchVerified        bool   `json:"killSwitchVerified"`
	ComplianceReplayVerified  bool   `json:"complianceReplayVerified"`
	AuditHashChainVerified    bool   `json:"auditHashChainVerified"`
	ExecutorEnabled           bool   `json:"executorEnabled"`
	RuntimeMutationAllowed    bool   `json:"runtimeMutationAllowed"`
	ClaimBoundary             string `json:"claimBoundary"`
}

type WindowsFluidVirtReplayedExecutorStep struct {
	StepID             string `json:"stepId"`
	StepType           string `json:"stepType"`
	StepState          string `json:"stepState"`
	DeterministicOrder int    `json:"deterministicOrder"`
	TouchesRuntime     bool   `json:"touchesRuntime"`
	ClaimBoundary      string `json:"claimBoundary"`
}

type WindowsFluidVirtReplayedValidation struct {
	ValidationID      string `json:"validationId"`
	ValidationState   string `json:"validationState"`
	ClaimBoundary     string `json:"claimBoundary"`
}

type WindowsFluidVirtReplayedGateEvaluation struct {
	GateID              string `json:"gateId"`
	GateState           string `json:"gateState"`
	RequiredBeforeRuntime bool `json:"requiredBeforeRuntime"`
	SourceBoundaryRef   string `json:"sourceBoundaryRef"`
	ClaimBoundary       string `json:"claimBoundary"`
}

type WindowsFluidVirtGuardedExecutorReplayedAuditEvent struct {
	EventID               string `json:"eventId"`
	EventType             string `json:"eventType"`
	EventState            string `json:"eventState"`
	DeterministicOrder    int    `json:"deterministicOrder"`
	ContainsSecretMaterial bool  `json:"containsSecretMaterial"`
	TouchesRuntime        bool   `json:"touchesRuntime"`
	ClaimBoundary         string `json:"claimBoundary"`
}

type WindowsFluidVirtExecutorFakeRuntimeReplayValidation struct {
	ValidationState                string `json:"validationState"`
	DeterministicReplay            bool   `json:"deterministicReplay"`
	FakeRuntimeOnly                bool   `json:"fakeRuntimeOnly"`
	RealRuntimeTouched             bool   `json:"realRuntimeTouched"`
	RealCgroupTouched              bool   `json:"realCgroupTouched"`
	RealQMPTouched                 bool   `json:"realQmpTouched"`
	RealQGATouched                 bool   `json:"realQgaTouched"`
	ExecutorBoundaryRespected      bool   `json:"executorBoundaryRespected"`
	InputBoundaryRespected         bool   `json:"inputBoundaryRespected"`
	OutputBoundaryRespected        bool   `json:"outputBoundaryRespected"`
	ControlledApplyBoundaryRespected bool `json:"controlledApplyBoundaryRespected"`
	AuditHashChainVerified         bool   `json:"auditHashChainVerified"`
	ComplianceReplayVerified       bool   `json:"complianceReplayVerified"`
	ForbiddenClaimsRespected       bool   `json:"forbiddenClaimsRespected"`
	ExecutorEnabled                bool   `json:"executorEnabled"`
	ControlledApplyReady           bool   `json:"controlledApplyReady"`
	RuntimeMutationEnabled         bool   `json:"runtimeMutationEnabled"`
	NextRequiredState              string `json:"nextRequiredState"`
}

type WindowsFluidVirtGuardedExecutorFakeRuntimeReplay struct {
	ExecutorFakeRuntimeReplayID              string `json:"executorFakeRuntimeReplayId"`
	ExecutorFakeRuntimeReplayVersion         string `json:"executorFakeRuntimeReplayVersion"`
	ReleaseTrack                             string `json:"releaseTrack"`
	LaneStatus                               string `json:"laneStatus"`
	ReplayMode                               string `json:"replayMode"`
	SourceProductModelRef                    string `json:"sourceProductModelRef"`
	SourceNodeActuatorBoundaryRef            string `json:"sourceNodeActuatorBoundaryRef"`
	SourceNodeActuatorReadonlyReplayRef      string `json:"sourceNodeActuatorReadonlyReplayRef"`
	SourceComplianceReplayRef                string `json:"sourceComplianceReplayRef"`
	SourceControlledApplyPlanBoundaryRef     string `json:"sourceControlledApplyPlanBoundaryRef"`
	SourceGuardedExecutorBoundaryRef         string `json:"sourceGuardedExecutorBoundaryRef"`
	ExecutorFakeRuntimeReplayAvailable       bool   `json:"executorFakeRuntimeReplayAvailable"`
	ExecutorFakeRuntimeReplayExecuted        bool   `json:"executorFakeRuntimeReplayExecuted"`
	ExecutorFakeRuntimeReplayDeterministic   bool   `json:"executorFakeRuntimeReplayDeterministic"`
	ExecutorEnabled                          bool   `json:"executorEnabled"`
	ExecutorRuntimeAvailable                 bool   `json:"executorRuntimeAvailable"`
	ExecutorExecuted                         bool   `json:"executorExecuted"`
	ControlledApplyEnabled                   bool   `json:"controlledApplyEnabled"`
	ControlledApplyExecuted                  bool   `json:"controlledApplyExecuted"`
	ControlledApplyReady                     bool   `json:"controlledApplyReady"`
	RuntimeMutationEnabled                   bool   `json:"runtimeMutationEnabled"`
	ActuatorRuntimeEnabled                   bool   `json:"actuatorRuntimeEnabled"`
	CgroupWriteEnabled                       bool   `json:"cgroupWriteEnabled"`
	CPUMaxMutationEnabled                    bool   `json:"cpuMaxMutationEnabled"`
	QMPCommandExecutionAllowed               bool   `json:"qmpCommandExecutionAllowed"`
	QGACommandExecutionAllowed               bool   `json:"qgaCommandExecutionAllowed"`
	RawRuntimeControlsExposed                bool   `json:"rawRuntimeControlsExposed"`
	ProductionMutationAllowed                bool   `json:"productionMutationAllowed"`
	AutonomousApplyAllowed                   bool   `json:"autonomousApplyAllowed"`
	EnforcementMode                          string `json:"enforcementMode"`
	WindowsGaClaimAllowed                    bool   `json:"windowsGaClaimAllowed"`
	WindowsProductionReadyClaimAllowed       bool   `json:"windowsProductionReadyClaimAllowed"`
	WindowsExecutionReadyByDefault           bool   `json:"windowsExecutionReadyByDefault"`
	VCPUHotplugClaimAllowed                  bool   `json:"vcpuHotplugClaimAllowed"`
	LogicalCPUScalingClaimAllowed            bool   `json:"logicalCpuScalingClaimAllowed"`
	PoolScalingClaimAllowed                  bool   `json:"poolScalingClaimAllowed"`
	LiveMigrationClaimAllowed                bool   `json:"liveMigrationClaimAllowed"`
	RebootRecreateRolloutMechanismAllowed    bool   `json:"rebootRecreateRolloutMechanismAllowed"`
	FakeRuntimeBoundary                      WindowsFluidVirtExecutorFakeRuntimeBoundary `json:"fakeRuntimeBoundary"`
	ExecutorReplayInput                      WindowsFluidVirtExecutorReplayInput `json:"executorReplayInput"`
	ExecutorReplayScenario                   WindowsFluidVirtExecutorReplayScenario `json:"executorReplayScenario"`
	ReplayedExecutorSteps                    []WindowsFluidVirtReplayedExecutorStep `json:"replayedExecutorSteps"`
	ReplayedInputValidations                 []WindowsFluidVirtReplayedValidation `json:"replayedInputValidations"`
	ReplayedOutputValidations                []WindowsFluidVirtReplayedValidation `json:"replayedOutputValidations"`
	ReplayedGateEvaluations                  []WindowsFluidVirtReplayedGateEvaluation `json:"replayedGateEvaluations"`
	ReplayedBlockingReasons                  []string `json:"replayedBlockingReasons"`
	ReplayedAuditEvents                      []WindowsFluidVirtGuardedExecutorReplayedAuditEvent `json:"replayedAuditEvents"`
	ExecutorFakeRuntimeReplayValidation      WindowsFluidVirtExecutorFakeRuntimeReplayValidation `json:"executorFakeRuntimeReplayValidation"`
	RequiredBeforeExecutorRuntime            []string `json:"requiredBeforeExecutorRuntime"`
	AllowedExecutorFakeRuntimeReplayActions  []string `json:"allowedExecutorFakeRuntimeReplayActions"`
	ForbiddenExecutorFakeRuntimeReplayActions []string `json:"forbiddenExecutorFakeRuntimeReplayActions"`
	ClaimBoundaries                          []string `json:"claimBoundaries"`
	EvidenceRefs                             []string `json:"evidenceRefs"`
	AuditTrail                               []string `json:"auditTrail"`
}

func NewWindowsFluidVirtGuardedExecutorFakeRuntimeReplayMinimal() WindowsFluidVirtGuardedExecutorFakeRuntimeReplay {
	return WindowsFluidVirtGuardedExecutorFakeRuntimeReplay{
		ExecutorFakeRuntimeReplayID:            "windows_fluidvirt_guarded_executor_fake_runtime_replay_v1",
		ExecutorFakeRuntimeReplayVersion:       "v1",
		ReleaseTrack:                           "technical_preview",
		LaneStatus:                             "gated_preview",
		ReplayMode:                             "guarded_executor_fake_runtime_replay_only",
		SourceProductModelRef:                  "windows-fluidvirt-product-model-v1",
		SourceNodeActuatorBoundaryRef:          "windows_fluidvirt_node_actuator_contract_boundary_v1",
		SourceNodeActuatorReadonlyReplayRef:    "windows_fluidvirt_node_actuator_mvp_readonly_replay_v1",
		SourceComplianceReplayRef:              "windows_fluidvirt_compliance_replay_audit_chain_v1",
		SourceControlledApplyPlanBoundaryRef:   "windows_fluidvirt_controlled_apply_plan_boundary_v1",
		SourceGuardedExecutorBoundaryRef:       "windows_fluidvirt_guarded_executor_boundary_v1",
		ExecutorFakeRuntimeReplayAvailable:     true,
		ExecutorFakeRuntimeReplayExecuted:      true,
		ExecutorFakeRuntimeReplayDeterministic: true,
		ExecutorEnabled:                        false,
		ExecutorRuntimeAvailable:               false,
		ExecutorExecuted:                       false,
		ControlledApplyEnabled:                 false,
		ControlledApplyExecuted:                false,
		ControlledApplyReady:                   false,
		RuntimeMutationEnabled:                 false,
		ActuatorRuntimeEnabled:                 false,
		CgroupWriteEnabled:                     false,
		CPUMaxMutationEnabled:                  false,
		QMPCommandExecutionAllowed:             false,
		QGACommandExecutionAllowed:             false,
		RawRuntimeControlsExposed:              false,
		ProductionMutationAllowed:              false,
		AutonomousApplyAllowed:                 false,
		EnforcementMode:                        "disabled",
		WindowsGaClaimAllowed:                  false,
		WindowsProductionReadyClaimAllowed:     false,
		WindowsExecutionReadyByDefault:         false,
		VCPUHotplugClaimAllowed:                false,
		LogicalCPUScalingClaimAllowed:          false,
		PoolScalingClaimAllowed:                false,
		LiveMigrationClaimAllowed:              false,
		RebootRecreateRolloutMechanismAllowed:  false,
		FakeRuntimeBoundary: WindowsFluidVirtExecutorFakeRuntimeBoundary{
			FakeRuntimeOnly:        true,
			DeterministicReplay:    true,
			SafeForCI:              true,
			UsesTemporaryFilesOnly: true,
			TouchesRealCgroup:      false,
			TouchesRealQMP:         false,
			TouchesRealQGA:         false,
			TouchesHostRuntime:     false,
			RequiresNoPrivileges:   true,
			RejectsSysFsCgroupPath: true,
			RejectsRawQMPInput:     true,
			RejectsRawQGAInput:     true,
			RejectsSecretMaterial:  true,
			ClaimBoundary:          "fake_runtime_only",
		},
		ExecutorReplayInput: WindowsFluidVirtExecutorReplayInput{
			InputID:                  "guarded_executor_replay_input_v1",
			InputVersion:             "v1",
			SourceFixtureRef:         "examples/windows-fluid-product-fixtures/guarded_executor_fake_runtime_replay_minimal.json",
			SourceExecutorBoundaryRef: "windows_fluidvirt_guarded_executor_boundary_v1",
			SourceControlledApplyPlanBoundaryRef: "windows_fluidvirt_controlled_apply_plan_boundary_v1",
			Deterministic:            true,
			ContainsSecretMaterial:   false,
			TouchesRuntime:           false,
			TouchesRealCgroup:        false,
			TouchesRealQMP:           false,
			TouchesRealQGA:           false,
			CandidateActionInputType: "controlled_apply_candidate_model",
			InputBoundaryRespected:   true,
			ClaimBoundary:            "input_boundary_respected",
		},
		ExecutorReplayScenario: WindowsFluidVirtExecutorReplayScenario{
			ScenarioID:              "guarded_executor_fake_runtime_replay_scenario_v1",
			ScenarioVersion:         "v1",
			ScenarioState:           "replay_completed",
			ScenarioMode:            "fake_runtime_only",
			CandidateActionState:    "planning_only",
			ManualApprovalPresent:   false,
			ActiveLeasePresent:      false,
			GuestWitnessPresent:     false,
			RollbackPlanProven:      false,
			ReturnToFloorPlanProven: false,
			KillSwitchVerified:      false,
			ComplianceReplayVerified: true,
			AuditHashChainVerified:  true,
			ExecutorEnabled:         false,
			RuntimeMutationAllowed:  false,
			ClaimBoundary:           "scenario_replay_only",
		},
		ReplayedExecutorSteps: []WindowsFluidVirtReplayedExecutorStep{
			{StepID: "s1", StepType: "replay_started", StepState: "replayed", DeterministicOrder: 1, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s2", StepType: "executor_boundary_loaded", StepState: "replayed", DeterministicOrder: 2, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s3", StepType: "input_boundary_checked", StepState: "replayed", DeterministicOrder: 3, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s4", StepType: "output_boundary_checked", StepState: "replayed", DeterministicOrder: 4, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s5", StepType: "compliance_replay_verified", StepState: "replayed", DeterministicOrder: 5, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s6", StepType: "audit_hash_chain_verified", StepState: "replayed", DeterministicOrder: 6, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s7", StepType: "candidate_action_loaded", StepState: "replayed", DeterministicOrder: 7, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s8", StepType: "approval_gate_evaluated", StepState: "replayed", DeterministicOrder: 8, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s9", StepType: "lease_gate_evaluated", StepState: "replayed", DeterministicOrder: 9, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s10", StepType: "witness_gate_evaluated", StepState: "replayed", DeterministicOrder: 10, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s11", StepType: "rollback_gate_evaluated", StepState: "replayed", DeterministicOrder: 11, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s12", StepType: "return_to_floor_gate_evaluated", StepState: "replayed", DeterministicOrder: 12, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s13", StepType: "kill_switch_gate_evaluated", StepState: "replayed", DeterministicOrder: 13, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s14", StepType: "runtime_boundary_checked", StepState: "replayed", DeterministicOrder: 14, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s15", StepType: "blocking_reasons_emitted", StepState: "replayed", DeterministicOrder: 15, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s16", StepType: "no_runtime_mutation_attested", StepState: "replayed", DeterministicOrder: 16, TouchesRuntime: false, ClaimBoundary: "replay_only"},
			{StepID: "s17", StepType: "replay_completed", StepState: "replayed", DeterministicOrder: 17, TouchesRuntime: false, ClaimBoundary: "replay_only"},
		},
		ReplayedInputValidations: []WindowsFluidVirtReplayedValidation{
			{ValidationID: "input_type_allowed_for_planning", ValidationState: "passed", ClaimBoundary: "input_boundary"},
			{ValidationID: "raw_cgroup_path_rejected", ValidationState: "passed", ClaimBoundary: "input_boundary"},
			{ValidationID: "raw_qmp_command_rejected", ValidationState: "passed", ClaimBoundary: "input_boundary"},
			{ValidationID: "raw_qga_command_rejected", ValidationState: "passed", ClaimBoundary: "input_boundary"},
			{ValidationID: "raw_shell_command_rejected", ValidationState: "passed", ClaimBoundary: "input_boundary"},
			{ValidationID: "secret_material_rejected", ValidationState: "passed", ClaimBoundary: "input_boundary"},
			{ValidationID: "unapproved_candidate_blocked", ValidationState: "passed", ClaimBoundary: "input_boundary"},
			{ValidationID: "unaudited_candidate_blocked", ValidationState: "passed", ClaimBoundary: "input_boundary"},
			{ValidationID: "production_auto_candidate_blocked", ValidationState: "passed", ClaimBoundary: "input_boundary"},
		},
		ReplayedOutputValidations: []WindowsFluidVirtReplayedValidation{
			{ValidationID: "boundary_validation_output_allowed", ValidationState: "passed", ClaimBoundary: "output_boundary"},
			{ValidationID: "blocking_reason_output_allowed", ValidationState: "passed", ClaimBoundary: "output_boundary"},
			{ValidationID: "audit_event_output_allowed", ValidationState: "passed", ClaimBoundary: "output_boundary"},
			{ValidationID: "no_runtime_apply_output", ValidationState: "passed", ClaimBoundary: "output_boundary"},
			{ValidationID: "no_cgroup_write_output", ValidationState: "passed", ClaimBoundary: "output_boundary"},
			{ValidationID: "no_qmp_command_output", ValidationState: "passed", ClaimBoundary: "output_boundary"},
			{ValidationID: "no_qga_command_output", ValidationState: "passed", ClaimBoundary: "output_boundary"},
			{ValidationID: "no_secret_output", ValidationState: "passed", ClaimBoundary: "output_boundary"},
			{ValidationID: "no_production_auto_output", ValidationState: "passed", ClaimBoundary: "output_boundary"},
		},
		ReplayedGateEvaluations: []WindowsFluidVirtReplayedGateEvaluation{
			{GateID: "manual_approval_gate_blocked", GateState: "blocked", RequiredBeforeRuntime: true, SourceBoundaryRef: "executorAuthorizationBoundary", ClaimBoundary: "gates_replayed"},
			{GateID: "active_lease_gate_blocked", GateState: "blocked", RequiredBeforeRuntime: true, SourceBoundaryRef: "executorLeaseBoundary", ClaimBoundary: "gates_replayed"},
			{GateID: "guest_witness_gate_blocked", GateState: "blocked", RequiredBeforeRuntime: true, SourceBoundaryRef: "executorWitnessBoundary", ClaimBoundary: "gates_replayed"},
			{GateID: "rollback_plan_gate_blocked", GateState: "blocked", RequiredBeforeRuntime: true, SourceBoundaryRef: "executorRollbackBoundary", ClaimBoundary: "gates_replayed"},
			{GateID: "return_to_floor_gate_blocked", GateState: "blocked", RequiredBeforeRuntime: true, SourceBoundaryRef: "executorReturnToFloorBoundary", ClaimBoundary: "gates_replayed"},
			{GateID: "kill_switch_gate_blocked", GateState: "blocked", RequiredBeforeRuntime: true, SourceBoundaryRef: "executorKillSwitchBoundary", ClaimBoundary: "gates_replayed"},
			{GateID: "node_allowlist_gate_blocked", GateState: "blocked", RequiredBeforeRuntime: true, SourceBoundaryRef: "sourceNodeActuatorBoundaryRef", ClaimBoundary: "gates_replayed"},
			{GateID: "cgroup_path_validation_gate_blocked", GateState: "blocked", RequiredBeforeRuntime: true, SourceBoundaryRef: "sourceNodeActuatorBoundaryRef", ClaimBoundary: "gates_replayed"},
			{GateID: "runtime_boundary_gate_passed_no_mutation", GateState: "passed", RequiredBeforeRuntime: true, SourceBoundaryRef: "executorRuntimeBoundary", ClaimBoundary: "gates_replayed"},
			{GateID: "compliance_replay_gate_passed", GateState: "passed", RequiredBeforeRuntime: true, SourceBoundaryRef: "sourceComplianceReplayRef", ClaimBoundary: "gates_replayed"},
			{GateID: "audit_hash_chain_gate_passed", GateState: "passed", RequiredBeforeRuntime: true, SourceBoundaryRef: "sourceComplianceReplayRef", ClaimBoundary: "gates_replayed"},
		},
		ReplayedBlockingReasons: []string{
			"executor_fake_runtime_replay_only",
			"executor_runtime_not_enabled",
			"controlled_apply_not_ready",
			"manual_approval_missing",
			"active_lease_missing",
			"guest_witness_missing",
			"rollback_plan_not_proven",
			"return_to_floor_plan_not_proven",
			"kill_switch_not_verified",
			"node_allowlist_missing",
			"cgroup_path_validation_missing",
			"cgroup_write_disabled",
			"qmp_execution_disabled",
			"qga_execution_disabled",
			"production_ready_claim_forbidden",
			"autonomous_apply_forbidden",
			"raw_runtime_control_forbidden",
		},
		ReplayedAuditEvents: []WindowsFluidVirtGuardedExecutorReplayedAuditEvent{
			{EventID: "a1", EventType: "fake_executor_replay_started", EventState: "replayed", DeterministicOrder: 1, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "audit_replay"},
			{EventID: "a2", EventType: "executor_boundary_loaded", EventState: "replayed", DeterministicOrder: 2, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "audit_replay"},
			{EventID: "a3", EventType: "input_validation_replayed", EventState: "replayed", DeterministicOrder: 3, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "audit_replay"},
			{EventID: "a4", EventType: "output_validation_replayed", EventState: "replayed", DeterministicOrder: 4, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "audit_replay"},
			{EventID: "a5", EventType: "gates_replayed", EventState: "replayed", DeterministicOrder: 5, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "audit_replay"},
			{EventID: "a6", EventType: "blocking_reasons_recorded", EventState: "replayed", DeterministicOrder: 6, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "audit_replay"},
			{EventID: "a7", EventType: "no_runtime_mutation_attested", EventState: "replayed", DeterministicOrder: 7, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "audit_replay"},
			{EventID: "a8", EventType: "fake_executor_replay_completed", EventState: "replayed", DeterministicOrder: 8, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "audit_replay"},
		},
		ExecutorFakeRuntimeReplayValidation: WindowsFluidVirtExecutorFakeRuntimeReplayValidation{
			ValidationState:                  "passed",
			DeterministicReplay:              true,
			FakeRuntimeOnly:                  true,
			RealRuntimeTouched:               false,
			RealCgroupTouched:                false,
			RealQMPTouched:                   false,
			RealQGATouched:                   false,
			ExecutorBoundaryRespected:        true,
			InputBoundaryRespected:           true,
			OutputBoundaryRespected:          true,
			ControlledApplyBoundaryRespected: true,
			AuditHashChainVerified:           true,
			ComplianceReplayVerified:         true,
			ForbiddenClaimsRespected:         true,
			ExecutorEnabled:                  false,
			ControlledApplyReady:             false,
			RuntimeMutationEnabled:           false,
			NextRequiredState:                "inventory_fluidshell_witness_or_guarded_runtime_mvp_boundary",
		},
		RequiredBeforeExecutorRuntime: []string{
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
		AllowedExecutorFakeRuntimeReplayActions: []string{
			"run_guarded_executor_fake_runtime_replay",
			"validate_executor_input_boundary",
			"validate_executor_output_boundary",
			"replay_executor_gate_evaluations",
			"emit_executor_blocking_reasons",
			"validate_no_runtime_mutation",
			"record_required_before_executor_runtime",
		},
		ForbiddenExecutorFakeRuntimeReplayActions: []string{
			"enable_executor",
			"execute_windows_executor_runtime",
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
			"fake-runtime replay does not imply executor runtime availability",
			"controlled apply remains disabled",
			"not windows ga",
			"not windows production ready",
		},
		EvidenceRefs: []string{
			"executor-boundary://windows_fluidvirt_guarded_executor_boundary_v1",
			"executor-replay://windows_fluidvirt_guarded_executor_fake_runtime_replay_v1",
		},
		AuditTrail: []string{
			"audit://guarded-executor-fake-runtime-replay-started",
			"audit://guarded-executor-gates-replayed",
			"audit://no-runtime-mutation-attested",
		},
	}
}

func ValidateWindowsFluidVirtGuardedExecutorFakeRuntimeReplay(replay WindowsFluidVirtGuardedExecutorFakeRuntimeReplay) error {
	if replay.ReplayMode != "guarded_executor_fake_runtime_replay_only" || !replay.ExecutorFakeRuntimeReplayDeterministic {
		return fmt.Errorf("replay must remain guarded fake-runtime deterministic mode")
	}
	if !replay.ExecutorFakeRuntimeReplayAvailable || !replay.ExecutorFakeRuntimeReplayExecuted {
		return fmt.Errorf("replay availability/execution flags must be true for replay model")
	}
	if replay.ExecutorEnabled || replay.ExecutorRuntimeAvailable || replay.ExecutorExecuted {
		return fmt.Errorf("executor runtime must remain disabled and unavailable")
	}
	if replay.ControlledApplyEnabled || replay.ControlledApplyExecuted || replay.ControlledApplyReady {
		return fmt.Errorf("controlled apply must remain disabled and not ready")
	}
	if replay.RuntimeMutationEnabled || replay.ActuatorRuntimeEnabled || replay.CgroupWriteEnabled || replay.CPUMaxMutationEnabled {
		return fmt.Errorf("runtime mutation/cgroup write/cpu mutation must remain disabled")
	}
	if replay.QMPCommandExecutionAllowed || replay.QGACommandExecutionAllowed || replay.RawRuntimeControlsExposed {
		return fmt.Errorf("qmp/qga execution and raw runtime controls must remain disabled")
	}
	if replay.ProductionMutationAllowed || replay.AutonomousApplyAllowed {
		return fmt.Errorf("production/autonomous mutation must remain disabled")
	}
	if replay.WindowsGaClaimAllowed || replay.WindowsProductionReadyClaimAllowed || replay.WindowsExecutionReadyByDefault {
		return fmt.Errorf("windows ga/production/execution-ready claims must remain disabled")
	}
	if !replay.ExecutorFakeRuntimeReplayValidation.FakeRuntimeOnly || replay.ExecutorFakeRuntimeReplayValidation.RealRuntimeTouched || replay.ExecutorFakeRuntimeReplayValidation.RealCgroupTouched || replay.ExecutorFakeRuntimeReplayValidation.RealQMPTouched || replay.ExecutorFakeRuntimeReplayValidation.RealQGATouched {
		return fmt.Errorf("validation result must attest fake-runtime only and no real runtime touch")
	}
	if !replay.ExecutorReplayInput.InputBoundaryRespected || replay.ExecutorReplayInput.ContainsSecretMaterial {
		return fmt.Errorf("input boundary must be respected and secret-free")
	}
	if !replay.FakeRuntimeBoundary.FakeRuntimeOnly || replay.FakeRuntimeBoundary.TouchesRealCgroup || replay.FakeRuntimeBoundary.TouchesRealQMP || replay.FakeRuntimeBoundary.TouchesRealQGA {
		return fmt.Errorf("fake runtime boundary must not touch real runtime surfaces")
	}
	for _, step := range replay.ReplayedExecutorSteps {
		if step.TouchesRuntime || step.StepState != "replayed" {
			return fmt.Errorf("replayed step %s violates replay-only constraints", step.StepID)
		}
	}
	for _, event := range replay.ReplayedAuditEvents {
		if event.ContainsSecretMaterial || event.TouchesRuntime || event.EventState != "replayed" {
			return fmt.Errorf("replayed audit event %s violates replay safety constraints", event.EventID)
		}
	}
	return nil
}

func ReplayGuardedExecutorFakeRuntime(replay WindowsFluidVirtGuardedExecutorFakeRuntimeReplay) (WindowsFluidVirtGuardedExecutorFakeRuntimeReplay, error) {
	if err := ValidateWindowsFluidVirtGuardedExecutorFakeRuntimeReplay(replay); err != nil {
		return WindowsFluidVirtGuardedExecutorFakeRuntimeReplay{}, err
	}
	return replay, nil
}

func hasForbiddenMaterial(raw string) bool {
	needle := strings.ToLower(raw)
	return strings.Contains(needle, "raw_qmp_payload") ||
		strings.Contains(needle, "raw_qga_payload") ||
		strings.Contains(needle, "\"token\":") ||
		strings.Contains(needle, "bearer ") ||
		strings.Contains(needle, "-----begin") ||
		strings.Contains(needle, "kubeconfig") ||
		strings.Contains(needle, "\"password\":")
}
