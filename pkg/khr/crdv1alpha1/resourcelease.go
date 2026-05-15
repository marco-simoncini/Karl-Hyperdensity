package crdv1alpha1

import "encoding/json"

// ResourceLease is a contract-aligned document for hyperdensity.karl.io/v1alpha1 ResourceLease.
type ResourceLease struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   ObjectMeta      `json:"metadata"`
	Spec       LeaseSpec       `json:"spec"`
	Status     json.RawMessage `json:"status,omitempty"`
}

// LeaseRef identifies donor or receiver (Shell or Cell).
type LeaseRef struct {
	APIGroup  string `json:"apiGroup,omitempty"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`
}

// LeaseSpec is the spec subset enforced by KHR Linux envelope dry-run.
type LeaseSpec struct {
	Donor                      LeaseRef         `json:"donor"`
	Receiver                   LeaseRef         `json:"receiver"`
	Resource                   string           `json:"resource"`
	Mode                       string           `json:"mode"`
	Amount                     json.RawMessage  `json:"amount,omitempty"`
	DurationSeconds            *int64           `json:"durationSeconds,omitempty"`
	TTLSeconds                 *int64           `json:"ttlSeconds,omitempty"`
	RollbackPlanRef            *RollbackPlanRef `json:"rollbackPlanRef,omitempty"`
	RollbackRequired           *bool            `json:"rollbackRequired,omitempty"`
	VerificationHooks          json.RawMessage  `json:"verificationHooks,omitempty"`
	NoRestart                  *bool            `json:"noRestart,omitempty"`
	GuestVisible               *bool            `json:"guestVisible,omitempty"`
	TelemetryConvergedRequired *bool            `json:"telemetryConvergedRequired,omitempty"`
	DryRunOnly                 *bool            `json:"dryRunOnly,omitempty"`
}

// RollbackPlanRef points at a rollback plan resource.
type RollbackPlanRef struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}
