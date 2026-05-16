# Hyperdensity Parent Fabric — workload pure candidates contract (Sprint 52)

## Summary

Sprint 52 copies **exactly three** stdlib-only functions from Dashboard `hyperdensity_parent_fabric_workload_helpers.go` into `pkg/hyperdensity/parentfabric/workload`:

| Dashboard function | Hyperdensity export |
|--------------------|---------------------|
| `hyperdensityAppsWorkloadResource` | `AppsWorkloadResource(kind) (resource string, ok bool)` |
| `hyperdensityPilotWorkloadTerm` | `PilotWorkloadTerm(kind) (term string, ok bool)` |
| `hyperdensityExecutionSupportsLiveApplyKind` | `ExecutionSupportsLiveApplyKind(kind) bool` |

**Full file verdict:** **`copy-deferred`** (unchanged).

## Not copied (Dashboard-owned)

| Category | Count | Status |
|----------|------:|--------|
| `api_path_builders` | 6 | Deferred |
| `observed_state_builders` | 21 | Deferred |
| `execution_apply_helpers` | 16 | Deferred |

## Artifacts

| Artifact | Path |
|----------|------|
| Implementation | `pkg/hyperdensity/parentfabric/workload/pure_candidates.go` |
| Golden contract | `workload/testdata/workload_pure_candidates_contract.golden.json` |
| Source manifest | `workload/testdata/workload_pure_candidates_source_manifest.golden.json` |
| Dashboard parity | `hyperdensity_parent_fabric_workload_pure_candidates_contract_test.go` |

## Rules

- **Stdlib only** — no K8s/KubeVirt imports or path literals.
- **No** Dashboard production import of `pkg/hyperdensity/parentfabric`.
- **No** API response / runtime behavior change.
- Semantics match Dashboard source **exactly** for the three functions.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_DEFERRED.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REAUDIT_CRITERIA.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_PURE_CANDIDATES_CONTRACT_M38.md`
