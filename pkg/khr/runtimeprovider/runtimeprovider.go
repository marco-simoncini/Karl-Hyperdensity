// Package runtimeprovider defines Linux provider stubs for KHR (Sprint 5).
package runtimeprovider

// LinuxSystemdProvider is a stub for systemd-backed Linux Cells (no execution).
type LinuxSystemdProvider struct {
	ID string `json:"id"`
}

func (p *LinuxSystemdProvider) Name() string {
	if p.ID == "" {
		return "linux.systemd.stub"
	}
	return p.ID
}

// LinuxCgroupEnvelopeProvider is a stub for cgroup envelope Linux Cells (no execution).
type LinuxCgroupEnvelopeProvider struct {
	ID string `json:"id"`
}

func (p *LinuxCgroupEnvelopeProvider) Name() string {
	if p.ID == "" {
		return "linux.cgroup.envelope.stub"
	}
	return p.ID
}
