# Real Evidence Bundle Summary

Bundle status: `blocked-preflight`

## Collected

- KubeVirt VM identity (`master-win11`, UID, runStrategy, status)
- Namespace-level VMIM check
- Namespace events snapshot
- Pool context-only snapshot (`win11-pool-*`)

## Missing for runnable bundle

- Live VMI identity and phase
- virt-launcher pod UID and restart counters
- QMP handshake and query evidence
- Guest fluidShell ACK telemetry
- Runtime CPU/RAM before/after pairs

## Governance posture

- If not certain, do not apply.
- This run did not attempt any runtime mutation.
