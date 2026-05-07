package windowsfluidvirt

import "time"

type RuntimeSourceMetadata struct {
	SourceKind      string `json:"sourceKind"`
	SourceName      string `json:"sourceName"`
	SourceNamespace string `json:"sourceNamespace"`
	ClusterContext  string `json:"clusterContext"`
}

type RuntimePolicyGates struct {
	Annotations            map[string]string `json:"annotations"`
	ExpectedRuntimeMode    string            `json:"expectedRuntimeMode"`
	PoolReplicaContextOnly bool              `json:"poolReplicaContextOnly"`
}

type RuntimeDeclaredObserved struct {
	DeclaredTarget ResourceQuantity `json:"declaredTarget"`
	RuntimeActual  ResourceQuantity `json:"runtimeActual"`
	Observed       ResourceQuantity `json:"observed"`
}

type DryRunLeaseIntent struct {
	ActionType         string           `json:"actionType"`
	ShellRef           string           `json:"shellRef"`
	Grant              ResourceQuantity `json:"grant"`
	RollbackReady      bool             `json:"rollbackReady"`
	ReturnToFloorReady bool             `json:"returnToFloorReady"`
}

type WindowsFluidRuntimeEvidenceBundle struct {
	Shell              WindowsFluidShell               `json:"shell"`
	KubeVirtBefore     KubeVirtRuntimeIdentityEvidence `json:"kubeVirtBefore"`
	KubeVirtAfter      KubeVirtRuntimeIdentityEvidence `json:"kubeVirtAfter"`
	QMP                *QMPEvidence                    `json:"qmp"`
	Guest              *GuestRuntimeEvidence           `json:"guest"`
	LeaseIntent        *DryRunLeaseIntent              `json:"leaseIntent"`
	PolicyGates        RuntimePolicyGates              `json:"policyGates"`
	ObservedBlockers   []string                        `json:"observedBlockers"`
	DeclaredObserved   RuntimeDeclaredObserved         `json:"declaredObserved"`
	Timestamps         map[string]string               `json:"timestamps"`
	SourceMetadata     RuntimeSourceMetadata           `json:"sourceMetadata"`
	SanitizationStatus string                          `json:"sanitizationStatus"`
}

func NewRuntimeEvidenceBundle(vmNamespace, vmName string) WindowsFluidRuntimeEvidenceBundle {
	now := time.Now().UTC().Format(time.RFC3339)
	return WindowsFluidRuntimeEvidenceBundle{
		KubeVirtBefore:     NewRuntimeIdentityEvidence(vmNamespace, vmName),
		KubeVirtAfter:      NewRuntimeIdentityEvidence(vmNamespace, vmName),
		ObservedBlockers:   nil,
		Timestamps:         map[string]string{"collectedAt": now},
		SanitizationStatus: "sanitized",
	}
}
