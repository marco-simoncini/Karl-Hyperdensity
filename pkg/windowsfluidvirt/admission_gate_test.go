package windowsfluidvirt

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestDefaultWindowsFluidPolicyPackIsConservative(t *testing.T) {
	policy := DefaultWindowsFluidPolicyPack()
	if policy.AllowMutationInThisPhase {
		t.Fatal("default policy must disable mutation")
	}
	if policy.AllowPoolReplicaModel {
		t.Fatal("default policy must deny pool replica model")
	}
	if policy.AllowGenericWindowsVm {
		t.Fatal("default policy must deny generic windows vm")
	}
	if policy.AllowedRuntimeMode != "in-place-qmp" {
		t.Fatalf("unexpected allowed runtime mode %s", policy.AllowedRuntimeMode)
	}
}

func TestBlockerPriorityPolicyTiers(t *testing.T) {
	policy := DefaultWindowsFluidPolicyPack()
	if BlockerPriorityForID(BlockerQemuPIDChanged, policy) != BlockerPriorityP0Quarantine {
		t.Fatal("qemu_pid_changed must be P0")
	}
	if BlockerPriorityForID(BlockerQMPSocketUnavailable, policy) != BlockerPriorityP1HardBlock {
		t.Fatal("qmp_socket_unavailable must be P1")
	}
	if BlockerPriorityForID(BlockerMemoryDriverUnverified, policy) != BlockerPriorityP2Capability {
		t.Fatal("memory_driver_unverified must be P2")
	}
	if BlockerPriorityForID(BlockerDashboard443TouchRisk, policy) != BlockerPriorityP3Environment {
		t.Fatal("dashboard_443_touch_risk must be P3")
	}
}

func TestAdmissionFixtureMatrix(t *testing.T) {
	fixtures := []string{
		"admission-master-win11-cpu.future-apply-admissible.json",
		"admission-master-win11-ram.memory-driver-blocked.json",
		"admission-master-win11-missing-return.blocked.json",
		"admission-master-win11-identity-change.quarantined.json",
		"admission-win11-pool.denied.json",
		"admission-generic-windows.denied.json",
		"admission-missing-qmp.blocked.json",
		"admission-missing-guest.blocked.json",
	}
	for _, fixtureName := range fixtures {
		t.Run(fixtureName, func(t *testing.T) {
			fixture, err := LoadAdmissionReplayFixture(admissionFixtureAbsPath(t, fixtureName))
			if err != nil {
				t.Fatalf("load fixture: %v", err)
			}
			result := EvaluateWindowsFluidAdmission(AdmissionEvaluationInput{
				Bundle:          fixture.Bundle,
				PolicyPack:      fixture.PolicyPack,
				RequestedAction: fixture.RequestedAction,
				EvaluationTime:  time.Date(2026, 5, 7, 14, 40, 0, 0, time.UTC),
			})
			if result.Decision.AdmissionPhase != fixture.ExpectedAdmissionPhase {
				t.Fatalf("phase mismatch expected=%s got=%s denial=%v blockers=%v",
					fixture.ExpectedAdmissionPhase, result.Decision.AdmissionPhase, result.Decision.DenialReasons, result.Decision.Blockers)
			}
			for _, expected := range fixture.ExpectedDecisionBlockers {
				assertHas(t, result.Decision.Blockers, expected)
			}
			if result.Decision.MutationAllowed {
				t.Fatal("mutationAllowed must always be false")
			}
			if result.Decision.ApplyAllowed {
				t.Fatal("applyAllowed must always be false")
			}
			if string(result.Decision.AdmissionPhase) == "ACTIVE" {
				t.Fatal("admission phase must never be ACTIVE")
			}
		})
	}
}

func TestPolicyRejectsLiveMigrationRecreateReboot(t *testing.T) {
	policy := DefaultWindowsFluidPolicyPack()
	bundle := NewRuntimeEvidenceBundle("karl", "master-win11")
	bundle.Shell = testShell()
	bundle.Shell.Spec.VMRef = "karl/master-win11"
	bundle.PolicyGates.Annotations = RequiredRuntimeAnnotations
	bundle.KubeVirtBefore = KubeVirtRuntimeIdentityEvidence{
		VMName: "master-win11", VMNamespace: "karl", VMIName: "master-win11", VMIUID: "vmi-master", VirtLauncherPodUID: "pod-master", NodeName: "karl-metal-01", QemuPID: "5120",
	}
	bundle.KubeVirtAfter = bundle.KubeVirtBefore
	bundle.KubeVirtAfter.MigrationRequired = true
	bundle.KubeVirtAfter.RecreateRequired = true
	bundle.QMP = &QMPEvidence{QMPConnected: true, QMPCapabilitiesNegotiated: true, QMPReadOnly: true, QMPCommandsExecuted: []string{"qmp_capabilities"}}
	bundle.Guest = &GuestRuntimeEvidence{
		GuestAck: true, PendingReboot: true, LastBootTime: "2026-05-01T11:00:00Z", MachineGUIDHash: "hash", MemoryAdapterVerified: true, ReturnToFloorReady: true,
	}
	bundle.Timestamps = map[string]string{"collectedAt": "2026-05-07T14:30:00Z"}

	result := EvaluateWindowsFluidAdmission(AdmissionEvaluationInput{
		Bundle:          bundle,
		PolicyPack:      &policy,
		RequestedAction: RequestedActionPrepareCPULease,
		EvaluationTime:  time.Date(2026, 5, 7, 14, 40, 0, 0, time.UTC),
	})
	if result.Decision.AdmissionPhase != AdmissionBlocked {
		t.Fatalf("expected BLOCKED, got %s", result.Decision.AdmissionPhase)
	}
	assertHas(t, result.PolicyViolations, "no_live_migration_required")
	assertHas(t, result.PolicyViolations, "no_recreate_required")
	assertHas(t, result.PolicyViolations, "no_reboot_required")
}

func TestAdmissionOutputDeterministicWithFixedTime(t *testing.T) {
	fixturePath := admissionFixtureAbsPath(t, "admission-master-win11-cpu.future-apply-admissible.json")
	fixture, err := LoadAdmissionReplayFixture(fixturePath)
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	evalTime := time.Date(2026, 5, 7, 14, 40, 0, 0, time.UTC)
	first := EvaluateWindowsFluidAdmission(AdmissionEvaluationInput{
		Bundle: fixture.Bundle, PolicyPack: fixture.PolicyPack, RequestedAction: fixture.RequestedAction, EvaluationTime: evalTime,
	})
	second := EvaluateWindowsFluidAdmission(AdmissionEvaluationInput{
		Bundle: fixture.Bundle, PolicyPack: fixture.PolicyPack, RequestedAction: fixture.RequestedAction, EvaluationTime: evalTime,
	})
	f1, _ := json.Marshal(first)
	f2, _ := json.Marshal(second)
	if string(f1) != string(f2) {
		t.Fatal("admission output must be deterministic with fixed evaluation time")
	}
}

func TestAdmissionCLIOutputsJSON(t *testing.T) {
	fixturePath := admissionFixtureAbsPath(t, "admission-master-win11-cpu.future-apply-admissible.json")
	cmd := exec.Command("go", "run", "./cmd/karl-fluid-admission", "-fixture", fixturePath, "-evaluation-time", "2026-05-07T14:40:00Z")
	cmd.Dir = admissionRepoRoot(t)
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("admission cli failed: %v", err)
	}
	var parsed AdmissionEvaluationResult
	if err := json.Unmarshal(output, &parsed); err != nil {
		t.Fatalf("admission cli output is not valid json: %v", err)
	}
}

func admissionFixtureAbsPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(admissionRepoRoot(t), "examples", "windows-fluid-admission-fixtures", name)
}

func admissionRepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
