# KARL Hyperdensity — Shell Passport Factory v1

**Contract ID:** `hyperdensity_shell_passport_factory_v1`  
**Milestone:** `hyperdensity_shell_passport_factory_v1`  
**Release track:** `technical_preview`

## Product definition

> A KARL Hyperdensity shell is a governed runtime envelope that can be evaluated as donor, receiver, protected, blocked or remediable before any resource movement is applied.

## Allowed Sprint 2 claim

> KARL Hyperdensity can enroll VM, DaaS pool replicas and containers as governed runtime shells with evidence-backed eligibility, blockers and remediation paths.

## Concepts

| Concept | Owner | Description |
|---------|-------|-------------|
| Shell passport | Karl-Hyperdensity | Per-shell governed envelope with identity, capability, eligibility, blockers, remediation, claim boundary |
| Shell factory | Karl-Hyperdensity | Enrollment rules, canonical kinds, required evidence, support profiles |
| Shell registry | Karl-Hyperdensity | Observed set of enrolled shells with aggregate counts |
| Enrollment | Karl-Hyperdensity | Identity verification + capability evidence load + claim boundary assignment |
| Runtime capability evidence | FluidVirt | Read-only actuator observations (cgroup, libvirt, guest, QGA) |
| Shell registry projection | Karl-Dashboard | Read-only ConfigMap surface; not source of truth |
| Identity / signals | Karl-Inventory | Endpoint and guest readiness signals only |

## Canonical shell kinds (Sprint 2)

- `container_linux`
- `container_linux_replica`
- `vm_linux`
- `vm_windows`
- `vm_windows_pool_replica`
- `daas_pool_replica`

## Excluded shell kinds (Sprint 2)

- `container_windows`
- `container_windows_replica`
- `vm_linux_pool_replica` (unless separately certified)
- Windows pool `master` / `template` / `controller` / `root` members
- Any shell without stable identity, `claimBoundary`, or `supportProfile`
- Any shell requiring reboot/recreate/migration for declared in-place movement

## Source-of-truth map

| Concern | Source of truth |
|---------|-----------------|
| Shell contracts & eligibility rules | Karl-Hyperdensity |
| Runtime capability evidence | FluidVirt |
| Registry projection | Karl-Dashboard (read-only) |
| Identity / access signals | Karl-Inventory |

## Safety invariants

- `productionMutationAllowed`: false
- `autoApplyAllowed`: false
- Blocked shells cannot be donor/receiver eligible
- No Windows total RAM hotplug or logical vCPU hotplug claims unless separately proven
