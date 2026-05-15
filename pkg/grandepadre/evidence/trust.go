package evidence

import (
	"strings"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence/integrity"
)

const annotationSignatureTrustTier = "khr.karl.io/signature-trust-tier"

// TrustInputs drives ComputeTrustTier (integrity semantics only; not apply authorization).
type TrustInputs struct {
	DigestMatch              bool
	RequireDigestMatch       bool
	Manifest                 *integrity.ArtifactManifest
	AllowUnsigned            bool
	AllowLocalDevSignature   bool
	MetadataAnnotations      map[string]string
	UnsignedDigestTrustLabel UnsignedDigestTrustPolicy
}

// ComputeTrustTier returns a trust tier for indexing. Digest mismatch always yields IntegrityFailed.
func ComputeTrustTier(in TrustInputs) TrustTier {
	if !in.DigestMatch {
		return TrustIntegrityFailed
	}
	if in.MetadataAnnotations != nil {
		if strings.TrimSpace(in.MetadataAnnotations[annotationSignatureTrustTier]) == "DevOnly" {
			return TrustDevOnly
		}
	}
	mode := "none"
	if in.Manifest != nil {
		mode = integrity.NormalizeSigningMode(in.Manifest.SigningMode)
	}
	if mode == "local-dev" {
		return TrustDevOnly
	}
	if in.Manifest != nil && in.Manifest.SignaturePresent && mode != "none" && mode != "local-dev" {
		return TrustUnknown
	}
	// Digest matched; no production PKI in this skeleton.
	if in.UnsignedDigestTrustLabel == UnsignedDigestAsUnsigned {
		return TrustUnsigned
	}
	return TrustIntegrityVerified
}
