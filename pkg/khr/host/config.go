package host

import (
	"encoding/json"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the karl-host-runtime local configuration (not a cluster CR).
type Config struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Metadata   struct {
		Name string `json:"name" yaml:"name"`
	} `json:"metadata" yaml:"metadata"`
	Spec HostSpec `json:"spec" yaml:"spec"`
}

// HostSpec holds host registration and sandbox policy.
type HostSpec struct {
	HostID              string            `json:"hostId" yaml:"hostId"`
	LinuxOnly           bool              `json:"linuxOnly" yaml:"linuxOnly"`
	SandboxMode         bool              `json:"sandboxMode" yaml:"sandboxMode"`
	SandboxApplyEnabled bool              `json:"sandboxApplyEnabled" yaml:"sandboxApplyEnabled"`
	AllowedNamespaces   []string          `json:"allowedNamespaces" yaml:"allowedNamespaces"`
	AllowedLabels       map[string]string `json:"allowedLabels" yaml:"allowedLabels"`
	CgroupRoot          string            `json:"cgroupRoot" yaml:"cgroupRoot"`
	AllowPathPrefixes        []string          `json:"allowPathPrefixes" yaml:"allowPathPrefixes"`
	ResourcePortLoopEnabled       bool  `json:"resourcePortLoopEnabled" yaml:"resourcePortLoopEnabled"`
	SandboxMaxMemoryDeltaBytes    int64 `json:"sandboxMaxMemoryDeltaBytes" yaml:"sandboxMaxMemoryDeltaBytes"`
	LaneDiscoveryEnabled              bool `json:"laneDiscoveryEnabled" yaml:"laneDiscoveryEnabled"`
	ResourceFutureSimulationEnabled   bool `json:"resourceFutureSimulationEnabled" yaml:"resourceFutureSimulationEnabled"`
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

// ValidateConfig returns validation errors.
func ValidateConfig(cfg *Config) []string {
	var errs []string
	if cfg == nil {
		return []string{"config is nil"}
	}
	if cfg.Spec.HostID == "" {
		errs = append(errs, "spec.hostId is required")
	}
	if !cfg.Spec.LinuxOnly {
		errs = append(errs, "spec.linuxOnly must be true")
	}
	if !cfg.Spec.SandboxMode {
		errs = append(errs, "spec.sandboxMode must be true for karl-host-runtime MVP")
	}
	return errs
}

// ValidateConfigForLaneDiscovery validates config for read-only multi-lane discovery (KHR-Q).
func ValidateConfigForLaneDiscovery(cfg *Config) []string {
	var errs []string
	if cfg == nil {
		return []string{"config is nil"}
	}
	if cfg.Spec.HostID == "" {
		errs = append(errs, "spec.hostId is required")
	}
	if !cfg.Spec.SandboxMode {
		errs = append(errs, "spec.sandboxMode must be true")
	}
	if !cfg.Spec.LaneDiscoveryEnabled {
		errs = append(errs, "spec.laneDiscoveryEnabled must be true for lane-discovery mode")
	}
	return errs
}

// ValidateConfigForResourceFutureSimulation validates config for KHR-R simulation mode.
func ValidateConfigForResourceFutureSimulation(cfg *Config) []string {
	var errs []string
	if cfg == nil {
		return []string{"config is nil"}
	}
	if cfg.Spec.HostID == "" {
		errs = append(errs, "spec.hostId is required")
	}
	if !cfg.Spec.SandboxMode {
		errs = append(errs, "spec.sandboxMode must be true")
	}
	if !cfg.Spec.LaneDiscoveryEnabled {
		errs = append(errs, "spec.laneDiscoveryEnabled must be true (simulation uses lane discovery input)")
	}
	if !cfg.Spec.ResourceFutureSimulationEnabled {
		errs = append(errs, "spec.resourceFutureSimulationEnabled must be true")
	}
	return errs
}
