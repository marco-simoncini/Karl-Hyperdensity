# Windows FluidVirt Product Model v1

`KARL Windows Prearmed Fluid Envelope v2` is integrated in Hyperdensity planning flow:

`observe -> assess compliance -> prepare lease -> build action slate -> dry-run actuator/QMP -> apply plan -> verify -> rollback/return-to-floor -> audit bundle append`

This contract milestone is model/evaluator only (no execution).

## WindowsHyperdensityTarget

`WindowsHyperdensityTarget` represents one Windows shell target (standalone or pool-child) with:

- runtime identity (`podUid`, `qemuPid`, `qemuStartTime`, `machineGuidHash`, `lastBootTime`)
- resource envelope (`cpu.floor/ceiling/current`, `memory.floor/ceiling/current`)
- mechanism constraints:
  - CPU: `cgroup-v2-cpu-max`
  - RAM: `qmp-balloon`
- continuity guarantees (`sameQemu`, `sameBoot`, `noLiveMigration`, `noRecreate`, `noRollout`)
- compliance binding (`compliancePhase`, `hyperdensityReady`)

## WindowsFluidResourceLease

`WindowsFluidResourceLease` models lease intent and safety:

- lease kinds: `cpu-entitlement`, `ram-balloon`, `combined-envelope`
- requested and previous state snapshots
- mandatory rollback and return-to-floor targets
- policy snapshot and audit refs
- status progression (`prepared` to `blocked/quarantined`)

Rules:

- CPU lease uses Node Fluid Actuator only
- RAM lease uses QMP balloon only
- no VM spec patch
- no vCPU hotplug
- no migration path as scaling mechanism

## WindowsFluidActionSlate

`WindowsFluidActionSlate` is planned action graph with required evidence and rollback links:

- `complianceReplay`
- `buildActuatorRequest`
- `actuatorDryRun`
- `cpuEntitlementApply`
- `cpuReturnToFloor`
- `qmpBalloonApply`
- `ramReturnToFloor`
- `guestVerify`
- `finalRestore`
- `auditBundleAppend`

`mutationAllowed` marks planned mutation-capable actions only; evaluators never execute them.

## Node Fluid Actuator Integration

CPU liquidity path integrates the hardened actuator contracts:

- `KARLNodeFluidActuatorRequest`
- `KARLNodeFluidActuatorAllowlist`
- `ValidateNodeFluidActuatorRequest`
- `KARLNodeFluidActuatorResult`

## Guest Witness Role

Guest evidence is required for ACK and continuity witness; missing guest witness blocks lease readiness.

## Rollback and Return-to-Floor

Both are mandatory preconditions and explicit action links in the slate.

## Audit Bundle Append

Uses deterministic hash-chain bundle model:

- validate existing chain
- append run
- enforce `previousRunHash`
- reject broken chain/duplicate runs
- optional write-back output path via replay CLI

## Pool-child semantics

- pool-child accepted as individual VM target
- pool scaling mechanism blocked

## Why no vCPU hotplug / logical CPU scaling / pool scaling

- Windows product path is entitlement liquidity (`cpu.max`) plus balloon liquidity (QMP), not topology mutation.
- Logical processor scaling claims are explicitly rejected.
- Pool scaling is provisioning context only, not runtime mechanism.
