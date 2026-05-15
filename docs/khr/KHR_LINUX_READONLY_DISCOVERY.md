# KHR Linux — read-only cgroup discovery (Sprint 7)

## Purpose

`discover-cgroups` is a **read-only** host inspection mode. It helps operators and automation (including future **Grande Padre** evidence bundles) answer:

- Which cgroup hierarchy version is visible?
- Where might a **Cell**’s cgroup slice live under a scan root?
- Which candidate paths are **policy-allowed** under an optional prefix gate?

This mode **never** writes cgroup knobs, never mutates systemd, and never talks to the Kubernetes API.

## Non-goals (explicit)

- **Not apply** — discovery does not change `cpu.max`, `memory.max`, or any controller file.
- **Not live transfer** — discovery does **not** enable Hyperdensity live resource transfer, lease execution, or reconciliation side effects.
- **Not authority** — output is **evidence / readiness** input; admission and blast-radius policy remain upstream.

## CLI

```bash
go run ./cmd/khr-linux-agent -mode discover-cgroups \
  -config examples/khr/khr-linux-agent-config.yaml \
  -cgroup-root /sys/fs/cgroup \
  -cell-input examples/khr/cell-linux-envelope-full.json \
  -allow-path-prefix /sys/fs/cgroup/karl.slice
```

When `-cgroup-root` is omitted, the agent defaults to `/sys/fs/cgroup` (Linux unified mount).

### Output fields (JSON)

| Field | Meaning |
|-------|---------|
| `discoveryMode` | Always `read-only`. |
| `scannedRoot` | Absolute directory used as the discovery anchor. |
| `allowedPathPrefix` | Optional operator gate; empty disables extra prefix checks. |
| `candidatePaths` | Ordered probe list (explicit handle paths first, then heuristics). |
| `selectedPath` | First candidate that passes existence + symlink + prefix checks (may be empty). |
| `blockedReasons` | Non-fatal reasons candidates were rejected or nothing matched. |
| `warnings` | Best-effort diagnostics (missing dirs, symlink resolution noise). |
| `mutationsForbidden` | Always `true` in Sprint 7 builds. |

## Cell mapping heuristics

If `spec.providerHandle.cgroupPath` is present, it is **validated** (existence, directory, policy) but never modified.  
If the handle path targets the canonical host mount (`/sys/fs/cgroup/...`) while you scan a fixture root in tests, the agent **rebases** the suffix onto `-cgroup-root` for matching.

Otherwise the agent probes (in order):

1. `karl.slice/karl-shell-<shellRef.name>.scope`
2. `karl.slice/karl-shell-<shellRef.name>` (without `.scope`)
3. `karl.slice/<metadata.name>`
4. `karl.slice` (broad; last resort)

## Tests and fixtures

Golden stdout fixtures live under `examples/khr/discovery/` and use the placeholder `__CGROUP_ROOT__` so CI can substitute temporary directories.

See also: `docs/khr/KHR_CGROUP_PATH_POLICY.md`.
