# KHR TP Live Scope-4 Guarded Apply Plan (KHR-BD)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-BD |
| **Cluster** | `karl-metal-01@ovh` |
| **Namespace** | `khr-runtime-sandbox` only |
| **Mode** | **Plan + read-only preflight** — **no guarded apply execution** in KHR-BD |

---

## Purpose

Prepare **live TP ResourceLease guarded apply** readiness (`karl-host-runtime -mode=resourcelease-guarded-apply`) without executing apply, cgroup mutation, or production changes. Scope-4 remains **operator-only** and **not active** until a dedicated execution sprint with explicit sign-off.

**Non-goals (KHR-BD):** live guarded apply, cgroup writes, autonomous scheduler, production enable, Dashboard mutating actions, ISO/systemd enable.

---

## Prerequisites

| ID | Prerequisite | Evidence |
|----|--------------|----------|
| P-01 | Scope-1 PASS | `docs/evidence/khr-tp-live-scope1/committed-scope1-khr-aw/verify-summary.json` |
| P-02 | Scope-2 manual loop | `docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/verify-summary.json` — `readyForScope2=manual-loop-pass` |
| P-03 | Scope-3 manual dry-run | `docs/evidence/khr-tp-live-scope3-dryrun/committed-scope3-dryrun-khr-bc/verify-summary.json` — `readyForScope3=manual-dryrun-pass` |
| P-04 | Dry-run decision allowed | `dryrun-output.json` — `dryRunDecision=allowed` |
| P-05 | Rollback plan ref | `dryrun-output.json` — `rollbackPlanRef` present |
| P-06 | Verification plan ref | `dryrun-output.json` — `verificationPlanRef` present |
| P-07 | Source ResourcePort ref | `dryrun-output.json` — `sourceResourcePortRef` present |
| P-08 | Native-live certified | `docs/evidence/khr-native-live-lane/certification-summary.json` — `status=certified` |
| P-09 | Provenance valid | `docs/evidence/khr-provenance/summary.json` — `readOnly=true`, `noAutonomousOrchestration=true` |
| P-10 | Operator confirmation | `KHR_TP_LIVE_SCOPE4_I_UNDERSTAND_GUARDED_APPLY=true` (execution sprint only) |
| P-11 | Enablement gates | `readyForScope0=true`, federation PASS |

---

## Mandatory guards (execution sprint — not KHR-BD)

| Guard | Requirement |
|-------|-------------|
| G-CTX | `kubectl` current-context = `karl-metal-01@ovh` |
| G-NS | `khr-runtime-sandbox` only |
| G-LABEL | Lease + namespace label `khr.karl.io/sandbox=true` |
| G-CONFIRM | `KHR_TP_LIVE_SCOPE4_I_UNDERSTAND_GUARDED_APPLY=true` |
| G-APPLY | `-mode=resourcelease-guarded-apply` **and** `-apply-resourcelease=true` **and** `-i-understand-this-is-sandbox` |
| G-DRY | Internal dry-run **allowed** before apply |
| G-ROLLBACK | `rollbackPlanRef` + captured baseline under `--sandbox-dir` |
| G-VERIFY | `verificationPlanRef` + post-apply verification |
| G-PROD | No mutation of production namespaces |
| G-DASH | Dashboard read-only; no apply buttons |

---

## Forbidden actions (KHR-BD and until sign-off)

| ID | Forbidden |
|----|-----------|
| F-01 | Live `resourcelease-guarded-apply` execution (KHR-BD) |
| F-02 | Cgroup mutation / envelope apply |
| F-03 | Production namespace targets (`karl`, `karl-system`, `default`, `kube-system`) |
| F-04 | Autonomous apply scheduler or background reconcile |
| F-05 | Persistent enablement of apply flags in committed sandbox config |
| F-06 | Apply without rollback plan |
| F-07 | Apply without verification plan |
| F-08 | Windows guarded apply |
| F-09 | KubeVirt template mutation |
| F-10 | ISO/systemd default enable of host-runtime |
| F-11 | Dashboard mutating controls |

---

## Preflight only (KHR-BD)

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
./scripts/khr_tp_live_scope4_preflight.sh
```

Output: `docs/evidence/khr-tp-live-scope4-preflight/<runId>/scope4-preflight-summary.json`

Expected:

| Field | KHR-BD value |
|-------|----------------|
| `status` | `PASS` |
| `readyForScope3` | `manual-dryrun-pass` |
| `readyForScope4` | `conditional/manual-preflight-pass` |
| `readyForScope4Active` | `false` |
| `guardedApplyExecuted` | `false` |
| `cgroupMutationObserved` | `false` |

`readyForScope4` is **never** boolean `true` or `active` until a dedicated guarded-apply execution sprint completes with sign-off.

---

## Scope-4 does **not** imply production enable

| Scope | After KHR-BD |
|-------|----------------|
| **scope-3** | `manual-dryrun-pass` — dry-run evidenced, not active |
| **scope-4** | `conditional/manual-preflight-pass` — apply **not** executed |

---

## Rollback / verification (execution sprint)

- Rollback: `resourcelease-rollback` mode with baseline captured pre-apply
- Verification: read-back cgroup limits, continuity checks, `noRestart` policy
- rdp-GW: no session revoke/disconnect from Hyperdensity apply path

---

## Related

- `KHR_TP_LIVE_SCOPE3_RESOURCELEASE_DRYRUN_PLAN.md`
- `RESOURCELEASE_GUARDED_APPLY_SANDBOX.md`
- `scripts/khr_tp_live_scope4_preflight.sh`
