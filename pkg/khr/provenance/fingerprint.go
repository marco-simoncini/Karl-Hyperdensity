package provenance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// FingerprintBytes hashes raw evidence bytes.
func FingerprintBytes(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	h := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(h[:])
}

// FingerprintJSON hashes canonical JSON (stable key order via marshal).
func FingerprintJSON(v any) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return FingerprintBytes(data), nil
}

// LineageHash derives a lineage anchor from source context and seed.
func LineageHash(cluster, namespace, lane, seed string) string {
	h := sha256.Sum256([]byte(cluster + "|" + namespace + "|" + lane + "|" + seed))
	return "lineage:" + hex.EncodeToString(h[:12])
}

// ProvenanceID is a stable provenance record identifier.
func ProvenanceID(generatedBy, fingerprint, lineage string) string {
	h := sha256.Sum256([]byte(generatedBy + "|" + fingerprint + "|" + lineage))
	return "khr-prov-" + hex.EncodeToString(h[:8])
}

// Match reports whether two provenance records align on fingerprint and lineage.
func Match(a, b Record) bool {
	if a.EvidenceFingerprint == "" || b.EvidenceFingerprint == "" {
		return false
	}
	if a.EvidenceFingerprint != b.EvidenceFingerprint {
		return false
	}
	if a.LineageHash != "" && b.LineageHash != "" && a.LineageHash != b.LineageHash {
		return false
	}
	if a.SourceCluster != "" && b.SourceCluster != "" && a.SourceCluster != b.SourceCluster {
		return false
	}
	if a.SourceLane != "" && b.SourceLane != "" && a.SourceLane != b.SourceLane {
		return false
	}
	return true
}
