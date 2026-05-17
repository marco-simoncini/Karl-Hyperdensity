# KHR Baremetal Installer Profile Expectations (KHR-CL / KHR-CM)

| Field | Value |
|-------|-------|
| **Profile ID** | `karl2-baremetal-khr-native` |
| **Sprint** | KHR-CL / KHR-CM / KHR-CN / KHR-CO / KHR-CP / **KHR-CQ** |
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
| `guardedApplyBlocked` | `true` (`guarded-apply` rejected) |
| Evidence | Karl-Installer `docs/evidence/khr-baremetal-khr-native-guarded-apply-preflight/committed-khr-cq-v1/` |

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
