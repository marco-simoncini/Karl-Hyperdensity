# Hyperdensity Ready Compliance Engine

Implemented `EvaluateWindowsHyperdensityReadyCompliance`:

- Input coverage: identity, QMP, guest fluidShell, RAM balloon, CPU actuator capability, policy annotations, pool context, remediation options.
- Output coverage: `compliancePhase`, `blockers`, `remediationActions`, `automatableActions`, `manualActions`, `risk`, `evidenceSummary`.

Phase model:

- `DISCOVERED_WINDOWS_VM`
- `ASSESSED_WINDOWS_VM`
- `BLOCKED_WITH_REMEDIATION`
- `FLUID_ENVELOPE_CANDIDATE`
- `HYPERDENSITY_READY_WINDOWS_SHELL`

Behavior:

- Standalone and pool-child VMs can become ready.
- Pool as scaling mechanism is blocked.
- Missing QMP/fluidShell/RAM balloon/CPU actuator produce remediation-required blocked state.
