# Windows Fluid Node Actuator v1

`KARLNodeFluidActuator` is the controlled node-local backend contract for CPU entitlement mutations in the Windows Prearmed Fluid Envelope model.

## Scope

- Node-local only; no cross-node mutation.
- Allowlisted shells only.
- VM mapping is mandatory: VM -> virt-launcher pod -> container scope -> cgroup path -> QEMU PID/start.
- Allowed writes are constrained to `cpu.max` and optionally `cpu.weight`.

## Contract highlights

- Request identity: `actuatorId`, `requestId`, `actuatorVersion`, `createdAt`.
- Runtime target identity: `vmRef`, `shellRef`, `namespace`, `virtLauncherPodRef`, `podUid`, `qemuPid`, `qemuStartTime`, `cgroupPath`.
- Policy controls: `allowlistDecision`, `policyDecision`, `mutationScope`, `allowedControllers`.
- Mutation values: `previousCpuMax`, `requestedCpuMax`, `appliedCpuMax`, `rollbackCpuMax`.
- Evidence chain: `beforeEvidence`, `afterEvidence`, `rollbackEvidence`, `returnToFloorEvidence`, `blockers`.

## MVP actuator request

`KARLNodeFluidActuatorRequest` fields:

- `requestId`, `requestVersion`, `action`
- identity: `namespace`, `vmName`, `vmUid`, `vmiUid`, `podName`, `podUid`, `nodeName`, `qemuPid`, `qemuStartTime`
- target: `cgroupPath`, `controller` (`cpu.max` only)
- values: `previousCpuMax`, `requestedCpuMax`, `rollbackCpuMax`, `minCpuMax`, `maxCpuMax`
- policy/evidence: `ttlSeconds`, `createdAt`, `expiresAt`, `reason`, `risk`, `evidenceRefs`, `policyVersion`, optional `attestationRef`

## MVP actuator allowlist

`KARLNodeFluidActuatorAllowlist` fields:

- `allowlistId`
- identity pinning: `nodeName`, `namespace`, `vmName`, `vmUid`, `podUid`, `qemuPid`, `qemuStartTime`
- target pinning: `allowedCgroupPath`, `allowedControllers`
- bounds and controls: `minCpuMax`, `maxCpuMax`, `allowParentCgroupWrite`, `allowArbitraryWrite`, `allowSymlinkTraversal`
- lifecycle: `allowedActions`, `createdAt`, `expiresAt`

## Safety validator

`ValidateNodeFluidActuatorRequest` enforces:

- stale request rejection (`ttlSeconds` + `createdAt` + `expiresAt`)
- kill-switch blocking
- allowlist identity match (node/namespace/vm/pod/qemu)
- exact cgroup path match
- symlink/path-escape rejection
- target file limited to `cpu.max`
- parent/arbitrary write rejection
- controller allowlist enforcement
- request and allowlist bounds enforcement
- `previousCpuMax` readback consistency
- safe audit output path check for CLI evidence output

## Lifecycle result model

`KARLNodeFluidActuatorResult` fields:

- `resultId`, `requestId`, `action`, `decision`
- `dryRun`, `mutationPerformed`
- `previousCpuMax`, `requestedCpuMax`, `observedBeforeCpuMax`, `observedAfterCpuMax`
- `rollbackCpuMax`, `returnToFloorCpuMax`
- `qemuPid`, `qemuStartTime`, `podUid`, `nodeName`
- `evidenceRefs`, `blockers`, deterministic `auditHash`, `createdAt`

Decision values:

- `accepted`, `rejected`, `applied`, `rolled_back`, `returned_to_floor`, `blocked`

## Dry-run/apply/rollback/return-to-floor

- `dry-run`: validates only, never mutates.
- `apply`: writes requested `cpu.max`, requires before/after readback.
- `rollback`: writes `rollbackCpuMax`.
- `return-to-floor`: restores floor target (request `previousCpuMax`).

## Kill switch and TTL

- kill switch file values (`true`, `1`, `block`, `blocked`) force rejection.
- TTL is mandatory (`ttlSeconds > 0`) and must match `createdAt -> expiresAt`.

## Audit evidence

- CLI can emit JSON evidence (`-evidence-out`) containing full lifecycle result.
- `auditHash` is deterministic over result payload and supports local append-only auditing.

## Why this is not vCPU hotplug

- No logical processor count change is requested.
- No vCPU add/remove operation is emitted.
- Entitlement changes happen at host cgroup level (`cpu.max`) while guest topology remains stable.

## Safety expectations

- Reject ambiguous targets.
- Reject non-owned cgroup path resolution.
- Reject stale or replayed requests.
- Keep panic-safe rollback and return-to-floor evidence mandatory.
