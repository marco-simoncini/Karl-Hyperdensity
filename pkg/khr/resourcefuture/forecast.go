package resourcefuture

import (
	"fmt"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/policygates"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/windowslane"
)

// PolicyContext enables KHR-V certification registry gating on forecasts.
type PolicyContext struct {
	Registry *certregistry.Registry
	Gates    policygates.Gates
	Now      time.Time
}

func (p PolicyContext) active() bool {
	return p.Registry != nil
}

func buildForecasts(cells []lanediscovery.DiscoveredCell, ports []lanediscovery.DiscoveredResourcePort, policy PolicyContext) (
	plans []CandidateScalePlan,
	saturation []SaturationForecastEntry,
	blocked []BlockedConstraint,
	compat []CompatibilityFallbackPrediction,
	liveInPlace []LiveInPlaceEligibility,
	restartPred []RestartRequiredPrediction,
	bundle ForecastBundle,
) {
	portByCell := map[string]lanediscovery.DiscoveredResourcePort{}
	for _, p := range ports {
		portByCell[p.CellRef] = p
	}

	for _, cell := range cells {
		port, ok := portByCell[cell.Ref]
		lane := ""
		class := ""
		provider := cell.ProviderBinding
		if ok {
			lane = port.Lane
			class = port.Classification
			provider = port.ProviderBinding
		} else {
			lane, provider, class, _, _ = inferFromCell(cell)
		}

		eligible, eligReason, eligState, blockedReason, stale, uncertified :=
			evaluateLiveInPlace(lane, class, cell.Running, policy)
		liveInPlace = append(liveInPlace, LiveInPlaceEligibility{
			TargetRef: cell.Ref, Lane: lane, Eligible: eligible, Reason: eligReason,
			EligibilityState: eligState, BlockedReason: blockedReason,
			StaleEvidence: stale, UncertifiedLane: uncertified,
		})

		restartReq, risk, restartReason := restartFor(lane, cell.OSFamily, cell.Running)
		restartPred = append(restartPred, RestartRequiredPrediction{
			TargetRef: cell.Ref, Lane: lane, Required: restartReq, RiskLevel: risk, Reason: restartReason,
		})

		if class == lanediscovery.ClassificationCompatibilityFallback ||
			provider == windowslane.ProviderKubevirtCompatibility {
			compat = append(compat, CompatibilityFallbackPrediction{
				TargetRef: cell.Ref, Lane: lane, ProviderBinding: provider,
				Likely: true, Reason: "compatibility provider path; live-in-place not guaranteed",
			})
			bundle.WindowsCompatibility = append(bundle.WindowsCompatibility, CompatibilityFallbackPrediction{
				TargetRef: cell.Ref, Lane: lane, ProviderBinding: provider,
				Likely: cell.OSFamily == "windows", Reason: restartReason,
			})
		}

		cpuPlan, cpuSat, cpuBlock := cpuForecast(cell, lane, eligible)
		if cpuPlan != nil {
			plans = append(plans, *cpuPlan)
		}
		if cpuSat != nil {
			saturation = append(saturation, *cpuSat)
			bundle.CPUScale = append(bundle.CPUScale, *cpuSat)
		}
		if cpuBlock != nil {
			blocked = append(blocked, *cpuBlock)
		}

		ramPlan, ramSat, ramBlock := ramForecast(cell, lane, eligible)
		if ramPlan != nil {
			plans = append(plans, *ramPlan)
		}
		if ramSat != nil {
			saturation = append(saturation, *ramSat)
			bundle.RAMScale = append(bundle.RAMScale, *ramSat)
		}
		if ramBlock != nil {
			blocked = append(blocked, *ramBlock)
		}

		if lane != lanediscovery.LaneLinuxContainerCgroup && cell.Running {
			mixed := SaturationForecastEntry{
				TargetRef: cell.Ref, Lane: lane, Resource: "cpu+memory",
				RiskLevel: riskLevelForMixed(lane, cell.OSFamily),
				Score:     mixedScore(lane, cell.OSFamily),
				Horizon:   "15m",
				Reason:    "mixed lane forecast: compatibility path may split CPU/RAM outcomes",
			}
			saturation = append(saturation, mixed)
			bundle.MixedLane = append(bundle.MixedLane, mixed)
		}
	}
	return plans, saturation, blocked, compat, liveInPlace, restartPred, bundle
}

func inferFromCell(cell lanediscovery.DiscoveredCell) (lane, provider, class string, liveScale bool, _ string) {
	h := lanediscovery.WorkloadHint{
		Name: cell.Name, Namespace: cell.Namespace, OSFamily: cell.OSFamily,
		VMType: cell.VMType, Running: cell.Running, SandboxPod: cell.VMType == "container",
		NativeLive: lanediscovery.IsNativeLiveWorkload(cell.Name, nil, cell.Namespace),
	}
	l, p, c, ls, _ := lanediscovery.ClassifyWorkload(h)
	return l, p, c, ls, ""
}

func evaluateLiveInPlace(lane, class string, running bool, policy PolicyContext) (
	eligible bool, reason, eligibilityState, blockedReason string, staleEvidence, uncertifiedLane bool,
) {
	if !running {
		return false, "workload not running", policygates.EligibilityBlocked, "workload not running", false, false
	}
	if policy.active() {
		if policygates.RequiresRegistry(lane) || class == lanediscovery.ClassificationNativeLive {
			out := policygates.Evaluate(lane, policy.Registry, policy.Gates, policy.Now)
			if out.Eligible {
				return true, "certified lane: policy gates passed (KHR-V)", policygates.EligibilityEligible, "", false, false
			}
			return false, out.BlockedReason, policygates.EligibilityBlocked, out.BlockedReason,
				out.StaleEvidence, out.UncertifiedLane
		}
		out := policygates.Evaluate(lane, policy.Registry, policy.Gates, policy.Now)
		if !out.Eligible {
			return false, out.BlockedReason, policygates.EligibilityBlocked, out.BlockedReason,
				out.StaleEvidence, true
		}
	}
	return liveInPlaceLegacy(lane, class)
}

func liveInPlaceLegacy(lane, class string) (bool, string, string, string, bool, bool) {
	if lane == lanediscovery.LaneNativeLive || class == lanediscovery.ClassificationNativeLive {
		return true, "native-live lane: live-in-place eligible (KHR-S)", policygates.EligibilityEligible, "", false, false
	}
	if lane == lanediscovery.LaneLinuxContainerCgroup && class == lanediscovery.ClassificationLiveInPlaceCapable {
		return true, "linux cgroup sandbox lane supports live-in-place simulation",
			policygates.EligibilityEligible, "", false, false
	}
	return false, "live-in-place not asserted on this lane", policygates.EligibilityBlocked,
		"live-in-place not asserted on this lane", false, true
}

func restartFor(lane, osFamily string, running bool) (required bool, risk, reason string) {
	if lane == lanediscovery.LaneNativeLive {
		return false, "low", "native-live lane: restart not required for cgroup live scale"
	}
	if osFamily == "windows" || lane == lanediscovery.LaneWindowsVMSession {
		return true, "high", "Windows/kubevirt compatibility may require guest restart for resource change"
	}
	if lane == lanediscovery.LaneLinuxVMCompatibility || lane == lanediscovery.LaneKubevirtCompatibility {
		if running {
			return false, "medium", "KubeVirt VM may require restart for some scale operations (simulation)"
		}
		return false, "low", "stopped VM: restart prediction deferred"
	}
	return false, "low", "cgroup lane: restart not required for simulation"
}

func cpuForecast(cell lanediscovery.DiscoveredCell, lane string, liveEligible bool) (*CandidateScalePlan, *SaturationForecastEntry, *BlockedConstraint) {
	risk := "low"
	score := 0.2
	if !cell.Running {
		risk = "high"
		score = 0.9
	} else 	if lane != lanediscovery.LaneNativeLive && lane != lanediscovery.LaneLinuxContainerCgroup {
		risk = "medium"
		score = 0.55
	}
	sat := &SaturationForecastEntry{
		TargetRef: cell.Ref, Lane: lane, Resource: "cpu",
		RiskLevel: risk, Score: score, Horizon: "15m",
		Reason: fmt.Sprintf("cpu scale forecast for %s/%s", cell.OSFamily, lane),
	}
	plan := &CandidateScalePlan{
		PlanID: fmt.Sprintf("cpu-scale-up-%s", sanitizeID(cell.Name)),
		TargetRef: cell.Ref, Lane: lane, Resource: "cpu", Direction: "scaleUp",
		Mode: "envelope", DeltaSummary: "+250m (simulated)",
		LiveInPlaceEligible: liveEligible, SimulationOnly: true,
	}
	var block *BlockedConstraint
	if !liveEligible {
		plan.Blocked = true
		plan.BlockedReason = "live-in-place not eligible"
		block = &BlockedConstraint{
			Constraint: "live-in-place-required", TargetRef: cell.Ref, Lane: lane,
			Reason: plan.BlockedReason,
		}
	}
	return plan, sat, block
}

func ramForecast(cell lanediscovery.DiscoveredCell, lane string, liveEligible bool) (*CandidateScalePlan, *SaturationForecastEntry, *BlockedConstraint) {
	risk := "low"
	score := 0.25
	if cell.OSFamily == "windows" {
		risk = "high"
		score = 0.85
	} else if lane == lanediscovery.LaneLinuxVMCompatibility {
		risk = "medium"
		score = 0.6
	}
	sat := &SaturationForecastEntry{
		TargetRef: cell.Ref, Lane: lane, Resource: "memory",
		RiskLevel: risk, Score: score, Horizon: "15m",
		Reason: fmt.Sprintf("ram scale forecast for %s/%s", cell.OSFamily, lane),
	}
	plan := &CandidateScalePlan{
		PlanID: fmt.Sprintf("ram-scale-up-%s", sanitizeID(cell.Name)),
		TargetRef: cell.Ref, Lane: lane, Resource: "memory", Direction: "scaleUp",
		Mode: "scaleUp", DeltaSummary: "+64Mi (simulated)",
		LiveInPlaceEligible: liveEligible, SimulationOnly: true,
	}
	var block *BlockedConstraint
	if cell.OSFamily == "windows" {
		plan.Blocked = true
		plan.BlockedReason = windowslane.BlockMemoryOnCompatibility().Reason
		block = &BlockedConstraint{
			Constraint: windowslane.BlockedRequiresRestart, TargetRef: cell.Ref, Lane: lane,
			Reason: plan.BlockedReason,
		}
	} else if !liveEligible {
		plan.Blocked = true
		plan.BlockedReason = "live-in-place not eligible for RAM scale"
		block = &BlockedConstraint{
			Constraint: "live-in-place-required", TargetRef: cell.Ref, Lane: lane,
			Reason: plan.BlockedReason,
		}
	}
	return plan, sat, block
}

func riskLevelForMixed(lane, osFamily string) string {
	if osFamily == "windows" {
		return "high"
	}
	if lane == lanediscovery.LaneLinuxVMCompatibility {
		return "medium"
	}
	return "low"
}

func mixedScore(lane, osFamily string) float64 {
	if osFamily == "windows" {
		return 0.8
	}
	if lane == lanediscovery.LaneLinuxVMCompatibility {
		return 0.5
	}
	return 0.3
}

func sanitizeID(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' {
			out = append(out, c)
		} else {
			out = append(out, '-')
		}
	}
	return string(out)
}
