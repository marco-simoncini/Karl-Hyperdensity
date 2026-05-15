package recommendation

import gpevidence "github.com/marco-simoncini/Karl-Hyperdensity/pkg/grandepadre/evidence"

// RecommendationReport is JSON output for khr-linux-agent -mode recommend-actions-local.
type RecommendationReport struct {
	GeneratedAt         string      `json:"generatedAt"`
	IndexedCount        int         `json:"indexedCount"`
	RecommendationCount int         `json:"recommendationCount"`
	ActionSlate         ActionSlate `json:"actionSlate"`
	MutationsForbidden  bool        `json:"mutationsForbidden"`
	ApplyAllowed        bool        `json:"applyAllowed"`
}

// BuildRecommendationReport runs the slate builder and attaches policy fields.
func BuildRecommendationReport(s *gpevidence.Store, o EngineOptions) RecommendationReport {
	slate := BuildActionSlate(s, o)
	idx := len(FilterIndicesByTenant(s.Snapshot(), o.Tenant))
	return RecommendationReport{
		GeneratedAt:         slate.GeneratedAt,
		IndexedCount:        idx,
		RecommendationCount: len(slate.Recommendations),
		ActionSlate:         slate,
		MutationsForbidden:  true,
		ApplyAllowed:        false,
	}
}
