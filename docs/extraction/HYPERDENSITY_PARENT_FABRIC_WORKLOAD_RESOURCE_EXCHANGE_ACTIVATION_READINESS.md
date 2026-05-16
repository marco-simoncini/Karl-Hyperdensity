# Hyperdensity Parent Fabric — resource exchange activation readiness (Sprint 85)

## Summary

**Sprint 85** is **readiness-only**: formal certification that a future dedicated sprint may flip `ResourceExchangeObservationWiredV1 = true`. **No flag changes, no runtime changes, no production call-site changes** in Sprint 85.

---

## 1. Scope

| Item | Sprint 85 |
|------|-----------|
| Activation readiness documentation + tests | **Yes** |
| `ResourceExchangeObservationWiredV1 = true` | **No** |
| Runtime / production changes | **No** |
| Broad `ObservationWiredV1` | **No** |

---

## 2. Non-goals

- Activation flip (reserved for Sprint 86+ when explicitly approved).
- Changing `ResourceExchangeObservationCandidateRuntimeUsedV1` (stays **true** from Sprint 84).
- Direct candidate calls in production.
- rollback, VM runtime, admission_guard wiring.
- Dashboard `parentfabric` import.
- Storage/network primitive implementation.

---

## 3. Preconditions from Sprint 78–84

| Sprint | Deliverable |
|--------|-------------|
| 78–82 | Candidates, wrappers, shadow matrices |
| 83 | 32 production call-sites wired to wrappers (8/12/12) |
| 84 | `CandidateRuntimeUsedV1 = true`; AND gate; runtime still **legacy** |

---

## 4. Current flag state

| Flag | Sprint 85 |
|------|-----------|
| `ResourceExchangeObservationCandidateRuntimeUsedV1` | **`true`** |
| `ResourceExchangeObservationWiredV1` | **`false`** |
| `ObservationWiredV1` | **`false`** |
| `ProductionWiredV1` | **`false`** |
| Candidate branch active | **no** (AND gate) |
| Effective wrapper path | **legacy** |

---

## 5. Activation target

Future activation sets **`ResourceExchangeObservationWiredV1 = true`** with **`CandidateRuntimeUsedV1` already true**.

AND gate: `Wired && CandidateUsed` → **candidate** branch at runtime.

---

## 6. Required activation invariants

- Production uses **wrappers only** (no direct candidate calls).
- Full-helper counts: wrapper **8/12/12**, legacy **0/0/0**.
- **wrapper ≡ candidate ≡ legacy** (24 shadow cases) before and expected after activation.
- `ObservationWiredV1` and `ProductionWiredV1` remain **false**.
- No apply observation helpers in `resource_exchange_*`.

---

## 7. Required post-activation hardening

Dedicated sprint after activation must include post-activation hardening tests and audit updates (mirror apply track Sprint 75–77).

---

## 8. Rollback plan

Minimum rollback after future activation:

1. Set `ResourceExchangeObservationWiredV1 = false`.
2. Re-run parity.
3. No production call-site changes required.

---

## 9. Risks

See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION_RISKS.md`.

---

## 10. Recommended next sprint

**Sprint 86:** activation flip (`ResourceExchangeObservationWiredV1 = true`) only, if explicitly approved — plus post-activation hardening.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION_RISKS.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_STAGING.md`
