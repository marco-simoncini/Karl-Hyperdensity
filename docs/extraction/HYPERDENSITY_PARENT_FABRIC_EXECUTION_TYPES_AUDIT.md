# Hyperdensity Parent Fabric — execution types source audit (Sprint 46–47)

## Source file

| | |
|--|--|
| **Dashboard path** | `pkg/server/hyperdensity_parent_fabric_execution_types.go` |
| **Package** | `server` |
| **Lines** | ~4571 |
| **`type` definitions** | **152** |
| **Functions** | **0** |

## Imports (source)

```go
import "time"
```

No `k8s.io/*`, `kubevirt.io/*`, `net/http`, `github.com/gorilla/*`, or Dashboard-external console imports in this file.

`time.Time` appears in a **small subset** of nested types (e.g. action-slate timestamps). Those structs are **not** copied in Sprint 46.

## Coupling notes

| Dependency | Detail |
|------------|--------|
| **Same package types** | `HyperdensityExecutionEngine` references **~150** nested surface types defined in this file |
| **Sibling file types** | `HyperdensityCPUQuantity` / `HyperdensityMemoryQuantity` live in `hyperdensity_parent_fabric.go` (duplicated in Hyperdensity `executiontypes` for summary contract) |
| **Runtime wiring** | Types are consumed by `hyperdensity_parent_fabric_execution.go` and handlers — **not** moved in Sprint 46 |

## Pure candidates copied (Sprint 46)

| Type | Rationale |
|------|-----------|
| `HyperdensityCPUQuantity` | Primitive DTO (from sibling file; required by summary) |
| `HyperdensityMemoryQuantity` | Primitive DTO (from sibling file; required by summary) |
| `HyperdensityExecutionSummary` | Pure counters / quantities + summary string |
| `HyperdensityExecutionEngineSpine` | Contract subset: `summary`, `supportsApply`, `supportedSurfaces`, `applyNotes` only |

## Not copied (deferred)

- Full `HyperdensityExecutionEngine` (**150** top-level JSON fields)
- Remaining **~148** nested surface structs (Linux shell, VM lane, Windows, admission, policy pack, …)
- Types using `time.Time` fields until a dedicated adapter sprint

Sample deferred type names (first 20 of 152):  
`HyperdensityFleetEquilibriumCandidateV1`, `HyperdensityFleetEquilibriumOnboardingV1`, `HyperdensityShellFactoryProfileV1`, `HyperdensityShellFactoryV1`, `HyperdensityShellClaimGeneratorV1`, `HyperdensityReleaseSupportMatrixV1`, `HyperdensityEvidenceBundleDemoScenarioPackV1`, `HyperdensityLiveResourceAuthorityV1`, `HyperdensityAdmissionGuardV1`, `HyperdensityPolicyPackV1`, `HyperdensityResourceEquilibriumV1`, `HyperdensityActionSlateV1`, …

## Decision

| | |
|--|--|
| **Mode** | **copy-contract / partial** — `pkg/hyperdensity/parentfabric/executiontypes` |
| **Runtime wiring** | **None** — Dashboard does **not** import `parentfabric/executiontypes` in production |
| **Equivalence claim** | **Spine + summary DTO only** — not full engine parity |

## Drift risk

- Dashboard may add fields to `HyperdensityExecutionSummary` or engine spine fields → golden contract test fails until Hyperdensity copy is updated.
- Quantity types in `hyperdensity_parent_fabric.go` could diverge from `executiontypes` duplicate.

## Rollback

- Revert Hyperdensity `executiontypes` package and docs; remove Dashboard contract test from parity runner.
- Dashboard source file unchanged — zero runtime rollback.

## Sprint 47 — drift guard (no new copy)

- Hyperdensity `SourceManifest` documents import set, type count, and json tags for the copied slice.
- Manifest does **not** read Dashboard at runtime; Dashboard AST test validates source locally.

## Related

- `HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_DRIFT_GUARD.md`
- `HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_CONTRACT.md`
- `pkg/hyperdensity/parentfabric/executiontypes/`
- Dashboard `docs/hyperdensity/HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_SOURCE_AUDIT_M31.md`
