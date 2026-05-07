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
