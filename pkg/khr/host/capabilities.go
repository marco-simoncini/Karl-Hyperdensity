package host

import (
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/runtimeprovider"
)

// CapabilitiesReport describes what this host can offer (Linux cgroup path only).
type CapabilitiesReport struct {
	HostID           string   `json:"hostId"`
	CgroupVersion    string   `json:"cgroupVersion"`
	RuntimeProviders []string `json:"runtimeProviders"`
	SupportedModes   []string `json:"supportedModes"`
	BlockedSurfaces  []string `json:"blockedSurfaces"`
}

// ReportCapabilities returns capability truth for ResourcePort planning.
func ReportCapabilities(cfg *Config) CapabilitiesReport {
	return CapabilitiesReport{
		HostID:        cfg.Spec.HostID,
		CgroupVersion: string(cgroup.DetectVersion()),
		RuntimeProviders: []string{
			(&runtimeprovider.LinuxSystemdProvider{}).Name(),
			(&runtimeprovider.LinuxCgroupEnvelopeProvider{}).Name(),
		},
		SupportedModes: []string{"envelope", "static"},
		BlockedSurfaces: []string{
			"kubevirt", "libvirt", "qmp", "windows", "vm-mutation", "autonomous-apply",
		},
	}
}
