package resourcefuture

import (
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
)

const ModeSimulate = "resourcefuture-simulate"

// SafetyPolicy attests simulation-only execution (KHR-R).
type SafetyPolicy struct {
	ReadOnly              bool `json:"readOnly"`
	NoMutation            bool `json:"noMutation"`
	NoApply               bool `json:"noApply"`
	NoRestart             bool `json:"noRestart"`
	NoRollout             bool `json:"noRollout"`
	NoRecreate            bool `json:"noRecreate"`
	NoAutonomousOrchestration bool `json:"noAutonomousOrchestration"`
	NoProductionMutation  bool `json:"noProductionMutation"`
	SimulationOnly        bool `json:"simulationOnly"`
}

// DefaultSafetyPolicy is the KHR-R default.
func DefaultSafetyPolicy() SafetyPolicy {
	return SafetyPolicy{
		ReadOnly: true, NoMutation: true, NoApply: true,
		NoRestart: true, NoRollout: true, NoRecreate: true,
		NoAutonomousOrchestration: true, NoProductionMutation: true,
		SimulationOnly: true,
	}
}

// CandidateScalePlan is a read-only simulated scale proposal (not applied).
type CandidateScalePlan struct {
	PlanID              string `json:"planId"`
	TargetRef           string `json:"targetRef"`
	Lane                string `json:"lane"`
	Resource            string `json:"resource"`
	Direction           string `json:"direction"`
	Mode                string `json:"mode"`
	DeltaSummary        string `json:"deltaSummary,omitempty"`
	LiveInPlaceEligible bool   `json:"liveInPlaceEligible"`
	SimulationOnly      bool   `json:"simulationOnly"`
	Blocked             bool   `json:"blocked"`
	BlockedReason       string `json:"blockedReason,omitempty"`
}

// SaturationForecastEntry predicts pressure on a target.
type SaturationForecastEntry struct {
	TargetRef     string  `json:"targetRef"`
	Lane          string  `json:"lane"`
	Resource      string  `json:"resource"`
	RiskLevel     string  `json:"riskLevel"`
	Score         float64 `json:"score"`
	Horizon       string  `json:"horizon"`
	Reason        string  `json:"reason,omitempty"`
}

// BlockedConstraint is a predicted blocker for scale simulation.
type BlockedConstraint struct {
	Constraint string `json:"constraint"`
	TargetRef  string `json:"targetRef,omitempty"`
	Lane       string `json:"lane,omitempty"`
	Reason     string `json:"reason"`
}

// CompatibilityFallbackPrediction forecasts compatibility-path outcomes.
type CompatibilityFallbackPrediction struct {
	TargetRef       string `json:"targetRef"`
	Lane            string `json:"lane"`
	ProviderBinding string `json:"providerBinding"`
	Likely          bool   `json:"likely"`
	Reason          string `json:"reason"`
}

// LiveInPlaceEligibility predicts live-in-place scale eligibility.
type LiveInPlaceEligibility struct {
	TargetRef string `json:"targetRef"`
	Lane      string `json:"lane"`
	Eligible  bool   `json:"eligible"`
	Reason    string `json:"reason,omitempty"`
}

// RestartRequiredPrediction forecasts restart/reboot requirement.
type RestartRequiredPrediction struct {
	TargetRef string `json:"targetRef"`
	Lane      string `json:"lane"`
	Required  bool   `json:"required"`
	RiskLevel string `json:"riskLevel"`
	Reason    string `json:"reason"`
}

// ForecastBundle groups simulation forecasts by type.
type ForecastBundle struct {
	CPUScale                []SaturationForecastEntry `json:"cpuScale"`
	RAMScale                []SaturationForecastEntry `json:"ramScale"`
	MixedLane               []SaturationForecastEntry `json:"mixedLane"`
	WindowsCompatibility    []CompatibilityFallbackPrediction `json:"windowsCompatibility"`
}

// PostureSummary is read-only posture input snapshot for simulation.
type PostureSummary struct {
	Sandbox               bool     `json:"sandbox"`
	Preview               bool     `json:"preview"`
	ProductionUnsupported bool     `json:"productionUnsupported"`
	LaneCount             int      `json:"laneCount"`
	BlockedStateCount     int      `json:"blockedStateCount"`
	Blockers              []string `json:"blockers,omitempty"`
}

// SimulateInput is aggregated read-only context for planning.
type SimulateInput struct {
	HostStatus      crdv1alpha1.Host           `json:"hostStatus"`
	LaneDiscovery   lanediscovery.Result       `json:"laneDiscovery"`
	ResourcePorts   []lanediscovery.DiscoveredResourcePort `json:"resourcePorts,omitempty"`
	ResourceLeases  []SimulatedLeaseRef        `json:"resourceLeases,omitempty"`
	PostureSummary  PostureSummary             `json:"postureSummary"`
}

// SimulatedLeaseRef is a lightweight lease observation for simulation.
type SimulatedLeaseRef struct {
	Ref      string `json:"ref"`
	Resource string `json:"resource,omitempty"`
	Mode     string `json:"mode,omitempty"`
	DryRunOnly bool `json:"dryRunOnly,omitempty"`
}

// SimulationResult is CLI JSON for resourcefuture-simulate.
type SimulationResult struct {
	Mode                            string                            `json:"mode"`
	Blocked                         bool                              `json:"blocked"`
	Reason                          string                            `json:"reason,omitempty"`
	ClusterContext                  string                            `json:"clusterContext,omitempty"`
	ObservedAt                      string                            `json:"observedAt"`
	Safety                          SafetyPolicy                      `json:"safety"`
	Input                           SimulateInput                     `json:"input"`
	CandidateScalePlans             []CandidateScalePlan              `json:"candidateScalePlans"`
	SaturationForecast              []SaturationForecastEntry         `json:"saturationForecast"`
	BlockedConstraints              []BlockedConstraint               `json:"blockedConstraints"`
	CompatibilityFallbackPrediction []CompatibilityFallbackPrediction `json:"compatibilityFallbackPrediction"`
	LiveInPlaceEligibility          []LiveInPlaceEligibility          `json:"liveInPlaceEligibility"`
	RestartRequiredPrediction       []RestartRequiredPrediction         `json:"restartRequiredPrediction"`
	Forecasts                       ForecastBundle                    `json:"forecasts"`
	Summary                         map[string]any                    `json:"summary,omitempty"`
}
