# admission_gate

Main helper: `EvaluateWindowsFluidAdmission`.

Input:

- runtime evidence bundle
- optional policy pack (default conservative policy is applied otherwise)
- requested action
- optional fixed evaluation time for deterministic replay

Output:

- admission decision model
- linked non-mutating dry-run action slate
- evidence score details
- policy violations
- blockers
- next safe step

Behavior highlights:

- pool replica and generic Windows VM are denied by default policy;
- P0 blockers force `QUARANTINED`;
- P1 blockers force `BLOCKED`;
- P2 blockers lead to `BLOCKED` or `NEEDS_MORE_EVIDENCE`;
- CPU intent can become `ADMITTED_FOR_FUTURE_APPLY` only with complete evidence and threshold score;
- memory intent remains blocked when memory safety is not proven.
