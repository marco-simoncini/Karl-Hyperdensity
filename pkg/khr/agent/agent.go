package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/audit"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/discovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence/integrity"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcelease"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/runtimeprovider"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/safety"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/telemetry"
)

// AgentVersion is embedded in JSON outputs (Sprint 10).
const AgentVersion = "0.0.1-sprint10"

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

// RunDiscoverCgroupsCLI performs read-only cgroup discovery (no mutations).
func RunDiscoverCgroupsCLI(cfg *Config, cell *crdv1alpha1.Cell, scannedRoot, allowPathPrefix string) *discovery.CgroupDiscoveryOutput {
	out := discovery.Run(cfg.Spec.AgentID, scannedRoot, allowPathPrefix, cell)
	out.Tool = "khr-linux-agent"
	out.Version = AgentVersion
	out.Mode = "discover-cgroups"
	return out
}

// RunReadTelemetryCLI collects read-only cgroup v2 metrics under a resolved path.
func RunReadTelemetryCLI(cfg *Config, cgroupPath, allowPathPrefix string, cell *crdv1alpha1.Cell, telemetryOutPath string) (*telemetry.ReadTelemetryOutput, error) {
	out := &telemetry.ReadTelemetryOutput{
		Tool:               "khr-linux-agent",
		Version:            AgentVersion,
		Mode:               "read-telemetry",
		AgentID:            cfg.Spec.AgentID,
		TelemetryMode:      "read-only",
		CgroupPath:         cgroupPath,
		AllowedPathPrefix:  strings.TrimSpace(allowPathPrefix),
		MutationsForbidden: true,
	}
	if cell != nil {
		out.CellRef = &telemetry.CellRef{
			APIVersion: cell.APIVersion,
			Kind:       cell.Kind,
			Namespace:  cell.Metadata.Namespace,
			Name:       cell.Metadata.Name,
		}
	}
	resolved, w, b := cgroup.ValidateCgroupPathForTelemetry(cgroupPath, allowPathPrefix)
	if len(b) > 0 {
		out.Metrics = telemetry.MetricsBundle{}
		out.Evidence = telemetry.BuildEvidence(w, b, out.Metrics)
		if err := writeTelemetryOutputOptional(out, telemetryOutPath); err != nil {
			return nil, err
		}
		return out, nil
	}
	out.CgroupPath = resolved
	m, w2, b2 := telemetry.ReadCgroupV2Metrics(resolved)
	allW := append(append([]string{}, w...), w2...)
	allB := append(append([]string{}, b...), b2...)
	telemetry.NormalizeMetricsEmpty(&m)
	out.Metrics = m
	out.Evidence = telemetry.BuildEvidence(allW, allB, m)
	if err := writeTelemetryOutputOptional(out, telemetryOutPath); err != nil {
		return nil, err
	}
	return out, nil
}

func writeTelemetryOutputOptional(out *telemetry.ReadTelemetryOutput, path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o600)
}

func writeCollectEvidenceOptional(bundle *evidence.CollectEvidenceBundle, path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	b, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o600)
}

// CollectEvidenceIntegrityOpts configures optional local manifest/digest sidecars (Sprint 10).
type CollectEvidenceIntegrityOpts struct {
	ManifestOutputPath string
	DigestOutputPath   string
	SigningMode        string
	SigningKeyFile     string
	ArtifactID         string
}

func shouldEmitEvidenceSidecars(o *CollectEvidenceIntegrityOpts) bool {
	if o == nil {
		return false
	}
	if strings.TrimSpace(o.ManifestOutputPath) != "" || strings.TrimSpace(o.DigestOutputPath) != "" {
		return true
	}
	return false
}

// RunCollectEvidenceCLI validates config, runs discovery, telemetry (when selectedPath is set),
// optionally ResourceLease dry-run when both lease and port payloads are provided, and returns
// a single local evidence bundle JSON shape (read-only; mutationsForbidden always true).
func RunCollectEvidenceCLI(cfg *Config, cell *crdv1alpha1.Cell, cgroupRoot, allowPrefix string, leaseRaw, portRaw []byte, cellCtx *resourcelease.CellContext, evidenceOutPath string, allowUnsafe bool, cpuDelta, memDelta string, integ *CollectEvidenceIntegrityOpts) (*evidence.CollectEvidenceBundle, error) {
	if errs := ValidateConfig(cfg); len(errs) > 0 {
		return nil, fmt.Errorf("invalid config: %s", strings.Join(errs, "; "))
	}
	disc := RunDiscoverCgroupsCLI(cfg, cell, cgroupRoot, allowPrefix)

	var cellRef *telemetry.CellRef
	if cell != nil {
		cellRef = &telemetry.CellRef{
			APIVersion: cell.APIVersion,
			Kind:       cell.Kind,
			Namespace:  cell.Metadata.Namespace,
			Name:       cell.Metadata.Name,
		}
	}

	var tSnap evidence.TelemetrySnapshot
	if strings.TrimSpace(disc.SelectedPath) != "" {
		tel, err := RunReadTelemetryCLI(cfg, disc.SelectedPath, allowPrefix, cell, "")
		if err != nil {
			return nil, err
		}
		tSnap = evidence.TelemetrySnapshotFrom(tel)
	} else {
		tSnap = evidence.TelemetrySnapshotSkipped(allowPrefix, cellRef, "discovery did not resolve selectedPath; telemetry skipped")
	}

	var partialWarn string
	var dry evidence.DryRunPayload
	leaseHas := len(bytes.TrimSpace(leaseRaw)) > 0
	portHas := len(bytes.TrimSpace(portRaw)) > 0
	switch {
	case leaseHas && portHas:
		dr, err := RunDryRunCLI(leaseRaw, portRaw, cellCtx, allowUnsafe, cpuDelta, memDelta)
		if err != nil {
			return nil, err
		}
		dry = evidence.DryRunPayloadFromResult(dr.ResourceLeaseDryRun, dr.CgroupEnvelopePlan, dr.MutationsForbidden, dr.UnsafeApplyFlagPresent, dr.FutureApplyGateRequired, dr.Audit)
	case leaseHas || portHas:
		partialWarn = "collect-evidence: dry-run skipped: both -lease-input and -resource-port-input are required when including ResourceLease simulation"
		dry = evidence.DryRunSkippedPayload(partialWarn)
	default:
		dry = evidence.DryRunSkippedPayload("no lease or resource port inputs provided for optional dry-run")
	}

	bundle := evidence.BuildCollectEvidenceBundle(AgentVersion, cfg.Spec.AgentID, cell, disc, tSnap, dry, partialWarn)
	if err := writeCollectEvidenceOptional(bundle, evidenceOutPath); err != nil {
		return nil, err
	}
	if integ != nil {
		if err := integrity.ValidateSigningMode(integ.SigningMode); err != nil {
			return nil, err
		}
		if err := integrity.RequireLocalDevKey(integ.SigningMode, integ.SigningKeyFile); err != nil {
			return nil, err
		}
		if integrity.NormalizeSigningMode(integ.SigningMode) == "local-dev" && strings.TrimSpace(integ.ManifestOutputPath) == "" {
			return nil, fmt.Errorf("collect-evidence: signing-mode=local-dev requires -evidence-manifest-output")
		}
		if shouldEmitEvidenceSidecars(integ) {
			if err := integrity.EmitEvidenceSidecars(
				bundle,
				bundle.AgentID,
				strings.TrimSpace(integ.ArtifactID),
				integ.ManifestOutputPath,
				integ.DigestOutputPath,
				integ.SigningMode,
				integ.SigningKeyFile,
			); err != nil {
				return nil, err
			}
		}
	}
	return bundle, nil
}
