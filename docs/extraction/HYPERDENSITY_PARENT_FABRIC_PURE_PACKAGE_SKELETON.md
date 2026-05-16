# Hyperdensity Parent Fabric — pure package skeleton (Sprint 45–48)

## Purpose

Record the **`pkg/hyperdensity/parentfabric/...`** **stdlib-only skeleton** added in Sprint 45: package comments, version constants, and subpackage placeholders (`summary`, `governance`, `evidence`, `recommendation`). **No** Dashboard logic was copied or moved.

## What exists (Sprint 45)

| Path | Role |
|------|------|
| `pkg/hyperdensity/parentfabric/doc.go` | Root package comment |
| `pkg/hyperdensity/parentfabric/version.go` | `ParentFabricPackageVersion`, `ParentFabricRuntimeOwnership`, `ParentFabricExtractionMode` |
| `pkg/hyperdensity/parentfabric/version_test.go` | Asserts constants are set; ownership remains **dashboard** string literal |
| `pkg/hyperdensity/parentfabric/{summary,governance,evidence,recommendation}/doc.go` | Reserved subtrees (comments only) |
| `pkg/hyperdensity/parentfabric/executiontypes/` | **Sprint 46–47** — partial copy-contract + drift manifest |
| `pkg/hyperdensity/parentfabric/workload/` | **Sprint 48** — placeholder only (`copy-deferred` audit for workload helpers) |

## Rules (unchanged from Sprint 44)

- **Dashboard** remains **runtime owner** for Parent Fabric handlers and I/O.
- **Phase 3 (Sprint 46):** copy-contract in Hyperdensity only — **no** Dashboard production import of `parentfabric`.
- **No** new runtime import of Hyperdensity `parentfabric` in Dashboard production (Sprint 45–46).
- **API / JSON ordering / apply / execution** paths are **untouched**.

## Enforcement

- `scripts/validate_parentfabric_pure_deps.sh` — static grep deny list for forbidden strings under `pkg/hyperdensity/parentfabric` (wired from `scripts/validate.sh`).

## Related

- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_PHASES.md`
- `HYPERDENSITY_PARENT_FABRIC_DEPENDENCY_GUARDS.md`
- `HYPERDENSITY_PARENT_FABRIC_PURE_CANDIDATE_AUDIT.md`
- Dashboard `docs/hyperdensity/HYPERDENSITY_PARENT_FABRIC_EXTRACTION_STATUS_M30.md`
