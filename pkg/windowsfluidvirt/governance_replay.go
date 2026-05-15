package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"time"
)

type GovernanceReplayFixture struct {
	Name                    string                            `json:"name"`
	RequestedAction         GovernanceRequestedAction         `json:"requestedAction"`
	AdmissionDecision       WindowsFluidAdmissionDecision     `json:"admissionDecision"`
	Bundle                  WindowsFluidRuntimeEvidenceBundle `json:"bundle"`
	PolicyPack              *WindowsFluidPolicyPack           `json:"policyPack,omitempty"`
	ExpectedGovernancePhase GovernancePhase                   `json:"expectedGovernancePhase"`
	ExpectedBlockers        []string                          `json:"expectedBlockers"`
}

func LoadGovernanceReplayFixture(path string) (GovernanceReplayFixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return GovernanceReplayFixture{}, err
	}
	var fixture GovernanceReplayFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return GovernanceReplayFixture{}, err
	}
	return fixture, nil
}

func EvaluateGovernanceReplayFixture(path string, evaluationTime time.Time) (ApplyGovernanceEvaluationResult, GovernanceReplayFixture, error) {
	fixture, err := LoadGovernanceReplayFixture(path)
	if err != nil {
		return ApplyGovernanceEvaluationResult{}, GovernanceReplayFixture{}, err
	}
	result := EvaluateWindowsFluidApplyGovernance(ApplyGovernanceEvaluationInput{
		AdmissionDecision: fixture.AdmissionDecision,
		Bundle:            fixture.Bundle,
		PolicyPack:        fixture.PolicyPack,
		RequestedAction:   fixture.RequestedAction,
		EvaluationTime:    evaluationTime,
	})
	return result, fixture, nil
}
