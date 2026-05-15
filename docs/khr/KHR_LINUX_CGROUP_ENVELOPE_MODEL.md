# KHR Linux — cgroup Envelope Model (Sprint 5)

## cgroup v2

The agent uses best-effort detection via `cgroup.controllers` at `/sys/fs/cgroup`. Unknown environments return `unknown` version and **must not** assume v2 semantics.

## Envelope semantics (planning only)

**Envelope** (Sprint 5 contract) means bounded CPU and memory pressure shaping using cgroup `cpu.max` / `memory.max` style controls **once writes are implemented and gated**. This sprint only returns:

- `cpuMaxDelta` / `memoryMaxDelta` echo strings for traceability
- `wouldWrite: false` unless `--allow-unsafe-apply` is passed (still discouraged)

## No-write default

`PlanEnvelope` never touches the filesystem unless `allowWrite` is true **and** the caller has not been overridden by the safety layer. The CLI clears `writePaths` whenever `--allow-unsafe-apply` is absent.

## Future work

- Map **Cell** identity to cgroup slice paths with KARL naming.
- Integrate observed usage from `pkg/khr/telemetry` (currently stub).
- Coordinate with Hyperdensity `ResourceLease` apply path (post–Sprint 5).
