# KHR evidence ingest contract (Sprint 11)

## Scope

This document defines the **API contract only** for moving local evidence produced by `khr-linux-agent` (`collect-evidence` + Sprint 10 integrity sidecars) into the **Grande Padre / Hyperdensity** plane. It does **not** mandate a production backend, transport, or kube apply.

## Core objects

| CRD | Role |
|-----|------|
| `EvidenceBundle` | Durable projection of a single evidence snapshot (hashes, readiness, refs). |
| `EvidenceIngestRequest` | Admission-shaped request carrying inline `bundle`, `manifest`, `digest`, policy, and `dryRunOnly`. |

CRD manifests live under `api/crds/hyperdensity.karl.io/`.

## Non-goals

- No required HTTP client in `khr-linux-agent`.
- No cgroup writes, systemd changes, or cluster mutations from this contract alone.
- No implication that ingest grants **apply** or **admission** authority.

## Ingest vs apply

**Ingest** records evidence for indexing, simulation, recommendations, and blast-radius context. **Apply** remains a separate, explicitly gated workflow (`mutate-preview-apply-dry-run` family). `EvidenceIngestRequest.spec.dryRunOnly` forces simulation-only handling on the consumer side.

## Integrity vs authorization

- `bundleSha256` / `digest` are **integrity** inputs (detect tampering or copy errors).
- `signingMode` / `signaturePresent` describe **local signing posture**; they are **not** RBAC or admission decisions.
- `local-dev` signatures are **DevOnly** trust tier (see `KHR_EVIDENCE_INTEGRITY_MODEL.md`).

## Grande Padre usage model

Grande Padre is expected to use evidence for **index**, **dry-run**, and **recommendation** surfaces—never as a sole automatic apply trigger. See `GRANDE_PADRE_EVIDENCE_INDEXING_MODEL.md`.

## Readiness signals

`confidence`, `readyForGrandePadre`, `blockedReasons`, and `warnings` (from the bundle) influence automation **readiness** and ranking; they do not override higher-level admission.

## Local CLI stub

`khr-linux-agent -mode prepare-ingest-request` materializes a YAML/JSON `EvidenceIngestRequest` from local files only (no network). See `KHR_LINUX_AGENT_RUNBOOK.md`.
