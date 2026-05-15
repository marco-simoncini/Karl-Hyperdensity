# Actuator Safety Model

Implemented `NodeFluidActuatorSafetyModel` + evaluator with controls:

- node-local blast radius
- target allowlist
- cgroup path validation
- PID/start-time validation
- pod UID / VM UID validation
- no symlink traversal
- no parent cgroup writes unless allowed
- no arbitrary file writes
- cpu.max bounds
- TTL enforcement
- stale request rejection
- replay protection (future-signable ready)
- kill switch and panic-safe rollback requirements
- audit log requirement
- dry-run mode support
- automatic return-to-floor requirement
