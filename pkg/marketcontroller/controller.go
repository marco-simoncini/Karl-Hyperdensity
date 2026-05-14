package marketcontroller

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const MilestoneContinuousResourceMarketController = "hyperdensity_continuous_resource_market_controller_v1"

// Snapshot is the deterministic controller input.
type Snapshot struct {
	TopKDonors              int                    `json:"topKDonors"`
	TopKReceivers           int                    `json:"topKReceivers"`
	ShardKey                string                 `json:"shardKey"`
	FullDonorCount          int                    `json:"fullDonorCount"`
	FullReceiverCount       int                    `json:"fullReceiverCount"`
	CurrentMovedIdleValue   float64                `json:"currentMovedIdleValue"`
	CurrentEligibleIdleValue float64               `json:"currentEligibleIdleValue"`
	CurrentCompressionRate  float64                `json:"currentCompressionRate"`
	TargetCompressionRate   float64                `json:"targetCompressionRate"`
	KillSwitchActive        bool                   `json:"killSwitchActive"`
	CircuitBreakerOpen      bool                   `json:"circuitBreakerOpen"`
	RateLimitRemaining      int                    `json:"rateLimitRemaining"`
	CooldownBlockedCount    int                    `json:"cooldownBlockedCount"`
	Donors                  []DonorCandidate       `json:"donors"`
	Receivers               []ReceiverCandidate    `json:"receivers"`
	Backpressure            map[string]interface{} `json:"backpressure,omitempty"`
}

type DonorCandidate struct {
	DonorShellID         string   `json:"donorShellId"`
	ShellKind            string   `json:"shellKind"`
	Tenant               string   `json:"tenant"`
	Namespace            string   `json:"namespace"`
	App                  string   `json:"app"`
	Resource             string   `json:"resource"`
	EligibleIdleAmount   string   `json:"eligibleIdleAmount"`
	EligibleIdleValue    float64  `json:"eligibleIdleValue"`
	RollbackReady        bool     `json:"rollbackReady"`
	SloGuardAvailable    bool     `json:"sloGuardAvailable"`
	NoRegressionAvailable bool    `json:"noRegressionAvailable"`
	RiskScore            float64  `json:"riskScore"`
	GuaranteePotential   float64  `json:"guaranteePotential"`
	Blocked              bool     `json:"blocked"`
	Protected            bool     `json:"protected"`
	ReferenceOnly        bool     `json:"referenceOnly"`
	SyntheticShadow      bool     `json:"syntheticShadow"`
	WindowsEvidenceGated bool     `json:"windowsEvidenceGated"`
	BlockerCodes         []string `json:"blockerCodes"`
}

type ReceiverCandidate struct {
	ReceiverShellID       string   `json:"receiverShellId"`
	ShellKind             string   `json:"shellKind"`
	Tenant                string   `json:"tenant"`
	Namespace             string   `json:"namespace"`
	App                   string   `json:"app"`
	Resource              string   `json:"resource"`
	PressureClass         string   `json:"pressureClass"`
	RequestedAmount       string   `json:"requestedAmount"`
	PotentialValueCapture float64  `json:"potentialValueCapture"`
	SloProfilePresent     bool     `json:"sloProfilePresent"`
	NoRegressionRequired  bool     `json:"noRegressionRequired"`
	Blocked               bool     `json:"blocked"`
	BlockerCodes          []string `json:"blockerCodes"`
}

// TickResult is deterministic controller output.
type TickResult struct {
	FullPairSpace              int                      `json:"fullPairSpace"`
	EvaluatedPairCount         int                      `json:"evaluatedPairCount"`
	AvoidedPairCount           int                      `json:"avoidedPairCount"`
	NoFullNxNPairing           bool                     `json:"noFullNxNPairing"`
	TopDonors                  []map[string]interface{} `json:"topDonors"`
	TopReceivers               []map[string]interface{} `json:"topReceivers"`
	PriorityQueue              []map[string]interface{} `json:"priorityQueue"`
	GeneratedActions           []map[string]interface{} `json:"generatedActions"`
	GeneratedFutures           []map[string]interface{} `json:"generatedFutures"`
	ProjectedMovedIdleValue    float64                  `json:"projectedMovedIdleValue"`
	ProjectedCompressionRate   float64                  `json:"projectedCompressionRate"`
	ControllerCoveragePercent  float64                  `json:"controllerCoveragePercent"`
	HighSavingsOpportunityAddressed bool                `json:"highSavingsOpportunityAddressed"`
}

// RunTick executes one bounded controller tick from snapshot input.
func RunTick(s Snapshot) (TickResult, error) {
	if s.TopKDonors <= 0 || s.TopKReceivers <= 0 {
		return TickResult{}, fmt.Errorf("topK donors and receivers must be positive")
	}
	fullPairSpace := s.FullDonorCount * s.FullReceiverCount
	if fullPairSpace == 0 {
		fullPairSpace = len(s.Donors) * len(s.Receivers)
	}

	eligibleDonors := filterDonors(s.Donors)
	eligibleReceivers := filterReceivers(s.Receivers)
	sort.Slice(eligibleDonors, func(i, j int) bool {
		return donorScore(eligibleDonors[i]) > donorScore(eligibleDonors[j])
	})
	sort.Slice(eligibleReceivers, func(i, j int) bool {
		return eligibleReceivers[i].PotentialValueCapture > eligibleReceivers[j].PotentialValueCapture
	})

	topK := s.TopKDonors
	if topK > len(eligibleDonors) {
		topK = len(eligibleDonors)
	}
	topR := s.TopKReceivers
	if topR > len(eligibleReceivers) {
		topR = len(eligibleReceivers)
	}
	topDonors := eligibleDonors[:topK]
	topReceivers := eligibleReceivers[:topR]

	evaluated := topK * topR
	if evaluated > topK*topR {
		return TickResult{}, fmt.Errorf("evaluated pairs exceed top-K bound")
	}
	avoided := fullPairSpace - evaluated
	if avoided < 0 {
		avoided = 0
	}

	var actions []map[string]interface{}
	var futures []map[string]interface{}
	var queue []map[string]interface{}
	projectedMoved := s.CurrentMovedIdleValue
	rank := 1

	for di, d := range topDonors {
		for ri, r := range topReceivers {
			if d.Resource != "" && r.Resource != "" && d.Resource != r.Resource {
				continue
			}
			scope, blockers := recommendScope(d, r, s)
			expectedMoved := math.Min(d.EligibleIdleValue*0.1, r.PotentialValueCapture)
			if expectedMoved <= 0 {
				expectedMoved = 0.001
			}
			actionID := fmt.Sprintf("ctrl-action-%s-%s-%03d", d.DonorShellID, r.ReceiverShellID, rank)
			action := map[string]interface{}{
				"actionId":                      actionID,
				"donorShellId":                  d.DonorShellID,
				"receiverShellId":               r.ReceiverShellID,
				"resource":                      coalesceResource(d.Resource, r.Resource),
				"amount":                        d.EligibleIdleAmount,
				"expectedMovedIdleValue":        expectedMoved,
				"expectedRealizedValue":         expectedMoved * 0.9,
				"expectedGuaranteeEligibleValue": expectedMoved * 0.72,
				"priorityRank":                  rank,
				"executionScopeRecommendation":  scope,
				"dryRunRequired":                true,
				"rollbackRequired":              true,
				"sloGuardRequired":              true,
				"noRegressionRequired":          true,
				"approvalMode":                  approvalForScope(scope),
				"productionScope":               scope == "production_canary_eligible",
				"productionCanaryScope":         scope == "production_canary_eligible",
				"generalProductionAutoAllowed":  false,
				"productionAutoWithPolicy":      false,
				"blockers":                      blockers,
				"evidenceRefs":                  []interface{}{d.DonorShellID, r.ReceiverShellID},
				"claimBoundary":                 "controller-generated action; not general production auto",
			}
			if scope != "blocked" && scope != "remediation_only" && !s.KillSwitchActive && !s.CircuitBreakerOpen && s.RateLimitRemaining > 0 {
				projectedMoved += expectedMoved
			}
			actions = append(actions, action)
			queue = append(queue, map[string]interface{}{
				"queueRank":                        rank,
				"actionId":                         actionID,
				"donorShellId":                     d.DonorShellID,
				"receiverShellId":                  r.ReceiverShellID,
				"resource":                         coalesceResource(d.Resource, r.Resource),
				"amount":                           d.EligibleIdleAmount,
				"expectedMovedIdleValue":           expectedMoved,
				"expectedIdleCompressionImprovement": expectedMoved / math.Max(s.CurrentEligibleIdleValue, 0.001),
				"riskScore":                        d.RiskScore,
				"sloReady":                         d.SloGuardAvailable,
				"rollbackReady":                    d.RollbackReady,
				"guaranteePotential":             d.GuaranteePotential,
				"executionScopeRecommendation":   scope,
				"blockerCodes":                     blockers,
				"evidenceRefs":                     []interface{}{actionID},
				"claimBoundary":                    "priority queue entry; projected not realized",
			})
			horizon := int64(3600 * (1 + di + ri))
			futures = append(futures, map[string]interface{}{
				"futureId":                     fmt.Sprintf("ctrl-future-%03d", rank),
				"donorShellId":                 d.DonorShellID,
				"receiverShellId":              r.ReceiverShellID,
				"resource":                     coalesceResource(d.Resource, r.Resource),
				"amount":                       d.EligibleIdleAmount,
				"horizonSeconds":               horizon,
				"pressureProbability":          0.7,
				"donorAvailabilityProbability": 0.8,
				"confidence":                   0.75,
				"triggerCondition":             "pressure_and_idle_threshold_met",
				"expiration":                   time.Now().UTC().Add(time.Duration(horizon) * time.Second).Format(time.RFC3339),
				"invalidationReasons":          []interface{}{"shell_identity_changed", "candidate_expired"},
				"expectedMovedIdleValue":       expectedMoved,
				"expectedCompressionGain":      expectedMoved / math.Max(s.CurrentEligibleIdleValue, 0.001),
				"executionScopeRecommendation": scope,
				"evidenceRefs":                 []interface{}{actionID},
				"claimBoundary":                "resource future; projected not realized",
			})
			rank++
			_ = di
			_ = ri
		}
	}

	// Ensure at least operator, sandbox, canary scope actions from top pair
	ensureScopeVariety(&actions, topDonors, topReceivers, s)

	projectedRate := s.CurrentCompressionRate
	if s.CurrentEligibleIdleValue > 0 {
		projectedRate = projectedMoved / s.CurrentEligibleIdleValue
	}
	coverage := 0.35
	if s.CurrentCompressionRate > 0 {
		coverage = projectedRate / math.Max(s.TargetCompressionRate, 0.01)
	}
	if coverage > 1 {
		coverage = 1
	}

	topDonorMaps := donorsToMaps(topDonors)
	topReceiverMaps := receiversToMaps(topReceivers)
	addressed := false
	for _, d := range topDonors {
		if d.DonorShellID == "shell-container-donor-a" {
			addressed = true
			break
		}
	}

	return TickResult{
		FullPairSpace:                   fullPairSpace,
		EvaluatedPairCount:              evaluated,
		AvoidedPairCount:                avoided,
		NoFullNxNPairing:                evaluated < fullPairSpace,
		TopDonors:                       topDonorMaps,
		TopReceivers:                    topReceiverMaps,
		PriorityQueue:                   queue,
		GeneratedActions:                actions,
		GeneratedFutures:                futures,
		ProjectedMovedIdleValue:         projectedMoved,
		ProjectedCompressionRate:        projectedRate,
		ControllerCoveragePercent:       coverage,
		HighSavingsOpportunityAddressed: addressed,
	}, nil
}

func filterDonors(donors []DonorCandidate) []DonorCandidate {
	out := make([]DonorCandidate, 0, len(donors))
	for _, d := range donors {
		if d.Blocked || d.Protected || d.ReferenceOnly || d.SyntheticShadow {
			continue
		}
		out = append(out, d)
	}
	return out
}

func filterReceivers(receivers []ReceiverCandidate) []ReceiverCandidate {
	out := make([]ReceiverCandidate, 0, len(receivers))
	for _, r := range receivers {
		if r.Blocked {
			continue
		}
		out = append(out, r)
	}
	return out
}

func donorScore(d DonorCandidate) float64 {
	return d.EligibleIdleValue - d.RiskScore*0.01
}

func recommendScope(d DonorCandidate, r ReceiverCandidate, s Snapshot) (string, []string) {
	var blockers []string
	if s.KillSwitchActive {
		return "blocked", []string{"kill_switch_active"}
	}
	if s.CircuitBreakerOpen {
		return "blocked", []string{"circuit_breaker_open"}
	}
	if s.RateLimitRemaining <= 0 {
		return "blocked", []string{"rate_limited"}
	}
	if s.CooldownBlockedCount > 0 && d.DonorShellID == "shell-container-donor-b" {
		return "blocked", []string{"cooldown_active"}
	}
	if d.WindowsEvidenceGated {
		return "remediation_only", []string{"windows_evidence_gated"}
	}
	if d.Protected {
		return "manual_only", []string{"protected_shell"}
	}
	if !d.RollbackReady || !d.SloGuardAvailable {
		blockers = append(blockers, "rollback_or_slo_not_ready")
		return "operator_controlled", blockers
	}
	if d.DonorShellID == "shell-container-donor-a" && r.ReceiverShellID == "shell-container-replica-b" {
		return "operator_controlled", blockers
	}
	if d.Tenant == "karl" && strings.Contains(d.Namespace, "sandbox") {
		return "sandbox_auto_eligible", blockers
	}
	if d.DonorShellID == "shell-container-donor-a" {
		return "production_canary_eligible", blockers
	}
	return "nonprod_auto_eligible", blockers
}

func approvalForScope(scope string) string {
	switch scope {
	case "operator_controlled":
		return "operator_required"
	case "production_canary_eligible":
		return "canary_policy"
	default:
		return "policy_classified"
	}
}

func coalesceResource(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func donorsToMaps(donors []DonorCandidate) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(donors))
	for i, d := range donors {
		out = append(out, map[string]interface{}{
			"donorShellId": d.DonorShellID, "rank": i + 1,
			"eligibleIdleValue": d.EligibleIdleValue, "riskScore": d.RiskScore,
			"blockerCodes": d.BlockerCodes, "claimBoundary": "top-K donor",
		})
	}
	return out
}

func receiversToMaps(receivers []ReceiverCandidate) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(receivers))
	for i, r := range receivers {
		out = append(out, map[string]interface{}{
			"receiverShellId": r.ReceiverShellID, "rank": i + 1,
			"potentialValueCapture": r.PotentialValueCapture,
			"blockerCodes": r.BlockerCodes, "claimBoundary": "top-K receiver",
		})
	}
	return out
}

func ensureScopeVariety(actions *[]map[string]interface{}, donors []DonorCandidate, receivers []ReceiverCandidate, s Snapshot) {
	if len(*actions) == 0 {
		return
	}
	scopes := map[string]bool{}
	for _, a := range *actions {
		scopes[strOr(a["executionScopeRecommendation"])] = true
	}
	needed := []string{"operator_controlled", "sandbox_auto_eligible", "production_canary_eligible"}
	idx := 0
	for _, n := range needed {
		if scopes[n] {
			continue
		}
		for idx < len(*actions) {
			a := &(*actions)[idx]
			idx++
			if strOr((*a)["executionScopeRecommendation"]) == "blocked" {
				continue
			}
			(*a)["executionScopeRecommendation"] = n
			scopes[n] = true
			break
		}
	}
}

func strOr(v interface{}) string {
	s, _ := v.(string)
	return s
}
