# KHR Technical Preview Reference Snapshot v1 (KHR-BT)

Versioned, **read-only** cross-repo index of committed Technical Preview live evidence on `karl-metal-01@ovh`. No new rollout, guarded apply, or cluster mutation in this sprint.

## Purpose

| Goal | Detail |
|------|--------|
| Single snapshot | `snapshot-summary.json` aggregates scope 1–4, governance, Dashboard LIVE_PASS, rdp-GW cluster-sandbox, Installer CRD/hybrid, ISO post-install |
| Cross-repo index | `cross-repo-evidence-index.json` lists every anchored artifact with repo + path |
| No mutation proof | `liveMutationPerformed=false`, `noNewRollout=true`, `globalDefaultsChanged=false` |

## Contract

| Artifact | Path |
|----------|------|
| JSON Schema (draft) | `docs/contracts/khr/khr-tp-reference-snapshot-v1.json` |
| Aggregator | `scripts/khr_tp_reference_snapshot_v1.sh` |
| Output | `docs/evidence/khr-tp-reference-snapshot-v1/<runId>/snapshot-summary.json` |

## Aggregated fields

| Field | Source |
|-------|--------|
| `contractSetId` | `khr-tp-contract-v1` (all repos) |
| `scopeReadiness` | Scope 1–4 committed verify/preflight summaries |
| `scope4CertificationState` | `committed-scope4-certification-khr-bf` + governance `KHR-BG` |
| `governanceState` | `committed-scope4-governance-khr-bg` |
| `providerProfile` | Dashboard LIVE_PASS (`khr-native`) |
| `dashboardLivePassRef` | Karl-Dashboard `committed-khr-bs-20260517T073046Z` |
| `rdpgwClusterSandboxRef` | rdp-GW `committed-cluster-sandbox-khr-ay` |
| `installerCrdFoundationRef` | Karl-Installer `karl2-khr-technical-preview` evidence |
| `hybridTransitionRef` | Karl-Installer `hybrid-transition` evidence |

## Operator run (read-only)

```bash
export KHR_TP_REFERENCE_SNAPSHOT_RUN_ID=committed-khr-bt-v1
# optional repo paths if not sibling checkout:
# export KHR_DASHBOARD_PATH=/path/to/Karl-Dashboard
# export KHR_INSTALLER_PATH=/path/to/Karl-Installer
# export KHR_ISO_PATH=/path/to/Karl-OS-ISO
# export KHR_RDP_GW_PATH=/path/to/rdp-GW

./scripts/khr_tp_reference_snapshot_v1.sh
```

## Committed snapshot (KHR-BT)

| Item | Value |
|------|-------|
| **runId** | `committed-khr-bt-v1` |
| **Dashboard LIVE_PASS** | `committed-khr-bs-20260517T073046Z` |
| **Image** | `registry.karl-technology.com/karl-evolution/console:khr-reference-df0ae11f9` |
| **Rollback verified** | `rollbackVerified=true` (post-evidence; cluster restored to `1.6.0`) |

## Constraints (unchanged)

- No production enable (`productionReady=false`)
- No autonomous orchestration
- No global ISO/installer/Dashboard defaults change
- Scope-4 guarded apply **not active** (`readyForScope4Active=false`)
- No Dashboard UI changes

## Related

- `KHR_TP_LIVE_REFERENCE_ENVIRONMENT.md`
- `TECHNICAL_PREVIEW_PACKAGE.md`
- Per-repo: `DASHBOARD_REFERENCE_SNAPSHOT_V1.md`, `INSTALLER_REFERENCE_SNAPSHOT_V1.md`, `ISO_REFERENCE_SNAPSHOT_V1.md`, `RDPGW_REFERENCE_SNAPSHOT_V1.md`
