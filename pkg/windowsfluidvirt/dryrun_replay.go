package windowsfluidvirt

import (
	"encoding/json"
	"os"
)

type DryRunReplayFixture struct {
	Name                   string                            `json:"name"`
	Bundle                 WindowsFluidRuntimeEvidenceBundle `json:"bundle"`
	ExpectedPhase          WindowsFluidPhase                 `json:"expectedPhase"`
	ExpectedClassification CertificationClassification       `json:"expectedClassification"`
	ExpectedBlockers       []string                          `json:"expectedBlockers"`
}

func LoadDryRunReplayFixture(path string) (DryRunReplayFixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DryRunReplayFixture{}, err
	}
	var fixture DryRunReplayFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return DryRunReplayFixture{}, err
	}
	return fixture, nil
}

func EvaluateDryRunReplayFixture(path string) (DryRunEvaluationResult, DryRunReplayFixture, error) {
	fixture, err := LoadDryRunReplayFixture(path)
	if err != nil {
		return DryRunEvaluationResult{}, DryRunReplayFixture{}, err
	}
	result := EvaluateWindowsFluidRuntimeDryRun(fixture.Bundle)
	return result, fixture, nil
}
