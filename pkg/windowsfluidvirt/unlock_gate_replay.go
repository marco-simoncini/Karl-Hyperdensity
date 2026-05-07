package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"time"
)

type UnlockGateReplayFixture struct {
	Name                    string                               `json:"name"`
	GateID                  UnlockGateID                         `json:"gateId,omitempty"`
	EvidenceBundle          *WindowsFluidRuntimeEvidenceBundle   `json:"evidenceBundle,omitempty"`
	GovernanceContract      *WindowsFluidApplyGovernanceContract `json:"governanceContract,omitempty"`
	ExecutorOutput          *FutureApplyExecutorEvaluationResult `json:"executorOutput,omitempty"`
	Attestation             *WindowsFluidPolicyAttestation       `json:"attestation,omitempty"`
	ExpectedGateStatus      UnlockGateStatus                     `json:"expectedGateStatus,omitempty"`
	ExpectedAggregateStatus UnlockGateSetAggregateStatus         `json:"expectedAggregateStatus,omitempty"`
	ExpectedBlockers        []string                             `json:"expectedBlockers,omitempty"`
}

func LoadUnlockGateReplayFixture(path string) (UnlockGateReplayFixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return UnlockGateReplayFixture{}, err
	}
	var fixture UnlockGateReplayFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return UnlockGateReplayFixture{}, err
	}
	return fixture, nil
}

func EvaluateUnlockGateReplayFixture(path string, evaluationTime time.Time) (WindowsFluidUnlockGateVerification, UnlockGateReplayFixture, error) {
	fixture, err := LoadUnlockGateReplayFixture(path)
	if err != nil {
		return WindowsFluidUnlockGateVerification{}, UnlockGateReplayFixture{}, err
	}
	result := EvaluateWindowsFluidUnlockGate(UnlockGateEvaluationInput{
		GateID:             fixture.GateID,
		EvidenceBundle:     fixture.EvidenceBundle,
		GovernanceContract: fixture.GovernanceContract,
		ExecutorOutput:     fixture.ExecutorOutput,
		Attestation:        fixture.Attestation,
		EvaluationTime:     evaluationTime,
	})
	return result, fixture, nil
}

func EvaluateUnlockGateSetReplayFixture(path string, evaluationTime time.Time) (WindowsFluidUnlockGateSetVerification, UnlockGateReplayFixture, error) {
	fixture, err := LoadUnlockGateReplayFixture(path)
	if err != nil {
		return WindowsFluidUnlockGateSetVerification{}, UnlockGateReplayFixture{}, err
	}
	result := EvaluateWindowsFluidUnlockGateSet(UnlockGateSetEvaluationInput{
		EvidenceBundle:     fixture.EvidenceBundle,
		GovernanceContract: fixture.GovernanceContract,
		ExecutorOutput:     fixture.ExecutorOutput,
		Attestation:        fixture.Attestation,
		EvaluationTime:     evaluationTime,
	})
	return result, fixture, nil
}
