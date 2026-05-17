# KHR Baremetal Installer Profile Expectations (KHR-CL / KHR-CM)

| Field | Value |
|-------|-------|
| **Profile ID** | `karl2-baremetal-khr-native` |
| **Sprint** | KHR-CL … KHR-CQ / **KHR-CR** / **KHR-CS** / **KHR-CT** / **KHR-CU** / **KHR-CV** / **KHR-CW** / **KHR-CX** / **KHR-CY** / **KHR-DC** / **KHR-DD** / **KHR-DE** / **KHR-DF** / **KHR-DG** / **KHR-DH** / **KHR-DI** / **KHR-DJ** / **KHR-DK** / **KHR-DL** / **KHR-DM** / **KHR-DP** / **KHR-DQ** / **KHR-DR** / **KHR-DS** / **KHR-DT** |
| **Mode** | Phased: plan, `crd-foundation`, `host-runtime-preview` on reference cluster |

---

## Expected installer outputs

| Output | Location | Invariants |
|--------|----------|------------|
| Auto-config plan JSON | `docs/evidence/khr-baremetal-khr-native-installer-profile/<runId>/auto-config-plan.json` | `runtimeMutation=false`, `planOnly=true` |
| Evidence summary | `.../summary.json` | `evidenceStatus=PLAN_PASS` |
| Hyperdensity phases | `hyperdensityPhases[]` (6 entries) | Read-only |
| Dashboard env plan | `dashboardEnvPlan.env` | Reference flags only |
| Inventory ingest plan | `inventoryIngestPlan` | Read-only sources |
| rdp-GW plan | `rdpGwPlan.deployMode=cluster-sandbox` | No production mutation |

---

## Mapping (normative)

| Field | Value |
|-------|-------|
| `installerProfile` | `karl2-baremetal-khr-native` |
| `dashboardProviderProfile` | `khr-native` |
| `targetEnvironment` | `baremetal-reference` |
| `kubevirtRequired` | `false` |
| `systemdEnable` | `false` |
| `globalDefaultsChanged` | `false` |

---

## Guards

| Env | Phase |
|-----|-------|
| `KARL_INSTALLER_BAREMETAL_KHR_NATIVE_I_UNDERSTAND=true` | All phases |
| `KARL_INSTALLER_BAREMETAL_KHR_NATIVE_APPLY_CRDS_I_UNDERSTAND=true` | `crd-foundation` apply only |
| Cluster context | `karl-metal-01@ovh` required for apply |

---

## Phase: host-runtime-preview (KHR-CN)

| Field | Value |
|-------|-------|
| `phase` | `host-runtime-preview` |
| `hostRuntimePreview` | `true` |
| `resourcePortLoopEnabled` | `false` |
| `resourceLeaseEnabled` | `false` |
| `systemdEnable` | `false` |
| `runtimeMutation` | `false` |
| Namespace | `khr-runtime-sandbox` (`khr.karl.io/sandbox=true`) |
| Evidence | `docs/evidence/khr-baremetal-khr-native-host-runtime-preview/committed-khr-cn-v1/` |

---

## Phase: resourceport-loop (KHR-CO)

| Field | Value |
|-------|-------|
| `phase` | `resourceport-loop` |
| `resourcePortLoopObserved` | `true` |
| `emissionMode` | `observed-json` |
| `resourcePortLoopEnabled` | `false` (persistent) |
| `resourceLeaseEnabled` | `false` |
| `emitCR` / `applyCR` | `false` |
| Evidence | `docs/evidence/khr-baremetal-khr-native-resourceport-loop/committed-khr-co-v1/` |

---

## Phase: resourcelease-dryrun (KHR-CP)

| Field | Value |
|-------|-------|
| `phase` | `resourcelease-dryrun` |
| `observedJsonSource` | KHR-CO `loop-output.json` |
| `applyAllowed` | `false` |
| `mutationObserved` | `false` |
| `persistentEnable` | `false` |
| `rollbackPlanPresent` | `true` (not executable) |
| `guardedApplyBlocked` | `true` (phase 5) |
| Evidence | Karl-Installer `docs/evidence/khr-baremetal-khr-native-resourcelease-dryrun/committed-khr-cp-v1/` |

---

## Phase: guarded-apply-preflight (KHR-CQ)

| Field | Value |
|-------|-------|
| `phase` | `guarded-apply-preflight` |
| `dryRunEvidenceRef` | KHR-CP `committed-khr-cp-v1/` |
| `applyAllowed` / `applyExecuted` | `false` |
| `runtimeMutation` | `false` |
| Evidence | Karl-Installer `docs/evidence/khr-baremetal-khr-native-guarded-apply-preflight/committed-khr-cq-v1/` |

---

## Phase: guarded-apply (KHR-CR)

| Field | Value |
|-------|-------|
| `phase` | `guarded-apply` |
| `singleTargetOnly` | `true` |
| `namespace` | `khr-runtime-sandbox` only |
| `applyAllowed` / `applyExecuted` | `true` |
| `runtimeMutation` | `true` (sandbox scope only) |
| `rollbackExecuted` / `rollbackVerified` | `true` |
| Evidence | Karl-Installer `docs/evidence/khr-baremetal-khr-native-guarded-apply/committed-khr-cr-v1/` |

---

## Phase: guarded-apply repeatability (KHR-CS)

| Field | Value |
|-------|-------|
| `phase` | `guarded-apply` (same as KHR-CR; **no new phase**) |
| `scope` | Repeatability + evidence hardening only — **not** scope expansion |
| `repeatabilityRuns` | `2` on same sandbox target |
| `rollbackProof` | `cpuMaxBefore`, `cpuMaxApplied`, `observedCpuMaxAfterRollback`, `expectedCpuMaxAfterRollback`, `baselineMatch` |
| `negativePath` | missing guard, wrong namespace, missing rollback plan, multi-target — all rejected |
| Evidence | Karl-Installer `docs/evidence/khr-baremetal-khr-native-guarded-apply-repeatability/committed-khr-cs-v1/` |

---

## Phase: audit snapshot + beta gate (KHR-CT)

| Field | Value |
|-------|-------|
| `phase` | audit aggregation only — **no new apply**, **no runtime mutation** |
| `betaRuntimeReady` | from `audit-snapshot.json` when CO–CS all PASS |
| `runtimeMutationScope` | `khr-runtime-sandbox` only |
| Evidence | Karl-Installer `docs/evidence/khr-baremetal-khr-native-audit-snapshot/committed-khr-ct-v1/` |

---

## Inventory ingest: audit snapshot (KHR-CU)

| Field | Value |
|-------|-------|
| `mode` | read-only observation ingest |
| `source` | Installer `committed-khr-ct-v1` |
| Evidence | Karl-Inventory `docs/evidence/khr-inventory-audit-snapshot-ingest/committed-khr-cu-v1/` |
| `enforcementEnabled` | `false` |

---

## Dashboard projection: Inventory audit (KHR-CV)

| Field | Value |
|-------|-------|
| `mode` | read-only backend projection |
| `source` | KHR-CU observation |
| Evidence | Karl-Dashboard `docs/evidence/khr-dashboard-inventory-audit-projection/committed-khr-cv-v1/` |

---

## Cockpit audit beta preview (KHR-CW)

| Field | Value |
|-------|-------|
| `component` | `audit-beta-gate` read-only status |
| Evidence | Karl-Dashboard `docs/evidence/khr-dashboard-cockpit-audit-beta-preview/committed-khr-cw-v1/` |

---

## Cockpit audit beta mounted preview (KHR-CX)

| Field | Value |
|-------|-------|
| `component` | `audit-beta-gate` mounted read-only tile in Cockpit preview shell |
| `runtimeMutation` | `false` — no rollout, no enforcement |
| Evidence | Karl-Dashboard `docs/evidence/khr-dashboard-cockpit-audit-beta-mounted-preview/committed-khr-cx-v1/` |

---

## KARL 2.0 live reference state (KHR-CY)

| Field | Value |
|-------|-------|
| `liveReferenceReady` | `true` on `karl-metal-01@ovh` reference path |
| `productionReady` | `false` — not production default |
| `rcReady` | `true` under sandbox/reference constraints only |
| Evidence | Karl-Installer `docs/evidence/karl2-live-reference-bundle/committed-khr-cy-v1/` |

No hyperdensity runtime mutation or rollout in CY.

---

## Enterprise module alignment (KHR-CY-E)

Migration-Factory, Licenziatore, DLP, Warden handoffs on branch `KHR`; `enterpriseAlignmentReady` required before KHR-CZ activation. Karl-App out of scope.

---

## Reference env activation read-only (KHR-CZ)

| Field | Value |
|-------|-------|
| **Status** | Read-only verification PASS on `karl-metal-01@ovh` |
| **Evidence** | Karl-Installer `docs/evidence/karl2-reference-env-activation-readonly/committed-khr-cz-v1/` |

No hyperdensity runtime mutation in CZ.

---

## Reference-preview flag-on dry smoke (KHR-DA)

Operator runbook and read-only flag-on/rollback smoke; no hyperdensity mutation.

---

## LIVE REFERENCE release declaration (KHR-DB)

KARL 2.0 declared **LIVE REFERENCE** on `karl-metal-01@ovh` (not production). Evidence: Karl-Installer `docs/evidence/karl2-live-reference-release-declaration/committed-khr-db-v1/`.

---

## LIVE REFERENCE operator acceptance (KHR-DC)

Operator acceptance pack PASS: activation, smoke, and rollback checklists; `productionReady=false`. Evidence: Karl-Installer `docs/evidence/karl2-live-reference-operator-acceptance/committed-khr-dc-v1/`.

---

## LIVE REFERENCE release index (KHR-DD)

Authoritative CO→DC audit manifest with pinned repo SHAs. Evidence: Karl-Installer `docs/evidence/karl2-live-reference-release-index/committed-khr-dd-v1/`.

---

## LIVE REFERENCE audit export (KHR-DE)

Operator/auditor export bundle. Evidence: Karl-Installer `docs/evidence/karl2-live-reference-audit-export/committed-khr-de-v1/`.

---

## LIVE REFERENCE operational rehearsal (KHR-DF)

Live reference-console rehearsal; no Hyperdensity runtime mutation. Evidence: Karl-Installer `docs/evidence/karl2-live-reference-operational-rehearsal/committed-khr-df-v1/`.

---

## LIVE REFERENCE bounded preview window (KHR-DG)

Bounded reference-preview window with periodic read-only probes; no Hyperdensity runtime mutation. Evidence: Karl-Installer `docs/evidence/karl2-live-reference-bounded-preview-window/committed-khr-dg-v1/`.

---

## LIVE module readiness scan (KHR-DH)

Read-only operating-posture scan; Hyperdensity focus in `hyperdensity-readiness.json`. Evidence: Karl-Installer `docs/evidence/karl2-live-module-readiness-scan/committed-khr-dh-v1/`.

---

## rdp-GW gateway posture (KHR-DI)

Access graph/continuity read-only via Dashboard projection; live gateway accepted inactive. Evidence: Karl-Installer `docs/evidence/karl2-rdpgw-live-reference-posture/committed-khr-di-v1/`. Posture: `LIVE_REFERENCE_PARTIAL_ACCEPTED`.

---

## Operating posture consolidation (KHR-DJ)

`hyperdensityOperating=true`, read-only. Evidence: Karl-Installer `docs/evidence/karl2-live-reference-operating-posture/committed-khr-dj-v1/`.

---

## Hyperdensity operating window (KHR-DK)

Bounded read-only stability window; `stableAcrossWindow=true`. Evidence: Karl-Installer `docs/evidence/karl2-hyperdensity-live-reference-operating-window/committed-khr-dk-v1/`.

---

## Promotion boundary (KHR-DL)

`liveReferenceStable=true`; `promotionAllowed=false`; `nextAllowedScope=reference-only`. Evidence: Karl-Installer `docs/evidence/karl2-live-reference-promotion-boundary/committed-khr-dl-v1/`.

---

## rdp-GW alignment plan (KHR-DM)

Path A LIVE REFERENCE only; production requires rdp-GW operating. Evidence: Karl-Installer `docs/evidence/karl2-rdpgw-alignment-plan/committed-khr-dm-v1/`. No-op for Hyperdensity.

---

## rdp-GW KHR route (KHR-DP)

Route validated on `khr-dp-v1`; reference gateway rolled back inactive. Evidence: Karl-Installer `docs/evidence/karl2-rdpgw-khr-route-bounded-activation/committed-khr-dp-v1/`. No-op.

---

## rdp-GW operating window (KHR-DQ)

`hyperdensityOperating=true` during rdpgw window; compatibility note only. Evidence: Karl-Installer `docs/evidence/karl2-rdpgw-live-reference-operating-window/committed-khr-dq-v1/`.

---

## rdp-GW sustained-readiness preflight (KHR-DR)

Schema hardening + preflight only; no-op for Hyperdensity. Evidence: Karl-Installer `docs/evidence/karl2-rdpgw-sustained-readiness-preflight/committed-khr-dr-v1/`.

---

## rdp-GW quota unblock guarded restore (KHR-DS)

Quota-safe manifest + guarded `karl-quota` restore; no-op for Hyperdensity. Evidence: Karl-Installer `docs/evidence/karl2-rdpgw-quota-unblock-guarded-restore/committed-khr-ds-v1/`.

---

## rdp-GW sustained LIVE REFERENCE enablement (KHR-DT)

Sustained LIVE REFERENCE operating posture; Hyperdensity operating during gate. Evidence: Karl-Installer `docs/evidence/karl2-rdpgw-sustained-live-reference-enable/committed-khr-dt-v1/`.

---

## Phase: crd-foundation (KHR-CM)

| Field | Value |
|-------|-------|
| `phase` | `crd-foundation` |
| `runtimeMutation` | `false` |
| `hostRuntimeEnabled` | `false` |
| `systemdEnable` | `false` |
| Evidence | `docs/evidence/khr-baremetal-khr-native-crd-foundation/committed-khr-cm-v1/` |

---

## Related

- Karl-Installer `INSTALLER_KARL2_BAREMETAL_KHR_NATIVE_PROFILE.md`
- `KHR_AUTO_CONFIGURATION_PLAN.md`
