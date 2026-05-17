# KHR beta readiness plan (from Snapshot v1)

Beta readiness is **evidence-backed** from **KHR TP Reference Snapshot v1** — not live cluster state at validate time.

| Field | Value |
|-------|-------|
| **Snapshot** | `committed-khr-bt-v1` |
| **Cluster (evidence origin)** | `karl-metal-01@ovh` |
| **contractSetId** | `khr-tp-contract-v1` |
| **providerProfile** | `khr-native` (reference env only) |

---

## Readiness gates (offline)

| Gate | Evidence | Status source |
|------|----------|---------------|
| Scope 1 | `committed-scope1-khr-aw` | snapshot `scopeReadiness.scope1` |
| Scope 2 loop | `committed-scope2-loop-khr-ba` | `manual-loop-pass` |
| Scope 3 dry-run | `committed-scope3-dryrun-khr-bc` | `dryrun-summary.json` PASS |
| Scope 4 certification | `committed-scope4-certification-khr-bf` | `certified-evidence-backed` |
| Governance | `committed-scope4-governance-khr-bg` | `certified` |
| Dashboard LIVE_PASS (backend) | `committed-khr-bs-20260517T073046Z` | snapshot `dashboardLivePassRef` |
| Dashboard UI preview LIVE_PASS | `committed-khr-cb-v1` | `docs/evidence/khr-dashboard-ui-projection-preview/` — KHR-CB |
| rdp-GW sandbox | `committed-cluster-sandbox-khr-ay` | snapshot `rdpgwClusterSandboxRef` |
| Installer CRD | `khr-installer-crd-foundation/20260517T070416Z` | snapshot |
| ISO post-install | `khr-post-install-verify/summary.json` | snapshot |
| Hybrid transition | `khr-hybrid-transition/20260516T195854Z` | snapshot |

Validate: `./scripts/khr_validate_reference_snapshot.sh` (included in default `./scripts/validate.sh`).

---

## Explicit non-goals (beta boundary)

- No production enable
- No autonomous orchestration
- No global Dashboard default change
- No mutating gateway revoke/disconnect automation
- No Scope-4 **active** enablement (`scope4Active=false` in snapshot)

---

## Cross-repo beta docs

| Repo | Doc |
|------|-----|
| Karl-Dashboard | `DASHBOARD_KHR_BETA_READINESS_FROM_SNAPSHOT.md` |
| Karl-Installer | `INSTALLER_KHR_BETA_READINESS_FROM_SNAPSHOT.md` |
| Karl-OS-ISO | `ISO_KHR_BETA_READINESS_FROM_SNAPSHOT.md` |
| rdp-GW | `RDPGW_KHR_BETA_READINESS_FROM_SNAPSHOT.md` |

---

## Inventory live ingest (KHR-BW / KHR-BX)

| Item | Status |
|------|--------|
| Beta blocker `inventory-live-ingest` | **Dashboard-visible** (KHR-BX projection) |
| Evidence | Karl-Inventory `docs/evidence/khr-inventory-live-ingest/committed-khr-bw-v1/` |
| Source | `live-readonly` (committed Hyperdensity snapshot + federation) |
| Dashboard | `inventoryObservationSummary` on khr-backend/projection (read-only) |
| Enforcement | **None** — `enforcement=false`, `mutationObserved=false` |

Expected fields: `observationId`, `observationSource`, `observedAt`, `snapshotRef`, `inventoryObservationSource`, `postureObserved`, `scopeReadinessObserved`, `continuityObserved`.

See `INVENTORY_LIVE_INGEST_EXPECTATIONS.md`, Karl-Dashboard `DASHBOARD_INVENTORY_LIVE_INGEST_PROJECTION.md`.

---

## Dashboard UI projection preview live (KHR-CB)

| Item | Status |
|------|--------|
| Endpoint | `GET /api/hyperdensity/khr-ui/projection-preview` |
| Evidence | Karl-Dashboard `docs/evidence/khr-dashboard-ui-projection-preview/committed-khr-cb-v1/` |
| `evidenceStatus` | `LIVE_PASS` (`source=live-readonly`, dual-flag active) |
| Rollback | `rollback-summary.json` → `rollbackVerified=true` |
| Distinct from | Backend projection LIVE_PASS (KHR-BS) — same rollout pattern, different endpoint |

See Karl-Dashboard `DASHBOARD_UI_KHR_PROJECTION_PREVIEW_EVIDENCE.md`.

---

## Baremetal standing reference profile (KHR-CC / KHR-CD)

| Item | Status |
|------|--------|
| Profile | `khr-native` on baremetal reference only (`karl-metal-01@ovh`) |
| Plan | Karl-Dashboard `docs/evidence/khr-dashboard-baremetal-standing-profile/committed-khr-cc-v1/` |
| Controlled apply evidence | Karl-Dashboard `docs/evidence/khr-dashboard-baremetal-standing-profile/committed-khr-cd-v1/` |
| `standingProfileState` | `rollback-verified` (post mandatory rollback) |
| `globalDefaultsChanged` | `false` |
| Public cloud | Remains `kubevirt.compatibility` primary for legacy fleets |

Anchors: KHR-BS + KHR-CB + KHR-CD (`VERIFY_PASS`, soak PASS, `rollbackVerified=true`) to `console:1.6.0`.

## Baremetal standing reference window (KHR-CE)

| Item | Status |
|------|--------|
| Window evidence | Karl-Dashboard `docs/evidence/khr-dashboard-baremetal-standing-window/committed-khr-ce-v1/` |
| `standingWindowState` | `window-closed-rollback-verified` |
| `standingWindowDurationSeconds` | `600` (committed); script default `1800` |
| `rollbackRequired` | `true` (executed) |
| `globalDefaultsChanged` | `false` |

See Karl-Dashboard `DASHBOARD_BAREMETAL_STANDING_WINDOW.md`.

## Baremetal standing window operations (KHR-CF)

| Item | Status |
|------|--------|
| Runbook | Karl-Dashboard `DASHBOARD_BAREMETAL_STANDING_WINDOW_OPERATIONS.md` |
| Ops evidence | `docs/evidence/khr-dashboard-baremetal-standing-window-ops/committed-khr-cf-v1/` |
| `standingWindowOperationalState` | `window-inactive-baseline` (post validation) |
| `operatorWindowLifecycle` | `inactive` |
| `abortRequired` | `false` at baseline |
| `rollbackIdempotent` | `true` |
| `globalDefaultsChanged` | `false` |

See Karl-Dashboard `DASHBOARD_BAREMETAL_KHR_NATIVE_STANDING_PROFILE.md`.

---

## Beta Candidate 0 (KHR-BV)

| Item | Path |
|------|------|
| Scope | `KHR_BETA_CANDIDATE_0_SCOPE.md` |
| Release marker | `KHR_BETA_CANDIDATE_0_RELEASE_MARKER.md` |
| Manifest | `docs/contracts/khr/khr-beta-candidate-0-manifest.json` |
| Check | `scripts/khr_beta_candidate_0_check.sh` |

Pseudo-tag: `khr-beta-candidate-0@committed-khr-bt-v1`

---

## Live re-validation (optional)

Operators may run `KHR_LIVE_VALIDATE=1 ./scripts/validate.sh` to refresh cluster checks. Live FAIL does not revoke snapshot v1 freeze; file a new snapshot version to promote new evidence.
