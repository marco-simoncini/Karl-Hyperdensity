package windowsfluidvirt

type WindowsFluidVirtReleaseTrack string

const (
	ReleaseTrackTechnicalPreview WindowsFluidVirtReleaseTrack = "technical_preview"
)

type WindowsFluidVirtLaneStatus string

const (
	LaneStatusTechnicalPreviewCandidate WindowsFluidVirtLaneStatus = "technical_preview_candidate"
	LaneStatusGatedPreview              WindowsFluidVirtLaneStatus = "gated_preview"
)

type WindowsFluidVirtReadinessState string

const (
	ReadinessStateModelOnly        WindowsFluidVirtReadinessState = "model_only"
	ReadinessStatePlanningOnly     WindowsFluidVirtReadinessState = "planning_only"
	ReadinessStateNotExecutionReady WindowsFluidVirtReadinessState = "not_execution_ready"
)

type WindowsFluidVirtSupportBoundary struct {
	ProductionMutationAllowed              bool   `json:"productionMutationAllowed"`
	AutonomousApplyAllowed                 bool   `json:"autonomousApplyAllowed"`
	EnforcementMode                        string `json:"enforcementMode"`
	WindowsGaClaimAllowed                  bool   `json:"windowsGaClaimAllowed"`
	WindowsProductionReadyClaimAllowed     bool   `json:"windowsProductionReadyClaimAllowed"`
	WindowsExecutionReadyByDefault         bool   `json:"windowsExecutionReadyByDefault"`
	VCPUHotplugClaimAllowed                bool   `json:"vcpuHotplugClaimAllowed"`
	LogicalCPUScalingClaimAllowed          bool   `json:"logicalCpuScalingClaimAllowed"`
	PoolScalingClaimAllowed                bool   `json:"poolScalingClaimAllowed"`
	LiveMigrationClaimAllowed              bool   `json:"liveMigrationClaimAllowed"`
	RebootRecreateRolloutMechanismAllowed  bool   `json:"rebootRecreateRolloutMechanismAllowed"`
	RawRuntimeControlsExposed              bool   `json:"rawRuntimeControlsExposed"`
}

type WindowsFluidVirtClaimBoundary struct {
	TechnicalPreviewCandidate bool `json:"technicalPreviewCandidate"`
	GatedPreview              bool `json:"gatedPreview"`
	NoGA                      bool `json:"noGa"`
	NoProductionReady         bool `json:"noProductionReady"`
	NoAutonomousApply         bool `json:"noAutonomousApply"`
	NoRawControls             bool `json:"noRawControls"`
	NoVCPUHotplug             bool `json:"noVcpuHotplug"`
	NoLogicalCPUScaling       bool `json:"noLogicalCpuScaling"`
	NoPoolScaling             bool `json:"noPoolScaling"`
	NoLiveMigrationVMIM       bool `json:"noLiveMigrationVmim"`
	NoRebootRecreateRollout   bool `json:"noRebootRecreateRollout"`
}

type WindowsFluidVirtCPUEntitlementLiquidityModel struct {
	Mechanism          string `json:"mechanism"`
	NodeActuatorNeeded bool   `json:"nodeActuatorNeeded"`
	ApplyEnabled       bool   `json:"applyEnabled"`
}

type WindowsFluidVirtRAMBalloonLiquidityModel struct {
	Mechanism    string `json:"mechanism"`
	QMPModelOnly bool   `json:"qmpModelOnly"`
	ApplyEnabled bool   `json:"applyEnabled"`
}

type WindowsFluidVirtGuestWitnessDependency struct {
	KarlAgentFluidShellRequired bool `json:"karlAgentFluidShellRequired"`
	QGARequired                 bool `json:"qgaRequired"`
	GuestAckRequired            bool `json:"guestAckRequired"`
}

type WindowsFluidVirtLeaseModel struct {
	PlanningOnly          bool `json:"planningOnly"`
	RollbackPlanRequired  bool `json:"rollbackPlanRequired"`
	ReturnToFloorRequired bool `json:"returnToFloorRequired"`
	AuditChainRequired    bool `json:"auditChainRequired"`
	ApplyEnabled          bool `json:"applyEnabled"`
}

type WindowsFluidVirtProductModel struct {
	ProductID                 string                                  `json:"productId"`
	ProductVersion            string                                  `json:"productVersion"`
	ReleaseTrack              WindowsFluidVirtReleaseTrack            `json:"releaseTrack"`
	LaneStatus                WindowsFluidVirtLaneStatus              `json:"laneStatus"`
	ReadinessState            WindowsFluidVirtReadinessState          `json:"readinessState"`
	SupportBoundary           WindowsFluidVirtSupportBoundary         `json:"supportBoundary"`
	ClaimBoundary             WindowsFluidVirtClaimBoundary           `json:"claimBoundary"`
	CPULiquidityModel         WindowsFluidVirtCPUEntitlementLiquidityModel `json:"cpuLiquidityModel"`
	RAMLiquidityModel         WindowsFluidVirtRAMBalloonLiquidityModel `json:"ramLiquidityModel"`
	GuestWitnessDependency    WindowsFluidVirtGuestWitnessDependency  `json:"guestWitnessDependency"`
	LeaseModel                WindowsFluidVirtLeaseModel              `json:"leaseModel"`
	ActionSlate               WindowsFluidVirtActionSlate             `json:"actionSlate"`
	EvidenceRefs              []string                                `json:"evidenceRefs,omitempty"`
	Blockers                  []WindowsFluidVirtBlocker               `json:"blockers,omitempty"`
}

func NewWindowsFluidVirtProductModel(version string) WindowsFluidVirtProductModel {
	if version == "" {
		version = "v1"
	}
	return WindowsFluidVirtProductModel{
		ProductID:      "windows-fluidvirt-product-model-v1",
		ProductVersion: version,
		ReleaseTrack:   ReleaseTrackTechnicalPreview,
		LaneStatus:     LaneStatusTechnicalPreviewCandidate,
		ReadinessState: ReadinessStatePlanningOnly,
		SupportBoundary: WindowsFluidVirtSupportBoundary{
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
			RawRuntimeControlsExposed:             false,
		},
		ClaimBoundary: WindowsFluidVirtClaimBoundary{
			TechnicalPreviewCandidate: true,
			GatedPreview:              true,
			NoGA:                      true,
			NoProductionReady:         true,
			NoAutonomousApply:         true,
			NoRawControls:             true,
			NoVCPUHotplug:             true,
			NoLogicalCPUScaling:       true,
			NoPoolScaling:             true,
			NoLiveMigrationVMIM:       true,
			NoRebootRecreateRollout:   true,
		},
		CPULiquidityModel: WindowsFluidVirtCPUEntitlementLiquidityModel{
			Mechanism:          "entitlement_liquidity_cgroup_v2_cpu_max",
			NodeActuatorNeeded: true,
			ApplyEnabled:       false,
		},
		RAMLiquidityModel: WindowsFluidVirtRAMBalloonLiquidityModel{
			Mechanism:    "balloon_liquidity_qmp_model",
			QMPModelOnly: true,
			ApplyEnabled: false,
		},
		GuestWitnessDependency: WindowsFluidVirtGuestWitnessDependency{
			KarlAgentFluidShellRequired: true,
			QGARequired:                 true,
			GuestAckRequired:            true,
		},
		LeaseModel: WindowsFluidVirtLeaseModel{
			PlanningOnly:          true,
			RollbackPlanRequired:  true,
			ReturnToFloorRequired: true,
			AuditChainRequired:    true,
			ApplyEnabled:          false,
		},
		ActionSlate: NewDefaultWindowsFluidVirtActionSlate(version),
	}
}

func (m WindowsFluidVirtProductModel) ValidateSafety() []WindowsFluidVirtBlocker {
	blockers := make([]WindowsFluidVirtBlocker, 0, 12)
	if m.SupportBoundary.RawRuntimeControlsExposed || m.ActionSlate.RawRuntimeControlsExposed {
		blockers = append(blockers, BlockerRawRuntimeControlForbidden)
	}
	if m.SupportBoundary.AutonomousApplyAllowed || m.ActionSlate.AutonomousApplyEnabled {
		blockers = append(blockers, BlockerAutonomousApplyForbidden)
	}
	if m.SupportBoundary.WindowsProductionReadyClaimAllowed {
		blockers = append(blockers, BlockerProductionReadyClaimForbidden)
	}
	if m.SupportBoundary.VCPUHotplugClaimAllowed {
		blockers = append(blockers, BlockerVCPUHotplugClaimForbidden)
	}
	if m.SupportBoundary.LogicalCPUScalingClaimAllowed {
		blockers = append(blockers, BlockerLogicalCPUScalingClaimForbidden)
	}
	if m.SupportBoundary.PoolScalingClaimAllowed {
		blockers = append(blockers, BlockerPoolScalingClaimForbidden)
	}
	if !m.LeaseModel.RollbackPlanRequired {
		blockers = append(blockers, BlockerMissingRollbackPlan)
	}
	if !m.LeaseModel.ReturnToFloorRequired {
		blockers = append(blockers, BlockerMissingReturnToFloorPlan)
	}
	if !m.LeaseModel.AuditChainRequired {
		blockers = append(blockers, BlockerMissingAuditChain)
	}
	if !m.GuestWitnessDependency.GuestAckRequired {
		blockers = append(blockers, BlockerMissingGuestWitness)
	}
	blockers = append(blockers, m.ActionSlate.PlanningOnlyBlockers()...)
	return dedupeBlockers(blockers)
}

func dedupeBlockers(values []WindowsFluidVirtBlocker) []WindowsFluidVirtBlocker {
	seen := make(map[WindowsFluidVirtBlocker]struct{}, len(values))
	result := make([]WindowsFluidVirtBlocker, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
