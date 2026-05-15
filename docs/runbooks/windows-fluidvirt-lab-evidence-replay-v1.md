# windows-fluidvirt-lab-evidence-replay-v1

Runbook for evidence-only replay of Windows FluidVirt dry-run cases.

## Scope and safety

- Read-only local replay only.
- No deploy, no cluster mutation, no VM mutation.
- No `kubectl apply/patch/rollout`, no `helm upgrade`.
- No frontend or dashboard work.

## Fixture location

- `examples/windows-fluid-dryrun-fixtures/`

Lab cases include:

- `master-win11-*` candidate VM evidence.
- `win11-pool-context-only.blocked.json` context-only pool case.
- `identity-change.quarantined.json`.
- lease dry-run cases.

## CLI usage

```bash
go run ./cmd/karl-fluid-dryrun \
  -fixture ./examples/windows-fluid-dryrun-fixtures/master-win11-certification-ready.ready.json \
  -evaluation-time 2026-05-07T14:33:00Z
```

Alternative input:

```bash
go run ./cmd/karl-fluid-dryrun -bundle /path/to/runtime-evidence-bundle.json
```

## Interpreting outcomes

- `READY`: certification evidence complete.
- `BLOCKED`: missing proof or failing safety gate.
- `QUARANTINED`: continuity break, isolate candidate.
- `LEASE_PREPARED`: lease intent validated in dry-run only.

`ACTIVE` is impossible in this runbook.

## `master-win11` handling

- Treat as single candidate VM for same-runtime continuity proofs.
- Confirm guest ACK + read-only QMP + stable identity before any future apply-phase planning.

## `win11-pool-*` handling

- Treat only as context signal.
- Never use pool replica behavior as FluidVirt success mechanism.
- Expected output stays `BLOCKED_POOL_REPLICA_MODEL`.

## Minimum proofs before first future +CPU phase

- Same node, same virt-launcher pod, same qemu process.
- No live migration, no recreate, no reboot.
- Guest ACK ready and QMP read-only ready.
- Rollback ready and return-to-floor ready.

## Still forbidden in this phase

- Runtime CPU/RAM apply.
- QMP mutating commands.
- Hotplug execution.
- Dashboard/frontend changes.
