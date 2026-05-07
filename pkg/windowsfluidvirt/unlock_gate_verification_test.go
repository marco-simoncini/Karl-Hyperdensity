package windowsfluidvirt

import (
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestUnlockGateFixtureMatrix(t *testing.T) {
	fixtures := []string{
		"gate0-executor-hard-disabled.passed.json",
		"gate0-executor-flag-true.failed.json",
		"gate1-master-win11-readonly-complete.passed.json",
		"gate1-master-win11-missing-qmp.blocked.json",
		"gate1-master-win11-missing-guest.blocked.json",
		"gate1-master-win11-identity-drift.quarantined.json",
		"gate1-win11-pool-target.blocked.json",
		"gate2-attestation-future-signable.passed.json",
		"gate2-attestation-malformed.blocked.json",
		"gate2-attestation-replayed-stale.blocked.json",
	}
	for _, fixtureName := range fixtures {
		t.Run(fixtureName, func(t *testing.T) {
			fixture, err := LoadUnlockGateReplayFixture(unlockGateFixtureAbsPath(t, fixtureName))
			if err != nil {
				t.Fatalf("load fixture: %v", err)
			}
			result := EvaluateWindowsFluidUnlockGate(UnlockGateEvaluationInput{
				GateID:             fixture.GateID,
				EvidenceBundle:     fixture.EvidenceBundle,
				GovernanceContract: fixture.GovernanceContract,
				ExecutorOutput:     fixture.ExecutorOutput,
				Attestation:        fixture.Attestation,
				ParityEvidence:     fixture.ParityEvidence,
				EvaluationTime:     time.Date(2026, 5, 7, 18, 20, 0, 0, time.UTC),
			})
			if result.GateStatus != fixture.ExpectedGateStatus {
				t.Fatalf("gate status mismatch expected=%s got=%s blockers=%v", fixture.ExpectedGateStatus, result.GateStatus, result.BlockerList)
			}
			for _, blocker := range fixture.ExpectedBlockers {
				assertHas(t, result.BlockerList, blocker)
			}
			if !result.ExecutorMustRemainDisabled || result.ApplyAllowed || result.MutationAllowed {
				t.Fatal("gate verification must stay non-executable")
			}
		})
	}
}

func TestUnlockGateSetFixtureMatrix(t *testing.T) {
	fixtures := []string{
		"gateset-0-1-2.passed.json",
		"gateset-quarantine.identity-drift.json",
	}
	for _, fixtureName := range fixtures {
		t.Run(fixtureName, func(t *testing.T) {
			fixture, err := LoadUnlockGateReplayFixture(unlockGateFixtureAbsPath(t, fixtureName))
			if err != nil {
				t.Fatalf("load fixture: %v", err)
			}
			result := EvaluateWindowsFluidUnlockGateSet(UnlockGateSetEvaluationInput{
				EvidenceBundle:     fixture.EvidenceBundle,
				GovernanceContract: fixture.GovernanceContract,
				ExecutorOutput:     fixture.ExecutorOutput,
				Attestation:        fixture.Attestation,
				ParityEvidence:     fixture.ParityEvidence,
				EvaluationTime:     time.Date(2026, 5, 7, 18, 20, 0, 0, time.UTC),
			})
			if result.AggregateStatus != fixture.ExpectedAggregateStatus {
				t.Fatalf("aggregate mismatch expected=%s got=%s blockers=%v", fixture.ExpectedAggregateStatus, result.AggregateStatus, result.Blockers)
			}
		})
	}
}

func TestUnlockGateNegativeMatrixMapping(t *testing.T) {
	evalTime := time.Date(2026, 5, 7, 18, 20, 0, 0, time.UTC)
	cases := []struct {
		name          string
		mutate        func(*UnlockGateSetEvaluationInput)
		expectedGate  UnlockGateID
		expectedState UnlockGateStatus
		expectedBlock string
		quarantine    bool
	}{
		{"missing QMP", func(in *UnlockGateSetEvaluationInput) { in.EvidenceBundle.QMP = nil }, Gate1LabReadOnlyEvidence, UnlockGateBlocked, BlockerQMPSocketUnavailable, false},
		{"wrong QMP socket", func(in *UnlockGateSetEvaluationInput) {
			in.EvidenceBundle.QMP.QMPSocketPath = ""
			in.EvidenceBundle.QMP.QMPConnected = false
		}, Gate1LabReadOnlyEvidence, UnlockGateBlocked, BlockerQMPSocketUnavailable, false},
		{"QEMU PID changed", func(in *UnlockGateSetEvaluationInput) { in.EvidenceBundle.KubeVirtAfter.QemuPID = "9999" }, Gate1LabReadOnlyEvidence, UnlockGateQuarantined, BlockerQemuPIDChanged, true},
		{"virt-launcher pod changed", func(in *UnlockGateSetEvaluationInput) {
			in.EvidenceBundle.KubeVirtAfter.VirtLauncherPodUID = "pod-changed"
		}, Gate1LabReadOnlyEvidence, UnlockGateQuarantined, BlockerVirtLauncherPodChanged, true},
		{"node changed", func(in *UnlockGateSetEvaluationInput) { in.EvidenceBundle.KubeVirtAfter.NodeName = "other-node" }, Gate1LabReadOnlyEvidence, UnlockGateQuarantined, BlockerNodeChanged, true},
		{"last boot changed", func(in *UnlockGateSetEvaluationInput) { in.EvidenceBundle.Guest.LastBootTime = "" }, Gate1LabReadOnlyEvidence, UnlockGateQuarantined, BlockerLastBootChanged, true},
		{"machine identity changed", func(in *UnlockGateSetEvaluationInput) { in.EvidenceBundle.Guest.MachineGUIDHash = "" }, Gate1LabReadOnlyEvidence, UnlockGateQuarantined, BlockerMachineGUIDChanged, true},
		{"pending reboot", func(in *UnlockGateSetEvaluationInput) { in.EvidenceBundle.Guest.PendingReboot = true }, Gate1LabReadOnlyEvidence, UnlockGateBlocked, BlockerPendingRebootDetected, false},
		{"guest ACK missing", func(in *UnlockGateSetEvaluationInput) { in.EvidenceBundle.Guest.GuestAck = false }, Gate1LabReadOnlyEvidence, UnlockGateBlocked, BlockerGuestAckMissing, false},
		{"agent module missing", func(in *UnlockGateSetEvaluationInput) {
			in.EvidenceBundle.PolicyGates.Annotations[AnnotationFluidRuntime] = "false"
		}, Gate1LabReadOnlyEvidence, UnlockGateBlocked, BlockerKarlAgentFluidModuleMissing, false},
		{"rollback not ready", func(in *UnlockGateSetEvaluationInput) { in.EvidenceBundle.LeaseIntent.RollbackReady = false }, Gate1LabReadOnlyEvidence, UnlockGateBlocked, BlockerRollbackNotReady, false},
		{"return-to-floor not ready", func(in *UnlockGateSetEvaluationInput) { in.EvidenceBundle.LeaseIntent.ReturnToFloorReady = false }, Gate1LabReadOnlyEvidence, UnlockGateBlocked, BlockerReturnToFloorNotReady, false},
		{"memory driver unverified", func(in *UnlockGateSetEvaluationInput) { in.EvidenceBundle.Guest.MemoryAdapterVerified = false }, Gate1LabReadOnlyEvidence, UnlockGateBlocked, BlockerMemoryDriverUnverified, false},
		{"critical Windows event", func(in *UnlockGateSetEvaluationInput) { in.EvidenceBundle.Guest.CriticalEventsDetected = true }, Gate1LabReadOnlyEvidence, UnlockGateQuarantined, BlockerCriticalWindowsEventDetected, true},
		{"LiveMigration object observed", func(in *UnlockGateSetEvaluationInput) {
			in.EvidenceBundle.KubeVirtAfter.LiveMigrationObjectsObserved = []string{"lm-1"}
		}, Gate1LabReadOnlyEvidence, UnlockGateBlocked, BlockerLiveMigrationRequired, false},
		{"VMIM observed", func(in *UnlockGateSetEvaluationInput) {
			in.EvidenceBundle.KubeVirtAfter.VMIMObjectsObserved = []string{"vmim-1"}
		}, Gate1LabReadOnlyEvidence, UnlockGateBlocked, BlockerLiveMigrationRequired, false},
		{"pool replica target", func(in *UnlockGateSetEvaluationInput) {
			in.EvidenceBundle.SourceMetadata.SourceName = "win11-pool-a"
			in.EvidenceBundle.PolicyGates.PoolReplicaContextOnly = true
		}, Gate1LabReadOnlyEvidence, UnlockGateBlocked, GateBlockerPoolContextOnly, false},
		{"generic Windows VM", func(in *UnlockGateSetEvaluationInput) {
			in.EvidenceBundle.PolicyGates.Annotations = map[string]string{}
		}, Gate1LabReadOnlyEvidence, UnlockGateBlocked, GateBlockerGenericWindowsVM, false},
		{"stale evidence", func(in *UnlockGateSetEvaluationInput) {
			in.EvidenceBundle.Timestamps["collectedAt"] = "2026-05-07T17:30:00Z"
		}, Gate1LabReadOnlyEvidence, UnlockGateBlocked, GateBlockerEvidenceStale, false},
		{"kill switch missing", func(in *UnlockGateSetEvaluationInput) {
			in.ExecutorOutput.KillSwitchSnapshot.KillSwitchID = ""
			in.ExecutorOutput.KillSwitchSnapshot.Enabled = false
		}, Gate0ExecutorHardDisabled, UnlockGateFailed, GateBlockerKillSwitchMissing, false},
		{"attestation missing", func(in *UnlockGateSetEvaluationInput) { in.Attestation = nil }, Gate2FutureSignableAttestation, UnlockGateBlocked, GateBlockerAttestationMissing, false},
		{"attestation malformed", func(in *UnlockGateSetEvaluationInput) {
			in.Attestation.Signature.Mode = "invalid-mode"
			in.Attestation.Signature.Value = "not-empty"
		}, Gate2FutureSignableAttestation, UnlockGateBlocked, GateBlockerAttestationMalformed, false},
		{"replayed old evidence", func(in *UnlockGateSetEvaluationInput) { in.Attestation.DecisionSnapshot["replayDetected"] = true }, Gate2FutureSignableAttestation, UnlockGateBlocked, GateBlockerAttestationReplayed, false},
		{"executor flag accidentally true", func(in *UnlockGateSetEvaluationInput) { in.ExecutorOutput.PreApplyGuard.ExecutorEnabled = true }, Gate0ExecutorHardDisabled, UnlockGateFailed, GateBlockerExecutorEnabledFlagTrue, false},
		{"mutationAllowed accidentally true", func(in *UnlockGateSetEvaluationInput) { in.GovernanceContract.MutationAllowed = true }, Gate0ExecutorHardDisabled, UnlockGateFailed, GateBlockerMutationAllowedTrue, false},
		{"applyAllowed accidentally true", func(in *UnlockGateSetEvaluationInput) { in.GovernanceContract.ApplyAllowed = true }, Gate0ExecutorHardDisabled, UnlockGateFailed, GateBlockerApplyAllowedTrue, false},
		{"command envelope contains QMP command", func(in *UnlockGateSetEvaluationInput) {
			in.ExecutorOutput.CommandEnvelope.QMPCommands = []string{"query-status"}
		}, Gate0ExecutorHardDisabled, UnlockGateFailed, GateBlockerEnvelopeNotEmpty, false},
		{"parity evidence missing", func(in *UnlockGateSetEvaluationInput) {
			in.ParityEvidence = nil
		}, GateHyperdensityParityComplete, UnlockGateBlocked, GateBlockerParityEvidenceMissing, false},
		{"parity partial success cpu down failed", func(in *UnlockGateSetEvaluationInput) {
			in.ParityEvidence.CPUScaleDown.GuestConfirmedActualState = false
		}, GateHyperdensityParityComplete, UnlockGateBlocked, GateBlockerHyperdensityParityPartialSuccessNotTotalFeasibility, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := baseUnlockGateSetInput(evalTime)
			tc.mutate(&input)
			setResult := EvaluateWindowsFluidUnlockGateSet(input)
			var gate WindowsFluidUnlockGateVerification
			found := false
			for _, candidate := range setResult.Gates {
				if candidate.GateID == tc.expectedGate {
					gate = candidate
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected gate %s not found", tc.expectedGate)
			}
			if gate.GateStatus != tc.expectedState {
				t.Fatalf("expected gate status %s got %s blockers=%v", tc.expectedState, gate.GateStatus, gate.BlockerList)
			}
			assertHas(t, gate.BlockerList, tc.expectedBlock)
			if tc.quarantine && gate.GateStatus != UnlockGateQuarantined {
				t.Fatalf("expected quarantine status, got %s", gate.GateStatus)
			}
			for _, gr := range setResult.Gates {
				if !gr.ExecutorMustRemainDisabled || gr.MutationAllowed || gr.ApplyAllowed {
					t.Fatal("gate results must never enable execution")
				}
			}
		})
	}
}

func TestHyperdensityParityGateRequirementMatrix(t *testing.T) {
	evalTime := time.Date(2026, 5, 7, 18, 20, 0, 0, time.UTC)
	cases := []struct {
		name          string
		mutate        func(*HyperdensityParityEvidence)
		expectedState UnlockGateStatus
		expectedBlock string
	}{
		{"all four proofs complete", nil, UnlockGatePassed, ""},
		{"missing cpu_scale_up", func(e *HyperdensityParityEvidence) { e.CPUScaleUp = nil }, UnlockGateBlocked, GateBlockerParityCPUScaleUpMissing},
		{"missing cpu_scale_down", func(e *HyperdensityParityEvidence) { e.CPUScaleDown = nil }, UnlockGateBlocked, GateBlockerParityCPUScaleDownMissing},
		{"missing ram_scale_up", func(e *HyperdensityParityEvidence) { e.RAMScaleUp = nil }, UnlockGateBlocked, GateBlockerParityRAMScaleUpMissing},
		{"missing ram_scale_down", func(e *HyperdensityParityEvidence) { e.RAMScaleDown = nil }, UnlockGateBlocked, GateBlockerParityRAMScaleDownMissing},
		{"cpu_scale_up guest confirm missing", func(e *HyperdensityParityEvidence) { e.CPUScaleUp.GuestConfirmedActualState = false }, UnlockGateBlocked, GateBlockerParityCPUScaleUpFailed},
		{"ram_scale_down return-to-floor missing", func(e *HyperdensityParityEvidence) { e.RAMScaleDown.ReturnToFloorVerified = false }, UnlockGateBlocked, GateBlockerParityRAMScaleDownFailed},
		{"any proof reboot true", func(e *HyperdensityParityEvidence) { e.RAMScaleUp.NoReboot = false }, UnlockGateBlocked, GateBlockerParityRAMScaleUpFailed},
		{"any proof rollout true", func(e *HyperdensityParityEvidence) { e.CPUScaleDown.NoRollout = false }, UnlockGateBlocked, GateBlockerParityCPUScaleDownFailed},
		{"any proof recreate true", func(e *HyperdensityParityEvidence) { e.CPUScaleDown.NoRecreate = false }, UnlockGateBlocked, GateBlockerParityCPUScaleDownFailed},
		{"any proof migration true", func(e *HyperdensityParityEvidence) { e.CPUScaleUp.NoLiveMigration = false }, UnlockGateBlocked, GateBlockerParityCPUScaleUpFailed},
		{"any proof sameQemu false", func(e *HyperdensityParityEvidence) { e.CPUScaleUp.SameQEMUProcess = false }, UnlockGateBlocked, GateBlockerParityCPUScaleUpFailed},
		{"partial CPU-only success", func(e *HyperdensityParityEvidence) { e.RAMScaleUp, e.RAMScaleDown = nil, nil }, UnlockGateBlocked, GateBlockerHyperdensityParityPartialSuccessNotTotalFeasibility},
		{"CPU up/down ok RAM missing", func(e *HyperdensityParityEvidence) { e.RAMScaleUp, e.RAMScaleDown = nil, nil }, UnlockGateBlocked, GateBlockerParityRAMScaleUpMissing},
		{"RAM up ok RAM down missing", func(e *HyperdensityParityEvidence) { e.RAMScaleDown = nil }, UnlockGateBlocked, GateBlockerParityRAMScaleDownMissing},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parity := completeParityEvidence()
			if tc.mutate != nil {
				tc.mutate(parity)
			}
			result := EvaluateWindowsFluidUnlockGate(UnlockGateEvaluationInput{
				GateID:         GateHyperdensityParityComplete,
				ParityEvidence: parity,
				EvaluationTime: evalTime,
			})
			if result.GateStatus != tc.expectedState {
				t.Fatalf("expected %s got %s blockers=%v", tc.expectedState, result.GateStatus, result.BlockerList)
			}
			if tc.expectedBlock != "" {
				assertHas(t, result.BlockerList, tc.expectedBlock)
			}
			if tc.expectedState == UnlockGateBlocked {
				assertHas(t, result.BlockerList, GateBlockerHyperdensityParityPartialSuccessNotTotalFeasibility)
			}
		})
	}
}

func TestUnlockGateCLIOutputDeterministic(t *testing.T) {
	fixturePath := unlockGateFixtureAbsPath(t, "gate0-executor-hard-disabled.passed.json")
	cmd := exec.Command("go", "run", "./cmd/karl-fluid-gates", "-fixture", fixturePath, "-mode", "gate", "-evaluation-time", "2026-05-07T18:20:00Z")
	cmd.Dir = governanceRepoRoot(t)
	out1, err := cmd.Output()
	if err != nil {
		t.Fatalf("first gate cli run failed: %v", err)
	}
	cmd = exec.Command("go", "run", "./cmd/karl-fluid-gates", "-fixture", fixturePath, "-mode", "gate", "-evaluation-time", "2026-05-07T18:20:00Z")
	cmd.Dir = governanceRepoRoot(t)
	out2, err := cmd.Output()
	if err != nil {
		t.Fatalf("second gate cli run failed: %v", err)
	}
	if string(out1) != string(out2) {
		t.Fatal("gate cli output must be deterministic with fixed evaluation-time")
	}
}

func baseUnlockGateSetInput(evalTime time.Time) UnlockGateSetEvaluationInput {
	bundle := &WindowsFluidRuntimeEvidenceBundle{
		Shell: WindowsFluidShell{
			Spec: WindowsFluidShellSpec{
				VMRef:             "karl/master-win11",
				RuntimeMode:       "in-place-qmp",
				MigrationRequired: false,
				RebootAllowed:     false,
				RecreateAllowed:   false,
				Floor:             ResourceQuantity{CPU: 4, Memory: 8192},
				Envelope:          WindowsFluidShellEnvelope{MaxCPU: 8, MaxMemory: 32768},
				RuntimeTarget:     ResourceQuantity{CPU: 6, Memory: 8192},
				RuntimeActual:     ResourceQuantity{CPU: 4, Memory: 8192},
				Guest: GuestSpec{
					AgentModule:            "fluidShell",
					RequireAck:             true,
					RequireNoPendingReboot: true,
				},
			},
			Status: WindowsFluidShellStatus{
				Phase:              StateReady,
				EvidenceRef:        "evidence/gates/master-win11",
				LastTransitionTime: evalTime.Add(-5 * time.Minute),
			},
		},
		KubeVirtBefore: KubeVirtRuntimeIdentityEvidence{
			VMName:              "master-win11",
			VMNamespace:         "karl",
			VMIName:             "master-win11",
			VMIUID:              "vmi-master",
			VirtLauncherPodUID:  "pod-master",
			NodeName:            "karl-metal-01",
			QemuPID:             "5120",
			VMIMObjectsObserved: []string{},
			Timestamps:          map[string]string{"collectedAt": evalTime.Add(-5 * time.Minute).Format(time.RFC3339)},
		},
		KubeVirtAfter: KubeVirtRuntimeIdentityEvidence{
			VMName:                       "master-win11",
			VMNamespace:                  "karl",
			VMIName:                      "master-win11",
			VMIUID:                       "vmi-master",
			VirtLauncherPodUID:           "pod-master",
			NodeName:                     "karl-metal-01",
			QemuPID:                      "5120",
			LiveMigrationObjectsObserved: []string{},
			VMIMObjectsObserved:          []string{},
			MigrationRequired:            false,
			RecreateRequired:             false,
			RolloutObserved:              false,
			Timestamps:                   map[string]string{"collectedAt": evalTime.Add(-5 * time.Minute).Format(time.RFC3339)},
		},
		QMP: &QMPEvidence{
			SidecarVersion:            "v0-readonly",
			QMPConnected:              true,
			QMPGreetingObserved:       true,
			QMPCapabilitiesNegotiated: true,
			QMPSocketPath:             "/var/run/qmp.sock",
			QemuPID:                   "5120",
			CPUTopologyObserved:       true,
			MaxCPUsObserved:           8,
			HotpluggableCPUsObserved:  true,
			MemoryDevicesObserved:     true,
			MemoryBackendsObserved:    true,
			QMPCommandsExecuted:       []string{"qmp_capabilities", "query-status"},
			QMPReadOnly:               true,
			QMPErrors:                 []string{},
		},
		Guest: &GuestRuntimeEvidence{
			GuestAck:              true,
			PendingReboot:         false,
			LastBootTime:          "2026-05-01T11:00:00Z",
			MachineGUIDHash:       "hash-master-win11",
			MemoryAdapterVerified: true,
			ReturnToFloorReady:    true,
		},
		LeaseIntent: &DryRunLeaseIntent{
			ActionType:         string(ActionPrepareCPULease),
			ShellRef:           "karl/master-win11",
			Grant:              ResourceQuantity{CPU: 2, Memory: 0},
			RollbackReady:      true,
			ReturnToFloorReady: true,
		},
		PolicyGates: RuntimePolicyGates{
			Annotations: map[string]string{
				AnnotationFluidRuntime:         "true",
				AnnotationNoLiveMigration:      "true",
				AnnotationNoReboot:             "required",
				AnnotationNoRecreate:           "required",
				AnnotationRuntimeMode:          "in-place-qmp",
				AnnotationSingleNodeCompatible: "true",
			},
			ExpectedRuntimeMode:    "in-place-qmp",
			PoolReplicaContextOnly: false,
		},
		Timestamps: map[string]string{
			"collectedAt": evalTime.Add(-5 * time.Minute).Format(time.RFC3339),
		},
		SourceMetadata: RuntimeSourceMetadata{
			SourceKind:      "vm",
			SourceName:      "master-win11",
			SourceNamespace: "karl",
			ClusterContext:  "karl-metal-01@ovh",
		},
	}

	governance := &WindowsFluidApplyGovernanceContract{
		ContractID:               "gc-master-win11-gates",
		ShellRef:                 "karl/master-win11",
		RequestedAction:          GovernanceFutureCPUApply,
		GovernancePhase:          GovernanceContractPrepared,
		MutationAllowed:          false,
		ApplyAllowed:             false,
		RuntimeMode:              "in-place-qmp",
		PolicyVersion:            "windows-fluid-admission-policy-v1",
		RollbackRequirement:      "mandatory",
		ReturnToFloorRequirement: "mandatory",
		Blockers:                 []string{},
		CreatedAt:                evalTime.Add(-5 * time.Minute).Format(time.RFC3339),
	}

	executor := &FutureApplyExecutorEvaluationResult{
		ExecutionResult: DisabledFutureApplyExecutionResult{
			ExecutorID:          "karl-fluid-executor-disabled-v1",
			ShellRef:            "karl/master-win11",
			RequestedAction:     GovernanceFutureCPUApply,
			ExecutionPhase:      ExecutionHardDisabled,
			ApplyAttempted:      false,
			MutationPerformed:   false,
			QMPCommandSent:      false,
			ClusterMutationSent: false,
			Blockers:            []string{BlockerFutureApplyExecutorDisabled},
			AttestationRefs:     []string{"att-master-win11-gates"},
			CreatedAt:           evalTime.Format(time.RFC3339),
		},
		PreApplyGuard: WindowsFluidPreApplyGuard{
			GuardID:                "guard-master-win11-gates",
			GovernanceContractRef:  "gc-master-win11-gates",
			RevalidationRef:        "reval-master-win11-gates",
			AttestationRef:         "att-master-win11-gates",
			KillSwitchReady:        true,
			ExecutorEnabled:        false,
			MutationWindowOpen:     false,
			QMPMutationAllowed:     false,
			ClusterMutationAllowed: false,
			GuardPhase:             GuardReadyButExecutorDisabled,
		},
		KillSwitchSnapshot: WindowsFluidKillSwitch{
			KillSwitchID:     "ks-master-win11-gates",
			Enabled:          true,
			Source:           "policy",
			Mode:             "hard-disabled",
			Reason:           "future apply executor disabled by policy",
			RequiredForApply: true,
			ObservedAt:       evalTime.Format(time.RFC3339),
		},
		CommandEnvelope: WindowsFluidExecutorCommandEnvelope{
			EnvelopeID:                "env-master-win11-gates",
			ShellRef:                  "karl/master-win11",
			RequestedAction:           GovernanceFutureCPUApply,
			CommandClass:              "cpu-lease",
			RuntimeMode:               "in-place-qmp",
			CommandPreviewOnly:        true,
			ContainsExecutableCommand: false,
			QMPCommands:               []string{},
			ClusterMutations:          []string{},
			GuestMutations:            []string{},
			DeniedReason:              BlockerFutureApplyExecutorDisabled,
			CreatedAt:                 evalTime.Format(time.RFC3339),
		},
	}

	attestation := &WindowsFluidPolicyAttestation{
		AttestationID:   "att-master-win11-gates",
		SubjectRef:      "gc-master-win11-gates",
		SubjectType:     "governance-contract",
		PolicyVersion:   "windows-fluid-admission-policy-v1",
		EvidenceRefs:    []string{"karl/master-win11"},
		BlockerSnapshot: []string{},
		InvariantSnapshot: map[string]bool{
			"same_node_invariant": true,
		},
		DecisionSnapshot: map[string]any{
			"governancePhase": GovernanceContractPrepared,
		},
		CreatedAt: evalTime.Add(-5 * time.Minute).Format(time.RFC3339),
		Attestor: Attestor{
			ComponentName:    "karl-fluid-governance-evaluator",
			ComponentVersion: "v1",
		},
		Signature: AttestationSignature{
			Mode:  "future-signable",
			Value: "",
		},
	}

	return UnlockGateSetEvaluationInput{
		EvidenceBundle:     bundle,
		GovernanceContract: governance,
		ExecutorOutput:     executor,
		Attestation:        attestation,
		ParityEvidence: &HyperdensityParityEvidence{
			CPUScaleUp: &HyperdensityParityOperationProof{
				Operation:                ParityCPUScaleUp,
				QMPConfirmedRuntimeState: true, GuestConfirmedActualState: true, SameVM: true, SameNamespace: true, SameNode: true,
				SameVirtLauncherPod: true, SameQEMUProcess: true, SameWindowsBoot: true, SameMachineIdentity: true,
				NoReboot: true, NoRollout: true, NoRecreate: true, NoLiveMigration: true, NoDestructiveMigration: true,
				RollbackVerified: true, ReturnToFloorVerified: true, EvidenceBackedAudit: true,
			},
			CPUScaleDown: &HyperdensityParityOperationProof{
				Operation:                ParityCPUScaleDown,
				QMPConfirmedRuntimeState: true, GuestConfirmedActualState: true, SameVM: true, SameNamespace: true, SameNode: true,
				SameVirtLauncherPod: true, SameQEMUProcess: true, SameWindowsBoot: true, SameMachineIdentity: true,
				NoReboot: true, NoRollout: true, NoRecreate: true, NoLiveMigration: true, NoDestructiveMigration: true,
				RollbackVerified: true, ReturnToFloorVerified: true, EvidenceBackedAudit: true,
			},
			RAMScaleUp: &HyperdensityParityOperationProof{
				Operation:                ParityRAMScaleUp,
				QMPConfirmedRuntimeState: true, GuestConfirmedActualState: true, SameVM: true, SameNamespace: true, SameNode: true,
				SameVirtLauncherPod: true, SameQEMUProcess: true, SameWindowsBoot: true, SameMachineIdentity: true,
				NoReboot: true, NoRollout: true, NoRecreate: true, NoLiveMigration: true, NoDestructiveMigration: true,
				RollbackVerified: true, ReturnToFloorVerified: true, EvidenceBackedAudit: true,
			},
			RAMScaleDown: &HyperdensityParityOperationProof{
				Operation:                ParityRAMScaleDown,
				QMPConfirmedRuntimeState: true, GuestConfirmedActualState: true, SameVM: true, SameNamespace: true, SameNode: true,
				SameVirtLauncherPod: true, SameQEMUProcess: true, SameWindowsBoot: true, SameMachineIdentity: true,
				NoReboot: true, NoRollout: true, NoRecreate: true, NoLiveMigration: true, NoDestructiveMigration: true,
				RollbackVerified: true, ReturnToFloorVerified: true, EvidenceBackedAudit: true,
			},
		},
		EvaluationTime: evalTime,
	}
}

func unlockGateFixtureAbsPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(governanceRepoRoot(t), "examples", "windows-fluid-unlock-gate-fixtures", name)
}

func completeParityEvidence() *HyperdensityParityEvidence {
	return &HyperdensityParityEvidence{
		CPUScaleUp:   completeParityProof(ParityCPUScaleUp),
		CPUScaleDown: completeParityProof(ParityCPUScaleDown),
		RAMScaleUp:   completeParityProof(ParityRAMScaleUp),
		RAMScaleDown: completeParityProof(ParityRAMScaleDown),
	}
}

func completeParityProof(op HyperdensityParityOperation) *HyperdensityParityOperationProof {
	return &HyperdensityParityOperationProof{
		Operation:                 op,
		QMPConfirmedRuntimeState:  true,
		GuestConfirmedActualState: true,
		SameVM:                    true,
		SameNamespace:             true,
		SameNode:                  true,
		SameVirtLauncherPod:       true,
		SameQEMUProcess:           true,
		SameWindowsBoot:           true,
		SameMachineIdentity:       true,
		NoReboot:                  true,
		NoRollout:                 true,
		NoRecreate:                true,
		NoLiveMigration:           true,
		NoDestructiveMigration:    true,
		RollbackVerified:          true,
		ReturnToFloorVerified:     true,
		EvidenceBackedAudit:       true,
	}
}
