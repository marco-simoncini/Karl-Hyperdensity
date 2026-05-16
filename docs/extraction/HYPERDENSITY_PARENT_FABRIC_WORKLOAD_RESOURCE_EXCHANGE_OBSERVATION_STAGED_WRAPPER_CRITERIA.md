# Hyperdensity Parent Fabric — resource exchange staged wrapper criteria (Sprint 80+)

Criteria for sprints **after** Sprint 80 staged wrappers, before production wiring.

---

## 1. Dedicated sprint per milestone

| Milestone | Requirement |
|-----------|-------------|
| Call-site wiring readiness | Dedicated sprint |
| Local helper shadow (optional) | Dedicated sprint before or with wiring |
| `CandidateRuntimeUsedV1` flip | Dedicated sprint |
| `WiredV1` activation | Dedicated sprint after staging |

---

## 2. Policy gates

| Gate | Required |
|------|----------|
| `ObservationWiredV1` | **false** until broad policy satisfied |
| `ProductionWiredV1` | **false** |
| No rollback / VM / admission wiring | **Yes** |
| wrapper == candidate == legacy | **Yes** (8+ cases) |
| No direct candidate in production | **Yes** |
| No wrapper in production until readiness PASS | **Yes** |
| Apply track unchanged | **Yes** |

---

## 3. Local ready/restart

Before `WiredV1=true`:

- Explicit decision: wire CPU only vs include ready/restart wrappers.
- If CPU-only: document acceptance of partial observation adapter coverage.
- If full: extend shadow matrix + staged wrappers for local helpers.

---

## 4. Rollback strategy

| Step | Action |
|------|--------|
| 1 | `ResourceExchangeObservationWiredV1 = false` |
| 2 | Restore legacy calls in resource_exchange (if wired) |
| 3 | `CandidateRuntimeUsedV1 = false` optional |
| 4 | Re-run parity + audit |

---

## 5. No API / payload drift

No HTTP/JSON changes. No cluster dependency for parity tests.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_STAGED_WRAPPERS.md`
