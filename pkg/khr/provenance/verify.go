package provenance

import (
	"fmt"
	"os"
	"time"
)

// VerifyFingerprint compares recorded fingerprint to recomputed evidence hash.
func VerifyFingerprint(record Record, evidence []byte) error {
	if record.EvidenceFingerprint == "" {
		return fmt.Errorf("provenance missing evidenceFingerprint")
	}
	want := FingerprintBytes(evidence)
	if record.EvidenceFingerprint != want {
		return fmt.Errorf("fingerprint mismatch: have %s want %s", record.EvidenceFingerprint, want)
	}
	return nil
}

// VerifyFingerprintFile loads a file and verifies its fingerprint.
func VerifyFingerprintFile(record Record, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return VerifyFingerprint(record, data)
}

// IsStaleProvenance reports whether generatedAt is older than maxAge.
func IsStaleProvenance(record Record, now time.Time, maxAgeSeconds int64) bool {
	if maxAgeSeconds <= 0 {
		maxAgeSeconds = DefaultStaleAfterSeconds
	}
	at, err := time.Parse(time.RFC3339, record.GeneratedAt)
	if err != nil {
		return true
	}
	return now.Sub(at) > time.Duration(maxAgeSeconds)*time.Second
}

// VerifyLineage compares expected lineage hash.
func VerifyLineage(record Record, expectedLineage string) error {
	if record.LineageHash == "" {
		return fmt.Errorf("provenance missing lineageHash")
	}
	if expectedLineage != "" && record.LineageHash != expectedLineage {
		return fmt.Errorf("lineage mismatch: have %s want %s", record.LineageHash, expectedLineage)
	}
	return nil
}

// VerifyApprovalProvenance ensures approval provenance matches certification/registry provenance.
func VerifyApprovalProvenance(approval, cert Record) error {
	if approval.ProvenanceID == "" || cert.ProvenanceID == "" {
		return fmt.Errorf("approval or certification missing provenance")
	}
	if !Match(approval, cert) {
		return fmt.Errorf("approval provenance mismatch with certification")
	}
	return nil
}

// ClassifyState returns provenance trust state from validation errors.
func ClassifyState(err error, stale bool) string {
	if err != nil {
		return StateMismatch
	}
	if stale {
		return StateStale
	}
	return StateTrusted
}
