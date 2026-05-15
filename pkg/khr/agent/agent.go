package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/audit"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcelease"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/runtimeprovider"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/safety"
)

// AgentVersion is embedded in JSON outputs (Sprint 6).
const AgentVersion = "0.0.1-sprint6"

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
	Tool                    string         `json:"tool"`
	Version                 string         `json:"version"`
	AgentID                 string         `json:"agentId"`
	CgroupVersion           string         `json:"cgroupVersion"`
	RuntimeProviders        []string       `json:"runtimeProviders"`
	MutationsForbidden      bool           `json:"mutationsForbidden"`
	UnsafeApplyFlagPresent  bool           `json:"unsafeApplyFlagPresent"`
	FutureApplyGateRequired bool           `json:"futureApplyGateRequired"`
	Audit                   []audit.Record `json:"audit,omitempty"`
}

// PrintCapabilitiesJSON returns structured capability JSON.
func PrintCapabilitiesJSON(cfg *Config, allowUnsafeApply bool) ([]byte, error) {
	var aud []audit.Record
	if safety.UnsafeApplyRequested(allowUnsafeApply) {
		aud = append(aud, audit.UnsafeApplyFlagNonOperational())
	}
	cp := Capabilities{
		Tool:                    "khr-linux-agent",
		Version:                 AgentVersion,
		AgentID:                 cfg.Spec.AgentID,
		CgroupVersion:           string(cgroup.DetectVersion()),
		RuntimeProviders:        []string{(&runtimeprovider.LinuxSystemdProvider{}).Name(), (&runtimeprovider.LinuxCgroupEnvelopeProvider{}).Name()},
		MutationsForbidden:      safety.MutationsForbidden(allowUnsafeApply),
		UnsafeApplyFlagPresent:  safety.UnsafeApplyRequested(allowUnsafeApply),
		FutureApplyGateRequired: safety.UnsafeApplyRequested(allowUnsafeApply),
		Audit:                   aud,
	}
	return json.MarshalIndent(cp, "", "  ")
}

// DryRunCLIResult is the stable JSON envelope for `dry-run` mode stdout.
type DryRunCLIResult struct {
	Tool                    string                     `json:"tool"`
	Version                 string                     `json:"version"`
	Mode                    string                     `json:"mode"`
	Audit                   []audit.Record             `json:"audit,omitempty"`
	ResourceLeaseDryRun     resourcelease.DryRunResult `json:"resourceLeaseDryRun"`
	CgroupEnvelopePlan      cgroup.EnvelopePlan        `json:"cgroupEnvelopePlan"`
	MutationsForbidden      bool                       `json:"mutationsForbidden"`
	UnsafeApplyFlagPresent  bool                       `json:"unsafeApplyFlagPresent"`
	FutureApplyGateRequired bool                       `json:"futureApplyGateRequired"`
}

// RunDryRunCLI parses lease/port JSON, evaluates dry-run rules, and returns a struct safe for golden tests.
func RunDryRunCLI(leaseRaw, portRaw []byte, ctx *resourcelease.CellContext, allowUnsafeApply bool, cpuDelta, memDelta string) (*DryRunCLIResult, error) {
	lease := &crdv1alpha1.ResourceLease{}
	if err := json.Unmarshal(leaseRaw, lease); err != nil {
		return nil, fmt.Errorf("resourcelease json: %w", err)
	}
	port := &crdv1alpha1.ResourcePort{}
	if err := json.Unmarshal(portRaw, port); err != nil {
		return nil, fmt.Errorf("resourceport json: %w", err)
	}
	leaseRes := resourcelease.DryRun(lease, port, ctx)
	plan := cgroup.PlanEnvelope(allowUnsafeApply, cpuDelta, memDelta)
	var aud []audit.Record
	if safety.UnsafeApplyRequested(allowUnsafeApply) {
		aud = append(aud, audit.UnsafeApplyFlagNonOperational())
	}
	return &DryRunCLIResult{
		Tool:                    "khr-linux-agent",
		Version:                 AgentVersion,
		Mode:                    "dry-run",
		Audit:                   aud,
		ResourceLeaseDryRun:     leaseRes,
		CgroupEnvelopePlan:      plan,
		MutationsForbidden:      safety.MutationsForbidden(allowUnsafeApply),
		UnsafeApplyFlagPresent:  safety.UnsafeApplyRequested(allowUnsafeApply),
		FutureApplyGateRequired: safety.UnsafeApplyRequested(allowUnsafeApply),
	}, nil
}
