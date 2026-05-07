package windowsfluidvirt

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

type UnlockGateID string

const (
	Gate0ExecutorHardDisabled      UnlockGateID = "GATE_0_EXECUTOR_HARD_DISABLED"
	Gate1LabReadOnlyEvidence       UnlockGateID = "GATE_1_LAB_READONLY_EVIDENCE_COMPLETE"
	Gate2FutureSignableAttestation UnlockGateID = "GATE_2_FUTURE_SIGNABLE_ATTESTATION_REPLAY"
)

type UnlockGateStatus string

const (
	UnlockGatePassed        UnlockGateStatus = "PASSED"
	UnlockGateFailed        UnlockGateStatus = "FAILED"
	UnlockGateBlocked       UnlockGateStatus = "BLOCKED"
	UnlockGateQuarantined   UnlockGateStatus = "QUARANTINED"
	UnlockGateNotApplicable UnlockGateStatus = "NOT_APPLICABLE"
)

type UnlockGateSetAggregateStatus string

const (
	GateSetPassed      UnlockGateSetAggregateStatus = "GATE_SET_PASSED"
	GateSetBlocked     UnlockGateSetAggregateStatus = "GATE_SET_BLOCKED"
	GateSetQuarantined UnlockGateSetAggregateStatus = "GATE_SET_QUARANTINED"
	GateSetFailed      UnlockGateSetAggregateStatus = "GATE_SET_FAILED"
)

const (
	GateBlockerMissingExecutorOutput        = "gate0_executor_output_missing"
	GateBlockerExecutorEnabledFlagTrue      = "gate0_executor_enabled_true"
	GateBlockerApplyAttemptedTrue           = "gate0_apply_attempted_true"
	GateBlockerMutationPerformedTrue        = "gate0_mutation_performed_true"
	GateBlockerQMPCommandSentTrue           = "gate0_qmp_command_sent_true"
	GateBlockerClusterMutationSentTrue      = "gate0_cluster_mutation_sent_true"
	GateBlockerMutationAllowedTrue          = "gate0_mutation_allowed_true"
	GateBlockerApplyAllowedTrue             = "gate0_apply_allowed_true"
	GateBlockerEnvelopeNotEmpty             = "gate0_command_envelope_not_empty"
	GateBlockerGuardMutationWindowOpen      = "gate0_mutation_window_open"
	GateBlockerGuardQMPMutationAllowed      = "gate0_qmp_mutation_allowed"
	GateBlockerGuardClusterMutationAllowed  = "gate0_cluster_mutation_allowed"
	GateBlockerGuardExecutorEnabled         = "gate0_guard_executor_enabled"
	GateBlockerMissingDisabledExecutorProof = "gate0_missing_disabled_executor_blocker"
	GateBlockerKillSwitchMissing            = "gate0_kill_switch_missing"

	GateBlockerEvidenceBundleMissing      = "gate1_evidence_bundle_missing"
	GateBlockerTargetNotMasterWin11       = "gate1_target_not_master_win11"
	GateBlockerPoolContextOnly            = "gate1_pool_context_only"
	GateBlockerIdentityEvidenceIncomplete = "gate1_identity_evidence_incomplete"
	GateBlockerEvidenceStale              = "gate1_evidence_stale"
	GateBlockerGenericWindowsVM           = "gate1_generic_windows_vm_not_certified"

	GateBlockerAttestationMissing       = "gate2_attestation_missing"
	GateBlockerAttestationMalformed     = "gate2_attestation_malformed"
	GateBlockerAttestationStale         = "gate2_attestation_stale"
	GateBlockerAttestationReplayed      = "gate2_attestation_replayed"
	GateBlockerAttestationSubjectDrift  = "gate2_attestation_subject_mismatch"
	GateBlockerAttestationEvidenceDrift = "gate2_attestation_evidence_ref_mismatch"
)

type WindowsFluidUnlockGateVerification struct {
	VerificationID             string           `json:"verificationId"`
	GateID                     UnlockGateID     `json:"gateId"`
	GateStatus                 UnlockGateStatus `json:"gateStatus"`
	ExecutorMustRemainDisabled bool             `json:"executorMustRemainDisabled"`
	MutationAllowed            bool             `json:"mutationAllowed"`
	ApplyAllowed               bool             `json:"applyAllowed"`
	CheckedInputs              []string         `json:"checkedInputs"`
	RequiredInputs             []string         `json:"requiredInputs"`
	MissingInputs              []string         `json:"missingInputs"`
	BlockerList                []string         `json:"blockerList"`
	NegativeTestRefs           []string         `json:"negativeTestRefs"`
	EvidenceRefs               []string         `json:"evidenceRefs"`
	AttestationRefs            []string         `json:"attestationRefs"`
	DeterministicHash          string           `json:"deterministicHash,omitempty"`
	EvaluatedAt                string           `json:"evaluatedAt"`
}

type UnlockGateEvaluationInput struct {
	GateID             UnlockGateID                         `json:"gateId"`
	EvidenceBundle     *WindowsFluidRuntimeEvidenceBundle   `json:"evidenceBundle,omitempty"`
	GovernanceContract *WindowsFluidApplyGovernanceContract `json:"governanceContract,omitempty"`
	ExecutorOutput     *FutureApplyExecutorEvaluationResult `json:"executorOutput,omitempty"`
	Attestation        *WindowsFluidPolicyAttestation       `json:"attestation,omitempty"`
	EvaluationTime     time.Time                            `json:"evaluationTime"`
}

type WindowsFluidUnlockGateSetVerification struct {
	VerificationSetID string                               `json:"verificationSetId"`
	Gates             []WindowsFluidUnlockGateVerification `json:"gates"`
	AggregateStatus   UnlockGateSetAggregateStatus         `json:"aggregateStatus"`
	NextSafeStep      string                               `json:"nextSafeStep"`
	MissingEvidence   []string                             `json:"missingEvidence"`
	Blockers          []string                             `json:"blockers"`
	EvaluatedAt       string                               `json:"evaluatedAt"`
}

type UnlockGateSetEvaluationInput struct {
	EvidenceBundle     *WindowsFluidRuntimeEvidenceBundle   `json:"evidenceBundle,omitempty"`
	GovernanceContract *WindowsFluidApplyGovernanceContract `json:"governanceContract,omitempty"`
	ExecutorOutput     *FutureApplyExecutorEvaluationResult `json:"executorOutput,omitempty"`
	Attestation        *WindowsFluidPolicyAttestation       `json:"attestation,omitempty"`
	EvaluationTime     time.Time                            `json:"evaluationTime"`
}

func EvaluateWindowsFluidUnlockGate(input UnlockGateEvaluationInput) WindowsFluidUnlockGateVerification {
	evaluationTime := input.EvaluationTime.UTC()
	if evaluationTime.IsZero() {
		evaluationTime = time.Now().UTC()
	}

	result := WindowsFluidUnlockGateVerification{
		VerificationID:             "windows-fluid-unlock-gate-" + strings.ToLower(string(input.GateID)) + "-" + evaluationTime.Format("20060102T150405"),
		GateID:                     input.GateID,
		GateStatus:                 UnlockGateBlocked,
		ExecutorMustRemainDisabled: true,
		MutationAllowed:            false,
		ApplyAllowed:               false,
		CheckedInputs:              []string{},
		RequiredInputs:             []string{},
		MissingInputs:              []string{},
		BlockerList:                []string{},
		NegativeTestRefs:           []string{},
		EvidenceRefs:               []string{},
		AttestationRefs:            []string{},
		EvaluatedAt:                evaluationTime.Format(time.RFC3339),
	}

	switch input.GateID {
	case Gate0ExecutorHardDisabled:
		return evaluateGate0(result, input, evaluationTime)
	case Gate1LabReadOnlyEvidence:
		return evaluateGate1(result, input, evaluationTime)
	case Gate2FutureSignableAttestation:
		return evaluateGate2(result, input, evaluationTime)
	default:
		result.GateStatus = UnlockGateNotApplicable
		return result
	}
}

func EvaluateWindowsFluidUnlockGateSet(input UnlockGateSetEvaluationInput) WindowsFluidUnlockGateSetVerification {
	evaluationTime := input.EvaluationTime.UTC()
	if evaluationTime.IsZero() {
		evaluationTime = time.Now().UTC()
	}

	gate0 := EvaluateWindowsFluidUnlockGate(UnlockGateEvaluationInput{
		GateID:             Gate0ExecutorHardDisabled,
		ExecutorOutput:     input.ExecutorOutput,
		GovernanceContract: input.GovernanceContract,
		EvaluationTime:     evaluationTime,
	})
	gate1 := EvaluateWindowsFluidUnlockGate(UnlockGateEvaluationInput{
		GateID:         Gate1LabReadOnlyEvidence,
		EvidenceBundle: input.EvidenceBundle,
		EvaluationTime: evaluationTime,
	})
	gate2 := EvaluateWindowsFluidUnlockGate(UnlockGateEvaluationInput{
		GateID:             Gate2FutureSignableAttestation,
		EvidenceBundle:     input.EvidenceBundle,
		GovernanceContract: input.GovernanceContract,
		Attestation:        input.Attestation,
		EvaluationTime:     evaluationTime,
	})

	gates := []WindowsFluidUnlockGateVerification{gate0, gate1, gate2}
	aggregate := GateSetPassed
	if gate0.GateStatus != UnlockGatePassed {
		aggregate = GateSetBlocked
	}
	for _, gate := range gates {
		if gate.GateStatus == UnlockGateQuarantined {
			aggregate = GateSetQuarantined
			break
		}
	}
	if aggregate != GateSetQuarantined {
		for _, gate := range gates {
			if gate.GateStatus == UnlockGateFailed {
				aggregate = GateSetFailed
				break
			}
			if gate.GateStatus == UnlockGateBlocked {
				aggregate = GateSetBlocked
			}
		}
	}

	missing := make([]string, 0, 8)
	blockers := make([]string, 0, 16)
	for _, gate := range gates {
		missing = append(missing, gate.MissingInputs...)
		blockers = append(blockers, gate.BlockerList...)
	}

	next := "keep_executor_hard_disabled_and_collect_missing_gate_evidence"
	switch aggregate {
	case GateSetPassed:
		next = "document_gate_set_pass_without_unlock_and_prepare_next_non_executable_gate_pack"
	case GateSetQuarantined:
		next = "freeze_candidate_and_rebuild_identity_continuity_proofs"
	case GateSetFailed:
		next = "fix_failed_gate_conditions_and_rerun_deterministic_gate_set_replay"
	}

	return WindowsFluidUnlockGateSetVerification{
		VerificationSetID: "windows-fluid-unlock-gate-set-" + evaluationTime.Format("20060102T150405"),
		Gates:             gates,
		AggregateStatus:   aggregate,
		NextSafeStep:      next,
		MissingEvidence:   dedupe(missing),
		Blockers:          dedupe(blockers),
		EvaluatedAt:       evaluationTime.Format(time.RFC3339),
	}
}

func evaluateGate0(
	result WindowsFluidUnlockGateVerification,
	input UnlockGateEvaluationInput,
	evaluationTime time.Time,
) WindowsFluidUnlockGateVerification {
	result.RequiredInputs = []string{"executor_output", "disabled_executor_proof"}
	if input.GovernanceContract != nil {
		result.RequiredInputs = append(result.RequiredInputs, "governance_contract_flags")
	}
	if input.ExecutorOutput == nil {
		result.MissingInputs = append(result.MissingInputs, "executor_output")
		result.BlockerList = append(result.BlockerList, GateBlockerMissingExecutorOutput)
		result.GateStatus = UnlockGateBlocked
		result.DeterministicHash = gateHash(result, evaluationTime)
		return result
	}

	executor := input.ExecutorOutput
	result.CheckedInputs = append(result.CheckedInputs, "executor_output", "preapply_guard", "command_envelope", "execution_result")
	if len(executor.ExecutionResult.AttestationRefs) > 0 {
		result.AttestationRefs = append(result.AttestationRefs, executor.ExecutionResult.AttestationRefs...)
	}
	if executor.ExecutionResult.ShellRef != "" {
		result.EvidenceRefs = append(result.EvidenceRefs, executor.ExecutionResult.ShellRef)
	}

	if executor.PreApplyGuard.ExecutorEnabled {
		result.BlockerList = append(result.BlockerList, GateBlockerExecutorEnabledFlagTrue, GateBlockerGuardExecutorEnabled)
	}
	if executor.PreApplyGuard.MutationWindowOpen {
		result.BlockerList = append(result.BlockerList, GateBlockerGuardMutationWindowOpen)
	}
	if executor.PreApplyGuard.QMPMutationAllowed {
		result.BlockerList = append(result.BlockerList, GateBlockerGuardQMPMutationAllowed)
	}
	if executor.PreApplyGuard.ClusterMutationAllowed {
		result.BlockerList = append(result.BlockerList, GateBlockerGuardClusterMutationAllowed)
	}
	if executor.ExecutionResult.ApplyAttempted {
		result.BlockerList = append(result.BlockerList, GateBlockerApplyAttemptedTrue)
	}
	if executor.ExecutionResult.MutationPerformed {
		result.BlockerList = append(result.BlockerList, GateBlockerMutationPerformedTrue)
	}
	if executor.ExecutionResult.QMPCommandSent {
		result.BlockerList = append(result.BlockerList, GateBlockerQMPCommandSentTrue)
	}
	if executor.ExecutionResult.ClusterMutationSent {
		result.BlockerList = append(result.BlockerList, GateBlockerClusterMutationSentTrue)
	}
	if executor.KillSwitchSnapshot.KillSwitchID == "" || !executor.KillSwitchSnapshot.Enabled {
		result.BlockerList = append(result.BlockerList, GateBlockerKillSwitchMissing)
	}
	if executor.CommandEnvelope.ContainsExecutableCommand ||
		len(executor.CommandEnvelope.QMPCommands) > 0 ||
		len(executor.CommandEnvelope.ClusterMutations) > 0 ||
		len(executor.CommandEnvelope.GuestMutations) > 0 {
		result.BlockerList = append(result.BlockerList, GateBlockerEnvelopeNotEmpty)
	}
	if !contains(executor.ExecutionResult.Blockers, BlockerFutureApplyExecutorDisabled) {
		result.BlockerList = append(result.BlockerList, GateBlockerMissingDisabledExecutorProof)
	}
	if input.GovernanceContract != nil {
		if input.GovernanceContract.MutationAllowed {
			result.BlockerList = append(result.BlockerList, GateBlockerMutationAllowedTrue)
		}
		if input.GovernanceContract.ApplyAllowed {
			result.BlockerList = append(result.BlockerList, GateBlockerApplyAllowedTrue)
		}
	}

	if len(result.BlockerList) == 0 {
		result.GateStatus = UnlockGatePassed
	} else {
		result.GateStatus = UnlockGateFailed
	}
	result.BlockerList = dedupe(result.BlockerList)
	result.DeterministicHash = gateHash(result, evaluationTime)
	return result
}

func evaluateGate1(
	result WindowsFluidUnlockGateVerification,
	input UnlockGateEvaluationInput,
	evaluationTime time.Time,
) WindowsFluidUnlockGateVerification {
	result.RequiredInputs = []string{
		"evidence_bundle",
		"master_win11_target",
		"identity_evidence",
		"qmp_evidence_readonly",
		"guest_evidence",
		"rollback_return_readiness",
		"fresh_evidence_window",
	}
	if input.EvidenceBundle == nil {
		result.MissingInputs = append(result.MissingInputs, "evidence_bundle")
		result.BlockerList = append(result.BlockerList, GateBlockerEvidenceBundleMissing)
		result.GateStatus = UnlockGateBlocked
		result.DeterministicHash = gateHash(result, evaluationTime)
		return result
	}

	bundle := *input.EvidenceBundle
	result.CheckedInputs = append(result.CheckedInputs, "evidence_bundle")
	result.EvidenceRefs = append(result.EvidenceRefs, bundle.Shell.Status.EvidenceRef, bundle.Shell.Spec.VMRef)

	if bundle.SourceMetadata.SourceName != "master-win11" {
		result.BlockerList = append(result.BlockerList, GateBlockerTargetNotMasterWin11)
	}
	if strings.HasPrefix(bundle.SourceMetadata.SourceName, "win11-pool-") || bundle.PolicyGates.PoolReplicaContextOnly {
		result.BlockerList = append(result.BlockerList, GateBlockerPoolContextOnly, BlockerLiveMigrationRequired)
	}
	if !identityComplete(bundle.KubeVirtBefore) || !identityComplete(bundle.KubeVirtAfter) {
		result.BlockerList = append(result.BlockerList, GateBlockerIdentityEvidenceIncomplete)
	}

	policy := DefaultWindowsFluidPolicyPack()
	dryRun := EvaluateWindowsFluidRuntimeDryRunWithOptions(bundle, DryRunEvaluationOptions{EvaluationTime: evaluationTime})
	result.BlockerList = append(result.BlockerList, dryRun.Blockers...)
	if dryRun.Classification == ClassificationBlockedGenericWindowsVM {
		result.BlockerList = append(result.BlockerList, GateBlockerGenericWindowsVM)
	}
	if bundle.Guest == nil {
		result.BlockerList = append(result.BlockerList, BlockerGuestAckMissing)
	}
	if bundle.QMP == nil {
		result.BlockerList = append(result.BlockerList, BlockerQMPSocketUnavailable)
	}
	if bundle.QMP != nil {
		for _, blocker := range ValidateQmpReadiness(*bundle.QMP) {
			result.BlockerList = append(result.BlockerList, blocker)
		}
	}
	if bundle.Guest != nil {
		for _, blocker := range EvaluateGuestReadiness(*bundle.Guest) {
			result.BlockerList = append(result.BlockerList, blocker)
		}
	}
	if !dryRun.Conditions["rollbackReady"] {
		result.BlockerList = append(result.BlockerList, BlockerRollbackNotReady)
	}
	if !dryRun.Conditions["returnToFloorReady"] {
		result.BlockerList = append(result.BlockerList, BlockerReturnToFloorNotReady)
	}
	if !evidenceFresh(bundle, policy, evaluationTime) {
		result.BlockerList = append(result.BlockerList, GateBlockerEvidenceStale)
	}
	if len(bundle.KubeVirtAfter.LiveMigrationObjectsObserved) > 0 ||
		len(bundle.KubeVirtAfter.VMIMObjectsObserved) > 0 ||
		bundle.KubeVirtAfter.MigrationRequired ||
		bundle.KubeVirtAfter.RecreateRequired ||
		bundle.KubeVirtAfter.RolloutObserved {
		result.BlockerList = append(result.BlockerList, BlockerLiveMigrationRequired)
	}

	result.BlockerList = dedupe(result.BlockerList)
	if result.GateStatus != UnlockGateQuarantined {
		for _, blocker := range result.BlockerList {
			if blocker == BlockerNodeChanged ||
				blocker == BlockerVirtLauncherPodChanged ||
				blocker == BlockerQemuPIDChanged ||
				blocker == BlockerLastBootChanged ||
				blocker == BlockerMachineGUIDChanged ||
				blocker == BlockerCriticalWindowsEventDetected {
				result.GateStatus = UnlockGateQuarantined
				break
			}
		}
	}
	if result.GateStatus == UnlockGateQuarantined {
		result.DeterministicHash = gateHash(result, evaluationTime)
		return result
	}
	if len(result.BlockerList) == 0 {
		result.GateStatus = UnlockGatePassed
	} else {
		result.GateStatus = UnlockGateBlocked
	}
	result.DeterministicHash = gateHash(result, evaluationTime)
	return result
}

func evaluateGate2(
	result WindowsFluidUnlockGateVerification,
	input UnlockGateEvaluationInput,
	evaluationTime time.Time,
) WindowsFluidUnlockGateVerification {
	result.RequiredInputs = []string{"attestation", "subject_ref_coherence", "evidence_ref_coherence", "fresh_attestation_window"}
	if input.Attestation == nil {
		result.MissingInputs = append(result.MissingInputs, "attestation")
		result.BlockerList = append(result.BlockerList, GateBlockerAttestationMissing)
		result.GateStatus = UnlockGateBlocked
		result.DeterministicHash = gateHash(result, evaluationTime)
		return result
	}

	att := *input.Attestation
	result.CheckedInputs = append(result.CheckedInputs, "attestation")
	result.AttestationRefs = append(result.AttestationRefs, att.AttestationID)
	result.EvidenceRefs = append(result.EvidenceRefs, att.EvidenceRefs...)

	if att.AttestationID == "" || att.SubjectRef == "" || att.SubjectType == "" || att.PolicyVersion == "" {
		result.BlockerList = append(result.BlockerList, GateBlockerAttestationMalformed)
	}
	if att.Signature.Mode != "future-signable" && att.Signature.Mode != "unsigned-dev" {
		result.BlockerList = append(result.BlockerList, GateBlockerAttestationMalformed)
	}
	if att.Signature.Value != "" {
		result.BlockerList = append(result.BlockerList, GateBlockerAttestationMalformed)
	}

	createdAt, err := time.Parse(time.RFC3339, att.CreatedAt)
	if err != nil || createdAt.After(evaluationTime) || int64(evaluationTime.Sub(createdAt).Seconds()) > DefaultWindowsFluidPolicyPack().MaxEvidenceAgeSeconds {
		result.BlockerList = append(result.BlockerList, GateBlockerAttestationStale)
	}

	if input.GovernanceContract != nil && att.SubjectRef != input.GovernanceContract.ContractID {
		result.BlockerList = append(result.BlockerList, GateBlockerAttestationSubjectDrift)
	}
	if input.EvidenceBundle != nil {
		shellRef := input.EvidenceBundle.Shell.Spec.VMRef
		sourceName := input.EvidenceBundle.SourceMetadata.SourceName
		if !contains(att.EvidenceRefs, shellRef) && !contains(att.EvidenceRefs, sourceName) {
			result.BlockerList = append(result.BlockerList, GateBlockerAttestationEvidenceDrift)
		}
	}
	if replayFlag(att.DecisionSnapshot, "replayDetected") || replayFlag(att.DecisionSnapshot, "replayedEvidence") {
		result.BlockerList = append(result.BlockerList, GateBlockerAttestationReplayed)
	}

	result.BlockerList = dedupe(result.BlockerList)
	if len(result.BlockerList) == 0 {
		result.GateStatus = UnlockGatePassed
	} else {
		result.GateStatus = UnlockGateBlocked
	}
	result.DeterministicHash = gateHash(result, evaluationTime)
	return result
}

func replayFlag(snapshot map[string]any, key string) bool {
	if snapshot == nil {
		return false
	}
	raw, ok := snapshot[key]
	if !ok {
		return false
	}
	value, ok := raw.(bool)
	return ok && value
}

func gateHash(result WindowsFluidUnlockGateVerification, evaluationTime time.Time) string {
	payload := map[string]any{
		"gateId":           result.GateID,
		"gateStatus":       result.GateStatus,
		"checkedInputs":    dedupe(result.CheckedInputs),
		"requiredInputs":   dedupe(result.RequiredInputs),
		"missingInputs":    dedupe(result.MissingInputs),
		"blockerList":      dedupe(result.BlockerList),
		"evidenceRefs":     dedupe(result.EvidenceRefs),
		"attestationRefs":  dedupe(result.AttestationRefs),
		"evaluatedAt":      evaluationTime.Format(time.RFC3339),
		"executorDisabled": result.ExecutorMustRemainDisabled,
		"mutationAllowed":  result.MutationAllowed,
		"applyAllowed":     result.ApplyAllowed,
	}
	raw, _ := json.Marshal(payload)
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
