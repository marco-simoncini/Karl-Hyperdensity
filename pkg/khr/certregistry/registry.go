// Package certregistry provides read-only lane certification registry (KHR-V).
package certregistry

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/nativelive"
)

const (
	RegistryID           = "khr-certification-registry-v1"
	CertStateCertified   = "certified-preview"
	CertStateFailed      = "failed"
	CertStateUncertified = "uncertified"

	DefaultValidForSeconds int64 = 86400 // 24h sandbox freshness window
)

// LaneAttestation captures gate inputs derived from certification evidence.
type LaneAttestation struct {
	NoRestart                  bool `json:"noRestart"`
	NoRollout                  bool `json:"noRollout"`
	NoRecreate                 bool `json:"noRecreate"`
	NoInterruption             bool `json:"noInterruption"`
	ShellContinuityPreserved   bool `json:"shellContinuityPreserved"`
	RollbackObserved           bool `json:"rollbackObserved"`
}

// LaneEntry is one registry row for a certified lane.
type LaneEntry struct {
	LaneID              string          `json:"laneId"`
	LaneType            string          `json:"laneType"`
	ProviderBinding     string          `json:"providerBinding"`
	CertificationState  string          `json:"certificationState"`
	ContinuityScore     float64         `json:"continuityScore"`
	LiveScaleConfidence string          `json:"liveScaleConfidence"`
	LastCertifiedAt     string          `json:"lastCertifiedAt"`
	EvidenceRef         string          `json:"evidenceRef"`
	ValidForSeconds     int64           `json:"validForSeconds"`
	Attestation         LaneAttestation `json:"attestation"`
}

// Registry is a read-only certification registry artifact.
type Registry struct {
	RegistryID                string      `json:"registryId"`
	Sprint                    string      `json:"sprint"`
	GeneratedAt               string      `json:"generatedAt"`
	ReadOnly                  bool        `json:"readOnly"`
	NoAutonomousOrchestration bool        `json:"noAutonomousOrchestration"`
	Entries                   []LaneEntry `json:"entries"`
}

// FindByLane returns the first entry matching lane id or lane type.
func (r *Registry) FindByLane(lane string) *LaneEntry {
	if r == nil {
		return nil
	}
	for i := range r.Entries {
		if r.Entries[i].LaneID == lane || r.Entries[i].LaneType == lane {
			return &r.Entries[i]
		}
	}
	return nil
}

// IsFresh reports whether entry evidence is within validForSeconds of now.
func IsFresh(entry LaneEntry, now time.Time) bool {
	if entry.ValidForSeconds <= 0 {
		return true
	}
	at, err := time.Parse(time.RFC3339, entry.LastCertifiedAt)
	if err != nil {
		return false
	}
	return now.Sub(at) <= time.Duration(entry.ValidForSeconds)*time.Second
}

// EntryFromCertificationSummary builds a registry row from native-live certification.
func EntryFromCertificationSummary(
	summary nativelive.CertificationSummary,
	evidenceRef string,
	validForSeconds int64,
	certifiedAt time.Time,
) LaneEntry {
	state := CertStateFailed
	if summary.Status == nativelive.CertificationCertified && !summary.RegressionDetected {
		state = CertStateCertified
	}
	if validForSeconds <= 0 {
		validForSeconds = DefaultValidForSeconds
	}
	rollbackObserved := summary.Metrics.RollbackLatencyMs.CPU > 0 ||
		summary.Metrics.RollbackLatencyMs.RAMUp > 0 ||
		summary.Metrics.RollbackLatencyMs.RAMDown > 0
	return LaneEntry{
		LaneID:              lanediscovery.LaneNativeLive,
		LaneType:            lanediscovery.LaneNativeLive,
		ProviderBinding:     "khr.native",
		CertificationState:  state,
		ContinuityScore:     summary.Scores.ContinuityScore,
		LiveScaleConfidence: summary.Scores.LiveScaleConfidence,
		LastCertifiedAt:     certifiedAt.UTC().Format(time.RFC3339),
		EvidenceRef:         evidenceRef,
		ValidForSeconds:     validForSeconds,
		Attestation: LaneAttestation{
			NoRestart:                summary.Invariants.NoRestart,
			NoRollout:                summary.Invariants.NoRollout,
			NoRecreate:               summary.Invariants.NoRecreate,
			NoInterruption:           !summary.Invariants.InterruptionDetected,
			ShellContinuityPreserved: summary.ContinuityProof.ShellContinuityPreserved,
			RollbackObserved:         rollbackObserved,
		},
	}
}

// GenerateFromSummary builds a KHR-V registry from the latest certification summary.
func GenerateFromSummary(
	sprint string,
	summary nativelive.CertificationSummary,
	evidenceRef string,
	validForSeconds int64,
	certifiedAt time.Time,
) Registry {
	return Registry{
		RegistryID:                RegistryID,
		Sprint:                    sprint,
		GeneratedAt:               time.Now().UTC().Format(time.RFC3339),
		ReadOnly:                  true,
		NoAutonomousOrchestration: true,
		Entries: []LaneEntry{
			EntryFromCertificationSummary(summary, evidenceRef, validForSeconds, certifiedAt),
		},
	}
}

// LoadJSON reads a registry from disk.
func LoadJSON(path string) (Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Registry{}, err
	}
	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return Registry{}, err
	}
	if reg.RegistryID == "" {
		return Registry{}, fmt.Errorf("registryId is required")
	}
	return reg, nil
}

// LoadSummaryJSON loads nativelive.CertificationSummary from path.
func LoadSummaryJSON(path string) (nativelive.CertificationSummary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nativelive.CertificationSummary{}, err
	}
	var s nativelive.CertificationSummary
	if err := json.Unmarshal(data, &s); err != nil {
		return nativelive.CertificationSummary{}, err
	}
	return s, nil
}
