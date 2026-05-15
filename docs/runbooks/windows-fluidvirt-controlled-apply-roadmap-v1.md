# Windows FluidVirt Controlled Apply Roadmap v1

This runbook explains the transition from dry-run planning to controlled apply, while keeping autonomous execution disabled.

## Current Milestone

- `cmd/karl-fluid-windows-executor` supports `plan`, `dry-run`, and `apply-plan-only`
- output is plan/evidence JSON only
- no cluster mutations, no QMP runtime calls, no cgroup writes

## Required Gates Before Controlled Apply

1. Compliance must be `HYPERDENSITY_READY_WINDOWS_SHELL`
2. Dry-run gate must pass
3. Manual approval must be `approved`
4. CPU/RAM apply gates must be enabled for the target lease
5. Kill switch state must be `allow`
6. Guest/workload verification must be required and present
7. Rollback and return-to-floor must be ready
8. Audit bundle append must be ready

## Plan States

- `awaiting_approval`: manual approval missing
- `apply_ready`: all required gates are satisfied
- `apply_blocked` / `blocked`: one or more hard blockers
- `quarantined`: continuity/safety identity blocker class

## Blast Radius Controls

- keep `maxBlastRadius=1`
- enforce `allowedNamespaces`, `allowedTargets`, `allowedLeaseKinds`
- reject ambiguous/expanded target scopes

## What Remains Forbidden

- autonomous apply
- vCPU hotplug/unplug
- logical CPU scaling claims
- pool/replica scaling as runtime mechanism
- LiveMigration/VMIM as mechanism
- reboot/recreate/rollout-driven resizing

## What Is Needed Before Production Enablement

- explicit policy unlock milestone
- signed attestation pipeline
- audited runtime executor with staged rollout guardrails
- continuous kill-switch and replay protection evidence
- formal SRE/on-call operational SOP
