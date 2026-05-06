# Runtime Authority No-Live-Migration Resilience v1

## Purpose

Define a safety contract for resilience actions when cluster capacity is constrained and live migration must not be forced.

## Scope

- Applies to Hyperdensity evidence/test environments and controlled single-node transitions.
- Covers node participation changes, VMI handling posture, and cleanup sequencing.

## Safety Invariants

1. Parent Fabric route must be authenticated and return HTTP 200 before any destructive step.
2. No forced live migration for VMIs.
3. No force-delete for pods.
4. No mutation of production workloads.
5. Core console/auth components must remain reachable and healthy.
6. Every destructive action must be allowlisted and auditable.

## Required Checks

- Node inventory for target and remaining nodes.
- Platform pod inventory on target node.
- VMI and VM ownership/classification on target node.
- Migration object state inventory.
- Unknown/high-risk object inventory with explicit disposition.

## Execution Rules

- Execute only actions classified as safe and reversible.
- Keep blocked items with precise reasons when classification or approval is missing.
- Re-validate Parent Fabric and Control Room routes after each action batch.

## Gate for Cordon/Drain/Remove

Cordon/drain/remove is allowed only when all conditions are true:

- Parent Fabric authenticated (HTTP 200) after latest action batch.
- Required platform workloads are safely relocated or explicitly drain-safe.
- VMIs on target are reduced to an approved safe set (ideally zero).
- Unknown/high-risk objects are reduced to zero or explicitly waived.

If any condition fails, verdict remains blocked and remediation must be documented.
# Runtime Authority No-Live-Migration Resilience v1

## Contract ID

- `hyperdensity_runtime_authority_no_live_migration_resilience_v1`
- `targetNodeName`: `karl-ids-metal-01`
- mode: controlled resilience proof

## Scope

This contract defines a controlled infrastructure resilience proof for KARL Live Resource Authority when a target node is made unavailable or considered for removal.

The goal is safety-correct behavior, not forced continuity.

## Required Safety Invariants

- No forced live migration.
- No unsafe fallback.
- No production mutation.
- No autonomous apply.
- Enforcement remains disabled.
- No raw runtime controls exposed (QMP/libvirt/QGA/cgroup).
- No false readiness claim for affected authority surfaces.
- Operator-controlled only.
- Windows lane remains out-of-scope/frozen.

## Expected Safe Outcomes

Any of the following are valid:

1. Controlled action proceeds safely without forced migration.
2. Controlled action is blocked by preflight with precise blockers.
3. Runtime authority degrades/blocks clearly with evidence and next operator action.

Blocked/degraded is a safe outcome and must not be treated as failure.

## Required Proof Shape

Suggested projection:

`executionEngine.hyperdensityLiveResourceAuthority.noLiveMigrationResilience`

Core fields:

- `resilienceId`
- `resilienceVersion`
- `targetNodeName`
- `operationMode=operator_controlled`
- `autonomousApplyAllowed=false`
- `enforcementMode=disabled`
- `productionMutationAllowed=false`
- `liveMigrationForced=false`
- `rawRuntimeControlsExposed=false`
- `preflightState`
- `nodeActionState`
- `liveMigrationPathState`
- `runtimeAuthorityState`
- `supportBoundaryState`
- `safetyGates[]`
- `blockers[]`
- `parentFabricBeforeState`
- `parentFabricAfterState`
- `resultSummary`
- `nextOperatorAction`

## Rollback and Rejoin

A rollback/rejoin plan is required before disruptive actions.
If such plan is missing or unsafe, node action must be blocked.

## VM RAM Authority Guardrail

VM RAM path remains runtime overlay authority where proven (virtio-mem / QMP / QOM evidence path) and does not broaden into generic KubeVirt template mutation claims.

## Release Posture

Technical Preview only. Not GA.
