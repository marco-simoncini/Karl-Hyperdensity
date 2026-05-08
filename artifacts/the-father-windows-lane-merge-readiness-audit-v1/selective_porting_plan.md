# Selective Porting Plan

## Integration Principle
Backend contracts first, runtime apply last, no dashboard import from stale branch, no direct merge.

## PR Sequence (Small, Ordered)

### PR 1 - Windows FluidVirt product contracts/models
Scope:
- `pkg/windowsfluidvirt/product_model.go`
- `pkg/windowsfluidvirt/action_slate.go`
- `pkg/windowsfluidvirt/blockers.go`
- `docs/contracts/windows-fluidvirt-product-model-v1.md`
- minimal fixtures from `examples/windows-fluid-product-fixtures/*`

Acceptance gates:
- contract schema stability
- explicit rejection paths for vCPU hotplug/logical CPU scaling/pool scaling
- no runtime mutation execution path

### PR 2 - Node Fluid Actuator MVP
Scope:
- `pkg/windowsfluidvirt/node_actuator_mvp.go`
- `pkg/windowsfluidvirt/node_actuator_compliance.go`
- `cmd/karl-node-fluid-actuator/main.go`
- `docs/contracts/windows-fluid-node-actuator-v1.md`
- minimal allowlist/request fixtures

Acceptance gates:
- cgroup path exact match and no path escape
- no parent/arbitrary writes
- TTL + replay + kill switch checks

### PR 3 - Compliance replay + audit hash chain
Scope:
- `pkg/windowsfluidvirt/compliance_replay_cli.go`
- `cmd/karl-fluid-compliance-replay/main.go`
- attestation/bundle index docs
- minimal replay fixtures

Acceptance gates:
- deterministic hash chain validation
- read-only replay behavior
- no real-signature production claim

### PR 4 - Controlled apply planning
Scope:
- `pkg/windowsfluidvirt/controlled_apply_plan.go`
- `docs/contracts/windows-fluidvirt-controlled-apply-plan-v1.md`
- controlled-apply planning fixtures

Acceptance gates:
- `autonomousApplyEnabled=false` default
- manual approval gate enforced
- dry-run + rollback + return-to-floor + audit blockers enforced

### PR 5 - Runtime executor fake-runtime / later
Scope:
- `cmd/karl-fluid-windows-executor/main.go`
- any `executor_disabled` planning-only evaluator pieces

Acceptance gates:
- planning-only mode confirmed
- runtime mutation remains disabled
- no cluster/QMP side effects

### PR 6 - Inventory fluidShell guest witness (separate repo)
Scope in `Karl-Inventory`:
- fluidShell module/contracts, worker integration, tests

Acceptance gates:
- witness/evidence-only role
- no actuator/scaler claim
- pending reboot/critical events/same-boot signals validated

### PR 7 - Dashboard contract visibility (later)
Scope:
- start from `The-Father` only (never from `The-Father-Windows`)
- read-only contract visibility widgets only

Acceptance gates:
- no raw runtime controls
- no 443/8888 production touch in this track
- backend contracts stable first
