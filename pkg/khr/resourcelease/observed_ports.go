package resourcelease

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
)

// LoadObservedResourcePortsFromFile reads ResourcePort objects from scope-2 observed-json
// loop evidence (full loop output or a JSON array of ports).
func LoadObservedResourcePortsFromFile(path string) ([]crdv1alpha1.ResourcePort, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var loop struct {
		Iterations []struct {
			ResourcePorts []crdv1alpha1.ResourcePort `json:"resourcePorts"`
		} `json:"iterations"`
	}
	if err := json.Unmarshal(raw, &loop); err == nil && len(loop.Iterations) > 0 {
		for _, it := range loop.Iterations {
			if len(it.ResourcePorts) > 0 {
				return it.ResourcePorts, nil
			}
		}
	}
	var list struct {
		Items []crdv1alpha1.ResourcePort `json:"items"`
	}
	if err := json.Unmarshal(raw, &list); err == nil && len(list.Items) > 0 {
		return list.Items, nil
	}
	var ports []crdv1alpha1.ResourcePort
	if err := json.Unmarshal(raw, &ports); err == nil && len(ports) > 0 {
		return ports, nil
	}
	return nil, fmt.Errorf("no ResourcePort entries in %s", path)
}
