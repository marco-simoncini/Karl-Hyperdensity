package windowsfluidvirt

import (
	"strings"
	"time"
)

type DryRunEvaluationResult struct {
	Phase                   WindowsFluidPhase           `json:"phase"`
	Classification          CertificationClassification `json:"classification"`
	Conditions              map[string]bool             `json:"conditions"`
	Blockers                []string                    `json:"blockers"`
	EvidenceSummary         map[string]any              `json:"evidenceSummary"`
	ActionSlate             WindowsFluidActionSlate     `json:"actionSlate"`
	RecommendedNextSafeStep string                      `json:"recommendedNextSafeStep"`
}

type DryRunEvaluationOptions struct {
	EvaluationTime time.Time
}

func EvaluateWindowsFluidRuntimeDryRun(bundle WindowsFluidRuntimeEvidenceBundle) DryRunEvaluationResult {
	return EvaluateWindowsFluidRuntimeDryRunWithOptions(bundle, DryRunEvaluationOptions{})
}

func EvaluateWindowsFluidRuntimeDryRunWithOptions(bundle WindowsFluidRuntimeEvidenceBundle, options DryRunEvaluationOptions) DryRunEvaluationResult {
	blockers := dedupe(bundle.ObservedBlockers)
	beforeIdentity := bundle.KubeVirtBefore
	afterIdentity := bundle.KubeVirtAfter

	var qmpEvidence QMPEvidence
	if bundle.QMP == nil {
		blockers = append(blockers, BlockerQMPSocketUnavailable)
		qmpEvidence = NewReadOnlyQMPEvidence("missing", "")
		qmpEvidence.QMPConnected = false
	} else {
		qmpEvidence = *bundle.QMP
	}

	var guestEvidence GuestRuntimeEvidence
	if bundle.Guest == nil {
		blockers = append(blockers, BlockerGuestAckMissing)
		guestEvidence = GuestRuntimeEvidence{
			GuestAck:              false,
			LastBootTime:          "missing-guest-evidence",
			MachineGUIDHash:       "missing-guest-evidence",
			MemoryAdapterVerified: true,
			ReturnToFloorReady:    true,
		}
	} else {
		guestEvidence = *bundle.Guest
	}

	if identityEvidenceMissing(beforeIdentity) || identityEvidenceMissing(afterIdentity) {
		blockers = append(blockers, BlockerVMIRecreateRequired)
		beforeIdentity = fillIdentityEvidence(beforeIdentity, bundle.SourceMetadata)
		afterIdentity = fillIdentityEvidence(afterIdentity, bundle.SourceMetadata)
	}

	runtimeGate := EvaluateFluidRuntimeGate(RuntimeGateInput{
		Annotations:    bundle.PolicyGates.Annotations,
		Shell:          bundle.Shell,
		BeforeIdentity: beforeIdentity,
		AfterIdentity:  afterIdentity,
		QMP:            qmpEvidence,
		Guest:          guestEvidence,
	})

	blockers = dedupe(append(blockers, runtimeGate.Blockers...))
	phase := runtimeGate.Phase
	classification := runtimeGate.Classification

	if bundle.PolicyGates.PoolReplicaContextOnly || isPoolReplica(bundle.SourceMetadata.SourceName) {
		blockers = append(blockers, BlockerLiveMigrationRequired)
		phase = StateBlocked
		classification = ClassificationBlockedPoolReplicaModel
	}
	if len(bundle.PolicyGates.Annotations) == 0 {
		phase = StateBlocked
		classification = ClassificationBlockedGenericWindowsVM
	}

	if bundle.LeaseIntent != nil {
		if phase == StateReady &&
			runtimeGate.Conditions["qmpReady"] &&
			runtimeGate.Conditions["guestAckReady"] &&
			bundle.LeaseIntent.RollbackReady &&
			bundle.LeaseIntent.ReturnToFloorReady {
			phase = StateLeasePrepared
		} else {
			phase = StateBlocked
			if !bundle.LeaseIntent.RollbackReady {
				blockers = append(blockers, BlockerRollbackNotReady)
			}
			if !bundle.LeaseIntent.ReturnToFloorReady {
				blockers = append(blockers, BlockerReturnToFloorNotReady)
				if bundle.LeaseIntent.ActionType == string(ActionPrepareMemoryLease) {
					blockers = append(blockers, BlockerMemoryReturnNotSafe)
				}
			}
		}
	}

	if hasQuarantineBlocker(blockers) {
		phase = StateQuarantined
		if classification == ClassificationSupportedCandidate || classification == ClassificationReadyForFluidShellCertification {
			classification = ClassificationQuarantinedIdentityChanged
		}
	}

	conditions := runtimeGate.Conditions
	evaluationTime := options.EvaluationTime.UTC()
	if evaluationTime.IsZero() {
		evaluationTime = time.Now().UTC()
	}
	actionSlate := BuildDryRunActionSlate(
		"windows-fluid-dryrun-"+evaluationTime.Format("20060102T150405"),
		bundle.Shell,
		phase,
		blockers,
		conditions,
		bundle.LeaseIntent,
	)

	return DryRunEvaluationResult{
		Phase:                   phase,
		Classification:          classification,
		Conditions:              conditions,
		Blockers:                dedupe(blockers),
		EvidenceSummary:         runtimeGate.EvidenceSummary,
		ActionSlate:             actionSlate,
		RecommendedNextSafeStep: recommendedStep(phase),
	}
}

func hasQuarantineBlocker(blockers []string) bool {
	for _, blocker := range blockers {
		switch blocker {
		case BlockerNodeChanged, BlockerVirtLauncherPodChanged, BlockerQemuPIDChanged, BlockerLastBootChanged, BlockerMachineGUIDChanged:
			return true
		}
	}
	return false
}

func recommendedStep(phase WindowsFluidPhase) string {
	switch phase {
	case StateReady:
		return "review-certification-gate-and-prepare-non-mutating-lease"
	case StateLeasePrepared:
		return "keep-dry-run-state-and-wait-for-explicit-future-apply-phase"
	case StateQuarantined:
		return "quarantine-runtime-and-rebuild-continuity-evidence"
	default:
		return "resolve-blockers-and-refresh-evidence"
	}
}

func isPoolReplica(vmName string) bool {
	return strings.HasPrefix(vmName, "win11-pool-")
}

func identityEvidenceMissing(identity KubeVirtRuntimeIdentityEvidence) bool {
	return identity.VMIUID == "" || identity.VirtLauncherPodUID == "" || identity.NodeName == "" || identity.QemuPID == ""
}

func fillIdentityEvidence(identity KubeVirtRuntimeIdentityEvidence, source RuntimeSourceMetadata) KubeVirtRuntimeIdentityEvidence {
	if identity.VMName == "" {
		identity.VMName = source.SourceName
	}
	if identity.VMNamespace == "" {
		identity.VMNamespace = source.SourceNamespace
	}
	if identity.VMIName == "" {
		identity.VMIName = identity.VMName
	}
	if identity.VMIUID == "" {
		identity.VMIUID = identity.VMName + "-missing-vmi"
	}
	if identity.VirtLauncherPodUID == "" {
		identity.VirtLauncherPodUID = identity.VMName + "-missing-pod"
	}
	if identity.NodeName == "" {
		identity.NodeName = "unknown-node"
	}
	if identity.QemuPID == "" {
		identity.QemuPID = "unknown-qemu-pid"
	}
	return identity
}
