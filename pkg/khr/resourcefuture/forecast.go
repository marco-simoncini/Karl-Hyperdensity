package resourcefuture

import (
	"fmt"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/windowslane"
)

func buildForecasts(cells []lanediscovery.DiscoveredCell, ports []lanediscovery.DiscoveredResourcePort) (
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

		eligible, eligReason := liveInPlaceFor(lane, class, cell.Running)
		liveInPlace = append(liveInPlace, LiveInPlaceEligibility{
			TargetRef: cell.Ref, Lane: lane, Eligible: eligible, Reason: eligReason,
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
	}
	l, p, c, ls, _ := lanediscovery.ClassifyWorkload(h)
	return l, p, c, ls, ""
}

func liveInPlaceFor(lane, class string, running bool) (bool, string) {
	if !running {
		return false, "workload not running"
	}
	if lane == lanediscovery.LaneLinuxContainerCgroup && class == lanediscovery.ClassificationLiveInPlaceCapable {
		return true, "linux cgroup sandbox lane supports live-in-place simulation"
	}
	return false, "live-in-place not asserted on this lane"
}

func restartFor(lane, osFamily string, running bool) (required bool, risk, reason string) {
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
	} else if lane != lanediscovery.LaneLinuxContainerCgroup {
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
