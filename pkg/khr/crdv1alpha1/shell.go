package crdv1alpha1

import "encoding/json"

// Shell is runtime.karl.io/v1alpha1 Shell (KHR-D foundation).
type Shell struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   ObjectMeta      `json:"metadata"`
	Spec       ShellSpec         `json:"spec"`
	Status     json.RawMessage `json:"status,omitempty"`
}

// ShellSpec is the KHR-D Shell spec subset.
type ShellSpec struct {
	ShellClassRef    LocalObjectReference `json:"shellClassRef"`
	ProviderBinding  string             `json:"providerBinding,omitempty"`
	RuntimeClass     string             `json:"runtimeClass,omitempty"`
	Owner            *KHROwner          `json:"owner,omitempty"`
	Resources        *KHRResources      `json:"resources,omitempty"`
	NetworkRefs      []KHRNamedRef      `json:"networkRefs,omitempty"`
	StorageRefs      []KHRStorageRef    `json:"storageRefs,omitempty"`
	DisplayName      string             `json:"displayName,omitempty"`
}

// ShellClass is runtime.karl.io/v1alpha1 ShellClass.
type ShellClass struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   ObjectMeta      `json:"metadata"`
	Spec       ShellClassSpec  `json:"spec"`
	Status     json.RawMessage `json:"status,omitempty"`
}

// ShellClassSpec is the KHR-D ShellClass spec subset.
type ShellClassSpec struct {
	ID                       string          `json:"id"`
	ComplianceFamily         string          `json:"complianceFamily"`
	DefaultProviderBinding   string          `json:"defaultProviderBinding,omitempty"`
	DefaultRuntimeClass      string          `json:"defaultRuntimeClass,omitempty"`
	DefaultResources         *KHRResources   `json:"defaultResources,omitempty"`
	DefaultNetworkRefs       []KHRNamedRef   `json:"defaultNetworkRefs,omitempty"`
	DefaultStorageRefs       []KHRStorageRef `json:"defaultStorageRefs,omitempty"`
	SupportedRuntimeProviders []string       `json:"supportedRuntimeProviders,omitempty"`
	DisplayName              string          `json:"displayName,omitempty"`
}

// KHROwner is user/tenant ownership.
type KHROwner struct {
	User   string `json:"user,omitempty"`
	Tenant string `json:"tenant,omitempty"`
}

// KHRResources is cpu/memory request/limit strings.
type KHRResources struct {
	CPU    *KHRResourceQuantity `json:"cpu,omitempty"`
	Memory *KHRResourceQuantity `json:"memory,omitempty"`
}

// KHRResourceQuantity holds request/limit.
type KHRResourceQuantity struct {
	Request string `json:"request,omitempty"`
	Limit   string `json:"limit,omitempty"`
}

// KHRNamedRef is a named network (or generic) ref.
type KHRNamedRef struct {
	Name string `json:"name"`
	Ref  string `json:"ref"`
}

// KHRStorageRef is a named storage ref with optional mode.
type KHRStorageRef struct {
	Name string `json:"name"`
	Ref  string `json:"ref"`
	Mode string `json:"mode,omitempty"`
}
