package windowsfluidvirt

import "time"

type HyperdensityParityOperation string

const (
	ParityCPUScaleUp   HyperdensityParityOperation = "cpu_scale_up"
	ParityCPUScaleDown HyperdensityParityOperation = "cpu_scale_down"
	ParityRAMScaleUp   HyperdensityParityOperation = "ram_scale_up"
	ParityRAMScaleDown HyperdensityParityOperation = "ram_scale_down"
)

type HyperdensityParityOperationProof struct {
	Operation                 HyperdensityParityOperation `json:"operation"`
	QMPConfirmedRuntimeState  bool                        `json:"qmpConfirmedRuntimeState"`
	GuestConfirmedActualState bool                        `json:"guestConfirmedActualState"`
	SameVM                    bool                        `json:"sameVm"`
	SameNamespace             bool                        `json:"sameNamespace"`
	SameNode                  bool                        `json:"sameNode"`
	SameVirtLauncherPod       bool                        `json:"sameVirtLauncherPod"`
	SameQEMUProcess           bool                        `json:"sameQemuProcess"`
	SameWindowsBoot           bool                        `json:"sameWindowsBoot"`
	SameMachineIdentity       bool                        `json:"sameMachineIdentity"`
	NoReboot                  bool                        `json:"noReboot"`
	NoRollout                 bool                        `json:"noRollout"`
	NoRecreate                bool                        `json:"noRecreate"`
	NoLiveMigration           bool                        `json:"noLiveMigration"`
	NoDestructiveMigration    bool                        `json:"noDestructiveMigration"`
	RollbackVerified          bool                        `json:"rollbackVerified"`
	ReturnToFloorVerified     bool                        `json:"returnToFloorVerified"`
	EvidenceBackedAudit       bool                        `json:"evidenceBackedAudit"`
}

type HyperdensityParityEvidence struct {
	CPUScaleUp   *HyperdensityParityOperationProof `json:"cpuScaleUp,omitempty"`
	CPUScaleDown *HyperdensityParityOperationProof `json:"cpuScaleDown,omitempty"`
	RAMScaleUp   *HyperdensityParityOperationProof `json:"ramScaleUp,omitempty"`
	RAMScaleDown *HyperdensityParityOperationProof `json:"ramScaleDown,omitempty"`
}

func (p HyperdensityParityOperationProof) Passed() bool {
	return p.QMPConfirmedRuntimeState &&
		p.GuestConfirmedActualState &&
		p.SameVM &&
		p.SameNamespace &&
		p.SameNode &&
		p.SameVirtLauncherPod &&
		p.SameQEMUProcess &&
		p.SameWindowsBoot &&
		p.SameMachineIdentity &&
		p.NoReboot &&
		p.NoRollout &&
		p.NoRecreate &&
		p.NoLiveMigration &&
		p.NoDestructiveMigration &&
		p.RollbackVerified &&
		p.ReturnToFloorVerified &&
		p.EvidenceBackedAudit
}

type HyperdensityParityGate struct {
	GateID UnlockGateID
}

func NewHyperdensityParityGate() HyperdensityParityGate {
	return HyperdensityParityGate{GateID: GateHyperdensityParityComplete}
}

func EvaluateHyperdensityParityGate(
	result WindowsFluidUnlockGateVerification,
	input UnlockGateEvaluationInput,
	evaluationTime time.Time,
) WindowsFluidUnlockGateVerification {
	gate := NewHyperdensityParityGate()
	return gate.Evaluate(result, input, evaluationTime)
}

func (g HyperdensityParityGate) Evaluate(
	result WindowsFluidUnlockGateVerification,
	input UnlockGateEvaluationInput,
	evaluationTime time.Time,
) WindowsFluidUnlockGateVerification {
	result.GateID = g.GateID
	result.RequiredInputs = []string{
		"parity_evidence",
		"cpu_scale_up_proof",
		"cpu_scale_down_proof",
		"ram_scale_up_proof",
		"ram_scale_down_proof",
		"qmp_guest_identity_audit_coherence",
		"rollback_and_return_to_floor_verification",
	}
	if input.ParityEvidence == nil {
		result.MissingInputs = append(result.MissingInputs, "parity_evidence")
		result.BlockerList = append(result.BlockerList, GateBlockerParityEvidenceMissing)
		result.GateStatus = UnlockGateBlocked
		result.DeterministicHash = gateHash(result, evaluationTime)
		return result
	}

	parity := input.ParityEvidence
	result.CheckedInputs = append(result.CheckedInputs, "parity_evidence")

	check := func(
		proof *HyperdensityParityOperationProof,
		op HyperdensityParityOperation,
		missingBlocker string,
		failedBlocker string,
		requiredInput string,
	) {
		if proof == nil {
			result.MissingInputs = append(result.MissingInputs, requiredInput)
			result.BlockerList = append(result.BlockerList, missingBlocker, GateBlockerHyperdensityParityPartialSuccessNotTotalFeasibility)
			return
		}
		result.CheckedInputs = append(result.CheckedInputs, string(op))
		if proof.Operation != op || !proof.Passed() {
			result.BlockerList = append(result.BlockerList, failedBlocker, GateBlockerHyperdensityParityPartialSuccessNotTotalFeasibility)
		}
	}

	check(parity.CPUScaleUp, ParityCPUScaleUp, GateBlockerParityCPUScaleUpMissing, GateBlockerParityCPUScaleUpFailed, "cpu_scale_up_proof")
	check(parity.CPUScaleDown, ParityCPUScaleDown, GateBlockerParityCPUScaleDownMissing, GateBlockerParityCPUScaleDownFailed, "cpu_scale_down_proof")
	check(parity.RAMScaleUp, ParityRAMScaleUp, GateBlockerParityRAMScaleUpMissing, GateBlockerParityRAMScaleUpFailed, "ram_scale_up_proof")
	check(parity.RAMScaleDown, ParityRAMScaleDown, GateBlockerParityRAMScaleDownMissing, GateBlockerParityRAMScaleDownFailed, "ram_scale_down_proof")

	result.BlockerList = dedupe(result.BlockerList)
	if len(result.BlockerList) == 0 {
		result.GateStatus = UnlockGatePassed
	} else {
		result.GateStatus = UnlockGateBlocked
	}
	result.DeterministicHash = gateHash(result, evaluationTime)
	return result
}
