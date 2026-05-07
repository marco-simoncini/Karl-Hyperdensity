# Windows Hyperdensity Ready Compliance v1

`EvaluateWindowsHyperdensityReadyCompliance` classifies Windows VMs into readiness phases and generates remediation.

## Readiness phases

- `DISCOVERED_WINDOWS_VM`
- `ASSESSED_WINDOWS_VM`
- `BLOCKED_WITH_REMEDIATION`
- `FLUID_ENVELOPE_CANDIDATE`
- `HYPERDENSITY_READY_WINDOWS_SHELL`

The engine reports blocked/not-ready-yet/ready states and never marks targets as "unsupported forever".

## Supported object origins

- Standalone VM.
- Pool-child VM.
- VM from golden image lineage.

Pool remains provisioning context only; it is not accepted as runtime scaling mechanism.

## Compliance inputs

- VM/VMI and runtime identity evidence.
- QMP availability evidence.
- Guest `fluidShell` and ACK evidence.
- RAM balloon capability evidence.
- CPU actuator capability evidence.
- Policy annotations and pool context.

## Remediation model

The engine binds blockers to remediation taxonomy IDs, split into:

- `automatableActions`
- `manualActions`

Each action declares required evidence, rollback direction, risk, and resolved blocker.

## Runtime boundaries

- No vCPU hotplug/unplug path.
- No VM spec patch mutation path.
- No LiveMigration/VMIM mechanism.
- No replica/pool scaling as runtime mechanism.

## Real Master-Win11 Replay Classification

Using evidence from `WINDOWS_PREARMED_FLUID_ENVELOPE_CONFIRMED`, `master-win11` is expected to classify as `HYPERDENSITY_READY_WINDOWS_SHELL` when:

- QMP evidence is present.
- `fluidShell` guest evidence is present with guest ACK.
- RAM balloon capability is present and return-to-floor is verified.
- Node-local CPU actuator capability is present and allowlisted.
- continuity invariants remain true (same QEMU, same boot, same node/pod, no migration/recreate/rollout).

## Prearmed Fluid Envelope v2 Target

Compliance target for Windows readiness is the Prearmed Fluid Envelope v2 model:

- CPU liquidity by entitlement lease (`cgroup v2 cpu.max` via authorized node-local actuator).
- RAM liquidity by QMP balloon lease.
- rollback and return-to-floor evidence mandatory.

This target does not claim logical CPU count scaling and does not claim vCPU hotplug success.

## Standalone and Pool-Child Semantics

- Standalone Windows VM: eligible for direct readiness classification.
- Pool-child Windows VM: eligible when treated as individual shell identity.
- Pool context is accepted only for provisioning lineage and blocked if requested as runtime scaling mechanism.
