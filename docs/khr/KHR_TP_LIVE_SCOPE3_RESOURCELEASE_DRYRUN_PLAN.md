# KHR TP Live Scope-3 ResourceLease Dry-Run Plan (KHR-BB)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-BB |
| **Cluster** | `karl-metal-01@ovh` |
| **Namespace** | `khr-runtime-sandbox` only |
| **Mode** | **Plan + read-only preflight** — **no live dry-run execution** in KHR-BB |

---

## Purpose

Prepare **live TP ResourceLease dry-run** readiness (`karl-host-runtime -mode=resourcelease-dryrun`) without executing dry-run or apply. Scope-3 remains **operator-only** and **blocked until explicit sign-off** in a dedicated execution sprint (KHR-BC or later).

**Non-goals (KHR-BB):** live ResourceLease dry-run, ResourceLease apply, cgroup mutation, autonomous scheduler, production enable, Dashboard mutating actions.

---

## Prerequisites

| ID | Prerequisite | Evidence |
|----|--------------|----------|
| P-01 | Scope-1 PASS | `docs/evidence/khr-tp-live-scope1/committed-scope1-khr-aw/verify-summary.json` |
| P-02 | Scope-2 manual loop | `docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/verify-summary.json` — `readyForScope2=manual-loop-pass` |
| P-03 | ResourcePort observation | Scope-2 `loop-summary.json` — `emissionMode=observed-json` |
| P-04 | Sample ResourceLease | `examples/khr/runtime-sandbox/resourcelease-dryrun-allowed.json` |
| P-05 | Rollback plan declared | Lease `governance.rollbackPlanRef` present (not executed in KHR-BB) |
| P-06 | Dry-run mode available | `karl-host-runtime -mode=resourcelease-dryrun` in tree (not invoked live) |
| P-07 | Sandbox guards | `khr.karl.io/sandbox=true`, `sandboxApplyEnabled=false` |
| P-08 | Enablement gates | `readyForScope0=true`, federation PASS |

---

## Mandatory guards (execution sprint — not KHR-BB)

| Guard | Requirement |
|-------|-------------|
| G-CTX | `kubectl` current-context = `karl-metal-01@ovh` |
| G-NS | `khr-runtime-sandbox` only |
| G-LABEL | Lease + namespace label `khr.karl.io/sandbox=true` |
| G-APPLY | `sandboxApplyEnabled=false`; no `-mode=resourcelease-guarded-apply` |
| G-DRY | Dry-run only; `noApply=true` in output |
| G-PROD | No mutation of production namespaces |
| G-NO-CGROUP | No cgroup writes from KHR paths in this scope |
| G-DASH | Dashboard read-only; no dry-run/apply buttons |

---

## Forbidden actions (KHR-BB and until sign-off)

| ID | Forbidden |
|----|-----------|
| F-01 | Live `resourcelease-dryrun` execution (KHR-BB) |
| F-02 | `resourcelease-guarded-apply` or any ResourceLease apply |
| F-03 | Cgroup mutation / envelope apply |
| F-04 | Autonomous dry-run scheduler or background reconcile |
| F-05 | Persistent enablement of apply flags in committed sandbox config |
| F-06 | Production namespace targets |
| F-07 | ISO/systemd default enable of host-runtime |
| F-08 | Dashboard mutating controls |

---

## Preflight only (KHR-BB)

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
./scripts/khr_tp_live_scope3_preflight.sh
```

Output: `docs/evidence/khr-tp-live-scope3-preflight/<runId>/scope3-preflight-summary.json`

Expected:

| Field | KHR-BB value |
|-------|----------------|
| `status` | `PASS` |
| `readyForScope2` | `manual-loop-pass` |
| `readyForScope3` | `conditional/manual-preflight-pass` |
| `readyForScope3Active` | `false` |
| `readyForScope4` | `false` |
| `resourceLeaseDryRunExecuted` | `false` |
| `resourceLeaseApplyEnabled` | `false` |

`readyForScope3` is **never** boolean `true` or `active` until a dedicated dry-run execution sprint completes with sign-off.

---

## Scope-3 does **not** imply Scope-4

| Scope | After KHR-BB |
|-------|----------------|
| **scope-3** | `conditional/manual-preflight-pass` — dry-run **not** executed |
| **scope-4** | **blocked** — no guarded apply sandbox |

---

## Rollback plan (execution sprint)

Rollback artifacts are **required in lease input** but **not executed** in KHR-BB. Execution sprint must run `resourcelease-rollback` only after explicit approval and evidence trail.

---

## Related

- `KHR_TP_LIVE_SCOPE2_RESOURCEPORT_LOOP_PLAN.md`
- `RESOURCELEASE_DRYRUN_AGAINST_RESOURCEPORT.md`
- Karl-Inventory `RESOURCELEASE_DRYRUN_OBSERVATION_PREP.md`
- Karl-OS-ISO `KHR_TP_LIVE_SCOPE1_SANDBOX.md` (Scope-3 boundary)
