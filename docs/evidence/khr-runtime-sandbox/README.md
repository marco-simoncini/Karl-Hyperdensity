# KHR runtime sandbox evidence (KHR-G)

| Field | Value |
|-------|-------|
| **Cluster** | `karl-metal-01@ovh` |
| **Namespace** | `khr-runtime-sandbox` |

## Latest summary

See [`summary.json`](summary.json) for the latest PASS run (`status`, `noProductionMutation`, artifact list).

Per-run artifacts live under `<runId>/` (e.g. `preflight.json`, `dry-run-lease.json`, `guarded-apply-sandbox.json`, `rollback-baseline.json`).

## Regenerate

```bash
KHR_RUNTIME_SANDBOX_LIVE=1 ./scripts/khr_runtime_sandbox_execute.sh
```
