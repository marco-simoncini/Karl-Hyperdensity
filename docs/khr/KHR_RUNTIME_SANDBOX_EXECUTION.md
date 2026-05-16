# KHR runtime sandbox execution (KHR-G)

| Field | Value |
|-------|-------|
| **Cluster** | `karl-metal-01@ovh` |
| **Namespace** | `khr-runtime-sandbox` |
| **Label allowlist** | `khr.karl.io/sandbox=true` |
| **Default apply** | `sandboxApplyEnabled: false` |

---

## Guardrails

| Rule | Enforcement |
|------|-------------|
| Dedicated sandbox namespace | Scripts exit if `namespace != khr-runtime-sandbox` |
| Label allowlist | Namespace + workload must carry `khr.karl.io/sandbox=true` |
| No production mutation | Blocklist: `karl-system`, `kube-system`, `default`, … |
| No VM/QMP/libvirt | Not invoked; capabilities report only |
| No Windows | Linux `CellContext` only |
| No autonomous apply | Operator runs `khr_runtime_sandbox_execute.sh` |
| Rollback required | `khr_runtime_sandbox_rollback.sh` |
| Flight recorder required | CLI `-mode=flight-recorder` during dry-run |

---

## Manifests

Under `examples/khr/runtime-sandbox/`:

| File | Purpose |
|------|---------|
| `namespace.yaml` | Sandbox namespace + labels |
| `configmap-karl-host-runtime.yaml` | Config with apply **disabled** |
| `karl-host-runtime-config-apply.yaml` | Apply test only (`sandboxApplyEnabled: true`) |
| `test-target-linux.yaml` | Pause pod Linux target |
| `resourceport-*.json` / `resourcelease-dry-run.json` | Local dry-run fixtures |

---

## Scripts

| Script | Step |
|--------|------|
| `scripts/khr_runtime_sandbox_preflight.sh` | Context, namespace, labels, apply manifests |
| `scripts/khr_runtime_sandbox_dry_run.sh` | ResourceLease dry-run + ResourcePort candidate |
| `scripts/khr_runtime_sandbox_guarded_apply.sh` | Prove blocked by default; apply with explicit config |
| `scripts/khr_runtime_sandbox_rollback.sh` | Baseline rollback |
| `scripts/khr_runtime_sandbox_collect_evidence.sh` | JSON/LOG artifacts + production proof |
| `scripts/khr_runtime_sandbox_execute.sh` | Full pipeline |

All scripts source `scripts/khr_runtime_sandbox_lib.sh` for cluster/namespace/label guards.

---

## Evidence

Artifacts: `docs/evidence/khr-runtime-sandbox/<runId>/`  
Summary: `docs/evidence/khr-runtime-sandbox/summary.json` (`status: PASS`)

---

## Run (operator)

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
chmod +x scripts/khr_runtime_sandbox_*.sh
KHR_RUNTIME_SANDBOX_LIVE=1 ./scripts/khr_runtime_sandbox_execute.sh
./scripts/validate_khr_runtime_sandbox.sh
```

Validate hook (requires committed evidence or live run):

```bash
KHR_RUNTIME_SANDBOX_LIVE=1 ./scripts/validate.sh
```

---

## ISO note

`karl-host-runtime.service` remains **disabled** on ISO. Sandbox validation runs from Hyperdensity checkout against the live cluster, not from provision flow.
