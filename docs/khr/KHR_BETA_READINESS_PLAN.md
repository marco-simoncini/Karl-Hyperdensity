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
| Dashboard LIVE_PASS | `committed-khr-bs-20260517T073046Z` | snapshot `dashboardLivePassRef` |
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
