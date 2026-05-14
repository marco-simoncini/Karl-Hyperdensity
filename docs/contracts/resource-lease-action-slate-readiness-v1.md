# KARL Hyperdensity — Resource Lease + Action Slate Runtime Readiness v1

**Contract ID:** `hyperdensity_resource_lease_action_slate_readiness_v1`  
**Milestone:** `hyperdensity_resource_lease_action_slate_readiness_v1`  
**Release track:** `technical_preview`

## Product definition

> KARL Hyperdensity can generate operator-controlled resource lease candidates and action slate entries from governed runtime shells, with dry-run, rollback, SLO precheck, risk and claim boundaries before any apply.

## Concepts

| Concept | Owner | Description |
|---------|-------|-------------|
| Donor index | Karl-Hyperdensity | Filtered lendable shells derived from shell registry |
| Receiver index | Karl-Hyperdensity | Filtered receivers and remediable receivers with pressure class |
| Resource lease candidate | Karl-Hyperdensity | Proposed lease with dry-run, rollback, SLO precheck gates |
| Action slate entry | Karl-Hyperdensity | Operator-facing action with readiness statuses |
| Action slate readiness | Karl-Hyperdensity | Top-K bounded pairing surface for Dashboard projection |
| Dry-run readiness | FluidVirt | Evidence-only actuator path resolution; no mutation executed |
| Shell registry | Karl-Hyperdensity | Upstream enrolled shell source |
| Registry projection | Karl-Dashboard | Read-only ConfigMap surface; not source of truth |

## Donor eligibility rules

A shell may appear as donor only if:

- `enrolled=true`
- `identityStable=true`
- `blocked=false`
- `protected=false`
- `donorEligibility=eligible`
- `supportProfile` exists
- `claimBoundary` exists
- `rollbackReadiness` is not `unknown`
- capability evidence exists (`evidenceRefs` non-empty)

## Receiver eligibility rules

A shell may appear as receiver only if:

- `enrolled=true`
- `identityStable=true`
- `blocked=false`
- `receiverEligibility=eligible` **or** remediable with explicit remediation path
- `supportProfile` exists
- `claimBoundary` exists
- capability evidence exists

## Action generation exclusions

Exclude from apply-oriented action generation when:

- `blocked=true` and no explicit remediation-only action
- `protected=true` and no explicit operator override
- missing `supportProfile`, `claimBoundary`, or `evidenceRefs`
- reboot/recreate/migration required for declared movement
- excluded shell kind
- Windows master/template/controller/root pool member

## No full N×N pairing

The action slate generator must not pair every donor with every receiver. Reference behavior uses top-K donors and receivers, tenant/namespace affinity, protected/blocked filtering, and `maxEvaluatedPairs`.

## Sprint 3 apply prohibition

- `autoApplyAllowed=false` on all lease candidates and action entries
- `productionMutationAllowed=false` on all lease candidates and action entries
- `readyForApply` may be `operator_controlled_ready` only — never autonomous apply
- Dashboard must not expose apply controls or mutate action state

## Source-of-truth map

| Domain | Owner |
|--------|-------|
| Lease / action slate contracts | Karl-Hyperdensity |
| Dry-run readiness evidence | FluidVirt |
| Action slate projection | Karl-Dashboard |
| Identity / signals | Karl-Inventory |

## Forbidden claims (Sprint 3)

- guaranteed savings active
- universal performance improvement
- production autonomous apply
- production mutation enabled
- Windows total RAM hotplug supported
- logical vCPU hotplug supported
- 1000 production workloads proven
- action slate auto-applies resources
- Dashboard applies runtime changes
