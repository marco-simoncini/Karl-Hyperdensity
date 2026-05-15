package evidenceingest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence/integrity"
)

const (
	annotationPreparationWarnings = "khr.karl.io/preparation-warnings"
	annotationSignatureTrustTier  = "khr.karl.io/signature-trust-tier"
	trustTierDevOnly              = "DevOnly"
)

// PrepareIngestRequest reads bundle, manifest, and digest files and returns
// EvidenceIngestRequest document bytes (YAML or JSON). No kube apply and no HTTP.
func PrepareIngestRequest(bundlePath, manifestPath, digestPath string, o PrepareOptions) ([]byte, error) {
	if strings.TrimSpace(bundlePath) == "" {
		return nil, fmt.Errorf("prepare-ingest-request: bundle input path is required (-bundle-input)")
	}
	if strings.TrimSpace(manifestPath) == "" {
		return nil, fmt.Errorf("prepare-ingest-request: manifest input path is required (-manifest-input)")
	}
	if strings.TrimSpace(digestPath) == "" {
		return nil, fmt.Errorf("prepare-ingest-request: digest input path is required (-digest-input)")
	}

	bundleRaw, err := os.ReadFile(bundlePath)
	if err != nil {
		return nil, fmt.Errorf("read bundle: %w", err)
	}
	manifestRaw, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	digestRaw, err := os.ReadFile(digestPath)
	if err != nil {
		return nil, fmt.Errorf("read digest: %w", err)
	}

	var bundle evidence.CollectEvidenceBundle
	if err := json.Unmarshal(bundleRaw, &bundle); err != nil {
		return nil, fmt.Errorf("parse bundle json: %w", err)
	}
	var manifest integrity.ArtifactManifest
	if err := json.Unmarshal(manifestRaw, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest json: %w", err)
	}

	canonical, err := integrity.CanonicalJSON(&bundle)
	if err != nil {
		return nil, fmt.Errorf("canonical bundle: %w", err)
	}
	computedSHA := integrity.SHA256Hex(canonical)
	digestLine := strings.TrimSpace(string(digestRaw))

	var warns []string
	if manifest.BundleSha256 != "" && manifest.BundleSha256 != computedSHA {
		warns = append(warns, fmt.Sprintf("manifest.bundleSha256 %q does not match recomputed canonical digest %q", manifest.BundleSha256, computedSHA))
	}
	if digestLine != "" && digestLine != computedSHA {
		warns = append(warns, fmt.Sprintf("digest file %q does not match recomputed canonical digest %q", digestLine, computedSHA))
	}
	if manifest.BundleBytes > 0 && int64(len(canonical)) != int64(manifest.BundleBytes) {
		warns = append(warns, fmt.Sprintf("manifest.bundleBytes=%d differs from len(canonical)=%d", manifest.BundleBytes, len(canonical)))
	}

	digestMatch := digestLine == computedSHA && (manifest.BundleSha256 == "" || manifest.BundleSha256 == computedSHA)

	var bundleObj interface{}
	if err := json.Unmarshal(bundleRaw, &bundleObj); err != nil {
		return nil, err
	}
	var manifestObj interface{}
	if err := json.Unmarshal(manifestRaw, &manifestObj); err != nil {
		return nil, err
	}

	ns := strings.TrimSpace(o.Namespace)
	if ns == "" {
		if bundle.CellRef != nil && strings.TrimSpace(bundle.CellRef.Namespace) != "" {
			ns = bundle.CellRef.Namespace
		} else {
			ns = "karl-sandbox"
		}
	}

	name := strings.TrimSpace(o.Name)
	if name == "" {
		base := strings.TrimSpace(manifest.ArtifactID)
		if base == "" {
			base = "ingest"
		}
		name = sanitizeName("khr-" + base + "-" + computedSHA[:12])
	}

	reqDigest := digestLine
	if reqDigest == "" {
		reqDigest = computedSHA
	}

	maxBytes := int64(len(canonical)) + 65536
	if int64(manifest.BundleBytes) > int64(len(canonical)) {
		maxBytes = int64(manifest.BundleBytes) + 65536
	}

	artifactID := strings.TrimSpace(manifest.ArtifactID)

	allowLocalDev := o.AllowLocalDevSignature
	if strings.TrimSpace(manifest.SigningMode) == "local-dev" {
		allowLocalDev = true
	}

	doc := map[string]interface{}{
		"apiVersion": "hyperdensity.karl.io/v1alpha1",
		"kind":       "EvidenceIngestRequest",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": ns,
			"labels": map[string]interface{}{
				"karl.io/source": "khr-linux-agent-prepare-ingest-request",
			},
		},
		"spec": map[string]interface{}{
			"artifactId": artifactID,
			"bundle":     bundleObj,
			"manifest":   manifestObj,
			"digest":     reqDigest,
			"source": map[string]interface{}{
				"agentId":  firstNonEmpty(manifest.AgentID, bundle.AgentID),
				"nodeName": strings.TrimSpace(o.NodeName),
				"hostId":   strings.TrimSpace(o.HostID),
				"tenant":   strings.TrimSpace(o.Tenant),
			},
			"policy": map[string]interface{}{
				"requireDigestMatch":     o.RequireDigestMatch,
				"allowUnsigned":          o.AllowUnsigned,
				"allowLocalDevSignature": allowLocalDev,
				"maxBundleBytes":         maxBytes,
			},
			"dryRunOnly": o.DryRunOnly,
			"ttlSeconds": int64(86400),
		},
		"status": map[string]interface{}{
			"phase":             "Pending",
			"digestMatch":       digestMatch,
			"signatureStatus":   signatureStatusFromManifest(&manifest),
			"rejectionReasons":  []string{},
			"evidenceBundleRef": nil,
		},
	}

	meta := doc["metadata"].(map[string]interface{})
	ann := map[string]interface{}{}
	if len(warns) > 0 {
		b, _ := json.Marshal(warns)
		ann[annotationPreparationWarnings] = string(b)
	}
	if strings.TrimSpace(manifest.SigningMode) == "local-dev" {
		ann[annotationSignatureTrustTier] = trustTierDevOnly
	}
	if len(ann) > 0 {
		meta["annotations"] = ann
	}

	return marshalDoc(doc, o.Format)
}

func marshalDoc(doc map[string]interface{}, format string) ([]byte, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "yaml", "yml":
		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf)
		enc.SetIndent(2)
		if err := enc.Encode(doc); err != nil {
			return nil, err
		}
		if err := enc.Close(); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case "json":
		b, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return nil, err
		}
		return append(b, '\n'), nil
	default:
		return nil, fmt.Errorf("unsupported ingest-request-format %q (yaml|json)", format)
	}
}

func signatureStatusFromManifest(m *integrity.ArtifactManifest) string {
	if m == nil {
		return "None"
	}
	switch strings.TrimSpace(m.SigningMode) {
	case "local-dev":
		if m.SignaturePresent {
			return "DevOnly"
		}
		return "Unsigned"
	case "none", "":
		if m.SignaturePresent {
			return "Invalid"
		}
		return "None"
	default:
		if m.SignaturePresent {
			return "Pending"
		}
		return "None"
	}
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

func sanitizeName(s string) string {
	var out strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			out.WriteRune(r)
		case r == '-' || r == '.':
			out.WriteRune(r)
		default:
			out.WriteRune('-')
		}
	}
	res := strings.Trim(out.String(), "-.")
	if res == "" {
		return "khr-ingest"
	}
	if len(res) > 253 {
		res = res[:253]
	}
	return res
}

// DefaultPrepareOptions returns conservative defaults for local file generation.
func DefaultPrepareOptions() PrepareOptions {
	return PrepareOptions{
		Format:                 "yaml",
		RequireDigestMatch:     true,
		AllowUnsigned:          true,
		AllowLocalDevSignature: false,
		DryRunOnly:             false,
	}
}

// WriteFile creates parent dirs and writes bytes with mode 0o600.
func WriteFile(path string, data []byte) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("output path is empty")
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0o600)
}
