# Actuator Safety Validator

Validator: `ValidateNodeFluidActuatorRequest`

Checks:

- request TTL and stale detection
- kill-switch enforcement
- node/namespace/vm/pod/qemu identity exact match
- qemu start-time exact match
- cgroup path exact allowlist match
- target file limited to `cpu.max`
- symlink/path traversal rejection
- parent cgroup write rejection
- arbitrary write rejection
- controller allowlist enforcement
- requested/rollback bounds inside request and allowlist bounds
- previous value readback match
- audit output path safety check
