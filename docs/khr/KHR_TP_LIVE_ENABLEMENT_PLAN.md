# KHR Technical Preview Live Enablement Plan (KHR-AV)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AV |
| **Cluster reference** | `karl-metal-01@ovh` |
| **Contract set** | `khr-tp-contract-v1` |
| **Status** | **Plan only** — no automatic enablement in KHR-AV |

---

## Purpose

Define **exactly** what may be turned on for KHR Technical Preview on the reference cluster, under which readiness gates, and which actions remain **forbidden**. This document is the executable blueprint for **KHR-TP-Live**; KHR-AV delivers the plan and read-only preflight only.

**Explicit non-goals (KHR-AV):** automatic enablement, production enable, autonomous orchestration, default systemd enable, destructive mutation, Dashboard mutating actions.

---

## Readiness gates (all scopes)

Every live scope requires these gates to **PASS** (verified by `scripts/khr_tp_live_enablement_preflight.sh`):

| Gate | Verification |
|------|----------------|
| **G1 CRD foundation** | 8 KHR CRDs on cluster **or** Installer `khr-installer-crd-foundation` evidence with `contractSetId` + verify |
| **G2 contractSetId** | `khr-tp-contract-v1` in contract manifest + evidence summaries |
| **G3 post-install** | `khr-tp-post-install-bundle` / ISO post-install verify **PASS** |
| **G4 native-live certification** | `docs/evidence/khr-native-live-lane/certification-summary.json` — `status=certified`, `readOnly=true`, `regressionDetected=false` |
| **G5 provenance** | `docs/evidence/khr-provenance/summary.json` — `readOnly=true`, `noAutonomousOrchestration=true` |
| **G6 federation** | Latest `federation-summary.json` — `status=PASS`, correlation match |
| **G7 rollback** | `scripts/khr_runtime_sandbox_rollback.sh` executable |
| **G8 no production namespace** | Enablement targets **sandbox only** (`khr-runtime-sandbox`, `khr-rdpgw-sandbox`) — never default `karl` production workloads |

---

## Live scopes

| Scope | Name | What may be enabled | KHR-AV preflight |
|-------|------|---------------------|------------------|
| **scope-0** | Read-only federation | Aggregate observation only; rdp-GW live-readonly evidence; federation check | `readyForScope0=true` when G1–G8 + federation PASS |
| **scope-1** | Runtime sandbox deploy | Manual deploy to `khr-runtime-sandbox` / rdpgw sandbox; **operator-only** | `readyForScope1=conditional` — gates PASS but **manual sprint required** |
| **scope-2** | ResourcePort loop | Host-runtime ResourcePort loop in sandbox | **blocked** until scope-1 deploy sprint |
| **scope-3** | ResourceLease dry-run | Dry-run / preview only in sandbox | **blocked** until scope-2 |
| **scope-4** | Guarded apply sandbox | Guarded apply with rollback evidence | **blocked** until scope-3 + explicit approval sprint |

### scope-0 detail

- Run `khr_runtime_observation_federation_check.sh`
- rdp-GW `live-readonly` or `fixture-readonly` evidence (live preferred)
- No cluster mutation

### scope-1 detail (next sprint: KHR-TP-Live runtime)

- Prerequisites: scope-0 + sandbox manifests + rollback rehearsed
- Manual: `kubectl apply` sandbox namespaces only
- systemd **remains disabled** on ISO hosts unless operator explicitly enables (out of band)
- rdpgw: `examples/khr/rdpgw-sandbox/` only

### scope-2–4

Deferred to dedicated sprints after scope-1 deploy evidence exists. Preflight reports `readyForScope2+ = false` with `blockedReason`.

---

## Forbidden actions (all scopes)

| ID | Forbidden |
|----|-----------|
| F-01 | Target **production** namespace (`karl` default workloads) for KHR enablement |
| F-02 | **Autonomous apply** or orchestration without operator evidence trail |
| F-03 | **systemd default enable** of `karl-host-runtime` on ISO provision |
| F-04 | Dashboard **mutating** actions (disconnect, revoke, apply buttons) |
| F-05 | `productionReady: true` or GA claims |
| F-06 | Auth enforcement / session mutation via KHR paths |
| F-07 | Destructive cluster mutation without rollback artifact |

---

## Preflight (read-only)

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
./scripts/khr_tp_live_enablement_preflight.sh
```

Output: `docs/evidence/khr-tp-live-enablement/<runId>/enablement-preflight-summary.json`

Expected KHR-AV reference run:

| Field | Expected |
|-------|----------|
| `readyForScope0` | `true` |
| `readyForScope1` | `conditional` |
| `readyForScope2` | `false` |
| `readyForScope3` | `false` |
| `readyForScope4` | `false` |
| `automaticEnablement` | `false` |

---

## Scope-4 certification (KHR-BF)

After KHR-BE guarded-apply evidence:

```bash
./scripts/khr_scope4_certification_check.sh
```

Output: `docs/evidence/khr-scope4-guarded-apply-certification/committed-scope4-certification-khr-bf/certification-summary.json`

| Field | Value |
|-------|----------|
| `scope4CertificationState` | `certified-evidence-backed` |
| `readyForScope4` | `manual-guarded-apply-pass` |
| `readyForScope4Active` | `false` |
| `guardedApplyEnabled` | `false` |
| `guardedApplyAutonomous` | `false` |

Failure semantics (read-only simulate): `SCOPE4_FAILURE_SEMANTICS.md` — **no live failure injection**.

---

## Related

| Repo | Document |
|------|----------|
| Karl-OS-ISO | `KHR_TECHNICAL_PREVIEW_PROFILE.md`, `KHR_POST_INSTALL_VERIFY.md` |
| Karl-Installer | `KHR_INSTALLER_CRD_FOUNDATION_EVIDENCE.md` |
| Karl-Dashboard | `TECHNICAL_PREVIEW_DASHBOARD_GUIDE.md` |
| rdp-GW | `RDPGW_SANDBOX_LIVE_EVIDENCE.md` |
| Hyperdensity | `RUNTIME_OBSERVATION_FEDERATION.md`, `TECHNICAL_PREVIEW_OPERATOR_RUNBOOK.md` |
