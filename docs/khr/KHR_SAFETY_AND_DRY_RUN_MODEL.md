# KHR Safety and Dry-Run Model (Sprint 5)

## Hard rules

1. **Real apply disabled by default** — no cgroup writes unless `--allow-unsafe-apply` is set (still experimental).
2. **`--allow-unsafe-apply` absent** — CLI clears any simulated `writePaths` from output paths for envelope plans.
3. **Mode gate** — `ResourceLease` dry-run **blocks** unless `spec.mode` is exactly `envelope` (case-insensitive).
4. **Resource gate** — only `cpu` and `memory` are admissible for the Linux MVP skeleton.
5. **Platform gate** — donor/receiver platforms must be `linux` (explicit context or defaulted in evaluator for Sprint 5 demos).
6. **ResourcePort gate** — envelope must appear in the relevant CPU or memory mode list; otherwise **blocked**.

## Dry-run output semantics

`allowed: true` means **“would proceed to next gated stage in a future controller”**, not that mutations occurred.

`expectedWrites` lists **hypothetical** cgroup files for operator review — never executed in Sprint 5 default paths.

## Rollback and verification (contractual)

Every successful dry-run includes textual `rollbackPlan` and `verificationPlan` arrays to align with Hyperdensity evidence discipline.

## Non-goals

- No automatic promotion to apply.
- No integration with cluster RBAC or service accounts in this sprint.
