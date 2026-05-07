package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"time"
)

type AdmissionReplayFixture struct {
	Name                     string                            `json:"name"`
	RequestedAction          RequestedAdmissionAction          `json:"requestedAction"`
	Bundle                   WindowsFluidRuntimeEvidenceBundle `json:"bundle"`
	PolicyPack               *WindowsFluidPolicyPack           `json:"policyPack,omitempty"`
	ExpectedAdmissionPhase   AdmissionPhase                    `json:"expectedAdmissionPhase"`
	ExpectedDecisionBlockers []string                          `json:"expectedDecisionBlockers"`
}

func LoadAdmissionReplayFixture(path string) (AdmissionReplayFixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return AdmissionReplayFixture{}, err
	}
	var fixture AdmissionReplayFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return AdmissionReplayFixture{}, err
	}
	return fixture, nil
}

func EvaluateAdmissionReplayFixture(path string, evaluationTime time.Time) (AdmissionEvaluationResult, AdmissionReplayFixture, error) {
	fixture, err := LoadAdmissionReplayFixture(path)
	if err != nil {
		return AdmissionEvaluationResult{}, AdmissionReplayFixture{}, err
	}
	result := EvaluateWindowsFluidAdmission(AdmissionEvaluationInput{
		Bundle:          fixture.Bundle,
		PolicyPack:      fixture.PolicyPack,
		RequestedAction: fixture.RequestedAction,
		EvaluationTime:  evaluationTime,
	})
	return result, fixture, nil
}
