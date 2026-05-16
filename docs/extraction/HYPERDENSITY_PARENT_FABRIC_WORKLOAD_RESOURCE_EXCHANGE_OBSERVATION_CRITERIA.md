# Hyperdensity Parent Fabric — resource exchange observation criteria (Sprint 78+)

Criteria for **future** sprints wiring resource_exchange observation. Sprint 78–80 = CPU path. Sprint 81 = local helper shadow (not wired).

---

## 1. Dedicated sprint per milestone

| Milestone | Requirement |
|-----------|-------------|
| Shadow matrix | Dedicated sprint; no bundling with apply/rollback/VM |
| Staged wrappers | Dedicated sprint after matrix PASS |
| Candidate runtime | Dedicated sprint; separate flag from wired |
| Activation flip | Dedicated sprint; explicit approval |

---

## 2. Policy gates (all future sprints)

| Gate | Required value |
|------|----------------|
| `ObservationWiredV1` | **`false`** until broad policy satisfied (not resource_exchange sprint) |
| `ProductionWiredV1` | **`false`** |
| `ApplyObservationWiredV1` | **unchanged** (`true` post Sprint 75) — resource_exchange must not depend on it |
| `ResourceExchangeObservationWiredV1` | **`false`** until dedicated activation sprint |
| Broad observation | **forbidden** in resource_exchange sprints |

---

## 3. Exclusions

- **No** rollback observation wiring in resource_exchange sprints.
- **No** VM runtime observation wiring.
- **No** `admission_guard_*` changes.
- **No** `apply.go` or apply wrapper/candidate reuse in `resource_exchange_*`.
- **No** Dashboard `pkg/hyperdensity/parentfabric` import.

---

## 4. resource_exchange-only golden

Future goldens must:

- Scan only `hyperdensity_parent_fabric_resource_exchange_*.go` (non-test).
- Record real call-site counts (no invented numbers).
- Assert `resourceExchangeUsesApplyWrappers == false`.
- Assert `resourceExchangeUsesApplyCandidates == false`.
- Assert `broadObservationAllowed == false`.

---

## 5. Audit script requirements

Extend or complement `audit_workload_resource_exchange_observation.sh`:

- Inventory listed observation helpers.
- Fail if `ResourceExchangeObservationWiredV1=true`, `ObservationWiredV1=true`, or `ProductionWiredV1=true`.
- Fail if apply wrappers/candidates appear in resource_exchange files.
- Fail on parentfabric import in resource_exchange files.
- Do **not** analyze `apply.go`.

---

## 6. Wrapper / candidate shadow matrix

Before wiring:

1. Define resource_exchange-scoped wrapper names (not apply names).
2. Define pure candidates in adapter/candidate files (test-only until flip).
3. Run legacy ≡ wrapper ≡ candidate matrix per helper.
4. Document local helpers (container ready/restart) explicitly.

---

## 7. Rollback strategy

| Step | Action |
|------|--------|
| 1 | `ResourceExchangeObservationWiredV1 = false` |
| 2 | Wrappers fall back to legacy branch |
| 3 | Optional: `CandidateRuntimeUsed` flag false for resource_exchange only |
| 4 | Re-run parity + resource_exchange audit |

Apply track rollback remains independent (Sprint 77 docs).

---

## 8. No API / payload drift

- No HTTP response shape changes.
- No JSON field ordering changes.
- No new production paths without explicit wiring sprint.
- Parity tests cluster-free.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_AUDIT.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CRITERIA.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_CRITERIA_M74.md`


---

## Sprint 83 (call-site wiring)

Sprint 83 wires all 32 production `resource_exchange_*` observation call sites to full-helper staged wrappers (8 CPU + 12 ready + 12 restart). `ResourceExchangeObservationWiredV1` and `ResourceExchangeObservationCandidateRuntimeUsedV1` remain **false**; effective runtime path is **legacy**. Direct candidate calls in production remain forbidden. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING.md` (Hyperdensity) and `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_M83.md` (Dashboard).


---

## Sprint 84 (candidate-runtime staging)

Sprint 84 sets `ResourceExchangeObservationCandidateRuntimeUsedV1=true` while `ResourceExchangeObservationWiredV1=false`. AND gate keeps effective runtime on legacy; candidate branch inactive. Production call-sites remain wrappers from Sprint 83. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_STAGING.md`.
