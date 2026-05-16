# Hyperdensity Parent Fabric — workload helpers re-audit criteria (Sprint 50)

## When `workload_helpers.go` may be re-audited

All of the following must be true before changing verdict from **`copy-deferred`**:

| # | Criterion |
|---|-----------|
| 1 | **API path builders** isolated in a documented Dashboard **adapter** (not in Hyperdensity pure-core). |
| 2 | **Observed-state builders** classified **runtime-bound** and covered by `WorkloadObservationAdapter` (doc + tests). |
| 3 | **Candidate functions** narrowed to explicit **pure allowlist** (currently 3 kind/mode helpers). |
| 4 | **`parentfabric/primitives`** stable — golden tests green in Hyperdensity `validate.sh`. |
| 5 | **`executiontypes`** drift manifest green (Dashboard AST test). |
| 6 | **Classification fixture** complete — every function in source file categorized (Sprint 50). |
| 7 | **Dashboard adapter classification test** PASS in parity runner. |
| 8 | **Dashboard adapter stub test** PASS (`TestHyperdensityParentFabricWorkloadAdapterStub`, Sprint 51). |
| 9 | **No production import** of `pkg/hyperdensity/parentfabric` until an explicit **wiring sprint** approves it. |

## What re-audit does **not** mean

- Automatic **`copy-approved`** for the full file.
- Moving KubeVirt/K8s path strings into Hyperdensity.
- Changing API responses, JSON ordering, or apply behavior.

## Sprint 50–51 outcome

| Sprint | Delivered | Verdict |
|--------|-----------|---------|
| **50** | Criteria documented + classification fixture/test | `copy-deferred` |
| **51** | Dashboard test-only adapter stubs + golden manifest | `copy-deferred` |
| **52** | Three pure-candidates in `parentfabric/workload` + golden + Dashboard parity | `copy-deferred` (full file) |
| **53** | Dashboard production-internal adapter v1 + tests — **not wired** | `copy-deferred` (full file) |
| **54** | Shadow tests legacy vs adapter v1 — **not wired** | `copy-deferred` (full file) |
| **55** | Wiring proposal + call-site inventory (51 sites) | `copy-deferred` (full file) |
| **56** | Path-only wiring (6 non-apply files) — `PathWiredV1=true` | `copy-deferred` (full file) |
| **57** | Pilot-only observation — `PilotObservationWiredV1=true`, `ObservationWiredV1=false` | `copy-deferred` (full file) |
| **58** | Pilot observation hardening + live proposal (no new wiring) | `copy-deferred` (full file) |
| **59** | Live observation staged in `live.go`; `LiveObservationWiredV1=false` | `copy-deferred` (full file) |
| **60** | Live wrapper shadow hardening (7 cases); flip not allowed | `copy-deferred` (full file) |
| **61** | `LiveObservationWiredV1=true`; true branch legacy-equivalent | `copy-deferred` (full file) |
| **62** | Semantic live candidate shadow; `CandidateRuntimeUsedV1=false` | `copy-deferred` (full file) |
| **63** | Branch swap: wrapper true branch → candidate; scoped to live | `copy-deferred` (full file) |
| **64** | Observation re-audit; `ObservationWiredV1=false` policy | `copy-deferred` (full file) |
| **65** | Apply-observation proposal + criteria; no `apply.go` wiring | `copy-deferred` (full file) |
| **66** | Apply-observation shadow matrix; candidate not runtime-used | `copy-deferred` (full file) |
| **67** | Apply-observation staged wrappers; `apply.go` still legacy | `copy-deferred` (full file) |
| **68** | Apply wrapper hardening 8×4; no `apply.go` wiring | `copy-deferred` (full file) |
| **69** | Apply wiring readiness | `copy-deferred` (full file) |
| **70** | Apply call-site wiring; flags false (legacy path) | `copy-deferred` (full file) |
| **71** | Apply post-wiring hardening 8×4 | `copy-deferred` (full file) |
| **72** | Apply flip criteria + risks (docs-only; no flip) | `copy-deferred` (full file) |
| **73** | Candidate-runtime readiness + branch logic (no flip) | `copy-deferred` (full file) |
| **74** | Candidate-runtime staging flip (`CandidateUsed=true`, Wired false) | `copy-deferred` (full file) |
| **75** | Apply observation activation (`Wired=true`, candidate branch active) | `copy-deferred` (full file) |
| **76** | Post-activation hardening 8×4 | `copy-deferred` (full file) |
| **77** | Migration boundary; apply track complete | `copy-deferred` (full file) |
| **78** | Resource exchange observation audit (8 listed call sites; no wiring) | `copy-deferred` (full file) |
| **79** | Resource exchange shadow matrix (candidate parity; no wiring) | `copy-deferred` (full file) |
| **80** | Resource exchange staged wrappers (wrapper parity; no production wiring) | `copy-deferred` (full file) |
| **81** | Local helper shadow matrix (ready/restart candidates; no production wiring) | `copy-deferred` (full file) |
| **82** | Full-helper staged wrappers + call-site wiring readiness | `copy-deferred` (full file) |

Re-audit for remaining functions unchanged. Sprint 57–82: apply track complete; full-helper readiness certified; **no** production wiring. **`ObservationWiredV1=false` is deliberate** (Sprint 64–78).

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_DEFERRED.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_BOUNDARY.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REAUDIT_CRITERIA_M36.md`


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
