package windowsfluidvirt

type WindowsFluidVirtNodeActuatorGate string

const (
	GateNodeAllowlistDefined           WindowsFluidVirtNodeActuatorGate = "node_allowlist_defined"
	GateCgroupV2CPUMaxPathValidated    WindowsFluidVirtNodeActuatorGate = "cgroup_v2_cpu_max_path_validated"
	GateHostScopeNodeLocalConfirmed    WindowsFluidVirtNodeActuatorGate = "host_scope_node_local_confirmed"
	GateManualApprovalRequired         WindowsFluidVirtNodeActuatorGate = "manual_approval_required"
	GateActiveLeaseRequired            WindowsFluidVirtNodeActuatorGate = "active_lease_required"
	GateLeaseTTLRequired               WindowsFluidVirtNodeActuatorGate = "lease_ttl_required"
	GateReturnToFloorPlanRequired      WindowsFluidVirtNodeActuatorGate = "return_to_floor_plan_required"
	GateRollbackPlanRequired           WindowsFluidVirtNodeActuatorGate = "rollback_plan_required"
	GateGuestWitnessRequired           WindowsFluidVirtNodeActuatorGate = "guest_witness_required"
	GateSameBootProofRequired          WindowsFluidVirtNodeActuatorGate = "same_boot_proof_required"
	GateSameQEMUProofRequired          WindowsFluidVirtNodeActuatorGate = "same_qemu_proof_required"
	GateAuditHashChainRequired         WindowsFluidVirtNodeActuatorGate = "audit_hash_chain_required"
	GateKillSwitchRequired             WindowsFluidVirtNodeActuatorGate = "kill_switch_required"
	GateNoRawControlExposure           WindowsFluidVirtNodeActuatorGate = "no_raw_control_exposure"
	GateNoAutonomousApply              WindowsFluidVirtNodeActuatorGate = "no_autonomous_apply"
	GateNoProductionApply              WindowsFluidVirtNodeActuatorGate = "no_production_apply"
	GateNoWindowsGAClaim               WindowsFluidVirtNodeActuatorGate = "no_windows_ga_claim"
)

type WindowsFluidVirtNodeActuatorBlocker string

const (
	NodeActuatorRuntimeNotPorted          WindowsFluidVirtNodeActuatorBlocker = "actuator_runtime_not_ported"
	NodeActuatorCgroupWritePathNotEnabled WindowsFluidVirtNodeActuatorBlocker = "cgroup_write_path_not_enabled"
	NodeActuatorNodeAllowlistMissing      WindowsFluidVirtNodeActuatorBlocker = "node_allowlist_missing"
	NodeActuatorManualApprovalMissing     WindowsFluidVirtNodeActuatorBlocker = "manual_approval_missing"
	NodeActuatorLeaseTTLMissing           WindowsFluidVirtNodeActuatorBlocker = "lease_ttl_missing"
	NodeActuatorReturnToFloorPlanMissing  WindowsFluidVirtNodeActuatorBlocker = "return_to_floor_plan_missing"
	NodeActuatorRollbackPlanMissing       WindowsFluidVirtNodeActuatorBlocker = "rollback_plan_missing"
	NodeActuatorGuestWitnessMissing       WindowsFluidVirtNodeActuatorBlocker = "guest_witness_missing"
	NodeActuatorSameBootProofMissing      WindowsFluidVirtNodeActuatorBlocker = "same_boot_proof_missing"
	NodeActuatorSameQEMUProofMissing      WindowsFluidVirtNodeActuatorBlocker = "same_qemu_proof_missing"
	NodeActuatorAuditChainMissing         WindowsFluidVirtNodeActuatorBlocker = "audit_chain_missing"
	NodeActuatorKillSwitchMissing         WindowsFluidVirtNodeActuatorBlocker = "kill_switch_missing"
	NodeActuatorProductionReadyForbidden  WindowsFluidVirtNodeActuatorBlocker = "production_ready_claim_forbidden"
	NodeActuatorAutonomousApplyForbidden  WindowsFluidVirtNodeActuatorBlocker = "autonomous_apply_forbidden"
	NodeActuatorRawRuntimeControlForbidden WindowsFluidVirtNodeActuatorBlocker = "raw_runtime_control_forbidden"
)

type WindowsFluidVirtNodeActuatorAllowedAction string

const (
	AllowedDefineBoundary         WindowsFluidVirtNodeActuatorAllowedAction = "define_node_actuator_contract_boundary"
	AllowedValidateSafetyDefaults WindowsFluidVirtNodeActuatorAllowedAction = "validate_node_actuator_safety_defaults"
	AllowedRecordRequiredGates    WindowsFluidVirtNodeActuatorAllowedAction = "record_required_gates"
	AllowedRecordBlockers         WindowsFluidVirtNodeActuatorAllowedAction = "record_blockers"
	AllowedRecordClaimBoundary    WindowsFluidVirtNodeActuatorAllowedAction = "record_claim_boundary"
	AllowedRecordFutureMVPPlan    WindowsFluidVirtNodeActuatorAllowedAction = "record_future_mvp_porting_plan"
)

type WindowsFluidVirtNodeActuatorForbiddenAction string

const (
	ForbiddenExecuteNodeActuator     WindowsFluidVirtNodeActuatorForbiddenAction = "execute_node_actuator"
	ForbiddenWriteCgroupCPUMax       WindowsFluidVirtNodeActuatorForbiddenAction = "write_cgroup_cpu_max"
	ForbiddenMutateQMPBalloon        WindowsFluidVirtNodeActuatorForbiddenAction = "mutate_qmp_balloon"
	ForbiddenExecuteQMPCommand       WindowsFluidVirtNodeActuatorForbiddenAction = "execute_qmp_command"
	ForbiddenExecuteQGACommand       WindowsFluidVirtNodeActuatorForbiddenAction = "execute_qga_command"
	ForbiddenEnableAutonomousApply   WindowsFluidVirtNodeActuatorForbiddenAction = "enable_autonomous_apply"
	ForbiddenEnableProductionAUTO    WindowsFluidVirtNodeActuatorForbiddenAction = "enable_production_auto"
	ForbiddenMarkExecutionReady      WindowsFluidVirtNodeActuatorForbiddenAction = "mark_windows_execution_ready"
	ForbiddenClaimWindowsGA          WindowsFluidVirtNodeActuatorForbiddenAction = "claim_windows_ga"
	ForbiddenClaimWindowsProdReady   WindowsFluidVirtNodeActuatorForbiddenAction = "claim_windows_production_ready"
	ForbiddenExposeRawRuntimeControl WindowsFluidVirtNodeActuatorForbiddenAction = "expose_raw_runtime_controls"
	ForbiddenCreateSecret            WindowsFluidVirtNodeActuatorForbiddenAction = "create_secret"
	ForbiddenCommitSecret            WindowsFluidVirtNodeActuatorForbiddenAction = "commit_secret"
	ForbiddenLogToken                WindowsFluidVirtNodeActuatorForbiddenAction = "log_token"
	ForbiddenPackageOSISO            WindowsFluidVirtNodeActuatorForbiddenAction = "package_os_iso"
	ForbiddenTouchDashboard          WindowsFluidVirtNodeActuatorForbiddenAction = "touch_dashboard"
	ForbiddenTouchInventory          WindowsFluidVirtNodeActuatorForbiddenAction = "touch_inventory"
)

type WindowsFluidVirtNodeActuatorCPUBoundary struct {
	Mechanism                   string `json:"mechanism"`
	MechanismState              string `json:"mechanismState"`
	HostScope                   string `json:"hostScope"`
	WritePathAllowed            bool   `json:"writePathAllowed"`
	CgroupWriteEnabled          bool   `json:"cgroupWriteEnabled"`
	CPUMaxMutationEnabled       bool   `json:"cpuMaxMutationEnabled"`
	RequiresNodeAllowlist       bool   `json:"requiresNodeAllowlist"`
	RequiresManualApproval      bool   `json:"requiresManualApproval"`
	RequiresGuestWitness        bool   `json:"requiresGuestWitness"`
	RequiresSameBootProof       bool   `json:"requiresSameBootProof"`
	RequiresSameQEMUProof       bool   `json:"requiresSameQemuProof"`
	RequiresReturnToFloorPlan   bool   `json:"requiresReturnToFloorPlan"`
	RequiresRollbackPlan        bool   `json:"requiresRollbackPlan"`
	RequiresAuditChain          bool   `json:"requiresAuditChain"`
	RequiresKillSwitch          bool   `json:"requiresKillSwitch"`
	RequiresLeaseTTL            bool   `json:"requiresLeaseTtl"`
}

type WindowsFluidVirtNodeActuatorRAMBoundary struct {
	Mechanism                                    string `json:"mechanism"`
	MechanismState                               string `json:"mechanismState"`
	QMPCommandExecutionAllowed                   bool   `json:"qmpCommandExecutionAllowed"`
	RawQMPControlExposed                         bool   `json:"rawQmpControlExposed"`
	MemoryApplyEnabled                           bool   `json:"memoryApplyEnabled"`
	GenericKubevirtVmRamTemplateMutationAllowed  bool   `json:"genericKubevirtVmRamTemplateMutationAllowed"`
}

type WindowsFluidVirtNodeActuatorContract struct {
	ActuatorContractID                        string                                   `json:"actuatorContractId"`
	ActuatorContractVersion                   string                                   `json:"actuatorContractVersion"`
	ReleaseTrack                              string                                   `json:"releaseTrack"`
	LaneStatus                                string                                   `json:"laneStatus"`
	ContractMode                              string                                   `json:"contractMode"`
	NodeActuatorMvpPorted                     bool                                     `json:"nodeActuatorMvpPorted"`
	ActuatorRuntimeEnabled                    bool                                     `json:"actuatorRuntimeEnabled"`
	RuntimeMutationEnabled                    bool                                     `json:"runtimeMutationEnabled"`
	ProductionMutationAllowed                 bool                                     `json:"productionMutationAllowed"`
	AutonomousApplyAllowed                    bool                                     `json:"autonomousApplyAllowed"`
	EnforcementMode                           string                                   `json:"enforcementMode"`
	RawRuntimeControlsExposed                 bool                                     `json:"rawRuntimeControlsExposed"`
	WindowsGaClaimAllowed                     bool                                     `json:"windowsGaClaimAllowed"`
	WindowsProductionReadyClaimAllowed        bool                                     `json:"windowsProductionReadyClaimAllowed"`
	WindowsExecutionReadyByDefault            bool                                     `json:"windowsExecutionReadyByDefault"`
	VCPUHotplugClaimAllowed                   bool                                     `json:"vcpuHotplugClaimAllowed"`
	LogicalCPUScalingClaimAllowed             bool                                     `json:"logicalCpuScalingClaimAllowed"`
	PoolScalingClaimAllowed                   bool                                     `json:"poolScalingClaimAllowed"`
	LiveMigrationClaimAllowed                 bool                                     `json:"liveMigrationClaimAllowed"`
	RebootRecreateRolloutMechanismAllowed     bool                                     `json:"rebootRecreateRolloutMechanismAllowed"`
	CPUBoundary                               WindowsFluidVirtNodeActuatorCPUBoundary  `json:"cpuBoundary"`
	RAMBoundary                               WindowsFluidVirtNodeActuatorRAMBoundary  `json:"ramBoundary"`
	RequiredGates                             []WindowsFluidVirtNodeActuatorGate       `json:"requiredGates"`
	Blockers                                  []WindowsFluidVirtNodeActuatorBlocker    `json:"blockers"`
	AllowedActions                            []WindowsFluidVirtNodeActuatorAllowedAction `json:"allowedActions"`
	ForbiddenActions                          []WindowsFluidVirtNodeActuatorForbiddenAction `json:"forbiddenActions"`
}

func NewWindowsFluidVirtNodeActuatorContractBoundary() WindowsFluidVirtNodeActuatorContract {
	return WindowsFluidVirtNodeActuatorContract{
		ActuatorContractID:                    "windows_fluidvirt_node_actuator_contract_boundary_v1",
		ActuatorContractVersion:               "v1",
		ReleaseTrack:                          "technical_preview",
		LaneStatus:                            "gated_preview",
		ContractMode:                          "boundary_only",
		NodeActuatorMvpPorted:                 false,
		ActuatorRuntimeEnabled:                false,
		RuntimeMutationEnabled:                false,
		ProductionMutationAllowed:             false,
		AutonomousApplyAllowed:                false,
		EnforcementMode:                       "disabled",
		RawRuntimeControlsExposed:             false,
		WindowsGaClaimAllowed:                 false,
		WindowsProductionReadyClaimAllowed:    false,
		WindowsExecutionReadyByDefault:        false,
		VCPUHotplugClaimAllowed:               false,
		LogicalCPUScalingClaimAllowed:         false,
		PoolScalingClaimAllowed:               false,
		LiveMigrationClaimAllowed:             false,
		RebootRecreateRolloutMechanismAllowed: false,
		CPUBoundary: WindowsFluidVirtNodeActuatorCPUBoundary{
			Mechanism:                 "cgroup_v2_cpu_max_entitlement_liquidity",
			MechanismState:            "contract_defined_not_runtime_enabled",
			HostScope:                 "node_local",
			WritePathAllowed:          false,
			CgroupWriteEnabled:        false,
			CPUMaxMutationEnabled:     false,
			RequiresNodeAllowlist:     true,
			RequiresManualApproval:    true,
			RequiresGuestWitness:      true,
			RequiresSameBootProof:     true,
			RequiresSameQEMUProof:     true,
			RequiresReturnToFloorPlan: true,
			RequiresRollbackPlan:      true,
			RequiresAuditChain:        true,
			RequiresKillSwitch:        true,
			RequiresLeaseTTL:          true,
		},
		RAMBoundary: WindowsFluidVirtNodeActuatorRAMBoundary{
			Mechanism:                                   "qmp_balloon_liquidity_model",
			MechanismState:                              "contract_reference_only",
			QMPCommandExecutionAllowed:                  false,
			RawQMPControlExposed:                        false,
			MemoryApplyEnabled:                          false,
			GenericKubevirtVmRamTemplateMutationAllowed: false,
		},
		RequiredGates: []WindowsFluidVirtNodeActuatorGate{
			GateNodeAllowlistDefined,
			GateCgroupV2CPUMaxPathValidated,
			GateHostScopeNodeLocalConfirmed,
			GateManualApprovalRequired,
			GateActiveLeaseRequired,
			GateLeaseTTLRequired,
			GateReturnToFloorPlanRequired,
			GateRollbackPlanRequired,
			GateGuestWitnessRequired,
			GateSameBootProofRequired,
			GateSameQEMUProofRequired,
			GateAuditHashChainRequired,
			GateKillSwitchRequired,
			GateNoRawControlExposure,
			GateNoAutonomousApply,
			GateNoProductionApply,
			GateNoWindowsGAClaim,
		},
		Blockers: []WindowsFluidVirtNodeActuatorBlocker{
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
		AllowedActions: []WindowsFluidVirtNodeActuatorAllowedAction{
			AllowedDefineBoundary,
			AllowedValidateSafetyDefaults,
			AllowedRecordRequiredGates,
			AllowedRecordBlockers,
			AllowedRecordClaimBoundary,
			AllowedRecordFutureMVPPlan,
		},
		ForbiddenActions: []WindowsFluidVirtNodeActuatorForbiddenAction{
			ForbiddenExecuteNodeActuator,
			ForbiddenWriteCgroupCPUMax,
			ForbiddenMutateQMPBalloon,
			ForbiddenExecuteQMPCommand,
			ForbiddenExecuteQGACommand,
			ForbiddenEnableAutonomousApply,
			ForbiddenEnableProductionAUTO,
			ForbiddenMarkExecutionReady,
			ForbiddenClaimWindowsGA,
			ForbiddenClaimWindowsProdReady,
			ForbiddenExposeRawRuntimeControl,
			ForbiddenCreateSecret,
			ForbiddenCommitSecret,
			ForbiddenLogToken,
			ForbiddenPackageOSISO,
			ForbiddenTouchDashboard,
			ForbiddenTouchInventory,
		},
	}
}
