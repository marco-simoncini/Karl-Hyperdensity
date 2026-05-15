# KHR Linux MVP — Design (Sprint 5)

**Repository:** `marco-simoncini/Karl-Hyperdensity`  
**Branch:** `KHR`  
**Status:** Skeleton only — **no production daemon**, **no ISO/installer changes**, **no automatic cgroup writes**.

## Scope

- **Linux only** — cgroup envelope planning and `ResourceLease` dry-run validation.
- **No VM / Windows / KubeVirt mutation** — KubeVirt remains untouched; this agent does not speak to the KubeVirt API.
- **Dry-run / simulation-first** — default CLI modes validate config and print plans as JSON.

## Components

| Path | Role |
|------|------|
| `cmd/khr-linux-agent` | CLI entry (`validate-config`, `dry-run`, `print-capabilities`) |
| `pkg/khr/agent` | Config load (YAML/JSON) and orchestration helpers |
| `pkg/khr/cgroup` | cgroup v2 detection (best-effort) and envelope **plan** (no write unless `--allow-unsafe-apply`) |
| `pkg/khr/resourcelease` | `ResourceLease`-shaped dry-run evaluation |
| `pkg/khr/runtimeprovider` | Linux provider **stubs** (systemd, cgroup envelope) |
| `pkg/khr/cellwatch` / `pkg/khr/telemetry` | No-op placeholders |
| `pkg/khr/safety` | Hard rule: mutations forbidden without explicit unsafe flag |

## CLI contract

See `docs/khr/KHR_LINUX_AGENT_RUNBOOK.md` and `examples/khr/`.

## References

- `docs/khr/KHR_LINUX_CGROUP_ENVELOPE_MODEL.md`
- `docs/khr/KHR_SAFETY_AND_DRY_RUN_MODEL.md`
- `marco-simoncini/Karl-Dashboard` Shell read-only contracts (Sprint 4)
