# Hyperdensity Parent Fabric — apply observation flip risks (Sprint 72)

## Summary

**Sprint 72** documents why flipping apply observation flags is **more risky** than Sprint 70 call-site wiring. Wiring repoints call sites to wrappers that still delegate to legacy; a flag flip changes **which implementation runs** at runtime.

---

## 1. Why flip is riskier than call-site wiring

| Aspect | Sprint 70 (call-site wiring) | Future flip sprint |
|--------|------------------------------|-------------------|
| Runtime behavior | **Unchanged** (legacy branch) | **May change** |
| Flags | All **false** | One or both **true** |
| Risk surface | Source structure only | Semantic output paths |
| Rollback | Restore 4 legacy calls | Flag false + possible call-site restore |
| Detection | 8×4 matrix (wrapper ≡ legacy) | Matrix + live apply payloads |

Call-site wiring proved **structural** equivalence. Flag flip proves **semantic** equivalence under production branch selection.

---

## 2. Flag granularity (do not conflate)

| State | Meaning |
|-------|---------|
| **Call-site wired** | `apply.go` calls `hyperdensityWorkloadApplyObserved*V1` wrappers |
| **Candidate runtime used** | `ApplyObservationCandidateRuntimeUsedV1=true` — wrapper may call candidate helpers |
| **Apply observation wired** | `ApplyObservationWiredV1=true` — apply observation adapter path active |
| **Broad observation wired** | `ObservationWiredV1=true` — **out of scope** for apply-only track |

Sprint 70 achieved **call-site wired** only. Sprint 71 confirmed legacy-equivalent behavior with flags false. Neither implies apply observation is "fully wired" at the adapter level.

---

## 3. Confusion risks between granular and broad flags

- Operators may assume `ApplyObservationWiredV1=true` enables all observation surfaces.
- `ObservationWiredV1=true` would wire **many** legacy call sites — **forbidden** in the same sprint as apply flip without dedicated broad policy.
- `ProductionWiredV1=true` is a separate production track — must remain **false**.

**Mitigation:** dedicated sprint titles; audit script guards; golden fields per subflag.

---

## 4. resource_exchange coupling risk

Apply observation helpers read **pod container resources** during apply plan construction. resource_exchange uses **different** stage-apply surfaces.

**Risk:** a flip sprint accidentally imports wrapper/candidate names into `resource_exchange_*` or shares flags with exchange admission.

**Mitigation:** explicit exclusion in flip criteria; audit grep for wrapper/candidate in resource_exchange files.

---

## 5. Rollback path

| Severity | Rollback |
|----------|----------|
| Flag flip only | Set `ApplyObservationCandidateRuntimeUsedV1` / `ApplyObservationWiredV1` to **false** |
| Wrapper regression | Restore four legacy calls in `apply.go` |
| Broad accidental flip | Set `ObservationWiredV1=false`; audit all observation call sites |

Rollback must **not** modify rollback.go, VM runtime, or resource_exchange as a side effect of apply rollback.

---

## 6. Dedicated sprint requirement

**Do not** combine in one sprint without explicit approval:

1. `ApplyObservationCandidateRuntimeUsedV1=true`
2. `ApplyObservationWiredV1=true`
3. `ObservationWiredV1=true`
4. resource_exchange changes
5. VM runtime observation changes

**Recommendation:** two flip sprints minimum — (A) candidate runtime used, (B) apply observation wired — each with 8×4 re-run and parity green.

---

## 7. Hyperdensity boundary

Hyperdensity receives **no** new Go adapter code for apply observation flip. Criteria and audit live in Dashboard + extraction docs. `workload_helpers.go` remains **`copy-deferred`**.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_CRITERIA.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_BROAD_OBSERVATION_DECISION.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_POST_WIRING_HARDENING.md`
