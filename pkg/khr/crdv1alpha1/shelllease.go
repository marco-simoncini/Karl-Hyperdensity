package crdv1alpha1

import "encoding/json"

// ShellLease is runtime.karl.io/v1alpha1 ShellLease (KHR-E).
type ShellLease struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   ObjectMeta      `json:"metadata"`
	Spec       ShellLeaseSpec  `json:"spec"`
	Status     json.RawMessage `json:"status,omitempty"`
}

// ShellLeaseSpec is the KHR-E ShellLease spec subset.
type ShellLeaseSpec struct {
	ShellRef      LocalObjectReference `json:"shellRef"`
	UserRef       string               `json:"userRef,omitempty"`
	Tenant        string               `json:"tenant,omitempty"`
	LeaseMode     string               `json:"leaseMode,omitempty"`
	AccessProfile string               `json:"accessProfile,omitempty"`
	ExpiresAt     string               `json:"expiresAt,omitempty"`
}

// GatewayRoute is gateway.karl.io/v1alpha1 GatewayRoute (KHR-E).
type GatewayRoute struct {
	APIVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Metadata   ObjectMeta         `json:"metadata"`
	Spec       GatewayRouteSpec   `json:"spec"`
	Status     json.RawMessage    `json:"status,omitempty"`
}

// GatewayRouteSpec is the KHR-E GatewayRoute spec subset.
type GatewayRouteSpec struct {
	ShellLeaseRef *LocalObjectReference `json:"shellLeaseRef,omitempty"`
	Protocol      string                `json:"protocol,omitempty"`
	TargetRef     string                `json:"targetRef,omitempty"`
	Gateway       string                `json:"gateway,omitempty"`
	PolicyRefs    []string              `json:"policyRefs,omitempty"`
	TokenRef      string                `json:"tokenRef,omitempty"`
	GatewayClass  string                `json:"gatewayClass,omitempty"`
	DisplayName   string                `json:"displayName,omitempty"`
	RDP           json.RawMessage       `json:"rdp,omitempty"`
}
