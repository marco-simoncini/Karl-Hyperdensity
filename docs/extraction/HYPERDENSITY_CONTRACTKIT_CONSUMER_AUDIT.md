# Hyperdensity contractkit — cross-consumer audit (Sprint 42)

## Purpose

Read-only inventory of **local clones** (when present) to see which Karl repositories **depend on** the nested Go module  
`github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit`. **No** API, Parent Fabric runtime, JSON ordering, or execution behavior changes are implied.

**Inspection date:** 2026-05-15 (authoritative for this revision).

## Method

- Enumerate `go.mod` files under each repo path with ripgrep / glob.
- Search for the substring `hyperdensity/contractkit` in all `go.mod` files and in `*.go` where needed.
- Record the **currently checked-out Git branch** at inspection time (may differ from `KHR`).

## Repositories explicitly requested

| Repository | Local path inspected | Branch @ audit | `go.mod` present | Imports `…/contractkit` (module root) | `…/contractkit/blockers` | `…/contractkit/contracts` | `…/contractkit/claimpolicy` | Pin policy (semver tag) | Status |
|------------|----------------------|----------------|------------------|----------------------------------------|---------------------------|----------------------------|-----------------------------|-------------------------|--------|
| **marco-simoncini/Karl-Dashboard** | `/home/m.simoncini/GitHub/Karl-Dashboard` | `KHR` | Yes (`Karl-Dashboard-dashboard/kubernetes-console/go.mod`) | Yes (require pin) | Yes (runtime + tests) | Yes (`*_test.go` only) | Yes (`*_test.go` only) | **Required** — exact pin `v0.1.9-khr-m1-m19` | **Consumer attivo** — **reference consumer** |
| **marco-simoncini/Karl-Hyperdensity** | `/home/m.simoncini/GitHub/Karl-Hyperdensity` | `KHR` | Yes (root + `pkg/hyperdensity/contractkit/go.mod`) | Yes (root `require` + **`replace`** → `./pkg/hyperdensity/contractkit`) | Yes (in-tree) | Yes (in-tree) | Yes (in-tree) | **N/A** — nested module is **authored here**; external consumers use tags | **Publisher / monorepo** (not a downstream semver consumer) |
| **marco-simoncini/Karl-OS-ISO** | `/home/m.simoncini/GitHub/Karl-OS-ISO` | `KHR` | No | No | No | No | No | No | **No consumer** |
| **marco-simoncini/Karl-Installer** | `/home/m.simoncini/GitHub/Karl-Installer` | `main` | Yes (root) | No | No | No | No | No | **No consumer** |
| **marco-simoncini/Karl-Inventory** | `/home/m.simoncini/GitHub/Karl-Inventory` | `KHR` | Yes (`inventory/collector/go.mod`) | No | No | No | No | No | **No consumer** |
| **marco-simoncini/Karl-Warden** | `/home/m.simoncini/GitHub/Karl-Warden` | `integration/identity-access-from-oidc` | Yes (root + `cmd/*/.docker-context/go.mod`) | No | No | No | No | No | **No consumer** |
| **marco-simoncini/FluidVirt** | `/home/m.simoncini/GitHub/FluidVirt` | `KHR` | Yes (root + staging modules) | No | No | No | No | No | **No consumer** |

## Additional Karl-named repos (quick `go.mod` scan)

Ripgrep across `~/GitHub/**/go.mod` for `hyperdensity/contractkit` on 2026-05-15 found **no** additional matches beyond **Karl-Dashboard** and **Karl-Hyperdensity** (root replace). Other Karl-prefixed directories under `/home/m.simoncini/GitHub` (e.g. `Karl-Genesi`, `Karl-Migration-Factory`, `KARL-APP`, `karl-directoryservice`) were not listed with a `go.mod` dependency on this nested module in that scan.

## Conclusions

1. **Karl-Dashboard** is the only inspected **downstream semver consumer** of the published nested module; treat it as the **reference consumer** for pins, CI env, and audit scripts (`HYPERDENSITY_CONTRACTKIT_REFERENCE_CONSUMER_M25.md` on Dashboard).
2. **Karl-Hyperdensity** **hosts** the module; consumers must follow **`HYPERDENSITY_CONTRACTKIT_CONSUMER_POLICY.md`** and **`HYPERDENSITY_CONTRACTKIT_CONSUMER_CI_HARDENING.md`**.
3. Repositories without a local clone, or on machines not listed above, should be recorded as **non ispezionabile / non presente localmente** in future audit revisions.

## Related

- `HYPERDENSITY_CONTRACTKIT_CONSUMER_POLICY.md`
- `HYPERDENSITY_CONTRACTKIT_CONSUMER_CI_HARDENING.md`
- `HYPERDENSITY_CONTRACTKIT_NO_REPUBLISH_POLICY.md`
- Dashboard `docs/hyperdensity/HYPERDENSITY_CONTRACTKIT_REFERENCE_CONSUMER_M25.md`
