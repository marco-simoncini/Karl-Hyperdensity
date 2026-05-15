# KHR Linux Agent — Runbook (Sprint 5–7)

**Audience:** platform engineers evaluating the skeleton locally.

## Build

```bash
go build -o bin/khr-linux-agent ./cmd/khr-linux-agent
```

## Modes

### `validate-config`

Validates `examples/khr/khr-linux-agent-config.yaml` (or your file):

```bash
go run ./cmd/khr-linux-agent -mode=validate-config -config=examples/khr/khr-linux-agent-config.yaml
```

Exit `0` when valid, non-zero when validation errors are present.

### `print-capabilities`

Emits JSON describing cgroup version, stub runtime providers, `mutationsForbidden` (always true in Sprint 6), optional `audit` when `--allow-unsafe-apply` is passed, and `futureApplyGateRequired` when that flag is present.

```bash
go run ./cmd/khr-linux-agent -mode=print-capabilities -config=examples/khr/khr-linux-agent-config.yaml
```

### `dry-run`

Requires `-lease-input` and `-resource-port-input`. Optional `-cell-context` (defaults to Linux/Linux when omitted inside the evaluator). Optional `-cpu-delta` / `-memory-delta` annotate the cgroup envelope plan text only (still no writes). Golden stdout fixtures live under `examples/khr/golden/` (tests set `KHR_TEST_CGROUP_VERSION` for deterministic cgroup version in JSON).

```bash
go run ./cmd/khr-linux-agent -mode=dry-run \
  -config=examples/khr/khr-linux-agent-config.yaml \
  -lease-input=examples/khr/resourcelease-linux-envelope-dry-run.json \
  -resource-port-input=examples/khr/resourceport-linux-envelope-for-dryrun.json \
  -cell-context=examples/khr/linux-cell-dry-run-input.json \
  -cpu-delta="-100000" \
  -memory-delta="-256Mi"
```

`--allow-unsafe-apply` is **non-operational** in Sprint 6: it never enables cgroup writes but emits an audit warning and sets `futureApplyGateRequired` in JSON. Do not treat the flag as a safety bypass; see `docs/khr/KHR_AUDIT_AND_APPLY_GATES.md`.

### `discover-cgroups`

Read-only cgroup discovery for mapping a `Cell` (optional) to candidate cgroup paths. Never writes.

```bash
go run ./cmd/khr-linux-agent -mode discover-cgroups \
  -config examples/khr/khr-linux-agent-config.yaml \
  -cgroup-root /sys/fs/cgroup \
  -cell-input examples/khr/cell-linux-envelope-full.json \
  -allow-path-prefix /sys/fs/cgroup/karl.slice
```

See `docs/khr/KHR_LINUX_READONLY_DISCOVERY.md` and `examples/khr/discovery/` for fixtures.

## Non-goals

- No systemd unit changes.
- No package install / postinst hooks.
- No privilege escalation beyond what the operator already has when (future) writes are enabled.
