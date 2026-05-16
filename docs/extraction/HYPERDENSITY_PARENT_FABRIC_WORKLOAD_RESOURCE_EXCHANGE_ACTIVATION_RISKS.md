# Hyperdensity Parent Fabric — resource exchange activation risks (Sprint 85)

## Summary

Risks associated with a future flip of `ResourceExchangeObservationWiredV1 = true`. Sprint 85 documents and guards; it does **not** execute the flip.

---

## Risks

| Risk | Mitigation |
|------|------------|
| Treating resource_exchange activation as **broad observation** | `ObservationWiredV1` / `ProductionWiredV1` stay false; dedicated sprint scope |
| Changing **candidate helper semantics** without shadow re-validation | Forbidden in activation sprint; 24-case matrix required |
| **Direct candidate** calls in production `resource_exchange_*` | Policy guards + audit; wrappers only |
| **CPU-only** partial rollback or wiring | Full-helper 8/12/12 enforced |
| Accidental touch of **rollback / VM / admission** | Out-of-scope guards in risk test + audit |
| Confusing `CandidateRuntimeUsedV1=true` with **candidate branch active** | AND gate requires **both** flags; Sprint 85 documents inactive branch |
| Expecting API/payload change on activation | wrapper ≡ candidate ≡ legacy → behavior unchanged |

---

## Minimum rollback

1. Set `ResourceExchangeObservationWiredV1 = false` (in candidate/adapter flag file per project convention).
2. Re-run `test_hyperdensity_parity.sh`.
3. **No** production call-site restoration required (wrappers remain).

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION_READINESS.md`


---

## Sprint 86 (resource_exchange activation)

Sprint 86 sets ResourceExchangeObservationWiredV1=true. Candidate branch active in resource_exchange wrappers only. ObservationWiredV1/ProductionWiredV1 remain false. See ACTIVATION.md and POST_ACTIVATION_HARDENING.md.


---

## Sprint 87 (resource_exchange boundary closure)

Sprint 87 closes resource_exchange observation Sprint 78–86 as boundary complete. No flag/runtime changes. Broad observation remains false. Next phase: KHR architecture memory and storage/network semantics. See MIGRATION_BOUNDARY.md, REMAINING_SURFACE_DECISION.md, KHR_ROADMAP_TRANSITION_NOTE.md.
