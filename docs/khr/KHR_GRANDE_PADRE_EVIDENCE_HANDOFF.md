# Grande Padre / Hyperdensity — local evidence bundle handoff (Sprint 9)

## What you receive

Operators or automation run `khr-linux-agent -mode collect-evidence` and ship the stdout JSON (or `-evidence-output` file) upstream. The bundle is designed for **ingest + triage**, not as an implicit apply approval.

## Ingest checklist

1. **Schema / version**: Read `version` and ensure your consumer supports that agent line.
2. **`evidenceSummary.readyForGrandePadre`**: When `false`, treat the bundle as **blocked for automation** until `blockedReasons` are cleared; still store for audit.
3. **`evidenceSummary.confidence`**: Floors across discovery, telemetry evidence, and (when present) ResourceLease dry-run outcome.
4. **`discovery.selectedPath`**: If empty, telemetry did not run against a live path; do not infer host pressure from `telemetry.metrics`.
5. **`dryRun`**: If `skipped`, check `skipReason` — partial lease/port inputs are operator mistakes, not host failures.
6. **`mutationsForbidden`**: Must remain `true` for this build; if a future agent ever sends `false`, your gate should hard-reject until explicitly redesigned.

## Recommended pairing

- Cluster admission / blast-radius context (not produced by this binary).
- Historical telemetry or SLO signals (this bundle is a point-in-time slice).

## References

- Local bundle field reference: `docs/khr/KHR_LOCAL_EVIDENCE_BUNDLE.md`
- Telemetry-only evidence: `docs/khr/KHR_TELEMETRY_EVIDENCE_MODEL.md`
- Agent runbook: `docs/khr/KHR_LINUX_AGENT_RUNBOOK.md`
