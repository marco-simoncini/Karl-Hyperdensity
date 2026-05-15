# State Machine

Source: `pkg/windowsfluidvirt/contracts.go`, `pkg/windowsfluidvirt/state_machine.go`

States implemented:

- EMPTY
- PREFLIGHT
- READY
- LEASE_PREPARED
- APPLYING
- VERIFYING
- ACTIVE
- RETURNING_TO_FLOOR
- ROLLED_BACK
- BLOCKED
- QUARANTINED

Gate rules implemented:

- `EvaluateLeaseCanBecomeActive`
  - requires QMP ACK + guest ACK
  - requires last boot unchanged + qemu pid unchanged
  - requires rollback ready + return-to-floor ready
  - continuity breaks can escalate to `QUARANTINED`
- `EvaluateReturnToFloorReadiness`
  - blocks when return-to-floor not ready
  - blocks when memory return is unsafe
- `EvaluateContinuityProofs`
  - qemu pid / last boot / machine guid / node / virt-launcher pod changes produce blockers
- missing ACKs resolve to `BLOCKED`
