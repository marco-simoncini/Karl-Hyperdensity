package actionapproval

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadJSON reads one ActionApproval from disk.
func LoadJSON(path string) (ActionApproval, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ActionApproval{}, err
	}
	var a ActionApproval
	if err := json.Unmarshal(data, &a); err != nil {
		return ActionApproval{}, err
	}
	if a.ActionID == "" {
		return ActionApproval{}, fmt.Errorf("actionId is required")
	}
	return a, nil
}

// SaveJSON writes one ActionApproval to disk.
func SaveJSON(path string, a ActionApproval) error {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

// LoadBundleJSON reads an approval bundle.
func LoadBundleJSON(path string) (Bundle, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Bundle{}, err
	}
	var b Bundle
	if err := json.Unmarshal(data, &b); err != nil {
		return Bundle{}, err
	}
	return b, nil
}

// SaveBundleJSON writes an approval bundle.
func SaveBundleJSON(path string, b Bundle) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
