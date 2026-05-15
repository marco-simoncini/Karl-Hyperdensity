# contractkit / claimpolicy — boundary + catalog (Sprint 35–36)

## Purpose

`pkg/hyperdensity/contractkit/claimpolicy` is a **stdlib-only** package for:

- **Sprint 35:** stable **posture** vocabulary (`PostureKind`, `KnownPosture`, `Postures`).
- **Sprint 36:** a minimal **claim-policy catalog** (`ClaimPolicyID`, `ClaimPolicyRule`, `Catalog`, `Known`, `Severity`, `RuntimeAllowed`, `MustKeepRuntimeDisabled`, `ForbiddenProductionClaimIDs`).

It complements `blockers` (M1 gate IDs) and `contracts` (DTOs / manifest) **without** changing API payloads, JSON ordering, or Parent Fabric runtime behavior. **There is no runtime enforcement in this package** — documentation and test parity only.

## Non-goals (explicit)

- **No** Karl-Dashboard **production runtime** import of `claimpolicy` (M17 freeze unchanged; parity tests may import in `*_test.go` only).
- **No** `contractkit/contracts` in Dashboard runtime (unchanged).
- **No** execution/apply path changes, Windows enablement, KubeVirt removal, or new Parent Fabric collectors.
- **No** change to `ContractKitVersion` (`v0.0.0-sprint26`) or manifest envelope (`hyperdensity.parity.manifest/v1`) for this slice — module semver tag bumps only.

## Sprint 35 vs Sprint 36

| Sprint | Delivered |
|--------|-----------|
| **35** | Boundary: `PackageVersion`, posture tokens, stable `Postures()` order. |
| **36** | Catalog: typed `ClaimPolicyID` constants, `ClaimPolicyRule` rows, lookup helpers, `ForbiddenProductionClaimIDs` (sorted), critical `MustKeepRuntimeDisabled` anchors. |

## Catalog semantics (high level)

- **Stable order:** `Catalog()` returns rules sorted **lexicographically by `ID`**.
- **Windows / apply:** Claims `no_windows_hyperdensity_apply` and `windows_lane_disabled` document **planning-only / disabled** posture; they do **not** enable Windows execution.
- **Autonomous apply:** `no_autonomous_apply` forbids unattended apply narratives on the Hyperdensity surface.
- **Production mutation:** `no_production_mutation` aligns with existing safety posture language.
- **KubeVirt:** `kubevirt_legacy_provider` is a **compatibility / legacy marker** (`RuntimeAllowed` true for documentation class); `no_generic_kubevirt_replacement` forbids implying a **generic non-KubeVirt replacement** — distinct IDs, distinct semantics.
- **contracts import:** `no_runtime_contracts_import` documents the Dashboard **M17** rule (parity / future gates), not a new importer.

## Package API (epoch `PackageVersion`)

| Area | Symbols |
|------|---------|
| Epoch | `PackageVersion` (`v0.0.0-sprint36`) |
| Posture (Sprint 35) | `PostureKind`, constants, `KnownPosture`, `Postures` |
| Catalog (Sprint 36) | `ClaimPolicyID`, `ClaimPolicyRule`, `Catalog`, `Known`, `Severity`, `RuntimeAllowed`, `MustKeepRuntimeDisabled`, `ForbiddenProductionClaimIDs` |

## Relationship to other contractkit packages

| Package | Role |
|---------|------|
| `blockers` | M1 gate / blocker ID catalog. |
| `contracts` | Summary DTO, manifest, golden helpers. |
| `claimpolicy` | Claim-policy + posture vocabulary for parity and **future** KHR apply / planning gates (no runtime wiring in Sprint 36). |

## Validation

```bash
( cd pkg/hyperdensity/contractkit && go test ./claimpolicy -count=1 )
./scripts/validate.sh
```

## Consumer pin

Dashboard bumps nested module `go.mod` when tag `pkg/hyperdensity/contractkit/v0.1.3-khr-m1-m14` (or newer) is published; **test-only** imports of `claimpolicy` remain in `*_test.go` only.

## Related

- `HYPERDENSITY_CONTRACTKIT_MODULE.md`
- `HYPERDENSITY_CONTRACTKIT_VERSION_MODEL.md`
- `HYPERDENSITY_CONTRACTKIT_RELEASE_TAGGING.md`
- `Karl-Dashboard/docs/hyperdensity/HYPERDENSITY_CONTRACTKIT_CLAIMPOLICY_M18.md`
- `Karl-Dashboard/docs/hyperdensity/HYPERDENSITY_CLAIMPOLICY_CATALOG_M19.md`
