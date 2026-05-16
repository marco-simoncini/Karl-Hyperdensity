# Hyperdensity Parent Fabric — pure candidate audit (repo-wide) (Sprint 45)

## Purpose

Describe how **Phase 1–2** (mechanical file audit + pure-helper identification) is executed **without** moving runtime code from Karl-Dashboard. Hyperdensity owns the **target packages** and **dependency guards**; Dashboard owns the **authoritative file list** and **CSV inventory**.

## Hyperdensity actions (Sprint 45)

1. **`pkg/hyperdensity/parentfabric/...` skeleton** — stdlib-only placeholders (see **`HYPERDENSITY_PARENT_FABRIC_PURE_PACKAGE_SKELETON.md`**).
2. **`scripts/validate_parentfabric_pure_deps.sh`** — fail CI if forbidden import strings appear under `parentfabric/`.
3. **Sprint 46:** `hyperdensity_parent_fabric_execution_types.go` → partial copy in `parentfabric/executiontypes` (see audit/contract docs). Dashboard source file **not** moved; production import **not** added.

## How to audit candidates (Phase 2)

| Step | Owner | Output |
|------|-------|--------|
| Enumerate `hyperdensity_parent_fabric*.go` | Dashboard | `list_parent_fabric_runtime_files.sh`, `parent_fabric_runtime_inventory_m29.csv` |
| Tag `category_guess` / `readiness_guess` | Dashboard script (heuristic) | CSV columns |
| Confirm “pure” with `go list -deps` / manual review | Future sprint | Per-file ADR or checklist row |

**PASS for Phase 2 (docs):** every non-test file appears in the CSV **or** the CSV header/doc states an explicit subset with rationale (Sprint 45 uses **full** inventory).

## Target package mapping (reminder)

| Heuristic bucket | Likely future `pkg/hyperdensity/parentfabric/...` |
|------------------|---------------------------------------------------|
| Summary surfaces | `summary` |
| Policy / matrix / ledger | `governance` |
| Evidence / collector / substrate | `evidence` |
| Action slate / futures | `recommendation` |
| Cross-cutting types / helpers | `parentfabric` (root) |

## Related

- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_PHASES.md`
- `HYPERDENSITY_PARENT_FABRIC_DEPENDENCY_GUARDS.md`
- Dashboard `docs/hyperdensity/HYPERDENSITY_PARENT_FABRIC_PURE_CANDIDATE_AUDIT_M29.md`
- Dashboard `docs/hyperdensity/parent_fabric_runtime_inventory_m29.csv`
