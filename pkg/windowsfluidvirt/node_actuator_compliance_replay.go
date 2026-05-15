package windowsfluidvirt

import (
	"encoding/json"
	"os"
)

type WindowsHyperdensityComplianceReplayFixture struct {
	Name                 string                                          `json:"name"`
	Input                EvaluateWindowsHyperdensityReadyComplianceInput `json:"input"`
	ExpectedPhase        WindowsHyperdensityCompliancePhase              `json:"expectedPhase"`
	ExpectedBlockers     []string                                        `json:"expectedBlockers"`
	ExpectedRemediations []string                                        `json:"expectedRemediations"`
}

type NodeActuatorSafetyReplayFixture struct {
	Name             string                        `json:"name"`
	Input            NodeFluidActuatorSafetyInput  `json:"input"`
	Model            *NodeFluidActuatorSafetyModel `json:"model,omitempty"`
	ExpectedAllowed  bool                          `json:"expectedAllowed"`
	ExpectedBlockers []string                      `json:"expectedBlockers"`
}

type WindowsCpuLeaseReplayFixture struct {
	Name             string                     `json:"name"`
	Lease            WindowsCpuEntitlementLease `json:"lease"`
	ExpectedStatus   string                     `json:"expectedStatus"`
	ExpectedBlockers []string                   `json:"expectedBlockers"`
}

func LoadWindowsHyperdensityComplianceReplayFixture(path string) (WindowsHyperdensityComplianceReplayFixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return WindowsHyperdensityComplianceReplayFixture{}, err
	}
	var fixture WindowsHyperdensityComplianceReplayFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return WindowsHyperdensityComplianceReplayFixture{}, err
	}
	return fixture, nil
}

func LoadNodeActuatorSafetyReplayFixture(path string) (NodeActuatorSafetyReplayFixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return NodeActuatorSafetyReplayFixture{}, err
	}
	var fixture NodeActuatorSafetyReplayFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return NodeActuatorSafetyReplayFixture{}, err
	}
	return fixture, nil
}

func LoadWindowsCpuLeaseReplayFixture(path string) (WindowsCpuLeaseReplayFixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return WindowsCpuLeaseReplayFixture{}, err
	}
	var fixture WindowsCpuLeaseReplayFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return WindowsCpuLeaseReplayFixture{}, err
	}
	return fixture, nil
}
