# KHR TP Live Scope-2 ResourcePort Loop Plan (KHR-AZ)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AZ |
| **Cluster** | `karl-metal-01@ovh` |
| **Namespace** | `khr-runtime-sandbox` only |
| **Mode** | **KHR-AZ:** preflight only — **KHR-BA:** bounded manual loop (`observed-json`) |

---

## Purpose

Prepare **live-readonly** readiness for the ResourcePort observation loop (`karl-host-runtime -mode=resourceport-loop`) without enabling it automatically. Scope-2 remains **operator-only** and **blocked until explicit sign-off** in a dedicated execution sprint.

**Non-goals (KHR-AZ):** permanent `resourcePortLoopEnabled=true`, ResourceLease dry-run/apply, production enable, autonomous orchestration, Dashboard mutating actions.

---

## Prerequisites (Scope-2)

| ID | Prerequisite | Evidence / check |
|----|--------------|------------------|
| P-01 | Scope-1 PASS | `docs/evidence/khr-tp-live-scope1/committed-scope1-khr-aw/verify-summary.json` — `readyForScope2=false`, `resourcePortLoopEnabled=false` |
| P-02 | Reference env scope-1 | `readyForScope1=true`, cluster-sandbox gateway preferred |
| P-03 | rdp-GW cluster-sandbox | `docs/evidence/khr-rdpgw-cluster-sandbox/committed-cluster-sandbox-khr-ay/verify-summary.json` — Scope-2 **visibility** depends on live-readonly gateway |
| P-04 | Enablement gates | `khr_tp_live_enablement_preflight.sh` — `readyForScope0=true` |
| P-05 | Sandbox namespace | `khr-runtime-sandbox` with `khr.karl.io/sandbox=true` |
| P-06 | CRD foundation | ResourcePort CRD present (`resourceports.runtime.karl.io`) |
| P-07 | Preview manifests | `examples/khr/tp-live-scope1/` — karl-host-runtime preview deployable |
| P-08 | Federation / provenance | Latest PASS federation + provenance summaries (read-only) |

---

## Mandatory guards (execution sprint — not KHR-AZ)

All guards must PASS before any loop run:

| Guard | Requirement |
|-------|-------------|
| G-CTX | `kubectl` current-context = `karl-metal-01@ovh` |
| G-NS | Target namespace `khr-runtime-sandbox` only; label `khr.karl.io/sandbox=true` |
| G-CFG | `resourcePortLoopEnabled=false` until operator sets opt-in config for execution sprint |
| G-APPLY | `sandboxApplyEnabled=false` — no ResourceLease apply |
| G-EMIT | Loop runs **observed-json** or **cr-preview** first; `apply-cr` only with `--i-understand-this-is-sandbox` |
| G-PROD | No mutation of `karl`, `default`, `kube-system`, `karl-system` |
| G-RDPGW | Production `karl/rdpgw` generation unchanged |
| G-DASH | Dashboard read-only; no loop/apply action buttons |

---

## Forbidden actions (KHR-AZ and until sign-off)

| ID | Forbidden |
|----|-----------|
| F-01 | Set `resourcePortLoopEnabled=true` in committed sandbox config (KHR-AZ) |
| F-02 | `karl-host-runtime -mode=resourceport-loop` with `-apply-cr=true` |
| F-03 | ResourceLease dry-run or apply |
| F-04 | Production namespace workload mutation |
| F-05 | Operator-less apply; `productionReady=true` claims |
| F-06 | ISO/systemd default enable of host-runtime |
| F-07 | Dashboard disconnect/revoke/apply controls |

---

## Dry-run / preflight only (KHR-AZ)

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
./scripts/khr_tp_live_scope2_preflight.sh
```

Output: `docs/evidence/khr-tp-live-scope2-preflight/<runId>/scope2-preflight-summary.json`

Expected:

| Field | KHR-AZ value |
|-------|----------------|
| `status` | `PASS` |
| `readyForScope1` | `true` |
| `readyForScope2` | `conditional/manual-preflight-pass` |
| `resourcePortLoopEnabled` | `false` |
| `sandboxApplyEnabled` | `false` |
| `resourceLeaseApplyEnabled` | `false` |
| `loopEnabled` | `false` |

`readyForScope2` is **never** boolean `true`. After KHR-BA manual loop PASS, use `readyForScope2=manual-loop-pass` (not “active”).

---

## Manual loop execution (KHR-BA)

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
export KHR_TP_LIVE_SCOPE2_I_UNDERSTAND_MANUAL_LOOP=true
export KHR_SCOPE2_LOOP_ITERATIONS=2          # 1-3 only
export KHR_SCOPE2_LOOP_TIMEOUT_SEC=120
export KHR_TP_LIVE_SCOPE2_LOOP_RUN_ID=committed-scope2-loop-khr-ba

./scripts/khr_tp_live_scope2_preflight.sh
./scripts/khr_tp_live_scope2_resourceport_loop_run.sh
./scripts/khr_tp_live_scope2_resourceport_loop_verify.sh
./scripts/khr_tp_live_scope2_resourceport_loop_cleanup.sh
```

| Guard | Value |
|-------|-------|
| Confirmation | `KHR_TP_LIVE_SCOPE2_I_UNDERSTAND_MANUAL_LOOP=true` |
| Emission | `observed-json` only (no `-emit-cr` / `-apply-cr`) |
| Iterations | `KHR_SCOPE2_LOOP_ITERATIONS` ∈ {1,2,3} |
| Timeout | `KHR_SCOPE2_LOOP_TIMEOUT_SEC` (default 120) |
| Persistent config | Cluster ConfigMap **must** keep `resourcePortLoopEnabled: false` |

Evidence: `docs/evidence/khr-tp-live-scope2-resourceport-loop/<runId>/` — `loop-summary.json`, `verify-summary.json`, `cleanup-summary.json`.

### Scope-2 does **not** imply Scope-3

| Scope | Status after KHR-BA |
|-------|---------------------|
| **scope-2** | `manual-loop-pass` — ResourcePort **observation** only |
| **scope-3** | **blocked** — no ResourceLease dry-run |
| **scope-4** | **blocked** — no guarded apply |

---

## Rollback / cleanup (execution sprint)

| Step | Action |
|------|--------|
| 1 | Stop loop process / scale preview deployment to zero |
| 2 | `karl-host-runtime -mode=resourceport-cleanup -namespace=khr-runtime-sandbox` (sandbox CRs only) |
| 3 | Remove preview deployment if applied: `khr_tp_live_scope1_rollback.sh` pattern |
| 4 | Verify `resourcePortLoopEnabled=false` in config |
| 5 | Confirm production namespaces unchanged |

---

## Required evidence (execution sprint — future)

| Artifact | Path pattern |
|----------|----------------|
| Preflight | `docs/evidence/khr-tp-live-scope2-preflight/<runId>/scope2-preflight-summary.json` |
| Loop observation | `docs/evidence/khr-resourceport-loop-live/<runId>/loop-summary.json` |
| CR preview | `docs/evidence/khr-resourceport-cr-preview/<runId>/` |
| Rollback | `scope2-rollback-summary.json` with `productionUntouched=true` |

KHR-AZ commits **preflight evidence only** (`committed-scope2-preflight-khr-az`).

---

## Related

- `KHR_TP_LIVE_ENABLEMENT_PLAN.md`
- `KHR_TP_LIVE_REFERENCE_ENVIRONMENT.md`
- `RESOURCEPORT_CR_PREVIEW.md`
- rdp-GW `RDPGW_REFERENCE_ENVIRONMENT.md` (cluster-sandbox Scope-1 dependency)
- Karl-OS-ISO `KHR_TP_LIVE_SCOPE1_SANDBOX.md` (Scope-2 manual boundary)
