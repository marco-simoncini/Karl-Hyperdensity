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

## Why this is not vCPU hotplug

- No logical processor count change is requested.
- No vCPU add/remove operation is emitted.
- Entitlement changes happen at host cgroup level (`cpu.max`) while guest topology remains stable.

## Safety expectations

- Reject ambiguous targets.
- Reject non-owned cgroup path resolution.
- Reject stale or replayed requests.
- Keep panic-safe rollback and return-to-floor evidence mandatory.
