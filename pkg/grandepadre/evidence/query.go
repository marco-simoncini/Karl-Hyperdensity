package evidence

import (
	"strings"
)

// GetByArtifactID returns index rows whose artifactId matches (trimmed).
func (s *Store) GetByArtifactID(id string) []EvidenceIndex {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	var out []EvidenceIndex
	for _, idx := range s.Snapshot() {
		if strings.TrimSpace(idx.ArtifactID) == id {
			out = append(out, idx)
		}
	}
	return out
}

// GetByCell returns rows whose cellRef matches namespace and name (both trimmed).
func (s *Store) GetByCell(namespace, name string) []EvidenceIndex {
	ns := strings.TrimSpace(namespace)
	nm := strings.TrimSpace(name)
	if ns == "" || nm == "" {
		return nil
	}
	var out []EvidenceIndex
	for _, idx := range s.Snapshot() {
		if idx.CellRef == nil {
			continue
		}
		if strings.TrimSpace(idx.CellRef.Namespace) == ns && strings.TrimSpace(idx.CellRef.Name) == nm {
			out = append(out, idx)
		}
	}
	return out
}

// ListReady returns bundles marked ready for Grande Padre with non-failed trust.
func (s *Store) ListReady() []EvidenceIndex {
	var out []EvidenceIndex
	for _, idx := range s.Snapshot() {
		if idx.ReadyForGrandePadre && idx.TrustTier != TrustIntegrityFailed && len(idx.BlockedReasons) == 0 {
			out = append(out, idx)
		}
	}
	return out
}

// ListBlocked returns bundles that are not ready, have blocked reasons, or failed integrity trust.
func (s *Store) ListBlocked() []EvidenceIndex {
	var out []EvidenceIndex
	for _, idx := range s.Snapshot() {
		if !idx.ReadyForGrandePadre || len(idx.BlockedReasons) > 0 || idx.TrustTier == TrustIntegrityFailed {
			out = append(out, idx)
		}
	}
	return out
}

// ListByConfidence filters by evidenceSummary.confidence (case-insensitive).
func (s *Store) ListByConfidence(level string) []EvidenceIndex {
	want := strings.ToLower(strings.TrimSpace(level))
	if want == "" {
		return nil
	}
	var out []EvidenceIndex
	for _, idx := range s.Snapshot() {
		if strings.ToLower(strings.TrimSpace(idx.Confidence)) == want {
			out = append(out, idx)
		}
	}
	return out
}

// ParseQueryKind maps CLI -query values to QueryKind.
func ParseQueryKind(s string) QueryKind {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "ready":
		return QueryReady
	case "blocked":
		return QueryBlocked
	case "by-confidence":
		return QueryByConfidence
	case "by-cell":
		return QueryByCell
	default:
		return QueryNone
	}
}
