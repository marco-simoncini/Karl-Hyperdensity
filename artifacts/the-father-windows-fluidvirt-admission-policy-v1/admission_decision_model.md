# admission_decision_model

`WindowsFluidAdmissionDecision` contains:

- immutable decision identity (`decisionId`, `shellRef`, `createdAt`);
- requested action (`certify-shell`, `prepare-cpu-lease`, `prepare-memory-lease`, `evidence-refresh`, `return-to-floor-check`, `quarantine`);
- admission phase:
  - `ADMITTED_FOR_FUTURE_APPLY`
  - `DENIED`
  - `BLOCKED`
  - `QUARANTINED`
  - `NEEDS_MORE_EVIDENCE`
- enforced runtime limits:
  - `mutationAllowed=false`
  - `applyAllowed=false`
  - `runtimeMode=in-place-qmp`
- evidence score and level;
- blockers, denial reasons, required additional evidence;
- blast radius, rollback, return-to-floor, TTL policy snapshots;
- audit references.
