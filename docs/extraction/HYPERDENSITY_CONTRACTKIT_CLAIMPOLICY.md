# contractkit / claimpolicy — boundary design (Sprint 35)

## Purpose

`pkg/hyperdensity/contractkit/claimpolicy` is a **stdlib-only** package that captures **stable posture / claim vocabulary** for Parent Fabric governance surfaces. It complements `blockers` (gate IDs) and `contracts` (DTOs) without coupling to Dashboard JSON field names in this sprint.

## Non-goals (Sprint 35)

- **No** Karl-Dashboard **runtime** import of `claimpolicy` (Parent Fabric handlers unchanged).
- **No** change to execution/apply semantics, API payloads, or JSON ordering.
- **No** new Parent Fabric surfaces; no Windows enablement; no KubeVirt removal.

## Package API (epoch `PackageVersion`)

| Symbol | Role |
|--------|------|
| `PackageVersion` | Design epoch string (`v0.0.0-sprint35`). |
| `PostureKind` | String alias for posture tokens. |
| `PostureEvidenceNamespace`, `PostureVisibilityOnly`, `PostureOperatorControlled` | Canonical posture literals. |
| `KnownPosture` | Membership test for defined postures. |
| `Postures` | Deterministic enumeration for tests/docs. |

## Relationship to other contractkit packages

| Package | Role |
|---------|------|
| `blockers` | M1 gate / blocker ID catalog. |
| `contracts` | Summary DTO, manifest, golden helpers (**test-only** on Dashboard today). |
| `claimpolicy` | Claim / policy posture vocabulary for **future** parity mapping (Sprint 35+). |

## Validation

```bash
( cd pkg/hyperdensity/contractkit && go test ./claimpolicy -count=1 )
./scripts/validate.sh
```

## Consumer pin

Dashboard parity tests may import `.../contractkit/claimpolicy` **only in `*_test.go`**, after the contractkit module tag that contains this package is published and `go.mod` is bumped (same flow as prior contractkit semver bumps).

## Related

- `HYPERDENSITY_CONTRACTKIT_MODULE.md`
- `Karl-Dashboard/docs/hyperdensity/HYPERDENSITY_CONTRACTKIT_CLAIMPOLICY_M18.md`
