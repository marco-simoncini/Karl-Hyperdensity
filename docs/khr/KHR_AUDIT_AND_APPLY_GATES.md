# KHR audit and apply gates (Sprint 6)

This document describes how the `khr-linux-agent` records **audit** signals and why **apply remains gated** even when operators pass `--allow-unsafe-apply`.

## Principles

1. **Dry-run and read-only by default** — the Linux MVP path must not mutate host cgroup files or systemd units.
2. **No silent escalation** — flags that sound dangerous must produce explicit **audit warnings** in structured JSON.
3. **Future apply gate** — real writes require a deliberate controller/operator workflow not implemented in Sprint 6.

## `--allow-unsafe-apply` (non-operational)

Sprint 6 keeps the flag **non-operational**:

- It does **not** enable cgroup writes, systemd mutations, or Kubernetes applies.
- When present, the agent emits an `audit` record with code `KHR_AUDIT_UNSAFE_APPLY_NON_OPERATIONAL` and sets `futureApplyGateRequired` to `true` in dry-run JSON output.
- `mutationsForbidden` remains `true` in JSON output regardless of the flag (writes are always disabled in this build).

This preserves CLI compatibility for future gates while avoiding a foot-gun where a flag name implies safety that does not yet exist.

## Planned apply gate (not implemented)

Before any real apply path ships, KHR expects at minimum:

- A versioned **exec contract** match between `RuntimeProvider.spec.execContractVersion` and the agent bundle.
- A **policy / admission** decision (simulation evidence, blast-radius, rollback hooks) recorded out-of-band.
- An explicit **two-person rule** or automation secret for production clusters (operator policy, not code in this repo).

## JSON fields

| Field | Meaning |
|-------|---------|
| `audit` | Ordered list of `audit.Record` entries (warnings, informational notices). |
| `mutationsForbidden` | Always `true` in Sprint 6 builds. |
| `unsafeApplyFlagPresent` | Mirrors CLI: operator passed `--allow-unsafe-apply`. |
| `futureApplyGateRequired` | `true` when the unsafe flag is present: reminds automation that a future gate must pass before writes. |

## Related documents

- `docs/khr/KHR_SAFETY_AND_DRY_RUN_MODEL.md`
- `docs/khr/KHR_LINUX_AGENT_RUNBOOK.md`
