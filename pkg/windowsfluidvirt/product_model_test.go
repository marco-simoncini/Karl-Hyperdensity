package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMasterWin11TargetReady(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	if !target.HyperdensityReady || len(target.Blockers) != 0 {
		t.Fatalf("expected ready target, blockers=%v", target.Blockers)
	}
}

func TestCombinedCPURAMLeasePrepared(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := PrepareWindowsFluidResourceLease(target, baseCombinedLease())
	if lease.Status != LeaseStatusPrepared {
		t.Fatalf("expected prepared lease, got %s blockers=%v", lease.Status, lease.Blockers)
	}
}

func TestActionSlateContainsRequiredActions(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := PrepareWindowsFluidResourceLease(target, baseCombinedLease())
	slate := BuildWindowsFluidActionSlate(target, lease)
	required := []WindowsFluidActionType{
		ActionActuatorDryRun,
		ActionCPUEntitlementApply,
		ActionCPUReturnToFloor,
		ActionQMPBalloonApply,
		ActionRAMReturnToFloor,
		ActionGuestVerify,
		ActionAuditBundleAppend,
	}
	for _, action := range required {
		if !hasActionTypeInSlate(slate, action) {
			t.Fatalf("missing action %s", action)
		}
	}
}

func TestPoolChildTargetReady(t *testing.T) {
	target := baseReadyTarget()
	target.TargetKind = TargetKindPoolChildWindowsVM
	target = EvaluateWindowsHyperdensityTarget(target)
	if !target.HyperdensityReady {
		t.Fatalf("pool child target should be ready, blockers=%v", target.Blockers)
	}
}

func TestPoolScalingMechanismBlocked(t *testing.T) {
	target := baseReadyTarget()
	target.PoolScalingRequested = true
	target = EvaluateWindowsHyperdensityTarget(target)
	assertHas(t, target.Blockers, BlockerPoolScalingAsMechanism)
}

func TestMissingActuatorBlocked(t *testing.T) {
	target := baseReadyTarget()
	target.CPU.ActuatorRequired = false
	target = EvaluateWindowsHyperdensityTarget(target)
	assertHas(t, target.Blockers, BlockerNodeFluidActuatorUnavailable)
}

func TestMissingQMPBalloonBlocked(t *testing.T) {
	target := baseReadyTarget()
	target.Memory.QMPRequired = false
	target = EvaluateWindowsHyperdensityTarget(target)
	assertHas(t, target.Blockers, BlockerRAMBalloonUnavailable)
}

func TestMissingGuestAckBlocked(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.GuestEvidenceRef = ""
	blockers := EvaluateWindowsFluidLeasePreconditions(target, lease)
	assertHas(t, blockers, BlockerGuestAckMissing)
}

func TestMissingReturnToFloorBlocked(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.ReturnToFloorTarget = WindowsFluidLeaseRequest{}
	lease = PrepareWindowsFluidResourceLease(target, lease)
	assertHas(t, lease.Blockers, BlockerReturnToFloorNotReady)
}

func TestVCPUHotplugLeaseRejected(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.RequestsVCPUHotplug = true
	lease = PrepareWindowsFluidResourceLease(target, lease)
	assertHas(t, lease.Blockers, BlockerLeaseRequestsVCPUHotplug)
}

func TestLogicalCPUScalingClaimRejected(t *testing.T) {
	target := baseReadyTarget()
	target.LogicalCPUScalingClaimed = true
	target = EvaluateWindowsHyperdensityTarget(target)
	assertHas(t, target.Blockers, BlockerLeaseRequestsVCPUHotplug)
}

func TestVMSpecPatchLeaseRejected(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.RequestsVMSpecPatch = true
	lease = PrepareWindowsFluidResourceLease(target, lease)
	assertHas(t, lease.Blockers, BlockerLeaseRequestsVMSpecPatch)
}

func TestNoRuntimeMutationExecutedByEvaluators(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := PrepareWindowsFluidResourceLease(target, baseCombinedLease())
	slate := BuildWindowsFluidActionSlate(target, lease)
	if err := RequireNoRuntimeMutation(slate); err != nil {
		t.Fatalf("expected no runtime mutation: %v", err)
	}
}

func TestAuditBundleAppendReady(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.AuditBundleRef = "bundle://master-win11"
	blockers := EvaluateWindowsFluidAuditReadiness(lease, BuildWindowsFluidActionSlate(target, lease))
	if len(blockers) > 0 {
		t.Fatalf("expected audit ready, blockers=%v", blockers)
	}
}

func TestWindowsFluidProductFixturesParse(t *testing.T) {
	fixtures := []string{
		"master-win11-target.ready.json",
		"master-win11-combined-lease.prepared.json",
		"pool-child-target.ready.json",
		"pool-scaling-target.blocked.json",
		"missing-actuator.blocked.json",
		"missing-qmp-balloon.blocked.json",
		"missing-guest-ack.blocked.json",
		"missing-return-to-floor.blocked.json",
		"vcpu-hotplug-lease.rejected.json",
		"logical-cpu-scaling-claim.rejected.json",
		"vm-spec-patch-lease.rejected.json",
		"audit-bundle-append.ready.json",
	}
	for _, name := range fixtures {
		data, err := os.ReadFile(filepath.Join(admissionRepoRoot(t), "examples", "windows-fluid-product-fixtures", name))
		if err != nil {
			t.Fatalf("read fixture %s: %v", name, err)
		}
		var generic map[string]any
		if err := json.Unmarshal(data, &generic); err != nil {
			t.Fatalf("parse fixture %s: %v", name, err)
		}
	}
}

func baseReadyTarget() WindowsHyperdensityTarget {
	return WindowsHyperdensityTarget{
		TargetID:           "target-master-win11",
		VMRef:              "master-win11",
		Namespace:          "karl",
		TargetKind:         TargetKindStandaloneWindowsVM,
		RuntimeMode:        RuntimeModePrearmedFluidEnvelopeV2,
		CompliancePhase:    ComplianceHyperdensityReadyWindowsShell,
		NodeName:           "karl-lab-metal-01",
		VirtLauncherPodRef: "virt-launcher-master-win11-kmwgg",
		PodUID:             "pod-uid",
		QemuPID:            "96",
		QemuStartTime:      "Thu May 7 18:58:03 2026",
		MachineGuidHash:    "hash",
		LastBootTime:       "/Date(1778180311500)/",
		CPU: WindowsCPUEnvelope{
			Mechanism:        "cgroup-v2-cpu-max",
			FloorCPUMax:      "300000 100000",
			CeilingCPUMax:    "600000 100000",
			CurrentCPUMax:    "600000 100000",
			ActuatorRequired: true,
		},
		Memory: WindowsMemoryEnvelope{
			Mechanism:    "qmp-balloon",
			FloorBytes:   12884901888,
			CeilingBytes: 13958643712,
			CurrentBytes: 12884901888,
			QMPRequired:  true,
		},
		Guest: WindowsGuestRequirements{
			FluidShellRequired: true,
			GuestAckRequired:   true,
		},
		Guarantees: WindowsRuntimeGuarantees{
			NoReboot:              true,
			NoRecreate:            true,
			NoRollout:             true,
			NoLiveMigration:       true,
			SameQEMU:              true,
			SameBoot:              true,
			RollbackRequired:      true,
			ReturnToFloorRequired: true,
		},
		EvidenceRefs: []string{"artifact://product-path"},
	}
}

func baseCombinedLease() WindowsFluidResourceLease {
	return WindowsFluidResourceLease{
		LeaseID:   "lease-master-win11-combined",
		TargetRef: "target-master-win11",
		LeaseKind: LeaseKindCombinedEnvelope,
		Requested: WindowsFluidLeaseRequest{
			CPUMax:      "600000 100000",
			MemoryBytes: 13958643712,
		},
		Previous: WindowsFluidLeaseRequest{
			CPUMax:      "300000 100000",
			MemoryBytes: 12884901888,
		},
		RollbackTarget: WindowsFluidLeaseRequest{
			CPUMax:      "300000 100000",
			MemoryBytes: 12884901888,
		},
		ReturnToFloorTarget: WindowsFluidLeaseRequest{
			CPUMax:      "300000 100000",
			MemoryBytes: 12884901888,
		},
		TTLSeconds:         int64((10 * time.Minute).Seconds()),
		Reason:             "prepare combined lease",
		Risk:               "medium",
		PolicySnapshot:     map[string]any{"phase": "ready"},
		ActionSlateRef:     "slate://master-win11",
		ActuatorRequestRef: "actuator://master-win11",
		QMPRequestRef:      "qmp://master-win11",
		GuestEvidenceRef:   "guest://ack",
		AuditBundleRef:     "bundle://master-win11",
		Status:             LeaseStatusPrepared,
		EvidenceRefs:       []string{"evidence://lease"},
	}
}

func hasActionTypeInSlate(slate WindowsFluidActionSlate, actionType WindowsFluidActionType) bool {
	for _, action := range slate.Actions {
		if action.ActionType == actionType {
			return true
		}
	}
	return false
}
