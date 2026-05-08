# Karl-OS-ISO Defer Plan

## Decision
No `Karl-OS-ISO` porting in this milestone.

## Why Defer
- packaging Node Fluid Actuator before backend/controller contract stabilization increases rollback and safety risk
- runtime packaging should follow, not lead, contract/governance merge
- current audit objective is integration readiness, not substrate rollout

## Future Prerequisites
Before any OS/ISO track starts, require:
1. stable Hyperdensity contracts and controlled apply gates merged
2. explicit node-local permissions model and allowlist contract
3. kill switch behavior documented and tested end-to-end
4. audit trail and replay integrity checks accepted
5. systemd/DaemonSet packaging strategy approved with blast-radius constraints

## Future Implementation Topics (Deferred)
- systemd unit/agent lifecycle hardening
- DaemonSet wiring and permissions
- filesystem/path allowlist policies
- kill switch propagation
- signed delivery and auditable release channel
