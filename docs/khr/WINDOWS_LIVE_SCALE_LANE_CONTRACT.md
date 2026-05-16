# Windows live scale lane contract (KHR-P)

KHR-P prepares the **Windows lane** for CPU/RAM live scale **without restart**. This sprint is **observation + contract + safety model only**. There is **no Windows runtime apply** in Hyperdensity, ISO, or host-runtime yet.

## Scope

| In scope | Out of scope |
|----------|----------------|
| ResourcePort capability contract | `resourcelease-guarded-apply` on Windows |
| Dry-run / observation JSON examples | MSI or Inventory agent changes |
| Blocked-state vocabulary | KubeVirt VM mutation |
| Provider binding semantics | Session drain enforcement |

## ResourcePort capabilities (`spec.windowsLiveScale`)

| Field | Type | Meaning |
|-------|------|---------|
| `cpuLiveScaleSupported` | bool | CPU hot-adjust without guest restart (target path) |
| `ramLiveScaleSupported` | bool | RAM hot-adjust without guest restart (target path) |
| `scaleUpSupported` | bool | Scale-up mode allowed on port |
| `scaleDownSupported` | bool | Scale-down mode allowed on port |
| `requiresRestart` | bool | **Target `false`** on `windows.host-runtime`; compatibility providers may report `true` |
| `providerBinding` | string | `windows.host-runtime` \| `kubevirt.compatibility` |

Schema: `docs/contracts/khr/windows-live-scale-resourceport.schema.json`

## Provider bindings

| Provider | Live scale posture | Notes |
|----------|-------------------|--------|
| `windows.host-runtime` | Target **live-in-place** | Native Windows session / host-runtime lane (future apply) |
| `kubevirt.compatibility` | **compatibility-fallback** | Observed-only; may imply VM-level restart for resource change |

## Blocked states

Dry-run and observation MUST surface blocked states instead of silent apply:

| State | When |
|-------|------|
| `requiresRestart` | Provider or guest requires process/VM restart for change |
| `requiresReboot` | Guest OS reboot required |
| `requiresSessionDrain` | Active RDP/session must drain before adjust |
| `providerUnsupported` | Provider binding cannot honor live-in-place scale |

Schema: `docs/contracts/khr/windows-live-scale-blocked-state.schema.json`

## ResourceLease observation

- **Observed** (`resourcelease-windows-cpu-dryrun-observed.json`): dry-run allowed, no mutation, documents expected observation writes.
- **Blocked** (`resourcelease-windows-memory-dryrun-blocked.json`): e.g. `requiresRestart` on compatibility path or over-limit delta.

## Safety model

1. **No restart target** — Hyperdensity live scale on Windows MUST NOT document restart/reboot as the default outcome.
2. **Sandbox / observation first** — contract examples are JSON fixtures only.
3. **Production block** — same production namespace blocklist as Linux KHR-O until Windows apply is explicitly enabled in a future sprint.

## Examples

| File | Role |
|------|------|
| `examples/khr/windows/resourceport-windows-session.json` | Session ResourcePort with `windows.host-runtime` capabilities |
| `examples/khr/windows/resourcelease-windows-cpu-dryrun-observed.json` | CPU dry-run observation (allowed) |
| `examples/khr/windows/resourcelease-windows-memory-dryrun-blocked.json` | Memory dry-run blocked (`requiresRestart`) |

## Validation

```bash
./scripts/validate_windows_live_scale_contract.sh
```
