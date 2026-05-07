# Gate Verification Model

Model: `WindowsFluidUnlockGateVerification`

- Gate IDs: Gate 0, Gate 1, Gate 2
- Statuses: PASSED, FAILED, BLOCKED, QUARANTINED, NOT_APPLICABLE
- Invariants always enforced: `executorMustRemainDisabled=true`, `mutationAllowed=false`, `applyAllowed=false`
- `PASSED` means gate condition verified, not unlock granted
