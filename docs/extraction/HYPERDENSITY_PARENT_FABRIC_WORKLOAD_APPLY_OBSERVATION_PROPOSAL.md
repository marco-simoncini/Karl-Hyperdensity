# Hyperdensity Parent Fabric — apply observation proposal (Sprint 65)

## Summary

**Sprint 65** is **proposal-only**: inventory and criteria for a future **`ApplyObservationWiredV1`** phase. **Sprint 66** shadow matrix. **Sprint 67** staged wrappers. **Sprint 68** wrapper hardening. **Sprint 69** wiring readiness. **Sprint 70** may wire `apply.go` call sites. **`apply.go` remains legacy in Sprint 69.** Hyperdensity receives **no** new Go adapter code.

---

## 1. Scope

| Item | Sprint 65 |
|------|-----------|
| Apply observation proposal + criteria docs | **Yes** |
| Dashboard apply-only audit script + policy test | **Yes** |
| `apply.go` wiring | **No** |
| Broad `ObservationWiredV1` flip | **No** |
| Hyperdensity `parentfabric` import in Dashboard | **No** |

---

## 2. Non-goals

- Wiring `apply.go` to adapter observation wrappers or candidates.
- Enabling `ApplyObservationWiredV1 = true`.
- Resource-exchange, rollback, VM runtime, or admission_guard observation.
- Moving runtime code from Dashboard to Hyperdensity.
- Changing API responses, JSON ordering, or path/pilot/live wiring (Sprints 56–63).

---

## 3. Current observation status

| Surface | Flag | Status |
|---------|------|--------|
| Path | `PathWiredV1` | **`true`** |
| Pilot | `PilotObservationWiredV1` | **`true`** |
| Live | `LiveObservationWiredV1` + candidate | **`true`** |
| Apply | `ApplyObservationWiredV1` | **`false`** (placeholder only) |
| Broad | `ObservationWiredV1` | **`false`** (deliberate) |
| Production | `ProductionWiredV1` | **`false`** |

`apply.go` uses **legacy** pod observation helpers only (4 call sites per Sprint 65 audit).

Full **`workload_helpers.go`** verdict remains **`copy-deferred`**.

---

## 4. Apply observation risks

| Risk | Level |
|------|-------|
| Apply path mutates plan targets from observed pod state | **High** |
| Confusion with live/pilot wiring completeness | **Medium** |
| Accidental broad `ObservationWiredV1` flip | **High** |
| Resource-exchange coupling if apply wired without boundary | **High** |

---

## 5. Proposed future flag

```text
hyperdensityWorkloadAdapterApplyObservationWiredV1 = false  // Sprint 65 placeholder
```

Activation requires a **dedicated sprint** after shadow matrix PASS. See **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CRITERIA.md`**.

---

## 6. Required shadow matrix

Before any flip, each apply observation call-site must have:

- Legacy vs adapter-wrapper triple-assert (or equivalent) per helper invocation.
- Apply-only golden fixture (no resource_exchange / rollback / VM).
- Parity runner green with apply audit script PASS.

---

## 7. Forbidden surfaces

Until explicitly approved in later sprints:

- `resource_exchange_*`
- `rollback.go` observed-state
- VM runtime files
- `admission_guard_*`
- execution/apply mode helpers (unchanged)
- Broad `ObservationWiredV1 = true`

---

## 8. Rollback

Sprint 65 is docs/audit only — no runtime rollback. Future apply wiring rollback: set `ApplyObservationWiredV1 = false` and revert `apply.go` call sites to legacy helpers only.

---

## 9. Test/audit requirements

| Artifact | Owner |
|----------|-------|
| `audit_workload_apply_observation.sh` | Dashboard |
| `hyperdensity_parent_fabric_workload_apply_observation_policy_test.go` | Dashboard |
| `testdata/..._apply_observation_proposal.golden.json` | Dashboard |

---

## 10. Risks

Premature apply wiring without isolated shadow coverage could change apply plan target fields while broad observation policy still forbids general flip.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CRITERIA.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_BROAD_OBSERVATION_DECISION.md`
