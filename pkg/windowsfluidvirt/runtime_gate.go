package windowsfluidvirt

const (
	AnnotationFluidRuntime         = "hyperdensity.karl.io/fluid-runtime"
	AnnotationNoLiveMigration      = "hyperdensity.karl.io/no-live-migration"
	AnnotationNoReboot             = "hyperdensity.karl.io/no-reboot"
	AnnotationNoRecreate           = "hyperdensity.karl.io/no-recreate"
	AnnotationRuntimeMode          = "hyperdensity.karl.io/runtime-mode"
	AnnotationSingleNodeCompatible = "hyperdensity.karl.io/single-node-compatible"
)

var RequiredRuntimeAnnotations = map[string]string{
	AnnotationFluidRuntime:         "true",
	AnnotationNoLiveMigration:      "true",
	AnnotationNoReboot:             "required",
	AnnotationNoRecreate:           "required",
	AnnotationRuntimeMode:          "in-place-qmp",
	AnnotationSingleNodeCompatible: "true",
}

type CertificationClassification string

const (
	ClassificationSupportedCandidate              CertificationClassification = "SUPPORTED_CANDIDATE"
	ClassificationReadyForFluidShellCertification CertificationClassification = "READY_FOR_FLUID_SHELL_CERTIFICATION"
	ClassificationBlockedGenericWindowsVM         CertificationClassification = "BLOCKED_GENERIC_WINDOWS_VM"
	ClassificationBlockedPoolReplicaModel         CertificationClassification = "BLOCKED_POOL_REPLICA_MODEL"
	ClassificationBlockedMissingQMP               CertificationClassification = "BLOCKED_MISSING_QMP"
	ClassificationBlockedMissingGuestAck          CertificationClassification = "BLOCKED_MISSING_GUEST_ACK"
	ClassificationBlockedLiveMigrationRequired    CertificationClassification = "BLOCKED_LIVE_MIGRATION_REQUIRED"
	ClassificationQuarantinedIdentityChanged      CertificationClassification = "QUARANTINED_IDENTITY_CHANGED"
)

type GuestRuntimeEvidence struct {
	GuestAck               bool
	PendingReboot          bool
	LastBootTime           string
	MachineGUIDHash        string
	MemoryAdapterVerified  bool
	ReturnToFloorReady     bool
	CriticalEventsDetected bool
}

type RuntimeGateInput struct {
	Annotations    map[string]string
	Shell          WindowsFluidShell
	BeforeIdentity KubeVirtRuntimeIdentityEvidence
	AfterIdentity  KubeVirtRuntimeIdentityEvidence
	QMP            QMPEvidence
	Guest          GuestRuntimeEvidence
}

type RuntimeGateEvaluation struct {
	Phase           WindowsFluidPhase
	Conditions      map[string]bool
	Blockers        []string
	EvidenceSummary map[string]any
	Classification  CertificationClassification
}

func EvaluateFluidRuntimeGate(input RuntimeGateInput) RuntimeGateEvaluation {
	conditions := map[string]bool{
		"fluidRuntimeReady":   true,
		"qmpReady":            true,
		"guestAckReady":       true,
		"noMigrationRequired": true,
		"noRebootProof":       true,
		"sameQemuProcess":     true,
		"sameNode":            true,
		"sameVirtLauncherPod": true,
		"returnToFloorReady":  true,
		"rollbackReady":       true,
	}

	blockers := ValidateRuntimeAnnotations(input.Annotations)
	blockers = append(blockers, EvaluateNoMigrationProof(input.AfterIdentity)...)
	blockers = append(blockers, EvaluateNoRecreateProof(input.AfterIdentity)...)
	blockers = append(blockers, EvaluateKubeVirtIdentityContinuity(input.BeforeIdentity, input.AfterIdentity)...)
	blockers = append(blockers, ValidateQmpReadiness(input.QMP)...)
	blockers = append(blockers, EvaluateGuestReadiness(input.Guest)...)
	blockers = append(blockers, ValidateWindowsFluidShell(input.Shell)...)

	for _, blocker := range dedupe(blockers) {
		switch blocker {
		case BlockerQMPSocketUnavailable, BlockerQMPAckMissing:
			conditions["qmpReady"] = false
		case BlockerGuestAckMissing, BlockerPendingRebootDetected, BlockerCriticalWindowsEventDetected:
			conditions["guestAckReady"] = false
		case BlockerLiveMigrationRequired:
			conditions["noMigrationRequired"] = false
		case BlockerLastBootChanged:
			conditions["noRebootProof"] = false
		case BlockerQemuPIDChanged:
			conditions["sameQemuProcess"] = false
		case BlockerNodeChanged:
			conditions["sameNode"] = false
		case BlockerVirtLauncherPodChanged:
			conditions["sameVirtLauncherPod"] = false
		case BlockerReturnToFloorNotReady, BlockerMemoryReturnNotSafe:
			conditions["returnToFloorReady"] = false
		case BlockerRollbackNotReady:
			conditions["rollbackReady"] = false
		}
	}
	if len(blockers) > 0 {
		conditions["fluidRuntimeReady"] = false
	}

	phase := StateReady
	if len(blockers) > 0 {
		phase = terminalPhaseForBlockers(blockers)
	}
	classification := EvaluateFluidShellCertificationReadiness(phase, blockers, input.AfterIdentity)

	return RuntimeGateEvaluation{
		Phase:      phase,
		Conditions: conditions,
		Blockers:   dedupe(blockers),
		EvidenceSummary: map[string]any{
			"kubevirtIdentity": RuntimeIdentitySummary(input.AfterIdentity),
			"qmp":              input.QMP,
			"guestAck":         input.Guest.GuestAck,
		},
		Classification: classification,
	}
}

func EvaluateGuestReadiness(guest GuestRuntimeEvidence) []string {
	var blockers []string
	if !guest.GuestAck {
		blockers = append(blockers, BlockerGuestAckMissing)
	}
	if guest.PendingReboot {
		blockers = append(blockers, BlockerPendingRebootDetected)
	}
	if guest.LastBootTime == "" {
		blockers = append(blockers, BlockerLastBootChanged)
	}
	if guest.MachineGUIDHash == "" {
		blockers = append(blockers, BlockerMachineGUIDChanged)
	}
	if !guest.MemoryAdapterVerified {
		blockers = append(blockers, BlockerMemoryDriverUnverified)
	}
	if !guest.ReturnToFloorReady {
		blockers = append(blockers, BlockerReturnToFloorNotReady)
	}
	if guest.CriticalEventsDetected {
		blockers = append(blockers, BlockerCriticalWindowsEventDetected)
	}
	return dedupe(blockers)
}

func EvaluateFluidShellCertificationReadiness(phase WindowsFluidPhase, blockers []string, identity KubeVirtRuntimeIdentityEvidence) CertificationClassification {
	if identity.VMName == "" || identity.VMIName == "" {
		return ClassificationBlockedGenericWindowsVM
	}
	if len(identity.LiveMigrationObjectsObserved) > 0 || len(identity.VMIMObjectsObserved) > 0 || identity.MigrationRequired {
		return ClassificationBlockedLiveMigrationRequired
	}
	for _, blocker := range blockers {
		switch blocker {
		case BlockerQMPSocketUnavailable, BlockerQMPAckMissing:
			return ClassificationBlockedMissingQMP
		case BlockerGuestAckMissing:
			return ClassificationBlockedMissingGuestAck
		case BlockerNodeChanged, BlockerVirtLauncherPodChanged, BlockerQemuPIDChanged, BlockerMachineGUIDChanged, BlockerLastBootChanged:
			return ClassificationQuarantinedIdentityChanged
		}
	}
	if phase == StateReady {
		return ClassificationReadyForFluidShellCertification
	}
	if phase == StateBlocked || phase == StateQuarantined {
		return ClassificationSupportedCandidate
	}
	return ClassificationBlockedPoolReplicaModel
}

func ValidateRuntimeAnnotations(annotations map[string]string) []string {
	var blockers []string
	for key, requiredValue := range RequiredRuntimeAnnotations {
		if annotations[key] != requiredValue {
			switch key {
			case AnnotationFluidRuntime:
				blockers = append(blockers, BlockerKarlAgentFluidModuleMissing)
			case AnnotationNoLiveMigration:
				blockers = append(blockers, BlockerLiveMigrationRequired)
			case AnnotationNoReboot:
				blockers = append(blockers, BlockerPendingRebootDetected)
			case AnnotationNoRecreate:
				blockers = append(blockers, BlockerVMIRecreateRequired)
			case AnnotationRuntimeMode:
				blockers = append(blockers, BlockerQMPSocketUnavailable)
			case AnnotationSingleNodeCompatible:
				blockers = append(blockers, BlockerGuestAckMissing)
			default:
				blockers = append(blockers, BlockerGuestAckMissing)
			}
		}
	}
	return dedupe(blockers)
}
