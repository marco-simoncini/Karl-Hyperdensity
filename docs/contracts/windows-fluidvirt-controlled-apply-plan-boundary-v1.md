# Windows FluidVirt Controlled Apply Plan Boundary v1

This milestone defines the controlled apply plan in **boundary-only / plan-only** mode.

## Boundary State

- controlled apply plan is defined
- controlled apply is not enabled
- controlled apply is not executed
- controlled apply is not ready
- executor is not enabled
- runtime actuator is not enabled
- runtime mutation is disabled

## Runtime Safety

- no real cgroup write
- no QMP execution
- no QGA execution
- no raw runtime control exposure
- no autonomous apply
- no production apply

## Claim Boundary

- no Windows GA
- no Windows production-ready
- no Windows execution-ready by default
- no vCPU hotplug claim
- no logical CPU scaling claim
- no pool scaling claim

## Preconditions Before Any Future Controlled Apply

- compliance replay + audit hash chain verified
- guest witness integration
- manual approval flow and operator identity audit
- active lease with TTL and expiry return-to-floor
- rollback and return-to-floor plans proven
- kill switch verified
- node allowlist and cgroup path validation strategy

This contract does not introduce any executable apply path.
