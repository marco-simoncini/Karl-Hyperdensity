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

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_OBSERVATION_REAUDIT.md`
