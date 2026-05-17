# KHR Beta Candidate 0 — release marker (KHR-BV)

Pseudo-tag and cross-repo manifest for **Beta Candidate 0**, derived from **Reference Snapshot v1**.

| Field | Value |
|-------|-------|
| **Marker id** | `khr-beta-candidate-0` |
| **Pseudo-tag** | `khr-beta-candidate-0@committed-khr-bt-v1` |
| **Manifest** | `docs/contracts/khr/khr-beta-candidate-0-manifest.json` |
| **Check** | `scripts/khr_beta_candidate_0_check.sh` |

---

## Snapshot anchor

| Item | Path |
|------|------|
| runId | `committed-khr-bt-v1` |
| Summary | `docs/evidence/khr-tp-reference-snapshot-v1/committed-khr-bt-v1/snapshot-summary.json` |
| Index | `.../cross-repo-evidence-index.json` |

---

## Validation mode

| Mode | Command | Cluster |
|------|---------|---------|
| **Offline (default)** | `./scripts/validate.sh` | Not required |
| **Snapshot + beta marker** | `./scripts/khr_beta_candidate_0_check.sh` | Not required |
| **Live (optional)** | `KHR_LIVE_VALIDATE=1 ./scripts/validate.sh` | `karl-metal-01@ovh` |

---

## Repo commit pins (manifest)

Pinned at Beta Candidate 0 cut (minimum ancestors; see manifest `repoCommits` for full SHAs):

| Repo | Role |
|------|------|
| Karl-Hyperdensity | Contracts, evidence, validation modes |
| Karl-Dashboard | LIVE_PASS, projection, provider profiles |
| Karl-Installer | CRD foundation, hybrid |
| Karl-OS-ISO | Post-install verify |
| rdp-GW | Gateway compatibility + sandbox |
| Karl-Inventory | Observation export stub |

---

## Beta acceptance criteria

| # | Criterion | Proof |
|---|-----------|-------|
| 1 | Offline validate PASS (all repos) | `validate.sh` / repo-specific validators |
| 2 | Live validate optional | `KHR_LIVE_VALIDATE=1` documented only |
| 3 | No production enable | snapshot `productionReady=false` |
| 4 | Rollback + certification present | Dashboard rollback evidence; scope4 certification + governance |
| 5 | Dashboard backend projection SoT | Hyperdensity contracts + `DASHBOARD_KHR_BACKEND_PROJECTION_API.md` |
| 6 | Provider profile propagation stable | Fixtures + `validate_khr_provider_profile_propagation.sh` |
| 7 | Reference snapshot PASS | `khr_validate_reference_snapshot.sh` |
| 8 | No new mutation in BV | Docs/manifest/check only |

---

## Evidence refs (from snapshot)

| Ref | Committed id |
|-----|----------------|
| Dashboard LIVE_PASS | `committed-khr-bs-20260517T073046Z` |
| rdp-GW sandbox | `committed-cluster-sandbox-khr-ay` |
| Installer CRD | `20260517T070416Z` |
| Hybrid | `20260516T195854Z` |
| Scope-4 cert | `committed-scope4-certification-khr-bf` |
| Governance | `committed-scope4-governance-khr-bg` |

---

## Explicit non-goals

- No runtime mutation in KHR-BV
- No rollout or default changes
- No production enable
- No autonomous orchestration
