package recommendation

import (
	"strings"

	gpevidence "github.com/marco-simoncini/Karl-Hyperdensity/pkg/grandepadre/evidence"
)

// DonorCandidates lists cells that may supply headroom in a future live resource market (dry-run only).
// Rules: ready + high confidence + trust not IntegrityFailed + trust not DevOnly.
func DonorCandidates(indices []gpevidence.EvidenceIndex) []gpevidence.CellRefLite {
	var out []gpevidence.CellRefLite
	seen := map[string]bool{}
	for _, idx := range indices {
		if !idx.ReadyForGrandePadre || len(idx.BlockedReasons) > 0 {
			continue
		}
		if idx.TrustTier == gpevidence.TrustIntegrityFailed || idx.TrustTier == gpevidence.TrustDevOnly {
			continue
		}
		if strings.ToLower(strings.TrimSpace(idx.Confidence)) != "high" {
			continue
		}
		if idx.CellRef == nil {
			continue
		}
		key := idx.CellRef.Namespace + "/" + idx.CellRef.Name
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, *idx.CellRef)
	}
	return out
}

// ReceiverCandidates lists cells that need remediation attention (blocked-like evidence).
func ReceiverCandidates(indices []gpevidence.EvidenceIndex) []gpevidence.CellRefLite {
	var out []gpevidence.CellRefLite
	seen := map[string]bool{}
	for _, idx := range indices {
		if !isBlockedLike(idx) {
			continue
		}
		if idx.CellRef == nil {
			continue
		}
		key := idx.CellRef.Namespace + "/" + idx.CellRef.Name
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, *idx.CellRef)
	}
	return out
}

func isBlockedLike(idx gpevidence.EvidenceIndex) bool {
	return !idx.ReadyForGrandePadre || len(idx.BlockedReasons) > 0 || idx.TrustTier == gpevidence.TrustIntegrityFailed
}
