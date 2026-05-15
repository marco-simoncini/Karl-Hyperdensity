package windowsfluidvirt

import (
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestExecutorFixtureMatrix(t *testing.T) {
	fixtures := []string{
		"executor-master-win11-cpu.hard-disabled.json",
		"executor-master-win11-stale.blocked.json",
		"executor-master-win11-identity-change.quarantined.json",
		"executor-master-win11-killswitch-missing.blocked.json",
		"executor-master-win11-killswitch-enabled.hard-disabled.json",
		"executor-win11-pool.denied.json",
		"executor-generic-windows.denied.json",
		"executor-command-envelope-empty.proof.json",
	}
	for _, fixtureName := range fixtures {
		t.Run(fixtureName, func(t *testing.T) {
			fixturePath := executorFixtureAbsPath(t, fixtureName)
			fixture, err := LoadExecutorReplayFixture(fixturePath)
			if err != nil {
				t.Fatalf("load fixture: %v", err)
			}
			result := EvaluateWindowsFluidFutureApplyExecutor(FutureApplyExecutorEvaluationInput{
				GovernanceContract: fixture.GovernanceContract,
				Revalidation:       fixture.Revalidation,
				Attestation:        fixture.Attestation,
				KillSwitch:         fixture.KillSwitch,
				EvaluationTime:     time.Date(2026, 5, 7, 16, 0, 0, 0, time.UTC),
			})
			if result.ExecutionResult.ExecutionPhase != fixture.ExpectedExecutionPhase {
				t.Fatalf("phase mismatch expected=%s got=%s", fixture.ExpectedExecutionPhase, result.ExecutionResult.ExecutionPhase)
			}
			for _, expected := range fixture.ExpectedBlockers {
				assertHas(t, result.ExecutionResult.Blockers, expected)
			}
			if result.ExecutionResult.ApplyAttempted {
				t.Fatal("applyAttempted must always be false")
			}
			if result.ExecutionResult.MutationPerformed {
				t.Fatal("mutationPerformed must always be false")
			}
			if result.ExecutionResult.QMPCommandSent {
				t.Fatal("qmpCommandSent must always be false")
			}
			if result.ExecutionResult.ClusterMutationSent {
				t.Fatal("clusterMutationSent must always be false")
			}
		})
	}
}

func TestExecutorNoActiveApplyingExecutedStates(t *testing.T) {
	fixture, err := LoadExecutorReplayFixture(executorFixtureAbsPath(t, "executor-master-win11-cpu.hard-disabled.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	result := EvaluateWindowsFluidFutureApplyExecutor(FutureApplyExecutorEvaluationInput{
		GovernanceContract: fixture.GovernanceContract,
		Revalidation:       fixture.Revalidation,
		Attestation:        fixture.Attestation,
		KillSwitch:         fixture.KillSwitch,
		EvaluationTime:     time.Date(2026, 5, 7, 16, 0, 0, 0, time.UTC),
	})
	if result.ExecutionResult.ExecutionPhase == ExecutionPhase("ACTIVE") ||
		result.ExecutionResult.ExecutionPhase == ExecutionPhase("APPLYING") ||
		result.ExecutionResult.ExecutionPhase == ExecutionPhase("EXECUTED") {
		t.Fatalf("forbidden execution phase emitted: %s", result.ExecutionResult.ExecutionPhase)
	}
}

func TestExecutorCommandEnvelopeIsPreviewOnly(t *testing.T) {
	fixture, err := LoadExecutorReplayFixture(executorFixtureAbsPath(t, "executor-command-envelope-empty.proof.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	result := EvaluateWindowsFluidFutureApplyExecutor(FutureApplyExecutorEvaluationInput{
		GovernanceContract: fixture.GovernanceContract,
		Revalidation:       fixture.Revalidation,
		Attestation:        fixture.Attestation,
		KillSwitch:         fixture.KillSwitch,
		EvaluationTime:     time.Date(2026, 5, 7, 16, 0, 0, 0, time.UTC),
	})
	if !result.CommandEnvelope.CommandPreviewOnly {
		t.Fatal("command envelope must be preview-only")
	}
	if result.CommandEnvelope.ContainsExecutableCommand {
		t.Fatal("command envelope must not contain executable commands")
	}
	if len(result.CommandEnvelope.QMPCommands) != 0 {
		t.Fatal("qmpCommands must be empty")
	}
	if len(result.CommandEnvelope.ClusterMutations) != 0 {
		t.Fatal("clusterMutations must be empty")
	}
	if len(result.CommandEnvelope.GuestMutations) != 0 {
		t.Fatal("guestMutations must be empty")
	}
}

func TestExecutorKillSwitchMissingBlocks(t *testing.T) {
	fixture, err := LoadExecutorReplayFixture(executorFixtureAbsPath(t, "executor-master-win11-killswitch-missing.blocked.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	result := EvaluateWindowsFluidFutureApplyExecutor(FutureApplyExecutorEvaluationInput{
		GovernanceContract: fixture.GovernanceContract,
		Revalidation:       fixture.Revalidation,
		Attestation:        fixture.Attestation,
		KillSwitch:         nil,
		EvaluationTime:     time.Date(2026, 5, 7, 16, 0, 0, 0, time.UTC),
	})
	if result.ExecutionResult.ExecutionPhase != ExecutionBlocked {
		t.Fatalf("expected blocked when kill switch missing, got %s", result.ExecutionResult.ExecutionPhase)
	}
}

func TestExecutorAttestationSignatureSafety(t *testing.T) {
	fixture, err := LoadExecutorReplayFixture(executorFixtureAbsPath(t, "executor-master-win11-killswitch-enabled.hard-disabled.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	result := EvaluateWindowsFluidFutureApplyExecutor(FutureApplyExecutorEvaluationInput{
		GovernanceContract: fixture.GovernanceContract,
		Revalidation:       fixture.Revalidation,
		Attestation:        fixture.Attestation,
		KillSwitch:         fixture.KillSwitch,
		EvaluationTime:     time.Date(2026, 5, 7, 16, 0, 0, 0, time.UTC),
	})
	mode := fixture.Attestation.Signature.Mode
	if mode != "unsigned-dev" && mode != "future-signable" {
		t.Fatalf("unexpected attestation mode %s", mode)
	}
	if fixture.Attestation.Signature.Value != "" {
		t.Fatal("attestation signature value must stay empty")
	}
	if result.ExecutionResult.ExecutionPhase != ExecutionHardDisabled {
		t.Fatalf("kill switch enabled must still hard-disable execution, got %s", result.ExecutionResult.ExecutionPhase)
	}
}

func TestExecutorCLIOutputDeterministic(t *testing.T) {
	fixturePath := executorFixtureAbsPath(t, "executor-master-win11-cpu.hard-disabled.json")
	cmd := exec.Command("go", "run", "./cmd/karl-fluid-executor", "-fixture", fixturePath, "-evaluation-time", "2026-05-07T16:00:00Z")
	cmd.Dir = governanceRepoRoot(t)
	out1, err := cmd.Output()
	if err != nil {
		t.Fatalf("first executor cli run failed: %v", err)
	}
	cmd = exec.Command("go", "run", "./cmd/karl-fluid-executor", "-fixture", fixturePath, "-evaluation-time", "2026-05-07T16:00:00Z")
	cmd.Dir = governanceRepoRoot(t)
	out2, err := cmd.Output()
	if err != nil {
		t.Fatalf("second executor cli run failed: %v", err)
	}
	if string(out1) != string(out2) {
		t.Fatal("executor cli output must be deterministic with fixed evaluation-time")
	}
}

func executorFixtureAbsPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(governanceRepoRoot(t), "examples", "windows-fluid-executor-fixtures", name)
}
