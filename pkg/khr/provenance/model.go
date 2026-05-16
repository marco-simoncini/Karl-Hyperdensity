// Package provenance provides trust/provenance semantics for KHR artifacts (KHR-Y).
package provenance

import "time"

const (
	ModelID = "khr-trust-provenance-v1"

	StateTrusted   = "trusted"
	StateMismatch  = "mismatch"
	StateStale     = "stale"
	StateUnknown   = "unknown"

	DefaultStaleAfterSeconds int64 = 86400
)

// Record is provenance metadata attached to evidence, certifications, approvals, and graph nodes.
type Record struct {
	ProvenanceID        string `json:"provenanceId"`
	GeneratedBy         string `json:"generatedBy"`
	GeneratedAt         string `json:"generatedAt"`
	SourceCluster       string `json:"sourceCluster"`
	SourceNamespace     string `json:"sourceNamespace"`
	SourceLane          string `json:"sourceLane"`
	LineageHash         string `json:"lineageHash"`
	EvidenceFingerprint string `json:"evidenceFingerprint"`
}

// SourceContext locates artifact origin.
type SourceContext struct {
	Cluster   string
	Namespace string
	Lane      string
}

// ValidationSummary is read-only provenance validation output.
type ValidationSummary struct {
	ModelID                   string   `json:"modelId"`
	Sprint                    string   `json:"sprint"`
	ValidatedAt               string   `json:"validatedAt"`
	ProvenanceState           string   `json:"provenanceState"`
	LineageIntegrity          bool     `json:"lineageIntegrity"`
	EvidenceFingerprint       string   `json:"evidenceFingerprint,omitempty"`
	RegistryIntegrity         bool     `json:"registryIntegrity"`
	ApprovalProvenanceValid   bool     `json:"approvalProvenanceValid"`
	GraphLineageVerified      bool     `json:"graphLineageVerified"`
	TrustWarnings             []string `json:"trustWarnings,omitempty"`
	ReadOnly                  bool     `json:"readOnly"`
	NoAutonomousOrchestration bool     `json:"noAutonomousOrchestration"`
	NoMutation                bool     `json:"noMutation"`
	NoApply                   bool     `json:"noApply"`
}

// NewRecord builds provenance for an artifact.
func NewRecord(generatedBy string, src SourceContext, lineageSeed string, evidence []byte, at time.Time) Record {
	if at.IsZero() {
		at = time.Now().UTC()
	}
	lineage := LineageHash(src.Cluster, src.Namespace, src.Lane, lineageSeed)
	fp := FingerprintBytes(evidence)
	pid := ProvenanceID(generatedBy, fp, lineage)
	return Record{
		ProvenanceID:        pid,
		GeneratedBy:         generatedBy,
		GeneratedAt:         at.UTC().Format(time.RFC3339),
		SourceCluster:       src.Cluster,
		SourceNamespace:     src.Namespace,
		SourceLane:          src.Lane,
		LineageHash:         lineage,
		EvidenceFingerprint: fp,
	}
}
