# Hyperdensity contractkit — version model (Sprint 29)

Three **independent** version layers; do not conflate them in manifests, tags, or tests.

| Layer | Constant / field | Current value | Purpose |
|-------|------------------|---------------|---------|
| **Module semver** | `ContractKitModuleVersion` / `go.mod` | `v0.1.9-khr-m1-m19` | Go module release; git tag `pkg/hyperdensity/contractkit/v0.1.9-khr-m1-m19` |
| **Contract schema** | `ContractKitVersion` / manifest `contractKitVersion` | `v0.0.0-sprint26` | Logical DTO/validator epoch; bump when mapping rules or contract shape changes |
| **Manifest envelope** | `FixtureManifestVersion` / manifest `manifestVersion` | `hyperdensity.parity.manifest/v1` | JSON manifest file format |

## API

```go
info := contracts.ReleaseInfo()
stable := contracts.CurrentStableReleaseInfo()
// info.ModuleVersion, GitTag, SchemaVersion, FixtureManifestVersion
// Sprint 40: stable aliases ContractKitCurrentStableModuleVersion / superseded list guards
```

`ValidateFixtureManifest` enforces:

- `manifestVersion == FixtureManifestVersion`
- `contractKitVersion == ContractKitVersion` (schema)

It does **not** embed module semver — Dashboard `go.mod` owns module pin.

## Rules

- **Test-only** — no runtime handler import.
- Module tag bump **without** schema change: update `ContractKitModuleVersion` + git tag; schema/manifest fields unchanged.
- Schema bump: update `ContractKitVersion`, manifests, and parity tests; new module tag recommended.
- Manifest format bump: update `FixtureManifestVersion` + new manifest files.

## No republish / no repoint (Sprint 40)

Go **module version strings** are effectively **immutable** after publication to the public module proxy or adoption in `go.sum`.

- **Do not** repoint or delete/recreate a released **`pkg/hyperdensity/contractkit/v…`** tag for the same semver.
- If a release must be corrected, publish a **strictly newer** semver and update pins.
- Superseded semvers: **`v0.1.5-khr-m1-m16`** (superseded by `v0.1.6`), **`v0.1.7-khr-m1-m18`** (superseded by `v0.1.8` lineage); see **`HYPERDENSITY_CONTRACTKIT_NO_REPUBLISH_POLICY.md`** and `contracts.ContractKitSupersededModuleVersions`.

## Related

- `HYPERDENSITY_CONTRACTKIT_RELEASE_TAGGING.md`
- `HYPERDENSITY_CONTRACTKIT_NO_REPUBLISH_POLICY.md`
- Dashboard `HYPERDENSITY_CONTRACTKIT_VERSION_M12.md`
