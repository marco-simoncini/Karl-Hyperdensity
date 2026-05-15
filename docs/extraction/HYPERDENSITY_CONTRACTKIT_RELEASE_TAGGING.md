# Hyperdensity contractkit — release tagging (Sprint 27)

## Submodule model

`github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit` is a **nested Go module** inside Karl-Hyperdensity. Consumers (Karl-Dashboard parity tests) pin it via `go.mod` pseudo-version or an optional **git tag** on the parent repo.

## Recommended first anchor tag

```text
contractkit/v0.1.0-khr-m1-m9
```

Naming: `<subdir-module>/v<semver>-khr-<milestone-range>` — compatible with Go module tags for subdirectories in a monorepo.

## Proposed commands (manual — not run automatically in sprint)

Replace `<commit>` with the contractkit-ready commit (e.g. Sprint 26 manifest + Sprint 27 case-ID helpers):

```bash
git tag contractkit/v0.1.0-khr-m1-m9 <commit>
git push origin contractkit/v0.1.0-khr-m1-m9
```

Dashboard may then pin:

```bash
go get github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit@v0.1.0-khr-m1-m9
```

## Policy note

**Do not tag automatically** in CI or sprint commits unless repo release policy explicitly allows it. **Pseudo-versions remain acceptable** until semver tags are approved (current KHR workflow).

## Compatibility promise (contractkit)

| Property | Commitment |
|----------|------------|
| Runtime | **None** — test/extraction helpers only |
| Dependencies | **Stdlib only** |
| Go version | **1.18** minimum |
| Breaking changes | Only with a **new tag** / bumped `ContractKitVersion` + manifest sync |
| Windows / apply | Anchors remain disabled / claim-safe (M1–M8) |

## When to tag

- After M1–M9 parity matrix, manifest, and Dashboard consumer checks are green on `KHR`.
- Before promoting contractkit beyond test-only imports (still not recommended for handlers without a dedicated sprint).

## Related

- `HYPERDENSITY_CONTRACTKIT_MODULE.md`
- `HYPERDENSITY_CONTRACTKIT_FIXTURE_MANIFEST.md`
- Dashboard `HYPERDENSITY_CONTRACTKIT_RELEASE_M10.md`
