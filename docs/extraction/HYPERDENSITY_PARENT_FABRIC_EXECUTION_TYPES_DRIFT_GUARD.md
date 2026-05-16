# Hyperdensity Parent Fabric — execution types drift guard (Sprint 47)

## Purpose

Strengthen the Sprint 46 **copy-contract** for `pkg/hyperdensity/parentfabric/executiontypes` with a **static source manifest** and tests — **without** copying new types, **without** Dashboard production imports, and **without** reading Dashboard files from Hyperdensity at runtime.

## Sprint 47 scope

| In scope | Out of scope |
|----------|----------------|
| `SourceManifest` + golden `execution_types_source_manifest.golden.json` | New helper copies (`workload_helpers.go`, …) |
| Reflect/json-tag validation in Hyperdensity tests | Dashboard → `parentfabric` production import |
| Dashboard AST drift test (`go/parser`) | API / handler / apply / execution behavior changes |

## SourceManifest (Hyperdensity)

Package: `pkg/hyperdensity/parentfabric/executiontypes`

| Field | Value (Sprint 47) |
|-------|-------------------|
| `sourceFile` | `pkg/server/hyperdensity_parent_fabric_execution_types.go` |
| `sourceImportSet` | `["time"]` |
| `sourceTypeDefinitionCount` | `152` |
| `executionSummaryFieldCount` | `21` |
| `executionSummaryJsonTags` | 21 stable json keys (see manifest golden) |
| `engineSpineJsonTags` | `summary`, `supportsApply`, `supportedSurfaces`, `applyNotes` |
| `copiedTypes` | CPU/Memory quantity, ExecutionSummary, ExecutionEngineSpine |

Functions: `DefaultSourceManifest()`, `ExecutionSummaryJSONTags()`, `EngineSpineJSONTags()`, `ValidateSourceManifest()`.

## Dashboard verification (local)

Dashboard does **not** import Hyperdensity root. It parses `hyperdensity_parent_fabric_execution_types.go` with **`go/ast`** and compares to a **mirrored** manifest golden under `pkg/server/testdata/`.

## Enforcement

- Hyperdensity: `go test ./pkg/hyperdensity/parentfabric/executiontypes/...` + `scripts/validate.sh`
- Dashboard: `TestHyperdensityParentFabricExecutionTypesDrift` in `test_hyperdensity_parity.sh`

## Sprint 48 note

`executiontypes` remains the **first stable** copy-contract. `workload_helpers.go` was audited separately and classified **`copy-deferred`** — see workload audit docs.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_AUDIT.md`
- `HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_CONTRACT.md`
- `HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_AUDIT.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_DRIFT_GUARD_M32.md`
