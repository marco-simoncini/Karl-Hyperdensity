# Actuator Lifecycle

Result model: `KARLNodeFluidActuatorResult`

Lifecycle decisions:

- `accepted` (dry-run valid)
- `rejected` (policy/validation reject)
- `applied` (apply success)
- `rolled_back` (rollback success)
- `returned_to_floor` (return-to-floor success)
- `blocked` (runtime IO/readback failure)

Rules:

- dry-run never mutates (`mutationPerformed=false`)
- apply requires before/after readback
- rollback restores `rollbackCpuMax`
- return-to-floor restores floor target
- blocked/rejected must include blockers
- each result includes deterministic `auditHash`
