package host

import (
	"os"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
)

const RuntimeVersion = "0.0.3-khr-n"

// Registration is the host identity reported to control plane (mock/local only).
type Registration struct {
	Tool              string            `json:"tool"`
	Version           string            `json:"version"`
	HostID            string            `json:"hostId"`
	Hostname          string            `json:"hostname"`
	Platform          string            `json:"platform"`
	CgroupVersion     string            `json:"cgroupVersion"`
	SandboxMode       bool              `json:"sandboxMode"`
	SandboxApply      bool              `json:"sandboxApplyEnabled"`
	AllowedNamespaces []string          `json:"allowedNamespaces,omitempty"`
	AllowedLabels     map[string]string `json:"allowedLabels,omitempty"`
	RegisteredAt      string            `json:"registeredAt"`
	ProductionSafe    bool              `json:"productionSafe"`
}

// RegisterHost builds a registration record (no network, no kube apply).
func RegisterHost(cfg *Config) Registration {
	host, _ := os.Hostname()
	return Registration{
		Tool:              "karl-host-runtime",
		Version:           RuntimeVersion,
		HostID:            cfg.Spec.HostID,
		Hostname:          host,
		Platform:          "linux",
		CgroupVersion:     string(cgroup.DetectVersion()),
		SandboxMode:       cfg.Spec.SandboxMode,
		SandboxApply:      cfg.Spec.SandboxApplyEnabled,
		AllowedNamespaces: append([]string(nil), cfg.Spec.AllowedNamespaces...),
		AllowedLabels:     copyLabels(cfg.Spec.AllowedLabels),
		RegisteredAt:      time.Now().UTC().Format(time.RFC3339),
		ProductionSafe:    false,
	}
}

func copyLabels(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
