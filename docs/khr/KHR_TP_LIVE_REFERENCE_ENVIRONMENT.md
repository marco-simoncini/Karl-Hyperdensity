# KHR TP Live Reference Environment (KHR-AX)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AX |
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
