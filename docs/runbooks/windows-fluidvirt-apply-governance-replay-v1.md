# windows-fluidvirt-apply-governance-replay-v1

Runbook for local replay of apply-governance contracts and formal transitions.

## Safety boundaries

- Local replay only.
- No cluster calls.
- No QMP calls.
- No runtime apply or mutation.

## Fixture replay

```bash
go run ./cmd/karl-fluid-governance \
  -fixture ./examples/windows-fluid-governance-fixtures/governance-master-win11-cpu.contract-prepared.json \
  -evaluation-time 2026-05-07T14:45:00Z
```

## Bundle + admission replay

```bash
go run ./cmd/karl-fluid-governance \
  -admission /path/to/admission-decision.json \
  -bundle /path/to/runtime-bundle.json \
  -policy /path/to/policy-pack.json \
  -requested-action future-cpu-apply
```

## Reading governance outputs

- `CONTRACT_PREPARED`: formal contract complete for future review only.
- `CONTRACT_BLOCKED`: hard governance prerequisites missing.
- `CONTRACT_QUARANTINED`: identity/critical invariant breach.
- `NEEDS_REVALIDATION`: stale evidence must be refreshed.

## Before first real +CPU test (future phase)

- admitted CPU decision from admission gate;
- contract prepared with no P0/P1 blockers;
- fresh identity/QMP/guest evidence;
- rollback and return-to-floor readiness proven;
- kill-switch readiness present.

## Still forbidden

- runtime CPU/RAM apply
- QMP mutating commands
- hotplug execution
- deploy operations
- dashboard/frontend changes
