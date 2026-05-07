package windowsfluidvirt

import "time"

type KubeVirtRuntimeIdentityEvidence struct {
	VMName                       string            `json:"vmName"`
	VMNamespace                  string            `json:"vmNamespace"`
	VMUID                        string            `json:"vmUid"`
	VMIName                      string            `json:"vmiName"`
	VMIUID                       string            `json:"vmiUid"`
	VMIPhase                     string            `json:"vmiPhase"`
	VirtLauncherPodName          string            `json:"virtLauncherPodName"`
	VirtLauncherPodUID           string            `json:"virtLauncherPodUid"`
	NodeName                     string            `json:"nodeName"`
	PodRestartCount              int64             `json:"podRestartCount"`
	ContainerIDs                 []string          `json:"containerIds"`
	QemuPID                      string            `json:"qemuPid"`
	QMPSocketPath                string            `json:"qmpSocketPath"`
	LiveMigrationObjectsObserved []string          `json:"liveMigrationObjectsObserved"`
	VMIMObjectsObserved          []string          `json:"vmimObjectsObserved"`
	MigrationRequired            bool              `json:"migrationRequired"`
	RecreateRequired             bool              `json:"recreateRequired"`
	RolloutObserved              bool              `json:"rolloutObserved"`
	Timestamps                   map[string]string `json:"timestamps"`
}

func EvaluateKubeVirtIdentityContinuity(before, after KubeVirtRuntimeIdentityEvidence) []string {
	var blockers []string
	if before.NodeName == "" || after.NodeName == "" || before.NodeName != after.NodeName {
		blockers = append(blockers, BlockerNodeChanged)
	}
	if before.VirtLauncherPodUID == "" || after.VirtLauncherPodUID == "" || before.VirtLauncherPodUID != after.VirtLauncherPodUID {
		blockers = append(blockers, BlockerVirtLauncherPodChanged)
	}
	if before.QemuPID == "" || after.QemuPID == "" || before.QemuPID != after.QemuPID {
		blockers = append(blockers, BlockerQemuPIDChanged)
	}
	if before.VMIUID == "" || after.VMIUID == "" || before.VMIUID != after.VMIUID {
		blockers = append(blockers, BlockerVMIRecreateRequired)
	}
	return dedupe(blockers)
}

func EvaluateNoMigrationProof(evidence KubeVirtRuntimeIdentityEvidence) []string {
	if evidence.MigrationRequired || len(evidence.LiveMigrationObjectsObserved) > 0 || len(evidence.VMIMObjectsObserved) > 0 {
		return []string{BlockerLiveMigrationRequired}
	}
	return nil
}

func EvaluateNoRecreateProof(evidence KubeVirtRuntimeIdentityEvidence) []string {
	var blockers []string
	if evidence.RecreateRequired {
		blockers = append(blockers, BlockerVMIRecreateRequired)
	}
	if evidence.RolloutObserved {
		blockers = append(blockers, BlockerVMIRecreateRequired)
	}
	return dedupe(blockers)
}

func RuntimeIdentitySummary(evidence KubeVirtRuntimeIdentityEvidence) map[string]any {
	return map[string]any{
		"vmRef":              evidence.VMNamespace + "/" + evidence.VMName,
		"vmiRef":             evidence.VMIName,
		"nodeName":           evidence.NodeName,
		"virtLauncherPod":    evidence.VirtLauncherPodName,
		"qemuPid":            evidence.QemuPID,
		"qmpSocketPath":      evidence.QMPSocketPath,
		"migrationObserved":  len(evidence.LiveMigrationObjectsObserved) > 0 || len(evidence.VMIMObjectsObserved) > 0,
		"recreateOrRollout":  evidence.RecreateRequired || evidence.RolloutObserved,
		"podRestartCount":    evidence.PodRestartCount,
		"evidenceTimestamps": evidence.Timestamps,
	}
}

func NewRuntimeIdentityEvidence(vmNamespace, vmName string) KubeVirtRuntimeIdentityEvidence {
	return KubeVirtRuntimeIdentityEvidence{
		VMNamespace: vmNamespace,
		VMName:      vmName,
		Timestamps: map[string]string{
			"collectedAt": time.Now().UTC().Format(time.RFC3339),
		},
	}
}
