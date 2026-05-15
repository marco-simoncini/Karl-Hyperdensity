// Package evidence is an in-memory Grande Padre evidence store skeleton (Sprint 12).
package evidence

import "time"

// TrustTier classifies integrity posture for indexing (not authorization).
type TrustTier string

const (
	TrustUnsigned            TrustTier = "Unsigned"
	TrustDevOnly             TrustTier = "DevOnly"
	TrustIntegrityVerified   TrustTier = "IntegrityVerified"
	TrustIntegrityFailed     TrustTier = "IntegrityFailed"
	TrustUnknown             TrustTier = "Unknown"
)

// UnsignedDigestTrustPolicy selects how digest-only (no PKI signature) matches are labeled.
// Default is IntegrityVerified: canonical bundle bytes matched the declared digest.
// Unsigned labels the same situation explicitly when operators want "no signature" semantics.
type UnsignedDigestTrustPolicy int

const (
	// UnsignedDigestAsIntegrityVerified is the default (digest channel verified; not production PKI).
	UnsignedDigestAsIntegrityVerified UnsignedDigestTrustPolicy = iota
	// UnsignedDigestAsUnsigned records digest match without a production/signature trust path.
	UnsignedDigestAsUnsigned
)

// CellRefLite is a minimal object reference for indexing (Cell or other).
type CellRefLite struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Name       string `json:"name,omitempty"`
}

// EvidenceIndex is a flattened index row for Hyperdensity queries.
type EvidenceIndex struct {
	ArtifactID           string       `json:"artifactId"`
	BundleSha256         string       `json:"bundleSha256"`
	CellRef              *CellRefLite `json:"cellRef,omitempty"`
	Confidence           string       `json:"confidence"`
	ReadyForGrandePadre  bool         `json:"readyForGrandePadre"`
	BlockedReasons       []string     `json:"blockedReasons"`
	Warnings             []string     `json:"warnings"`
	TrustTier            TrustTier    `json:"trustTier"`
	IndexedAt            string       `json:"indexedAt"`
}

// BlockedRemediableIndex aggregates blocked evidence per cell for recommendations.
type BlockedRemediableIndex struct {
	CellRef            *CellRefLite `json:"cellRef,omitempty"`
	BlockedReasons     []string     `json:"blockedReasons"`
	RemediationHints   []string     `json:"remediationHints"`
	LastArtifactID     string       `json:"lastArtifactId"`
	Confidence         string       `json:"confidence"`
	TrustTier          TrustTier    `json:"trustTier"`
}

// IngestOutcome captures counters from a single ingest operation.
type IngestOutcome struct {
	IndexedCount   int
	DuplicateCount int
	Index          EvidenceIndex
}

// LocalIndexReport is JSON output for khr-linux-agent -mode index-evidence-local.
type LocalIndexReport struct {
	IndexedCount       int             `json:"indexedCount"`
	DuplicateCount     int             `json:"duplicateCount"`
	Ready              []EvidenceIndex `json:"ready"`
	Blocked            []EvidenceIndex `json:"blocked"`
	QueryResult        interface{}     `json:"queryResult"`
	MutationsForbidden bool            `json:"mutationsForbidden"`
}

// QueryKind selects optional CLI filtering.
type QueryKind string

const (
	QueryNone          QueryKind = ""
	QueryReady         QueryKind = "ready"
	QueryBlocked       QueryKind = "blocked"
	QueryByConfidence  QueryKind = "by-confidence"
	QueryByCell        QueryKind = "by-cell"
)

// NowFunc is injectable time (tests).
var NowFunc = func() time.Time { return time.Now().UTC() }
