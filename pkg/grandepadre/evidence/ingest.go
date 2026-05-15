package evidence

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence/integrity"
)

// IngestDocument is a parsed EvidenceIngestRequest (subset used by the local store).
type IngestDocument struct {
	MetadataName        string
	MetadataNamespace   string
	MetadataAnnotations map[string]string
	SpecArtifactID      string
	Bundle              evidence.CollectEvidenceBundle
	Manifest            integrity.ArtifactManifest
	DigestLine          string
	RequireDigestMatch  bool
	AllowUnsigned       bool
	AllowLocalDevSig    bool
}

// ParseIngestDocument parses YAML or JSON EvidenceIngestRequest bytes.
func ParseIngestDocument(data []byte) (*IngestDocument, error) {
	data = trimBOM(data)
	var root map[string]interface{}
	if err := yaml.Unmarshal(data, &root); err != nil {
		if jerr := json.Unmarshal(data, &root); jerr != nil {
			return nil, fmt.Errorf("parse ingest document: yaml: %w; json: %v", err, jerr)
		}
	}
	meta, _ := root["metadata"].(map[string]interface{})
	var name, ns string
	var ann map[string]string
	if meta != nil {
		name, _ = meta["name"].(string)
		ns, _ = meta["namespace"].(string)
		if a, ok := meta["annotations"].(map[string]interface{}); ok {
			ann = map[string]string{}
			for k, v := range a {
				if s, ok := v.(string); ok {
					ann[k] = s
				}
			}
		}
	}
	spec, ok := root["spec"].(map[string]interface{})
	if !ok || spec == nil {
		return nil, fmt.Errorf("ingest document: missing spec")
	}
	bundleObj := spec["bundle"]
	if bundleObj == nil {
		return nil, fmt.Errorf("ingest document: missing spec.bundle")
	}
	bundleJSON, err := json.Marshal(bundleObj)
	if err != nil {
		return nil, fmt.Errorf("ingest bundle marshal: %w", err)
	}
	var bundle evidence.CollectEvidenceBundle
	if err := json.Unmarshal(bundleJSON, &bundle); err != nil {
		return nil, fmt.Errorf("ingest bundle decode: %w", err)
	}
	var manifest integrity.ArtifactManifest
	if mObj, ok := spec["manifest"]; ok && mObj != nil {
		mb, err := json.Marshal(mObj)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(mb, &manifest); err != nil {
			return nil, fmt.Errorf("ingest manifest decode: %w", err)
		}
	}
	digest, _ := spec["digest"].(string)
	digest = strings.TrimSpace(digest)
	policy, _ := spec["policy"].(map[string]interface{})
	requireDM := true
	allowUn := true
	allowLD := false
	if policy != nil {
		if v, ok := policy["requireDigestMatch"].(bool); ok {
			requireDM = v
		}
		if v, ok := policy["allowUnsigned"].(bool); ok {
			allowUn = v
		}
		if v, ok := policy["allowLocalDevSignature"].(bool); ok {
			allowLD = v
		}
	}
	art, _ := spec["artifactId"].(string)
	if strings.TrimSpace(art) == "" {
		art = strings.TrimSpace(manifest.ArtifactID)
	}
	return &IngestDocument{
		MetadataName:        strings.TrimSpace(name),
		MetadataNamespace:   strings.TrimSpace(ns),
		MetadataAnnotations: ann,
		SpecArtifactID:      strings.TrimSpace(art),
		Bundle:              bundle,
		Manifest:            manifest,
		DigestLine:          digest,
		RequireDigestMatch:  requireDM,
		AllowUnsigned:       allowUn,
		AllowLocalDevSig:    allowLD,
	}, nil
}

func trimBOM(b []byte) []byte {
	if len(b) >= 3 && b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		return b[3:]
	}
	return b
}

func manifestPtr(m *integrity.ArtifactManifest) *integrity.ArtifactManifest {
	if m == nil {
		return nil
	}
	if strings.TrimSpace(m.ArtifactID) == "" && strings.TrimSpace(m.BundleSha256) == "" &&
		strings.TrimSpace(m.SigningMode) == "" && !m.SignaturePresent {
		return nil
	}
	return m
}

// BuildEvidenceIndex builds an index row from a parsed ingest document.
func BuildEvidenceIndex(doc *IngestDocument, unsignedLabel UnsignedDigestTrustPolicy) (EvidenceIndex, error) {
	canonical, err := integrity.CanonicalJSON(&doc.Bundle)
	if err != nil {
		return EvidenceIndex{}, fmt.Errorf("canonical bundle: %w", err)
	}
	computed := integrity.SHA256Hex(canonical)
	digestLine := strings.TrimSpace(doc.DigestLine)
	digestMatch := digestLine == computed && (doc.Manifest.BundleSha256 == "" || doc.Manifest.BundleSha256 == computed)

	tier := ComputeTrustTier(TrustInputs{
		DigestMatch:              digestMatch,
		RequireDigestMatch:       doc.RequireDigestMatch,
		Manifest:                 manifestPtr(&doc.Manifest),
		AllowUnsigned:            doc.AllowUnsigned,
		AllowLocalDevSignature:   doc.AllowLocalDevSig,
		MetadataAnnotations:      doc.MetadataAnnotations,
		UnsignedDigestTrustLabel: unsignedLabel,
	})

	summary := doc.Bundle.EvidenceSummary
	var cell *CellRefLite
	ref := doc.Bundle.CellRef
	if ref == nil && doc.Bundle.Telemetry.CellRef != nil {
		ref = doc.Bundle.Telemetry.CellRef
	}
	if ref != nil {
		cell = &CellRefLite{
			APIVersion: ref.APIVersion,
			Kind:       ref.Kind,
			Namespace:  ref.Namespace,
			Name:       ref.Name,
		}
	}
	artifactID := doc.SpecArtifactID
	if artifactID == "" {
		artifactID = doc.Manifest.ArtifactID
	}
	return EvidenceIndex{
		ArtifactID:          artifactID,
		BundleSha256:        computed,
		CellRef:             cell,
		Confidence:          summary.Confidence,
		ReadyForGrandePadre: summary.ReadyForGrandePadre,
		BlockedReasons:      copyStrs(summary.BlockedReasons),
		Warnings:            copyStrs(summary.Warnings),
		TrustTier:           tier,
		IndexedAt:           NowFunc().Format(time.RFC3339),
	}, nil
}

func copyStrs(s []string) []string {
	if len(s) == 0 {
		return []string{}
	}
	return append([]string(nil), s...)
}

// Ingest parses a request document, builds an index row, and stores it.
func (s *Store) Ingest(data []byte, unsignedLabel UnsignedDigestTrustPolicy) (IngestOutcome, error) {
	doc, err := ParseIngestDocument(data)
	if err != nil {
		return IngestOutcome{}, err
	}
	idx, err := BuildEvidenceIndex(doc, unsignedLabel)
	if err != nil {
		return IngestOutcome{}, err
	}
	dupBefore := s.DuplicateTotal()
	s.StoreBundle(idx)
	dupDelta := s.DuplicateTotal() - dupBefore
	return IngestOutcome{
		IndexedCount:   1,
		DuplicateCount: dupDelta,
		Index:          idx,
	}, nil
}
