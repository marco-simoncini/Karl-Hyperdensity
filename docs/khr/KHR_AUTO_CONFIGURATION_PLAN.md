# KHR Auto-Configuration Plan — KARL 2.0 Baremetal Reference (KHR-CK)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-CK / KHR-CL / **KHR-CM** |
| **Scope** | **CO:** bounded ResourcePort observed-json loop; ResourceLease dry-run pending |
| **Primary cluster** | `karl-metal-01@ovh` (baremetal reference) |
| **First auto-configured module** | **Hyperdensity** |
| **Production** | **NOT production ready** |

---

## Purpose

Define how **KARL 2.0** baremetal reference environments will be **auto-configured** in a future sprint sequence: ordered bootstrap, cross-repo contracts, and read-only evidence gates — **without** enabling global defaults, production claims, or public-cloud profile changes today.

---

## Explicit non-goals (KHR-CK)

| Non-goal | Status |
|----------|--------|
| Global auto-enable of KHR flags | **Forbidden** — reference env only |
| Production / GA claims | **Forbidden** |
| Public-cloud default profile change | **Forbidden** — `karl1-kubevirt-legacy` remains default |
| Runtime mutation in this sprint | **None** — documentation + guards only |
| Live rollout beyond CRD foundation | **None** — phases 2–6 remain plan-only |

---

## KARL 2.0 baremetal reference model

| Layer | Role in auto-configuration |
|-------|---------------------------|
| **Karl-OS-ISO** | Prepares CRD assets, profile manifest, boundary metadata — **does not** enable runtime or Dashboard flags globally |
| **Karl-Installer** | Profile `karl2-baremetal-khr-native` → `providerProfile=khr-native`, `targetEnvironment=baremetal-reference` (KHR-CL: plan/dry-run) |
| **Karl-Hyperdensity** | **First module** in bootstrap order — CRD foundation through governance evidence |
| **Karl-Dashboard** | Reference env flags + cockpit `shell-workload-list` as first migrated component |
| **Karl-Inventory** | Read-only auto ingest from snapshot / committed evidence |
| **rdp-GW** | `cluster-sandbox` gateway as reference module for federation visibility |

---

## Hyperdensity bootstrap order (normative)

Hyperdensity is the **first** auto-configured module. Later repos consume its evidence; nothing in CK applies cluster changes.

```text
Phase 1  CRD foundation          → Installer karl2 CRD apply + contract manifest verify
Phase 2  Host-runtime preview    → Sandbox NS, preview manifests, host-runtime disabled by default
Phase 3  ResourcePort loop       → Scope-2 preflight + bounded manual loop evidence (sandbox only)
Phase 4  ResourceLease dry-run   → Scope-3 dry-run evidence (no production apply)
Phase 5  Guarded apply policy    → Scope-4 certification + guarded-apply evidence (policy only)
Phase 6  Governance              → Scope-4 governance bundle + snapshot v1 aggregation
```

### Phase detail

| Phase | Hyperdensity anchor | Evidence / script (existing) | Auto-config gate (future) |
|-------|---------------------|------------------------------|---------------------------|
| **1 — CRD foundation** | `KHR_BAREMETAL_INSTALLER_PROFILE_EXPECTATIONS.md` | `khr_baremetal_khr_native_crd_foundation_evidence.sh` → `committed-khr-cm-v1` | `phase=crd-foundation`, `runtimeMutation=false`, `hostRuntimeEnabled=false`, `crdDiffEmpty=true` |
| **2 — Host-runtime preview** | `KHR_BAREMETAL_INSTALLER_PROFILE_EXPECTATIONS.md` | `khr_baremetal_khr_native_host_runtime_preview_evidence.sh` → `committed-khr-cn-v1` | `phase=host-runtime-preview`, `hostRuntimePreview=true`, `resourcePortLoopEnabled=false`, `resourceLeaseEnabled=false`, `systemdEnable=false` |
| **3 — ResourcePort loop** | `KHR_BAREMETAL_INSTALLER_PROFILE_EXPECTATIONS.md` | `khr_baremetal_khr_native_resourceport_loop_evidence.sh` → `committed-khr-co-v1` | `phase=resourceport-loop`, `emissionMode=observed-json`, `resourcePortLoopObserved=true`, `persistentLoopEnabled=false` |
| **4 — ResourceLease dry-run** | `KHR_BAREMETAL_INSTALLER_PROFILE_EXPECTATIONS.md` | `khr_baremetal_khr_native_resourcelease_dryrun_evidence.sh` → `committed-khr-cp-v1` | `phase=resourcelease-dryrun`, `applyAllowed=false`, `mutationObserved=false`, consumes CO observed-json |
| **4b — Guarded-apply preflight** | `KHR_BAREMETAL_INSTALLER_PROFILE_EXPECTATIONS.md` | `khr_baremetal_khr_native_guarded_apply_preflight_evidence.sh` → `committed-khr-cq-v1` | `phase=guarded-apply-preflight`, `applyExecuted=false`, consumes CP evidence |
| **4c — TP dry-run (reference)** | `KHR_TP_LIVE_SCOPE3_RESOURCELEASE_DRYRUN_PLAN.md` | `committed-scope3-dryrun-khr-bc` | plan-only TP anchor |
| **5 — Guarded apply (live)** | `KHR_BAREMETAL_INSTALLER_PROFILE_EXPECTATIONS.md` | `khr_baremetal_khr_native_guarded_apply_evidence.sh` → `committed-khr-cr-v1` | `phase=guarded-apply`, sandbox-only, single-target, rollback verified |
| **5a — Guarded apply repeatability** | `KHR_BAREMETAL_INSTALLER_PROFILE_EXPECTATIONS.md` | `khr_baremetal_khr_native_guarded_apply_repeatability_evidence.sh` → `committed-khr-cs-v1` | two cycles, rollback-proof baseline fields, negative-path checks; **no fleet / no dashboard rollout** |
| **5c — Audit snapshot + beta gate** | `KHR_BAREMETAL_INSTALLER_PROFILE_EXPECTATIONS.md`, `KHR_BETA_READINESS_PLAN.md` | `khr_baremetal_khr_native_audit_snapshot_evidence.sh` → `committed-khr-ct-v1` | aggregates CO–CS; `betaRuntimeReady`; **no apply** |
| **5b — Scope-4 certification (TP)** | `KHR_TP_LIVE_SCOPE4_GUARDED_APPLY_PLAN.md` | `committed-scope4-certification-khr-bf` | TP reference anchor |
| **6 — Governance** | `KHR_SCOPE4_OPERATIONAL_GOVERNANCE.md` | `committed-scope4-governance-khr-bg`, `committed-khr-bt-v1` snapshot | `scope4Active=false` in snapshot |

Each phase is **read-only observable** before the next phase may be marked auto-ready in a future sprint.

---

## Cross-repo auto-configuration map

| Repo | CK deliverable | Consumes Hyperdensity phase |
|------|----------------|----------------------------|
| Karl-Installer | `INSTALLER_KARL2_BAREMETAL_KHR_NATIVE_PROFILE.md` | Phase 1 |
| Karl-OS-ISO | `ISO_KARL2_AUTO_CONFIGURATION_BOUNDARY.md` | Phase 1 assets only |
| Karl-Dashboard | `DASHBOARD_KARL2_AUTO_CONFIGURATION.md` | Phases 1–6 via projection + cockpit mount |
| Karl-Inventory | `INVENTORY_KARL2_AUTO_INGEST.md` | Phase 6 snapshot + federation |
| rdp-GW | `RDPGW_KARL2_REFERENCE_CONFIGURATION.md` | Phase 3 visibility (cluster-sandbox) |

---

## Dashboard first migrated component

| Component | Status (reference) | CK note |
|-----------|-------------------|---------|
| `shell-workload-list` | **Mounted live** (`committed-khr-cj-v1`) | First cockpit component in KARL 2.0 auto-config plan |
| Other cockpit surfaces | Preview or legacy | Remain pending — no global mount |

See Karl-Dashboard `DASHBOARD_COCKPIT_SHELL_WORKLOAD_LIST_MOUNT_PLAN.md` and `DASHBOARD_COCKPIT_SHELL_WORKLOAD_LIST_LIVE_EVIDENCE.md`.

---

## Reference env flags (not global defaults)

Future auto-configuration on baremetal reference **may** set (operator/reference deployment only):

| Variable | Reference value | Global default |
|----------|-----------------|----------------|
| `HYPERDENSITY_KHR_BACKEND_PROJECTION_ENABLED` | `true` | `false` |
| `HYPERDENSITY_KHR_UI_PROJECTION_ENABLED` | `true` | `false` |
| `HYPERDENSITY_KHR_UI_COMPONENT_PREVIEW_ENABLED` | `true` | `false` |
| `HYPERDENSITY_KHR_COCKPIT_SHELL_LIST_MOUNT_ENABLED` | `true` | `false` |
| `HYPERDENSITY_KHR_TP_REFERENCE_ENV` | `true` | `false` |
| `HYPERDENSITY_KHR_PROVIDER_PROFILE` | `khr-native` | unset → `public-cloud-kubevirt-compatibility` |
| `KARL_INSTALLER_PROFILE` | `karl2-baremetal-khr-native` (plan) | unset → `karl1-kubevirt-legacy` |

---

## Validation (KHR-CK)

```bash
./scripts/validate_khr_auto_configuration_plan.sh
./scripts/validate.sh   # includes doc scope guard
```

Fixture: `examples/khr/karl2-baremetal-auto-configuration-plan.json`

---

## Related

- `KHR_BOOTSTRAP_CONSUMER_EXPECTATIONS.md`
- `KHR_INSTALLER_PROFILE_EXPECTATIONS.md`
- `KHR_TP_REFERENCE_SNAPSHOT_V1.md`
- `KHR_BETA_READINESS_PLAN.md`
- Cross-repo CK docs (Installer, ISO, Dashboard, Inventory, rdp-GW)
