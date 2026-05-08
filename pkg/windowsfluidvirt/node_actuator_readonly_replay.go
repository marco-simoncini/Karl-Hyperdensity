package windowsfluidvirt

type WindowsFluidVirtReplayCheckState string

const (
	ReplayCheckPassed  WindowsFluidVirtReplayCheckState = "passed"
	ReplayCheckBlocked WindowsFluidVirtReplayCheckState = "blocked"
	ReplayCheckReplayed WindowsFluidVirtReplayCheckState = "replayed"
)

type WindowsFluidVirtReplayTransitionState string

const (
	TransitionReplayedModelOnly WindowsFluidVirtReplayTransitionState = "replayed_model_only"
)

type WindowsFluidVirtReplayEventState string

const (
	ReplayEventReplayed WindowsFluidVirtReplayEventState = "replayed"
)

type WindowsFluidVirtFakeRuntimeBoundary struct {
	FakeRuntimeOnly        bool   `json:"fakeRuntimeOnly"`
	UsesTemporaryFilesOnly bool   `json:"usesTemporaryFilesOnly"`
	TouchesRealCgroup      bool   `json:"touchesRealCgroup"`
	TouchesRealQMP         bool   `json:"touchesRealQmp"`
	TouchesRealQGA         bool   `json:"touchesRealQga"`
	TouchesHostRuntime     bool   `json:"touchesHostRuntime"`
	RequiresNoPrivileges   bool   `json:"requiresNoPrivileges"`
	DeterministicReplay    bool   `json:"deterministicReplay"`
	SafeForCI              bool   `json:"safeForCi"`
	ClaimBoundary          string `json:"claimBoundary"`
}

type WindowsFluidVirtReplayScenario struct {
	ScenarioID              string `json:"scenarioId"`
	ScenarioVersion         string `json:"scenarioVersion"`
	ScenarioState           string `json:"scenarioState"`
	SourceFixtureRef        string `json:"sourceFixtureRef"`
	CPUMechanism            string `json:"cpuMechanism"`
	RAMMechanism            string `json:"ramMechanism"`
	SameBootRequired        bool   `json:"sameBootRequired"`
	SameQEMURequired        bool   `json:"sameQemuRequired"`
	GuestWitnessRequired    bool   `json:"guestWitnessRequired"`
	ReturnToFloorRequired   bool   `json:"returnToFloorRequired"`
	RollbackRequired        bool   `json:"rollbackRequired"`
	AuditRequired           bool   `json:"auditRequired"`
	ManualApprovalRequired  bool   `json:"manualApprovalRequired"`
	LeaseTTLRequired        bool   `json:"leaseTtlRequired"`
	RuntimeMutationAllowed  bool   `json:"runtimeMutationAllowed"`
	ClaimBoundary           string `json:"claimBoundary"`
}

type WindowsFluidVirtReplayedCPUEntitlementTransition struct {
	TransitionID            string                                `json:"transitionId"`
	TransitionState         WindowsFluidVirtReplayTransitionState `json:"transitionState"`
	FromCPUMax              string                                `json:"fromCpuMax"`
	ToCPUMax                string                                `json:"toCpuMax"`
	AppliedToRealCgroup     bool                                  `json:"appliedToRealCgroup"`
	FakeRuntimeOnly         bool                                  `json:"fakeRuntimeOnly"`
	RequiresApproval        bool                                  `json:"requiresApproval"`
	RequiresRollback        bool                                  `json:"requiresRollback"`
	RequiresReturnToFloor   bool                                  `json:"requiresReturnToFloor"`
	ValidationState         string                                `json:"validationState"`
	ClaimBoundary           string                                `json:"claimBoundary"`
}

type WindowsFluidVirtReplayedSafetyCheck struct {
	CheckID                 string                              `json:"checkId"`
	CheckState              WindowsFluidVirtReplayCheckState   `json:"checkState"`
	SourceContractGate      WindowsFluidVirtNodeActuatorGate   `json:"sourceContractGate"`
	RequiredBeforeRuntimeMVP bool                              `json:"requiredBeforeRuntimeMvp"`
	ClaimBoundary           string                              `json:"claimBoundary"`
}

type WindowsFluidVirtReplayedAuditEvent struct {
	EventID                 string                           `json:"eventId"`
	EventType               string                           `json:"eventType"`
	EventState              WindowsFluidVirtReplayEventState `json:"eventState"`
	TimestampPolicy         string                           `json:"timestampPolicy"`
	ContainsSecretMaterial  bool                             `json:"containsSecretMaterial"`
	ClaimBoundary           string                           `json:"claimBoundary"`
}

type WindowsFluidVirtReplayValidationResult struct {
	ValidationState         string `json:"validationState"`
	FakeRuntimeOnly         bool   `json:"fakeRuntimeOnly"`
	RealRuntimeTouched      bool   `json:"realRuntimeTouched"`
	RealCgroupTouched       bool   `json:"realCgroupTouched"`
	RealQMPTouched          bool   `json:"realQmpTouched"`
	RealQGATouched          bool   `json:"realQgaTouched"`
	ContractBoundaryRespected bool `json:"contractBoundaryRespected"`
	ForbiddenClaimsRespected bool  `json:"forbiddenClaimsRespected"`
	ReplayDeterministic     bool   `json:"replayDeterministic"`
	RuntimeMVPReady         bool   `json:"runtimeMvpReady"`
	NextRequiredState       string `json:"nextRequiredState"`
}

type WindowsFluidVirtNodeActuatorReadonlyReplay struct {
	ReplayID                             string                                  `json:"replayId"`
	ReplayVersion                        string                                  `json:"replayVersion"`
	ReleaseTrack                         string                                  `json:"releaseTrack"`
	LaneStatus                           string                                  `json:"laneStatus"`
	ReplayMode                           string                                  `json:"replayMode"`
	SourceContractRef                    string                                  `json:"sourceContractRef"`
	NodeActuatorMvpReplayAvailable       bool                                    `json:"nodeActuatorMvpReplayAvailable"`
	NodeActuatorMvpPortedAsRuntime       bool                                    `json:"nodeActuatorMvpPortedAsRuntime"`
	ActuatorRuntimeEnabled               bool                                    `json:"actuatorRuntimeEnabled"`
	RuntimeMutationEnabled               bool                                    `json:"runtimeMutationEnabled"`
	CgroupWriteEnabled                   bool                                    `json:"cgroupWriteEnabled"`
	CPUMaxMutationEnabled                bool                                    `json:"cpuMaxMutationEnabled"`
	QMPCommandExecutionAllowed           bool                                    `json:"qmpCommandExecutionAllowed"`
	QGACommandExecutionAllowed           bool                                    `json:"qgaCommandExecutionAllowed"`
	RawRuntimeControlsExposed            bool                                    `json:"rawRuntimeControlsExposed"`
	ProductionMutationAllowed            bool                                    `json:"productionMutationAllowed"`
	AutonomousApplyAllowed               bool                                    `json:"autonomousApplyAllowed"`
	EnforcementMode                      string                                  `json:"enforcementMode"`
	WindowsGaClaimAllowed                bool                                    `json:"windowsGaClaimAllowed"`
	WindowsProductionReadyClaimAllowed   bool                                    `json:"windowsProductionReadyClaimAllowed"`
	WindowsExecutionReadyByDefault       bool                                    `json:"windowsExecutionReadyByDefault"`
	VCPUHotplugClaimAllowed              bool                                    `json:"vcpuHotplugClaimAllowed"`
	LogicalCPUScalingClaimAllowed        bool                                    `json:"logicalCpuScalingClaimAllowed"`
	PoolScalingClaimAllowed              bool                                    `json:"poolScalingClaimAllowed"`
	LiveMigrationClaimAllowed            bool                                    `json:"liveMigrationClaimAllowed"`
	RebootRecreateRolloutMechanismAllowed bool                                   `json:"rebootRecreateRolloutMechanismAllowed"`
	ReplayScenario                       WindowsFluidVirtReplayScenario          `json:"replayScenario"`
	FakeRuntimeBoundary                  WindowsFluidVirtFakeRuntimeBoundary     `json:"fakeRuntimeBoundary"`
	ReplayedCPUEntitlementTransitions    []WindowsFluidVirtReplayedCPUEntitlementTransition `json:"replayedCpuEntitlementTransitions"`
	ReplayedSafetyChecks                 []WindowsFluidVirtReplayedSafetyCheck   `json:"replayedSafetyChecks"`
	ReplayedBlockers                     []WindowsFluidVirtNodeActuatorBlocker   `json:"replayedBlockers"`
	ReplayedAuditEvents                  []WindowsFluidVirtReplayedAuditEvent    `json:"replayedAuditEvents"`
	ReplayValidationResult               WindowsFluidVirtReplayValidationResult  `json:"replayValidationResult"`
	RequiredBeforeRuntimeMVP             []string                                `json:"requiredBeforeRuntimeMvp"`
	AllowedReplayActions                 []string                                `json:"allowedReplayActions"`
	ForbiddenReplayActions               []string                                `json:"forbiddenReplayActions"`
	ClaimBoundaries                      []string                                `json:"claimBoundaries"`
	EvidenceRefs                         []string                                `json:"evidenceRefs"`
	AuditTrail                           []string                                `json:"auditTrail"`
}

func NewWindowsFluidVirtNodeActuatorReadonlyReplayMinimal() WindowsFluidVirtNodeActuatorReadonlyReplay {
	return WindowsFluidVirtNodeActuatorReadonlyReplay{
		ReplayID:                         "windows_fluidvirt_node_actuator_mvp_readonly_replay_v1",
		ReplayVersion:                    "v1",
		ReleaseTrack:                     "technical_preview",
		LaneStatus:                       "gated_preview",
		ReplayMode:                       "readonly_replay_only",
		SourceContractRef:                "windows_fluidvirt_node_actuator_contract_boundary_v1",
		NodeActuatorMvpReplayAvailable:   true,
		NodeActuatorMvpPortedAsRuntime:   false,
		ActuatorRuntimeEnabled:           false,
		RuntimeMutationEnabled:           false,
		CgroupWriteEnabled:               false,
		CPUMaxMutationEnabled:            false,
		QMPCommandExecutionAllowed:       false,
		QGACommandExecutionAllowed:       false,
		RawRuntimeControlsExposed:        false,
		ProductionMutationAllowed:        false,
		AutonomousApplyAllowed:           false,
		EnforcementMode:                  "disabled",
		WindowsGaClaimAllowed:            false,
		WindowsProductionReadyClaimAllowed: false,
		WindowsExecutionReadyByDefault:   false,
		VCPUHotplugClaimAllowed:          false,
		LogicalCPUScalingClaimAllowed:    false,
		PoolScalingClaimAllowed:          false,
		LiveMigrationClaimAllowed:        false,
		RebootRecreateRolloutMechanismAllowed: false,
		ReplayScenario: WindowsFluidVirtReplayScenario{
			ScenarioID:             "node_actuator_replay_minimal",
			ScenarioVersion:        "v1",
			ScenarioState:          "replay_completed",
			SourceFixtureRef:       "examples/windows-fluid-product-fixtures/node_actuator_readonly_replay_minimal.json",
			CPUMechanism:           "cgroup_v2_cpu_max_entitlement_liquidity",
			RAMMechanism:           "qmp_balloon_liquidity_model_reference_only",
			SameBootRequired:       true,
			SameQEMURequired:       true,
			GuestWitnessRequired:   true,
			ReturnToFloorRequired:  true,
			RollbackRequired:       true,
			AuditRequired:          true,
			ManualApprovalRequired: true,
			LeaseTTLRequired:       true,
			RuntimeMutationAllowed: false,
			ClaimBoundary:          "technical_preview_boundary_active",
		},
		FakeRuntimeBoundary: DefaultWindowsFluidVirtFakeRuntimeBoundary(),
		ReplayedCPUEntitlementTransitions: []WindowsFluidVirtReplayedCPUEntitlementTransition{
			NewModelOnlyTransition("t1", "floor_to_ceiling_candidate", "100000 100000", "200000 100000", true),
			NewModelOnlyTransition("t2", "ceiling_to_floor_return", "200000 100000", "100000 100000", true),
			NewModelOnlyTransition("t3", "blocked_without_approval", "100000 100000", "200000 100000", false),
			NewModelOnlyTransition("t4", "blocked_without_guest_witness", "100000 100000", "200000 100000", false),
			NewModelOnlyTransition("t5", "blocked_without_return_to_floor", "200000 100000", "100000 100000", false),
			NewModelOnlyTransition("t6", "blocked_without_audit_chain", "100000 100000", "200000 100000", false),
		},
		ReplayedSafetyChecks: BuildDefaultReplayedSafetyChecks(),
		ReplayedBlockers: []WindowsFluidVirtNodeActuatorBlocker{
			NodeActuatorRuntimeNotPorted,
			NodeActuatorCgroupWritePathNotEnabled,
			NodeActuatorNodeAllowlistMissing,
			NodeActuatorManualApprovalMissing,
			NodeActuatorLeaseTTLMissing,
			NodeActuatorReturnToFloorPlanMissing,
			NodeActuatorRollbackPlanMissing,
			NodeActuatorGuestWitnessMissing,
			NodeActuatorSameBootProofMissing,
			NodeActuatorSameQEMUProofMissing,
			NodeActuatorAuditChainMissing,
			NodeActuatorKillSwitchMissing,
			NodeActuatorProductionReadyForbidden,
			NodeActuatorAutonomousApplyForbidden,
			NodeActuatorRawRuntimeControlForbidden,
		},
		ReplayedAuditEvents: []WindowsFluidVirtReplayedAuditEvent{
			{EventID: "e1", EventType: "replay_started", EventState: ReplayEventReplayed, TimestampPolicy: "deterministic_or_not_available", ContainsSecretMaterial: false, ClaimBoundary: "technical_preview_boundary_active"},
			{EventID: "e2", EventType: "boundary_loaded", EventState: ReplayEventReplayed, TimestampPolicy: "deterministic_or_not_available", ContainsSecretMaterial: false, ClaimBoundary: "technical_preview_boundary_active"},
			{EventID: "e3", EventType: "fake_runtime_initialized", EventState: ReplayEventReplayed, TimestampPolicy: "deterministic_or_not_available", ContainsSecretMaterial: false, ClaimBoundary: "technical_preview_boundary_active"},
			{EventID: "e4", EventType: "candidate_transition_replayed", EventState: ReplayEventReplayed, TimestampPolicy: "deterministic_or_not_available", ContainsSecretMaterial: false, ClaimBoundary: "technical_preview_boundary_active"},
			{EventID: "e5", EventType: "blocker_evaluated", EventState: ReplayEventReplayed, TimestampPolicy: "deterministic_or_not_available", ContainsSecretMaterial: false, ClaimBoundary: "technical_preview_boundary_active"},
			{EventID: "e6", EventType: "return_to_floor_replayed", EventState: ReplayEventReplayed, TimestampPolicy: "deterministic_or_not_available", ContainsSecretMaterial: false, ClaimBoundary: "technical_preview_boundary_active"},
			{EventID: "e7", EventType: "replay_completed", EventState: ReplayEventReplayed, TimestampPolicy: "deterministic_or_not_available", ContainsSecretMaterial: false, ClaimBoundary: "technical_preview_boundary_active"},
			{EventID: "e8", EventType: "no_runtime_mutation_attested", EventState: ReplayEventReplayed, TimestampPolicy: "deterministic_or_not_available", ContainsSecretMaterial: false, ClaimBoundary: "technical_preview_boundary_active"},
		},
		ReplayValidationResult: WindowsFluidVirtReplayValidationResult{
			ValidationState:          "passed",
			FakeRuntimeOnly:          true,
			RealRuntimeTouched:       false,
			RealCgroupTouched:        false,
			RealQMPTouched:           false,
			RealQGATouched:           false,
			ContractBoundaryRespected: true,
			ForbiddenClaimsRespected:  true,
			ReplayDeterministic:      true,
			RuntimeMVPReady:          false,
			NextRequiredState:        "node_actuator_mvp_port_with_guarded_fake_runtime_or_compliance_replay",
		},
		RequiredBeforeRuntimeMVP: []string{
			"explicit_runtime_mvp_milestone",
			"node_allowlist_design",
			"cgroup_path_validation_strategy",
			"manual_approval_flow",
			"guest_witness_integration",
			"same_boot_same_qemu_evidence",
			"rollback_plan",
			"return_to_floor_plan",
			"audit_hash_chain",
			"kill_switch",
			"lease_ttl",
			"compliance_replay",
			"no_raw_runtime_control_surface",
		},
		AllowedReplayActions: []string{
			"run_readonly_replay",
			"run_fake_runtime_simulation",
			"validate_contract_boundary",
			"validate_replayed_blockers",
			"validate_replayed_audit_events",
			"validate_no_runtime_mutation",
			"record_runtime_mvp_requirements",
		},
		ForbiddenReplayActions: []string{
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
			"readonly replay does not imply runtime actuator enablement",
			"no production apply",
			"no autonomous apply",
			"not windows ga",
			"not windows production ready",
		},
		EvidenceRefs: []string{
			"contract://windows_fluidvirt_node_actuator_contract_boundary_v1",
			"replay://node_actuator_readonly_replay",
		},
		AuditTrail: []string{
			"audit://replay-started",
			"audit://replay-completed",
			"audit://no-runtime-mutation-attested",
		},
	}
}

func BuildDefaultReplayedSafetyChecks() []WindowsFluidVirtReplayedSafetyCheck {
	return []WindowsFluidVirtReplayedSafetyCheck{
		{CheckID: "contract_boundary_loaded", CheckState: ReplayCheckPassed, SourceContractGate: GateNodeAllowlistDefined, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "node_allowlist_required", CheckState: ReplayCheckReplayed, SourceContractGate: GateNodeAllowlistDefined, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "manual_approval_required", CheckState: ReplayCheckReplayed, SourceContractGate: GateManualApprovalRequired, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "active_lease_required", CheckState: ReplayCheckReplayed, SourceContractGate: GateActiveLeaseRequired, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "lease_ttl_required", CheckState: ReplayCheckReplayed, SourceContractGate: GateLeaseTTLRequired, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "return_to_floor_plan_required", CheckState: ReplayCheckReplayed, SourceContractGate: GateReturnToFloorPlanRequired, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "rollback_plan_required", CheckState: ReplayCheckReplayed, SourceContractGate: GateRollbackPlanRequired, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "guest_witness_required", CheckState: ReplayCheckReplayed, SourceContractGate: GateGuestWitnessRequired, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "same_boot_proof_required", CheckState: ReplayCheckReplayed, SourceContractGate: GateSameBootProofRequired, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "same_qemu_proof_required", CheckState: ReplayCheckReplayed, SourceContractGate: GateSameQEMUProofRequired, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "audit_hash_chain_required", CheckState: ReplayCheckReplayed, SourceContractGate: GateAuditHashChainRequired, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "kill_switch_required", CheckState: ReplayCheckReplayed, SourceContractGate: GateKillSwitchRequired, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "no_raw_control_exposure", CheckState: ReplayCheckPassed, SourceContractGate: GateNoRawControlExposure, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "no_autonomous_apply", CheckState: ReplayCheckPassed, SourceContractGate: GateNoAutonomousApply, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "no_production_apply", CheckState: ReplayCheckPassed, SourceContractGate: GateNoProductionApply, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
		{CheckID: "no_windows_ga_claim", CheckState: ReplayCheckPassed, SourceContractGate: GateNoWindowsGAClaim, RequiredBeforeRuntimeMVP: true, ClaimBoundary: "technical_preview_boundary_active"},
	}
}

func NewModelOnlyTransition(id, validationState, fromCPUMax, toCPUMax string, requiresApproval bool) WindowsFluidVirtReplayedCPUEntitlementTransition {
	return WindowsFluidVirtReplayedCPUEntitlementTransition{
		TransitionID:          id,
		TransitionState:       TransitionReplayedModelOnly,
		FromCPUMax:            fromCPUMax,
		ToCPUMax:              toCPUMax,
		AppliedToRealCgroup:   false,
		FakeRuntimeOnly:       true,
		RequiresApproval:      requiresApproval,
		RequiresRollback:      true,
		RequiresReturnToFloor: true,
		ValidationState:       validationState,
		ClaimBoundary:         "technical_preview_boundary_active",
	}
}
