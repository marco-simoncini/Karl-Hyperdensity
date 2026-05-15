# KHR Linux Agent — Runbook (Sprint 5)

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

Emits JSON describing cgroup version, stub runtime providers, and whether mutations are locked.

```bash
go run ./cmd/khr-linux-agent -mode=print-capabilities -config=examples/khr/khr-linux-agent-config.yaml
```

### `dry-run`

Requires `-lease-input` and `-resource-port-input` (Sprint 5 policy). Optional `-cell-context` (defaults to Linux/Linux when omitted inside evaluator). Optional `-cpu-delta` / `-memory-delta` for envelope plan text.

```bash
go run ./cmd/khr-linux-agent -mode=dry-run \
  -config=examples/khr/khr-linux-agent-config.yaml \
  -lease-input=examples/khr/resourcelease-linux-envelope-dry-run.json \
  -resource-port-input=examples/khr/resourceport-linux-envelope-for-dryrun.json \
  -cell-context=examples/khr/linux-cell-dry-run-input.json \
  -cpu-delta="-100000" \
  -memory-delta="-256Mi"
```

**Never** pass `--allow-unsafe-apply` outside an isolated lab; Sprint 5 tests do not enable it.

## Non-goals

- No systemd unit changes.
- No package install / postinst hooks.
- No privilege escalation beyond what the operator already has when (future) writes are enabled.
