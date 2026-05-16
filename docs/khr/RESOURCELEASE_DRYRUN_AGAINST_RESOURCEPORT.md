# ResourceLease dry-run against ResourcePort CR (KHR-L)

Validate **ResourceLease** specs against **cluster-scoped ResourcePort CRs** applied in `khr-runtime-sandbox`. **No apply**, **no mutation**, **no autonomous ResourceLease reconcile**.

## Mode

```bash
karl-host-runtime -mode=resourcelease-dryrun \
  -config=examples/khr/runtime-sandbox/karl-host-runtime-config-loop.yaml \
  -lease-input=examples/khr/runtime-sandbox/resourcelease-dryrun-allowed.json \
  -namespace=khr-runtime-sandbox \
  -cluster-context=karl-metal-01@ovh
```

## Decision output

| Field | Description |
|-------|-------------|
| `dryRunDecision` | `allowed` \| `blocked` |
| `allowed` / `blocked` | Boolean gates |
| `reason` / `blockedReason` | Human-readable block cause |
| `baseline` | Sandbox baseline capture (no writes) |
| `rollbackPlan` / `verificationPlan` | Planned steps (read-only) |
| `rollbackPlanRef` / `rollbackPlanStatus` | From lease governance when present |
| `verificationPlanRef` / `verificationPlanStatus` | Planned verification |
| `sourceResourcePortRef` | Matched `cluster/ResourcePort/<name>` |
| `noMutation` / `noApply` | Always true for this mode |

## Lease evaluation

| Input | Checked |
|-------|---------|
| `leaseKind` | `transfer` or `runtime` |
| `transfer.resource` | `cpu` \| `memory` |
| `transfer.mode` | `envelope` (via inner dry-run) |
| `transfer.amount` | Sandbox limits (cpu ≤ 500 milliCpu) |
| Target | `khr.karl.io/resource-port-ref` annotation, or `shellRef` / `cellRef` match on port CR |
| Donor/receiver | Required kinds/names |

## Safety gates

- `sandboxMode` + `linuxOnly`
- Namespace allowlist (`khr-runtime-sandbox`)
- Label allowlist on **lease** metadata (`khr.karl.io/sandbox=true`)
- Cluster context guard (`karl-metal-01@ovh`)
- Production namespace blocklist
- ResourcePorts listed by `karl.io/sandbox-namespace=<ns>`

## Evidence

`./scripts/khr_resourcelease_dryrun_evidence.sh` — apply sandbox ResourcePorts, allowed dry-run, blocked cases, cleanup, production mutation proof.

## Related

- KHR-K: `docs/khr/RESOURCEPORT_CR_PREVIEW.md`
- KHR-J: `docs/khr/RESOURCEPORT_REPORTING_LOOP.md`
