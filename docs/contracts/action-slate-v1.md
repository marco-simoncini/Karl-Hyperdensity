# action-slate-v1

Contract ID: `hyperdensity_action_slate_v1`

Defines the KARL Hyperdensity Action Slate v1 projection for prevalidated, operator-controlled resource rebalancing candidates.

## Product principle

- Action Slate is recommendation/projection only.
- KARL prepares candidate moves before pressure becomes outage.
- Perceived speed comes from executing prepared moves, not discovering moves during incidents.
- No autonomous apply.
- No enforcement.
- No production mutation by default.
- No raw runtime controls.
- No raw resource creation.

## Scope and safety invariants

- `actionSlateId=hyperdensity_action_slate_v1`
- `actionSlateVersion=v1`
- `actionSlateMode=prevalidated_recommendation_only`
- `releaseTrack=technical_preview`
- `operatorControlled=true`
- `autonomousApplyAllowed=false`
- `enforcementMode=disabled`
- `productionMutationAllowed=false`
- `evidenceScope=evidence_namespace_only`
- `pairingStrategy=top_k_pairing`
- `batchingMode=rate_limited_operator_controlled`
- `fullFleetScanRequired=false`

## Required action model

Each action must expose:

- identity and kind (`actionId`, `actionKind`)
- readiness state (`actionState`)
- bounded resource scope (`resource`, `shellKind`)
- donor/receiver context and amount (`donor`, `receiver`, `amount`)
- recommendation-only semantics (`recommendationMode`)
- explicit dry-run and rollback state (`dryRunState`, `rollbackState`)
- explicit risk (`risk`)
- support boundary wording (`supportBoundary`)
- evidence references (`evidenceRefs`)
- TTL (`validForSeconds`, `expiresAt`)
- operator next-step (`operatorAction`)
- forbidden positive actions (`forbiddenActions`)
- immutable safety booleans (`autonomousApplyAllowed=false`, `productionMutationAllowed=false`)

## Required safety copy

Serialized safety wording must include:

- "Autonomous apply is disabled."
- "Enforcement is disabled."
- "No production mutation."
- "No raw runtime controls are exposed."
- "No raw resource creation."
- "Dry-run is not mutation."
- "Preview is not apply."
- "Simulation is not enforcement."
- "Support claims are evidence-backed only."
- "Technical Preview boundary active."

## Performance and predictive claims

- Sub-second decisioning and seconds-level transfer values in this contract are targets unless measured.
- Predictive Resource Futures are out-of-scope for this milestone.

## Out-of-scope and support boundaries

- Windows remains out-of-scope.
- No GA claim.
- No HA claim from single-node proof.
- VM RAM wording remains `virtio-mem/QMP/QOM runtime overlay` (not generic KubeVirt template mutation).
