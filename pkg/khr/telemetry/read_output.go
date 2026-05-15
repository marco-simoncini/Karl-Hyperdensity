package telemetry

// CellRef identifies an optional Cell document used for correlation only.
type CellRef struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace,omitempty"`
	Name       string `json:"name"`
}

// ReadTelemetryOutput is the stable JSON envelope for `read-telemetry` mode.
type ReadTelemetryOutput struct {
	Tool               string        `json:"tool"`
	Version            string        `json:"version"`
	Mode               string        `json:"mode"`
	AgentID            string        `json:"agentId"`
	TelemetryMode      string        `json:"telemetryMode"`
	CgroupPath         string        `json:"cgroupPath"`
	AllowedPathPrefix  string        `json:"allowedPathPrefix"`
	CellRef            *CellRef      `json:"cellRef,omitempty"`
	Metrics            MetricsBundle `json:"metrics"`
	Evidence           Evidence      `json:"evidence"`
	MutationsForbidden bool          `json:"mutationsForbidden"`
}

// NormalizeMetricsEmpty clears empty maps/slices for cleaner JSON omitempty.
func NormalizeMetricsEmpty(m *MetricsBundle) {
	if len(m.CPUStat) == 0 {
		m.CPUStat = nil
	}
	if len(m.MemoryEvents) == 0 {
		m.MemoryEvents = nil
	}
	if len(m.IOStat) == 0 {
		m.IOStat = nil
	}
}
