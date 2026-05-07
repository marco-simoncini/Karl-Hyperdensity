package windowsfluidvirt

func ValidateWindowsFluidShell(shell WindowsFluidShell) []string {
	var blockers []string
	if shell.Spec.RuntimeMode != "in-place-qmp" {
		blockers = append(blockers, BlockerQMPSocketUnavailable)
	}
	if shell.Spec.MigrationRequired {
		blockers = append(blockers, BlockerLiveMigrationRequired)
	}
	if shell.Spec.RebootAllowed {
		blockers = append(blockers, BlockerPendingRebootDetected)
	}
	if shell.Spec.RecreateAllowed {
		blockers = append(blockers, BlockerVMIRecreateRequired)
	}
	if shell.Spec.Guest.AgentModule != "fluidShell" {
		blockers = append(blockers, BlockerKarlAgentFluidModuleMissing)
	}
	if !shell.Spec.Guest.RequireAck {
		blockers = append(blockers, BlockerGuestAckMissing)
	}

	if !orderedQuantity(shell.Spec.Floor.CPU, shell.Spec.RuntimeTarget.CPU, shell.Spec.Envelope.MaxCPU) {
		blockers = append(blockers, BlockerCPUTopologyNotConfirmed)
	}
	if !orderedQuantity(shell.Spec.Floor.Memory, shell.Spec.RuntimeTarget.Memory, shell.Spec.Envelope.MaxMemory) {
		blockers = append(blockers, BlockerGuestMemoryNotConfirmed)
	}
	if shell.Spec.RuntimeActual.CPU > shell.Spec.Envelope.MaxCPU {
		blockers = append(blockers, BlockerCPUTopologyNotConfirmed)
	}
	if shell.Spec.RuntimeActual.Memory > shell.Spec.Envelope.MaxMemory {
		blockers = append(blockers, BlockerGuestMemoryNotConfirmed)
	}
	if shell.Status.Phase == StateReady && shell.Status.EvidenceRef == "" {
		blockers = append(blockers, BlockerGuestAckMissing)
	}

	return dedupe(blockers)
}

func ValidateFluidResourceLease(lease FluidResourceLease) []string {
	var blockers []string
	if lease.Spec.Mode != "in-place" {
		blockers = append(blockers, BlockerLiveMigrationRequired)
	}
	if lease.Spec.TTLSeconds <= 0 {
		blockers = append(blockers, BlockerRollbackNotReady)
	}
	if !lease.Guarantees.NoLiveMigration {
		blockers = append(blockers, BlockerLiveMigrationRequired)
	}
	if !lease.Guarantees.NoReboot {
		blockers = append(blockers, BlockerPendingRebootDetected)
	}
	if !lease.Guarantees.NoRecreate {
		blockers = append(blockers, BlockerVMIRecreateRequired)
	}
	if !lease.Guarantees.SameNode {
		blockers = append(blockers, BlockerNodeChanged)
	}
	if !lease.Guarantees.SameVirtLauncher {
		blockers = append(blockers, BlockerVirtLauncherPodChanged)
	}
	if !lease.Guarantees.SameQemuProcess {
		blockers = append(blockers, BlockerQemuPIDChanged)
	}
	if !lease.Guarantees.SameMachineGUID {
		blockers = append(blockers, BlockerMachineGUIDChanged)
	}
	if !lease.Guarantees.SameLastBoot {
		blockers = append(blockers, BlockerLastBootChanged)
	}
	if !lease.Guarantees.GuestAckRequired {
		blockers = append(blockers, BlockerGuestAckMissing)
	}
	if !lease.Guarantees.QMPAckRequired {
		blockers = append(blockers, BlockerQMPAckMissing)
	}
	if lease.Status.Phase == StateActive {
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
	}
	return dedupe(blockers)
}

func ValidateWindowsFluidEvidence(evidence WindowsFluidEvidence) []string {
	var blockers []string
	if evidence.LastBootBefore == "" || evidence.LastBootAfter == "" {
		blockers = append(blockers, BlockerLastBootChanged)
	}
	if evidence.QemuPIDBefore == "" || evidence.QemuPIDAfter == "" {
		blockers = append(blockers, BlockerQemuPIDChanged)
	}
	if evidence.NodeBefore == "" || evidence.NodeAfter == "" {
		blockers = append(blockers, BlockerNodeChanged)
	}
	if evidence.VirtLauncherPodBefore == "" || evidence.VirtLauncherPodAfter == "" {
		blockers = append(blockers, BlockerVirtLauncherPodChanged)
	}
	if evidence.MachineGUIDBefore == "" || evidence.MachineGUIDAfter == "" {
		blockers = append(blockers, BlockerMachineGUIDChanged)
	}

	if evidence.NoReboot != (evidence.LastBootBefore == evidence.LastBootAfter && evidence.LastBootBefore != "") {
		blockers = append(blockers, BlockerLastBootChanged)
	}
	if evidence.NodeBefore != evidence.NodeAfter {
		blockers = append(blockers, BlockerNodeChanged)
	}
	if evidence.VirtLauncherPodBefore != evidence.VirtLauncherPodAfter {
		blockers = append(blockers, BlockerVirtLauncherPodChanged)
	}
	if evidence.QemuPIDBefore != evidence.QemuPIDAfter {
		blockers = append(blockers, BlockerQemuPIDChanged)
	}
	if evidence.MachineGUIDBefore != evidence.MachineGUIDAfter {
		blockers = append(blockers, BlockerMachineGUIDChanged)
	}
	if !evidence.NoMigration {
		blockers = append(blockers, BlockerLiveMigrationRequired)
	}

	return dedupe(blockers)
}

func orderedQuantity(floor, target, envelope int64) bool {
	return floor <= target && target <= envelope
}
