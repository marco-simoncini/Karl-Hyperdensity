package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcelease"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/runtimeprovider"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/safety"
)

// Config is minimal agent configuration (YAML or JSON).
type Config struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Metadata   struct {
		Name string `json:"name" yaml:"name"`
	} `json:"metadata" yaml:"metadata"`
	Spec struct {
		AgentID       string `json:"agentId" yaml:"agentId"`
		LinuxOnly     bool   `json:"linuxOnly" yaml:"linuxOnly"`
		LogLevel      string `json:"logLevel" yaml:"logLevel"`
		CellWatchMode string `json:"cellWatchMode" yaml:"cellWatchMode"`
	} `json:"spec" yaml:"spec"`
}

// LoadConfig reads YAML or JSON from path.
func LoadConfig(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if strings.HasSuffix(strings.ToLower(path), ".json") {
		if err := json.Unmarshal(raw, cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	if err := yaml.Unmarshal(raw, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// ValidateConfig returns errors for invalid config.
func ValidateConfig(cfg *Config) []string {
	var errs []string
	if cfg == nil {
		return []string{"config is nil"}
	}
	if cfg.Spec.AgentID == "" {
		errs = append(errs, "spec.agentId is required")
	}
	if !cfg.Spec.LinuxOnly {
		errs = append(errs, "spec.linuxOnly must be true for khr-linux-agent MVP")
	}
	return errs
}

// Capabilities describes stub runtime providers and host signals.
type Capabilities struct {
	AgentID                string   `json:"agentId"`
	CgroupVersion          string   `json:"cgroupVersion"`
	RuntimeProviders       []string `json:"runtimeProviders"`
	ApplyLocked            bool     `json:"applyLocked"`
	UnsafeApplyFlagPresent bool     `json:"unsafeApplyFlagPresent"`
}

// PrintCapabilitiesJSON returns structured capability JSON.
func PrintCapabilitiesJSON(cfg *Config, allowUnsafeApply bool) ([]byte, error) {
	cp := Capabilities{
		AgentID:                cfg.Spec.AgentID,
		CgroupVersion:          string(cgroup.DetectVersion()),
		RuntimeProviders:       []string{(&runtimeprovider.LinuxSystemdProvider{}).Name(), (&runtimeprovider.LinuxCgroupEnvelopeProvider{}).Name()},
		ApplyLocked:            safety.MutationsForbidden(allowUnsafeApply),
		UnsafeApplyFlagPresent: allowUnsafeApply,
	}
	return json.MarshalIndent(cp, "", "  ")
}

// DryRunEnvelopePlan returns cgroup envelope plan JSON (never writes unless allowUnsafeApply).
func DryRunEnvelopePlan(allowUnsafeApply bool, cpuDelta, memDelta string) ([]byte, error) {
	plan := cgroup.PlanEnvelope(allowUnsafeApply, cpuDelta, memDelta)
	if !allowUnsafeApply {
		plan.WouldWrite = false
		plan.WritePaths = nil
	}
	return json.MarshalIndent(plan, "", "  ")
}

// DryRunLease evaluates lease JSON with optional port JSON and cell context.
func DryRunLease(leaseRaw, portRaw []byte, ctx *resourcelease.CellContext) ([]byte, error) {
	lease := &resourcelease.LeaseInput{}
	if err := json.Unmarshal(leaseRaw, lease); err != nil {
		return nil, fmt.Errorf("lease json: %w", err)
	}
	var port *resourcelease.ResourcePortInput
	if len(portRaw) > 0 {
		port = &resourcelease.ResourcePortInput{}
		if err := json.Unmarshal(portRaw, port); err != nil {
			return nil, fmt.Errorf("resourceport json: %w", err)
		}
	}
	out := resourcelease.DryRun(lease, port, ctx)
	return json.MarshalIndent(out, "", "  ")
}
