# Windows FluidVirt lane — reconciliation (post Sprint 14)

## Context

Sprint 14 merged **The-Father-Windows** runtime types into `KHR` and dropped
FVI-only tests that targeted a different `WindowsFluidVirt` planning-only API.
Sprint 15 **reconciles** the lane: the same TF-W structs remain the source of
truth; planning-only semantics are expressed via **`safety.go`** helpers and
new tests (`planning_safety_test.go`, `blocker_catalog_test.go`,
`no_windows_claims_test.go`).

## What changed

- **`pkg/windowsfluidvirt/safety.go`**: read-only mapping from TF-W
  `WindowsFluidActionSlate` (+ optional target/lease) to a conceptual
  planning-safety model (no parallel Hyperdensity product fork).
- **`BuildWindowsFluidActionSlate`**: every `WindowsFluidAction` now keeps
  `mutationAllowed: false` at JSON level so the default slate is a **typed
  readiness graph** without execution flags on steps. Action **types** still
  name future gated phases (CPU entitlement apply, QMP balloon, etc.).

## What did not change

- Linux Hyperdensity / KHR cgroup evidence path is untouched.
- No Dashboard runtime changes.
- No Windows apply enablement, no autonomous apply, no raw runtime control
  surface in this package beyond what TF-W already models under disabled
  executor / governance gates elsewhere.

## FVI fields not on TF-W

The older FVI struct had explicit `PlanningOnly`, `ApplyEnabled`,
`AutonomousApplyEnabled`, `RawRuntimeControlsExposed`, etc. Those names do not
exist on TF-W types; `DerivePlanningSafety` documents the mapping and derives
booleans for tests and audits without inventing new persisted fields.
