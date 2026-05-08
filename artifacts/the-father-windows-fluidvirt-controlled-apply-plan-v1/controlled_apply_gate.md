# Controlled Apply Gate

`WindowsFluidControlledApplyGate` defaults to safe deny for apply:

- manual approval required
- dry-run required
- autonomous apply disabled
- kill switch required
- rollback/return/audit/workload verification required

Apply can be enabled only by explicit gate fixture + approval + allowlists.
