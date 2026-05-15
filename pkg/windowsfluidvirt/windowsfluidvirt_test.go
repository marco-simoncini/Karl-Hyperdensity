package windowsfluidvirt

import (
	"testing"
	"time"
)

func TestCanonicalBlockersContainMandatoryIDs(t *testing.T) {
	required := []string{
		BlockerQMPSocketUnavailable,
		BlockerGuestAgentUnavailable,
		BlockerKarlAgentFluidModuleMissing,
		BlockerPendingRebootDetected,
		BlockerQemuPIDChanged,
		BlockerLastBootChanged,
		BlockerMachineGUIDChanged,
		BlockerLiveMigrationRequired,
		BlockerVMIRecreateRequired,
		BlockerVirtLauncherPodChanged,
		BlockerNodeChanged,
		BlockerMemoryDriverUnverified,
		BlockerMemoryReturnNotSafe,
		BlockerCPUTopologyNotConfirmed,
		BlockerGuestMemoryNotConfirmed,
		BlockerRollbackNotReady,
		BlockerReturnToFloorNotReady,
		BlockerQMPAckMissing,
		BlockerGuestAckMissing,
		BlockerHotplugErrorDetected,
		BlockerCriticalWindowsEventDetected,
		BlockerDashboard443TouchRisk,
		BlockerCandidate8888Unavailable,
		BlockerWindowsAgentRepoNotPresentInTarget,
		BlockerFutureApplyExecutorDisabled,
	}
	for _, id := range required {
		if _, ok := LookupBlocker(id); !ok {
			t.Fatalf("missing canonical blocker id %s", id)
		}
	}
}

func TestValidateWindowsFluidShellBlocksInvalidRuntimeMode(t *testing.T) {
	shell := testShell()
	shell.Spec.RuntimeMode = "other-mode"
	blockers := ValidateWindowsFluidShell(shell)
	assertHas(t, blockers, BlockerQMPSocketUnavailable)
}

func TestEvaluateLeaseCanBecomeActiveBlocksMissingAcks(t *testing.T) {
	lease := testLease()
	lease.Status.QMPAck = false
	lease.Status.GuestAck = false
	evidence := testEvidence()
	result := EvaluateLeaseCanBecomeActive(lease, evidence)
	if result.Allowed {
		t.Fatal("lease should not become active without acknowledgments")
	}
	assertHas(t, result.Blockers, BlockerQMPAckMissing)
	assertHas(t, result.Blockers, BlockerGuestAckMissing)
}

func TestEvaluateLeaseCanBecomeActiveQuarantinesOnPIDChange(t *testing.T) {
	lease := testLease()
	evidence := testEvidence()
	evidence.QemuPIDAfter = "pid-2"
	result := EvaluateLeaseCanBecomeActive(lease, evidence)
	if result.Next != StateQuarantined {
		t.Fatalf("expected QUARANTINED, got %s", result.Next)
	}
	assertHas(t, result.Blockers, BlockerQemuPIDChanged)
}

func TestEvaluateReturnToFloorReadinessBlocksUnsafeMemoryReturn(t *testing.T) {
	lease := testLease()
	lease.Status.ReturnToFloorReady = false
	result := EvaluateReturnToFloorReadiness(lease)
	if result.Allowed {
		t.Fatal("return-to-floor must be blocked when not ready")
	}
	assertHas(t, result.Blockers, BlockerReturnToFloorNotReady)
	assertHas(t, result.Blockers, BlockerMemoryReturnNotSafe)
}

func TestEvaluateWindowsFluidReadinessBlockedWhenEvidenceMissing(t *testing.T) {
	shell := testShell()
	evidence := testEvidence()
	evidence.LastBootBefore = ""
	result := EvaluateWindowsFluidReadiness(shell, evidence)
	if result.Allowed {
		t.Fatal("readiness must be blocked with missing evidence")
	}
	assertHas(t, result.Blockers, BlockerLastBootChanged)
}

func testShell() WindowsFluidShell {
	return WindowsFluidShell{
		Spec: WindowsFluidShellSpec{
			VMRef:             "vm-01",
			RuntimeMode:       "in-place-qmp",
			MigrationRequired: false,
			RebootAllowed:     false,
			RecreateAllowed:   false,
			Floor:             ResourceQuantity{CPU: 2, Memory: 4096},
			Envelope:          WindowsFluidShellEnvelope{MaxCPU: 8, MaxMemory: 32768},
			RuntimeTarget:     ResourceQuantity{CPU: 4, Memory: 8192},
			RuntimeActual:     ResourceQuantity{CPU: 4, Memory: 8192},
			FluidDevices: FluidDeviceSpec{
				CPU:    CPUFluidDeviceSpec{Mode: "qmp-hotplug", MaxCPU: 8},
				Memory: MemoryFluidDeviceSpec{Mode: "karl-fluid-memory-envelope", BlockSize: 512},
			},
			Guest: GuestSpec{
				AgentModule:            "fluidShell",
				RequireAck:             true,
				RequireNoPendingReboot: true,
			},
		},
		Status: WindowsFluidShellStatus{
			Phase:              StateReady,
			EvidenceRef:        "evidence/windows-vm-01",
			LastTransitionTime: time.Now().UTC(),
		},
	}
}

func testLease() FluidResourceLease {
	return FluidResourceLease{
		Spec: FluidResourceLeaseSpec{
			ShellRef:       "windows-shell/vm-01",
			Mode:           "in-place",
			Grant:          ResourceQuantity{CPU: 2, Memory: 4096},
			TTLSeconds:     300,
			RollbackTarget: ResourceQuantity{CPU: 2, Memory: 4096},
		},
		Guarantees: FluidResourceLeaseGuarantees{
			NoLiveMigration:  true,
			NoReboot:         true,
			NoRecreate:       true,
			SameNode:         true,
			SameVirtLauncher: true,
			SameQemuProcess:  true,
			SameMachineGUID:  true,
			SameLastBoot:     true,
			GuestAckRequired: true,
			QMPAckRequired:   true,
		},
		Status: FluidResourceLeaseStatus{
			Phase:              StateActive,
			QMPAck:             true,
			GuestAck:           true,
			LastBootUnchanged:  true,
			QemuPIDUnchanged:   true,
			RollbackReady:      true,
			ReturnToFloorReady: true,
			EvidenceRef:        "evidence/windows-vm-01",
		},
	}
}

func testEvidence() WindowsFluidEvidence {
	return WindowsFluidEvidence{
		BeforeCPU:             2,
		AfterCPU:              4,
		BeforeMemory:          4096,
		AfterMemory:           8192,
		RuntimeTarget:         ResourceQuantity{CPU: 4, Memory: 8192},
		RuntimeActual:         ResourceQuantity{CPU: 4, Memory: 8192},
		QMPEvidence:           map[string]any{"ack": true},
		GuestEvidence:         map[string]any{"guestAck": true},
		QemuPIDBefore:         "pid-1",
		QemuPIDAfter:          "pid-1",
		VirtLauncherPodBefore: "pod-a",
		VirtLauncherPodAfter:  "pod-a",
		NodeBefore:            "node-1",
		NodeAfter:             "node-1",
		LastBootBefore:        "2026-05-07T10:00:00Z",
		LastBootAfter:         "2026-05-07T10:00:00Z",
		MachineGUIDBefore:     "guid-1",
		MachineGUIDAfter:      "guid-1",
		VMIUIDBefore:          "vmi-1",
		VMIUIDAfter:           "vmi-1",
		NoReboot:              true,
		NoRecreate:            true,
		NoMigration:           true,
		RollbackResult:        map[string]any{"ready": true},
		ReturnToFloorResult:   map[string]any{"ready": true},
		BlockerList:           nil,
		Timestamps: map[string]time.Time{
			"collectedAt": time.Now().UTC(),
		},
	}
}

func assertHas(t *testing.T, items []string, expected string) {
	t.Helper()
	for _, item := range items {
		if item == expected {
			return
		}
	}
	t.Fatalf("expected %s in %v", expected, items)
}
