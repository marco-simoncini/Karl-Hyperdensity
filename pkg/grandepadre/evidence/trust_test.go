package evidence

import (
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence/integrity"
)

func TestComputeTrustTierDigestFailed(t *testing.T) {
	tier := ComputeTrustTier(TrustInputs{DigestMatch: false})
	if tier != TrustIntegrityFailed {
		t.Fatal(tier)
	}
}

func TestComputeTrustTierAnnotationDevOnly(t *testing.T) {
	tier := ComputeTrustTier(TrustInputs{
		DigestMatch: true,
		MetadataAnnotations: map[string]string{
			annotationSignatureTrustTier: "DevOnly",
		},
	})
	if tier != TrustDevOnly {
		t.Fatal(tier)
	}
}

func TestComputeTrustTierUnknownSignatureMode(t *testing.T) {
	m := &integrity.ArtifactManifest{
		SigningMode:       "hsm-pki",
		SignaturePresent:  true,
		BundleSha256:      "a",
	}
	tier := ComputeTrustTier(TrustInputs{DigestMatch: true, Manifest: m})
	if tier != TrustUnknown {
		t.Fatal(tier)
	}
}
