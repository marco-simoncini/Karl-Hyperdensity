# KHR Baremetal Installer Profile Expectations (KHR-CL / KHR-CM)

| Field | Value |
|-------|-------|
| **Profile ID** | `karl2-baremetal-khr-native` |
| **Sprint** | KHR-CL … KHR-CQ / **KHR-CR** / **KHR-CS** / **KHR-CT** / **KHR-CU** / **KHR-CV** / **KHR-CW** / **KHR-CX** / **KHR-CY** / **KHR-DC** |
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
