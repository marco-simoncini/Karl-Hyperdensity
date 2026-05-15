package contracts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FixtureManifest lists Hyperdensity parity fixture cases (M1–M7/M8). Test-only metadata.
type FixtureManifest struct {
	ManifestVersion    string        `json:"manifestVersion"`
	ContractKitVersion string        `json:"contractKitVersion"`
	Cases              []FixtureCase `json:"cases"`
}

// FixtureCase describes one parity anchor: Dashboard-shaped fixture and contract golden.
type FixtureCase struct {
	ID                     string `json:"id"`
	Milestone              string `json:"milestone"`
	DashboardFixture       string `json:"dashboardFixture"`
	ContractGolden         string `json:"contractGolden"`
	Edge                   string `json:"edge"`
	SupportsApply          *bool  `json:"supportsApply,omitempty"`
	WindowsEnabled         bool   `json:"windowsEnabled"`
	KubeVirtLegacyRequired bool   `json:"kubeVirtLegacyRequired"`
	ClaimSafe              bool   `json:"claimSafe"`
}

// ParseFixtureManifest unmarshals JSON into FixtureManifest.
func ParseFixtureManifest(data []byte) (FixtureManifest, error) {
	var m FixtureManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return FixtureManifest{}, fmt.Errorf("parse fixture manifest: %w", err)
	}
	return m, nil
}

// ValidateFixtureManifest checks manifest shape and M1–M7 claim-safe invariants.
func ValidateFixtureManifest(m FixtureManifest) error {
	if strings.TrimSpace(m.ManifestVersion) == "" {
		return fmt.Errorf("manifestVersion is required")
	}
	if strings.TrimSpace(m.ContractKitVersion) == "" {
		return fmt.Errorf("contractKitVersion is required")
	}
	if m.ContractKitVersion != ContractKitVersion {
		return fmt.Errorf("contractKitVersion %q incompatible with ContractKitVersion %q",
			m.ContractKitVersion, ContractKitVersion)
	}
	if len(m.Cases) == 0 {
		return fmt.Errorf("cases must not be empty")
	}
	seen := make(map[string]struct{}, len(m.Cases))
	for i, c := range m.Cases {
		if err := validateFixtureCase(c, i); err != nil {
			return err
		}
		if _, ok := seen[c.ID]; ok {
			return fmt.Errorf("duplicate case id %q", c.ID)
		}
		seen[c.ID] = struct{}{}
	}
	return nil
}

func validateFixtureCase(c FixtureCase, index int) error {
	prefix := fmt.Sprintf("cases[%d]", index)
	if strings.TrimSpace(c.ID) == "" {
		return fmt.Errorf("%s: id is required", prefix)
	}
	if strings.TrimSpace(c.Milestone) == "" {
		return fmt.Errorf("%s: milestone is required", prefix)
	}
	if strings.TrimSpace(c.DashboardFixture) == "" {
		return fmt.Errorf("%s: dashboardFixture is required", prefix)
	}
	if strings.TrimSpace(c.ContractGolden) == "" {
		return fmt.Errorf("%s: contractGolden is required", prefix)
	}
	if !c.ClaimSafe {
		return fmt.Errorf("%s: claimSafe must be true for M1–M7 anchors", prefix)
	}
	if c.WindowsEnabled {
		return fmt.Errorf("%s: windowsEnabled must be false for M1–M7 anchors", prefix)
	}
	if !c.KubeVirtLegacyRequired {
		return fmt.Errorf("%s: kubeVirtLegacyRequired must be true for M1–M7 anchors", prefix)
	}
	return nil
}
