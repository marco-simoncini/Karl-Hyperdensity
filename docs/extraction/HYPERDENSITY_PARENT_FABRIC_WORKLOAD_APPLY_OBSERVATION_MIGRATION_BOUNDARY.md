# Hyperdensity Parent Fabric — apply observation migration boundary (Sprint 77)

## Summary

**Sprint 77** formally closes the **apply-observation track** (Sprint 65–76). Apply observation is **complete** in scope: four `apply.go` helpers wired through staged wrappers with candidate branch active and **legacy-equivalent** behavior. **Broad observation remains disabled.**

---

## 1. Scope

| Item | Sprint 77 |
|------|-----------|
| Migration boundary documentation | **Yes** |
| Next-surface decision documentation | **Yes** |
| Boundary golden + audit/test | **Yes** (Dashboard) |
| Flag flips | **No** |
| Runtime / wiring changes | **No** |

---

## 2. Non-goals

- `ObservationWiredV1 = true` (broad observation).
- `ProductionWiredV1 = true`.
- Wiring resource_exchange, rollback, or VM runtime observation.
- Hyperdensity Go adapter code or Dashboard `parentfabric` import.
- Changing `workload_helpers.go` verdict (stays **`copy-deferred`**).

---

## 3. Apply observation completed scope

| Surface | Status |
|---------|--------|
| Four `apply.go` pod resource observation helpers | **Complete** (wrappers + candidate branch) |
| `ApplyObservationWiredV1` | **`true`** |
| `ApplyObservationCandidateRuntimeUsedV1` | **`true`** |
| wrapper ≡ candidate ≡ legacy (8×4) | **PASS** |
| Direct legacy observation calls in `apply.go` | **0** |
| Direct candidate calls in `apply.go` | **0** |

---

## 4. Completed sprint chain Sprint 65–76

| Sprint | Deliverable |
|--------|-------------|
| **65** | Proposal + criteria |
| **66** | Shadow matrix |
| **67** | Staged wrappers |
| **68** | Wrapper hardening 8×4 |
| **69** | Wiring readiness |
| **70** | `apply.go` call-site wiring (4 wrappers) |
| **71** | Post-wiring hardening |
| **72** | Flip criteria + risks |
| **73** | Candidate-runtime readiness |
| **74** | `CandidateRuntimeUsedV1 = true` (staging) |
| **75** | `ApplyObservationWiredV1 = true` (activation) |
| **76** | Post-activation hardening |

---

## 5. Current flags

| Flag | Value |
|------|-------|
| `ApplyObservationWiredV1` | **`true`** |
| `ApplyObservationCandidateRuntimeUsedV1` | **`true`** |
| `ObservationWiredV1` | **`false`** |
| `ProductionWiredV1` | **`false`** |
| Candidate branch (apply wrappers) | **active** |

---

## 6. Why broad observation remains false

- Apply activation wires **only** the four apply observation helpers.
- Broad `ObservationWiredV1` would affect **many** legacy call sites across pilot, live, apply, resource_exchange, rollback, VM, etc.
- Sprint 64 policy and Sprint 77 boundary explicitly forbid inferring broad flip from apply completion.
- Remaining legacy observation surfaces require **dedicated** tracks with separate criteria and shadow matrices.

---

## 7. Remaining observation surfaces

See **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_OBSERVATION_NEXT_SURFACE_DECISION.md`**.

| Track | Status |
|-------|--------|
| resource_exchange | legacy — separate sprint |
| rollback | legacy — safety-critical — audit only |
| VM runtime | legacy — separate sprint |
| usage.go / other-review | classification needed |

---

## 8. Migration boundary decision

**Decision:** Apply-observation migration boundary is **COMPLETE**. The next sprint must **not** automatically set `ObservationWiredV1=true`.

**Recommended next track:** `resource_exchange_observation_audit` (proposal/audit only) or `usage.go` classification — not broad observation.

**Sprint 78 update:** `resource_exchange_observation_audit` **complete** (audit/proposal only). See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_AUDIT.md`. **No** resource_exchange wiring. Recommended next: `resource_exchange_observation_shadow_matrix`.

---

## 9. Rollback posture

| Step | Action |
|------|--------|
| 1 | `ApplyObservationWiredV1 = false` — wrappers use legacy branch |
| 2 | `ApplyObservationCandidateRuntimeUsedV1 = false` — optional |
| 3 | Restore four legacy calls in `apply.go` — only if undoing Sprint 70 wiring |
| 4 | Re-run parity + audits |

Boundary closure does not remove rollback documentation from prior sprints.

---

## 10. Risks

- Treating apply track complete as permission for broad flip.
- Bundling resource_exchange wiring with apply boundary closure.
- Candidate helper drift breaking equivalence while flags stay true.

---

## 11. Recommended next tracks

1. **resource_exchange observation audit** — high risk; shadow matrix before any wiring.
2. **usage.go classification** — medium risk; inventory only.
3. **rollback observation audit** — very high / safety-critical; no wiring before resource_exchange policy.
4. **VM runtime observation audit** — high risk; no Windows runtime claim.

**Not recommended:** `ObservationWiredV1=true` as next sprint.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_OBSERVATION_NEXT_SURFACE_DECISION.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_BROAD_OBSERVATION_DECISION.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_POST_ACTIVATION_HARDENING.md`
