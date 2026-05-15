# Repo Readiness Matrix

| Repo | Baseline Branch | Windows Branch | Status | Porting Decision | Merge Strategy | Blockers | Required Tests | Risk |
|---|---|---|---|---|---|---|---|---|
| `Karl-Dashboard` | `The-Father` (`5e8e601d`) | `The-Father-Windows` (`030216b3`) | Windows branch stale (`The-Father` ahead 17 commits) | **Do not port from Windows branch** | No merge, no cherry-pick from stale UI branch | stale branch, UI drift risk, forbidden 443/8888 touch | none in this milestone; later UI contract-read tests only | High |
| `Karl-Hyperdensity` | `main` (`094afd17`) | `The-Father-Windows` (`970ae03d`) | Windows branch contains complete proof/model stack; heavy diff footprint | **Selective backend-first porting** | New integration branch from `main`; cherry-pick/manual port by PR slices | artifact bloat, contract duplication risk, claim drift risk | `go test` targeted packages, fixture contract checks, safety assertions | Medium-High |
| `Karl-Inventory` | integration baseline (`67824b7`, mainline baseline to confirm) | `The-Father-Windows` (`123f72e1`) | Focused fluidShell witness delta | **Separate PR track, witness-only** | No merge; isolated PR(s) after Hyperdensity contract stabilization | branch-base alignment pending, generated/build outputs present in working tree | `dotnet test` for `KarlInventoryAgent.FluidShell.Tests`, config/schema validation | Medium |
| `Karl-OS-ISO` | integration baseline (`a17d7bd2`) | `The-Father-Windows` (`a17d7bd2`) | No effective delta in audited clone | **Defer entirely** | No merge, no porting in this milestone | packaging premature before backend controller merge | none now | Low-Medium |

## Hyperdensity Candidate Set Assessment
Candidate files requested for evaluation in `Karl-Hyperdensity` are present only on Windows branch and absent on `main`, including:
- `pkg/windowsfluidvirt/*` core product/contracts/evaluators
- `cmd/karl-node-fluid-actuator/main.go`
- `cmd/karl-fluid-compliance-replay/main.go`
- `cmd/karl-fluid-windows-executor/main.go`
- Windows contracts and fixtures under `docs/contracts/` and `examples/windows-fluid-*`

Integration suitability: strong for selective PR slicing; unsuitable for direct merge due to payload size and mixed readiness levels.
