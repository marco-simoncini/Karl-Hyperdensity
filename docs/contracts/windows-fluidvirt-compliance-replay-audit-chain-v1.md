# Windows FluidVirt Compliance Replay Audit Chain v1

This milestone introduces deterministic compliance replay and audit hash chain modeling for Windows FluidVirt.

## Scope

- deterministic compliance replay only
- audit hash chain modeling and verification
- model/replay-only functions in package code
- no operational CLI porting in this milestone

## Explicit Non-Scope

- no runtime actuator enablement
- no real cgroup write
- no QMP execution
- no QGA execution
- no controlled apply enablement
- no executor enablement
- no production apply
- no autonomous apply
- no Windows GA claim
- no Windows production-ready claim

## Semantics

- `complianceReplayExecuted=true` is replay/model verification only
- it does not imply runtime readiness or controlled apply readiness
- replay and audit chain remain deterministic and non-mutative

## Safety

- runtime mutation flags remain disabled
- raw runtime controls remain disabled
- replay events and input explicitly mark no runtime touch and no secret material
