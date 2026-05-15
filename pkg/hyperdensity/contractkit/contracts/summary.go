// Package contracts holds minimal Hyperdensity DTOs extracted from Dashboard
// parent-fabric ?view=summary (M1). No HTTP, no Kubernetes clients.
package contracts

import (
	"encoding/json"
	"fmt"
	"strings"
)

const SummaryAPIVersion = "hyperdensity.karl.io/parent-fabric-summary/v1"

// ParentFabricSummary is a permissive subset of Dashboard parent-fabric summary JSON.
type ParentFabricSummary struct {
	APIVersion      string                 `json:"apiVersion"`
	GeneratedAt     string                 `json:"generatedAt"`
	Source          string                 `json:"source"`
	ParentPool      ParentPoolSummary      `json:"parentPool"`
	ExecutionEngine ExecutionEngineSummary `json:"executionEngine"`
	WindowsLane     WindowsLaneSummary     `json:"windowsLane"`
	KubeVirtLegacy  KubeVirtLegacySummary  `json:"kubeVirtLegacy"`
	Hyperdensity    HyperdensityPosture    `json:"hyperdensity"`
}

type ParentPoolSummary struct {
	UsageSummary  string `json:"usageSummary,omitempty"`
	DonorCount    int    `json:"donorCount,omitempty"`
	ReceiverCount int    `json:"receiverCount,omitempty"`
}

type ExecutionEngineSummary struct {
	Mode            string `json:"mode,omitempty"`
	AutonomousMode  bool   `json:"autonomousMode,omitempty"`
	ApplyAllowed    bool   `json:"applyAllowed,omitempty"`
	DryRunSupported bool   `json:"dryRunSupported,omitempty"`
}

type WindowsLaneSummary struct {
	Enabled  bool     `json:"enabled"`
	Reason   string   `json:"reason,omitempty"`
	Blockers []string `json:"blockers,omitempty"`
}

type KubeVirtLegacySummary struct {
	Present      bool   `json:"present"`
	ProviderMode string `json:"providerMode,omitempty"`
}

type HyperdensityPosture struct {
	RecommendationOnly bool `json:"recommendationOnly"`
	OperatorControlled bool `json:"operatorControlled"`
}

// ParseParentFabricSummary unmarshals JSON into ParentFabricSummary.
func ParseParentFabricSummary(data []byte) (ParentFabricSummary, error) {
	var s ParentFabricSummary
	if err := json.Unmarshal(data, &s); err != nil {
		return ParentFabricSummary{}, fmt.Errorf("parse parent fabric summary: %w", err)
	}
	return s, nil
}

// ValidateSummary checks required M1 fields and basic consistency.
func ValidateSummary(s ParentFabricSummary) error {
	if strings.TrimSpace(s.APIVersion) == "" {
		return fmt.Errorf("apiVersion is required")
	}
	if strings.TrimSpace(s.GeneratedAt) == "" {
		return fmt.Errorf("generatedAt is required")
	}
	if strings.TrimSpace(s.Source) == "" {
		return fmt.Errorf("source is required")
	}
	if !s.ExecutionEngine.DryRunSupported {
		return fmt.Errorf("executionEngine.dryRunSupported must be true for M1 anchor")
	}
	if !s.KubeVirtLegacy.Present {
		return fmt.Errorf("kubeVirtLegacy.present must be true for M1 anchor")
	}
	if !s.Hyperdensity.RecommendationOnly {
		return fmt.Errorf("hyperdensity.recommendationOnly must be true for M1 anchor")
	}
	return nil
}

// ValidateNoForbiddenClaims rejects postures that imply production apply, Windows enablement,
// or autonomous broad execution on the M1 anchor surface.
func ValidateNoForbiddenClaims(s ParentFabricSummary) error {
	if s.WindowsLane.Enabled {
		return fmt.Errorf("windows lane must not be enabled")
	}
	if s.ExecutionEngine.ApplyAllowed && !s.Hyperdensity.OperatorControlled {
		return fmt.Errorf("applyAllowed requires operatorControlled posture")
	}
	if s.ExecutionEngine.AutonomousMode {
		return fmt.Errorf("autonomousMode must be false")
	}
	if strings.EqualFold(s.ExecutionEngine.Mode, "autonomous") {
		return fmt.Errorf("executionEngine.mode must not be autonomous")
	}
	for _, id := range s.WindowsLane.Blockers {
		if id == "windows_lane_enabled" || id == "production_ready_windows" {
			return fmt.Errorf("forbidden windows blocker claim: %s", id)
		}
	}
	return nil
}
