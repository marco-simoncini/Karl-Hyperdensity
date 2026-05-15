package windowsfluidvirt

import "time"

type AdmissionEvaluationInput struct {
	Bundle         WindowsFluidRuntimeEvidenceBundle `json:"bundle"`
	PolicyPack     *WindowsFluidPolicyPack           `json:"policyPack"`
	RequestedAction RequestedAdmissionAction         `json:"requestedAction"`
	EvaluationTime time.Time                         `json:"evaluationTime"`
}

type AdmissionEvaluationResult struct {
	Decision        WindowsFluidAdmissionDecision `json:"decision"`
	ActionSlate     WindowsFluidActionSlate       `json:"actionSlate"`
	EvidenceScore   AdmissionEvidenceScore        `json:"evidenceScore"`
	PolicyViolations []string                     `json:"policyViolations"`
	Blockers        []string                      `json:"blockers"`
	NextSafeStep    string                        `json:"nextSafeStep"`
}

func EvaluateWindowsFluidAdmission(input AdmissionEvaluationInput) AdmissionEvaluationResult {
	policy := DefaultWindowsFluidPolicyPack()
	if input.PolicyPack != nil {
		policy = *input.PolicyPack
	}
	requestedAction := input.RequestedAction
	if requestedAction == "" {
		requestedAction = deriveRequestedAction(input.Bundle.LeaseIntent)
	}

	options := DryRunEvaluationOptions{EvaluationTime: input.EvaluationTime}
	dryRun := EvaluateWindowsFluidRuntimeDryRunWithOptions(input.Bundle, options)
	score := EvaluateAdmissionEvidenceScore(input.Bundle, dryRun, policy, input.EvaluationTime)

	blockers := dedupe(append([]string{}, dryRun.Blockers...))
	policyViolations := evaluatePolicyViolations(input.Bundle, dryRun, policy, requestedAction)
	denialReasons := make([]string, 0, 8)
	requiredEvidence := append([]string{}, score.MissingEvidence...)
	phase := AdmissionNeedsMoreEvidence

	if isPoolReplica(input.Bundle.SourceMetadata.SourceName) && !policy.AllowPoolReplicaModel {
		phase = AdmissionDenied
		denialReasons = append(denialReasons, "pool_replica_model_denied")
	}
	if dryRun.Classification == ClassificationBlockedGenericWindowsVM && !policy.AllowGenericWindowsVm {
		phase = AdmissionDenied
		denialReasons = append(denialReasons, "generic_windows_vm_denied")
	}
	if phase != AdmissionDenied && (dryRun.Phase == StateQuarantined || hasPriority(blockers, policy, BlockerPriorityP0Quarantine)) {
		phase = AdmissionQuarantined
		denialReasons = append(denialReasons, "identity_or_runtime_quarantine")
	}
	if phase != AdmissionDenied && phase != AdmissionQuarantined && hasPriority(blockers, policy, BlockerPriorityP1HardBlock) {
		phase = AdmissionBlocked
		denialReasons = append(denialReasons, "hard_blocker_present")
	}

	if phase != AdmissionDenied && phase != AdmissionQuarantined && phase != AdmissionBlocked {
		if hasPriority(blockers, policy, BlockerPriorityP2Capability) {
			phase = AdmissionNeedsMoreEvidence
			denialReasons = append(denialReasons, "capability_evidence_missing")
		}
	}

	if dryRun.Phase != StateReady && dryRun.Phase != StateLeasePrepared {
		if phase == AdmissionNeedsMoreEvidence {
			denialReasons = append(denialReasons, "dryrun_not_ready")
		}
	}

	if len(policyViolations) > 0 && phase != AdmissionDenied && phase != AdmissionQuarantined {
		phase = AdmissionBlocked
		denialReasons = append(denialReasons, "policy_violation")
		requiredEvidence = append(requiredEvidence, policyViolations...)
	}

	if phase == AdmissionNeedsMoreEvidence || phase == AdmissionBlocked {
		if requestedAction == RequestedActionPrepareMemoryLease && !input.Bundle.KubeVirtAfter.MigrationRequired {
			if input.Bundle.Guest == nil || !input.Bundle.Guest.MemoryAdapterVerified {
				requiredEvidence = append(requiredEvidence, "memory_driver_verification_required")
			}
		}
	}

	if phase != AdmissionDenied && phase != AdmissionBlocked && phase != AdmissionQuarantined {
		if requestedAction == RequestedActionPrepareCPULease &&
			(dryRun.Phase == StateReady || dryRun.Phase == StateLeasePrepared) &&
			score.Score >= policy.MinEvidenceScoreForFutureApply &&
			score.EvidenceLevel == EvidenceLevelFutureApplyAdmissible &&
			len(score.HardBlockers) == 0 {
			phase = AdmissionAdmittedForFutureApply
		} else if requestedAction == RequestedActionPrepareMemoryLease {
			if input.Bundle.Guest == nil || !input.Bundle.Guest.MemoryAdapterVerified || !dryRun.Conditions["returnToFloorReady"] {
				phase = AdmissionBlocked
				denialReasons = append(denialReasons, "memory_safety_not_proven")
				blockers = append(blockers, BlockerMemoryDriverUnverified, BlockerMemoryReturnNotSafe)
			} else {
				phase = AdmissionNeedsMoreEvidence
				denialReasons = append(denialReasons, "memory_action_requires_future_phase_review")
			}
		}
	}

	if !policy.AllowMutationInThisPhase {
		denialReasons = append(denialReasons, "mutation_forbidden_in_current_phase")
	}

	evaluationTime := input.EvaluationTime.UTC()
	if evaluationTime.IsZero() {
		evaluationTime = time.Now().UTC()
	}
	decision := WindowsFluidAdmissionDecision{
		DecisionID:                 "windows-fluid-admission-" + evaluationTime.Format("20060102T150405"),
		ShellRef:                   input.Bundle.Shell.Spec.VMRef,
		RequestedAction:            requestedAction,
		AdmissionPhase:             phase,
		MutationAllowed:            false,
		ApplyAllowed:               false,
		RuntimeMode:                "in-place-qmp",
		EvidenceScore:              score.Score,
		EvidenceLevel:              score.EvidenceLevel,
		PolicyVersion:              policy.PolicyVersion,
		Blockers:                   dedupe(blockers),
		DenialReasons:              dedupe(denialReasons),
		RequiredAdditionalEvidence: dedupe(requiredEvidence),
		BlastRadius:                policy.BlastRadiusPolicy,
		RollbackPolicy:             policy.RollbackPolicy,
		ReturnToFloorPolicy:        policy.ReturnToFloorPolicy,
		TTLPolicy:                  policy.TTLPolicy,
		AuditRefs:                  []string{input.Bundle.Shell.Status.EvidenceRef},
		CreatedAt:                  evaluationTime.Format(time.RFC3339),
	}

	nextStep := "collect-missing-evidence"
	switch phase {
	case AdmissionAdmittedForFutureApply:
		nextStep = "prepare-separate-future-apply-governance-review"
	case AdmissionDenied:
		nextStep = "use-certified-single-vm-candidate-only"
	case AdmissionBlocked:
		nextStep = "resolve-hard-blockers-before-next-admission-check"
	case AdmissionQuarantined:
		nextStep = "freeze-candidate-and-rebuild-identity-continuity"
	}

	return AdmissionEvaluationResult{
		Decision:         decision,
		ActionSlate:      dryRun.ActionSlate,
		EvidenceScore:    score,
		PolicyViolations: dedupe(policyViolations),
		Blockers:         dedupe(blockers),
		NextSafeStep:     nextStep,
	}
}

func deriveRequestedAction(intent *DryRunLeaseIntent) RequestedAdmissionAction {
	if intent == nil {
		return RequestedActionEvidenceRefresh
	}
	switch intent.ActionType {
	case string(ActionPrepareCPULease):
		return RequestedActionPrepareCPULease
	case string(ActionPrepareMemoryLease):
		return RequestedActionPrepareMemoryLease
	default:
		return RequestedActionEvidenceRefresh
	}
}

func evaluatePolicyViolations(
	bundle WindowsFluidRuntimeEvidenceBundle,
	dryRun DryRunEvaluationResult,
	policy WindowsFluidPolicyPack,
	requestedAction RequestedAdmissionAction,
) []string {
	violations := make([]string, 0, 8)
	if bundle.Shell.Spec.RuntimeMode != policy.AllowedRuntimeMode {
		violations = append(violations, "runtime_mode_not_allowed")
	}
	if policy.RequireCertifiedFluidShell && dryRun.Classification != ClassificationReadyForFluidShellCertification && dryRun.Phase != StateLeasePrepared {
		violations = append(violations, "shell_not_certified_for_admission")
	}
	if policy.RequireQmpAck && !dryRun.Conditions["qmpReady"] {
		violations = append(violations, "qmp_ack_required")
	}
	if policy.RequireGuestAck && !dryRun.Conditions["guestAckReady"] {
		violations = append(violations, "guest_ack_required")
	}
	if policy.RequireRollbackReady && !dryRun.Conditions["rollbackReady"] {
		violations = append(violations, "rollback_ready_required")
	}
	if policy.RequireReturnToFloorReady && !dryRun.Conditions["returnToFloorReady"] {
		violations = append(violations, "return_to_floor_ready_required")
	}
	if policy.RequireNoLiveMigration && !dryRun.Conditions["noMigrationRequired"] {
		violations = append(violations, "no_live_migration_required")
	}
	if policy.RequireNoRecreate && (bundle.KubeVirtAfter.RecreateRequired || bundle.KubeVirtAfter.RolloutObserved) {
		violations = append(violations, "no_recreate_required")
	}
	if policy.RequireNoReboot && bundle.Guest != nil && bundle.Guest.PendingReboot {
		violations = append(violations, "no_reboot_required")
	}
	if policy.RequireSameNode && !dryRun.Conditions["sameNode"] {
		violations = append(violations, "same_node_required")
	}
	if policy.RequireSameVirtLauncherPod && !dryRun.Conditions["sameVirtLauncherPod"] {
		violations = append(violations, "same_virtlauncher_pod_required")
	}
	if policy.RequireSameQemuProcess && !dryRun.Conditions["sameQemuProcess"] {
		violations = append(violations, "same_qemu_required")
	}
	if policy.RequireSameLastBoot && bundle.Guest != nil && bundle.Guest.LastBootTime == "" {
		violations = append(violations, "same_last_boot_required")
	}
	if policy.RequireSameMachineIdentity && bundle.Guest != nil && bundle.Guest.MachineGUIDHash == "" {
		violations = append(violations, "same_machine_identity_required")
	}
	if requestedAction == RequestedActionPrepareMemoryLease && policy.ReturnToFloorPolicy.RequireMemorySafety {
		if bundle.Guest == nil || !bundle.Guest.MemoryAdapterVerified {
			violations = append(violations, "memory_driver_verification_required")
		}
	}
	return dedupe(violations)
}

func hasPriority(blockers []string, policy WindowsFluidPolicyPack, target BlockerPriority) bool {
	for _, blocker := range blockers {
		if BlockerPriorityForID(blocker, policy) == target {
			return true
		}
	}
	return false
}
