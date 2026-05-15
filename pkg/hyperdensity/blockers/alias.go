// Package blockers re-exports the contractkit blocker catalog for in-repo callers.
// Prefer importing github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit/blockers externally.
package blockers

import ck "github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit/blockers"

const (
	SeverityCritical = ck.SeverityCritical
	SeverityHigh     = ck.SeverityHigh
	SeverityMedium   = ck.SeverityMedium
	SeverityInfo     = ck.SeverityInfo

	IDNoWindowsLane                  = ck.IDNoWindowsLane
	IDNoProductionMutation           = ck.IDNoProductionMutation
	IDKeepWindowsLaneDisabled        = ck.IDKeepWindowsLaneDisabled
	IDWindowsDisabled                = ck.IDWindowsDisabled
	IDDryRunOnly                     = ck.IDDryRunOnly
	IDRuntimeApplyDisabled           = ck.IDRuntimeApplyDisabled
	IDUnsupportedBroadVMExecution    = ck.IDUnsupportedBroadVMExecution
	IDUnsupportedBroadMemoryExec     = ck.IDUnsupportedBroadMemoryExec
	IDUnsupportedMultiContainerWiden = ck.IDUnsupportedMultiContainerWiden
	IDUnsupportedBroadAutomation     = ck.IDUnsupportedBroadAutomation
)

type Blocker = ck.Blocker

var (
	Known    = ck.Known
	Severity = ck.Severity
	Catalog  = ck.Catalog
)
