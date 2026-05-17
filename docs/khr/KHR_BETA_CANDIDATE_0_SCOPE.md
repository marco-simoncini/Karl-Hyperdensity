# KHR Beta Candidate 0 — scope (KHR-BV)

Cross-repo scope definition for **Beta Candidate 0**, anchored on **Reference Snapshot v1** (`committed-khr-bt-v1`).

| Field | Value |
|-------|-------|
| **Release marker** | `khr-beta-candidate-0` |
| **Pseudo-tag** | `khr-beta-candidate-0@committed-khr-bt-v1` |
| **contractSetId** | `khr-tp-contract-v1` |
| **Cluster (evidence origin)** | `karl-metal-01@ovh` |

---

## In scope (Beta Candidate 0)

| Area | Deliverable |
|------|-------------|
| **Contracts & CRDs** | ShellLease, GatewayRoute, ResourceLease, ResourcePort, EvidenceBundle (schema + examples) |
| **Hyperdensity** | Scope 1–4 committed evidence; certification + governance; offline validation modes (KHR-BU) |
| **Dashboard** | KHR backend projection API; provider profiles; reference-env LIVE_PASS evidence; rollback verified |
| **Installer** | karl2 CRD foundation evidence; hybrid transition evidence |
| **ISO** | Post-install verify; profile manifest; host-runtime disabled on ISO |
| **rdp-GW** | Legacy pool/replica preserved; read-only ShellLease/GatewayRoute resolvers; cluster-sandbox evidence |
| **Inventory** | TP observation export stub; scope observation prep docs (offline) |
| **Validation** | Default offline `validate.sh`; optional `KHR_LIVE_VALIDATE=1` |

---

## TP-only (not promoted to beta product)

| Item | Reason |
|------|--------|
| Scope-4 **active** enablement | `scope4Active=false` in snapshot |
| Reference-env console image rollout | Evidence on reference env only; global default unchanged |
| Manual ResourcePort loop / dry-run paths | Operator-gated; not autonomous |
| Grand Padre evidence ingest automation | Observation contracts only |
| Production canary / auto-apply | Explicitly out of beta-0 |

---

## Experimental (documented, not beta-gated)

| Item | Notes |
|------|-------|
| Access graph continuity (live-readonly) | rdp-GW `/karl-gw/v1/accessgraph/session` |
| Session identity resolve | KHR-AN read-only correlation |
| Federation observation | Cross-repo metadata only |
| Windows FluidVirt lane | Planning-only in Hyperdensity |
| Inventory live ingest | Stub/export only — **beta blocker** for full beta |

---

## Out of scope (Beta Candidate 0)

| Item |
|------|
| Production enable / GA |
| Autonomous orchestration |
| Global Dashboard default change (`khr-native` not default) |
| Mutating gateway revoke / disconnect |
| New Dashboard ShellLease/GatewayRoute UI |
| Inventory mandatory agent / live cluster ingest |
| Windows parity production claims |
| Release packaging / installer GA |
| Persistence / disaster-recovery productization |

---

## Acceptance criteria

See `KHR_BETA_CANDIDATE_0_RELEASE_MARKER.md` and `docs/contracts/khr/khr-beta-candidate-0-manifest.json`.

Validate: `./scripts/khr_beta_candidate_0_check.sh`

---

## Beta blockers (post candidate-0)

| Blocker | Owner repo |
|---------|------------|
| Inventory live ingest | Karl-Inventory |
| Dashboard UI consumption of GatewayRoute/ShellLease | Karl-Dashboard |
| Windows parity | Hyperdensity / Windows lane |
| Persistence / recovery | Hyperdensity |
| Release packaging | Installer / ISO |

---

## Related

- `KHR_BETA_READINESS_PLAN.md`
- `KHR_SNAPSHOT_V1_FREEZE_POLICY.md`
- `KHR_VALIDATION_MODES.md`
- Per-repo: `../Karl-Dashboard/docs/khr/DASHBOARD_BETA_CANDIDATE_0_SCOPE.md` (and siblings)
