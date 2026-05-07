# Future Apply Executor Interface

`WindowsFluidFutureApplyExecutor` defines only conceptual hard-disabled methods:

- `EvaluatePreApplyGuard`
- `BuildCommandEnvelope`
- `DenyExecution`
- `EvaluateKillSwitch`
- `EmitExecutionDeniedEvidence`

No CPU/RAM apply methods and no mutating QMP methods are exposed.
