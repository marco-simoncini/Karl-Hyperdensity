# Hyperdensity M1–M7 extraction readiness (Sprint 24)

## Summary

After M1–M7, Karl-Hyperdensity holds **real, tested** `pkg/hyperdensity/blockers` and `pkg/hyperdensity/contracts` helpers plus redacted goldens. Karl-Dashboard consumes them **only in `pkg/server` tests** via a pinned module version. This is **extraction scaffolding**, not production wiring.

## Ready for next micro-step (code)

- **Continue test-only import path** already proven: Dashboard tests import `blockers` and `contracts` without touching handlers.
- **Shared test-only fixture suite:** committed JSON under `kubernetes-console/pkg/server/testdata/` plus Hyperdensity `testdata/dashboard/` goldens — safe to extend with new edges under the same allowlist/policy helpers.
- **Optional package boundary:** a future `parentfabric/summary` (or similar) subpackage in Hyperdensity could hold only DTO + validators to reduce **module-wide import** surface — still **no handler import** until a deliberate later sprint.

## Not ready

- **Runtime import in Dashboard handlers** — would change authoritative API assembly, couple release trains, and blur claim-safe mapping vs live `HyperdensityParentFabricSummaryV1` JSON.
- **KHR as execution arm** — apply, Windows enablement, KubeVirt removal, or new Parent Fabric are out of scope.
- **Replacing Dashboard summary transport** — contract JSON is not a drop-in for native `?view=summary` responses.

## Recommended next step

1. **Package split (optional):** extract `blockers` + `contracts` into a minimal Go module or submodule with **`go 1.16`/`1.22` alignment doc** so Dashboard is not forced to `1.22` by unrelated Hyperdensity packages (e.g. `windowsfluidvirt` using `any`).
2. **Shared fixture suite:** single manifest listing all Dashboard testdata paths + Hyperdensity goldens (still test-only).
3. **CI guard:** Dashboard `test_hyperdensity_parity.sh` in GitHub Actions (Sprint 24) — no cluster, no npm.

## Risks

| Risk | Mitigation |
|------|------------|
| **Go 1.22 coupling** | Dashboard `kubernetes-console/go.mod` at `go 1.22` while importing full Karl-Hyperdensity module | Target minimal module or `replace` local during dev; document in M3/M7 anchors |
| **Module-wide import** | `go get` pulls entire Hyperdensity tree | Split `pkg/hyperdensity/{blockers,contracts}` module when ready |
| **Dashboard runtime divergence** | Live API evolves; fixtures stale | Manual redacted capture script + env-guarded golden refresh; parity CI on every KHR PR touching console tests |
| **Claim drift** | `supportsApply` vs `applyAllowed` | Keep `ValidateApplySemantics` + matrix M2/M3 in parity suite |

## Matrix reference

See `HYPERDENSITY_PARITY_MATRIX_M1_M7.md` for milestone-by-milestone artifacts and guarantees.
