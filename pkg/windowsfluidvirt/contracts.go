package windowsfluidvirt

import "time"

type WindowsFluidPhase string

const (
	StateEmpty            WindowsFluidPhase = "EMPTY"
	StatePreflight        WindowsFluidPhase = "PREFLIGHT"
	StateReady            WindowsFluidPhase = "READY"
	StateLeasePrepared    WindowsFluidPhase = "LEASE_PREPARED"
	StateApplying         WindowsFluidPhase = "APPLYING"
	StateVerifying        WindowsFluidPhase = "VERIFYING"
	StateActive           WindowsFluidPhase = "ACTIVE"
	StateReturningToFloor WindowsFluidPhase = "RETURNING_TO_FLOOR"
	StateRolledBack       WindowsFluidPhase = "ROLLED_BACK"
	StateBlocked          WindowsFluidPhase = "BLOCKED"
	StateQuarantined      WindowsFluidPhase = "QUARANTINED"
)

type ResourceQuantity struct {
	CPU    int64 `json:"cpu"`
	Memory int64 `json:"memory"`
}

type WindowsFluidShell struct {
	Spec   WindowsFluidShellSpec   `json:"spec"`
	Status WindowsFluidShellStatus `json:"status"`
}

type WindowsFluidShellSpec struct {
	VMRef             string                    `json:"vmRef"`
	RuntimeMode       string                    `json:"runtimeMode"`
	MigrationRequired bool                      `json:"migrationRequired"`
	RebootAllowed     bool                      `json:"rebootAllowed"`
	RecreateAllowed   bool                      `json:"recreateAllowed"`
	Floor             ResourceQuantity          `json:"floor"`
	Envelope          WindowsFluidShellEnvelope `json:"envelope"`
	RuntimeTarget     ResourceQuantity          `json:"runtimeTarget"`
	RuntimeActual     ResourceQuantity          `json:"runtimeActual"`
	FluidDevices      FluidDeviceSpec           `json:"fluidDevices"`
	Guest             GuestSpec                 `json:"guest"`
}

type WindowsFluidShellEnvelope struct {
	MaxCPU    int64 `json:"maxCpu"`
	MaxMemory int64 `json:"maxMemory"`
}

type FluidDeviceSpec struct {
	CPU    CPUFluidDeviceSpec    `json:"cpu"`
	Memory MemoryFluidDeviceSpec `json:"memory"`
}

type CPUFluidDeviceSpec struct {
	Mode   string `json:"mode"`
	MaxCPU int64  `json:"maxCpu"`
}

type MemoryFluidDeviceSpec struct {
	Mode      string `json:"mode"`
	BlockSize int64  `json:"blockSize"`
}

type GuestSpec struct {
	AgentModule            string `json:"agentModule"`
	RequireAck             bool   `json:"requireAck"`
	RequireNoPendingReboot bool   `json:"requireNoPendingReboot"`
}

type WindowsFluidShellStatus struct {
	Phase              WindowsFluidPhase `json:"phase"`
	Conditions         []string          `json:"conditions"`
	Blockers           []string          `json:"blockers"`
	EvidenceRef        string            `json:"evidenceRef"`
	LastTransitionTime time.Time         `json:"lastTransitionTime"`
}

type FluidResourceLease struct {
	Spec       FluidResourceLeaseSpec       `json:"spec"`
	Guarantees FluidResourceLeaseGuarantees `json:"guarantees"`
	Status     FluidResourceLeaseStatus     `json:"status"`
}

type FluidResourceLeaseSpec struct {
	ShellRef       string           `json:"shellRef"`
	Mode           string           `json:"mode"`
	Grant          ResourceQuantity `json:"grant"`
	TTLSeconds     int64            `json:"ttlSeconds"`
	RollbackTarget ResourceQuantity `json:"rollbackTarget"`
}

type FluidResourceLeaseGuarantees struct {
	NoLiveMigration  bool `json:"noLiveMigration"`
	NoReboot         bool `json:"noReboot"`
	NoRecreate       bool `json:"noRecreate"`
	SameNode         bool `json:"sameNode"`
	SameVirtLauncher bool `json:"sameVirtLauncherPod"`
	SameQemuProcess  bool `json:"sameQemuProcess"`
	SameMachineGUID  bool `json:"sameMachineGuid"`
	SameLastBoot     bool `json:"sameLastBoot"`
	GuestAckRequired bool `json:"guestAckRequired"`
	QMPAckRequired   bool `json:"qmpAckRequired"`
}

type FluidResourceLeaseStatus struct {
	Phase              WindowsFluidPhase `json:"phase"`
	QMPAck             bool              `json:"qmpAck"`
	GuestAck           bool              `json:"guestAck"`
	LastBootUnchanged  bool              `json:"lastBootUnchanged"`
	QemuPIDUnchanged   bool              `json:"qemuPidUnchanged"`
	RollbackReady      bool              `json:"rollbackReady"`
	ReturnToFloorReady bool              `json:"returnToFloorReady"`
	Blockers           []string          `json:"blockers"`
	EvidenceRef        string            `json:"evidenceRef"`
}

type WindowsFluidEvidence struct {
	BeforeCPU             int64                `json:"beforeCpu"`
	AfterCPU              int64                `json:"afterCpu"`
	BeforeMemory          int64                `json:"beforeMemory"`
	AfterMemory           int64                `json:"afterMemory"`
	RuntimeTarget         ResourceQuantity     `json:"runtimeTarget"`
	RuntimeActual         ResourceQuantity     `json:"runtimeActual"`
	QMPEvidence           map[string]any       `json:"qmpEvidence"`
	GuestEvidence         map[string]any       `json:"guestEvidence"`
	QemuPIDBefore         string               `json:"qemuPidBefore"`
	QemuPIDAfter          string               `json:"qemuPidAfter"`
	VirtLauncherPodBefore string               `json:"virtLauncherPodBefore"`
	VirtLauncherPodAfter  string               `json:"virtLauncherPodAfter"`
	NodeBefore            string               `json:"nodeBefore"`
	NodeAfter             string               `json:"nodeAfter"`
	LastBootBefore        string               `json:"lastBootBefore"`
	LastBootAfter         string               `json:"lastBootAfter"`
	MachineGUIDBefore     string               `json:"machineGuidBefore"`
	MachineGUIDAfter      string               `json:"machineGuidAfter"`
	VMIUIDBefore          string               `json:"vmiUidBefore"`
	VMIUIDAfter           string               `json:"vmiUidAfter"`
	NoReboot              bool                 `json:"noReboot"`
	NoRecreate            bool                 `json:"noRecreate"`
	NoMigration           bool                 `json:"noMigration"`
	RollbackResult        map[string]any       `json:"rollbackResult"`
	ReturnToFloorResult   map[string]any       `json:"returnToFloorResult"`
	BlockerList           []string             `json:"blockerList"`
	Timestamps            map[string]time.Time `json:"timestamps"`
}
