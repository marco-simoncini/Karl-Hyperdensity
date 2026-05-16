package nativelive

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/shellcontinuity"
)

const (
	CertificationID      = "khr-native-live-certification-v1"
	LaneNativeLive       = "native-live"
	CertificationCertified = "certified"
	CertificationFailed    = "failed"

	ConfidenceHigh   = "high"
	ConfidenceMedium = "medium"
	ConfidenceLow    = "low"
)

// ApplyLatencyMs records per-step apply durations (milliseconds).
type ApplyLatencyMs struct {
	CPU    int64 `json:"cpu,omitempty"`
	RAMUp  int64 `json:"ramUp,omitempty"`
	RAMDown int64 `json:"ramDown,omitempty"`
}

// RollbackLatencyMs records per-step rollback durations (milliseconds).
type RollbackLatencyMs struct {
	RAMDown int64 `json:"ramDown,omitempty"`
	RAMUp   int64 `json:"ramUp,omitempty"`
	CPU     int64 `json:"cpu,omitempty"`
}

// RunMetrics is one native-live evidence run (cluster-observed, read-only).
type RunMetrics struct {
	RunIndex              int   `json:"runIndex,omitempty"`
	RestartCountBefore    int64 `json:"restartCountBefore"`
	RestartCountAfter     int64 `json:"restartCountAfter"`
	RestartCountDelta     int64 `json:"restartCountDelta"`
	RolloutCount          int64 `json:"rolloutCount"`
	RolloutDetected       bool  `json:"rolloutDetected"`
	RecreateDetected      bool  `json:"recreateDetected"`
	InterruptionDetected  bool  `json:"interruptionDetected"`
	InterruptionWindowMs  int64 `json:"interruptionWindowMs"`
	ApplyLatencyMs        ApplyLatencyMs    `json:"applyLatencyMs"`
	RollbackLatencyMs     RollbackLatencyMs `json:"rollbackLatencyMs"`
	NativeLiveLaneCount   int   `json:"nativeLiveLaneCount"`
	LiveInPlaceEligible            bool                    `json:"liveInPlaceEligible"`
	RollbackPass                   bool                    `json:"rollbackPass"`
	ShellContinuityPreserved       bool                    `json:"shellContinuityPreserved"`
	AppContinuityPreserved         bool                    `json:"appContinuityPreserved"`
	UserSessionContinuityPreserved bool                    `json:"userSessionContinuityObserved"`
	SessionContinuityPreserved     bool                    `json:"sessionContinuityPreserved"`
	ContinuityState                string                  `json:"continuityState,omitempty"`
	ContinuityEvidence             shellcontinuity.Evidence `json:"continuityEvidence,omitempty"`
}

// Invariants is the no-interruption contract for native-live certification.
type Invariants struct {
	NoRestart              bool  `json:"noRestart"`
	NoRollout              bool  `json:"noRollout"`
	NoRecreate             bool  `json:"noRecreate"`
	InterruptionWindowMs   int64 `json:"interruptionWindowMs"`
	InterruptionDetected   bool  `json:"interruptionDetected"`
}

// CertificationMetrics aggregates observability across runs.
type CertificationMetrics struct {
	RestartCountDelta    int64             `json:"restartCountDelta"`
	RolloutCount         int64             `json:"rolloutCount"`
	RecreateDetected     bool              `json:"recreateDetected"`
	ApplyLatencyMs       ApplyLatencyMs    `json:"applyLatencyMs"`
	RollbackLatencyMs    RollbackLatencyMs `json:"rollbackLatencyMs"`
	InterruptionWindowMs int64             `json:"interruptionWindowMs"`
	RunCount             int               `json:"runCount"`
}

// ContinuityCertificationProof groups KHR-U shell/session/app continuity proofs.
type ContinuityCertificationProof struct {
	ResourceContinuityPreserved bool `json:"resourceContinuityPreserved"`
	ShellContinuityPreserved    bool `json:"shellContinuityPreserved"`
	AppContinuityPreserved      bool `json:"appContinuityPreserved"`
	SessionContinuityPreserved  bool `json:"sessionContinuityPreserved"`
}

// CertificationScores are read-only continuity/confidence signals (KHR-T/U).
type CertificationScores struct {
	ResourceContinuityScore float64 `json:"resourceContinuityScore"`
	SessionContinuityScore  float64 `json:"sessionContinuityScore"`
	ContinuityScore         float64 `json:"continuityScore"`
	LiveScaleConfidence     string  `json:"liveScaleConfidence"`
}

// RunFingerprint is a deterministic, volatility-stripped run signature.
type RunFingerprint struct {
	RunIndex    int    `json:"runIndex"`
	Fingerprint string `json:"fingerprint"`
}

// CertificationSummary is the native-live certification JSON artifact.
type CertificationSummary struct {
	CertificationID      string               `json:"certificationId"`
	Sprint               string               `json:"sprint"`
	Lane                 string               `json:"lane"`
	Status               string               `json:"status"`
	RegressionDetected   bool                 `json:"regressionDetected"`
	RegressionReasons    []string             `json:"regressionReasons,omitempty"`
	Invariants           Invariants           `json:"invariants"`
	Metrics              CertificationMetrics `json:"metrics"`
	Scores               CertificationScores  `json:"scores"`
	RunFingerprints      []RunFingerprint            `json:"runFingerprints"`
	ContinuityProof      ContinuityCertificationProof  `json:"continuityProof"`
	BaselineMatch        bool                          `json:"baselineMatch"`
	BaselineDiff         []string             `json:"baselineDiff,omitempty"`
	ReadOnly             bool                 `json:"readOnly"`
	NoAutomation         bool                 `json:"noAutomation"`
	NoAutonomousOrchestration bool            `json:"noAutonomousOrchestration"`
}

// BaselineCertification is the committed baseline for compare (deterministic fields only).
type BaselineCertification struct {
	CertificationID    string     `json:"certificationId"`
	Lane               string     `json:"lane"`
	Status             string     `json:"status"`
	Invariants         Invariants `json:"invariants"`
	Scores             CertificationScores `json:"scores"`
	ExpectedRunCount   int        `json:"expectedRunCount,omitempty"`
}

// NormalizeRunMetrics fills derived fields on a run snapshot.
func NormalizeRunMetrics(m *RunMetrics) {
	if m == nil {
		return
	}
	m.RestartCountDelta = m.RestartCountAfter - m.RestartCountBefore
	if m.RestartCountDelta < 0 {
		m.RestartCountDelta = 0
	}
	if m.RolloutDetected && m.RolloutCount == 0 {
		m.RolloutCount = 1
	}
}

// AggregateRuns builds a certification summary from one or more runs.
func AggregateRuns(sprint string, runs []RunMetrics) CertificationSummary {
	for i := range runs {
		runs[i].RunIndex = i + 1
		NormalizeRunMetrics(&runs[i])
	}
	inv := aggregateInvariants(runs)
	metrics := aggregateMetrics(runs)
	scores := ScoreCertification(inv, metrics, runs)
	reasons := RegressionReasons(runs, inv)
	regression := len(reasons) > 0
	status := CertificationCertified
	if regression {
		status = CertificationFailed
	}
	fps := make([]RunFingerprint, 0, len(runs))
	for _, r := range runs {
		fps = append(fps, RunFingerprint{
			RunIndex:    r.RunIndex,
			Fingerprint: FingerprintRun(r),
		})
	}
	cp := aggregateContinuityProof(runs, inv)
	return CertificationSummary{
		CertificationID:             CertificationID,
		Sprint:                      sprint,
		Lane:                        LaneNativeLive,
		Status:                      status,
		RegressionDetected:          regression,
		RegressionReasons:           reasons,
		Invariants:                  inv,
		Metrics:                     metrics,
		Scores:                      scores,
		RunFingerprints:             fps,
		ContinuityProof:             cp,
		ReadOnly:                    true,
		NoAutomation:                true,
		NoAutonomousOrchestration:   true,
	}
}

func aggregateInvariants(runs []RunMetrics) Invariants {
	inv := Invariants{
		NoRestart: true, NoRollout: true, NoRecreate: true,
		InterruptionWindowMs: 0, InterruptionDetected: false,
	}
	for _, r := range runs {
		if r.RestartCountDelta > 0 || r.RestartCountAfter > r.RestartCountBefore {
			inv.NoRestart = false
		}
		if r.RolloutDetected || r.RolloutCount > 0 {
			inv.NoRollout = false
		}
		if r.RecreateDetected {
			inv.NoRecreate = false
		}
		if r.InterruptionDetected {
			inv.InterruptionDetected = true
		}
		if r.InterruptionWindowMs > inv.InterruptionWindowMs {
			inv.InterruptionWindowMs = r.InterruptionWindowMs
		}
	}
	return inv
}

func aggregateMetrics(runs []RunMetrics) CertificationMetrics {
	m := CertificationMetrics{RunCount: len(runs)}
	for _, r := range runs {
		if r.RestartCountDelta > m.RestartCountDelta {
			m.RestartCountDelta = r.RestartCountDelta
		}
		if r.RolloutCount > m.RolloutCount {
			m.RolloutCount = r.RolloutCount
		}
		if r.RolloutDetected {
			m.RolloutCount = max64(m.RolloutCount, 1)
		}
		if r.RecreateDetected {
			m.RecreateDetected = true
		}
		if r.InterruptionWindowMs > m.InterruptionWindowMs {
			m.InterruptionWindowMs = r.InterruptionWindowMs
		}
		m.ApplyLatencyMs = maxApplyLatency(m.ApplyLatencyMs, r.ApplyLatencyMs)
		m.RollbackLatencyMs = maxRollbackLatency(m.RollbackLatencyMs, r.RollbackLatencyMs)
	}
	return m
}

func max64(a, b int64) int64 {
	if b > a {
		return b
	}
	return a
}

func maxApplyLatency(a, b ApplyLatencyMs) ApplyLatencyMs {
	if b.CPU > a.CPU {
		a.CPU = b.CPU
	}
	if b.RAMUp > a.RAMUp {
		a.RAMUp = b.RAMUp
	}
	if b.RAMDown > a.RAMDown {
		a.RAMDown = b.RAMDown
	}
	return a
}

func maxRollbackLatency(a, b RollbackLatencyMs) RollbackLatencyMs {
	if b.CPU > a.CPU {
		a.CPU = b.CPU
	}
	if b.RAMUp > a.RAMUp {
		a.RAMUp = b.RAMUp
	}
	if b.RAMDown > a.RAMDown {
		a.RAMDown = b.RAMDown
	}
	return a
}

// RegressionReasons lists human-readable regression causes.
func RegressionReasons(runs []RunMetrics, inv Invariants) []string {
	var reasons []string
	for _, r := range runs {
		if r.RestartCountDelta > 0 {
			reasons = append(reasons, fmt.Sprintf("run %d: restart count delta %d", r.RunIndex, r.RestartCountDelta))
		}
		if r.RolloutDetected || r.RolloutCount > 0 {
			reasons = append(reasons, fmt.Sprintf("run %d: rollout detected", r.RunIndex))
		}
		if r.RecreateDetected {
			reasons = append(reasons, fmt.Sprintf("run %d: recreate detected", r.RunIndex))
		}
		if r.InterruptionDetected {
			reasons = append(reasons, fmt.Sprintf("run %d: interruption detected", r.RunIndex))
		}
		if r.InterruptionWindowMs > 0 {
			reasons = append(reasons, fmt.Sprintf("run %d: interruption window %dms", r.RunIndex, r.InterruptionWindowMs))
		}
		if !r.ShellContinuityPreserved {
			reasons = append(reasons, fmt.Sprintf("run %d: shell continuity interrupted", r.RunIndex))
		}
		if !r.AppContinuityPreserved {
			reasons = append(reasons, fmt.Sprintf("run %d: app continuity interrupted", r.RunIndex))
		}
		if !r.SessionContinuityPreserved {
			reasons = append(reasons, fmt.Sprintf("run %d: session continuity interrupted", r.RunIndex))
		}
	}
	if !inv.NoRestart {
		reasons = append(reasons, "invariant: noRestart violated")
	}
	if !inv.NoRollout {
		reasons = append(reasons, "invariant: noRollout violated")
	}
	if !inv.NoRecreate {
		reasons = append(reasons, "invariant: noRecreate violated")
	}
	if inv.InterruptionDetected {
		reasons = append(reasons, "invariant: interruptionDetected")
	}
	if inv.InterruptionWindowMs > 0 {
		reasons = append(reasons, fmt.Sprintf("invariant: interruptionWindowMs=%d", inv.InterruptionWindowMs))
	}
	sort.Strings(reasons)
	return dedupeStrings(reasons)
}

// CheckRegression returns an error when certification invariants fail.
func CheckRegression(summary CertificationSummary) error {
	if !summary.RegressionDetected && summary.Status == CertificationCertified {
		return nil
	}
	if len(summary.RegressionReasons) == 0 {
		return fmt.Errorf("native-live certification regression detected")
	}
	return fmt.Errorf("native-live certification regression: %s", summary.RegressionReasons[0])
}

func aggregateContinuityProof(runs []RunMetrics, inv Invariants) ContinuityCertificationProof {
	cp := ContinuityCertificationProof{
		ResourceContinuityPreserved: inv.NoRestart && inv.NoRollout && inv.NoRecreate && !inv.InterruptionDetected,
		ShellContinuityPreserved:    true,
		AppContinuityPreserved:      true,
		SessionContinuityPreserved:  true,
	}
	for _, r := range runs {
		if !r.ShellContinuityPreserved {
			cp.ShellContinuityPreserved = false
		}
		if !r.AppContinuityPreserved {
			cp.AppContinuityPreserved = false
		}
		if !r.SessionContinuityPreserved {
			cp.SessionContinuityPreserved = false
		}
	}
	return cp
}

// ScoreCertification computes resource + session continuity scores (read-only).
func ScoreCertification(inv Invariants, metrics CertificationMetrics, runs []RunMetrics) CertificationScores {
	resource := 1.0
	if !inv.NoRestart {
		resource -= 0.4
	}
	if !inv.NoRollout {
		resource -= 0.3
	}
	if !inv.NoRecreate {
		resource -= 0.2
	}
	if inv.InterruptionDetected || inv.InterruptionWindowMs > 0 {
		resource -= 0.5
	}
	if metrics.RestartCountDelta > 0 {
		resource -= 0.2
	}
	if resource < 0 {
		resource = 0
	}
	session := 1.0
	for _, r := range runs {
		if !r.ShellContinuityPreserved || !r.AppContinuityPreserved || !r.SessionContinuityPreserved {
			session = 0
			break
		}
	}
	combined := resource
	if session < combined {
		combined = session
	}
	confidence := ConfidenceHigh
	if combined < 0.85 {
		confidence = ConfidenceMedium
	}
	if combined < 0.5 {
		confidence = ConfidenceLow
	}
	for _, r := range runs {
		if !r.LiveInPlaceEligible || r.NativeLiveLaneCount < 1 || !r.RollbackPass {
			confidence = ConfidenceLow
			break
		}
	}
	return CertificationScores{
		ResourceContinuityScore: resource,
		SessionContinuityScore:  session,
		ContinuityScore:         combined,
		LiveScaleConfidence:     confidence,
	}
}

// FingerprintRun returns a deterministic hash for a run (excludes wall-clock noise).
func FingerprintRun(r RunMetrics) string {
	payload := struct {
		RestartDelta         int64 `json:"restartCountDelta"`
		Rollout              bool  `json:"rolloutDetected"`
		Recreate             bool  `json:"recreateDetected"`
		Interruption         bool  `json:"interruptionDetected"`
		InterruptionWindowMs int64 `json:"interruptionWindowMs"`
		NativeLiveLanes      int   `json:"nativeLiveLaneCount"`
		LiveEligible         bool  `json:"liveInPlaceEligible"`
		RollbackPass         bool  `json:"rollbackPass"`
		ShellContinuity      bool  `json:"shellContinuityPreserved"`
		AppContinuity        bool  `json:"appContinuityPreserved"`
		SessionContinuity    bool  `json:"sessionContinuityPreserved"`
	}{
		r.RestartCountDelta, r.RolloutDetected, r.RecreateDetected,
		r.InterruptionDetected, r.InterruptionWindowMs,
		r.NativeLiveLaneCount, r.LiveInPlaceEligible, r.RollbackPass,
		r.ShellContinuityPreserved, r.AppContinuityPreserved, r.SessionContinuityPreserved,
	}
	b, _ := json.Marshal(payload)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:8])
}

// CompareBaseline checks deterministic certification fields against a baseline.
func CompareBaseline(summary CertificationSummary, baseline BaselineCertification) (match bool, diffs []string) {
	if baseline.CertificationID != "" && summary.CertificationID != baseline.CertificationID {
		diffs = append(diffs, "certificationId mismatch")
	}
	if baseline.Lane != "" && summary.Lane != baseline.Lane {
		diffs = append(diffs, "lane mismatch")
	}
	if baseline.Status != "" && summary.Status != baseline.Status {
		diffs = append(diffs, fmt.Sprintf("status: got %s want %s", summary.Status, baseline.Status))
	}
	if baseline.ExpectedRunCount > 0 && summary.Metrics.RunCount != baseline.ExpectedRunCount {
		diffs = append(diffs, fmt.Sprintf("runCount: got %d want %d", summary.Metrics.RunCount, baseline.ExpectedRunCount))
	}
	diffs = append(diffs, compareInvariants(summary.Invariants, baseline.Invariants)...)
	diffs = append(diffs, compareScores(summary.Scores, baseline.Scores)...)
	return len(diffs) == 0, diffs
}

func compareInvariants(got, want Invariants) []string {
	var d []string
	if got.NoRestart != want.NoRestart {
		d = append(d, "invariants.noRestart")
	}
	if got.NoRollout != want.NoRollout {
		d = append(d, "invariants.noRollout")
	}
	if got.NoRecreate != want.NoRecreate {
		d = append(d, "invariants.noRecreate")
	}
	if got.InterruptionWindowMs != want.InterruptionWindowMs {
		d = append(d, "invariants.interruptionWindowMs")
	}
	if got.InterruptionDetected != want.InterruptionDetected {
		d = append(d, "invariants.interruptionDetected")
	}
	return d
}

func compareScores(got, want CertificationScores) []string {
	var d []string
	if want.LiveScaleConfidence != "" && got.LiveScaleConfidence != want.LiveScaleConfidence {
		d = append(d, "scores.liveScaleConfidence")
	}
	if want.ContinuityScore > 0 && got.ContinuityScore < want.ContinuityScore {
		d = append(d, "scores.continuityScore below baseline")
	}
	if want.ResourceContinuityScore > 0 && got.ResourceContinuityScore < want.ResourceContinuityScore {
		d = append(d, "scores.resourceContinuityScore below baseline")
	}
	if want.SessionContinuityScore > 0 && got.SessionContinuityScore < want.SessionContinuityScore {
		d = append(d, "scores.sessionContinuityScore below baseline")
	}
	return d
}

// FingerprintsMatch reports whether all run fingerprints are identical (repeatable runs).
func FingerprintsMatch(summary CertificationSummary) bool {
	if len(summary.RunFingerprints) < 2 {
		return true
	}
	first := summary.RunFingerprints[0].Fingerprint
	for _, fp := range summary.RunFingerprints[1:] {
		if fp.Fingerprint != first {
			return false
		}
	}
	return true
}

func dedupeStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
