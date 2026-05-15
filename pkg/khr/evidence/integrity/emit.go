package integrity

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// EmitEvidenceSidecars writes digest and/or manifest for a bundle value v (typically *evidence.CollectEvidenceBundle).
// Canonical JSON is always the basis for bundleSha256 and bundleBytes.
func EmitEvidenceSidecars(v any, agentID, artifactID, manifestPath, digestPath, signingMode, keyFile string) error {
	mode := NormalizeSigningMode(signingMode)
	if err := ValidateSigningMode(signingMode); err != nil {
		return err
	}
	if err := RequireLocalDevKey(mode, keyFile); err != nil {
		return err
	}
	mp := strings.TrimSpace(manifestPath)
	dp := strings.TrimSpace(digestPath)
	if mode == "local-dev" && mp == "" {
		return fmt.Errorf("collect-evidence: signing-mode=local-dev requires -evidence-manifest-output (signature metadata is stored only in the manifest)")
	}
	if mp == "" && dp == "" {
		return nil
	}

	canonical, err := CanonicalJSON(v)
	if err != nil {
		return fmt.Errorf("canonical json: %w", err)
	}
	sha := SHA256Hex(canonical)

	var sigB64, sigAlg string
	if mode == "local-dev" {
		sig, err := SignLocalDev(canonical, keyFile)
		if err != nil {
			return fmt.Errorf("local-dev sign: %w", err)
		}
		sigB64 = base64.StdEncoding.EncodeToString(sig)
		sigAlg = SignatureAlgorithmLocalDev
	}

	m := BuildManifest(agentID, artifactID, mode, canonical, sha, sigB64, sigAlg)

	if dp != "" {
		if err := WriteDigestFile(dp, sha); err != nil {
			return err
		}
	}
	if mp != "" {
		if err := WriteManifestFile(mp, m); err != nil {
			return err
		}
	}
	return nil
}
