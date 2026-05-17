# KHR TP Live Reference Environment (KHR-AX)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AX / **KHR-BN** |
| **Cluster** | `karl-metal-01@ovh` |
| **Scopes** | **0 and 1 only** — scope-2+ **blocked** |
| **Mode** | Read-only observation — **not production** |

---

## Definition

The **reference environment** is the stabilized Technical Preview live posture on `karl-metal-01@ovh` where scope-0 federation and scope-1 sandbox evidence have PASS, without enabling ResourcePort loop, ResourceLease apply, or production namespaces.

| Property | Value |
|----------|-------|
| **Cluster context** | `karl-metal-01@ovh` |
| **Runtime namespace** | `khr-runtime-sandbox` (`khr.karl.io/sandbox=true`) |
| **Gateway namespace** | `khr-rdpgw-sandbox` (`khr.karl.io/sandbox=true`) |
| **contractSetId** | `khr-tp-contract-v1` |
| **productionReady** | `false` (always) |
| **noAutonomousOrchestration** | `true` (always) |

---

## Required evidence (stabilized)

| Artifact | Path | Requirement |
|----------|------|-------------|
| Enablement preflight | `docs/evidence/khr-tp-live-enablement/<runId>/enablement-preflight-summary.json` | `status=PASS`, `readyForScope1=true` |
| Scope-1 verify | `docs/evidence/khr-tp-live-scope1/committed-scope1-khr-aw/verify-summary.json` (or latest PASS) | `accessGraphLiveReadonly=true`, `readyForScope2=false` |
| Federation | `docs/evidence/khr-runtime-observation-federation/*/federation-summary.json` | `status=PASS` |
| rdp-GW continuity | `docs/evidence/khr-accessgraph-continuity/*/summary.json` | `source=live-readonly` preferred |
| rdp-GW cluster-sandbox | `docs/evidence/khr-rdpgw-cluster-sandbox/committed-cluster-sandbox-khr-ay/verify-summary.json` | **Preferred:** `deployMode=cluster-sandbox`, `accessGraphLiveReadonly=true` |
| Scope-2 preflight | `docs/evidence/khr-tp-live-scope2-preflight/committed-scope2-preflight-khr-az/scope2-preflight-summary.json` | `readyForScope2=conditional/manual-preflight-pass`, `resourcePortLoopEnabled=false` |
| Scope-2 manual loop | `docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/verify-summary.json` | `readyForScope2=manual-loop-pass`, `emissionMode=observed-json` |
| Scope-3 preflight | `docs/evidence/khr-tp-live-scope3-preflight/committed-scope3-preflight-khr-bb/scope3-preflight-summary.json` | preflight PASS; apply disabled |
| Scope-3 manual dry-run | `docs/evidence/khr-tp-live-scope3-dryrun/committed-scope3-dryrun-khr-bc/verify-summary.json` | `readyForScope3=manual-dryrun-pass`, `dryRunObserved=true`, `noMutation=true`, `noApply=true`, not active |
| Scope-4 preflight | `docs/evidence/khr-tp-live-scope4-preflight/committed-scope4-preflight-khr-bd/scope4-preflight-summary.json` | preflight PASS; apply disabled until execution sprint |
| Scope-4 guarded apply | `docs/evidence/khr-tp-live-scope4-guarded-apply/committed-scope4-guarded-apply-khr-be/verify-summary.json` | `readyForScope4=manual-guarded-apply-pass`, rollback verified, sandbox-only, not active |
| rdp-GW deploy mode (fallback) | `docs/evidence/khr-rdpgw-scope1/*/deploy-summary.json` | `deployMode=local-gateway` with warning if cluster-sandbox unavailable |

---

## Rollback policy

| Rule | Detail |
|------|--------|
| **Scope-1 rollback** | `khr_tp_live_scope1_rollback.sh` removes sandbox **deployments** only |
| **Namespaces** | Retained with `khr.karl.io/sandbox=true` labels |
| **Production** | `karl`, `default`, `kube-system`, `karl-system` — **never** targeted for KHR live enablement |
| **ISO/systemd** | `karl-host-runtime.service` remains **disabled** on ISO provision |

Re-run scope-1 deploy after rollback only with `KHR_TP_LIVE_SCOPE1_I_UNDERSTAND_SANDBOX=true`.

---

## Reference env check

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
./scripts/khr_tp_live_reference_env_check.sh
```

Output: `docs/evidence/khr-tp-live-reference-env/<runId>/reference-env-summary.json`

---

## Dashboard TP readiness (reference)

| Item | Value |
|------|-------|
| Endpoint | `GET /api/hyperdensity/tp-readiness` |
| Reference env flags | `HYPERDENSITY_KHR_TP_READINESS_ENABLED=true` + `HYPERDENSITY_KHR_TP_REFERENCE_ENV=true` |
| Fixture | `examples/khr-dashboard/tp-readiness-reference-env.json` |

See Karl-Dashboard `DASHBOARD_TP_READINESS_REFERENCE_ENV.md`.

---

## KHR-native reference activation (KHR-BN)

Reference env may activate **KHR-native** provider profile on Dashboard **without** changing global defaults:

```bash
export HYPERDENSITY_KHR_BACKEND_PROJECTION_ENABLED=true
export HYPERDENSITY_KHR_TP_REFERENCE_ENV=true
export HYPERDENSITY_KHR_PROVIDER_PROFILE=khr-native
```

| Profile | Reference env | Global default / legacy |
|---------|---------------|-------------------------|
| `khr-native` | Explicit activation on `karl-metal-01@ovh` | Future greenfield only |
| `public-cloud-kubevirt-compatibility` | Remains compatibility profile for legacy installs | Unset env → `default-documented` |

**No production claim:** `productionReady=false`, `autonomousOrchestration=false`, read-only projection only.

Fixture: Karl-Dashboard `reference-env-khr-native-activation.json`  
Doc: `DASHBOARD_REFERENCE_ENV_ACTIVATION_PROFILE.md`

### Dashboard activation evidence (KHR-BO / KHR-BP / **KHR-BQ**)

| Level | `evidenceStatus` | `source` | Requirement |
|-------|------------------|----------|-------------|
| **Live activation** | `LIVE_PASS` | `live-readonly` | KHR-enabled console image + activation env; `GET .../khr-backend/projection` HTTP 200, `providerProfile=khr-native` |
| **Fixture proof** | `PASS` | `fixture-readonly` | CI/offline contract validation |
| **Remediation** | `REMEDIATION_PASS` | `remediation-readonly` | Live port-forward + env audit; `remediation-plan.md` when env/route missing |
| **Rollout blocked** | `ROLLOUT_BLOCKED_IMAGE_MISSING` | `rollout-blocked-readonly` | No KHR-enabled image — **no partial env patch** (KHR-BQ) |
| **Rollback ready** | `ROLLBACK_READY` | plan/rollback JSON | Pre-rollout snapshot in `rollout-plan.json` / `rollback-plan.json` |
| **Rollback verified** | `ROLLBACK_VERIFIED` | `rollback-summary.json` | Image `1.6.0` restored, activation env cleared, route inactive |

**Committed LIVE_PASS (KHR-BS):** Karl-Dashboard `docs/evidence/khr-dashboard-reference-env-activation/committed-khr-bs-20260517T073046Z/`.

### KHR TP Reference Snapshot v1 (KHR-BT)

Cross-repo read-only index: `KHR_TP_REFERENCE_SNAPSHOT_V1.md` → `docs/evidence/khr-tp-reference-snapshot-v1/committed-khr-bt-v1/snapshot-summary.json`.

| Artifact | Path |
|----------|------|
| Rollout plan | `.../rollout-plan.json` (KHR-BQ) |
| Live summary | `.../live-summary.json` (KHR-BQ verify or LIVE_PASS live evidence) |
| Rollback plan | `.../rollback-plan.json` (KHR-BQ) |
| Summary | `.../summary.json` (KHR-BP live evidence) |
| Live connectivity | `.../live-connectivity.json` |
| Deployment env audit | `.../deployment-env-audit.json` |
| Remediation plan | `.../remediation-plan.md` (when required) |

**Console image resolve (KHR-BR):**

| `evidenceStatus` | Meaning |
|------------------|---------|
| `IMAGE_RESOLVED` | Pullable KHR-enabled console image identified |
| `IMAGE_UNRESOLVED_BUILD_READY` | Sources buildable; image not published yet |
| `IMAGE_BUILD_SOURCE_MISSING` | `kubernetes-console` tree incomplete |
| `IMAGE_PULL_BLOCKED` | Image ref set but not pullable |
| `LIVE_PASS` | After rollout + live projection verify |

Script: `khr_dashboard_khr_enabled_image_resolve.sh` → `docs/evidence/khr-dashboard-console-image/<runId>/image-resolution-summary.json`

**Controlled rollout (KHR-BQ):**

```bash
export KHR_DASHBOARD_REFERENCE_ENV_ROLLOUT_I_UNDERSTAND=true
export KHR_DASHBOARD_KHR_ENABLED_CONSOLE_IMAGE=<khr-console-image>  # required for apply
cd ../Karl-Dashboard && ./scripts/khr_dashboard_reference_env_rollout_plan.sh
export KHR_DASHBOARD_ROLLOUT_RUN_ID=<runId>
./scripts/khr_dashboard_reference_env_rollout_verify.sh
./scripts/khr_dashboard_reference_env_rollout_rollback.sh  # after evidence
```

**Live evidence (KHR-BP):**

```bash
export KHR_DASHBOARD_REFERENCE_ENV_LIVE_I_UNDERSTAND=true
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
cd ../Karl-Dashboard && ./scripts/khr_dashboard_reference_env_live_evidence.sh
```

Uses `kubectl port-forward svc/karl-console-next-oidc -n karl-system` when `DASHBOARD_BASE_URL` is unreachable. Rollout scripts patch **only** approved reference Deployment; live evidence does **not** auto-patch.

See Karl-Dashboard `DASHBOARD_KHR_ENABLED_CONSOLE_ROLLOUT.md` and `DASHBOARD_REFERENCE_ENV_ACTIVATION_EVIDENCE.md`.

---

## Baremetal KHR-native standing reference profile (KHR-CC / KHR-CD)

| Dimension | Baremetal reference (`karl-metal-01@ovh`) | Public cloud / legacy installs |
|-----------|---------------------------------------------|--------------------------------|
| **Primary target** | `khr-native` (Shell/Cell/Lease-first) | `public-cloud-kubevirt-compatibility` |
| **KubeVirt** | Optional compatibility provider | Required compatibility binding for VM fleets |
| **Standing profile state** | `rollback-verified` (CD apply + soak + mandatory rollback) | N/A |
| **Rollback baseline** | `console:1.6.0` (KHR-CD proven) | Unchanged installer default |
| **globalDefaultsChanged** | `false` | `false` |

**Temporary rollout evidence (done):** KHR-BS backend projection `LIVE_PASS` + KHR-CB UI preview `LIVE_PASS`, both with mandatory rollback.

**Standing profile (CD):** Karl-Dashboard `docs/evidence/khr-dashboard-baremetal-standing-profile/committed-khr-cd-v1/` — controlled apply on reference console, 300s soak (`applied-soak-verified`), mandatory rollback (`rollback-verified`). Anchors CC plan `committed-khr-cc-v1`.

**Global default:** ISO and unset Dashboard env still resolve to `public-cloud-kubevirt-compatibility`. Standing profile is **not** a global switch.

---

## Forbidden (all reference env operations)

- scope-2 ResourcePort loop enable
- scope-3/4 ResourceLease dry-run/apply
- Production namespace mutation
- Autonomous orchestration
- Dashboard mutating actions / action buttons
- ISO default systemd enable

---

## Related

- `KHR_TP_LIVE_ENABLEMENT_PLAN.md`
- `KHR_TP_LIVE_SCOPE1_SANDBOX.md` (Hyperdensity + ISO)
- rdp-GW `RDPGW_REFERENCE_ENVIRONMENT.md`
