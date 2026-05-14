# Runtime Market Tick With Real Shell Inputs + No-Apply Production Observation v1

**Milestone:** `hyperdensity_runtime_market_tick_real_shell_inputs_no_apply_v1`

## Purpose

Sprint 18 proved a controlled smoke tick can run.  
Sprint 19 adds a fail-closed **real-input production observation market tick** driven by observed shell, idle, pressure, SLO, rollback, and risk inputs — with no apply and no production movement.

## Core Product Rule

Real observed input means shell identity from Shell Registry / Inventory / Kubernetes objects; idle and pressure from measurement surfaces; SLO from performance proof source; rollback from rollback readiness source. Synthetic/reference records are explicitly separated and cannot count as production evidence.

Sprint 19 is no-apply: actions and futures may be generated; `selectedForExecution=false`; `productionMovementExecuted=false`.

## Allowed Sprint 19 claim

“KARL Hyperdensity can run a no-apply production observation market tick using real observed shell, idle, pressure, SLO and rollback inputs to generate bounded actions and futures while keeping general production auto disabled.”
