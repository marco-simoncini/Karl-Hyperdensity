package crdv1alpha1

import "encoding/json"

// Cell is runtime.karl.io/v1alpha1 Cell.
type Cell struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   ObjectMeta      `json:"metadata"`
	Spec       CellSpec        `json:"spec"`
	Status     json.RawMessage `json:"status,omitempty"`
}

// LocalObjectReference is a name-only ref within namespace.
type LocalObjectReference struct {
	Name string `json:"name"`
}

// CellSpec is the KHR-D Cell spec subset.
type CellSpec struct {
	ShellRef               LocalObjectReference  `json:"shellRef"`
	HostRef                *LocalObjectReference `json:"hostRef,omitempty"`
	RuntimeProviderRef     LocalObjectReference  `json:"runtimeProviderRef"`
	ProviderBinding        string                `json:"providerBinding,omitempty"`
	RuntimeClass           string                `json:"runtimeClass,omitempty"`
	Owner                  *KHROwner             `json:"owner,omitempty"`
	Resources              *KHRResources         `json:"resources,omitempty"`
	NetworkRefs            []KHRNamedRef         `json:"networkRefs,omitempty"`
	StorageRefs            []KHRStorageRef       `json:"storageRefs,omitempty"`
	ResourcePortProfileRef *LocalObjectReference `json:"resourcePortProfileRef,omitempty"`
	ResourcePorts          json.RawMessage       `json:"resourcePorts,omitempty"`
	ProviderHandle         json.RawMessage       `json:"providerHandle,omitempty"`
	PriorityClass          string                `json:"priorityClass,omitempty"`
}
