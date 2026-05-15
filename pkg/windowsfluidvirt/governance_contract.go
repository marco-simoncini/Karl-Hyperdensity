package windowsfluidvirt

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

type GovernanceRequestedAction string

const (
	GovernanceFutureCPUApply             GovernanceRequestedAction = "future-cpu-apply"
	GovernanceFutureMemoryApply          GovernanceRequestedAction = "future-memory-apply"
	GovernanceFutureReturnToFloor        GovernanceRequestedAction = "future-return-to-floor"
	GovernanceFutureRollback             GovernanceRequestedAction = "future-rollback"
	GovernanceEvidenceRefreshBeforeApply GovernanceRequestedAction = "evidence-refresh-before-apply"
)

type GovernancePhase string

const (
	GovernanceContractPrepared    GovernancePhase = "CONTRACT_PREPARED"
	GovernanceContractBlocked     GovernancePhase = "CONTRACT_BLOCKED"
	GovernanceContractQuarantined GovernancePhase = "CONTRACT_QUARANTINED"
	GovernanceNeedsRevalidation   GovernancePhase = "NEEDS_REVALIDATION"
)

type TransitionPhase string

const (
	TransitionDryRunReady                    TransitionPhase = "DRYRUN_READY"
	TransitionLeasePrepared                  TransitionPhase = "LEASE_PREPARED"
	TransitionAdmittedForFutureApply         TransitionPhase = "ADMITTED_FOR_FUTURE_APPLY"
	TransitionContractPrepared               TransitionPhase = "CONTRACT_PREPARED"
	TransitionNeedsRevalidation              TransitionPhase = "NEEDS_REVALIDATION"
	TransitionFutureApplyEligibleTheoretical TransitionPhase = "FUTURE_APPLY_ELIGIBLE"
	TransitionContractBlocked                TransitionPhase = "CONTRACT_BLOCKED"
	TransitionContractQuarantined            TransitionPhase = "CONTRACT_QUARANTINED"
)

type RuntimeInvariantCheck struct {
	ID                 string `json:"id"`
	Description        string `json:"description"`
	Required           bool   `json:"required"`
	Observed           bool   `json:"observed"`
	Passed             bool   `json:"passed"`
	BlockerIfFailed    string `json:"blockerIfFailed"`
	QuarantineIfFailed bool   `json:"quarantineIfFailed"`
	EvidenceRef        string `json:"evidenceRef"`
}

type WindowsFluidRuntimeInvariantSet struct {
	InvariantSetID                      string                `json:"invariantSetId"`
	NoLiveMigrationInvariant            RuntimeInvariantCheck `json:"noLiveMigrationInvariant"`
	NoVMIMInvariant                     RuntimeInvariantCheck `json:"noVMIMInvariant"`
	NoRecreateInvariant                 RuntimeInvariantCheck `json:"noRecreateInvariant"`
	NoRolloutInvariant                  RuntimeInvariantCheck `json:"noRolloutInvariant"`
	NoRebootInvariant                   RuntimeInvariantCheck `json:"noRebootInvariant"`
	SameNodeInvariant                   RuntimeInvariantCheck `json:"sameNodeInvariant"`
	SameVirtLauncherPodInvariant        RuntimeInvariantCheck `json:"sameVirtLauncherPodInvariant"`
	SameQemuProcessInvariant            RuntimeInvariantCheck `json:"sameQemuProcessInvariant"`
	SameWindowsBootInvariant            RuntimeInvariantCheck `json:"sameWindowsBootInvariant"`
	SameMachineIdentityInvariant        RuntimeInvariantCheck `json:"sameMachineIdentityInvariant"`
	QMPAckInvariant                     RuntimeInvariantCheck `json:"qmpAckInvariant"`
	GuestAckInvariant                   RuntimeInvariantCheck `json:"guestAckInvariant"`
	RollbackReadyInvariant              RuntimeInvariantCheck `json:"rollbackReadyInvariant"`
	ReturnToFloorReadyInvariant         RuntimeInvariantCheck `json:"returnToFloorReadyInvariant"`
	KillSwitchReadyInvariant            RuntimeInvariantCheck `json:"killSwitchReadyInvariant"`
	EvidenceFreshnessInvariant          RuntimeInvariantCheck `json:"evidenceFreshnessInvariant"`
	QMPReadOnlyUntilApplyPhaseInvariant RuntimeInvariantCheck `json:"qmpReadOnlyUntilApplyPhaseInvariant"`
}

type StaleEvidencePolicy string

const (
	StalePolicyNeedsRevalidation StaleEvidencePolicy = "NEEDS_REVALIDATION"
	StalePolicyBlocked           StaleEvidencePolicy = "BLOCKED"
)

type RevalidationOutput string

const (
	RevalidationReady       RevalidationOutput = "REVALIDATION_READY"
	RevalidationBlocked     RevalidationOutput = "REVALIDATION_BLOCKED"
	RevalidationQuarantined RevalidationOutput = "REVALIDATION_QUARANTINED"
	RevalidationStale       RevalidationOutput = "REVALIDATION_STALE"
)

type RequiredFreshEvidence struct {
	KubeVirtIdentity       bool `json:"kubeVirtIdentity"`
	QMPEvidence            bool `json:"qmpEvidence"`
	GuestEvidence          bool `json:"guestEvidence"`
	RollbackReadiness      bool `json:"rollbackReadiness"`
	ReturnToFloorReadiness bool `json:"returnToFloorReadiness"`
	KillSwitchReadiness    bool `json:"killSwitchReadiness"`
}

type RequiredComparisons struct {
	NodeUnchanged            bool `json:"nodeUnchanged"`
	VirtLauncherPodUnchanged bool `json:"virtLauncherPodUnchanged"`
	QemuPIDUnchanged         bool `json:"qemuPidUnchanged"`
	LastBootUnchanged        bool `json:"lastBootUnchanged"`
	MachineIdentityUnchanged bool `json:"machineIdentityUnchanged"`
	NoVMIMObserved           bool `json:"noVmimObserved"`
	NoLiveMigrationObserved  bool `json:"noLiveMigrationObserved"`
	NoRecreateObserved       bool `json:"noRecreateObserved"`
}

type WindowsFluidPreApplyRevalidationContract struct {
	RevalidationID        string                `json:"revalidationId"`
	ShellRef              string                `json:"shellRef"`
	GovernanceContractRef string                `json:"governanceContractRef"`
	MaxEvidenceAgeSeconds int64                 `json:"maxEvidenceAgeSeconds"`
	RequiredFreshEvidence RequiredFreshEvidence `json:"requiredFreshEvidence"`
	RequiredComparisons   RequiredComparisons   `json:"requiredComparisons"`
	StaleEvidencePolicy   StaleEvidencePolicy   `json:"staleEvidencePolicy"`
	OutputAllowed         RevalidationOutput    `json:"outputAllowed"`
}

type AttestationSignature struct {
	Mode  string `json:"mode"`
	Value string `json:"value"`
}

type Attestor struct {
	ComponentName    string `json:"componentName"`
	ComponentVersion string `json:"componentVersion"`
}

type WindowsFluidPolicyAttestation struct {
	AttestationID     string               `json:"attestationId"`
	SubjectRef        string               `json:"subjectRef"`
	SubjectType       string               `json:"subjectType"`
	PolicyVersion     string               `json:"policyVersion"`
	EvidenceRefs      []string             `json:"evidenceRefs"`
	BlockerSnapshot   []string             `json:"blockerSnapshot"`
	InvariantSnapshot map[string]bool      `json:"invariantSnapshot"`
	DecisionSnapshot  map[string]any       `json:"decisionSnapshot"`
	CreatedAt         string               `json:"createdAt"`
	Attestor          Attestor             `json:"attestor"`
	Signature         AttestationSignature `json:"signature"`
	AuditHash         string               `json:"auditHash,omitempty"`
}

type WindowsFluidTransitionProof struct {
	ProofID         string          `json:"proofId"`
	FromPhase       TransitionPhase `json:"fromPhase"`
	ToPhase         TransitionPhase `json:"toPhase"`
	Allowed         bool            `json:"allowed"`
	Reason          string          `json:"reason"`
	RequiredInputs  []string        `json:"requiredInputs"`
	ObservedInputs  []string        `json:"observedInputs"`
	MissingInputs   []string        `json:"missingInputs"`
	BlockerList     []string        `json:"blockerList"`
	InvariantChecks map[string]bool `json:"invariantChecks"`
	Timestamp       string          `json:"timestamp"`
	ProofHash       string          `json:"proofHash,omitempty"`
	AuditRefs       []string        `json:"auditRefs"`
}

type WindowsFluidApplyGovernanceContract struct {
	ContractID                     string                    `json:"contractId"`
	ShellRef                       string                    `json:"shellRef"`
	SourceAdmissionDecisionID      string                    `json:"sourceAdmissionDecisionId"`
	SourceActionSlateID            string                    `json:"sourceActionSlateId"`
	RequestedAction                GovernanceRequestedAction `json:"requestedAction"`
	GovernancePhase                GovernancePhase           `json:"governancePhase"`
	MutationAllowed                bool                      `json:"mutationAllowed"`
	ApplyAllowed                   bool                      `json:"applyAllowed"`
	RuntimeMode                    string                    `json:"runtimeMode"`
	PolicyVersion                  string                    `json:"policyVersion"`
	EvidenceScoreAtAdmission       int64                     `json:"evidenceScoreAtAdmission"`
	RequiredEvidenceAtApplyTime    []string                  `json:"requiredEvidenceAtApplyTime"`
	RequiredFreshnessWindowSeconds int64                     `json:"requiredFreshnessWindowSeconds"`
	RequiredPreApplyRevalidation   bool                      `json:"requiredPreApplyRevalidation"`
	RequiredPostApplyVerification  []string                  `json:"requiredPostApplyVerification"`
	RollbackRequirement            string                    `json:"rollbackRequirement"`
	ReturnToFloorRequirement       string                    `json:"returnToFloorRequirement"`
	BlastRadiusLimit               string                    `json:"blastRadiusLimit"`
	KillSwitchRequirement          bool                      `json:"killSwitchRequirement"`
	AuditRequirement               string                    `json:"auditRequirement"`
	InvariantSetRef                string                    `json:"invariantSetRef"`
	TransitionProofRef             string                    `json:"transitionProofRef"`
	Blockers                       []string                  `json:"blockers"`
	DenialReasons                  []string                  `json:"denialReasons"`
	CreatedAt                      string                    `json:"createdAt"`
}

type ApplyGovernanceEvaluationInput struct {
	AdmissionDecision WindowsFluidAdmissionDecision     `json:"admissionDecision"`
	Bundle            WindowsFluidRuntimeEvidenceBundle `json:"bundle"`
	PolicyPack        *WindowsFluidPolicyPack           `json:"policyPack"`
	RequestedAction   GovernanceRequestedAction         `json:"requestedAction"`
	EvaluationTime    time.Time                         `json:"evaluationTime"`
}

type ApplyGovernanceEvaluationResult struct {
	GovernanceContract   WindowsFluidApplyGovernanceContract      `json:"governanceContract"`
	TransitionProof      WindowsFluidTransitionProof              `json:"transitionProof"`
	RuntimeInvariantSet  WindowsFluidRuntimeInvariantSet          `json:"runtimeInvariantSet"`
	PreApplyRevalidation WindowsFluidPreApplyRevalidationContract `json:"preApplyRevalidation"`
	PolicyAttestation    WindowsFluidPolicyAttestation            `json:"policyAttestation"`
	FinalGovernancePhase GovernancePhase                          `json:"finalGovernancePhase"`
	NextSafeStep         string                                   `json:"nextSafeStep"`
}

func EvaluateWindowsFluidApplyGovernance(input ApplyGovernanceEvaluationInput) ApplyGovernanceEvaluationResult {
	policy := DefaultWindowsFluidPolicyPack()
	if input.PolicyPack != nil {
		policy = *input.PolicyPack
	}
	evaluationTime := input.EvaluationTime.UTC()
	if evaluationTime.IsZero() {
		evaluationTime = time.Now().UTC()
	}

	admission := input.AdmissionDecision
	if admission.DecisionID == "" {
		admissionEval := EvaluateWindowsFluidAdmission(AdmissionEvaluationInput{
			Bundle:          input.Bundle,
			PolicyPack:      &policy,
			RequestedAction: mapGovernanceActionToAdmission(input.RequestedAction),
			EvaluationTime:  evaluationTime,
		})
		admission = admissionEval.Decision
	}
	dryRun := EvaluateWindowsFluidRuntimeDryRunWithOptions(input.Bundle, DryRunEvaluationOptions{EvaluationTime: evaluationTime})

	requestedAction := input.RequestedAction
	if requestedAction == "" {
		requestedAction = mapAdmissionActionToGovernance(admission.RequestedAction)
	}
	invariantSet, invariantBlockers, invariantQuarantine := buildRuntimeInvariantSet(input.Bundle, dryRun, policy, evaluationTime)
	combinedBlockers := dedupe(append(admission.Blockers, invariantBlockers...))

	revalidation := buildPreApplyRevalidationContract(input.Bundle, policy, invariantSet, evaluationTime)

	phase := GovernanceContractPrepared
	denialReasons := make([]string, 0, 8)
	requiredEvidence := make([]string, 0, 12)

	switch admission.AdmissionPhase {
	case AdmissionQuarantined:
		phase = GovernanceContractQuarantined
		denialReasons = append(denialReasons, "source_admission_quarantined")
	case AdmissionAdmittedForFutureApply:
	default:
		phase = GovernanceContractBlocked
		denialReasons = append(denialReasons, "source_admission_not_admitted")
	}

	if isPoolReplica(input.Bundle.SourceMetadata.SourceName) || input.Bundle.PolicyGates.PoolReplicaContextOnly {
		phase = GovernanceContractBlocked
		denialReasons = append(denialReasons, "pool_replica_context_only")
	}
	if dryRun.Classification == ClassificationBlockedGenericWindowsVM {
		phase = GovernanceContractBlocked
		denialReasons = append(denialReasons, "generic_windows_vm_not_certified")
	}
	if invariantQuarantine || hasPriority(combinedBlockers, policy, BlockerPriorityP0Quarantine) {
		phase = GovernanceContractQuarantined
		denialReasons = append(denialReasons, "quarantine_invariant_or_p0_blocker")
	}
	if phase != GovernanceContractQuarantined && hasPriority(combinedBlockers, policy, BlockerPriorityP1HardBlock) {
		phase = GovernanceContractBlocked
		denialReasons = append(denialReasons, "p1_hard_blocker_present")
	}
	if revalidation.OutputAllowed == RevalidationStale && phase == GovernanceContractPrepared {
		phase = GovernanceNeedsRevalidation
		denialReasons = append(denialReasons, "stale_evidence_requires_revalidation")
	}
	if revalidation.OutputAllowed == RevalidationBlocked && phase == GovernanceContractPrepared {
		phase = GovernanceContractBlocked
		denialReasons = append(denialReasons, "preapply_revalidation_blocked")
	}
	if revalidation.OutputAllowed == RevalidationQuarantined {
		phase = GovernanceContractQuarantined
		denialReasons = append(denialReasons, "preapply_revalidation_quarantined")
	}
	if !invariantSet.RollbackReadyInvariant.Passed {
		phase = GovernanceContractBlocked
		combinedBlockers = append(combinedBlockers, BlockerRollbackNotReady)
		denialReasons = append(denialReasons, "rollback_requirement_missing")
	}
	if !invariantSet.ReturnToFloorReadyInvariant.Passed {
		phase = GovernanceContractBlocked
		combinedBlockers = append(combinedBlockers, BlockerReturnToFloorNotReady)
		denialReasons = append(denialReasons, "return_to_floor_requirement_missing")
	}
	if !invariantSet.KillSwitchReadyInvariant.Passed {
		phase = GovernanceContractBlocked
		denialReasons = append(denialReasons, "kill_switch_requirement_missing")
	}
	if requestedAction == GovernanceFutureMemoryApply {
		if !invariantSet.ReturnToFloorReadyInvariant.Passed || !invariantSet.QMPReadOnlyUntilApplyPhaseInvariant.Passed {
			phase = GovernanceContractBlocked
			denialReasons = append(denialReasons, "memory_safety_not_proven")
		}
		if input.Bundle.Guest == nil || !input.Bundle.Guest.MemoryAdapterVerified {
			phase = GovernanceContractBlocked
			combinedBlockers = append(combinedBlockers, BlockerMemoryDriverUnverified)
			denialReasons = append(denialReasons, "memory_driver_not_verified")
		}
	}

	requiredEvidence = append(requiredEvidence,
		"fresh_kubevirt_identity_evidence",
		"fresh_qmp_evidence",
		"fresh_guest_evidence",
		"fresh_rollback_readiness",
		"fresh_return_to_floor_readiness",
		"kill_switch_readiness_proof",
	)
	if phase == GovernanceNeedsRevalidation || revalidation.OutputAllowed == RevalidationStale {
		requiredEvidence = append(requiredEvidence, "revalidation_evidence_refresh")
	}

	transitionProof := buildTransitionProof(admission, phase, combinedBlockers, invariantSet, evaluationTime)
	contract := WindowsFluidApplyGovernanceContract{
		ContractID:                     "windows-fluid-governance-" + evaluationTime.Format("20060102T150405"),
		ShellRef:                       input.Bundle.Shell.Spec.VMRef,
		SourceAdmissionDecisionID:      admission.DecisionID,
		SourceActionSlateID:            "windows-fluid-action-slate-reference",
		RequestedAction:                requestedAction,
		GovernancePhase:                phase,
		MutationAllowed:                false,
		ApplyAllowed:                   false,
		RuntimeMode:                    "in-place-qmp",
		PolicyVersion:                  policy.PolicyVersion,
		EvidenceScoreAtAdmission:       admission.EvidenceScore,
		RequiredEvidenceAtApplyTime:    dedupe(requiredEvidence),
		RequiredFreshnessWindowSeconds: policy.MaxEvidenceAgeSeconds,
		RequiredPreApplyRevalidation:   true,
		RequiredPostApplyVerification: []string{
			"continuity_post_apply_proof",
			"rollback_readiness_post_apply",
			"return_to_floor_proof_post_apply",
			"no_migration_no_recreate_no_reboot_post_apply",
		},
		RollbackRequirement:      "mandatory",
		ReturnToFloorRequirement: "mandatory",
		BlastRadiusLimit:         policy.BlastRadiusPolicy.Scope,
		KillSwitchRequirement:    !policy.AllowMutationInThisPhase,
		AuditRequirement:         "mandatory-audit-trail",
		InvariantSetRef:          invariantSet.InvariantSetID,
		TransitionProofRef:       transitionProof.ProofID,
		Blockers:                 dedupe(combinedBlockers),
		DenialReasons:            dedupe(denialReasons),
		CreatedAt:                evaluationTime.Format(time.RFC3339),
	}
	revalidation.GovernanceContractRef = contract.ContractID
	attestation := buildPolicyAttestation(contract, transitionProof, invariantSet, policy, evaluationTime)

	nextStep := "refresh_evidence_and_revalidate_contract"
	switch phase {
	case GovernanceContractPrepared:
		nextStep = "hold_contract_as_non_executable_and_prepare_future_phase_review"
	case GovernanceContractBlocked:
		nextStep = "resolve_governance_blockers_before_next_contract_evaluation"
	case GovernanceContractQuarantined:
		nextStep = "isolate_candidate_and_rebuild_identity_continuity_proofs"
	case GovernanceNeedsRevalidation:
		nextStep = "refresh_stale_evidence_and_rerun_preapply_revalidation"
	}

	return ApplyGovernanceEvaluationResult{
		GovernanceContract:   contract,
		TransitionProof:      transitionProof,
		RuntimeInvariantSet:  invariantSet,
		PreApplyRevalidation: revalidation,
		PolicyAttestation:    attestation,
		FinalGovernancePhase: phase,
		NextSafeStep:         nextStep,
	}
}

func buildRuntimeInvariantSet(
	bundle WindowsFluidRuntimeEvidenceBundle,
	dryRun DryRunEvaluationResult,
	policy WindowsFluidPolicyPack,
	evaluationTime time.Time,
) (WindowsFluidRuntimeInvariantSet, []string, bool) {
	evidenceRef := bundle.Shell.Status.EvidenceRef
	invariantSet := WindowsFluidRuntimeInvariantSet{
		InvariantSetID: "windows-fluid-invariant-set-" + evaluationTime.Format("20060102T150405"),
		NoLiveMigrationInvariant: createInvariant(
			"no_live_migration_invariant",
			"No LiveMigration objects are observed",
			true,
			!bundle.KubeVirtAfter.MigrationRequired && len(bundle.KubeVirtAfter.LiveMigrationObjectsObserved) == 0,
			BlockerLiveMigrationRequired,
			false,
			evidenceRef,
		),
		NoVMIMInvariant: createInvariant(
			"no_vmim_invariant",
			"No VMIM objects are observed",
			true,
			len(bundle.KubeVirtAfter.VMIMObjectsObserved) == 0,
			BlockerLiveMigrationRequired,
			false,
			evidenceRef,
		),
		NoRecreateInvariant: createInvariant(
			"no_recreate_invariant",
			"VMI recreate is not required",
			true,
			!bundle.KubeVirtAfter.RecreateRequired,
			BlockerVMIRecreateRequired,
			false,
			evidenceRef,
		),
		NoRolloutInvariant: createInvariant(
			"no_rollout_invariant",
			"No rollout is observed",
			true,
			!bundle.KubeVirtAfter.RolloutObserved,
			BlockerVMIRecreateRequired,
			false,
			evidenceRef,
		),
		NoRebootInvariant: createInvariant(
			"no_reboot_invariant",
			"Windows guest has no pending reboot",
			true,
			bundle.Guest != nil && !bundle.Guest.PendingReboot,
			BlockerPendingRebootDetected,
			false,
			evidenceRef,
		),
		SameNodeInvariant: createInvariant(
			"same_node_invariant",
			"Node continuity is preserved",
			policy.RequireSameNode,
			dryRun.Conditions["sameNode"],
			BlockerNodeChanged,
			true,
			evidenceRef,
		),
		SameVirtLauncherPodInvariant: createInvariant(
			"same_virtlauncher_pod_invariant",
			"virt-launcher pod continuity is preserved",
			policy.RequireSameVirtLauncherPod,
			dryRun.Conditions["sameVirtLauncherPod"],
			BlockerVirtLauncherPodChanged,
			true,
			evidenceRef,
		),
		SameQemuProcessInvariant: createInvariant(
			"same_qemu_process_invariant",
			"QEMU process continuity is preserved",
			policy.RequireSameQemuProcess,
			dryRun.Conditions["sameQemuProcess"],
			BlockerQemuPIDChanged,
			true,
			evidenceRef,
		),
		SameWindowsBootInvariant: createInvariant(
			"same_windows_boot_invariant",
			"Windows last boot continuity is preserved",
			policy.RequireSameLastBoot,
			bundle.Guest != nil && bundle.Guest.LastBootTime != "",
			BlockerLastBootChanged,
			true,
			evidenceRef,
		),
		SameMachineIdentityInvariant: createInvariant(
			"same_machine_identity_invariant",
			"Machine identity continuity is preserved",
			policy.RequireSameMachineIdentity,
			bundle.Guest != nil && bundle.Guest.MachineGUIDHash != "",
			BlockerMachineGUIDChanged,
			true,
			evidenceRef,
		),
		QMPAckInvariant: createInvariant(
			"qmp_ack_invariant",
			"QMP ACK is available and healthy",
			policy.RequireQmpAck,
			dryRun.Conditions["qmpReady"],
			BlockerQMPAckMissing,
			false,
			evidenceRef,
		),
		GuestAckInvariant: createInvariant(
			"guest_ack_invariant",
			"Guest ACK is available and healthy",
			policy.RequireGuestAck,
			dryRun.Conditions["guestAckReady"],
			BlockerGuestAckMissing,
			false,
			evidenceRef,
		),
		RollbackReadyInvariant: createInvariant(
			"rollback_ready_invariant",
			"Rollback readiness is proven",
			policy.RequireRollbackReady,
			dryRun.Conditions["rollbackReady"],
			BlockerRollbackNotReady,
			false,
			evidenceRef,
		),
		ReturnToFloorReadyInvariant: createInvariant(
			"return_to_floor_ready_invariant",
			"Return-to-floor readiness is proven",
			policy.RequireReturnToFloorReady,
			dryRun.Conditions["returnToFloorReady"],
			BlockerReturnToFloorNotReady,
			false,
			evidenceRef,
		),
		KillSwitchReadyInvariant: createInvariant(
			"kill_switch_ready_invariant",
			"Kill-switch governance is enabled for this phase",
			true,
			!policy.AllowMutationInThisPhase,
			BlockerRollbackNotReady,
			false,
			evidenceRef,
		),
		EvidenceFreshnessInvariant: createInvariant(
			"evidence_freshness_invariant",
			"Evidence is fresh inside policy window",
			true,
			evidenceFresh(bundle, policy, evaluationTime),
			"",
			false,
			evidenceRef,
		),
		QMPReadOnlyUntilApplyPhaseInvariant: createInvariant(
			"qmp_readonly_until_apply_phase_invariant",
			"QMP stays read-only in governance phase",
			true,
			bundle.QMP != nil && bundle.QMP.QMPReadOnly,
			BlockerHotplugErrorDetected,
			true,
			evidenceRef,
		),
	}

	allChecks := []RuntimeInvariantCheck{
		invariantSet.NoLiveMigrationInvariant,
		invariantSet.NoVMIMInvariant,
		invariantSet.NoRecreateInvariant,
		invariantSet.NoRolloutInvariant,
		invariantSet.NoRebootInvariant,
		invariantSet.SameNodeInvariant,
		invariantSet.SameVirtLauncherPodInvariant,
		invariantSet.SameQemuProcessInvariant,
		invariantSet.SameWindowsBootInvariant,
		invariantSet.SameMachineIdentityInvariant,
		invariantSet.QMPAckInvariant,
		invariantSet.GuestAckInvariant,
		invariantSet.RollbackReadyInvariant,
		invariantSet.ReturnToFloorReadyInvariant,
		invariantSet.KillSwitchReadyInvariant,
		invariantSet.EvidenceFreshnessInvariant,
		invariantSet.QMPReadOnlyUntilApplyPhaseInvariant,
	}

	blockers := make([]string, 0, len(allChecks))
	quarantine := false
	for _, check := range allChecks {
		if !check.Passed && check.BlockerIfFailed != "" {
			blockers = append(blockers, check.BlockerIfFailed)
		}
		if !check.Passed && check.QuarantineIfFailed {
			quarantine = true
		}
	}
	return invariantSet, dedupe(blockers), quarantine
}

func createInvariant(
	id, description string,
	required bool,
	observed bool,
	blocker string,
	quarantine bool,
	evidenceRef string,
) RuntimeInvariantCheck {
	passed := observed
	if !required {
		passed = true
	}
	return RuntimeInvariantCheck{
		ID:                 id,
		Description:        description,
		Required:           required,
		Observed:           observed,
		Passed:             passed,
		BlockerIfFailed:    blocker,
		QuarantineIfFailed: quarantine,
		EvidenceRef:        evidenceRef,
	}
}

func buildPreApplyRevalidationContract(
	bundle WindowsFluidRuntimeEvidenceBundle,
	policy WindowsFluidPolicyPack,
	invariants WindowsFluidRuntimeInvariantSet,
	evaluationTime time.Time,
) WindowsFluidPreApplyRevalidationContract {
	output := RevalidationReady
	if !invariants.EvidenceFreshnessInvariant.Passed {
		output = RevalidationStale
	}
	if !invariants.SameNodeInvariant.Passed ||
		!invariants.SameVirtLauncherPodInvariant.Passed ||
		!invariants.SameQemuProcessInvariant.Passed ||
		!invariants.SameWindowsBootInvariant.Passed ||
		!invariants.SameMachineIdentityInvariant.Passed {
		output = RevalidationQuarantined
	}
	if output != RevalidationQuarantined && (!invariants.GuestAckInvariant.Passed ||
		!invariants.QMPAckInvariant.Passed ||
		!invariants.RollbackReadyInvariant.Passed ||
		!invariants.ReturnToFloorReadyInvariant.Passed) {
		output = RevalidationBlocked
	}

	return WindowsFluidPreApplyRevalidationContract{
		RevalidationID:        "windows-fluid-revalidation-" + evaluationTime.Format("20060102T150405"),
		ShellRef:              bundle.Shell.Spec.VMRef,
		GovernanceContractRef: "",
		MaxEvidenceAgeSeconds: policy.MaxEvidenceAgeSeconds,
		RequiredFreshEvidence: RequiredFreshEvidence{
			KubeVirtIdentity:       true,
			QMPEvidence:            true,
			GuestEvidence:          true,
			RollbackReadiness:      true,
			ReturnToFloorReadiness: true,
			KillSwitchReadiness:    true,
		},
		RequiredComparisons: RequiredComparisons{
			NodeUnchanged:            true,
			VirtLauncherPodUnchanged: true,
			QemuPIDUnchanged:         true,
			LastBootUnchanged:        true,
			MachineIdentityUnchanged: true,
			NoVMIMObserved:           true,
			NoLiveMigrationObserved:  true,
			NoRecreateObserved:       true,
		},
		StaleEvidencePolicy: StalePolicyNeedsRevalidation,
		OutputAllowed:       output,
	}
}

func buildPolicyAttestation(
	contract WindowsFluidApplyGovernanceContract,
	proof WindowsFluidTransitionProof,
	invariants WindowsFluidRuntimeInvariantSet,
	policy WindowsFluidPolicyPack,
	evaluationTime time.Time,
) WindowsFluidPolicyAttestation {
	invariantSnapshot := map[string]bool{
		invariants.NoLiveMigrationInvariant.ID:            invariants.NoLiveMigrationInvariant.Passed,
		invariants.NoVMIMInvariant.ID:                     invariants.NoVMIMInvariant.Passed,
		invariants.NoRecreateInvariant.ID:                 invariants.NoRecreateInvariant.Passed,
		invariants.NoRolloutInvariant.ID:                  invariants.NoRolloutInvariant.Passed,
		invariants.NoRebootInvariant.ID:                   invariants.NoRebootInvariant.Passed,
		invariants.SameNodeInvariant.ID:                   invariants.SameNodeInvariant.Passed,
		invariants.SameVirtLauncherPodInvariant.ID:        invariants.SameVirtLauncherPodInvariant.Passed,
		invariants.SameQemuProcessInvariant.ID:            invariants.SameQemuProcessInvariant.Passed,
		invariants.SameWindowsBootInvariant.ID:            invariants.SameWindowsBootInvariant.Passed,
		invariants.SameMachineIdentityInvariant.ID:        invariants.SameMachineIdentityInvariant.Passed,
		invariants.QMPAckInvariant.ID:                     invariants.QMPAckInvariant.Passed,
		invariants.GuestAckInvariant.ID:                   invariants.GuestAckInvariant.Passed,
		invariants.RollbackReadyInvariant.ID:              invariants.RollbackReadyInvariant.Passed,
		invariants.ReturnToFloorReadyInvariant.ID:         invariants.ReturnToFloorReadyInvariant.Passed,
		invariants.KillSwitchReadyInvariant.ID:            invariants.KillSwitchReadyInvariant.Passed,
		invariants.EvidenceFreshnessInvariant.ID:          invariants.EvidenceFreshnessInvariant.Passed,
		invariants.QMPReadOnlyUntilApplyPhaseInvariant.ID: invariants.QMPReadOnlyUntilApplyPhaseInvariant.Passed,
	}
	return WindowsFluidPolicyAttestation{
		AttestationID:     "windows-fluid-attestation-" + evaluationTime.Format("20060102T150405"),
		SubjectRef:        contract.ContractID,
		SubjectType:       "governance-contract",
		PolicyVersion:     policy.PolicyVersion,
		EvidenceRefs:      []string{contract.ShellRef},
		BlockerSnapshot:   contract.Blockers,
		InvariantSnapshot: invariantSnapshot,
		DecisionSnapshot: map[string]any{
			"governancePhase": contract.GovernancePhase,
			"requestedAction": contract.RequestedAction,
			"mutationAllowed": contract.MutationAllowed,
			"applyAllowed":    contract.ApplyAllowed,
		},
		CreatedAt: evaluationTime.Format(time.RFC3339),
		Attestor: Attestor{
			ComponentName:    "karl-fluid-governance-evaluator",
			ComponentVersion: "v1",
		},
		Signature: AttestationSignature{
			Mode:  "unsigned-dev",
			Value: "",
		},
		AuditHash: proof.ProofHash,
	}
}

func buildTransitionProof(
	admission WindowsFluidAdmissionDecision,
	phase GovernancePhase,
	blockers []string,
	invariants WindowsFluidRuntimeInvariantSet,
	evaluationTime time.Time,
) WindowsFluidTransitionProof {
	from := TransitionAdmittedForFutureApply
	switch admission.AdmissionPhase {
	case AdmissionAdmittedForFutureApply:
		from = TransitionAdmittedForFutureApply
	case AdmissionBlocked:
		from = TransitionLeasePrepared
	case AdmissionDenied:
		from = TransitionDryRunReady
	case AdmissionQuarantined:
		from = TransitionLeasePrepared
	default:
		from = TransitionDryRunReady
	}

	to := TransitionContractPrepared
	switch phase {
	case GovernanceContractPrepared:
		to = TransitionContractPrepared
	case GovernanceNeedsRevalidation:
		to = TransitionNeedsRevalidation
	case GovernanceContractBlocked:
		to = TransitionContractBlocked
	case GovernanceContractQuarantined:
		to = TransitionContractQuarantined
	}
	allowed := phase == GovernanceContractPrepared
	requiredInputs := []string{
		"admission_decision",
		"runtime_evidence_bundle",
		"policy_pack",
		"invariant_set",
		"preapply_revalidation",
	}
	observedInputs := []string{
		"admission_decision",
		"runtime_evidence_bundle",
		"policy_pack",
	}
	missingInputs := make([]string, 0, 4)
	for _, required := range requiredInputs {
		if !contains(observedInputs, required) {
			missingInputs = append(missingInputs, required)
		}
	}
	invariantChecks := map[string]bool{
		invariants.SameNodeInvariant.ID:                   invariants.SameNodeInvariant.Passed,
		invariants.SameVirtLauncherPodInvariant.ID:        invariants.SameVirtLauncherPodInvariant.Passed,
		invariants.SameQemuProcessInvariant.ID:            invariants.SameQemuProcessInvariant.Passed,
		invariants.SameWindowsBootInvariant.ID:            invariants.SameWindowsBootInvariant.Passed,
		invariants.SameMachineIdentityInvariant.ID:        invariants.SameMachineIdentityInvariant.Passed,
		invariants.QMPReadOnlyUntilApplyPhaseInvariant.ID: invariants.QMPReadOnlyUntilApplyPhaseInvariant.Passed,
	}

	body := strings.Join([]string{
		string(from),
		string(to),
		evaluationTime.Format(time.RFC3339),
		strings.Join(dedupe(blockers), ","),
	}, "|")
	sum := sha256.Sum256([]byte(body))

	reason := "transition_allowed_for_non_executable_contract"
	if !allowed {
		reason = "transition_denied_by_governance_blockers"
	}
	return WindowsFluidTransitionProof{
		ProofID:         "windows-fluid-transition-proof-" + evaluationTime.Format("20060102T150405"),
		FromPhase:       from,
		ToPhase:         to,
		Allowed:         allowed,
		Reason:          reason,
		RequiredInputs:  requiredInputs,
		ObservedInputs:  observedInputs,
		MissingInputs:   missingInputs,
		BlockerList:     dedupe(blockers),
		InvariantChecks: invariantChecks,
		Timestamp:       evaluationTime.Format(time.RFC3339),
		ProofHash:       hex.EncodeToString(sum[:]),
		AuditRefs:       admission.AuditRefs,
	}
}

func mapGovernanceActionToAdmission(action GovernanceRequestedAction) RequestedAdmissionAction {
	switch action {
	case GovernanceFutureCPUApply:
		return RequestedActionPrepareCPULease
	case GovernanceFutureMemoryApply:
		return RequestedActionPrepareMemoryLease
	case GovernanceFutureReturnToFloor:
		return RequestedActionReturnToFloorCheck
	case GovernanceFutureRollback:
		return RequestedActionQuarantine
	default:
		return RequestedActionEvidenceRefresh
	}
}

func mapAdmissionActionToGovernance(action RequestedAdmissionAction) GovernanceRequestedAction {
	switch action {
	case RequestedActionPrepareCPULease:
		return GovernanceFutureCPUApply
	case RequestedActionPrepareMemoryLease:
		return GovernanceFutureMemoryApply
	case RequestedActionReturnToFloorCheck:
		return GovernanceFutureReturnToFloor
	case RequestedActionQuarantine:
		return GovernanceFutureRollback
	default:
		return GovernanceEvidenceRefreshBeforeApply
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
