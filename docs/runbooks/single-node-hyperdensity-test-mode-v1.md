# Single-Node Hyperdensity Test Mode v1

## Objective

Run Hyperdensity validation in a temporary single-node posture while preserving operator safety and auditability.

## Preconditions

- Parent Fabric authentication validated (HTTP 200).
- Control Room routes reachable.
- Remaining node capacity validated for temporary concentration.
- Explicit operator approvals for any VMI stop/delete actions.

## Hard Safety Rules

- No blind delete operations.
- No forced live migration.
- No force-delete pod operations.
- No production workload mutations.
- Keep console/auth path healthy throughout execution.

## Recommended Flow

1. Capture baseline inventory (nodes, pods, VMIs, migrations, unknown/high-risk objects).
2. Relocate only clearly safe controller-managed platform workloads.
3. Stop/delete only explicitly approved test VMIs/VMs.
4. Delete only allowlisted non-essential stale objects.
5. Re-run post-action validation (Parent Fabric + routes + inventories).
6. Evaluate node-removal gate.
7. Perform cordon -> drain -> remove only when gate is fully green.

## Blocked Outcome Handling

When gate is blocked:

- Produce a precise blocker summary.
- Preserve all before/after artifacts.
- Publish remediation steps with owner and required approval.

## Minimum Evidence Set

- Auth + Parent Fabric before/after.
- Route validation before/after.
- Platform/VMI/migration inventories before/after.
- Cleanup action logs.
- Final gate decision and remediation plan.
# Single-Node Hyperdensity Test Mode v1

## Purpose

This runbook defines a controlled single-node development/test mode for KARL Hyperdensity.

It is designed to continue product/API/UI work when one node is intentionally removed or made non-participating.

## Posture

- Not HA.
- Not GA.
- Development/test only.
- Operator-controlled only.

## Mandatory Safety Rules

- No production mutation.
- No autonomous apply.
- Enforcement remains disabled.
- No forced live migration.
- No raw runtime controls.
- Allowlist-driven cleanup only (no blind cluster-wide deletes).

## Required Pre-Destructive Checks

Before any cleanup/cordon/drain/remove action:

1. Restore authenticated Parent Fabric proof (`HTTP 200`).
2. Validate auth path (`/auth/login` OIDC flow).
3. Capture pre-action inventory and rollback artifacts.
4. Classify objects by risk and ownership.
5. Block unknown/high-risk object deletion.

## VM/VMI Handling

- Every VMI on target node requires per-VM controlled handling.
- Prefer graceful stop when explicitly safe.
- Do not force live migration.
- Do not delete unknown VMIs.
- Preserve KubeVirt CRDs/controllers required for VM evidence/testing.

## Expected Runtime Behavior

- Degraded/no-live-migration states are acceptable and should be surfaced clearly.
- Safety blocks are correct outcomes when destructive actions are unsafe.
- No false readiness claims.

## Cleanup and Node Action Flow

1. Auth restore + Parent Fabric `HTTP 200`.
2. Platform relocation plan (console/auth/ingress/kubevirt core).
3. VMI/VM plan and controlled execution.
4. Allowlist cleanup execution only.
5. Re-run removal preflight.
6. Cordon/drain/remove only if gate explicitly allows.

## Validation Requirements

- Console/auth path remains reachable.
- Parent Fabric remains `HTTP 200`.
- Control Room routes remain reachable.
- No forced migration observed.
- No autonomous/enforcement changes.
- No production mutation.
