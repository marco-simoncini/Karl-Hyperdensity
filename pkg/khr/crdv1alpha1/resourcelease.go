package crdv1alpha1

import "encoding/json"

// ResourceLease is a contract-aligned document for hyperdensity.karl.io/v1alpha1 ResourceLease (ADR-0005 unified).
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

// ShellBinding is the Shell section of a unified ResourceLease.
type ShellBinding struct {
	Kind       string `json:"kind,omitempty"`
	Ref        string `json:"ref,omitempty"`
	Experience string `json:"experience,omitempty"`
}

// CellBinding is the Cell section of a unified ResourceLease.
type CellBinding struct {
	Ref          string          `json:"ref,omitempty"`
	HostSelector json.RawMessage `json:"hostSelector,omitempty"`
	RuntimeClass string          `json:"runtimeClass,omitempty"`
}

// TransferLeaseSpec is the transfer block (leaseKind=transfer).
type TransferLeaseSpec struct {
	Donor    LeaseRef        `json:"donor"`
	Receiver LeaseRef        `json:"receiver"`
	Resource string          `json:"resource"`
	Mode     string          `json:"mode"`
	Amount   json.RawMessage `json:"amount,omitempty"`
}

// LeaseGovernance holds dry-run, rollback, and telemetry gates.
type LeaseGovernance struct {
	DryRunOnly                 *bool            `json:"dryRunOnly,omitempty"`
	RollbackPlanRef            *RollbackPlanRef `json:"rollbackPlanRef,omitempty"`
	RollbackRequired           *bool            `json:"rollbackRequired,omitempty"`
	VerificationHooks          json.RawMessage  `json:"verificationHooks,omitempty"`
	NoRestart                  *bool            `json:"noRestart,omitempty"`
	GuestVisible               *bool            `json:"guestVisible,omitempty"`
	TelemetryConvergedRequired *bool            `json:"telemetryConvergedRequired,omitempty"`
	DurationSeconds            *int64           `json:"durationSeconds,omitempty"`
	TTLSeconds                 *int64           `json:"ttlSeconds,omitempty"`
}

// LeaseSpec is the unified ResourceLease spec (runtime + transfer).
type LeaseSpec struct {
	LeaseKind string          `json:"leaseKind,omitempty"`
	Shell     ShellBinding    `json:"shell,omitempty"`
	Cell      CellBinding     `json:"cell,omitempty"`
	Provider  string          `json:"provider,omitempty"`
	Transfer  *TransferLeaseSpec `json:"transfer,omitempty"`
	Governance *LeaseGovernance `json:"governance,omitempty"`

	// Deprecated inline transfer fields (v1alpha1 fixtures); prefer spec.transfer.
	Donor                      LeaseRef         `json:"donor,omitempty"`
	Receiver                   LeaseRef         `json:"receiver,omitempty"`
	Resource                   string           `json:"resource,omitempty"`
	Mode                       string           `json:"mode,omitempty"`
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

// EffectiveTransfer returns donor/receiver/resource/mode from unified or legacy spec.
func (s LeaseSpec) EffectiveTransfer() (donor, receiver LeaseRef, resource, mode string, ok bool) {
	if s.Transfer != nil {
		donor, receiver = s.Transfer.Donor, s.Transfer.Receiver
		resource, mode = s.Transfer.Resource, s.Transfer.Mode
		return donor, receiver, resource, mode, donor.Name != "" && receiver.Name != ""
	}
	if s.Donor.Name != "" && s.Receiver.Name != "" {
		return s.Donor, s.Receiver, s.Resource, s.Mode, true
	}
	return LeaseRef{}, LeaseRef{}, "", "", false
}

// EffectiveDryRunOnly returns governance or legacy dryRunOnly flag.
func (s LeaseSpec) EffectiveDryRunOnly() bool {
	if s.Governance != nil && s.Governance.DryRunOnly != nil {
		return *s.Governance.DryRunOnly
	}
	if s.DryRunOnly != nil {
		return *s.DryRunOnly
	}
	return false
}
