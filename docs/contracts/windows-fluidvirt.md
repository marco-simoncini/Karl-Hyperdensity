# windows-fluidvirt-v1

Evidence-only contract pack for certified Windows Fluid Shells in KARL Hyperdensity.

## Scope

- No CPU/RAM mutation in this phase.
- No QMP mutating commands in this phase.
- No live migration/recreate/reboot flow in this phase.
- Readiness resolves only from continuity evidence + blockers.

## Contracts

- `WindowsFluidShell` (`hyperdensity_windows_fluid_shell_v1`)
  - Defines certified shell identity and fluid envelope.
  - Enforces `runtimeMode=in-place-qmp`, no reboot/recreate/migration.
  - Tracks `runtimeTarget` separately from `runtimeActual`.
- `FluidResourceLease` (`hyperdensity_fluid_resource_lease_v1`)
  - Captures granted in-place lease and safety guarantees.
  - Prevents `ACTIVE` when ack/continuity/rollback/return-to-floor proofs are missing.
- `WindowsFluidEvidence` (`hyperdensity_windows_fluid_evidence_v1`)
  - Collects before/after continuity proofs for qemu pid, pod, node, last boot, machine GUID.
  - Carries guest evidence payload emitted by `modules.fluidShell`.
- `WindowsFluidBlocker` (`hyperdensity_windows_fluid_blocker_v1`)
  - Canonical blocker taxonomy with stable IDs.

## State Machine

- `EMPTY -> PREFLIGHT -> READY -> LEASE_PREPARED -> APPLYING -> VERIFYING -> ACTIVE`
- `ACTIVE -> RETURNING_TO_FLOOR -> EMPTY`
- Any continuity break (`qemu pid`, `last boot`, `machine GUID`, `node`, `virt-launcher pod`) moves to `QUARANTINED`.
- Missing mandatory ack/proofs moves to `BLOCKED`.

## Mandatory Gate Rules

- `APPLYING -> ACTIVE` requires:
  - QMP ack
  - guest ack
  - no reboot proof
  - no recreate proof
  - no migration proof
  - same qemu proof
- `RETURNING_TO_FLOOR -> EMPTY` requires:
  - resource return verified
  - guest confirmation
  - QMP confirmation
  - no reboot proof
  - same qemu proof
- Memory downscale gate:
  - blocked when `memory_return_not_safe`
  - blocked when `memory_driver_unverified`
  - blocked when `guest_memory_not_confirmed`

## Integration Source

Guest-side evidence source is `Karl-Inventory modules.fluidShell`.
Integration mapping is defined in `docs/contracts/windows-fluidvirt-integration-contract-v1.md`.

## Runtime Gate and QMP discovery

- Runtime gate contract: `docs/contracts/windows-fluid-runtime-gate-v1.md`
- KubeVirt identity model: `docs/contracts/windows-fluid-kubevirt-identity-v1.md`
- QMP evidence contract: `docs/contracts/windows-fluid-qmp-evidence-v1.md`
- QMP command policy: `docs/contracts/windows-fluid-qmp-command-policy-v1.md`
- Dry-run evaluation contract: `docs/contracts/windows-fluid-dryrun-evaluation-v1.md`
- Admission policy contract: `docs/contracts/windows-fluid-admission-policy-v1.md`
- Read-only cluster discovery runbook: `docs/runbooks/windows-fluidvirt-readonly-cluster-discovery-v1.md`
- Lab evidence replay runbook: `docs/runbooks/windows-fluidvirt-lab-evidence-replay-v1.md`
- Admission replay runbook: `docs/runbooks/windows-fluidvirt-admission-replay-v1.md`
