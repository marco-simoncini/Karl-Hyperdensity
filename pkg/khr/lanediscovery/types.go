package lanediscovery

// Lane identifiers (KHR-Q).
const (
	LaneNativeLive             = "native-live"
	LaneLinuxContainerCgroup   = "linux-container-cgroup"
	LaneLinuxVMCompatibility   = "linux-vm-compatibility"
	LaneWindowsVMSession       = "windows-vm-session"
	LaneKubevirtCompatibility  = "kubevirt-compatibility"

	LabelNativeLive = "khr.karl.io/native-live"
)

// Classification values for discovered workloads.
const (
	ClassificationNativeLive            = "native-live"
	ClassificationLiveInPlaceCapable    = "live-in-place-capable"
	ClassificationObservationOnly       = "observation-only"
	ClassificationCompatibilityFallback = "compatibility-fallback"
	ClassificationUnsupported         = "unsupported"
)

const ModeLaneDiscovery = "lane-discovery"

// DiscoveredHost is a cluster node observed read-only.
type DiscoveredHost struct {
	HostID           string            `json:"hostId"`
	NodeName         string            `json:"nodeName"`
	Provider         string            `json:"provider"`
	RuntimeMode      string            `json:"runtimeMode,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`
	Ready            bool              `json:"ready"`
	ObservedAt       string            `json:"observedAt,omitempty"`
}

// DiscoveredShell is a projected Shell observation.
type DiscoveredShell struct {
	Ref           string `json:"ref"`
	Namespace     string `json:"namespace"`
	Name          string `json:"name"`
	RuntimeClass  string `json:"runtimeClass,omitempty"`
	OSFamily      string `json:"osFamily,omitempty"`
	VMType        string `json:"vmType,omitempty"`
	ProviderBinding string `json:"providerBinding,omitempty"`
}

// DiscoveredCell is a projected Cell observation.
type DiscoveredCell struct {
	Ref             string `json:"ref"`
	ShellRef        string `json:"shellRef"`
	Namespace       string `json:"namespace"`
	Name            string `json:"name"`
	VMType          string `json:"vmType,omitempty"`
	OSFamily        string `json:"osFamily,omitempty"`
	SessionType     string `json:"sessionType,omitempty"`
	NodeName        string `json:"nodeName,omitempty"`
	Running         bool   `json:"running"`
	ProviderBinding string `json:"providerBinding,omitempty"`
}

// DiscoveredResourcePort is a projected or cluster-observed ResourcePort.
type DiscoveredResourcePort struct {
	Ref                        string `json:"ref"`
	ShellRef                   string `json:"shellRef"`
	CellRef                    string `json:"cellRef"`
	Lane                       string `json:"lane"`
	ProviderBinding            string `json:"providerBinding"`
	Classification             string `json:"classification"`
	LiveScaleCapabilityObserved bool   `json:"liveScaleCapabilityObserved"`
	ClusterObserved            bool   `json:"clusterObserved,omitempty"`
}

// LaneCapability summarizes a lane on the cluster.
type LaneCapability struct {
	Lane                        string `json:"lane"`
	Classification              string `json:"classification"`
	ProviderBinding             string `json:"providerBinding"`
	LiveScaleCapabilityObserved bool   `json:"liveScaleCapabilityObserved"`
	WorkloadCount               int    `json:"workloadCount"`
}

// BlockedState is a read-only blocked posture (no apply).
type BlockedState struct {
	State      string `json:"state"`
	Reason     string `json:"reason"`
	TargetRef  string `json:"targetRef,omitempty"`
	Lane       string `json:"lane,omitempty"`
}

// SafetyPolicy attests read-only discovery.
type SafetyPolicy struct {
	ReadOnly        bool `json:"readOnly"`
	NoPatch         bool `json:"noPatch"`
	NoApply         bool `json:"noApply"`
	NoRestart       bool `json:"noRestart"`
	NoRollout       bool `json:"noRollout"`
	NoRecreate      bool `json:"noRecreate"`
	NoProductionApply bool `json:"noProductionMutation"`
}

// DefaultSafetyPolicy is the KHR-Q default.
func DefaultSafetyPolicy() SafetyPolicy {
	return SafetyPolicy{
		ReadOnly: true, NoPatch: true, NoApply: true,
		NoRestart: true, NoRollout: true, NoRecreate: true,
		NoProductionApply: true,
	}
}

// Result is CLI JSON for lane-discovery mode.
type Result struct {
	Mode                   string                   `json:"mode"`
	Blocked                bool                     `json:"blocked"`
	Reason                 string                   `json:"reason,omitempty"`
	ClusterContext         string                   `json:"clusterContext"`
	ObservedAt             string                   `json:"observedAt"`
	Safety                 SafetyPolicy             `json:"safety"`
	DiscoveredHosts        []DiscoveredHost         `json:"discoveredHosts"`
	DiscoveredShells       []DiscoveredShell        `json:"discoveredShells"`
	DiscoveredCells        []DiscoveredCell         `json:"discoveredCells"`
	DiscoveredResourcePorts []DiscoveredResourcePort `json:"discoveredResourcePorts"`
	LaneCapabilities       []LaneCapability         `json:"laneCapabilities"`
	BlockedStates          []BlockedState           `json:"blockedStates"`
	Summary                map[string]int           `json:"summary,omitempty"`
}
