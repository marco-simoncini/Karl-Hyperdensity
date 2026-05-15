# Windows FluidVirt — planning-only safety model

## Strategic framing

The **Windows FluidVirt** lane is a **strategic asset**: it carries contracts,
evidence, dry-run, and readiness work for Windows VM workloads in the KARL
ecosystem. It is **not** an assertion of **active Hyperdensity production** on
Windows, and it must **not** be read as a GA or production-ready claim.

## Hard boundaries (lane policy)

| Policy | Meaning on `KHR` |
|--------|-------------------|
| No production Hyperdensity claim | `HyperdensityReady` and related gates are evidence-driven; lane docs and helpers never declare Windows GA. |
| No apply | Default slates keep `applyAllowed: false` at envelope level; executor paths remain hard-disabled until future KHR apply gates. |
| No raw runtime controls | No new exposure of arbitrary QMP / cgroup / cluster knobs beyond existing TF-W contracts under governance. |
| No Windows CPU/RAM hotplug promise | Target and lease fields that imply vCPU hotplug, logical CPU scaling-as-mechanism, or pool-scaling-as-mechanism are **rejected** with canonical blockers. |
| Planning / readiness / evidence | Default `BuildWindowsFluidActionSlate` is a **planning graph**: action **types** may reference future gated steps, but **`mutationAllowed` is false on every step** and envelope mutation/apply stay false. |

## TF-W → planning-only mapping

Implemented in `pkg/windowsfluidvirt/safety.go`:

- **`ApplyEnabled`** ↔ `WindowsFluidActionSlate.ApplyAllowed` (must be `false`).
- **`RuntimeMutationEnabled`** ↔ `RuntimeMutationExecuted \|\| MutationAllowed` on the envelope slate (both must be `false` in lane).
- **`AutonomousApplyEnabled`**, **`RawRuntimeControlsExposed`**, **`ProductionReadyClaim`** ↔ no distinct TF-W fields; derived layer keeps them **false** for this lane (documentation + tests).
- **Forbidden claims** (`VCPUHotplugClaim`, `LogicalCPUScalingClaim`, `PoolScalingClaim`) ↔ `WindowsHyperdensityTarget` and `WindowsFluidResourceLease` flags; `DerivePlanningSafetyWithTargetLease` marks `PlanningOnly` false when any are set.

## Violations API

`EvaluatePlanningSafetyViolations(slate)` returns strings when envelope or step
flags violate the planning lane. Tests use it to prove unsafe manual mutation
is classified as blocked at the model layer.

## Related code

- `product_model.go` — `BuildWindowsFluidActionSlate`, lease/target evaluators.
- `blockers.go` — `CanonicalBlockers` catalog.
- `unlock_gate_verification.go` / `governance_contract.go` — executor remains
  disabled in current phase (see existing gate tests).
