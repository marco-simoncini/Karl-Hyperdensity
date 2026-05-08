package windowsfluidvirt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type WindowsFluidVirtComplianceReplayInput struct {
	InputID                string `json:"inputId"`
	InputVersion           string `json:"inputVersion"`
	SourceFixtureRef       string `json:"sourceFixtureRef"`
	SourceReplayRef        string `json:"sourceReplayRef"`
	SourceContractRef      string `json:"sourceContractRef"`
	Deterministic          bool   `json:"deterministic"`
	ContainsSecretMaterial bool   `json:"containsSecretMaterial"`
	TouchesRuntime         bool   `json:"touchesRuntime"`
	TouchesRealCgroup      bool   `json:"touchesRealCgroup"`
	TouchesRealQMP         bool   `json:"touchesRealQmp"`
	TouchesRealQGA         bool   `json:"touchesRealQga"`
	ClaimBoundary          string `json:"claimBoundary"`
}

type WindowsFluidVirtComplianceReplayedEvent struct {
	EventID                string `json:"eventId"`
	EventType              string `json:"eventType"`
	EventState             string `json:"eventState"`
	DeterministicOrder     int    `json:"deterministicOrder"`
	PreviousHash           string `json:"previousHash"`
	EventHash              string `json:"eventHash"`
	ContainsSecretMaterial bool   `json:"containsSecretMaterial"`
	TouchesRuntime         bool   `json:"touchesRuntime"`
	ClaimBoundary          string `json:"claimBoundary"`
}

type WindowsFluidVirtComplianceCheck struct {
	CheckID       string `json:"checkId"`
	CheckState    string `json:"checkState"`
	Required      bool   `json:"required"`
	FailureReason string `json:"failureReason,omitempty"`
	SourceRef     string `json:"sourceRef"`
	ClaimBoundary string `json:"claimBoundary"`
}

type WindowsFluidVirtComplianceReplayValidationResult struct {
	ValidationState                string `json:"validationState"`
	DeterministicReplay            bool   `json:"deterministicReplay"`
	AuditHashChainVerified         bool   `json:"auditHashChainVerified"`
	ContractBoundaryRespected      bool   `json:"contractBoundaryRespected"`
	ReadonlyReplayBoundaryRespected bool  `json:"readonlyReplayBoundaryRespected"`
	ForbiddenClaimsRespected       bool   `json:"forbiddenClaimsRespected"`
	RealRuntimeTouched             bool   `json:"realRuntimeTouched"`
	RealCgroupTouched              bool   `json:"realCgroupTouched"`
	RealQMPTouched                 bool   `json:"realQmpTouched"`
	RealQGATouched                 bool   `json:"realQgaTouched"`
	RuntimeMutationEnabled         bool   `json:"runtimeMutationEnabled"`
	ActuatorRuntimeEnabled         bool   `json:"actuatorRuntimeEnabled"`
	ControlledApplyReady           bool   `json:"controlledApplyReady"`
	RuntimeMVPReady                bool   `json:"runtimeMvpReady"`
	NextRequiredState              string `json:"nextRequiredState"`
}

type WindowsFluidVirtComplianceReplay struct {
	ComplianceReplayID                string                                      `json:"complianceReplayId"`
	ComplianceReplayVersion           string                                      `json:"complianceReplayVersion"`
	ReleaseTrack                      string                                      `json:"releaseTrack"`
	LaneStatus                        string                                      `json:"laneStatus"`
	ReplayMode                        string                                      `json:"replayMode"`
	SourceProductModelRef             string                                      `json:"sourceProductModelRef"`
	SourceNodeActuatorBoundaryRef     string                                      `json:"sourceNodeActuatorBoundaryRef"`
	SourceNodeActuatorReadonlyReplayRef string                                    `json:"sourceNodeActuatorReadonlyReplayRef"`
	ComplianceReplayAvailable         bool                                        `json:"complianceReplayAvailable"`
	ComplianceReplayExecuted          bool                                        `json:"complianceReplayExecuted"`
	ComplianceReplayDeterministic     bool                                        `json:"complianceReplayDeterministic"`
	AuditHashChainAvailable           bool                                        `json:"auditHashChainAvailable"`
	AuditHashChainVerified            bool                                        `json:"auditHashChainVerified"`
	RuntimeMutationEnabled            bool                                        `json:"runtimeMutationEnabled"`
	ActuatorRuntimeEnabled            bool                                        `json:"actuatorRuntimeEnabled"`
	CgroupWriteEnabled                bool                                        `json:"cgroupWriteEnabled"`
	QMPCommandExecutionAllowed        bool                                        `json:"qmpCommandExecutionAllowed"`
	QGACommandExecutionAllowed        bool                                        `json:"qgaCommandExecutionAllowed"`
	RawRuntimeControlsExposed         bool                                        `json:"rawRuntimeControlsExposed"`
	ProductionMutationAllowed         bool                                        `json:"productionMutationAllowed"`
	AutonomousApplyAllowed            bool                                        `json:"autonomousApplyAllowed"`
	EnforcementMode                   string                                      `json:"enforcementMode"`
	WindowsGaClaimAllowed             bool                                        `json:"windowsGaClaimAllowed"`
	WindowsProductionReadyClaimAllowed bool                                       `json:"windowsProductionReadyClaimAllowed"`
	WindowsExecutionReadyByDefault    bool                                        `json:"windowsExecutionReadyByDefault"`
	ReplayInput                       WindowsFluidVirtComplianceReplayInput       `json:"replayInput"`
	ReplayedEvents                    []WindowsFluidVirtComplianceReplayedEvent   `json:"replayedEvents"`
	ComplianceChecks                  []WindowsFluidVirtComplianceCheck           `json:"complianceChecks"`
	AuditHashChain                    WindowsFluidVirtAuditHashChain              `json:"auditHashChain"`
	ReplayValidationResult            WindowsFluidVirtComplianceReplayValidationResult `json:"replayValidationResult"`
	RequiredBeforeControlledApply     []string                                    `json:"requiredBeforeControlledApply"`
	AllowedComplianceReplayActions    []string                                    `json:"allowedComplianceReplayActions"`
	ForbiddenComplianceReplayActions  []string                                    `json:"forbiddenComplianceReplayActions"`
	ClaimBoundaries                   []string                                    `json:"claimBoundaries"`
	EvidenceRefs                      []string                                    `json:"evidenceRefs"`
	AuditTrail                        []string                                    `json:"auditTrail"`
}

func NewWindowsFluidVirtComplianceReplayMinimal() WindowsFluidVirtComplianceReplay {
	events := []WindowsFluidVirtComplianceReplayedEvent{
		{EventID: "e1", EventType: "compliance_replay_started", EventState: "replayed", DeterministicOrder: 1, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e2", EventType: "product_model_loaded", EventState: "replayed", DeterministicOrder: 2, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e3", EventType: "actuator_boundary_loaded", EventState: "replayed", DeterministicOrder: 3, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e4", EventType: "readonly_replay_loaded", EventState: "replayed", DeterministicOrder: 4, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e5", EventType: "fake_runtime_boundary_verified", EventState: "replayed", DeterministicOrder: 5, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e6", EventType: "blocker_taxonomy_verified", EventState: "replayed", DeterministicOrder: 6, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e7", EventType: "candidate_transition_verified", EventState: "replayed", DeterministicOrder: 7, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e8", EventType: "return_to_floor_requirement_verified", EventState: "replayed", DeterministicOrder: 8, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e9", EventType: "rollback_requirement_verified", EventState: "replayed", DeterministicOrder: 9, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e10", EventType: "audit_hash_chain_built", EventState: "replayed", DeterministicOrder: 10, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e11", EventType: "audit_hash_chain_verified", EventState: "replayed", DeterministicOrder: 11, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e12", EventType: "compliance_replay_completed", EventState: "replayed", DeterministicOrder: 12, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
		{EventID: "e13", EventType: "no_runtime_mutation_attested", EventState: "replayed", DeterministicOrder: 13, ContainsSecretMaterial: false, TouchesRuntime: false, ClaimBoundary: "technical_preview_boundary_active"},
	}
	chain, replayedWithHashes, _ := BuildWindowsFluidVirtAuditHashChain(events)
	return WindowsFluidVirtComplianceReplay{
		ComplianceReplayID:                 "windows_fluidvirt_compliance_replay_audit_chain_v1",
		ComplianceReplayVersion:            "v1",
		ReleaseTrack:                       "technical_preview",
		LaneStatus:                         "gated_preview",
		ReplayMode:                         "deterministic_compliance_replay_only",
		SourceProductModelRef:              "windows-fluidvirt-product-model-v1",
		SourceNodeActuatorBoundaryRef:      "windows_fluidvirt_node_actuator_contract_boundary_v1",
		SourceNodeActuatorReadonlyReplayRef: "windows_fluidvirt_node_actuator_mvp_readonly_replay_v1",
		ComplianceReplayAvailable:          true,
		ComplianceReplayExecuted:           true,
		ComplianceReplayDeterministic:      true,
		AuditHashChainAvailable:            true,
		AuditHashChainVerified:             true,
		RuntimeMutationEnabled:             false,
		ActuatorRuntimeEnabled:             false,
		CgroupWriteEnabled:                 false,
		QMPCommandExecutionAllowed:         false,
		QGACommandExecutionAllowed:         false,
		RawRuntimeControlsExposed:          false,
		ProductionMutationAllowed:          false,
		AutonomousApplyAllowed:             false,
		EnforcementMode:                    "disabled",
		WindowsGaClaimAllowed:              false,
		WindowsProductionReadyClaimAllowed: false,
		WindowsExecutionReadyByDefault:     false,
		ReplayInput: WindowsFluidVirtComplianceReplayInput{
			InputID:                "compliance_replay_input_minimal",
			InputVersion:           "v1",
			SourceFixtureRef:       "examples/windows-fluid-product-fixtures/compliance_replay_minimal.json",
			SourceReplayRef:        "windows_fluidvirt_node_actuator_mvp_readonly_replay_v1",
			SourceContractRef:      "windows_fluidvirt_node_actuator_contract_boundary_v1",
			Deterministic:          true,
			ContainsSecretMaterial: false,
			TouchesRuntime:         false,
			TouchesRealCgroup:      false,
			TouchesRealQMP:         false,
			TouchesRealQGA:         false,
			ClaimBoundary:          "technical_preview_boundary_active",
		},
		ReplayedEvents: replayedWithHashes,
		ComplianceChecks: []WindowsFluidVirtComplianceCheck{
			{CheckID: "product_model_claim_boundary_valid", CheckState: "passed", Required: true, SourceRef: "windows-fluidvirt-product-model-v1", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "actuator_contract_boundary_valid", CheckState: "passed", Required: true, SourceRef: "windows_fluidvirt_node_actuator_contract_boundary_v1", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "readonly_replay_boundary_valid", CheckState: "passed", Required: true, SourceRef: "windows_fluidvirt_node_actuator_mvp_readonly_replay_v1", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "fake_runtime_boundary_valid", CheckState: "passed", Required: true, SourceRef: "fake_runtime_boundary", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "cgroup_write_disabled", CheckState: "passed", Required: true, SourceRef: "node_actuator_boundary", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "qmp_execution_disabled", CheckState: "passed", Required: true, SourceRef: "node_actuator_boundary", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "qga_execution_disabled", CheckState: "passed", Required: true, SourceRef: "node_actuator_boundary", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "autonomous_apply_disabled", CheckState: "passed", Required: true, SourceRef: "safety_defaults", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "production_mutation_disabled", CheckState: "passed", Required: true, SourceRef: "safety_defaults", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "windows_ga_claim_disabled", CheckState: "passed", Required: true, SourceRef: "claim_boundary", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "windows_production_ready_claim_disabled", CheckState: "passed", Required: true, SourceRef: "claim_boundary", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "windows_execution_ready_default_disabled", CheckState: "passed", Required: true, SourceRef: "claim_boundary", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "vcpu_hotplug_claim_disabled", CheckState: "passed", Required: true, SourceRef: "claim_boundary", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "logical_cpu_scaling_claim_disabled", CheckState: "passed", Required: true, SourceRef: "claim_boundary", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "pool_scaling_claim_disabled", CheckState: "passed", Required: true, SourceRef: "claim_boundary", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "raw_runtime_controls_disabled", CheckState: "passed", Required: true, SourceRef: "claim_boundary", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "blocker_taxonomy_present", CheckState: "passed", Required: true, SourceRef: "node_actuator_boundary_blockers", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "audit_hash_chain_verified", CheckState: "passed", Required: true, SourceRef: "audit_hash_chain", ClaimBoundary: "technical_preview_boundary_active"},
			{CheckID: "no_secret_material_present", CheckState: "passed", Required: true, SourceRef: "replayed_events", ClaimBoundary: "technical_preview_boundary_active"},
		},
		AuditHashChain: chain,
		ReplayValidationResult: WindowsFluidVirtComplianceReplayValidationResult{
			ValidationState:                 "passed",
			DeterministicReplay:             true,
			AuditHashChainVerified:          true,
			ContractBoundaryRespected:       true,
			ReadonlyReplayBoundaryRespected: true,
			ForbiddenClaimsRespected:        true,
			RealRuntimeTouched:              false,
			RealCgroupTouched:               false,
			RealQMPTouched:                  false,
			RealQGATouched:                  false,
			RuntimeMutationEnabled:          false,
			ActuatorRuntimeEnabled:          false,
			ControlledApplyReady:            false,
			RuntimeMVPReady:                 false,
			NextRequiredState:               "controlled_apply_plan_boundary_or_guarded_fake_runtime_mvp",
		},
		RequiredBeforeControlledApply: []string{
			"controlled_apply_plan_boundary",
			"operator_manual_approval_flow",
			"active_lease_model",
			"lease_ttl",
			"guest_witness_integration",
			"same_boot_same_qemu_evidence",
			"rollback_plan",
			"return_to_floor_plan",
			"audit_hash_chain_required",
			"kill_switch",
			"node_allowlist",
			"no_raw_runtime_control_surface",
			"compliance_replay_required_for_every_candidate",
		},
		AllowedComplianceReplayActions: []string{
			"run_deterministic_compliance_replay",
			"validate_product_model_boundary",
			"validate_actuator_boundary",
			"validate_readonly_replay_boundary",
			"build_audit_hash_chain",
			"verify_audit_hash_chain",
			"validate_no_runtime_mutation",
			"record_required_before_controlled_apply",
		},
		ForbiddenComplianceReplayActions: []string{
			"execute_node_actuator_runtime",
			"write_cgroup_cpu_max",
			"mutate_qmp_balloon",
			"execute_qmp_command",
			"execute_qga_command",
			"touch_real_cgroup",
			"touch_real_qmp",
			"touch_real_qga",
			"enable_runtime_actuator",
			"enable_controlled_apply",
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
			"deterministic compliance replay is model-only",
			"compliance replay execution does not imply runtime readiness",
			"not windows ga",
			"not windows production ready",
		},
		EvidenceRefs: []string{
			"contract://windows_fluidvirt_node_actuator_contract_boundary_v1",
			"replay://windows_fluidvirt_node_actuator_mvp_readonly_replay_v1",
			"compliance://windows_fluidvirt_compliance_replay_audit_chain_v1",
		},
		AuditTrail: []string{
			"audit://compliance-replay-started",
			"audit://audit-hash-chain-verified",
			"audit://compliance-replay-completed",
		},
	}
}

func ValidateWindowsFluidVirtComplianceReplay(replay WindowsFluidVirtComplianceReplay) error {
	if replay.ReplayMode != "deterministic_compliance_replay_only" || !replay.ComplianceReplayDeterministic {
		return fmt.Errorf("compliance replay must remain deterministic replay-only")
	}
	if replay.ActuatorRuntimeEnabled || replay.RuntimeMutationEnabled || replay.CgroupWriteEnabled {
		return fmt.Errorf("runtime mutation surfaces must remain disabled")
	}
	if replay.QMPCommandExecutionAllowed || replay.QGACommandExecutionAllowed {
		return fmt.Errorf("qmp/qga execution must remain disabled")
	}
	if replay.AutonomousApplyAllowed || replay.ProductionMutationAllowed {
		return fmt.Errorf("autonomous/production mutation must remain disabled")
	}
	if replay.WindowsGaClaimAllowed || replay.WindowsProductionReadyClaimAllowed || replay.WindowsExecutionReadyByDefault {
		return fmt.Errorf("windows ga/production/execution-ready claims must remain disabled")
	}
	if replay.ReplayValidationResult.ControlledApplyReady || replay.ReplayValidationResult.RuntimeMVPReady {
		return fmt.Errorf("controlled apply and runtime MVP readiness must remain false")
	}
	if !VerifyWindowsFluidVirtAuditHashChain(replay.AuditHashChain, replay.ReplayedEvents) {
		return fmt.Errorf("audit hash chain verification failed")
	}
	return nil
}

func LoadComplianceReplayFixtureFromTemporaryFile(path string) (WindowsFluidVirtComplianceReplay, error) {
	clean := filepath.Clean(path)
	if strings.HasPrefix(clean, "/sys/fs/cgroup") {
		return WindowsFluidVirtComplianceReplay{}, fmt.Errorf("real cgroup path forbidden in compliance replay fixture loading")
	}
	if !strings.Contains(clean, os.TempDir()) {
		return WindowsFluidVirtComplianceReplay{}, fmt.Errorf("compliance replay fixture must be loaded from temporary path")
	}
	raw, err := os.ReadFile(clean)
	if err != nil {
		return WindowsFluidVirtComplianceReplay{}, err
	}
	var replay WindowsFluidVirtComplianceReplay
	if err := json.Unmarshal(raw, &replay); err != nil {
		return WindowsFluidVirtComplianceReplay{}, err
	}
	if err := ValidateWindowsFluidVirtComplianceReplay(replay); err != nil {
		return WindowsFluidVirtComplianceReplay{}, err
	}
	return replay, nil
}
