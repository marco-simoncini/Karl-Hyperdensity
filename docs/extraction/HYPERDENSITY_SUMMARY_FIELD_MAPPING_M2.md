# Parent-fabric summary field mapping — apply semantics (M2 / Sprint 18)

This document clarifies how **Dashboard runtime JSON** relates to the **Hyperdensity M1 contract** (`ParentFabricSummary`). It does **not** change Dashboard handlers or API responses.

---

## `supportsApply` (Dashboard) vs `applyAllowed` (Hyperdensity contract)

| Field | Where | Meaning |
|-------|--------|---------|
| **`executionEngine.supportsApply`** | Dashboard `HyperdensityExecutionEngineSummaryV1` (`?view=summary` and full snapshot) | **Technical capability flag**: the execution engine surface can represent or route apply-style operations when policy and gates allow. Does **not** alone mean “production apply is on” or “autonomous apply is allowed”. |
| **`executionEngine.applyAllowed`** | Hyperdensity `ParentFabricSummary` (M1 golden / contracts package) | **Product / contract claim**: whether the **exported contract view** asserts that apply is allowed as a general product posture. M1 golden uses **`false`** to mean “this anchor does not authorize apply as a broad product claim”. |

### Critical rule (Sprint 18)

> **`applyAllowed: false` in the M1 golden does NOT mean Dashboard has no technical apply path.**

Dashboard may still expose POST `/api/hyperdensity/parent-fabric/execution`, dry-run categories, operator-controlled guest-assisted modes, and `supportsApply: true` in summary while **gates** block production/autonomous/broad execution.

The M1 contract is a **reduced, claim-safe projection** for extraction tests — not a full mirror of Dashboard execution posture.

---

## Apply posture vocabulary

| Term | Typical Dashboard signal | Hyperdensity contract interpretation |
|------|--------------------------|--------------------------------------|
| **Runtime technically supports apply** | `supportsApply: true`, execution routes exist, dry-run validated | Not modeled as `applyAllowed: true` unless contract explicitly claims it |
| **Product allows apply now** | Production auto enablement, wave certification passed | Would require `applyAllowed: true` **and** `recommendationOnly: false` in a future contract version (not M1) |
| **Operator controlled apply** | Guest-assisted modes, auth on POST execution, approval objects | Maps to `hyperdensity.operatorControlled: true`; compatible with `applyAllowed: false` in M1 |
| **Dry-run only** | Execution category `dry_run_only`, `dryRunSupported` | Maps to `dryRunSupported: true`; aligns with blocker `dry_run_only` |
| **Autonomous apply forbidden** | Auto policy disabled, `unsupported_broad_automation`, production auto gates | `autonomousMode: false`; blocker catalog includes `unsupported_broad_automation` |

---

## Related Dashboard fields (reference)

- `executionEngine.summary` — human-readable execution summary (category often `dry_run_only`).
- `hyperdensityProductionAutoPolicy` / enablement plan — product gates; not copied into M1 minimal struct.
- Gate contributions: `no_production_mutation`, `no_windows_lane` — see `pkg/hyperdensity/blockers`.

---

## Mapping table (M1 subset)

| Hyperdensity `ParentFabricSummary` | Dashboard analogue |
|-----------------------------------|-------------------|
| `executionEngine.applyAllowed` | **Not** equal to `!supportsApply`; contract claim only |
| `executionEngine.dryRunSupported` | Dry-run category + certification posture |
| `executionEngine.mode` | e.g. `operator_controlled` vs wave/autonomous labels in enablement surfaces |
| `hyperdensity.recommendationOnly` | Readonly/planning surfaces; no production GA claim |
| `hyperdensity.operatorControlled` | Auth + approval + guest-assisted operator paths |
| `windowsLane.enabled` | `windowsLaneEnabled` on VM collector surfaces (expected false) |

---

## Next steps

- **Sprint 18:** Dashboard compile-only test imports `pkg/hyperdensity/blockers` for ID parity.
- **Future M3:** Optional golden generator from redacted `?view=summary` with explicit field mapping function (Dashboard → contract).
