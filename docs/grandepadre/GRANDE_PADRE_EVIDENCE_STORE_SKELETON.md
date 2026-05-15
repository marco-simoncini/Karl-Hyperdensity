# Grande Padre — local evidence store skeleton (Sprint 12)

This document describes the **in-memory** evidence store under `pkg/grandepadre/evidence/`. It is a **contract and indexing skeleton** for Hyperdensity: it does **not** run a controller, HTTP server, Kubernetes apply, or any host mutation.

## Purpose

- Parse `EvidenceIngestRequest` documents (YAML/JSON) produced by KHR (`prepare-ingest-request` or hand-authored examples).
- Derive canonical bundle digest, classify **trust tier** (integrity-oriented, not authorization), and build **EvidenceIndex** rows.
- Support **queries** (`ListReady`, `ListBlocked`, `ListByConfidence`, `GetByCell`, …) and a **blocked/remediable** aggregation for recommendations.

## Non-goals

- No admission decisions, no apply gates, no lease enforcement.
- No persistence: restarting the process drops the store.
- No transport: the CLI mode `index-evidence-local` reads a file and prints JSON only.

## CLI

See `docs/khr/KHR_LINUX_AGENT_RUNBOOK.md` (`-mode=index-evidence-local`).

## Related contracts

- `docs/ingest/KHR_EVIDENCE_INGEST_CONTRACT.md`
- `docs/ingest/EVIDENCE_INGEST_SECURITY_BOUNDARIES.md`
