# admission-guard-enforce-simulation-v1

Contract ID: `hyperdensity_admission_guard_enforce_simulation_v1`

Simulation-only projection for future Admission Guard enforce impact.

## Purpose

This contract surfaces the expected impact **if** future enforcement were enabled, while keeping current behavior unchanged:

- admission guard remains `audit_only`
- mutate preview remains `audit_preview_only`
- enforcement remains `disabled`
- autonomous apply remains `false`
- production mutation remains forbidden

## Required top-level fields

- `simulationId`
- `simulationVersion`
- `simulationMode`
- `enforcementMode`
- `admissionGuardMode`
- `mutatePreviewMode`
- `autonomousApplyAllowed`
- `policyPackId`
- `policyConsistencyRequired`
- `evidenceScope`
- `productionMutationAllowed`
- `simulatedObjects`
- `simulatedDecisions`
- `summary`
- `safetyGates`
- `nextSimulationAction`
- `simulationBlocker`

## Required simulated object fields

- `objectRef`
- `namespace`
- `kind`
- `name`
- `factoryManaged`
- `shellKind`
- `shellProfile`
- `readinessState`
- `policyPackCovered`
- `currentAdmissionDecision`
- `simulatedEnforcementDecision`
- `wouldReject`
- `wouldAllow`
- `wouldWarn`
- `wouldMutate`
- `mutatePreviewAvailable`
- `rejectionReasons`
- `warningReasons`
- `remediationHints`
- `sourceRules`
- `safetyNotes`

## Contract semantics

- `simulation_only` means no webhook enforcement activation and no resource mutation.
- Raw/non-conforming objects can simulate `wouldReject=true`.
- Factory-managed compliant or warming objects must not be falsely marked ready.
- `warming_up`, `partial`, and `blocked` are not ready states.
- Windows remains out-of-scope.
- Evidence namespace is the only live proof namespace.
