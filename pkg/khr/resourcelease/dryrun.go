package resourcelease

import (
	"fmt"
	"strings"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
)

// CellContext carries donor/receiver hints for Linux-only enforcement and optional CRD snapshots.
type CellContext struct {
	DonorPlatform           string                       `json:"donorPlatform"`
	ReceiverPlatform        string                       `json:"receiverPlatform"`
	DonorCell               *crdv1alpha1.Cell            `json:"donorCell,omitempty"`
	ReceiverCell            *crdv1alpha1.Cell            `json:"receiverCell,omitempty"`
	DonorRuntimeProvider    *crdv1alpha1.RuntimeProvider `json:"donorRuntimeProvider,omitempty"`
	ReceiverRuntimeProvider *crdv1alpha1.RuntimeProvider `json:"receiverRuntimeProvider,omitempty"`
}

// DryRunResult is structured JSON output for CLI and tests.
type DryRunResult struct {
	Allowed           bool     `json:"allowed"`
	Blocked           bool     `json:"blocked"`
	Reason            string   `json:"reason,omitempty"`
	ExpectedWrites    []string `json:"expectedWrites,omitempty"`
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

// DryRun evaluates a v1alpha1 ResourceLease with safety rules (no writes).
func DryRun(lease *crdv1alpha1.ResourceLease, port *crdv1alpha1.ResourcePort, ctx *CellContext) DryRunResult {
	res := DryRunResult{
		ExpectedWrites:   nil,
		RollbackPlan:     []string{"revert cgroup cpu.max/memory.max to prior values captured pre-apply"},
		VerificationPlan: []string{"read back cgroup limits", "compare cgroup.events pressure", "confirm process RSS/CPU throttle metrics"},
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
	ls := lease.Spec
	donor, receiver, resource, mode, ok := ls.EffectiveTransfer()
	if !ok {
		res.Blocked = true
		res.Reason = "donor/receiver kind and name are required (spec.transfer or legacy inline fields)"
		return res
	}
	if resource != "cpu" && resource != "memory" {
		res.Blocked = true
		res.Reason = fmt.Sprintf("resource %q not supported in linux envelope MVP", resource)
		return res
	}
	if strings.ToLower(mode) != "envelope" {
		res.Blocked = true
		res.Reason = fmt.Sprintf("mode %q blocked: only envelope mode allowed in KHR Linux MVP", mode)
		return res
	}
	if donor.Kind == "" || donor.Name == "" || receiver.Kind == "" || receiver.Name == "" {
		res.Blocked = true
		res.Reason = "donor/receiver kind and name are required"
		return res
	}
	if port == nil {
		res.Blocked = true
		res.Reason = "ResourcePort input required for dry-run verification"
		return res
	}
	isCPU := resource == "cpu"
	var portOK bool
	if isCPU {
		portOK = modeSliceContains(port.Spec.Ports.CPU.Modes, "envelope")
	} else {
		portOK = modeSliceContains(port.Spec.Ports.Memory.Modes, "envelope")
	}
	if !portOK {
		res.Blocked = true
		res.Reason = "ResourcePort does not allow envelope for requested resource"
		res.ResourcePortCheck = "incompatible"
		return res
	}
	res.ResourcePortCheck = "compatible"

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
