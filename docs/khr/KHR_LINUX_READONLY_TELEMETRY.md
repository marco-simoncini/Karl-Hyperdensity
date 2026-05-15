# KHR Linux — read-only cgroup telemetry (Sprint 8)

## Purpose

`read-telemetry` samples **cgroup v2** controller files from a **directory path already chosen** by the operator or by `discover-cgroups` (Sprint 7). It emits structured JSON suitable as **evidence** for Hyperdensity / **Grande Padre** readiness pipelines.

## Non-goals (explicit)

- **Not apply** — no writes to `cpu.max`, `memory.max`, `memory.high`, or any cgroup knob.
- **Not a decision engine** — telemetry does not approve or deny `ResourceLease` objects and does not trigger apply.
- **Not live transfer** — this mode does not move resources between Cells or reconcile controllers.

## CLI

```bash
go run ./cmd/khr-linux-agent -mode read-telemetry \
  -config examples/khr/khr-linux-agent-config.yaml \
  -cgroup-path /sys/fs/cgroup/karl.slice/karl-shell-example.scope \
  -allow-path-prefix /sys/fs/cgroup/karl.slice \
  -cell-input examples/khr/telemetry/telemetry-input-cell.json \
  -telemetry-output /tmp/khr-telemetry.json
```

`-telemetry-output` is optional; when set, the same JSON written to stdout is also written to the given path (handy for evidence bundles).

## Safety behaviour

- The cgroup path must resolve to a **directory**; symlink chains that escape the resolved directory or optional `-allow-path-prefix` are **blocked** (see `pkg/khr/cgroup/file_read_policy.go`).
- Missing individual metric files produce **warnings** and partial metrics.
- If **both** `cpu.stat` and `memory.current` are missing/unreadable, the reader adds a **blockedReason** explaining that the sample is unusable.

## Fixture placeholders

Example JSON under `examples/khr/telemetry/` uses `__ROOT__` as a stand-in for a temporary directory root in CI (tests substitute it with `t.TempDir()`), avoiding accidental doubling when a cgroup path already includes subdirectories.

## Relationship to discovery

Run `discover-cgroups` first when you do not know the slice path, then pass `selectedPath` into `-cgroup-path` for telemetry.

See also: `docs/khr/KHR_TELEMETRY_EVIDENCE_MODEL.md`.
