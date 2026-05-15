package recommendation

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	gpevidence "github.com/marco-simoncini/Karl-Hyperdensity/pkg/grandepadre/evidence"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func fixedClock(t *testing.T, ts time.Time) {
	t.Helper()
	prev := gpevidence.NowFunc
	gpevidence.NowFunc = func() time.Time { return ts.UTC() }
	t.Cleanup(func() { gpevidence.NowFunc = prev })
}

func loadStore(t *testing.T, relPaths ...string) *gpevidence.Store {
	t.Helper()
	s := gpevidence.NewStore()
	var files []string
	root := repoRoot(t)
	for _, p := range relPaths {
		files = append(files, filepath.Join(root, p))
	}
	if err := IngestAllIntoStore(s, files, gpevidence.UnsignedDigestAsIntegrityVerified); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestReadyHighObserveAndPrepareLease(t *testing.T) {
	fixedClock(t, time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC))
	s := loadStore(t, "examples/grandepadre/evidence-store/ingest-request-ready.yaml")
	rep := BuildRecommendationReport(s, EngineOptions{DryRunOnly: true, UnsignedLabel: gpevidence.UnsignedDigestAsIntegrityVerified})
	if rep.IndexedCount != 1 {
		t.Fatalf("indexed=%d", rep.IndexedCount)
	}
	var sawObserve, sawLease bool
	for _, r := range rep.ActionSlate.Recommendations {
		if r.ActionType == ActionObserve {
			sawObserve = true
		}
		if r.ActionType == ActionPrepareResourceLease {
			sawLease = true
		}
	}
	if !sawObserve || !sawLease {
		t.Fatalf("recommendations=%+v", rep.ActionSlate.Recommendations)
	}
	if len(rep.ActionSlate.DonorCandidates) != 1 {
		t.Fatalf("donors=%d", len(rep.ActionSlate.DonorCandidates))
	}
}

func TestBlockedRemediate(t *testing.T) {
	fixedClock(t, time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC))
	s := loadStore(t, "examples/grandepadre/evidence-store/ingest-request-blocked.yaml")
	rep := BuildRecommendationReport(s, EngineOptions{DryRunOnly: true, UnsignedLabel: gpevidence.UnsignedDigestAsIntegrityVerified})
	var sawRemediate bool
	for _, r := range rep.ActionSlate.Recommendations {
		if r.ActionType == ActionRemediate && r.Risk != RiskBlocked {
			sawRemediate = true
		}
	}
	if !sawRemediate {
		t.Fatalf("recommendations=%+v", rep.ActionSlate.Recommendations)
	}
	if len(rep.ActionSlate.ReceiverCandidates) == 0 {
		t.Fatal("expected receiver candidates")
	}
}

func TestLowConfidenceCollectMore(t *testing.T) {
	fixedClock(t, time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC))
	s := loadStore(t, "examples/grandepadre/recommendation/recommendation-input-low-confidence.yaml")
	rep := BuildRecommendationReport(s, EngineOptions{DryRunOnly: true, UnsignedLabel: gpevidence.UnsignedDigestAsIntegrityVerified})
	var sawCollect bool
	for _, r := range rep.ActionSlate.Recommendations {
		if r.ActionType == ActionCollectMoreEvidence {
			sawCollect = true
		}
	}
	if !sawCollect {
		t.Fatalf("recommendations=%+v", rep.ActionSlate.Recommendations)
	}
}

func TestIntegrityFailedRiskBlocked(t *testing.T) {
	fixedClock(t, time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC))
	root := repoRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, "examples/grandepadre/evidence-store/ingest-request-ready.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	corrupt := regexp.MustCompile(`(?m)^(\s*digest:\s*)[0-9a-f]{64}\s*$`).
		ReplaceAllString(string(raw), "${1}0000000000000000000000000000000000000000000000000000000000000000")
	s := gpevidence.NewStore()
	if _, err := s.Ingest([]byte(corrupt), gpevidence.UnsignedDigestAsIntegrityVerified); err != nil {
		t.Fatal(err)
	}
	rep := BuildRecommendationReport(s, EngineOptions{DryRunOnly: true, UnsignedLabel: gpevidence.UnsignedDigestAsIntegrityVerified})
	var sawBlocked bool
	for _, r := range rep.ActionSlate.Recommendations {
		if r.Risk == RiskBlocked {
			sawBlocked = true
		}
	}
	if !sawBlocked {
		t.Fatalf("recommendations=%+v", rep.ActionSlate.Recommendations)
	}
}

func TestDevOnlyNonProdObserve(t *testing.T) {
	fixedClock(t, time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC))
	s := loadStore(t, "examples/grandepadre/recommendation/recommendation-input-devonly.yaml")
	rep := BuildRecommendationReport(s, EngineOptions{DryRunOnly: true, UnsignedLabel: gpevidence.UnsignedDigestAsIntegrityVerified})
	var sawDev bool
	for _, r := range rep.ActionSlate.Recommendations {
		for _, reason := range r.Reasons {
			if strings.Contains(reason, "DevOnly") || strings.Contains(strings.ToLower(reason), "non-production") {
				sawDev = true
			}
		}
		for _, p := range r.Prerequisites {
			if strings.Contains(p, "DevOnly") {
				sawDev = true
			}
		}
	}
	if !sawDev {
		t.Fatalf("recommendations=%+v", rep.ActionSlate.Recommendations)
	}
	if len(rep.ActionSlate.DonorCandidates) != 0 {
		t.Fatalf("DevOnly rows must not become donor candidates, got %v", rep.ActionSlate.DonorCandidates)
	}
}

func TestDonorReceiverAndSummary(t *testing.T) {
	fixedClock(t, time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC))
	s := loadStore(t,
		"examples/grandepadre/evidence-store/ingest-request-ready.yaml",
		"examples/grandepadre/evidence-store/ingest-request-blocked.yaml",
	)
	rep := BuildRecommendationReport(s, EngineOptions{DryRunOnly: true, UnsignedLabel: gpevidence.UnsignedDigestAsIntegrityVerified})
	if rep.ActionSlate.Summary.RecommendationCount != len(rep.ActionSlate.Recommendations) {
		t.Fatalf("summary mismatch %d vs %d", rep.ActionSlate.Summary.RecommendationCount, len(rep.ActionSlate.Recommendations))
	}
	if rep.ActionSlate.Summary.DonorCandidates != len(rep.ActionSlate.DonorCandidates) {
		t.Fatal("summary donor count")
	}
	if rep.ActionSlate.Summary.ReceiverCandidates != len(rep.ActionSlate.ReceiverCandidates) {
		t.Fatal("summary receiver count")
	}
}

func TestTenantFilter(t *testing.T) {
	fixedClock(t, time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC))
	s := loadStore(t, "examples/grandepadre/evidence-store/ingest-request-ready.yaml")
	rep := BuildRecommendationReport(s, EngineOptions{
		Tenant:        "other-ns",
		DryRunOnly:    true,
		UnsignedLabel: gpevidence.UnsignedDigestAsIntegrityVerified,
	})
	if rep.IndexedCount != 0 {
		t.Fatalf("indexed=%d", rep.IndexedCount)
	}
}
