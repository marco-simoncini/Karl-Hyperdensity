# Windows FluidVirt Guarded Executor Boundary v1

This milestone defines the Windows FluidVirt executor contract boundary in disabled mode.

## Boundary State

- boundary-only (`guarded_executor_boundary_only`)
- executor boundary defined
- executor disabled
- executor runtime unavailable
- executor not executed

## Runtime Safety

- controlled apply disabled
- runtime actuator disabled
- runtime mutation disabled
- no real cgroup write
- no QMP/QGA execution
- no raw runtime controls

## Claim Boundary

- no Windows GA
- no Windows production-ready
- no Windows execution-ready by default
- no autonomous apply
- no production apply

## Required Preconditions Before Any Future Executor Enablement

- controlled apply plan ready
- manual approval + operator identity audit
- active lease + TTL
- guest witness integration (fluidShell/QGA, same-boot, same-QEMU)
- rollback baseline/plan and return-to-floor plan
- compliance replay + audit hash chain verification
- kill switch verification
- node allowlist + cgroup path validation strategy

Future guarded fake-runtime executor MVP remains a separate milestone.
