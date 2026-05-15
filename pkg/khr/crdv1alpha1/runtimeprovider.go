package crdv1alpha1

import "encoding/json"

// RuntimeProvider is runtime.karl.io/v1alpha1 RuntimeProvider.
type RuntimeProvider struct {
	APIVersion string              `json:"apiVersion"`
	Kind       string              `json:"kind"`
	Metadata   ObjectMeta          `json:"metadata"`
	Spec       RuntimeProviderSpec `json:"spec"`
	Status     json.RawMessage     `json:"status,omitempty"`
}

// RuntimeProviderSpec is the spec subset for linux cgroup / systemd drivers.
type RuntimeProviderSpec struct {
	ID                            string                `json:"id"`
	DisplayName                   string                `json:"displayName,omitempty"`
	Driver                        string                `json:"driver"`
	SupportedShellClassIDs        []string              `json:"supportedShellClassIds,omitempty"`
	DefaultResourcePortProfileRef *LocalObjectReference `json:"defaultResourcePortProfileRef,omitempty"`
	ExecContractVersion           string                `json:"execContractVersion,omitempty"`
	Config                        json.RawMessage       `json:"config,omitempty"`
}
