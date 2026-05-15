package recommendation

import (
	"fmt"
	"strings"

	gpevidence "github.com/marco-simoncini/Karl-Hyperdensity/pkg/grandepadre/evidence"
)

// ActionType is a non-mutating recommendation category.
type ActionType string

const (
	ActionObserve              ActionType = "observe"
	ActionRemediate            ActionType = "remediate"
	ActionPrepareResourceLease ActionType = "prepare-resourcelease"
	ActionCollectMoreEvidence  ActionType = "collect-more-evidence"
)

// Resource hints for Hyperdensity envelope planning (skeleton).
type Resource string

const (
	ResourceCPU     Resource = "cpu"
	ResourceMemory  Resource = "memory"
	ResourceUnknown Resource = "unknown"
)

// EnvelopeMode hints cgroup envelope vs unknown.
type EnvelopeMode string

const (
	ModeEnvelope EnvelopeMode = "envelope"
	ModeUnknown  EnvelopeMode = "unknown"
)

// ActionRecommendation is a single dry-run suggestion (never apply authorization).
type ActionRecommendation struct {
	ActionID        string                  `json:"actionId"`
	ActionType      ActionType              `json:"actionType"`
	TargetCellRef   *gpevidence.CellRefLite `json:"targetCellRef,omitempty"`
	DonorCellRef    *gpevidence.CellRefLite `json:"donorCellRef,omitempty"`
	ReceiverCellRef *gpevidence.CellRefLite `json:"receiverCellRef,omitempty"`
	Resource        Resource                `json:"resource"`
	Mode            EnvelopeMode            `json:"mode"`
	Confidence      string                  `json:"confidence"`
	Risk            Risk                    `json:"risk"`
	Priority        Priority                `json:"priority"`
	Reasons         []string                `json:"reasons"`
	Prerequisites   []string                `json:"prerequisites"`
	DryRunOnly      bool                    `json:"dryRunOnly"`
	ApplyAllowed    bool                    `json:"applyAllowed"`
}

// SlateSummary aggregates counts for operators (still dry-run).
type SlateSummary struct {
	RecommendationCount int            `json:"recommendationCount"`
	DonorCandidates     int            `json:"donorCandidates"`
	ReceiverCandidates  int            `json:"receiverCandidates"`
	ByRisk              map[string]int `json:"byRisk"`
	ByActionType        map[string]int `json:"byActionType"`
	Notes               string         `json:"notes"`
}

// ActionSlate is the structured output of the local recommendation skeleton.
type ActionSlate struct {
	GeneratedAt        string                              `json:"generatedAt"`
	Source             string                              `json:"source"`
	Recommendations    []ActionRecommendation              `json:"recommendations"`
	Blocked            []gpevidence.EvidenceIndex          `json:"blocked"`
	Remediable         []gpevidence.BlockedRemediableIndex `json:"remediable"`
	DonorCandidates    []gpevidence.CellRefLite            `json:"donorCandidates"`
	ReceiverCandidates []gpevidence.CellRefLite            `json:"receiverCandidates"`
	Summary            SlateSummary                        `json:"summary"`
}

func summarizeSlate(slate *ActionSlate) {
	byRisk := map[string]int{}
	byType := map[string]int{}
	for _, r := range slate.Recommendations {
		byRisk[string(r.Risk)]++
		byType[string(r.ActionType)]++
	}
	slate.Summary = SlateSummary{
		RecommendationCount: len(slate.Recommendations),
		DonorCandidates:     len(slate.DonorCandidates),
		ReceiverCandidates:  len(slate.ReceiverCandidates),
		ByRisk:              byRisk,
		ByActionType:        byType,
		Notes:               "Local decision skeleton only; no apply. Hyperdensity commercial vision is predictive live resource market — future KHR apply gate required before mutations.",
	}
}

func actionID(prefix string, i int) string {
	return fmt.Sprintf("%s-%03d", strings.TrimSpace(prefix), i)
}
