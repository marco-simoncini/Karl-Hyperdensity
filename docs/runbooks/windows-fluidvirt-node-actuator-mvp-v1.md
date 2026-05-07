# Windows FluidVirt Node Actuator MVP Runbook v1

## Goal

Execute governed CPU entitlement liquidity using local MVP actuator semantics without cluster/QMP calls.

## Prepare request

1. Build `KARLNodeFluidActuatorRequest`.
2. Set action: `dry-run`, `apply`, `rollback`, or `return-to-floor`.
3. Pin identity fields (`namespace`, `vmName`, `vmUid`, `vmiUid`, `podUid`, `qemuPid`, `qemuStartTime`).
4. Set `controller=cpu.max`.
5. Set `previousCpuMax`, `requestedCpuMax`, `rollbackCpuMax`.
6. Set `minCpuMax`, `maxCpuMax`.
7. Set `ttlSeconds`, `createdAt`, `expiresAt`.

## Prepare allowlist

1. Build `KARLNodeFluidActuatorAllowlist`.
2. Pin exact identity and `allowedCgroupPath`.
3. Keep:
   - `allowParentCgroupWrite=false`
   - `allowArbitraryWrite=false`
   - `allowSymlinkTraversal=false`
4. Allow `cpu.max` only and required actions only.

## Dry-run

```bash
go run ./cmd/karl-node-fluid-actuator \
  -mode dry-run \
  -request request.json \
  -allowlist allowlist.json \
  -kill-switch kill_switch.txt \
  -evaluation-time 2026-05-07T22:05:00Z \
  -evidence-out dry_run.json
```

Expected:

- `decision=accepted`
- `mutationPerformed=false`

## Apply

Run same command with `-mode apply`.

Expected:

- `decision=applied`
- before/after readback is present
- `mutationPerformed=true`

## Rollback

Run with `-mode rollback`.

Expected:

- `decision=rolled_back`
- target returns to `rollbackCpuMax`

## Return-to-floor

Run with `-mode return-to-floor`.

Expected:

- `decision=returned_to_floor`
- target returns to floor (`previousCpuMax` in request)

## Validate audit

- confirm `auditHash` is present in each result
- persist request, allowlist, and result JSON files together
- use compliance replay bundle append for deterministic multi-run audit chain

## Forbidden

- vCPU hotplug/unplug
- logical CPU scaling claim
- pool scaling as mechanism
- LiveMigration/VMIM as mechanism
- writes outside allowlisted `cpu.max`
- parent cgroup writes
- arbitrary file writes
