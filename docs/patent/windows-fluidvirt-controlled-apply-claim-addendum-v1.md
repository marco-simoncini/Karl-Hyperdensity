# Windows FluidVirt Controlled Apply Claim Addendum v1

Draft tecnico per revisione brevettuale.

## Focus

- Hyperdensity Ready Windows target qualification
- feature-gated controlled apply planning
- CPU entitlement lease via node actuator (`cpu.max`)
- RAM entitlement lease via QMP balloon
- guest workload verification as mandatory witness layer
- rollback and return-to-floor as mandatory safeguards
- deterministic audit-chain append requirements
- autonomous apply disabled until explicit policy unlock

## Claim-Oriented Technical Framing

1. A control system receives a Windows Hyperdensity target and a fluid resource lease request.
2. A policy gate snapshot is evaluated before any apply readiness assertion.
3. Dry-run and manual approval are required by policy, with autonomous apply disabled by default.
4. Apply readiness is rejected when forbidden mechanisms are requested (vCPU hotplug, logical CPU scaling, VM spec patch, pool scaling mechanism).
5. Verification, rollback, return-to-floor, and audit bundle append plans are mandatory readiness dimensions.
6. The system emits deterministic plan/evidence output without executing runtime mutation in the planning milestone.

## Guarded Unlock Principle

Controlled apply can be made eligible only by explicit policy unlock and approval proof; no implicit autonomous transition is permitted.
