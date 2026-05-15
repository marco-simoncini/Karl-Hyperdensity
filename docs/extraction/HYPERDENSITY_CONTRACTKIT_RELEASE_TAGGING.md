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

## Fourth anchor tag (Sprint 36 — claimpolicy catalog)

**Go module version:** `v0.1.3-khr-m1-m14`  
**Git tag:**

```text
pkg/hyperdensity/contractkit/v0.1.3-khr-m1-m14
```

Completes `claimpolicy` as a **minimal claim-policy catalog** (Sprint 36): `Catalog`, `Known`, `ForbiddenProductionClaimIDs`, etc. **No** `ContractKitVersion` / manifest bump; Dashboard **test-only** parity expands; runtime import freeze unchanged (M17).

```bash
git tag pkg/hyperdensity/contractkit/v0.1.3-khr-m1-m14 <sprint-36-commit>
git push origin pkg/hyperdensity/contractkit/v0.1.3-khr-m1-m14
```

Dashboard pin:

```bash
go get github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit@v0.1.3-khr-m1-m14
```

## Fifth anchor tag (Sprint 37 — claimpolicy surface mapping)

**Go module version:** `v0.1.4-khr-m1-m15`  
**Git tag:**

```text
pkg/hyperdensity/contractkit/v0.1.4-khr-m1-m15
```

Adds **Parent Fabric surface mapping** for `claimpolicy` (`SurfaceMappings`, `ValidateSurfaceMappings`). **No** `ContractKitVersion` / manifest bump; Dashboard **test-only** parity; M17 runtime import freeze unchanged.

```bash
git tag pkg/hyperdensity/contractkit/v0.1.4-khr-m1-m15 <sprint-37-commit>
git push origin pkg/hyperdensity/contractkit/v0.1.4-khr-m1-m15
```

Dashboard pin:

```bash
go get github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit@v0.1.4-khr-m1-m15
```

## Sixth anchor tag (Sprint 38 — claimpolicy Dashboard file traceability)

**Go module version:** `v0.1.6-khr-m1-m17`  
**Git tag:**

```text
pkg/hyperdensity/contractkit/v0.1.6-khr-m1-m17
```

Adds **`DashboardFiles`** traceability on `SurfaceClaimMapping` rows (`ValidateDashboardFileTraceability`), including **corrected** `windows_lane_disabled` file anchors. **No** `ContractKitVersion` / manifest bump; Dashboard **test-only** parity; M17 runtime import freeze unchanged.

**Note:** `v0.1.5-khr-m1-m16` tagged the first traceability slice but listed one path without the `windows_lane_disabled` token; consumers should pin **`v0.1.6-khr-m1-m17`** or newer.

```bash
git tag pkg/hyperdensity/contractkit/v0.1.6-khr-m1-m17 <sprint-38-commit>
git push origin pkg/hyperdensity/contractkit/v0.1.6-khr-m1-m17
```

Dashboard pin:

```bash
go get github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit@v0.1.6-khr-m1-m17
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
