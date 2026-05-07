package windowsfluidvirt

type ActionType string

const (
	ActionCertifyShell       ActionType = "certify-shell"
	ActionPrepareCPULease    ActionType = "prepare-cpu-lease"
	ActionPrepareMemoryLease ActionType = "prepare-memory-lease"
	ActionEvidenceRefresh    ActionType = "evidence-refresh"
	ActionQuarantine         ActionType = "quarantine"
	ActionBlocked            ActionType = "blocked"
)

type DryRunState string

const (
	DryRunReady       DryRunState = "ready"
	DryRunBlocked     DryRunState = "blocked"
	DryRunQuarantined DryRunState = "quarantined"
	DryRunIncomplete  DryRunState = "incomplete"
)

type WindowsFluidActionSlate struct {
	ActionID           string      `json:"actionId"`
	ShellRef           string      `json:"shellRef"`
	ActionType         ActionType  `json:"actionType"`
	DryRunState        DryRunState `json:"dryRunState"`
	RuntimeMode        string      `json:"runtimeMode"`
	MutationAllowed    bool        `json:"mutationAllowed"`
	ApplyAllowed       bool        `json:"applyAllowed"`
	RollbackReady      bool        `json:"rollbackReady"`
	ReturnToFloorReady bool        `json:"returnToFloorReady"`
	QMPReady           bool        `json:"qmpReady"`
	GuestAckReady      bool        `json:"guestAckReady"`
	NoRebootProof      bool        `json:"noRebootProof"`
	SameQemuProof      bool        `json:"sameQemuProof"`
	SameNodeProof      bool        `json:"sameNodeProof"`
	SamePodProof       bool        `json:"samePodProof"`
	Blockers           []string    `json:"blockers"`
	EvidenceRefs       []string    `json:"evidenceRefs"`
	ValidForSeconds    int64       `json:"validForSeconds"`
	CreatedAt          string      `json:"createdAt"`
}

func BuildDryRunActionSlate(
	actionID string,
	shell WindowsFluidShell,
	phase WindowsFluidPhase,
	blockers []string,
	conditions map[string]bool,
	leaseIntent *DryRunLeaseIntent,
) WindowsFluidActionSlate {
	actionType := ActionEvidenceRefresh
	switch phase {
	case StateReady:
		actionType = ActionCertifyShell
	case StateLeasePrepared:
		if leaseIntent != nil && leaseIntent.ActionType == string(ActionPrepareMemoryLease) {
			actionType = ActionPrepareMemoryLease
		} else {
			actionType = ActionPrepareCPULease
		}
	case StateQuarantined:
		actionType = ActionQuarantine
	case StateBlocked:
		actionType = ActionBlocked
	}

	dryRunState := DryRunReady
	switch phase {
	case StateBlocked:
		dryRunState = DryRunBlocked
	case StateQuarantined:
		dryRunState = DryRunQuarantined
	case StatePreflight:
		dryRunState = DryRunIncomplete
	}

	return WindowsFluidActionSlate{
		ActionID:           actionID,
		ShellRef:           shell.Spec.VMRef,
		ActionType:         actionType,
		DryRunState:        dryRunState,
		RuntimeMode:        "in-place-qmp",
		MutationAllowed:    false,
		ApplyAllowed:       false,
		RollbackReady:      conditions["rollbackReady"],
		ReturnToFloorReady: conditions["returnToFloorReady"],
		QMPReady:           conditions["qmpReady"],
		GuestAckReady:      conditions["guestAckReady"],
		NoRebootProof:      conditions["noRebootProof"],
		SameQemuProof:      conditions["sameQemuProcess"],
		SameNodeProof:      conditions["sameNode"],
		SamePodProof:       conditions["sameVirtLauncherPod"],
		Blockers:           dedupe(blockers),
		EvidenceRefs:       []string{shell.Status.EvidenceRef},
		ValidForSeconds:    300,
		CreatedAt:          shell.Status.LastTransitionTime.UTC().Format("2006-01-02T15:04:05Z"),
	}
}
