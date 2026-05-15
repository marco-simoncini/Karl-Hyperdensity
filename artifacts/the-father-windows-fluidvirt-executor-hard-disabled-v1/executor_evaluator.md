# Executor Evaluator

`EvaluateWindowsFluidFutureApplyExecutor` takes governance contract, revalidation, attestation, optional kill switch, and optional evaluation time.

It returns:

- disabled execution result
- pre-apply guard
- kill switch snapshot
- preview-only command envelope
- blocker set and next safe step

Execution is always denied and non-mutating.
