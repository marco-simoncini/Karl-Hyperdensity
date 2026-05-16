package windowslane

import "fmt"

const (
	ProviderWindowsHostRuntime   = "windows.host-runtime"
	ProviderKubevirtCompatibility = "kubevirt.compatibility"

	LiveScaleTargetInPlace    = "live-in-place"
	LiveScaleTargetCompat     = "compatibility-fallback"

	BlockedRequiresRestart      = "requiresRestart"
	BlockedRequiresReboot         = "requiresReboot"
	BlockedRequiresSessionDrain   = "requiresSessionDrain"
	BlockedProviderUnsupported    = "providerUnsupported"
)

// Capabilities is the KHR-P Windows ResourcePort live scale matrix.
type Capabilities struct {
	CPULiveScaleSupported bool   `json:"cpuLiveScaleSupported"`
	RAMLiveScaleSupported bool   `json:"ramLiveScaleSupported"`
	ScaleUpSupported      bool   `json:"scaleUpSupported"`
	ScaleDownSupported    bool   `json:"scaleDownSupported"`
	RequiresRestart       bool   `json:"requiresRestart"`
	LiveScaleTarget       string `json:"liveScaleTarget,omitempty"`
	ObservationOnly       bool   `json:"observationOnly,omitempty"`
}

// DefaultHostRuntimeCapabilities is the target windows.host-runtime posture (KHR-P).
func DefaultHostRuntimeCapabilities() Capabilities {
	return Capabilities{
		CPULiveScaleSupported: true,
		RAMLiveScaleSupported: true,
		ScaleUpSupported:      true,
		ScaleDownSupported:    true,
		RequiresRestart:       false,
		LiveScaleTarget:       LiveScaleTargetInPlace,
		ObservationOnly:       true,
	}
}

// CapabilitiesForProvider returns observation capabilities per provider binding.
func CapabilitiesForProvider(providerBinding string) (Capabilities, error) {
	switch providerBinding {
	case ProviderWindowsHostRuntime:
		return DefaultHostRuntimeCapabilities(), nil
	case ProviderKubevirtCompatibility:
		return Capabilities{
			CPULiveScaleSupported: false,
			RAMLiveScaleSupported: false,
			ScaleUpSupported:      false,
			ScaleDownSupported:    false,
			RequiresRestart:         true,
			LiveScaleTarget:         LiveScaleTargetCompat,
			ObservationOnly:         true,
		}, nil
	default:
		return Capabilities{}, fmt.Errorf("unsupported providerBinding %q", providerBinding)
	}
}

// ValidateCapabilities enforces KHR-P safety (requiresRestart target false on host-runtime).
func ValidateCapabilities(c Capabilities, providerBinding string) error {
	if providerBinding == ProviderWindowsHostRuntime && c.RequiresRestart {
		return fmt.Errorf("windows.host-runtime target requires requiresRestart=false")
	}
	if c.LiveScaleTarget == "" {
		return fmt.Errorf("liveScaleTarget required")
	}
	if c.ObservationOnly != true {
		return fmt.Errorf("KHR-P windows lane must be observationOnly")
	}
	return nil
}

// BlockedObservation is a dry-run blocked outcome for Windows lane.
type BlockedObservation struct {
	Blocked       bool   `json:"blocked"`
	BlockedState  string `json:"blockedState"`
	BlockedReason string `json:"blockedReason,omitempty"`
	Reason        string `json:"reason,omitempty"`
	NoMutation    bool   `json:"noMutation"`
	NoApply       bool   `json:"noApply"`
	LiveScaleTarget string `json:"liveScaleTarget,omitempty"`
}

// BlockMemoryOnCompatibility returns the standard kubevirt.compatibility block for RAM scaleUp.
func BlockMemoryOnCompatibility() BlockedObservation {
	msg := "kubevirt.compatibility: memory scaleUp requires VM restart; not live-in-place"
	return BlockedObservation{
		Blocked:         true,
		BlockedState:    BlockedRequiresRestart,
		BlockedReason:   msg,
		Reason:          msg,
		NoMutation:      true,
		NoApply:         true,
		LiveScaleTarget: LiveScaleTargetCompat,
	}
}
