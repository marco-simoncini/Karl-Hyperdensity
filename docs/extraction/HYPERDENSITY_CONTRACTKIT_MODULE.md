# Hyperdensity contractkit module (Sprint 25)

## Purpose

Minimal **Go submodule** so Karl-Dashboard (and other consumers) can import **only** parent-fabric parity helpers — `blockers` + `contracts` — without pulling the full Karl-Hyperdensity module (e.g. `pkg/windowsfluidvirt`, KHR agent, CRD tooling).

## Module path

```text
github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit
```

Root: `pkg/hyperdensity/contractkit/go.mod`

## Packages included

| Package | Role |
|---------|------|
| `contractkit/blockers` | Gate/blocker ID catalog (`Known`, `Severity`, `Catalog`) |
| `contractkit/contracts` | `ParentFabricSummary` DTO, mapping, golden JSON, fixture policy validators |
| `contractkit/claimpolicy` | Posture vocabulary (**Sprint 35**) + minimal claim-policy catalog (**Sprint 36**: `ClaimPolicyID`, `Catalog`, `ForbiddenProductionClaimIDs`, …); Dashboard **runtime** must not import |

Implementation lives here; `pkg/hyperdensity/blockers` and `pkg/hyperdensity/contracts` at repo root are **thin re-export aliases** for in-repo compatibility.

## What contractkit does NOT include

- `pkg/windowsfluidvirt` (uses Go 1.18+ `any`; not imported)
- KHR Linux agent, telemetry, runtime provider
- Kubernetes clients, HTTP servers, CRD apply
- Dashboard handlers or Parent Fabric runtime

## Go version

- **contractkit:** `go 1.18` — stdlib only; no `any`; no generics requirement beyond what stdlib uses.
- **Root Karl-Hyperdensity module:** remains `go 1.22` for `windowsfluidvirt` and other packages.
- **Why not 1.16:** root repo and unrelated packages need 1.18+; contractkit targets **1.18** as the lowest practical floor for consumers without coupling to root `1.22`.

## Coupling reduction

Before Sprint 25, Dashboard `go.mod` required `github.com/marco-simoncini/Karl-Hyperdensity` (entire module) → forced **Go 1.22** alignment and large transitive surface.

After Sprint 25, Dashboard parity tests require **only** `.../contractkit` → smaller module graph; Dashboard may keep `go 1.22` for other console deps but is **no longer tied** to Hyperdensity root module version.

## Import policy

- **Dashboard:** test-only imports from `contractkit/blockers` and `contractkit/contracts` (`*_test.go` only).
- **Runtime:** no handler or API import of contractkit until a deliberate later extraction sprint.

## Validation

```bash
./scripts/validate.sh   # includes (cd pkg/hyperdensity/contractkit && go test ./...)
```

## Related

- M1–M7 matrix: `HYPERDENSITY_PARITY_MATRIX_M1_M7.md`
- Dashboard import doc: `docs/hyperdensity/HYPERDENSITY_CONTRACTKIT_IMPORT_M8.md` (Karl-Dashboard)
