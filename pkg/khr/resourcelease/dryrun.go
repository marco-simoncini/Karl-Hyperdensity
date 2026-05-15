package resourcelease

import (
	"encoding/json"
	"fmt"
	"strings"
)

// LeaseInput is a minimal ResourceLease-shaped document for dry-run (Sprint 5).
type LeaseInput struct {
	Spec struct {
		Donor struct {
			Kind      string `json:"kind"`
			Namespace string `json:"namespace"`
			Name      string `json:"name"`
		} `json:"donor"`
		Receiver struct {
			Kind      string `json:"kind"`
			Namespace string `json:"namespace"`
			Name      string `json:"name"`
		} `json:"receiver"`
		Resource string          `json:"resource"`
		Mode     string          `json:"mode"`
		Amount   json.RawMessage `json:"amount,omitempty"`
	} `json:"spec"`
}

// ResourcePortInput describes allowed modes for cpu/memory (subset of CRD contract).
type ResourcePortInput struct {
	Spec struct {
		Ports struct {
			CPU struct {
				Modes []string `json:"modes"`
			} `json:"cpu"`
			Memory struct {
				Modes []string `json:"memory"`
			} `json:"memory"`
		} `json:"ports"`
	} `json:"spec"`
}

// CellContext carries donor/receiver runtime hints for Linux-only enforcement.
type CellContext struct {
	DonorPlatform    string `json:"donorPlatform"`
	ReceiverPlatform string `json:"receiverPlatform"`
}

// DryRunResult is structured JSON output for CLI and tests.
type DryRunResult struct {
	Allowed           bool     `json:"allowed"`
	Blocked           bool     `json:"blocked"`
	Reason            string   `json:"reason,omitempty"`
	ExpectedWrites    []string `json:"expectedWrites"`
	RollbackPlan      []string `json:"rollbackPlan"`
	VerificationPlan  []string `json:"verificationPlan"`
	ResourcePortCheck string   `json:"resourcePortCheck,omitempty"`
	Notes             []string `json:"notes,omitempty"`
}

func modeSliceContains(modes []string, want string) bool {
	for _, m := range modes {
		if strings.EqualFold(m, want) {
			return true
		}
	}
	return false
}

// DryRun evaluates a ResourceLease-shaped input with safety rules (no writes).
func DryRun(lease *LeaseInput, port *ResourcePortInput, ctx *CellContext) DryRunResult {
	res := DryRunResult{
		ExpectedWrites:   []string{},
		RollbackPlan:     []string{"revert cgroup cpu.max/memory.max to prior values captured pre-apply"},
		VerificationPlan: []string{"read back cgroup limits", "compare cgroup.events pressure", "confirm process RSS/CPU throttle metrics"},
		Notes:            []string{},
	}
	if lease == nil {
		res.Blocked = true
		res.Reason = "lease input is nil"
		return res
	}
	if ctx == nil {
		ctx = &CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"}
	}
	if strings.ToLower(ctx.DonorPlatform) != "linux" || strings.ToLower(ctx.ReceiverPlatform) != "linux" {
		res.Blocked = true
		res.Reason = "non-linux donor/receiver blocked in KHR Linux MVP"
		return res
	}
	if lease.Spec.Resource != "cpu" && lease.Spec.Resource != "memory" {
		res.Blocked = true
		res.Reason = fmt.Sprintf("resource %q not supported in linux envelope MVP", lease.Spec.Resource)
		return res
	}
	if strings.ToLower(lease.Spec.Mode) != "envelope" {
		res.Blocked = true
		res.Reason = fmt.Sprintf("mode %q blocked: only envelope mode allowed in Sprint 5", lease.Spec.Mode)
		return res
	}
	if lease.Spec.Donor.Kind == "" || lease.Spec.Donor.Name == "" || lease.Spec.Receiver.Kind == "" || lease.Spec.Receiver.Name == "" {
		res.Blocked = true
		res.Reason = "donor/receiver kind and name are required"
		return res
	}
	if port == nil {
		res.Blocked = true
		res.Reason = "ResourcePort input required for Sprint 5 dry-run verification"
		return res
	}
	{
		isCPU := lease.Spec.Resource == "cpu"
		var ok bool
		if isCPU {
			ok = modeSliceContains(port.Spec.Ports.CPU.Modes, "envelope")
		} else {
			ok = modeSliceContains(port.Spec.Ports.Memory.Modes, "envelope")
		}
		if !ok {
			res.Blocked = true
			res.Reason = "ResourcePort does not allow envelope for requested resource"
			res.ResourcePortCheck = "incompatible"
			return res
		}
		res.ResourcePortCheck = "compatible"
	}

	res.Allowed = true
	res.Blocked = false
	res.Reason = "dry-run only: no cgroup writes performed"
	res.ExpectedWrites = []string{
		"(simulated) write cgroup cpu.max delta on donor slice",
		"(simulated) write cgroup memory.max delta on donor slice",
		"(simulated) mirror inverse delta on receiver slice",
	}
	return res
}
