# KHR validation modes (KHR-BU)

Stabilize Technical Preview validation after **Reference Snapshot v1** (`committed-khr-bt-v1`). Separate **offline** checks (default) from **live cluster** checks (opt-in).

| Mode | Env | Cluster | Use |
|------|-----|---------|-----|
| **Offline (default)** | *(unset)* | Not required | CI, local dev, beta readiness gate |
| **Live** | `KHR_LIVE_VALIDATE=1` | `karl-metal-01@ovh` (or configured context) | Operator re-validation, fresh evidence |

---

## Offline validation (default)

`./scripts/validate.sh` runs without `KHR_LIVE_VALIDATE`.

| Layer | Script / artifact |
|-------|-------------------|
| Go tests, schemas, contracts | `validate.sh` (existing) |
| Reference snapshot v1 | `scripts/khr_validate_reference_snapshot.sh` |
| Committed scope evidence | `scripts/khr_validate_committed_evidence.sh` |
| Snapshot JSON | `docs/evidence/khr-tp-reference-snapshot-v1/committed-khr-bt-v1/snapshot-summary.json` |

**Scope-3 offline rule:** `dryrun-summary.json` (`status=PASS`) is authoritative. `verify-summary.json` may show `FAIL` after a stale live re-run against current cluster state — offline validate does **not** re-execute `khr_tp_live_scope3_dryrun_verify.sh`.

**No mutation:** offline mode does not run enablement preflight, ResourcePort loop, scope dry-run verify, or reference-env cluster checks.

---

## Live validation (opt-in)

```bash
export KHR_LIVE_VALIDATE=1
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
./scripts/validate.sh
```

Runs kubectl-backed scripts:

- `khr_tp_live_enablement_preflight.sh`
- Scope 2/3/4 preflight, loop, dry-run verify, reference-env check

Use only when intentionally refreshing live evidence. Live failures do not block offline beta readiness if committed snapshot remains valid.

---

## Reference snapshot validation

```bash
./scripts/khr_validate_reference_snapshot.sh
# or
KHR_TP_REFERENCE_SNAPSHOT_RUN_ID=committed-khr-bt-v1 ./scripts/validate_khr_tp_reference_snapshot_v1.sh
```

Re-aggregates snapshot (read-only) and validates cross-repo evidence index + committed PASS summaries.

---

## Related

- `KHR_TP_REFERENCE_SNAPSHOT_V1.md`
- `KHR_SNAPSHOT_V1_FREEZE_POLICY.md`
- `KHR_BETA_READINESS_PLAN.md`
