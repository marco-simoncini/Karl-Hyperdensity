# Hyperdensity Parent Fabric — resource exchange observation migration boundary (Sprint 87)

## Summary

**Sprint 87** formally closes the **resource_exchange observation track (Sprint 78–86)** as **boundary complete**. No flag flips, no runtime changes, no production call-site changes. Broad observation remains **disabled**.

---

## 1. Scope

| Item | Sprint 87 |
|------|-----------|
| resource_exchange observation track | **Complete** |
| Flag changes | **None** |
| Runtime / API changes | **None** |
| Documentation + audit closure | **Yes** |

Completed chain: audit → shadow matrix CPU → local helper classification → local helper shadow matrix → staged wrappers → full-helper readiness → call-site wiring → candidate-runtime staging → activation readiness → activation → post-activation hardening → **migration boundary**.

---

## 2. Non-goals

- Broad observation (`ObservationWiredV1` / `ProductionWiredV1`).
- rollback, VM runtime, admission_guard wiring.
- Dashboard `pkg/hyperdensity/parentfabric` import.
- `workload_helpers.go` copy (remains **copy-deferred**).

---

## 3. Completed Sprint 78–86 chain

| Sprint | Milestone |
|--------|-----------|
| 78 | resource_exchange observation audit |
| 79 | shadow matrix (CPU) |
| 80 | staged wrappers |
| 81 | local helper classification + shadow matrix |
| 82 | full-helper staged wrappers + wiring readiness |
| 83 | production call-site wiring (8/12/12 wrappers) |
| 84 | candidate-runtime staging |
| 85 | activation readiness |
| 86 | `ResourceExchangeObservationWiredV1=true` activation |
| 87 | **boundary closure** (this document) |

---

## 4. Final resource_exchange state

| Flag / invariant | Value |
|------------------|-------|
| `ResourceExchangeObservationCandidateRuntimeUsedV1` | **true** |
| `ResourceExchangeObservationWiredV1` | **true** |
| Candidate branch (wrappers only) | **active** |
| Wrapper runtime path | **candidate** |
| Wrapper ≡ candidate ≡ legacy (24 cases) | **PASS** |
| Wrapper production counts | CPU **8**, ready **12**, restart **12** |
| Legacy production counts | **0/0/0** |
| Direct candidate production calls | **0** |
| `ObservationWiredV1` | **false** |
| `ProductionWiredV1` | **false** |

---

## 5. Why broad observation remains false

Apply and resource_exchange tracks completing does **not** authorize `ObservationWiredV1=true`. Each surface requires its own audit, shadow matrix, policy, and dedicated sprint chain. Sprint 87 explicitly records that completion gates do **not** cascade to broad observation.

---

## 6. Remaining observation surfaces

See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REMAINING_SURFACE_DECISION.md`:

- **rollback** — legacy, safety-critical
- **VM runtime** — legacy, high risk
- **admission_guard** — legacy, policy-critical
- **usage.go / other-review** — classification pending
- **broad observation** — explicitly disabled

---

## 7. Rollback posture

Resource_exchange rollback: set `ResourceExchangeObservationWiredV1=false`; re-run parity. No call-site restore required.

rollback **surface** (observed-state) remains a separate legacy track — not in scope for Sprint 87.

---

## 8. Risks

- Mistaking track completion for permission to flip broad observation.
- Premature rollback/VM/admission wiring without dedicated shadow matrices.
- Semantic drift in candidate helpers over time.

---

## 9. Next roadmap boundary

After resource_exchange closure, KHR must shift from indefinite micro-sprint adapter work toward **architecture / project memory** canonization — especially **storage** and **network/OVN** semantics. See `HYPERDENSITY_KHR_ROADMAP_TRANSITION_NOTE.md`.

---

## 10. Recommended next sprint

**Sprint 88 — KHR architecture memory**: consolidate storage primitives (EphemeralDisk, ephemeralOverlay, ephemeralClone, scratch, readonly, persistent, discardPolicy, promote-to-image) and network primitives (KARLNetwork, CellNetwork, ShellNetwork, NetworkAttachment, NetworkLease, NetworkPolicy) with OVN/SDN Dashboard compatibility mapping.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_POST_ACTIVATION_HARDENING.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REMAINING_SURFACE_DECISION.md`
- `HYPERDENSITY_KHR_ROADMAP_TRANSITION_NOTE.md`


---

## Sprint 88 (KHR architecture memory)

Sprint 88 canonizes KHR/KARL architecture memory, storage semantics, and network/OVN semantics. No runtime/adapter changes. KubeVirt remains compatibility provider and public-cloud fallback. See HYPERDENSITY_KHR_ARCHITECTURE_MEMORY.md and related Sprint 88 docs.


---

## Sprint 89 (ResourceLease minimal contract)

Sprint 89 adds ResourceLease minimal contract sketch (storage/network/provider/examples). No CRD, no controller, no runtime. See HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md and related Sprint 89 docs.
