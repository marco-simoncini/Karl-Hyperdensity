package windowsfluidvirt

import "time"

type QMPEvidence struct {
	SidecarVersion            string   `json:"sidecarVersion"`
	QMPConnected              bool     `json:"qmpConnected"`
	QMPGreetingObserved       bool     `json:"qmpGreetingObserved"`
	QMPCapabilitiesNegotiated bool     `json:"qmpCapabilitiesNegotiated"`
	QMPSocketPath             string   `json:"qmpSocketPath"`
	QemuPID                   string   `json:"qemuPid"`
	QemuProcessStartTime      string   `json:"qemuProcessStartTime"`
	CPUTopologyObserved       bool     `json:"cpuTopologyObserved"`
	MaxCPUsObserved           int64    `json:"maxCpusObserved"`
	HotpluggableCPUsObserved  bool     `json:"hotpluggableCpusObserved"`
	MemoryDevicesObserved     bool     `json:"memoryDevicesObserved"`
	MemoryBackendsObserved    bool     `json:"memoryBackendsObserved"`
	QMPCommandsExecuted       []string `json:"qmpCommandsExecuted"`
	QMPReadOnly               bool     `json:"qmpReadOnly"`
	QMPErrors                 []string `json:"qmpErrors"`
	Timestamps                []string `json:"timestamps"`
}

var QMPReadOnlyAllowlist = map[string]struct{}{
	"qmp_capabilities":          {},
	"query-status":              {},
	"query-cpus-fast":           {},
	"query-hotpluggable-cpus":   {},
	"query-memory-devices":      {},
	"query-memory-size-summary": {},
	"query-machines":            {},
	"query-version":             {},
}

var QMPForbiddenCommands = map[string]struct{}{
	"device_add":       {},
	"device_del":       {},
	"qom-set":          {},
	"set_link":         {},
	"system_powerdown": {},
	"stop":             {},
	"cont":             {},
	"quit":             {},
	"migrate":          {},
	"migrate_cancel":   {},
	"object-add":       {},
	"object-del":       {},
	"cpu-add":          {},
	"balloon":          {},
	"memsave":          {},
	"pmemsave":         {},
}

func IsQMPCommandReadOnly(command string) bool {
	_, forbidden := QMPForbiddenCommands[command]
	if forbidden {
		return false
	}
	_, allowed := QMPReadOnlyAllowlist[command]
	return allowed
}

func ValidateQmpReadiness(qmp QMPEvidence) []string {
	var blockers []string
	if !qmp.QMPConnected {
		blockers = append(blockers, BlockerQMPSocketUnavailable)
	}
	if !qmp.QMPCapabilitiesNegotiated {
		blockers = append(blockers, BlockerQMPAckMissing)
	}
	if !qmp.QMPReadOnly {
		blockers = append(blockers, BlockerHotplugErrorDetected)
	}
	for _, command := range qmp.QMPCommandsExecuted {
		if !IsQMPCommandReadOnly(command) {
			blockers = append(blockers, BlockerHotplugErrorDetected)
			break
		}
	}
	if len(qmp.QMPErrors) > 0 {
		blockers = append(blockers, BlockerQMPAckMissing)
	}
	return dedupe(blockers)
}

func NewReadOnlyQMPEvidence(sidecarVersion, socketPath string) QMPEvidence {
	return QMPEvidence{
		SidecarVersion:      sidecarVersion,
		QMPSocketPath:       socketPath,
		QMPReadOnly:         true,
		QMPCommandsExecuted: make([]string, 0, 8),
		QMPErrors:           nil,
		Timestamps:          []string{time.Now().UTC().Format(time.RFC3339)},
	}
}
