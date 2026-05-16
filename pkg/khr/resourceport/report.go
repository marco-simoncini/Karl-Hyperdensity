package resourceport

import (
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

// Candidate is a ResourcePort-shaped observation (not applied to cluster).
type Candidate struct {
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   crdv1alpha1.ObjectMeta `json:"metadata"`
	Spec       CandidateSpec          `json:"spec"`
	Status     CandidateStatus        `json:"status,omitempty"`
}

// CandidateSpec aligns with KHR-C ResourcePort observation fields.
type CandidateSpec struct {
	Provider     string                  `json:"provider"`
	ShellRef     string                  `json:"shellRef"`
	CellRef      string                  `json:"cellRef"`
	Capabilities []string                `json:"capabilities"`
	Hotplug      crdv1alpha1.ResourcePortHotplug `json:"hotplug"`
	Ports        crdv1alpha1.ResourcePortsMatrix `json:"ports"`
}

// CandidateStatus is minimal observation status.
type CandidateStatus struct {
	ObservedAt string `json:"observedAt"`
	Phase      string `json:"phase"`
}

// ReportCandidate emits a ResourcePort candidate from host capabilities.
func ReportCandidate(cfg *host.Config, shellRef, cellRef, namespace, name string) Candidate {
	caps := host.ReportCapabilities(cfg)
	now := time.Now().UTC().Format(time.RFC3339)
	return Candidate{
		APIVersion: "runtime.karl.io/v1alpha1",
		Kind:       "ResourcePort",
		Metadata: crdv1alpha1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"khr.karl.io/sandbox": "true",
				"karl.io/host-id":     cfg.Spec.HostID,
			},
		},
		Spec: CandidateSpec{
			Provider:     "khr.native",
			ShellRef:     shellRef,
			CellRef:      cellRef,
			Capabilities: []string{"cpu.envelope", "memory.envelope"},
			Hotplug: crdv1alpha1.ResourcePortHotplug{
				CPU: false, Memory: false, Disk: false, Network: false,
			},
			Ports: crdv1alpha1.ResourcePortsMatrix{
				CPU:    crdv1alpha1.ResourceModes{Modes: caps.SupportedModes},
				Memory: crdv1alpha1.ResourceModes{Modes: caps.SupportedModes},
			},
		},
		Status: CandidateStatus{
			ObservedAt: now,
			Phase:      "Observed",
		},
	}
}
