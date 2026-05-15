package crdv1alpha1

// ResourcePort is runtime.karl.io/v1alpha1 ResourcePort.
type ResourcePort struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   ObjectMeta       `json:"metadata"`
	Spec       ResourcePortSpec `json:"spec"`
}

// ResourcePortSpec carries capability truth matrix for cpu/memory and optional resources.
type ResourcePortSpec struct {
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
