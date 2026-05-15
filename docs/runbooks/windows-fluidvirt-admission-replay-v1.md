# windows-fluidvirt-admission-replay-v1

Runbook for replaying admission policy decisions in local read-only mode.

## Safety rules

- Local replay only.
- No deploy, no cluster mutation, no VM mutation.
- No frontend or dashboard changes.

## CLI usage

```bash
go run ./cmd/karl-fluid-admission \
  -fixture ./examples/windows-fluid-admission-fixtures/admission-master-win11-cpu.future-apply-admissible.json \
  -evaluation-time 2026-05-07T14:40:00Z
```

Optional custom policy:

```bash
go run ./cmd/karl-fluid-admission \
  -bundle /path/to/runtime-evidence-bundle.json \
  -policy /path/to/policy-pack.json \
  -requested-action prepare-cpu-lease
```

## Reading decisions

- `ADMITTED_FOR_FUTURE_APPLY`: governance-approved for future apply review only.
- `DENIED`: model/context rejected (for example pool replica or generic Windows VM).
- `BLOCKED`: hard blockers or policy violations.
- `QUARANTINED`: identity/critical continuity break.
- `NEEDS_MORE_EVIDENCE`: missing capability/freshness signals.

## Evidence score interpretation

- Score is conservative and coupled with blocker tiers.
- High score does not bypass hard blockers.
- Review `missingEvidence`, `hardBlockers`, `softUnknowns` before any next step.

## BLOCKED vs QUARANTINED

- `BLOCKED`: remediable evidence/policy gaps.
- `QUARANTINED`: continuity integrity breach; isolate candidate and re-establish identity proofs.

## Before first future +CPU review

- Dry-run `READY` or `LEASE_PREPARED`
- QMP ready + guest ACK
- same node/pod/qemu/boot/machine proofs
- no migration/recreate/reboot
- rollback and return-to-floor ready

## Still forbidden in this phase

- CPU/RAM runtime apply
- QMP mutating commands
- hotplug execution
- deploy operations
