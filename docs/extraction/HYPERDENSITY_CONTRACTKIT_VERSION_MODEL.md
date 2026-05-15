# Hyperdensity contractkit — version model (Sprint 29)

Three **independent** version layers; do not conflate them in manifests, tags, or tests.

| Layer | Constant / field | Current value | Purpose |
|-------|------------------|---------------|---------|
| **Module semver** | `ContractKitModuleVersion` / `go.mod` | `v0.1.2-khr-m1-m13` | Go module release; git tag `pkg/hyperdensity/contractkit/v0.1.2-khr-m1-m13` |
| **Contract schema** | `ContractKitVersion` / manifest `contractKitVersion` | `v0.0.0-sprint26` | Logical DTO/validator epoch; bump when mapping rules or contract shape changes |
| **Manifest envelope** | `FixtureManifestVersion` / manifest `manifestVersion` | `hyperdensity.parity.manifest/v1` | JSON manifest file format |

## API

```go
info := contracts.ReleaseInfo()
// info.ModuleVersion, GitTag, SchemaVersion, FixtureManifestVersion
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

## Related

- `HYPERDENSITY_CONTRACTKIT_RELEASE_TAGGING.md`
- Dashboard `HYPERDENSITY_CONTRACTKIT_VERSION_M12.md`
