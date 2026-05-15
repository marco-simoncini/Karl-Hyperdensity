package evidence

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/discovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence/integrity"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/telemetry"
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

func fixedNow(t *testing.T) {
	t.Helper()
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	NowFunc = func() time.Time { return now }
	t.Cleanup(func() {
		NowFunc = func() time.Time { return time.Now().UTC() }
	})
}

func TestIngestAcceptedFromExampleFixture(t *testing.T) {
	fixedNow(t)
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "examples", "grandepadre", "evidence-store", "ingest-request-ready.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	s := NewStore()
	if _, err := s.Ingest(raw, UnsignedDigestAsIntegrityVerified); err != nil {
		t.Fatal(err)
	}
	s.DeduplicateBySha256()
	if s.Len() != 1 {
		t.Fatalf("len=%d", s.Len())
	}
	ready := s.ListReady()
	if len(ready) != 1 {
		t.Fatalf("ready=%v", ready)
	}
	if ready[0].TrustTier != TrustIntegrityVerified {
		t.Fatalf("trust=%s", ready[0].TrustTier)
	}
}

func TestDigestMismatchIntegrityFailed(t *testing.T) {
	fixedNow(t)
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "examples", "grandepadre", "evidence-store", "ingest-request-ready.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	digestCorrupt := regexp.MustCompile(`(?m)^(\s*digest:\s*)[0-9a-f]{64}\s*$`).
		ReplaceAllString(string(raw), "${1}0000000000000000000000000000000000000000000000000000000000000000")
	s := NewStore()
	if _, err := s.Ingest([]byte(digestCorrupt), UnsignedDigestAsIntegrityVerified); err != nil {
		t.Fatal(err)
	}
	idx := s.Snapshot()[0]
	if idx.TrustTier != TrustIntegrityFailed {
		t.Fatalf("trust=%s", idx.TrustTier)
	}
	if len(s.ListReady()) != 0 {
		t.Fatal("expected no ready rows")
	}
}

func TestLocalDevDevOnly(t *testing.T) {
	fixedNow(t)
	t.Setenv("KHR_TEST_COLLECTED_AT", "2026-05-15T16:00:00Z")
	t.Setenv("KHR_TEST_TELEMETRY_NOW", "2026-05-15T15:00:00Z")
	cell := &crdv1alpha1.Cell{
		APIVersion: "runtime.karl.io/v1alpha1",
		Kind:       "Cell",
		Metadata: crdv1alpha1.ObjectMeta{Name: "c", Namespace: "ns"},
	}
	disc := &discovery.CgroupDiscoveryOutput{
		AgentID: "a", CgroupVersion: "v2", DiscoveryMode: "read-only",
		ScannedRoot: "/r", SelectedPath: "/r/s", MutationsForbidden: true,
	}
	tel := evidence.TelemetrySnapshot{
		TelemetryMode: "read-only", CgroupPath: "/r/s",
		Metrics: telemetry.MetricsBundle{CPUStat: map[string]int64{"usage_usec": 1}, MemoryCurrent: "1"},
		Evidence: telemetry.Evidence{
			ObservedAt: "2026-05-15T15:00:00Z", Source: "cgroup-v2", Confidence: "medium",
		},
		MutationsForbidden: true,
	}
	b := evidence.BuildCollectEvidenceBundle("0.0.1-sprint13", "a", cell, disc, tel, evidence.DryRunSkippedPayload("skip"), "")
	canonical, err := integrity.CanonicalJSON(b)
	if err != nil {
		t.Fatal(err)
	}
	sha := integrity.SHA256Hex(canonical)
	man := integrity.BuildManifest("a", "art", "local-dev", canonical, sha, "sig", "ed25519")
	doc := &IngestDocument{
		SpecArtifactID:     "art",
		Bundle:             *b,
		Manifest:           man,
		DigestLine:         sha,
		RequireDigestMatch: true,
		AllowUnsigned:      true,
		AllowLocalDevSig:   true,
	}
	idx, err := BuildEvidenceIndex(doc, UnsignedDigestAsIntegrityVerified)
	if err != nil {
		t.Fatal(err)
	}
	if idx.TrustTier != TrustDevOnly {
		t.Fatalf("trust=%s", idx.TrustTier)
	}
}

func TestUnsignedDigestPolicyIntegrityVerifiedVsUnsigned(t *testing.T) {
	fixedNow(t)
	t.Setenv("KHR_TEST_COLLECTED_AT", "2026-05-15T16:00:00Z")
	t.Setenv("KHR_TEST_TELEMETRY_NOW", "2026-05-15T15:00:00Z")
	cell := &crdv1alpha1.Cell{
		APIVersion: "runtime.karl.io/v1alpha1",
		Kind:       "Cell",
		Metadata: crdv1alpha1.ObjectMeta{Name: "c", Namespace: "ns"},
	}
	disc := &discovery.CgroupDiscoveryOutput{
		AgentID: "a", CgroupVersion: "v2", DiscoveryMode: "read-only",
		ScannedRoot: "/r", SelectedPath: "/r/s", MutationsForbidden: true,
	}
	tel := evidence.TelemetrySnapshot{
		TelemetryMode: "read-only", CgroupPath: "/r/s",
		Metrics: telemetry.MetricsBundle{CPUStat: map[string]int64{"usage_usec": 2}, MemoryCurrent: "2"},
		Evidence: telemetry.Evidence{
			ObservedAt: "2026-05-15T15:00:00Z", Source: "cgroup-v2", Confidence: "low",
		},
		MutationsForbidden: true,
	}
	b := evidence.BuildCollectEvidenceBundle("0.0.1-sprint13", "a", cell, disc, tel, evidence.DryRunSkippedPayload("skip"), "")
	canonical, err := integrity.CanonicalJSON(b)
	if err != nil {
		t.Fatal(err)
	}
	sha := integrity.SHA256Hex(canonical)
	man := integrity.BuildManifest("a", "art2", "none", canonical, sha, "", "")
	doc := &IngestDocument{
		SpecArtifactID:     "art2",
		Bundle:             *b,
		Manifest:           man,
		DigestLine:         sha,
		RequireDigestMatch: true,
		AllowUnsigned:      true,
	}
	idx, err := BuildEvidenceIndex(doc, UnsignedDigestAsIntegrityVerified)
	if err != nil {
		t.Fatal(err)
	}
	if idx.TrustTier != TrustIntegrityVerified {
		t.Fatalf("default trust=%s", idx.TrustTier)
	}
	idx2, err := BuildEvidenceIndex(doc, UnsignedDigestAsUnsigned)
	if err != nil {
		t.Fatal(err)
	}
	if idx2.TrustTier != TrustUnsigned {
		t.Fatalf("unsigned label trust=%s", idx2.TrustTier)
	}
}

func TestDuplicateSha256Deduplicates(t *testing.T) {
	fixedNow(t)
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "examples", "grandepadre", "evidence-store", "ingest-request-ready.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	s := NewStore()
	if _, err := s.Ingest(raw, UnsignedDigestAsIntegrityVerified); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Ingest(raw, UnsignedDigestAsIntegrityVerified); err != nil {
		t.Fatal(err)
	}
	if s.DuplicateTotal() != 1 {
		t.Fatalf("dup=%d", s.DuplicateTotal())
	}
	removed := s.DeduplicateBySha256()
	if removed != 1 || s.Len() != 1 {
		t.Fatalf("removed=%d len=%d", removed, s.Len())
	}
}

func TestQueryReadyBlockedByCell(t *testing.T) {
	fixedNow(t)
	root := repoRoot(t)
	readyRaw, err := os.ReadFile(filepath.Join(root, "examples", "grandepadre", "evidence-store", "ingest-request-ready.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	blockedRaw, err := os.ReadFile(filepath.Join(root, "examples", "grandepadre", "evidence-store", "ingest-request-blocked.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	s := NewStore()
	if _, err := s.Ingest(readyRaw, UnsignedDigestAsIntegrityVerified); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Ingest(blockedRaw, UnsignedDigestAsIntegrityVerified); err != nil {
		t.Fatal(err)
	}
	s.DeduplicateBySha256()
	if len(s.ListReady()) != 1 {
		t.Fatalf("ready n=%d", len(s.ListReady()))
	}
	if len(s.ListBlocked()) != 1 {
		t.Fatalf("blocked n=%d", len(s.ListBlocked()))
	}
	cellHits := s.GetByCell("karl-sandbox", "demo-cell")
	if len(cellHits) != 2 {
		t.Fatalf("cell hits=%d", len(cellHits))
	}
}

func TestRunLocalIndexQueryByConfidence(t *testing.T) {
	fixedNow(t)
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "examples", "grandepadre", "evidence-store", "ingest-request-ready.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	s := NewStore()
	rep, err := RunLocalIndex(s, raw, LocalIndexParams{
		UnsignedLabel: UnsignedDigestAsIntegrityVerified,
		Query:         QueryByConfidence,
		Confidence:    "high",
	})
	if err != nil {
		t.Fatal(err)
	}
	qr, ok := rep.QueryResult.([]EvidenceIndex)
	if !ok || len(qr) != 1 {
		t.Fatalf("queryResult=%T %#v", rep.QueryResult, rep.QueryResult)
	}
}

func TestBuildBlockedRemediableIndex(t *testing.T) {
	fixedNow(t)
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "examples", "grandepadre", "evidence-store", "ingest-request-blocked.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	s := NewStore()
	if _, err := s.Ingest(raw, UnsignedDigestAsIntegrityVerified); err != nil {
		t.Fatal(err)
	}
	br := BuildBlockedRemediableIndex(s.Snapshot())
	if len(br) != 1 {
		t.Fatalf("remediable groups=%d", len(br))
	}
	if len(br[0].RemediationHints) == 0 {
		t.Fatal("expected remediation hints")
	}
	if br[0].LastArtifactID == "" {
		t.Fatal("last artifact")
	}
}
