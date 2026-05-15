# KHR telemetry evidence model (Sprint 8)

## Evidence envelope

Each `read-telemetry` response includes an `evidence` object:

| Field | Role |
|-------|------|
| `observedAt` | UTC timestamp (RFC3339) when the read finished. |
| `source` | Always `cgroup-v2` for this sprint. |
| `confidence` | `high` when both CPU (`cpu.stat` keys) and memory (`memory.current`) are present; `medium` when only one axis is present; `low` otherwise. |
| `warnings` | Non-fatal issues (missing optional files, parse skips, symlink resolution noise). |
| `blockedReasons` | Policy failures or unusable samples (for example prefix violations, or missing both `cpu.stat` and `memory.current`). |

## Consumer guidance (Grande Padre / Hyperdensity)

1. Treat telemetry JSON as **read-only evidence**, not as an admission decision.
2. Combine with `ResourcePort`, `ResourceLease` simulation output, and operator context before any future apply.
3. Down-rank or reject automation actions when `confidence` is `low` or `blockedReasons` is non-empty.

## Mutations

`mutationsForbidden` is always `true` in Sprint 8 builds: cgroup writes remain out of scope for this binary.
