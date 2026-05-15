# Grande Padre — recommendation engine (Sprint 13)

This sprint adds `pkg/grandepadre/recommendation/`, a **dry-run decision skeleton** that reads indexed evidence from the local in-memory store (`pkg/grandepadre/evidence`) and emits:

- **ActionRecommendation** rows (observe, remediate, prepare-resourcelease, collect-more-evidence),
- **ActionSlate** aggregates (recommendations, blocked evidence, remediable groups, donor/receiver candidates, summary),
- **Risk** and **priority** labels for triage only.

## Non-goals (explicit)

- **Recommendation is not apply.** No cgroup, systemd, Kubernetes, or network mutation is performed.
- **Action slate does not change cluster state.** It is an offline planning artifact.
- **No HTTP server** and **no production controller** in this sprint.
- A future **KHR apply gate** is required before any real mutation path; this engine only prepares intent strings.

## Commercial context (vision)

Hyperdensity as a **predictive live resource market** is the long-term product direction. Sprint 13 stays strictly local: it does not connect to live marketplaces, pricing, or execution.

## CLI

`khr-linux-agent -mode=recommend-actions-local` — see `docs/khr/KHR_LINUX_AGENT_RUNBOOK.md`.
