package windowsfluidvirt

import (
	"errors"
	"testing"
)

func TestEvaluateFluidRuntimeGateReady(t *testing.T) {
	input := RuntimeGateInput{
		Annotations: RequiredRuntimeAnnotations,
		Shell:       testShell(),
		BeforeIdentity: KubeVirtRuntimeIdentityEvidence{
			VMName:              "win-01",
			VMNamespace:         "karl-system",
			VMIName:             "win-01-vmi",
			VMIUID:              "vmi-uid",
			VirtLauncherPodName: "virt-launcher-1",
			VirtLauncherPodUID:  "pod-uid",
			NodeName:            "node-1",
			QemuPID:             "1234",
		},
		AfterIdentity: KubeVirtRuntimeIdentityEvidence{
			VMName:              "win-01",
			VMNamespace:         "karl-system",
			VMIName:             "win-01-vmi",
			VMIUID:              "vmi-uid",
			VirtLauncherPodName: "virt-launcher-1",
			VirtLauncherPodUID:  "pod-uid",
			NodeName:            "node-1",
			QemuPID:             "1234",
		},
		QMP: QMPEvidence{
			QMPConnected:              true,
			QMPCapabilitiesNegotiated: true,
			QMPReadOnly:               true,
			QMPCommandsExecuted:       []string{"qmp_capabilities", "query-status"},
		},
		Guest: GuestRuntimeEvidence{
			GuestAck:               true,
			LastBootTime:           "2026-05-07T10:00:00Z",
			MachineGUIDHash:        "machine-hash",
			MemoryAdapterVerified:  true,
			ReturnToFloorReady:     true,
			CriticalEventsDetected: false,
		},
	}

	result := EvaluateFluidRuntimeGate(input)
	if result.Phase != StateReady {
		t.Fatalf("expected READY phase, got %s blockers=%v", result.Phase, result.Blockers)
	}
	if result.Classification != ClassificationReadyForFluidShellCertification {
		t.Fatalf("expected READY_FOR_FLUID_SHELL_CERTIFICATION classification, got %s", result.Classification)
	}
}

func TestEvaluateFluidRuntimeGateBlockedMissingAnnotations(t *testing.T) {
	input := RuntimeGateInput{
		Annotations: map[string]string{},
		Shell:       testShell(),
		BeforeIdentity: KubeVirtRuntimeIdentityEvidence{
			VMName:              "win-01",
			VMNamespace:         "karl-system",
			VMIName:             "win-01-vmi",
			VMIUID:              "vmi-uid",
			VirtLauncherPodName: "virt-launcher-1",
			VirtLauncherPodUID:  "pod-uid",
			NodeName:            "node-1",
			QemuPID:             "1234",
		},
		AfterIdentity: KubeVirtRuntimeIdentityEvidence{
			VMName:              "win-01",
			VMNamespace:         "karl-system",
			VMIName:             "win-01-vmi",
			VMIUID:              "vmi-uid",
			VirtLauncherPodName: "virt-launcher-1",
			VirtLauncherPodUID:  "pod-uid",
			NodeName:            "node-1",
			QemuPID:             "1234",
		},
		QMP: QMPEvidence{
			QMPConnected:              true,
			QMPCapabilitiesNegotiated: true,
			QMPReadOnly:               true,
			QMPCommandsExecuted:       []string{"qmp_capabilities"},
		},
		Guest: GuestRuntimeEvidence{
			GuestAck:              true,
			LastBootTime:          "2026-05-07T10:00:00Z",
			MachineGUIDHash:       "machine-hash",
			MemoryAdapterVerified: true,
			ReturnToFloorReady:    true,
		},
	}
	result := EvaluateFluidRuntimeGate(input)
	if result.Phase != StateBlocked {
		t.Fatalf("expected BLOCKED phase, got %s", result.Phase)
	}
}

func TestEvaluateFluidRuntimeGateBlockedMissingQMP(t *testing.T) {
	input := RuntimeGateInput{
		Annotations: RequiredRuntimeAnnotations,
		Shell:       testShell(),
		BeforeIdentity: KubeVirtRuntimeIdentityEvidence{
			VMName:              "win-01",
			VMNamespace:         "karl-system",
			VMIName:             "win-01-vmi",
			VMIUID:              "vmi-uid",
			VirtLauncherPodName: "virt-launcher-1",
			VirtLauncherPodUID:  "pod-uid",
			NodeName:            "node-1",
			QemuPID:             "1234",
		},
		AfterIdentity: KubeVirtRuntimeIdentityEvidence{
			VMName:              "win-01",
			VMNamespace:         "karl-system",
			VMIName:             "win-01-vmi",
			VMIUID:              "vmi-uid",
			VirtLauncherPodName: "virt-launcher-1",
			VirtLauncherPodUID:  "pod-uid",
			NodeName:            "node-1",
			QemuPID:             "1234",
		},
		QMP: QMPEvidence{
			QMPConnected:              false,
			QMPCapabilitiesNegotiated: false,
			QMPReadOnly:               true,
			QMPCommandsExecuted:       []string{},
		},
		Guest: GuestRuntimeEvidence{
			GuestAck:              true,
			LastBootTime:          "2026-05-07T10:00:00Z",
			MachineGUIDHash:       "machine-hash",
			MemoryAdapterVerified: true,
			ReturnToFloorReady:    true,
		},
	}

	result := EvaluateFluidRuntimeGate(input)
	if result.Phase != StateBlocked {
		t.Fatalf("expected BLOCKED phase with missing qmp, got %s", result.Phase)
	}
	assertHas(t, result.Blockers, BlockerQMPSocketUnavailable)
}

func TestEvaluateNoMigrationProof(t *testing.T) {
	blockers := EvaluateNoMigrationProof(KubeVirtRuntimeIdentityEvidence{
		MigrationRequired: true,
	})
	assertHas(t, blockers, BlockerLiveMigrationRequired)
}

func TestEvaluateQmpReadinessMutatingRejected(t *testing.T) {
	blockers := ValidateQmpReadiness(QMPEvidence{
		QMPConnected:              true,
		QMPCapabilitiesNegotiated: true,
		QMPReadOnly:               true,
		QMPCommandsExecuted:       []string{"device_add"},
	})
	assertHas(t, blockers, BlockerHotplugErrorDetected)
}

func TestEvaluateGuestReadinessBlockers(t *testing.T) {
	blockers := EvaluateGuestReadiness(GuestRuntimeEvidence{
		GuestAck:              false,
		PendingReboot:         true,
		LastBootTime:          "",
		MachineGUIDHash:       "",
		MemoryAdapterVerified: false,
		ReturnToFloorReady:    false,
	})
	assertHas(t, blockers, BlockerGuestAckMissing)
	assertHas(t, blockers, BlockerPendingRebootDetected)
	assertHas(t, blockers, BlockerMachineGUIDChanged)
	assertHas(t, blockers, BlockerLastBootChanged)
	assertHas(t, blockers, BlockerMemoryDriverUnverified)
	assertHas(t, blockers, BlockerReturnToFloorNotReady)
}

func TestGuestEvidencePendingRebootMapsToBlocker(t *testing.T) {
	blockers := EvaluateGuestReadiness(GuestRuntimeEvidence{
		GuestAck:              true,
		PendingReboot:         true,
		LastBootTime:          "2026-05-07T10:00:00Z",
		MachineGUIDHash:       "machine-hash",
		MemoryAdapterVerified: true,
		ReturnToFloorReady:    true,
	})
	assertHas(t, blockers, BlockerPendingRebootDetected)
}

func TestGuestEvidenceMissingMachineGuidMapsToBlocker(t *testing.T) {
	blockers := EvaluateGuestReadiness(GuestRuntimeEvidence{
		GuestAck:              true,
		LastBootTime:          "2026-05-07T10:00:00Z",
		MachineGUIDHash:       "",
		MemoryAdapterVerified: true,
		ReturnToFloorReady:    true,
	})
	assertHas(t, blockers, BlockerMachineGUIDChanged)
}

func TestGuestEvidenceMissingLastBootMapsToBlocker(t *testing.T) {
	blockers := EvaluateGuestReadiness(GuestRuntimeEvidence{
		GuestAck:              true,
		LastBootTime:          "",
		MachineGUIDHash:       "machine-hash",
		MemoryAdapterVerified: true,
		ReturnToFloorReady:    true,
	})
	assertHas(t, blockers, BlockerLastBootChanged)
}

func TestGuestEvidenceMemoryAdapterUnverifiedMapsToBlocker(t *testing.T) {
	blockers := EvaluateGuestReadiness(GuestRuntimeEvidence{
		GuestAck:              true,
		LastBootTime:          "2026-05-07T10:00:00Z",
		MachineGUIDHash:       "machine-hash",
		MemoryAdapterVerified: false,
		ReturnToFloorReady:    true,
	})
	assertHas(t, blockers, BlockerMemoryDriverUnverified)
}

func TestGuestEvidenceReturnToFloorFalseMapsToBlocker(t *testing.T) {
	blockers := EvaluateGuestReadiness(GuestRuntimeEvidence{
		GuestAck:              true,
		LastBootTime:          "2026-05-07T10:00:00Z",
		MachineGUIDHash:       "machine-hash",
		MemoryAdapterVerified: true,
		ReturnToFloorReady:    false,
	})
	assertHas(t, blockers, BlockerReturnToFloorNotReady)
}

func TestGuestEvidenceGuestAckFalseMapsToBlocker(t *testing.T) {
	blockers := EvaluateGuestReadiness(GuestRuntimeEvidence{
		GuestAck:              false,
		LastBootTime:          "2026-05-07T10:00:00Z",
		MachineGUIDHash:       "machine-hash",
		MemoryAdapterVerified: true,
		ReturnToFloorReady:    true,
	})
	assertHas(t, blockers, BlockerGuestAckMissing)
}

func TestQmpEvidenceErrorMapsToAckMissing(t *testing.T) {
	blockers := ValidateQmpReadiness(QMPEvidence{
		QMPConnected:              true,
		QMPCapabilitiesNegotiated: true,
		QMPReadOnly:               true,
		QMPCommandsExecuted:       []string{"qmp_capabilities", "query-status"},
		QMPErrors:                 []string{errors.New("broken qmp command").Error()},
	})
	assertHas(t, blockers, BlockerQMPAckMissing)
}
