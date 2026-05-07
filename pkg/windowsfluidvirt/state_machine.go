package windowsfluidvirt

type TransitionEvaluation struct {
	Allowed  bool
	Blockers []string
	Next     WindowsFluidPhase
}

func EvaluateWindowsFluidReadiness(shell WindowsFluidShell, evidence WindowsFluidEvidence) TransitionEvaluation {
	blockers := ValidateWindowsFluidShell(shell)
	blockers = append(blockers, ValidateWindowsFluidEvidence(evidence)...)
	blockers = append(blockers, EvaluateContinuityProofs(evidence)...)
	if len(blockers) > 0 {
		return TransitionEvaluation{Allowed: false, Blockers: dedupe(blockers), Next: StateBlocked}
	}
	return TransitionEvaluation{Allowed: true, Next: StateReady}
}

func EvaluateLeaseCanBecomeActive(lease FluidResourceLease, evidence WindowsFluidEvidence) TransitionEvaluation {
	blockers := ValidateFluidResourceLease(lease)
	if !lease.Status.QMPAck {
		blockers = append(blockers, BlockerQMPAckMissing)
	}
	if !lease.Status.GuestAck {
		blockers = append(blockers, BlockerGuestAckMissing)
	}
	if !lease.Status.LastBootUnchanged {
		blockers = append(blockers, BlockerLastBootChanged)
	}
	if !lease.Status.QemuPIDUnchanged {
		blockers = append(blockers, BlockerQemuPIDChanged)
	}
	if !lease.Status.RollbackReady {
		blockers = append(blockers, BlockerRollbackNotReady)
	}
	if !lease.Status.ReturnToFloorReady {
		blockers = append(blockers, BlockerReturnToFloorNotReady)
	}
	blockers = append(blockers, EvaluateContinuityProofs(evidence)...)
	if len(blockers) > 0 {
		return TransitionEvaluation{Allowed: false, Blockers: dedupe(blockers), Next: terminalPhaseForBlockers(blockers)}
	}
	return TransitionEvaluation{Allowed: true, Next: StateActive}
}

func EvaluateReturnToFloorReadiness(lease FluidResourceLease) TransitionEvaluation {
	var blockers []string
	if !lease.Status.QMPAck {
		blockers = append(blockers, BlockerQMPAckMissing)
	}
	if !lease.Status.GuestAck {
		blockers = append(blockers, BlockerGuestAckMissing)
	}
	if !lease.Status.QemuPIDUnchanged {
		blockers = append(blockers, BlockerQemuPIDChanged)
	}
	if !lease.Status.RollbackReady {
		blockers = append(blockers, BlockerRollbackNotReady)
	}
	if !lease.Status.ReturnToFloorReady {
		blockers = append(blockers, BlockerReturnToFloorNotReady, BlockerMemoryReturnNotSafe)
	}
	if len(blockers) > 0 {
		return TransitionEvaluation{Allowed: false, Blockers: dedupe(blockers), Next: terminalPhaseForBlockers(blockers)}
	}
	return TransitionEvaluation{Allowed: true, Next: StateEmpty}
}

func EvaluateContinuityProofs(evidence WindowsFluidEvidence) []string {
	var blockers []string
	if evidence.QemuPIDBefore == "" || evidence.QemuPIDAfter == "" || evidence.QemuPIDBefore != evidence.QemuPIDAfter {
		blockers = append(blockers, BlockerQemuPIDChanged)
	}
	if evidence.LastBootBefore == "" || evidence.LastBootAfter == "" || evidence.LastBootBefore != evidence.LastBootAfter {
		blockers = append(blockers, BlockerLastBootChanged)
	}
	if evidence.MachineGUIDBefore == "" || evidence.MachineGUIDAfter == "" || evidence.MachineGUIDBefore != evidence.MachineGUIDAfter {
		blockers = append(blockers, BlockerMachineGUIDChanged)
	}
	if evidence.NodeBefore == "" || evidence.NodeAfter == "" || evidence.NodeBefore != evidence.NodeAfter {
		blockers = append(blockers, BlockerNodeChanged)
	}
	if evidence.VirtLauncherPodBefore == "" || evidence.VirtLauncherPodAfter == "" || evidence.VirtLauncherPodBefore != evidence.VirtLauncherPodAfter {
		blockers = append(blockers, BlockerVirtLauncherPodChanged)
	}
	if !evidence.NoReboot {
		blockers = append(blockers, BlockerLastBootChanged)
	}
	if !evidence.NoMigration {
		blockers = append(blockers, BlockerLiveMigrationRequired)
	}
	return dedupe(blockers)
}

func terminalPhaseForBlockers(blockers []string) WindowsFluidPhase {
	for _, blocker := range blockers {
		def, ok := LookupBlocker(blocker)
		if !ok {
			continue
		}
		if def.ResultingPhase == PhaseQuarantined {
			return StateQuarantined
		}
	}
	return StateBlocked
}

func dedupe(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
