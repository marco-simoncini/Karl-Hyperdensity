package windowsfluidvirt

import "time"

type AdmissionEvidenceScore struct {
	Score           int64         `json:"score"`
	EvidenceLevel   EvidenceLevel `json:"evidenceLevel"`
	MissingEvidence []string      `json:"missingEvidence"`
	HardBlockers    []string      `json:"hardBlockers"`
	SoftUnknowns    []string      `json:"softUnknowns"`
}

func EvaluateAdmissionEvidenceScore(
	bundle WindowsFluidRuntimeEvidenceBundle,
	dryRun DryRunEvaluationResult,
	policy WindowsFluidPolicyPack,
	evaluationTime time.Time,
) AdmissionEvidenceScore {
	score := int64(0)
	missingEvidence := make([]string, 0, 8)
	softUnknowns := make([]string, 0, 8)
	hardBlockers := make([]string, 0, 8)

	count := func(ok bool, reason string) {
		if ok {
			score += 5
			return
		}
		missingEvidence = append(missingEvidence, reason)
	}

	count(annotationsComplete(bundle.PolicyGates.Annotations), "required_annotations_missing")
	count(len(ValidateWindowsFluidShell(bundle.Shell)) == 0, "shell_contract_invalid")
	count(identityComplete(bundle.KubeVirtAfter), "identity_incomplete")
	count(bundle.QMP != nil && bundle.QMP.QMPConnected, "qmp_not_connected")
	count(bundle.QMP != nil && bundle.QMP.QMPCapabilitiesNegotiated, "qmp_ack_missing")
	count(bundle.QMP != nil && bundle.QMP.QMPReadOnly, "qmp_readonly_not_enforced")
	count(bundle.Guest != nil && bundle.Guest.GuestAck, "guest_ack_missing")
	count(bundle.Guest != nil && !bundle.Guest.PendingReboot, "pending_reboot_detected")
	count(bundle.Guest != nil && bundle.Guest.LastBootTime != "", "last_boot_proof_missing")
	count(bundle.Guest != nil && bundle.Guest.MachineGUIDHash != "", "machine_identity_proof_missing")
	count(dryRun.Conditions["sameNode"], "same_node_proof_missing")
	count(dryRun.Conditions["sameVirtLauncherPod"], "same_pod_proof_missing")
	count(dryRun.Conditions["sameQemuProcess"], "same_qemu_proof_missing")
	count(dryRun.Conditions["noMigrationRequired"], "no_migration_proof_missing")
	count(dryRun.Conditions["rollbackReady"], "rollback_ready_missing")
	count(dryRun.Conditions["returnToFloorReady"], "return_to_floor_ready_missing")
	count(bundle.Guest != nil && bundle.Guest.MemoryAdapterVerified, "memory_driver_unverified")
	count(bundle.Guest != nil && !bundle.Guest.CriticalEventsDetected, "critical_windows_event_detected")
	count(evidenceFresh(bundle, policy, evaluationTime), "evidence_stale")
	count(!bundle.KubeVirtAfter.RecreateRequired && !bundle.KubeVirtAfter.RolloutObserved, "recreate_or_rollout_detected")

	for _, blocker := range dryRun.Blockers {
		switch BlockerPriorityForID(blocker, policy) {
		case BlockerPriorityP0Quarantine, BlockerPriorityP1HardBlock:
			hardBlockers = append(hardBlockers, blocker)
		case BlockerPriorityP2Capability:
			softUnknowns = append(softUnknowns, blocker)
		}
	}

	level := EvidenceLevelInsufficient
	switch {
	case score >= 90 && len(hardBlockers) == 0:
		level = EvidenceLevelFutureApplyAdmissible
	case score >= 75:
		level = EvidenceLevelDryRunReady
	case score >= 50:
		level = EvidenceLevelPartial
	default:
		level = EvidenceLevelInsufficient
	}

	return AdmissionEvidenceScore{
		Score:           score,
		EvidenceLevel:   level,
		MissingEvidence: dedupe(missingEvidence),
		HardBlockers:    dedupe(hardBlockers),
		SoftUnknowns:    dedupe(softUnknowns),
	}
}

func annotationsComplete(annotations map[string]string) bool {
	if len(annotations) == 0 {
		return false
	}
	for k, v := range RequiredRuntimeAnnotations {
		if annotations[k] != v {
			return false
		}
	}
	return true
}

func identityComplete(identity KubeVirtRuntimeIdentityEvidence) bool {
	return identity.VMName != "" &&
		identity.VMNamespace != "" &&
		identity.VMIName != "" &&
		identity.VMIUID != "" &&
		identity.NodeName != "" &&
		identity.VirtLauncherPodUID != "" &&
		identity.QemuPID != ""
}

func evidenceFresh(bundle WindowsFluidRuntimeEvidenceBundle, policy WindowsFluidPolicyPack, evaluationTime time.Time) bool {
	collectedAt, ok := bundle.Timestamps["collectedAt"]
	if !ok || collectedAt == "" {
		return false
	}
	parsed, err := time.Parse(time.RFC3339, collectedAt)
	if err != nil {
		return false
	}
	et := evaluationTime.UTC()
	if et.IsZero() {
		et = time.Now().UTC()
	}
	if parsed.After(et) {
		return false
	}
	age := et.Sub(parsed)
	return int64(age.Seconds()) <= policy.MaxEvidenceAgeSeconds
}
