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

Re-audit may proceed only after criteria 1–9 satisfied for **remaining** functions; Sprint 52 satisfies **pure allowlist copy** only (3/46). Sprint 53 satisfies **production adapter shell** but not criterion 1 (full adapter wired).

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_DEFERRED.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_BOUNDARY.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REAUDIT_CRITERIA_M36.md`
