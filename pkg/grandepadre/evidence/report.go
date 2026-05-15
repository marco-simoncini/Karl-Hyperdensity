package evidence

import (
	"fmt"
)

// LocalIndexParams configures RunLocalIndex (CLI wiring).
type LocalIndexParams struct {
	UnsignedLabel UnsignedDigestTrustPolicy
	Query         QueryKind
	CellNamespace string
	CellName      string
	Confidence    string
}

// RunLocalIndex ingests one request document, deduplicates by SHA-256, and builds the report.
func RunLocalIndex(s *Store, ingest []byte, p LocalIndexParams) (LocalIndexReport, error) {
	if len(ingest) == 0 {
		return LocalIndexReport{}, fmt.Errorf("empty ingest document")
	}
	if _, err := s.Ingest(ingest, p.UnsignedLabel); err != nil {
		return LocalIndexReport{}, err
	}
	s.DeduplicateBySha256()
	ready := s.ListReady()
	blocked := s.ListBlocked()
	if ready == nil {
		ready = []EvidenceIndex{}
	}
	if blocked == nil {
		blocked = []EvidenceIndex{}
	}
	rep := LocalIndexReport{
		IndexedCount:       s.Len(),
		DuplicateCount:     s.DuplicateTotal(),
		Ready:              ready,
		Blocked:            blocked,
		MutationsForbidden: true,
	}
	switch p.Query {
	case QueryReady:
		rep.QueryResult = ready
	case QueryBlocked:
		rep.QueryResult = blocked
	case QueryByConfidence:
		qr := s.ListByConfidence(p.Confidence)
		if qr == nil {
			qr = []EvidenceIndex{}
		}
		rep.QueryResult = qr
	case QueryByCell:
		qr := s.GetByCell(p.CellNamespace, p.CellName)
		if qr == nil {
			qr = []EvidenceIndex{}
		}
		rep.QueryResult = qr
	default:
		rep.QueryResult = nil
	}
	return rep, nil
}
