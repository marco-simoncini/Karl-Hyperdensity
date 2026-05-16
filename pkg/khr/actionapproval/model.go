// Package actionapproval provides read-only operator approval workflow (KHR-W).
package actionapproval

import (
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/policygates"
)

const (
	WorkflowID = "khr-operator-action-approval-v1"

	StatePending  = "pending"
	StateApproved = "approved"
	StateRejected = "rejected"
	StateExpired  = "expired"

	DefaultTTLSeconds int64 = 3600
)

// ActionApproval is a local evidence record for operator consent (no apply).
type ActionApproval struct {
	ActionID           string                    `json:"actionId"`
	ResourceFutureRef  string                    `json:"resourceFutureRef"`
	ResourceLeaseRef   string                    `json:"resourceLeaseRef"`
	LaneID             string                    `json:"laneId"`
	CertificationRef   string                    `json:"certificationRef"`
	PolicyGateResult   policygates.EligibilityOutcome `json:"policyGateResult"`
	ApprovalState      string                    `json:"approvalState"`
	ApprovedBy         string                    `json:"approvedBy,omitempty"`
	ApprovedAt         string                    `json:"approvedAt,omitempty"`
	ExpiresAt          string                    `json:"expiresAt"`
	Reason             string                    `json:"reason,omitempty"`
	ReadOnly           bool                      `json:"readOnly"`
	NoApply            bool                      `json:"noApply"`
	NoMutation         bool                      `json:"noMutation"`
	NoAutonomousOrchestration bool               `json:"noAutonomousOrchestration"`
}

// Bundle groups approvals for evidence export.
type Bundle struct {
	WorkflowID                string           `json:"workflowId"`
	Sprint                    string           `json:"sprint"`
	GeneratedAt               string           `json:"generatedAt"`
	ReadOnly                  bool             `json:"readOnly"`
	NoAutonomousOrchestration bool             `json:"noAutonomousOrchestration"`
	Approvals                 []ActionApproval `json:"approvals"`
}
