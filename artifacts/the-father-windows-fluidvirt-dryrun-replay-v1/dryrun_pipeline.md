# dryrun_pipeline

Pipeline function: `EvaluateWindowsFluidRuntimeDryRunWithOptions`.

Input:

- `WindowsFluidRuntimeEvidenceBundle`
- optional evaluation timestamp (`DryRunEvaluationOptions`)

Output:

- phase (`READY`, `BLOCKED`, `QUARANTINED`, `LEASE_PREPARED`)
- certification classification
- conditions/blockers/evidence summary
- non-mutating action slate
- recommended next safe step

Constraints enforced:

- no `ACTIVE` output path
- no `APPLYING` output path
- lease intent can only reach `LEASE_PREPARED`
- lease intent is blocked when rollback/return-to-floor is not proven
- pool replica model remains context-only and blocked
