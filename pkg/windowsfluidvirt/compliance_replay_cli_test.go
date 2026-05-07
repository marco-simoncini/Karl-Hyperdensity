package windowsfluidvirt

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"
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
