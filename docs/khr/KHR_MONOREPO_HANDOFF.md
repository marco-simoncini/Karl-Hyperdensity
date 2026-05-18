# KHR Monorepo Handoff

**Canonical KHR product repository:** [marco-simoncini/KHR](https://github.com/marco-simoncini/KHR)

## Status

- This source repository is **preserved** for history and integration-specific adapters.
- **No destructive migration** was performed; no files were deleted from this repo.
- **No Windows** qcow2/raw/ISO images or **secrets** were moved into KHR.

## Policy

1. Future KHR product code, schemas, runtime contracts, and sprint evidence metadata should land in **KHR** first.
2. This repo remains valid for deployment-specific wiring until an explicit cutover sprint.
3. Provenance for KHR-REPO-A import is recorded in KHR under `provenance/migration-manifest.json`.

## Source SHA at consolidation

See KHR `provenance/source-shas.json` (KHR-REPO-A, committed-khr-repo-a-v1).

