# Grande Padre evidence indexing model

## Intent

Grande Padre treats `EvidenceBundle` / `EvidenceIngestRequest` as **observation inputs** for:

1. **Indexing** — correlate Cell, node, tenant, and historical pressure.
2. **Dry-run** — combine with `ResourceLease` / `ResourcePort` simulations already in Hyperdensity contracts.
3. **Recommendation** — surface readiness, blast-radius, and next actions (`recommendedNextAction` from the bundle).

## Not automatic apply

Indexing and recommendations **must not** imply apply. Execution paths require explicit gates (`mutate-preview-apply-dry-run-v1` and operator controls).

## dryRunOnly

When `EvidenceIngestRequest.spec.dryRunOnly` is `true`, Grande Padre must keep the artifact in **simulation-only** paths: no promotion to mutation controllers, no lease activation, regardless of `readyForGrandePadre`.

## Confidence and blockers

- `confidence: low` or non-empty `blockedReasons` should **down-rank** automation suggestions.
- `readyForGrandePadre: false` is a hard signal that upstream automation should not treat the snapshot as “green” for execution planning.

## DevOnly integrity

Bundles with `local-dev` signatures are indexed with **DevOnly** trust tier; they must not satisfy production integrity requirements unless a future PKI contract explicitly replaces this tier.
