package crdv1alpha1

import "encoding/json"

// ResourcePort is runtime.karl.io/v1alpha1 ResourcePort.
type ResourcePort struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   ObjectMeta       `json:"metadata"`
	Spec       ResourcePortSpec `json:"spec"`
	Status     json.RawMessage  `json:"status,omitempty"`
}

// ResourcePortHotplug declares hotplug posture per resource class (KHR-C observation).
type ResourcePortHotplug struct {
	CPU     bool `json:"cpu"`
	Memory  bool `json:"memory"`
	Disk    bool `json:"disk"`
	Network bool `json:"network"`
}

// ResourcePortSpec carries capability truth matrix and optional observation binding.
type ResourcePortSpec struct {
	Provider                    string              `json:"provider,omitempty"`
	ShellRef                    string              `json:"shellRef,omitempty"`
	CellRef                     string              `json:"cellRef,omitempty"`
	Capabilities                []string            `json:"capabilities,omitempty"`
	Hotplug                     *ResourcePortHotplug `json:"hotplug,omitempty"`
	Constraints                 json.RawMessage     `json:"constraints,omitempty"`
	Description                 string              `json:"description,omitempty"`
	AppliesToRuntimeProviderIDs []string            `json:"appliesToRuntimeProviderIds,omitempty"`
	Ports                       ResourcePortsMatrix `json:"ports"`
	Notes                       string              `json:"notes,omitempty"`
}

// ResourcePortsMatrix lists supported modes per resource class.
type ResourcePortsMatrix struct {
	CPU     ResourceModes  `json:"cpu"`
	Memory  ResourceModes  `json:"memory"`
	Disk    *ResourceModes `json:"disk,omitempty"`
	Network *ResourceModes `json:"network,omitempty"`
	GPU     *ResourceModes `json:"gpu,omitempty"`
}

// ResourceModes is the modes array for one resource class.
type ResourceModes struct {
	Modes []string `json:"modes"`
}
