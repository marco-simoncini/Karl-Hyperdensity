# Windows FluidVirt Controlled Apply Plan v1

This contract introduces the first controlled apply planning pipeline for Windows FluidVirt, with explicit feature gates and manual approval.

## Scope

- Input models: `WindowsHyperdensityTarget` + `WindowsFluidResourceLease`
- New gate model: `WindowsFluidControlledApplyGate`
- New planning model: `WindowsFluidControlledApplyPlan`
- Evaluators produce plan/evidence only
- No autonomous apply and no runtime mutation in this milestone

## Controlled Apply Gate

`WindowsFluidControlledApplyGate` defines explicit policy toggles:

- requires dry-run first
- requires manual approval
- keeps `autonomousApplyEnabled=false`
- separately gates CPU apply, RAM apply, node actuator apply, and QMP balloon apply
- requires guest verification and workload verification
- requires rollback, return-to-floor, and audit bundle
- requires kill switch allow-state
- enforces blast-radius and allowlists (namespace, target, lease kind)

Default gate is safety-first and blocks apply unless explicit controlled fixture enables all required apply switches.

## Manual Approval

`EvaluateWindowsFluidManualApproval` supports:

- `not_required`
- `required`
- `approved`
- `rejected`
- `expired`

Apply can never be allowed when approval is missing/rejected/expired and manual approval is required.

## Dry-Run First

`EvaluateWindowsFluidDryRunGate` enforces dry-run prerequisites before readiness:

- dry-run action must exist in slate
- lease/planning blockers must be clear

## Apply Readiness

`EvaluateWindowsFluidApplyReadiness` blocks for:

- non-ready compliance
- pool scaling request
- autonomous apply enabled
- missing manual approval
- missing dry-run proof
- missing guest/workload verification
- missing rollback/return-to-floor/audit
- kill switch blocked
- missing CPU/RAM apply capability toggles
- forbidden mechanisms (vCPU hotplug, logical CPU scaling, VM spec patch)

## Verification, Rollback, Return-to-Floor, Audit

Dedicated evaluators build explicit plan sections:

- `EvaluateWindowsFluidVerificationPlan`
- `EvaluateWindowsFluidRollbackPlan`
- `EvaluateWindowsFluidReturnToFloorPlan`
- `EvaluateWindowsFluidAuditBundlePlan`

These are mandatory guardrails for controlled apply readiness.

## Why Autonomous Apply Remains Disabled

This milestone is policy-gated planning only. Autonomous apply remains disabled to prevent unsupervised mutation and to enforce explicit operator intent.

## Why No vCPU Hotplug / Logical CPU Scaling

Windows FluidVirt product path uses:

- CPU entitlement liquidity via node actuator (`cpu.max`)
- RAM balloon liquidity via QMP balloon

It does not rely on guest topology mutation (vCPU hotplug/unplug) or logical CPU scaling claims.
