package windowsfluidvirt

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestDryRunFixtureMatrix(t *testing.T) {
	fixtures := []string{
		"master-win11-missing-annotations.blocked.json",
		"master-win11-missing-qmp.blocked.json",
		"master-win11-missing-guest.blocked.json",
		"master-win11-certification-ready.ready.json",
		"win11-pool-context-only.blocked.json",
		"identity-change.quarantined.json",
		"lease-prepared-cpu.dryrun.json",
		"lease-memory-return-not-ready.blocked.json",
	}

	for _, fixtureName := range fixtures {
		t.Run(fixtureName, func(t *testing.T) {
			fixturePath := fixtureAbsPath(t, fixtureName)
			fixture, err := LoadDryRunReplayFixture(fixturePath)
			if err != nil {
				t.Fatalf("load fixture: %v", err)
			}
			result := EvaluateWindowsFluidRuntimeDryRunWithOptions(fixture.Bundle, DryRunEvaluationOptions{
				EvaluationTime: time.Date(2026, 5, 7, 14, 33, 0, 0, time.UTC),
			})
			if result.Phase != fixture.ExpectedPhase {
				t.Fatalf("phase mismatch expected=%s got=%s blockers=%v", fixture.ExpectedPhase, result.Phase, result.Blockers)
			}
			if result.Classification != fixture.ExpectedClassification {
				t.Fatalf("classification mismatch expected=%s got=%s", fixture.ExpectedClassification, result.Classification)
			}
			for _, expectedBlocker := range fixture.ExpectedBlockers {
				assertHas(t, result.Blockers, expectedBlocker)
			}
			if result.ActionSlate.MutationAllowed {
				t.Fatal("mutationAllowed must be false in dry-run")
			}
			if result.ActionSlate.ApplyAllowed {
				t.Fatal("applyAllowed must be false in dry-run")
			}
			slateJSON := strings.ToLower(stringifyActionSlate(t, result.ActionSlate))
			for forbidden := range QMPForbiddenCommands {
				if strings.Contains(slateJSON, "\""+forbidden+"\"") {
					t.Fatalf("action slate must not include forbidden qmp command %s", forbidden)
				}
			}
		})
	}
}

func TestDryRunIncompleteBundleBlocked(t *testing.T) {
	bundle := NewRuntimeEvidenceBundle("karl-system", "master-win11")
	bundle.Shell = testShell()
	bundle.PolicyGates.Annotations = map[string]string{}
	bundle.Guest = nil
	bundle.QMP = nil

	result := EvaluateWindowsFluidRuntimeDryRun(bundle)
	if result.Phase != StateBlocked {
		t.Fatalf("expected BLOCKED for incomplete bundle, got %s", result.Phase)
	}
	assertHas(t, result.Blockers, BlockerGuestAckMissing)
	assertHas(t, result.Blockers, BlockerQMPSocketUnavailable)
}

func TestDryRunPendingRebootBlocked(t *testing.T) {
	bundle := NewRuntimeEvidenceBundle("karl-system", "master-win11")
	shell := testShell()
	bundle.Shell = shell
	bundle.KubeVirtBefore = KubeVirtRuntimeIdentityEvidence{
		VMName: "master-win11", VMNamespace: "karl-system", VMIName: "master-win11", VMIUID: "vmi", VirtLauncherPodUID: "pod", NodeName: "karl-metal-01", QemuPID: "5120",
	}
	bundle.KubeVirtAfter = bundle.KubeVirtBefore
	bundle.PolicyGates.Annotations = RequiredRuntimeAnnotations
	bundle.QMP = &QMPEvidence{
		QMPConnected: true, QMPCapabilitiesNegotiated: true, QMPReadOnly: true, QMPCommandsExecuted: []string{"qmp_capabilities"},
	}
	bundle.Guest = &GuestRuntimeEvidence{
		GuestAck: true, PendingReboot: true, LastBootTime: "2026-05-01T11:00:00Z", MachineGUIDHash: "hash", MemoryAdapterVerified: true, ReturnToFloorReady: true,
	}
	bundle.SourceMetadata = RuntimeSourceMetadata{SourceKind: "vm", SourceName: "master-win11", SourceNamespace: "karl-system"}

	result := EvaluateWindowsFluidRuntimeDryRun(bundle)
	if result.Phase != StateBlocked {
		t.Fatalf("expected BLOCKED for pending reboot, got %s", result.Phase)
	}
	assertHas(t, result.Blockers, BlockerPendingRebootDetected)
}

func TestDryRunCLIEmitsDeterministicJSON(t *testing.T) {
	fixturePath := fixtureAbsPath(t, "master-win11-certification-ready.ready.json")
	cmd := exec.Command("go", "run", "./cmd/karl-fluid-dryrun", "-fixture", fixturePath, "-evaluation-time", "2026-05-07T14:33:00Z")
	cmd.Dir = repoRoot(t)
	out1, err := cmd.Output()
	if err != nil {
		t.Fatalf("first cli run failed: %v", err)
	}

	cmd = exec.Command("go", "run", "./cmd/karl-fluid-dryrun", "-fixture", fixturePath, "-evaluation-time", "2026-05-07T14:33:00Z")
	cmd.Dir = repoRoot(t)
	out2, err := cmd.Output()
	if err != nil {
		t.Fatalf("second cli run failed: %v", err)
	}
	if string(out1) != string(out2) {
		t.Fatal("dry-run cli output must be deterministic with fixed evaluation-time")
	}
}

func stringifyActionSlate(t *testing.T, slate WindowsFluidActionSlate) string {
	t.Helper()
	data, err := json.Marshal(slate)
	if err != nil {
		t.Fatalf("marshal action slate: %v", err)
	}
	return string(data)
}

func fixtureAbsPath(t *testing.T, fixture string) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "examples", "windows-fluid-dryrun-fixtures", fixture)
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
