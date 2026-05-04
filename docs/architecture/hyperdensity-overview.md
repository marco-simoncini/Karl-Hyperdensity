# Hyperdensity Overview

Hyperdensity is KARL's resource-control plane (Grande Padre) for Linux shells.

## Scope

- VM Linux shell
- Container Linux shell

## Compliance Semantics

A shell is compliant only when live bidirectional CPU and RAM behavior is proven with non-disruptive guards:
- CPU up/down live
- RAM up/down live
- no reboot
- no VMI recreate
- no rollout
- no destructive migration
- same runtime continuity
- rollback proof

## Runtime Overlay Principle

Hyperdensity separates three planes:
1. Declared (template/spec envelope and baseline)
2. Runtime (desired/applied by Hyperdensity)
3. Observed (usage/evidence telemetry)

Declared != Runtime != Observed.

## Policy Posture

- Recommendation mode: `recommendation_only`
- Apply mode: `operator_controlled`
- Autonomous mode: disabled
