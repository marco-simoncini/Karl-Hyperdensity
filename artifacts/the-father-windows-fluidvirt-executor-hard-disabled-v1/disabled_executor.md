# Disabled Executor

`DisabledWindowsFluidApplyExecutor` always returns a non-executing result:

- `executionPhase` in `{EXECUTION_HARD_DISABLED, EXECUTION_BLOCKED, EXECUTION_QUARANTINED, EXECUTION_DENIED}`
- `applyAttempted=false`
- `mutationPerformed=false`
- `qmpCommandSent=false`
- `clusterMutationSent=false`
- blocker includes `future_apply_executor_disabled`
