# Hyperdensity Parent Fabric — broad observation decision (Sprint 64)

## Decision (recommended)

**`hyperdensityWorkloadAdapterObservationWiredV1` remains `false`.**

This is **deliberate policy**, not incomplete work by mistake.

---

## Rationale

| Approach | Verdict |
|----------|---------|
| Granular subflags (`PathWired`, `PilotObservationWired`, `LiveObservationWired`, future `ApplyObservationWired`, etc.) | **Preferred** |
| Broad `ObservationWiredV1 = true` when all surfaces migrated + shadowed | **Future only** |
| Broad flag permanently false as semantic guard | **Acceptable long-term** |

---

## Preconditions before broad flip (if ever)

1. Zero legacy observation call sites in forbidden categories (apply, resource_exchange, rollback, VM runtime).
2. Per-surface shadow tests PASS.
3. Dedicated sprint with explicit golden `broadObservationRecommended: true`.
4. Parity runner green.
5. Rollback documented.

---

## Sprint 64–65 status

| Item | Value |
|------|-------|
| `broadObservationRecommended` | **`false`** |
| `ObservationWiredV1` | **`false`** |
| `ProductionWiredV1` | **`false`** |
| `ApplyObservationWiredV1` | **`false`** (Sprint 65 placeholder) |
| `apply.go` wiring | **None** (Sprint 65–66 proposal + shadow only) |
| Apply shadow ready | **Yes** (Sprint 66) |
| Apply staged wrappers | **Yes** (Sprint 67; not used by `apply.go`) |
| Apply wrapper hardening | **Yes** (Sprint 68; 8×4 matrix) |
| Apply wiring readiness | **Yes** (Sprint 69) |
| Apply call-site wiring | **Yes** (Sprint 70; flags **false**) |
| Apply post-wiring hardening | **Yes** (Sprint 71) |
| Apply flip criteria | **Yes** (Sprint 72; docs-only) |
| Apply candidate-runtime readiness | **Yes** (Sprint 73; docs-only) |
| Apply post-activation hardening | **Yes** (Sprint 76) |
| Apply migration boundary complete | **Yes** (Sprint 77) |
| Resource exchange audit complete | **Yes** (Sprint 78, audit-only) |
| Resource exchange wired | **No** |
| Resource exchange shadow matrix (Sprint 79) | **Yes** (candidate only) |
| Resource exchange staged wrappers (Sprint 80) | **Yes** (not used in production) |
| Local helper shadow matrix (Sprint 81) | **Yes** (candidates only) |

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_OBSERVATION_REAUDIT.md`


---

## Sprint 83 (call-site wiring)

Sprint 83 wires all 32 production `resource_exchange_*` observation call sites to full-helper staged wrappers (8 CPU + 12 ready + 12 restart). `ResourceExchangeObservationWiredV1` and `ResourceExchangeObservationCandidateRuntimeUsedV1` remain **false**; effective runtime path is **legacy**. Direct candidate calls in production remain forbidden. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING.md` (Hyperdensity) and `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_M83.md` (Dashboard).


---

## Sprint 84 (candidate-runtime staging)

Sprint 84 sets `ResourceExchangeObservationCandidateRuntimeUsedV1=true` while `ResourceExchangeObservationWiredV1=false`. AND gate keeps effective runtime on legacy; candidate branch inactive. Production call-sites remain wrappers from Sprint 83. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_STAGING.md`.


---

## Sprint 85 (activation readiness)

Sprint 85 is readiness-only for `ResourceExchangeObservationWiredV1=true`. No flag changes. Sprint 86 may execute activation flip if approved. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION_READINESS.md`.


---

## Sprint 86 (resource_exchange activation)

Sprint 86 sets ResourceExchangeObservationWiredV1=true. Candidate branch active in resource_exchange wrappers only. ObservationWiredV1/ProductionWiredV1 remain false. See ACTIVATION.md and POST_ACTIVATION_HARDENING.md.


---

## Sprint 87 (resource_exchange boundary closure)

Sprint 87 closes resource_exchange observation Sprint 78–86 as boundary complete. No flag/runtime changes. Broad observation remains false. Next phase: KHR architecture memory and storage/network semantics. See MIGRATION_BOUNDARY.md, REMAINING_SURFACE_DECISION.md, KHR_ROADMAP_TRANSITION_NOTE.md.
