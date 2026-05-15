package windowsfluidvirt

import "time"

type ExecutionPhase string

const (
	ExecutionHardDisabled ExecutionPhase = "EXECUTION_HARD_DISABLED"
	ExecutionDenied       ExecutionPhase = "EXECUTION_DENIED"
	ExecutionBlocked      ExecutionPhase = "EXECUTION_BLOCKED"
	ExecutionQuarantined  ExecutionPhase = "EXECUTION_QUARANTINED"
)

type PreApplyGuardPhase string

const (
	GuardReadyButExecutorDisabled PreApplyGuardPhase = "GUARD_READY_BUT_EXECUTOR_DISABLED"
	GuardBlocked                  PreApplyGuardPhase = "GUARD_BLOCKED"
	GuardQuarantined              PreApplyGuardPhase = "GUARD_QUARANTINED"
	GuardNeedsRevalidation        PreApplyGuardPhase = "GUARD_NEEDS_REVALIDATION"
)

type WindowsFluidFutureApplyExecutor interface {
	EvaluatePreApplyGuard(
		governanceContract WindowsFluidApplyGovernanceContract,
		revalidation WindowsFluidPreApplyRevalidationContract,
		attestation WindowsFluidPolicyAttestation,
		killSwitch WindowsFluidKillSwitch,
		evaluationTime time.Time,
	) WindowsFluidPreApplyGuard
	BuildCommandEnvelope(
		governanceContract WindowsFluidApplyGovernanceContract,
		guard WindowsFluidPreApplyGuard,
		evaluationTime time.Time,
	) WindowsFluidExecutorCommandEnvelope
	DenyExecution(
		governanceContract WindowsFluidApplyGovernanceContract,
		guard WindowsFluidPreApplyGuard,
		envelope WindowsFluidExecutorCommandEnvelope,
		attestation WindowsFluidPolicyAttestation,
		evaluationTime time.Time,
	) DisabledFutureApplyExecutionResult
	EvaluateKillSwitch(
		governanceContract WindowsFluidApplyGovernanceContract,
		evaluationTime time.Time,
	) WindowsFluidKillSwitch
	EmitExecutionDeniedEvidence(
		result DisabledFutureApplyExecutionResult,
		guard WindowsFluidPreApplyGuard,
		killSwitch WindowsFluidKillSwitch,
		envelope WindowsFluidExecutorCommandEnvelope,
		evaluationTime time.Time,
	) map[string]any
}

type WindowsFluidPreApplyGuard struct {
	GuardID                      string             `json:"guardId"`
	GovernanceContractRef        string             `json:"governanceContractRef"`
	RevalidationRef              string             `json:"revalidationRef"`
	AttestationRef               string             `json:"attestationRef"`
	KillSwitchReady              bool               `json:"killSwitchReady"`
	ExecutorEnabled              bool               `json:"executorEnabled"`
	MutationWindowOpen           bool               `json:"mutationWindowOpen"`
	QMPMutationAllowed           bool               `json:"qmpMutationAllowed"`
	ClusterMutationAllowed       bool               `json:"clusterMutationAllowed"`
	RequiredFreshEvidencePresent bool               `json:"requiredFreshEvidencePresent"`
	IdentityContinuityPassed     bool               `json:"identityContinuityPassed"`
	RollbackReady                bool               `json:"rollbackReady"`
	ReturnToFloorReady           bool               `json:"returnToFloorReady"`
	Blockers                     []string           `json:"blockers"`
	GuardPhase                   PreApplyGuardPhase `json:"guardPhase"`
}

type WindowsFluidKillSwitch struct {
	KillSwitchID     string `json:"killSwitchId"`
	Enabled          bool   `json:"enabled"`
	Source           string `json:"source"`
	Mode             string `json:"mode"`
	Reason           string `json:"reason"`
	RequiredForApply bool   `json:"requiredForApply"`
	ObservedAt       string `json:"observedAt"`
}

type WindowsFluidExecutorCommandEnvelope struct {
	EnvelopeID                string                    `json:"envelopeId"`
	ShellRef                  string                    `json:"shellRef"`
	RequestedAction           GovernanceRequestedAction `json:"requestedAction"`
	CommandClass              string                    `json:"commandClass"`
	RuntimeMode               string                    `json:"runtimeMode"`
	CommandPreviewOnly        bool                      `json:"commandPreviewOnly"`
	ContainsExecutableCommand bool                      `json:"containsExecutableCommand"`
	QMPCommands               []string                  `json:"qmpCommands"`
	ClusterMutations          []string                  `json:"clusterMutations"`
	GuestMutations            []string                  `json:"guestMutations"`
	RequiredEvidenceRefs      []string                  `json:"requiredEvidenceRefs"`
	RequiredAttestationRefs   []string                  `json:"requiredAttestationRefs"`
	DeniedReason              string                    `json:"deniedReason"`
	CreatedAt                 string                    `json:"createdAt"`
}

type DisabledFutureApplyExecutionResult struct {
	ExecutorID            string                    `json:"executorId"`
	ShellRef              string                    `json:"shellRef"`
	RequestedAction       GovernanceRequestedAction `json:"requestedAction"`
	ExecutionPhase        ExecutionPhase            `json:"executionPhase"`
	ApplyAttempted        bool                      `json:"applyAttempted"`
	MutationPerformed     bool                      `json:"mutationPerformed"`
	QMPCommandSent        bool                      `json:"qmpCommandSent"`
	ClusterMutationSent   bool                      `json:"clusterMutationSent"`
	Reason                string                    `json:"reason"`
	Blockers              []string                  `json:"blockers"`
	RequiredFutureUnlocks []string                  `json:"requiredFutureUnlocks"`
	AttestationRefs       []string                  `json:"attestationRefs"`
	CreatedAt             string                    `json:"createdAt"`
}

type FutureApplyExecutorEvaluationInput struct {
	GovernanceContract WindowsFluidApplyGovernanceContract      `json:"governanceContract"`
	Revalidation       WindowsFluidPreApplyRevalidationContract `json:"revalidation"`
	Attestation        WindowsFluidPolicyAttestation            `json:"attestation"`
	KillSwitch         *WindowsFluidKillSwitch                  `json:"killSwitch"`
	EvaluationTime     time.Time                                `json:"evaluationTime"`
}

type FutureApplyExecutorEvaluationResult struct {
	ExecutionResult    DisabledFutureApplyExecutionResult  `json:"executionResult"`
	PreApplyGuard      WindowsFluidPreApplyGuard           `json:"preApplyGuard"`
	KillSwitchSnapshot WindowsFluidKillSwitch              `json:"killSwitchSnapshot"`
	CommandEnvelope    WindowsFluidExecutorCommandEnvelope `json:"commandEnvelope"`
	Blockers           []string                            `json:"blockers"`
	NextSafeStep       string                              `json:"nextSafeStep"`
	DeniedEvidence     map[string]any                      `json:"deniedEvidence"`
}

type DisabledWindowsFluidApplyExecutor struct {
	ExecutorID string
}

func NewDisabledWindowsFluidApplyExecutor() DisabledWindowsFluidApplyExecutor {
	return DisabledWindowsFluidApplyExecutor{ExecutorID: "karl-fluid-executor-disabled-v1"}
}

func (d DisabledWindowsFluidApplyExecutor) EvaluateKillSwitch(
	governanceContract WindowsFluidApplyGovernanceContract,
	evaluationTime time.Time,
) WindowsFluidKillSwitch {
	ts := evaluationTime.UTC()
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	return WindowsFluidKillSwitch{
		KillSwitchID:     "windows-fluid-killswitch-" + ts.Format("20060102T150405"),
		Enabled:          true,
		Source:           "default-safe",
		Mode:             "hard-disabled",
		Reason:           "future apply executor disabled by policy",
		RequiredForApply: true,
		ObservedAt:       ts.Format(time.RFC3339),
	}
}

func (d DisabledWindowsFluidApplyExecutor) EvaluatePreApplyGuard(
	governanceContract WindowsFluidApplyGovernanceContract,
	revalidation WindowsFluidPreApplyRevalidationContract,
	attestation WindowsFluidPolicyAttestation,
	killSwitch WindowsFluidKillSwitch,
	evaluationTime time.Time,
) WindowsFluidPreApplyGuard {
	ts := evaluationTime.UTC()
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	blockers := dedupe(append([]string{}, governanceContract.Blockers...))
	if !killSwitch.Enabled || !killSwitch.RequiredForApply {
		blockers = append(blockers, BlockerFutureApplyExecutorDisabled)
	}
	if killSwitch.KillSwitchID == "" {
		blockers = append(blockers, BlockerFutureApplyExecutorDisabled)
	}
	if attestation.Signature.Mode != "unsigned-dev" && attestation.Signature.Mode != "future-signable" {
		blockers = append(blockers, BlockerFutureApplyExecutorDisabled)
	}
	if attestation.Signature.Value != "" {
		blockers = append(blockers, BlockerFutureApplyExecutorDisabled)
	}
	phase := GuardReadyButExecutorDisabled
	if governanceContract.GovernancePhase == GovernanceContractQuarantined || revalidation.OutputAllowed == RevalidationQuarantined {
		phase = GuardQuarantined
	}
	if revalidation.OutputAllowed == RevalidationStale {
		phase = GuardNeedsRevalidation
	}
	if governanceContract.GovernancePhase == GovernanceContractBlocked || revalidation.OutputAllowed == RevalidationBlocked || len(blockers) > 0 {
		if phase != GuardQuarantined && phase != GuardNeedsRevalidation {
			phase = GuardBlocked
		}
	}
	identityPassed := phase != GuardQuarantined
	return WindowsFluidPreApplyGuard{
		GuardID:                      "windows-fluid-preapply-guard-" + ts.Format("20060102T150405"),
		GovernanceContractRef:        governanceContract.ContractID,
		RevalidationRef:              revalidation.RevalidationID,
		AttestationRef:               attestation.AttestationID,
		KillSwitchReady:              killSwitch.Enabled && killSwitch.RequiredForApply,
		ExecutorEnabled:              false,
		MutationWindowOpen:           false,
		QMPMutationAllowed:           false,
		ClusterMutationAllowed:       false,
		RequiredFreshEvidencePresent: revalidation.OutputAllowed == RevalidationReady,
		IdentityContinuityPassed:     identityPassed,
		RollbackReady:                governanceContract.RollbackRequirement == "mandatory" && !contains(governanceContract.Blockers, BlockerRollbackNotReady),
		ReturnToFloorReady:           governanceContract.ReturnToFloorRequirement == "mandatory" && !contains(governanceContract.Blockers, BlockerReturnToFloorNotReady),
		Blockers:                     dedupe(blockers),
		GuardPhase:                   phase,
	}
}

func (d DisabledWindowsFluidApplyExecutor) BuildCommandEnvelope(
	governanceContract WindowsFluidApplyGovernanceContract,
	guard WindowsFluidPreApplyGuard,
	evaluationTime time.Time,
) WindowsFluidExecutorCommandEnvelope {
	ts := evaluationTime.UTC()
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	commandClass := "evidence-refresh"
	switch governanceContract.RequestedAction {
	case GovernanceFutureCPUApply:
		commandClass = "cpu-lease"
	case GovernanceFutureMemoryApply:
		commandClass = "memory-lease"
	case GovernanceFutureReturnToFloor:
		commandClass = "return-to-floor"
	case GovernanceFutureRollback:
		commandClass = "rollback"
	}
	return WindowsFluidExecutorCommandEnvelope{
		EnvelopeID:                "windows-fluid-command-envelope-" + ts.Format("20060102T150405"),
		ShellRef:                  governanceContract.ShellRef,
		RequestedAction:           governanceContract.RequestedAction,
		CommandClass:              commandClass,
		RuntimeMode:               "in-place-qmp",
		CommandPreviewOnly:        true,
		ContainsExecutableCommand: false,
		QMPCommands:               []string{},
		ClusterMutations:          []string{},
		GuestMutations:            []string{},
		RequiredEvidenceRefs:      governanceContract.RequiredEvidenceAtApplyTime,
		RequiredAttestationRefs:   []string{guard.AttestationRef},
		DeniedReason:              BlockerFutureApplyExecutorDisabled,
		CreatedAt:                 ts.Format(time.RFC3339),
	}
}

func (d DisabledWindowsFluidApplyExecutor) DenyExecution(
	governanceContract WindowsFluidApplyGovernanceContract,
	guard WindowsFluidPreApplyGuard,
	envelope WindowsFluidExecutorCommandEnvelope,
	attestation WindowsFluidPolicyAttestation,
	evaluationTime time.Time,
) DisabledFutureApplyExecutionResult {
	ts := evaluationTime.UTC()
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	executionPhase := ExecutionHardDisabled
	switch guard.GuardPhase {
	case GuardQuarantined:
		executionPhase = ExecutionQuarantined
	case GuardBlocked:
		executionPhase = ExecutionBlocked
	case GuardNeedsRevalidation:
		executionPhase = ExecutionBlocked
	}
	if governanceContract.GovernancePhase == GovernanceContractBlocked {
		executionPhase = ExecutionBlocked
	}
	if governanceContract.GovernancePhase == GovernanceContractQuarantined {
		executionPhase = ExecutionQuarantined
	}
	if governanceContract.GovernancePhase == GovernanceNeedsRevalidation {
		executionPhase = ExecutionBlocked
	}
	if governanceContract.GovernancePhase == GovernanceContractBlocked &&
		(contains(governanceContract.DenialReasons, "pool_replica_context_only") || contains(governanceContract.DenialReasons, "generic_windows_vm_not_certified")) {
		executionPhase = ExecutionDenied
	}
	blockers := dedupe(append(append([]string{}, guard.Blockers...), governanceContract.Blockers...))
	blockers = append(blockers, BlockerFutureApplyExecutorDisabled)
	return DisabledFutureApplyExecutionResult{
		ExecutorID:            d.ExecutorID,
		ShellRef:              governanceContract.ShellRef,
		RequestedAction:       governanceContract.RequestedAction,
		ExecutionPhase:        executionPhase,
		ApplyAttempted:        false,
		MutationPerformed:     false,
		QMPCommandSent:        false,
		ClusterMutationSent:   false,
		Reason:                "executor hard-disabled: no runtime mutation",
		Blockers:              dedupe(blockers),
		RequiredFutureUnlocks: []string{"separate_executor_unlock_milestone", "formal_signed_attestation_pipeline", "explicit_runtime_mutation_authorization"},
		AttestationRefs:       []string{attestation.AttestationID},
		CreatedAt:             ts.Format(time.RFC3339),
	}
}

func (d DisabledWindowsFluidApplyExecutor) EmitExecutionDeniedEvidence(
	result DisabledFutureApplyExecutionResult,
	guard WindowsFluidPreApplyGuard,
	killSwitch WindowsFluidKillSwitch,
	envelope WindowsFluidExecutorCommandEnvelope,
	evaluationTime time.Time,
) map[string]any {
	ts := evaluationTime.UTC()
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	return map[string]any{
		"evidenceType":              "execution_denied",
		"executorId":                result.ExecutorID,
		"shellRef":                  result.ShellRef,
		"executionPhase":            result.ExecutionPhase,
		"applyAttempted":            result.ApplyAttempted,
		"mutationPerformed":         result.MutationPerformed,
		"qmpCommandSent":            result.QMPCommandSent,
		"clusterMutationSent":       result.ClusterMutationSent,
		"guardPhase":                guard.GuardPhase,
		"killSwitchMode":            killSwitch.Mode,
		"commandPreviewOnly":        envelope.CommandPreviewOnly,
		"containsExecutableCommand": envelope.ContainsExecutableCommand,
		"reason":                    result.Reason,
		"blockers":                  result.Blockers,
		"createdAt":                 ts.Format(time.RFC3339),
	}
}

func EvaluateWindowsFluidFutureApplyExecutor(input FutureApplyExecutorEvaluationInput) FutureApplyExecutorEvaluationResult {
	evaluationTime := input.EvaluationTime.UTC()
	if evaluationTime.IsZero() {
		evaluationTime = time.Now().UTC()
	}
	executor := NewDisabledWindowsFluidApplyExecutor()
	killSwitch := input.KillSwitch
	if killSwitch == nil {
		missingKillSwitch := WindowsFluidKillSwitch{
			KillSwitchID:     "",
			Enabled:          false,
			Source:           "environment",
			Mode:             "hard-disabled",
			Reason:           "kill switch proof missing or unobservable",
			RequiredForApply: true,
			ObservedAt:       evaluationTime.Format(time.RFC3339),
		}
		killSwitch = &missingKillSwitch
	}
	guard := executor.EvaluatePreApplyGuard(
		input.GovernanceContract,
		input.Revalidation,
		input.Attestation,
		*killSwitch,
		evaluationTime,
	)
	envelope := executor.BuildCommandEnvelope(input.GovernanceContract, guard, evaluationTime)
	execution := executor.DenyExecution(input.GovernanceContract, guard, envelope, input.Attestation, evaluationTime)
	deniedEvidence := executor.EmitExecutionDeniedEvidence(execution, guard, *killSwitch, envelope, evaluationTime)
	nextSafeStep := "keep-executor-hard-disabled-and-refresh-governance-evidence"
	switch execution.ExecutionPhase {
	case ExecutionQuarantined:
		nextSafeStep = "restore-identity-continuity-before-any-future-unlock-review"
	case ExecutionBlocked:
		nextSafeStep = "resolve-guard-and-revalidation-blockers-while-executor-stays-disabled"
	case ExecutionDenied:
		nextSafeStep = "use-certified-single-vm-candidate-and-deny-pool-generic-models"
	}
	return FutureApplyExecutorEvaluationResult{
		ExecutionResult:    execution,
		PreApplyGuard:      guard,
		KillSwitchSnapshot: *killSwitch,
		CommandEnvelope:    envelope,
		Blockers:           execution.Blockers,
		NextSafeStep:       nextSafeStep,
		DeniedEvidence:     deniedEvidence,
	}
}
