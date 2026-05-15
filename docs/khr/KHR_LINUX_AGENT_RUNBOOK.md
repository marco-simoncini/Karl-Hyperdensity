# KHR Linux Agent — Runbook (Sprint 5–9)

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

### `read-telemetry`

Read-only cgroup v2 metrics under `-cgroup-path` (typically from `discover-cgroups`). Optional `-allow-path-prefix`, `-cell-input`, `-telemetry-output`.

```bash
go run ./cmd/khr-linux-agent -mode read-telemetry \
  -config examples/khr/khr-linux-agent-config.yaml \
  -cgroup-path /sys/fs/cgroup/karl.slice/karl-shell-example.scope \
  -allow-path-prefix /sys/fs/cgroup/karl.slice
```

See `docs/khr/KHR_LINUX_READONLY_TELEMETRY.md` and `examples/khr/telemetry/`.

### `collect-evidence`

Single-shot **local evidence bundle** (validate config → discover → telemetry when a path is selected → optional ResourceLease dry-run). Requires `-config` and `-cell-input`. Optional `-cgroup-root`, `-allow-path-prefix`, `-lease-input`, `-resource-port-input`, `-cell-context`, `-evidence-output` (same lease/port/cell-context semantics as `dry-run` for the optional slice).

```bash
go run ./cmd/khr-linux-agent -mode collect-evidence \
  -config examples/khr/khr-linux-agent-config.yaml \
  -cell-input examples/khr/evidence/collect-evidence-input-cell.json \
  -cgroup-root /sys/fs/cgroup
```

See `docs/khr/KHR_LOCAL_EVIDENCE_BUNDLE.md`, `docs/khr/KHR_GRANDE_PADRE_EVIDENCE_HANDOFF.md`, `docs/khr/KHR_EVIDENCE_INTEGRITY_MODEL.md`, and `examples/khr/evidence/`.

#### Integrity sidecars (Sprint 10)

Optional `-evidence-manifest-output` and `-evidence-digest-output` write a **manifest** and **SHA-256** line computed from **canonical** bundle JSON (not necessarily identical to pretty-printed stdout). `-signing-mode` defaults to `none`; `local-dev` uses `-signing-key-file` (Ed25519 PKCS#8 PEM) and **requires** `-evidence-manifest-output`. This is developer workflow only, not production PKI.

```bash
go run ./cmd/khr-linux-agent -mode collect-evidence \
  -config examples/khr/khr-linux-agent-config.yaml \
  -cell-input examples/khr/evidence/collect-evidence-input-cell.json \
  -cgroup-root /tmp/khr-evidence-example-root \
  -evidence-manifest-output /tmp/evidence-manifest.json \
  -evidence-digest-output /tmp/evidence-bundle.sha256
```

Reproducing `examples/khr/evidence-integrity/*` fixtures uses fixed test clocks (`KHR_TEST_COLLECTED_AT`, `KHR_TEST_TELEMETRY_NOW`, `KHR_TEST_INTEGRITY_NOW`) and `KHR_TEST_INTEGRITY_CHAIN_STUB=1`; the sample tree under `/tmp/khr-evidence-example-root` must exist with the same cgroup metric files as in the Sprint 9 telemetry goldens.

## Non-goals

- No systemd unit changes.
- No package install / postinst hooks.
- No privilege escalation beyond what the operator already has when (future) writes are enabled.
