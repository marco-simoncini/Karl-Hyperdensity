# Inventory fluidShell Integration Plan

## Decision
Integrate `Karl-Inventory` fluidShell as **separate PR track** and **witness/evidence-only component**.

## Functional Boundary
Allowed role:
- guest ACK
- same-boot witness
- machine identity hash evidence
- pending reboot and critical events evidence
- memory/state evidence for return-to-floor and post-apply verification

Forbidden role:
- no actuator role
- no scaling engine role
- no autonomous apply role

## Candidate Scope (Windows branch)
- `docs/fluid-shell-module.md`
- `inventory/agent-windows/installer/config.json`
- `inventory/agent-windows/service/src/KarlInventoryAgent.FluidShell/*`
- `inventory/agent-windows/service/src/KarlInventoryAgent/Services/FluidShell/*`
- related config/program/worker wiring and tests

## Integration Checks
- verify baseline branch compatibility (mainline branch naming currently requires confirmation in target repo)
- remove/ignore generated build outputs from commits
- run `dotnet test` for `KarlInventoryAgent.FluidShell.Tests`
- enforce claim boundary in docs and release notes

## PR Timing
Start Inventory PR after Hyperdensity PR 1-3 contracts are stable, so witness payload binds to final contract IDs.
