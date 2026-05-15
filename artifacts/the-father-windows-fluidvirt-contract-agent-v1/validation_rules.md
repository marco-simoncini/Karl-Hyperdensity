# Validation Rules

Implemented helpers in `pkg/windowsfluidvirt`:

- `ValidateWindowsFluidShell`
- `ValidateFluidResourceLease`
- `ValidateWindowsFluidEvidence`
- `EvaluateWindowsFluidReadiness`
- `EvaluateLeaseCanBecomeActive`
- `EvaluateContinuityProofs`
- `EvaluateReturnToFloorReadiness`

Key enforced rules:

- shell runtime mode must be `in-place-qmp`
- no migration/reboot/recreate allowed in shell contract
- floor <= runtimeTarget <= envelope
- runtimeActual must stay within envelope
- `agentModule=fluidShell` and guest ACK required
- lease mode must be `in-place` with positive TTL
- all safety guarantees must remain true
- ACTIVE lease requires QMP+guest ack + continuity + rollback + return-to-floor readiness
- evidence continuity proofs require same pid/node/pod/boot/machine
- incomplete evidence cannot produce READY/ACTIVE without blockers
