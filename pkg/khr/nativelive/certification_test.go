package nativelive

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func passRun() RunMetrics {
	return RunMetrics{
		RestartCountBefore:   0,
		RestartCountAfter:    0,
		RolloutDetected:      false,
		RecreateDetected:     false,
		InterruptionDetected: false,
		InterruptionWindowMs: 0,
		NativeLiveLaneCount:  1,
		LiveInPlaceEligible:  true,
		RollbackPass:         true,
		ApplyLatencyMs:       ApplyLatencyMs{CPU: 10, RAMUp: 8, RAMDown: 7},
		RollbackLatencyMs:    RollbackLatencyMs{CPU: 5, RAMUp: 4, RAMDown: 3},
	}
}

func TestAggregateRuns_DeterministicFingerprint(t *testing.T) {
	r1 := passRun()
	r2 := passRun()
	s1 := AggregateRuns("KHR-T", []RunMetrics{r1})
	s2 := AggregateRuns("KHR-T", []RunMetrics{r2})
	if s1.RunFingerprints[0].Fingerprint != s2.RunFingerprints[0].Fingerprint {
		t.Fatalf("fingerprints differ: %s vs %s", s1.RunFingerprints[0].Fingerprint, s2.RunFingerprints[0].Fingerprint)
	}
}

func TestCheckRegression_NoInterruptionInvariant(t *testing.T) {
	summary := AggregateRuns("KHR-T", []RunMetrics{passRun()})
	if err := CheckRegression(summary); err != nil {
		t.Fatal(err)
	}
	if summary.Invariants.InterruptionWindowMs != 0 {
		t.Fatalf("interruptionWindowMs=%d", summary.Invariants.InterruptionWindowMs)
	}
}

func TestCheckRegression_RestartDetected(t *testing.T) {
	r := passRun()
	r.RestartCountBefore = 0
	r.RestartCountAfter = 1
	summary := AggregateRuns("KHR-T", []RunMetrics{r})
	if err := CheckRegression(summary); err == nil {
		t.Fatal("expected regression on restart")
	}
	if !summary.RegressionDetected {
		t.Fatal("regressionDetected=false")
	}
}

func TestCheckRegression_RolloutDetected(t *testing.T) {
	r := passRun()
	r.RolloutDetected = true
	summary := AggregateRuns("KHR-T", []RunMetrics{r})
	if err := CheckRegression(summary); err == nil {
		t.Fatal("expected regression on rollout")
	}
}

func TestCheckRegression_RecreateDetected(t *testing.T) {
	r := passRun()
	r.RecreateDetected = true
	summary := AggregateRuns("KHR-T", []RunMetrics{r})
	if err := CheckRegression(summary); err == nil {
		t.Fatal("expected regression on recreate")
	}
}

func TestCheckRegression_InterruptionWindow(t *testing.T) {
	r := passRun()
	r.InterruptionWindowMs = 1
	summary := AggregateRuns("KHR-T", []RunMetrics{r})
	if err := CheckRegression(summary); err == nil {
		t.Fatal("expected regression on interruption window")
	}
}

func TestFingerprintsMatch_RepeatableRuns(t *testing.T) {
	summary := AggregateRuns("KHR-T", []RunMetrics{passRun(), passRun()})
	if !FingerprintsMatch(summary) {
		t.Fatal("expected matching fingerprints across repeatable runs")
	}
}

func TestCompareBaseline_Pass(t *testing.T) {
	summary := AggregateRuns("KHR-T", []RunMetrics{passRun()})
	baseline := BaselineCertification{
		CertificationID: CertificationID,
		Lane:            LaneNativeLive,
		Status:          CertificationCertified,
		Invariants: Invariants{
			NoRestart: true, NoRollout: true, NoRecreate: true,
			InterruptionWindowMs: 0, InterruptionDetected: false,
		},
		Scores:           CertificationScores{ContinuityScore: 1.0, LiveScaleConfidence: ConfidenceHigh},
		ExpectedRunCount: 1,
	}
	match, diffs := CompareBaseline(summary, baseline)
	if !match {
		t.Fatalf("baseline mismatch: %v", diffs)
	}
}

func TestCertificationSummary_GoldenFixture(t *testing.T) {
	path := filepath.Join("testdata", "certification-summary-pass.json")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var summary CertificationSummary
	if err := json.Unmarshal(b, &summary); err != nil {
		t.Fatal(err)
	}
	if err := CheckRegression(summary); err != nil {
		t.Fatal(err)
	}
	if summary.Scores.ContinuityScore < 1.0 {
		t.Fatalf("continuityScore=%v", summary.Scores.ContinuityScore)
	}
}
