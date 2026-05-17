# KHR Baremetal Installer Profile Expectations (KHR-CL)

| Field | Value |
|-------|-------|
| **Profile ID** | `karl2-baremetal-khr-native` |
| **Sprint** | KHR-CL |
| **Mode** | Plan/dry-run + guarded reference only |

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

## Guard

`KARL_INSTALLER_BAREMETAL_KHR_NATIVE_I_UNDERSTAND=true` required. Without it, installer fails safe.

---

## Related

- Karl-Installer `INSTALLER_KARL2_BAREMETAL_KHR_NATIVE_PROFILE.md`
- `KHR_AUTO_CONFIGURATION_PLAN.md`
