# KHR local evidence bundle (`collect-evidence`, Sprint 9)

## Purpose

`khr-linux-agent -mode collect-evidence` produces **one JSON document** that chains, in order:

1. **Config validation** (same rules as `-mode validate-config`; fatal errors abort with non-zero exit and no bundle).
2. **Cgroup discovery** (`discover-cgroups` logic): read-only scan under `-cgroup-root` (default host root when empty inside the library) with optional `-allow-path-prefix`.
3. **Telemetry** (`read-telemetry` logic): if `discovery.selectedPath` is non-empty, metrics are sampled under that directory; otherwise telemetry is **skipped** with `telemetry.skipped: true` and a stable `skipReason` (bundle remains JSON-valid).
4. **ResourceLease dry-run** (optional): runs only when **both** `-lease-input` and `-resource-port-input` are provided and readable. If exactly one is provided, dry-run is skipped with a **warning** in `evidenceSummary.warnings` and matching `dryRun.skipReason`. When neither is provided, dry-run is skipped with an informational `skipReason` (not a lease/port contract warning).

There are **no writes**: no cgroup mutation, no systemd, no cluster APIs. Top-level `mutationsForbidden` is always `true`.

## CLI flags

| Flag | Required | Notes |
|------|----------|--------|
| `-config` | yes | Agent config (YAML/JSON). |
| `-cell-input` | yes | Cell JSON for discovery correlation and optional `cellRef`. |
| `-cgroup-root` | no | Discovery scan root. |
| `-allow-path-prefix` | no | Path prefix policy for discovery + telemetry. |
| `-lease-input` | no | Path to ResourceLease JSON (dry-run slice). |
| `-resource-port-input` | no | Path to ResourcePort JSON. |
| `-cell-context` | no | Optional CellContext JSON (same as `dry-run` mode); used only when both lease and port inputs are present. |
| `-cpu-delta` / `-memory-delta` | no | Passed through to envelope plan text (simulation only). |
| `-allow-unsafe-apply` | no | Non-operational: audit only, same as other modes. |
| `-evidence-output` | no | If set, writes the same JSON as stdout to a file (mode `0o600`). |
| `-evidence-manifest-output` | no | Artifact manifest JSON (Sprint 10); see `docs/khr/KHR_EVIDENCE_INTEGRITY_MODEL.md`. |
| `-evidence-digest-output` | no | Single-line SHA-256 hex of **canonical** bundle JSON. |
| `-signing-mode` | no | `none` (default) or `local-dev` (requires manifest + key file). |
| `-signing-key-file` | no | Ed25519 PEM for `local-dev` only. |
| `-artifact-id` | no | Optional `artifactId` in manifest. |

## Top-level JSON fields

| Field | Description |
|-------|-------------|
| `tool` | `khr-linux-agent` |
| `version` | Agent build string (e.g. `0.0.1-sprint13`). |
| `mode` | `collect-evidence` |
| `agentId` | From config `spec.agentId`. |
| `collectedAt` | RFC3339 UTC (`KHR_TEST_COLLECTED_AT` overrides in tests). |
| `cellRef` | Optional Cell identity when a Cell was parsed. |
| `discovery` | Snapshot without nested `tool`/`version`/`mode` from standalone discovery. |
| `telemetry` | Telemetry-shaped snapshot; may include `skipped` / `skipReason`. |
| `dryRun` | Optional dry-run payload; may be `skipped: true` with `skipReason`. |
| `evidenceSummary` | Aggregated `confidence`, `readyForGrandePadre`, `blockedReasons`, `warnings`, `recommendedNextAction`. |
| `mutationsForbidden` | Always `true`. |

Implementation lives in `pkg/khr/evidence/` (`CollectEvidenceBundle` and summarizer). Examples: `examples/khr/evidence/`. Integrity sidecars: `pkg/khr/evidence/integrity/`, `examples/khr/evidence-integrity/`.

## Deterministic tests

Golden tests set:

- `KHR_TEST_CGROUP_VERSION`
- `KHR_TEST_COLLECTED_AT`
- `KHR_TEST_TELEMETRY_NOW`

Placeholders in fixtures use `__CGROUP_ROOT__` for the temp scan root, replaced at test time (same pattern as discovery goldens).
