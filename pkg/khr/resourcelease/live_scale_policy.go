package resourcelease

import (
	"fmt"
	"strings"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

const (
	ModeScaleUp   = "scaleUp"
	ModeScaleDown = "scaleDown"
	ModeEnvelope  = "envelope"

	AnnotationRestartRequired = "khr.karl.io/restart-required"
	AnnotationRolloutRequired = "khr.karl.io/rollout-required"
	AnnotationRecreateRequired = "khr.karl.io/recreate-required"
)

// LiveScalePolicy captures no-restart / no-rollout sandbox rules.
type LiveScalePolicy struct {
	NoRestart      bool `json:"noRestart"`
	NoRollout      bool `json:"noRollout"`
	NoRecreate     bool `json:"noRecreate"`
	NoProduction   bool `json:"noProductionMutation"`
	LiveInPlace    bool `json:"liveInPlace"`
}

// DefaultLiveScalePolicy is the KHR-O sandbox default.
func DefaultLiveScalePolicy() LiveScalePolicy {
	return LiveScalePolicy{
		NoRestart:    true,
		NoRollout:    true,
		NoRecreate:   true,
		NoProduction: true,
		LiveInPlace:  true,
	}
}

// ValidateLiveScaleLease blocks restart/rollout/recreate semantics on leases.
func ValidateLiveScaleLease(lease *crdv1alpha1.ResourceLease) error {
	if lease == nil {
		return fmt.Errorf("lease is nil")
	}
	ann := lease.Metadata.Annotations
	checks := []struct {
		key, label string
	}{
		{AnnotationRestartRequired, "restart"},
		{AnnotationRolloutRequired, "rollout"},
		{AnnotationRecreateRequired, "recreate"},
	}
	for _, c := range checks {
		if strings.EqualFold(ann[c.key], "true") {
			return fmt.Errorf("%s required: blocked — KHR-O live-in-place only (no %s)", c.label, c.label)
		}
	}
	return nil
}

// ValidateTransferMode checks resource/mode pairing for dry-run and apply.
func ValidateTransferMode(resource, mode string) error {
	modeNorm := strings.ToLower(strings.TrimSpace(mode))
	switch resource {
	case "cpu":
		if modeNorm != strings.ToLower(ModeEnvelope) {
			return fmt.Errorf("cpu requires mode envelope, got %q", mode)
		}
	case "memory":
		if modeNorm != strings.ToLower(ModeScaleUp) &&
			modeNorm != strings.ToLower(ModeScaleDown) &&
			modeNorm != strings.ToLower(ModeEnvelope) {
			return fmt.Errorf("memory requires mode scaleUp, scaleDown, or envelope, got %q", mode)
		}
	default:
		return fmt.Errorf("unsupported resource %q", resource)
	}
	return nil
}

// SandboxMaxMemoryDelta returns configured max memory delta bytes.
func SandboxMaxMemoryDelta(cfg *host.Config) int64 {
	if cfg != nil && cfg.Spec.SandboxMaxMemoryDeltaBytes > 0 {
		return cfg.Spec.SandboxMaxMemoryDeltaBytes
	}
	return SandboxMaxMemoryBytes
}
