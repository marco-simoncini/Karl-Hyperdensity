# ResourceLease guarded apply in sandbox (KHR-M)

Opt-in **CPU envelope** apply under `khr-runtime-sandbox` after **ResourcePort CR** binding and **dry-run allowed**. No automatic apply. No production mutation.

## Modes

| Mode | Purpose |
|------|---------|
| `resourcelease-dryrun` | Evaluate lease against cluster ResourcePorts (KHR-L) |
| `resourcelease-guarded-apply` | Dry-run + optional cgroup `cpu.max` write |
| `resourcelease-rollback` | Restore baseline `cpu.max` |

## Guarded apply gates (all required)

1. `--apply-resourcelease=true`
2. `--i-understand-this-is-sandbox`
3. Namespace `khr-runtime-sandbox` (allowlist + production blocklist)
4. Label allowlist on lease metadata
5. Cluster context `karl-metal-01@ovh`
6. Internal dry-run **allowed**
7. `rollbackPlanRef` + non-empty `rollbackPlan`
8. `verificationPlanRef` + non-empty `verificationPlan`
9. Resource **cpu** only; `milliCpu` ≤ **500**

## Apply behavior

- Captures **baseline** (`cpu.max` before) under `--sandbox-dir`
- Writes cgroup v2 `cpu.max` under `allowPathPrefixes` / sandbox fallback
- **Verify**: read-back `cpu.max`, `noRestart: true`
- **Flight recorder** events on CLI output

## Rollback

```bash
karl-host-runtime -mode=resourcelease-rollback \
  -config=examples/khr/runtime-sandbox/karl-host-runtime-config-guarded-apply.yaml \
  -sandbox-dir=/tmp/khr-resourcelease-guarded-apply \
  -baseline-id=sandbox-default
```

## Evidence

`./scripts/khr_resourcelease_guarded_apply_evidence.sh`

## Safety

- No ResourceLease `kubectl apply` to cluster
- No production namespaces
- ISO: preview disabled — see Karl-OS-ISO `KHR_HOST_RUNTIME_PREVIEW.md`
