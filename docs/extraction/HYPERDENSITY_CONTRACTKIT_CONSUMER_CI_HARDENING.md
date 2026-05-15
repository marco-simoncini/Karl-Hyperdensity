# Hyperdensity contractkit — consumer CI hardening (Sprint 41)

## Purpose

This guide documents how **CI and private forks** should fetch the nested Go module  
`github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit` **without** changing any HTTP API, Parent Fabric runtime behavior, JSON ordering, or execution paths. It is **documentation + environment variables** only.

## Nested module identity

| Item | Value |
|------|--------|
| **Module path** | `github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit` |
| **Current stable (module semver)** | `v0.1.9-khr-m1-m19` |
| **Git tag on parent repo** | `pkg/hyperdensity/contractkit/v0.1.9-khr-m1-m19` |

The tag **must** include the subdirectory prefix (`pkg/hyperdensity/contractkit/…`); Go resolves nested modules from the **monorepo root** using that tag shape.

## Superseded module versions (do not pin for new work)

| Version | Notes |
|---------|--------|
| `v0.1.5-khr-m1-m16` | Superseded by `v0.1.6-khr-m1-m17` lineage. |
| `v0.1.7-khr-m1-m18` | Superseded by `v0.1.8-khr-m1-m18`; **avoid** due to proxy/cache immutability hazards (see `HYPERDENSITY_CONTRACTKIT_NO_REPUBLISH_POLICY.md`). |

`v0.1.8-khr-m1-m18` remains **historically valid** (Sprint 39) but is **not** the current pin after Sprint 40; consumers should stay on **`v0.1.9-khr-m1-m19`** unless a newer sprint explicitly bumps the pin.

## Private repo / CI: `GOPRIVATE` and checksum DB

For `github.com/marco-simoncini/*` modules, configure the Go toolchain so the **sum verifier** and **module proxy** do not block or mis-resolve private tags:

```bash
export GOPRIVATE=github.com/marco-simoncini/*
export GONOSUMDB=github.com/marco-simoncini/*
```

- **`GOPRIVATE`** — modules under this prefix are treated as private (short module path resolution, different default proxy/sumdb behavior).
- **`GONOSUMDB`** — do not consult `sum.golang.org` for checksum lines for these module paths (avoids **404 / unknown revision** when the public sum DB has not yet seen a fresh tag, or for private visibility).

### Optional: bypass proxy for the same prefix

If your organization’s proxy cannot see the GitHub repo or lags tag ingestion:

```bash
export GONOPROXY=github.com/marco-simoncini/*
```

With **`GONOPROXY`**, Go may **fetch git/vcs directly** from origin (subject to network credentials on the runner). Use when the corporate proxy is the bottleneck; prefer **`GOPRIVATE` + `GONOSUMDB`** first.

## When `GOSUMDB=off` is acceptable

Setting **`GOSUMDB=off`** disables checksum verification **globally** for that invocation. Use it **only** as a last resort when:

- the CI job cannot use **`GONOSUMDB`** with a scoped prefix (older tooling), or  
- a misconfigured runner still hits sumdb errors after **`GOPRIVATE` / `GONOSUMDB`** are set.

**Prefer** `GOPRIVATE` and `GONOSUMDB` scoped to `github.com/marco-simoncini/*` before turning off the sum database entirely.

## Why not republish or repoint tags

Published **module@version** content is treated as **immutable** by the ecosystem (public proxy, downstream `go.sum`). **Never**:

- delete and recreate a released tag for the same semver, or  
- move a tag (`git tag -f`) to a different commit while keeping the same version string.

If a release is wrong, publish a **strictly newer** module semver, push a **new** `pkg/hyperdensity/contractkit/v…` tag, bump `ContractKitModuleVersion` in Karl-Hyperdensity, and update consumer pins. See **`HYPERDENSITY_CONTRACTKIT_NO_REPUBLISH_POLICY.md`**.

## Karl-Dashboard alignment

- Parity workflow sets **`GOPRIVATE`** / **`GONOSUMDB`** (and **`GONOPROXY`** where useful) on the Go step.  
- `scripts/hyperdensity/audit_contractkit_module_pin.sh` enforces the **exact** `go.mod` pin and rejects superseded pins and **pseudo-versions** for this module.

## Related

- `HYPERDENSITY_CONTRACTKIT_NO_REPUBLISH_POLICY.md`
- `HYPERDENSITY_CONTRACTKIT_RELEASE_TAGGING.md`
- `HYPERDENSITY_CONTRACTKIT_VERSION_MODEL.md`
- `HYPERDENSITY_CONTRACTKIT_CONSUMER_AUDIT.md`
- `HYPERDENSITY_CONTRACTKIT_CONSUMER_POLICY.md`
- `templates/audit_contractkit_module_pin.sh`
- Dashboard `docs/hyperdensity/HYPERDENSITY_CONTRACTKIT_CONSUMER_CI_M24.md`
