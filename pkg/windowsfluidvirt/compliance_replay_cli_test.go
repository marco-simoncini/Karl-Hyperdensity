package windowsfluidvirt

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestComplianceReplayCLIReadyStandalone(t *testing.T) {
	output := runComplianceReplayCLI(t,
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
	)
	if output.CompliancePhase != string(ComplianceHyperdensityReadyWindowsShell) {
		t.Fatalf("expected ready phase, got %s", output.CompliancePhase)
	}
	if !output.HyperdensityReady {
		t.Fatal("expected hyperdensityReady=true")
	}
}

func TestComplianceReplayCLIPoolChildReady(t *testing.T) {
	output := runComplianceReplayCLI(t,
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-pool-child-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
	)
	if output.CompliancePhase != string(ComplianceHyperdensityReadyWindowsShell) {
		t.Fatalf("expected ready phase, got %s", output.CompliancePhase)
	}
	if !output.PoolContext.IsPoolChild {
		t.Fatal("expected pool child context")
	}
	if output.PoolScalingMechanismBlocked {
		t.Fatal("pool child as provisioning context must not be blocked")
	}
}

func TestComplianceReplayCLIPoolScalingBlocked(t *testing.T) {
	output := runComplianceReplayCLI(t,
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-pool-scaling-mechanism.blocked.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
	)
	if output.CompliancePhase != string(ComplianceBlockedWithRemediation) {
		t.Fatalf("expected blocked phase, got %s", output.CompliancePhase)
	}
	if !output.PoolScalingMechanismBlocked {
		t.Fatal("expected poolScalingMechanismBlocked=true")
	}
	assertHas(t, output.Blockers, BlockerPoolScalingAsMechanism)
}

func TestComplianceReplayCLIDeterministicOutput(t *testing.T) {
	args := []string{
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
	}
	raw1 := runComplianceReplayCLIRaw(t, args...)
	raw2 := runComplianceReplayCLIRaw(t, args...)
	if raw1 != raw2 {
		t.Fatal("cli output must be deterministic with fixed evaluation-time")
	}
}

func TestComplianceReplayHashAndEvidenceHashStable(t *testing.T) {
	args := []string{
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
	}
	out1 := runComplianceReplayCLI(t, args...)
	out2 := runComplianceReplayCLI(t, args...)
	if out1.EvidenceHash != out2.EvidenceHash {
		t.Fatalf("evidence hash must be stable: %s vs %s", out1.EvidenceHash, out2.EvidenceHash)
	}
	if out1.ReplayHash != out2.ReplayHash {
		t.Fatalf("replay hash must be stable: %s vs %s", out1.ReplayHash, out2.ReplayHash)
	}
}

func TestComplianceReplayAttestationFutureSignable(t *testing.T) {
	output := runComplianceReplayCLI(t,
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
		"-emit-attestation",
		"-attestation-mode", string(AttestationModeFutureSignable),
	)
	if output.Attestation == nil {
		t.Fatal("expected attestation to be emitted")
	}
	if output.Attestation.Signature.Mode != AttestationModeFutureSignable {
		t.Fatalf("unexpected attestation mode %s", output.Attestation.Signature.Mode)
	}
	if output.Attestation.Signature.Value != "" {
		t.Fatal("attestation signature value must stay empty")
	}
}

func TestComplianceReplayInvalidAttestationModeRejected(t *testing.T) {
	cmd := exec.Command(
		"go", "run", "./cmd/karl-fluid-compliance-replay",
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
		"-emit-attestation",
		"-attestation-mode", "invalid-mode",
	)
	cmd.Dir = admissionRepoRoot(t)
	if _, err := cmd.CombinedOutput(); err == nil {
		t.Fatal("expected invalid attestation mode to fail")
	}
}

func TestComplianceReplayNoMutationFlagsTrue(t *testing.T) {
	output := runComplianceReplayCLI(t,
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
	)
	if output.MutationFlags.RuntimeMutation || output.MutationFlags.CPUApply || output.MutationFlags.RAMApply || output.MutationFlags.ActuatorApply || output.MutationFlags.ClusterCalls || output.MutationFlags.QMPCalls {
		t.Fatalf("all mutation flags must stay false: %+v", output.MutationFlags)
	}
}

func TestComplianceReplayBundleSingleRunDeterministic(t *testing.T) {
	args := []string{
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
		"-emit-attestation",
		"-attestation-mode", "future-signable",
		"-emit-bundle-index",
		"-bundle-subject", "windows-shell/karl/master-win11",
	}
	out1 := runComplianceReplayCLI(t, args...)
	out2 := runComplianceReplayCLI(t, args...)
	if out1.BundleIndex == nil || out2.BundleIndex == nil {
		t.Fatal("expected bundle index output")
	}
	if out1.BundleIndex.Chain.LatestRunHash != out2.BundleIndex.Chain.LatestRunHash {
		t.Fatal("bundle latest run hash must be deterministic")
	}
	if out1.BundleIndex.BundleID != out2.BundleIndex.BundleID {
		t.Fatal("bundle id must be deterministic with fixed evaluation-time")
	}
}

func TestComplianceReplayBundleContainsHashes(t *testing.T) {
	output := runComplianceReplayCLI(t,
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
		"-emit-attestation",
		"-attestation-mode", "future-signable",
		"-emit-bundle-index",
	)
	if output.BundleIndex == nil {
		t.Fatal("bundle index missing")
	}
	if output.BundleIndex.RunCount != 1 {
		t.Fatalf("expected 1 run, got %d", output.BundleIndex.RunCount)
	}
	run := output.BundleIndex.Runs[0]
	if run.EvidenceHash == "" || run.ReplayHash == "" || run.AttestationHash == "" || run.RunHash == "" {
		t.Fatalf("run hashes must be populated: %+v", run)
	}
}

func TestComplianceReplayRunHashIncludesPreviousHash(t *testing.T) {
	output := runComplianceReplayCLI(t,
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
		"-emit-bundle-index",
		"-previous-run-hash", "abc123previous",
	)
	run := output.BundleIndex.Runs[0]
	if run.PreviousRunHash != "abc123previous" {
		t.Fatalf("expected previous hash to be wired, got %s", run.PreviousRunHash)
	}
	withoutPrev, err := BuildWindowsComplianceReplayBundleRun(output.WindowsComplianceReplayOutput, output.Attestation, "", time.Date(2026, 5, 7, 20, 45, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build run without previous hash: %v", err)
	}
	if run.RunHash == withoutPrev.RunHash {
		t.Fatal("runHash must change when previousRunHash changes")
	}
}

func TestComplianceReplayFirstRunPreviousHashEmpty(t *testing.T) {
	output := runComplianceReplayCLI(t,
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
		"-emit-bundle-index",
	)
	run := output.BundleIndex.Runs[0]
	if run.PreviousRunHash != "" {
		t.Fatalf("first run previous hash must be empty, got %s", run.PreviousRunHash)
	}
}

func TestComplianceReplayTwoRunChainValid(t *testing.T) {
	evalTime := time.Date(2026, 5, 7, 20, 45, 0, 0, time.UTC)
	in1 := mustLoadComplianceFixture(t, "master-win11-real-evidence.ready.json").Input
	replay1, err := EvaluateWindowsComplianceReplay(in1, "run1", evalTime)
	if err != nil {
		t.Fatalf("evaluate replay1: %v", err)
	}
	att1, err := BuildWindowsComplianceReplayAttestation(replay1, AttestationModeFutureSignable, evalTime)
	if err != nil {
		t.Fatalf("attestation1: %v", err)
	}
	run1, err := BuildWindowsComplianceReplayBundleRun(replay1, &att1, "", evalTime)
	if err != nil {
		t.Fatalf("bundle run1: %v", err)
	}
	in2 := mustLoadComplianceFixture(t, "master-win11-pool-child-real-evidence.ready.json").Input
	replay2, err := EvaluateWindowsComplianceReplay(in2, "run2", evalTime)
	if err != nil {
		t.Fatalf("evaluate replay2: %v", err)
	}
	att2, err := BuildWindowsComplianceReplayAttestation(replay2, AttestationModeFutureSignable, evalTime)
	if err != nil {
		t.Fatalf("attestation2: %v", err)
	}
	run2, err := BuildWindowsComplianceReplayBundleRun(replay2, &att2, run1.RunHash, evalTime)
	if err != nil {
		t.Fatalf("bundle run2: %v", err)
	}
	bundle, err := BuildWindowsComplianceReplayBundleIndex("windows-shell/karl/master-win11", "windows-fluid-compliance-replay-bundle-index-v1", []WindowsComplianceReplayBundleRun{run1, run2}, evalTime)
	if err != nil {
		t.Fatalf("build bundle: %v", err)
	}
	if !bundle.Chain.ChainValid {
		t.Fatalf("expected chain valid, got notes=%s", bundle.Chain.Notes)
	}
	if bundle.Chain.LatestRunHash != run2.RunHash {
		t.Fatal("latest run hash mismatch")
	}
}

func TestComplianceReplayTwoRunChainReadyThenBlocked(t *testing.T) {
	evalTime := time.Date(2026, 5, 7, 20, 45, 0, 0, time.UTC)
	inReady := mustLoadComplianceFixture(t, "master-win11-real-evidence.ready.json").Input
	replayReady, _ := EvaluateWindowsComplianceReplay(inReady, "ready", evalTime)
	runReady, _ := BuildWindowsComplianceReplayBundleRun(replayReady, nil, "", evalTime)
	inBlocked := mustLoadComplianceFixture(t, "master-win11-pool-scaling-mechanism.blocked.json").Input
	replayBlocked, _ := EvaluateWindowsComplianceReplay(inBlocked, "blocked", evalTime)
	runBlocked, _ := BuildWindowsComplianceReplayBundleRun(replayBlocked, nil, runReady.RunHash, evalTime)
	bundle, err := BuildWindowsComplianceReplayBundleIndex("windows-shell/karl/master-win11", "windows-fluid-compliance-replay-bundle-index-v1", []WindowsComplianceReplayBundleRun{runReady, runBlocked}, evalTime)
	if err != nil {
		t.Fatalf("build bundle: %v", err)
	}
	if bundle.LatestHyperdensityReady {
		t.Fatal("latest status must be false when latest run is blocked")
	}
	if bundle.LatestCompliancePhase != string(ComplianceBlockedWithRemediation) {
		t.Fatalf("unexpected latest phase %s", bundle.LatestCompliancePhase)
	}
}

func TestComplianceReplayBrokenPreviousHashInvalid(t *testing.T) {
	evalTime := time.Date(2026, 5, 7, 20, 45, 0, 0, time.UTC)
	in := mustLoadComplianceFixture(t, "master-win11-real-evidence.ready.json").Input
	replay, _ := EvaluateWindowsComplianceReplay(in, "run1", evalTime)
	run1, _ := BuildWindowsComplianceReplayBundleRun(replay, nil, "", evalTime)
	run2, _ := BuildWindowsComplianceReplayBundleRun(replay, nil, "wrong-previous", evalTime)
	index := WindowsComplianceReplayBundleIndex{
		BundleID:      "test",
		BundleVersion: "windows-fluid-compliance-replay-bundle-index-v1",
		SubjectRef:    "windows-shell/karl/master-win11",
		SubjectType:   "windows-hyperdensity-ready-compliance-replay-bundle",
		RunCount:      2,
		Runs:          []WindowsComplianceReplayBundleRun{run1, run2},
	}
	validated := ValidateWindowsComplianceReplayBundleIndex(index)
	if validated.Chain.ChainValid {
		t.Fatal("expected broken chain")
	}
}

func TestComplianceReplayRunCountMismatchInvalid(t *testing.T) {
	index := WindowsComplianceReplayBundleIndex{
		BundleID:      "test",
		BundleVersion: "windows-fluid-compliance-replay-bundle-index-v1",
		RunCount:      2,
		Runs:          []WindowsComplianceReplayBundleRun{{RunID: "only-one"}},
	}
	validated := ValidateWindowsComplianceReplayBundleIndex(index)
	if validated.Chain.ChainValid {
		t.Fatal("expected invalid chain for runCount mismatch")
	}
}

func TestComplianceReplayBundleAttestationModeAndSignatureConstraints(t *testing.T) {
	evalTime := time.Date(2026, 5, 7, 20, 45, 0, 0, time.UTC)
	in := mustLoadComplianceFixture(t, "master-win11-real-evidence.ready.json").Input
	replay, _ := EvaluateWindowsComplianceReplay(in, "run1", evalTime)
	run, _ := BuildWindowsComplianceReplayBundleRun(replay, nil, "", evalTime)
	run.AttestationMode = "invalid"
	index := WindowsComplianceReplayBundleIndex{
		BundleID:      "test",
		BundleVersion: "windows-fluid-compliance-replay-bundle-index-v1",
		RunCount:      1,
		Runs:          []WindowsComplianceReplayBundleRun{run},
	}
	validated := ValidateWindowsComplianceReplayBundleIndex(index)
	if validated.Chain.ChainValid {
		t.Fatal("expected invalid chain for invalid attestation mode")
	}

	run.AttestationMode = AttestationModeFutureSignable
	run.AttestationSignatureValue = "not-empty"
	index.Runs = []WindowsComplianceReplayBundleRun{run}
	validated = ValidateWindowsComplianceReplayBundleIndex(index)
	if validated.Chain.ChainValid {
		t.Fatal("expected invalid chain for non-empty signature value")
	}
}

func TestComplianceReplayCLIEmitsBundleIndexDeterministically(t *testing.T) {
	args := []string{
		"-input", windowsComplianceFixtureAbsPath(t, "master-win11-real-evidence.ready.json"),
		"-evaluation-time", "2026-05-07T20:45:00Z",
		"-emit-attestation",
		"-attestation-mode", "future-signable",
		"-emit-bundle-index",
		"-bundle-subject", "windows-shell/karl/master-win11",
	}
	raw1 := runComplianceReplayCLIRaw(t, args...)
	raw2 := runComplianceReplayCLIRaw(t, args...)
	if raw1 != raw2 {
		t.Fatal("bundle index cli output must be deterministic")
	}
}

func runComplianceReplayCLI(t *testing.T, args ...string) WindowsComplianceReplayCLIOutput {
	t.Helper()
	raw := runComplianceReplayCLIRaw(t, args...)
	var output WindowsComplianceReplayCLIOutput
	if err := json.Unmarshal([]byte(raw), &output); err != nil {
		t.Fatalf("parse cli output: %v, raw=%s", err, raw)
	}
	return output
}

func runComplianceReplayCLIRaw(t *testing.T, args ...string) string {
	t.Helper()
	cmdArgs := append([]string{"run", "./cmd/karl-fluid-compliance-replay"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = admissionRepoRoot(t)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("compliance replay cli failed: %v\n%s", err, string(out))
	}
	return string(out)
}

func complianceReplayArtifactAbsPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(admissionRepoRoot(t), "artifacts", "the-father-windows-compliance-replay-cli-attestation-v1", name)
}
