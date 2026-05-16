# Hyperdensity Parent Fabric — resource exchange candidate-runtime revalidation (Sprint 84)

## Summary

After Sprint 84 sets `ResourceExchangeObservationCandidateRuntimeUsedV1=true`, shadow re-validation confirms **wrapper ≡ candidate ≡ legacy** for all **24** full-helper cases. Production remains on **wrappers only**; candidate branch stays **inactive** until `ResourceExchangeObservationWiredV1=true`.

---

## Matrix

| Helper | Cases |
|--------|-------|
| CPU | **8** |
| ready | **8** |
| restart | **8** |
| **Total** | **24** |

---

## Assertions

- `wrapper == candidate == legacy` for every case.
- Production `resource_exchange_*` uses wrappers, **not** direct candidate helpers.
- `ResourceExchangeObservationWiredV1 = false` → candidate branch **inactive**.
- `ResourceExchangeObservationCandidateRuntimeUsedV1 = true` does **not** change effective runtime path.

---

## Next activation sprint

Flipping `ResourceExchangeObservationWiredV1=true` requires **separate** sprint approval, post-flip hardening, and parity — not bundled with Sprint 84.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_STAGING.md`


---

## Sprint 85 (activation readiness)

Sprint 85 is readiness-only for `ResourceExchangeObservationWiredV1=true`. No flag changes. Sprint 86 may execute activation flip if approved. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION_READINESS.md`.


---

## Sprint 86 (resource_exchange activation)

Sprint 86 sets ResourceExchangeObservationWiredV1=true. Candidate branch active in resource_exchange wrappers only. ObservationWiredV1/ProductionWiredV1 remain false. See ACTIVATION.md and POST_ACTIVATION_HARDENING.md.
