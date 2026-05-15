# Hyperdensity contractkit — reusable consumer policy (Sprint 42)

## Scope

Normative checklist for **any** Go module that imports  
`github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit` **as an external dependency** (semver tag from GitHub). **Karl-Hyperdensity** itself uses a **`replace`** to the in-tree nested module; this policy primarily targets **downstream repos** (today: **Karl-Dashboard** as **reference consumer**).

**No** HTTP API, runtime handler, JSON ordering, or execution-path changes are mandated by this document.

## Current stable module version

| Item | Value |
|------|--------|
| **Module path** | `github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit` |
| **Pinned semver (stable)** | `v0.1.9-khr-m1-m19` |
| **Git tag (parent repo)** | `pkg/hyperdensity/contractkit/v0.1.9-khr-m1-m19` |

## Recommended CI environment

```bash
export GOPRIVATE=github.com/marco-simoncini/*
export GONOSUMDB=github.com/marco-simoncini/*
# Optional when proxy lags or cannot see private tags:
export GONOPROXY=github.com/marco-simoncini/*
```

Prefer scoped **`GONOSUMDB`** over global **`GOSUMDB=off`**. Rationale: `HYPERDENSITY_CONTRACTKIT_CONSUMER_CI_HARDENING.md`.

## Forbidden for stable downstream consumers

| Rule | Rationale |
|------|-----------|
| **Pseudo-version** pins (`v0.0.0-…`, `vX.Y.Z-0.<time>-<commit>`) for `contractkit` | Non-reproducible drift; bypasses release hygiene. |
| **Semver** `v0.1.5-khr-m1-m16`, `v0.1.7-khr-m1-m18` | **Superseded** — see `contracts.ContractKitSupersededModuleVersions` and `HYPERDENSITY_CONTRACTKIT_NO_REPUBLISH_POLICY.md`. |
| **Tag repoint / delete-recreate** for a published `pkg/hyperdensity/contractkit/v…` | Breaks Go module immutability expectations. |

Downstream consumers **may** additionally reject pins older than current stable (e.g. previous `v0.1.8-khr-m1-m18`) via project-specific audit scripts.

## Import surface policy (declare per consumer)

| Consumer class | `contractkit/blockers` | `contractkit/contracts` | `contractkit/claimpolicy` |
|----------------|------------------------|-------------------------|---------------------------|
| **Karl-Dashboard production** (`pkg/server` non-test) | **Allowed** (only subpath today) | **Forbidden** | **Forbidden** |
| **Karl-Dashboard tests** | Allowed | Allowed (`*_test.go`) | Allowed (`*_test.go`) |
| **Other repos** | Must document **runtime vs test-only** in their own ADR; default **test-only** unless an explicit sprint expands runtime. | Same | Same |

## Required checks (downstream)

1. **`go.mod` exact pin** — one `require` line for the module path at the **current stable** semver (no pseudo-version).
2. **No superseded pins** — at minimum block `v0.1.5-khr-m1-m16` and `v0.1.7-khr-m1-m18`.
3. **No pseudo-version** for this module path.
4. **`go mod verify`** — run in CI where compatible (checksum verification for `go.sum`).
5. **Parity / audit scripts** — Dashboard-like consumers should run `audit_contractkit_module_pin.sh`, runtime import audit, and parity tests; copy/adapt from **Karl-Dashboard** `kubernetes-console/scripts/hyperdensity/`.

## Template (copy, do not submodule)

Generic bash template (not wired into other repos by default):

- `docs/extraction/templates/audit_contractkit_module_pin.sh`

Update **`EXPECTED_CONTRACTKIT_VERSION`** only in a **named sprint** when Hyperdensity bumps `ContractKitModuleVersion` and publishes a new tag.

## Second consumer adoption (Sprint 43)

Before adding a new downstream `go.mod` consumer, follow **`HYPERDENSITY_CONTRACTKIT_SECOND_CONSUMER_PLAYBOOK.md`** and file a completed **`templates/CONTRACTKIT_CONSUMER_DECISION_RECORD.md`** in the consumer repo.

## Related

- `HYPERDENSITY_CONTRACTKIT_SECOND_CONSUMER_PLAYBOOK.md`
- `templates/CONTRACTKIT_CONSUMER_DECISION_RECORD.md`
- `HYPERDENSITY_CONTRACTKIT_CONSUMER_AUDIT.md`
- `HYPERDENSITY_CONTRACTKIT_CONSUMER_CI_HARDENING.md`
- `HYPERDENSITY_CONTRACTKIT_NO_REPUBLISH_POLICY.md`
- `HYPERDENSITY_CONTRACTKIT_RELEASE_TAGGING.md`
- `HYPERDENSITY_CONTRACTKIT_VERSION_MODEL.md`
- `templates/audit_contractkit_module_pin.sh`
- Dashboard `docs/hyperdensity/HYPERDENSITY_CONTRACTKIT_REFERENCE_CONSUMER_M25.md`
- Dashboard `docs/hyperdensity/HYPERDENSITY_CONTRACTKIT_REFERENCE_CONSUMER_CHECKLIST_M26.md`
