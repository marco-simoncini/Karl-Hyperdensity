# Hyperdensity Parent Fabric — execution types contract (Sprint 46)

## Scope

Sprint 46 opens **Phase 3** in a **minimal, controlled** way:

- **Copy-contract only** in `pkg/hyperdensity/parentfabric/executiontypes`
- **Golden file:** `pkg/hyperdensity/parentfabric/executiontypes/testdata/execution_types_contract.golden.json`
- **Dashboard** remains **runtime owner** — source file untouched, **no** production import of `parentfabric`

## What the contract covers

| Artifact | Purpose |
|----------|---------|
| `ExecutionTypesPackageVersion` | `v0.0.0-sprint46` |
| `HyperdensityExecutionSummary` | Zero-value JSON shape + field set (21 fields) |
| `HyperdensityExecutionEngineSpine` | Four top-level engine fields only |
| `ContractDocument` | Metadata + deferred note + source stats |

## What the contract does **not** cover

- Full `HyperdensityExecutionEngine` (150 top-level fields)
- Nested VM / Windows / apply executor surfaces
- Handler or API response behavior

## Enforcement

- `go test ./pkg/hyperdensity/parentfabric/executiontypes/...`
- `scripts/validate_parentfabric_pure_deps.sh` (includes `executiontypes/`)
- Dashboard test-only: `TestHyperdensityParentFabricExecutionTypesContract` (local golden + source invariants; **no** Hyperdensity root module import)

## Drift

Update Hyperdensity copy + golden when Dashboard changes copied DTOs. Do **not** widen Dashboard runtime imports without an allowlist sprint.

## Related

- `HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_AUDIT.md`
- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_PHASES.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_CONTRACT_M31.md`
