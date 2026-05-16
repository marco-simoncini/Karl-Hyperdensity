# RAM live scale sandbox lane (KHR-O)

KHR-O adds a **memory** live-scale lane beside the existing CPU cgroup lane. All operations are **sandbox-only**, **live-in-place**, and explicitly forbid restart, rollout, recreate, and production mutation.

## Dry-run

`resourcelease-dryrun` accepts:

| Field | Values |
|-------|--------|
| `resource` | `memory` |
| `mode` | `scaleUp`, `scaleDown`, `envelope` |
| `amount.bytes` | positive delta (scaleUp/scaleDown) or envelope target |

ResourcePort `spec.ports.memory.modes` must include the requested mode (or `envelope`).

## Guarded apply

When `--apply-resourcelease` and `--i-understand-this-is-sandbox` are set against `khr-runtime-sandbox`:

1. Capture baseline `memory.max` and `memory.high`.
2. Compute target from mode + delta (`cgroup.ComputeMemoryTarget`).
3. Write `memory.high` then `memory.max` under the lease cgroup path.
4. Read back and verify live update.
5. Persist baseline JSON for `resourcelease-rollback`.

CPU apply (`resource=cpu`, `mode=envelope`) is unchanged from KHR-M.

## Safety gates

- `sandboxMode` + `linuxOnly` + namespace allowlist + label allowlist
- `sandboxMaxMemoryDeltaBytes` (default 512Mi) caps per-request delta
- Production namespaces blocked
- Lease annotations `khr.karl.io/restart-required`, `rollout-required`, `recreate-required` => **blocked**
- No pod restart, no Deployment rollout, no VM reboot on this path

## Evidence

```bash
./scripts/khr_ram_live_scale_evidence.sh
```

Produces CPU+RAM combined scenario, live cgroup proof, and no-restart/no-rollout attestation under `docs/evidence/khr-ram-live-scale/`.
