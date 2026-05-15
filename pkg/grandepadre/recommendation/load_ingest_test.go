package recommendation

import (
	"path/filepath"
	"testing"

	gpevidence "github.com/marco-simoncini/Karl-Hyperdensity/pkg/grandepadre/evidence"
)

func TestCollectIngestPathsDirAndFiles(t *testing.T) {
	root := repoRoot(t)
	dir := filepath.Join(root, "examples", "grandepadre", "recommendation")
	p1 := filepath.Join(dir, "recommendation-input-ready.yaml")
	paths, err := CollectIngestPaths([]string{p1}, dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) < 2 {
		t.Fatalf("expected dir+file merge, got %v", paths)
	}
}

func TestDonorCandidatesSkipsNonHighAndDevOnly(t *testing.T) {
	ref := &gpevidence.CellRefLite{
		APIVersion: "runtime.karl.io/v1alpha1",
		Kind:       "Cell",
		Namespace:  "karl-sandbox",
		Name:       "demo-cell",
	}
	indices := []gpevidence.EvidenceIndex{
		{CellRef: ref, Confidence: "high", ReadyForGrandePadre: true, TrustTier: gpevidence.TrustIntegrityVerified},
		{CellRef: ref, Confidence: "medium", ReadyForGrandePadre: true, TrustTier: gpevidence.TrustIntegrityVerified},
		{CellRef: ref, Confidence: "high", ReadyForGrandePadre: true, TrustTier: gpevidence.TrustDevOnly},
	}
	d := DonorCandidates(indices)
	if len(d) != 1 {
		t.Fatalf("donors=%v", d)
	}
}
