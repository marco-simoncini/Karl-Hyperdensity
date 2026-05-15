# Hyperdensity contractkit — release tagging (Sprint 27–28)

## Submodule model

`github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit` is a **nested Go module** inside Karl-Hyperdensity. Consumers pin via **semver git tag** (preferred after Sprint 28) or pseudo-version.

## First anchor tag (applied Sprint 28)

**Go module version:** `v0.1.0-khr-m1-m9`  
**Git tag on parent repo (required prefix for nested modules):**

```text
pkg/hyperdensity/contractkit/v0.1.0-khr-m1-m9
```

**Commit:** `c03ef68c939a42349688a28600e4a4531413f44b` (Sprint 27 case-ID helpers + manifest)

### Commands used

```bash
git fetch --all --prune
git checkout KHR
git pull --ff-only
git tag pkg/hyperdensity/contractkit/v0.1.0-khr-m1-m9 c03ef68c939a42349688a28600e4a4531413f44b
git push origin pkg/hyperdensity/contractkit/v0.1.0-khr-m1-m9
```

Dashboard pin:

```bash
go get github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit@v0.1.0-khr-m1-m9
```

### Note on tag prefix

Go requires the git tag name to include the **subdirectory path** to the module root (`pkg/hyperdensity/contractkit/...`). A short alias tag alone is **not** resolved by `go get`.

## Second anchor tag (Sprint 30 — version model + clean semver)

**Go module version:** `v0.1.1-khr-m1-m12`  
**Git tag:**

```text
pkg/hyperdensity/contractkit/v0.1.1-khr-m1-m12
```

**Base:** Sprint 29 release version model (`aaeafba`) + `ContractKitModuleVersion` alignment commit.

```bash
git tag pkg/hyperdensity/contractkit/v0.1.1-khr-m1-m12 <sprint-30-commit>
git push origin pkg/hyperdensity/contractkit/v0.1.1-khr-m1-m12
```

Dashboard pin (replaces pseudo-version `v0.1.0-khr-m1-m9.0.20260515175117-aaeafba...`):

```bash
go get github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit@v0.1.1-khr-m1-m12
```

## Third anchor tag (Sprint 35 — claimpolicy package, module bump only)

**Go module version:** `v0.1.2-khr-m1-m13`  
**Git tag:**

```text
pkg/hyperdensity/contractkit/v0.1.2-khr-m1-m13
```

Adds `pkg/hyperdensity/contractkit/claimpolicy` (stdlib vocabulary; **no** schema / manifest / API change). Dashboard may bump `go.mod` for test-only parity; runtime import freeze unchanged (M17).

```bash
git tag pkg/hyperdensity/contractkit/v0.1.2-khr-m1-m13 <sprint-35-commit>
git push origin pkg/hyperdensity/contractkit/v0.1.2-khr-m1-m13
```

Dashboard pin:

```bash
go get github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit@v0.1.2-khr-m1-m13
```

## Compatibility promise (contractkit)

| Property | Commitment |
|----------|------------|
| Runtime | **None** — test/extraction helpers only |
| Dependencies | **Stdlib only** |
| Go version | **1.18** minimum |
| Breaking changes | Only with a **new tag** / bumped `ContractKitVersion` + manifest sync |
| Windows / apply | Anchors remain disabled / claim-safe (M1–M8) |

## Related

- `HYPERDENSITY_CONTRACTKIT_MODULE.md`
- `HYPERDENSITY_CONTRACTKIT_FIXTURE_MANIFEST.md`
- Dashboard `HYPERDENSITY_CONTRACTKIT_RELEASE_M11.md`
